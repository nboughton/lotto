package main

import (
	"log"

	"github.com/pilu/traffic"
)

func handlerDataRange(w traffic.ResponseWriter, r *traffic.Request) {
	first, last, err := db.getDataRange()
	if err != nil {
		log.Println("handlerDataRange:", err.Error())
	}

	w.WriteJSON(map[string]int64{"first": first.Unix(), "last": last.Unix()})
}
