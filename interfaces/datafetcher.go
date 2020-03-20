package interfaces

import (
	"context"
	"sync"

	"github.com/mhereman/cryptotrader/types"
)

type IDataFether interface {
	Register(context.Context, types.Symbol, types.Timeframe) (types.SeriesChannel, error)
	RunAsync(context.Context, *sync.WaitGroup)
}
