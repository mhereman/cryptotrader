package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) GetOrderBook(ctx context.Context, symbol types.Symbol) (book types.OrderBook, err error) {
	var ds *bin.DepthService
	var binanceSymbol string
	var response *bin.DepthResponse
	var numBids, numAsks, index int
	var bid bin.Bid
	var ask bin.Ask
	var bids, asks []types.OrderBookEntry

	if binanceSymbol, err = b.symbolToBinance(symbol); err != nil {
		logger.Errorf("Binance::GetOrderBook Error: %v\n", err)
		return
	}

	ds = b.client.NewDepthService()
	ds.Symbol(binanceSymbol)
	if response, err = ds.Do(ctx); err != nil {
		logger.Errorf("Binance::GetOrderBook Error: %v\n", err)
		return
	}

	numBids = len(response.Bids)
	numAsks = len(response.Asks)
	bids = make([]types.OrderBookEntry, numBids, numBids)
	asks = make([]types.OrderBookEntry, numAsks, numAsks)

	for index, bid = range response.Bids {
		bids[index] = types.NewOrderBookEntry(b.toFloat(bid.Price), b.toFloat(bid.Quantity))
	}
	for index, ask = range response.Asks {
		asks[index] = types.NewOrderBookEntry(b.toFloat(ask.Price), b.toFloat(ask.Quantity))
	}
	book = types.NewOrderBook(bids, asks)
	return
}
