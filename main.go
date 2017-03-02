package main

import (
	"flag"

	"github.com/pilu/traffic"
)

var (
	router *traffic.Router
	db     = connectDB("results.db")
)

func init() {
	p := flag.Int("p", 3002, "Set the port the application listens on")
	flag.Parse()

	traffic.SetPort(*p)

	router = traffic.New()
	router.Get("/", handlerRoot)
}

func main() {
	router.Run()
}
