package main

import (
	"fmt"
	//"sort"

	"github.com/gonum/stat"
	//"log"
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
	regressionLinear = "linear"
	regressionPoly   = "polynomial"
	graphTypeScatter = "scatter"
	graphTypeBar     = "bar"
	graphTypeLine    = "line"
)

type dataset struct {
	X           []string  `json:"x"`
	Y           []float64 `json:"y"`
	Z           []int     `json:"z"` // number
	Name        string    `json:"name"`
	Mode        string    `json:"mode"`
	Type        string    `json:"type"`
	Line        line      `json:"line"`
	Marker      marker    `json:"marker"`
	ConnectGaps bool      `json:"connectgaps"`
}

type datasetB struct {
	X           []string  `json:"x"`
	Y           []float64 `json:"y"`
	Z           []int     `json:"z"` // number
	Name        string    `json:"name"`
	Mode        string    `json:"mode"`
	Type        string    `json:"type"`
	Line        line      `json:"line"`
	Marker      markerB   `json:"marker"`
	ConnectGaps bool      `json:"connectgaps"`
}

type marker struct {
	Colour  string  `json:"color"`
	Size    float64 `json:"size"`
	Line    line    `json:"line"`
	Opacity float64 `json:"opacity"`
	Symbol  string  `json:"symbol"`
}

type markerB struct {
	Colour   string    `json:"color"`
	Size     []float64 `json:"size"`
	SizeMode string    `json:"sizemode"`
	SizeRef  float64   `json:"sizeref"`
	Line     line      `json:"line"`
	Opacity  float64   `json:"opacity"`
	Symbol   string    `json:"symbol"`
}

type line struct {
	Width   float64 `json:"width"`
	Colour  string  `json:"color"`
	Shape   string  `json:"shape"`
	Dash    string  `json:"dash"`
	Opacity float64 `json:"opacity"`
}

func graphTimeSeries(records <-chan dbRow, bestFit bool, t string) []dataset {
	data := make([]dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {

				switch t {
				case graphTypeScatter:
					data[ball] = dataset{
						Mode: "markers",
						Marker: marker{
							Colour:  colors[ball],
							Size:    6,
							Opacity: 1,
						},
					}

				case graphTypeLine:
					data[ball] = dataset{
						Mode: "lines",
						Line: line{
							Width: 1,
						},
					}
				} // END SWITCH

				data[ball].Name = label(ball)
			}
			data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			data[ball].Y = append(data[ball].Y, float64(row.Num[ball]))
		}

		i++
	}

	if bestFit && t == graphTypeScatter {
		data = append(data, regressionSets(data, regressionLinear)...)
	}

	return data
}

func graphFreqDist(records <-chan dbRow, bestFit bool, t string) []dataset {
	data := make([]dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {

				switch t {
				case graphTypeScatter:
					data[ball] = dataset{
						Mode: "markers",
						Marker: marker{
							Colour:  colors[ball],
							Size:    6,
							Opacity: 1,
						},
						X: freqDistXLabels(),
						Y: make([]float64, maxBallNum),
					}

				case graphTypeBar:
					data[ball] = dataset{
						Type: graphTypeBar,
						X:    freqDistXLabels(),
						Y:    make([]float64, maxBallNum),
					}

				} // END SWITCH

				data[ball].Name = label(ball)
			}

			data[ball].Y[row.Num[ball]-1]++
		}
		i++
	}

	/*
		if bestFit && t != graphTypeBar {
			data = append(data, regressionSets(data, regressionPoly)...)
		}
	*/

	return data
}

func graphScatter3D(records <-chan dbRow) []dataset {
	data := make([]dataset, balls)

	// x: machine
	// y: set
	// z: ball result

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {
				data[ball] = dataset{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    4,
						Opacity: 0.9,
						Line:    line{Width: 0.1},
					},
				}

				data[ball].Name = label(ball)
			}
			data[ball].X = append(data[ball].X, row.Machine)
			data[ball].Y = append(data[ball].Y, float64(row.Set))
			data[ball].Z = append(data[ball].Z, row.Num[ball])
		}
		i++
	}

	return data
}

func graphMSFreqDist(records <-chan dbRow) []dataset {
	data := make([]dataset, 12)
	// @TODO: Machine/Set frequency dist.
	/*
		i := 0
		for row := records {
			if i == 0 {
				data[]
			}
		}
	*/

	return data
}

func regressionSets(data []dataset, t string) []dataset {
	// Calculate and append regressionLinear regressions for each set
	r, rX := make([]dataset, balls), make([]float64, len(data[0].Y))

	// Generate numerical X axis data
	for i := range rX {
		rX[i] = float64(i) + 1
	}

	// Iterate existing sets and create new regression sets
	for i, set := range data {
		r[i] = dataset{
			Name: set.Name,
			Mode: "lines",
			Line: line{
				Dash:   "dot",
				Width:  2,
				Colour: colors[i],
			},
			ConnectGaps: true,
			X:           set.X,
		}

		switch t {
		case regressionLinear:
			r[i].Y = linRegYData(rX, set.Y, false)
		case regressionPoly:
			r[i].Y = polRegYData(rX, set.Y, false)
		}
	}

	return r
}

func polRegYData(prX, prY []float64, subzero bool) []float64 {
	// @TODO: polynomial regression sets for non-linear best fits
	y := make([]float64, len(prX))
	//a, b := stat.LinearRegression(prX, prY, nil, false)
	//r2 := stat.RSquared(prX, prY, nil, a, b)
	//log.Println(r2)

	return y
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

func freqDistXLabels() []string {
	var x []string
	for i := 0; i < maxBallNum; i++ {
		x = append(x, fmt.Sprintf("%d", i+1))
	}
	return x
}

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
