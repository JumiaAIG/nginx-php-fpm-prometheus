package main

import (
)
import (
	"net/http"
	"log"
	"io/ioutil"
	"strconv"
)

func httpRequest(url string) string {
	resp, err := http.Get(url)

	scrapeFailures := 0
	scrapeFailuresBody := "scrape failure:\t" + strconv.Itoa(scrapeFailures)

	if err != nil {
		log.Println(err)
		scrapeFailures = scrapeFailures+1
		scrapeFailuresBody := "scrape failure:\t" + strconv.Itoa(scrapeFailures)
		return scrapeFailuresBody
	}

	if (resp.StatusCode != http.StatusOK){
		log.Println("Status code is not OK.")
		scrapeFailures = scrapeFailures+1
		scrapeFailuresBody := "scrape failure:\t" + strconv.Itoa(scrapeFailures)
		return scrapeFailuresBody
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		scrapeFailures = scrapeFailures+1
		scrapeFailuresBody := "scrape failure:\t" + strconv.Itoa(scrapeFailures)
		return scrapeFailuresBody
	}

	resp.Body.Close()

	return string(body) + scrapeFailuresBody
}
