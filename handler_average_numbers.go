package main

import (
	"strconv"
	//"time"

	"github.com/pilu/traffic"
)

var (
	jsTimeFormat = "2006-01-02"
)

func handlerAverageNumbers(w traffic.ResponseWriter, r *traffic.Request) {
	var p queryParams
	p.Start = r.Param("start")
	p.End = r.Param("end")
	p.Set, _ = strconv.Atoi(r.Param("set"))
	p.Machine = r.Param("machine")

	res, err := db.getAverageNumbers(p)
	if err != nil {
		//w.WriteJSON("Invalid machine/set combination")
		w.WriteJSON(err.Error())
	} else {
		w.WriteJSON(res)
	}
}
