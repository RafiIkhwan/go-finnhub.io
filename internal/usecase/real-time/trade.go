package usecase

import (
	"context"
	"time"

	"github.com/evrone/go-clean-template/internal/entity"
)

type RealTimeTradeRepo interface {
	Create(ctx context.Context, trade entity.RealTimeTrade) error
	GetBySymbol(ctx context.Context, symbol string, start, end time.Time) ([]entity.RealTimeTrade, error)
}

type RealTimeTrade struct {
	repo RealTimeTradeRepo
}

func NewRealTimeTrade(repo RealTimeTradeRepo) *RealTimeTrade {
	return &RealTimeTrade{repo: repo}
}

func (uc *RealTimeTrade) CreateTrade(ctx context.Context, trade entity.RealTimeTrade) error {
	return uc.repo.Create(ctx, trade)
}

func (uc *RealTimeTrade) GetTrades(ctx context.Context, symbol string, start, end time.Time) ([]entity.RealTimeTrade, error) {
	return uc.repo.GetBySymbol(ctx, symbol, start, end)
}
