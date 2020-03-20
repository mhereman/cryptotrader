package cryptotrader

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/mhereman/cryptotrader/algorithms"
	"github.com/mhereman/cryptotrader/logger"

	"github.com/mhereman/cryptotrader/exchange"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/types"
)

type AssetConfig struct {
	Symbol    types.Symbol
	Timeframe types.Timeframe
}

func NewAssetConfigFromFlags(base, quote, timeframe string) (ac AssetConfig, err error) {
	ac.Symbol = types.NewSymbol(base, quote)
	if ac.Timeframe, err = types.NewTimeframeFromString(timeframe); err != nil {
		logger.Errorf("Error parsing timeframe: %s %v\n", timeframe, err)
		return
	}
	return
}

type ExchangeConfig struct {
	Name   string
	ArgMap map[string]string
}

func NewExchangeConfigFromFlags(name string, args map[string]string) (ec ExchangeConfig, err error) {
	ec.Name = strings.ToLower(name)
	ec.ArgMap = args
	return
}

type AlgorithmConfig struct {
	Name   string
	Config types.AlgorithmConfig
}

func NewAlgorithmConfigFromFlags(name string, args map[string]string) (ac AlgorithmConfig, err error) {
	ac.Name = name
	ac.Config = types.AlgorithmConfig(args)
	return
}

type TradeVolumeType int

const (
	TVTFixed TradeVolumeType = iota
	TVTPercent
)

func TradeVolumeTypeFromString(in string) (out TradeVolumeType, err error) {
	switch strings.ToLower(in) {
	case "fixed":
		out = TVTFixed
	case "pct", "percent":
		out = TVTPercent
	default:
		err = fmt.Errorf("Invalid tradevolume type: %s", in)
	}
	return
}

func (tvt TradeVolumeType) String() string {
	if tvt == TVTFixed {
		return "fixed"
	}
	return "percent"
}

type TradeConfig struct {
	TradeVolumeType TradeVolumeType
	Volume          float64
	Reduce          bool
	Paper           bool
}

func NewTradeConfigFromFlags(tvt string, volume float64, reduce bool, paper bool) (tc TradeConfig, err error) {
	if tc.TradeVolumeType, err = TradeVolumeTypeFromString(tvt); err != nil {
		return
	}
	if tc.TradeVolumeType == TVTPercent {
		if volume > 1.0 {
			volume = 1.0
		}
	}
	tc.Volume = volume
	tc.Reduce = reduce
	tc.Paper = paper
	return
}

type CryptoTrader struct {
	ctx            context.Context
	cancelFn       context.CancelFunc
	wg             *sync.WaitGroup
	signalChannel  types.SignalChannel
	assetCfg       AssetConfig
	exchangeCfg    ExchangeConfig
	algoCfg        AlgorithmConfig
	tradeCfg       TradeConfig
	openTrades     map[string]string
	exchangeDriver interfaces.IExchangeDriver
	dataFetcher    interfaces.IDataFether
	algorithm      interfaces.IAlgorithm
}

func New(assetConfig AssetConfig, exchangeConfig ExchangeConfig, algorithmConfig AlgorithmConfig, tradeConfig TradeConfig) (ct *CryptoTrader) {
	ct = new(CryptoTrader)
	ct.ctx, ct.cancelFn = context.WithCancel(context.Background())
	ct.wg = &sync.WaitGroup{}
	ct.signalChannel = make(types.SignalChannel)
	ct.assetCfg = assetConfig
	ct.exchangeCfg = exchangeConfig
	ct.algoCfg = algorithmConfig
	ct.tradeCfg = tradeConfig
	ct.openTrades = make(map[string]string)
	return
}

