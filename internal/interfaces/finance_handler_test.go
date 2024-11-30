package interfaces_test

import (
	"context"
	"errors"
	"fincraft-finance/api/finance"
	"fincraft-finance/internal/interfaces"
	"fincraft-finance/internal/usecases/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockIncomeService, *interfaces.FinanceHandler) {
	ctrl := gomock.NewController(t)
	mockUsecase := mocks.NewMockIncomeService(ctrl)
	handler := interfaces.NewFinanceHandler(mockUsecase)

	return ctrl, mockUsecase, handler
}

func TestFinanceHandler_AddIncome_WhenValidInput_ShouldReturnResponse(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.AddIncomeRequest{
		UserId:      1,
		CategoryId:  2,
		Amount:      100.50,
		Description: "Test income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, 1, 2, 100.50, "Test income").
		Return(nil)

	resp, err := handler.AddIncome(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestFinanceHandler_AddIncome_WhenUseCaseFails_ShouldReturnInternalError(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.AddIncomeRequest{
		UserId:      1,
		CategoryId:  2,
		Amount:      100.50,
		Description: "Test income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, 1, 2, 100.50, "Test income").
		Return(nil, errors.New("db error"))

	resp, err := handler.AddIncome(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "failed to add income: db error")
}

func TestFinanceHandler_AddIncome_WhenInvalidInput_ShouldReturnValidationError(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.AddIncomeRequest{
		UserId:      1,
		CategoryId:  2,
		Amount:      -100.50,
		Description: "Negative income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, 1, 2, -100.50, "Negative income").
		Return(nil, errors.New("validation failed: amount must be greater than 0"))

	resp, err := handler.AddIncome(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "validation failed")
}
