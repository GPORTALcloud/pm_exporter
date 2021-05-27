package api

import (
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 15 * time.Second}

type Client interface {
	GetHost() string
	FetchInventoryMetrics() error
}
