package services

import (
	"errors"
	"testing"
	"time"

	"url-shortener/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// Mocks
//

type MockURLRepo struct {
	mock.Mock
}

func (m *MockURLRepo) Save(url *entity.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepo) FindByID(id string) (*entity.URL, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.URL), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockURLRepo) IncrementClick(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockStatsRepo struct {
	mock.Mock
}

func (m *MockStatsRepo) Save(stat *entity.URLStat) error {
	args := m.Called(stat)
	return args.Error(0)
}

type MockIDGenerator struct {
	mock.Mock
}

func (m *MockIDGenerator) Generate() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

//
// Tests
//

func TestShorten_Success(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	urlInput := "https://example.com"
	ownerID := "123"
	shortID := "abc123"
	now := time.Now()

	idGen.On("Generate").Return(shortID, nil)
	repo.On("Save", mock.MatchedBy(func(u *entity.URL) bool {
		return u.OriginalURL == urlInput && u.ID == shortID && u.OwnerID == ownerID
	})).Return(nil)

	url, err := service.Shorten(urlInput, ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, shortID, url.ID)
	assert.Equal(t, urlInput, url.OriginalURL)
	assert.WithinDuration(t, now, url.CreatedAt, time.Second)
}

func TestShorten_ErrorOnIDGen(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	idGen.On("Generate").Return("", errors.New("generate id failed"))

	_, err := service.Shorten("https://example.com", "123")
	assert.Error(t, err)
}

func TestShorten_ErrorOnSave(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	shortID := "abc123"
	idGen.On("Generate").Return(shortID, nil)
	repo.On("Save", mock.AnythingOfType("*entity.URL")).Return(errors.New("failed to save"))

	_, err := service.Shorten("https://example.com", "123")
	assert.Error(t, err)
}

func TestResolve_Success(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	urlID := "abc123"
	originalURL := "https://example.com"
	ip := "127.0.0.1"
	userAgent := "Go-http-client/1.1"
	referer := "https://google.com"

	repo.On("FindByID", urlID).Return(&entity.URL{ID: urlID, OriginalURL: originalURL}, nil)
	repo.On("IncrementClick", urlID).Return(nil)
	statsRepo.On("Save", mock.AnythingOfType("*entity.URLStat")).Return(nil)

	url, err := service.Resolve(urlID, ip, userAgent, referer)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, urlID, url.ID)
	assert.Equal(t, originalURL, url.OriginalURL)
}

func TestResolve_NotFound(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	urlID := "abc123"
	repo.On("FindByID", urlID).Return(nil, errors.New("not found"))

	_, err := service.Resolve(urlID, "127.0.0.1", "UA", "ref")
	assert.Error(t, err)
}

func TestResolve_ErrorOnIncrementClick(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	urlID := "abc123"
	repo.On("FindByID", urlID).Return(&entity.URL{ID: urlID, OriginalURL: "https://example.com"}, nil)
	repo.On("IncrementClick", urlID).Return(errors.New("failed to increment"))

	_, err := service.Resolve(urlID, "127.0.0.1", "UA", "ref")
	assert.Error(t, err)
}

func TestResolve_ErrorOnSaveStats(t *testing.T) {
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen, statsRepo)

	urlID := "abc123"
	repo.On("FindByID", urlID).Return(&entity.URL{ID: urlID, OriginalURL: "https://example.com"}, nil)
	repo.On("IncrementClick", urlID).Return(nil)
	statsRepo.On("Save", mock.AnythingOfType("*entity.URLStat")).Return(errors.New("failed to save stats"))

	_, err := service.Resolve(urlID, "127.0.0.1", "UA", "ref")
	assert.Error(t, err)
}
