package binance

import (
	"context"
	"fmt"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) GetSymbolInfo(ctx context.Context, symbol types.Symbol) (info types.SymbolInfo, err error) {
	var eis *bin.ExchangeInfoService
	var response *bin.ExchangeInfo
	var symInfo bin.Symbol
	var BinanceSymbol string

	if BinanceSymbol, err = b.symbolToBinance(symbol); err != nil {
		logger.Errorf("Binance::GetSymbolInfo Error %v\n", err)
		return
	}

	eis = b.client.NewExchangeInfoService()
	if response, err = eis.Do(ctx); err != nil {
		logger.Errorf("Binance::GetSymbolInfo Error %v\n", err)
		return
	}

	for _, symInfo = range response.Symbols {
		if symInfo.Symbol == BinanceSymbol {
			info.Symbol = symbol
			info.MinPrice = symInfo.PriceFilter().MinPrice
			info.MinLotQuantity = symInfo.LotSizeFilter().MinQuantity
			return
		}
	}

	err = fmt.Errorf("Failed to find symbol: %s\n", BinanceSymbol)
	logger.Errorf("Binance::GetSymbolInfo Error %v\n", err)
	return
}
