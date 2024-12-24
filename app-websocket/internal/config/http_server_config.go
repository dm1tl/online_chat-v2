package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type HTTPServerConfig struct {
	Address        string        `yaml:"address"`
	MaxHeaderBytes int           `yaml:"maxheaderbytes"`
	ReadTimeout    time.Duration `yaml:"readtimeout"`
	WriteTimeout   time.Duration `yaml:"writetimeout"`
	IdleTimeout    time.Duration `yaml:"idletimeout"`
}

func NewHTTPServerConfig() (HTTPServerConfig, error) {
	var cfg HTTPServerConfig
	filepath := os.Getenv("HTTP_SERVER_CONFIG_PATH")
	if filepath == "" {
		return cfg, fmt.Errorf("HTTP_SERVER_CONFIG_PATH is empty")
	}
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		return cfg, fmt.Errorf("couldn't load config file #%v", err)
	}

	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return cfg, fmt.Errorf("couldn't parse config file into model #%v", err)
	}
	return cfg, nil
}
