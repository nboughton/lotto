package main

import (
	"github.com/pilu/traffic"
)

func handlerListSets(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getSetList(parseQueryParams(r))
	if err != nil {
		w.WriteJSON(err)
	} else {
		w.WriteJSON(res)
	}
}

func handlerListMachines(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getMachineList(parseQueryParams(r))
	if err != nil {
		w.WriteJSON(err)
	} else {
		w.WriteJSON(res)
	}
}
