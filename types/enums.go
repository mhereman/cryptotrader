package types

type Side int

const (
	Buy Side = iota
	Sell
)

func (s Side) String() string {
	if int(s) == int(Buy) {
		return "Buy"
	}
	return "Sell"
}

type OrderType int

const (
	Limit OrderType = iota
	Market
	StopLoss
	StopLossLimit
	TakeProfit
	TakeProfitLimit
	LimitMaker
)

type OrderStatus int

const (
	StatusNew OrderStatus = iota
	StatusPartiallyFilled
	StatusFilled
	StatusCanceled
	StatusPendingCancel
	StatusRejected
	StatusExpired
)

type TimeInForce int

const (
	GoodTillCancel TimeInForce = iota
	ImmediateOrCancel
	FillOrCancel
)
