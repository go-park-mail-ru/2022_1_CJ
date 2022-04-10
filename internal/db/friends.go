package db

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	GetFriendsByID(ctx context.Context, FriendsID string) ([]string, error)
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
	filter := bson.M{"user_id": UserID, "requests": PersonID}
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
	filter := bson.M{"user_id": UserID, "friends": PersonID}

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
	update := bson.M{"$push": bson.D{{Key: "requests", Value: PersonID}}}

	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeFriends(ctx context.Context, UserID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": UserID}
	update := bson.M{"$push": bson.D{{Key: "friends", Value: PersonID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	// second friend
	filter = bson.M{"user_id": PersonID}
	update = bson.M{"$push": bson.D{{Key: "friends", Value: UserID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteRequest(ctx context.Context, UserID string, PersonID string) error {
	// first friend
	filter := bson.M{"user_id": UserID}
	update := bson.M{"$pull": bson.M{"requests": PersonID}} // Проверить на работу
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// -------------------------DELETE
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

func (repo *friendsRepositoryImpl) GetFriendsByID(ctx context.Context, FriendsID string) ([]string, error) {
	friends := new(core.Friends)
	filter := bson.M{"_id": FriendsID}
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
