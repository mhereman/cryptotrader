package types

import (
	"fmt"
	"strconv"
	"strings"
)

type SymbolInfo struct {
	Symbol         Symbol
	MinPrice       string
	MinLotQuantity string

	priceDecimalPlaces int
	priceInc           int
}

func NewSymbolInfo(symbol Symbol, minPriceDecimalPlaces string, minLotQuantity string) (si SymbolInfo) {
	si = SymbolInfo{
		Symbol:         symbol,
		MinPrice:       minPriceDecimalPlaces,
		MinLotQuantity: minLotQuantity,
	}
	si.calculateLimits()
	return si
}

func (si *SymbolInfo) calculateLimits() {
	var trimmed, integerPart, decimalPart string
	var parts []string

	trimmed = strings.Trim(si.MinPrice, "0")
	parts = strings.Split(trimmed, ".")

	integerPart = parts[0]
	decimalPart = parts[1]

	si.priceDecimalPlaces = len(decimalPart)
	if si.priceDecimalPlaces == 0 {
		si.priceDecimalPlaces = len(integerPart) * -1
		si.priceInc, _ = strconv.Atoi(string(integerPart[0]))
	} else {
		si.priceInc, _ = strconv.Atoi(string(decimalPart[len(decimalPart)-1]))
	}
}

func (si SymbolInfo) ClampPrice(price float64) (clampedPrice string, err error) {
	var format string
	var length, digit, decPlac int

	if si.priceDecimalPlaces > 0 {
		format = fmt.Sprintf("%%.%df", si.priceDecimalPlaces)
		clampedPrice = fmt.Sprintf(format, price)
		length = len(clampedPrice)
		if digit, err = strconv.Atoi(string(clampedPrice[length-1])); err != nil {
			return
		}
		clampedPrice = clampedPrice[:length-1] + fmt.Sprintf("%d", (digit/si.priceInc)*si.priceInc)
		return
	}

	format = "%.f"
	clampedPrice = fmt.Sprintf(format, price)
	length = len(clampedPrice)
	decPlac = si.priceDecimalPlaces * -1

	format = fmt.Sprintf("%%0%dd", decPlac-1)
	clampedPrice = (clampedPrice[:length-(decPlac-1)] + fmt.Sprintf(format, 0))[:length]
	if digit, err = strconv.Atoi(string(clampedPrice[length-decPlac : length-decPlac+1])); err != nil {
		return
	}
	clampedPrice = clampedPrice[:length-decPlac] + fmt.Sprintf("%d", ((digit+(si.priceInc-1))/si.priceInc)*si.priceInc) + clampedPrice[length-decPlac+1:]
	return
}
