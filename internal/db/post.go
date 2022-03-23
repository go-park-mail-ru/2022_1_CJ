package db

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *core.Post) (*core.Post, error)
}

type postRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (repo *postRepositoryImpl) CreatePost(ctx context.Context, post *core.Post) (*core.Post, error) {
	if err := repo.InitPost(post); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, post)
	return post, err
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
