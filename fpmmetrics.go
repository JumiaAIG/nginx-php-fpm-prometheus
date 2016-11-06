package main

import (
	"strconv"
	"regexp"
)

var (
	fpmStatusLineRegexp = regexp.MustCompile(`(?m)^(.*):\s+(.*)$`)
)

// NewMetricsFromMatches creates a new Metrics instance and populates it with given data.
func fpmNewMetricsFromMatches(matches [][]string) *Metrics {
	metrics := &Metrics{}
	metrics.fpmPopulateFromMatches(matches)
	return metrics
}

func (m *Metrics) fpmPopulateFromMatches(matches [][]string) {
	for _, match := range matches {
		key := match[1]
		value := match[2]
		switch key {
		case "start since":
			m.FpmStartSince, _ = strconv.Atoi(value)
		case "accepted conn":
			m.FpmAcceptedConn, _ = strconv.Atoi(value)
		case "listen queue":
			m.FpmListenQueue, _ = strconv.Atoi(value)
		case "max listen queue":
			m.FpmMaxListenQueue, _ = strconv.Atoi(value)
		case "listen queue len":
			m.FpmListenQueueLength, _ = strconv.Atoi(value)
		case "idle processes":
			m.FpmIdleProcesses, _ = strconv.Atoi(value)
		case "active processes":
			m.FpmActiveProcesses, _ = strconv.Atoi(value)
		case "total processes":
			m.FpmTotalProcesses, _ = strconv.Atoi(value)
		case "max active processes":
			m.FpmMaxActiveProcesses, _ = strconv.Atoi(value)
		case "max children reached":
			m.FpmMaxChildrenReached, _ = strconv.Atoi(value)
		case "slow requests":
			m.FpmSlowRequests, _ = strconv.Atoi(value)
		case "scrape failure":
			m.FpmScrapeFailures, _ = strconv.Atoi(value)
		}

	}
}

func fpmParseBody(body string) [][]string {
	matches := fpmStatusLineRegexp.FindAllStringSubmatch(string(body), -1)
	return matches
}
