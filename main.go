package main

import (
	"encoding/json"
	"fmt"
	"stock-server/wallstreet"
	"time"
)

func main() {

	stockManager := wallstreet.NewStockManager()
	go func() {
		for {
			value, _ := json.Marshal(<-stockManager.StockUpdateChannel)
			fmt.Println(string(value))
		}
	}()

	stockManager.AddStock("CHUNT", "Conner Hunt's Hats", 420.69)
	stockManager.StartSimulateStocks(3 * time.Second)

	for {
		time.Sleep(100 * time.Second)
	}

}
