package types

import "github.com/google/uuid"

// Order represents an order on the exchange
type Order struct {
	// UserReference user reference of the order
	UserReference uuid.UUID

	// Symbol of the order
	Symbol Symbol

	// Side of the order
	Side Side

	// Type of the order
	Type OrderType

	// TimeInForce of the order
	TimeInForce TimeInForce

	// Quantiy in base asset of the order
	Quantity float64

	// Price in quote asset of the order
	Price float64

	// StopPrice in quote asset of the order
	StopPrice float64
}

// NewLimitOrder creates a new Limit order instance
func NewLimitOrder(symbol Symbol, side Side, timeInForce TimeInForce, quantity float64, price float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = Limit
	o.TimeInForce = timeInForce
	o.Quantity = quantity
	o.Price = price
	return
}

// NewMarketOrder creates a new Market order instance
func NewMarketOrder(symbol Symbol, side Side, quantity float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = Market
	o.Quantity = quantity
	return
}

// NewStopLossOrder creates a new StopLoss order instance
func NewStopLossOrder(symbol Symbol, side Side, quantity float64, stopPrice float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = StopLoss
	o.Quantity = quantity
	o.StopPrice = stopPrice
	return o
}

// NewStopLossLimitOrder creates a new StopLossLimit order instance
func NewStopLossLimitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64, price float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = StopLossLimit
	o.Quantity = quantity
	o.Price = price
	o.StopPrice = stopPrice
	return o
}

// NewTakeProfitOrder creates a new TakeProfit order instance
func NewTakeProfitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = TakeProfit
	o.Quantity = quantity
	o.StopPrice = stopPrice
	return o
}

// NewTakeProfitLimitOrder creates a new TakeProfitLimit order instance
func NewTakeProfitLimitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64, price float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = TakeProfitLimit
	o.Quantity = quantity
	o.Price = price
	o.StopPrice = stopPrice
	return o
}

// NewLimitMakerOrder creates a new LimitMaker order instance
func NewLimitMakerOrder(symbol Symbol, side Side, timeInForce TimeInForce, quantity float64, price float64) (o Order) {
	o.UserReference = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = LimitMaker
	o.TimeInForce = timeInForce
	o.Quantity = quantity
	o.Price = price
	return
}
