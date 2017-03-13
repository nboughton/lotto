package main

import "fmt"

type graphData struct {
	Labels   []string  `json:"labels"`
	Datasets []dataset `json:"datasets"`
}

type dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

type plotlyData struct {
	Data []plotlyDatasetLine `json:"data"`
}

type plotlyDatasetLine struct {
	X    []string  `json:"x"`
	Y    []float64 `json:"y"`
	Name string    `json:"name"`
	Mode string    `json:"mode"`
}

func parseResultsForGraph(records <-chan dbRow) graphData {
	var gd graphData
	gd.Datasets = make([]dataset, 7)

	labelsSet := false
	for row := range records {
		gd.Labels = append(gd.Labels, fmt.Sprintf("%d:%s:%s", row.Set, row.Machine, row.Date.Format(formatYYYYMMDD)))
		for i := 0; i < 7; i++ {
			if !labelsSet {
				label := ""
				if i < 6 {
					label = fmt.Sprintf("Ball %d", i+1)
				} else {
					label = "Bonus Ball"
				}
				gd.Datasets[i].Label = label
			}
			gd.Datasets[i].Data = append(gd.Datasets[i].Data, row.Num[i])
		}
		labelsSet = true
	}

	return gd
}

func parseResultsForGraphPlotly(records <-chan dbRow) plotlyData {
	var gd plotlyData
	gd.Data = make([]plotlyDatasetLine, 7)

	i := 0
	for row := range records {
		for j := 0; j < 7; j++ {
			if i == 0 {
				gd.Data[j].Mode = "line"
				if j < 6 {
					gd.Data[j].Name = fmt.Sprintf("Ball %d", j+1)
				} else {
					gd.Data[j].Name = "Bonus Ball"
				}
			}
			gd.Data[j].X = append(gd.Data[j].X, fmt.Sprintf("%d:%s:%s", row.Set, row.Machine, row.Date.Format(formatYYYYMMDD)))
			gd.Data[j].Y = append(gd.Data[j].Y, float64(row.Num[j]))
		}
		i++
	}

	return gd
}
