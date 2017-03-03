package main

import (
//"fmt"
//"time"
)

type graphData struct {
	Labels   []string  `json:"labels"`
	Datasets []dataset `json:"datasets"`
}

type dataset struct {
	Label string `json:"label"`
	Data  []int  `json:"data"`
}
