package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendFriendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.SendFriendRequestRequest
		userID string
	}

	type InputIsUniqRequest struct {
		userID   string
		personID string
	}

	type OutputIsUniqRequest struct {
		err error
	}

	type InputIsNotFriend struct {
		userID   string
		personID string
	}

	type OutputIsNotFriend struct {
		err error
	}

	type InputMakeOutcomingRequest struct {
		userID   string
		personID string
	}

	type OutputMakeOutcomingRequest struct {
		err error
	}

	type InputMakeIncomingRequest struct {
		userID   string
		personID string
	}

	type OutputMakeIncomingRequest struct {
		err error
	}

	type Output struct {
		res *dto.SendFriendRequestResponse
		err error
	}

	var tests = []struct {
		name                       string
		input                      Input
		inputIsUniqRequest         InputIsUniqRequest
		outputIsUniqRequest        OutputIsUniqRequest
		inputIsNotFriend           InputIsNotFriend
		outputIsNotFriend          OutputIsNotFriend
		inputMakeOutcomingRequest  InputMakeOutcomingRequest
		outputMakeOutcomingRequest OutputMakeOutcomingRequest
		inputMakeIncomingRequest   InputMakeIncomingRequest
		outputMakeIncomingRequest  OutputMakeIncomingRequest
		output                     Output
	}{
		{
			name:                       "Success send request",
			input:                      Input{info: &dto.SendFriendRequestRequest{UserID: "2"}, userID: "1"},
			inputIsUniqRequest:         InputIsUniqRequest{userID: "2", personID: "1"},
			outputIsUniqRequest:        OutputIsUniqRequest{err: nil},
			inputIsNotFriend:           InputIsNotFriend{userID: "2", personID: "1"},
			outputIsNotFriend:          OutputIsNotFriend{err: nil},
			inputMakeOutcomingRequest:  InputMakeOutcomingRequest{userID: "2", personID: "1"},
			outputMakeOutcomingRequest: OutputMakeOutcomingRequest{err: nil},
			inputMakeIncomingRequest:   InputMakeIncomingRequest{userID: "2", personID: "1"},
			outputMakeIncomingRequest:  OutputMakeIncomingRequest{err: nil},
			output:                     Output{res: &dto.SendFriendRequestResponse{}, err: nil},
		},
		{
			name:                "Double send request error",
			input:               Input{info: &dto.SendFriendRequestRequest{UserID: "3"}, userID: "4"},
			inputIsUniqRequest:  InputIsUniqRequest{userID: "3", personID: "4"},
			outputIsUniqRequest: OutputIsUniqRequest{err: constants.ErrRequestAlreadyExist},
			output:              Output{res: nil, err: constants.ErrRequestAlreadyExist},
		},
		{
			name:                "Already friends",
			input:               Input{info: &dto.SendFriendRequestRequest{UserID: "5"}, userID: "6"},
			inputIsUniqRequest:  InputIsUniqRequest{userID: "5", personID: "6"},
			outputIsUniqRequest: OutputIsUniqRequest{err: nil},
			inputIsNotFriend:    InputIsNotFriend{userID: "5", personID: "6"},
			outputIsNotFriend:   OutputIsNotFriend{err: constants.ErrAlreadyFriends},
			output:              Output{res: nil, err: constants.ErrAlreadyFriends},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[0].inputIsUniqRequest.personID, tests[0].inputIsUniqRequest.userID).Return(tests[0].outputIsUniqRequest.err),
		testRepo.mockFriendsR.EXPECT().IsNotFriend(ctx, tests[0].inputIsNotFriend.personID, tests[0].inputIsNotFriend.userID).Return(tests[0].outputIsNotFriend.err),
		testRepo.mockFriendsR.EXPECT().MakeOutcomingRequest(ctx, tests[0].inputMakeOutcomingRequest.personID, tests[0].inputMakeOutcomingRequest.userID).Return(tests[0].outputMakeOutcomingRequest.err),
		testRepo.mockFriendsR.EXPECT().MakeIncomingRequest(ctx, tests[0].inputMakeIncomingRequest.personID, tests[0].inputMakeIncomingRequest.userID).Return(tests[0].outputMakeIncomingRequest.err),

		//second
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[1].inputIsUniqRequest.personID, tests[1].inputIsUniqRequest.userID).Return(tests[1].outputIsUniqRequest.err),

		//third
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[2].inputIsUniqRequest.personID, tests[2].inputIsUniqRequest.userID).Return(tests[2].outputIsUniqRequest.err),
		testRepo.mockFriendsR.EXPECT().IsNotFriend(ctx, tests[2].inputIsNotFriend.personID, tests[2].inputIsNotFriend.userID).Return(tests[2].outputIsNotFriend.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.SendFriendRequest(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestAcceptFriendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.AcceptFriendRequestRequest
		userID string
	}

	type InputMakeFriends struct {
		userID   string
		personID string
	}

	type OutputMakeFriends struct {
		err error
	}

	type InputDeleteOutcomingRequest struct {
		userID   string
		personID string
	}

	type OutputDeleteOutcomingRequest struct {
		err error
	}

	type InputDeleteIncomingRequest struct {
		userID   string
		personID string
	}

	type OutputDeleteIncomingRequest struct {
		err error
	}

	type InputGetOutcomingRequestsByUserID struct {
		userID string
	}

	type OutputGetOutcomingRequestsByUserID struct {
		res []string
		err error
	}

	type Output struct {
		res *dto.AcceptFriendRequestResponse
		err error
	}

	var tests = []struct {
		name                               string
		input                              Input
		inputMakeFriends                   InputMakeFriends
		outputMakeFriends                  OutputMakeFriends
		inputDeleteOutcomingRequest        InputDeleteOutcomingRequest
		outputDeleteOutcomingRequest       OutputDeleteOutcomingRequest
		inputDeleteIncomingRequest         InputDeleteIncomingRequest
		outputDeleteIncomingRequest        OutputDeleteIncomingRequest
		inputGetOutcomingRequestsByUserID  InputGetOutcomingRequestsByUserID
		outputGetOutcomingRequestsByUserID OutputGetOutcomingRequestsByUserID
		output                             Output
	}{
		{
			name:   "Error make yourself friend",
			input:  Input{info: &dto.AcceptFriendRequestRequest{UserID: "1", IsAccepted: true}, userID: "1"},
			output: Output{res: nil, err: constants.ErrAddYourself},
		},
		{
			name:                               "Success",
			input:                              Input{info: &dto.AcceptFriendRequestRequest{UserID: "5", IsAccepted: true}, userID: "2"},
			inputMakeFriends:                   InputMakeFriends{userID: "2", personID: "5"},
			outputMakeFriends:                  OutputMakeFriends{err: nil},
			inputDeleteOutcomingRequest:        InputDeleteOutcomingRequest{userID: "2", personID: "5"},
			outputDeleteIncomingRequest:        OutputDeleteIncomingRequest{err: nil},
			inputDeleteIncomingRequest:         InputDeleteIncomingRequest{userID: "2", personID: "5"},
			outputDeleteOutcomingRequest:       OutputDeleteOutcomingRequest{err: nil},
			inputGetOutcomingRequestsByUserID:  InputGetOutcomingRequestsByUserID{userID: "2"},
			outputGetOutcomingRequestsByUserID: OutputGetOutcomingRequestsByUserID{res: []string{"123"}, err: nil},
			output:                             Output{res: &dto.AcceptFriendRequestResponse{RequestsID: []string{"123"}}, err: nil},
		},
	}

	gomock.InOrder(
		// second
		testRepo.mockFriendsR.EXPECT().MakeFriends(ctx, tests[1].inputMakeFriends.userID, tests[1].inputMakeFriends.personID).Return(tests[1].outputMakeFriends.err),
		testRepo.mockFriendsR.EXPECT().DeleteOutcomingRequest(ctx, tests[1].inputDeleteOutcomingRequest.userID, tests[1].inputDeleteOutcomingRequest.personID).Return(tests[1].outputDeleteOutcomingRequest.err),
		testRepo.mockFriendsR.EXPECT().DeleteIncomingRequest(ctx, tests[1].inputDeleteIncomingRequest.userID, tests[1].inputDeleteIncomingRequest.personID).Return(tests[1].outputDeleteIncomingRequest.err),
		testRepo.mockFriendsR.EXPECT().GetOutcomingRequestsByUserID(ctx, tests[1].inputGetOutcomingRequestsByUserID.userID).Return(tests[1].outputGetOutcomingRequestsByUserID.res, tests[1].outputGetOutcomingRequestsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.AcceptFriendRequest(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestDeleteFriend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.DeleteFriendRequest
		userID string
	}

	type InputDeleteFriend struct {
		userID   string
		personID string
	}

	type OutputDeleteFriend struct {
		err error
	}

	type InputGetFriendsByUserID struct {
		userID string
	}

	type OutputGetFriendsByUserID struct {
		friends []string
		err     error
	}

	type Output struct {
		res *dto.DeleteFriendResponse
		err error
	}

	var err = errors.Errorf("Can't delete friends")
	var tests = []struct {
		name                     string
		input                    Input
		inputDeleteFriend        InputDeleteFriend
		outputDeleteFriend       OutputDeleteFriend
		inputGetFriendsByUserID  InputGetFriendsByUserID
		outputGetFriendsByUserID OutputGetFriendsByUserID
		output                   Output
	}{
		{
			name:               "Error delete friends",
			input:              Input{info: &dto.DeleteFriendRequest{ExFriendID: "123"}, userID: "345"},
			inputDeleteFriend:  InputDeleteFriend{userID: "345", personID: "123"},
			outputDeleteFriend: OutputDeleteFriend{err: err},
			output:             Output{res: nil, err: err},
		},
		{
			name:                     "Success",
			input:                    Input{info: &dto.DeleteFriendRequest{ExFriendID: "456"}, userID: "787"},
			inputDeleteFriend:        InputDeleteFriend{userID: "787", personID: "456"},
			outputDeleteFriend:       OutputDeleteFriend{err: nil},
			inputGetFriendsByUserID:  InputGetFriendsByUserID{userID: "787"},
			outputGetFriendsByUserID: OutputGetFriendsByUserID{friends: []string{"12345", "32345"}, err: nil},
			output:                   Output{res: &dto.DeleteFriendResponse{FriendsID: []string{"12345", "32345"}}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().DeleteFriend(ctx, tests[0].inputDeleteFriend.userID, tests[0].inputDeleteFriend.personID).Return(tests[0].outputDeleteFriend.err),

		// second
		testRepo.mockFriendsR.EXPECT().DeleteFriend(ctx, tests[1].inputDeleteFriend.userID, tests[1].inputDeleteFriend.personID).Return(tests[1].outputDeleteFriend.err),
		testRepo.mockFriendsR.EXPECT().GetFriendsByUserID(ctx, tests[1].inputGetFriendsByUserID.userID).Return(tests[1].outputGetFriendsByUserID.friends, tests[1].outputGetFriendsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.DeleteFriend(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetFriendsByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID string
	}

	type InputGetFriendsByUserID struct {
		userID string
	}

	type OutputGetFriendsByUserID struct {
		friends []string
		err     error
	}

	type Output struct {
		res *dto.GetFriendsResponse
		err error
	}

	testFriends := make([]string, 2)
	testFriends[0] = "12345"
	testFriends[1] = "123456"
	var tests = []struct {
		name                     string
		input                    Input
		inputGetFriendsByUserID  InputGetFriendsByUserID
		outputGetFriendsByUserID OutputGetFriendsByUserID
		output                   Output
	}{
		{
			name:                     "Success",
			input:                    Input{userID: "787"},
			inputGetFriendsByUserID:  InputGetFriendsByUserID{userID: "787"},
			outputGetFriendsByUserID: OutputGetFriendsByUserID{friends: testFriends, err: nil},
			output:                   Output{res: &dto.GetFriendsResponse{FriendsID: testFriends}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetFriendsByUserID(ctx, tests[0].inputGetFriendsByUserID.userID).Return(tests[0].outputGetFriendsByUserID.friends, tests[0].outputGetFriendsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.GetFriendsByUserID(dbUserImpl, ctx, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetOutcomingRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID string
	}

	type InputGetOutcomingRequestsByUserID struct {
		userID string
	}

	type OutputGetOutcomingRequestsByUserID struct {
		requests []string
		err      error
	}

	type Output struct {
		res *dto.GetOutcomingRequestsResponse
		err error
	}

	testRequests := make([]string, 2)
	testRequests[0] = "12345"
	testRequests[1] = "123456"
	var tests = []struct {
		name                               string
		input                              Input
		inputGetOutcomingRequestsByUserID  InputGetOutcomingRequestsByUserID
		outputGetOutcomingRequestsByUserID OutputGetOutcomingRequestsByUserID
		output                             Output
	}{
		{
			name:                               "Success",
			input:                              Input{userID: "557"},
			inputGetOutcomingRequestsByUserID:  InputGetOutcomingRequestsByUserID{userID: "557"},
			outputGetOutcomingRequestsByUserID: OutputGetOutcomingRequestsByUserID{requests: testRequests, err: nil},
			output:                             Output{res: &dto.GetOutcomingRequestsResponse{RequestIDs: testRequests}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetOutcomingRequestsByUserID(ctx, tests[0].inputGetOutcomingRequestsByUserID.userID).Return(tests[0].outputGetOutcomingRequestsByUserID.requests, tests[0].outputGetOutcomingRequestsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.GetOutcomingRequests(dbUserImpl, ctx, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetIncomingRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID string
	}

	type InputGetIncomingRequestsByUserID struct {
		userID string
	}

	type OutputGetIncomingRequestsByUserID struct {
		requests []string
		err      error
	}

	type Output struct {
		res *dto.GetIncomingRequestsResponse
		err error
	}

	testRequests := make([]string, 2)
	testRequests[0] = "12345"
	testRequests[1] = "123456"
	var tests = []struct {
		name                              string
		input                             Input
		inputGetIncomingRequestsByUserID  InputGetIncomingRequestsByUserID
		outputGetIncomingRequestsByUserID OutputGetIncomingRequestsByUserID
		output                            Output
	}{
		{
			name:                              "Success",
			input:                             Input{userID: "557"},
			inputGetIncomingRequestsByUserID:  InputGetIncomingRequestsByUserID{userID: "557"},
			outputGetIncomingRequestsByUserID: OutputGetIncomingRequestsByUserID{requests: testRequests, err: nil},
			output:                            Output{res: &dto.GetIncomingRequestsResponse{RequestIDs: testRequests}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetIncomingRequestsByUserID(ctx, tests[0].inputGetIncomingRequestsByUserID.userID).Return(tests[0].outputGetIncomingRequestsByUserID.requests, tests[0].outputGetIncomingRequestsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := service.FriendsService.GetIncomingRequests(dbUserImpl, ctx, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
