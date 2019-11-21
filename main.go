package main

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/app"
)

func main() {
	//flag.Parse()
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Log.Fatal("could not create CPU profile: ", err)
	//	}
	//	if err := pprof.StartCPUProfile(f); err != nil {
	//		log.Log.Fatal("could not start CPU profile: ", err)
	//	}
	//	defer pprof.StopCPUProfile()
	//}
	app.App()
	//select {}
	//if *memprofile != "" {
	//	f, err := os.Create(*memprofile)
	//	if err != nil {
	//		log.Log.Fatal("could not create memory profile: ", err)
	//	}
	//	runtime.GC() // get up-to-date statistics
	//	if err := pprof.WriteHeapProfile(f); err != nil {
	//		log.Log.Fatal("could not write memory profile: ", err)
	//	}
	//	f.Close()
	//}

}
