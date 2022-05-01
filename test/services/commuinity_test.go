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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCommunity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommunityImpl := service.NewCommunityService(TestLogger(t), TestBD)

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

			res, errRes := service.CommunityService.CreateCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
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
	dbCommunityImpl := service.NewCommunityService(TestLogger(t), TestBD)

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
			name: "Didn't find user admin in db",
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

			res, errRes := service.CommunityService.GetCommunity(dbCommunityImpl, ctx, test.input.info)
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
	dbCommunityImpl := service.NewCommunityService(TestLogger(t), TestBD)

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
			name: "Didn't find user in db",
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

			res, errRes := service.CommunityService.GetUserManageCommunities(dbCommunityImpl, ctx, test.input.info)
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
	dbCommunityImpl := service.NewCommunityService(TestLogger(t), TestBD)

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

			res, errRes := service.CommunityService.GetCommunities(dbCommunityImpl, ctx, test.input.info)
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
	dbCommunityImpl := service.NewCommunityService(TestLogger(t), TestBD)

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

			res, errRes := service.CommunityService.JoinCommunity(dbCommunityImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
