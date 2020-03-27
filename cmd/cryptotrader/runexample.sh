#!/usr/bin/env bash

# Logging Configuration
#######################

# The log level.
# Can be: debug, error, warning, info, none
LOGLEVEL='info'

# Exchange Configuration
########################

# The exchange to use.
# Currently only Binance is supported.
EXCHANGE='binance'

# Your Binance API Key.
API_KEY=''

# Your Binance API Secret.
API_SECRET=''


# Asset Configuration
#####################

# The Base Asset.
BASE_ASSET='btc'

# The Quote Asset.
QUOTE_ASSET='usdt'

# The timeframe the bot should work on.
# This defines the timeframe of the candles to retrieve and the update interval.
TIME_FRAME='4h'


# Algorithm Configuration
#########################

# The name of the algorithm to use.
# Currently only Ema/Sma is supported.
ALGO='Ema/Sma'

# The algorithm arguments
# Ema/Sma.sma_len = 14
# Ema/Sma.ema_len = 7
# Ema/Sma.rsi_len = 14
# Ema/Sma.rsi_buy_max = 90.0
ALGO_ARGS='Ema/Sma.sma_len=14;Ema/Sma.ema_len=7;Ema/Sma.rsi_len=14;Ema/Sma.rsi_buy_max=90.0;Ema/Sma.backtest=false'


# Trade Configuration
#####################

# Type of trade to perform.
# fixed uses a fixed amount of quote asset to trade, configure the ammount in the volume setting.
# percent uses a percent of the availabel quote asset to trade, configure the ammount in the volume setting.
TRADE_TYPE='fixed'

# The amount to trade.
# If the trade type is fixed, this is an absolute number expressed in the quote asset.
# If the trade type is percent, this is a percentage of the available quote asset, the max value then = 1.0.
VOLUME='100.0'

# Reduce the trade amount if inssufficient funds are available.
# This is only used in combination with the fixed trade type.
REDUCE='true'

# Perform paper trading
# If this is set to true the trades will not be executed but a message is printed to the console.
PAPER_TRADING='true'

# Max slippage on buy orders
# Use 0 to disable
MAX_SLIPPAGE='0.001'



# Notifier Configuration
########################

# The notifier to enable
# If empty string, no notifier is configured
# At the moment only the Proximus SMS api is available as a notifier ('proximus-sms')
# The configuration of the Proximus SMS api needs a 'apiToken' and 'destination' entry
NOTIFIER=''

# The configuration arguments for the notifier
NOTIFIER_CONFIG=''



# Run cryptotrader
cryptotrader \
    -loglevel=${LOGLEVEL} \
    -exchange=${EXCHANGE} \
    -exchangeargs="apiKey=${API_KEY};apiSecret=${API_SECRET}" \
    -base=${BASE_ASSET} \
    -quote=${QUOTE_ASSET} \
    -timeframe=${TIME_FRAME} \
    -algo=${ALGO} \
    -algoargs=${ALGO_ARGS} \
    -tradetype=${TRADE_TYPE} \
    -volume=${VOLUME} \
    -reduce=${REDUCE} \
    -papertrading=${PAPER_TRADING} \
    -maxslippage=${MAX_SLIPPAGE} \
    -notifier=${NOTIFIER} \
    -notifierargs=${NOTIFIER_CONFIG}
