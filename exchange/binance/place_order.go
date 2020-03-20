package binance

import (
	"context"
	"fmt"

	bin "github.com/adshao/go-binance"
	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) PlaceOrder(ctx context.Context, order types.Order) (info types.OrderInfo, err error) {
	var cos *bin.CreateOrderService
	var binanceSymbol string
	var response *bin.CreateOrderResponse
	var numFills, index int
	var fill *bin.Fill
	var orderFills []types.OrderFill

	if binanceSymbol, err = b.symbolToBinance(order.Symbol); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}

	cos = b.client.NewCreateOrderService()
	cos.NewClientOrderID(order.Uuid.String())
	cos.Symbol(binanceSymbol)
	cos.Side(b.sideToBinance(order.Side))
	cos.Type(b.orderTypeToBinance(order.Type))

	switch order.Type {
	case types.Limit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.Price(fmt.Sprintf("%f", order.Price))
	case types.Market:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
	case types.StopLoss:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.StopPrice(fmt.Sprintf("%f", order.StopPrice))
	case types.StopLossLimit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.StopPrice(fmt.Sprintf("%f", order.StopPrice))
		cos.Price(fmt.Sprintf("%f", order.Price))
	case types.TakeProfit:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.StopPrice(fmt.Sprintf("%f", order.StopPrice))
	case types.TakeProfitLimit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.StopPrice(fmt.Sprintf("%f", order.StopPrice))
		cos.Price(fmt.Sprintf("%f", order.Price))
	case types.LimitMaker:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		cos.Price(fmt.Sprintf("%f", order.Price))
	}

	if response, err = cos.Do(ctx); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}

	if info.Uuid, err = uuid.Parse(response.ClientOrderID); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}

	info.ExchangeOrderID = response.OrderID
	if info.Symbol, err = b.toSymbol(response.Symbol); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}
	info.TransactionTime = b.toTime(response.TransactTime)
	info.OriginalQuantity = b.toFloat(response.OrigQuantity)
	info.ExecutedQuantity = b.toFloat(response.ExecutedQuantity)
	info.Status = b.toStatus(response.Status)
	info.TimeInForce = b.toTimeInForce(response.TimeInForce)
	info.OrderType = b.toOrderType(response.Type)
	info.Side = b.toSide(response.Side)

	numFills = len(response.Fills)
	if numFills > 0 {
		orderFills = make([]types.OrderFill, numFills, numFills)
		for index, fill = range response.Fills {
			orderFills[index] = types.NewOrderFill(
				b.toFloat(fill.Price),
				b.toFloat(fill.Quantity),
				b.toFloat(fill.Commission),
				fill.CommissionAsset,
			)
		}
		info.Fills = orderFills
	}

	return
}
