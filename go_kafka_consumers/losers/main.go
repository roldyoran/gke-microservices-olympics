package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type StudentMessage struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Faculty    string `json:"faculty"`
	Discipline int    `json:"discipline"`
}

func consumeMessages(brokers []string, topic string, stopChan chan struct{}, rdb *redis.Client) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "student-consumer-group",
		Topic:    topic,
		MinBytes: 10e3,            // 10KB
		MaxBytes: 10e6,            // 10MB
		MaxWait:  1 * time.Second, // Maximum wait of 1 second per message
	})

	defer reader.Close()

	fmt.Printf("Consuming messages from topic: %s\n", topic)
	for {
		select {
		case <-stopChan:
			fmt.Printf("Stopping message consumption from topic: %s\n", topic)
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				// Handle timeout error generally
				if errors.Is(err, context.DeadlineExceeded) {
					// No more messages available, wait a moment
					time.Sleep(2 * time.Second)
					continue
				}
				log.Printf("Unexpected error receiving message: %v\n", err)
				continue
			}

			fmt.Printf("Message received on topic %s: %s\n", msg.Topic, string(msg.Value))

			// Process the message and store it in Redis
			var student StudentMessage
			// Simulate reading the message as a JSON string
			if err := json.Unmarshal([]byte(msg.Value), &student); err != nil {
				log.Printf("Error decoding message: %v\n", err)
				continue
			}

			// Update counts in Redis
			if student.Status == "lost" {
				rdb.Incr(ctx, fmt.Sprintf("faculty:%s:losers", student.Faculty))
			}
			// add total students per faculty regardless of whether they won or lost
			rdb.Incr(ctx, fmt.Sprintf("faculty:%s:total", student.Faculty))
		}
	}
}

func main() {
	brokers := []string{"my-cluster-kafka-bootstrap:9092"}
	stopChan := make(chan struct{})

	// Redis configuration
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-service:6379", // Redis service name in Kubernetes
	})

	go consumeMessages(brokers, "student-losers", stopChan, rdb)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	fmt.Println("Interrupt signal received, shutting down consumer...")
	close(stopChan)
	time.Sleep(2 * time.Second)
}
