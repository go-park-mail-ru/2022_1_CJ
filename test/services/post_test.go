package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/convert"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewPostService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.CreatePostRequest
		userID string
	}
	type InputCreatePost struct {
		post *core.Post
	}
	type OutputCreatePost struct {
		post *core.Post
		err  error
	}

	type InputUserAddPost struct {
		userID string
		postID string
	}
	type OutputUserAddPost struct {
		err error
	}

	type Output struct {
		res *dto.CreatePostResponse
		err error
	}
	var err = errors.Errorf("Can't insert")
	tests := []struct {
		name              string
		input             Input
		inputCreatePost   InputCreatePost
		outputCreatePost  OutputCreatePost
		inputUserAddPost  InputUserAddPost
		outputUserAddPost OutputUserAddPost
		output            Output
	}{
		{
			name: "Can't insert in postRepo",
			input: Input{info: &dto.CreatePostRequest{
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "0"},
			inputCreatePost: InputCreatePost{post: &core.Post{
				AuthorID: "0",
				Message:  "It's my first post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreatePost: OutputCreatePost{post: nil,
				err: err},
			output: Output{nil, err},
		},
		{
			name: "Can't insert in userRepo",
			input: Input{info: &dto.CreatePostRequest{
				Message: "It's my second post!",
				Images:  []string{"src/img.jpg"}},
				userID: "1"},
			inputCreatePost: InputCreatePost{post: &core.Post{
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreatePost: OutputCreatePost{post: &core.Post{
				ID:       "1",
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}},
				err: nil},
			inputUserAddPost:  InputUserAddPost{userID: "1", postID: "1"},
			outputUserAddPost: OutputUserAddPost{err: err},
			output:            Output{nil, err},
		},
		{
			name: "Success",
			input: Input{info: &dto.CreatePostRequest{
				Message: "It's my second post!",
				Images:  []string{"src/img.jpg"}},
				userID: "1"},
			inputCreatePost: InputCreatePost{post: &core.Post{
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreatePost: OutputCreatePost{post: &core.Post{
				ID:       "1",
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}},
				err: nil},
			inputUserAddPost:  InputUserAddPost{userID: "1", postID: "1"},
			outputUserAddPost: OutputUserAddPost{err: nil},
			output:            Output{&dto.CreatePostResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockPostR.EXPECT().CreatePost(ctx, tests[0].inputCreatePost.post).Return(tests[0].outputCreatePost.post, tests[0].outputCreatePost.err),
		testRepo.mockPostR.EXPECT().CreatePost(ctx, tests[1].inputCreatePost.post).Return(tests[1].outputCreatePost.post, tests[1].outputCreatePost.err),
		testRepo.mockUserR.EXPECT().UserAddPost(ctx, tests[1].inputUserAddPost.userID, tests[1].inputUserAddPost.postID).Return(tests[1].outputUserAddPost.err),
		testRepo.mockPostR.EXPECT().CreatePost(ctx, tests[2].inputCreatePost.post).Return(tests[2].outputCreatePost.post, tests[2].outputCreatePost.err),
		testRepo.mockUserR.EXPECT().UserAddPost(ctx, tests[2].inputUserAddPost.userID, tests[2].inputUserAddPost.postID).Return(tests[2].outputUserAddPost.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.PostService.CreatePost(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewPostService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetPostRequest
	}

	type Output struct {
		res *dto.GetPostResponse
		err error
	}
	type OutputGetPost struct {
		post *core.Post
		err  error
	}

	tests := []struct {
		name          string
		input         Input
		outputGetPost OutputGetPost
		output        Output
	}{
		{
			name: "Can't find in postRepo",
			input: Input{info: &dto.GetPostRequest{
				PostID: "0"}},
			outputGetPost: OutputGetPost{post: nil, err: constants.ErrDBNotFound},

			output: Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{info: &dto.GetPostRequest{
				PostID: "677be1d2-9b64-48e9-9341-5ba0c2f57686"}},
			outputGetPost: OutputGetPost{post: &core.Post{
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}, err: nil},
			output: Output{&dto.GetPostResponse{Post: convert.Post2DTO(&core.Post{
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}, &core.User{ID: "1"})}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[0].input.info.PostID).Return(tests[0].outputGetPost.post, tests[0].outputGetPost.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[1].input.info.PostID).Return(tests[1].outputGetPost.post, tests[1].outputGetPost.err),
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].outputGetPost.post.AuthorID).Return(&core.User{ID: "1"}, nil),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.PostService.GetPost(dbUserImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestEditPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewPostService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditPostRequest
		userID string
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputUserCheckPost struct {
		user   *core.User
		postID string
	}
	type OutputUserCheckPost struct {
		err error
	}

	type InputGetPostByID struct {
		postID string
	}

	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}

	type InputEditPost struct {
		post *core.Post
	}

	type OutputEditPost struct {
		post *core.Post
		err  error
	}

	type Output struct {
		res *dto.EditPostResponse
		err error
	}

	tests := []struct {
		name                string
		input               Input
		inputGetUserByID    InputGetUserByID
		outputGetUserByID   OutputGetUserByID
		inputUserCheckPost  InputUserCheckPost
		outputUserCheckPost OutputUserCheckPost
		inputGetPostByID    InputGetPostByID
		outputGetPostByID   OutputGetPostByID
		inputEditPost       InputEditPost
		outputEditPost      OutputEditPost
		output              Output
	}{
		{
			name: "Can't find User in db",
			input: Input{info: &dto.EditPostRequest{
				PostID:  "1",
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "0"},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Don't find posts in user by UserId",
			input: Input{info: &dto.EditPostRequest{
				PostID:  "1",
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "1"},
			inputGetUserByID:    InputGetUserByID{userID: "1"},
			outputGetUserByID:   OutputGetUserByID{user: &core.User{ID: "1"}, err: nil},
			inputUserCheckPost:  InputUserCheckPost{user: &core.User{ID: "1"}, postID: "1"},
			outputUserCheckPost: OutputUserCheckPost{err: constants.ErrDBNotFound},
			output:              Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Don't find post by postID",
			input: Input{info: &dto.EditPostRequest{
				PostID:  "1",
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "2"},
			inputGetUserByID:    InputGetUserByID{userID: "2"},
			outputGetUserByID:   OutputGetUserByID{user: &core.User{ID: "2"}, err: nil},
			inputUserCheckPost:  InputUserCheckPost{user: &core.User{ID: "2"}, postID: "1"},
			outputUserCheckPost: OutputUserCheckPost{nil},
			inputGetPostByID:    InputGetPostByID{postID: "1"},
			outputGetPostByID:   OutputGetPostByID{nil, constants.ErrDBNotFound},
			output:              Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{info: &dto.EditPostRequest{
				PostID:  "1",
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "3"},
			inputGetUserByID:    InputGetUserByID{userID: "3"},
			outputGetUserByID:   OutputGetUserByID{user: &core.User{ID: "3"}, err: nil},
			inputUserCheckPost:  InputUserCheckPost{user: &core.User{ID: "3"}, postID: "1"},
			outputUserCheckPost: OutputUserCheckPost{nil},
			inputGetPostByID:    InputGetPostByID{postID: "1"},
			outputGetPostByID:   OutputGetPostByID{&core.Post{ID: "1"}, nil},
			inputEditPost: InputEditPost{&core.Post{
				ID:       "1",
				Message:  "It's my first post!",
				Images:   []string{"src/img.jpg"},
				AuthorID: "3"}},
			outputEditPost: OutputEditPost{post: &core.Post{
				ID:       "1",
				Message:  "It's my first post!",
				Images:   []string{"src/img.jpg"},
				AuthorID: "3"},
				err: nil},
			output: Output{&dto.EditPostResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[1].inputGetUserByID.userID).Return(tests[1].outputGetUserByID.user, tests[1].outputGetUserByID.err),
		testRepo.mockUserR.EXPECT().UserCheckPost(ctx, tests[1].inputUserCheckPost.user, tests[1].inputUserCheckPost.postID).Return(tests[1].outputUserCheckPost.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[2].inputGetUserByID.userID).Return(tests[2].outputGetUserByID.user, tests[2].outputGetUserByID.err),
		testRepo.mockUserR.EXPECT().UserCheckPost(ctx, tests[2].inputUserCheckPost.user, tests[2].inputUserCheckPost.postID).Return(tests[2].outputUserCheckPost.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[2].inputGetPostByID.postID).Return(tests[2].outputGetPostByID.post, tests[2].outputGetPostByID.err),

		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[3].inputGetUserByID.userID).Return(tests[3].outputGetUserByID.user, tests[3].outputGetUserByID.err),
		testRepo.mockUserR.EXPECT().UserCheckPost(ctx, tests[3].inputUserCheckPost.user, tests[3].inputUserCheckPost.postID).Return(tests[3].outputUserCheckPost.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[3].inputGetPostByID.postID).Return(tests[3].outputGetPostByID.post, tests[3].outputGetPostByID.err),
		testRepo.mockPostR.EXPECT().EditPost(ctx, tests[3].inputEditPost.post).Return(tests[3].outputEditPost.post, tests[3].outputEditPost.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.PostService.EditPost(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbUserImpl := service.NewPostService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.DeletePostRequest
		userID string
	}
	type InputGetPostByID struct {
		postID string
	}
	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}

	type InputDeletePost struct {
		postID string
	}
	type OutputDeletePost struct {
		err error
	}

	type InputUserDeletePost struct {
		userID string
		postID string
	}

	type OutputUserDeletePost struct {
		err error
	}

	type Output struct {
		res *dto.DeletePostResponse
		err error
	}
	var err = errors.Errorf("Can't delete")
	tests := []struct {
		name                 string
		input                Input
		inputGetPostByID     InputGetPostByID
		outputGetPostByID    OutputGetPostByID
		inputDeletePost      InputDeletePost
		outputDeletePost     OutputDeletePost
		inputUserDeletePost  InputUserDeletePost
		outputUserDeletePost OutputUserDeletePost
		output               Output
	}{
		{
			name:              "Can't find post in db",
			input:             Input{info: &dto.DeletePostRequest{PostID: "0"}, userID: "0"},
			inputGetPostByID:  InputGetPostByID{postID: "0"},
			outputGetPostByID: OutputGetPostByID{post: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name:              "Don't match userId and authorId",
			input:             Input{info: &dto.DeletePostRequest{PostID: "1"}, userID: "1"},
			inputGetPostByID:  InputGetPostByID{postID: "1"},
			outputGetPostByID: OutputGetPostByID{post: &core.Post{ID: "1", AuthorID: "2"}, err: nil},
			output:            Output{nil, constants.ErrAuthorIDMismatch},
		},
		{
			name:              "Can't delete post",
			input:             Input{info: &dto.DeletePostRequest{PostID: "2"}, userID: "2"},
			inputGetPostByID:  InputGetPostByID{postID: "2"},
			outputGetPostByID: OutputGetPostByID{post: &core.Post{ID: "2", AuthorID: "2"}, err: nil},
			inputDeletePost:   InputDeletePost{postID: "2"},
			outputDeletePost:  OutputDeletePost{err: err},
			output:            Output{res: nil, err: err},
		},
		{
			name:                 "Can't delete post in user",
			input:                Input{info: &dto.DeletePostRequest{PostID: "3"}, userID: "3"},
			inputGetPostByID:     InputGetPostByID{postID: "3"},
			outputGetPostByID:    OutputGetPostByID{post: &core.Post{ID: "3", AuthorID: "3"}, err: nil},
			inputDeletePost:      InputDeletePost{postID: "3"},
			outputDeletePost:     OutputDeletePost{err: nil},
			inputUserDeletePost:  InputUserDeletePost{userID: "3", postID: "3"},
			outputUserDeletePost: OutputUserDeletePost{err: err},
			output:               Output{res: nil, err: err},
		},
		{
			name:                 "Success",
			input:                Input{info: &dto.DeletePostRequest{PostID: "4"}, userID: "4"},
			inputGetPostByID:     InputGetPostByID{postID: "4"},
			outputGetPostByID:    OutputGetPostByID{post: &core.Post{ID: "4", AuthorID: "4"}, err: nil},
			inputDeletePost:      InputDeletePost{postID: "4"},
			outputDeletePost:     OutputDeletePost{err: nil},
			inputUserDeletePost:  InputUserDeletePost{userID: "4", postID: "4"},
			outputUserDeletePost: OutputUserDeletePost{err: nil},
			output:               Output{res: &dto.DeletePostResponse{}, err: nil},
		},
	}

	gomock.InOrder(
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[0].inputGetPostByID.postID).Return(tests[0].outputGetPostByID.post, tests[0].outputGetPostByID.err),

		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[1].inputGetPostByID.postID).Return(tests[1].outputGetPostByID.post, tests[1].outputGetPostByID.err),

		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[2].inputGetPostByID.postID).Return(tests[2].outputGetPostByID.post, tests[2].outputGetPostByID.err),
		testRepo.mockPostR.EXPECT().DeletePost(ctx, tests[2].inputDeletePost.postID).Return(tests[2].outputDeletePost.err),

		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[3].inputGetPostByID.postID).Return(tests[3].outputGetPostByID.post, tests[3].outputGetPostByID.err),
		testRepo.mockPostR.EXPECT().DeletePost(ctx, tests[3].inputDeletePost.postID).Return(tests[3].outputDeletePost.err),
		testRepo.mockUserR.EXPECT().UserDeletePost(ctx, tests[3].inputUserDeletePost.userID, tests[3].inputUserDeletePost.postID).Return(tests[3].outputUserDeletePost.err),

		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[4].inputGetPostByID.postID).Return(tests[4].outputGetPostByID.post, tests[4].outputGetPostByID.err),
		testRepo.mockPostR.EXPECT().DeletePost(ctx, tests[4].inputDeletePost.postID).Return(tests[4].outputDeletePost.err),
		testRepo.mockUserR.EXPECT().UserDeletePost(ctx, tests[4].inputUserDeletePost.userID, tests[4].inputUserDeletePost.postID).Return(tests[4].outputUserDeletePost.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := service.PostService.DeletePost(dbUserImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
