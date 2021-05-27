package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/GPORTALcloud/pm_exporter/pkg/config"
	"github.com/GPORTALcloud/pm_exporter/pkg/metric"
	"github.com/GPORTALcloud/pm_exporter/pkg/updater"
)

var (
	workerCount     int
	refreshInterval time.Duration
	persistDuration time.Duration
	listenAddr      string
	enableLifecycle bool
	configPath      string
)

func init() {
	flag.StringVar(&configPath, "config.file", "/etc/pm_exporter.yml", "Defines the path to the platform management config")
	flag.StringVar(&listenAddr, "web.listen-address", "0.0.0.0:9096", "Address the exporter listens on")
	flag.BoolVar(&enableLifecycle, "web.enable-lifecycle", false, "With this parameter set calls to /-/reload are reloading the config")
	flag.IntVar(&workerCount, "worker.count", 10, "Worker processes for calling platform management API's")
	flag.DurationVar(&refreshInterval, "metric.refresh_interval", time.Second*60, "Interval the exporter refresh the metrics")
	flag.DurationVar(&persistDuration, "metric.persist_duration", time.Second*90, "Duration collected metrics persist before being invalidated")
	flag.Parse()
	config.Prepare(configPath)
}

func main() {

	metric.SetPersistDuration(persistDuration)

	go updater.Run(workerCount, refreshInterval)
	http.HandleFunc("/metrics", metric.MetricHttpHandler)

	if enableLifecycle {
		http.HandleFunc("/-/reload", config.ReloadConfigHandler)
	}

	log.Println("Starting listening on: " + listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Printf("Error starting http server: %v", err)
	}
}
