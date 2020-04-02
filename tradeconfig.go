package cryptotrader

import (
	"fmt"
	"strings"
)

// TradeVolumeType represents the way to calculate the trade volume
type TradeVolumeType int

const (
	// TVTFixed use a fixed volume (quote asset) for trading
	TVTFixed TradeVolumeType = iota

	// TVTPercent use a percentage of the available quote asset for trading
	TVTPercent
)

const maxPctVolume = 0.999

// TradeVolumeTypeFromString creates a new TradeVolumeType from it's string representation
func TradeVolumeTypeFromString(in string) (out TradeVolumeType, err error) {
	switch strings.ToLower(in) {
	case "fixed":
		out = TVTFixed
	case "pct", "percent":
		out = TVTPercent
	default:
		err = fmt.Errorf("Invalid tradevolume type: %s", in)
	}
	return
}

// String return the string representation fo the TradeVolumeType
func (tvt TradeVolumeType) String() string {
	if tvt == TVTFixed {
		return "fixed"
	}
	return "percent"
}

// TradeConfig represents the config for the trades to place
type TradeConfig struct {
	// TradeVolumeType how to calculate the volume of the buy/sell orders
	TradeVolumeType TradeVolumeType

	// The volume to trade (depends on TradeVolumeType)
	// If TradeVolumeType == TVTFixed the Volume value is the actual quote asset quantity to trade
	// If TradeVolumeType == TVTPercent the Volume value represents the percentage of the available quote asset quantity to trade (max Volume value = 1.0)
	Volume float64

	// Reduce if true reduces the Volume to the available quantity if the TradeVolumeType == TVTFixed and the available asset quantity is insufficient
	Reduce bool

	// Paper perform paper trading only, do not issue any orer on the exchange
	Paper bool

	// Max slippage in percent
	// 0.1% = 0.001
	MaxSlippage float64
}

// NewTradeConfigFromFlags creates a new TradeConfig insance from the cmdline argument values
func NewTradeConfigFromFlags(tvt string, volume float64, reduce bool, paper bool, maxSlippage float64) (tc TradeConfig, err error) {
	if tc.TradeVolumeType, err = TradeVolumeTypeFromString(tvt); err != nil {
		return
	}
	if tc.TradeVolumeType == TVTPercent {
		if volume > maxPctVolume {
			volume = maxPctVolume
		}
	}
	tc.Volume = volume
	tc.Reduce = reduce
	tc.Paper = paper
	tc.MaxSlippage = maxSlippage
	return
}
