package services_test

import (
	"errors"
	"testing"
	"time"

	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/services"
	"url-shortener/internal/services/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Save(u *entity.User) (*entity.User, error) {
	args := m.Called(u)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) FindByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepo) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.User), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(pw string) (string, error) {
	args := m.Called(pw)
	return args.String(0), args.Error(1)
}

func (m *MockHasher) CheckPasswordHash(pw, hash string) bool {
	args := m.Called(pw, hash)
	return args.Bool(0)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// --- Tests ---

func TestUserService_Save_Success(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.UserInput{Name: "John", Email: "john@example.com", Password: "123"}
	hashed := "hashed123"
	savedUser := &entity.User{
		ID:             "1",
		Name:           "John",
		Email:          "john@example.com",
		HashedPassword: hashed,
		CreatedAt:      time.Now(),
	}

	repo.On("FindByEmail", input.Email).Return(nil, nil)
	hasher.On("HashPassword", input.Password).Return(hashed, nil)
	repo.On("Save", mock.AnythingOfType("*entity.User")).Return(savedUser, nil)

	result, err := svc.Save(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, savedUser.ID, result.ID)
	assert.Equal(t, savedUser.Email, result.Email)
}

func TestUserService_Save_EmailExists(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.UserInput{Name: "John", Email: "john@example.com", Password: "123"}
	existingUser := &entity.User{ID: "1", Email: input.Email}

	repo.On("FindByEmail", input.Email).Return(existingUser, nil)

	result, err := svc.Save(input)

	assert.ErrorIs(t, err, exceptions.ErrEmailAlreadyExists)
	assert.Nil(t, result)
}

func TestUserService_Save_HashError(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.UserInput{Name: "John", Email: "john@example.com", Password: "123"}
	repo.On("FindByEmail", input.Email).Return(nil, nil)
	hasher.On("HashPassword", input.Password).Return("", errors.New("hash error"))

	result, err := svc.Save(input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "hash error")
}

func TestUserService_Save_SaveError(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.UserInput{Name: "John", Email: "john@example.com", Password: "123"}
	hashed := "hashed123"

	repo.On("FindByEmail", input.Email).Return(nil, nil)
	hasher.On("HashPassword", input.Password).Return(hashed, nil)
	repo.On("Save", mock.AnythingOfType("*entity.User")).Return(nil, errors.New("save error"))

	result, err := svc.Save(input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "save error")
}

func TestUserService_LoginUser_Success(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.LoginUserInput{Email: "john@example.com", Password: "123"}
	user := &entity.User{ID: "1", Email: input.Email, HashedPassword: "hashed123"}

	repo.On("FindByEmail", input.Email).Return(user, nil)
	hasher.On("CheckPasswordHash", input.Password, user.HashedPassword).Return(true)
	token.On("GenerateToken", user.ID).Return("jwt-token", nil)

	result, err := svc.LoginUser(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "jwt-token", result.Token)
}

func TestUserService_LoginUser_InvalidPassword(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.LoginUserInput{Email: "john@example.com", Password: "wrong"}
	user := &entity.User{ID: "1", Email: input.Email, HashedPassword: "hashed123"}

	repo.On("FindByEmail", input.Email).Return(user, nil)
	hasher.On("CheckPasswordHash", input.Password, user.HashedPassword).Return(false)

	result, err := svc.LoginUser(input)

	assert.ErrorIs(t, err, exceptions.ErrInvalidCredentials)
	assert.Nil(t, result)
}

func TestUserService_LoginUser_UserNotFound(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.LoginUserInput{Email: "notfound@example.com", Password: "123"}
	repo.On("FindByEmail", input.Email).Return(nil, errors.New("not found"))

	result, err := svc.LoginUser(input)

	assert.ErrorIs(t, err, exceptions.ErrInvalidCredentials)
	assert.Nil(t, result)
}

func TestUserService_LoginUser_GenerateTokenError(t *testing.T) {
	repo := new(MockUserRepo)
	hasher := new(MockHasher)
	token := new(MockTokenService)

	svc := services.NewUserService(repo, hasher, token)

	input := &dto.LoginUserInput{Email: "john@example.com", Password: "123"}
	user := &entity.User{ID: "1", Email: input.Email, HashedPassword: "hashed123"}

	repo.On("FindByEmail", input.Email).Return(user, nil)
	hasher.On("CheckPasswordHash", input.Password, user.HashedPassword).Return(true)
	token.On("GenerateToken", user.ID).Return("", errors.New("token error"))

	result, err := svc.LoginUser(input)

	assert.ErrorIs(t, err, exceptions.ErrInvalidCredentials)
	assert.Nil(t, result)
}
