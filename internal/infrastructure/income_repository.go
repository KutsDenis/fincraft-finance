package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"fincraft-finance/api/finance"
)

// IncomeRepository реализует методы для работы с доходами
type IncomeRepository struct {
	db *sql.DB
}

// NewIncomeRepository создает новый экземпляр IncomeRepository
func NewIncomeRepository(db *sql.DB) *IncomeRepository {
	return &IncomeRepository{db: db}
}

// AddIncome добавляет новый доход в базу данных
func (r *IncomeRepository) AddIncome(ctx context.Context, userID int64, categoryID int32, amount int64, description string) error {
	_, err := r.db.ExecContext(ctx, `
		SELECT * FROM add_income($1, $2, $3, $4)
	`, userID, categoryID, amount, description)

	return err
}

// GetIncomesForPeriod возвращает список доходов за указанный период по категориям
func (r *IncomeRepository) GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*finance.CategoryIncome, error) {
	// retrieve incomes
	rows, err := r.db.QueryContext(ctx, `
		SELECT * FROM get_incomes_for_period($1, $2, $3)
	`, userID, startDate, endDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get incomes")
	}

	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	categoryIncomesMap, err := CategoryIncomesMapFromRows(rows)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map rows to incomes")
	}

	// convert map to slice
	categoryIncomes := CategoryIncomesFromMapOfCategoryIncomes(categoryIncomesMap)

	return categoryIncomes, nil
}

// CategoryIncomesMapFromRows возвращает мап категорий и доходов из бд
func CategoryIncomesMapFromRows(rows *sql.Rows) (map[int32][]*finance.Income, error) {
	categoryIncomesMap := make(map[int32][]*finance.Income)

	for rows.Next() {
		var income finance.Income
		var categoryID int32

		if err := rows.Scan(&categoryID, &income.Amount, &income.Description); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		categoryIncomesMap[categoryID] = append(categoryIncomesMap[categoryID], &income)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating rows")
	}

	return categoryIncomesMap, nil
}

// CategoryIncomesFromMapOfCategoryIncomes возвращает слайс категорий и доходов из мапы категорий и доходов
func CategoryIncomesFromMapOfCategoryIncomes(categoryIncomesMap map[int32][]*finance.Income) []*finance.CategoryIncome {
	categoryIncomes := make([]*finance.CategoryIncome, 0, len(categoryIncomesMap))

	for categoryID, incomes := range categoryIncomesMap {
		categoryIncomes = append(categoryIncomes, &finance.CategoryIncome{
			CategoryId: categoryID,
			Incomes:    incomes,
		})
	}

	return categoryIncomes
}
