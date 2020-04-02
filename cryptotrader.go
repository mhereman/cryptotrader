package cryptotrader

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/mhereman/cryptotrader/algorithms"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/notifiers"

	"github.com/mhereman/cryptotrader/exchange"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/types"
)

type CryptoTrader struct {
	ctx            context.Context
	cancelFn       context.CancelFunc
	wg             *sync.WaitGroup
	signalChannel  types.SignalChannel
	assetCfg       AssetConfig
	exchangeCfg    ExchangeConfig
	algoCfg        AlgorithmConfig
	tradeCfg       TradeConfig
	notifierCfg    NotifierConfig
	openTrades     map[string]string
	exchangeDriver interfaces.IExchangeDriver
	dataFetcher    interfaces.IDataFetcher
	algorithm      interfaces.IAlgorithm
	notifier       interfaces.INotifier
	buyFn          func(types.AccountInfo, types.Symbol) error
	closeFn        func(types.AccountInfo, types.Symbol) error
	my_var         int
}

func New(assetConfig AssetConfig, exchangeConfig ExchangeConfig, algorithmConfig AlgorithmConfig, tradeConfig TradeConfig, notifierConfig NotifierConfig) (ct *CryptoTrader) {
	ct = new(CryptoTrader)
	ct.ctx, ct.cancelFn = context.WithCancel(context.Background())
	ct.wg = &sync.WaitGroup{}
	ct.signalChannel = make(types.SignalChannel)
	ct.assetCfg = assetConfig
	ct.exchangeCfg = exchangeConfig
	ct.algoCfg = algorithmConfig
	ct.tradeCfg = tradeConfig
	ct.notifierCfg = notifierConfig
	ct.openTrades = make(map[string]string)

	if ct.tradeCfg.Paper {
		if ct.tradeCfg.TradeVolumeType == TVTFixed {
			ct.buyFn = ct.paperBuyFixed
		} else {
			ct.buyFn = ct.paperBuyPercent
		}
		ct.closeFn = ct.paperClosePosition
	} else {
		if ct.tradeCfg.TradeVolumeType == TVTFixed {
			ct.buyFn = ct.liveBuyFixed
		} else {
			ct.buyFn = ct.liveBuyPercent
		}
		ct.closeFn = ct.liveClosePosition
	}
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

	if err = ct.initNotifier(); err != nil {
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
		ct.notifier.Notify(ct.ctx, []byte("Cryptotrader running in paper trading mode"))
	} else {
		logger.Infoln("Cryptotrader running in live mode")
		ct.notifier.Notify(ct.ctx, []byte("Cryptotrader running in live mode"))
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

func (ct *CryptoTrader) initNotifier() (err error) {
	if ct.notifierCfg.Name == "" {
		ct.notifierCfg.Name = "noop"
		ct.notifierCfg.ArgMap = make(map[string]string)
	}

	if ct.notifier, err = notifiers.GetNotifier(ct.ctx, ct.notifierCfg.Name, ct.notifierCfg.ArgMap); err != nil {
		logger.Errorf("Error configuring notifier: %v\n", err)
		return
	}
	logger.Infof("Notifier '%s' initialized\n", ct.notifierCfg.Name)
	return
}

func (ct *CryptoTrader) executeSignal(signal types.Signal) (err error) {
	var accountInfo types.AccountInfo

	if accountInfo, err = ct.exchangeDriver.GetAccountInfo(ct.ctx); err != nil {
		logger.Errorf("Execute Signal Error: %v\n", err)
		return
	}

	if signal.IsBacktest {
		logger.Infoln(signal.String())
		return
	}

	if signal.Side == types.Buy {
		err = ct.buyFn(accountInfo, signal.Symbol)
	} else {
		err = ct.closeFn(accountInfo, signal.Symbol)
	}

	if err != nil {
		logger.Errorf("Execute Signal Error: %v\n", err)
	}
	return
}

func (ct *CryptoTrader) buyMarket(symbol types.Symbol, orderQuantity float64) (orderInfo types.OrderInfo, averagePrice float64, err error) {
	var orderBook types.OrderBook
	var baseQuantity float64

	if orderBook, err = ct.exchangeDriver.GetOrderBook(ct.ctx, symbol); err != nil {
		err = fmt.Errorf("buyMarket Failed to retrieve order book for symbol: %s %v", symbol.String(), err)
		return
	}
	baseQuantity, averagePrice = orderBook.GetBuyVolumeAndAveragePrice(orderQuantity)
	baseQuantity = normalizeQuantity(baseQuantity)

	if orderInfo, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewMarketOrder(symbol, types.Buy, baseQuantity), nil); err != nil {
		err = fmt.Errorf("buyMarket Failed to place order for symbol: %s %v", symbol.String(), err)
		return
	}

	return
}

func (ct *CryptoTrader) buyLimit(symbol types.Symbol, orderQuantity float64) (orderInfo types.OrderInfo, limitPrice float64, err error) {
	var marketPrice, baseQuantity float64
	var symbolInfo types.SymbolInfo

	if marketPrice, err = ct.exchangeDriver.Ticker(ct.ctx, symbol); err != nil {
		err = fmt.Errorf("buyLimit Failed to retrieve market price for symbol %s %v", symbol.String(), err)
		return
	}

	if symbolInfo, err = ct.exchangeDriver.GetSymbolInfo(ct.ctx, symbol); err != nil {
		err = fmt.Errorf("buyLimit Failed to retrieve symbol info for symbol %s %v", symbol.String(), err)
		return
	}

	limitPrice = marketPrice * (1.0 + ct.tradeCfg.MaxSlippage)
	baseQuantity = orderQuantity / limitPrice
	baseQuantity = normalizeQuantity(baseQuantity)

	if orderInfo, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewLimitOrder(symbol, types.Buy, types.ImmediateOrCancel, baseQuantity, limitPrice), &symbolInfo); err != nil {
		err = fmt.Errorf("buyLimit Failed to place order for symbol: %s %v", symbol.String(), err)
		return
	}

	return
}

