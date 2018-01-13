package main

import (

	"time"
	"github.com/stock-simulator-server/app"
)

func main() {
	app.RunApp()
	for{
		time.Sleep(1 * time.Second)
	}
}
