package types

// SeriesChannel channel to report series on
type SeriesChannel chan Series

// Series represents a series of candles of a symbol
type Series struct {
	// Symbol of the series
	Symbol Symbol

	// Timeframe of the series
	Timeframe Timeframe

	// Candles of the series (last candle in the array is the current candle)
	Candles []OHLC
}

// NewSeries creates a Series instance
func NewSeries(symbol Symbol, timeframe Timeframe, candles []OHLC) Series {
	return Series{
		Symbol:    symbol,
		Timeframe: timeframe,
		Candles:   candles,
	}
}

// CurrentOpen price (the active candle)
func (s Series) CurrentOpen() float64 {
	return s.Candles[len(s.Candles)-1].Open
}

// PreviousOpen price (the last finished candle)
func (s Series) PreviousOpen() float64 {
	return s.Candles[len(s.Candles)-2].Open
}

// Open prices
func (s Series) Open() (res []float64) {
	var numCandles, index int
	var candle OHLC

	numCandles = len(s.Candles)
	res = make([]float64, numCandles, numCandles)
	for index, candle = range s.Candles {
		res[index] = candle.Open
	}
	return
}

// OpenRange prices
func (s Series) OpenRange(start int, size int) (res []float64) {
	var index int

	if (start + size) > len(s.Candles) {
		size = len(s.Candles) - start
	}

	res = make([]float64, size, size)
	for index = start; index < (start + size); index++ {
		res[index-start] = s.Candles[index].Open
	}
	return
}

// OpenLastN prices
func (s Series) OpenLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.OpenRange(start, size)
	return
}

// CurrentHigh price (the active candle)
func (s Series) CurrentHigh() float64 {
	return s.Candles[len(s.Candles)-1].High
}

// PreviousHigh price (the last finished candle)
func (s Series) PreviousHigh() float64 {
	return s.Candles[len(s.Candles)-2].High
}

// High prices
func (s Series) High() (res []float64) {
	var numCandles, index int
	var candle OHLC

	numCandles = len(s.Candles)
	res = make([]float64, numCandles, numCandles)
	for index, candle = range s.Candles {
		res[index] = candle.High
	}
	return
}

// HighRange prices
func (s Series) HighRange(start int, size int) (res []float64) {
	var index int

	if (start + size) > len(s.Candles) {
		size = len(s.Candles) - start
	}

	res = make([]float64, size, size)
	for index = start; index < (start + size); index++ {
		res[index-start] = s.Candles[index].High
	}
	return
}

// HighLastN prices
func (s Series) HighLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.HighRange(start, size)
	return
}

// CurrentLow price (the active candle)
func (s Series) CurrentLow() float64 {
	return s.Candles[len(s.Candles)-1].Low
}

// PreviousLow price (the last finished candle)
func (s Series) PreviousLow() float64 {
	return s.Candles[len(s.Candles)-2].Low
}

// Low prices
func (s Series) Low() (res []float64) {
	var numCandles, index int
	var candle OHLC

	numCandles = len(s.Candles)
	res = make([]float64, numCandles, numCandles)
	for index, candle = range s.Candles {
		res[index] = candle.Low
	}
	return
}

// LowRange prices
func (s Series) LowRange(start int, size int) (res []float64) {
	var index int

	if (start + size) > len(s.Candles) {
		size = len(s.Candles) - start
	}

	res = make([]float64, size, size)
	for index = start; index < (start + size); index++ {
		res[index-start] = s.Candles[index].Low
	}
	return
}

// LowLastN prices
func (s Series) LowLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.LowRange(start, size)
	return
}

// CurrentClose price (the active candle)
func (s Series) CurrentClose() float64 {
	return s.Candles[len(s.Candles)-1].Close
}

// PreviousClose price (the last finished candle)
func (s Series) PreviousClose() float64 {
	return s.Candles[len(s.Candles)-2].Close
}

// Close prices
func (s Series) Close() (res []float64) {
	var numCandles, index int
	var candle OHLC

	numCandles = len(s.Candles)
	res = make([]float64, numCandles, numCandles)
	for index, candle = range s.Candles {
		res[index] = candle.Close
	}
	return
}

// CloseRange prices
func (s Series) CloseRange(start int, size int) (res []float64) {
	var index int

	if (start + size) > len(s.Candles) {
		size = len(s.Candles) - start
	}

	res = make([]float64, size, size)
	for index = start; index < (start + size); index++ {
		res[index-start] = s.Candles[index].Close
	}
	return
}

// CloseLastN prices
func (s Series) CloseLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.CloseRange(start, size)
	return
}
