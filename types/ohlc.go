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

// NewOHLC ...
func NewOHLC(o, h, l, c, v float64, ot, ct time.Time) OHLC {
	return OHLC{
		Open:      o,
		High:      h,
		Low:       l,
		Close:     c,
		Volume:    v,
		OpenTime:  ot,
		CloseTime: ct,
	}
}

func (o OHLC) String() string {
	return fmt.Sprintf("[%v] O: %f  H: %f  L: %f  C: %f", o.OpenTime, o.Open, o.High, o.Low, o.Close)
}
