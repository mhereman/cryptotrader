package types

import "strings"

// AccountInfo represents information of the account on the exchange
type AccountInfo struct {
	// MakeCommission commission for maker orders
	MakerCommission float64

	// TakerCommission commission for taker orders
	TakerCommission float64

	// BuyerCommission commission for buy orders
	BuyerCommission float64

	// SellerCommission commission for sell orders
	SellerCommission float64

	// Balances the balances of the account
	Balances []AccountBalance
}

// NewAccountInfo creates a new AccountInfo instance
func NewAccountInfo(makerCommission float64, takerCommission float64, buyerCommission float64, sellerCommission float64, balances []AccountBalance) AccountInfo {
	return AccountInfo{
		MakerCommission:  makerCommission,
		TakerCommission:  takerCommission,
		BuyerCommission:  buyerCommission,
		SellerCommission: sellerCommission,
		Balances:         balances,
	}
}

// GetAssetQuantity returns the balance information for the requested asset
func (ai AccountInfo) GetAssetQuantity(asset string) (free float64, locked float64) {
	var balance AccountBalance

	asset = strings.ToUpper(asset)
	for _, balance = range ai.Balances {
		if balance.Asset == asset {
			free = balance.Free
			locked = balance.Locked
			return
		}
	}
	return
}
