package main

import (
	"log"
	"net/http"

	jweb "github.com/nboughton/go-utils/json/web"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	MainTable  []TableRow `json:"mainTable"`
	TimeSeries graphData  `json:"timeSeries"`
	FreqDist   graphData  `json:"freqDist"`
}

// TableRow contains data used in the top table
type TableRow struct {
	Label string `json:"label"`
	Num   []int  `json:"num"`
}

func handlerQuery(w http.ResponseWriter, r *http.Request) {
	p := params(r)

	jweb.New(http.StatusOK,
		PageData{
			MainTable:  createMainTableData(p),
			TimeSeries: graphTimeSeries(db.getResults(p)),
			FreqDist:   graphFreqDist(db.getResults(p)),
		},
	).Write(w)
}

func handlerListSets(w http.ResponseWriter, r *http.Request) {
	res, err := db.getSetList(params(r))
	if err != nil {
		jweb.New(http.StatusInternalServerError, err).Write(w)
	} else {
		jweb.New(http.StatusOK, res).Write(w)
	}
}

func handlerListMachines(w http.ResponseWriter, r *http.Request) {
	res, err := db.getMachineList(params(r))
	if err != nil {
		jweb.New(http.StatusInternalServerError, err).Write(w)
	} else {
		jweb.New(http.StatusOK, res).Write(w)
	}
}

func handlerDataRange(w http.ResponseWriter, r *http.Request) {
	f, l, err := db.getDataRange()
	if err != nil {
		log.Println("handlerDataRange:", err.Error())
	}

	jweb.New(http.StatusOK, map[string]int64{"first": f.Unix(), "last": l.Unix()}).Write(w)
}
