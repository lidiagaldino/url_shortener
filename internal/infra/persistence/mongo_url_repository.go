package persistence

import (
	"context"
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/infra/persistence/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoURLRepository struct {
	collection *mongo.Collection
}

func NewMongoURLRepository(db *mongo.Database) *MongoURLRepository {
	return &MongoURLRepository{
		collection: db.Collection("urls"),
	}
}

func (r *MongoURLRepository) Save(url *entity.URL) error {
	_, err := r.collection.InsertOne(context.TODO(), fromModelUrl(url))
	return err
}

func (r *MongoURLRepository) FindByID(id string) (*entity.URL, error) {
	var url model.URL
	err := r.collection.FindOne(context.TODO(), map[string]string{"_id": id}).Decode(&url)
	if err != nil {
		return nil, err
	}
	return toEntityUrl(&url), nil
}

func fromModelUrl(url *entity.URL) *model.URL {
	return &model.URL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		CreatedAt:   time.Now(),
	}
}

func toEntityUrl(url *model.URL) *entity.URL {
	return &entity.URL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		CreatedAt:   url.CreatedAt,
	}
}
