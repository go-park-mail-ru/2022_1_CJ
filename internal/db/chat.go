package db

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	IsUniqDialog(ctx context.Context, userID1 string, userID2 string) error
	IsDialogExist(ctx context.Context, userID1 string, userID2 string) (string, error)
	CreateDialog(ctx context.Context, userID string, name string, authorIDs []string) (*core.Dialog, error)
	IsChatExist(ctx context.Context, dialogID string) error
	SendMessage(ctx context.Context, message core.Message, dialogID string) error
	ReadMessage(ctx context.Context, userID string, messageID string, dialogID string) error
	GetDialogByID(ctx context.Context, dialogID string) (*core.Dialog, error)
}

type chatRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewChatRepository(db *mongo.Database) (*chatRepositoryImpl, error) {
	return &chatRepositoryImpl{db: db, coll: db.Collection("chats")}, nil
}

// NewUserRepositoryTest for Tests (bad)
func NewChatRepositoryTest(collection *mongo.Collection) (*chatRepositoryImpl, error) {
	return &chatRepositoryImpl{coll: collection}, nil
}

func (repo *chatRepositoryImpl) IsUniqDialog(ctx context.Context, userID1 string, userID2 string) error {
	filter := bson.M{"participants": bson.D{{Key: "$all", Value: bson.A{userID1, userID2}}}}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrDialogAlreadyExist
		} else {
			return err
		}
	}
	return nil
}

func (repo *chatRepositoryImpl) IsDialogExist(ctx context.Context, userID1 string, userID2 string) (string, error) {
	chat := &core.Dialog{}
	filter := bson.M{"participants": bson.D{{Key: "$all", Value: bson.A{userID1, userID2}}}}
	if err := repo.coll.FindOne(ctx, filter).Decode(&chat); err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}
	return chat.ID, nil
}

func (repo *chatRepositoryImpl) CreateDialog(ctx context.Context, userID string, name string, authorIDs []string) (*core.Dialog, error) {
	dialog := new(core.Dialog)
	if err := repo.InitDialog(dialog, userID, authorIDs, name); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, dialog)

	// Sanitize
	p := bluemonday.UGCPolicy()

	dialog.Name = p.Sanitize(dialog.Name)

	for i := range dialog.Participants {
		dialog.Participants[i] = p.Sanitize(dialog.Participants[i])
	}

	return dialog, err
}

func (repo *chatRepositoryImpl) IsChatExist(ctx context.Context, dialogID string) error {
	filter := bson.M{"_id": dialogID}
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

func (repo *chatRepositoryImpl) ReadMessage(ctx context.Context, userID string, messageID string, dialogID string) error {

	filter := bson.M{"_id": dialogID}

	update := bson.M{"$set": bson.D{{Key: "messages.$[n1].is_participants_read.$[n2].is_read", Value: true}}}
	var s []interface{}
	ch1 := bson.D{{Key: "n1._id", Value: messageID}}
	ch2 := bson.D{{Key: "n2._id", Value: userID}}
	s = append(s, ch1)
	s = append(s, ch2)
	arrayFilter := options.FindOneAndUpdateOptions{ArrayFilters: &options.ArrayFilters{Filters: s}}

	if err := repo.coll.FindOneAndUpdate(ctx, filter, update, &arrayFilter).Err(); err != nil {
		return wrapError(err)
	}

	return nil

}

func (repo *chatRepositoryImpl) GetDialogByID(ctx context.Context, DialogID string) (*core.Dialog, error) {
	dialog := new(core.Dialog)
	filter := bson.M{"_id": DialogID}
	err := repo.coll.FindOne(ctx, filter).Decode(dialog)

	// Sanitize
	p := bluemonday.UGCPolicy()

	dialog.Name = p.Sanitize(dialog.Name)

	for i := range dialog.Messages {
		dialog.Messages[i].Body = p.Sanitize(dialog.Messages[i].Body)
	}

	for i := range dialog.Participants {
		dialog.Participants[i] = p.Sanitize(dialog.Participants[i])
	}

	return dialog, wrapError(err)
}

func (repo *chatRepositoryImpl) InitDialog(dialog *core.Dialog, userID string, authorIDs []string, name string) error {
	id, err := core.GenUUID()
	if err != nil {
		return err
	}
	dialog.CreatedAt = time.Now().Unix()
	dialog.ID = id
	dialog.Name = name
	dialog.Participants = append(dialog.Participants, userID)
	for _, authorID := range authorIDs {
		dialog.Participants = append(dialog.Participants, authorID)
	}
	return nil
}
