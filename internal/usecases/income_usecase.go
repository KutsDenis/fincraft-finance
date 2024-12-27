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

	if err != nil {
		return nil, err
	}

	if len(incomes) == 0 {
		return nil, NewNotFoundError("incomes")
	}

	return incomes, nil
}

func validateGetIncomesForPeriodInput(userID int64, startTime, endTime time.Time) error {
	if userID <= 0 {
		return NewValidationError(
			fmt.Sprintf("invalid user ID: %d", userID))
	}

	if startTime.IsZero() {
		return NewValidationError(
			fmt.Sprintf("invalid start time: %s", startTime))
	}

	if endTime.IsZero() {
		return NewValidationError(
			fmt.Sprintf("invalid end time: %s", endTime))
	}

	if startTime.After(endTime) {
		return NewValidationError(
			fmt.Sprintf("start time must be before end time"))
	}

	return nil
}
