package service

import (
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	AuthService      AuthService
	OAuthService     OAuthService
	UserService      UserService
	PostService      PostService
	FriendsService   FriendsService
	StaticService    StaticService
	ChatService      ChatService
	LikeService      LikeService
	CommunityService CommunityService
	CommentService   CommentService
}

func NewRegistry(log *logrus.Entry, repository *db.Repository) *Registry {
	registry := new(Registry)

	registry.AuthService = NewAuthService(log, repository)
	registry.OAuthService = NewOAuthService(log, repository)
	registry.UserService = NewUserService(log, repository)
	registry.FriendsService = NewFriendsService(log, repository)
	registry.PostService = NewPostService(log, repository)
	registry.StaticService = NewStaticService(log, repository)
	registry.ChatService = NewChatService(log, repository)
	registry.LikeService = NewLikeService(log, repository)
	registry.CommunityService = NewCommunityService(log, repository)
	registry.CommentService = NewCommentService(log, repository)

	return registry
}
