package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	tests := []struct {
		name              string
		input             string
		resultGetUserByID error
		output            error
	}{
		{
			name:              "Don't found in BD",
			input:             "0",
			resultGetUserByID: constants.ErrDBNotFound,
			output:            constants.ErrDBNotFound,
		},
		{
			name:              "Found in BD",
			input:             "677be1d2-9b64-48e9-9341-5ba0c2f57686",
			resultGetUserByID: nil,
			output:            nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gomock.InOrder(
				testRepo.mockUserR.EXPECT().GetUserByID(ctx, test.input).Return(&core.User{}, test.resultGetUserByID),
			)
			_, res := service.UserService.GetUserData(dbUserImpl, ctx, test.input)
			if !assert.Equal(t, test.output, res) {
				t.Error("got : ", res, " expected :", test.output)
			}
		})
	}
}

func TestGetUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	tests := []struct {
		name                 string
		input                string
		resultGetUserByID    error
		resultGetPostsByUser error
		output               error
	}{
		{
			name:                 "Don't found in BD",
			input:                "0",
			resultGetUserByID:    constants.ErrDBNotFound,
			resultGetPostsByUser: constants.ErrDBNotFound,
			output:               constants.ErrDBNotFound,
		},
		{
			name:                 "Found in BD",
			input:                "677be1d2-9b64-48e9-9341-5ba0c2f57686",
			resultGetUserByID:    nil,
			resultGetPostsByUser: nil,
			output:               nil,
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].input).Return(&core.User{}, tests[0].resultGetUserByID),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].input).Return(&core.User{}, tests[1].resultGetUserByID),
		testRepo.mockUserR.EXPECT().GetPostsByUser(ctx, tests[1].input).Return([]string{}, tests[1].resultGetPostsByUser),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			_, res := service.UserService.GetUserPosts(dbUserImpl, ctx, test.input)
			if !assert.Equal(t, test.output, res) {
				t.Error("got : ", res, " expected :", test.output)
			}
		})
	}
}

func TestGetFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	tests := []struct {
		name              string
		input             string
		resultGetUserByID error
		resultGetFeed     error
		output            error
	}{
		{
			name:              "Don't found in BD",
			input:             "0",
			resultGetUserByID: constants.ErrDBNotFound,
			resultGetFeed:     nil,
			output:            constants.ErrDBNotFound,
		},
		{
			name:              "Found in BD",
			input:             "677be1d2-9b64-48e9-9341-5ba0c2f57686",
			resultGetUserByID: nil,
			resultGetFeed:     nil,
			output:            nil,
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].input).Return(&core.User{}, tests[0].resultGetUserByID),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].input).Return(&core.User{}, tests[1].resultGetUserByID),
		testRepo.mockPostR.EXPECT().GetFeed(ctx, tests[1].input).Return([]string{}, tests[1].resultGetFeed),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			_, res := service.UserService.GetFeed(dbUserImpl, ctx, test.input)
			if !assert.Equal(t, test.output, res) {
				t.Error("got : ", res, " expected :", test.output)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()
	type GetUserByID struct {
		user *core.User
		err  error
	}
	type FriendsByID struct {
		friends []string
		err     error
	}
	type Output struct {
		res *dto.GetProfileResponse
		err error
	}
	tests := []struct {
		name                 string
		input                *dto.GetProfileRequest
		resultGetUserByID    GetUserByID
		resultGetFriendsByID FriendsByID
		output               Output
	}{
		{
			name:                 "Don't found in BD",
			input:                &dto.GetProfileRequest{UserID: "0"},
			resultGetUserByID:    GetUserByID{&core.User{}, constants.ErrDBNotFound},
			resultGetFriendsByID: FriendsByID{nil, nil},
			output:               Output{nil, constants.ErrDBNotFound},
		},
		{
			name:                 "Found in BD but not in friends Repo",
			input:                &dto.GetProfileRequest{UserID: "1"},
			resultGetUserByID:    GetUserByID{&core.User{FriendsID: "0"}, nil},
			resultGetFriendsByID: FriendsByID{nil, constants.ErrDBNotFound},
			output:               Output{nil, constants.ErrDBNotFound},
		},
		{
			name:                 "Found in BD and in friendsRepo",
			input:                &dto.GetProfileRequest{UserID: "2"},
			resultGetUserByID:    GetUserByID{&core.User{FriendsID: "3"}, nil},
			resultGetFriendsByID: FriendsByID{[]string{"4"}, nil},
			output: Output{&dto.GetProfileResponse{
				UserProfile: convert.Profile2DTO(
					&core.User{FriendsID: "3"},
					[]string{"4"})}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].input.UserID).Return(tests[0].resultGetUserByID.user, tests[0].resultGetUserByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].input.UserID).Return(tests[1].resultGetUserByID.user, tests[1].resultGetUserByID.err),
		testRepo.mockFriensR.EXPECT().GetFriendsByID(ctx, tests[1].resultGetUserByID.user.FriendsID).Return(tests[1].resultGetFriendsByID.friends, tests[1].resultGetFriendsByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].input.UserID).Return(tests[2].resultGetUserByID.user, tests[2].resultGetUserByID.err),
		testRepo.mockFriensR.EXPECT().GetFriendsByID(ctx, tests[2].resultGetUserByID.user.FriendsID).Return(tests[2].resultGetFriendsByID.friends, tests[2].resultGetFriendsByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.UserService.GetProfile(dbUserImpl, ctx, test.input)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestEditProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditProfileRequest
		userID string
	}

	type EditInfo struct {
		user *core.User
		err  error
	}

	type FriendsByID struct {
		friends []string
		err     error
	}

	type Output struct {
		res *dto.EditProfileResponse
		err error
	}
	tests := []struct {
		name                 string
		input                Input
		resultEditInfo       EditInfo
		resultGetFriendsByID FriendsByID
		output               Output
	}{
		{
			name: "Don't found in BD",
			input: Input{info: &dto.EditProfileRequest{
				NewInfo: dto.EditProfile{
					Name:     common.UserName{First: "John", Last: "Doe"},
					Avatar:   "fmt/img/avatar.jpg",
					Phone:    "+8(800)-555-35-35",
					Location: "Moscow",
					BirthDay: "01.02.2018"}},
				userID: "0"},
			resultEditInfo:       EditInfo{nil, constants.ErrDBNotFound},
			resultGetFriendsByID: FriendsByID{nil, nil},
			output:               Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Found in BD",
			input: Input{info: &dto.EditProfileRequest{
				NewInfo: dto.EditProfile{
					Name:     common.UserName{First: "John", Last: "Doe"},
					Avatar:   "fmt/img/avatar.jpg",
					Phone:    "+8(800)-555-35-35",
					Location: "Moscow",
					BirthDay: "01.02.2018"}},
				userID: "1"},
			resultEditInfo:       EditInfo{&core.User{}, nil},
			resultGetFriendsByID: FriendsByID{nil, nil},
			output:               Output{&dto.EditProfileResponse{UserProfile: dto.UserProfile{}}, nil},
		},
	}
	tes1 := convert.EditProfile2Core(&tests[0].input.info.NewInfo)
	tes2 := convert.EditProfile2Core(&tests[1].input.info.NewInfo)
	gomock.InOrder(
		testRepo.mockUserR.EXPECT().EditInfo(ctx, &tes1, tests[0].input.userID).Return(tests[0].resultEditInfo.user, tests[0].resultEditInfo.err),
		testRepo.mockUserR.EXPECT().EditInfo(ctx, &tes2, tests[1].input.userID).Return(tests[1].resultEditInfo.user, tests[1].resultEditInfo.err),
		testRepo.mockFriensR.EXPECT().GetFriendsByUserID(ctx, tests[1].input.userID).Return(tests[1].resultGetFriendsByID.friends, tests[1].resultGetFriendsByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.UserService.EditProfile(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
