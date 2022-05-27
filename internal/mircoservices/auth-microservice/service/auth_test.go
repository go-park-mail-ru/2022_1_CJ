package auth_service

import (
	"context"
	auth_core "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/core"
	authdto "github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *authdto.LoginUserRequest
	}
	type InputGetUserByEmail struct {
		email string
	}
	type OutputGetUserByEmail struct {
		authUser *auth_core.User
		err      error
	}

	type Output struct {
		res *authdto.LoginUserResponse
		err error
	}
	var err = errors.Errorf("GetUserByEmail error")
	tests := []struct {
		name                 string
		input                Input
		inputGetUserByEmail  InputGetUserByEmail
		outputGetUserByEmail OutputGetUserByEmail
		output               Output
	}{
		{
			name: "GetUserByEmail error",
			input: Input{info: &authdto.LoginUserRequest{
				Email:    "email@e",
				Password: "1234"}},
			inputGetUserByEmail: InputGetUserByEmail{
				email: "email@e",
			},
			outputGetUserByEmail: OutputGetUserByEmail{authUser: nil,
				err: err},
			output: Output{nil, err},
		},
	}
	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByEmail(ctx, tests[0].inputGetUserByEmail.email).Return(tests[0].outputGetUserByEmail.authUser, tests[0].outputGetUserByEmail.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := AuthService.LoginUser(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *authdto.SignupUserRequest
	}
	type InputCheckUserEmailExistence struct {
		email string
	}
	type OutputCheckUserEmailExistence struct {
		res bool
		err error
	}
	// Unused struct
	type InputCreateUser struct {
		//user *auth_core.User
	}
	type OutputCreateUser struct {
		//res string
		//err error
	}

	type Output struct {
		res *authdto.SignupUserResponse
		err error
	}
	var err = errors.Errorf("CheckUserEmailExistence error")
	tests := []struct {
		name                          string
		input                         Input
		inputCheckUserEmailExistence  InputCheckUserEmailExistence
		outputCheckUserEmailExistence OutputCheckUserEmailExistence
		inputCreateUser               InputCreateUser
		outputCreateUser              OutputCreateUser
		output                        Output
	}{
		{
			name: "GetUserByEmail error",
			input: Input{info: &authdto.SignupUserRequest{
				Email:    "email@e",
				Password: "1234"}},
			inputCheckUserEmailExistence: InputCheckUserEmailExistence{
				email: "email@e",
			},
			outputCheckUserEmailExistence: OutputCheckUserEmailExistence{res: false,
				err: err},
			output: Output{nil, err},
		},
	}
	gomock.InOrder(
		testRepo.mockUserR.EXPECT().CheckUserEmailExistence(ctx, tests[0].inputCheckUserEmailExistence.email).Return(tests[0].outputCheckUserEmailExistence.res, tests[0].outputCheckUserEmailExistence.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := AuthService.SignupUser(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
