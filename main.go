package main

import (
	"encoding/json"
	"fmt"
	"stock-server/wallstreet"
	"time"
)

func main() {
	exchange:= wallstreet.BuildExchange(time.Second)
	exchange.AddStock("CHUNT", "Chunt's Hats", 69, 420, time.Second * 10)
	exchange.AddStock("KING", "Paddle King", 10, 100, time.Second * 2)
	exchange.AddStock("CBIO", "Sebio's Streaming Services", 10, 100, time.Second)

	exchange.AddPortfolio()

	exchange.RunExchange()


}
