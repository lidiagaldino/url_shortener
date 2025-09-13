package services

import (
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/repository"
	"url-shortener/pkg"
)

type URLService struct {
	repo        repository.URLRepository
	statsRepo   repository.URLStatsRepository
	idGenerator pkg.IDGenerator
}

func NewURLService(repo repository.URLRepository, idGen pkg.IDGenerator, statsRepo repository.URLStatsRepository) *URLService {
	return &URLService{
		repo:        repo,
		statsRepo:   statsRepo,
		idGenerator: idGen,
	}
}

func (s *URLService) Shorten(originalURL, ownerID string) (*entity.URL, error) {
	id, err := s.idGenerator.Generate()

	if err != nil {
		return nil, err
	}

	url := entity.URL{
		ID:          id,
		OriginalURL: originalURL,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Save(&url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (s *URLService) Resolve(id, ip, userAgent, referer string) (*entity.URL, error) {
	url, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.IncrementClick(id); err != nil {
		return nil, err
	}

	stat := &entity.URLStat{
		URLID:     id,
		ClickedAt: time.Now(),
		IP:        ip,
		UserAgent: userAgent,
		Referer:   referer,
	}
	if err := s.statsRepo.Save(stat); err != nil {
		return nil, err
	}

	return url, nil
}
