package main

import (
	"net/http"

	"github.com/burningsunrise/tplink-exporter/collector"
	"github.com/burningsunrise/tplink-exporter/formatter"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	format := &formatter.Formatter{
		HideKeys:      false,
		FieldsOrder:   []string{"topic", "message"},
		ShowFullLevel: true,
		TrimMessages:  true,
	}

	log.SetLevel(log.InfoLevel)
	log.SetFormatter(format)
}

func main() {
	tplinkCollector := collector.NewTplinkCollector()
	prometheus.MustRegister(tplinkCollector)

	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :9797")
	log.Fatal(http.ListenAndServe(":9797", nil))
}
