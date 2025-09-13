package services

import (
	"net"
	"net/url"
	"strings"
	"time"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/domain/repository"
	"url-shortener/internal/services/dto"
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

func (s *URLService) isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate()
}

func (s *URLService) validateDestination(raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return exceptions.ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return exceptions.ErrInvalidURL
	}
	host := u.Hostname()
	if host == "" {
		return exceptions.ErrInvalidURL
	}
	lower := strings.ToLower(host)
	if lower == "localhost" {
		return exceptions.ErrInvalidURL
	}
	if ip := net.ParseIP(host); ip != nil {
		if s.isPrivateIP(ip) {
			return exceptions.ErrInvalidURL
		}
	}

	return nil
}

func (s *URLService) Shorten(originalURL, ownerID string) (*entity.URL, error) {
	if err := s.validateDestination(originalURL); err != nil {
		return nil, exceptions.ErrInvalidURL
	}

	id, err := s.idGenerator.Generate()
	if err != nil {
		return nil, err
	}

	urlEntity := entity.URL{
		ID:          id,
		OriginalURL: originalURL,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Save(&urlEntity); err != nil {
		return nil, err
	}

	return &urlEntity, nil
}

func (s *URLService) Resolve(id, ip, userAgent, referer string) (*entity.URL, error) {
	url, err := s.repo.FindByID(id)
	if err != nil {
		return nil, exceptions.ErrURLNotFound
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

func (s *URLService) Stats(id, ownerID string) (*dto.URLStats, error) {
	url, err := s.repo.FindByID(id)
	if err != nil {
		return nil, exceptions.ErrURLNotFound
	}

	if url.OwnerID != ownerID {
		return nil, exceptions.ErrUnauthorizedURLStatistics
	}

	stats, err := s.statsRepo.FindByURLID(id)
	if err != nil {
		return nil, exceptions.ErrURLNotFound
	}

	var statsData []dto.Data

	for _, sEnt := range stats {
		statsData = append(statsData, dto.Data{
			ID:        sEnt.ID,
			URLID:     sEnt.URLID,
			ClickedAt: sEnt.ClickedAt,
			IP:        sEnt.IP,
			UserAgent: sEnt.UserAgent,
			Referer:   sEnt.Referer,
		})
	}

	resume := dto.Resume{
		Clicks:    url.ClickCount,
		LastClick: url.LastClick,
	}

	result := &dto.URLStats{
		StatsResume: resume,
		StatsData:   statsData,
	}

	return result, nil
}
