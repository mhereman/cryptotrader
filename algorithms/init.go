package algorithms

import (
	"fmt"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
)

var algoFactory map[string]func() (interfaces.IAlgorithm, error) = make(map[string]func() (interfaces.IAlgorithm, error))

func RegisterAlgorithm(name string, factory func() (interfaces.IAlgorithm, error)) {
	algoFactory[name] = factory
	logger.Printf("Registered algorithm: %s\n", name)
}

func GetAlgorithm(name string) (algo interfaces.IAlgorithm, err error) {
	var ok bool
	var fn func() (interfaces.IAlgorithm, error)

	if fn, ok = algoFactory[name]; !ok {
		err = fmt.Errorf("Algorithm %s does not exist")
		return
	}

	algo, err = fn()
	return
}
