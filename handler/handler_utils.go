package handler

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/nboughton/lotto/db"
)

type numFreq struct {
	num  int
	freq int
}

type sortByFreq []numFreq

func (b sortByFreq) Len() int           { return len(b) }
func (b sortByFreq) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b sortByFreq) Less(i, j int) bool { return b[i].freq < b[j].freq }

func params(r *http.Request) db.QueryParams {
	var (
		p        = r.URL.Query()
		set, _   = strconv.Atoi(p["set"][0])
		start, _ = time.Parse(time.RFC3339, p["start"][0])
		end, _   = time.Parse(time.RFC3339, p["end"][0])
	)
	return db.QueryParams{
		Start:   start,
		End:     end,
		Set:     set,
		Machine: p["machine"][0],
	}
}

func createMainTableData(e *Env, p db.QueryParams) []TableRow {
	var (
		sBalls, sBonus, byPosition                = getSortedResultsByFreq(e, p)
		mostOverall, leastOverall, mostByPosition = getMostAndleast(sBalls, sBonus, byPosition)
	)

	last, err := e.DB.LastDraw()
	if err != nil {
		log.Println(err)
	}

	return []TableRow{
		TableRow{Label: "Most Recent", Num: last},
		TableRow{Label: "Most Frequent (overall)", Num: mostOverall},
		TableRow{Label: "Most Frequent (position)", Num: mostByPosition},
		TableRow{Label: "Least Frequent (overall)", Num: leastOverall},
		TableRow{Label: "Random Set", Num: drawRandomSet()},
	}
}

func getMostAndleast(sBalls, sBonus sortByFreq, byPosition []sortByFreq) (mostOverall, leastOverall, mostByPosition []int) {
	// mostOverall = first 6 and leastOverall = last non-zero 6
	for i, j := 0, len(sBalls)-1; i < j; i, j = i+1, j-1 {
		if sBalls[i].num != 0 && len(mostOverall) < 6 {
			mostOverall = append(mostOverall, sBalls[i].num)
		}
		if sBalls[j].num != 0 && len(leastOverall) < 6 {
			leastOverall = append(leastOverall, sBalls[j].num)
		}
		if len(mostOverall) == 6 && len(leastOverall) == 6 {
			break
		}
	}

	// Sort the results, this is largely cosmetic.
	sort.Ints(mostOverall)
	sort.Ints(leastOverall)

	// Add the bonus ball for mostOverall/leastOverall frequent, don't duplicate numbers
	for i, j := len(sBonus)-1, 0; i > j; i, j = i-1, j+1 {
		if sBonus[i].num != 0 && len(mostOverall) < db.BALLS && !containsInt(mostOverall, sBonus[i].num) {
			mostOverall = append(mostOverall, sBonus[i].num)
		}
		if sBonus[j].num != 0 && len(leastOverall) < db.BALLS && !containsInt(leastOverall, sBonus[j].num) {
			leastOverall = append(leastOverall, sBonus[j].num)
		}
		if len(mostOverall) == db.BALLS && len(leastOverall) == db.BALLS {
			break
		}
	}

	// get most and least byPosition set
	mostByPosition = make([]int, db.BALLS)

	for i := range byPosition {
		mostByPosition[i] = byPosition[i][0].num
	}

	mostByPosition[6] = mostOverall[6]

	return mostOverall, leastOverall, mostByPosition
}

// Returns sorted ball results for query p by frequency as well as collated numbers for mode checking
func getSortedResultsByFreq(e *Env, p db.QueryParams) (sortByFreq, sortByFreq, []sortByFreq) {
	var (
		sBalls     = make(sortByFreq, db.MAXBALLNUM+1)
		sBonus     = make(sortByFreq, db.MAXBALLNUM+1)
		byPosition = make([]sortByFreq, db.BALLS)
	)

	// Set up byPosition arrays of numbers
	for i := range byPosition {
		byPosition[i] = make(sortByFreq, db.MAXBALLNUM+1)
	}

	for row := range e.DB.Results(p) {
		for ball := 0; ball < db.BALLS; ball++ {
			n := row.Num[ball]
			if ball < 6 {
				// Collate total frequncies for first 6
				sBalls[n].num = n
				sBalls[n].freq++

				// Collate sorted results by position
				byPosition[ball][n].num = n
				byPosition[ball][n].freq++
			} else {
				// Collate frequencies for bonus ball separately
				sBonus[n].num = n
				sBonus[n].freq++
			}
		}
	}

	// Sort byPosition arrays
	for i := range byPosition {
		sort.Sort(sort.Reverse(byPosition[i]))
	}

	// Sort both lists
	sort.Sort(sort.Reverse(sBalls))
	sort.Sort(sBonus)
	return sBalls, sBonus, byPosition
}

func containsInt(a []int, t int) bool {
	for _, n := range a {
		if n == t {
			return true
		}
	}

	return false
}
