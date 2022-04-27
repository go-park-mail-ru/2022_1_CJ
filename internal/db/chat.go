package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ChatRepository interface {
	IsUniqDialog(ctx context.Context, firstUserID string, secondUserID string) error
	CreateDialog(ctx context.Context, userID string, authorIDs []string) (*core.Dialog, error)
	IsChatExist(ctx context.Context, dialogID string) error
	SendMessage(ctx context.Context, message core.Message, dialogID string) error
	ReadMessage(ctx context.Context, messageID string, dialogID string) error
	GetDialogByID(ctx context.Context, dialogID string) (*core.Dialog, error)
}

type chatRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *chatRepositoryImpl) IsUniqDialog(ctx context.Context, firstUserID string, secondUserID string) error {
	filter := bson.M{"$and": bson.A{
		bson.D{{"author_ids.2", bson.D{{"$exists", false}}}},
		bson.D{{"author_ids", firstUserID}},
		bson.D{{"author_ids", secondUserID}},
	}}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrDialogAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *chatRepositoryImpl) CreateDialog(ctx context.Context, userID string, authorIDs []string) (*core.Dialog, error) {
	dialog := new(core.Dialog)
	if err := initDialog(dialog, userID, authorIDs); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, dialog)
	return dialog, err
}

func (repo *chatRepositoryImpl) IsChatExist(ctx context.Context, dialogID string) error {
	filter := bson.D{{"_id", dialogID}}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func (repo *chatRepositoryImpl) SendMessage(ctx context.Context, message core.Message, dialogID string) error {
	filter := bson.M{"_id": dialogID}
	update := bson.M{"$push": bson.D{{Key: "messages", Value: message}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

func (repo *chatRepositoryImpl) ReadMessage(ctx context.Context, messageID string, dialogID string) error {
	filter := bson.M{"_id": dialogID, "messages._id": messageID}
	update := bson.M{"$set": bson.D{{Key: "messages.$.is_read", Value: true}}}
	if _, err := repo.coll.UpdateOne(ctx, filter, update); err != nil {
		return wrapError(err)
	}
	return nil
}

func (repo *chatRepositoryImpl) GetDialogByID(ctx context.Context, DialogID string) (*core.Dialog, error) {
	dialog := new(core.Dialog)
	filter := bson.D{{"_id", DialogID}}
	err := repo.coll.FindOne(ctx, filter).Decode(dialog)
	return dialog, wrapError(err)
}

func initDialog(dialog *core.Dialog, userID string, authorIDs []string) error {
	id, err := core.GenUUID()
	if err != nil {
		return err
	}
	dialog.CreatedAt = time.Now().Unix()
	dialog.ID = id

	dialog.Participants = append(dialog.Participants, userID)
	for _, authorID := range authorIDs {
		dialog.Participants = append(dialog.Participants, authorID)
	}
	return nil
}

func NewChatRepository(db *mongo.Database) (*chatRepositoryImpl, error) {
	return &chatRepositoryImpl{db: db, coll: db.Collection("chats")}, nil
}
