package cryptotrader

import (
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

// AssetConfig represents the config for the asset to trade
type AssetConfig struct {
	// Symbol to trade
	Symbol types.Symbol

	// Timeframe to trade
	Timeframe types.Timeframe
}

// NewAssetConfigFromFlags creates a new AssetConfig insance from the cmdline argument values
func NewAssetConfigFromFlags(base string, quote string, timeframe string) (ac AssetConfig, err error) {
	ac.Symbol = types.NewSymbol(base, quote)
	if ac.Timeframe, err = types.NewTimeframeFromString(timeframe); err != nil {
		logger.Errorf("Error parsing timeframe: %s %v\n", timeframe, err)
		return
	}
	return
}
