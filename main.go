package main

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/stock-simulator-server/src/run"

	"github.com/stock-simulator-server/src/log"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func RunApp() {

}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	run.App()
	select {}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}

}
