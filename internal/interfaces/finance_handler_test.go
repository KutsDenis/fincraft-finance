package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"fincraft-finance/api/finance"
	"fincraft-finance/internal/interfaces"
	"fincraft-finance/internal/usecases/mocks"
)

func setupTest(t *testing.T) (*gomock.Controller, *mocks.MockIncomeService, *interfaces.FinanceHandler) {
	ctrl := gomock.NewController(t)
	mockUsecase := mocks.NewMockIncomeService(ctrl)
	handler := interfaces.NewFinanceHandler(mockUsecase)

	return ctrl, mockUsecase, handler
}

func Test_FinanceHandler_AddIncome_ReturnsNoError_WhenValidInput(t *testing.T) {
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
		AddIncome(ctx, int64(1), 2, 100.50, "Test income").
		Return(nil)

	resp, err := handler.AddIncome(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func Test_FinanceHandler_AddIncome_ReturnsInternalError_WhenUseCaseFails(t *testing.T) {
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
		AddIncome(ctx, int64(1), 2, 100.50, "Test income").
		Return(errors.New("db error"))

	resp, err := handler.AddIncome(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "failed to add income: db error")
}

func Test_FinanceHandler_AddIncome_ReturnsValidationError_WhenInvalidAmount(t *testing.T) {
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
		AddIncome(ctx, int64(1), 2, -100.50, "Negative income").
		Return(errors.New("validation failed: amount must be greater than 0"))

	resp, err := handler.AddIncome(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "validation failed")
}

func Test_FinanceHandler_GetIncomeForPeriod_ReturnsTotalIncome_WhenValidInput(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomeForPeriodRequest{
		UserId:    1,
		StartDate: "2024-12-01",
		EndDate:   "2024-12-31",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomeForPeriod(ctx, int64(1), "2024-12-01", "2024-12-31").
		Return(100.50, nil)

	resp, err := handler.GetIncomeForPeriod(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 100.50, resp.TotalIncome)
}
