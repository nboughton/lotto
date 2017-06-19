package main

import (
	"fmt"
)

const (
	maxBallNum       = 59
	balls            = 7
	graphTypeScatter = "scatter"
	graphTypeBar     = "bar"
	graphTypeLine    = "line"
)

type graphDataset struct {
	Label string    `json:"label"`
	Type  string    `json:"type"`
	Data  []float64 `json:"data"`
}

type graphData struct {
	Labels   []string       `json:"labels"`
	Datasets []graphDataset `json:"datasets"`
}

var formatTSLabel = "06/01/02"

func graphTimeSeries(records <-chan dbRow) graphData {
	var d graphData
	d.Datasets = make([]graphDataset, balls)

	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if d.Datasets[ball].Label == "" {
				d.Datasets[ball].Label = label(ball)
				d.Datasets[ball].Type = graphTypeLine
			}

			d.Datasets[ball].Data = append(d.Datasets[ball].Data, float64(row.Num[ball]))
		}
		d.Labels = append(d.Labels, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatTSLabel), row.Machine[:3], row.Set))
	}

	return d
}

func graphFreqDist(records <-chan dbRow) graphData {
	var d graphData
	d.Datasets = make([]graphDataset, balls)
	// Populate Labels
	for i := 0; i < maxBallNum; i++ {
		d.Labels = append(d.Labels, fmt.Sprintf("%d", i+1))
	}

	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if d.Datasets[ball].Label == "" { // Set Label and create Data
				d.Datasets[ball].Label = label(ball)
				d.Datasets[ball].Type = graphTypeBar
				d.Datasets[ball].Data = make([]float64, maxBallNum)
			}

			d.Datasets[ball].Data[row.Num[ball]-1]++
		}
	}

	return d
}

func graphMachineSetDist(records <-chan dbRow) graphData {
	var (
		d graphData
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
