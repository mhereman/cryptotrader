package binance

import (
	"context"
	"strconv"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// GetAccountInfo executes the get account info request
func (b *Binance) GetAccountInfo(ctx context.Context) (info types.AccountInfo, err error) {
	var as *bin.GetAccountService
	var response *bin.Account
	var numBalances, index int
	var balance bin.Balance
	var balances []types.AccountBalance
	var free, locked float64

	as = b.client.NewGetAccountService()
	if response, err = as.Do(ctx); err != nil {
		logger.Errorf("Binance::GetAccount Error: %v\n", err)
		return
	}

	numBalances = len(response.Balances)
	balances = make([]types.AccountBalance, numBalances, numBalances)
	for index, balance = range response.Balances {
		if free, err = strconv.ParseFloat(balance.Free, 64); err != nil {
			logger.Errorf("Binance::GetAccount Error: %v\n", err)
			return
		}
		if locked, err = strconv.ParseFloat(balance.Locked, 64); err != nil {
			logger.Errorf("Binance::GetAccount Error: %v\n", err)
			return
		}

		balances[index] = types.NewAccountBalance(balance.Asset, free, locked)
	}

	info.MakerCommission = float64(response.MakerCommission) * 0.0001
	info.TakerCommission = float64(response.TakerCommission) * 0.0001
	info.BuyerCommission = float64(response.BuyerCommission) * 0.0001
	info.SellerCommission = float64(response.SellerCommission) * 0.0001
	info.Balances = balances

	return
}
