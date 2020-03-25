package types

import (
	"fmt"
	"strings"
)

// Symbol represents a tradable asset pair
type Symbol struct {
	base  string
	quote string
}

// NewSymbol creates a Symbol instance
func NewSymbol(baseAsset string, quoteAsset string) Symbol {
	return Symbol{
		base:  strings.ToUpper(baseAsset),
		quote: strings.ToUpper(quoteAsset),
	}
}

// NewSymbolFromString creates a Symbol instance from a string
// The string format should be:
//	base/quote
func NewSymbolFromString(in string) (sym Symbol, err error) {
	var parts []string
	parts = strings.Split(in, "/")

	if len(parts) != 2 {
		err = fmt.Errorf("Invalid symbol string")
		return
	}

	sym = Symbol{
		base:  strings.ToUpper(parts[0]),
		quote: strings.ToUpper(parts[1]),
	}
	return
}

// Base returns the base asset
func (s Symbol) Base() string {
	return s.base
}

// Quote returns the quote asset
func (s Symbol) Quote() string {
	return s.quote
}

// String returns the string version of the symbol
func (s Symbol) String() string {
	return fmt.Sprintf("%s/%s", s.base, s.quote)
}
