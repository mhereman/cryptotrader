package main

import (
	"log"

	"github.com/mhereman/cryptotrader"

	"github.com/mhereman/cryptotrader/logger"
)

func main() {
	var assetCfg cryptotrader.AssetConfig
	var exchangeCfg cryptotrader.ExchangeConfig
	var algoCfg cryptotrader.AlgorithmConfig
	var tradeCfg cryptotrader.TradeConfig
	var notifierCfg cryptotrader.NotifierConfig
	var trader *cryptotrader.CryptoTrader
	var err error

	if assetCfg, exchangeCfg, algoCfg, tradeCfg, notifierCfg, err = cryptotrader.ReadFlags(); err != nil {
		log.Fatalf("Error %v\n", err)
	}

	logger.Debugf("API_KEY: %s", exchangeCfg.ArgMap["apiKey"])
	logger.Debugf("API_SECRET: %s", exchangeCfg.ArgMap["apiSecret"])

	logger.Infoln("Starting cryptotrader")
	trader = cryptotrader.New(assetCfg, exchangeCfg, algoCfg, tradeCfg, notifierCfg)
	if err = trader.Run(); err != nil {
		logger.Fatalf("Error %v\n", err)
	}
}
