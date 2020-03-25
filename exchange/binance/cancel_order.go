package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// CancelOrder executes the cancel order request
func (b Binance) CancelOrder(ctx context.Context, order types.Order, newUUID uuid.UUID) (info types.OrderInfo, err error) {
	var cos *bin.CancelOrderService
	var response *bin.CancelOrderResponse
	var binanceSymbol string

	if binanceSymbol, err = b.symbolToBinance(order.Symbol); err != nil {
		logger.Errorf("Binance::CancelOrder Error %v\n", err)
		return
	}

	cos = b.client.NewCancelOrderService()
	cos.Symbol(binanceSymbol)
	cos.OrigClientOrderID(order.UserReference.String())
	cos.NewClientOrderID(newUUID.String())
	if response, err = cos.Do(ctx); err != nil {
		logger.Errorf("Binance::CancelOrder Error %v\n", err)
		return
	}

	if info.UserReference, err = uuid.Parse(response.OrigClientOrderID); err != nil {
		logger.Errorf("Binance::CancelOrder Error %v\n", err)
		return
	}
	if info.CancelUserReference, err = uuid.Parse(response.ClientOrderID); err != nil {
		logger.Errorf("Binance::CancelOrder Error %v\n", err)
		return
	}
	info.ExchangeOrderID = response.OrderID
	if info.Symbol, err = b.toSymbol(response.Symbol); err != nil {
		logger.Errorf("Binance::CancelOrder Error %v\n", err)
		return
	}
	info.TransactionTime = b.toTime(response.TransactTime)
	info.OriginalQuantity = b.toFloat(response.OrigQuantity)
	info.ExecutedQuantity = b.toFloat(response.ExecutedQuantity)
	info.Price = b.toFloat(response.Price)
	info.Status = b.toStatus(response.Status)
	info.TimeInForce = b.toTimeInForce(response.TimeInForce)
	info.OrderType = b.toOrderType(response.Type)
	info.Side = b.toSide(response.Side)

	return
}
