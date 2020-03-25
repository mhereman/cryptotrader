package binance

import (
	"context"
	"time"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/logger"
)

// GetServerTime executes the get server time request
func (b Binance) GetServerTime(ctx context.Context) (serverTime time.Time, err error) {
	var sts *bin.ServerTimeService
	var response int64

	sts = b.client.NewServerTimeService()
	if response, err = sts.Do(ctx); err != nil {
		logger.Errorf("Binance::GetServerTime Error: %v\n", err)
		return
	}

	serverTime = b.toTime(response)

	return
}
