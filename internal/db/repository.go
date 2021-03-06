package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	UserRepo      UserRepository
	FriendsRepo   FriendsRepository
	PostRepo      PostRepository
	ChatRepo      ChatRepository
	LikeRepo      LikeRepository
	CommunityRepo CommunityRepository
	CommentRepo   CommentRepository
}

func NewRepository(dbConn *mongo.Database) (*Repository, error) {
	var err error
	repository := new(Repository)

	repository.UserRepo, err = NewUserRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	repository.FriendsRepo, err = NewFriendsRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create friends repository: %w", err)
	}

	repository.PostRepo, err = NewPostRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create post repository: %w", err)
	}

	repository.ChatRepo, err = NewChatRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create chats repository: %w", err)
	}

	repository.LikeRepo, err = NewLikeRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create like repository: %w", err)
	}

	repository.CommunityRepo, err = NewCommunityRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create community repository: %w", err)
	}

	repository.CommentRepo, err = NewCommentRepository(dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment repository: %w", err)
	}

	return repository, nil
}
