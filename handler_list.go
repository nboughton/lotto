package main

import (
	"github.com/pilu/traffic"
)

func handlerListSets(w traffic.ResponseWriter, r *traffic.Request) {
	res, _ := db.getSetList(parseQueryParams(r))
	w.WriteJSON(res)
}

func handlerListMachines(w traffic.ResponseWriter, r *traffic.Request) {
	res, _ := db.getMachineList(parseQueryParams(r))
	w.WriteJSON(res)
}
