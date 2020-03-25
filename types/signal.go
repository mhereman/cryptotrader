package types

import (
	"fmt"
	"time"
)

// SignalChannel channel to report signals on
type SignalChannel chan Signal

// Signal represents a buy or sell signal
type Signal struct {
	// AlgorithmName of the algorithm issuing the signal
	AlgorithmName string

	// Symbol on which the signal applies
	Symbol Symbol

	// Side of the signal
	Side Side

	// Time of the signal
	SignalTime time.Time
}

// NewSignal creates a new Signal instance
func NewSignal(algoName string, symbol Symbol, side Side) Signal {
	return Signal{
		AlgorithmName: algoName,
		Symbol:        symbol,
		Side:          side,
		SignalTime:    time.Now(),
	}
}

// String returns a string representation of the signal
func (s Signal) String() string {
	return fmt.Sprintf("%s - %s: %s %v", s.AlgorithmName, s.Symbol.String(), s.Side.String(), s.SignalTime)
}
