package main

const (
	// SatInBTC represents number of satoshis in 1 bitcoin
	SatInBTC = uint64(100000000)

	// PricesURL is URL for crypo prices
	PricesURL = "https://min-api.cryptocompare.com/data/price?fsym=WAVES&tsyms=BTC,ETH,EUR,HRK"

	// TelAnonOps group for error logging
	TelAnonOps = -1001213539865

	// TelPollerTimeout is Telegram poller timeout in seconds
	TelPollerTimeout = 30

	// TokenID is funding token Waves ID
	TokenID = "66DUhUoJaoZcstkKpcoN3FUcqjB6v8VJd5ZQd6RsPxhv"

	// TokenAddress is funding token Waves address
	TokenAddress = "3PBmmxKhFcDhb8PrDdCdvw2iGMPnp7VuwPy"

	// WavesMonitorTick interval in seconds
	WavesMonitorTick = 10
)
