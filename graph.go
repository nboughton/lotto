package main

import (
	"fmt"

	"github.com/gonum/stat"
)

var (
	maxBallNum = 59
	balls      = 7
	colors     = []string{
		"rgba(31,119,180,1)",
		"rgba(255,127,14,1)",
		"rgba(44,160,44,1)",
		"rgba(214,39,40,1)",
		"rgba(148,103,189,1)",
		"rgba(140,86,75,1)",
		"rgba(227,119,194,1)",
	}
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
	X      []string `json:"x"` // machine
	Y      []int    `json:"y"` // set
	Z      []int    `json:"z"` // number
	Name   string   `json:"name"`
	Mode   string   `json:"mode"`
	Type   string   `json:"type"`
	Marker marker   `json:"marker"`
	Line   line     `json:"line"`
}

type marker struct {
	Colour  string  `json:"color"`
	Size    float64 `json:"size"`
	Line    line    `json:"line"`
	Opacity float64 `json:"opacity"`
	Symbol  string  `json:"symbol"`
}

type line struct {
	Width   float64 `json:"width"`
	Colour  string  `json:"color"`
	Shape   string  `json:"shape"`
	Dash    string  `json:"dash"`
	Opacity float64 `json:"opacity"`
}

func graphTimeSeries(records <-chan dbRow, bestFit bool, t string) []dataset2D {
	data := make([]dataset2D, balls)

	// Distribution over time
	switch t {
	case "scatter":
		i := 0
		for row := range records {
			for ball := 0; ball < balls; ball++ {
				if i == 0 {
					data[ball] = dataset2D{
						Mode: "markers+lines",
						Marker: marker{
							Colour:  colors[ball],
							Size:    8,
							Opacity: 1,
						},
						Line: line{
							Shape: "spline",
							Dash:  "dot",
							Width: 0.5,
						},
					}

					data[ball].Name = label(ball)
				}
				data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
				data[ball].Y = append(data[ball].Y, float64(row.Num[ball]))
			}
			i++
		}
	}

	return data
}

func graphFreqDist(records <-chan dbRow, bestFit bool, t string) []dataset2D {
	data := make([]dataset2D, balls)

	switch t {
	case "scatter":
		// Create scatter data
		i := 0
		for row := range records {
			for ball := 0; ball < balls; ball++ {
				if i == 0 {
					data[ball] = dataset2D{
						Mode: "markers+lines",
						Marker: marker{
							Colour:  colors[ball],
							Size:    8,
							Opacity: 1,
						},
						Line: line{
							Dash:  "dot",
							Width: 0.5,
						},
						X: freqDistXLabels(),
						Y: make([]float64, maxBallNum),
					}

					data[ball].Name = label(ball)
				}

				data[ball].Y[row.Num[ball]-1]++
			}
			i++
		}

	case "bar":
		i := 0
		for row := range records {
			for ball := 0; ball < balls; ball++ {
				if i == 0 {
					data[ball] = dataset2D{
						Type: "bar",
						X:    freqDistXLabels(),
						Y:    make([]float64, maxBallNum),
					}
					data[ball].Name = label(ball)
				}

				data[ball].Y[row.Num[ball]-1]++
			}
			i++
		}
	}

	return data
}

func graphScatter3D(records <-chan dbRow) []dataset3D {
	data := make([]dataset3D, balls)

	// x: machine
	// y: set
	// z: ball result

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {
				data[ball] = dataset3D{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    4,
						Opacity: 0.9,
						Line:    line{Width: 0.1},
					},
					/*
						Line: line{
							Width:   0.5,
							Opacity: 0.5,
						},
					*/
				}

				data[ball].Name = label(ball)
			}
			data[ball].X = append(data[ball].X, row.Machine)
			data[ball].Y = append(data[ball].Y, row.Set)
			data[ball].Z = append(data[ball].Z, row.Num[ball])
		}
		i++
	}

	return data
}

func generateLinRegSets(data []dataset2D) []dataset2D {
	// Calculate and append linear regressions for each set
	linReg, lrX := make([]dataset2D, balls), make([]float64, maxBallNum)
	for i := range lrX {
		lrX[i] = float64(i) + 1
	}

	for i, set := range data {
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
			Y:           linRegYData(lrX, set.Y, false),
		}
	}

	return linReg
}

func freqDistXLabels() []string {
	var x []string
	for i := 0; i < maxBallNum; i++ {
		x = append(x, fmt.Sprintf("%d", i+1))
	}
	return x
}

func linRegYData(lrX, lrY []float64, subzero bool) []float64 {
	a, b := stat.LinearRegression(lrX, lrY, nil, false)

	y := make([]float64, len(lrX))
	for idx, x := range lrX {
		n := a + (b * x)
		if subzero {
			y[idx] = n
		} else if n < 0 {
			y[idx] = 0
		} else {
			y[idx] = n
		}

	}
	return y
}

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
