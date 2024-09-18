package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log-ingester/internal/db"
	"log-ingester/pkg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func MessagHandler(ctx context.Context, msg *kafka.Message, db *db.MongoDB) {
	message := msg.Value
	var res pkg.Message
	err := json.Unmarshal(message, &res)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}
	err = db.InsertLogToDB(ctx, res)
	if err != nil {
		log.Printf("Error inserting log to MongoDB: %v", err)
		return
	}
	fmt.Println("Successfully inserted log into MongoDB")
}
