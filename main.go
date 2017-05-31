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

	r.HandleFunc("/api/sets", handlerListSets).Methods("GET")
	r.HandleFunc("/api/machines", handlerListMachines).Methods("GET")
	r.HandleFunc("/api/query", handlerQuery).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r))
}
