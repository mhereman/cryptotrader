package types

import (
	"time"

	"github.com/google/uuid"
)

type OrderInfo struct {
	Uuid             uuid.UUID
	CancelUuid       uuid.UUID
	ExchangeOrderID  int64
	Symbol           Symbol
	TransactionTime  time.Time
	OriginalQuantity float64
	ExecutedQuantity float64
	Price            float64
	StopPrice        float64
	Status           OrderStatus
	TimeInForce      TimeInForce
	OrderType        OrderType
	Side             Side
	Fills            []OrderFill
}

func NewOrderInfo(u uuid.UUID, eID int64, s Symbol, t time.Time, oq, eq, p, sp float64, st OrderStatus, timeInForce TimeInForce, ot OrderType, side Side) (info OrderInfo) {
	info.Uuid = u
	info.ExchangeOrderID = eID
	info.Symbol = s
	info.TransactionTime = t
	info.OriginalQuantity = oq
	info.ExecutedQuantity = eq
	info.Price = p
	info.StopPrice = sp
	info.Status = st
	info.TimeInForce = timeInForce
	info.OrderType = ot
	info.Side = side
	return
}

func (info *OrderInfo) AppendFills(fils ...OrderFill) {
	info.Fills = append(info.Fills, fils...)
}
