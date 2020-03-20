package types

type OrderBookEntry struct {
	Price    float64
	Quantity float64
}

func NewOrderBookEntry(p, q float64) OrderBookEntry {
	return OrderBookEntry{
		Price:    p,
		Quantity: q,
	}
}

type OrderBook struct {
	Bids []OrderBookEntry
	Asks []OrderBookEntry
}

func NewOrderBook(b, a []OrderBookEntry) OrderBook {
	return OrderBook{
		Bids: b,
		Asks: a,
	}
}

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
