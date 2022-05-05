package auth_db

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/mircoservices/auth-microservice/model/core"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *auth_core.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*auth_core.User, error)
	CheckUserEmailExistence(ctx context.Context, email string) (bool, error)
}

type authRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) (*authRepositoryImpl, error) {
	return &authRepositoryImpl{db: db, coll: db.Collection("auth")}, nil
}

// CreateUser tries to insert given user to the db:
// returns error if the email is already taken, otherwise inserts.
func (repo *authRepositoryImpl) CreateUser(ctx context.Context, user *auth_core.User) (string, error) {
	filter := bson.M{"email": user.Email}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return auth_constants.Nothing, status.Error(codes.Internal, auth_constants.ErrEmailAlreadyTaken.Error())
		} else {
			return auth_constants.Nothing, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
		}
	}

	if err := repo.InitUser(user); err != nil {
		return auth_constants.Nothing, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}

	_, err := repo.coll.InsertOne(ctx, user)
	if err != nil {
		return user.ID, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}
	return user.ID, err
}

func (repo *authRepositoryImpl) GetUserByID(ctx context.Context, ID string) (*auth_core.User, error) {
	user := new(auth_core.User)
	filter := bson.M{"_id": ID}
	err := repo.coll.FindOne(ctx, filter).Decode(user)

	// Sanitize
	userSanitize(user)
	if err != nil {
		return user, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}
	return user, err
}

// GetUserByEmail looks up in the db for user with the provided email.
func (repo *authRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*auth_core.User, error) {
	user := new(auth_core.User)
	filter := bson.M{"email": email}
	err := repo.coll.FindOne(ctx, filter).Decode(user)

	// Sanitize
	userSanitize(user)

	if err != nil {
		return user, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}
	return user, err
}

// CheckUserEmailExistence checks whether user with given email exists. Returns true if email is already taken.
func (repo *authRepositoryImpl) CheckUserEmailExistence(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	err := repo.coll.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, status.Error(codes.Internal, fmt.Errorf("%s", err).Error())
	}

	return true, nil
}

// Help func for defense from XSS attacks
func userSanitize(user *auth_core.User) {
	p := bluemonday.UGCPolicy()
	user.Email = p.Sanitize(user.Email)
}

func (repo *authRepositoryImpl) InitUser(user *auth_core.User) error {
	uid, err := auth_core.GenUUID()
	if err != nil {
		return err
	}
	user.ID = uid
	user.CreatedAt = time.Now().Unix()
	return nil
}
