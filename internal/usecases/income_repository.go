package usecases

import (
	"context"
	"time"

	"fincraft-finance/api/finance"
)

//go:generate mockgen -source=income_repository.go -destination=mocks/income_repository_mock.go -package=mocks

// IncomeRepository репозиторий для работы с доходами
type IncomeRepository interface {
	AddIncome(ctx context.Context, userID int64, categoryID int32, amount int64, description string) error
	GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*finance.CategoryIncome, error)
}
