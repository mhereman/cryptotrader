package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) GetSeries(ctx context.Context, symbol types.Symbol, timeframe types.Timeframe) (series types.Series, err error) {
	var ks *bin.KlinesService
	var binanceSymbol, binanceInterval string
	var response []*bin.Kline
	var numKlines, index int
	var kline *bin.Kline
	var ohlc []types.OHLC

	if binanceSymbol, err = b.symbolToBinance(symbol); err != nil {
		logger.Errorf("Binance::GetSeries Error: %v\n", err)
		return
	}

	if binanceInterval, err = b.timeframeToBinance(timeframe); err != nil {
		logger.Errorf("Binance:GetSeries Error: %v\n", err)
		return
	}

	ks = b.client.NewKlinesService()
	ks.Symbol(binanceSymbol)
	ks.Interval(binanceInterval)
	if response, err = ks.Do(ctx); err != nil {
		logger.Errorf("Binance::GetSeries Error: %v\n", err)
		return
	}

	numKlines = len(response)
	ohlc = make([]types.OHLC, numKlines, numKlines)
	for index, kline = range response {
		ohlc[index] = types.NewOHLC(
			b.toFloat(kline.Open),
			b.toFloat(kline.High),
			b.toFloat(kline.Low),
			b.toFloat(kline.Close),
			b.toFloat(kline.Volume),
			b.toTime(kline.OpenTime),
			b.toTime(kline.CloseTime),
		)
	}

	series = types.NewSeries(symbol, timeframe, ohlc)
	return
}
