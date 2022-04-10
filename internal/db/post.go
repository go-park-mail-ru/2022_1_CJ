package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *core.Post) (*core.Post, error)

	GetPostByID(ctx context.Context, postID string) (*core.Post, error)
	GetPostsByUserID(ctx context.Context, userID string) ([]core.Post, error)

	EditPost(ctx context.Context, post *core.Post) (*core.Post, error)
	DeletePost(ctx context.Context, postID string) error

	GetFeed(ctx context.Context, userID string) ([]core.Post, error)
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

func (repo *postRepositoryImpl) GetPostByID(ctx context.Context, postID string) (*core.Post, error) {
	post := new(core.Post)
	filter := bson.M{"_id": postID}
	err := repo.coll.FindOne(ctx, filter).Decode(post)
	return post, wrapError(err)
}

// TODO: add pagination
func (repo *postRepositoryImpl) GetPostsByUserID(ctx context.Context, userID string) ([]core.Post, error) {
	var posts []core.Post
	filter := bson.M{"author_id": userID}
	cursor, err := repo.coll.Find(ctx, filter)
	if err != nil {
		return posts, err
	} else {
		err = cursor.All(ctx, &posts)
	}
	return posts, err
}

func (repo *postRepositoryImpl) EditPost(ctx context.Context, post *core.Post) (*core.Post, error) {
	filter := bson.M{"_id": post.ID}
	_, err := repo.coll.ReplaceOne(ctx, filter, post)
	return post, err
}

func (repo *postRepositoryImpl) DeletePost(ctx context.Context, postID string) error {
	filter := bson.M{"_id": postID}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

// TODO: refactor
func (repo *postRepositoryImpl) GetFeed(ctx context.Context, userID string) ([]core.Post, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := repo.coll.Find(ctx, bson.M{}, opts)

	var posts []core.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, err
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
	post.CreatedAt = time.Now().Unix()
	return nil
}
