package binance

import (
	"context"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
)

func (b Binance) TestConnectivity(ctx context.Context) (ok bool, err error) {
	var ps *bin.PingService

	ps = b.client.NewPingService()
	if err = ps.Do(ctx); err != nil {
		logger.Errorf("Binance::TestConnectivity Error: %v\n", err)
		return
	}
	ok = true
	return
}
