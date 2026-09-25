package api

import (
	"net/http"

	"github.com/BugLrd/HoneyPot-DashBoard/internal/database"
)

type Server struct {
	db *database.Database
}

func NewServer(db *database.Database) *Server {
	return &Server{db: db}
}

func (s *Server) Start(addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: s.Handler(),
	}
	return server.ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.handleHealth)
	mux.HandleFunc("/events", s.handleEvents)
	mux.HandleFunc("/event", s.handleEvent)
	mux.HandleFunc("/event/{value}", s.handleEvent)
	mux.HandleFunc("/stats", s.handleStats)
	mux.HandleFunc("/ip", s.handleIpStats)
	mux.HandleFunc("/ip/{ip}", s.handleIpStats)
	mux.HandleFunc("/command", s.handleCommandStats)
	mux.HandleFunc("/command/{command}", s.handleCommandStats)

	return mux
}
