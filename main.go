package main

import (
	"flag"

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
	router.Get("/api/results", handlerResults)
	router.Get("/api/results/average", handlerResultsAverage)
	router.Get("/api/sets", handlerListSets)
	//router.Get("/api/machines", handlerListMachines) // Not currently in use
}

func main() {
	router.Run()
}
