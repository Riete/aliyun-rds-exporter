package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log"
	"net/http"

	"github.com/riete/aliyun-rds-exporter/exporter"
)

const ListenPort string = "10001"

func main() {
	rds := exporter.RdsExporter{}
	rds.InitGauge()
	registry := prometheus.NewRegistry()
	registry.MustRegister(rds)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", ListenPort), nil))
}
