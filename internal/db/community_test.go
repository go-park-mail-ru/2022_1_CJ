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

func TestCreateCommunity(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		community := TestNullCommunity(t)
		ctx := context.Background()
		communityActual, err := communityCollection.CreateCommunity(ctx, community)
		assert.Nil(t, err)
		assert.NotNil(t, community.ID)
		assert.NotNil(t, community.CreatedAt)
		assert.NotNil(t, communityActual)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		community := TestNullCommunity(t)
		ctx := context.Background()
		_, err := communityCollection.CreateCommunity(ctx, community)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestGetCommunityByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		expectedCommunity := TestCommunity(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedCommunity.ID},
			{"name", expectedCommunity.Name},
			{"image", expectedCommunity.Image},
			{"info", expectedCommunity.Info},
			{"followers", expectedCommunity.FollowerIDs},
			{"admins", expectedCommunity.AdminIDs},
			{"posts", expectedCommunity.PostIDs},
			{"created_at", expectedCommunity.CreatedAt},
		}))
		ctx := context.Background()
		community, err := communityCollection.GetCommunityByID(ctx, TestPost(t).ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedCommunity, community)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		community, err := communityCollection.GetCommunityByID(ctx, "0")
		testNullUser := TestNullCommunity(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullUser, community)
	})
}

func TestDeleteCommunity(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.DeleteCommunity(ctx, TestCommunity(t).ID)
		assert.Nil(t, err)
	})
}

func TestEditCommunity(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.EditCommunity(ctx, TestCommunity(t))
		assert.Nil(t, err)
	})
}

func TestAddFollower(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.AddFollower(ctx, TestCommunity(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestCommunityAddPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.CommunityAddPost(ctx, TestCommunity(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestCommunityDeletePost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.CommunityDeletePost(ctx, TestCommunity(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestDeleteAdmin(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.DeleteAdmin(ctx, TestCommunity(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestDeleteFollower(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse(), mtest.CreateSuccessResponse())

		ctx := context.Background()
		err := communityCollection.DeleteFollower(ctx, TestCommunity(t).ID, TestUser(t).ID)
		assert.Nil(t, err)
	})
}

func TestGetAllCommunities(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		expectedCommunity1 := TestCommunity(t)

		expectedCommunity2 := TestCommunity(t)
		expectedCommunity2.ID = "1"

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedCommunity1.ID},
			{"name", expectedCommunity1.Name},
			{"image", expectedCommunity1.Image},
			{"info", expectedCommunity1.Info},
			{"followers", expectedCommunity1.FollowerIDs},
			{"admins", expectedCommunity1.AdminIDs},
			{"posts", expectedCommunity1.PostIDs},
			{"created_at", expectedCommunity1.CreatedAt},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedCommunity2.ID},
			{"name", expectedCommunity2.Name},
			{"image", expectedCommunity2.Image},
			{"info", expectedCommunity2.Info},
			{"followers", expectedCommunity2.FollowerIDs},
			{"admins", expectedCommunity2.AdminIDs},
			{"posts", expectedCommunity2.PostIDs},
			{"created_at", expectedCommunity2.CreatedAt},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors, mtest.CreateSuccessResponse())

		ctx := context.Background()
		communities, _, err := communityCollection.GetAllCommunities(ctx, -1, 1)
		assert.Nil(t, err)
		assert.Equal(t, []core.Community{{
			ID:          expectedCommunity1.ID,
			Name:        expectedCommunity1.Name,
			Image:       expectedCommunity1.Image,
			Info:        expectedCommunity1.Info,
			FollowerIDs: expectedCommunity1.FollowerIDs,
			AdminIDs:    expectedCommunity1.AdminIDs,
			PostIDs:     expectedCommunity1.PostIDs,
			CreatedAt:   expectedCommunity1.CreatedAt,
		}, {
			ID:          expectedCommunity2.ID,
			Name:        expectedCommunity2.Name,
			Image:       expectedCommunity2.Image,
			Info:        expectedCommunity2.Info,
			FollowerIDs: expectedCommunity2.FollowerIDs,
			AdminIDs:    expectedCommunity2.AdminIDs,
			PostIDs:     expectedCommunity2.PostIDs,
			CreatedAt:   expectedCommunity2.CreatedAt,
		}}, communities)
	})
}

func TestSearchCommunities(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		expectedCommunity1 := TestCommunity(t)

		expectedCommunity2 := TestCommunity(t)
		expectedCommunity2.ID = "1"

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{"_id", expectedCommunity1.ID},
			{"name", expectedCommunity1.Name},
			{"image", expectedCommunity1.Image},
			{"info", expectedCommunity1.Info},
			{"followers", expectedCommunity1.FollowerIDs},
			{"admins", expectedCommunity1.AdminIDs},
			{"posts", expectedCommunity1.PostIDs},
			{"created_at", expectedCommunity1.CreatedAt},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{"_id", expectedCommunity2.ID},
			{"name", expectedCommunity2.Name},
			{"image", expectedCommunity2.Image},
			{"info", expectedCommunity2.Info},
			{"followers", expectedCommunity2.FollowerIDs},
			{"admins", expectedCommunity2.AdminIDs},
			{"posts", expectedCommunity2.PostIDs},
			{"created_at", expectedCommunity2.CreatedAt},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(first, second, killCursors, mtest.CreateSuccessResponse())

		ctx := context.Background()
		communities, _, err := communityCollection.SearchCommunities(ctx, "community", -1, 1)
		assert.Nil(t, err)
		assert.Equal(t, []core.Community{{
			ID:          expectedCommunity1.ID,
			Name:        expectedCommunity1.Name,
			Image:       expectedCommunity1.Image,
			Info:        expectedCommunity1.Info,
			FollowerIDs: expectedCommunity1.FollowerIDs,
			AdminIDs:    expectedCommunity1.AdminIDs,
			PostIDs:     expectedCommunity1.PostIDs,
			CreatedAt:   expectedCommunity1.CreatedAt,
		}, {
			ID:          expectedCommunity2.ID,
			Name:        expectedCommunity2.Name,
			Image:       expectedCommunity2.Image,
			Info:        expectedCommunity2.Info,
			FollowerIDs: expectedCommunity2.FollowerIDs,
			AdminIDs:    expectedCommunity2.AdminIDs,
			PostIDs:     expectedCommunity2.PostIDs,
			CreatedAt:   expectedCommunity2.CreatedAt,
		}}, communities)
	})
}

func TestInitCommunity(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		communityCollection, _ := NewCommunityRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		community := TestNullCommunity(t)
		err := communityCollection.InitCommunity(community)
		assert.Nil(t, err)
		assert.NotNil(t, community.ID)
		assert.NotNil(t, community.CreatedAt)
	})
}
