package main

import (
	"github.com/pilu/traffic"
)

func handlerListSets(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getSetList(r.Param("machine"))
	if err != nil {
		res, _ = db.getSetList("all")
	}

	w.WriteJSON(res)
}
