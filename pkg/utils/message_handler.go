package utils

import (
	"context"
	"fmt"
	"log-ingester/internal/db"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func MessagHandler(ctx context.Context, msg *kafka.Message, db *db.MongoDB) {
	message := msg.Value
	fmt.Println(message)
}
