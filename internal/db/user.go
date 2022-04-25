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

	UserAddPost(ctx context.Context, userID string, postID string) error
	UserCheckPost(ctx context.Context, user *core.User, postID string) error
	UserDeletePost(ctx context.Context, userID string, postID string) error

	SelectUsers(ctx context.Context, selector string) ([]core.User, error)

	AddDialog(ctx context.Context, dialogID string, userID string) error
	GetUserDialogs(ctx context.Context, userID string) ([]string, error)
	UserCheckDialog(ctx context.Context, dialogID string, userID string) error
}

type userRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

// NewUserRepository creates a new instance of userRepositoryImpl
func NewUserRepository(db *mongo.Database) (*userRepositoryImpl, error) {
	return &userRepositoryImpl{db: db, coll: db.Collection("users")}, nil
}

// NewUserRepositoryTest for Tests (bad)
func NewUserRepositoryTest(collection *mongo.Collection) (*userRepositoryImpl, error) {
	return &userRepositoryImpl{coll: collection}, nil
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

// UserAddPost Add new user post
func (repo *userRepositoryImpl) UserAddPost(ctx context.Context, userID string, postID string) error {
	if _, err := repo.coll.UpdateByID(ctx, userID, bson.M{"$push": bson.D{{Key: "posts", Value: postID}}}); err != nil {
		return err
	}
	return nil
}

//UserCheckPost Check existing post in posts by User
func (repo *userRepositoryImpl) UserCheckPost(ctx context.Context, user *core.User, postID string) error {
	filter := bson.M{"_id": user.ID, "posts": postID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	return nil
}

// UserDeletePost Add new user post
func (repo *userRepositoryImpl) UserDeletePost(ctx context.Context, userID string, postID string) error {
	filter := bson.M{"_id": userID, "posts": postID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	if _, err := repo.coll.UpdateByID(ctx, userID, bson.M{"$pull": bson.M{"posts": postID}}); err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes from the db user with the provided email
func (repo *userRepositoryImpl) DeleteUser(ctx context.Context, user *core.User) error {
	filter := bson.M{"email": user.Email}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

func (repo *userRepositoryImpl) SelectUsers(ctx context.Context, selector string) ([]core.User, error) {
	var users []core.User

	fuzzy := bson.M{"$regex": selector, "$options": "i"}
	filter := bson.M{"$or": []bson.M{
		{"name.first": fuzzy},
		{"name.last": fuzzy}},
	}

	cursor, err := repo.coll.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return users, nil
		}
		return nil, err
	} else {
		err = cursor.All(ctx, &users)
	}

	return users, err
}

func (repo *userRepositoryImpl) AddDialog(ctx context.Context, dialogID string, userID string) error {
	if _, err := repo.coll.UpdateByID(ctx, userID, bson.M{"$push": bson.D{{Key: "dialog_ids", Value: dialogID}}}); err != nil {
		return err
	}
	return nil
}

func (repo *userRepositoryImpl) UserCheckDialog(ctx context.Context, dialogID string, userID string) error {
	filter := bson.M{"_id": userID, "dialog_ids": dialogID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	return nil
}

func (repo *userRepositoryImpl) GetUserDialogs(ctx context.Context, userID string) ([]string, error) {
	user := new(core.User)
	filter := bson.M{"_id": userID}
	err := repo.coll.FindOne(ctx, filter).Decode(user)
	return user.DialogIDs, wrapError(err)
}

func (repo *userRepositoryImpl) InitUser(user *core.User) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	user.ID = uid
	user.Image = "default.jpeg"
	user.CreatedAt = time.Now().Unix()
	return nil
}
