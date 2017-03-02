package main

import (
	"github.com/pilu/traffic"
)

func handlerRoot(w traffic.ResponseWriter, r *traffic.Request) {
	w.Render("index")
}
