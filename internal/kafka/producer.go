package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

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
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		return nil, err
	}
	defer a.Close()

	durarion, err := time.ParseDuration("60s")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.CreateTopics(ctx, []kafka.TopicSpecification{{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}},
		kafka.SetAdminOperationTimeout(durarion))

	if err != nil {
		return nil, err
	}

	log.Println("Successfully created topic", res)

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
	p.producer.Flush(1 * 1000)
	p.producer.Close()
	return nil
}
