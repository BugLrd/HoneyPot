package database

import (
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/event"
)

// Database is a wrapper around *sql.DB that handles the event schema.
// It depends only on the shared event type.

type Database struct {
	db *sql.DB
}

type Stats struct {
	TotalEvents      int64 `json:"total_events"`
	TotalSessions    int64 `json:"total_sessions"`
	UniqueIPs        int64 `json:"unique_ips"`
	SuccessfulLogins int64 `json:"successful_logins"`
	FailedLogins     int64 `json:"failed_logins"`
	TotalCommands    int64 `json:"total_commands"`
	UniqueCommands   int64 `json:"unique_commands"`
}

func New(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}

	if err := d.createTables(); err != nil {
		db.Close()
		return nil, err
	}

	return d, nil
}

func (d *Database) Ping() error {
	return d.db.Ping()
}

func (d *Database) createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	timestamp DATETIME NOT NULL,
    	session_id TEXT,
    	protocol TEXT,
    	source_ip TEXT,
    	source_port INTEGER,
    	dest_ip TEXT,
    	dest_port INTEGER,
    	event_id TEXT,
    	sensor TEXT,
    	username TEXT,
    	password TEXT,
    	command TEXT,
    	client_version TEXT,
    	hassh TEXT,
    	architecture TEXT,
    	duration INTEGER
	);
	CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
	CREATE INDEX IF NOT EXISTS idx_events_session_id ON events(session_id);
	CREATE INDEX IF NOT EXISTS idx_events_event_id ON events(event_id);
	`
	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) SaveEvent(event event.Event) (int64, error) {
	result, err := d.db.Exec(`
		INSERT INTO events (
			timestamp,
			session_id,
			protocol,
			source_ip,
			source_port,
			dest_ip,
			dest_port,
			event_id,
			sensor,
			username,
			password,
			command,
			client_version,
			hassh,
			architecture,
			duration
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		event.Timestamp,
		event.SessionID,
		event.Protocol,
		event.SourceIP,
		event.SourcePort,
		event.DestIP,
		event.DestPort,
		event.EventID,
		event.Sensor,
		event.Username,
		event.Password,
		event.Command,
		event.ClientVersion,
		event.HASSH,
		event.Architecture,
		event.DurationMS,
	)

	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
func (d *Database) GetEvent(id string) (*event.Event, error) {
	event := &event.Event{}

	err := d.db.QueryRow(`
        SELECT
            id,
            timestamp,
            session_id,
            protocol,
            source_ip,
            source_port,
            dest_ip,
            dest_port,
            event_id,
            sensor,
            username,
            password,
            command,
            client_version,
            hassh,
            architecture,
            duration
        FROM events
        WHERE id = ?
    `, id).Scan(
		&event.ID,
		&event.Timestamp,
		&event.SessionID,
		&event.Protocol,
		&event.SourceIP,
		&event.SourcePort,
		&event.DestIP,
		&event.DestPort,
		&event.EventID,
		&event.Sensor,
		&event.Username,
		&event.Password,
		&event.Command,
		&event.ClientVersion,
		&event.HASSH,
		&event.Architecture,
		&event.DurationMS,
	)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (d *Database) AllEvents() ([]event.Event, error) {
	rows, err := d.db.Query(`
        SELECT
            id,
            timestamp,
            session_id,
            protocol,
            source_ip,
            source_port,
            dest_ip,
            dest_port,
            event_id,
            sensor,
            username,
            password,
            command,
            client_version,
            hassh,
            architecture,
            duration
        FROM events
        ORDER BY id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(
			&e.ID,
			&e.Timestamp,
			&e.SessionID,
			&e.Protocol,
			&e.SourceIP,
			&e.SourcePort,
			&e.DestIP,
			&e.DestPort,
			&e.EventID,
			&e.Sensor,
			&e.Username,
			&e.Password,
			&e.Command,
			&e.ClientVersion,
			&e.HASSH,
			&e.Architecture,
			&e.DurationMS,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (d *Database) Stats() (Stats, error) {
	var stats Stats

	err := d.db.QueryRow(`
		SELECT
			COUNT(*) AS total_events,

			COUNT(DISTINCT session_id) AS total_sessions,

			COUNT(DISTINCT source_ip) AS unique_ips,

			COUNT(*) FILTER (
				WHERE event_id = 'cowrie.login.success'
			) AS successful_logins,

			COUNT(*) FILTER (
				WHERE event_id = 'cowrie.login.failed'
			) AS failed_logins,

			COUNT(*) FILTER (
				WHERE event_id = 'cowrie.command.input'
			) AS total_commands,

			COUNT(DISTINCT command) FILTER (
				WHERE event_id = 'cowrie.command.input'
				AND command IS NOT NULL
				AND command != ''
			) AS unique_commands

		FROM events
	`).Scan(
		&stats.TotalEvents,
		&stats.TotalSessions,
		&stats.UniqueIPs,
		&stats.SuccessfulLogins,
		&stats.FailedLogins,
		&stats.TotalCommands,
		&stats.UniqueCommands,
	)

	if err != nil {
		return Stats{}, err
	}

	return stats, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
