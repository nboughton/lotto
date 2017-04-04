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
	router.Get("/api/sets", handlerListSets)
	router.Get("/api/machines", handlerListMachines)
	router.Get("/api/numbers", handlerNumbers)
	router.Get("/api/graph/results/freqdist/:type", handlerResultsFreqDist)
	router.Get("/api/graph/results/timeseries/:type", handlerResultsTimeSeries)
	router.Get("/api/graph/results/raw/scatter3d", handlerResultsScatter3D)
	router.Get("/api/graph/ms/freqdist/:type", handlerMSFreqDist)

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
