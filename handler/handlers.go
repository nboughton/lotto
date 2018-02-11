package handler

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	jweb "github.com/nboughton/go-utils/json/web"
	"github.com/nboughton/lotto/graph"
	"github.com/nboughton/stalotto/db"
	"github.com/nboughton/stalotto/lotto"
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
				TimeSeries: graph.TimeSeries(e.DB.Results(p.Start, p.End, p.Machines, p.Sets)),
				FreqDist:   graph.FreqDist(e.DB.Results(p.Start, p.End, p.Machines, p.Sets)),
			},
		).Write(w)
	})
}

// ListSets returns a list of available ball sets
func ListSets(e *Env) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := params(r)

		res, err := e.DB.Sets(p.Start, p.End, p.Sets)
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
		p := params(r)
		res, err := e.DB.Machines(p.Start, p.End, p.Machines)
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

type queryParams struct {
	Start    time.Time
	End      time.Time
	Sets     []int
	Machines []string
}

func params(r *http.Request) queryParams {
	p := r.URL.Query()

	set, _ := strconv.Atoi(p["set"][0])
	sets := []int{}
	if set != 0 {
		sets = []int{set}
	}

	machine := p["machine"][0]
	machines := []string{}
	if machine != "all" {
		machines = []string{machine}
	}

	start, _ := time.Parse(time.RFC3339, p["start"][0])
	end, _ := time.Parse(time.RFC3339, p["end"][0])

	return queryParams{
		Start:    start,
		End:      end,
		Sets:     sets,
		Machines: machines,
	}
}

func createMainTableData(e *Env, p queryParams) []TableRow {
	set := lotto.ResultSet{}
	for res := range e.DB.Results(p.Start, p.End, p.Machines, p.Sets) {
		set = append(set, res)
	}
	balls, bonus := set.ByDrawFrequency()

	most := balls.Desc()[:6]
	sort.Ints(most)
	most = append(most, bonus.Desc()[0])

	least := balls.Asc()[:6]
	sort.Ints(least)
	least = append(least, bonus.Asc()[0])

	last := set[len(set)-1].Balls
	sort.Ints(last)
	last = append(last, set[len(set)-1].Bonus)

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
