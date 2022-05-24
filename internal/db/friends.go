package db

import (
	"context"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FriendsRepository interface {
	CreateFriends(ctx context.Context, userID string) error
	IsUniqRequest(ctx context.Context, from string, to string) error
	IsNotFriend(ctx context.Context, from string, to string) error

	CreateRequest(ctx context.Context, from string, to string) error
	DeleteRequest(ctx context.Context, from string, to string) error

	MakeFriends(ctx context.Context, userID1 string, userID2 string) error
	DeleteFriend(ctx context.Context, userID1 string, userID2 string) error

	GetFriends(ctx context.Context, userID string) ([]string, error)
	GetIncomingRequests(ctx context.Context, userID string) ([]string, error)
	GetOutcomingRequests(ctx context.Context, userID string) ([]string, error)
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

func (repo *friendsRepositoryImpl) IsUniqRequest(ctx context.Context, from string, to string) error {
	filter := bson.M{"_id": from, "requests": to}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrRequestAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) IsNotFriend(ctx context.Context, from string, to string) error {
	filter := bson.M{"_id": from, "friends": to}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrAlreadyFriends
		} else {
			return err
		}
	}
	return nil
}

func (repo *friendsRepositoryImpl) CreateRequest(ctx context.Context, from string, to string) error {
	filter := bson.M{"_id": to}
	update := bson.M{"$push": bson.M{"incoming_requests": from}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"_id": from}
	update = bson.M{"$push": bson.M{"outcoming_requests": to}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteRequest(ctx context.Context, from string, to string) error {
	filter := bson.M{"_id": to}
	update := bson.M{"$pull": bson.M{"incoming_requests": from}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"_id": from}
	update = bson.M{"$pull": bson.M{"outcoming_requests": to}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) MakeFriends(ctx context.Context, userID1 string, userID2 string) error {
	filter := bson.M{"_id": userID1}
	update := bson.M{"$push": bson.M{"friends": userID2}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"_id": userID2}
	update = bson.M{"$push": bson.M{"friends": userID1}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) DeleteFriend(ctx context.Context, userID1 string, userID2 string) error {
	filter := bson.M{"_id": userID2}
	update := bson.M{"$pull": bson.M{"friends": userID1}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	filter = bson.M{"_id": userID1}
	update = bson.M{"$pull": bson.M{"friends": userID2}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (repo *friendsRepositoryImpl) GetFriends(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(&friends)
	return friends.Friends, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetIncomingRequests(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	opts := options.FindOne().SetProjection(bson.M{"incoming_requests": 1})
	err := repo.coll.FindOne(ctx, filter, opts).Decode(&friends)
	return friends.IncomingRequests, wrapError(err)
}

func (repo *friendsRepositoryImpl) GetOutcomingRequests(ctx context.Context, userID string) ([]string, error) {
	friends := core.Friends{}
	filter := bson.M{"_id": userID}
	opts := options.FindOne().SetProjection(bson.M{"outcoming_requests": 1})
	err := repo.coll.FindOne(ctx, filter, opts).Decode(&friends)
	return friends.OutcomingRequests, wrapError(err)
}

func NewFriendsRepository(db *mongo.Database) (*friendsRepositoryImpl, error) {
	return &friendsRepositoryImpl{db: db, coll: db.Collection("friends")}, nil
}

func NewFriendsRepositoryTest(collection *mongo.Collection) (*friendsRepositoryImpl, error) {
	return &friendsRepositoryImpl{coll: collection}, nil
}
