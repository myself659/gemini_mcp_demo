package main

import (
	"log"

	"ip-store/backend/internal/api"
	"ip-store/backend/internal/database"
)

func main() {
	// Initialize the database
	database.InitDB("ip-store.db")

	r := api.SetupRouter(database.DB)
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
