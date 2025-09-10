package exceptions

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
	ErrUserNotFound       = errors.New("usuário não encontrado")
)
