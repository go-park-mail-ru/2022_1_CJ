package db

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *core.User) error

	GetUserByID(ctx context.Context, ID string) (*core.User, error)
	GetUserByEmail(ctx context.Context, email string) (*core.User, error)

	CheckUserEmailExistence(ctx context.Context, email string) (bool, error)

	UpdateUser(ctx context.Context, user *core.User) error

	DeleteUser(ctx context.Context, user *core.User) error
	UserAddPost(ctx context.Context, user *core.User, postID string) error
}

type userRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

// CreateUser tries to insert given user to the db:
// returns error if the email is already taken, otherwise inserts.
func (repo *userRepositoryImpl) CreateUser(ctx context.Context, user *core.User) error {
	filter := bson.M{"email": user.Email}
	if err := repo.coll.FindOne(ctx, filter).Err(); err != mongo.ErrNoDocuments {
		if err == nil {
			return constants.ErrEmailAlreadyTaken
		} else {
			return err
		}
	}

	if err := repo.InitUser(user); err != nil {
		return err
	}

	_, err := repo.coll.InsertOne(ctx, user)
	return err
}

func (repo *userRepositoryImpl) GetUserByID(ctx context.Context, ID string) (*core.User, error) {
	user := new(core.User)
	filter := bson.M{"_id": ID}
	err := repo.coll.FindOne(ctx, filter).Decode(user)
	return user, wrapError(err)
}

// GetUserByEmail looks up in the db for user with the provided email.
func (repo *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*core.User, error) {
	user := new(core.User)
	filter := bson.M{"email": email}
	err := repo.coll.FindOne(ctx, filter).Decode(user)
	return user, wrapError(err)
}

// CheckUserEmailExistence checks whether user with given email exists. Returns true if email is already taken.
func (repo *userRepositoryImpl) CheckUserEmailExistence(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	err := repo.coll.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// UpdateUser updates user.
func (repo *userRepositoryImpl) UpdateUser(ctx context.Context, user *core.User) error {
	filter := bson.M{"_id": user.ID}
	_, err := repo.coll.ReplaceOne(ctx, filter, user)
	return err
}

// Add new user post
func (repo *userRepositoryImpl) UserAddPost(ctx context.Context, user *core.User, postID string) error {
	// TODO find in []string, if not exist than add
	//filter := bson.M{"posts": user.Posts}
	//err := repo.coll.FindOne(ctx, postID).Err()
	//if err != nil {
	//	if err == mongo.ErrNoDocuments {
	//		return false, nil
	//	}
	//	return false, err
	//}
	return nil
}

// DeleteUser deletes from the db user with the provided email
func (repo *userRepositoryImpl) DeleteUser(ctx context.Context, user *core.User) error {
	filter := bson.M{"email": user.Email}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

// NewUserRepository creates a new instance of userRepositoryImpl
func NewUserRepository(db *mongo.Database) (*userRepositoryImpl, error) {
	return &userRepositoryImpl{db: db, coll: db.Collection("users")}, nil
}

func (repo *userRepositoryImpl) InitUser(user *core.User) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	user.ID = uid
	user.CreatedAt = time.Now().Unix()
	return nil
}
