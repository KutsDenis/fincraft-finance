package usecases

import (
	"context"
	"fmt"
	"time"

	"fincraft-finance/internal/domain"
)

//go:generate mockgen -source=income_usecase.go -destination=mocks/income_usecase_mock.go -package=mocks

// IncomeService контракт сервиса для работы с доходами
type IncomeService interface {
	AddIncome(ctx context.Context, userID int64, categoryID int, amount float64, description string) error
	GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate string) ([]*domain.Income, error)
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

// GetIncomesForPeriod возвращает список доходов за указанный период
func (u *IncomeUseCase) GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate string) ([]*domain.Income, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("validation failed : invalid user ID: %d", userID)
	}
	startTime, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("validation failed : failed to parse start date: %w", err)
	}
	endTime, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("validation failed : failed to parse end date: %w", err)
	}

	if startTime.After(endTime) {
		return nil, fmt.Errorf("validation failed : start date is after end date")
	}

	return u.repo.GetIncomesForPeriod(ctx, userID, startDate, endDate)
}
