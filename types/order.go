package types

import "github.com/google/uuid"

type Order struct {
	Uuid        uuid.UUID
	Symbol      Symbol
	Side        Side
	Type        OrderType
	TimeInForce TimeInForce
	Quantity    float64
	Price       float64
	StopPrice   float64
}

func NewLimitOrder(symbol Symbol, side Side, timeInForce TimeInForce, quantity float64, price float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = Limit
	o.TimeInForce = timeInForce
	o.Quantity = quantity
	o.Price = price
	return
}

func NewMarketOrder(symbol Symbol, side Side, quantity float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = Market
	o.Quantity = quantity
	return
}

func NewStopLossOrder(symbol Symbol, side Side, quantity float64, stopPrice float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = StopLoss
	o.Quantity = quantity
	o.StopPrice = stopPrice
	return o
}

func NewStopLossLimitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64, price float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = StopLossLimit
	o.Quantity = quantity
	o.Price = price
	o.StopPrice = stopPrice
	return o
}

func NewTakeProfitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = TakeProfit
	o.Quantity = quantity
	o.StopPrice = stopPrice
	return o
}

func NewTakeProfitLimitOrder(symbol Symbol, side Side, quantity float64, stopPrice float64, price float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = TakeProfitLimit
	o.Quantity = quantity
	o.Price = price
	o.StopPrice = stopPrice
	return o
}

func NewLimitMakerOrder(symbol Symbol, side Side, timeInForce TimeInForce, quantity float64, price float64) (o Order) {
	o.Uuid = uuid.New()
	o.Symbol = symbol
	o.Side = side
	o.Type = LimitMaker
	o.TimeInForce = timeInForce
	o.Quantity = quantity
	o.Price = price
	return
}
