package main

import (
	"log"
	"net/http"
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

	JSON{
		Status: http.StatusOK,
		Data: PageData{
			MainTable:  createMainTableData(p),
			TimeSeries: graphTimeSeries(db.getResults(p)),
			FreqDist:   graphFreqDist(db.getResults(p)),
		},
	}.write(w)
}

func handlerListSets(w http.ResponseWriter, r *http.Request) {
	res, err := db.getSetList(params(r))
	if err != nil {
		JSON{Status: http.StatusInternalServerError, Data: err}.write(w)
	} else {
		JSON{Status: http.StatusOK, Data: res}.write(w)
	}
}

func handlerListMachines(w http.ResponseWriter, r *http.Request) {
	res, err := db.getMachineList(params(r))
	if err != nil {
		JSON{Status: http.StatusInternalServerError, Data: err}.write(w)
	} else {
		JSON{Status: http.StatusOK, Data: res}.write(w)
	}
}

func handlerDataRange(w http.ResponseWriter, r *http.Request) {
	f, l, err := db.getDataRange()
	if err != nil {
		log.Println("handlerDataRange:", err.Error())
	}

	JSON{Status: http.StatusOK, Data: map[string]int64{"first": f.Unix(), "last": l.Unix()}}.write(w)
}
