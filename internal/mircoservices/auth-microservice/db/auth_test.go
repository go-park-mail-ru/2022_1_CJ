package auth_db

import (
	"context"
	auth_core "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/core"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func AuthUser(t *testing.T) *auth_core.User {
	t.Helper()

	return &auth_core.User{
		ID:        "123",
		Email:     "rinat@mail.com",
		Password:  auth_core.UserPassword{},
		CreatedAt: 12,
	}
}

func NullAuthUser(t *testing.T) *auth_core.User {
	t.Helper()

	return &auth_core.User{}
}

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		authUser := AuthUser(t)
		ctx := context.Background()
		_, err := authCollection.CreateUser(ctx, authUser)
		assert.NotNil(t, err)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		authUser := AuthUser(t)
		ctx := context.Background()
		_, err := authCollection.CreateUser(ctx, authUser)
		assert.NotNil(t, err)
	})
}

func TestGetUserByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)
		expectedlyUser := AuthUser(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedlyUser.ID},
			{"email", expectedlyUser.Email},
			{"password", expectedlyUser.Password},
			{"created_at", expectedlyUser.CreatedAt},
		}))
		ctx := context.Background()
		authUser, err := authCollection.GetUserByID(ctx, AuthUser(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedlyUser, authUser)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		post, err := authCollection.GetUserByID(ctx, "0")
		testNullUser := NullAuthUser(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, post)
	})
}

func TestGetUserByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)
		expectedlyUser := AuthUser(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedlyUser.ID},
			{"email", expectedlyUser.Email},
			{"password", expectedlyUser.Password},
			{"created_at", expectedlyUser.CreatedAt},
		}))
		ctx := context.Background()
		authUser, err := authCollection.GetUserByEmail(ctx, AuthUser(t).Email)
		assert.Nil(t, err)
		assert.Equal(t, expectedlyUser, authUser)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		post, err := authCollection.GetUserByEmail(ctx, "0")
		testNullUser := NullAuthUser(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, post)
	})
}

func TestCheckEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		actual, err := authCollection.CheckUserEmailExistence(ctx, AuthUser(t).Email)
		assert.NotNil(t, err)
		assert.Equal(t, false, actual)
	})
}

func TestInitUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		authCollection, _ := NewAuthRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		authUser := NullAuthUser(t)
		err := authCollection.InitUser(authUser)
		assert.Nil(t, err)
		assert.NotNil(t, authUser.ID)
		assert.NotNil(t, authUser.CreatedAt)
	})
}
