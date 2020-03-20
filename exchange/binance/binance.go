package binance

import (
	"context"
	"fmt"

	"github.com/mhereman/cryptotrader/exchange"
	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"

	bin "github.com/adshao/go-binance"
)

const exchangeName = "binance"

func init() {
	exchange.RegisterExchange(exchangeName, createBinance)
}

type Binance struct {
	client     *bin.Client
	allSymbols map[string][]string
}

func New(ctx context.Context, config map[string]string) (driver *Binance, err error) {
	var apiKey, apiSecret string
	var ok bool

	if apiKey, ok = config["apiKey"]; !ok {
		err = fmt.Errorf("Binance config error: 'apiKey' entry not found")
		return
	}

	if apiSecret, ok = config["apiSecret"]; !ok {
		err = fmt.Errorf("Binance config error: 'apiSecret' entry not found")
		return
	}

	driver = new(Binance)
	driver.client = bin.NewClient(apiKey, apiSecret)
	driver.allSymbols = make(map[string][]string)

	if err = driver.getAllSymbols(ctx); err != nil {
		logger.Errorf("Binance::New Error: %v\n", err)
		return
	}
	return
}

func createBinance(ctx context.Context, config map[string]string) (driver interfaces.IExchangeDriver, err error) {
	driver, err = New(ctx, config)
	return
}

func (b Binance) Name() string {
	return exchangeName
}

func (b *Binance) getAllSymbols(ctx context.Context) (err error) {
	var eis *bin.ExchangeInfoService
	var response *bin.ExchangeInfo
	var symbol bin.Symbol

	eis = b.client.NewExchangeInfoService()
	if response, err = eis.Do(ctx); err != nil {
		return
	}

	for _, symbol = range response.Symbols {
		b.allSymbols[symbol.Symbol] = []string{symbol.BaseAsset, symbol.QuoteAsset}
	}
	return
}
