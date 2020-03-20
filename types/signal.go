package types

import (
	"fmt"
	"time"
)

type Signal struct {
	AlgorithmName string
	Symbol        Symbol
	Side          Side
	SignalTime    time.Time
}

func (s Signal) String() string {
	return fmt.Sprintf("%s - %s: %s %v", s.AlgorithmName, s.Symbol.String(), s.Side.String(), s.SignalTime)
}

type SignalChannel chan Signal
