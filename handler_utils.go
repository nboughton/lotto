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
		sBalls, sBonus = getSortedResultsByFreq(p)
		most, least    = getMostAndleast(sBalls, sBonus)
	)

	last, err := db.getLastDraw()
	if err != nil {
		log.Println(err)
	}

	return []TableRow{
		TableRow{
			Label: "Last Draw",
			Num:   last,
		},
		TableRow{
			Label: "Most Frequent",
			Num:   most,
		},
		TableRow{
			Label: "Least Frequent",
			Num:   least,
		},
		TableRow{
			Label: "Random Set",
			Num:   drawRandomSet(),
		},
	}
}

func getMostAndleast(sBalls, sBonus sortByFreq) (most, least []int) {
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
func getSortedResultsByFreq(p queryParams) (sortByFreq, sortByFreq) {
	var (
		sBalls = make(sortByFreq, maxBallNum+1)
		sBonus = make(sortByFreq, maxBallNum+1)
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
		}
	}

	// Sort both lists
	sort.Sort(sort.Reverse(sBalls))
	sort.Sort(sBonus)
	return sBalls, sBonus
}

func containsInt(a []int, t int) bool {
	for _, n := range a {
		if n == t {
			return true
		}
	}

	return false
}
