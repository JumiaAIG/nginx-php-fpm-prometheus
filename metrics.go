package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)


// Metrics contains the status data collected from PHP-FPM.
type Metrics struct {
	FpmStartSince         		int `help:"FPM: Seconds since FPM start" type:"counter" name:"php_fpm_start_since"`
	FpmAcceptedConn       		int `help:"FPM: Total of accepted connections" type:"counter" name:"php_fpm_accepted_conn"`
	FpmListenQueue        		int `help:"FPM: Number of connections that have been initiated but not yet accepted" type:"gauge" name:"php_fpm_listen_queue"`
	FpmMaxListenQueue     		int `help:"FPM: Max. connections the listen queue has reached since FPM start" type:"counter" name:"php_fpm_max_listen_queue"`
	FpmListenQueueLength  		int `help:"FPM: Maximum number of connections that can be queued" type:"gauge" name:"php_fpm_listen_queue_length"`
	FpmIdleProcesses      		int `help:"FPM: Idle process count" type:"gauge" name:"php_fpm_idle_processes"`
	FpmActiveProcesses    		int `help:"FPM: Active process count" type:"gauge" name:"php_fpm_active_processes"`
	FpmTotalProcesses     		int `help:"FPM: Total process count" type:"gauge" name:"php_fpm_total_processes"`
	FpmMaxActiveProcesses 		int `help:"FPM: Maximum active process count" type:"counter" name:"php_fpm_max_active_processes"`
	FpmMaxChildrenReached 		int `help:"FPM: Number of times the process limit has been reached" type:"counter" name:"php_fpm_max_children_reached"`
	FpmSlowRequests       		int `help:"FPM: Number of requests that exceed request_slowlog_timeout" type:"counter" name:"php_fpm_slow_requests"`
	FpmScrapeFailures     		int `help:"FPM: Number of errors while scraping php_fpm" type:"counter" name:"php_fpm_exporter_scrape_failures_total"` //gauge?
	NginxActiveConnections		int `help:"NGINX: Number of active client connections including Waiting connections." type:"gauge" name:"nginx_active_connections"`
	NginxAcceptedConnections	int `help:"NGINX: Total number of accepted client connections" type:"counter" name:"nginx_accepted_connections"`
	NginxHandledConnections		int `help:"NGINX: Total number of handled connections. Generally, the parameter value is the same as accepts unless some resource limits have been reached" type:"counter" name:"nginx_handled_connections"`
	NginxNumberOfRequests		int `help:"NGINX: Total number of client requests" type:"counter" name:"nginx_number_of_requests"`
	NginxConnectionsReading		int `help:"NGINX: Number of connections where nginx is reading the request header" type:"gauge" name:"nginx_connections_reading"`
	NginxConnectionsWriting		int `help:"NGINX: Number of connections where nginx is writing the response back to the client" type:"gauge" name:"nginx_connections_writing"`
	NginxConnectionsWaiting		int `help:"NGINX: Number of idle client connections waiting for a request" type:"gauge" name:"nginx_connections_waiting"`
	NginxScrapeFailures    		int `help:"NGINX: Number of errors while scraping Nginx" type:"counter" name:"nginx_exporter_scrape_failures_total"` //gauge?
}

func (m *Metrics) WriteTo(w io.Writer) {
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
