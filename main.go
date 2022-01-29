package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Starting Smart Hub 2 Exporter")

	ip := os.Getenv("SMARTHUB2_IP")
	if len(ip) == 0 {
		ip = "192.168.1.254"
	}

	registry := prometheus.NewRegistry()
	collector := newSmartHub2Collector(ip)

	registry.MustRegister(collector)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(`<html><head><title>Smart Hub 2 Exporter</title></head>` +
			`<body><a href="metrics">Metrics</a></body></html>`))
	})
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Println("Beginning to serve on port :9906")
	log.Fatal(http.ListenAndServe(":9906", nil))
}
