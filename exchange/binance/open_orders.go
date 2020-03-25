package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// OpenOrders executes the open orders request
func (b Binance) OpenOrders(ctx context.Context, symbol types.Symbol) (orders []types.OrderInfo, err error) {
	var loos *bin.ListOpenOrdersService
	var response []*bin.Order
	var binanceSymbol string
	var numOrders, index int
	var order *bin.Order
	var tmpUUID uuid.UUID
	var tmpSymbol types.Symbol

	if binanceSymbol, err = b.symbolToBinance(symbol); err != nil {
		logger.Errorf("Binance::OpenOrders Error %v\n", err)
		return
	}

	loos = b.client.NewListOpenOrdersService()
	loos.Symbol(binanceSymbol)
	if response, err = loos.Do(ctx); err != nil {
		logger.Errorf("Binance::OpenOrders Error %v\n", err)
		return
	}

	numOrders = len(response)
	orders = make([]types.OrderInfo, numOrders, numOrders)
	for index, order = range response {
		if tmpUUID, err = uuid.Parse(order.ClientOrderID); err != nil {
			logger.Errorf("Binance::OpenOrders Error %v\n", err)
			return
		}
		if tmpSymbol, err = b.toSymbol(order.Symbol); err != nil {
			logger.Errorf("Binance::OpenOrders Error %v\n", err)
			return
		}

		orders[index] = types.NewOrderInfo(
			tmpUUID,
			order.OrderID,
			tmpSymbol,
			b.toTime(order.Time),
			b.toFloat(order.OrigQuantity),
			b.toFloat(order.ExecutedQuantity),
			b.toFloat(order.Price),
			b.toFloat(order.StopPrice),
			b.toStatus(order.Status),
			b.toTimeInForce(order.TimeInForce),
			b.toOrderType(order.Type),
			b.toSide(order.Side),
		)
	}

	return
}
