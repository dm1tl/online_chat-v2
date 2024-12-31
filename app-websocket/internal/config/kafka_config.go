package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	BrokerList    []string `yaml:"brokers" env-required:"true"`
	Topic         string   `yaml:"topic" env-required:"true"`
	ConsumerGroup string   `yaml:"consumer_group" env-required:"true"`
}

func NewKafkaConfig() (*KafkaConfig, error) {
	var config KafkaConfig
	filepath := os.Getenv("KAFKA_CONFIG_PATH")
	if filepath == "" {
		return nil, fmt.Errorf("KAFKA_CONFIG_PATH is empty")
	}
	configFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("couldn't load config file %w", err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, fmt.Errorf("couldn't parse config file into model %w", err)
	}
	return &config, nil
}
