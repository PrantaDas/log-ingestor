package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBroker    string
	KafkaTopic     string
	KafkaGroupID   string
	MongoDBURI     string
	CollectionName string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func LoadConfig() *Config {
	return &Config{
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "error_logs"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "log_group"),
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://mongo:27017"),
		CollectionName: getEnv("COLLECTION_NAME", "errors"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
