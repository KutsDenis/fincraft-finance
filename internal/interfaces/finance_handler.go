package interfaces

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"fincraft-finance/api/finance"
	"fincraft-finance/internal/domain"
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
	incomes, err := h.usecase.GetIncomesForPeriod(ctx, req.UserId, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		if errors.Is(err, usecases.ValidationError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, usecases.NotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to get incomes: %v", err)
	}

	return &finance.GetIncomesForPeriodResponse{Incomes: mapIncomes(incomes)}, nil
}

func mapIncomes(incomes []domain.Income) []*finance.Income {
	incomeResponses := make([]*finance.Income, 0, len(incomes))
	for _, income := range incomes {
		incomeResponses = append(incomeResponses, &finance.Income{
			CategoryId:  int32(income.CategoryID),
			Amount:      int64(income.Amount),
			Description: income.Description,
		})

	}
	return incomeResponses
}
