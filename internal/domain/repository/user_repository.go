package repository

import "url-shortener/internal/domain/entity"

type UserRepository interface {
	Save(user *entity.User) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}
