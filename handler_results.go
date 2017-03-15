package main

import (
	"strconv"

	"github.com/pilu/traffic"
)

func handlerResults(w traffic.ResponseWriter, r *traffic.Request) {
	p := parseQueryParams(r)
	res := []dbRow{}

	for row := range db.getResults(p) {
		res = append(res, row)
	}

	w.WriteJSON(res)
}

func handlerResultsLineGraph(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphLine(db.getResults(parseQueryParams(r))))
}

func handlerResultsScatter3DGraph(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphScatter3D(db.getResults(parseQueryParams(r))))
}

func handlerResultsAverage(w traffic.ResponseWriter, r *traffic.Request) {
	p := parseQueryParams(r)

	res, err := db.getResultsAverage(p)
	if err != nil {
		//w.WriteJSON("Invalid machine/set combination")
		w.WriteJSON(err.Error())
	} else {
		w.WriteJSON(res)
	}
}

func parseQueryParams(r *traffic.Request) queryParams {
	var p queryParams
	p.Start = r.Param("start")
	p.End = r.Param("end")
	p.Set, _ = strconv.Atoi(r.Param("set"))
	p.Machine = r.Param("machine")

	return p
}
