package main

import (

	"time"
	"github.com/stock-simulator-server/app"
	"github.com/stock-simulator-server/web"
	"github.com/stock-simulator-server/account"
)

func main() {

	go app.RunApp()
	web.StartHandlers()

	for{
		time.Sleep(1 * time.Second)
	}
}
