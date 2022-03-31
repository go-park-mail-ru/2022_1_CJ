package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Скорее всего Update by id не будет работать
type FriendsRepository interface {
	CreateFriends(ctx context.Context, FriendsID string, UserID string) error
	IsUniqRequest(ctx context.Context, UserID string, PersonID string) error
	IsNotFriend(ctx context.Context, UserID string, PersonID string) error

	MakeRequest(ctx context.Context, UserID string, PersonID string) error
	MakeFriends(ctx context.Context, UserID string, PersonID string) error

	DeleteRequest(ctx context.Context, UserID string, PersonID string) error
	DeleteFriend(ctx context.Context, ExFriendID1 string, ExFriendID2 string) error

	GetRequestsByUserID(ctx context.Context, UserID string) ([]string, error)
	GetFriendsByUserID(ctx context.Context, UserID string) ([]string, error)
}

func (repo *friendsRepositoryImpl) CreateFriends(ctx context.Context, FriendsID string, UserID string) error {
	friends := new(core.Friends)
	friends.ID = FriendsID
	friends.UserID = UserID

	_, err := repo.coll.InsertOne(ctx, friends)
	return err
}

// --------------REQUESTS
func (repo *friendsRepositoryImpl) IsUniqRequest(ctx context.Context, UserID string, PersonID string) error {
	filter := bson.M{"user_id": UserID, "requests": bson.M{"$in": PersonID}}

	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrRequestAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) IsNotFriend(ctx context.Context, UserID string, PersonID string) error {
	filter := bson.M{"user_id": UserID, "friends": bson.M{"$in": PersonID}}

	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrAlreadyFriends
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeRequest(ctx context.Context, UserID string, PersonID string) error {
	filter := bson.M{"user_id": UserID}
	update := bson.M{"requests": bson.M{"$push": PersonID}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeFriends(ctx context.Context, UserID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": UserID}
	update := bson.M{"friends": bson.M{"$push": PersonID}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	// second friend
	filter = bson.M{"user_id": PersonID}
	update = bson.M{"friends": bson.M{"$push": UserID}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteRequest(ctx context.Context, UserID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": UserID}
	update := bson.M{"$pull": bson.M{"requests": bson.M{"$in": PersonID}}} // Проверить на работу
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// -------------------------DELETE
func (repo *friendsRepositoryImpl) DeleteFriend(ctx context.Context, ExFriendID1 string, ExFriendID2 string) error {
	filter := bson.M{"user_id": ExFriendID2}
	update := bson.M{"$pull": bson.M{"friends": bson.M{"$in": ExFriendID1}}} // Проверить на работу
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"user_id": ExFriendID1}
	update = bson.M{"$pull": bson.M{"friends": bson.M{"$in": ExFriendID2}}} // Проверить на работу
	if _, err := repo.coll.UpdateByID(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) GetRequestsByUserID(ctx context.Context, UserID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"user_id": UserID}
	err := repo.coll.FindOne(ctx, filter).Decode(friends)
	return friends.Requests, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetFriendsByUserID(ctx context.Context, UserID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"user_id": UserID}
	err := repo.coll.FindOne(ctx, filter).Decode(friends)
	return friends.Friends, wrapError(err)
}

type friendsRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewFriendsRepository(db *mongo.Database) (*friendsRepositoryImpl, error) {
	return &friendsRepositoryImpl{db: db, coll: db.Collection("friends")}, nil
}
