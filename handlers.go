package main

import (
	"log"
	"strconv"
	"time"

	"github.com/pilu/traffic"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	Machines   []string
	Sets       []int
	Start, End string
}

func handlerRoot(w traffic.ResponseWriter, r *traffic.Request) {
	s, e, err := db.getDataRange()
	if err != nil {
		log.Println(err)
	}

	s, _ = time.Parse(formatYYYYMMDD, "2015-10-10") // Because reasons.
	q := queryParams{
		Start:   s.Format(formatYYYYMMDD),
		End:     e.Format(formatYYYYMMDD),
		Machine: "all",
		Set:     0,
	}

	sets, err := db.getSetList(q)
	if err != nil {
		log.Println(err)
	}

	machines, err := db.getMachineList(q)
	if err != nil {
		log.Println(err)
	}

	w.Render("index", &PageData{machines, sets, q.Start, q.End})
}

func handlerResultsScatterGraph(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphScatter(db.getResults(parseQueryParams(r)), true))
}

func handlerResultsScatter3DGraph(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphScatter3D(db.getResults(parseQueryParams(r))))
}

func handlerResultsBarGraph(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphBar(db.getResults(parseQueryParams(r))))
}

func handlerMachineSetsCombos(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(db.getMachineSetCombinations(parseQueryParams(r)))
}

func handlerResultsAverage(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getResultsAverage(parseQueryParams(r))
	if err != nil {
		//w.WriteJSON("Invalid machine/set combination")
		w.WriteJSON(err.Error())
	} else {
		w.WriteJSON(res)
	}
}

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

func handlerDataRange(w traffic.ResponseWriter, r *traffic.Request) {
	f, l, err := db.getDataRange()
	if err != nil {
		log.Println("handlerDataRange:", err.Error())
	}

	w.WriteJSON(map[string]int64{"first": f.Unix(), "last": l.Unix()})
}

func parseQueryParams(r *traffic.Request) queryParams {
	var p queryParams
	p.Start = r.Param("start")
	p.End = r.Param("end")
	p.Set, _ = strconv.Atoi(r.Param("set"))
	p.Machine = r.Param("machine")

	return p
}
