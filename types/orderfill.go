package types

type OrderFill struct {
	Price           float64
	Quantity        float64
	Commission      float64
	CommissionAsset string
}

func NewOrderFill(p, q, c float64, ca string) (f OrderFill) {
	f.Price = p
	f.Quantity = q
	f.Commission = c
	f.CommissionAsset = ca
	return
}
