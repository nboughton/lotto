package graph

import (
	"fmt"
	"strconv"

	plotly "github.com/nboughton/go-plotlytypes"
	"github.com/nboughton/stalotto/lotto"
)

const (
	tScatter = "scatter"
	tBar     = "bar"
	tLine    = "line"
)

// Data groups datasets together to exported and turned into a nice graph
type Data struct {
	Datasets []plotly.Dataset `json:"datasets"`
}

var formatTSLabel = "06/01/02"

// TimeSeries creates a Data struct for a time series graph
func TimeSeries(set lotto.ResultSet) Data {
	var d Data
	d.Datasets = make([]plotly.Dataset, lotto.BALLS+1)

	for _, row := range set {
		for ball := 0; ball < lotto.BALLS; ball++ {
			if d.Datasets[ball].Name == "" {
				d.Datasets[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				d.Datasets[ball].Type = tLine
			}

			d.Datasets[ball].Y = append(d.Datasets[ball].Y, strconv.Itoa(row.Balls[ball]))
			d.Datasets[ball].X = append(d.Datasets[ball].X, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
		}

		if d.Datasets[lotto.BALLS].Name == "" {
			d.Datasets[lotto.BALLS].Name = "Bonus"
		}

		d.Datasets[lotto.BALLS].Y = append(d.Datasets[lotto.BALLS].Y, strconv.Itoa(row.Bonus))
		d.Datasets[lotto.BALLS].X = append(d.Datasets[lotto.BALLS].X, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
	}

	return d
}

// FreqDist creates a Data struct for a frequency distribution graph
func FreqDist(set lotto.ResultSet) Data {
	var d Data
	d.Datasets = make([]plotly.Dataset, lotto.BALLS+1)
	// Populate Labels
	var labels []string
	for i := 0; i < lotto.MAXBALLVAL; i++ {
		labels = append(labels, fmt.Sprintf("%d", i+1))
	}

	for _, row := range set {
		for ball := 0; ball < lotto.BALLS; ball++ {
			if d.Datasets[ball].Name == "" { // Set Name and create Data
				d.Datasets[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				d.Datasets[ball].Type = tBar
				d.Datasets[ball].Y = make([]string, lotto.MAXBALLVAL)
				d.Datasets[ball].X = labels
			}

			d.Datasets[ball].Y[row.Balls[ball]-1] = addY(d.Datasets[ball].Y[row.Balls[ball]-1], 1)
		}

		if d.Datasets[lotto.BALLS].Name == "" {
			d.Datasets[lotto.BALLS].Name = "Bonus"
			d.Datasets[lotto.BALLS].Type = tBar
			d.Datasets[lotto.BALLS].Y = make([]string, lotto.MAXBALLVAL)
		}

		d.Datasets[lotto.BALLS].Y[row.Bonus-1] = addY(d.Datasets[lotto.BALLS].Y[row.Bonus-1], 1)
	}

	return d
}

// MachineSetDist creates a dataset appropriate for a machine/set distribution bubble graph
/*
func MachineSetDist(records <-chan lotto.Result) Data {
	var (
		d Data
		m = make(map[string]int)
	)

	for row := range records {
		v := fmt.Sprintf("%s:%d", row.Machine, row.Set)
		m[v]++
	}

	return d
}
*/
func addY(cur string, inc int) string {
	a, err := strconv.Atoi(cur)
	if err != nil {
		return cur // effectively ignore errors for now
	}

	return strconv.Itoa(a + inc)
}

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
