package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fincraft-finance/api/finance"
	"fincraft-finance/internal/domain"
)

const (
	minUserID = 1
)

//go:generate mockgen -source=income_usecase.go -destination=mocks/income_usecase_mock.go -package=mocks

// IncomeService контракт сервиса для работы с доходами
type IncomeService interface {
	AddIncome(ctx context.Context, userID int64, categoryID int32, amount int64, description string) error
	GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*finance.CategoryIncome, error)
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
func (u *IncomeUseCase) AddIncome(ctx context.Context, userID int64, categoryID int32, amount int64, description string) error {
	income := &domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.Money(amount),
		Description: description,
	}

	if err := income.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return u.repo.AddIncome(ctx, userID, categoryID, amount, description)
}

// GetIncomesForPeriod возвращает список доходов по категориям за указанный период времени
func (u *IncomeUseCase) GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*finance.CategoryIncome, error) {
	if err := validateGetIncomesForPeriodInput(ctx, userID, startDate, endDate); err != nil {
		return nil, err
	}

	incomes, err := u.repo.GetIncomesForPeriod(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get incomes: %w", err)
	}

	return incomes, nil
}

// MARK: Validations

// validateGetIncomesForPeriodInput валидирует входные данные для GetIncomesForPeriod
func validateGetIncomesForPeriodInput(ctx context.Context, userID int64, startDate, endDate time.Time) error {
	var errors []string

	if ctx == nil {
		errors = append(errors, "context cannot be nil")
	}

	if userID < minUserID {
		errors = append(errors, "user ID must be positive")
	}

	if startDate.IsZero() || endDate.IsZero() {
		errors = append(errors, "dates cannot be zero")
	}

	if !startDate.IsZero() && !endDate.IsZero() && startDate.After(endDate) {
		errors = append(errors, "start date must be before or equal to end date")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, "; "))
	}
	return nil
}
