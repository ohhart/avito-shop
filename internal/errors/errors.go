package errors

import "errors"

var (
	ErrInternalServer     = errors.New("внутренняя ошибка сервера")
	ErrInvalidRequest     = errors.New("неверный запрос")
	ErrUnauthorized       = errors.New("неавторизованный доступ")
	ErrNotFound           = errors.New("ресурс не найден")
	ErrInsufficientFunds  = errors.New("недостаточно средств")
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrItemNotFound       = errors.New("товар не найден")
	ErrUserNotFound       = errors.New("пользователь не найден")
)
