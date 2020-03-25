package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/types"
)

// IExchangeDriver rerpresents the interface for an exchange plugin
type IExchangeDriver interface {
	// Name returns the name of the exchange plugin
	Name() string

	// GetAccountInfo executes the get account info request
	GetAccountInfo(context.Context) (types.AccountInfo, error)

	// TestConnectivity tests exchange connectivity
	TestConnectivity(context.Context) (bool, error)

	// GetServerTime executes the get server time request
	GetServerTime(context.Context) (time.Time, error)

	// GetOrderBook executes the get orderbook request
	GetOrderBook(context.Context, types.Symbol) (types.OrderBook, error)

	// GetSeries executes the get series request
	GetSeries(context.Context, types.Symbol, types.Timeframe) (types.Series, error)

	// Ticker executes the ticker request
	Ticker(context.Context, types.Symbol) (float64, error)

	// PlaceOrder executes the place order request
	PlaceOrder(context.Context, types.Order) (types.OrderInfo, error)

	// GetOrder executes the get order request
	GetOrder(context.Context, types.Order) (types.OrderInfo, error)

	// CancelOrder executes the cancel order requests
	CancelOrder(context.Context, types.Order, uuid.UUID) (types.OrderInfo, error)

	// OpenOrders executes the open orders request
	OpenOrders(context.Context, types.Symbol) ([]types.OrderInfo, error)

	// GetOrderTrades executs the get order trades request
	GetOrderTrades(context.Context, types.OrderInfo) ([]types.Trade, error)
}
