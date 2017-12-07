package graph

import (
	"fmt"

	"github.com/nboughton/lotto/db"
)

const (
	typeScatter = "scatter"
	typeBar     = "bar"
	typeLine    = "line"
)

// Dataset wraps an individual dataset that can be exported to json
type Dataset struct {
	Label string    `json:"label"`
	Type  string    `json:"type"`
	Data  []float64 `json:"data"`
}

// Data groups datasets together to exported and turned into a nice graph
type Data struct {
	Labels   []string  `json:"labels"`
	Datasets []Dataset `json:"datasets"`
}

var formatTSLabel = "06/01/02"

// TimeSeries creates a Data struct for a time series graph
func TimeSeries(records <-chan db.Record) Data {
	var d Data
	d.Datasets = make([]Dataset, db.BALLS)

	for row := range records {
		for ball := 0; ball < db.BALLS; ball++ {
			if d.Datasets[ball].Label == "" {
				d.Datasets[ball].Label = label(ball)
				d.Datasets[ball].Type = typeLine
			}

			d.Datasets[ball].Data = append(d.Datasets[ball].Data, float64(row.Num[ball]))
		}
		d.Labels = append(d.Labels, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
	}

	return d
}

// FreqDist creates a Data struct for a frequency distribution graph
func FreqDist(records <-chan db.Record) Data {
	var d Data
	d.Datasets = make([]Dataset, db.BALLS)
	// Populate Labels
	for i := 0; i < db.MAXBALLNUM; i++ {
		d.Labels = append(d.Labels, fmt.Sprintf("%d", i+1))
	}

	for row := range records {
		for ball := 0; ball < db.BALLS; ball++ {
			if d.Datasets[ball].Label == "" { // Set Label and create Data
				d.Datasets[ball].Label = label(ball)
				d.Datasets[ball].Type = typeBar
				d.Datasets[ball].Data = make([]float64, db.MAXBALLNUM)
			}

			d.Datasets[ball].Data[row.Num[ball]-1]++
		}
	}

	return d
}

// MachineSetDist creates a dataset appropriate for a machine/set distribution bubble graph
func MachineSetDist(records <-chan db.Record) Data {
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

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
