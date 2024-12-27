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

func setupTestTime() time.Time {
	return time.Now().UTC()
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
	testCases := []struct {
		name         string
		userID       int64
		incomesCount int
		categoryID   int
		amount       float64
		description  string
	}{
		{
			name:         "returns_incomes_when_valid_input",
			userID:       1,
			incomesCount: 4,
			categoryID:   1,
			amount:       100.50,
			description:  "test income",
		},
		{
			name:         "returns_incomes_for_another_user",
			userID:       2,
			incomesCount: 3,
			categoryID:   1,
			amount:       20.99,
			description:  "test income",
		},
		{
			name:         "returns_incomes_for_large_amount",
			userID:       3,
			incomesCount: 5,
			categoryID:   1,
			amount:       1000.00,
			description:  "test income large amount",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
					t.Fatal(err)
				}
			}()

			startTime := setupTestTime()

			seedUserIncomes(t, domain.Income{
				UserID:      tc.userID,
				CategoryID:  tc.categoryID,
				Amount:      domain.NewMoneyFromFloat(tc.amount),
				Description: tc.description,
			}, tc.incomesCount)

			endTime := setupTestTime()

			repo := infrastructure.NewIncomeRepository(testdb.DB)
			ctx := context.Background()

			incomes, err := repo.GetIncomesForPeriod(ctx, tc.userID, startTime, endTime)

			require.NoError(t, err)
			require.Len(t, incomes, tc.incomesCount)
			for _, income := range incomes {
				assert.Equal(t, tc.userID, income.UserID)
				assert.Equal(t, tc.categoryID, income.CategoryID)
				assert.Equal(t, tc.amount, income.Amount.ToFloat())
				assert.Equal(t, tc.description, income.Description)
			}
		})
	}
}
func Test_IncomeRepository_GetIncomesForPeriod_ReturnsEmptyIncomes_WhenNoIncomesInPeriod(t *testing.T) {
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

	startTime := setupTestTime().Add(-2 * time.Hour)
	endTime := setupTestTime().Add(-1 * time.Hour)
	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	incomes, err := repo.GetIncomesForPeriod(ctx, userID, startTime, endTime)

	require.NoError(t, err)
	assert.Len(t, incomes, 0)
	assert.NotNil(t, incomes)
	assert.Empty(t, incomes)

}
func Test_IncomeRepository_GetIncomesForPeriod_ReturnsPartialIncomes_WhenSomeIncomesInPeriod(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	userID := int64(1)
	categoryID := 1
	amount := 20.99
	targetIncomesCount := 3

	user := testdb.UserParams{ID: int(userID)}
	err := user.SeedUser(testdb.DB)
	require.NoError(t, err)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	// Create old incomes (outside target period)
	for i := 0; i < 2; i++ {
		err = repo.AddIncome(ctx, userID, categoryID, amount, "old test income")
		require.NoError(t, err)
	}

	startTime := setupTestTime()
	// Create target period incomes
	for i := 0; i < targetIncomesCount; i++ {
		err = repo.AddIncome(ctx, userID, categoryID, amount, "target test income")
		require.NoError(t, err)
	}
	endTime := setupTestTime()

	// Create future incomes (outside target period)
	for i := 0; i < 2; i++ {
		err = repo.AddIncome(ctx, userID, categoryID, amount, "future test income")
		require.NoError(t, err)
	}

	// Get incomes for our target period
	incomes, err := repo.GetIncomesForPeriod(ctx, userID, startTime, endTime)

	require.NoError(t, err)
	assert.NotNil(t, incomes)
	assert.Len(t, incomes, targetIncomesCount)

	for _, income := range incomes {
		assert.Equal(t, "target test income", income.Description)
		assert.Equal(t, userID, income.UserID)
		assert.Equal(t, categoryID, income.CategoryID)
		assert.Equal(t, amount, income.Amount.ToFloat())
	}
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenInvalidUserID(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	validUserID := int64(1)
	incomesCount := 1
	categoryID := 1
	amount := 20.99

	startTime := setupTestTime()

	seedUserIncomes(t, domain.Income{
		UserID:      validUserID,
		CategoryID:  categoryID,
		Amount:      domain.NewMoneyFromFloat(amount),
		Description: "test income",
	}, incomesCount)

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	endTime := setupTestTime()
	invalidUserID := int64(999)
	_, err := repo.GetIncomesForPeriod(ctx, invalidUserID, startTime, endTime)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "User does not exist")
}

func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenUserNotExist(t *testing.T) {
	defer func() {
		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
			t.Fatal(err)
		}
	}()

	repo := infrastructure.NewIncomeRepository(testdb.DB)
	ctx := context.Background()

	startTime := setupTestTime()
	endTime := setupTestTime()
	notExistingUserID := int64(999)
	_, err := repo.GetIncomesForPeriod(ctx, notExistingUserID, startTime, endTime)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "User does not exist")
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
			startTime: setupTestTime(),
			endTime:   setupTestTime().Add(-1 * time.Hour),
		},
		{
			name:      "end_date_one_minute_before_start",
			startTime: setupTestTime(),
			endTime:   setupTestTime().Add(-1 * time.Minute),
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
	startTime := setupTestTime()
	endTime := setupTestTime()
	incomes, err := repo.GetIncomesForPeriod(ctx, 1, startTime, endTime)

	require.NoError(t, err)
	assert.Empty(t, incomes)
}
