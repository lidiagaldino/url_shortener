package security

import "github.com/golang-jwt/jwt/v5"

type TokenService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}
