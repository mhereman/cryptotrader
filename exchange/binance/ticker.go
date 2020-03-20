package binance

import (
	"context"
	"fmt"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b *Binance) Ticker(ctx context.Context, symbol types.Symbol) (price float64, err error) {
	var ts *bin.ListPricesService
	var response []*bin.SymbolPrice
	var binanceSymbol string

	if binanceSymbol, err = b.symbolToBinance(symbol); err != nil {
		logger.Errorf("Binance::Ticker Error: %v\n", err)
		return
	}

	ts = b.client.NewListPricesService()
	ts.Symbol(binanceSymbol)
	if response, err = ts.Do(ctx); err != nil {
		logger.Errorf("Binance::Ticker Error: %v\n", err)
		return
	}

	if len(response) != 1 {
		err = fmt.Errorf("Invalid ticker response, expected 1 value, received %d values", len(response))
		logger.Errorf("Binance::Ticker Error: %v\n", err)
		return
	}

	price = b.toFloat(response[0].Price)
	return
}
