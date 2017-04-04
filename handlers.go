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

func handlerNumbers(w traffic.ResponseWriter, r *traffic.Request) {
	switch r.Param("type") {
	case "average":
		res, err := db.getResultsAverage(params(r))
		if err != nil {
			w.WriteJSON(err.Error())
		} else {
			w.WriteJSON(res)
		}

	case "ranges":
		res, err := db.getResultsAverageRanges(params(r))
		if err != nil {
			w.WriteJSON(err.Error())
		} else {
			w.WriteJSON(res)
		}

	case "frequent":
		bSort, res := make(ballSortByFreq, maxBallNum+1), make([]int, balls)

		for row := range db.getResults(params(r)) {
			for ball := 0; ball < balls; ball++ {
				n := row.Num[ball]
				bSort[n].num = n
				bSort[n].freq++
			}
		}
		sort.Sort(sort.Reverse(bSort))

		for i, b := range bSort[:7] {
			res[i] = b.num
		}
		sort.Ints(res)

		w.WriteJSON(res)

	} // END SWITCH

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
