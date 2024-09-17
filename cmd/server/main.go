package main

import (
	"context"
	"fmt"
	"log"
	"log-ingester/config"
	"log-ingester/internal/kafka"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type contextKey string

type Request struct {
	mu      sync.Mutex
	Options *RequestOption
}

type RequestOption struct {
	KafkaProducer *kafka.KafkaProducer
	Err           error
	RequestId     string
}

const (
	optionKey contextKey = "optionKey"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config := config.LoadConfig()

	kafkaProducer, err := kafka.NewKafkaProducer(config.KafkaBroker, config.KafkaTopic)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	defer kafkaProducer.Close()

	router := mux.NewRouter()

	router.HandleFunc("/simulate-error", func(w http.ResponseWriter, r *http.Request) {
		err := simulatError()
		reqOption := &Request{
			Options: &RequestOption{
				KafkaProducer: kafkaProducer,
			},
		}

		if err != nil {
			reqOption.mu.Lock()
			defer reqOption.mu.Unlock()
			reqOption.Options.Err = err
			reqOption.Options.RequestId = uuid.New().String()
			ctx := context.WithValue(r.Context(), optionKey, reqOption)
			errHandler(ctx)
			http.Error(w, "An error occurred while processing request", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sucess"))
	})

	log.Printf("Starting HTTP Server on port %s...\n", config.Port)
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("127.0.0.1:%s", config.Port),
	}

	log.Fatal(server.ListenAndServe())

	<-ctx.Done()
	log.Println("Shutting down HTTP Server gracefully...")
}

func errHandler(ctx context.Context) {
	reqOption, ok := ctx.Value(optionKey).(*Request)
	if !ok {
		return
	}
	reqOption.mu.Lock()
	defer reqOption.mu.Unlock()
	if reqOption.Options.Err != nil {
		kafkaProducer := reqOption.Options.KafkaProducer
		message := map[string]interface{}{
			"requestId":    reqOption.Options.RequestId,
			"error":        reqOption.Options.Err.Error(),
			"source":       "HTTP Server",
			"additionInfo": "Dummy Additional Info",
		}
		err := kafkaProducer.PushLog(ctx, message)
		if err != nil {
			log.Printf("Error pushing error log to Kafka: %v", err)
		}
	}
}

func simulatError() error {
	return fmt.Errorf("this is a simulated error")
}
