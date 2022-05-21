package db

import (
	"context"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/utils"
	"github.com/microcosm-cc/bluemonday"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *core.Post) (*core.Post, error)

	GetPostByID(ctx context.Context, postID string) (*core.Post, error)
	GetPostsByUserID(ctx context.Context, userID string, pageNumber int64, limit int64) ([]core.Post, *common.PageResponse, error)

	EditPost(ctx context.Context, post *core.Post) (*core.Post, error)
	DeletePost(ctx context.Context, postID string) error

	GetFeed(ctx context.Context, userID string, pageNumber int64, limit int64) ([]core.Post, *common.PageResponse, error)

	PostAddComment(ctx context.Context, postID string, commentID string) error
	PostCheckComment(ctx context.Context, post *core.Post, commentID string) error
	PostDeleteComment(ctx context.Context, postID string, commentID string) error
}

type postRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewPostRepository(db *mongo.Database) (*postRepositoryImpl, error) {
	return &postRepositoryImpl{db: db, coll: db.Collection("posts")}, nil
}

// NewUserRepositoryTest for Tests (bad)
func NewPostRepositoryTest(collection *mongo.Collection) (*postRepositoryImpl, error) {
	return &postRepositoryImpl{coll: collection}, nil
}

func (repo *postRepositoryImpl) CreatePost(ctx context.Context, post *core.Post) (*core.Post, error) {
	if err := repo.InitPost(post); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, post)

	// Sanitize
	p := bluemonday.UGCPolicy()
	post.Message = p.Sanitize(post.Message)

	return post, err
}

// PostAddComment Add new comment
func (repo *postRepositoryImpl) PostAddComment(ctx context.Context, postID string, commentID string) error {
	if _, err := repo.coll.UpdateByID(ctx, postID, bson.M{"$push": bson.D{{Key: "comment_ids", Value: commentID}}}); err != nil {
		return err
	}
	return nil
}

// PostDeleteComment Delete comment
func (repo *postRepositoryImpl) PostDeleteComment(ctx context.Context, postID string, commentID string) error {
	filter := bson.M{"_id": postID, "comment_ids": commentID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	if _, err := repo.coll.UpdateByID(ctx, postID, bson.M{"$pull": bson.M{"comment_ids": commentID}}); err != nil {
		return err
	}
	return nil
}

//PostCheckComment Check existing comment in post
func (repo *postRepositoryImpl) PostCheckComment(ctx context.Context, post *core.Post, commentID string) error {
	filter := bson.M{"_id": post.ID, "comment_ids": commentID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	return nil
}

func (repo *postRepositoryImpl) GetPostByID(ctx context.Context, postID string) (*core.Post, error) {
	post := new(core.Post)
	filter := bson.M{"_id": postID}
	err := repo.coll.FindOne(ctx, filter).Decode(post)
	return post, wrapError(err)
}

func (repo *postRepositoryImpl) GetPostsByUserID(ctx context.Context, userID string, pageNumber int64, limit int64) ([]core.Post, *common.PageResponse, error) {
	var posts []core.Post
	filter := bson.M{"author_id": userID}

	findOptions := options.Find()

	findOptions.SetSort(bson.D{{"created_at", -1}})

	if limit != -1 {
		findOptions.SetSkip((pageNumber - 1) * limit)
		findOptions.SetLimit(limit)
	}

	cursor, err := repo.coll.Find(ctx, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return posts, &common.PageResponse{}, nil
		}
		return nil, nil, err
	} else {
		err = cursor.All(ctx, &posts)
	}

	total, _ := repo.coll.CountDocuments(ctx, filter)
	res := &common.PageResponse{
		Total:       total,
		AmountPages: total/limit + utils.IsLarge(total%limit > 0),
	}
	if limit == -1 {
		res.AmountPages = 1
	}

	// Sanitize
	p := bluemonday.UGCPolicy()
	for i, _ := range posts {
		posts[i].Message = p.Sanitize(posts[i].Message)
	}

	return posts, res, err
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

func (repo *postRepositoryImpl) GetFeed(ctx context.Context, userID string, pageNumber int64, limit int64) ([]core.Post, *common.PageResponse, error) {
	filter := bson.M{}
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	if limit != -1 {
		opts.SetSkip((pageNumber - 1) * limit)
		opts.SetLimit(limit)
	}

	cursor, err := repo.coll.Find(ctx, filter, opts)

	var posts []core.Post
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return posts, &common.PageResponse{}, nil
		}
		return nil, nil, err
	} else {
		err = cursor.All(ctx, &posts)
	}

	total, _ := repo.coll.CountDocuments(ctx, filter)
	res := &common.PageResponse{
		Total:       total,
		AmountPages: total/limit + utils.IsLarge(total%limit > 0),
	}
	if limit == -1 {
		res.AmountPages = 1
	}

	// Sanitize
	p := bluemonday.UGCPolicy()
	for i, _ := range posts {
		posts[i].Message = p.Sanitize(posts[i].Message)
	}

	return posts, res, err
}

func (repo *postRepositoryImpl) InitPost(post *core.Post) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	post.ID = uid
	post.CreatedAt = time.Now().Unix()
	post.CommentsIDs = []string{}
	return nil
}
