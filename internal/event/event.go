package event

import "time"

// Event represents a parsed event from Cowrie logs.
// It is shared between the collector and database packages.
// Keep the fields identical to the original collector.Event definition.
// If new fields are added, update both packages accordingly.

// Note: this file intentionally contains only the struct definition to avoid
// import cycles.

type Event struct {
	ID        int64
	Timestamp time.Time

	SessionID string
	Protocol  string

	SourceIP   string
	SourcePort int
	DestIP     string
	DestPort   int

	EventID string
	Sensor  string

	Username string
	Password string

	Command string

	ClientVersion string
	HASSH         string

	Architecture string

	DurationMS int64
}
