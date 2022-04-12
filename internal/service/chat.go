package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type ChatService interface {
	CreateDialog(ctx context.Context, request *dto.CreateDialogRequest) (*dto.CreateDialogResponse, error)
	SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetDialogs(ctx context.Context, request *dto.GetDialogsRequest) (*dto.GetDialogsResponse, error)
}

type chatServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *chatServiceImpl) CreateDialog(ctx context.Context, request *dto.CreateDialogRequest) (*dto.CreateDialogResponse, error) {
	switch {
	case len(request.AuthorIDs) < 1:
		svc.log.Errorf("CreateDialog error: %s", constants.ErrSingleChat)
		return nil, constants.ErrSingleChat
	case len(request.AuthorIDs) == 1:
		if err := svc.db.ChatRepo.IsUniqDialog(ctx, request.UserID, request.AuthorIDs[0]); err != nil {
			svc.log.Errorf("CreateDialog error: %s", err)
			return nil, err
		}
	}

	dialog, err := svc.db.ChatRepo.CreateDialog(ctx, request.UserID, request.AuthorIDs)
	if err != nil {
		svc.log.Errorf("CreateDialog error: %s", err)
		return nil, err
	}
	svc.log.Debug("Create dialog success")
	if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, request.UserID); err != nil {
		svc.log.Errorf("CreateDialog error: %s", err)
		return nil, err
	}

	for _, id := range request.AuthorIDs {
		if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, id); err != nil {
			svc.log.Errorf("CreateDialog error: %s", err)
			return nil, err
		}
	}
	svc.log.Debug("Add dialog to users success")
	return &dto.CreateDialogResponse{DialogID: dialog.ID}, nil
}

// Сделать получатель всех чатов

func (svc *chatServiceImpl) SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	message := &common.MessageInfo{Text: request.MessageInfo.Text,
		AuthorID: request.MessageInfo.AuthorID,
		DialogID: request.MessageInfo.DialogID}
	svc.log.Debugf("Text: %s; DialogID: %s; AuthorID: %s", message.Text, message.DialogID, message.AuthorID)

	if err := svc.db.ChatRepo.IsChatExist(ctx, message.DialogID); err != nil {
		svc.log.Errorf("Chat not exist error: %s", err)
		return nil, err
	}

	svc.log.Debug("Chat was founded")

	if err := svc.db.ChatRepo.SendMessage(ctx, *message); err != nil {
		svc.log.Errorf("SendMessage error: %s", err)
		return nil, err
	}

	svc.log.Debug("Message was sent successful")

	return &dto.SendMessageResponse{}, nil
}

func (svc *chatServiceImpl) GetDialogs(ctx context.Context, request *dto.GetDialogsRequest) (*dto.GetDialogsResponse, error) {
	ids, err := svc.db.UserRepo.GetUserDialogs(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserDialogs error: %s", err)
		return nil, err
	}
	var DialogsInfo []common.DialogInfo
	for _, id := range ids {
		dInf, err := svc.db.ChatRepo.GetDialogInfo(ctx, id)
		if err != nil {
			svc.log.Errorf("GetDialogInfo error: %s", err)
		}
		DialogsInfo = append(DialogsInfo, dInf)
	}
	return &dto.GetDialogsResponse{DialogsInfo: DialogsInfo}, err
}

func NewChatService(log *logrus.Entry, db *db.Repository) ChatService {
	return &chatServiceImpl{log: log, db: db}
}
