package updater

import (
	"log"
	"time"

	"github.com/GPORTALcloud/pm_exporter/pkg/api"
	"github.com/GPORTALcloud/pm_exporter/pkg/config"
)

var clients = make(chan api.Client)

// worker Function used for starting a worker pool
func worker(clients <-chan api.Client) {
	for c := range clients {
		log.Printf("Fetching metrics from %s", c.GetHost())
		err := c.FetchInventoryMetrics()
		if err != nil {
			log.Printf("Error updating %s: %v", c.GetHost(), err)
		}
	}
}

// Run Start the worker pool and add api clients to the chan every $refreshInterval
func Run(workerCount int, refreshInterval time.Duration) {
	for w := 0; w < workerCount; w++ {
		go worker(clients)
	}
	for {
		cfg := config.GetConfig()
		for i, _ := range cfg.PlatformManagements {
			switch pm := cfg.PlatformManagements[i].Type; pm {
			case "IDRAC":
				fallthrough
			case "IDRAC9":
				fallthrough
			case "IDRAC8":
				client := api.NewRedfishAPI(cfg.PlatformManagements[i].Host)
				client.SetUser(cfg.PlatformManagements[i].Username, cfg.PlatformManagements[i].Password)
				clients <- client
			case "ILO5":
				client := api.NewHpeApi(cfg.PlatformManagements[i].Host)
				client.SetUser(cfg.PlatformManagements[i].Username, cfg.PlatformManagements[i].Password)
				clients <- client
				continue
			default:
				log.Printf("Error fetching inventory for %v - type unknown (%v)\n", cfg.PlatformManagements[i].Host, cfg.PlatformManagements[i].Type)
			}
		}
		time.Sleep(refreshInterval)
	}
}
