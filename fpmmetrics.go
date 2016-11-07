package main

import (
	"strconv"
	"regexp"
	"io"
	"reflect"
	"bytes"
	"fmt"
)

var (
	fpmStatusLineRegexp = regexp.MustCompile(`(?m)^(.*):\s+(.*)$`)
)

// NewMetricsFromMatches creates a new Metrics instance and populates it with given data.
func fpmNewMetricsFromMatches(matches [][]string) *FpmMetrics {
	metrics := &FpmMetrics{}
	metrics.fpmPopulateFromMatches(matches)
	return metrics
}

func (m *FpmMetrics) fpmPopulateFromMatches(matches [][]string) {
	for _, match := range matches {
		key := match[1]
		value := match[2]
		switch key {
		case "start since":
			m.StartSince, _ = strconv.Atoi(value)
		case "accepted conn":
			m.AcceptedConn, _ = strconv.Atoi(value)
		case "listen queue":
			m.ListenQueue, _ = strconv.Atoi(value)
		case "max listen queue":
			m.MaxListenQueue, _ = strconv.Atoi(value)
		case "listen queue len":
			m.ListenQueueLength, _ = strconv.Atoi(value)
		case "idle processes":
			m.IdleProcesses, _ = strconv.Atoi(value)
		case "active processes":
			m.ActiveProcesses, _ = strconv.Atoi(value)
		case "total processes":
			m.TotalProcesses, _ = strconv.Atoi(value)
		case "max active processes":
			m.MaxActiveProcesses, _ = strconv.Atoi(value)
		case "max children reached":
			m.MaxChildrenReached, _ = strconv.Atoi(value)
		case "slow requests":
			m.SlowRequests, _ = strconv.Atoi(value)
		case "scrape failure":
			m.ScrapeFailures, _ = strconv.Atoi(value)
		}

	}
}

func fpmParseBody(body string) [][]string {
	matches := fpmStatusLineRegexp.FindAllStringSubmatch(string(body), -1)
	return matches
}

type FpmMetrics struct {
	StartSince         int `help:"FPM: Seconds since FPM start" type:"counter" name:"php_fpm_start_since"`
	AcceptedConn       int `help:"FPM: Total of accepted connections" type:"counter" name:"php_fpm_accepted_conn"`
	ListenQueue        int `help:"FPM: Number of connections that have been initiated but not yet accepted" type:"gauge" name:"php_fpm_listen_queue"`
	MaxListenQueue     int `help:"FPM: Max. connections the listen queue has reached since FPM start" type:"counter" name:"php_fpm_max_listen_queue"`
	ListenQueueLength  int `help:"FPM: Maximum number of connections that can be queued" type:"gauge" name:"php_fpm_listen_queue_length"`
	IdleProcesses      int `help:"FPM: Idle process count" type:"gauge" name:"php_fpm_idle_processes"`
	ActiveProcesses    int `help:"FPM: Active process count" type:"gauge" name:"php_fpm_active_processes"`
	TotalProcesses     int `help:"FPM: Total process count" type:"gauge" name:"php_fpm_total_processes"`
	MaxActiveProcesses int `help:"FPM: Maximum active process count" type:"counter" name:"php_fpm_max_active_processes"`
	MaxChildrenReached int `help:"FPM: Number of times the process limit has been reached" type:"counter" name:"php_fpm_max_children_reached"`
	SlowRequests       int `help:"FPM: Number of requests that exceed request_slowlog_timeout" type:"counter" name:"php_fpm_slow_requests"`
	ScrapeFailures     int `help:"FPM: Number of errors while scraping php_fpm" type:"counter" name:"php_fpm_exporter_scrape_failures_total"` //gauge?
}

func (m *FpmMetrics) FpmWriteTo(w io.Writer) {
	typ := reflect.TypeOf(*m)
	val := reflect.ValueOf(*m)
	buf := &bytes.Buffer{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("name")
		buf.WriteString(fmt.Sprintf("# HELP %s %s\n", name, field.Tag.Get("help")))
		buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", name, field.Tag.Get("type")))
		buf.WriteString(fmt.Sprintf("%s %d\n", name, val.Field(i).Int()))
	}

	io.Copy(w, buf)
}