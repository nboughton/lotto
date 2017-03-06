package main

import "fmt"

var formatYYYYMMDD = "2006-01-02"

type graphData struct {
	Labels   []string  `json:"labels"`
	Datasets []dataset `json:"datasets"`
}

type dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}

func parseResultsForGraph(records <-chan dbRow) graphData {
	var gd graphData
	gd.Datasets = make([]dataset, 7)

	labelsSet := false
	for row := range records {
		gd.Labels = append(gd.Labels, row.Date.Format(formatYYYYMMDD))
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
