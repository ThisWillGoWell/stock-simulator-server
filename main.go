package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/app"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/database"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/session"
	"github.com/stock-simulator-server/src/valuable"
	"github.com/stock-simulator-server/src/web"
	"github.com/stock-simulator-server/src/wires"
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
	disableDb := os.Getenv("DISABLE_DB") == "True"
	serveLog := os.Getenv("SERVE_LOG") == "True"
	autoLoad := os.Getenv("AUTO_LOAD") == "True"

	//start DB
	if !disableDb {
		database.InitDatabase()
	}
	if serveLog {
		filepath := os.Getenv("FILE_SERVE")

		web.ServePath(filepath)
	}
	//valuable.ValuablesLock.EnableDebug()
	//ledger.EntriesLock.EnableDebug()

	//Wiring of system
	wires.ConnectWires(disableDb)
	//this takes the subscribe output and converts it to a message
	change.StartDetectChanges()
	session.StartSessionCleaner()
	account.RunUserSend()

	order.Run()
	valuable.StartStockStimulation()

	if autoLoad {
		go app.LoadConfig()
	}
	//go app.LoadVars()
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
