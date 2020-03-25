package types

// OrderBookEntry represents an entry on the order book
type OrderBookEntry struct {
	// Price on the book
	Price float64

	// Quantity (in base asset) on the book
	Quantity float64
}

// NewOrderBookEntry creates a new OrderBookEntry instance
func NewOrderBookEntry(price float64, quantity float64) OrderBookEntry {
	return OrderBookEntry{
		Price:    price,
		Quantity: quantity,
	}
}
