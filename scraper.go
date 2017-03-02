package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseURL    = "https://www.lottery.co.uk"
	archiveURL = "%s/lotto/results/archive-%d"
	start      = 1994
)

// scraper is a special kind of evil. It will return a stream of data as it
// pulls all lotto result data from 1994 onwards including ball machine and
// ball set used.
func scraper() <-chan dbRow {
	var c = make(chan dbRow)

	go func() {
		// Iterate each year of the archives from start
		for y := start; y <= time.Now().Year(); y++ {
			doc, err := goquery.NewDocument(fmt.Sprintf(archiveURL, baseURL, y))
			if err != nil {
				log.Fatal(err.Error())
				break
			}

			doc.Find(".lotto tbody tr").Each(func(i int, s *goquery.Selection) {
				// skip the first row th
				if s.Children().First().Is("td") {
					var (
						row dbRow
						d   = s.Children().First() // Date
						n   = d.Next()             // Numbers
					)

					// Get details page to extract ball set and ball machine
					href, _ := d.Find("a").Attr("href")
					row.date, _ = time.Parse("02-01-2006", strings.Replace(href, "/lotto/results-", "", -1))

					// Extract machine and set
					mDoc, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseURL, href))
					if err != nil {
						log.Println(err)
					}
					mDoc.Find(".lotto tbody tr").Each(func(i int, s *goquery.Selection) {
						sText := s.Children().First().Text()
						if strings.Contains(sText, "Used:") {
							if strings.Contains(sText, "Set") {
								row.ballSet, _ = strconv.Atoi(strings.Split(sText, ": ")[1])
							} else {
								row.ballMachine = strings.Split(sText, ": ")[1]
							}
						}
					})

					// Extract numbers
					n.Find("div .result").Each(func(i int, s *goquery.Selection) {
						num, _ := strconv.Atoi(s.Text())
						row.num = append(row.num, num)
					})

					log.Println(row.date, row.num, row.ballMachine, row.ballSet)
					c <- row
				}
			})
		}
		close(c)
	}()

	return c
}
