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
	IsUniqRequest(ctx context.Context, userID string, personID string) error
	IsNotFriend(ctx context.Context, userID string, personID string) error

	MakeOutcomingRequest(ctx context.Context, userID string, personID string) error
	MakeIncomingRequest(ctx context.Context, userID string, personID string) error
	MakeFriends(ctx context.Context, userID string, personID string) error

	DeleteOutcomingRequest(ctx context.Context, userID string, personID string) error
	DeleteIncomingRequest(ctx context.Context, userID string, personID string) error
	DeleteFriend(ctx context.Context, exFriendID1 string, exFriendID2 string) error

	GetOutcomingRequestsByUserID(ctx context.Context, userID string) ([]string, error)
	GetIncomingRequestsByUserID(ctx context.Context, userID string) ([]string, error)
	GetFriendsByUserID(ctx context.Context, userID string) ([]string, error)

	GetFriendsByID(ctx context.Context, userID string) ([]string, error)
}

type friendsRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *friendsRepositoryImpl) CreateFriends(ctx context.Context, userID string) error {
	friends := core.Friends{ID: userID}
	_, err := repo.coll.InsertOne(ctx, friends)
	return err
}

func (repo *friendsRepositoryImpl) IsUniqRequest(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": userID, "requests": personID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrRequestAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) IsNotFriend(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": userID, "friends": personID}

	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrAlreadyFriends
		} else {
			return err
		}
	}

	return nil
}

func (repo *friendsRepositoryImpl) MakeOutcomingRequest(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$push": bson.D{{Key: "requests", Value: personID}}}

	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeIncomingRequest(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": personID}
	update := bson.M{"$push": bson.D{{Key: "incoming_requests", Value: userID}}}

	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) MakeFriends(ctx context.Context, userID string, personID string) error {
	// first friend
	filter := bson.M{"_id": userID}
	update := bson.M{"$push": bson.D{{Key: "friends", Value: personID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	// second friend
	filter = bson.M{"_id": personID}
	update = bson.M{"$push": bson.D{{Key: "friends", Value: userID}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteOutcomingRequest(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$pull": bson.M{"requests": personID}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) DeleteIncomingRequest(ctx context.Context, userID string, personID string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$pull": bson.M{"incoming_requests": personID}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *friendsRepositoryImpl) DeleteFriend(ctx context.Context, exFriendID1 string, exFriendID2 string) error {
	filter := bson.M{"_id": exFriendID2}
	update := bson.M{"$pull": bson.M{"friends": exFriendID1}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"_id": exFriendID1}
	update = bson.M{"$pull": bson.M{"friends": exFriendID2}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) GetOutcomingRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(&friends)
	return friends.OutcomingRequests, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetIncomingRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(&friends)
	return friends.IncomingRequest, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetFriendsByUserID(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(&friends)
	return friends.Friends, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetFriendsByID(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(&friends)
	return friends.Friends, wrapError(err)
}

func NewFriendsRepository(db *mongo.Database) (*friendsRepositoryImpl, error) {
	return &friendsRepositoryImpl{db: db, coll: db.Collection("friends")}, nil
}
