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
	DBName         string
	Port           string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
}

func LoadConfig() *Config {
	return &Config{
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:31001"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "error_logs"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "log_group"),
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		CollectionName: getEnv("COLLECTION_NAME", "errors"),
		DBName:         getEnv("DB_NAME", "log-ingester"),
		Port:           getEnv("PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
