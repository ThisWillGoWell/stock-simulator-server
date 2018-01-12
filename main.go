package main

import (

	"time"
	"stock-server/app"
)

func main() {
	app.RunApp()
	for{
		time.Sleep(1 * time.Second)
	}
}
