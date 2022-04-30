package db

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"time"

	"github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommunityRepository interface {
	CreateCommunity(ctx context.Context, community *core.Community) (*core.Community, error)
	EditCommunity(ctx context.Context, community *core.Community) error
	GetCommunityByID(ctx context.Context, communityID string) (*core.Community, error)
	DeleteCommunity(ctx context.Context, communityID string) error
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