func (ct *CryptoTrader) Run() (err error) {
	var signal types.Signal
	var accountInfo types.AccountInfo
	var balance types.AccountBalance
	var seriesChannel types.SeriesChannel
	var mainLoop bool

	defer func() {
		if r := recover(); r != nil {
			logger.Printf("Panic: %v\n", r)
			ct.cancelFn()
		}
	}()
	defer func() {
		if err != nil {
			ct.cancelFn()
		}
	}()

	setupCloseHandler(ct.ctx, ct.cancelFn)

	if err = ct.initExchangeDriver(); err != nil {
		return
	}

	if seriesChannel, err = ct.initDataFetcher(); err != nil {
		return
	}

	if err = ct.initAlgorithm(); err != nil {
		return
	}

	if accountInfo, err = ct.exchangeDriver.GetAccountInfo(ct.ctx); err != nil {
		logger.Errorf("Failed to retrieve account info: %v\n", err)
		return
	}
	logger.Infof("Maker commission:  %f\n", accountInfo.MakerCommission)
	logger.Infof("Taker commission:  %f\n", accountInfo.TakerCommission)
	logger.Infof("Buyer commission:  %f\n", accountInfo.BuyerCommission)
	logger.Infof("Seller commission: %f\n", accountInfo.SellerCommission)
	logger.Infof("Balances:")
	for _, balance = range accountInfo.Balances {
		if balance.Free > 0 || balance.Locked > 0 {
			logger.Infof(" . %s: Free: %f  Locked: %f\n", balance.Asset, balance.Free, balance.Locked)
		} else {
			logger.Debugf(" . %s: Free: %f  Locked: %f\n", balance.Asset, balance.Free, balance.Locked)
		}
	}

	ct.dataFetcher.RunAsync(ct.ctx, ct.wg)
	ct.algorithm.RunAsync(ct.ctx, ct.algoCfg.Config, seriesChannel, ct.signalChannel, ct.wg)

	if ct.tradeCfg.Paper {
		logger.Infoln("Cryptotrader running in paper trading mode")
	} else {
		logger.Infoln("Cryptotrader running in live mode")
	}
	logger.Infof(" . Symbol: %s[%s], Algorithm: %s, Ordersize: %s: %f", ct.assetCfg.Symbol.String(), ct.assetCfg.Timeframe.String(), ct.algoCfg.Name, ct.tradeCfg.TradeVolumeType.String(), ct.tradeCfg.Volume)
	showDonations()

	mainLoop = true
	for mainLoop {
		select {
		case <-ct.ctx.Done():
			logger.Infoln("Crptotrader shutting down")
			mainLoop = false
			break
		case signal = <-ct.signalChannel:
			if signal.Side == types.Sell {
				logger.Debugln("SELL")
			} else {
				logger.Debugln("BUY")
			}
			logger.Infof("Received signal: %s\n", signal.String())
			ct.executeSignal(signal)
		}
	}

	ct.wg.Done()
	return
}

func (ct *CryptoTrader) initExchangeDriver() (err error) {
	if ct.exchangeDriver, err = exchange.GetExchange(ct.ctx, ct.exchangeCfg.Name, ct.exchangeCfg.ArgMap); err != nil {
		logger.Errorf("Error configuring exchange: %v\n", err)
		return
	}
	logger.Infof("Exchange '%s' initialized\n", ct.exchangeCfg.Name)
	return
}

func (ct *CryptoTrader) initDataFetcher() (seriesChannel types.SeriesChannel, err error) {
	ct.dataFetcher = NewDataFetcher(ct.exchangeDriver)
	if seriesChannel, err = ct.dataFetcher.Register(ct.ctx, ct.assetCfg.Symbol, ct.assetCfg.Timeframe); err != nil {
		logger.Errorf("Error registering asset: %s[%s] %v\n", ct.assetCfg.Symbol.String(), ct.assetCfg.Timeframe.String(), err)
		return
	}
	logger.Infof("DataCacher initialized for asset %s[%s]\n", ct.assetCfg.Symbol.String(), ct.assetCfg.Timeframe.String())
	return
}

func (ct *CryptoTrader) initAlgorithm() (err error) {
	if ct.algorithm, err = algorithms.GetAlgorithm(ct.algoCfg.Name); err != nil {
		logger.Errorf("Error configuring algorithm: %v\n", err)
		return
	}
	logger.Infof("Algorithm '%s' initialized\n", ct.algoCfg.Name)
	return
}

func (ct *CryptoTrader) executeSignal(signal types.Signal) (err error) {
	var accountInfo types.AccountInfo
	var ok bool

	if accountInfo, err = ct.exchangeDriver.GetAccountInfo(ct.ctx); err != nil {
		logger.Errorf("Execute Signal Error: %v\n", err)
		return
	}

	if signal.Side == types.Buy {
		if _, ok = ct.openTrades[signal.Symbol.String()]; !ok {
			if ct.tradeCfg.TradeVolumeType == TVTFixed {
				err = ct.buyFixed(accountInfo, signal.Symbol)
			} else {
				err = ct.buyPercent(accountInfo, signal.Symbol)
			}
		}
	} else {
		if _, ok = ct.openTrades[signal.Symbol.String()]; ok {
			err = ct.closePosition(accountInfo, signal.Symbol)
		}
	}

	if err != nil {
		logger.Errorf("Execute Signal Error: %v\n", err)
	}
	return
}

