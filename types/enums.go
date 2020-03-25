package types

// Side of an order
type Side int

const (
	// Buy order
	Buy Side = iota

	// Sell order
	Sell
)

// String returns the string name of the Side
func (s Side) String() string {
	if int(s) == int(Buy) {
		return "Buy"
	}
	return "Sell"
}

// OrderType of an order
type OrderType int

const (
	// Limit order
	Limit OrderType = iota

	// Market order
	Market

	// StopLoss order
	StopLoss

	// StopLossLimit order
	StopLossLimit

	// TakeProfit order
	TakeProfit

	// TakeProfitLimit order
	TakeProfitLimit

	// LimitMaker order
	LimitMaker
)

// OrderStatus of an order
type OrderStatus int

const (
	// StatusNew ...
	StatusNew OrderStatus = iota

	// StatusPartiallyFilled ...
	StatusPartiallyFilled

	// StatusFilled ...
	StatusFilled

	// StatusCanceled ...
	StatusCanceled

	// StatusPendingCancel ...
	StatusPendingCancel

	// StatusRejected ...
	StatusRejected

	// StatusExpired ...
	StatusExpired
)

// TimeInForce of an order
type TimeInForce int

const (
	// GoodTillCancel ...
	GoodTillCancel TimeInForce = iota

	// ImmediateOrCancel ...
	ImmediateOrCancel

	// FillOrCancel ...
	FillOrCancel
)
