package main

import "fmt"

type plotlyDatasetLine struct {
	X    []string  `json:"x"`
	Y    []float64 `json:"y"`
	Name string    `json:"name"`
	Mode string    `json:"mode"`
}

func parseResultsForLineGraph(records <-chan dbRow) []plotlyDatasetLine {
	data := make([]plotlyDatasetLine, 7)

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				data[ball].Mode = "line"
				if ball < 6 {
					data[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				} else {
					data[ball].Name = "Bonus Ball"
				}
			}
			data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			data[ball].Y = append(data[ball].Y, float64(row.Num[ball]))
		}
		i++
	}

	return data
}

type plotlyDataset3DScatter struct {
	X      []string `json:"x"` // machine:set
	Y      []string `json:"y"` // date
	Z      []int    `json:"z"` // number
	Name   string   `json:"name"`
	Mode   string   `json:"mode"`
	Type   string   `json:"type"`
	Marker marker   `json:"marker"`
}

type marker struct {
	//Colour string `json:"colour"`
	Size    int     `json:"size"`
	Line    line    `json:"line"`
	Opacity float64 `json:"opacity"`
	Symbol  string  `json:"symbol"`
}

type line struct {
	Width float64 `json:"width"`
}

func parseResultsFor3DScatterGraph(records <-chan dbRow) []plotlyDataset3DScatter {
	data := make([]plotlyDataset3DScatter, 7)

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				set := plotlyDataset3DScatter{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    3,
						Line:    line{Width: 0.5},
						Opacity: 0.8,
						Symbol:  "circle",
					},
				}

				data[ball] = set
				if ball < 6 {
					data[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				} else {
					data[ball].Name = "Bonus Ball"
				}
			}
			data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d", row.Machine, row.Set))
			data[ball].Y = append(data[ball].Y, row.Date.Format(formatYYYYMMDD))
			data[ball].Z = append(data[ball].Z, row.Num[ball])
		}
		i++
	}

	return data
}

/*
func parseResultsFor3DScatterGraph(records <-chan dbRow) []plotlyDataset3DScatter {
	data := make([]plotlyDataset3DScatter, 7)

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				set := plotlyDataset3DScatter{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    12,
						Line:    line{Width: 0.5},
						Opacity: 0.7,
						Symbol:  "circle",
					},
				}

				data[ball] = set
				if ball < 6 {
					data[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				} else {
					data[ball].Name = "Bonus Ball"
				}
			}
			data[ball].X = append(data[ball].X, row.Set)
			data[ball].Y = append(data[ball].Y, row.Num[ball])
			data[ball].Z = append(data[ball].Z, row.Machine)
		}
		i++
	}

	return data
}
*/
