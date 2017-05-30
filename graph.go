package main

import (
	"fmt"
	//"sort"
	"strconv"
	//"strings"
	//"github.com/gonum/stat"
	//pt "github.com/nboughton/go-plotlytypes"
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

type graphDataset struct {
	Label string   `json:"label"`
	Data  []string `json:"data"`
}

type graphData struct {
	Labels   []string       `json:"labels"`
	Datasets []graphDataset `json:"datasets"`
}

func lineGraph(records <-chan dbRow) graphData {
	var d graphData
	d.Datasets = make([]graphDataset, balls)

	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if d.Datasets[ball].Label == "" {
				d.Datasets[ball].Label = fmt.Sprintf("Ball %d", ball+1)
			}

			d.Datasets[ball].Data = append(d.Datasets[ball].Data, strconv.Itoa(row.Num[ball]))
		}
		d.Labels = append(d.Labels, fmt.Sprintf("%s:%s:%d", row.Date.Format(formatYYYYMMDD), row.Machine, row.Set))
	}

	return d
}

/*
func graphResultsTimeSeries(records <-chan dbRow, bestFit bool, t string) []pt.Dataset {
	data := make([]pt.Dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {

				switch t {
				case graphTypeScatter:
					data[ball] = pt.Dataset{
						Mode: "markers",
						Marker: pt.Marker{
							Colour:  colors[ball],
							Size:    avgMarkerSize,
							Opacity: 1,
						},
					}

				case graphTypeLine:
					data[ball] = pt.Dataset{
						Mode: "lines",
						Line: pt.Line{
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

func graphResultsFreqDist(records <-chan dbRow, bestFit bool, t string) []pt.Dataset {
	data := make([]pt.Dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {

				switch t {
				case graphTypeScatter:
					data[ball] = pt.Dataset{
						Mode: "markers",
						Marker: pt.Marker{
							Colour:  colors[ball],
							Size:    avgMarkerSize,
							Opacity: 1,
						},
						X: freqDistXLabels(),
						Y: make([]string, maxBallNum),
					}

				case graphTypeBar:
					data[ball] = pt.Dataset{
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
	//

	return data
}

func graphResultsRawScatter3D(records <-chan dbRow) []pt.Dataset {
	data := make([]pt.Dataset, balls)

	i := 0
	for row := range records {
		for ball := 0; ball < balls; ball++ {
			if i == 0 {
				data[ball] = pt.Dataset{
					Mode: "markers",
					Type: "scatter3d",
					Marker: pt.Marker{
						Size:    avgMarkerSize / 2,
						Opacity: 0.9,
						Line:    pt.Line{Width: 0.1},
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

func graphMSFreqDistScatter3D(m map[string]int) []pt.Dataset {
	data := pt.Dataset{
		Type: "scatter3d",
		Mode: "markers",
		Marker: pt.Marker{
			Size:    avgMarkerSize,
			Opacity: 0.9,
			Line: pt.Line{
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

	return []pt.Dataset{data}
}

func graphMSFreqDistBubble(m map[string]int) []pt.DatasetB {
	data := pt.DatasetB{
		Type: "scatter",
		Mode: "markers",
		Marker: pt.MarkerB{
			Opacity: 0.8,
			Line:    pt.Line{Width: 0.1},
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

	return []pt.DatasetB{data}
}

func regressionSet(data []pt.Dataset, t string) []pt.Dataset {
	// Calculate and append regressionLinear regressions for each set
	r, rX := make([]pt.Dataset, balls), make([]float64, len(data[0].Y))

	// Generate numerical X axis data
	for i := range rX {
		rX[i] = float64(i) + 1
	}

	// Iterate existing sets and create new regression sets
	for i, set := range data {
		r[i] = pt.Dataset{
			Name: set.Name,
			Mode: "lines",
			Line: pt.Line{
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
*/
