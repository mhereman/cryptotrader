package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// GetOrder executes the get order request
func (b Binance) GetOrder(ctx context.Context, order types.Order) (info types.OrderInfo, err error) {
	var gos *bin.GetOrderService
	var response *bin.Order
	var binanceSymbol string

	if binanceSymbol, err = b.symbolToBinance(order.Symbol); err != nil {
		logger.Errorf("Binance::GetOrder Error %v\n", err)
		return
	}

	gos = b.client.NewGetOrderService()
	gos.Symbol(binanceSymbol)
	gos.OrigClientOrderID(order.UserReference.String())
	if response, err = gos.Do(ctx); err != nil {
		logger.Errorf("Binance::GetOrder Error %v\n", err)
		return
	}

	if info.UserReference, err = uuid.Parse(response.ClientOrderID); err != nil {
		logger.Errorf("Binance::GetOrder Error %v\n", err)
		return
	}
	info.ExchangeOrderID = response.OrderID
	if info.Symbol, err = b.toSymbol(response.Symbol); err != nil {
		logger.Errorf("Binance::GetOrder Error %v\n", err)
		return
	}
	info.TransactionTime = b.toTime(response.Time)
	info.OriginalQuantity = b.toFloat(response.OrigQuantity)
	info.ExecutedQuantity = b.toFloat(response.ExecutedQuantity)
	info.Price = b.toFloat(response.Price)
	info.StopPrice = b.toFloat(response.StopPrice)
	info.Status = b.toStatus(response.Status)
	info.TimeInForce = b.toTimeInForce(response.TimeInForce)
	info.OrderType = b.toOrderType(response.Type)
	info.Side = b.toSide(response.Side)

	return
}
