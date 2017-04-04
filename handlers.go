package main

import (
	"log"
	"sort"
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

// NumbersData contains tidbits of information regarding numbers
type NumbersData struct {
	Frequent []int    `json:"frequent"`
	Ranges   []string `json:"ranges"`
	MeanAvg  []int    `json:"meanAvg"`
}

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

	bSort, bbSort, resFreq := make(ballSortByFreq, maxBallNum+1), make(ballSortByFreq, maxBallNum+1), []int{}
	for row := range db.getResults(p) {
		for ball := 0; ball < balls; ball++ {
			n := row.Num[ball]
			if ball < 6 {
				bSort[n].num = n
				bSort[n].freq++
			} else {
				bbSort[n].num = n
				bbSort[n].freq++
			}
		}

	}
	sort.Sort(sort.Reverse(bSort))
	sort.Sort(bbSort)

	for _, b := range bSort[:6] {
		resFreq = append(resFreq, b.num)
	}
	sort.Ints(resFreq)

	resFreq = append(resFreq, bbSort[len(bbSort)-1].num)

	w.WriteJSON(NumbersData{MeanAvg: resAvg, Ranges: resRange, Frequent: resFreq})

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
