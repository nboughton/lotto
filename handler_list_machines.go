package main

import (
	"strconv"

	"github.com/pilu/traffic"
)

func handlerListMachines(w traffic.ResponseWriter, r *traffic.Request) {
	set, _ := strconv.Atoi(r.Param("set"))
	res, _ := db.getMachineList(set)
	w.WriteJSON(res)
}
