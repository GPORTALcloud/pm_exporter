package redfish

import (
	"fmt"
	"github.com/g-portal/redfish_exporter/pkg/drivers/redfish/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stmcginnis/gofish"
)

type Redfish struct {
	client *gofish.APIClient
}

func (rf *Redfish) Connect(host, username, password string, verifyTLS bool) error {
	var err error

	cfg := gofish.ClientConfig{
		Endpoint:            host,
		Username:            username,
		Password:            password,
		Insecure:            !verifyTLS,
		TLSHandshakeTimeout: 30,
	}
	// Debug
	//cfg.DumpWriter = os.Stdout

	rf.client, err = gofish.Connect(cfg)

	if err != nil {

		return fmt.Errorf("error connecting to redfish: %v", err)
	}

	return err
}

func (rf *Redfish) GetMetrics() (*prometheus.Registry, error) {
	m := metrics.NewMetrics(rf.client)
	err := m.Collect()
	if err != nil {
		return nil, fmt.Errorf("error collecting metrics: %v", err)
	}

	return m.Registry(), nil
}

func (rf *Redfish) Disconnect() error {
	rf.client.Logout()

	return nil
}
