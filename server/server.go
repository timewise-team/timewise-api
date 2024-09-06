package server

import (
	"api/config"
	"api/feature"
	"log"
)

func RegisterServer() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	// Initialize router
	r := feature.RegisterHandlerV1()
	// Start server
	log.Printf("Server is running on port %s", cfg.ServerPort)
	if err := r.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
