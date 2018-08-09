package main

import "os"
import "time"

const (
	markets_DB  = "altcoin"
	count       = 10
	API_KEY     = ""
	API_SECRET  = ""

	interval = 10 * time.Second
)

var (
	// EXCHANGES
	_poloniex  = false
	_binance   = true
	_bittrex   = true
	_cryptopia = true

	// MARKETS
	_BTC  = true
	_USDT = true
	
	// LOGIN
	username   = os.Getenv("DBUSER")
	password   = os.Getenv("DBPASS")
	server_url = os.Getenv("SERVERDB")
)
