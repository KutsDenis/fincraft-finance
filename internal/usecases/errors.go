package usecases

import (
	"errors"
	"fmt"
)

var (
	// ValidationError представляет ошибку неверного входного значения, сбоя проверки данных или отсутствия определённых полей
	ValidationError = errors.New("validation error")
	// NotFoundError представляет ошибку, возникающую, если данные не найдены
	NotFoundError = errors.New("not found error")
)

func NewValidationError(msg string) error {
	return fmt.Errorf("%w: %s", ValidationError, msg)
}

func NewNotFoundError(entity string) error {
	return fmt.Errorf("%w: %s", NotFoundError, entity)
}
