package metric

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/GPORTALcloud/pm_exporter/pkg/config"
)

const (
	PMPlatformManagementUp = "pm_platform_management_up"
	PMPowerSupplyHealth    = "pm_power_supply_health"
	PMBatteryHealth        = "pm_battery_health"
	PMCpuHealth            = "pm_cpu_health"
	PMFanHealth            = "pm_fan_health"
	PMStorageHealth        = "pm_storage_health"
	PMTemperatureHealth    = "pm_temperature_health"
	PMIntrusionHealth      = "pm_intrusion_health"
	PMLicenceHealth        = "pm_license_health"
	PMMemoryHealth         = "pm_memory_health"
	PMOverallHealth        = "pm_overall_health"
)

var metricStore = Storage{
	Mutex:           &sync.Mutex{},
	persistDuration: 60,
	Metrics:         map[string]Metric{},
}

func init() {
	RegisterMetric(PMPlatformManagementUp, "gauge", "Platform Management reachable")
	RegisterMetric(PMPowerSupplyHealth, "gauge", "Power Supply Health status")
	RegisterMetric(PMBatteryHealth, "gauge", "Battery Health status")
	RegisterMetric(PMCpuHealth, "gauge", "CPU Health status")
	RegisterMetric(PMFanHealth, "gauge", "Fan Health status")
	RegisterMetric(PMStorageHealth, "gauge", "Storage Health status")
	RegisterMetric(PMTemperatureHealth, "gauge", "Temperature Health status")
	RegisterMetric(PMIntrusionHealth, "gauge", "Intrusion Health status")
	RegisterMetric(PMLicenceHealth, "gauge", "License Health status")
	RegisterMetric(PMMemoryHealth, "gauge", "Memory Health status")
	RegisterMetric(PMOverallHealth, "gauge", "Overall Health status")
}

// SetPersistDuration overwrite duration the metrics are getting persisted within the metricStore
func SetPersistDuration(duration time.Duration) {
	metricStore.persistDuration = duration
}

// RegisterMetric used for adding new metric types to the metricStore
func RegisterMetric(name string, t string, help string) {
	metricStore.Mutex.Lock()
	defer metricStore.Mutex.Unlock()
	if _, ok := metricStore.Metrics[name]; ok {
		return
	}
	metricStore.Metrics[name] = Metric{Name: name, Type: t, Help: help, Vars: map[string]MetricVar{}}
}

// UpdateMetric adds new metric values to the metricStore
func UpdateMetric(name string, host string, val interface{}) {
	metricStore.Mutex.Lock()
	defer metricStore.Mutex.Unlock()
	metricStore.Metrics[name].Vars[host] = MetricVar{
		Value:     val,
		ExpiresAt: time.Now().Add(metricStore.persistDuration).Unix(),
	}
}

// UnsetMetric removes specific host metrics from the metricStore
func UnsetMetric(name string, host string) {
	delete(metricStore.Metrics[name].Vars, host)
}

// MetricHttpHandler HTTP handler for returning the metrics formatted for prometheus
func MetricHttpHandler(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte(dump()))
	if err != nil {
		log.Printf("Error writing to ResponseWriter: %v\n", err)
	}
}

// dump Utility function used by the HTTP handler for returning a proper prometheus format
func dump() string {
	metricStore.Mutex.Lock()
	defer metricStore.Mutex.Unlock()
	responseString := ""
	for i, m := range metricStore.Metrics {
		responseString += fmt.Sprintf("# HELP %v %v\n", m.Name, m.Help)
		responseString += fmt.Sprintf("# TYPE %v %v\n", m.Name, m.Type)
		for h, d := range metricStore.Metrics[i].Vars {
			nodeConfig := config.GetManagementConfig(h)
			if d.ExpiresAt <= time.Now().Unix() {
				UnsetMetric(m.Name, h)
				continue
			}
			labels := map[string]string{}
			if nodeConfig == nil {
				UnsetMetric(m.Name, h)
				continue
			} else {
				labels = nodeConfig.Labels
			}
			labels["host"] = h
			labelGroups := []string{}
			for key, value := range labels {
				labelGroups = append(labelGroups, fmt.Sprintf("%v=\"%v\"", key, value))
			}
			labelString := fmt.Sprintf("{%v}", strings.Join(labelGroups, ","))
			responseString += fmt.Sprintf("%v%v %v\n", m.Name, string(labelString), d.Value)
		}
	}
	return responseString
}
