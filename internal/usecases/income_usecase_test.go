package usecases_test

import (
	"context"
	"errors"
	"testing"

	"fincraft-finance/internal/usecases"
	"fincraft-finance/internal/usecases/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockIncomeRepository, *usecases.IncomeUseCase) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIncomeRepository(ctrl)
	useCase := usecases.NewIncomeUseCase(mockRepo)
	return ctrl, mockRepo, useCase
}

func TestIncomeUseCase_AddIncome_WhenValidInput_ShouldSucceed(t *testing.T) {
	ctrl, mockRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().AddIncome(ctx, 1, 2, 100.50, "Test income").Return(nil)

	err := useCase.AddIncome(ctx, 1, 2, 100.50, "Test income")

	assert.NoError(t, err)
}

func TestIncomeUseCase_AddIncome_WhenInvalidInput_ShouldReturnValidationError(t *testing.T) {
	_, _, useCase := setupTest(t)

	tests := []struct {
		name   string
		userID int64
		catID  int
		amount float64
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

func TestIncomeUseCase_AddIncome_WhenRepoFails_ShouldReturnError(t *testing.T) {
	ctrl, mockRepo, useCase := setupTest(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo.EXPECT().AddIncome(ctx, 1, 2, 100.0, "Test income").Return(errors.New("db error"))

	err := useCase.AddIncome(ctx, 1, 2, 100, "Test income")

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
}
