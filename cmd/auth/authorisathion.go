package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/controller"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/handler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"time"
)

const (
	configPathEnvVar = "CONFIG_AUTH_PATH"
	configHost       = "microservice_auth.host"
	configPost       = "microservice_auth.port"
	configNetwork    = "microservice_auth.network"
)

func main() {
	// -------------------- Set up viper (config) -------------------- //

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(configPathEnvVar))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("fatal error config file: %s \n", err)
		os.Exit(1)
	}
	// -------------------- Set up logging -------------------- //

	log := logrus.New()

	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	switch viper.GetString("logging.level") {
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "notice":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	log.SetFormatter(&formatter)

	log.Infof("log level: %s", log.Level.String())

	// -------------------- Set database -------------------- //

	clientOptions := options.Client().ApplyURI(viper.GetString("db.connection_string"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("failed to establish mongo db connection: %s", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("failed to check (ping) mongo db connection: %s", err)
	}

	log.Info("connected to MongoDB")
	mongoDB := client.Database(viper.GetString("db.database"))

	//-------------------- Set up service -------------------- //

	listenAddr := viper.GetString(configHost) + ":" + viper.GetString(configPost)
	lis, err := net.Listen(viper.GetString(configNetwork), listenAddr)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
	)

	handler.RegisterUserAuthServer(server, controller.CreateAuthServer(mongoDB, logrus.NewEntry(log)))

	log.Infof("Server started listen at", listenAddr)

	err = server.Serve(lis)
	if err != nil {
		return
	}
}
