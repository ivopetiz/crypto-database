package main

import (
	//"fmt"
	"log"
	"time"

	"github.com/toorop/go-bittrex"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	balance_DB = "balance"
	username   = os.Getenv("DBUSER")
	password   = os.Getenv("DBPASS")
	server_url = os.Getenv("SERVERDB")
	API_KEY    = os.Getenv("BITTREX_API_KEY")
	API_SECRET = os.Getenv("BITTREX_API_SECRET")
)


func main() {

	var BTC, USDT, total float64

	// BITTREX
	bittrex := bittrex.New(API_KEY, API_SECRET)

	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:	  server_url,
		Username: username,
		Password: password,
	})
	if err != nil {log.Fatal(err)}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  balance_DB,
		Precision: "s",
	})
	if err != nil {log.Fatal(err)}

	balances, err := bittrex.GetBalances()
	if err != nil {log.Fatal(err)}

	btc_val, err := bittrex.GetTicker("USDT-BTC")
	if err != nil {log.Fatal(err)}

	for _, coin := range balances {
		if coin.Currency == "BTC" {
			BTC = coin.Balance
			total += BTC
		} else if coin.Currency == "USDT" {
			USDT = coin.Balance
			total += USDT/btc_val.Last
		} else if coin.Balance > 0 {
			//fmt.Println(coin.Currency, coin.Balance)
			//fmt.Println(bittrex.GetTicker("BTC-"+coin.Currency))
			val, err := bittrex.GetTicker("BTC-"+coin.Currency)
			if err != nil {log.Fatal()}

			// Create a point and add to batch
			tags := map[string]string{"Coin": coin.Currency}
			fields := map[string]interface{}{
				"Balance": coin.Balance ,
				"In BTC": coin.Balance*val.Last ,
			}

			//fmt.Println(err, marketSummaries)
			pt, err := client.NewPoint("balance", tags, fields, time.Now())
			if err != nil {log.Println(err)}
			bp.AddPoint(pt)

			// Write the batch
			if err := c.Write(bp); err != nil {log.Println(err)}
				
			total += val.Last * coin.Balance
		}
	}

	// Create a point and add to batch
	tags := map[string]string{}
	fields := map[string]interface{}{
		"Balance": total ,
		"BTC Value": btc_val.Last ,
		"Total USD": btc_val.Last * total,
	}

	//log.Println(err, marketSummaries)
	pt, err := client.NewPoint("summary", tags, fields, time.Now())
	if err != nil {log.Println(err)}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {log.Println(err)}
}