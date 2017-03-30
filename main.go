package main

import (
	"flag"
	"time"

	"github.com/pilu/traffic"
)

var (
	router *traffic.Router
	db     = connectDB("results.db")
)

func init() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	flag.Parse()

	traffic.SetPort(*p)

	router = traffic.New()
	router.Get("/", handlerRoot)
	router.Get("/api/range", handlerDataRange)
	router.Get("/api/results/average", handlerResultsAverage)
	router.Get("/api/results/graph/freqdist/:type", handlerResultsFreqDist)
	router.Get("/api/results/graph/timeseries/:type", handlerResultsTimeSeries)
	router.Get("/api/results/graph/3d/scatter", handlerResultsScatter3D)
	router.Get("/api/sets", handlerListSets)
	router.Get("/api/machines", handlerListMachines)
	//router.Get("/api/machines/sets/combos", handlerMachineSetsCombos)

	// Update every 24 hours
	go func() {
		for t := range time.NewTicker(time.Minute).C {
			if t.Hour() == 0 && t.Minute() == 0 {
				db.updateDB()
			}
		}
	}()
}

func main() {
	router.Run()
}
