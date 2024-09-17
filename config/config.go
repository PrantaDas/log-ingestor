package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBroker    string
	KafkaTopic     string
	KafkaGroupID   string
	MongoDBURI     string
	CollectionName string
	DBName         string
	Port           int
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
}

func LoadConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Invalid PORT environment variable")
	}
	return &Config{
		KafkaBroker:    getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:     getEnv("KAFKA_TOPIC", "error_logs"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "log_group"),
		MongoDBURI:     getEnv("MONGODB_URI", "mongodb://mongo:27017"),
		CollectionName: getEnv("COLLECTION_NAME", "errors"),
		DBName:         getEnv("DB_NAME", "log-ingester"),
		Port:           port,
	}
}

func getEnv(key string, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
