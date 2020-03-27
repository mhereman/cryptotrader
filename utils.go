package cryptotrader

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mhereman/cryptotrader/logger"
)

func ReadFlags() (assetCfg AssetConfig, exchangeCfg ExchangeConfig, algoConfig AlgorithmConfig, tradeConfig TradeConfig, notifierConfig NotifierConfig, err error) {
	var base, quote, timeFrame, exchange, exchangeArgsString, algo, algoConfigString, tradeType, logLevel, notifier, notifierConfigString *string
	var volume, maxSlippage *float64
	var reduce, paperTrading *bool

	base = flag.String("base", "btc", "Base asset to trade")
	quote = flag.String("quote", "usdt", "Quote asset to trade")
	timeFrame = flag.String("timeframe", "4h", "Timeframe to trade, unit in ['s', 'm', 'h', 'd', 'w', 'M']")

	exchange = flag.String("exchange", "binance", "Exchange to trade on, valid exchanges: ['binance']")
	exchangeArgsString = flag.String("exchangeargs", "apiKey=abc;apiSecret=def", "Exchange arguments, e.g. apiKey, apiSecret, ...")

	algo = flag.String("algo", "Ema/Sma", "Algorithm to trade, valid algorithms: ['Ema/Sma']")
	algoConfigString = flag.String("algoargs", "Key=value;Key2=value2", "Algorithm arguments")

	tradeType = flag.String("tradetype", "pct", "How to calculate trade volume, valid: ['pct', 'fixed']")
	volume = flag.Float64("volume", 1.0, "Trade volume. If tradetype = pct, the volume is the percentage of the availabel quote asset, otherwise the fixed volume of the trade asset.")
	reduce = flag.Bool("reduce", true, "Reduce the trade volume if not sufficient funds are available")
	paperTrading = flag.Bool("papertrading", false, "Papertrading enabled or not")
	maxSlippage = flag.Float64("maxslippage", 0.001, "Max slippage on buy orders; if set to 0 no max slippage is configured. Sell orders are always market orders.")

	logLevel = flag.String("loglevel", "info", "Log leve to use, valid (most verbose to less): ['debug', 'error', warning', 'info', 'none'")

	notifier = flag.String("notifier", "", "If set, the notifier service to use, valid notifiers: ['', 'proximus-sms']")
	notifierConfigString = flag.String("notifierargs", "Key=value;key2=value", "Notifier arguments")

	flag.Parse()

	logger.SetLogLevel(logger.NewLogLevelFromString(*logLevel))

	if assetCfg, err = NewAssetConfigFromFlags(*base, *quote, *timeFrame); err != nil {
		return
	}

	if exchangeCfg, err = NewExchangeConfigFromFlags(*exchange, buildArgMap(*exchangeArgsString)); err != nil {
		return
	}

	if algoConfig, err = NewAlgorithmConfigFromFlags(*algo, buildArgMap(*algoConfigString)); err != nil {
		return
	}

	if tradeConfig, err = NewTradeConfigFromFlags(*tradeType, *volume, *reduce, *paperTrading, *maxSlippage); err != nil {
		return
	}

	if notifierConfig, err = NewNotifierConfigFromFlags(*notifier, buildArgMap(*notifierConfigString)); err != nil {
		return
	}
	return
}

func buildArgMap(in string) (out map[string]string) {
	var parts, kv []string
	var part string

	out = make(map[string]string)

	parts = strings.Split(in, ";")
	for _, part = range parts {
		kv = strings.Split(part, "=")
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return
}

func setupCloseHandler(ctx context.Context, cancelFn context.CancelFunc) {
	var ch chan os.Signal
	ch = make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ch
		cancelFn()
	}()
}

func showDonations() {
	var donationOptions [][]string = [][]string{
		[]string{"Bitcoin (BTC) ", "bc1q07vep45su43azxzpd4x04a9f7sku7szua29ed6"},
		[]string{"Ethereum (ETH)", "0x0e4825331f704697c44729012ce2608493bcf60E"},
		[]string{"Litecoin (LTC)", "MREwCamJLXAu6gJGd11A8i9qYESbKRTkH7"},
		[]string{"Dash (DASH)   ", "XodChdVy5JpDvSmEfHHYZwLCRTQuKoJJCv"},
		[]string{"Zcash (ZEC)   ", "t1WFQRxKsEJiWdrLNfjGyKzFcBVvGQKMgKm"},
	}

	fmt.Println()
	fmt.Println("If you like Cryptotrader consider giving a donation to support the developers.")
	fmt.Println("Donations can be given by means of the following crypto assets and their corresponding addresses:")
	fmt.Println()
	for _, option := range donationOptions {
		fmt.Printf("  %s: %s\n", option[0], option[1])
	}
	fmt.Println()
}
