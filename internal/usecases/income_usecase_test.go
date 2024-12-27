package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsNoError_WhenValidInput(t *testing.T) {
	ctrl, mockRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().GetIncomesForPeriod(ctx, int64(1), gomock.Any(), gomock.Any()).Return(nil, nil)

	now := time.Now()
	startDate := now.Add(-time.Hour * 24)
	endDate := now.Add(time.Hour * 24)
	_, err := useCase.GetIncomesForPeriod(ctx, int64(1), startDate, endDate)

	assert.NoError(t, err)
}

func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenInvalidInput(t *testing.T) {
	ctrl, _, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	tests := []struct {
		name      string
		userID    int64
		startDate time.Time
		endDate   time.Time
		errMsg    string
		ctx       context.Context
	}{
		{"Zero UserID", 0, time.Now(), time.Now(), "validation failed: user ID must be positive", ctx},
		{"StartDate After EndDate", 1, time.Now().Add(time.Hour * 24), time.Now(), "validation failed: start date must be before or equal to end date", ctx},
		{"Zero Dates", 1, time.Time{}, time.Time{}, "validation failed: dates cannot be zero", ctx},
		{"Nil Context", 1, time.Now(), time.Now(), "validation failed: context cannot be nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := useCase.GetIncomesForPeriod(tt.ctx, tt.userID, tt.startDate, tt.endDate)

			assert.Error(t, err)
			assert.EqualError(t, err, tt.errMsg)
		})
	}
}

func Test_IncomeUseCase_GetIncomesForPeriod_ReturnsError_WhenRepoFails(t *testing.T) {
	ctrl, mockRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().GetIncomesForPeriod(ctx, int64(1), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

	now := time.Now()
	_, err := useCase.GetIncomesForPeriod(ctx, int64(1), now, now)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get incomes:")
}
