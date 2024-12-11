package infrastructure

import (
	"context"
	"database/sql"
	"fincraft-finance/internal/domain"
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
func (r *IncomeRepository) AddIncome(ctx context.Context, userID int64, categoryID int, amount float64, description string) error {
	_, err := r.db.ExecContext(ctx, `
		SELECT * FROM add_income($1, $2, $3, $4)
	`, userID, categoryID, amount, description)

	return err
}

// GetIncomesForPeriod возвращает список доходов за указанный период
func (r *IncomeRepository) GetIncomesForPeriod(ctx context.Context, userID int64, startDate, endDate string) ([]*domain.Income, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT * FROM get_incomes_for_period($1, $2, $3)
	`, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incomes := []*domain.Income{}
	var amount float64
	for rows.Next() {
		var income domain.Income
		if err := rows.Scan(&income.UserID, &income.CategoryID, &amount, &income.Description); err != nil {
			return nil, err
		}
		income.Amount = domain.NewMoneyFromFloat(amount)
		incomes = append(incomes, &income)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return incomes, nil
}
