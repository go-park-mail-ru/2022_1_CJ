package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreateFriends(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		friends := TestFriends(t)
		ctx := context.Background()
		err := friendsCollection.CreateFriends(ctx, friends.ID)
		assert.Nil(t, err)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		friends := TestFriends(t)
		ctx := context.Background()
		err := friendsCollection.CreateFriends(ctx, friends.ID)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestIsUniqRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("already friends", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.IsUniqRequest(ctx, "123", "234")
		assert.NotNil(t, err)
	})

	mt.Run("uniq request", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		err := friendsCollection.IsUniqRequest(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestIsNotFriend(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("already friends", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.IsNotFriend(ctx, "123", "234")
		assert.NotNil(t, err)
	})

	mt.Run("not friends", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		err := friendsCollection.IsNotFriend(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestMakeOutcomingRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.MakeOutcomingRequest(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestMakeIncomingRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.MakeIncomingRequest(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestMakeFriends(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.MakeFriends(ctx, "123", "234")
		assert.Nil(t, err)
	})

	mt.Run("UpdateOne error", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		ctx := context.Background()
		err := friendsCollection.MakeFriends(ctx, "123", "234")
		assert.NotNil(t, err)
	})
}

func TestDeleteOutcomingRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("delete error", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		ctx := context.Background()
		err := friendsCollection.DeleteOutcomingRequest(ctx, "123", "234")
		assert.NotNil(t, err)
	})

	mt.Run("success delete", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})
		ctx := context.Background()
		err := friendsCollection.DeleteOutcomingRequest(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestDeleteIncomingRequest(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("delete error", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		ctx := context.Background()
		err := friendsCollection.DeleteIncomingRequest(ctx, "123", "234")
		assert.NotNil(t, err)
	})

	mt.Run("success delete", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})
		ctx := context.Background()
		err := friendsCollection.DeleteIncomingRequest(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestDeleteFriend(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("delete error", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		ctx := context.Background()
		err := friendsCollection.DeleteFriend(ctx, "123", "234")
		assert.NotNil(t, err)
	})

	mt.Run("success delete", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := friendsCollection.DeleteFriend(ctx, "123", "234")
		assert.Nil(t, err)
	})
}

func TestGetOutcomingRequestsByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success return outcoming requests", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		friends := TestFriends(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", friends.ID},
			{"requests", friends.OutcomingRequests},
			{"incoming_requests", friends.IncomingRequest},
			{"friends", friends.Friends},
		}))
		requests, err := friendsCollection.GetOutcomingRequestsByUserID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, friends.OutcomingRequests, requests)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		requests, err := friendsCollection.GetOutcomingRequestsByUserID(ctx, TestPost(t).ID)
		testNullFriends := TestFriendsNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullFriends.OutcomingRequests, requests)
	})
}

func TestGetIncomingRequestsByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success return incoming requests", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		friends := TestFriends(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", friends.ID},
			{"requests", friends.OutcomingRequests},
			{"incoming_requests", friends.IncomingRequest},
			{"friends", friends.Friends},
		}))
		requests, err := friendsCollection.GetIncomingRequestsByUserID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, friends.IncomingRequest, requests)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		requests, err := friendsCollection.GetIncomingRequestsByUserID(ctx, TestPost(t).ID)
		testNullFriends := TestFriendsNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullFriends.IncomingRequest, requests)
	})
}

func TestGetFriendsByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success return friends", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		friends := TestFriends(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", friends.ID},
			{"requests", friends.OutcomingRequests},
			{"incoming_requests", friends.IncomingRequest},
			{"friends", friends.Friends},
		}))
		requests, err := friendsCollection.GetFriendsByUserID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, friends.Friends, requests)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		friendsCollection, _ := NewFriendsRepositoryTest(mt.Coll)
		ctx := context.Background()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		requests, err := friendsCollection.GetFriendsByUserID(ctx, TestPost(t).ID)
		testNullFriends := TestFriendsNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullFriends.Friends, requests)
	})
}
