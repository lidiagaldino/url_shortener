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

func (r *MongoURLStatsRepository) FindByURLID(urlID string) ([]entity.URLStat, error) {
	ctx := context.TODO()

	filter := map[string]interface{}{"url_id": urlID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var stats []entity.URLStat
	for cursor.Next(ctx) {
		var m model.URLStat
		if err := cursor.Decode(&m); err != nil {
			return nil, err
		}

		stats = append(stats, entity.URLStat{
			ID:        m.ID.Hex(),
			URLID:     m.URLID,
			ClickedAt: m.ClickedAt,
			IP:        m.IP,
			UserAgent: m.UserAgent,
			Referer:   m.Referer,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return stats, nil
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
