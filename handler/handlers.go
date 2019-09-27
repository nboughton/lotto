package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/nboughton/lotto/graph"
	"github.com/nboughton/stalotto/db"
	"github.com/nboughton/stalotto/lotto"
)

var jsonH = struct {
	key string
	val string
}{
	"Content-Type",
	"application/json; charset=utf-8",
}

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

type request struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Sets     []int     `json:"sets"`
	Machines []string  `json:"machines"`
}

// Query handles the main page query and returns all relevant data for the page.
func Query(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No request body", 400)
			return
		}

		var p request
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Malformed JSON request", 500)
			return
		}

		set := lotto.ResultSet{}
		for res := range e.DB.Results(p.Start, p.End, p.Machines, p.Sets, false) {
			set = append(set, res)
		}

		w.Header().Set(jsonH.key, jsonH.val)
		json.NewEncoder(w).Encode(PageData{
			MainTable:  createMainTableData(e, p),
			TimeSeries: graph.TimeSeries(set),
			FreqDist:   graph.FreqDist(set),
		})
	})
}

// ListSets returns a list of available ball sets
func ListSets(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No request body", 400)
			return
		}

		var p request
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Malformed JSON request", 500)
			return
		}

		res, err := e.DB.Sets(p.Start, p.End, p.Machines)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set(jsonH.key, jsonH.val)
		json.NewEncoder(w).Encode(res)
	})
}

// ListMachines returns a list of available lotto machines
func ListMachines(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "No request body", 400)
			return
		}

		var p request
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Malformed JSON request", 500)
			return
		}

		res, err := e.DB.Machines(p.Start, p.End, p.Sets)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set(jsonH.key, jsonH.val)
		json.NewEncoder(w).Encode(res)
	})
}

// DataRange returns the first and last record dates
func DataRange(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, l, err := e.DB.DataRange()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		res := struct {
			First int64
			Last  int64
		}{
			f.Unix(),
			l.Unix(),
		}

		w.Header().Set(jsonH.key, jsonH.val)
		json.NewEncoder(w).Encode(res)
	})
}

func createMainTableData(e *Env, p request) []TableRow {
	set := lotto.ResultSet{}
	for res := range e.DB.Results(p.Start, p.End, p.Machines, p.Sets) {
		set = append(set, res)
	}
	balls, bonus := set.ByDrawFrequency()

	most := balls.Prune().Desc().Balls()[:6]
	sort.Ints(most)
	most = append(most, bonus.Prune().Desc().Balls()[0])

	least := balls.Prune().Asc().Balls()[:6]
	sort.Ints(least)
	least = append(least, bonus.Prune().Asc().Balls()[0])

	last := set[0].Balls
	sort.Ints(last)
	last = append(last, set[0].Bonus)

	numbers := []int{}
	for i := 1; i <= lotto.MAXBALLVAL; i++ {
		numbers = append(numbers, i)
	}

	return []TableRow{
		TableRow{Label: "Most Recent", Num: last},
		TableRow{Label: "Most Frequent (overall)", Num: most},
		TableRow{Label: "Least Frequent (overall)", Num: least},
		TableRow{Label: "Random Set", Num: lotto.Draw(numbers, lotto.BALLS+1)},
	}
}
