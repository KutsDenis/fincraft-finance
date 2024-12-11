package interfaces

import (
	"context"
	"strings"
	"time"

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

// GetIncomesForPeriod возвращает список доходов за указанный период
func (h *FinanceHandler) GetIncomesForPeriod(ctx context.Context, req *finance.GetIncomesForPeriodRequest) (*finance.GetIncomesForPeriodResponse, error) {
	incomes, err := h.usecase.GetIncomesForPeriod(ctx, req.UserId, req.StartDate.AsTime().Format(time.RFC3339), req.EndDate.AsTime().Format(time.RFC3339))
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if strings.Contains(err.Error(), "does not exist") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get incomes: %v", err)
	}

	var incomeResponses []*finance.Income
	for _, income := range incomes {
		incomeResponses = append(incomeResponses, &finance.Income{
			UserId:      income.UserID,
			CategoryId:  int32(income.CategoryID),
			Amount:      income.Amount.ToFloat(),
			Description: income.Description,
		})
	}

	return &finance.GetIncomesForPeriodResponse{Incomes: incomeResponses}, nil
}
