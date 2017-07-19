package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type config struct {
	Port int
	Log  bool
}

var (
	cfg config
	db  = connectDB("results.db")
)

func init() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	l := flag.Bool("l", true, "Log requests to STDOUT")
	flag.Parse()

	cfg.Port = *p
	cfg.Log = *l

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

	r.HandleFunc("/sets", handlerListSets).Methods("GET")
	r.HandleFunc("/machines", handlerListMachines).Methods("GET")
	r.HandleFunc("/query", handlerQuery).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	if cfg.Log {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", cfg.Port),
				handlers.LoggingHandler(os.Stdout, r),
			),
		)
	} else {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", cfg.Port),
				r,
			),
		)
	}

}