func (ct *CryptoTrader) buyFixed(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var freeQuantity, orderQuantity, baseQuantity, averagePrice float64
	var orderBook types.OrderBook
	var orderInfo types.OrderInfo
	var ok bool

	if _, ok = ct.openTrades[symbol.String()]; ok {
		logger.Warningf("BuyFixed: Trade still open for symbol: %s\n", symbol.String())
		return
	}

	orderQuantity = ct.tradeCfg.Volume
	if ct.tradeCfg.Paper == false {
		freeQuantity, _ = accountInfo.GetAssetQuantity(symbol.Quote())
		if (freeQuantity * 0.995) < orderQuantity {
			if ct.tradeCfg.Reduce {
				orderQuantity = freeQuantity * 0.995
			} else {
				logger.Warningf("Insufficient funds to initiate buy position for symbol: %s (required: %f, available: %f)", symbol.String(), orderQuantity, (freeQuantity * 0.995))
				return
			}
		}
	}

	if orderBook, err = ct.exchangeDriver.GetOrderBook(ct.ctx, symbol); err != nil {
		logger.Errorf("BuyFixed Failed to retrieve order book for symbol: %s %v\n", symbol.String(), err)
		return
	}
	baseQuantity, averagePrice = orderBook.GetBuyVolumeAndAveragePrice(orderQuantity)
	baseQuantity = math.Floor(baseQuantity*1000000) / 1000000

	if ct.tradeCfg.Paper {
		ct.openTrades[symbol.String()] = orderInfo.Uuid.String()
		logger.Printf("PaperTrade: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbol.String(), baseQuantity, averagePrice, ct.openTrades[symbol.String()])
		return
	}

	if orderInfo, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewMarketOrder(symbol, types.Buy, baseQuantity)); err != nil {
		logger.Warningf("BuyFixed Error: %v\n", err)
		logger.Printf("Unable to initiate buy position for symbol: %s", symbol.String())
		return
	}

	ct.openTrades[symbol.String()] = orderInfo.Uuid.String()

	return
}

func (ct *CryptoTrader) buyPercent(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var freeQuantity, orderQuantity, baseQuantity, averagePrice float64
	var orderBook types.OrderBook
	var orderInfo types.OrderInfo
	var ok bool

	if _, ok = ct.openTrades[symbol.String()]; ok {
		logger.Warningf("BuyFixed: Trade still open for symbol: %s\n", symbol.String())
		return
	}

	freeQuantity, _ = accountInfo.GetAssetQuantity(symbol.Quote())
	orderQuantity = freeQuantity * ct.tradeCfg.Volume * 0.995

	if orderBook, err = ct.exchangeDriver.GetOrderBook(ct.ctx, symbol); err != nil {
		logger.Errorf("BuyPercent Failed to retrieve order book for symbol: %s %v\n", symbol.String(), err)
		return
	}
	baseQuantity, averagePrice = orderBook.GetBuyVolumeAndAveragePrice(orderQuantity)
	baseQuantity = math.Floor(baseQuantity*1000000) / 1000000

	if ct.tradeCfg.Paper {
		ct.openTrades[symbol.String()] = orderInfo.Uuid.String()
		logger.Printf("PaperTrade: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbol.String(), baseQuantity, averagePrice, ct.openTrades[symbol.String()])
		return
	}

	if orderInfo, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewMarketOrder(symbol, types.Buy, baseQuantity)); err != nil {
		logger.Warningf("BuyPercent Error: %v\n", err)
		logger.Printf("Unable to initiate buy position for symbol: %s", symbol.String())
		return
	}

	ct.openTrades[symbol.String()] = orderInfo.Uuid.String()

	return
}

func (ct *CryptoTrader) closePosition(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var origTradeID string
	var origTradeUUID uuid.UUID
	var orderInfo types.OrderInfo
	var trades []types.Trade
	var trade types.Trade
	var baseQuantity, price float64
	var ok bool

	if origTradeID, ok = ct.openTrades[symbol.String()]; !ok {
		logger.Warningf("ClosePosition: No trade to close for symbol: %s\n", symbol.String())
		return
	}
	if origTradeUUID, err = uuid.Parse(origTradeID); err != nil {
		logger.Errorf("ClosePosition: Invalid original trade uuid: %s %v\n", origTradeID, err)
		return
	}

	if ct.tradeCfg.Paper {
		ct.openTrades[symbol.String()] = orderInfo.Uuid.String()

		if price, err = ct.exchangeDriver.Ticker(ct.ctx, symbol); err != nil {
			logger.Errorf("ClosePosition: Failed to retrieve market price for symbol: %s %v\n", err)
			return
		}

		logger.Printf("PaperTrade: Sell %s [Market Price: %f; UUID: %s]", symbol.String(), baseQuantity, price, ct.openTrades[symbol.String()])
		delete(ct.openTrades, symbol.String())
		return
	}

	if orderInfo, err = ct.exchangeDriver.GetOrder(ct.ctx, types.Order{
		Symbol: symbol,
		Uuid:   origTradeUUID,
	}); err != nil {
		logger.Errorf("ClosePosition Error: %v\n", err)
		return
	}

	if trades, err = ct.exchangeDriver.GetOrderTrades(ct.ctx, orderInfo); err != nil {
		logger.Errorf("ClosePosition Error: %v\n", err)
		return
	}

	for _, trade = range trades {
		baseQuantity += trade.Quantity
		if trade.CommissionAsset == symbol.Base() {
			baseQuantity -= trade.Commission
		}
	}
	baseQuantity = math.Floor(baseQuantity*1000000) / 1000000

	logger.Infof("Close position: sell quantity: %f\n", baseQuantity)
	if _, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewMarketOrder(symbol, types.Sell, baseQuantity)); err != nil {
		logger.Warningf("ClosePosition Error: %v\n", err)
		logger.Printf("Unable to close position for symbol: %s", symbol.String())
	}

	delete(ct.openTrades, symbol.String())

	return
}
