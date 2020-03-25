package types

import (
	"time"

	"github.com/google/uuid"
)

// OrderInfo represents information about an order on the exchange
type OrderInfo struct {
	// UserReference the user ref of the Order
	UserReference uuid.UUID

	// UUID of the cancel order
	CancelUserReference uuid.UUID

	// ExchangeOrderID is the exchange provided order id
	ExchangeOrderID int64

	// Symbol of the order
	Symbol Symbol

	// TransactinoTime of the order
	TransactionTime time.Time

	// OriginalQuantity of the order
	OriginalQuantity float64

	// ExecutedQuantity of the order
	ExecutedQuantity float64

	// Price of the order
	Price float64

	// StopPrice of the order
	StopPrice float64

	// Status of the order
	Status OrderStatus

	// TimeInForce of the order
	TimeInForce TimeInForce

	// Type of the order
	OrderType OrderType

	// Side of the order
	Side Side

	// Fills of the order
	Fills []OrderFill
}

// NewOrderInfo creates a new OrderInfo instance
func NewOrderInfo(userRef uuid.UUID, exchangeOrderID int64, symbol Symbol, transactionTime time.Time, originalQuantity float64, executedQuantity float64, price float64, stopPrice float64, status OrderStatus, timeInForce TimeInForce, orderType OrderType, side Side) (info OrderInfo) {
	info.UserReference = userRef
	info.ExchangeOrderID = exchangeOrderID
	info.Symbol = symbol
	info.TransactionTime = transactionTime
	info.OriginalQuantity = originalQuantity
	info.ExecutedQuantity = executedQuantity
	info.Price = price
	info.StopPrice = stopPrice
	info.Status = status
	info.TimeInForce = timeInForce
	info.OrderType = orderType
	info.Side = side
	return
}

// AppendFills adds fills to the order info
func (info *OrderInfo) AppendFills(fils ...OrderFill) {
	info.Fills = append(info.Fills, fils...)
}
