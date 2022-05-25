package service

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.CreateCommunityRequest
		userID string
	}
	type InputCreateCommunity struct {
		community *core.Community
	}
	type OutputCreateCommunity struct {
		community *core.Community
		err       error
	}

	type InputUserAddCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserAddCommunity struct {
		err error
	}

	type Output struct {
		res *dto.CreateCommunityResponse
		err error
	}
	var errInsert = errors.Errorf("Can't insert")
	var errAddCommunity = errors.Errorf("Can't push")
	tests := []struct {
		name                   string
		input                  Input
		inputCreateCommunity   InputCreateCommunity
		outputCreateCommunity  OutputCreateCommunity
		inputUserAddCommunity  InputUserAddCommunity
		outputUserAddCommunity OutputUserAddCommunity
		output                 Output
	}{
		{
			name: "Can't insert in community",
			input: Input{
				info: &dto.CreateCommunityRequest{
					Name:   "New community",
					Image:  "image.jpg",
					Info:   "New community info",
					Admins: nil,
				},
				userID: "0",
			},
			inputCreateCommunity: InputCreateCommunity{community: &core.Community{
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"0"},
				AdminIDs:    []string{"0"},
			}},
			outputCreateCommunity: OutputCreateCommunity{community: nil, err: errInsert},
			output:                Output{nil, errInsert},
		},
		{
			name: "Can't add community",
			input: Input{
				info: &dto.CreateCommunityRequest{
					Name:   "New community",
					Image:  "image.jpg",
					Info:   "New community info",
					Admins: nil,
				},
				userID: "1",
			},
			inputCreateCommunity: InputCreateCommunity{community: &core.Community{
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"1"},
				AdminIDs:    []string{"1"},
			}},
			outputCreateCommunity: OutputCreateCommunity{community: &core.Community{
				ID:          "1",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"1"},
				AdminIDs:    []string{"1"},
				PostIDs:     nil,
				CreatedAt:   123451,
			}, err: nil},
			inputUserAddCommunity: InputUserAddCommunity{
				userID:      "1",
				communityID: "1",
			},
			outputUserAddCommunity: OutputUserAddCommunity{err: errAddCommunity},
			output:                 Output{nil, errAddCommunity},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.CreateCommunityRequest{
					Name:   "New community",
					Image:  "image.jpg",
					Info:   "New community info",
					Admins: nil,
				},
				userID: "2",
			},
			inputCreateCommunity: InputCreateCommunity{community: &core.Community{
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"2"},
				AdminIDs:    []string{"2"},
			}},
			outputCreateCommunity: OutputCreateCommunity{community: &core.Community{
				ID:          "2",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"2"},
				AdminIDs:    []string{"2"},
				PostIDs:     nil,
				CreatedAt:   123451,
			}, err: nil},
			inputUserAddCommunity: InputUserAddCommunity{
				userID:      "2",
				communityID: "2",
			},
			outputUserAddCommunity: OutputUserAddCommunity{err: nil},
			output:                 Output{&dto.CreateCommunityResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().CreateCommunity(ctx, tests[0].inputCreateCommunity.community).Return(tests[0].outputCreateCommunity.community, tests[0].outputCreateCommunity.err),

		testRepo.mockCommunityR.EXPECT().CreateCommunity(ctx, tests[1].inputCreateCommunity.community).Return(tests[1].outputCreateCommunity.community, tests[1].outputCreateCommunity.err),
		testRepo.mockUserR.EXPECT().UserAddCommunity(ctx, tests[1].inputUserAddCommunity.userID, tests[1].inputUserAddCommunity.communityID).Return(tests[1].outputUserAddCommunity.err),

		testRepo.mockCommunityR.EXPECT().CreateCommunity(ctx, tests[2].inputCreateCommunity.community).Return(tests[2].outputCreateCommunity.community, tests[2].outputCreateCommunity.err),
		testRepo.mockUserR.EXPECT().UserAddCommunity(ctx, tests[2].inputUserAddCommunity.userID, tests[2].inputUserAddCommunity.communityID).Return(tests[2].outputUserAddCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.CreateCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetCommunityRequest
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}

	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type Output struct {
		res *dto.GetCommunityResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetCommunityByID  InputGetCommunityByID
		outputGetCommunityByID OutputGetCommunityByID
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		output                 Output
	}{
		{
			name: "Didn't find community in db",
			input: Input{
				info: &dto.GetCommunityRequest{CommunityID: "0"},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "0"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Didn't find post admin in db",
			input: Input{
				info: &dto.GetCommunityRequest{CommunityID: "1"},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "1",
					Name:        "New community",
					Image:       "image.jpg",
					Info:        "New community info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
					PostIDs:     nil,
					CreatedAt:   123451,
				},
				err: nil,
			},
			inputGetUserByID: InputGetUserByID{
				userID: "1",
			},
			outputGetUserByID: OutputGetUserByID{user: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.GetCommunityRequest{CommunityID: "2"},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "New community",
					Image:       "image.jpg",
					Info:        "New community info",
					FollowerIDs: []string{"2"},
					AdminIDs:    []string{"2"},
					PostIDs:     nil,
					CreatedAt:   123451,
				},
				err: nil,
			},
			inputGetUserByID: InputGetUserByID{
				userID: "2",
			},
			outputGetUserByID: OutputGetUserByID{user: &core.User{
				ID: "2",
				Name: common.UserName{
					First: "CJ",
					Last:  "CJ",
				},
				Image: "image.jpg",
				Email: "cj@cj.com",
			}, err: nil},
			output: Output{&dto.GetCommunityResponse{Community: convert.Community2DTOprofile(&core.Community{
				ID:          "2",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"2"},
				AdminIDs:    []string{"2"},
				PostIDs:     nil,
				CreatedAt:   123451,
			}, []dto.User{{
				ID: "2",
				Name: common.UserName{
					First: "CJ",
					Last:  "CJ",
				},
				Image: "image.jpg",
				Email: "cj@cj.com"}})}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[0].inputGetCommunityByID.communityID).Return(tests[0].outputGetCommunityByID.community, tests[0].outputGetCommunityByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetCommunity(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetUserManageCommunities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetUserManageCommunitiesRequest
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}

	type Output struct {
		res *dto.GetUserManageCommunitiesResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		inputGetCommunityByID  InputGetCommunityByID
		outputGetCommunityByID OutputGetCommunityByID
		output                 Output
	}{
		{
			name: "Didn't find post in db",
			input: Input{
				info: &dto.GetUserManageCommunitiesRequest{
					UserID: "0",
					Limit:  -1,
					Page:   1,
				},
			},
			inputGetUserByID: InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{
				user: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Didn't find community in db",
			input: Input{
				info: &dto.GetUserManageCommunitiesRequest{
					UserID: "1",
					Limit:  -1,
					Page:   1,
				},
			},
			inputGetUserByID: InputGetUserByID{userID: "1"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{
				ID: "1",
				Name: common.UserName{
					First: "CJ",
					Last:  "CJ",
				},
				Image:        "image.jpg",
				Email:        "cj@cj.com",
				CommunityIDs: []string{"1"},
			}, err: nil},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "success",
			input: Input{
				info: &dto.GetUserManageCommunitiesRequest{
					UserID: "2",
					Limit:  -1,
					Page:   1,
				},
			},
			inputGetUserByID: InputGetUserByID{userID: "2"},
			outputGetUserByID: OutputGetUserByID{user: &core.User{
				ID: "2",
				Name: common.UserName{
					First: "CJ",
					Last:  "CJ",
				},
				Image:        "image.jpg",
				Email:        "cj@cj.com",
				CommunityIDs: []string{"2"},
			}, err: nil},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "New community",
					Image:       "image.jpg",
					Info:        "New community info",
					FollowerIDs: []string{"2"},
					AdminIDs:    []string{"2"},
					PostIDs:     nil,
					CreatedAt:   123451,
				},
				err: nil,
			},
			output: Output{&dto.GetUserManageCommunitiesResponse{Communities: []dto.Community{convert.Community2DTO(&core.Community{
				ID:          "2",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"2"},
				AdminIDs:    []string{"2"},
				PostIDs:     nil,
				CreatedAt:   123451,
			})}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetUserManageCommunities(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetCommunities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetCommunitiesRequest
	}
	type InputGetAllCommunities struct {
		limit      int64
		pageNumber int64
	}
	type OutputGetAllCommunities struct {
		communities []core.Community
		page        *common.PageResponse
		err         error
	}
	type Output struct {
		res *dto.GetCommunitiesResponse
		err error
	}

	tests := []struct {
		name                    string
		input                   Input
		inputGetAllCommunities  InputGetAllCommunities
		outputGetAllCommunities OutputGetAllCommunities
		output                  Output
	}{
		{
			name: "GetAllCommunities error",
			input: Input{
				info: &dto.GetCommunitiesRequest{
					Limit: -1,
					Page:  1,
				}},
			inputGetAllCommunities: InputGetAllCommunities{
				limit:      -1,
				pageNumber: 1,
			},
			outputGetAllCommunities: OutputGetAllCommunities{
				communities: nil,
				page:        nil,
				err:         constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.GetCommunitiesRequest{
					Limit: -1,
					Page:  1,
				}},
			inputGetAllCommunities: InputGetAllCommunities{
				limit:      -1,
				pageNumber: 1,
			},
			outputGetAllCommunities: OutputGetAllCommunities{
				communities: []core.Community{{
					ID:          "1",
					Name:        "New community",
					Image:       "image.jpg",
					Info:        "New community info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
					PostIDs:     nil,
					CreatedAt:   123451,
				}},
				page: &common.PageResponse{
					Total:       1,
					AmountPages: 1,
				},
				err: nil,
			},
			output: Output{&dto.GetCommunitiesResponse{Communities: []dto.Community{convert.Community2DTO(&core.Community{
				ID:          "1",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"1"},
				AdminIDs:    []string{"1"},
				PostIDs:     nil,
				CreatedAt:   123451,
			})}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().GetAllCommunities(ctx, tests[0].inputGetAllCommunities.limit, tests[0].inputGetAllCommunities.pageNumber).Return(tests[0].outputGetAllCommunities.communities, tests[0].outputGetAllCommunities.page, tests[0].outputGetAllCommunities.err),

		testRepo.mockCommunityR.EXPECT().GetAllCommunities(ctx, tests[1].inputGetAllCommunities.limit, tests[1].inputGetAllCommunities.pageNumber).Return(tests[1].outputGetAllCommunities.communities, tests[1].outputGetAllCommunities.page, tests[1].outputGetAllCommunities.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetCommunities(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestJoinCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.JoinCommunityRequest
		userID string
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}
	type InputAddFollower struct {
		communityID string
		userID      string
	}
	type OutputAddFollower struct {
		err error
	}
	type InputUserAddCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserAddCommunity struct {
		err error
	}
	type Output struct {
		res *dto.JoinCommunityResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		inputAddFollower       InputAddFollower
		outputAddFollower      OutputAddFollower
		inputUserAddCommunity  InputUserAddCommunity
		outputUserAddCommunity OutputUserAddCommunity
		output                 Output
	}{
		{
			name: "Didn't find in db",
			input: Input{
				info:   &dto.JoinCommunityRequest{CommunityID: "0"},
				userID: "0",
			},
			inputGetUserByID: InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{
				user: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Error already follower",
			input: Input{
				info:   &dto.JoinCommunityRequest{CommunityID: "1"},
				userID: "1",
			},
			inputGetUserByID: InputGetUserByID{userID: "1"},
			outputGetUserByID: OutputGetUserByID{
				user: &core.User{
					ID:           "1",
					Name:         common.UserName{},
					Image:        "image.jpg",
					CommunityIDs: []string{"1"},
				},
				err: nil,
			},
			output: Output{nil, constants.ErrAlreadyFollower},
		},
		{
			name: "Success",
			input: Input{
				info:   &dto.JoinCommunityRequest{CommunityID: "2"},
				userID: "2",
			},
			inputGetUserByID: InputGetUserByID{userID: "2"},
			outputGetUserByID: OutputGetUserByID{
				user: &core.User{
					ID:           "2",
					Name:         common.UserName{},
					Image:        "image.jpg",
					CommunityIDs: []string{"3"},
				},
				err: nil,
			},
			inputAddFollower: InputAddFollower{
				communityID: "2",
				userID:      "2",
			},
			outputUserAddCommunity: OutputUserAddCommunity{err: nil},
			inputUserAddCommunity: InputUserAddCommunity{
				userID:      "2",
				communityID: "2",
			},
			outputAddFollower: OutputAddFollower{err: nil},
			output:            Output{&dto.JoinCommunityResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
		testRepo.mockCommunityR.EXPECT().AddFollower(ctx, tests[2].inputAddFollower.communityID, tests[2].inputAddFollower.userID).Return(tests[2].outputAddFollower.err),
		testRepo.mockUserR.EXPECT().UserAddCommunity(ctx, tests[2].inputUserAddCommunity.communityID, tests[2].inputUserAddCommunity.userID).Return(tests[2].outputUserAddCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.JoinCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestLeaveCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.LeaveCommunityRequest
		userID string
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}
	type InputDeleteAdmin struct {
		communityID string
		userID      string
	}
	type OutputDeleteAdmin struct {
		err error
	}
	type InputDeleteFollower struct {
		communityID string
		userID      string
	}
	type OutputDeleteFollower struct {
		err error
	}
	type InputUserDeleteCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserDeleteCommunity struct {
		err error
	}
	type Output struct {
		res *dto.LeaveCommunityResponse
		err error
	}

	var errDeleteAdmin = errors.Errorf("Can't delete admin")
	var errDeleteFollower = errors.Errorf("Can't delete follower")
	var errDeleteUserCommunity = errors.Errorf("Can't delete post community")
	tests := []struct {
		name                      string
		input                     Input
		inputGetCommunityByID     InputGetCommunityByID
		outputGetCommunityByID    OutputGetCommunityByID
		inputDeleteAdmin          InputDeleteAdmin
		outputDeleteAdmin         OutputDeleteAdmin
		inputDeleteFollower       InputDeleteFollower
		outputDeleteFollower      OutputDeleteFollower
		inputUserDeleteCommunity  InputUserDeleteCommunity
		outputUserDeleteCommunity OutputUserDeleteCommunity
		output                    Output
	}{
		{
			name: "Didn't find community in db",
			input: Input{
				info:   &dto.LeaveCommunityRequest{CommunityID: "0"},
				userID: "0",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "0"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Didn't can't delete admin of community in db",
			input: Input{
				info:   &dto.LeaveCommunityRequest{CommunityID: "1"},
				userID: "1",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "1",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: nil,
					AdminIDs:    []string{"1", "2"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputDeleteAdmin: InputDeleteAdmin{
				communityID: "1",
				userID:      "1",
			},
			outputDeleteAdmin: OutputDeleteAdmin{err: errDeleteAdmin},
			output:            Output{nil, errDeleteAdmin},
		},
		{
			name: "Didn't can't delete follower of community in db",
			input: Input{
				info:   &dto.LeaveCommunityRequest{CommunityID: "2"},
				userID: "2",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: nil,
					AdminIDs:    []string{"1", "2"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputDeleteAdmin: InputDeleteAdmin{
				communityID: "2",
				userID:      "2",
			},
			outputDeleteAdmin: OutputDeleteAdmin{err: nil},
			inputDeleteFollower: InputDeleteFollower{
				communityID: "2",
				userID:      "2",
			},
			outputDeleteFollower: OutputDeleteFollower{err: errDeleteFollower},
			output:               Output{nil, errDeleteFollower},
		},
		{
			name: "UserDeleteCommunity error",
			input: Input{
				info:   &dto.LeaveCommunityRequest{CommunityID: "3"},
				userID: "3",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "3"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: nil,
					AdminIDs:    []string{"3", "2"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputDeleteAdmin: InputDeleteAdmin{
				communityID: "3",
				userID:      "3",
			},
			outputDeleteAdmin: OutputDeleteAdmin{err: nil},
			inputDeleteFollower: InputDeleteFollower{
				communityID: "3",
				userID:      "3",
			},
			outputDeleteFollower: OutputDeleteFollower{err: nil},
			inputUserDeleteCommunity: InputUserDeleteCommunity{
				userID:      "3",
				communityID: "3",
			},
			outputUserDeleteCommunity: OutputUserDeleteCommunity{err: errDeleteUserCommunity},
			output:                    Output{nil, errDeleteUserCommunity},
		},
		{
			name: "Success",
			input: Input{
				info:   &dto.LeaveCommunityRequest{CommunityID: "4"},
				userID: "4",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "4"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "4",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: nil,
					AdminIDs:    []string{"3", "4"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputDeleteAdmin: InputDeleteAdmin{
				communityID: "4",
				userID:      "4",
			},
			outputDeleteAdmin: OutputDeleteAdmin{err: nil},
			inputDeleteFollower: InputDeleteFollower{
				communityID: "4",
				userID:      "4",
			},
			outputDeleteFollower: OutputDeleteFollower{err: nil},
			inputUserDeleteCommunity: InputUserDeleteCommunity{
				userID:      "4",
				communityID: "4",
			},
			outputUserDeleteCommunity: OutputUserDeleteCommunity{err: nil},
			output:                    Output{&dto.LeaveCommunityResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[0].inputGetCommunityByID.communityID).Return(tests[0].outputGetCommunityByID.community, tests[0].outputGetCommunityByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().DeleteAdmin(ctx, tests[1].inputDeleteAdmin.communityID, tests[1].inputDeleteAdmin.userID).Return(tests[1].outputDeleteAdmin.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().DeleteAdmin(ctx, tests[2].inputDeleteAdmin.communityID, tests[2].inputDeleteAdmin.userID).Return(tests[2].outputDeleteAdmin.err),
		testRepo.mockCommunityR.EXPECT().DeleteFollower(ctx, tests[2].inputDeleteFollower.communityID, tests[2].inputDeleteFollower.userID).Return(tests[2].outputDeleteFollower.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[3].inputGetCommunityByID.communityID).Return(tests[3].outputGetCommunityByID.community, tests[3].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().DeleteAdmin(ctx, tests[3].inputDeleteAdmin.communityID, tests[3].inputDeleteAdmin.userID).Return(tests[3].outputDeleteAdmin.err),
		testRepo.mockCommunityR.EXPECT().DeleteFollower(ctx, tests[3].inputDeleteFollower.communityID, tests[3].inputDeleteFollower.userID).Return(tests[3].outputDeleteFollower.err),
		testRepo.mockUserR.EXPECT().UserDeleteCommunity(ctx, tests[3].inputUserDeleteCommunity.userID, tests[3].inputUserDeleteCommunity.communityID).Return(tests[3].outputUserDeleteCommunity.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[4].inputGetCommunityByID.communityID).Return(tests[4].outputGetCommunityByID.community, tests[4].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().DeleteAdmin(ctx, tests[4].inputDeleteAdmin.communityID, tests[4].inputDeleteAdmin.userID).Return(tests[4].outputDeleteAdmin.err),
		testRepo.mockCommunityR.EXPECT().DeleteFollower(ctx, tests[4].inputDeleteFollower.communityID, tests[4].inputDeleteFollower.userID).Return(tests[4].outputDeleteFollower.err),
		testRepo.mockUserR.EXPECT().UserDeleteCommunity(ctx, tests[4].inputUserDeleteCommunity.userID, tests[4].inputUserDeleteCommunity.communityID).Return(tests[4].outputUserDeleteCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.LeaveCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetFollowers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetFollowersRequest
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type Output struct {
		res *dto.GetFollowersResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetCommunityByID  InputGetCommunityByID
		outputGetCommunityByID OutputGetCommunityByID
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		output                 Output
	}{
		{
			name: "Didn't find community in db",
			input: Input{
				info: &dto.GetFollowersRequest{CommunityID: "0"},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "0"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Didn't can't delete admin of community in db",
			input: Input{
				info: &dto.GetFollowersRequest{CommunityID: "1", Limit: -1, Page: 1},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "1",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputGetUserByID: InputGetUserByID{
				userID: "1",
			},
			outputGetUserByID: OutputGetUserByID{user: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.GetFollowersRequest{CommunityID: "2", Limit: -1, Page: 1},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: []string{"2"},
					AdminIDs:    []string{"2"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputGetUserByID: InputGetUserByID{
				userID: "2",
			},
			outputGetUserByID: OutputGetUserByID{user: &core.User{
				ID:    "2",
				Name:  common.UserName{},
				Image: "image.jpg",
				Email: "123",
			}, err: nil},
			output: Output{&dto.GetFollowersResponse{Amount: 1, Followers: []dto.User{{ID: "2",
				Name:  common.UserName{},
				Image: "image.jpg",
				Email: "123"}}, AmountPages: 1, Total: 1}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[0].inputGetCommunityByID.communityID).Return(tests[0].outputGetCommunityByID.community, tests[0].outputGetCommunityByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetFollowers(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetUserCommunities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetUserCommunitiesRequest
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}
	type InputCommunityByID struct {
		community string
	}
	type OutputCommunityByID struct {
		community *core.Community
		err       error
	}

	type Output struct {
		res *dto.GetUserCommunitiesResponse
		err error
	}

	tests := []struct {
		name                string
		input               Input
		inputGetUserByID    InputGetUserByID
		outputGetUserByID   OutputGetUserByID
		inputCommunityByID  InputCommunityByID
		outputCommunityByID OutputCommunityByID
		output              Output
	}{
		{
			name: "Didn't find post in db",
			input: Input{
				info: &dto.GetUserCommunitiesRequest{UserID: "0"},
			},
			inputGetUserByID: InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{
				user: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "GetCommunityByID error",
			input: Input{
				info: &dto.GetUserCommunitiesRequest{UserID: "1", Limit: -1, Page: 1},
			},
			inputGetUserByID: InputGetUserByID{userID: "1"},
			outputGetUserByID: OutputGetUserByID{
				user: &core.User{
					ID:           "1",
					CommunityIDs: []string{"1"},
				},
				err: nil,
			},
			inputCommunityByID: InputCommunityByID{
				community: "1",
			},
			outputCommunityByID: OutputCommunityByID{community: nil, err: constants.ErrDBNotFound},
			output:              Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.GetUserCommunitiesRequest{UserID: "2", Limit: -1, Page: 1},
			},
			inputGetUserByID: InputGetUserByID{userID: "2"},
			outputGetUserByID: OutputGetUserByID{
				user: &core.User{
					ID:           "2",
					CommunityIDs: []string{"2"},
				},
				err: nil,
			},
			inputCommunityByID: InputCommunityByID{
				community: "2",
			},
			outputCommunityByID: OutputCommunityByID{community: &core.Community{
				ID:    "2",
				Name:  "Community",
				Image: "image.jpg",
			}, err: nil},
			output: Output{&dto.GetUserCommunitiesResponse{Communities: []dto.Community{{
				ID:    "2",
				Name:  "Community",
				Image: "image.jpg"}}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputCommunityByID.community).Return(tests[1].outputCommunityByID.community, tests[1].outputCommunityByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputCommunityByID.community).Return(tests[2].outputCommunityByID.community, tests[2].outputCommunityByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetUserCommunities(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetMutualFriends(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.GetMutualFriendsRequest
		userID string
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type Output struct {
		res *dto.GetMutualFriendsResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		inputGetCommunityByID  InputGetCommunityByID
		outputGetCommunityByID OutputGetCommunityByID
		output                 Output
	}{
		{
			name: "Didn't find community in db",
			input: Input{
				info:   &dto.GetMutualFriendsRequest{CommunityID: "0", Limit: -1, Page: 1},
				userID: "0",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "0"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Didn't find post in db",
			input: Input{
				info:   &dto.GetMutualFriendsRequest{CommunityID: "1", Limit: -1, Page: 1},
				userID: "1",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "1",
					Name:        "Community",
					Image:       "Image",
					Info:        "Info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
				},
				err: nil,
			},
			inputGetUserByID: InputGetUserByID{userID: "1"},
			outputGetUserByID: OutputGetUserByID{
				user: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[0].inputGetCommunityByID.communityID).Return(tests[0].outputGetCommunityByID.community, tests[0].outputGetCommunityByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetMutualFriends(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestSearchCommunities(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.SearchCommunitiesRequest
	}
	type InputSearchCommunities struct {
		selector   string
		limit      int64
		pageNumber int64
	}
	type OutputSearchCommunities struct {
		communities []core.Community
		page        *common.PageResponse
		err         error
	}
	type Output struct {
		res *dto.SearchCommunitiesResponse
		err error
	}

	tests := []struct {
		name                    string
		input                   Input
		inputSearchCommunities  InputSearchCommunities
		outputSearchCommunities OutputSearchCommunities
		output                  Output
	}{
		{
			name: "SearchCommunities error",
			input: Input{
				info: &dto.SearchCommunitiesRequest{
					Selector: "My",
					Limit:    -1,
					Page:     1,
				}},
			inputSearchCommunities: InputSearchCommunities{
				selector:   "My",
				limit:      -1,
				pageNumber: 1,
			},
			outputSearchCommunities: OutputSearchCommunities{
				communities: nil,
				page:        nil,
				err:         constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info: &dto.SearchCommunitiesRequest{
					Selector: "Community",
					Limit:    -1,
					Page:     1,
				}},
			inputSearchCommunities: InputSearchCommunities{
				selector:   "Community",
				limit:      -1,
				pageNumber: 1,
			},
			outputSearchCommunities: OutputSearchCommunities{
				communities: []core.Community{{
					ID:          "1",
					Name:        "New community",
					Image:       "image.jpg",
					Info:        "New community info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
					PostIDs:     nil,
					CreatedAt:   123451,
				}},
				page: &common.PageResponse{
					Total:       1,
					AmountPages: 1,
				},
				err: nil,
			},
			output: Output{&dto.SearchCommunitiesResponse{Communities: []dto.Community{convert.Community2DTO(&core.Community{
				ID:          "1",
				Name:        "New community",
				Image:       "image.jpg",
				Info:        "New community info",
				FollowerIDs: []string{"1"},
				AdminIDs:    []string{"1"},
				PostIDs:     nil,
				CreatedAt:   123451,
			})}, Total: 1, AmountPages: 1}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommunityR.EXPECT().SearchCommunities(ctx, tests[0].inputSearchCommunities.selector, tests[0].inputSearchCommunities.limit, tests[0].inputSearchCommunities.pageNumber).Return(tests[0].outputSearchCommunities.communities, tests[0].outputSearchCommunities.page, tests[0].outputSearchCommunities.err),
		testRepo.mockCommunityR.EXPECT().SearchCommunities(ctx, tests[1].inputSearchCommunities.selector, tests[1].inputSearchCommunities.limit, tests[1].inputSearchCommunities.pageNumber).Return(tests[1].outputSearchCommunities.communities, tests[1].outputSearchCommunities.page, tests[1].outputSearchCommunities.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.SearchCommunities(dbCommunityImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestUpdatePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.UpdatePhotoCommunityRequest
		url    string
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}
	type InputEditCommunity struct {
		community *core.Community
	}
	type OutputEditCommunity struct {
		err error
	}
	type Output struct {
		res *dto.UpdatePhotoCommunityResponse
		err error
	}
	var errEdit = errors.Errorf("Can't edit")
	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		inputGetCommunityByID    InputGetCommunityByID
		outputGetCommunityByID   OutputGetCommunityByID
		inputEditCommunity       InputEditCommunity
		outputEditCommunity      OutputEditCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info:   &dto.UpdatePhotoCommunityRequest{CommunityID: "0"},
				url:    "12143245344.jpg",
				userID: "0"},
			inputUserCheckCommunity: InputUserCheckCommunity{
				userID:      "0",
				communityID: "0",
			},
			outputUserCheckCommunity: OutputUserCheckCommunity{err: constants.ErrDBNotFound},
			output:                   Output{nil, fmt.Errorf("UserCheckCommunity: %w", constants.ErrDBNotFound)},
		},
		{
			name: "GetCommunityByID error",
			input: Input{
				info:   &dto.UpdatePhotoCommunityRequest{CommunityID: "1"},
				url:    "12143245344.jpg",
				userID: "1"},
			inputUserCheckCommunity: InputUserCheckCommunity{
				userID:      "1",
				communityID: "1",
			},
			outputUserCheckCommunity: OutputUserCheckCommunity{err: nil},
			inputGetCommunityByID:    InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, fmt.Errorf("GetCommunityByID: %w", constants.ErrDBNotFound)},
		},
		{
			name: "ErrAuthorIDMismatch error",
			input: Input{
				info:   &dto.UpdatePhotoCommunityRequest{CommunityID: "2"},
				url:    "12143245344.jpg",
				userID: "2"},
			inputUserCheckCommunity: InputUserCheckCommunity{
				userID:      "2",
				communityID: "2",
			},
			outputUserCheckCommunity: OutputUserCheckCommunity{err: nil},
			inputGetCommunityByID:    InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "Community",
					Image:       "img",
					Info:        "Info",
					FollowerIDs: []string{"3"},
					AdminIDs:    []string{"3"},
					CreatedAt:   0,
				},
				err: nil,
			},
			output: Output{nil, constants.ErrAuthorIDMismatch},
		},
		{
			name: "EditCommunity error",
			input: Input{
				info:   &dto.UpdatePhotoCommunityRequest{CommunityID: "3"},
				url:    "12143245344.jpg",
				userID: "3"},
			inputUserCheckCommunity: InputUserCheckCommunity{
				userID:      "3",
				communityID: "3",
			},
			outputUserCheckCommunity: OutputUserCheckCommunity{err: nil},
			inputGetCommunityByID:    InputGetCommunityByID{communityID: "3"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "3",
					Name:        "Community",
					Image:       "img",
					Info:        "Info",
					FollowerIDs: []string{"3"},
					AdminIDs:    []string{"3"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputEditCommunity: InputEditCommunity{community: &core.Community{
				ID:          "3",
				Name:        "Community",
				Image:       "12143245344.jpg",
				Info:        "Info",
				FollowerIDs: []string{"3"},
				AdminIDs:    []string{"3"},
				CreatedAt:   0,
			}},
			outputEditCommunity: OutputEditCommunity{err: errEdit},
			output:              Output{nil, errEdit},
		},
		{
			name: "Success",
			input: Input{
				info:   &dto.UpdatePhotoCommunityRequest{CommunityID: "4"},
				url:    "12143245344.jpg",
				userID: "4"},
			inputUserCheckCommunity: InputUserCheckCommunity{
				userID:      "4",
				communityID: "4",
			},
			outputUserCheckCommunity: OutputUserCheckCommunity{err: nil},
			inputGetCommunityByID:    InputGetCommunityByID{communityID: "4"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "4",
					Name:        "Community",
					Image:       "img",
					Info:        "Info",
					FollowerIDs: []string{"4"},
					AdminIDs:    []string{"4"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputEditCommunity: InputEditCommunity{community: &core.Community{
				ID:          "4",
				Name:        "Community",
				Image:       "12143245344.jpg",
				Info:        "Info",
				FollowerIDs: []string{"4"},
				AdminIDs:    []string{"4"},
				CreatedAt:   0,
			}},
			outputEditCommunity: OutputEditCommunity{err: nil},
			output:              Output{&dto.UpdatePhotoCommunityResponse{URL: "12143245344.jpg"}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),

		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[1].inputUserCheckCommunity.userID, tests[1].inputUserCheckCommunity.communityID).Return(tests[1].outputUserCheckCommunity.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),

		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[2].inputUserCheckCommunity.userID, tests[2].inputUserCheckCommunity.communityID).Return(tests[2].outputUserCheckCommunity.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),

		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[3].inputUserCheckCommunity.userID, tests[3].inputUserCheckCommunity.communityID).Return(tests[3].outputUserCheckCommunity.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[3].inputGetCommunityByID.communityID).Return(tests[3].outputGetCommunityByID.community, tests[3].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().EditCommunity(ctx, tests[3].inputEditCommunity.community).Return(tests[3].outputEditCommunity.err),

		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[4].inputUserCheckCommunity.userID, tests[4].inputUserCheckCommunity.communityID).Return(tests[4].outputUserCheckCommunity.err),
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[4].inputGetCommunityByID.communityID).Return(tests[4].outputGetCommunityByID.community, tests[4].outputGetCommunityByID.err),
		testRepo.mockCommunityR.EXPECT().EditCommunity(ctx, tests[4].inputEditCommunity.community).Return(tests[4].outputEditCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.UpdatePhoto(dbCommunityImpl, ctx, test.input.info, test.input.url, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetCommunityPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.GetCommunityPostsRequest
		userID string
	}
	type InputGetCommunityByID struct {
		communityID string
	}
	type OutputGetCommunityByID struct {
		community *core.Community
		err       error
	}
	type InputGetPostByID struct {
		postID string
	}
	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}
	type InputGetLikeBySubjectID struct {
		postID string
	}
	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}
	type Output struct {
		res *dto.GetCommunityPostsResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetCommunityByID    InputGetCommunityByID
		outputGetCommunityByID   OutputGetCommunityByID
		inputGetPostByID         InputGetPostByID
		outputGetPostByID        OutputGetPostByID
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		output                   Output
	}{
		{
			name: "Didn't find community in db",
			input: Input{
				info: &dto.GetCommunityPostsRequest{CommunityID: "0", Limit: -1, Page: 1},
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "0"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: nil,
				err:       constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "GetPostByID error",
			input: Input{
				info:   &dto.GetCommunityPostsRequest{CommunityID: "1", Limit: -1, Page: 1},
				userID: "0",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "1"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "1",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: []string{"1"},
					AdminIDs:    []string{"1"},
					PostIDs:     []string{"1"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputGetPostByID: InputGetPostByID{
				postID: "1",
			},
			outputGetPostByID: OutputGetPostByID{post: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "GetLikeBySubjectID error",
			input: Input{
				info:   &dto.GetCommunityPostsRequest{CommunityID: "2", Limit: -1, Page: 1},
				userID: "1",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "2"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "2",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: []string{"2"},
					AdminIDs:    []string{"2"},
					PostIDs:     []string{"2"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputGetPostByID: InputGetPostByID{
				postID: "2",
			},
			outputGetPostByID: OutputGetPostByID{post: &core.Post{
				ID:        "2",
				AuthorID:  "2",
				Message:   "Message",
				Images:    nil,
				CreatedAt: 12341,
				Type:      "Community",
			}, err: nil},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{postID: "2"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{
				info:   &dto.GetCommunityPostsRequest{CommunityID: "3", Limit: -1, Page: 1},
				userID: "2",
			},
			inputGetCommunityByID: InputGetCommunityByID{communityID: "3"},
			outputGetCommunityByID: OutputGetCommunityByID{
				community: &core.Community{
					ID:          "3",
					Name:        "My community",
					Image:       "image.jpg",
					Info:        "Info",
					FollowerIDs: []string{"3"},
					AdminIDs:    []string{"3"},
					PostIDs:     []string{"3"},
					CreatedAt:   0,
				},
				err: nil,
			},
			inputGetPostByID: InputGetPostByID{
				postID: "3",
			},
			outputGetPostByID: OutputGetPostByID{post: &core.Post{
				ID:        "3",
				AuthorID:  "3",
				Message:   "Message",
				Images:    nil,
				CreatedAt: 12341,
				Type:      "Community",
			}, err: nil},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{postID: "3"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{
					ID:        "3",
					Subject:   "3",
					Amount:    0,
					UserIDs:   nil,
					CreatedAt: 0,
				},
				err: nil,
			},
			output: Output{&dto.GetCommunityPostsResponse{Posts: []dto.GetPosts{{Post: convert.Post2DTOByCommunity(&core.Post{
				ID:        "3",
				AuthorID:  "3",
				Message:   "Message",
				Images:    nil,
				CreatedAt: 12341,
				Type:      "Community",
			}, &core.Community{
				ID:          "3",
				Name:        "My community",
				Image:       "image.jpg",
				Info:        "Info",
				FollowerIDs: []string{"3"},
				AdminIDs:    []string{"3"},
				PostIDs:     []string{"3"},
				CreatedAt:   0,
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
		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[0].inputGetCommunityByID.communityID).Return(tests[0].outputGetCommunityByID.community, tests[0].outputGetCommunityByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[1].inputGetCommunityByID.communityID).Return(tests[1].outputGetCommunityByID.community, tests[1].outputGetCommunityByID.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[1].inputGetPostByID.postID).Return(tests[1].outputGetPostByID.post, tests[1].outputGetPostByID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[2].inputGetCommunityByID.communityID).Return(tests[2].outputGetCommunityByID.community, tests[2].outputGetCommunityByID.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[2].inputGetPostByID.postID).Return(tests[2].outputGetPostByID.post, tests[2].outputGetPostByID.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[2].inputGetLikeBySubjectID.postID).Return(tests[2].outputGetLikeBySubjectID.like, tests[2].outputGetLikeBySubjectID.err),

		testRepo.mockCommunityR.EXPECT().GetCommunityByID(ctx, tests[3].inputGetCommunityByID.communityID).Return(tests[3].outputGetCommunityByID.community, tests[3].outputGetCommunityByID.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[3].inputGetPostByID.postID).Return(tests[3].outputGetPostByID.post, tests[3].outputGetPostByID.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[3].inputGetLikeBySubjectID.postID).Return(tests[3].outputGetLikeBySubjectID.like, tests[3].outputGetLikeBySubjectID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.GetCommunityPosts(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestEditCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditCommunityRequest
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}

	type Output struct {
		res *dto.EditCommunityResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info: &dto.EditCommunityRequest{
					CommunityID: "0",
					Name:        "0",
					Image:       "0",
					Info:        "0",
					Admins:      nil,
				},
				userID: "0",
			},
			inputUserCheckCommunity: InputUserCheckCommunity{userID: "0", communityID: "0"},
			outputUserCheckCommunity: OutputUserCheckCommunity{
				err: constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.EditCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestCreatePostCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.CreatePostCommunityRequest
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}

	type Output struct {
		res *dto.CreatePostCommunityResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info: &dto.CreatePostCommunityRequest{
					CommunityID: "0",
					Message:     "0",
					Images:      nil,
				},
				userID: "0",
			},
			inputUserCheckCommunity: InputUserCheckCommunity{userID: "0", communityID: "0"},
			outputUserCheckCommunity: OutputUserCheckCommunity{
				err: constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.CreatePostCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestEditPostCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditPostCommunityRequest
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}

	type Output struct {
		res *dto.EditPostCommunityResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info: &dto.EditPostCommunityRequest{
					CommunityID: "0",
					PostID:      "0",
					Message:     "0",
					Images:      nil,
				},
				userID: "0",
			},
			inputUserCheckCommunity: InputUserCheckCommunity{userID: "0", communityID: "0"},
			outputUserCheckCommunity: OutputUserCheckCommunity{
				err: constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.EditPostCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestDeletePostCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.DeletePostCommunityRequest
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}

	type Output struct {
		res *dto.DeletePostCommunityResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info: &dto.DeletePostCommunityRequest{
					CommunityID: "0",
					PostID:      "0",
				},
				userID: "0",
			},
			inputUserCheckCommunity: InputUserCheckCommunity{userID: "0", communityID: "0"},
			outputUserCheckCommunity: OutputUserCheckCommunity{
				err: constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.DeletePostCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestDeleteCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := NewCommunityService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.DeleteCommunityRequest
		userID string
	}
	type InputUserCheckCommunity struct {
		userID      string
		communityID string
	}
	type OutputUserCheckCommunity struct {
		err error
	}

	type Output struct {
		res *dto.DeleteCommunityResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputUserCheckCommunity  InputUserCheckCommunity
		outputUserCheckCommunity OutputUserCheckCommunity
		output                   Output
	}{
		{
			name: "UserCheckCommunity error",
			input: Input{
				info: &dto.DeleteCommunityRequest{
					CommunityID: "0",
				},
				userID: "0",
			},
			inputUserCheckCommunity: InputUserCheckCommunity{userID: "0", communityID: "0"},
			outputUserCheckCommunity: OutputUserCheckCommunity{
				err: constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().UserCheckCommunity(ctx, tests[0].inputUserCheckCommunity.userID, tests[0].inputUserCheckCommunity.communityID).Return(tests[0].outputUserCheckCommunity.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommunityService.DeleteCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
