package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
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
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
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

func TestGetUserByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
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
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
		}))
		user := TestUser(t)
		ctx := context.Background()
		_, err := userCollection.GetUserByID(ctx, user.ID)
		assert.Nil(t, err)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		user, err := userCollection.GetUserByID(ctx, "0")
		testNullUser := TestUserNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, user)
	})

}

func TestGetUserByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
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
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
		}))
		user := TestUser(t)
		ctx := context.Background()
		_, err := userCollection.GetUserByEmail(ctx, user.Email)
		assert.Nil(t, err)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		user, err := userCollection.GetUserByEmail(ctx, "0")
		testNullUser := TestUserNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, user)
	})
}

func TestCheckUserEmailExistence(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
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
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
		}))
		user := TestUser(t)
		ctx := context.Background()
		res, err := userCollection.CheckUserEmailExistence(ctx, user.ID)
		assert.Nil(t, err)
		assert.True(t, res)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		res, err := userCollection.CheckUserEmailExistence(ctx, "0")
		assert.Nil(t, err)
		assert.True(t, !res)
	})
}

func TestUpdateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		user := TestUser(t)
		ctx := context.Background()
		err := userCollection.UpdateUser(ctx, user)
		assert.Nil(t, err)
	})
}

func TestUserAddPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := userCollection.UserAddPost(ctx, "1234567890", "123")
		assert.Nil(t, err)
	})
}

func TestUserCheckPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := userCollection.UserCheckPost(ctx, TestUser(t), "123")
		assert.Nil(t, err)
	})
}

func TestUserDeletePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)
		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())
		ctx := context.Background()
		err := userCollection.UserDeletePost(ctx, TestUser(t).ID, TestPost(t).ID)
		assert.Nil(t, err)
	})

	mt.Run("don't find post from user in db", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))

		ctx := context.Background()
		err := userCollection.UserDeletePost(ctx, TestUser(t).ID, TestPost(t).ID)
		assert.NotNil(t, err)
		assert.Equal(t, constants.ErrDBNotFound, err)
	})
}

func TestUserDelete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)
		mt.AddMockResponses(bson.D{{"ok", 1}, {"acknowledged", true}, {"n", 1}})

		ctx := context.Background()
		err := userCollection.DeleteUser(ctx, TestUser(t))
		assert.Nil(t, err)
	})
}

func TestSelectUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		expectedUser1 := core.User{
			ID: "1234567890",
			Name: common.UserName{
				First: "Sasha",
				Last:  "Userov",
			},
		}
		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedUser1.ID},
			{"name", expectedUser1.Name},
		})

		expectedUser2 := core.User{
			ID: "201029481",
			Name: common.UserName{
				First: "Sashaureolov",
				Last:  "Userov",
			},
		}
		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedUser2.ID},
			{"name", expectedUser2.Name},
		})

		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors)

		ctx := context.Background()
		users, err := userCollection.SelectUsers(ctx, "Sash")
		assert.Nil(t, err)
		assert.Equal(t, []core.User{
			{ID: expectedUser1.ID, Name: expectedUser1.Name},
			{ID: expectedUser2.ID, Name: expectedUser2.Name},
		}, users)
	})
}

func TestAddDialog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := userCollection.AddDialog(ctx, TestDialog(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestUserCheckDialog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
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
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
		}))
		ctx := context.Background()
		err := userCollection.UserCheckDialog(ctx, TestUser(t).ID, TestDialog(t).ID)
		assert.Nil(t, err)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		err := userCollection.UserCheckDialog(ctx, "0", "0")
		assert.NotNil(t, err)
		assert.Equal(t, constants.ErrDBNotFound, err)
	})
}

func TestGetUserDialogs(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		expectedUser := core.User{
			ID: "1234567890",
			Name: common.UserName{
				First: "Sasha",
				Last:  "Userov",
			},
			Email:     "user@example.org",
			Image:     "src/img.jpg",
			Phone:     "+8(800)-555-35-35",
			DialogIDs: []string{"12432536443"},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedUser.ID},
			{"name", expectedUser.Name},
			{"email", expectedUser.Email},
			{"images", expectedUser.Image},
			{"phone", expectedUser.Phone},
			{"dialog_ids", expectedUser.DialogIDs},
		}))
		user := TestUser(t)
		ctx := context.Background()
		friends, err := userCollection.GetUserDialogs(ctx, user.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedUser.DialogIDs, friends)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		friends, err := userCollection.GetUserDialogs(ctx, "0")
		testNullUser := TestUserNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser.DialogIDs, friends)
	})
}

func TestInitUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		userCollection, _ := db.NewUserRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		user := TestUserNull(t)
		err := userCollection.InitUser(user)
		assert.Nil(t, err)
		assert.NotNil(t, user.ID)
		assert.NotNil(t, user.CreatedAt)
		assert.Equal(t, "default.jpeg", user.Image)
	})
}
