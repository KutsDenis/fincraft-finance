package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"fincraft-finance/internal/usecases"
	"fincraft-finance/internal/usecases/mocks"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockIncomeRepository, *usecases.IncomeUseCase) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIncomeRepository(ctrl)
	useCase := usecases.NewIncomeUseCase(mockRepo)
	return ctrl, mockRepo, useCase
}

func Test_IncomeUseCase_AddIncome_ReturnsNoError_WhenValidInput(t *testing.T) {
	ctrl, mockRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().AddIncome(ctx, int64(1), int32(2), int64(10050), "Test income").Return(nil)

	err := useCase.AddIncome(ctx, int64(1), int32(2), int64(10050), "Test income")

	assert.NoError(t, err)
}

func Test_IncomeUseCase_AddIncome_ReturnsValidationError_WhenInvalidInput(t *testing.T) {
	_, _, useCase := setupTest(t)

	tests := []struct {
		name   string
		userID int64
		catID  int32
		amount int64
		desc   string
		errMsg string
	}{
		{"Negative Amount", 1, 2, -100, "Invalid income", "validation failed: amount must be greater than 0"},
		{"Zero UserID", 0, 2, 100, "Invalid income", "validation failed: user ID must be valid"},
		{"Zero CategoryID", 1, 0, 100, "Invalid income", "validation failed: category ID must be valid"},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := useCase.AddIncome(ctx, tt.userID, tt.catID, tt.amount, tt.desc)

			assert.Error(t, err)
			assert.EqualError(t, err, tt.errMsg)
		})
	}
}

// func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsIncomes_WhenValidInput(t *testing.T) {
// 	ctrl, mockRepo, useCase := setupTest(t)
// 	defer ctrl.Finish()

// 	ctx := context.Background()
// 	startDate := "2023-01-01T00:00:00Z"
// 	endDate := "2023-01-31T23:59:59Z"
// 	expectedIncomes := []*domain.Income{
// 		{UserID: 1, CategoryID: 2, Amount: domain.Money(1000), Description: "Income 1"},
// 		{UserID: 1, CategoryID: 3, Amount: domain.Money(2000), Description: "Income 2"},
// 	}

// 	mockRepo.EXPECT().GetIncomesForPeriod(ctx, int64(1), startDate, endDate).Return(expectedIncomes, nil)

// 	incomes, err := useCase.GetIncomesForPeriod(ctx, int64(1), startDate, endDate)

// 	require.NoError(t, err)
// 	require.Equal(t, expectedIncomes, incomes)
// }

// func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenInvalidUserID(t *testing.T) {
// 	_, _, useCase := setupTest(t)

// 	ctx := context.Background()
// 	startDate := "2023-01-01T00:00:00Z"
// 	endDate := "2023-01-31T23:59:59Z"

// 	_, err := useCase.GetIncomesForPeriod(ctx, int64(0), startDate, endDate)

// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), "invalid user ID: 0")
// }

// func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenInvalidStartDate(t *testing.T) {
// 	_, _, useCase := setupTest(t)

// 	ctx := context.Background()
// 	startDate := "invalid-date"
// 	endDate := "2023-01-31T23:59:59Z"

// 	_, err := useCase.GetIncomesForPeriod(ctx, int64(1), startDate, endDate)

// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), "failed to parse start date")
// }

// func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenInvalidEndDate(t *testing.T) {
// 	_, _, useCase := setupTest(t)

// 	ctx := context.Background()
// 	startDate := "2023-01-01T00:00:00Z"
// 	endDate := "invalid-date"

// 	_, err := useCase.GetIncomesForPeriod(ctx, int64(1), startDate, endDate)

// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), "failed to parse end date")
// }

// func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenStartDateIsAfterEndDate(t *testing.T) {
// 	_, _, useCase := setupTest(t)

// 	ctx := context.Background()
// 	startDate := "2023-01-31T23:59:59Z"
// 	endDate := "2023-01-01T00:00:00Z"

// 	_, err := useCase.GetIncomesForPeriod(ctx, int64(1), startDate, endDate)

// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), "start date is after end date")
// }
