package db

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type FriendsRepository interface {
	CreateFriends(ctx context.Context, userID string) error
	IsUniqRequest(ctx context.Context, userID string, PersonID string) error
	IsNotFriend(ctx context.Context, userID string, PersonID string) error

	MakeRequest(ctx context.Context, userID string, PersonID string) error
	MakeFriends(ctx context.Context, userID string, PersonID string) error

	DeleteRequest(ctx context.Context, userID string, PersonID string) error
	DeleteFriend(ctx context.Context, ExFriendID1 string, ExFriendID2 string) error

	GetRequestsByUserID(ctx context.Context, userID string) ([]string, error)
	GetFriendsByUserID(ctx context.Context, userID string) ([]string, error)

	GetFriendsByID(ctx context.Context, userID string) ([]string, error)
}

type friendsRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *friendsRepositoryImpl) CreateFriends(ctx context.Context, userID string) error {
	friends := new(core.Friends)
	friends.ID = userID
	_, err := repo.coll.InsertOne(ctx, friends)
	return err
}

func (repo *friendsRepositoryImpl) IsUniqRequest(ctx context.Context, userID string, PersonID string) error {
	filter := bson.M{"user_id": userID, "requests": PersonID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrRequestAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) IsNotFriend(ctx context.Context, userID string, PersonID string) error {
	filter := bson.M{"user_id": userID, "friends": PersonID}

	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrAlreadyFriends
		} else {
			return err
		}
	}

	return nil
}

func (repo *friendsRepositoryImpl) MakeRequest(ctx context.Context, userID string, PersonID string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$push": bson.D{{Key: "requests", Value: PersonID}}}

	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeFriends(ctx context.Context, userID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": userID}
	update := bson.M{"$push": bson.D{{Key: "friends", Value: PersonID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	// second friend
	filter = bson.M{"user_id": PersonID}
	update = bson.M{"$push": bson.D{{Key: "friends", Value: userID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteRequest(ctx context.Context, userID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": userID}
	update := bson.M{"$pull": bson.M{"requests": PersonID}} // Проверить на работу
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) DeleteFriend(ctx context.Context, ExFriendID1 string, ExFriendID2 string) error {
	filter := bson.M{"user_id": ExFriendID2}
	update := bson.M{"$pull": bson.M{"friends": ExFriendID1}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"user_id": ExFriendID1}
	update = bson.M{"$pull": bson.M{"friends": ExFriendID2}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) GetRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"user_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(friends)
	return friends.Requests, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetFriendsByUserID(ctx context.Context, userID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"user_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(friends)
	return friends.Friends, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetFriendsByID(ctx context.Context, userID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(friends)
	return friends.Friends, wrapError(err)
}

func NewFriendsRepository(db *mongo.Database) (*friendsRepositoryImpl, error) {
	return &friendsRepositoryImpl{db: db, coll: db.Collection("friends")}, nil
}
