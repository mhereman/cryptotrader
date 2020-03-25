package types

import "strings"

// AccountBalance represents the current account balance of an asset
type AccountBalance struct {
	// The asset
	Asset string

	// The free balance
	Free float64

	// The locked balance (can be locked because of it beining used in an open order, withdrawal, etc...)
	Locked float64
}

// NewAccountBalance creates a new AccountBalance instance
func NewAccountBalance(asset string, free float64, locked float64) AccountBalance {
	return AccountBalance{
		Asset:  strings.ToUpper(asset),
		Free:   free,
		Locked: locked,
	}
}
