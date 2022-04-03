package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *core.Post) (*core.Post, error)
	EditPost(ctx context.Context, post *core.Post) (*core.Post, error)
	GetPostByID(ctx context.Context, ID string) (*core.Post, error)
	DeletePost(ctx context.Context, post *core.Post) error
	GetFeed(ctx context.Context, UserID string) ([]string, error)
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

func (repo *postRepositoryImpl) GetFeed(ctx context.Context, UserID string) ([]string, error) {
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := repo.coll.Find(ctx, filter, opts)
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	var res []string

	for _, result := range results {
		res = append(res, result["created_at"].(string))
	}
	return res, err
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
