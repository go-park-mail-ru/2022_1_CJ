package db

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/constants"
	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/common"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommunityRepository interface {
	CreateCommunity(ctx context.Context, community *core.Community) (*core.Community, error)
	EditCommunity(ctx context.Context, community *core.Community) error
	GetCommunityByID(ctx context.Context, communityID string) (*core.Community, error)
	DeleteCommunity(ctx context.Context, communityID string) error

	SearchCommunities(ctx context.Context, selector string, limit, pageNumber int64) ([]core.Community, *common.PageResponse, error)
	GetAllCommunities(ctx context.Context, limit, pageNumber int64) ([]core.Community, *common.PageResponse, error)

	AddFollower(ctx context.Context, communityID string, userID string) error
	DeleteFollower(ctx context.Context, communityID string, userID string) error
	DeleteAdmin(ctx context.Context, communityID string, userID string) error

	CommunityAddPost(ctx context.Context, communityID string, postID string) error
	CommunityDeletePost(ctx context.Context, communityID string, postID string) error
}

type comunnityRepositoryImpl struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewCommunityRepository(db *mongo.Database) (*comunnityRepositoryImpl, error) {
	return &comunnityRepositoryImpl{db: db, coll: db.Collection("community")}, nil
}

// NewCommunityRepositoryTest for Tests (bad)
func NewCommunityRepositoryTest(collection *mongo.Collection) (*comunnityRepositoryImpl, error) {
	return &comunnityRepositoryImpl{coll: collection}, nil
}

// CommunityAddPost Add new community post
func (repo *comunnityRepositoryImpl) CommunityAddPost(ctx context.Context, communityID string, postID string) error {
	if _, err := repo.coll.UpdateByID(ctx, communityID, bson.M{"$push": bson.D{{Key: "posts", Value: postID}}}); err != nil {
		return err
	}
	return nil
}

// CommunityDeletePost ...
func (repo *comunnityRepositoryImpl) CommunityDeletePost(ctx context.Context, communityID string, postID string) error {
	filter := bson.M{"_id": communityID, "posts": postID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	if _, err := repo.coll.UpdateByID(ctx, communityID, bson.M{"$pull": bson.M{"posts": postID}}); err != nil {
		return err
	}
	return nil
}

func (repo *comunnityRepositoryImpl) CreateCommunity(ctx context.Context, community *core.Community) (*core.Community, error) {
	if err := repo.InitCommunity(community); err != nil {
		return nil, err
	}
	_, err := repo.coll.InsertOne(ctx, community)

	comunnitySanitize(community)

	return community, err
}

func (repo *comunnityRepositoryImpl) EditCommunity(ctx context.Context, community *core.Community) error {
	filter := bson.M{"_id": community.ID}
	_, err := repo.coll.ReplaceOne(ctx, filter, community)
	return err
}

func (repo *comunnityRepositoryImpl) AddFollower(ctx context.Context, communityID string, userID string) error {
	if _, err := repo.coll.UpdateByID(ctx, communityID, bson.M{"$push": bson.D{{Key: "followers", Value: userID}}}); err != nil {
		return err
	}
	return nil
}
func (repo *comunnityRepositoryImpl) DeleteFollower(ctx context.Context, communityID string, userID string) error {
	filter := bson.M{"_id": communityID, "followers": userID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	if _, err := repo.coll.UpdateByID(ctx, communityID, bson.M{"$pull": bson.M{"followers": userID}}); err != nil {
		return err
	}
	return nil
}
func (repo *comunnityRepositoryImpl) DeleteAdmin(ctx context.Context, communityID string, userID string) error {
	filter := bson.M{"_id": communityID, "followers": userID}
	if err := repo.coll.FindOne(ctx, filter).Err(); err == mongo.ErrNoDocuments {
		return constants.ErrDBNotFound
	}
	if _, err := repo.coll.UpdateByID(ctx, communityID, bson.M{"$pull": bson.M{"admins": userID}}); err != nil {
		return err
	}
	return nil
}

func (repo *comunnityRepositoryImpl) DeleteCommunity(ctx context.Context, communityID string) error {
	filter := bson.M{"_id": communityID}
	_, err := repo.coll.DeleteOne(ctx, filter)
	return err
}

func (repo *comunnityRepositoryImpl) GetCommunityByID(ctx context.Context, communityID string) (*core.Community, error) {
	community := new(core.Community)
	filter := bson.M{"_id": communityID}
	err := repo.coll.FindOne(ctx, filter).Decode(community)

	comunnitySanitize(community)
	return community, wrapError(err)
}

func (repo *comunnityRepositoryImpl) GetAllCommunities(ctx context.Context, limit, pageNumber int64) ([]core.Community, *common.PageResponse, error) {
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit != -1 {
		opts.SetSkip((pageNumber - 1) * limit)
		opts.SetLimit(limit)
	}
	cursor, err := repo.coll.Find(ctx, bson.M{}, opts)
	var communities []core.Community
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return communities, &common.PageResponse{}, nil
		}
		return nil, nil, err
	} else {
		err = cursor.All(ctx, &communities)
	}
	total, _ := repo.coll.CountDocuments(ctx, bson.M{})

	isLarge := func(res bool) int64 {
		if res {
			return 1
		} else {
			return 0
		}
	}
	res := &common.PageResponse{
		Total:       total,
		AmountPages: total/limit + isLarge(total%limit > 0),
	}
	if limit == -1 {
		res.AmountPages = 1
	}
	for _, comm := range communities {
		comunnitySanitize(&comm)
	}
	return communities, res, err
}

func (repo *comunnityRepositoryImpl) SearchCommunities(ctx context.Context, selector string, limit, pageNumber int64) ([]core.Community, *common.PageResponse, error) {
	fuzzy := bson.M{"$regex": selector, "$options": "i"}
	filter := bson.D{{Key: "name", Value: fuzzy}}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	if limit != -1 {
		opts.SetSkip((pageNumber - 1) * limit)
		opts.SetLimit(limit)
	}
	cursor, err := repo.coll.Find(ctx, filter, opts)
	var communities []core.Community
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return communities, &common.PageResponse{}, nil
		}
		return nil, nil, err
	} else {
		err = cursor.All(ctx, &communities)
	}
	total, _ := repo.coll.CountDocuments(ctx, filter)

	isLarge := func(res bool) int64 {
		if res {
			return 1
		} else {
			return 0
		}
	}
	res := &common.PageResponse{
		Total:       total,
		AmountPages: total/limit + isLarge(total%limit > 0),
	}
	if limit == -1 {
		res.AmountPages = 1
	}
	for _, comm := range communities {
		comunnitySanitize(&comm)
	}
	return communities, res, err
}

func (repo *comunnityRepositoryImpl) InitCommunity(community *core.Community) error {
	uid, err := core.GenUUID()
	if err != nil {
		return err
	}
	community.ID = uid
	community.CreatedAt = time.Now().Unix()
	community.PostIDs = []string{}
	return nil
}

// Help func for defense from XSS attacks
func comunnitySanitize(community *core.Community) {
	p := bluemonday.UGCPolicy()
	community.Name = p.Sanitize(community.Name)
	community.Image = p.Sanitize(community.Image)
	community.Info = p.Sanitize(community.Info)
}
