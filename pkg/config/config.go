package config

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	pmConfig     *PlatformManagementConfig
	pmConfigPath string
)

// Prepare function used for setting the config path (after flags are parsed)
func Prepare(config_path string) {
	pmConfig = &PlatformManagementConfig{}

	pmConfigPath = os.Getenv("CONFIG_PATH")
	if pmConfigPath == "" {
		pmConfigPath = config_path
	}

	if _, err := os.Stat(pmConfigPath); os.IsNotExist(err) {
		log.Fatalf("Error loading config file: %v", err)
		os.Exit(1)
	}

	pmConfig.reload()
}

// reload Utility function for reloading the config file
func (c *PlatformManagementConfig) reload() {
	configContent, err := ioutil.ReadFile(pmConfigPath)
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configContent, c)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
		os.Exit(1)
	}
	log.Println("Config file reloaded")
}

// GetConfig Returns the current config instance
func GetConfig() *PlatformManagementConfig {
	return pmConfig
}

// GetManagementConfig Shortcut for getting node specific configs (if present)
func GetManagementConfig(host string) *NodeConfig {
	for i, _ := range pmConfig.PlatformManagements {
		if pmConfig.PlatformManagements[i].Host == host {
			return &pmConfig.PlatformManagements[i]
		}
	}
	return nil
}

// ReloadConfigHandler HTTP handler for /-/reload endpoint (if enabled)
func ReloadConfigHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" && req.Method != "POST" {
		_, _ = w.Write([]byte("Only POST or PUT requests allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	pmConfig.reload()
	_, _ = w.Write([]byte("OK"))
}
