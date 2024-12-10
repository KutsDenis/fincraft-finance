package infrastructure

import (
	"context"
	"database/sql"
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

// GetIncomeForPeriod возвращает общий доход за период
func (r *IncomeRepository) GetIncomeForPeriod(ctx context.Context, userID int64, startDate string, endDate string) (float64, error) {
	var totalIncome float64
	err := r.db.QueryRowContext(ctx, `
		SELECT get_income_for_period($1, $2, $3)
	`, userID, startDate, endDate).Scan(&totalIncome)

	return totalIncome, err
}
