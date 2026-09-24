package collector

import (
	"bufio"
	"io"
	"log"
	"os"
	"time"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/database"
)

type Reader struct {
	path string
	db   *database.Database
}

func NewReader(path string, db *database.Database) *Reader {
	return &Reader{
		path: path,
		db:   db,
	}
}

// Run processes the file. Set follow to true if you want to tail the log continuously.
func (r *Reader) Run(follow bool) (int, error) {
	count := 0

	file, err := os.Open(r.path)
	if err != nil {
		return count, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			r.processLine(line)
			count++
		}

		if err != nil {
			if err == io.EOF {
				if !follow {
					break
				}
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return count, err
		}
	}
	return count, nil
}

func (r *Reader) processLine(line []byte) {
	event, err := ParseCowrieEvent(line)
	if err != nil {
		log.Printf("[PARSE ERROR] %v | line: %s", err, string(line))
		return
	}

	id, err := r.db.SaveEvent(event)
	if err != nil {
		log.Printf("[DB INSERT ERROR] %v | event: %+v", err, event)
		return
	}

	log.Printf("[SUCCESS] Inserted row id: %d for event: %s", id, event.EventID)
}
