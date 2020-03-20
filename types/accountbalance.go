package types

import "strings"

type AccountBalance struct {
	Asset  string
	Free   float64
	Locked float64
}

func NewAccountBalance(asset string, free float64, locked float64) AccountBalance {
	return AccountBalance{
		Asset:  strings.ToUpper(asset),
		Free:   free,
		Locked: locked,
	}
}
