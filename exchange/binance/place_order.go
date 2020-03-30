package binance

import (
	"context"
	"fmt"

	bin "github.com/adshao/go-binance"
	"github.com/google/uuid"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// PlaceOrder executes the place order request
func (b Binance) PlaceOrder(ctx context.Context, order types.Order, symbolInfo *types.SymbolInfo) (info types.OrderInfo, err error) {
	var cos *bin.CreateOrderService
	var binanceSymbol string
	var response *bin.CreateOrderResponse
	var numFills, index int
	var fill *bin.Fill
	var orderFills []types.OrderFill
	var strPrice string

	if binanceSymbol, err = b.symbolToBinance(order.Symbol); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}

	cos = b.client.NewCreateOrderService()
	cos.NewClientOrderID(order.UserReference.String())
	cos.Symbol(binanceSymbol)
	cos.Side(b.sideToBinance(order.Side))
	cos.Type(b.orderTypeToBinance(order.Type))

	switch order.Type {
	case types.Limit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.Price); err != nil {
			return
		}
		cos.Price(strPrice)
	case types.Market:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
	case types.StopLoss:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.StopPrice); err != nil {
			return
		}
		cos.StopPrice(strPrice)
	case types.StopLossLimit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.StopPrice); err != nil {
			return
		}
		cos.StopPrice(strPrice)
		if strPrice, err = symbolInfo.ClampPrice(order.Price); err != nil {
			return
		}
		cos.Price(strPrice)
	case types.TakeProfit:
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.StopPrice); err != nil {
			return
		}
		cos.StopPrice(strPrice)
	case types.TakeProfitLimit:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.StopPrice); err != nil {
			return
		}
		cos.StopPrice(strPrice)
		if strPrice, err = symbolInfo.ClampPrice(order.Price); err != nil {
			return
		}
		cos.Price(strPrice)
	case types.LimitMaker:
		cos.TimeInForce(b.timeInForceToBinance(order.TimeInForce))
		cos.Quantity(fmt.Sprintf("%f", order.Quantity))
		if strPrice, err = symbolInfo.ClampPrice(order.Price); err != nil {
			return
		}
		cos.Price(strPrice)
	}

	if response, err = cos.Do(ctx); err != nil {
		logger.Errorf("Binance::PlaceOrder Error %v\n", err)
		return
	}

	if info.UserReference, err = uuid.Parse(response.ClientOrderID); err != nil {
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
