package main

// Register plugins here
import (
	// Exchanges
	_ "github.com/mhereman/cryptotrader/exchange/binance"

	// Algorithms
	_ "github.com/mhereman/cryptotrader/algorithms/emasmav1"

	// Notifiers
	_ "github.com/mhereman/cryptotrader/notifiers/noop"
	_ "github.com/mhereman/cryptotrader/notifiers/proximussms"
)
