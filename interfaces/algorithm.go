package interfaces

import (
	"context"
	"sync"

	"github.com/mhereman/cryptotrader/types"
)

// IAlgorithm represents the Algorithm plugin interface
type IAlgorithm interface {
	// Name returns the name of the algorithm
	Name() string

	// DefaultConfig returns the default configuration of the algorithm
	DefaultConfig() types.AlgorithmConfig

	// Config returns the current configuration of the algorithm
	Config() types.AlgorithmConfig

	// RunAsync runs the algorithm in a goroutine
	RunAsync(context.Context, types.AlgorithmConfig, types.SeriesChannel, types.SignalChannel, *sync.WaitGroup) error
}
