package main

import "fmt"

type plotlyDatasetLine struct {
	X    []string  `json:"x"`
	Y    []float64 `json:"y"`
	Name string    `json:"name"`
	Mode string    `json:"mode"`
}

func parseResultsForGraph(records <-chan dbRow) []plotlyDatasetLine {
	data := make([]plotlyDatasetLine, 7)

	i := 0
	for row := range records {
		for j := 0; j < 7; j++ {
			if i == 0 {
				data[j].Mode = "line"
				if j < 6 {
					data[j].Name = fmt.Sprintf("Ball %d", j+1)
				} else {
					data[j].Name = "Bonus Ball"
				}
			}
			data[j].X = append(data[j].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			data[j].Y = append(data[j].Y, float64(row.Num[j]))
		}
		i++
	}

	return data
}
