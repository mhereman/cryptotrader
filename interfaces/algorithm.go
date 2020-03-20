package interfaces

import (
	"context"
	"sync"

	"github.com/mhereman/cryptotrader/types"
)

type IAlgorithm interface {
	Name() string
	DefaultConfig() types.AlgorithmConfig
	Config() types.AlgorithmConfig
	RunAsync(context.Context, types.AlgorithmConfig, types.SeriesChannel, types.SignalChannel, *sync.WaitGroup) error
}
