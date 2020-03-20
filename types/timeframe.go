package types

import (
	"fmt"
	"strconv"
	"time"
)

type TimeframeUnit string

const (
	TuSec   = TimeframeUnit("s")
	TuMin   = TimeframeUnit("m")
	TuHour  = TimeframeUnit("h")
	TuDay   = TimeframeUnit("d")
	TuWeek  = TimeframeUnit("w")
	TuMonth = TimeframeUnit("M")
)

type Timeframe struct {
	Value int
	Unit  TimeframeUnit
}

func NewTimeframe(v int, u TimeframeUnit) Timeframe {
	return Timeframe{
		Value: v,
		Unit:  u,
	}
}

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

func (tf Timeframe) String() string {
	return fmt.Sprintf("%d%s", tf.Value, tf.Unit)
}

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
