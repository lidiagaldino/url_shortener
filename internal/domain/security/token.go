package security

type TokenService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (string, error)
}
