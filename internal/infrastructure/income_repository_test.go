package infrastructure_test

import (
	"context"
	"database/sql"
	"fincraft-finance/internal/infrastructure"
	"fincraft-finance/internal/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

func TestIntegration_IncomeRepository_AddIncome_Success(t *testing.T) {
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

func TestIntegration_IncomeRepository_AddIncome_HandlesDBError(t *testing.T) {
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

func TestIntegration_IncomeRepository_AddIncome_HandlesInvalidData(t *testing.T) {
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

func TestIntegration_IncomeRepository_AddIncome_HandlesDBConnectionError(t *testing.T) {
	invalidDB, err := sql.Open("postgres",
		"postgres://invalid_user:invalid_password@localhost:5434/invalid_db?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	//noinspection GoUnhandledErrorResult
	defer invalidDB.Close()

	repo := infrastructure.NewIncomeRepository(invalidDB)

	ctx := context.Background()
	err = repo.AddIncome(ctx, 1, 2, 100.50, "Test income")
	assert.Error(t, err)
	assert.Contains(t, err.Error(),
		"dial tcp [::1]:5434: connectex: "+
			"No connection could be made because the target machine actively refused it.")
}