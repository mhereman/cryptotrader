package types

import "strings"

type AccountInfo struct {
	MakerCommission  float64
	TakerCommission  float64
	BuyerCommission  float64
	SellerCommission float64
	Balances         []AccountBalance
}

func NewAccountInfo(mc, tc, bc, sc float64, balances []AccountBalance) AccountInfo {
	return AccountInfo{
		MakerCommission:  mc,
		TakerCommission:  tc,
		BuyerCommission:  bc,
		SellerCommission: sc,
		Balances:         balances,
	}
}

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
