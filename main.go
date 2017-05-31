package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"fmt"
	"github.com/gorilla/mux"
)

type config struct {
	Port int
}

var (
	cfg config
	db  = connectDB("results.db")
)

func init() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	flag.Parse()

	cfg.Port = *p

	// Update at 21:30 every night
	go func() {
		for t := range time.NewTicker(time.Minute).C {
			if t.Hour() == 21 && t.Minute() == 30 {
				db.updateDB()
			}
		}
	}()
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/range", handlerDataRange).Methods("GET")
	r.HandleFunc("/api/sets", handlerListSets).Methods("GET")
	r.HandleFunc("/api/machines", handlerListMachines).Methods("GET")
	r.HandleFunc("/api/query", handlerQuery).Methods("GET")
	/*
		r.HandleFunc("/api/numbers", handlerNumbers).Methods("GET")
		r.HandleFunc("/api/graph/results/freqdist/{type}", handlerResultsFreqDist).Methods("GET")
		r.HandleFunc("/api/graph/results/timeseries/{type}", handlerResultsTimeSeries).Methods("GET")
		r.HandleFunc("/api/graph/results/raw/scatter3d", handlerResultsScatter3D).Methods("GET")
		r.HandleFunc("/api/graph/ms/freqdist/{type}", handlerMSFreqDist).Methods("GET")
	*/
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