func (ct *CryptoTrader) liveBuyFixed(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString string
	var orderQuantity, freeQuantity, baseQuantity, price float64
	var orderInfo types.OrderInfo
	var ok bool

	symbolString = symbol.String()
	if _, ok = ct.openTrades[symbolString]; ok {
		logger.Warningf("liveBuyFixed: Trade still open for symbol: %s\n", symbolString)
		return
	}

	orderQuantity = ct.tradeCfg.Volume
	freeQuantity, _ = accountInfo.GetAssetQuantity(symbol.Quote())
	if (freeQuantity * maxPctVolume) < orderQuantity {
		if ct.tradeCfg.Reduce == false {
			logger.Warningf("liveBuyFixed: Insufficient funds to initiate buy position for symbol %s (required: %f, available: %f)", symbolString, orderQuantity, (freeQuantity * 0.995))
			return
		}
		orderQuantity = freeQuantity * maxPctVolume
	}

	if ct.tradeCfg.MaxSlippage > 0.0 {
		orderInfo, price, err = ct.buyLimit(symbol, orderQuantity)
	} else {
		orderInfo, price, err = ct.buyMarket(symbol, orderQuantity)
	}
	if err != nil {
		logger.Errorf("liveBuyFixed: %v\n", err)
		return
	}

	ct.openTrades[symbolString] = orderInfo.UserReference.String()
	logger.Infof("liveBuyFixed: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbolString, baseQuantity, price, ct.openTrades[symbolString])
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Buy", symbolString)))

	return
}

func (ct *CryptoTrader) paperBuyFixed(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString string
	var orderQuantity, baseQuantity, averagePrice float64
	var orderBook types.OrderBook
	var ok bool

	symbolString = symbol.String()
	if _, ok = ct.openTrades[symbolString]; ok {
		logger.Warningf("paperFixed: Trade still open for symbol: %s\n", symbolString)
		return
	}

	orderQuantity = ct.tradeCfg.Volume
	if orderBook, err = ct.exchangeDriver.GetOrderBook(ct.ctx, symbol); err != nil {
		logger.Errorf("paperBuyFixed Failed to retrieve order book for symbol: %s %v\n", symbolString, err)
		return
	}
	baseQuantity, averagePrice = orderBook.GetBuyVolumeAndAveragePrice(orderQuantity)
	baseQuantity = normalizeQuantity(baseQuantity)

	ct.openTrades[symbolString] = uuid.New().String()
	logger.Infof("paperBuyFixed: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbolString, baseQuantity, averagePrice, ct.openTrades[symbolString])
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Buy", symbolString)))

	return
}

func (ct *CryptoTrader) liveBuyPercent(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString string
	var freeQuantity, orderQuantity, baseQuantity, price float64
	var orderInfo types.OrderInfo
	var ok bool

	symbolString = symbol.String()
	if _, ok = ct.openTrades[symbolString]; ok {
		logger.Warningf("liveBuyPercent: Trade still open for symbol: %s\n", symbolString)
		return
	}

	freeQuantity, _ = accountInfo.GetAssetQuantity(symbol.Quote())
	orderQuantity = freeQuantity * ct.tradeCfg.Volume

	if ct.tradeCfg.MaxSlippage > 0 {
		orderInfo, price, err = ct.buyLimit(symbol, orderQuantity)
	} else {
		orderInfo, price, err = ct.buyMarket(symbol, orderQuantity)
	}
	if err != nil {
		logger.Errorf("liveBuyPercent: %v\n", err)
		return
	}

	ct.openTrades[symbolString] = orderInfo.UserReference.String()
	logger.Infof("liveBuyPercent: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbolString, baseQuantity, price, ct.openTrades[symbolString])
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Buy", symbolString)))
	return
}

func (ct *CryptoTrader) paperBuyPercent(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString string
	var freeQuantity, orderQuantity, baseQuantity, averagePrice float64
	var orderBook types.OrderBook
	var ok bool

	symbolString = symbol.String()
	if _, ok = ct.openTrades[symbolString]; ok {
		logger.Warningf("paperFixed: Trade still open for symbol: %s\n", symbolString)
		return
	}

	freeQuantity, _ = accountInfo.GetAssetQuantity(symbol.Quote())
	orderQuantity = freeQuantity * ct.tradeCfg.Volume

	if orderBook, err = ct.exchangeDriver.GetOrderBook(ct.ctx, symbol); err != nil {
		logger.Errorf("BuyPercent Failed to retrieve order book for symbol: %s %v\n", symbol.String(), err)
		return
	}
	baseQuantity, averagePrice = orderBook.GetBuyVolumeAndAveragePrice(orderQuantity)
	baseQuantity = normalizeQuantity(baseQuantity)

	ct.openTrades[symbolString] = uuid.New().String()
	logger.Infof("paperBuyPercent: Buy %s [Amount: %f; Average Price: %f; UUID: %s]", symbolString, baseQuantity, averagePrice, ct.openTrades[symbolString])
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Buy", symbolString)))
	return
}

func (ct *CryptoTrader) liveClosePosition(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString, origTradeID string
	var origTradeUUID uuid.UUID
	var orderInfo types.OrderInfo
	var trades []types.Trade
	var trade types.Trade
	var baseQuantity float64
	var ok bool

	symbolString = symbol.String()
	if origTradeID, ok = ct.openTrades[symbol.String()]; !ok {
		logger.Warningf("liveClosePosition: No trade to close for symbol: %s\n", symbol.String())
		return
	}
	if origTradeUUID, err = uuid.Parse(origTradeID); err != nil {
		logger.Errorf("liveClosePosition: Invalid original trade uuid: %s %v\n", origTradeID, err)
		return
	}

	if orderInfo, err = ct.exchangeDriver.GetOrder(ct.ctx, types.Order{
		Symbol:        symbol,
		UserReference: origTradeUUID,
	}); err != nil {
		logger.Errorf("liveClosePosition Error: %v\n", err)
		return
	}

	if trades, err = ct.exchangeDriver.GetOrderTrades(ct.ctx, orderInfo); err != nil {
		logger.Errorf("liveClosePosition Error: %v\n", err)
		return
	}

	for _, trade = range trades {
		baseQuantity += trade.Quantity
		if trade.CommissionAsset == symbol.Base() {
			baseQuantity -= trade.Commission
		}
	}
	baseQuantity = normalizeQuantity(baseQuantity)

	if _, err = ct.exchangeDriver.PlaceOrder(ct.ctx, types.NewMarketOrder(symbol, types.Sell, baseQuantity), nil); err != nil {
		logger.Warningf("liveClosePosition Unable to close position for symbol: %s %v\n", symbolString, err)
	}

	logger.Infof("liveClosePosition: Symbol %s sell quantity: %f\n", symbolString, baseQuantity)
	delete(ct.openTrades, symbol.String())
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Sell", symbol.String())))

	return
}

func (ct *CryptoTrader) paperClosePosition(accountInfo types.AccountInfo, symbol types.Symbol) (err error) {
	var symbolString string
	var price float64
	var ok bool

	symbolString = symbol.String()
	if _, ok = ct.openTrades[symbol.String()]; !ok {
		logger.Warningf("paperClosePosition: No trade to close for symbol: %s\n", symbolString)
		return
	}

	if price, err = ct.exchangeDriver.Ticker(ct.ctx, symbol); err != nil {
		logger.Errorf("paperClosePosition: Failed to retrieve market price for symbol: %s %v\n", symbolString, err)
		return
	}

	logger.Printf("PaperTrade: Sell %s [Market Price: %f; UUID: %s]", symbolString, price, ct.openTrades[symbolString])
	delete(ct.openTrades, symbolString)
	ct.notifier.Notify(ct.ctx, []byte(fmt.Sprintf("Signal %s Sell", symbolString)))

	return
}
