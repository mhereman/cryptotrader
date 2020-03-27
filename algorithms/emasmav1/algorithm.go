package emasmav1

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	cfgBacktest  = "Ema/Sma.backtest"
)

var defaultConfig types.AlgorithmConfig = types.AlgorithmConfig{
	cfgSmaLen:    "14",
	cfgEmaLen:    "7",
	cfgRsiLen:    "14",
	cfgRsiBuyMax: "90.0",
	cfgBacktest:  "false",
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
	backtest      bool
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
		cfgBacktest:  fmt.Sprintf("%t", a.backtest),
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
	var calcSeries types.Series

	calcSeries = series.SubSeries(0, series.Length()-1)

	sma = talib.Sma(calcSeries.Close(), a.smaLen)
	ema = talib.Ema(calcSeries.Close(), a.emaLen)
	rsi = talib.Rsi(calcSeries.Close(), a.rsiLen)
	open = calcSeries.CurrentOpen()
	close = calcSeries.CurrentClose()

	buySignal = talib.Crossover(ema, sma) && rsi[len(rsi)-1] < a.rsiBuyMax && close > open
	sellSignal = talib.Crossunder(ema, sma)

	if buySignal {
		logger.Debugf("EMIT BUY")
		if a.backtest {
			a.emit(types.NewBacktestSignal(name, series.Symbol, types.Buy, calcSeries.CurrentCandleTime()))
		} else {
			a.emit(types.NewSignal(name, series.Symbol, types.Buy))
		}
	}

	if sellSignal {
		logger.Debugf("EMIT SELL")
		if a.backtest {
			a.emit(types.NewBacktestSignal(name, series.Symbol, types.Sell, calcSeries.CurrentCandleTime()))
		} else {
			a.emit(types.NewSignal(name, series.Symbol, types.Sell))
		}
	}
}

func (a *Algorithm) checkBacktest(ctx context.Context, series types.Series) {
	var minSampleLen, length int
	var subSeries types.Series

	minSampleLen = int(math.Max(math.Max(float64(a.smaLen), float64(a.emaLen)), float64(a.rsiLen))) + 2

	for length = minSampleLen; length <= series.Length(); length++ {
		subSeries = series.SubSeries(0, length)
		a.check(ctx, subSeries)
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
		case cfgBacktest:
			a.backtest = strings.ToLower(value) == "true"
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
			logger.Debugf("Algorithm :%s received new data\n", a.Name())
			if a.backtest {
				a.checkBacktest(ctx, series)
				continue
			}
			a.check(ctx, series)
		}
	}
}
