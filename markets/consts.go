package main

import "os"
import "time"

const (
	markets_DB = "altcoin"
	count      = 10
	apiKey     = ""
	apiPass    = ""

	interval = 10 * time.Second
)

var (
	// EXCHANGES
	vPoloniex  = false
	vBinance   = true
	vBittrex   = true
	vCryptopia = true

	// MARKETS
	vBTC  = true
	vUSDT = true

	// LOGIN
	username  = os.Getenv("DBUSER")
	password  = os.Getenv("DBPASS")
	serverURL = os.Getenv("SERVERDB")
)
