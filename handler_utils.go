package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type numFreq struct {
	num  int
	freq int
}

type sortByFreq []numFreq

func (b sortByFreq) Len() int           { return len(b) }
func (b sortByFreq) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b sortByFreq) Less(i, j int) bool { return b[i].freq < b[j].freq }

func params(r *http.Request) queryParams {
	var (
		p        = r.URL.Query()
		set, _   = strconv.Atoi(p["set"][0])
		start, _ = time.Parse(time.RFC3339, p["start"][0])
		end, _   = time.Parse(time.RFC3339, p["end"][0])
	)
	return queryParams{
		Start:   start,
		End:     end,
		Set:     set,
		Machine: p["machine"][0],
	}
}

func createMainTableData(p queryParams) []TableRow {
	var (
		sBalls, sBonus, byPosition                = getSortedResultsByFreq(p)
		mostOverall, leastOverall, mostByPosition = getMostAndleast(sBalls, sBonus, byPosition)
	)

	last, err := db.getLastDraw()
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
		if sBonus[i].num != 0 && len(mostOverall) < balls && !containsInt(mostOverall, sBonus[i].num) {
			mostOverall = append(mostOverall, sBonus[i].num)
		}
		if sBonus[j].num != 0 && len(leastOverall) < balls && !containsInt(leastOverall, sBonus[j].num) {
			leastOverall = append(leastOverall, sBonus[j].num)
		}
		if len(mostOverall) == balls && len(leastOverall) == balls {
			break
		}
	}

	// get most and least byPosition set
	mostByPosition = make([]int, balls)

	for i := range byPosition {
		mostByPosition[i] = byPosition[i][0].num
	}

	mostByPosition[6] = mostOverall[6]

	return mostOverall, leastOverall, mostByPosition
}

// Returns sorted ball results for query p by frequency as well as collated numbers for mode checking
func getSortedResultsByFreq(p queryParams) (sortByFreq, sortByFreq, []sortByFreq) {
	var (
		sBalls     = make(sortByFreq, maxBallNum+1)
		sBonus     = make(sortByFreq, maxBallNum+1)
		byPosition = make([]sortByFreq, balls)
	)

	// Set up byPosition arrays of numbers
	for i := range byPosition {
		byPosition[i] = make(sortByFreq, maxBallNum+1)
	}

	for row := range db.getResults(p) {
		for ball := 0; ball < balls; ball++ {
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
