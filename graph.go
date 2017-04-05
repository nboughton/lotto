package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gonum/stat"
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
	avgMarkerSize    = 8
	regressionLinear = "linear"
	regressionPoly   = "polynomial"
	graphTypeScatter = "scatter"
	graphTypeBar     = "bar"
	graphTypeLine    = "line"
)

/*
   Why are all axes strings? Because it's more flexible. I'm passing this data to
	 Plotly.js which interprets it accordingly and where strconv has no value it returns
	 an empty string which plotly ignores but maintains the correct number of elements
	 in the set. This removes useless visual noise from the graph and allows me to make
	 my axes anything. Yes, there's a performance overhead to the conversion but it is
	 on the scale of a fraction of a fraction of a second, which I see as a reasonable
	 compromise
*/

// standard dataset for most graphs
type dataset struct {
	X           []string `json:"x"`
	Y           []string `json:"y"`
	Z           []string `json:"z"`
	Name        string   `json:"name"`
	Mode        string   `json:"mode"`
	Type        string   `json:"type"`
	Line        line     `json:"line"`
	Marker      marker   `json:"marker"`
	ConnectGaps bool     `json:"connectgaps"`
}

// datasetB is specifically for bubble graphs because it requires a distinctly different
// marker object
type datasetB struct {
	X           []string `json:"x"`
	Y           []string `json:"y"`
	Z           []string `json:"z"`
	Name        string   `json:"name"`
	Mode        string   `json:"mode"`
	Type        string   `json:"type"`
	Line        line     `json:"line"`
	Marker      markerB  `json:"marker"`
	ConnectGaps bool     `json:"connectgaps"`
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

func graphResultsTimeSeries(records <-chan dbRow, bestFit bool, t string) []dataset {
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
							Size:    avgMarkerSize,
							Opacity: 1,
						},
					}

				case graphTypeLine:
					data[ball] = dataset{
						Mode: "lines",
						Line: line{
							Width: 1,
							Shape: "spline",
						},
					}
				} // END SWITCH

				data[ball].Name = label(ball)
			}
			data[ball].X = append(data[ball].X, fmt.Sprintf("%s:%d:%s", row.Date.Format(formatYYYYMMDD), row.Set, row.Machine))
			data[ball].Y = append(data[ball].Y, strconv.Itoa(row.Num[ball]))
		}

		i++
	}

	// Best fit lines only really suit scatter graphs so far
	if bestFit && t == graphTypeScatter {
		data = append(data, regressionSet(data, regressionLinear)...)
	}

	return data
}

func graphResultsFreqDist(records <-chan dbRow, bestFit bool, t string) []dataset {
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
							Size:    avgMarkerSize,
							Opacity: 1,
						},
						X: freqDistXLabels(),
						Y: make([]string, maxBallNum),
					}

				case graphTypeBar:
					data[ball] = dataset{
						Type: graphTypeBar,
						X:    freqDistXLabels(),
						Y:    make([]string, maxBallNum),
					}

				} // END SWITCH

				data[ball].Name = label(ball)
			}

			n, _ := strconv.Atoi(data[ball].Y[row.Num[ball]-1]) // I don't care about this error because it actually eliminates useless 0 values
			data[ball].Y[row.Num[ball]-1] = strconv.Itoa(n + 1)
		}
		i++
	}

	/* currently not working because I have yet to work out how to do non-linear best fit lines
	if bestFit && t != graphTypeBar {
		data = append(data, regressionSet(data, regressionPoly)...)
	}
	*/

	return data
}

func graphResultsRawScatter3D(records <-chan dbRow) []dataset {
	data := make([]dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {
				data[ball] = dataset{
					Mode: "markers",
					Type: "scatter3d",
					Marker: marker{
						Size:    avgMarkerSize / 2,
						Opacity: 0.9,
						Line:    line{Width: 0.1},
					},
				}

				data[ball].Name = label(ball)
			}
			data[ball].X = append(data[ball].X, row.Machine)
			data[ball].Y = append(data[ball].Y, strconv.Itoa(row.Set))
			data[ball].Z = append(data[ball].Z, strconv.Itoa(row.Num[ball]))
		}
		i++
	}

	return data
}

func graphMSFreqDistScatter3D(m map[string]int) []dataset {
	data := dataset{
		Type: "scatter3d",
		Mode: "markers",
		Marker: marker{
			Size:    avgMarkerSize,
			Opacity: 0.9,
			Line: line{
				Width: 0.1,
			},
		},
	}

	l := []string{}
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)

	for _, k := range l {
		s := strings.Split(k, ":")

		data.X = append(data.X, s[0])               // Machine
		data.Y = append(data.Y, s[1])               // Set
		data.Z = append(data.Z, strconv.Itoa(m[k])) // Frequency
	}

	return []dataset{data}
}

func graphMSFreqDistBubble(m map[string]int) []datasetB {
	data := datasetB{
		Type: "scatter",
		Mode: "markers",
		Marker: markerB{
			Opacity: 0.8,
			Line:    line{Width: 0.1},
		},
	}

	l := []string{}
	for k := range m {
		l = append(l, k)
	}
	sort.Strings(l)

	for _, k := range l {
		s := strings.Split(k, ":")

		data.X = append(data.X, s[0])                                // Machine
		data.Y = append(data.Y, s[1])                                // Set
		data.Marker.Size = append(data.Marker.Size, float64(m[k])*2) // Frequency
	}

	return []datasetB{data}
}

func regressionSet(data []dataset, t string) []dataset {
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

		// Translate Y axis to numerical data
		rY := make([]float64, len(data[0].Y))
		for j := range rY {
			rY[j], _ = strconv.ParseFloat(data[i].Y[j], 64)
		}
		switch t {
		case regressionLinear:
			r[i].Y = linRegYData(rX, rY, false)

		case regressionPoly:
			r[i].Y = polRegYData(rX, rY, false) // Not currently working because I suck at maths

		}
	}

	return r
}

func polRegYData(prX, prY []float64, subzero bool) []string {
	// @TODO: polynomial regression sets for non-linear best fits
	y := make([]string, len(prX))
	//a, b := stat.LinearRegression(prX, prY, nil, false)
	//r2 := stat.RSquared(prX, prY, nil, a, b)
	//log.Println(r2)

	return y
}

func linRegYData(lrX, lrY []float64, subzero bool) []string {
	a, b := stat.LinearRegression(lrX, lrY, nil, false)

	y := make([]string, len(lrX))
	for idx, x := range lrX {
		n := a + (b * x)
		f := strconv.FormatFloat(n, 'f', -1, 64)
		if subzero {
			y[idx] = f
		} else if n < 0 {
			y[idx] = "0"
		} else {
			y[idx] = f
		}

	}
	return y
}

func freqDistXLabels() []string {
	var x []string
	for i := 0; i < maxBallNum; i++ {
		x = append(x, strconv.Itoa(i+1))
	}
	return x
}

func label(ball int) string {
	if ball < 6 {
		return fmt.Sprintf("Ball %d", ball+1)
	}

	return "Bonus"
}
