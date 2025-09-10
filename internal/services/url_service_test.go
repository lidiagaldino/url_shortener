package services

import (
	"errors"
	"testing"
	"time"

	"url-shortener/internal/domain/entity"
)

// --- Mocks ---

type mockRepo struct {
	saveFunc     func(url *entity.URL) error
	findByIDFunc func(id string) (*entity.URL, error)
}

func (m *mockRepo) Save(url *entity.URL) error {
	return m.saveFunc(url)
}

func (m *mockRepo) FindByID(id string) (*entity.URL, error) {
	return m.findByIDFunc(id)
}

type mockIDGen struct {
	generateFunc func() (string, error)
}

func (m *mockIDGen) Generate() (string, error) {
	return m.generateFunc()
}

// --- Testes ---

func TestShorten_Success(t *testing.T) {
	mockRepo := &mockRepo{
		saveFunc: func(url *entity.URL) error { return nil },
	}
	mockIDGen := &mockIDGen{
		generateFunc: func() (string, error) { return "abc123", nil },
	}

	service := NewURLService(mockRepo, mockIDGen)

	url, err := service.Shorten("https://example.com", "123")
	if err != nil {
		t.Fatalf("expected success, but it failed: %v", err)
	}

	if url.ID != "abc123" {
		t.Errorf("expcted ID to be 'abc123', but was '%s'", url.ID)
	}
	if url.OriginalURL != "https://example.com" {
		t.Errorf("expected origial URL to be 'https://example.com', but was '%s'", url.OriginalURL)
	}
	if time.Since(url.CreatedAt) > time.Second {
		t.Errorf("createdAt is invalid")
	}
}

func TestShorten_ErrorOnIDGen(t *testing.T) {
	mockRepo := &mockRepo{}
	mockIDGen := &mockIDGen{
		generateFunc: func() (string, error) { return "", errors.New("generate id failed") },
	}

	service := NewURLService(mockRepo, mockIDGen)

	_, err := service.Shorten("https://example.com", "123")
	if err == nil {
		t.Fatal("expected error, but was nil")
	}
}

func TestShorten_ErrorOnSave(t *testing.T) {
	mockRepo := &mockRepo{
		saveFunc: func(url *entity.URL) error { return errors.New("failed to save") },
	}
	mockIDGen := &mockIDGen{
		generateFunc: func() (string, error) { return "abc123", nil },
	}

	service := NewURLService(mockRepo, mockIDGen)

	_, err := service.Shorten("https://example.com", "123")
	if err == nil {
		t.Fatal("expected error, but was nil")
	}
}

func TestResolve_Success(t *testing.T) {
	mockRepo := &mockRepo{
		findByIDFunc: func(id string) (*entity.URL, error) {
			return &entity.URL{ID: id, OriginalURL: "https://example.com"}, nil
		},
	}
	mockIDGen := &mockIDGen{}

	service := NewURLService(mockRepo, mockIDGen)

	url, err := service.Resolve("abc123")
	if err != nil {
		t.Fatalf("expected success, but it failed: %v", err)
	}

	if url.ID != "abc123" {
		t.Errorf("expected ID to be 'abc123', but was '%s'", url.ID)
	}
}

func TestResolve_NotFound(t *testing.T) {
	mockRepo := &mockRepo{
		findByIDFunc: func(id string) (*entity.URL, error) {
			return nil, errors.New("not found")
		},
	}
	mockIDGen := &mockIDGen{}

	service := NewURLService(mockRepo, mockIDGen)

	_, err := service.Resolve("abc123")
	if err == nil {
		t.Fatal("expected not found, but was nil")
	}
}
