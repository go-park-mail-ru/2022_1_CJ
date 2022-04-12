package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/db"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

const (
	database          = "cj_test"
	connection_string = "mongodb://root:rootpassword@mongodb_test:27017/"
)

func TestMain(m *testing.M) {
	// start test db

	os.Exit(m.Run())
}

func TestConnectToDB(t *testing.T) (*db.Repository, error) {
	t.Helper()
	clientOptions := options.Client().ApplyURI(connection_string)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Info("connected to MongoDB")
	mongoDB := client.Database(database)
	return db.NewRepository(mongoDB)
}

func TestUser(t *testing.T) *core.User {
	t.Helper()

	return &core.User{
		Name:  common.UserName{First: "User", Last: "Userov"},
		Email: "user@example.org",
		Image: "src/img.jpg",
	}
}
