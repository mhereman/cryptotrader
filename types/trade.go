package types

import "time"

type Trade struct {
	Symbol          Symbol
	ID              int64
	OrderID         int64
	Price           float64
	Quantity        float64
	QuoteQuantity   float64
	Commission      float64
	CommissionAsset string
	Time            time.Time
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}

func NewTrade(symbol Symbol, id, orderId int64, price, quantity, quoteQuantity, commission float64, commissionAsset string, tm time.Time, isBuyer, isMaker, isBestMacth bool) (t Trade) {
	return Trade{
		Symbol:          symbol,
		ID:              id,
		OrderID:         orderId,
		Price:           price,
		Quantity:        quantity,
		QuoteQuantity:   quoteQuantity,
		Commission:      commission,
		CommissionAsset: commissionAsset,
		Time:            tm,
		IsBuyer:         isBuyer,
		IsMaker:         isMaker,
		IsBestMatch:     isBestMacth,
	}
}
