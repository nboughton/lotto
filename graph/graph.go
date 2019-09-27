package graph

import (
	"fmt"

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
	Datasets []plotly.Dataset `json:"data"`
}

var (
	formatTSLabel = "06/01/02"
	line          = plotly.Line{
		Width: 1.5,
	}
	marker = plotly.Marker{
		Line: line,
	}
)

// TimeSeries creates a Data struct for a time series graph
func TimeSeries(set lotto.ResultSet) Data {
	var d Data
	d.Datasets = make([]plotly.Dataset, lotto.BALLS+1)

	for _, row := range set {
		for ball := 0; ball < lotto.BALLS; ball++ {
			if d.Datasets[ball].Name == "" {
				d.Datasets[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				d.Datasets[ball].Type = tLine
				d.Datasets[ball].Line = line
				d.Datasets[ball].ConnectGaps = true
			}

			d.Datasets[ball].Y = d.Datasets[ball].Y.AppendInt(row.Balls[ball])
			d.Datasets[ball].X = d.Datasets[ball].X.AppendStr(fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
		}

		if d.Datasets[lotto.BALLS].Name == "" {
			d.Datasets[lotto.BALLS].Name = "Bonus"
			d.Datasets[lotto.BALLS].Type = tLine
			d.Datasets[lotto.BALLS].Line = line
		}

		d.Datasets[lotto.BALLS].Y = d.Datasets[lotto.BALLS].Y.AppendInt(row.Bonus)
		d.Datasets[lotto.BALLS].X = d.Datasets[lotto.BALLS].X.AppendStr(fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
	}

	return d
}

// FreqDist creates a Data struct for a frequency distribution graph
func FreqDist(set lotto.ResultSet) Data {
	var d Data
	d.Datasets = make([]plotly.Dataset, lotto.BALLS+1)
	// Populate Labels
	var labels plotly.Axis
	for i := 0; i < lotto.MAXBALLVAL; i++ {
		labels = labels.AppendStr(fmt.Sprintf("%d", i+1))
	}

	for _, row := range set {
		for ball := 0; ball < lotto.BALLS; ball++ {
			if d.Datasets[ball].Name == "" { // Set Name and create Data
				d.Datasets[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				d.Datasets[ball].Type = tBar
				d.Datasets[ball].Y = make(plotly.Axis, lotto.MAXBALLVAL)
				d.Datasets[ball].X = labels
			}

			var err error
			if d.Datasets[ball].Y, err = d.Datasets[ball].Y.AddInt(row.Balls[ball]-1, 1); err != nil {
				fmt.Println(err)
			}
		}

		if d.Datasets[lotto.BALLS].Name == "" {
			d.Datasets[lotto.BALLS].Name = "Bonus"
			d.Datasets[lotto.BALLS].Type = tBar
			d.Datasets[lotto.BALLS].Y = make(plotly.Axis, lotto.MAXBALLVAL)
			d.Datasets[lotto.BALLS].X = labels
		}

		var err error
		if d.Datasets[lotto.BALLS].Y, err = d.Datasets[lotto.BALLS].Y.AddInt(row.Bonus-1, 1); err != nil {
			fmt.Println(err)
		}
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

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
