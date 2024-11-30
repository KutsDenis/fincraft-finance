package domain

import "errors"

// Income представляет бизнес-объект дохода
type Income struct {
	UserID      int64
	CategoryID  int
	Amount      Money
	Description string
}

// Validate проверяет бизнес-правила для дохода
func (i *Income) Validate() error {
	if i.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if i.UserID <= 0 {
		return errors.New("user ID must be valid")
	}
	if i.CategoryID <= 0 {
		return errors.New("category ID must be valid")
	}
	return nil
}
