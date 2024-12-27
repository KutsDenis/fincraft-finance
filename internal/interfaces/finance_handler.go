package interfaces

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"fincraft-finance/api/finance"
	"fincraft-finance/internal/usecases"
)

// FinanceHandler обрабатывает запросы к сервису финансов
type FinanceHandler struct {
	finance.UnimplementedFinanceServiceServer
	usecase usecases.IncomeService
}

// NewFinanceHandler создает новый экземпляр FinanceHandler
func NewFinanceHandler(usecase usecases.IncomeService) *FinanceHandler {
	return &FinanceHandler{usecase: usecase}
}

// AddIncome добавляет доход
func (h *FinanceHandler) AddIncome(ctx context.Context, req *finance.AddIncomeRequest) (*emptypb.Empty, error) {
	err := h.usecase.AddIncome(ctx, req.UserId, req.CategoryId, req.Amount, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add income: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetIncomesForPeriod возвращает список доходов за указанный период по категориям
func (h *FinanceHandler) GetIncomesForPeriod(ctx context.Context, req *finance.GetIncomesForPeriodRequest) (*finance.GetIncomesForPeriodResponse, error) {
	incomes, err := h.usecase.GetIncomesForPeriod(ctx, req.UserId, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "validation failed"):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case strings.Contains(err.Error(), "does not exist"):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "failed to get incomes: %v", err)
		}
	}

	return &finance.GetIncomesForPeriodResponse{Categories: incomes}, nil
}
