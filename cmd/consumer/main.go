package main

import (
	"context"
	"log"
	"log-ingester/config"
	"log-ingester/internal/db"
	"log-ingester/internal/kafka"
	"log-ingester/pkg/utils"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config := config.LoadConfig()

	db, err := db.ConnectToMongoDB(ctx, config)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	consuemer, err := kafka.NewKafkaConsumer(config.KafkaBroker, config.KafkaGroupID, []string{config.KafkaTopic}, db)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}

	defer consuemer.Close()

	consuemer.StartConsuming(ctx, utils.MessagHandler)

	<-ctx.Done()
	log.Println("Shuttiing down gracefully...")
}
