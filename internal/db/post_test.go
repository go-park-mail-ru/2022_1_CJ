package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreatePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		post := TestPost(t)
		ctx := context.Background()
		postActual, err := postCollection.CreatePost(ctx, post)
		assert.Nil(t, err)
		assert.NotNil(t, post.ID)
		assert.NotNil(t, post.CreatedAt)
		assert.NotNil(t, postActual)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		post := TestPost(t)
		ctx := context.Background()
		_, err := postCollection.CreatePost(ctx, post)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestGetPostByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		expectedPost := TestPost(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPost.ID},
			{Key: "author_id", Value: expectedPost.AuthorID},
			{Key: "message", Value: expectedPost.Message},
			{Key: "files", Value: expectedPost.Files},
			{Key: "created_at", Value: expectedPost.CreatedAt},
		}))
		ctx := context.Background()
		post, err := postCollection.GetPostByID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedPost, post)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		post, err := postCollection.GetPostByID(ctx, "0")
		testNullUser := TestPostNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, post)
	})
}

func TestGetPostsByUserID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		expectedPost1 := core.Post{
			ID:        "12345678",
			AuthorID:  "123456789",
			Files:     []string{"src/image1.jpg"},
			CreatedAt: 1323123,
		}
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPost1.ID},
			{Key: "author_id", Value: expectedPost1.AuthorID},
			{Key: "files", Value: expectedPost1.Files},
			{Key: "created_at", Value: expectedPost1.CreatedAt},
		})

		expectedPost2 := core.Post{
			ID:        "101010101",
			AuthorID:  "123456789",
			Files:     []string{"src/image2.jpg"},
			CreatedAt: 1323123,
		}
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expectedPost2.ID},
			{Key: "author_id", Value: expectedPost2.AuthorID},
			{Key: "files", Value: expectedPost2.Files},
			{Key: "created_at", Value: expectedPost2.CreatedAt},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		ctx := context.Background()
		posts, _, err := postCollection.GetPostsByUserID(ctx, "123456789", 1, -1)
		assert.Nil(t, err)
		assert.Equal(t, []core.Post{
			{ID: expectedPost1.ID, AuthorID: expectedPost1.AuthorID, Files: expectedPost1.Files, CreatedAt: expectedPost1.CreatedAt},
			{ID: expectedPost2.ID, AuthorID: expectedPost2.AuthorID, Files: expectedPost2.Files, CreatedAt: expectedPost2.CreatedAt},
		}, posts)
	})
}

func TestDeletePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 1}, {Key: "acknowledged", Value: true}, {Key: "n", Value: 1}})

		ctx := context.Background()
		err := postCollection.DeletePost(ctx, TestPost(t).ID)
		assert.Nil(t, err)
	})
}

func TestEditPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		expectedPost := TestPost(t)
		ctx := context.Background()
		post, err := postCollection.EditPost(ctx, expectedPost)
		assert.Nil(t, err)
		assert.Equal(t, expectedPost, post)
	})
}

func TestGetFeed(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		expectedPost1 := core.Post{
			ID:        "12345678",
			AuthorID:  "123456789",
			Files:     []string{"src/image1.jpg"},
			CreatedAt: 1323123,
		}
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedPost1.ID},
			{Key: "author_id", Value: expectedPost1.AuthorID},
			{Key: "files", Value: expectedPost1.Files},
			{Key: "created_at", Value: expectedPost1.CreatedAt},
		})

		expectedPost2 := core.Post{
			ID:        "101010101",
			AuthorID:  "123456789",
			Files:     []string{"src/image2.jpg"},
			CreatedAt: 1323123,
		}
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: expectedPost2.ID},
			{Key: "author_id", Value: expectedPost2.AuthorID},
			{Key: "files", Value: expectedPost2.Files},
			{Key: "created_at", Value: expectedPost2.CreatedAt},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		ctx := context.Background()
		posts, _, err := postCollection.GetFeed(ctx, "12", 1, -1)
		assert.Nil(t, err)
		assert.Equal(t, []core.Post{
			{ID: expectedPost1.ID, AuthorID: expectedPost1.AuthorID, Files: expectedPost1.Files, CreatedAt: expectedPost1.CreatedAt},
			{ID: expectedPost2.ID, AuthorID: expectedPost2.AuthorID, Files: expectedPost2.Files, CreatedAt: expectedPost2.CreatedAt},
		}, posts)
	})
}

func TestInitPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		post := TestPostNull(t)
		err := postCollection.InitPost(post)
		assert.Nil(t, err)
		assert.NotNil(t, post.ID)
		assert.NotNil(t, post.CreatedAt)
	})
}

func TestPostAddComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		post := TestPost(t)
		comment := TestComment(t)
		err := postCollection.PostAddComment(ctx, post.ID, comment.ID)
		assert.Nil(t, err)
	})
}

func TestPostDeleteComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())
		ctx := context.Background()
		post := TestPost(t)
		comment := TestComment(t)
		err := postCollection.PostDeleteComment(ctx, post.ID, comment.ID)
		assert.Nil(t, err)
	})
}

func TestPostCheckComment(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		postCollection, _ := NewPostRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())
		ctx := context.Background()
		post := TestPost(t)
		comment := TestComment(t)
		err := postCollection.PostCheckComment(ctx, post, comment.ID)
		assert.Nil(t, err)
	})
}
