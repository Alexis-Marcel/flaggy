package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) Stream(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
	w.WriteHeader(http.StatusOK)
	rc.Flush() // Force send headers immediately

	events, unsub := s.broadcaster.Subscribe()
	defer unsub()

	// Send initial connection event
	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"ok\"}\n\n")
	rc.Flush()

	// Keepalive ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-events:
			if !ok {
				return // Broadcaster closed
			}
			data, _ := json.Marshal(event.Data)
			if event.ID != "" {
				fmt.Fprintf(w, "id: %s\n", event.ID)
			}
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
			rc.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			rc.Flush()
		}
	}
}
