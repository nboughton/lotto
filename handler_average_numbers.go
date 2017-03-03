package main

import (
	"strconv"
	//"time"

	"github.com/pilu/traffic"
	"log"
)

var (
	jsTimeFormat = "2006-01-02"
)

func handlerAverageNumbers(w traffic.ResponseWriter, r *traffic.Request) {
	var (
		p       queryParams
		start   = r.Param("start")
		end     = r.Param("end")
		set, _  = strconv.Atoi(r.Param("set"))
		machine = r.Param("machine")
	)

	log.Println(p, start, end, set, machine)
}
