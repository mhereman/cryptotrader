package types

// OrderFill represents a filled trade of an order
// And order can have multiple fill's before being completely
// executed
type OrderFill struct {
	// Price of the fill
	Price float64

	// Quantity filled (in base asset)
	Quantity float64

	// Commission payed
	Commission float64

	// Asset in which the commission is payed
	CommissionAsset string
}

// NewOrderFill creates a new OrderFill instance
func NewOrderFill(price float64, quantity float64, commission float64, commissionAsset string) (f OrderFill) {
	f.Price = price
	f.Quantity = quantity
	f.Commission = commission
	f.CommissionAsset = commissionAsset
	return
}
