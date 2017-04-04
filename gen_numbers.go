package main

import (
	"math/rand"
	"sort"
	"time"
)

func drawRandomSet() []int {
	var (
		pool    []int
		results []int
	)

	// Populate pool
	for i := 1; i <= 59; i++ {
		pool = append(pool, i)
	}

	// Select balls and remove each selected ball from the pool
	for i := 0; i < 6; i++ {
		pool, results = drawBall(pool, results)
	}

	sort.Ints(results)

	// One last time for the bonus ball
	pool, results = drawBall(pool, results)

	return results
}

func drawBall(p, r []int) ([]int, []int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(p))
	// Append the new result
	r = append(r, p[n])
	// Remove n from pool
	p = append(p[:n], p[n+1:]...)

	return p, r
}
