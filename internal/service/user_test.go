package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestGetUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

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
			_, res := UserService.GetUserData(dbUserImpl, ctx, test.input)
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
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type InputGetUserByID struct {
		userID string
	}

	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputGetPostsByUserID struct {
		userID     string
		pageNumber int64
		limit      int64
	}

	type OutputGetPostsByUserID struct {
		post  []core.Post
		pages *common.PageResponse
		err   error
	}

	type InputGetLikeBySubjectID struct {
		postIDs []string
	}

	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type Output struct {
		res *dto.GetUserPostsResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    *dto.GetUserPostsRequest
		inputGetUserByID         InputGetUserByID
		outputGetUserByID        OutputGetUserByID
		inputGetPostsByUserID    InputGetPostsByUserID
		outputGetPostsByUserID   OutputGetPostsByUserID
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		output                   Output
	}{
		{
			name:              "Didn't find user in db",
			input:             &dto.GetUserPostsRequest{UserID: "0", Limit: -1, Page: 1},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{}, err: constants.ErrDBNotFound},
			output:            Output{res: nil, err: constants.ErrDBNotFound},
		},
		{
			name:                  "Success",
			input:                 &dto.GetUserPostsRequest{UserID: "677be1d2", Limit: -1, Page: 1},
			inputGetUserByID:      InputGetUserByID{userID: "677be1d2"},
			outputGetUserByID:     OutputGetUserByID{user: &core.User{ID: "677be1d2", Posts: []string{"123", "234"}}, err: nil},
			inputGetPostsByUserID: InputGetPostsByUserID{userID: "677be1d2", pageNumber: 1, limit: -1},
			outputGetPostsByUserID: OutputGetPostsByUserID{post: []core.Post{{
				ID:          "3",
				AuthorID:    "677be1d2",
				Message:     "Message",
				Attachments: nil,
				CreatedAt:   12341,
				Type:        "User",
			}}, pages: &common.PageResponse{Total: 1, AmountPages: 1}, err: nil},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{postIDs: []string{"3"}},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{like: &core.Like{
				ID:        "3",
				Subject:   "3",
				Amount:    0,
				UserIDs:   nil,
				CreatedAt: 0,
			}, err: nil,
			},
			output: Output{&dto.GetUserPostsResponse{Posts: []dto.GetPosts{{Post: convert.Post2DTOByUser(&core.Post{
				ID:          "3",
				AuthorID:    "677be1d2",
				Message:     "Message",
				Attachments: nil,
				CreatedAt:   12341,
				Type:        "User",
			}, &core.User{
				ID: "677be1d2",
			}), Likes: convert.Like2DTO(&core.Like{
				ID:        "3",
				Subject:   "3",
				Amount:    0,
				UserIDs:   nil,
				CreatedAt: 0,
			}, "2")}}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		//second
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockPostR.EXPECT().GetPostsByUserID(ctx, tests[1].inputGetPostsByUserID.userID, tests[1].inputGetPostsByUserID.pageNumber, tests[1].inputGetPostsByUserID.limit).Return(tests[1].outputGetPostsByUserID.post, tests[1].outputGetPostsByUserID.pages, tests[1].outputGetPostsByUserID.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[1].inputGetLikeBySubjectID.postIDs[0]).Return(tests[1].outputGetLikeBySubjectID.like, tests[1].outputGetLikeBySubjectID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, err := UserService.GetUserPosts(dbUserImpl, ctx, test.input)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}

			if !assert.Equal(t, test.output.err, err) {
				t.Error("got : ", err, " expected :", test.output)
			}
		})
	}
}

