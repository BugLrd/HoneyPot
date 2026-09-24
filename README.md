# HoneyPot‑Dashboard

> A lightweight Go application that ingests Cowrie honeypot logs into a SQLite database and exposes a REST API for analysis.

## 📋 Overview

The project is split into three logical layers:

| Layer | Responsibility | Location |
|-------|----------------|----------|
| **Collector** | Parses raw Cowrie JSON events, converts them to the internal `event.Event` struct, and stores them in SQLite. | `internal/collector` |
| **Database** | A thin wrapper around `database/sql` that creates the `events` table, provides CRUD helpers and a statistics endpoint. | `internal/database` |
| **API** | Serves HTTP endpoints for health checks, event lists, individual events and aggregated statistics. | `internal/api` |

The **command‑line tools** in `cmd/` bootstrap the collector and the web server.

## ⚙️ Prerequisites

* Go 1.22 or newer
* Cowrie honeypot installed and a JSON log file (default path: `/home/kai/Projects/cowrie/var/log/cowrie/cowrie.json` – you can edit `cmd/server/main.go` to point to your own log file.)
* An empty folder `data/` for the SQLite database (the application will create `data/events.db` automatically).

## 🚀 Getting Started

```bash
# Run the collector and API server in one process
go run ./cmd/server
```

> **Tip:** The server listens on `:8080` by default. You can change the port by editing `cmd/server/main.go` or adding an `-addr` flag in a future release.

### Importing Events

The server automatically tails the Cowrie JSON file and imports every new event. To run the import process separately:

```bash
go run ./cmd/export
```

This reads all events from the database and writes them to `data/exported_events.json`.

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | `GET` | Health check – returns `{"status":"ok"}` when the DB is reachable. |
| `/events` | `GET` | Returns a JSON array of all events. |
| `/event` | `GET` | Query param `id` – returns a single event by its numeric ID. |
| `/event/{id}` | `GET` | Path param `id` – same as above. |
| `/stats` | `GET` | Returns aggregated statistics (`total_events`, `unique_ips`, etc.). |

#### Example

```bash
curl http://localhost:8080/events | jq '.[0:5]'
```

### Testing

```bash
go test ./...
```

All unit tests live under `internal/database/database_test.go` and cover database creation, CRUD and statistics queries.

## 📦 Project Structure

```
HoneyPot‑Dashboard/
├── cmd/
│   ├── export/main.go   # Exports all events to JSON
│   └── server/main.go   # Starts collector + HTTP API
├── data/
│   └── events.db        # SQLite database (auto‑created)
├── internal/
│   ├── api/          # HTTP handlers
│   ├── collector/    # Log parsing & ingestion
│   ├── database/     # DB wrapper & schema
│   └── event/        # Shared event struct
└── README.md
```

---

> For any questions, open an issue or reach out to the maintainer.
