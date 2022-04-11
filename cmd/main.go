package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core/chat"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/api"
)

const (
	configPathEnvVar = "CONFIG_PATH"
	defaultAddress   = "0.0.0.0"
	defaultPort      = "8080"
)

func main() {
	// -------------------- Set up viper (config) -------------------- //

	viper.AutomaticEnv()

	viper.SetConfigFile(viper.GetString(configPathEnvVar))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("fatal error config file: %s \n", err)
		os.Exit(1)
	}

	viper.SetDefault("service.bind.address", defaultAddress)
	viper.SetDefault("service.bind.port", defaultPort)

	// -------------------- Set up logging -------------------- //

	log := logrus.New()

	formatter := logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	var debug bool
	switch viper.GetString("logging.level") {
	case "warning":
		log.SetLevel(logrus.WarnLevel)
	case "notice":
		log.SetLevel(logrus.InfoLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		debug = true
		formatter.PrettyPrint = true
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

	// -------------------- Hub for WebSocket -------------------- //

	hub := chat.NewHub()
	go hub.Run()

	// -------------------- Set up service -------------------- //

	svc, err := api.NewAPIService(hub, logrus.NewEntry(log), mongoDB, debug)
	if err != nil {
		log.Fatalf("error creating service instance: %s", err)
	}

	go svc.Serve()

	// -------------------- Listen for Interruption signal and shutdown -------------------- //

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(viper.GetInt("service.shutdown_timeout"))*time.Second,
	)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
