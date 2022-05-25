package service

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/sirupsen/logrus"
)

type ChatService interface {
	CreateChat(ctx context.Context, request *dto.CreateChatRequest) (*dto.CreateChatResponse, error)
	GetDialogs(ctx context.Context, request *dto.GetDialogsRequest) (*dto.GetDialogsResponse, error)
	GetDialog(ctx context.Context, request *dto.GetDialogRequest) (*dto.GetDialogResponse, error)
	GetDialogByUserID(ctx context.Context, request *dto.GetDialogByUserIDRequest, currentUserID string) (*dto.GetDialogByUserIDResponse, error)

	SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	ReadMessage(ctx context.Context, request *dto.ReadMessageRequest) (*dto.ReadMessageResponse, error)
	CheckDialog(ctx context.Context, request *dto.CheckDialogRequest) error
}

type chatServiceImpl struct {
	log *logrus.Entry
	db  *db.Repository
}

func (svc *chatServiceImpl) CreateChat(ctx context.Context, request *dto.CreateChatRequest) (*dto.CreateChatResponse, error) {
	if len(request.AuthorIDs) == 1 {
		if request.AuthorIDs[0] == request.UserID {
			return nil, constants.ErrSingleChat
		}
		if err := svc.db.ChatRepo.IsUniqDialog(ctx, request.UserID, request.AuthorIDs[0]); err != nil {
			return nil, err
		}
	}

	dialog, err := svc.db.ChatRepo.CreateDialog(ctx, request.UserID, request.Name, request.AuthorIDs)
	if err != nil {
		svc.log.Errorf("CreateDialog error: %s", err)
		return nil, err
	}

	if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, request.UserID); err != nil {
		svc.log.Errorf("AddDialog error: %s", err)
		return nil, err
	}

	svc.log.Debug("Create dialog success")
	for _, id := range request.AuthorIDs {
		if id != request.UserID {
			if err := svc.db.UserRepo.AddDialog(ctx, dialog.ID, id); err != nil {
				svc.log.Errorf("AddDialog error: %s", err)
				return nil, err
			}
		}
	}
	svc.log.Debug("Add dialog to users success")
	return &dto.CreateChatResponse{DialogID: dialog.ID}, nil
}

func (svc *chatServiceImpl) SendMessage(ctx context.Context, request *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	dialog, err := svc.db.ChatRepo.GetDialogByID(ctx, request.Message.DialogID)
	if err != nil {
		return nil, fmt.Errorf("GetDialogByID: %w", err)
	}

	var isRead []core.IsRead
	for _, id := range dialog.Participants {
		if id != request.Message.AuthorID {
			isRead = append(isRead, core.IsRead{Participant: id, IsRead: false})
		}
	}

	message := core.Message{
		Body:        request.Message.Body,
		AuthorID:    request.Message.AuthorID,
		IsRead:      isRead,
		ID:          request.Message.ID,
		Attachments: request.Message.Attachments,
		CreatedAt:   request.Message.CreatedAt,
	}

	if err := svc.db.ChatRepo.SendMessage(ctx, message, request.Message.DialogID); err != nil {
		return nil, fmt.Errorf("SendMessage: %w", err)
	}

	return &dto.SendMessageResponse{}, nil
}

func (svc *chatServiceImpl) ReadMessage(ctx context.Context, request *dto.ReadMessageRequest) (*dto.ReadMessageResponse, error) {

	if err := svc.db.ChatRepo.IsChatExist(ctx, request.Message.DialogID); err != nil {
		svc.log.Errorf("Chat not exist error: %s", err)
		return nil, err
	}

	if err := svc.db.ChatRepo.ReadMessage(ctx, request.Message.AuthorID, request.Message.Body, request.Message.DialogID); err != nil {
		svc.log.Errorf("SendMessage error: %s", err)
		return nil, err
	}

	svc.log.Debug("Message was read successful")
	return &dto.ReadMessageResponse{}, nil
}

func (svc *chatServiceImpl) GetDialogs(ctx context.Context, request *dto.GetDialogsRequest) (*dto.GetDialogsResponse, error) {
	ids, err := svc.db.UserRepo.GetUserDialogs(ctx, request.UserID)
	if err != nil {
		svc.log.Errorf("GetUserDialogs error: %s", err)
		return nil, err
	}

	ids, total, page := utils.GetLimitArray(&ids, request.Limit, request.Page)

	var dialogs []dto.Dialog
	for _, id := range ids {
		dInf, err := svc.db.ChatRepo.GetDialogByID(ctx, id)
		if err != nil {
			svc.log.Errorf("GetDialogInfo error: %s", err)
		}

		dialogs = append(dialogs, convert.Dialog2DTO(dInf, request.UserID))
	}

	for i, dialog := range dialogs {
		if len(dialog.Participants) == 1 {
			participant, err := svc.db.UserRepo.GetUserByID(ctx, dialog.Participants[0])
			if err != nil {
				return nil, err
			}
			dialogs[i].Name = participant.Name.Full()
		}
	}

	return &dto.GetDialogsResponse{Dialogs: dialogs, Total: total, AmountPages: page}, err
}

func (svc *chatServiceImpl) GetDialogByUserID(ctx context.Context, request *dto.GetDialogByUserIDRequest, currentUserID string) (*dto.GetDialogByUserIDResponse, error) {
	dialogID, err := svc.db.ChatRepo.IsDialogExist(ctx, request.UserID, currentUserID)
	if err != nil {
		svc.log.Errorf("IsDialogExist error: %s", err)
		return nil, err
	}
	return &dto.GetDialogByUserIDResponse{DialogID: dialogID}, nil
}

func (svc *chatServiceImpl) CheckDialog(ctx context.Context, request *dto.CheckDialogRequest) error {
	dialog, err := svc.db.ChatRepo.GetDialogByID(ctx, request.DialogID)
	if err != nil {
		svc.log.Errorf("GetDialogInfo error: %s", err)
		return err
	}

	svc.log.Info("Dialog: %s in User: %s", request.DialogID, request.UserID)
	err = svc.db.UserRepo.UserCheckDialog(ctx, dialog.ID, request.UserID)
	if err != nil {
		svc.log.Errorf("Don't found in db")
		return constants.ErrDBNotFound
	}
	return nil
}

func (svc *chatServiceImpl) GetDialog(ctx context.Context, request *dto.GetDialogRequest) (*dto.GetDialogResponse, error) {
	err := svc.db.UserRepo.UserCheckDialog(ctx, request.DialogID, request.UserID)
	if err != nil {
		return nil, constants.ErrDBNotFound
	}

	dialogCore, err := svc.db.ChatRepo.GetDialogByID(ctx, request.DialogID)
	if err != nil {
		return nil, err
	}

	var total int64
	var page int64
	dialogCore.Messages, total, page = utils.GetLimitMessage(&dialogCore.Messages, request.Limit, request.Page)

	dialog := convert.Dialog2DTO(dialogCore, request.UserID)
	if len(dialog.Participants) == 1 {
		participant, err := svc.db.UserRepo.GetUserByID(ctx, dialog.Participants[0])
		if err != nil {
			return nil, err
		}
		dialog.Name = participant.Name.Full()
	}

	return &dto.GetDialogResponse{Dialog: dialog, Messages: convert.Messages2DTO(dialogCore.Messages, request.UserID), Total: total, AmountPages: page}, err
}

func NewChatService(log *logrus.Entry, db *db.Repository) ChatService {
	return &chatServiceImpl{log: log, db: db}
}
