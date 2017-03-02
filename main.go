package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	baseURL = "https://www.lottery.co.uk/lotto/results/archive-1994"
)

/* @TODO: rewrite generator to create num array and remove selected numbers every time
write scraper and db code to generate aggregate dataset with web front end to
view data.
*/

func main() {
	max, res := 59, [][]int{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < 1000; i++ {
		m := []int{}
		for len(m) < 6 {
			n := rand.Intn(max) + 1
			if !contains(n, m) {
				m = append(m, n)
			}
		}
		res = append(res, m)
		//time.Sleep(time.Second)
	}
	fmt.Println(res[rand.Intn(len(res))])
}

func contains(n int, s []int) bool {
	for _, v := range s {
		if v == n {
			return true
		}
	}
	return false
}
