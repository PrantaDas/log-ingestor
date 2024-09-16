package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

type Producer interface {
	PushLog(ctx context.Context, message map[string]interface{}) error
	Close() error
}

func NewKafkaProducer(broker string, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) PushLog(ctx context.Context, message map[string]interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)

	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	p.producer.Flush(15 * 1000)
	p.producer.Close()
	return nil
}
