package main

import (
	"log"
	"net/http"
	//"sort"
	"strconv"
	"time"
	//"github.com/gonum/stat"
	//"github.com/gorilla/mux"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	Machines   []string
	Sets       []int
	Start, End time.Time
}

// NumbersData contains tidbits of information regarding numbers
type NumbersData struct {
	Frequent []int     `json:"frequent"`
	Least    []int     `json:"least"`
	Ranges   []string  `json:"ranges"`
	MeanAvg  []int     `json:"meanAvg"`
	ModeAvg  []float64 `json:"modeAvg"`
	Random   []int     `json:"random"`
	Last     []int     `json:"last"`
}

type numFreq struct {
	num  int
	freq int
}

type ballSortByFreq []numFreq

func (b ballSortByFreq) Len() int           { return len(b) }
func (b ballSortByFreq) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ballSortByFreq) Less(i, j int) bool { return b[i].freq < b[j].freq }

func handlerQuery(w http.ResponseWriter, r *http.Request) {
	p := params(r)
	switch p.Query {
	case 1: // Create a LineChart using the old numbers code
		JSON{Status: http.StatusOK, Data: lineGraph(db.getResults(p))}.write(w)
	}
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

func params(r *http.Request) queryParams {
	var (
		p        = r.URL.Query()
		query, _ = strconv.Atoi(p["query"][0])
		set, _   = strconv.Atoi(p["set"][0])
		start, _ = time.Parse(time.RFC3339, p["start"][0])
		end, _   = time.Parse(time.RFC3339, p["end"][0])
	)
	return queryParams{
		Query:   query,
		Start:   start,
		End:     end,
		Set:     set,
		Machine: p["machine"][0],
	}
}

/*
func handlerRoot(w http.ResponseWriter, r *http.Request) {
	s, e, err := db.getDataRange()
	if err != nil {
		log.Println(err)
	}

	s, _ = time.Parse(formatYYYYMMDD, "2015-10-10") // Because reasons.
	q := queryParams{
		Start:   s,
		End:     e,
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

	JSON{Status: http.StatusOK, Data: PageData{machines, sets, q.Start, q.End}}.write(w)
}

func handlerResultsScatter3D(w http.ResponseWriter, r *http.Request) {
	JSON{Status: http.StatusOK, Data: graphResultsRawScatter3D(db.getResults(params(r)))}.write(w)
}

func handlerResultsFreqDist(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	JSON{Status: http.StatusOK, Data: graphResultsFreqDist(db.getResults(params(r)), true, p["type"])}.write(w)
}

func handlerMSFreqDist(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	switch p["type"] {
	case "bubble":
		JSON{Status: http.StatusOK, Data: graphMSFreqDistBubble(db.getMachineSetCombinations(params(r)))}.write(w)
	case "scatter3d":
		JSON{Status: http.StatusOK, Data: graphMSFreqDistScatter3D(db.getMachineSetCombinations(params(r)))}.write(w)
	}
}

func handlerResultsTimeSeries(w http.ResponseWriter, r *http.Request) {
	p := mux.Vars(r)
	JSON{Status: http.StatusOK, Data: graphResultsTimeSeries(db.getResults(params(r)), true, p["type"])}.write(w)
}

func handlerNumbers(w http.ResponseWriter, r *http.Request) {
	var (
		p                     = params(r)
		sBalls, sBonus, modes = getSortedResultsByFreq(p)
		most, least           = getMostAndleast(sBalls, sBonus)
	)

	// Average results by totals
	resAvg, err := db.getResultsAverage(p)
	if err != nil {
		JSON{Status: http.StatusInternalServerError, Data: err.Error()}.write(w)
		return
	}

	// Average range of each sorted ball field
	resRange, err := db.getResultsAverageRanges(p)
	if err != nil {
		JSON{Status: http.StatusInternalServerError, Data: err.Error()}.write(w)
		return
	}

	// Create Mode sets
	m := make([]float64, balls)
	for i, set := range modes {
		m[i], _ = stat.Mode(set, nil)
	}

	last, err := db.getLastDraw()
	if err != nil {
		log.Println(err)
	}

	JSON{
		Status: http.StatusOK,
		Data: NumbersData{
			MeanAvg:  resAvg,
			ModeAvg:  m,
			Ranges:   resRange,
			Frequent: most,
			Least:    least,
			Random:   drawRandomSet(),
			Last:     last,
		},
	}.write(w)

}

func getMostAndleast(sBalls, sBonus ballSortByFreq) (most, least []int) {
	// most = first 6 and least = last non-zero 6
	for i, j := 0, len(sBalls)-1; i < j; i, j = i+1, j-1 {
		if sBalls[i].num != 0 && len(most) < 6 {
			most = append(most, sBalls[i].num)
		}
		if sBalls[j].num != 0 && len(least) < 6 {
			least = append(least, sBalls[j].num)
		}
		if len(most) == 6 && len(least) == 6 {
			break
		}
	}

	// Sort the results, this is largely cosmetic.
	sort.Ints(most)
	sort.Ints(least)

	// Add the bonus ball for most/least frequent, don't duplicate numbers
	for i, j := len(sBonus)-1, 0; i > j; i, j = i-1, j+1 {
		if sBonus[i].num != 0 && len(most) < balls && !containsInt(most, sBonus[i].num) {
			most = append(most, sBonus[i].num)
		}
		if sBonus[j].num != 0 && len(least) < balls && !containsInt(least, sBonus[j].num) {
			least = append(least, sBonus[j].num)
		}
		if len(most) == balls && len(least) == balls {
			break
		}
	}
	return most, least
}

// Returns sorted ball results for query p by frequency as well as collated numbers for mode checking
func getSortedResultsByFreq(p queryParams) (ballSortByFreq, ballSortByFreq, [][]float64) {
	var (
		sBalls = make(ballSortByFreq, maxBallNum+1)
		sBonus = make(ballSortByFreq, maxBallNum+1)
		modes  = make([][]float64, balls)
	)
	for row := range db.getResults(p) {
		for ball := 0; ball < balls; ball++ {
			n := row.Num[ball]
			if ball < 6 {
				// Collate total frequncies for first 6
				sBalls[n].num = n
				sBalls[n].freq++
			} else {
				// Collate frequencies for bonus ball separately
				sBonus[n].num = n
				sBonus[n].freq++
			}
			// Collate raw numbers for mode
			modes[ball] = append(modes[ball], float64(n))
		}
	}

	// Sort both lists
	sort.Sort(sort.Reverse(sBalls))
	sort.Sort(sBonus)
	return sBalls, sBonus, modes
}

func containsInt(a []int, t int) bool {
	for _, n := range a {
		if n == t {
			return true
		}
	}

	return false
}
*/
