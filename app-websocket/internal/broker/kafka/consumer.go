package kafka

import (
	"app-websocket/internal/config"
	"app-websocket/internal/domain"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

const (
	consumerGroup  = "my-group"
	sessionTimeout = 500
	noTimeout      = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	topic    string
	handler  domain.MessageHandler
	stop     bool
}

func NewConsumer(cfg config.KafkaConfig) (*Consumer, error) {
	op := "kafka.NewProducer"
	conf := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(cfg.BrokerList, ","),
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       false,
	}
	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("op: %s: %w", op, err)
	}
	if err := c.Subscribe(cfg.Topic, nil); err != nil {
		return nil, fmt.Errorf("op: %s: %w", op, err)
	}
	return &Consumer{
		consumer: c,
		topic:    cfg.Topic,
	}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			logrus.Error(err)
		}
		if kafkaMsg == nil {
			continue
		}
		msg, err := domain.NewMessage(kafkaMsg.Value)
		if err != nil {
			logrus.Error(err)
			continue
		}
		if err := c.handler(*msg); err != nil {
			logrus.Error("failed to process message", err)
			continue
		}
		if _, err := c.consumer.CommitMessage(kafkaMsg); err != nil {
			logrus.Error("failed to commit message", err)
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	return c.consumer.Close()
}
