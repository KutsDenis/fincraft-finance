package infrastructure_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"fincraft-finance/internal/infrastructure"
	"fincraft-finance/internal/testdb"
)

func TestMain(m *testing.M) {
	if err := testdb.SetupTestDB(); err != nil {
		panic(err)
	}
	defer testdb.CloseTestDB()

	m.Run()
}

func seedDefaultUser(t *testing.T) {
	user := testdb.UserParams{}
	err := user.SeedUser(testdb.DB)
	require.NoError(t, err)
}

type UserIncomes struct {
	UserID                int64
	Email                 string
	CategoryIncomesByDate map[time.Time][]CategoryIncomes
}

type CategoryIncomes struct {
	CategoryID int32
	Count      int
	Amount     int64
}

// seedUsersAndIncomes создает пользователей и добавляет доходы в базу данных
func seedUsersAndIncomes(t *testing.T, userIncomes []UserIncomes) {
	if len(userIncomes) == 0 {
		t.Error("userIncomes must not be empty")
		return
	}

	// For each user
	for _, userIncome := range userIncomes {
		user := testdb.UserParams{ID: userIncome.UserID, Email: userIncome.Email}
		err := user.SeedUser(testdb.DB)
		require.NoError(t, err)

		for date, categoryIncomes := range userIncome.CategoryIncomesByDate {
			for idx, categoryIncome := range categoryIncomes {
				income := testdb.IncomeParams{
					UserID:     user.ID,
					CategoryID: categoryIncome.CategoryID,
					Amount:     categoryIncome.Amount,
					CreatedAt:  date,
				}

				err := income.SeedIncome(testdb.DB)
				require.NoErrorf(t, err, "failed to seed income %s", idx)
			}
		}
	}
}

func Test_IncomeRepository_AddIncome_ReturnsNoError_WhenValidInput(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	seedDefaultUser(t)
	repo := infrastructure.NewIncomeRepository(testdb.DB)

	ctx := context.Background()
	err := repo.AddIncome(ctx, 1, 2, 10050, "test income")
	assert.NoError(t, err)
}

func Test_IncomeRepository_AddIncome_ReturnsError_WhenUserInvalid(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	repo := infrastructure.NewIncomeRepository(testdb.DB)

	ctx := context.Background()
	err := repo.AddIncome(ctx, 999, 2, 10050, "Invalid user")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violates foreign key constraint")
}

func Test_IncomeRepository_AddIncome_ReturnsError_WhenInvalidAmount(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	seedDefaultUser(t)
	repo := infrastructure.NewIncomeRepository(testdb.DB)

	ctx := context.Background()
	err := repo.AddIncome(ctx, 1, 2, -10050, "Negative amount")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violates check constraint")
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsSliceOfCategoryIncomes_WhenValidInput(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	now := time.Now()
	yesterday := now.Add(-time.Hour * 24)
	user1 := int64(1)
	user1Email := "user1@example.com"
	user2 := int64(2)
	user2Email := "user2@example.com"
	seedUsersAndIncomes(t, []UserIncomes{
		{UserID: user1, Email: user1Email, CategoryIncomesByDate: map[time.Time][]CategoryIncomes{
			now: {
				{CategoryID: 1, Count: 1, Amount: 100},
				{CategoryID: 1, Count: 3, Amount: 33},
				{CategoryID: 2, Count: 2, Amount: 200},
			},
			yesterday: {
				{CategoryID: 1, Count: 1, Amount: 50},
				{CategoryID: 2, Count: 1, Amount: 100},
			},
		}},
		{UserID: user2, Email: user2Email, CategoryIncomesByDate: map[time.Time][]CategoryIncomes{
			now: {
				{CategoryID: 1, Count: 1, Amount: 100},
				{CategoryID: 1, Count: 3, Amount: 33},
			},
		}},
	})

	repo := infrastructure.NewIncomeRepository(testdb.DB)

	ctx := context.Background()
	startDate := yesterday.Add(-time.Hour * 24)
	endDate := now.Add(time.Hour * 24)

	categoryIncomes, err := repo.GetIncomesForPeriod(ctx, user1, startDate, endDate)

	require.NoError(t, err, "failed to get incomes")
	require.Len(t, categoryIncomes, 2, "unexpected number of category incomes")
}

func Test_IncomeRepository_categoryIncomesMapFromRows_ReturnsError_WhenInvalidRows(t *testing.T) {
	rows, err := testdb.DB.QueryContext(context.Background(), "SELECT 1")
	assert.NoError(t, err)

	_, err = infrastructure.CategoryIncomesMapFromRows(rows)
	assert.Error(t, err)
}
func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenRowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := infrastructure.NewIncomeRepository(db)

	mock.ExpectQuery(`SELECT (.+) FROM get_incomes_for_period\(\$1, \$2, \$3\)`).
		WillReturnRows(sqlmock.NewRows([]string{"category_id", "amount", "description", "created_at"}).
			AddRow(1, 100, "test", time.Now()).
			RowError(0, fmt.Errorf("row error")))

	ctx := context.Background()
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	_, err = repo.GetIncomesForPeriod(ctx, 1, startDate, endDate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error iterating rows")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenQueryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := infrastructure.NewIncomeRepository(db)

	mock.ExpectQuery(`SELECT (.+) FROM get_incomes_for_period\(\$1, \$2, \$3\)`).
		WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf("database error"))

	ctx := context.Background()
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	_, err = repo.GetIncomesForPeriod(ctx, 1, startDate, endDate)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get incomes")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
