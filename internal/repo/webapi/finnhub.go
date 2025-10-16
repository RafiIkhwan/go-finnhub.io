package webapi

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

type FinnhubRepo interface {
	Subscribe(ctx context.Context, symbol string) (<-chan entity.TradeDetails, error)
	Unsubscribe(symbol string) error
	Close() error
}