package binance

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mhereman/cryptotrader/logger"

	bin "github.com/adshao/go-binance"
	"github.com/mhereman/cryptotrader/types"
)

func (b Binance) symbolToBinance(symbol types.Symbol) (binanceSymbol string, err error) {
	var ok bool

	binanceSymbol = fmt.Sprintf("%s%s", symbol.Base(), symbol.Quote())
	if _, ok = b.allSymbols[binanceSymbol]; !ok {
		err = fmt.Errorf("Symbol '%s' is not available on Binance", binanceSymbol)
		binanceSymbol = ""
		return
	}
	return
}

func (b Binance) toSymbol(in string) (symbol types.Symbol, err error) {
	var baseAndQuote []string
	var ok bool

	if baseAndQuote, ok = b.allSymbols[strings.ToUpper(in)]; !ok {
		err = fmt.Errorf("Symbol '%s' is not available on Binance", in)
		return
	}
	symbol = types.NewSymbol(baseAndQuote[0], baseAndQuote[1])
	return
}

func (b Binance) timeframeToBinance(timeframe types.Timeframe) (binanceTimeframe string, err error) {
	switch timeframe.Unit {
	case types.TuSec:
		err = fmt.Errorf("Timeframe %ds is not valid on Binance", timeframe.Value)
	case types.TuMin:
		switch timeframe.Value {
		case 1, 3, 5, 15, 30:
			binanceTimeframe = fmt.Sprintf("%dm", timeframe.Value)
		default:
			err = fmt.Errorf("Timeframe %dm is not valid on Binance", timeframe.Value)
		}
	case types.TuHour:
		switch timeframe.Value {
		case 1, 2, 4, 6, 8, 12:
			binanceTimeframe = fmt.Sprintf("%dh", timeframe.Value)
		default:
			err = fmt.Errorf("Timeframe %dh is not valid on Binance", timeframe.Value)
		}
	case types.TuDay:
		switch timeframe.Value {
		case 1, 3:
			binanceTimeframe = fmt.Sprintf("%dd", timeframe.Value)
		default:
			err = fmt.Errorf("Timeframe %dd is not valid on Binance", timeframe.Value)
		}
	case types.TuWeek:
		if timeframe.Value == 1 {
			binanceTimeframe = fmt.Sprintf("%dw", timeframe.Value)
		} else {
			err = fmt.Errorf("Timeframe %dw is not valid on Binance", timeframe.Value)
		}
	case types.TuMonth:
		if timeframe.Value == 1 {
			binanceTimeframe = fmt.Sprintf("%dM", timeframe.Value)
		} else {
			err = fmt.Errorf("Timeframe %dW is not valid on Binance", timeframe.Value)
		}
	}
	return
}

func (b Binance) sideToBinance(s types.Side) bin.SideType {
	if s == types.Buy {
		return bin.SideTypeBuy
	}
	return bin.SideTypeSell
}

func (b Binance) toSide(in bin.SideType) types.Side {
	if in == bin.SideTypeBuy {
		return types.Buy
	}
	return types.Sell
}

func (b Binance) orderTypeToBinance(t types.OrderType) bin.OrderType {
	switch t {
	case types.Limit:
		return bin.OrderTypeLimit
	case types.Market:
		return bin.OrderTypeMarket
	case types.StopLoss:
		return bin.OrderTypeStopLoss
	case types.StopLossLimit:
		return bin.OrderTypeStopLossLimit
	case types.TakeProfit:
		return bin.OrderTypeTakeProfit
	case types.TakeProfitLimit:
		return bin.OrderTypeTakeProfitLimit
	case types.LimitMaker:
		return bin.OrderTypeLimitMaker
	}
	return bin.OrderTypeMarket
}

func (b Binance) toOrderType(in bin.OrderType) types.OrderType {
	switch in {
	case bin.OrderTypeLimit:
		return types.Limit
	case bin.OrderTypeMarket:
		return types.Market
	case bin.OrderTypeStopLoss:
		return types.StopLoss
	case bin.OrderTypeStopLossLimit:
		return types.StopLossLimit
	case bin.OrderTypeTakeProfit:
		return types.TakeProfit
	case bin.OrderTypeTakeProfitLimit:
		return types.TakeProfitLimit
	case bin.OrderTypeLimitMaker:
		return types.LimitMaker
	}
	return types.Market
}

func (b Binance) timeInForceToBinance(t types.TimeInForce) bin.TimeInForceType {
	switch t {
	case types.GoodTillCancel:
		return bin.TimeInForceTypeGTC
	case types.ImmediateOrCancel:
		return bin.TimeInForceTypeIOC
	case types.FillOrCancel:
		return bin.TimeInForceTypeFOK
	}
	return bin.TimeInForceTypeGTC
}

func (b Binance) toTimeInForce(in bin.TimeInForceType) types.TimeInForce {
	switch in {
	case bin.TimeInForceTypeGTC:
		return types.GoodTillCancel
	case bin.TimeInForceTypeIOC:
		return types.ImmediateOrCancel
	case bin.TimeInForceTypeFOK:
		return types.FillOrCancel
	}
	return types.GoodTillCancel
}

func (b Binance) toFloat(in string) (flt float64) {
	var err error

	if flt, err = strconv.ParseFloat(in, 64); err != nil {
		logger.Warningf("Binance::toFloat Error %v\n", err)
		flt = math.NaN()
	}
	return
}

func (b Binance) toTime(in int64) (t time.Time) {
	var secs, nano int64

	secs = in / 1000
	nano = (in % 1000) * 1000000
	t = time.Unix(secs, nano)
	return
}

func (b Binance) toStatus(in bin.OrderStatusType) types.OrderStatus {
	switch in {
	case bin.OrderStatusTypeNew:
		return types.StatusNew
	case bin.OrderStatusTypePartiallyFilled:
		return types.StatusPartiallyFilled
	case bin.OrderStatusTypeFilled:
		return types.StatusFilled
	case bin.OrderStatusTypeCanceled:
		return types.StatusCanceled
	case bin.OrderStatusTypePendingCancel:
		return types.StatusPendingCancel
	case bin.OrderStatusTypeRejected:
		return types.StatusRejected
	case bin.OrderStatusTypeExpired:
		return types.StatusExpired
	}
	return types.StatusRejected
}
