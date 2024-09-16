package kafka

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

type Consumer interface {
	StartConsuming(ctx context.Context, msgHandler func(msg *kafka.Message))
	Close() error
}

func NewKafkaConsumer(broker string, topic string, groupID string, topics []string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: c,
		topic:    topic,
	}, nil
}

func (k *KafkaConsumer) StartConsuming(ctx context.Context, msgHandler func(msg *kafka.Message)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := k.consumer.ReadMessage(time.Second)
				if err == nil {
					msgHandler(msg)
				} else {
					log.Printf("Consumer error: %v", err)
				}
			}
		}
	}()
}

func (k *KafkaConsumer) Close() error {
	k.consumer.Close()
	return nil
}
