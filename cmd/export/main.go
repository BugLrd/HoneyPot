package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/database"
)

func main() {
	dbPath := "data/events.db"
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	events, err := db.AllEvents()
	if err != nil {
		log.Fatalf("all events: %v", err)
	}

	out, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		log.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile("data/exported_events.json", out, 0644); err != nil {
		log.Fatalf("write: %v", err)
	}
	log.Printf("exported %d events", len(events))
}
