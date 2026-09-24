package database

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/event"
)

// CowrieEvent mirrors the raw JSON structure.
// It is defined locally to keep the test self‑contained.
// In production code this is in internal/collector.

type CowrieEvent struct {
	Session    string `json:"session"`
	Protocol   string `json:"protocol"`
	SrcIP      string `json:"src_ip"`
	SrcPort    int    `json:"src_port"`
	DstIP      string `json:"dst_ip"`
	DstPort    int    `json:"dst_port"`
	EventID    string `json:"eventid"`
	Sensor     string `json:"sensor"`
	UUID       string `json:"uuid"`
	Timestamp  string `json:"timestamp"`
	Message    string `json:"message"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Input      string `json:"input"`
	Version    string `json:"version"`
	HASSH      string `json:"hassh"`
	Arch       string `json:"arch"`
	DurationMS int64  `json:"duration_ms"`
}

func toEvent(raw CowrieEvent) (event.Event, error) {
	ts, err := time.Parse(time.RFC3339Nano, raw.Timestamp)
	if err != nil {
		return event.Event{}, err
	}
	return event.Event{
		ID:            0,
		Timestamp:     ts,
		SessionID:     raw.Session,
		Protocol:      raw.Protocol,
		SourceIP:      raw.SrcIP,
		SourcePort:    raw.SrcPort,
		DestIP:        raw.DstIP,
		DestPort:      raw.DstPort,
		EventID:       raw.EventID,
		Sensor:        raw.Sensor,
		Username:      raw.Username,
		Password:      raw.Password,
		Command:       raw.Input,
		ClientVersion: raw.Version,
		HASSH:         raw.HASSH,
		Architecture:  raw.Arch,
		DurationMS:    raw.DurationMS,
	}, nil
}

func TestImportAndPrintAllEvents(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	jsonPath := filepath.Join("..", "..", "data", "cowrie_events.json")
	f, err := os.Open(jsonPath)
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	for {
		var raw CowrieEvent
		if err := dec.Decode(&raw); err == io.EOF {
			break
		} else if err != nil {
			t.Fatalf("decode: %v", err)
		}
		ev, err := toEvent(raw)
		if err != nil {
			t.Fatalf("convert: %v", err)
		}
		if _, err := db.SaveEvent(ev); err != nil {
			t.Fatalf("save: %v", err)
		}
	}

	rows, err := db.db.Query("SELECT id, timestamp, session_id, protocol, source_ip, source_port, dest_ip, dest_port, event_id, sensor, username, password, command, client_version, hassh, architecture, duration FROM events")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var e event.Event
		var id int64
		var ts time.Time
		var sess, proto, srcIP, destIP, eventID, sensor, uname, pwd, cmd, ver, hassh, arch string
		var srcPort, destPort, dur int
		if err := rows.Scan(&id, &ts, &sess, &proto, &srcIP, &srcPort, &destIP, &destPort, &eventID, &sensor, &uname, &pwd, &cmd, &ver, &hassh, &arch, &dur); err != nil {
			t.Fatalf("scan: %v", err)
		}
		e.ID = id
		e.Timestamp = ts
		e.SessionID = sess
		e.Protocol = proto
		e.SourceIP = srcIP
		e.SourcePort = srcPort
		e.DestIP = destIP
		e.DestPort = destPort
		e.EventID = eventID
		e.Sensor = sensor
		e.Username = uname
		e.Password = pwd
		e.Command = cmd
		e.ClientVersion = ver
		e.HASSH = hassh
		e.Architecture = arch
		e.DurationMS = int64(dur)
		t.Logf("%+v", e)
		count++
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows error: %v", err)
	}
	t.Logf("total events imported: %d", count)
}

func TestAllEvents(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	ev := event.Event{
		Timestamp:     time.Now().UTC(),
		SessionID:     "session-1",
		Protocol:      "ssh",
		SourceIP:      "127.0.0.1",
		SourcePort:    2222,
		DestIP:        "0.0.0.0",
		DestPort:      22,
		EventID:       "cowrie.session.connect",
		Sensor:        "sensor-1",
		Username:      "alice",
		Password:      "secret",
		Command:       "ls -la",
		ClientVersion: "OpenSSH_9.0",
		HASSH:         "abc123",
		Architecture:  "x86_64",
		DurationMS:    123,
	}

	if _, err := db.SaveEvent(ev); err != nil {
		t.Fatalf("save event: %v", err)
	}

	events, err := db.AllEvents()
	if err != nil {
		t.Fatalf("all events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Username != ev.Username {
		t.Fatalf("expected username %q, got %q", ev.Username, events[0].Username)
	}
	if events[0].Command != ev.Command {
		t.Fatalf("expected command %q, got %q", ev.Command, events[0].Command)
	}
}
