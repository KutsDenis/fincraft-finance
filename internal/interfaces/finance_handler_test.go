package interfaces_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

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
		Amount:      10050,
		Description: "Test income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, int64(1), int32(2), int64(10050), "Test income").
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
		Amount:      10050,
		Description: "Test income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, int64(1), int32(2), int64(10050), "Test income").
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
		Amount:      -10050,
		Description: "Negative income",
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		AddIncome(ctx, int64(1), int32(2), int64(-10050), "Negative income").
		Return(errors.New("validation failed: amount must be greater than 0"))

	resp, err := handler.AddIncome(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Contains(t, err.Error(), "validation failed")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsNoError_WhenValidInput(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: &timestamppb.Timestamp{Seconds: 1612137600},
		EndDate:   &timestamppb.Timestamp{Seconds: 1612224000},
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(1), req.StartDate.AsTime(), req.EndDate.AsTime()).
		Return([]*finance.CategoryIncome{}, nil)

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
