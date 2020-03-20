package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/types"
)

type IExchangeDriver interface {
	Name() string

	GetAccountInfo(context.Context) (types.AccountInfo, error)

	TestConnectivity(context.Context) (bool, error)
	GetServerTime(context.Context) (time.Time, error)
	GetOrderBook(context.Context, types.Symbol) (types.OrderBook, error)
	GetSeries(context.Context, types.Symbol, types.Timeframe) (types.Series, error)
	Ticker(context.Context, types.Symbol) (float64, error)

	PlaceOrder(context.Context, types.Order) (types.OrderInfo, error)
	GetOrder(context.Context, types.Order) (types.OrderInfo, error)
	CancelOrder(context.Context, types.Order, uuid.UUID) (types.OrderInfo, error)
	OpenOrders(context.Context, types.Symbol) ([]types.OrderInfo, error)
	GetOrderTrades(context.Context, types.OrderInfo) ([]types.Trade, error)
}
