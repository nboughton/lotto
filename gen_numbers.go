package main

import (
	"math/rand"
	"time"
)

func generateNumbers(ranges []string) []int {
	var (
		pool    []int
		results []int
	)

	// Populate pool
	for i := 1; i <= 59; i++ {
		pool = append(pool, i)
	}

	// Set random seed
	rand.Seed(time.Now().UnixNano())

	// Select balls and remove each selected ball from the pool
	for i := 0; i < 7; i++ {
		// Select an index from the pool randomly
		r := rand.Intn(len(pool))
		// Get the value
		n := pool[r]
		// Append the new result
		results = append(results, n)
		// Remove n from pool
		pool = append(pool[:n], pool[n+1:]...)
	}

	return results
}
