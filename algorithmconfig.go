package cryptotrader

import "github.com/mhereman/cryptotrader/types"

// AlgorithmConfig represents the config for the algorithm to use to make buy and sell descissions
type AlgorithmConfig struct {
	// Name of the algorithm
	Name string

	// Config of the algorithm
	Config types.AlgorithmConfig
}

// NewAlgorithmConfigFromFlags creates a new AlgorithmConfig insance from the cmdline argument values
func NewAlgorithmConfigFromFlags(name string, args map[string]string) (ac AlgorithmConfig, err error) {
	ac.Name = name
	ac.Config = types.AlgorithmConfig(args)
	return
}
