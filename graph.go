package main

import "fmt"

type datasetLine struct {
	X    []string  `json:"x"`
	Y    []float64 `json:"y"`
	Name string    `json:"name"`
	Mode string    `json:"mode"`
	Line line      `json:"line"`
}

type datasetScatter3D struct {
	X      []string `json:"x"` // machine:set
	Y      []string `json:"y"` // date
	Z      []int    `json:"z"` // number
	Name   string   `json:"name"`
	Mode   string   `json:"mode"`
	Type   string   `json:"type"`
	Marker marker   `json:"marker"`
}

type marker struct {
	Colour  string  `json:"colour"`
	Size    float64 `json:"size"`
	Line    line    `json:"line"`
	Opacity float64 `json:"opacity"`
	Symbol  string  `json:"symbol"`
}

type line struct {
	Width  float64 `json:"width"`
	Colour string  `json:"colour"`
	Shape  string  `json:"shape"`
}

func graphLine(records <-chan dbRow) []datasetLine {
	data := make([]datasetLine, 7)

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				data[ball] = datasetLine{
					Mode: "line",
					Line: line{Shape: "spline", Width: 1.5},
				}

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

func graphScatter3D(records <-chan dbRow) []datasetScatter3D {
	data := make([]datasetScatter3D, 7)

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				set := datasetScatter3D{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    4,
						Opacity: 1,
						//Line:    line{Width: 0.2},
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
