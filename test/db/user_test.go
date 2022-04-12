package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch), mtest.CreateSuccessResponse())
		user := TestUser(t)
		ctx := context.Background()
		err := userCollection.CreateUser(ctx, user)
		assert.Nil(t, err)
		assert.NotNil(t, user.ID)
		assert.NotNil(t, user.FriendsID)
		assert.NotNil(t, user.CreatedAt)
	})

	mt.Run("Find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		expectedUser := core.User{
			ID: "1234567890",
			Name: common.UserName{
				First: "Sasha",
				Last:  "Userov",
			},
			Email: "user@example.org",
			Image: "src/img.jpg",
			Phone: "+8(800)-555-35-35",
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedUser.ID},
			{"name", expectedUser.Name},
			{"email", expectedUser.Email},
			{"email", expectedUser.Image},
			{"email", expectedUser.Phone},
		}))
		user := TestUser(t)
		ctx := context.Background()
		err := userCollection.CreateUser(ctx, user)
		assert.NotNil(t, err)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch), mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		user := TestUser(t)
		ctx := context.Background()
		err := userCollection.CreateUser(ctx, user)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})

}
