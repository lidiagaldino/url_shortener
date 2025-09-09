package services

import (
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/repository"
	"url-shortener/pkg"
)

type URLService struct {
	repo        repository.URLRepository
	idGenerator pkg.IDGenerator
}

func NewURLService(repo repository.URLRepository, idGen pkg.IDGenerator) *URLService {
	return &URLService{
		repo:        repo,
		idGenerator: idGen,
	}
}

func (s *URLService) Shorten(originalURL string) (*entity.URL, error) {
	id, err := s.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	url := entity.URL{
		ID:          id,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Save(url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (s *URLService) Resolve(id string) (*entity.URL, error) {
	return s.repo.FindByID(id)
}
