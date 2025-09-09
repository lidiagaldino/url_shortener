package infra

import (
	"context"
	"url-shortener/internal/domain/entity"

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

func (r *MongoURLRepository) Save(url entity.URL) error {
	_, err := r.collection.InsertOne(context.TODO(), url)
	return err
}

func (r *MongoURLRepository) FindByID(id string) (*entity.URL, error) {
	var url entity.URL
	err := r.collection.FindOne(context.TODO(), map[string]string{"_id": id}).Decode(&url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}
