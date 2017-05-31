package main

import (
	"fmt"
	"strconv"
)

var (
	colors = []string{
		"rgba(31,119,180,1)",
		"rgba(255,127,14,1)",
		"rgba(44,160,44,1)",
		"rgba(214,39,40,1)",
		"rgba(148,103,189,1)",
		"rgba(140,86,75,1)",
		"rgba(227,119,194,1)",
	}
)

const (
	maxBallNum       = 59
	balls            = 7
	graphTypeScatter = "scatter"
	graphTypeBar     = "bar"
	graphTypeLine    = "line"
)

type graphDataset struct {
	Label string   `json:"label"`
	Type  string   `json:"type"`
	Data  []string `json:"data"`
}

type graphData struct {
	Labels   []string       `json:"labels"`
	Datasets []graphDataset `json:"datasets"`
}

func graphTimeSeries(records <-chan dbRow) graphData {
	var d graphData
	d.Datasets = make([]graphDataset, balls)

	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if d.Datasets[ball].Label == "" {
				d.Datasets[ball].Label = label(ball)
				d.Datasets[ball].Type = graphTypeLine
			}

			d.Datasets[ball].Data = append(d.Datasets[ball].Data, strconv.Itoa(row.Num[ball]))
		}
		d.Labels = append(d.Labels, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatYYYYMMDD), row.Machine, row.Set))
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
				d.Datasets[ball].Data = make([]string, maxBallNum)
			}

			n, _ := strconv.Atoi(d.Datasets[ball].Data[row.Num[ball]-1])
			d.Datasets[ball].Data[row.Num[ball]-1] = strconv.Itoa(n + 1)
		}
	}

	return d
}

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
