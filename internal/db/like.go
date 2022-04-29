package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepository interface {
	CreateLike(ctx context.Context, like *core.Like) (*core.Like, error)
	GetLikeBySubjectID(ctx context.Context, subjectID string) (*core.Like, error)
	IncreaseLike(ctx context.Context, subjectID string, userID string) error
	ReduceLike(ctx context.Context, subjectID string, userID string) error
}

type likeRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewLikeRepository(db *mongo.Database) (*likeRepositoryImpl, error) {
	return &likeRepositoryImpl{db: db, coll: db.Collection("likes")}, nil
}

// NewLikeRepositoryTest for Tests (bad)
func NewLikeRepositoryTest(collection *mongo.Collection) (*likeRepositoryImpl, error) {
	return &likeRepositoryImpl{coll: collection}, nil
}

func (repo *likeRepositoryImpl) CreateLike(ctx context.Context, like *core.Like) (*core.Like, error) {
	if err := repo.InitLike(like); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, like)
	return like, err
}

func (repo *likeRepositoryImpl) GetLikeBySubjectID(ctx context.Context, subjectID string) (*core.Like, error) {
	like := new(core.Like)
	filter := bson.M{"subject_id": subjectID}
	err := repo.coll.FindOne(ctx, filter).Decode(like)
	return like, wrapError(err)
}

func (repo *likeRepositoryImpl) IncreaseLike(ctx context.Context, subjectID string, userID string) error {
	filter1 := bson.M{"subject_id": subjectID}
	update1 := bson.M{"$inc": bson.D{{Key: "amount", Value: 1}}}
	if _, err := repo.coll.UpdateOne(ctx, filter1, update1); err != nil {
		return err
	}

	filter2 := bson.M{"subject_id": subjectID}
	update2 := bson.M{"$push": bson.D{{Key: "user_ids", Value: userID}}}
	if err := repo.coll.FindOneAndUpdate(ctx, filter2, update2).Err(); err != nil {
		return wrapError(err)
	}
	return nil
}

func (repo *likeRepositoryImpl) ReduceLike(ctx context.Context, subjectID string, userID string) error {
	filter1 := bson.M{"subject_id": subjectID}
	update1 := bson.M{"$inc": bson.D{{Key: "amount", Value: -1}}}
	if _, err := repo.coll.UpdateOne(ctx, filter1, update1); err != nil {
		return err
	}

	filter2 := bson.M{"subject_id": subjectID}
	update2 := bson.M{"$pull": bson.D{{Key: "user_ids", Value: userID}}}
	if err := repo.coll.FindOneAndUpdate(ctx, filter2, update2).Err(); err != nil {
		return wrapError(err)
	}
	return nil
}

func (repo *likeRepositoryImpl) InitLike(like *core.Like) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	like.ID = uid
	like.Amount = 0
	like.UserIDs = []string{}
	like.CreatedAt = time.Now().Unix()
	return nil
}
