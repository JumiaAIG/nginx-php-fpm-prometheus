package main

import (
	"strconv"
	"regexp"
)

var (
	match1LineRegexp = regexp.MustCompile(`Active connections:\s+(\d+)`)
	match2LineRegexp = regexp.MustCompile(`\s*(\d+)\s+(\d+)\s+(\d+)`)
	match3LineRegexp = regexp.MustCompile(`Reading:\s*(\d+)\s*Writing:\s*(\d+)\s*Waiting:\s*(\d+)`)
	match4LineRegexp = regexp.MustCompile(`scrape failure:\s*(\d+)`)
)

// NewMetricsFromMatches creates a new Metrics instance and populates it with given data.
func nginxNewMetricsFromMatches(matches [][]string) *Metrics {
	metrics := &Metrics{}
	metrics.nginxPopulateFromMatches(matches)
	return metrics
}

func (m *Metrics) nginxPopulateFromMatches(matches [][]string) {
	for _, match := range matches {
		key := match[0]
		value := match[1]
		switch key {
		case "active connections":
			m.NginxActiveConnections, _ = strconv.Atoi(value)
		case "accepted connections":
			m.NginxAcceptedConnections, _ = strconv.Atoi(value)
		case "handled connections":
			m.NginxHandledConnections, _ = strconv.Atoi(value)
		case "listen queue":
			m.NginxNumberOfRequests, _ = strconv.Atoi(value)
		case "max listen queue":
			m.NginxConnectionsReading, _ = strconv.Atoi(value)
		case "listen queue len":
			m.NginxConnectionsWriting, _ = strconv.Atoi(value)
		case "idle processes":
			m.NginxConnectionsWaiting, _ = strconv.Atoi(value)
		case "scrape failure":
			m.NginxScrapeFailures, _ = strconv.Atoi(value)
		}

	}
}

func nginxParseBody(body string) [][]string {
	match1 := match1LineRegexp.FindAllStringSubmatch(string(body), -1)
	matches := [][]string{{"active connections",match1[0][1]}}

	match2 := match2LineRegexp.FindAllStringSubmatch(string(body), -1)
	matches = append(matches,[]string{"accepted connections",match2[0][1]})
	matches = append(matches,[]string{"handled connections",match2[0][2]})
	matches = append(matches,[]string{"number of requests",match2[0][3]})

	match3 := match3LineRegexp.FindAllStringSubmatch(string(body), -1)
	matches = append(matches,[]string{"connections reading",match3[0][1]})
	matches = append(matches,[]string{"connections writing",match3[0][2]})
	matches = append(matches,[]string{"connections waiting",match3[0][3]})

	match4 := match4LineRegexp.FindAllStringSubmatch(string(body), -1)
	matches = append(matches,[]string{"scrape failure",match4[0][1]})

	return matches
}
