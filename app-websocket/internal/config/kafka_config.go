package config

type KafkaConfig struct {
	BrokerList    []string `yaml:"brokers" env-required:"true"`
	Topic         string   `yaml:"topic" env-required:"true"`
	ConsumerGroup string   `yaml:"consumer_group" env-required:"true"`
}
