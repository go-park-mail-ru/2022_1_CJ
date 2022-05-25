package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.SignupUserRequest
		userID string
		token  string
	}

	type InputCreateUser struct {
		user *core.User
	}

	type OutputCreateUser struct {
		err error
	}
	type InputCreateFriends struct {
		userID string
	}

	type OutputCreateFriends struct {
		err error
	}
	type Output struct {
		res *dto.SignupUserResponse
		err error
	}
	expected, _ := utils.GenerateCSRFToken("2")

	tests := []struct {
		name                string
		input               Input
		inputCreateUser     InputCreateUser
		outputCreateUser    OutputCreateUser
		inputCreateFriends  InputCreateFriends
		outputCreateFriends OutputCreateFriends
		output              Output
	}{
		{
			name: "Can't create user",
			input: Input{info: &dto.SignupUserRequest{Name: common.UserName{First: "Sasha", Last: "Web"},
				Email:    "wrong",
				Password: "1234",
			}, userID: "0", token: "0"},
			inputCreateUser: InputCreateUser{user: &core.User{
				ID:    "0",
				Name:  common.UserName{First: "Sasha", Last: "Web"},
				Email: "wrong",
			}},
			outputCreateUser: OutputCreateUser{err: constants.ErrEmailAlreadyTaken},
			output:           Output{nil, constants.ErrEmailAlreadyTaken},
		},
		{
			name: "CreateFriends error",
			input: Input{info: &dto.SignupUserRequest{Name: common.UserName{First: "Sasha", Last: "Web"},
				Email:    "wrong",
				Password: "1234",
			}, userID: "1", token: "0"},
			inputCreateUser: InputCreateUser{user: &core.User{
				ID:    "1",
				Name:  common.UserName{First: "Sasha", Last: "Web"},
				Email: "wrong",
			}},
			outputCreateUser: OutputCreateUser{err: nil},
			inputCreateFriends: InputCreateFriends{
				userID: "1",
			},
			outputCreateFriends: OutputCreateFriends{err: constants.ErrEmailAlreadyTaken},
			output:              Output{nil, constants.ErrEmailAlreadyTaken},
		},
		{
			name: "Success",
			input: Input{info: &dto.SignupUserRequest{Name: common.UserName{First: "Sasha", Last: "Web"},
				Email:    "wrong",
				Password: "1234",
			}, userID: "2", token: "0"},
			inputCreateUser: InputCreateUser{user: &core.User{
				ID:    "2",
				Name:  common.UserName{First: "Sasha", Last: "Web"},
				Email: "wrong",
			}},
			outputCreateUser: OutputCreateUser{err: nil},
			inputCreateFriends: InputCreateFriends{
				userID: "2",
			},
			outputCreateFriends: OutputCreateFriends{err: nil},
			output:              Output{&dto.SignupUserResponse{AuthToken: "0", CSRFToken: expected}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().CreateUser(ctx, tests[0].inputCreateUser.user).Return(tests[0].outputCreateUser.err),
		testRepo.mockUserR.EXPECT().CreateUser(ctx, tests[1].inputCreateUser.user).Return(tests[1].outputCreateUser.err),
		testRepo.mockFriendsR.EXPECT().CreateFriends(ctx, tests[1].inputCreateFriends.userID).Return(tests[1].outputCreateFriends.err),
		testRepo.mockUserR.EXPECT().CreateUser(ctx, tests[2].inputCreateUser.user).Return(tests[2].outputCreateUser.err),
		testRepo.mockFriendsR.EXPECT().CreateFriends(ctx, tests[2].inputCreateFriends.userID).Return(tests[2].outputCreateFriends.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := AuthService.SignupUser(dbUserImpl, ctx, test.input.info, test.input.userID, test.input.token)
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

	TestBD, _ := TestRepositories(t, ctrl)
	dbUserImpl := NewAuthService(TestLogger(t), TestBD)

	ctx := context.Background()
	expected, _ := utils.GenerateCSRFToken("0")
	t.Run("Success", func(t *testing.T) {

		res, _ := AuthService.LoginUser(dbUserImpl, ctx, "0", "1")
		if !assert.Equal(t, res.CSRFToken, expected) {
			t.Error("got : ", res, " expected :", res.CSRFToken)
		}
	})

}
