package kafka

import (
	"app-websocket/internal/config"
	"app-websocket/internal/domain"
	byteencoding "app-websocket/pkg/byte_encoding"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

var errUnknownType = errors.New("unknown event type")

const (
	flushTimeout = 5000
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	op := "kafka.NewProducer"
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.BrokerList, ","),
	}
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("op: %s: %w", op, err)
	}
	return &Producer{
		producer: p,
	}, nil
}

func (p *Producer) Produce(msg domain.Message, topic string, key int64) error {
	op := "kafka.Produce"
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("op: %s: %w", op, err)
	}
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: jsonMsg,
		Key:   byteencoding.Int64ToBytes(key),
	}
	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return err
	}
	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case *kafka.Error:
		return ev
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
