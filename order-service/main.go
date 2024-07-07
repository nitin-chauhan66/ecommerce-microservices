package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

type Order struct {
	ID         string   `json:"id"`
	ProductIDs []string `json:"product_ids"`
	Total      float64  `json:"total"`
}

var producer sarama.AsyncProducer

func main() {
	// Retry settings
	const maxRetries = 5
	const initialRetryInterval = 5 * time.Second

	// Create Kafka producer with retry
	err := createKafkaProducerWithRetry([]string{"kafka:9092"}, maxRetries, initialRetryInterval)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Start a goroutine to handle successes and errors
	go handleProducerEvents()

	r := mux.NewRouter()
	r.HandleFunc("/orders", createOrder).Methods("POST")

	log.Println("Order Service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func createKafkaProducerWithRetry(brokerList []string, maxRetries int, initialRetryInterval time.Duration) error {
	retries := 0
	retryInterval := initialRetryInterval

	for {
		var err error
		producer, err = sarama.NewAsyncProducer(brokerList, nil)
		if err == nil {
			return nil // Successfully created producer
		}

		retries++
		if retries > maxRetries {
			return err // Max retries exceeded
		}

		log.Printf("Error creating Kafka producer: %v. Retrying in %v...", err, retryInterval)
		time.Sleep(retryInterval)
		retryInterval += 5 * time.Second // increase intervla time by 5 secs
	}

}

func handleProducerEvents() {
	for {
		select {
		case success := <-producer.Successes():
			log.Printf("Message sent to partition %d at offset %d", success.Partition, success.Offset)
		case err := <-producer.Errors():
			log.Printf("Failed to send message: %v", err)
		}
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid order format", http.StatusBadRequest)
		return
	}

	// Basic validation
	if order.ID == "" || len(order.ProductIDs) == 0 || order.Total <= 0 {
		http.Error(w, "Invalid order data", http.StatusBadRequest)
		return
	}

	// Send order to Kafka
	orderJSON, _ := json.Marshal(order)
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "new-orders",
		Value: sarama.StringEncoder(orderJSON),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
