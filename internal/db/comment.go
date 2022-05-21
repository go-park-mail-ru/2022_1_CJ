package db

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *core.Comment) (*core.Comment, error)
	GetCommentByID(ctx context.Context, commentID string) (*core.Comment, error)
	EditComment(ctx context.Context, comment *core.Comment) (*core.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
}

type commentRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) (*commentRepositoryImpl, error) {
	return &commentRepositoryImpl{db: db, coll: db.Collection("comments")}, nil
}

// NewUserRepositoryTest for Tests (bad)
func NewCommentRepositoryTest(collection *mongo.Collection) (*commentRepositoryImpl, error) {
	return &commentRepositoryImpl{coll: collection}, nil
}

func (repo *commentRepositoryImpl) CreateComment(ctx context.Context, comment *core.Comment) (*core.Comment, error) {
	if err := repo.InitComment(comment); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, comment)

	// Sanitize
	p := bluemonday.UGCPolicy()
	comment.Message = p.Sanitize(comment.Message)

	return comment, err
}

func (repo *commentRepositoryImpl) GetCommentByID(ctx context.Context, commentID string) (*core.Comment, error) {
	Comment := new(core.Comment)
	filter := bson.M{"_id": commentID}
	err := repo.coll.FindOne(ctx, filter).Decode(Comment)
	return Comment, wrapError(err)
}

func (repo *commentRepositoryImpl) EditComment(ctx context.Context, comment *core.Comment) (*core.Comment, error) {
	filter := bson.M{"_id": comment.ID}
	_, err := repo.coll.ReplaceOne(ctx, filter, comment)
	return comment, err
}

func (repo *commentRepositoryImpl) DeleteComment(ctx context.Context, commentID string) error {
	filter := bson.M{"_id": commentID}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

func (repo *commentRepositoryImpl) InitComment(comment *core.Comment) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	comment.ID = uid
	comment.CreatedAt = time.Now().Unix()
	return nil
}
