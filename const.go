package main

const (
	// SatInBTC represents number of satoshis in 1 bitcoin
	SatInBTC = uint64(100000000)

	// PricesURL is URL for crypo prices
	PricesURL = "https://min-api.cryptocompare.com/data/price?fsym=WAVES&tsyms=BTC,ETH,EUR,HRK"

	// TelAnonOps group for error logging
	TelAnonOps = -1001213539865

	// TelAnonTeam group for team messages
	TelAnonTeam = -1001280228955

	// TelPollerTimeout is Telegram poller timeout in seconds
	TelPollerTimeout = 30

	// TokenID is funding token Waves ID
	TokenID = "66DUhUoJaoZcstkKpcoN3FUcqjB6v8VJd5ZQd6RsPxhv"

	// TokenAddress is funding token Waves address
	TokenAddress = "3PBmmxKhFcDhb8PrDdCdvw2iGMPnp7VuwPy"

	// WavesMonitorTick interval in seconds
	WavesMonitorTick = 10

	// WavesNodeURL is an URL for Waves Node
	WavesNodeURL = "https://nodes.wavesnodes.com"

	// MatcherNodeURL is an URL for Waves Node
	MatcherNodeURL = "https://matcher.waves.exchange"

	// MatcherPublicKey represents Waves matcher public key
	MatcherPublicKey = "9cpfKN9suPNvfeUNphzxXMjcnn974eme8ZhWUjaktzU5"

	// WavesFee is Waves regular fee amount
	WavesFee = 100000

	// WavesExchangeFee is Waves exchange fee amount
	WavesExchangeFee = 300000

	// Port represents a port the app will listen on
	Port = 5002

	// AHRKId is AHRK asset id
	AHRKId = string("Gvs59WEEXVAQiRZwisUosG7fVNr8vnzS8mjkgqotrERT")

	// AHRKDec represents number of decimals in AHRK
	AHRKDec = uint64(1000000)
)
