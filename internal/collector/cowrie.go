package collector

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/event"
)

type CowrieEvent struct {
	Session    string `json:"session"`
	Protocol   string `json:"protocol"`
	SrcIP      string `json:"src_ip"`
	SrcPort    int    `json:"src_port"`
	DstIP      string `json:"dst_ip"`
	DstPort    int    `json:"dst_port"`
	EventID    string `json:"eventid"`
	Sensor     string `json:"sensor"`
	Timestamp  string `json:"timestamp"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Input      string `json:"input"`
	Version    string `json:"version"`
	HASSH      string `json:"hassh"`
	Arch       string `json:"arch"`
	DurationMS int64  `json:"duration_ms"`
}

func ParseCowrieEvent(raw []byte) (event.Event, error) {
	var cowrie CowrieEvent
	if err := json.Unmarshal(raw, &cowrie); err != nil {
		return event.Event{}, err
	}

	var ts time.Time
	var err error
	if strings.TrimSpace(cowrie.Timestamp) != "" {
		ts, err = time.Parse(time.RFC3339Nano, cowrie.Timestamp)
		if err != nil {
			ts = time.Now().UTC()
		}
	} else {
		ts = time.Now().UTC()
	}

	return event.Event{
		Timestamp:     ts,
		SessionID:     cowrie.Session,
		Protocol:      cowrie.Protocol,
		SourceIP:      cowrie.SrcIP,
		SourcePort:    cowrie.SrcPort,
		DestIP:        cowrie.DstIP,
		DestPort:      cowrie.DstPort,
		EventID:       cowrie.EventID,
		Sensor:        cowrie.Sensor,
		Username:      cowrie.Username,
		Password:      cowrie.Password,
		Command:       cowrie.Input,
		ClientVersion: cowrie.Version,
		HASSH:         cowrie.HASSH,
		Architecture:  cowrie.Arch,
		DurationMS:    cowrie.DurationMS,
	}, nil
}
