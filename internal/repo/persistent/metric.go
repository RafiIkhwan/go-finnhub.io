package persistent

import (
    "context"

    "github.com/Masterminds/squirrel"
    "github.com/evrone/go-clean-template/internal/entity"
    "github.com/evrone/go-clean-template/pkg/postgres"
)

type MetricRepo struct {
    db *postgres.Postgres
}

func NewMetricRepo(db *postgres.Postgres) *MetricRepo {
    return &MetricRepo{db: db}
}

func (r *MetricRepo) Create(ctx context.Context, metric entity.TradeDetails) error {
		// First insert into real_time_trades table
		tradeQuery, tradeArgs, err := squirrel.Insert("real_time_trades").
			Columns("type").
			Values("trade"). // or use metric.Type if available
			Suffix("RETURNING id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return err
		}
		
		var realTimeTradeID int
		err = r.db.Pool.QueryRow(ctx, tradeQuery, tradeArgs...).Scan(&realTimeTradeID)
		if err != nil {
			return err
		}
		
		// Then insert into trade_details table
		query, args, err := squirrel.Insert("trade_details").
			Columns("real_time_trade_id", "symbol", "price", "timestamp_ms", "volume", "conditions").
			Values(realTimeTradeID, metric.Symbol, metric.Price, metric.Timestamp, metric.Volume, metric.Conditions).
			PlaceholderFormat(squirrel.Dollar).
			ToSql()
		if err != nil {
			return err
		}

		return r.db.Pool.QueryRow(ctx, query, args...).Scan(&realTimeTradeID)
}