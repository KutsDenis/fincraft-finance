package interfaces

import (
	"context"

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
	err := h.usecase.AddIncome(ctx, req.UserId, int(req.CategoryId), req.Amount, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add income: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetIncomeForPeriod возвращает общий доход за период
func (h *FinanceHandler) GetIncomeForPeriod(ctx context.Context, req *finance.GetIncomeForPeriodRequest) (*finance.GetIncomeForPeriodResponse, error) {
	total_income_for_period, err := h.usecase.GetIncomeForPeriod(ctx, req.UserId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get income for period: %v", err)
	}

	return &finance.GetIncomeForPeriodResponse{TotalIncome: total_income_for_period}, nil
}
