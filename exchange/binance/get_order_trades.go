package binance

import (
	"context"
	"strings"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) GetOrderTrades(ctx context.Context, orderInfo types.OrderInfo) (trades []types.Trade, err error) {
	var gots *bin.ListTradesService
	var response []*bin.TradeV3
	var trade *bin.TradeV3
	var binanceSymbol string

	if binanceSymbol, err = b.symbolToBinance(orderInfo.Symbol); err != nil {
		logger.Errorf("Binance:GetOrderTrades Error %v\n", err)
		return
	}

	gots = b.client.NewListTradesService()
	gots.Symbol(binanceSymbol)
	gots.StartTime(orderInfo.TransactionTime.UnixNano() / 1000000)
	if response, err = gots.Do(ctx); err != nil {
		logger.Errorf("Binance:GetOrderTrades Error %v\n", err)
		return
	}

	trades = make([]types.Trade, 0, 0)
	for _, trade = range response {
		if trade.OrderID == orderInfo.ExchangeOrderID {
			var symbol types.Symbol
			if symbol, err = b.toSymbol(trade.Symbol); err != nil {
				logger.Errorf("Binance::GetOrderTrades Error %v\n", err)
				return
			}
			trades = append(trades, types.NewTrade(
				symbol,
				trade.ID,
				trade.OrderID,
				b.toFloat(trade.Price),
				b.toFloat(trade.Quantity),
				b.toFloat(trade.QuoteQuantity),
				b.toFloat(trade.Commission),
				strings.ToUpper(trade.CommissionAsset),
				b.toTime(trade.Time),
				trade.IsBuyer,
				trade.IsMaker,
				trade.IsBestMatch,
			))
		}
	}
	return
}
