package types

type SeriesChannel chan Series

type Series struct {
	Symbol    Symbol
	Timeframe Timeframe
	Candles   []OHLC
}

// First Candle in the array is the most recent one
func NewSeries(s Symbol, tf Timeframe, c []OHLC) Series {
	return Series{
		Symbol:    s,
		Timeframe: tf,
		Candles:   c,
	}
}

func (s Series) LastOpen() float64 {
	return s.Candles[len(s.Candles)-1].Open
}

func (s Series) PrevOpen() float64 {
	return s.Candles[len(s.Candles)-2].Open
}

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

func (s Series) OpenLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.OpenRange(start, size)
	return
}

func (s Series) LastHigh() float64 {
	return s.Candles[len(s.Candles)-1].High
}

func (s Series) PrevHigh() float64 {
	return s.Candles[len(s.Candles)-2].High
}

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

func (s Series) HighLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.HighRange(start, size)
	return
}

func (s Series) LastLow() float64 {
	return s.Candles[len(s.Candles)-1].Low
}

func (s Series) PrevLow() float64 {
	return s.Candles[len(s.Candles)-2].Low
}

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

func (s Series) LowLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.LowRange(start, size)
	return
}

func (s Series) LastClose() float64 {
	return s.Candles[len(s.Candles)-1].Close
}

func (s Series) PrevClose() float64 {
	return s.Candles[len(s.Candles)-2].Close
}

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

func (s Series) CloseLastN(size int) (res []float64) {
	var start int

	start = len(s.Candles) - size
	res = s.CloseRange(start, size)
	return
}
