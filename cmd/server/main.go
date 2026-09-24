package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/api"
	"github.com/BugLrd/HoneyPot-DashBoard/internal/collector"
	"github.com/BugLrd/HoneyPot-DashBoard/internal/database"
)

func main() {
	dbPath := "data/events.db"
	jsonPath := "/home/kai/Projects/cowrie/var/log/cowrie/cowrie.json"

	// Ensure the parent directory exists before sqlite tries to write to it
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		log.Fatalf("json file %s does not exist", jsonPath)
	}

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	r := collector.NewReader(jsonPath, db)
	go func() {
		count, err := r.Run(true)
		if err != nil {
			log.Fatalf("error reading events: %v", err)
		}
		log.Printf("%d events imported successfully", count)
	}()

	server := api.NewServer(db)
	if err := server.Start(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
