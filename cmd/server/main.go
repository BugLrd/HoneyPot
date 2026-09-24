package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/collector"
	"github.com/BugLrd/HoneyPot-DashBoard/internal/database"
)

func main() {
	dbPath := "data/events.db"
	jsonPath := "data/cowrie_events.json"

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
	count, err := r.Run(false)
	if err != nil {
		log.Fatalf("error reading events: %v", err)
	}

	log.Printf("%d events imported successfully", count)
}
