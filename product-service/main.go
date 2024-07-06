package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var products []Product

func main() {
	// Initialize products (in a real scenario, this would come from a database)
	products = []Product{
		{ID: "1", Name: "Laptop", Price: 999.99},
		{ID: "2", Name: "Smartphone", Price: 499.99},
	}

	r := mux.NewRouter()
	r.HandleFunc("/products", getProducts).Methods("GET")

	log.Println("Product Service is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
