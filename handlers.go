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

func handlerNumbers(w traffic.ResponseWriter, r *traffic.Request) {
	p := params(r)
	resAvg, err := db.getResultsAverage(p)
	if err != nil {
		w.WriteJSON(err.Error())
		return
	}

	resRange, err := db.getResultsAverageRanges(p)
	if err != nil {
		w.WriteJSON(err.Error())
		return
	}

	// Treat the bonus ball as a separate entity as it is selected in isolation from
	// the first six. Hence bSort and bbSort.
	var (
		bSort               = make(ballSortByFreq, maxBallNum+1)
		bbSort              = make(ballSortByFreq, maxBallNum+1)
		modes               = make([][]float64, balls)
		mostFreq, leastFreq []int
	)
	for row := range db.getResults(p) {
		for ball := 0; ball < balls; ball++ {
			n := row.Num[ball]
			if ball < 6 {
				// Collate total frequncies for first 6
				bSort[n].num = n
				bSort[n].freq++
			} else {
				// Collate frequencies for bonus ball separately
				bbSort[n].num = n
				bbSort[n].freq++
			}
			// Collate raw numbers for mode
			modes[ball] = append(modes[ball], float64(n))
		}
		//last = row.Num
	}

	// Sort both lists
	sort.Sort(sort.Reverse(bSort))
	sort.Sort(bbSort)

	// Pick out most frequent first 6
	for _, b := range bSort[:6] {
		mostFreq = append(mostFreq, b.num)
	}

	// Pick out least frequent last six, ignoring any 0s
	for i := len(bSort) - 1; i > 0; i-- {
		if len(leastFreq) == 6 {
			break
		}
		if bSort[i].num != 0 {
			leastFreq = append(leastFreq, bSort[i].num)
		}
	}

	// Sort the results, this is largely cosmetic.
	sort.Ints(mostFreq)
	sort.Ints(leastFreq)

	// Add the bonus ball for most frequent, don't duplicate numbers
	for i := len(bbSort) - 1; i > 0; i-- {
		if bbSort[i].num != 0 && !containsInt(mostFreq, bbSort[i].num) {
			mostFreq = append(mostFreq, bbSort[i].num)
			break
		}
	}

	// Add the bonus ball for least frequent, ensuring no duplicate numbers
	for _, b := range bbSort {
		if b.num != 0 && !containsInt(leastFreq, b.num) {
			leastFreq = append(leastFreq, b.num)
			break
		}
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
		Frequent: mostFreq,
		Least:    leastFreq,
		Random:   drawRandomSet(),
		Last:     last,
	})

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
