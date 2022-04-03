package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *core.Post) (*core.Post, error)
	EditPost(ctx context.Context, post *core.Post) (*core.Post, error)
	GetPostByID(ctx context.Context, ID string) (*core.Post, error)
	DeletePost(ctx context.Context, post *core.Post) error
}

type postRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *postRepositoryImpl) GetPostByID(ctx context.Context, ID string) (*core.Post, error) {
	post := new(core.Post)
	filter := bson.M{"_id": ID}
	err := repo.coll.FindOne(ctx, filter).Decode(post)
	return post, wrapError(err)
}

func (repo *postRepositoryImpl) CreatePost(ctx context.Context, post *core.Post) (*core.Post, error) {
	if err := repo.InitPost(post); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, post)
	return post, err
}

func (repo *postRepositoryImpl) EditPost(ctx context.Context, post *core.Post) (*core.Post, error) {
	filter := bson.M{"_id": post.ID}
	_, err := repo.coll.ReplaceOne(ctx, filter, post)
	return post, err
}

func (repo *postRepositoryImpl) DeletePost(ctx context.Context, post *core.Post) error {
	filter := bson.M{"_id": post.ID}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

func NewPostRepository(db *mongo.Database) (*postRepositoryImpl, error) {
	return &postRepositoryImpl{db: db, coll: db.Collection("posts")}, nil
}

func (repo *postRepositoryImpl) InitPost(post *core.Post) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	post.ID = uid
	post.AuthorID = post.AuthorID
	post.CreatedAt = time.Now().Unix()
	return nil
}
