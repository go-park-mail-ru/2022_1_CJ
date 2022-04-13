package api

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/api/controllers"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIService struct {
	log    *logrus.Entry
	router *echo.Echo
	debug  bool
}

func (svc *APIService) Serve() {
	svc.log.Info("Starting HTTP server")
	listenAddr := viper.GetString("service.bind.address") + ":" + viper.GetString("service.bind.port")
	svc.log.Fatal(svc.router.Start(listenAddr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	if err := svc.router.Shutdown(ctx); err != nil {
		svc.log.Fatal(err)
	}
	return nil
}

func NewAPIService(log *logrus.Entry, dbConn *mongo.Database, debug bool) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: echo.New(),
		debug:  debug,
	}

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()

	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	registry := service.NewRegistry(log, repository)

	authCtrl := controllers.NewAuthController(log, registry)
	userCtrl := controllers.NewUserController(log, registry)
	friendsCtrl := controllers.NewFriendsController(log, registry)
	postCtrl := controllers.NewPostController(log, registry)
	staticCtrl := controllers.NewStaticController(log, registry)
	chatCtrl := controllers.NewChatController(log, repository, registry)

	svc.router.HTTPErrorHandler = svc.httpErrorHandler
	svc.router.Use(svc.XRequestIDMiddleware(), svc.LoggingMiddleware())

	api := svc.router.Group("/api")

	authAPI := api.Group("/auth")

	authAPI.POST("/signup", authCtrl.SignupUser)
	authAPI.POST("/login", authCtrl.LoginUser)
	authAPI.DELETE("/logout", authCtrl.LogoutUser)

	userAPI := api.Group("/user", svc.AuthMiddleware())

	userAPI.GET("/get", userCtrl.GetUserData)
	userAPI.GET("/posts", userCtrl.GetUserPosts)
	userAPI.GET("/feed", userCtrl.GetFeed)
	userAPI.POST("/update_photo", userCtrl.UpdatePhoto)
	userAPI.GET("/profile", userCtrl.GetProfile)
	userAPI.POST("/profile/edit", userCtrl.EditProfile)
	userAPI.GET("/search", userCtrl.SearchUsers)

	friendsAPI := api.Group("/friends", svc.AuthMiddleware())

	friendsAPI.POST("/request", friendsCtrl.SendFriendRequest)
	friendsAPI.POST("/accept", friendsCtrl.AcceptFriendRequest)
	friendsAPI.GET("/requests", friendsCtrl.GetFriendRequests)
	friendsAPI.GET("/get", friendsCtrl.GetFriendsByUserID)
	friendsAPI.DELETE("/delete", friendsCtrl.DeleteFriend)

	postAPI := api.Group("/post", svc.AuthMiddleware())

	postAPI.POST("/create", postCtrl.CreatePost)
	postAPI.GET("/get", postCtrl.GetPost)
	postAPI.PUT("/edit", postCtrl.EditPost)
	postAPI.DELETE("/delete", postCtrl.DeletePost)

	static := api.Group("/static")

	static.POST("/upload", staticCtrl.UploadImage, svc.AuthMiddleware())

	chatAPI := api.Group("/messenger", svc.AuthMiddleware())

	chatAPI.GET("/dialogs", chatCtrl.GetDialogs)
	chatAPI.POST("/create", chatCtrl.CreateDialog)
	chatAPI.POST("/ws", chatCtrl.WsHandler)
	return svc, nil
}
