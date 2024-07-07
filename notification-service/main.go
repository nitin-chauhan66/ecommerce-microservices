package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Order struct {
	ID         string   `json:"id"`
	ProductIDs []string `json:"product_ids"`
	Total      float64  `json:"total"`
}

const (
	maxRetries           = 5
	initialRetryInterval = 5 * time.Second
)

func main() {
	consumer, err := createConsumerWithRetry([]string{"kafka:9092"}, maxRetries, initialRetryInterval)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to the "new-orders" topic
	partitionConsumer, err := consumer.ConsumePartition("new-orders", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	log.Println("Notification Service is running and listening for new orders")

	// Consume messages
	for msg := range partitionConsumer.Messages() {
		var order Order
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}

		// In a real scenario, you'd send an actual notification here
		log.Printf("Sending notification for order %s, total: $%.2f", order.ID, order.Total)
	}
}

func createConsumerWithRetry(brokerList []string, maxRetries int, initialRetryInterval time.Duration) (sarama.Consumer, error) {
	retries := 0
	retryInterval := initialRetryInterval

	for {
		var err error
		consumer, err := sarama.NewConsumer(brokerList, nil)
		if err == nil {
			return consumer, nil // Successfully created consumer
		}

		retries++
		if retries > maxRetries {
			return nil, err // Max retries exceeded
		}

		log.Printf("Error creating Kafka consumer: %v. Retrying in %v...", err, retryInterval)
		time.Sleep(retryInterval)
		retryInterval += 5 * time.Second // increase intervla time by 5 secs
	}
}
