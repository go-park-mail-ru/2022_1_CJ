package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func TestCreateDialog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		dialog := TestDialog(t)
		ctx := context.Background()
		dialogActual, err := chatCollection.CreateDialog(ctx, dialog.Name, dialog.Name, dialog.Participants)
		assert.Nil(t, err)
		assert.NotNil(t, dialogActual.ID)
		assert.NotNil(t, dialogActual.CreatedAt)
		assert.NotNil(t, dialogActual)
	})

	mt.Run("custom error duplicate in insert", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))
		dialog := TestDialog(t)
		ctx := context.Background()
		_, err := chatCollection.CreateDialog(ctx, dialog.Name, dialog.Name, dialog.Participants)
		assert.NotNil(t, err)
		rte := mongo.IsDuplicateKeyError(err)
		assert.True(t, rte)
	})
}

func TestIsChatExist(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		expectedDialog := TestDialog(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedDialog.ID},
			{Key: "name", Value: expectedDialog.Name},
			{Key: "participants", Value: expectedDialog.Participants},
			{Key: "messages", Value: expectedDialog.Messages},
			{Key: "created_at", Value: expectedDialog.CreatedAt},
		}))
		ctx := context.Background()
		err := chatCollection.IsChatExist(ctx, TestDialog(t).ID)
		assert.Nil(t, err)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		err := chatCollection.IsChatExist(ctx, "0")
		assert.NotNil(t, err)
	})
}

func TestIsUniqDialog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("not uniq dialog", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		expectedDialog := TestDialog(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedDialog.ID},
			{Key: "name", Value: expectedDialog.Name},
			{Key: "participants", Value: expectedDialog.Participants},
			{Key: "messages", Value: expectedDialog.Messages},
			{Key: "created_at", Value: expectedDialog.CreatedAt},
		}))
		ctx := context.Background()
		err := chatCollection.IsUniqDialog(ctx, TestDialog(t).Participants[0], TestDialog(t).Participants[1])
		assert.Equal(t, constants.ErrDialogAlreadyExist, err)
	})

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		err := chatCollection.IsUniqDialog(ctx, "0", "1")
		assert.NotNil(t, err)
	})
}

func TestSendMessage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		dialog := TestDialog(t)
		message := TestMessage(t)
		ctx := context.Background()

		err := chatCollection.SendMessage(ctx, *message, dialog.ID)

		assert.Nil(t, err)
	})
}

func TestReadMessage(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("not success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		dialog := TestDialog(t)
		message := TestMessage(t)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.Background()

		err := chatCollection.ReadMessage(ctx, message.IsRead[0].Participant, message.ID, dialog.ID)

		assert.Equal(t, constants.ErrDBNotFound, err)
	})
}

func TestGetDialogByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		expectedDialog := TestDialog(t)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedDialog.ID},
			{Key: "name", Value: expectedDialog.Name},
			{Key: "participants", Value: expectedDialog.Participants},
			{Key: "messages", Value: expectedDialog.Messages},
			{Key: "created_at", Value: expectedDialog.CreatedAt},
		}))
		ctx := context.Background()
		dialog, err := chatCollection.GetDialogByID(ctx, expectedDialog.ID)
		assert.Nil(t, err)
		assert.Equal(t, expectedDialog, dialog)
	})

	mt.Run("don't find in collection", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateCursorResponse(0, "foo.bar", mtest.FirstBatch))
		ctx := context.Background()
		post, err := chatCollection.GetDialogByID(ctx, "0")
		testNullDialog := TestDialogNull(t)
		assert.NotNil(t, err)
		assert.Equal(t, testNullDialog, post)
	})
}

func TestInitDialog(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		chatCollection, _ := NewChatRepositoryTest(mt.Coll)

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		dialog := TestDialogNull(t)
		err := chatCollection.InitDialog(dialog, TestUser(t).ID, TestDialog(t).Participants, TestDialog(t).Name)
		assert.Nil(t, err)
		assert.NotNil(t, dialog.ID)
		assert.NotNil(t, dialog.CreatedAt)
		assert.NotNil(t, dialog.Participants)
		assert.NotNil(t, dialog.Name)
	})
}
