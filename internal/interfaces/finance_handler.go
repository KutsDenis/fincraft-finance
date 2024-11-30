package interfaces

import (
	"context"
	"fincraft-finance/api/finance"
	"fincraft-finance/internal/usecases"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
