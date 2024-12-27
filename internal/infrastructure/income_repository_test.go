package infrastructure_test

import (
	"context"
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

// seedUsersAndIncomes создает пользователей и добавляет доходы в базу данных
func seedUsersAndIncomes(t *testing.T, now time.Time) {
	user := testdb.UserParams{ID: 1}
	err := user.SeedUser(testdb.DB)
	require.NoError(t, err)

	// user 2
	user = testdb.UserParams{ID: 2, Email: "test2@test.com"}
	err = user.SeedUser(testdb.DB)
	require.NoError(t, err)

	income := testdb.IncomeParams{UserID: 1, CreatedAt: now}

	// user 1 with 4 incomes
	err = income.SeedIncomes(testdb.DB, now, now.Add(time.Hour), now.Add(2*time.Hour), now.Add(3*time.Hour))
	require.NoError(t, err)

	// user 2 with 4 incomes
	income = testdb.IncomeParams{UserID: 2, Amount: 20.99}
	err = income.SeedIncomes(testdb.DB, now, now.Add(time.Hour), now.Add(2*time.Hour), now.Add(3*time.Hour))
	require.NoError(t, err)
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

// func Test_IncomeRepository_GetIncomesForPeriod_ReturnsSliceOfIncomes_WhenValidInput(t *testing.T) {
// 	defer func() {
// 		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	now := time.Now()
// 	seedUsersAndIncomes(t, now)
// 	repo := infrastructure.NewIncomeRepository(testdb.DB)

// 	ctx := context.Background()
// 	startTime := now.Format(time.RFC3339)
// 	endDate := now.Add(3 * time.Hour).Format(time.RFC3339)
// 	incomes, err := repo.GetIncomesForPeriod(ctx, 1, startTime, endDate)

// 	require.NoError(t, err)
// 	assert.Len(t, incomes, 3)

// 	assert.Equal(t, incomes[0].UserID, int64(1))
// 	assert.Equal(t, incomes[1].UserID, int64(1))
// 	assert.Equal(t, incomes[2].UserID, int64(1))

// 	assert.Equal(t, incomes[0].Amount, 10050)
// }

// func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenInvalidUserID(t *testing.T) {
// 	defer func() {
// 		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	now := time.Now()
// 	seedUsersAndIncomes(t, now)
// 	repo := infrastructure.NewIncomeRepository(testdb.DB)

// 	ctx := context.Background()
// 	startTime := now.Format(time.RFC3339)
// 	endDate := now.Add(3 * time.Hour).Format(time.RFC3339)
// 	_, err := repo.GetIncomesForPeriod(ctx, 999, startTime, endDate)

// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "does not exist")
// }

// func Test_IncomeRepository_GetIncomesForPeriod_ReturnsError_WhenInvalidDates(t *testing.T) {
// 	defer func() {
// 		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	seedUsersAndIncomes(t, time.Now())
// 	repo := infrastructure.NewIncomeRepository(testdb.DB)

// 	ctx := context.Background()
// 	_, err := repo.GetIncomesForPeriod(ctx, 1, "invalid date", "invalid date")

// 	require.Error(t, err)
// 	assert.Error(t, err, "invalid input syntax for type timestamp")
// }

// func Test_IncomeRepository_GetIncomesForPeriod_ReturnsEmptySlice_WhenNoIncomes(t *testing.T) {
// 	defer func() {
// 		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	now := time.Now()
// 	seedDefaultUser(t)
// 	repo := infrastructure.NewIncomeRepository(testdb.DB)

// 	ctx := context.Background()
// 	startTime := now.Format(time.RFC3339)
// 	endDate := now.Add(3 * time.Hour).Format(time.RFC3339)
// 	incomes, err := repo.GetIncomesForPeriod(ctx, 1, startTime, endDate)

// 	require.NoError(t, err)
// 	assert.Empty(t, incomes)
// }

// func Test_IncomeRepository_GetIncomesForPeriod_ReturnsIncomesForDifferentUsers(t *testing.T) {
// 	defer func() {
// 		if err := testdb.TruncateTables(testdb.DB, testdb.UsersTable, testdb.IncomesTable); err != nil {
// 			t.Fatal(err)
// 		}
// 	}()

// 	now := time.Now()
// 	seedUsersAndIncomes(t, now)
// 	repo := infrastructure.NewIncomeRepository(testdb.DB)

// 	ctx := context.Background()
// 	startTime := now.Format(time.RFC3339)
// 	endDate := now.Add(3 * time.Hour).Format(time.RFC3339)
// 	incomes, err := repo.GetIncomesForPeriod(ctx, 2, startTime, endDate)

// 	require.NoError(t, err)
// 	assert.Len(t, incomes, 3)

// 	assert.Equal(t, incomes[0].UserID, int64(2))
// 	assert.Equal(t, incomes[1].UserID, int64(2))
// 	assert.Equal(t, incomes[2].UserID, int64(2))

// 	assert.Equal(t, incomes[0].Amount, 2099)
// }
