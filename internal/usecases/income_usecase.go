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
	GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]domain.Income, error)
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
func (u *IncomeUseCase) GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]domain.Income, error) {
	if err := validateGetIncomesForPeriodInput(userID, startDate, endDate); err != nil {
		return nil, err
	}

	incomes, err := u.repo.GetIncomesForPeriod(ctx, userID, startDate, endDate)
	if len(incomes) == 0 {
		return nil, &NotFoundError{Entity: "Incomes"}
	}
	if err != nil {
		return nil, err
	}

	return incomes, nil
}

func validateGetIncomesForPeriodInput(userID int64, startTime, endTime time.Time) error {
	if userID <= 0 {
		return &ValidationError{Message: fmt.Sprintf("invalid user ID: %d", userID)}
	}

	if startTime.IsZero() {
		return &ValidationError{Message: "start time cannot be zero"}
	}

	if endTime.IsZero() {
		return &ValidationError{Message: "end time cannot be zero"}
	}

	if startTime.After(endTime) {
		return &ValidationError{Message: "start time must be before end time"}
	}

	return nil
}
