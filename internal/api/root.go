package api

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/monitoring"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/api/controllers"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/cl"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/handler"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/service"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type APIService struct {
	log     *logrus.Entry
	router  *echo.Echo
	debug   bool
	metrics *monitoring.PrometheusMetrics
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

func NewAPIService(log *logrus.Entry, dbConn *mongo.Database, debug bool, grpcConn *grpc.ClientConn) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: echo.New(),
		debug:  debug,
	}

	authService := cl.NewAuthRepository(log, handler.NewUserAuthClient(grpcConn))

	svc.router.Validator = NewValidator()
	svc.router.Binder = NewBinder()

	repository, err := db.NewRepository(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	registry := service.NewRegistry(log, repository)

	authCtrl := controllers.NewAuthController(log, registry, authService)
	userCtrl := controllers.NewUserController(log, registry)
	friendsCtrl := controllers.NewFriendsController(log, registry)
	postCtrl := controllers.NewPostController(log, registry)
	staticCtrl := controllers.NewStaticController(log, registry)
	likeCtrl := controllers.NewLikeController(log, registry)
	communitiesCtrl := controllers.NewCommunityController(log, registry)
	commentCtrl := controllers.NewCommentController(log, registry)
	chatCtrl := controllers.NewChatController(log, repository, registry)

	svc.router.HTTPErrorHandler = svc.httpErrorHandler

	svc.metrics = monitoring.RegisterMonitoring(svc.router)
	svc.router.Use(svc.XRequestIDMiddleware(), svc.LoggingMiddleware(), svc.AccessLogMiddleware())

	api := svc.router.Group("/api")

	authAPI := api.Group("/auth")

	authAPI.POST("/signup", authCtrl.SignupUser)
	authAPI.POST("/login", authCtrl.LoginUser)
	authAPI.DELETE("/logout", authCtrl.LogoutUser)

	userAPI := api.Group("/user", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	userAPI.GET("/get", userCtrl.GetUserData)
	userAPI.GET("/posts", userCtrl.GetUserPosts)
	userAPI.GET("/feed", userCtrl.GetFeed)
	userAPI.POST("/update_photo", userCtrl.UpdatePhoto)
	userAPI.GET("/search", userCtrl.SearchUsers)
	userAPI.GET("/profile", userCtrl.GetProfile)
	userAPI.POST("/profile/edit", userCtrl.EditProfile)

	friendsAPI := api.Group("/friends", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	friendsAPI.POST("/request", friendsCtrl.SendFriendRequest)
	friendsAPI.POST("/accept", friendsCtrl.AcceptFriendRequest)
	friendsAPI.GET("/requests/outcoming", friendsCtrl.GetOutcomingRequests)
	friendsAPI.GET("/requests/incoming", friendsCtrl.GetIncomingRequests)
	friendsAPI.GET("/get", friendsCtrl.GetFriendsByUserID)
	friendsAPI.DELETE("/delete", friendsCtrl.DeleteFriend)

	postAPI := api.Group("/post", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	postAPI.POST("/create", postCtrl.CreatePost)
	postAPI.GET("/get", postCtrl.GetPost)
	postAPI.PUT("/edit", postCtrl.EditPost)
	postAPI.DELETE("/delete", postCtrl.DeletePost)

	likeAPI := api.Group("/like", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	likeAPI.POST("/increase", likeCtrl.IncreaseLike)
	likeAPI.POST("/reduce", likeCtrl.ReduceLike)
	likeAPI.GET("/post/get", likeCtrl.GetLikePost)
	likeAPI.GET("/photo/get", likeCtrl.GetLikePhoto)

	static := api.Group("/static")

	static.POST("/upload", staticCtrl.UploadImage, svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	chatAPI := api.Group("/messenger", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	chatAPI.GET("/dialogs", chatCtrl.GetDialogs)
	chatAPI.GET("/get", chatCtrl.GetDialog)
	chatAPI.GET("/user_dialog", chatCtrl.GetDialogByUserID)
	chatAPI.POST("/create", chatCtrl.CreateDialog)
	chatAPI.GET("/ws", chatCtrl.WsHandler)

	communitiesAPI := api.Group("/communities", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	communitiesAPI.GET("/get", communitiesCtrl.GetCommunity)
	communitiesAPI.GET("/posts", communitiesCtrl.GetCommunityPosts)
	communitiesAPI.GET("/list", communitiesCtrl.GetUserCommunities)
	communitiesAPI.GET("/managed_list", communitiesCtrl.GetUserManageCommunities)
	communitiesAPI.GET("/full_list", communitiesCtrl.GetCommunities)
	communitiesAPI.GET("/search", communitiesCtrl.SearchCommunities)
	communitiesAPI.GET("/join", communitiesCtrl.JoinCommunity)
	communitiesAPI.GET("/leave", communitiesCtrl.LeaveCommunity)
	communitiesAPI.GET("/followers", communitiesCtrl.GetFollowers)
	communitiesAPI.GET("/mutual_friends", communitiesCtrl.GetMutualFriends)
	communitiesAPI.POST("/create", communitiesCtrl.CreateCommunity)
	communitiesAPI.PUT("/edit", communitiesCtrl.EditCommunity)
	communitiesAPI.DELETE("/delete", communitiesCtrl.DeleteCommunity)
	communitiesAPI.POST("/update_photo", communitiesCtrl.UpdatePhotoCommunity)

	communitiesPostAPI := communitiesAPI.Group("/post")

	communitiesPostAPI.POST("/create", communitiesCtrl.CreatePostCommunity)
	communitiesPostAPI.PUT("/edit", communitiesCtrl.EditPostCommunity)
	communitiesPostAPI.DELETE("/delete", communitiesCtrl.DeletePostCommunity)

	commentAPI := api.Group("/comment", svc.AuthMiddlewareMicro(authService), svc.CSRFMiddleware())

	commentAPI.POST("/create", commentCtrl.CreateComment)
	commentAPI.GET("/get", commentCtrl.GetComments)
	commentAPI.PUT("/edit", commentCtrl.EditComment)
	commentAPI.DELETE("/delete", commentCtrl.DeleteComment)

	return svc, nil
}
