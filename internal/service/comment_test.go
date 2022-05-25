package service

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/dto"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommentImpl := NewCommentService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.CreateCommentRequest
		userID string
	}
	type InputCreateComment struct {
		comment *core.Comment
	}
	type OutputCreateComment struct {
		comment *core.Comment
		err     error
	}

	type InputPostAddComment struct {
		postID    string
		commentID string
	}
	type OutputPostAddComment struct {
		err error
	}

	type Output struct {
		res *dto.CreateCommentResponse
		err error
	}
	var err = errors.Errorf("Can't insert")
	tests := []struct {
		name                 string
		input                Input
		inputCreateComment   InputCreateComment
		outputCreateComment  OutputCreateComment
		inputPostAddComment  InputPostAddComment
		outputPostAddComment OutputPostAddComment
		output               Output
	}{
		{
			name: "Can't insert in commentRepo",
			input: Input{info: &dto.CreateCommentRequest{
				PostID:  "0",
				Message: "It's my first post!",
				Images:  []string{"src/img.jpg"}},
				userID: "0"},
			inputCreateComment: InputCreateComment{comment: &core.Comment{
				AuthorID: "0",
				Message:  "It's my first post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreateComment: OutputCreateComment{comment: nil,
				err: err},
			output: Output{nil, err},
		},
		{
			name: "Can't insert in postRepo",
			input: Input{info: &dto.CreateCommentRequest{
				PostID:  "1",
				Message: "It's my second post!",
				Images:  []string{"src/img.jpg"}},
				userID: "1"},
			inputCreateComment: InputCreateComment{comment: &core.Comment{
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreateComment: OutputCreateComment{comment: &core.Comment{
				ID:       "1",
				AuthorID: "1",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}},
				err: nil},
			inputPostAddComment: InputPostAddComment{
				postID:    "1",
				commentID: "1",
			},
			outputPostAddComment: OutputPostAddComment{err: err},
			output:               Output{nil, err},
		},
		{
			name: "Success",
			input: Input{info: &dto.CreateCommentRequest{
				PostID:  "2",
				Message: "It's my second post!",
				Images:  []string{"src/img.jpg"}},
				userID: "2"},
			inputCreateComment: InputCreateComment{comment: &core.Comment{
				AuthorID: "2",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}}},
			outputCreateComment: OutputCreateComment{comment: &core.Comment{
				ID:       "2",
				AuthorID: "2",
				Message:  "It's my second post!",
				Images:   []string{"src/img.jpg"}},
				err: nil},
			inputPostAddComment: InputPostAddComment{
				postID:    "2",
				commentID: "2",
			},
			outputPostAddComment: OutputPostAddComment{err: nil},
			output:               Output{&dto.CreateCommentResponse{}, nil},
		},
	}

	gomock.InOrder(
		testRepo.mockCommentR.EXPECT().CreateComment(ctx, tests[0].inputCreateComment.comment).Return(tests[0].outputCreateComment.comment, tests[0].outputCreateComment.err),

		testRepo.mockCommentR.EXPECT().CreateComment(ctx, tests[1].inputCreateComment.comment).Return(tests[1].outputCreateComment.comment, tests[1].outputCreateComment.err),
		testRepo.mockPostR.EXPECT().PostAddComment(ctx, tests[1].inputPostAddComment.postID, tests[1].inputPostAddComment.commentID).Return(tests[1].outputPostAddComment.err),

		testRepo.mockCommentR.EXPECT().CreateComment(ctx, tests[2].inputCreateComment.comment).Return(tests[2].outputCreateComment.comment, tests[2].outputCreateComment.err),
		testRepo.mockPostR.EXPECT().PostAddComment(ctx, tests[2].inputPostAddComment.postID, tests[2].inputPostAddComment.commentID).Return(tests[2].outputPostAddComment.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommentService.CreateComment(dbCommentImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestGetComments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommentImpl := NewCommentService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info *dto.GetCommentsRequest
	}

	type Output struct {
		res *dto.GetCommentsResponse
		err error
	}

	type InputGetPostByID struct {
		postID string
	}

	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}

	tests := []struct {
		name              string
		input             Input
		inputGetPostByID  InputGetPostByID
		outputGetPostByID OutputGetPostByID
		output            Output
	}{
		{
			name: "Can't find in postRepo",
			input: Input{info: &dto.GetCommentsRequest{
				PostID: "0", Limit: -1, Page: 1}},
			inputGetPostByID:  InputGetPostByID{postID: "0"},
			outputGetPostByID: OutputGetPostByID{nil, constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
		{
			name: "Success",
			input: Input{info: &dto.GetCommentsRequest{
				PostID: "1", Limit: -1, Page: 1}},
			inputGetPostByID: InputGetPostByID{postID: "1"},
			outputGetPostByID: OutputGetPostByID{&core.Post{
				ID:          "1",
				AuthorID:    "1",
				CommentsIDs: nil,
			}, nil},
			output: Output{
				res: &dto.GetCommentsResponse{Comments: nil, AmountPages: 1, Total: 0},
				err: nil,
			},
		},
	}

	gomock.InOrder(
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[0].inputGetPostByID.postID).Return(tests[0].outputGetPostByID.post, tests[0].outputGetPostByID.err),
		testRepo.mockPostR.EXPECT().GetPostByID(ctx, tests[1].inputGetPostByID.postID).Return(tests[1].outputGetPostByID.post, tests[1].outputGetPostByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommentService.GetComments(dbCommentImpl, ctx, test.input.info)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommentImpl := NewCommentService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.DeleteCommentRequest
		userID string
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputGetCommentByID struct {
		commentID string
	}
	type OutputGetCommentByID struct {
		comment *core.Comment
		err     error
	}

	type InputGetPostByID struct {
		postID string
	}

	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}

	type InputPostCheckComment struct {
		post      *core.Post
		commentID string
	}

	type OutputPostCheckComment struct {
		err error
	}

	type InputDeleteComment struct {
		commentID string
	}

	type OutputDeleteComment struct {
		err error
	}

	type InputPostDeleteComment struct {
		postID    string
		commentID string
	}

	type OutputPostDeleteComment struct {
		err error
	}

	type Output struct {
		res *dto.DeleteCommentResponse
		err error
	}

	tests := []struct {
		name                    string
		input                   Input
		inputGetUserByID        InputGetUserByID
		outputGetUserByID       OutputGetUserByID
		inputUserCheckPost      InputGetCommentByID
		outputUserCheckPost     OutputGetCommentByID
		inputGetPostByID        InputGetPostByID
		outputGetPostByID       OutputGetPostByID
		onputPostCheckComment   InputPostCheckComment
		outputPostCheckComment  OutputPostCheckComment
		inputDeleteComment      InputDeleteComment
		outputDeleteComment     OutputDeleteComment
		inputPostDeleteComment  InputPostDeleteComment
		outputPostDeleteComment OutputPostDeleteComment
		output                  Output
	}{
		{
			name: "Can't find User in db",
			input: Input{info: &dto.DeleteCommentRequest{
				PostID:    "0",
				CommentID: "0"},
				userID: "0"},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommentService.DeleteComment(dbCommentImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}

func TestEditComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	TestBD, testRepo := TestRepositories(t, ctrl)
	dbCommentImpl := NewCommentService(TestLogger(t), TestBD)

	ctx := context.Background()

	type Input struct {
		info   *dto.EditCommentRequest
		userID string
	}
	type InputGetUserByID struct {
		userID string
	}
	type OutputGetUserByID struct {
		user *core.User
		err  error
	}

	type InputGetCommentByID struct {
		commentID string
	}
	type OutputGetCommentByID struct {
		comment *core.Comment
		err     error
	}

	type InputGetPostByID struct {
		postID string
	}

	type OutputGetPostByID struct {
		post *core.Post
		err  error
	}

	type InputPostCheckComment struct {
		post      *core.Post
		commentID string
	}

	type OutputPostCheckComment struct {
		err error
	}

	type InputEditComment struct {
		comment *core.Comment
	}

	type OutputEditComment struct {
		comment *core.Comment
		err     error
	}

	type Output struct {
		res *dto.EditCommentResponse
		err error
	}

	tests := []struct {
		name                   string
		input                  Input
		inputGetUserByID       InputGetUserByID
		outputGetUserByID      OutputGetUserByID
		inputUserCheckPost     InputGetCommentByID
		outputUserCheckPost    OutputGetCommentByID
		inputGetPostByID       InputGetPostByID
		outputGetPostByID      OutputGetPostByID
		onputPostCheckComment  InputPostCheckComment
		outputPostCheckComment OutputPostCheckComment
		inputEditPost          InputEditComment
		outputEditPost         OutputEditComment
		output                 Output
	}{
		{
			name: "Can't find User in db",
			input: Input{info: &dto.EditCommentRequest{
				PostID:    "0",
				CommentID: "0",
				Message:   "It's my first post!",
				Images:    []string{"src/img.jpg"}},
				userID: "0"},
			inputGetUserByID:  InputGetUserByID{userID: "0"},
			outputGetUserByID: OutputGetUserByID{user: nil, err: constants.ErrDBNotFound},
			output:            Output{nil, constants.ErrDBNotFound},
		},
	}

	gomock.InOrder(
		testRepo.mockUserR.EXPECT().GetUserByID(ctx, tests[0].inputGetUserByID.userID).Return(tests[0].outputGetUserByID.user, tests[0].outputGetUserByID.err),
	)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			res, errRes := CommentService.EditComment(dbCommentImpl, ctx, test.input.info, test.input.userID)
			if !assert.Equal(t, test.output.res, res) {
				t.Error("got : ", res, " expected :", test.output.res)
			}
			if !assert.Equal(t, test.output.err, errRes) {
				t.Error("got : ", errRes, " expected :", test.output.err)
			}
		})
	}
}
