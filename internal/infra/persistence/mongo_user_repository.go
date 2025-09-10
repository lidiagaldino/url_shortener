package persistence

import (
	"context"
	"log"
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/infra/persistence/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Save(user *entity.User) (*entity.User, error) {
	result, err := r.collection.InsertOne(context.TODO(), fromModelUser(user))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	savedUser, err := r.FindByID(result.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return savedUser, nil
}

func (r *MongoUserRepository) FindByID(id string) (*entity.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.collection.FindOne(context.TODO(), map[string]any{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return toEntityUser(&user), nil
}

func (r *MongoUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user model.User
	err := r.collection.FindOne(context.TODO(), map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return toEntityUser(&user), nil
}

func fromModelUser(user *entity.User) *model.User {
	return &model.User{
		ID:             primitive.NewObjectID(),
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		CreatedAt:      time.Now(),
	}
}

func toEntityUser(user *model.User) *entity.User {
	return &entity.User{
		ID:             user.ID.Hex(),
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
	}
}
