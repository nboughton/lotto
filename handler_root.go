package main

import (
	"log"

	"github.com/pilu/traffic"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	Machines []string
	Sets     []int
}

func handlerRoot(w traffic.ResponseWriter, r *traffic.Request) {
	sets, err := db.getSetList()
	if err != nil {
		log.Println(err)
	}

	machines, err := db.getMachineList()
	if err != nil {
		log.Println(err)
	}

	w.Render("index", &PageData{Machines: machines, Sets: sets})
}
