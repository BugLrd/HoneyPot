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
func (d *Database) GetEvent(id int64) (*event.Event, error) {
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

func (d *Database) Close() error {
	return d.db.Close()
}
