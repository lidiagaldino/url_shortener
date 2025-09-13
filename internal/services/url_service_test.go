package services_test

import (
	"errors"
	"testing"
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockURLRepo struct {
	mock.Mock
}

func (m *MockURLRepo) Save(u *entity.URL) error {
	args := m.Called(u)
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

func (m *MockStatsRepo) FindByURLID(urlID string) ([]entity.URLStat, error) {
	args := m.Called(urlID)
	if args.Get(0) != nil {
		return args.Get(0).([]entity.URLStat), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockIDGen struct {
	mock.Mock
}

func (m *MockIDGen) Generate() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// ----------------- Shorten -----------------

func TestURLService_Shorten_Success(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	idGen.On("Generate").Return("abc123", nil)
	urlRepo.On("Save", mock.AnythingOfType("*entity.URL")).Return(nil)

	urlEntity, err := svc.Shorten("https://example.com", "owner1")

	assert.NoError(t, err)
	assert.NotNil(t, urlEntity)
	assert.Equal(t, "abc123", urlEntity.ID)
	assert.Equal(t, "owner1", urlEntity.OwnerID)
	assert.Equal(t, "https://example.com", urlEntity.OriginalURL)
}

func TestURLService_Shorten_InvalidURL(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	_, err := svc.Shorten("invalid-url", "owner1")
	assert.ErrorIs(t, err, exceptions.ErrInvalidURL)
}

func TestURLService_Shorten_PrivateIP(t *testing.T) {
	idGen := new(MockIDGen)
	repo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)

	svc := services.NewURLService(repo, idGen, statsRepo)

	privateURL := "http://192.168.0.1/test"
	ownerID := "user1"

	idGen.On("Generate").Return("some-id", nil)

	result, err := svc.Shorten(privateURL, ownerID)

	assert.ErrorIs(t, err, exceptions.ErrInvalidURL)
	assert.Nil(t, result)
}

func TestURLService_Shorten_GenerateError(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	idGen.On("Generate").Return("", errors.New("generate error"))

	_, err := svc.Shorten("https://example.com", "owner1")
	assert.Error(t, err)
	assert.EqualError(t, err, "generate error")
}

func TestURLService_Shorten_SaveError(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	idGen.On("Generate").Return("abc123", nil)
	urlRepo.On("Save", mock.AnythingOfType("*entity.URL")).Return(errors.New("save error"))

	_, err := svc.Shorten("https://example.com", "owner1")
	assert.Error(t, err)
	assert.EqualError(t, err, "save error")
}

// ----------------- Resolve -----------------

func TestURLService_Resolve_Success(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlEntity := &entity.URL{
		ID:          "abc123",
		OriginalURL: "https://example.com",
	}

	urlRepo.On("FindByID", "abc123").Return(urlEntity, nil)
	urlRepo.On("IncrementClick", "abc123").Return(nil)
	statsRepo.On("Save", mock.AnythingOfType("*entity.URLStat")).Return(nil)

	res, err := svc.Resolve("abc123", "1.2.3.4", "user-agent", "referer")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "https://example.com", res.OriginalURL)
}

func TestURLService_Resolve_NotFound(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlRepo.On("FindByID", "abc123").Return(nil, errors.New("not found"))

	res, err := svc.Resolve("abc123", "1.2.3.4", "user-agent", "referer")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "not found")
}

func TestURLService_Resolve_IncrementClickError(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlEntity := &entity.URL{ID: "abc123", OriginalURL: "https://example.com"}

	urlRepo.On("FindByID", "abc123").Return(urlEntity, nil)
	urlRepo.On("IncrementClick", "abc123").Return(errors.New("increment error"))

	res, err := svc.Resolve("abc123", "1.2.3.4", "user-agent", "referer")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "increment error")
}

func TestURLService_Resolve_StatsSaveError(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlEntity := &entity.URL{ID: "abc123", OriginalURL: "https://example.com"}

	urlRepo.On("FindByID", "abc123").Return(urlEntity, nil)
	urlRepo.On("IncrementClick", "abc123").Return(nil)
	statsRepo.On("Save", mock.AnythingOfType("*entity.URLStat")).Return(errors.New("stats save error"))

	res, err := svc.Resolve("abc123", "1.2.3.4", "user-agent", "referer")

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.EqualError(t, err, "stats save error")
}

// ----------------- Stats -----------------

// func TestURLService_Stats_Success(t *testing.T) {
// 	urlRepo := new(MockURLRepo)
// 	statsRepo := new(MockStatsRepo)
// 	idGen := new(MockIDGen)

// 	svc := services.NewURLService(urlRepo, idGen, statsRepo)

// 	urlEntity := &entity.URL{
// 		ID:         "url1",
// 		OwnerID:    "owner1",
// 		ClickCount: 5,
// 		LastClick:  time.Now(),
// 	}

// 	statsEntities := []entity.URLStat{
// 		{
// 			ID:        "stat1",
// 			URLID:     "url1",
// 			ClickedAt: time.Now(),
// 			IP:        "1.2.3.4",
// 			UserAgent: "agent1",
// 			Referer:   "ref1",
// 		},
// 	}

// 	urlRepo.On("FindByID", "url1").Return(urlEntity, nil)
// 	statsRepo.On("FindByURLID", "url1").Return(statsEntities, nil)

// 	result, err := svc.Stats("url1", "owner1")

// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(1, len(result.StatsData))
// 	assert.Equal("url1", result.StatsData[0].URLID)
// 	assert.Equal(6, result.StatsResume.Clicks) // ClickCount + 1
// }

func TestURLService_Stats_Unauthorized(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlEntity := &entity.URL{
		ID:      "url1",
		OwnerID: "owner1",
	}

	urlRepo.On("FindByID", "url1").Return(urlEntity, nil)

	result, err := svc.Stats("url1", "otherUser")

	assert.ErrorIs(t, err, exceptions.ErrUnauthorizedURLStatistics)
	assert.Nil(t, result)
}

func TestURLService_Stats_URLNotFound(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlRepo.On("FindByID", "url1").Return(nil, errors.New("not found"))

	result, err := svc.Stats("url1", "owner1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
}

func TestURLService_Stats_StatsRepoError(t *testing.T) {
	urlRepo := new(MockURLRepo)
	statsRepo := new(MockStatsRepo)
	idGen := new(MockIDGen)

	svc := services.NewURLService(urlRepo, idGen, statsRepo)

	urlEntity := &entity.URL{
		ID:      "url1",
		OwnerID: "owner1",
	}

	urlRepo.On("FindByID", "url1").Return(urlEntity, nil)
	statsRepo.On("FindByURLID", "url1").Return(nil, errors.New("stats error"))

	result, err := svc.Stats("url1", "owner1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "stats error")
}
