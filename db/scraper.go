package db

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
	resultsURL = "%s/lotto/results"
	start      = 1994
)

// ScrapeMostRecent grabs just the most recent lotto results and passes them out.
func ScrapeMostRecent() <-chan Record {
	c := make(chan Record)

	go func() {
		// Grab the results page and then start traipsing through their awful HTML to find results links
		doc, err := goquery.NewDocument(fmt.Sprintf(resultsURL, baseURL))
		if err != nil {
			log.Println(err.Error())
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if s.HasClass("button-blue") {
				href, ok := s.Attr("href")

				if ok && strings.Contains(href, "/lotto/results-") {
					var row Record

					row.Date, err = time.Parse("02-01-2006", strings.Replace(href, "/lotto/results-", "", -1))
					if err != nil {
						log.Println("Date parse error: ", err.Error())
					}

					page, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseURL, href))
					if err != nil {
						log.Println("Page get error: ", err.Error())
					}

					page.Find("span .lotto-ball").Each(func(i int, s *goquery.Selection) {
						n, _ := strconv.Atoi(s.Text())
						row.Num = append(row.Num, n)
					})

					n, _ := strconv.Atoi(page.Find("span .lotto-bonus-ball").First().Text())
					row.Num = append(row.Num, n)

					page.Find(".lotto tr td").Each(func(i int, s *goquery.Selection) {
						if strings.Contains(s.Text(), "Used:") {
							if strings.Contains(s.Text(), "Machine") {
								row.Machine = strings.Split(s.Text(), ": ")[1]
							} else {
								row.Set, _ = strconv.Atoi(strings.Split(s.Text(), ": ")[1])
							}
						}
					})

					c <- row
				}
			}
		})

		close(c)
	}()

	return c
}

// ScrapeFullArchive scraper is a special kind of evil. It will return a stream of data as it
// pulls all lotto result data from 1994 onwards including ball machine and
// ball set used.
func ScrapeFullArchive() <-chan Record {
	c := make(chan Record)

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
						row Record
						d   = s.Children().First() // Date
						n   = d.Next()             // Numbers
					)

					// Get details page to extract ball set and ball machine
					href, _ := d.Find("a").Attr("href")
					row.Date, _ = time.Parse("02-01-2006", strings.Replace(href, "/lotto/results-", "", -1))

					// Extract machine and set
					mDoc, err := goquery.NewDocument(fmt.Sprintf("%s%s", baseURL, href))
					if err != nil {
						log.Println(err)
					}
					mDoc.Find(".lotto tbody tr").Each(func(i int, s *goquery.Selection) {
						sText := s.Children().First().Text()
						if strings.Contains(sText, "Used:") {
							if strings.Contains(sText, "Set") {
								row.Set, _ = strconv.Atoi(strings.Split(sText, ": ")[1])
							} else {
								row.Machine = strings.Split(sText, ": ")[1]
							}
						}
					})

					// Extract numbers
					n.Find("div .result").Each(func(i int, s *goquery.Selection) {
						num, _ := strconv.Atoi(s.Text())
						row.Num = append(row.Num, num)
					})

					c <- row
				}
			})
		}
		close(c)
	}()

	return c
}
