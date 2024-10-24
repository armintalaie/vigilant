package main

import (
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vigilant/internal/events"
	"vigilant/internal/grpc"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	kProducer := internal.NewKProducer("logs", producer)
	go func() {
		address := ":50051"
		log.Printf("Starting gRPC server on %s", address)
		if err := grpc.StartServer(address, kProducer); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	consumer, err := internal.NewConsumer([]string{"localhost:9092"}, "logs", "cmd/desktop/logs.db")
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	go func() {
		if err := consumer.Start(); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	select {}
}
