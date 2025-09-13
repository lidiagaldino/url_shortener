package persistence

import (
	"context"
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/infra/persistence/model"

	"go.mongodb.org/mongo-driver/bson"
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

func (r *MongoURLRepository) IncrementClick(id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{"click_count": 1},
		"$set": bson.M{"last_click": time.Now()},
	}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func fromModelUrl(url *entity.URL) *model.URL {
	return &model.URL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		OwnerID:     url.OwnerID,
		CreatedAt:   time.Now(),
	}
}

func toEntityUrl(url *model.URL) *entity.URL {
	return &entity.URL{
		ID:          url.ID,
		OriginalURL: url.OriginalURL,
		OwnerID:     url.OwnerID,
		CreatedAt:   url.CreatedAt,
		ClickCount:  url.ClickCount,
		LastClick:   url.LastClick,
	}
}
