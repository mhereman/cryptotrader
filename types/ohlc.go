package types

import (
	"fmt"
	"time"
)

// OHLC represents a candlestick
type OHLC struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	OpenTime  time.Time
	CloseTime time.Time
}

// NewOHLC creates an OHLC instance
func NewOHLC(openPrice float64, highPrice float64, lowPrice float64, closePrice float64, volume float64, openTime time.Time, closeTime time.Time) OHLC {
	return OHLC{
		Open:      openPrice,
		High:      highPrice,
		Low:       lowPrice,
		Close:     closePrice,
		Volume:    volume,
		OpenTime:  openTime,
		CloseTime: closeTime,
	}
}

// String returns a string version of the candle
func (o OHLC) String() string {
	return fmt.Sprintf("[%v] O: %f  H: %f  L: %f  C: %f", o.OpenTime, o.Open, o.High, o.Low, o.Close)
}
