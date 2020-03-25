package types

// OrderBook represents an assets orderbook on the exchange
type OrderBook struct {
	// Symbol of the order book
	Symbol Symbol

	// Bids on the order book
	Bids []OrderBookEntry

	// Asks on the order book
	Asks []OrderBookEntry
}

// NewOrderBook creates a new OrderBook instance
func NewOrderBook(symbol Symbol, bids []OrderBookEntry, asks []OrderBookEntry) OrderBook {
	return OrderBook{
		Symbol: symbol,
		Bids:   bids,
		Asks:   asks,
	}
}

// GetBuyVolumeAndAveragePrice calculates the buy volume and average price in base asset from the
// provided volume in quote asset
func (ob OrderBook) GetBuyVolumeAndAveragePrice(quoteAmount float64) (baseAmount float64, averagePrice float64) {
	var ask OrderBookEntry
	var price, quantity, multipliedPrice float64

	baseAmount = 0
	multipliedPrice = 0
	for _, ask = range ob.Asks {
		price = ask.Price
		quantity = ask.Quantity

		if (quantity * price) >= quoteAmount {
			//baseAmount += quoteAmount / price
			baseAmount += quantity * (quoteAmount / (quantity * price))
			multipliedPrice += (quantity * (quoteAmount / (quantity * price))) * price
			break
		}
		baseAmount += quantity
		quoteAmount -= (quantity * price)
		multipliedPrice += price * quantity
	}
	averagePrice = multipliedPrice / baseAmount

	return
}