func TestGetFeed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		userID string
		info   *dto.GetUserFeedRequest
	}

	type InputGetUserByID struct {
		userID string
	}

	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputGetFeed struct {
		userID     string
		pageNumber int64
		limit      int64
	}

	type OutputGetFeed struct {
		post  []core.Post
		pages *common.PageResponse
		err   error
	}

	type InputAuthorID struct {
		userID string
	}

	type OutputAuthorID struct {
		user *core.User
		err  error
	}

	type InputGetLikeBySubjectID struct {
		postIDs []string
	}

	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type Output struct {
		res *dto.GetUserFeedResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetUserByID         InputGetUserByID
		outputGetUserByID        OutputGetUserByID
		inputGetFeed             InputGetFeed
		outputGetFeed            OutputGetFeed
		inputAuthorID            InputAuthorID
		outputAuthorID           OutputAuthorID
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		output                   Output
	}{
		{
			name:              "Didn't find user in db",
			input:             Input{info: &dto.GetUserFeedRequest{Limit: -1, Page: 1}, userID: "0"},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{}, err: constants.ErrDBNotFound},
			output:            Output{res: nil, err: constants.ErrDBNotFound},
		},
		{
			name:              "Success",
			input:             Input{info: &dto.GetUserFeedRequest{Limit: -1, Page: 1}, userID: "677be1d2"},
			inputGetUserByID:  InputGetUserByID{userID: "677be1d2"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{ID: "677be1d2", Posts: []string{"123", "234"}}, err: nil},
			inputGetFeed:      InputGetFeed{userID: "677be1d2", pageNumber: 1, limit: -1},
			outputGetFeed: OutputGetFeed{post: []core.Post{{
				ID:          "3",
				AuthorID:    "1234",
				Message:     "Message",
				Attachments: nil,
				CreatedAt:   12341,
				Type:        "user",
			}}, pages: &common.PageResponse{Total: 1, AmountPages: 1}, err: nil},
			inputAuthorID:           InputAuthorID{userID: "1234"},
			outputAuthorID:          OutputAuthorID{user: &core.User{ID: "1234", Posts: []string{"3"}}, err: nil},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{postIDs: []string{"3"}},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{like: &core.Like{
				ID:        "3",
				Subject:   "3",
				Amount:    0,
				UserIDs:   nil,
				CreatedAt: 0,
			}, err: nil,
			},
			output: Output{&dto.GetUserFeedResponse{Posts: []dto.GetPosts{{Post: convert.Post2DTOByUser(&core.Post{
				ID:          "3",
				AuthorID:    "1234",
				Message:     "Message",
				Attachments: nil,
				CreatedAt:   12341,
				Type:        "User",
			}, &core.User{
				ID:    "1234",
				Posts: []string{"3"},
			}), Likes: convert.Like2DTO(&core.Like{
				ID:        "3",
				Subject:   "3",
				Amount:    0,
				UserIDs:   nil,
				CreatedAt: 0,
			}, "2")}}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		//second
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockPostR.EXPECT().GetFeed(ctx, tests[1].inputGetFeed.userID, tests[1].inputGetFeed.pageNumber, tests[1].inputGetFeed.limit).Return(tests[1].outputGetFeed.post, tests[1].outputGetFeed.pages, tests[1].outputGetFeed.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[1].inputGetLikeBySubjectID.postIDs[0]).Return(tests[1].outputGetLikeBySubjectID.like, tests[1].outputGetLikeBySubjectID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputAuthorID.userID).Return(tests[1].outputAuthorID.user, tests[1].outputAuthorID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, err := UserService.GetFeed(dbUserImpl, ctx, test.input.userID, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}

			if !assert.Equal(t, test.output.err, err) {
				t.Error("got : ", err, " expected :", test.output)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetProfileRequest
	}

	type InputGetUserByID struct {
		userID string
	}

	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type Output struct {
		res *dto.GetProfileResponse
		err error
	}

	tests := []struct {
		name              string
		input             Input
		inputGetUserByID  InputGetUserByID
		outputGetUserByID OutputGetUserByID
		output            Output
	}{
		{
			name:              "don't found in database",
			input:             Input{info: &dto.GetProfileRequest{UserID: "0"}},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{}, err: mongo.ErrNoDocuments},
			output:            Output{res: nil, err: mongo.ErrNoDocuments},
		},
		{
			name:              "success",
			input:             Input{info: &dto.GetProfileRequest{UserID: "123"}},
			inputGetUserByID:  InputGetUserByID{userID: "123"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{ID: "123", Email: "123@bk", Phone: "891"}, err: nil},
			output:            Output{res: &dto.GetProfileResponse{UserProfile: convert.Profile2DTO(&core.User{ID: "123", Email: "123@bk", Phone: "891"})}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		//second
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, err := UserService.GetProfile(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}

			if !assert.Equal(t, test.output.err, err) {
				t.Error("got : ", err, " expected :", test.output)
			}
		})
	}
}

func TestEditProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditProfileRequest
		userID string
	}

	type EditInfo struct {
		user *core.User
		err  error
	}

	type UpdateUser struct {
		err error
	}

	type Output struct {
		res *dto.EditProfileResponse
		err error
	}
	tests := []struct {
		name             string
		input            Input
		resultEditInfo   EditInfo
		resultUpdateUser UpdateUser
		output           Output
	}{
		{
			name: "Don't found in BD",
			input: Input{info: &dto.EditProfileRequest{
				Name:     common.UserName{First: "John", Last: "Doe"},
				Avatar:   "fmt/img/avatar.jpg",
				Phone:    "+8(800)-555-35-35",
				Location: "Moscow",
				BirthDay: "01.02.2018"}, userID: "0"},
			resultEditInfo: EditInfo{nil, constants.ErrDBNotFound},
			output:         Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Found in BD",
			input: Input{info: &dto.EditProfileRequest{
				Name:     common.UserName{First: "John", Last: "Doe"},
				Avatar:   "fmt/img/avatar.jpg",
				Phone:    "+8(800)-555-35-35",
				Location: "Moscow",
				BirthDay: "01.02.2018"}, userID: "1"},
			resultEditInfo:   EditInfo{&core.User{}, nil},
			resultUpdateUser: UpdateUser{nil},
			output:           Output{&dto.EditProfileResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].input.userID).Return(tests[0].resultEditInfo.user, tests[0].resultEditInfo.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].input.userID).Return(tests[1].resultEditInfo.user, tests[1].resultEditInfo.err),
		testRepo.mockUserR.EXPECT().UpdateUser(ctx, &core.User{
			Image:    "fmt/img/avatar.jpg",
			Name:     common.UserName{First: "John", Last: "Doe"},
			Phone:    "+8(800)-555-35-35",
			Location: "Moscow",
			BirthDay: "01.02.2018"}).Return(tests[1].resultUpdateUser.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := UserService.EditProfile(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestUpdatePhotoUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewUserService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info string
		url  string
	}

	type InputGetUserByID struct {
		userID string
	}

	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputUpdatePhoto struct {
		user *core.User
	}

	type OutputUpdatePhoto struct {
		err error
	}

	type Output struct {
		res *dto.UpdatePhotoResponse
		err error
	}

	tests := []struct {
		name              string
		input             Input
		inputGetUserByID  InputGetUserByID
		outputGetUserByID OutputGetUserByID
		inputUpdatePhoto  InputUpdatePhoto
		outputUpdatePhoto OutputUpdatePhoto
		output            Output
	}{
		{
			name:              "don't found in database",
			input:             Input{info: "0", url: "-"},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{}, err: mongo.ErrNoDocuments},
			output:            Output{res: nil, err: mongo.ErrNoDocuments},
		},
		{
			name:              "success",
			input:             Input{info: "123", url: "http/"},
			inputGetUserByID:  InputGetUserByID{userID: "123"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{ID: "123", Email: "123@bk", Phone: "891"}, err: nil},
			inputUpdatePhoto:  InputUpdatePhoto{user: &core.User{ID: "123", Email: "123@bk", Phone: "891", Image: "http/"}},
			outputUpdatePhoto: OutputUpdatePhoto{err: nil},
			output:            Output{res: &dto.UpdatePhotoResponse{URL: "http/"}, err: nil},
		},
	}

	gomock.InOrder(
		// first
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		//second
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockUserR.EXPECT().UpdateUser(ctx, tests[1].inputUpdatePhoto.user).Return(tests[1].outputUpdatePhoto.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, err := UserService.UpdatePhoto(dbUserImpl, ctx, test.input.url, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}

			if !assert.Equal(t, test.output.err, err) {
				t.Error("got : ", err, " expected :", test.output)
			}
		})
	}
}
