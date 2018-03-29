package main

import (
	"flag"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/database"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/trade"
	"github.com/stock-simulator-server/src/valuable"
	"github.com/stock-simulator-server/src/web"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//start DB
	database.InitDatabase()
	//Wiring of system
	change.SubscribeUpdateInputs.RegisterInput(portfolio.PortfoliosUpdateChannel.GetBufferedOutput(10))
	change.SubscribeUpdateInputs.RegisterInput(ledger.EntriesUpdate.GetBufferedOutput(100))
	change.SubscribeUpdateInputs.RegisterInput(valuable.ValuableUpdateChannel.GetBufferedOutput(10))

	//this takes the subscribe output and converts it to a message
	client.BroadcastMessageBuilder()
	change.StartDetectChanges()

	trade.RunTrader()
	valuable.StartStockStimulation()

	// go app.RunApp()
	go web.StartHandlers()
	select {}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

}
