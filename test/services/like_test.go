package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncreaseLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewLikeService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.IncreaseLikeRequest
		userID string
	}

	type InputGetLikeBySubjectID struct {
		subjectID string
	}
	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type InputIncreaseLike struct {
		subjectID string
		userID    string
	}
	type OutputIncreaseLike struct {
		err error
	}

	type Output struct {
		res *dto.IncreaseLikeResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		inputCreateDialog        InputIncreaseLike
		outputIncreaseLike       OutputIncreaseLike
		output                   Output
	}{
		{
			name: "ErrBadJson error",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID:  "",
				PhotoID: "",
			}, userID: "0"},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "ErrBadJson error 1",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID:  "1234",
				PhotoID: "12543",
			}, userID: "0"},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "ErrBadJson error 2",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID:  "1234",
				PhotoID: "12543",
			}},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "GetLikeBySubjectID error",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID: "1",
			}, userID: "1"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "1"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "IncreaseLike error",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID: "2",
			}, userID: "2"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "2"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			inputCreateDialog: InputIncreaseLike{
				subjectID: "2",
				userID:    "2",
			},
			outputIncreaseLike: OutputIncreaseLike{err: constants.ErrDBNotFound},
			output:             Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "success",
			input: Input{info: &dto.IncreaseLikeRequest{
				PostID: "3",
			}, userID: "3"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "3"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			inputCreateDialog: InputIncreaseLike{
				subjectID: "3",
				userID:    "3",
			},
			outputIncreaseLike: OutputIncreaseLike{err: nil},
			output:             Output{&dto.IncreaseLikeResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[3].inputGetLikeBySubjectID.subjectID).Return(tests[3].outputGetLikeBySubjectID.like, tests[3].outputGetLikeBySubjectID.err),

		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[4].inputGetLikeBySubjectID.subjectID).Return(tests[4].outputGetLikeBySubjectID.like, tests[4].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().IncreaseLike(ctx, tests[4].inputCreateDialog.subjectID, tests[4].inputCreateDialog.userID).Return(tests[4].outputIncreaseLike.err),

		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[5].inputGetLikeBySubjectID.subjectID).Return(tests[5].outputGetLikeBySubjectID.like, tests[5].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().IncreaseLike(ctx, tests[5].inputCreateDialog.subjectID, tests[5].inputCreateDialog.userID).Return(tests[5].outputIncreaseLike.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.LikeService.IncreaseLike(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetLikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewLikeService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.GetLikePostRequest
		userID string
	}

	type InputGetLikeBySubjectID struct {
		subjectID string
	}
	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type Output struct {
		res *dto.GetLikePostResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		output                   Output
	}{
		{
			name: "GetLikeBySubjectID error",
			input: Input{info: &dto.GetLikePostRequest{
				PostID: "1",
			}, userID: "1"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "1"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},

		{
			name: "success",
			input: Input{info: &dto.GetLikePostRequest{
				PostID: "3",
			}, userID: "3"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "3"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			output: Output{&dto.GetLikePostResponse{Likes: convert.Like2DTO(&core.Like{}, "3")}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[0].inputGetLikeBySubjectID.subjectID).Return(tests[0].outputGetLikeBySubjectID.like, tests[0].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[1].inputGetLikeBySubjectID.subjectID).Return(tests[1].outputGetLikeBySubjectID.like, tests[1].outputGetLikeBySubjectID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.LikeService.GetLikePost(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetLikePhoto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewLikeService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.GetLikePhotoRequest
		userID string
	}

	type InputGetLikeBySubjectID struct {
		subjectID string
	}
	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type Output struct {
		res *dto.GetLikePhotoResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		output                   Output
	}{
		{
			name: "GetLikeBySubjectID error",
			input: Input{info: &dto.GetLikePhotoRequest{
				PhotoID: "1",
			}, userID: "1"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "1"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},

		{
			name: "success",
			input: Input{info: &dto.GetLikePhotoRequest{
				PhotoID: "3",
			}, userID: "3"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "3"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			output: Output{&dto.GetLikePhotoResponse{Likes: convert.Like2DTO(&core.Like{}, "3")}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[0].inputGetLikeBySubjectID.subjectID).Return(tests[0].outputGetLikeBySubjectID.like, tests[0].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[1].inputGetLikeBySubjectID.subjectID).Return(tests[1].outputGetLikeBySubjectID.like, tests[1].outputGetLikeBySubjectID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.LikeService.GetLikePhoto(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestReduceLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewLikeService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.ReduceLikeRequest
		userID string
	}

	type InputGetLikeBySubjectID struct {
		subjectID string
	}
	type OutputGetLikeBySubjectID struct {
		like *core.Like
		err  error
	}

	type InputReduceLike struct {
		subjectID string
		userID    string
	}
	type OutputReduceLike struct {
		err error
	}

	type Output struct {
		res *dto.ReduceLikeResponse
		err error
	}

	tests := []struct {
		name                     string
		input                    Input
		inputGetLikeBySubjectID  InputGetLikeBySubjectID
		outputGetLikeBySubjectID OutputGetLikeBySubjectID
		inputReduceLike          InputReduceLike
		outputReduceLike         OutputReduceLike
		output                   Output
	}{
		{
			name: "ErrBadJson error",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID:  "",
				PhotoID: "",
			}, userID: "0"},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "ErrBadJson error 1",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID:  "1234",
				PhotoID: "12543",
			}, userID: "0"},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "ErrBadJson error 2",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID:  "1234",
				PhotoID: "12543",
			}},
			output: Output{nil, constants.ErrBadJson},
		},
		{
			name: "GetLikeBySubjectID error",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID: "1",
			}, userID: "1"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "1"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: nil,
				err:  constants.ErrDBNotFound,
			},
			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "IncreaseLike error",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID: "2",
			}, userID: "2"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "2"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			inputReduceLike: InputReduceLike{
				subjectID: "2",
				userID:    "2",
			},
			outputReduceLike: OutputReduceLike{err: constants.ErrDBNotFound},
			output:           Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "success",
			input: Input{info: &dto.ReduceLikeRequest{
				PostID: "3",
			}, userID: "3"},
			inputGetLikeBySubjectID: InputGetLikeBySubjectID{subjectID: "3"},
			outputGetLikeBySubjectID: OutputGetLikeBySubjectID{
				like: &core.Like{},
				err:  nil,
			},
			inputReduceLike: InputReduceLike{
				subjectID: "3",
				userID:    "3",
			},
			outputReduceLike: OutputReduceLike{err: nil},
			output:           Output{&dto.ReduceLikeResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[3].inputGetLikeBySubjectID.subjectID).Return(tests[3].outputGetLikeBySubjectID.like, tests[3].outputGetLikeBySubjectID.err),

		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[4].inputGetLikeBySubjectID.subjectID).Return(tests[4].outputGetLikeBySubjectID.like, tests[4].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().ReduceLike(ctx, tests[4].inputReduceLike.subjectID, tests[4].inputReduceLike.userID).Return(tests[4].outputReduceLike.err),

		testRepo.mockLikeR.EXPECT().GetLikeBySubjectID(ctx, tests[5].inputGetLikeBySubjectID.subjectID).Return(tests[5].outputGetLikeBySubjectID.like, tests[5].outputGetLikeBySubjectID.err),
		testRepo.mockLikeR.EXPECT().ReduceLike(ctx, tests[5].inputReduceLike.subjectID, tests[5].inputReduceLike.userID).Return(tests[5].outputReduceLike.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.LikeService.ReduceLike(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
