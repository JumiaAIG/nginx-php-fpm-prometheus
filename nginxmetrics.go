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
	match1LineRegexp = regexp.MustCompile(`Active connections:\s+(\d+)`)
	match2LineRegexp = regexp.MustCompile(`\s*(\d+)\s+(\d+)\s+(\d+)`)
	match3LineRegexp = regexp.MustCompile(`Reading:\s*(\d+)\s*Writing:\s*(\d+)\s*Waiting:\s*(\d+)`)
	match4LineRegexp = regexp.MustCompile(`scrape failure:\s*(\d+)`)
)

// NewMetricsFromMatches creates a new Metrics instance and populates it with given data.
func nginxNewMetricsFromMatches(matches [][]string) *NginxMetrics {
	metrics := &NginxMetrics{}
	metrics.nginxPopulateFromMatches(matches)
	return metrics
}

func (m *NginxMetrics) nginxPopulateFromMatches(matches [][]string) {
	for _, match := range matches {
		key := match[0]
		value := match[1]
		switch key {
		case "active connections":
			m.ActiveConnections, _ = strconv.Atoi(value)
		case "accepted connections":
			m.AcceptedConnections, _ = strconv.Atoi(value)
		case "handled connections":
			m.HandledConnections, _ = strconv.Atoi(value)
		case "number of requests":
			m.NumberOfRequests, _ = strconv.Atoi(value)
		case "connections reading":
			m.ConnectionsReading, _ = strconv.Atoi(value)
		case "connections writing":
			m.ConnectionsWriting, _ = strconv.Atoi(value)
		case "connections waiting":
			m.ConnectionsWaiting, _ = strconv.Atoi(value)
		case "scrape failure":
			m.ScrapeFailures, _ = strconv.Atoi(value)
		}

	}
}

func nginxParseBody(body string) [][]string {
	// http://stackoverflow.com/questions/25025467/catching-panics-in-golang
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

type NginxMetrics struct {
	ActiveConnections	int `help:"NGINX: Number of active client connections including Waiting connections." type:"gauge" name:"nginx_active_connections"`
	AcceptedConnections	int `help:"NGINX: Total number of accepted client connections" type:"counter" name:"nginx_accepted_connections"`
	HandledConnections	int `help:"NGINX: Total number of handled connections. Generally, the parameter value is the same as accepts unless some resource limits have been reached" type:"counter" name:"nginx_handled_connections"`
	NumberOfRequests	int `help:"NGINX: Total number of client requests" type:"counter" name:"nginx_number_of_requests"`
	ConnectionsReading	int `help:"NGINX: Number of connections where nginx is reading the request header" type:"gauge" name:"nginx_connections_reading"`
	ConnectionsWriting	int `help:"NGINX: Number of connections where nginx is writing the response back to the client" type:"gauge" name:"nginx_connections_writing"`
	ConnectionsWaiting	int `help:"NGINX: Number of idle client connections waiting for a request" type:"gauge" name:"nginx_connections_waiting"`
	ScrapeFailures    	int `help:"NGINX: Number of errors while scraping Nginx" type:"counter" name:"nginx_exporter_scrape_failures_total"` //gauge?
}

func (m *NginxMetrics) NginxWriteTo(w io.Writer) {
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
