package usecase

import (
    "context"

    "github.com/evrone/go-clean-template/internal/entity"
    repo "github.com/evrone/go-clean-template/internal/repo/webapi"
)

type TradeStream struct {
    repo repo.FinnhubRepo
}

func NewTradeStream(repo repo.FinnhubRepo) *TradeStream {
    return &TradeStream{repo: repo}
}

func (uc *TradeStream) SubscribeToTrades(ctx context.Context, symbol string) (<-chan entity.TradeDetails, error) {
    return uc.repo.Subscribe(ctx, symbol)
}

func (uc *TradeStream) UnsubscribeFromTrades(symbol string) error {
    return uc.repo.Unsubscribe(symbol)
}