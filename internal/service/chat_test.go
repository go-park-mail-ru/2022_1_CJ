package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestCreateDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.CreateDialogRequest
	}

	type InputIsUniqDialog struct {
		firstUserID  string
		secondUserID string
	}
	type OutputIsUniqDialog struct {
		err error
	}

	type InputCreateDialog struct {
		userID    string
		name      string
		authorIDs []string
	}
	type OutputCreateDialog struct {
		dialog *core.Dialog
		err    error
	}

	type InputAddDialogForUserID struct {
		dialogID string
		userID   string
	}
	type OutputAddDialogForUserID struct {
		err error
	}
	type InputAddDialogForAuthorIDs struct {
		dialogID string
		userID   string
	}
	type OutputAddDialogForAuthorIDs struct {
		err error
	}

	type Output struct {
		res *dto.CreateDialogResponse
		err error
	}

	tests := []struct {
		name                        string
		input                       Input
		inputIsUniqDialog           InputIsUniqDialog
		outputIsUniqDialog          OutputIsUniqDialog
		inputCreateDialog           InputCreateDialog
		outputCreateDialog          OutputCreateDialog
		inputAddDialogForUserID     InputAddDialogForUserID
		outputAddDialogForUserID    OutputAddDialogForUserID
		inputAddDialogForAuthorIDs  InputAddDialogForAuthorIDs
		outputAddDialogForAuthorIDs OutputAddDialogForAuthorIDs
		output                      Output
	}{
		{
			name:   "AuthorIDs < 1",
			input:  Input{info: &dto.CreateDialogRequest{UserID: "0", Name: "nice chat", AuthorIDs: nil}},
			output: Output{nil, constants.ErrSingleChat},
		},
		{
			name:   "AuthorIDs == 1 and AuthorIDs[0] == request.UserID",
			input:  Input{info: &dto.CreateDialogRequest{UserID: "1", Name: "nice chat", AuthorIDs: []string{"1"}}},
			output: Output{nil, constants.ErrSingleChat},
		},
		{
			name:  "AuthorIDs == 1 and Uniq Dialog exist ",
			input: Input{info: &dto.CreateDialogRequest{UserID: "2", Name: "nice chat", AuthorIDs: []string{"3"}}},
			inputIsUniqDialog: InputIsUniqDialog{
				firstUserID:  "2",
				secondUserID: "3",
			},
			outputIsUniqDialog: OutputIsUniqDialog{err: constants.ErrDialogAlreadyExist},
			output:             Output{nil, constants.ErrDialogAlreadyExist},
		},
		{
			name:  "success",
			input: Input{info: &dto.CreateDialogRequest{UserID: "3", Name: "nice chat", AuthorIDs: []string{"4"}}},
			inputIsUniqDialog: InputIsUniqDialog{
				firstUserID:  "3",
				secondUserID: "4",
			},
			outputIsUniqDialog: OutputIsUniqDialog{err: nil},
			inputCreateDialog: InputCreateDialog{
				userID:    "3",
				name:      "nice chat",
				authorIDs: []string{"4"},
			},
			outputCreateDialog: OutputCreateDialog{
				dialog: &core.Dialog{ID: "5"},
				err:    nil,
			},
			inputAddDialogForUserID: InputAddDialogForUserID{
				dialogID: "5",
				userID:   "3",
			},
			outputAddDialogForUserID: OutputAddDialogForUserID{
				err: nil,
			},
			inputAddDialogForAuthorIDs: InputAddDialogForAuthorIDs{
				dialogID: "5",
				userID:   "4",
			},
			outputAddDialogForAuthorIDs: OutputAddDialogForAuthorIDs{
				err: nil,
			},
			output: Output{&dto.CreateDialogResponse{DialogID: "5"}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockChatR.EXPECT().IsUniqDialog(ctx, tests[2].inputIsUniqDialog.firstUserID,
			tests[2].inputIsUniqDialog.secondUserID).Return(tests[2].outputIsUniqDialog.err),

		testRepo.mockChatR.EXPECT().IsUniqDialog(ctx, tests[3].inputIsUniqDialog.firstUserID,
			tests[3].inputIsUniqDialog.secondUserID).Return(tests[3].outputIsUniqDialog.err),
		testRepo.mockChatR.EXPECT().CreateDialog(ctx, tests[3].inputCreateDialog.userID,
			tests[3].inputCreateDialog.name,
			tests[3].inputCreateDialog.authorIDs).Return(tests[3].outputCreateDialog.dialog, tests[3].outputCreateDialog.err),
		testRepo.mockUserR.EXPECT().AddDialog(ctx, tests[3].inputAddDialogForUserID.dialogID,
			tests[3].inputAddDialogForUserID.userID).Return(tests[3].outputAddDialogForUserID.err),
		testRepo.mockUserR.EXPECT().AddDialog(ctx, tests[3].inputAddDialogForAuthorIDs.dialogID,
			tests[3].inputAddDialogForAuthorIDs.userID).Return(tests[3].outputAddDialogForAuthorIDs.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := ChatService.CreateDialog(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.SendMessageRequest
	}

	type InputGetDialogByID struct {
		dialogID string
	}
	type OutputGetDialogByID struct {
		*core.Dialog
		err error
	}

	type InputSendMessage struct {
		message  core.Message
		dialogID string
	}
	type OutputSendMessage struct {
		err error
	}

	type Output struct {
		res *dto.SendMessageResponse
		err error
	}
	var err = errors.Errorf("Don't found in DB")
	tests := []struct {
		name                string
		input               Input
		inputGetDialogByID  InputGetDialogByID
		outputGetDialogByID OutputGetDialogByID
		inputSendMessage    InputSendMessage
		outputSendMessage   OutputSendMessage
		output              Output
	}{
		{
			name: "Don't found in DB",
			input: Input{info: &dto.SendMessageRequest{Message: dto.Message{
				ID:       "0",
				DialogID: "0",
				AuthorID: "0",
			}}},
			inputGetDialogByID: InputGetDialogByID{dialogID: "0"},
			outputGetDialogByID: OutputGetDialogByID{
				Dialog: nil,
				err:    err,
			},
			output: Output{nil, err},
		},
		{
			name: "success",
			input: Input{info: &dto.SendMessageRequest{Message: dto.Message{
				ID:       "1",
				DialogID: "1",
				AuthorID: "1",
				Body:     "hi",
				Event:    "send",
			}}},
			inputGetDialogByID: InputGetDialogByID{dialogID: "1"},
			outputGetDialogByID: OutputGetDialogByID{
				Dialog: &core.Dialog{
					ID:           "1",
					Participants: []string{"1", "2"},
				},
				err: nil,
			},
			inputSendMessage: InputSendMessage{
				message: core.Message{
					ID:       "1",
					AuthorID: "1",
					Body:     "hi",
					IsRead:   []core.IsRead{{Participant: "2", IsRead: false}},
				},
				dialogID: "1",
			},
			outputSendMessage: OutputSendMessage{err: nil},
			output:            Output{&dto.SendMessageResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockChatR.EXPECT().GetDialogByID(ctx, tests[0].inputGetDialogByID.dialogID).Return(tests[0].outputGetDialogByID.Dialog,
			tests[0].outputGetDialogByID.err),

		testRepo.mockChatR.EXPECT().GetDialogByID(ctx, tests[1].inputGetDialogByID.dialogID).Return(tests[1].outputGetDialogByID.Dialog,
			tests[1].outputGetDialogByID.err),
		testRepo.mockChatR.EXPECT().SendMessage(ctx, tests[1].inputSendMessage.message, tests[1].inputSendMessage.dialogID).Return(tests[1].outputSendMessage.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := ChatService.SendMessage(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestReadMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.ReadMessageRequest
	}

	type InputIsChatExist struct {
		dialogID string
	}
	type OutputIsChatExist struct {
		err error
	}

	type InputReadMessage struct {
		userID    string
		messageID string
		dialogID  string
	}
	type OutputReadMessage struct {
		err error
	}

	type Output struct {
		res *dto.ReadMessageResponse
		err error
	}

	tests := []struct {
		name              string
		input             Input
		inputIsChatExist  InputIsChatExist
		outputIsChatExist OutputIsChatExist
		inputReadMessage  InputReadMessage
		outputReadMessage OutputReadMessage
		output            Output
	}{
		{
			name: "Don't found in DB",
			input: Input{info: &dto.ReadMessageRequest{
				Message: dto.Message{
					ID:       "0",
					DialogID: "0",
					Event:    "read",
					AuthorID: "0",
					DestinID: "",
					Body:     "0",
				},
			}},
			inputIsChatExist:  InputIsChatExist{dialogID: "0"},
			outputIsChatExist: OutputIsChatExist{err: mongo.ErrNoDocuments},
			output:            Output{nil, mongo.ErrNoDocuments},
		},
		{
			name: "success",
			input: Input{info: &dto.ReadMessageRequest{
				Message: dto.Message{
					ID:       "1",
					DialogID: "1",
					Event:    "read",
					AuthorID: "1",
					DestinID: "",
					Body:     "1",
				},
			}},
			inputIsChatExist:  InputIsChatExist{dialogID: "1"},
			outputIsChatExist: OutputIsChatExist{nil},
			inputReadMessage: InputReadMessage{
				userID:    "1",
				messageID: "1",
				dialogID:  "1",
			},
			outputReadMessage: OutputReadMessage{err: nil},
			output:            Output{&dto.ReadMessageResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockChatR.EXPECT().IsChatExist(ctx, tests[0].inputIsChatExist.dialogID).Return(tests[0].outputIsChatExist.err),

		testRepo.mockChatR.EXPECT().IsChatExist(ctx, tests[1].inputIsChatExist.dialogID).Return(tests[1].outputIsChatExist.err),
		testRepo.mockChatR.EXPECT().ReadMessage(ctx, tests[1].inputReadMessage.userID, tests[1].inputReadMessage.messageID, tests[1].inputReadMessage.dialogID).Return(tests[1].outputReadMessage.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := ChatService.ReadMessage(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetDialogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetDialogsRequest
	}

	type InputGetUserDialogs struct {
		postID string
	}
	type OutputGetUserDialogs struct {
		ids []string
		err error
	}

	type InputGetDialogByID struct {
		dialogID string
	}
	type OutputGetDialogByID struct {
		dialog *core.Dialog
		err    error
	}

	type Output struct {
		res *dto.GetDialogsResponse
		err error
	}

	tests := []struct {
		name                 string
		input                Input
		inputGetUserDialogs  InputGetUserDialogs
		outputGetUserDialogs OutputGetUserDialogs
		inputGetDialogByID   InputGetDialogByID
		outputGetDialogByID  OutputGetDialogByID
		output               Output
	}{
		{
			name:                 "Don't found in DB",
			input:                Input{info: &dto.GetDialogsRequest{UserID: "0"}},
			inputGetUserDialogs:  InputGetUserDialogs{postID: "0"},
			outputGetUserDialogs: OutputGetUserDialogs{ids: nil, err: mongo.ErrNoDocuments},
			output:               Output{nil, mongo.ErrNoDocuments},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserDialogs(ctx, tests[0].inputGetUserDialogs.postID).Return(tests[0].outputGetUserDialogs.ids,
			tests[0].outputGetUserDialogs.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := ChatService.GetDialogs(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestCheckDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.CheckDialogRequest
	}

	type InputGetDialogByID struct {
		dialogID string
	}
	type OutputGetDialogByID struct {
		dialog *core.Dialog
		err    error
	}

	type InputUserCheckDialog struct {
		dialogID string
		userID   string
	}

	type OutputUserCheckDialog struct {
		err error
	}

	type Output struct {
		err error
	}

	tests := []struct {
		name                  string
		input                 Input
		inputGetDialogByID    InputGetDialogByID
		outputGetDialogByID   OutputGetDialogByID
		inputUserCheckDialog  InputUserCheckDialog
		outputUserCheckDialog OutputUserCheckDialog
		output                Output
	}{
		{
			name:               "Don't found in DB",
			input:              Input{info: &dto.CheckDialogRequest{UserID: "0", DialogID: "0"}},
			inputGetDialogByID: InputGetDialogByID{dialogID: "0"},
			outputGetDialogByID: OutputGetDialogByID{
				dialog: nil,
				err:    mongo.ErrNoDocuments,
			},
			output: Output{mongo.ErrNoDocuments},
		},
		{
			name:               "Don't found post in dialog",
			input:              Input{info: &dto.CheckDialogRequest{UserID: "1", DialogID: "1"}},
			inputGetDialogByID: InputGetDialogByID{dialogID: "1"},
			outputGetDialogByID: OutputGetDialogByID{
				dialog: &core.Dialog{
					ID:   "1",
					Name: "chat",
				},
				err: nil,
			},
			inputUserCheckDialog: InputUserCheckDialog{
				dialogID: "1",
				userID:   "1",
			},
			outputUserCheckDialog: OutputUserCheckDialog{err: mongo.ErrNoDocuments},
			output:                Output{constants.ErrDBNotFound},
		},
		{
			name:               "success",
			input:              Input{info: &dto.CheckDialogRequest{UserID: "2", DialogID: "2"}},
			inputGetDialogByID: InputGetDialogByID{dialogID: "2"},
			outputGetDialogByID: OutputGetDialogByID{
				dialog: &core.Dialog{
					ID:   "2",
					Name: "chat",
				},
				err: nil,
			},
			inputUserCheckDialog: InputUserCheckDialog{
				dialogID: "2",
				userID:   "2",
			},
			outputUserCheckDialog: OutputUserCheckDialog{err: nil},
			output:                Output{nil},
		},
	}

	gomock.InOrder(
		testRepo.mockChatR.EXPECT().GetDialogByID(ctx, tests[0].inputGetDialogByID.dialogID).Return(tests[0].outputGetDialogByID.dialog, tests[0].outputGetDialogByID.err),

		testRepo.mockChatR.EXPECT().GetDialogByID(ctx, tests[1].inputGetDialogByID.dialogID).Return(tests[1].outputGetDialogByID.dialog, tests[1].outputGetDialogByID.err),
		testRepo.mockUserR.EXPECT().UserCheckDialog(ctx, tests[1].inputUserCheckDialog.dialogID, tests[1].inputUserCheckDialog.userID).Return(tests[1].outputUserCheckDialog.err),

		testRepo.mockChatR.EXPECT().GetDialogByID(ctx, tests[2].inputGetDialogByID.dialogID).Return(tests[2].outputGetDialogByID.dialog, tests[2].outputGetDialogByID.err),
		testRepo.mockUserR.EXPECT().UserCheckDialog(ctx, tests[2].inputUserCheckDialog.dialogID, tests[2].inputUserCheckDialog.userID).Return(tests[2].outputUserCheckDialog.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			errRes := ChatService.CheckDialog(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewChatService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetDialogRequest
	}

	type InputUserCheckDialog struct {
		dialogID string
		postID   string
	}

	type OutputUserCheckDialog struct {
		err error
	}

	type InputGetDialogByID struct {
		dialogID string
		err      error
	}

	type OutputGetDialogByID struct {
		dialog *core.Dialog
		err    error
	}

	type Output struct {
		res *dto.GetDialogResponse
		err error
	}

	tests := []struct {
		name                  string
		input                 Input
		inputUserCheckDialog  InputUserCheckDialog
		outputUserCheckDialog OutputUserCheckDialog
		inputGetDialogByID    InputGetDialogByID
		outputGetDialogByID   OutputGetDialogByID
		output                Output
	}{
		{
			name:                 "Don't found in DB",
			input:                Input{info: &dto.GetDialogRequest{UserID: "0", DialogID: "0"}},
			inputUserCheckDialog: InputUserCheckDialog{dialogID: "0", postID: "0"},
			outputUserCheckDialog: OutputUserCheckDialog{
				err: constants.ErrDBNotFound,
			},
			output: Output{res: nil, err: constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckDialog(ctx, tests[0].inputUserCheckDialog.dialogID, tests[0].inputUserCheckDialog.postID).Return(tests[0].outputUserCheckDialog.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := ChatService.GetDialog(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
