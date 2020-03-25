package types

import (
	"fmt"
	"strconv"
	"time"
)

// TimeframeUnit represents the unit for a timeframe
type TimeframeUnit string

const (
	// TuSec second unit
	TuSec = TimeframeUnit("s")

	// TuMin minute unit
	TuMin = TimeframeUnit("m")

	// TuHour hour unit
	TuHour = TimeframeUnit("h")

	// TuDay day unit
	TuDay = TimeframeUnit("d")

	// TuWeek week unit
	TuWeek = TimeframeUnit("w")

	// TuMonth month unit
	TuMonth = TimeframeUnit("M")
)

// Timeframe represents a timeframe on which to retrieve candles and execute the algorithm on
type Timeframe struct {
	// Value of the timeframe
	Value int

	// Unit of the timeframe's value
	Unit TimeframeUnit
}

// NewTimeframe creates a new Timeframe instance
func NewTimeframe(value int, unit TimeframeUnit) Timeframe {
	return Timeframe{
		Value: value,
		Unit:  unit,
	}
}

// NewTimeframeFromString creates a new Timeframe instance from its string representation
func NewTimeframeFromString(in string) (tf Timeframe, err error) {
	var unitPart, valuePart string
	unitPart = in[len(in)-1:]
	valuePart = in[:len(in)-1]

	switch TimeframeUnit(unitPart) {
	case TuSec, TuMin, TuHour, TuDay, TuWeek, TuMonth:
		tf.Unit = TimeframeUnit(unitPart)
	default:
		err = fmt.Errorf("Invalid unit: %s", unitPart)
		return
	}

	if tf.Value, err = strconv.Atoi(valuePart); err != nil {
		return
	}
	return
}

// String returns the string version of the timeframe
func (tf Timeframe) String() string {
	return fmt.Sprintf("%d%s", tf.Value, tf.Unit)
}

// NextOpen calculates the next open time base on the current open time
func (tf Timeframe) NextOpen(currentOpen time.Time, currentClose time.Time) (nextOpen time.Time) {
	switch tf.Unit {
	case TuSec:
		nextOpen = currentOpen.Add(time.Second * time.Duration(tf.Value))
	case TuMin:
		nextOpen = currentOpen.Add(time.Minute + time.Duration(tf.Value))
	case TuHour:
		nextOpen = currentOpen.Add(time.Hour * time.Duration(tf.Value))
	case TuDay:
		nextOpen = currentOpen.AddDate(0, 0, 1)
	case TuWeek:
		nextOpen = currentOpen.AddDate(0, 0, 7)
	case TuMonth:
		nextOpen = currentOpen.AddDate(0, 1, 0)
	default:
		nextOpen = currentClose
	}
	return
}
