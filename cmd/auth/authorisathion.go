package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/api"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/controller"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/handler"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

var (
	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	info = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "grpcMetrics",
		Help: "help",
	}, []string{"name"})
)

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, info)
	info.WithLabelValues("Test")
}

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

	//-------------------- Set up metrics -------------------- //
	s := echo.New()
	s.Validator = api.NewValidator()
	s.Binder = api.NewBinder()

	s.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 5 * time.Minute}),
	)

	grpc_prometheus.EnableHandlingTimeHistogram()

	grpcMetrics.InitializeMetrics(server)

	log.Info("success init metrics: auth gRPC")

	listAdr := viper.GetString(configHost) + ":" + "9082"
	go func() {
		err := s.Start(listAdr)
		if err != nil {
			log.Fatalln("cant start metrics", err)
		}
	}()

	//-------------------- Set up service -------------------- //

	listenAddr := viper.GetString(configHost) + ":" + viper.GetString(configPost)
	lis, err := net.Listen(viper.GetString(configNetwork), listenAddr)
	if err != nil {
		log.Fatalln("cant listen port", err)
	}

	handler.RegisterUserAuthServer(server, controller.CreateAuthServer(mongoDB, logrus.NewEntry(log)))

	log.Infof("Server started listen at: %s", listenAddr)

	err = server.Serve(lis)
	if err != nil {
		return
	}
}
