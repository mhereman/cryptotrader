package cryptotrader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
	"github.com/mhereman/cryptotrader/types"
)

const (
	checkInterval time.Duration = (time.Millisecond * 1000)
	sleepInterval time.Duration = (time.Millisecond * 500)
)

type DataFetcher struct {
	driver        interfaces.IExchangeDriver
	channels      map[string]map[string]types.SeriesChannel
	refreshTime   map[string]map[string]time.Time
	lastCheckTime time.Time
	mux           sync.RWMutex
}

func NewDataFetcher(driver interfaces.IExchangeDriver) (dc *DataFetcher) {
	dc = &DataFetcher{
		driver:        driver,
		channels:      make(map[string]map[string]types.SeriesChannel),
		refreshTime:   make(map[string]map[string]time.Time),
		lastCheckTime: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		mux:           sync.RWMutex{},
	}
	return
}

func (dc *DataFetcher) RunAsync(ctx context.Context, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go fetchRoutine(ctx, waitGroup, dc)
}

func (dc *DataFetcher) Register(ctx context.Context, symbol types.Symbol, timeFrame types.Timeframe) (channel types.SeriesChannel, err error) {
	dc.mux.Lock()
	defer dc.mux.Unlock()

	var ok bool
	var subMap1 map[string]types.SeriesChannel
	var subMap2 map[string]time.Time

	var nextRefreshTime time.Time

	if subMap1, ok = dc.channels[symbol.String()]; !ok {
		subMap1 = make(map[string]types.SeriesChannel)
	}
	if subMap2, ok = dc.refreshTime[symbol.String()]; !ok {
		subMap2 = make(map[string]time.Time)
	}

	if _, ok = subMap1[timeFrame.String()]; ok {
		err = fmt.Errorf("DataCacher::Register Error Symbol %s%s with timeframe %d%s already registered", symbol.Base(), symbol.Quote(), timeFrame.Value, timeFrame.Unit)
		return
	}

	nextRefreshTime = nextRefreshTime.Add(time.Second * 2)
	channel = make(types.SeriesChannel)

	subMap1[timeFrame.String()] = channel
	subMap2[timeFrame.String()] = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
	dc.channels[symbol.String()] = subMap1
	dc.refreshTime[symbol.String()] = subMap2

	return
}

func (dc *DataFetcher) fetchExpiredData(ctx context.Context) {
	var subMap1 map[string]time.Time
	var subMap2 map[string]types.SeriesChannel
	var symbolString, timeframeString string
	var symbol types.Symbol
	var timeFrame types.Timeframe
	var refreshTime, nextRefreshTime time.Time
	var series types.Series
	var seriesChannel types.SeriesChannel
	var err error
	var ok bool
	var refTime time.Time

	refTime = time.Now()
	if refTime.Sub(dc.lastCheckTime) < checkInterval {
		return
	}

	dc.mux.Lock()
	defer dc.mux.Unlock()

	for symbolString, subMap1 = range dc.refreshTime {
		if symbol, err = types.NewSymbolFromString(symbolString); err != nil {
			logger.Errorf("DataCacher::FetchExpiredData Error %v\n", err)
			return
		}

		for timeframeString, refreshTime = range subMap1 {
			if timeFrame, err = types.NewTimeframeFromString(timeframeString); err != nil {
				logger.Errorf("DataCacher::FetchExpiredData error %v\n", err)
				return
			}

			if refreshTime.Before(refTime) {
				logger.Debugf("Fetching new data: %v\n", refTime)
				if series, nextRefreshTime, err = dc.fetchData(ctx, symbol, timeFrame); err != nil {
					logger.Errorf("DataCacher::FetchExpiredData Error %v\n", err)
					continue
				}

				if subMap2, ok = dc.channels[symbolString]; !ok {
					err = fmt.Errorf("No channel for series %s[%s]", symbolString, timeframeString)
					logger.Errorf("DataCacher::FetchExpiredData Error %v\n", err)
					continue
				}
				if seriesChannel, ok = subMap2[timeframeString]; !ok {
					err = fmt.Errorf("No channel for series %s[%s]", symbolString, timeframeString)
					logger.Errorf("DataCacher::FetchExpiredData Error %v\n", err)
					continue
				}

				// We want to be sure the next candle has been produced completely so we ad 100milli slip time
				nextRefreshTime = nextRefreshTime.Add(time.Second * 2)

				logger.Debugf("Pushing new data")
				seriesChannel <- series
				subMap1[timeframeString] = nextRefreshTime
				dc.refreshTime[symbolString] = subMap1
			}
		}
	}
	dc.lastCheckTime = refTime
}

func (dc *DataFetcher) fetchData(ctx context.Context, symbol types.Symbol, timeFrame types.Timeframe) (series types.Series, nextRefreshTime time.Time, err error) {
	var lastCandle types.OHLC

	if series, err = dc.driver.GetSeries(ctx, symbol, timeFrame); err != nil {
		logger.Errorf("DataCacher::fetchData Error %v\n", err)
		return
	}

	lastCandle = series.Candles[len(series.Candles)-1]
	nextRefreshTime = timeFrame.NextOpen(lastCandle.OpenTime, lastCandle.CloseTime)
	return
}

func fetchRoutine(ctx context.Context, wg *sync.WaitGroup, dc *DataFetcher) {
	defer wg.Done()

	var runLoop bool
	runLoop = true

	for runLoop {
		select {
		case <-ctx.Done():
			runLoop = false
		default:
			dc.fetchExpiredData(ctx)
			time.Sleep(sleepInterval)
		}
	}
}
