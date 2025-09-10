package services

import (
	"url-shortener/internal/domain/entity"
	"url-shortener/internal/domain/exceptions"
	"url-shortener/internal/domain/repository"
	"url-shortener/internal/domain/security"
	"url-shortener/internal/services/dto"
	"url-shortener/pkg"
)

type UserService struct {
	repo   repository.UserRepository
	hasher pkg.Hasher
	token  security.TokenService
}

func NewUserService(repo repository.UserRepository, hasher pkg.Hasher, token security.TokenService) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
		token:  token,
	}
}

func (s *UserService) Save(user *dto.UserInput) (*dto.UserOutput, error) {
	userEmail, _ := s.repo.FindByEmail(user.Email)
	if userEmail != nil {
		return nil, exceptions.ErrEmailAlreadyExists
	}
	hashedPassword, err := s.hasher.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	userEntity := entity.User{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}

	savedUser, err := s.repo.Save(&userEntity)
	if err != nil {
		return nil, err
	}

	result := dto.UserOutput{
		ID:        savedUser.ID,
		Name:      savedUser.Name,
		Email:     savedUser.Email,
		CreatedAt: savedUser.CreatedAt,
	}

	return &result, nil
}

func (s *UserService) LoginUser(input *dto.LoginUserInput) (*dto.LoginUserOutput, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, exceptions.ErrInvalidCredentials
	}

	if !s.hasher.CheckPasswordHash(input.Password, user.HashedPassword) {
		return nil, exceptions.ErrInvalidCredentials
	}

	t, err := s.token.GenerateToken(user.ID)
	if err != nil {
		return nil, exceptions.ErrInvalidCredentials
	}

	result := dto.LoginUserOutput{
		Token: t,
	}

	return &result, nil
}
