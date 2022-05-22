package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreateComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		comment := TestComment(t)
		ctx := context.Background()
		commentActual, err := commentCollection.CreateComment(ctx, comment)
		assert.Nil(t, err)
		assert.NotNil(t, comment.ID)
		assert.NotNil(t, comment.CreatedAt)
		assert.NotNil(t, commentActual)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		comment := TestComment(t)
		ctx := context.Background()
		_, err := commentCollection.CreateComment(ctx, comment)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestGetCommentByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		expectedComment := TestComment(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedComment.ID},
			{"author_id", expectedComment.AuthorID},
			{"message", expectedComment.Message},
			{"images", expectedComment.Images},
			{"created_at", expectedComment.CreatedAt},
		}))
		ctx := context.Background()
		comment, err := commentCollection.GetCommentByID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedComment, comment)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		comment, err := commentCollection.GetCommentByID(ctx, "0")
		testNullComment := TestCommentNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullComment, comment)
	})
}

func TestDeleteComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		ctx := context.Background()
		err := commentCollection.DeleteComment(ctx, TestPost(t).ID)
		assert.Nil(t, err)
	})
}

func TestEditComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		expectedComment := TestComment(t)
		ctx := context.Background()
		comment, err := commentCollection.EditComment(ctx, expectedComment)
		assert.Nil(t, err)
		assert.Equal(t, expectedComment, comment)
	})
}

func TestInitComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		commentCollection, _ := NewCommentRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		comment := TestCommentNull(t)
		err := commentCollection.InitComment(comment)
		assert.Nil(t, err)
		assert.NotNil(t, comment.ID)
		assert.NotNil(t, comment.CreatedAt)
	})
}
