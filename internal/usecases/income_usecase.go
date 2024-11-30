package usecases

import (
	"context"
	"fincraft-finance/internal/domain"
	"fmt"
)

//go:generate mockgen -source=income_usecase.go -destination=mocks/income_usecase_mock.go -package=mocks

// IncomeService контракт сервиса для работы с доходами
type IncomeService interface {
	AddIncome(ctx context.Context, userID int64, categoryID int, amount float64, description string) error
}

// IncomeUseCase use-case для работы с доходами
type IncomeUseCase struct {
	repo IncomeRepository
}

// NewIncomeUseCase создает новый экземпляр IncomeUseCase
func NewIncomeUseCase(repo IncomeRepository) *IncomeUseCase {
	return &IncomeUseCase{repo: repo}
}

// AddIncome добавляет новый доход в хранилище данных
func (u *IncomeUseCase) AddIncome(ctx context.Context, userID int64, categoryID int, amount float64, description string) error {
	income := &domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: description,
	}

	if err := income.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.repo.AddIncome(ctx, userID, categoryID, amount, description)
}
