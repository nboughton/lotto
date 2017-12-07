package handler

import (
	"log"
	"net/http"

	jweb "github.com/nboughton/go-utils/json/web"
	"github.com/nboughton/lotto/db"
	"github.com/nboughton/lotto/graph"
)

// Env allows for persistent data to be passed into route handlers, such as DB handles etc
type Env struct {
	DB *db.AppDB
}

// PageData is used by the index template to populate things and stuff
type PageData struct {
	MainTable  []TableRow `json:"mainTable"`
	TimeSeries graph.Data `json:"timeSeries"`
	FreqDist   graph.Data `json:"freqDist"`
}

// TableRow contains data used in the top table
type TableRow struct {
	Label string `json:"label"`
	Num   []int  `json:"num"`
}

// Query handles the main page query and returns all relevant data for the page.
func Query(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := params(r)

		jweb.New(http.StatusOK,
			PageData{
				MainTable:  createMainTableData(e, p),
				TimeSeries: graph.TimeSeries(e.DB.Results(p)),
				FreqDist:   graph.FreqDist(e.DB.Results(p)),
			},
		).Write(w)
	})
}

// ListSets returns a list of available ball sets
func ListSets(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := e.DB.Sets(params(r))
		if err != nil {
			jweb.New(http.StatusInternalServerError, err).Write(w)
		} else {
			jweb.New(http.StatusOK, res).Write(w)
		}
	})
}

// ListMachines returns a list of available lotto machines
func ListMachines(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := e.DB.Machines(params(r))
		if err != nil {
			jweb.New(http.StatusInternalServerError, err).Write(w)
		} else {
			jweb.New(http.StatusOK, res).Write(w)
		}
	})
}

// DataRange returns the first and last record dates
func DataRange(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, l, err := e.DB.DataRange()
		if err != nil {
			log.Println("handlerDataRange:", err.Error())
		}

		jweb.New(http.StatusOK, map[string]int64{"first": f.Unix(), "last": l.Unix()}).Write(w)
	})
}
