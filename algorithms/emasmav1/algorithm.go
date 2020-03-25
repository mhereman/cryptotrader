package emasmav1

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/markcheno/go-talib"

	"github.com/mhereman/cryptotrader/algorithms"
	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

const (
	name string = "Ema/Sma"

	cfgSmaLen    = "Ema/Sma.sma_len"
	cfgEmaLen    = "Ema/Sma.ema_len"
	cfgRsiLen    = "Ema/Sma.rsi_len"
	cfgRsiBuyMax = "Ema/Sma.rsi_buy_max"
)

var defaultConfig types.AlgorithmConfig = types.AlgorithmConfig{
	cfgSmaLen:    "14",
	cfgEmaLen:    "7",
	cfgRsiLen:    "14",
	cfgRsiBuyMax: "90.0",
}

func init() {
	algorithms.RegisterAlgorithm(name, createAlgorithm)
}

// Algorithm represents the Ema/Sma algorithm
type Algorithm struct {
	smaLen        int
	emaLen        int
	rsiLen        int
	rsiBuyMax     float64
	seriesChannel types.SeriesChannel
	signalChannel types.SignalChannel
}

// NewAlgorithm creates a new Ema/Sma algorithm
func NewAlgorithm() (a *Algorithm, err error) {
	a = new(Algorithm)
	if err = a.configure(defaultConfig); err != nil {
		a = nil
		return
	}
	return
}

func createAlgorithm() (algo interfaces.IAlgorithm, err error) {
	algo, err = NewAlgorithm()
	return
}

// Name returns the name of the algorithm
func (a Algorithm) Name() string {
	return name
}

// DefaultConfig returns the default configuration of the algorithm
func (a Algorithm) DefaultConfig() types.AlgorithmConfig {
	return defaultConfig
}

// Config returns the current configuration of the algorithm
func (a Algorithm) Config() types.AlgorithmConfig {
	return types.AlgorithmConfig{
		cfgSmaLen:    fmt.Sprintf("%d", a.smaLen),
		cfgEmaLen:    fmt.Sprintf("%d", a.emaLen),
		cfgRsiLen:    fmt.Sprintf("%d", a.rsiLen),
		cfgRsiBuyMax: fmt.Sprintf("%f", a.rsiBuyMax),
	}
}

// RunAsync runs the algorithm in a goroutine
func (a *Algorithm) RunAsync(ctx context.Context, config types.AlgorithmConfig, seriesChannel types.SeriesChannel, signalChannel types.SignalChannel, waitGroup *sync.WaitGroup) (err error) {
	a.seriesChannel = seriesChannel
	a.signalChannel = signalChannel

	if err = a.configure(config); err != nil {
		logger.Errorf("Algorithm[%s]::RunAsync Error %v", name, err)
		return
	}

	waitGroup.Add(1)
	go runRoutine(ctx, waitGroup, a.seriesChannel, a)
	return
}

func (a *Algorithm) emit(signal types.Signal) {
	a.signalChannel <- signal
}

func (a *Algorithm) check(ctx context.Context, series types.Series) {
	var sma, ema, rsi []float64
	var open, close float64
	var buySignal, sellSignal bool

	sma = talib.Sma(series.Close(), a.smaLen)
	ema = talib.Ema(series.Close(), a.emaLen)
	rsi = talib.Rsi(series.Close(), a.rsiLen)
	open = series.PreviousOpen()
	close = series.PreviousClose()

	buySignal = talib.Crossover(ema, sma) && rsi[len(rsi)-1] < a.rsiBuyMax && close > open
	sellSignal = talib.Crossunder(ema, sma)

	if buySignal {
		logger.Debugf("EMIT BUY")
		a.emit(types.NewSignal(name, series.Symbol, types.Buy))
	}

	if sellSignal {
		logger.Debugf("EMIT SELL")
		a.emit(types.NewSignal(name, series.Symbol, types.Sell))
	}
}

func (a *Algorithm) configure(config types.AlgorithmConfig) (err error) {
	var key, value string
	for key, value = range config {
		switch key {
		case cfgSmaLen:
			if a.smaLen, err = strconv.Atoi(value); err != nil {
				return
			}
		case cfgEmaLen:
			if a.emaLen, err = strconv.Atoi(value); err != nil {
				return
			}
		case cfgRsiLen:
			if a.rsiLen, err = strconv.Atoi(value); err != nil {
				return
			}
		case cfgRsiBuyMax:
			if a.rsiBuyMax, err = strconv.ParseFloat(value, 64); err != nil {
				return
			}
		}
	}
	return
}

func runRoutine(ctx context.Context, wg *sync.WaitGroup, seriesChannel types.SeriesChannel, a *Algorithm) {
	defer wg.Done()

	var runLoop bool
	var series types.Series

	runLoop = true
	for runLoop {
		select {
		case <-ctx.Done():
			runLoop = false
		case series = <-seriesChannel:
			a.check(ctx, series)
		}
	}
}
