package interfaces

import (
	"context"
	"sync"

	"github.com/mhereman/cryptotrader/types"
)

// IDataFetcher represents the datafetcher
type IDataFetcher interface {
	// Register registers a symbol and timeframe to fetch
	Register(context.Context, types.Symbol, types.Timeframe) (types.SeriesChannel, error)

	// RunAsync runs the DataFetcher in a goroutine
	RunAsync(context.Context, *sync.WaitGroup)
}
