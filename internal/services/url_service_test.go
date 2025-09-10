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
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen)

	urlInput := "https://example.com"
	OwnerID := "123"
	shortID := "abc123"
	now := time.Now()

	idGen.On("Generate").Return(shortID, nil)
	repo.On("Save", mock.MatchedBy(func(u *entity.URL) bool {
		return u.OriginalURL == urlInput && u.ID == shortID && u.OwnerID == OwnerID
	})).Return(nil)

	url, err := service.Shorten(urlInput, OwnerID)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, shortID, url.ID)
	assert.Equal(t, urlInput, url.OriginalURL)
	assert.WithinDuration(t, now, url.CreatedAt, time.Second)
}

func TestShorten_ErrorOnIDGen(t *testing.T) {
	repo := new(MockURLRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen)

	idGen.On("Generate").Return("", errors.New("generate id failed"))

	_, err := service.Shorten("https://example.com", "123")
	assert.Error(t, err)
}

func TestShorten_ErrorOnSave(t *testing.T) {
	repo := new(MockURLRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen)

	shortID := "abc123"
	idGen.On("Generate").Return(shortID, nil)
	repo.On("Save", mock.AnythingOfType("*entity.URL")).Return(errors.New("failed to save"))

	_, err := service.Shorten("https://example.com", "123")
	assert.Error(t, err)
}

func TestResolve_Success(t *testing.T) {
	repo := new(MockURLRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen)

	urlID := "abc123"
	originalURL := "https://example.com"
	repo.On("FindByID", urlID).Return(&entity.URL{ID: urlID, OriginalURL: originalURL}, nil)

	url, err := service.Resolve(urlID)
	assert.NoError(t, err)
	assert.NotNil(t, url)
	assert.Equal(t, urlID, url.ID)
	assert.Equal(t, originalURL, url.OriginalURL)
}

func TestResolve_NotFound(t *testing.T) {
	repo := new(MockURLRepo)
	idGen := new(MockIDGenerator)
	service := NewURLService(repo, idGen)

	urlID := "abc123"
	repo.On("FindByID", urlID).Return(nil, errors.New("not found"))

	_, err := service.Resolve(urlID)
	assert.Error(t, err)
}
