package main

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
)

type Order struct {
	ID         string   `json:"id"`
	ProductIDs []string `json:"product_ids"`
	Total      float64  `json:"total"`
}

func main() {
	// Set up Kafka consumer
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
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
