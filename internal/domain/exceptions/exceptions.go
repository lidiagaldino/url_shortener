package exceptions

import "errors"

var (
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrUserNotFound              = errors.New("user not found")
	ErrURLNotFound               = errors.New("url not found")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrInvalidURL                = errors.New("invalid url")
	ErrUnauthorizedURLStatistics = errors.New("unauthorized to access url statistics")
)
