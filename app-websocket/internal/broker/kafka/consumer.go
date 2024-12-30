package kafka

import (
	"app-websocket/internal/config"
	"app-websocket/internal/domain"
	"context"
	"encoding/json"
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

func (c *Consumer) Consume(ctx context.Context, handler domain.MessageHandler) error {
	res := make(chan error)
	defer close(res)

	go func() {
		for {
			if c.stop {
				break
			}
			kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsFatal() {
					res <- fmt.Errorf("critical Kafka error: %w", err)
					return
				}
				logrus.Error("Non-critical Kafka error", err)
				continue
			}
			if kafkaMsg == nil {
				continue
			}
			var domainMsg domain.Message
			if err := json.Unmarshal(kafkaMsg.Value, &domainMsg); err != nil {
				logrus.Error("Failed to parse message", err)
				continue
			}
			if err := handler(domainMsg); err != nil {
				logrus.Error("Failed to process message", err)
				continue
			}
			if _, err := c.consumer.CommitMessage(kafkaMsg); err != nil {
				logrus.Error("Failed to commit message", err)
				res <- fmt.Errorf("failed to commit message: %w", err)
				return
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-res:
			return err
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	return c.consumer.Close()
}
