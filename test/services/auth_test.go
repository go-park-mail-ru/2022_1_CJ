package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.SignupUserRequest
	}
	type InputCheckUserEmailExistence struct {
		email string
	}
	type OutputCheckUserEmailExistence struct {
		exists bool
		err    error
	}

	type InputCreateUser struct {
		user *core.User
	}

	type OutputCreateUser struct {
		err error
	}
	type InputCreateFriends struct {
		friendsID string
		userID    string
	}

	type OutputCreateFriends struct {
		err error
	}
	type Output struct {
		res *dto.SignupUserResponse
		err error
	}
	tests := []struct {
		name                          string
		input                         Input
		inputCheckUserEmailExistence  InputCheckUserEmailExistence
		outputCheckUserEmailExistence OutputCheckUserEmailExistence
		inputCreateUser               InputCreateUser
		outputCreateUser              OutputCreateUser
		inputCreateFriends            InputCreateFriends
		outputCreateFriends           OutputCreateFriends
		output                        Output
	}{
		{
			name: "With taken email",
			input: Input{info: &dto.SignupUserRequest{Name: common.UserName{First: "Sasha", Last: "Web"},
				Email:    "wrong",
				Password: "1234"}},
			inputCheckUserEmailExistence:  InputCheckUserEmailExistence{email: "wrong"},
			outputCheckUserEmailExistence: OutputCheckUserEmailExistence{exists: true, err: nil},
			output:                        Output{nil, constants.ErrEmailAlreadyTaken},
		},
		//{
		//	name: "Can't insert in userRepo",
		//	input: Input{info: &dto.SignupUserRequest{Name: common.UserName{First: "Sasha", Last: "Web"},
		//											Email:    "SashaWeb@mail.ru",
		//											Password: "12 34"}},
		//	inputCheckUserEmailExistence:  InputCheckUserEmailExistence{email: "SashaWeb@mail.ru"},
		//	outputCheckUserEmailExistence: OutputCheckUserEmailExistence{exists: false, err: nil},
		//	inputCreateUser: InputCreateUser{post: &core.User{
		//											Name: common.UserName{First: "Sasha", Last: "Web"},
		//											Email: "SashaWeb@mail.ru"}},
		//	outputCreateUser: OutputCreateUser{},
		//	output:           Output{nil, nil},
		//},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().CheckUserEmailExistence(ctx, tests[0].inputCheckUserEmailExistence.email).Return(
			tests[0].outputCheckUserEmailExistence.exists, tests[0].outputCheckUserEmailExistence.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.AuthService.SignupUser(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.LoginUserRequest
	}
	type InputGetUserByEmail struct {
		email string
	}
	type OutputGetUserByEmail struct {
		user *core.User
		err  error
	}

	type Output struct {
		res *dto.LoginUserResponse
		err error
	}
	tests := []struct {
		name                  string
		input                 Input
		inputGetUserByEmail   InputGetUserByEmail
		outputGetUserByEmaile OutputGetUserByEmail
		output                Output
	}{
		{
			name:                  "With wrong post",
			input:                 Input{info: &dto.LoginUserRequest{Email: "wrong", Password: "1234"}},
			inputGetUserByEmail:   InputGetUserByEmail{email: "wrong"},
			outputGetUserByEmaile: OutputGetUserByEmail{user: nil, err: constants.ErrDBNotFound},
			output:                Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByEmail(ctx, tests[0].inputGetUserByEmail.email).Return(
			tests[0].outputGetUserByEmaile.user, tests[0].outputGetUserByEmaile.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.AuthService.LoginUser(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
