package types

import (
	"fmt"
	"strings"
)

type Symbol struct {
	base  string
	quote string
}

func NewSymbol(b, q string) Symbol {
	return Symbol{
		base:  strings.ToUpper(b),
		quote: strings.ToUpper(q),
	}
}

func NewSymbolFromString(in string) (sym Symbol, err error) {
	var parts []string
	parts = strings.Split(in, "/")

	if len(parts) != 2 {
		err = fmt.Errorf("Invalid symbol string")
		return
	}

	sym = Symbol{
		base:  parts[0],
		quote: parts[1],
	}
	return
}

func (s Symbol) Base() string {
	return s.base
}

func (s Symbol) Quote() string {
	return s.quote
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s/%s", s.base, s.quote)
}
