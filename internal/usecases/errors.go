package usecases

import "fmt"

// ValidationError представляет ошибку неверного входного значения, сбоя проверки данных или отсутствия определённых полей
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NotFoundError представляет ошибку, возникающую, если данные не найдены
type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Entity)
}
