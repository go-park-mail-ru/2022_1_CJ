package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.SendFriendRequestRequest
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

	type InputCreateRequest struct {
		userID   string
		personID string
	}

	type OutputCreateRequest struct {
		err error
	}

	type Output struct {
		res *dto.SendFriendRequestResponse
		err error
	}

	var tests = []struct {
		name                string
		input               Input
		inputIsUniqRequest  InputIsUniqRequest
		outputIsUniqRequest OutputIsUniqRequest
		inputIsNotFriend    InputIsNotFriend
		outputIsNotFriend   OutputIsNotFriend
		inputCreateRequest  InputCreateRequest
		outputCreateRequest OutputCreateRequest
		output              Output
	}{
		{
			name:                "Success send request",
			input:               Input{info: &dto.SendFriendRequestRequest{From: "1", To: "2"}},
			inputIsUniqRequest:  InputIsUniqRequest{userID: "2", personID: "1"},
			outputIsUniqRequest: OutputIsUniqRequest{err: nil},
			inputIsNotFriend:    InputIsNotFriend{userID: "2", personID: "1"},
			outputIsNotFriend:   OutputIsNotFriend{err: nil},
			inputCreateRequest:  InputCreateRequest{userID: "2", personID: "1"},
			outputCreateRequest: OutputCreateRequest{err: nil},
			output:              Output{res: &dto.SendFriendRequestResponse{}, err: nil},
		},
		{
			name:                "Double send request error",
			input:               Input{info: &dto.SendFriendRequestRequest{From: "3", To: "4"}},
			inputIsUniqRequest:  InputIsUniqRequest{userID: "4", personID: "3"},
			outputIsUniqRequest: OutputIsUniqRequest{err: constants.ErrRequestAlreadyExist},
			output:              Output{res: nil, err: constants.ErrRequestAlreadyExist},
		},
		{
			name:                "Already friends",
			input:               Input{info: &dto.SendFriendRequestRequest{From: "5", To: "6"}},
			inputIsUniqRequest:  InputIsUniqRequest{userID: "6", personID: "5"},
			outputIsUniqRequest: OutputIsUniqRequest{err: nil},
			inputIsNotFriend:    InputIsNotFriend{userID: "6", personID: "5"},
			outputIsNotFriend:   OutputIsNotFriend{err: constants.ErrAlreadyFriends},
			output:              Output{res: nil, err: constants.ErrAlreadyFriends},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[0].inputIsUniqRequest.personID, tests[0].inputIsUniqRequest.userID).Return(tests[0].outputIsUniqRequest.err),
		testRepo.mockFriendsR.EXPECT().IsNotFriend(ctx, tests[0].inputIsNotFriend.personID, tests[0].inputIsNotFriend.userID).Return(tests[0].outputIsNotFriend.err),
		testRepo.mockFriendsR.EXPECT().CreateRequest(ctx, tests[0].inputCreateRequest.personID, tests[0].inputCreateRequest.userID).Return(tests[0].outputCreateRequest.err),

		//second
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[1].inputIsUniqRequest.personID, tests[1].inputIsUniqRequest.userID).Return(tests[1].outputIsUniqRequest.err),

		//third
		testRepo.mockFriendsR.EXPECT().IsUniqRequest(ctx, tests[2].inputIsUniqRequest.personID, tests[2].inputIsUniqRequest.userID).Return(tests[2].outputIsUniqRequest.err),
		testRepo.mockFriendsR.EXPECT().IsNotFriend(ctx, tests[2].inputIsNotFriend.personID, tests[2].inputIsNotFriend.userID).Return(tests[2].outputIsNotFriend.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.SendRequest(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestAcceptRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.AcceptFriendRequestRequest
	}

	type InputDeleteRequest struct {
		userID   string
		personID string
	}

	type OutputDeleteRequest struct {
		err error
	}

	type InputMakeFriends struct {
		userID   string
		personID string
	}

	type OutputMakeFriends struct {
		err error
	}

	type Output struct {
		res *dto.AcceptFriendRequestResponse
		err error
	}

	var tests = []struct {
		name                string
		input               Input
		inputDeleteRequest  InputDeleteRequest
		outputDeleteRequest OutputDeleteRequest
		inputMakeFriends    InputMakeFriends
		outputMakeFriends   OutputMakeFriends
		output              Output
	}{
		{
			name:                "Success",
			input:               Input{info: &dto.AcceptFriendRequestRequest{From: "1", To: "2"}},
			inputMakeFriends:    InputMakeFriends{userID: "1", personID: "2"},
			outputMakeFriends:   OutputMakeFriends{err: nil},
			inputDeleteRequest:  InputDeleteRequest{userID: "1", personID: "2"},
			outputDeleteRequest: OutputDeleteRequest{err: nil},
			output:              Output{res: &dto.AcceptFriendRequestResponse{}, err: nil},
		},
	}

	gomock.InOrder(
		// second
		testRepo.mockFriendsR.EXPECT().DeleteRequest(ctx, tests[0].inputDeleteRequest.userID, tests[0].inputDeleteRequest.personID).Return(tests[0].outputDeleteRequest.err),
		testRepo.mockFriendsR.EXPECT().MakeFriends(ctx, tests[0].inputMakeFriends.userID, tests[0].inputMakeFriends.personID).Return(tests[0].outputMakeFriends.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.AcceptRequest(dbUserImpl, ctx, test.input.info)
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
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.DeleteFriendRequest
	}

	type InputDeleteFriend struct {
		userID   string
		personID string
	}

	type OutputDeleteFriend struct {
		err error
	}

	type InputCreateRequest struct {
		userID   string
		personID string
	}

	type OutputCreateRequest struct {
		err error
	}

	type Output struct {
		res *dto.DeleteFriendResponse
		err error
	}

	var err = errors.Errorf("Can't delete friends")
	var tests = []struct {
		name                string
		input               Input
		inputDeleteFriend   InputDeleteFriend
		outputDeleteFriend  OutputDeleteFriend
		inputCreateRequest  InputCreateRequest
		outputCreateRequest OutputCreateRequest
		output              Output
	}{
		{
			name:               "Error delete friends",
			input:              Input{info: &dto.DeleteFriendRequest{UserID: "1", FriendID: "2"}},
			inputDeleteFriend:  InputDeleteFriend{userID: "1", personID: "2"},
			outputDeleteFriend: OutputDeleteFriend{err: err},
			output:             Output{res: nil, err: err},
		},
		{
			name:                "Success",
			input:               Input{info: &dto.DeleteFriendRequest{UserID: "3", FriendID: "4"}},
			inputDeleteFriend:   InputDeleteFriend{userID: "3", personID: "4"},
			outputDeleteFriend:  OutputDeleteFriend{err: nil},
			inputCreateRequest:  InputCreateRequest{userID: "4", personID: "3"},
			outputCreateRequest: OutputCreateRequest{err: nil},
			output:              Output{res: &dto.DeleteFriendResponse{}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().DeleteFriend(ctx, tests[0].inputDeleteFriend.userID, tests[0].inputDeleteFriend.personID).Return(tests[0].outputDeleteFriend.err),

		// second
		testRepo.mockFriendsR.EXPECT().DeleteFriend(ctx, tests[1].inputDeleteFriend.userID, tests[1].inputDeleteFriend.personID).Return(tests[1].outputDeleteFriend.err),
		testRepo.mockFriendsR.EXPECT().CreateRequest(ctx, tests[1].inputCreateRequest.userID, tests[1].inputCreateRequest.personID).Return(tests[1].outputCreateRequest.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.DeleteFriend(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID *dto.GetFriendsRequest
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
			input:                    Input{userID: &dto.GetFriendsRequest{UserID: "1"}},
			inputGetFriendsByUserID:  InputGetFriendsByUserID{userID: "1"},
			outputGetFriendsByUserID: OutputGetFriendsByUserID{friends: testFriends, err: nil},
			output:                   Output{res: &dto.GetFriendsResponse{FriendIDs: testFriends}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetFriends(ctx, tests[0].inputGetFriendsByUserID.userID).Return(tests[0].outputGetFriendsByUserID.friends, tests[0].outputGetFriendsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.GetFriends(dbUserImpl, ctx, test.input.userID)
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
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID *dto.GetOutcomingRequestsRequest
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
			input:                              Input{userID: &dto.GetOutcomingRequestsRequest{UserID: "1"}},
			inputGetOutcomingRequestsByUserID:  InputGetOutcomingRequestsByUserID{userID: "1"},
			outputGetOutcomingRequestsByUserID: OutputGetOutcomingRequestsByUserID{requests: testRequests, err: nil},
			output:                             Output{res: &dto.GetOutcomingRequestsResponse{RequestIDs: testRequests}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetOutcomingRequests(ctx, tests[0].inputGetOutcomingRequestsByUserID.userID).Return(tests[0].outputGetOutcomingRequestsByUserID.requests, tests[0].outputGetOutcomingRequestsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.GetOutcomingRequests(dbUserImpl, ctx, test.input.userID)
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
	dbUserImpl := NewFriendsService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID *dto.GetIncomingRequestsRequest
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
			input:                             Input{userID: &dto.GetIncomingRequestsRequest{UserID: "1"}},
			inputGetIncomingRequestsByUserID:  InputGetIncomingRequestsByUserID{userID: "1"},
			outputGetIncomingRequestsByUserID: OutputGetIncomingRequestsByUserID{requests: testRequests, err: nil},
			output:                            Output{res: &dto.GetIncomingRequestsResponse{RequestIDs: testRequests}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockFriendsR.EXPECT().GetIncomingRequests(ctx, tests[0].inputGetIncomingRequestsByUserID.userID).Return(tests[0].outputGetIncomingRequestsByUserID.requests, tests[0].outputGetIncomingRequestsByUserID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, errRes := FriendsService.GetIncomingRequests(dbUserImpl, ctx, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
