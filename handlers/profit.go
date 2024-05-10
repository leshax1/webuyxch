package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"webuyxch/database"

	"go.mongodb.org/mongo-driver/mongo"
)

type ProfitHandler struct {
	DB *mongo.Client
}

func (h *ProfitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	profit, err := database.Profit(h.DB)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error calculating profit", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]int16{"profit": profit} // Create a map to be serialized

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}