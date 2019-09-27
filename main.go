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
	"github.com/nboughton/lotto/handler"
	"github.com/nboughton/stalotto/db"
)

func main() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	l := flag.Bool("l", true, "Log requests to STDOUT")
	flag.Parse()

	e := &handler.Env{DB: db.Connect("./results.db")}
	e.DB.Update()

	// Update at 21:30 every night
	go func() {
		for t := range time.NewTicker(time.Minute).C {
			if t.Hour() == 21 && t.Minute() == 30 {
				e.DB.Update()
			}
		}
	}()

	log.Println("Configuring router")
	r := mux.NewRouter()
	r.Handle("/sets", handler.ListSets(e)).Methods("POST")
	r.Handle("/machines", handler.ListMachines(e)).Methods("POST")
	r.Handle("/query", handler.Query(e)).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))

	log.Println("Starting server")
	if *l {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", *p),
				handlers.LoggingHandler(os.Stdout, r),
			),
		)
	} else {
		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", *p),
				r,
			),
		)
	}

}
