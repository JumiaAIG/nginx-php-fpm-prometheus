package main

import (
	"flag"
	"log"
	"gopkg.in/tylerb/graceful.v1"
	"time"
	"net/http"
)

var (
	fpmStatusURL     = ""
	nginxStatusURL   = ""
)

func main() {
	fpmUrl := flag.String("fpm-status-url", "", "PHP-FPM status URL")
	nginxUrl := flag.String("nginx-status-url", "", "Nginx status URL")
	addr := flag.String("addr", "0.0.0.0:8080", "IP/port for the HTTP server")
	flag.Parse()

	if *fpmUrl == "" {
		log.Fatal("The '-fpm-status-url' flag is required.")
	} else {
		fpmStatusURL = *fpmUrl
	}

	if *nginxUrl == "" {
		log.Fatal("The '-nginx-status-url' flag is required.")
	} else {
		nginxStatusURL = *nginxUrl
	}

	server := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:        *addr,
			ReadTimeout: time.Duration(5) * time.Second,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				
				fpmBody := httpRequest(fpmStatusURL)
				fpmBodyParsed := fpmParseBody(fpmBody)
				fpmMetrics := fpmNewMetricsFromMatches(fpmBodyParsed)
				fpmMetrics.WriteTo(w)

				nginxBody := httpRequest(nginxStatusURL)
				nginxBodyParsed := nginxParseBody(nginxBody)
				nginxMetrics := nginxNewMetricsFromMatches(nginxBodyParsed)
				nginxMetrics.WriteTo(w)
			}),
		},
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
