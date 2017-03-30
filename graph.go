package main

import (
	"fmt"

	"github.com/gonum/stat"
)

type dataset2D struct {
	X           []string  `json:"x"`
	Y           []float64 `json:"y"`
	Name        string    `json:"name"`
	Mode        string    `json:"mode"`
	Type        string    `json:"type"`
	Line        line      `json:"line"`
	Marker      marker    `json:"marker"`
	ConnectGaps bool      `json:"connectgaps"`
}

type dataset3D struct {
	X      []string `json:"x"` // machine:set
	Y      []string `json:"y"` // date
	Z      []int    `json:"z"` // number
	Name   string   `json:"name"`
	Mode   string   `json:"mode"`
	Type   string   `json:"type"`
	Marker marker   `json:"marker"`
}

type marker struct {
	Colour  string  `json:"color"`
	Size    float64 `json:"size"`
	Line    line    `json:"line"`
	Opacity float64 `json:"opacity"`
	Symbol  string  `json:"symbol"`
}

type line struct {
	Width  float64 `json:"width"`
	Colour string  `json:"color"`
	Shape  string  `json:"shape"`
	Dash   string  `json:"dash"`
}

func graphScatter(records <-chan dbRow, bestFit bool) []dataset2D {
	data := make([]dataset2D, 7)

	// x: date:machine:set
	// y: ball result

	// Use an array to sync markers with regression lines
	colors := []string{"rgba(31,119,180,1)", "rgba(255,127,14,1)", "rgba(44,160,44,1)", "rgba(214,39,40,1)", "rgba(148,103,189,1)", "rgba(140,86,75,1)", "rgba(227,119,194,1)"}

	// Create a float64 numeric representation of linear regression 'X' axis
	lrX := []float64{}

	// Create scatter data
	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				data[ball] = dataset2D{
					Mode: "markers",
					Marker: marker{
						Colour:  colors[ball],
						Size:    8,
						Opacity: 1,
					},
				}

				if ball < 6 {
					data[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				} else {
					data[ball].Name = "Bonus"
				}
			}
			data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			data[ball].Y = append(data[ball].Y, float64(row.Num[ball]))
		}
		lrX = append(lrX, float64(i))
		i++
	}

	// Calculate and append linear regressions for each set
	if bestFit {
		linReg := make([]dataset2D, 7)
		for i, set := range data {
			a, b := stat.LinearRegression(lrX, set.Y, nil, false)

			y := make([]float64, len(lrX))
			for idx, x := range lrX {
				y[idx] = a + (b * x)
			}

			linReg[i] = dataset2D{
				Name: set.Name,
				Mode: "lines",
				Line: line{
					Dash:   "dot",
					Width:  1.5,
					Colour: colors[i],
				},
				ConnectGaps: true,
				X:           set.X,
				Y:           y,
			}
		}

		data = append(data, linReg...)
	}

	return data
}

func graphScatter3D(records <-chan dbRow) []dataset3D {
	data := make([]dataset3D, 7)

	// x: machine:set
	// y: date
	// z: ball result

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				set := dataset3D{
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
					data[ball].Name = "Bonus"
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

func graphBar(records <-chan dbRow) []dataset2D {
	data := make([]dataset2D, 7)

	// x: numbers 1..60
	// y: frequency

	i := 0
	for row := range records {
		for ball := 0; ball < 7; ball++ {
			if i == 0 {
				// Populate X labels
				var x []string
				for j := 0; j < 60; j++ {
					x = append(x, fmt.Sprintf("%d", j+1))
				}

				data[ball] = dataset2D{
					Type: "bar",
					X:    x,
					Y:    make([]float64, 60),
				}

				if ball < 6 {
					data[ball].Name = fmt.Sprintf("Ball %d", ball+1)
				} else {
					data[ball].Name = "Bonus"
				}
			}

			data[ball].Y[row.Num[ball]]++
		}
		i++
	}

	return data
}
