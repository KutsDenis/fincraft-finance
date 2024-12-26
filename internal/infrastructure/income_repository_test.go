package infrastructure_test

import (
	"context"
	"fincraft-finance/internal/domain"
	"testing"
	"time"

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

// seedUserIncomes создает пользователей и добавляет доходы в базу данных
func seedUserIncomes(t *testing.T, income domain.Income, incomesCount int) {
	user := testdb.UserParams{ID: int(income.UserID)}
	err := user.SeedUser(testdb.DB)
	require.NoError(t, err)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	for i := 0; i < incomesCount; i++ {
		err = repo.AddIncome(ctx, income.UserID, income.CategoryID, income.Amount.ToFloat(), income.Description)
		require.NoError(t, err)
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
	err := repo.AddIncome(ctx, 1, 2, 100.50, "test income")
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
	err := repo.AddIncome(ctx, 999, 2, 100.50, "Invalid user")
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
	err := repo.AddIncome(ctx, 1, 2, -100.50, "Negative amount")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "violates check constraint")
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsIncomes_WhenValidInput(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	userID := int64(1)
	incomesCount := 4
	categoryID := 1
	amount := 100.50

	seedUserIncomes(t, domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: "test income",
	}, incomesCount)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	startTime := time.Now()
	endTime := startTime.Add(3 * time.Hour)
	incomes, err := repo.GetIncomesForPeriod(ctx, userID, startTime, endTime)
	require.NoError(t, err)

	assert.Len(t, incomes, incomesCount)
	for _, income := range incomes {
		assert.Equal(t, income.UserID, userID)
		assert.Equal(t, income.CategoryID, categoryID)
		assert.Equal(t, income.Amount.ToFloat(), amount)
	}
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsIncomes_WhenValidInputAndAnotherUser(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	userID := int64(2)
	incomesCount := 3
	categoryID := 1
	amount := 20.99

	seedUserIncomes(t, domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: "test income",
	}, incomesCount)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour)
	incomes, err := repo.GetIncomesForPeriod(ctx, userID, startTime, endTime)
	require.NoError(t, err)

	assert.Len(t, incomes, incomesCount)
	for _, income := range incomes {
		assert.Equal(t, income.UserID, userID)
		assert.Equal(t, income.CategoryID, categoryID)
		assert.Equal(t, income.Amount.ToFloat(), amount)
	}
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenInvalidUserID(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	userID := int64(999)
	incomesCount := 1
	categoryID := 1
	amount := 20.99

	seedUserIncomes(t, domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: "test income",
	}, incomesCount)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour)
	_, err := repo.GetIncomesForPeriod(ctx, userID, startTime, endTime)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenInvalidDates(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	userID := int64(1)
	incomesCount := 1
	categoryID := 1
	amount := 20.99

	seedUserIncomes(t, domain.Income{
		UserID:      userID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: "test income",
	}, incomesCount)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	testCases := []struct {
		name      string
		startTime time.Time
		endTime   time.Time
	}{
		{
			name:      "end_date_one_day_before_start",
			startTime: time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "end_date_one_hour_before_start",
			startTime: time.Now(),
			endTime:   time.Now().Add(-1 * time.Hour),
		},
		{
			name:      "end_date_one_minute_before_start",
			startTime: time.Now(),
			endTime:   time.Now().Add(-1 * time.Minute),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.GetIncomesForPeriod(ctx, userID, tc.startTime, tc.endTime)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "End date cannot be before start date")
		})
	}
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsEmptySlice_WhenNoIncomes(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	seedDefaultUser(t)
	repo := infrastructure.NewIncomeRepository(testdb.DB)

	ctx := context.Background()
	startTime := time.Now()
	endTime := startTime.Add(3 * time.Hour)
	incomes, err := repo.GetIncomesForPeriod(ctx, 1, startTime, endTime)

	require.NoError(t, err)
	assert.Empty(t, incomes)
}
