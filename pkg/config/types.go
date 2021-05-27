package config

type NodeConfig struct {
	Host     string            `yaml:"host"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
	Type     string            `yaml:"type"`
	Labels   map[string]string `yaml:"labels"`
}

type PlatformManagementConfig struct {
	PlatformManagements []NodeConfig `yaml:"platform_managements"`
}
