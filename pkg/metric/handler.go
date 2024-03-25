package metric

import (
	"github.com/g-portal/redfish_exporter/pkg/api"
	"github.com/g-portal/redfish_exporter/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func Handle(c *gin.Context) {
	params := extractCollectorParams(c.Request)

	client, err := api.NewClient(params.Host, params.Username, params.Password, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("finished with handling")
	defer func() {
		err := client.Disconnect()
		if err != nil {
			log.Printf("error disconnecting: %v", err)
		}
	}()

	// Get metrics from the client.
	registry, err := client.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

type collectorParams struct {
	Username string
	Password string
	Host     string
}

func extractCollectorParams(r *http.Request) collectorParams {
	cfg := config.GetConfig()
	params := collectorParams{}
	if host := r.URL.Query().Get("host"); host != "" {
		params.Host = host
	}
	if username := r.URL.Query().Get("username"); username != "" {
		params.Username = username
	} else {
		params.Username = cfg.Redfish.Username
	}
	if password := r.URL.Query().Get("password"); password != "" {
		params.Password = password
	} else {
		params.Password = cfg.Redfish.Password
	}

	return params
}
