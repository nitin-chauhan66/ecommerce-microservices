package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

type Order struct {
	ID         string   `json:"id"`
	ProductIDs []string `json:"product_ids"`
	Total      float64  `json:"total"`
}

var producer sarama.SyncProducer

func main() {
	// Set up Kafka producer
	var err error
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	r := mux.NewRouter()
	r.HandleFunc("/orders", createOrder).Methods("POST")

	log.Println("Order Service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real scenario, you'd validate the order, check inventory, etc.

	// Send order to Kafka
	orderJSON, _ := json.Marshal(order)
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "new-orders",
		Value: sarama.StringEncoder(orderJSON),
	})
	if err != nil {
		http.Error(w, "Error publishing order to Kafka", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
