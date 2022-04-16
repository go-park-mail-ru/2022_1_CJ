package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/sirupsen/logrus"
)

type ChatService interface {
	CreateDialog(ctx context.Context, request *dto.CreateDialogRequest) (*dto.CreateDialogResponse, error)
	SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetDialogs(ctx context.Context, request *dto.GetDialogsRequest) (*dto.GetDialogsResponse, error)

	GetDialog(ctx context.Context, request *dto.GetDialogRequest) (*dto.GetDialogResponse, error)
}

type chatServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *chatServiceImpl) CreateDialog(ctx context.Context, request *dto.CreateDialogRequest) (*dto.CreateDialogResponse, error) {
	switch {
	case len(request.AuthorIDs) < 1:
		svc.log.Errorf("%s", constants.ErrSingleChat)
		return nil, constants.ErrSingleChat
	case len(request.AuthorIDs) == 1:
		// TODO: don't correct working IsUniqDialog
		//if err := svc.db.ChatRepo.IsUniqDialog(ctx, request.UserID, request.AuthorIDs[0]); err != nil {
		//	svc.log.Errorf("IsUniqDialog error: %s", err)
		//	return nil, err
		//}
	}

	dialog, err := svc.db.ChatRepo.CreateDialog(ctx, request.UserID, request.AuthorIDs)
	if err != nil {
		svc.log.Errorf("CreateDialog error: %s", err)
		return nil, err
	}

	svc.log.Debug("Create dialog success")
	if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, request.UserID); err != nil {
		svc.log.Errorf("AddDialog error: %s", err)
		return nil, err
	}

	for _, id := range request.AuthorIDs {
		if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, id); err != nil {
			svc.log.Errorf("AddDialog error: %s", err)
			return nil, err
		}
	}
	svc.log.Debug("Add dialog to users success")
	return &dto.CreateDialogResponse{DialogID: dialog.ID}, nil
}

// Сделать получатель всех чатов

func (svc *chatServiceImpl) SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	message := &core.Message{Body: request.Message.Body,
		AuthorID: request.Message.AuthorID, IsRead: false}

	svc.log.Debugf("Text: %s; DialogID: %s; AuthorID: %s", message.Body, request.Message.DialogID, message.AuthorID)

	if err := svc.db.ChatRepo.IsChatExist(ctx, request.Message.DialogID); err != nil {
		svc.log.Errorf("Chat not exist error: %s", err)
		return nil, err
	}

	if err := svc.db.ChatRepo.SendMessage(ctx, message, request.Message.DialogID); err != nil {
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
	var dialogs []dto.Dialog
	for _, id := range ids {
		dInf, err := svc.db.ChatRepo.GetDialogByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetDialogInfo error: %s", err)
		}
		dialogs = append(dialogs, convert.Dialog2DTO(dInf, request.UserID))
	}
	return &dto.GetDialogsResponse{Dialogs: dialogs}, err
}

func (svc *chatServiceImpl) GetDialog(ctx context.Context, request *dto.GetDialogRequest) (*dto.GetDialogResponse, error) {
	ids, err := svc.db.UserRepo.GetUserDialogs(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserDialogs error: %s", err)
		return nil, err
	}

	check := func(s []string, str string) bool {
		for _, v := range s {
			if v == str {
				return true
			}
		}
		return false
	}
	if isHave := check(ids, request.DialogID); isHave != true {
		svc.log.Errorf("User don't have dialog")
		return nil, constants.ErrDBNotFound
	}

	dialog, err := svc.db.ChatRepo.GetDialogByID(ctx, request.DialogID)
	if err != nil {
		svc.log.Errorf("GetDialogByID error: %s", err)
		return nil, err
	}

	return &dto.GetDialogResponse{Dialog: convert.Dialog2DTO(dialog, request.UserID), Messages: convert.Messages2DTO(dialog.Messages)}, err
}

func NewChatService(log *logrus.Entry, db *db.Repository) ChatService {
	return &chatServiceImpl{log: log, db: db}
}
