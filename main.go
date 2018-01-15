package main

import (

	"time"
	"github.com/stock-simulator-server/src/app"
	"github.com/stock-simulator-server/src/web"
)

func main() {
	go app.RunApp()
	web.StartHandlers()

	for{
		time.Sleep(1 * time.Second)
	}
}
