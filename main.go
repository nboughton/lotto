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
	"github.com/nboughton/lotto/db"
	"github.com/nboughton/lotto/handler"
)

type config struct {
	Port int
	Log  bool
}

func main() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	l := flag.Bool("l", true, "Log requests to STDOUT")
	flag.Parse()

	c := config{*p, *l}
	e := &handler.Env{DB: db.Connect("./results.db")}

	// Update at 21:30 every night
	go func() {
		for t := range time.NewTicker(time.Minute).C {
			if t.Hour() == 21 && t.Minute() == 30 {
				e.DB.Update()
			}
		}
	}()

	r := mux.NewRouter()

	r.Handle("/sets", handler.ListSets(e)).Methods("GET")
	r.Handle("/machines", handler.ListMachines(e)).Methods("GET")
	r.Handle("/query", handler.Query(e)).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	if c.Log {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", c.Port),
				handlers.LoggingHandler(os.Stdout, r),
			),
		)
	} else {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", c.Port),
				r,
			),
		)
	}

}
