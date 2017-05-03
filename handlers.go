package main

import (
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/gonum/stat"
	"github.com/pilu/traffic"
)

// PageData is used by the index template to populate things and stuff
type PageData struct {
	Machines   []string
	Sets       []int
	Start, End string
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

func handlerResultsScatter3D(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphResultsRawScatter3D(db.getResults(params(r))))
}

func handlerResultsFreqDist(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphResultsFreqDist(db.getResults(params(r)), true, r.Param("type")))
}

func handlerMSFreqDist(w traffic.ResponseWriter, r *traffic.Request) {
	switch r.Param("type") {
	case "bubble":
		w.WriteJSON(graphMSFreqDistBubble(db.getMachineSetCombinations(params(r))))
	case "scatter3d":
		w.WriteJSON(graphMSFreqDistScatter3D(db.getMachineSetCombinations(params(r))))
	}
}

func handlerResultsTimeSeries(w traffic.ResponseWriter, r *traffic.Request) {
	w.WriteJSON(graphResultsTimeSeries(db.getResults(params(r)), true, r.Param("type")))
}

func handlerNumbers(w traffic.ResponseWriter, r *traffic.Request) {
	var (
		p                     = params(r)
		sBalls, sBonus, modes = getSortedResultsByFreq(p)
		most, least           = getMostAndleast(sBalls, sBonus)
	)

	// Average results by totals
	resAvg, err := db.getResultsAverage(p)
	if err != nil {
		w.WriteJSON(err.Error())
		return
	}

	// Average range of each sorted ball field
	resRange, err := db.getResultsAverageRanges(p)
	if err != nil {
		w.WriteJSON(err.Error())
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

	w.WriteJSON(NumbersData{
		MeanAvg:  resAvg,
		ModeAvg:  m,
		Ranges:   resRange,
		Frequent: most,
		Least:    least,
		Random:   drawRandomSet(),
		Last:     last,
	})

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

func handlerListSets(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getSetList(params(r))
	if err != nil {
		w.WriteJSON(err)
	} else {
		w.WriteJSON(res)
	}
}

func handlerListMachines(w traffic.ResponseWriter, r *traffic.Request) {
	res, err := db.getMachineList(params(r))
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

func params(r *traffic.Request) queryParams {
	set, _ := strconv.Atoi(r.Param("set"))
	return queryParams{
		Start:   r.Param("start"),
		End:     r.Param("end"),
		Set:     set,
		Machine: r.Param("machine"),
	}
}
