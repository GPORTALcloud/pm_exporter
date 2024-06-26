package main

import (
	"flag"
	"github.com/g-portal/redfish_exporter/pkg/config"
	"github.com/g-portal/redfish_exporter/pkg/metric"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	listenAddr string
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config.file", "/etc/redfish_exporter/config.yml", "Defines the path to the platform management config")
	flag.StringVar(&listenAddr, "web.listen-address", "0.0.0.0:9096", "Address the exporter listens on")
	flag.Parse()

	config.SetPath(configPath)
}

func main() {
	r := gin.Default()

	r.GET("/metrics", metric.Handle)
	log.Println("Starting listening on: " + listenAddr)
	err := r.Run(listenAddr)
	if err != nil {
		log.Printf("Error starting http server: %v", err)
	}
}
