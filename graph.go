package main

import "fmt"

type plotlyData struct {
	Data []plotlyDatasetLine `json:"data"`
}

type plotlyDatasetLine struct {
	X    []string  `json:"x"`
	Y    []float64 `json:"y"`
	Name string    `json:"name"`
	Mode string    `json:"mode"`
}

func parseResultsForGraph(records <-chan dbRow) plotlyData {
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
			gd.Data[j].X = append(gd.Data[j].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			gd.Data[j].Y = append(gd.Data[j].Y, float64(row.Num[j]))
		}
		i++
	}

	return gd
}
