package main

import (
	"log"

	"github.com/pilu/traffic"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	Machines   []string
	Sets       []int
	Start, End string
}

func handlerRoot(w traffic.ResponseWriter, r *traffic.Request) {
	sets, err := db.getSetList("all")
	if err != nil {
		log.Println(err)
	}

	machines, err := db.getMachineList(0)
	if err != nil {
		log.Println(err)
	}

	s, e, err := db.getDataRange()
	if err != nil {
		log.Println(err)
	}

	w.Render("index", &PageData{machines, sets, s.String(), e.String()})
}
