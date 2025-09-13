package persistence

import (
	"context"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/repository"
	"url-shortener/internal/infra/persistence/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoURLStatsRepository struct {
	collection *mongo.Collection
}

func NewMongoURLStatsRepository(db *mongo.Database) repository.URLStatsRepository {
	return &MongoURLStatsRepository{
		collection: db.Collection("url_stats"),
	}
}

func (r *MongoURLStatsRepository) Save(stat *entity.URLStat) error {
	urlStat := fromModelUrlStats(stat)

	_, err := r.collection.InsertOne(context.TODO(), urlStat)
	return err
}

func fromModelUrlStats(url *entity.URLStat) *model.URLStat {
	return &model.URLStat{
		ID:        primitive.NewObjectID(),
		URLID:     url.URLID,
		ClickedAt: url.ClickedAt,
		IP:        url.IP,
		UserAgent: url.UserAgent,
		Referer:   url.Referer,
	}
}
