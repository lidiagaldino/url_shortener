package repository

import "url-shortener/internal/domain/entity"

type URLRepository interface {
	Save(url *entity.URL) error
	FindByID(id string) (*entity.URL, error)
	IncrementClick(id string) error
}
