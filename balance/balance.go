package main

import (
	//"fmt"
	"log"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/toorop/go-bittrex"
)

const (
	balanceDB = "balance"
	username  = os.Getenv("DBUSER")
	password  = os.Getenv("DBPASS")
	serverURL = os.Getenv("SERVERDB")
	apiKey    = os.Getenv("BITTREX_API_KEY")
	apiPass   = os.Getenv("BITTREX_API_SECRET")
)

func main() {

	var BTC, USDT, total float64

	// BITTREX
	bittrex := bittrex.New(apiKey, apiPass)

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     serverURL,
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  balanceDB,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	balances, err := bittrex.GetBalances()
	if err != nil {
		log.Fatal(err)
	}

	btcVal, err := bittrex.GetTicker("USDT-BTC")
	if err != nil {
		log.Fatal(err)
	}

	for _, coin := range balances {
		if coin.Currency == "BTC" {
			BTC = coin.Balance
			total += BTC
		} else if coin.Currency == "USDT" {
			USDT = coin.Balance
			total += USDT / btcVal.Last
		} else if coin.Balance > 0 {
			//fmt.Println(coin.Currency, coin.Balance)
			//fmt.Println(bittrex.GetTicker("BTC-"+coin.Currency))
			val, err := bittrex.GetTicker("BTC-" + coin.Currency)
			if err != nil {
				log.Fatal()
			}

			// Create a point and add to batch
			tags := map[string]string{"Coin": coin.Currency}
			fields := map[string]interface{}{
				"Balance": coin.Balance,
				"In BTC":  coin.Balance * val.Last,
			}

			//fmt.Println(err, marketSummaries)
			pt, err := client.NewPoint("balance", tags, fields, time.Now())
			if err != nil {
				log.Println(err)
			}
			bp.AddPoint(pt)

			// Write the batch
			if err := c.Write(bp); err != nil {
				log.Println(err)
			}

			total += val.Last * coin.Balance
		}
	}

	// Create a point and add to batch
	tags := map[string]string{}
	fields := map[string]interface{}{
		"Balance":   total,
		"BTC Value": btcVal.Last,
		"Total USD": btcVal.Last * total,
	}

	//log.Println(err, marketSummaries)
	pt, err := client.NewPoint("summary", tags, fields, time.Now())
	if err != nil {
		log.Println(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Println(err)
	}
}
