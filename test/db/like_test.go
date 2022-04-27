package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreateLike(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		like := TestNullLike(t)
		ctx := context.Background()
		likeActual, err := likeCollection.CreateLike(ctx, like)
		assert.Nil(t, err)
		assert.NotNil(t, like.ID)
		assert.NotNil(t, like.CreatedAt)
		assert.NotNil(t, likeActual)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		like := TestNullLike(t)
		ctx := context.Background()
		_, err := likeCollection.CreateLike(ctx, like)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestGetLikeBySubjectID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		expectedLike := TestLike(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedLike.ID},
			{"subject_id", expectedLike.Subject},
			{"amount", expectedLike.Amount},
			{"user_ids", expectedLike.UserIDs},
			{"created_at", expectedLike.CreatedAt},
		}))
		ctx := context.Background()
		like, err := likeCollection.GetLikeBySubjectID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedLike, like)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		like, err := likeCollection.GetLikeBySubjectID(ctx, "0")
		testNullUser := TestNullLike(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, like)
	})
}

func TestIncreaseLike(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := likeCollection.IncreaseLike(ctx, TestLike(t).Subject, "1")
		assert.NotNil(t, err)
	})
}

func TestReduceLike(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := likeCollection.ReduceLike(ctx, TestLike(t).Subject, "1")
		assert.NotNil(t, err)
	})
}

func TestInitLike(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		likeCollection, _ := db.NewLikeRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		like := TestNullLike(t)
		err := likeCollection.InitLike(like)
		assert.Nil(t, err)
		assert.NotNil(t, like.ID)
		assert.NotNil(t, like.CreatedAt)
	})
}
