package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type ChatRepository interface {
	IsUniqDialog(ctx context.Context, fUserID string, sUserID string) error
	CreateDialog(ctx context.Context, UserId string, AuthorIDs []string) (*core.Dialog, error)
	IsChatExist(ctx context.Context, DialogID string) error
	SendMessage(ctx context.Context, Message common.MessageInfo) error
	GetDialogInfo(ctx context.Context, DialogID string) (common.DialogInfo, error)
}

type chatRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *chatRepositoryImpl) IsUniqDialog(ctx context.Context, fUserID string, sUserID string) error {
	filter := bson.M{"$size": bson.D{{Key: "author_ids", Value: 2}},
		"author_ids": bson.A{fUserID, sUserID}} // Проверка на правильность
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrDialogAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *chatRepositoryImpl) CreateDialog(ctx context.Context, UserId string, AuthorIDs []string) (*core.Dialog, error) {
	dialog := new(core.Dialog)
	if err := initDialog(dialog, UserId, AuthorIDs); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, dialog)
	return dialog, err
}

func (repo *chatRepositoryImpl) IsChatExist(ctx context.Context, DialogID string) error {
	filter := bson.D{{"_id", DialogID}}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func (repo *chatRepositoryImpl) SendMessage(ctx context.Context, Message common.MessageInfo) error {
	newMassage := new(core.Message)
	if err := initNewMassage(newMassage, Message); err != nil {
		return err
	}
	update := bson.D{{"$put", newMassage}}
	if _, err := repo.coll.UpdateByID(ctx, Message.DialogID, update); err != nil {
		return err
	}
	return nil
}

func (repo *chatRepositoryImpl) GetDialogInfo(ctx context.Context, DialogID string) (common.DialogInfo, error) {
	dialog := new(core.Dialog)
	filter := bson.D{{"_id", DialogID}}
	err := repo.coll.FindOne(ctx, filter).Decode(dialog)
	return common.DialogInfo{DialogID: DialogID,
			Title: strings.Join(dialog.AuthorIDs[:], ",")},
		wrapError(err)
}

func initDialog(dialog *core.Dialog, userID string, authorIDs []string) error {
	did, err := core.GenUUID()
	if err != nil {
		return err
	}
	dialog.ID = did
	dialog.AuthorIDs = append(dialog.AuthorIDs, userID)
	for _, id := range authorIDs {
		dialog.AuthorIDs = append(dialog.AuthorIDs, id)
	}
	return nil
}

func initNewMassage(newMsg *core.Message, m common.MessageInfo) error {
	msgID, err := core.GenUUID()
	if err != nil {
		return err
	}
	newMsg.ID = msgID
	newMsg.Text = m.Text
	newMsg.AuthorID = m.AuthorID
	newMsg.CreatedAt = time.Now().Unix()
	return nil
}

func NewChatRepository(db *mongo.Database) (*chatRepositoryImpl, error) {
	return &chatRepositoryImpl{db: db, coll: db.Collection("chats")}, nil
}
