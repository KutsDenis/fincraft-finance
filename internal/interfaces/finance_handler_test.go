package interfaces_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"fincraft-finance/api/finance"
	"fincraft-finance/internal/domain"
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

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsSliceOfIncomes_WhenValidInput(t *testing.T) {
	ctrl, mockUseCase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: timestamppb.New(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	expectedIncomesFromUseCase := []*domain.Income{
		{UserID: 1, CategoryID: 2, Amount: domain.NewMoneyFromFloat(100.0), Description: "Income 1"},
		{UserID: 1, CategoryID: 3, Amount: domain.NewMoneyFromFloat(200.0), Description: "Income 2"},
	}

	expectedIncomes := []*finance.Income{
		{UserId: 1, CategoryId: 2, Amount: 100.0, Description: "Income 1"},
		{UserId: 1, CategoryId: 3, Amount: 200.0, Description: "Income 2"},
	}

	ctx := context.Background()
	mockUseCase.EXPECT().GetIncomesForPeriod(ctx, int64(1), "2024-12-01T00:00:00Z", "2024-12-31T23:59:59Z").Return(expectedIncomesFromUseCase, nil)

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.NoError(t, err)
	require.Equal(t, expectedIncomes, resp.Incomes)
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsInternalError_WhenUseCaseFails(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: timestamppb.New(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(1), "2024-12-01T00:00:00Z", "2024-12-31T23:59:59Z").
		Return(nil, errors.New("db error"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.Internal, status.Code(err))
	require.Contains(t, err.Error(), "failed to get incomes: db error")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsInvalidArgumentError_WhenInvalidUserID(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    0,
		StartDate: timestamppb.New(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(0), "2024-12-01T00:00:00Z", "2024-12-31T23:59:59Z").
		Return(nil, errors.New("validation failed : invalid user ID: 0"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, err.Error(), "invalid user ID: 0")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsNotFoundError_WhenNotExistingUserID(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    999,
		StartDate: timestamppb.New(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(999), "2024-12-01T00:00:00Z", "2024-12-31T23:59:59Z").
		Return(nil, errors.New("User with ID 999 does not exist"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
	require.Contains(t, err.Error(), "User with ID 999 does not exist")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsInvalidArgumentError_WhenInvalidStartDate(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: timestamppb.New(time.Date(1969, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(1), "1969-12-01T00:00:00Z", "2024-12-31T23:59:59Z").
		Return(nil, errors.New("validation failed : failed to parse start date"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, err.Error(), "failed to parse start date")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsInvalidArgumentError_WhenInvalidEndDate(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: timestamppb.New(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(1), "2024-12-01T00:00:00Z", "1969-12-31T23:59:59Z").
		Return(nil, errors.New("validation failed : failed to parse end date"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, err.Error(), "failed to parse end date")
}

func Test_FinanceHandler_GetIncomesForPeriod_ReturnsInvalidArgumentError_WhenStartDateAfterEndDate(t *testing.T) {
	ctrl, mockUsecase, handler := setupTest(t)
	defer ctrl.Finish()

	req := &finance.GetIncomesForPeriodRequest{
		UserId:    1,
		StartDate: timestamppb.New(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
		EndDate:   timestamppb.New(time.Date(2024, 12, 1, 23, 59, 59, 0, time.UTC)),
	}

	ctx := context.Background()
	mockUsecase.EXPECT().
		GetIncomesForPeriod(ctx, int64(1), "2024-12-31T00:00:00Z", "2024-12-01T23:59:59Z").
		Return(nil, errors.New("validation failed : start date is after end date"))

	resp, err := handler.GetIncomesForPeriod(ctx, req)

	require.Nil(t, resp)
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, err.Error(), "start date is after end date")
}
