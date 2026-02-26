package api

import (
	"encoding/json"
	"io"
	"net/http"
)

const maxBodySize = 1 << 20 // 1 MB

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodySize)
	defer io.Copy(io.Discard, r.Body)
	return json.NewDecoder(r.Body).Decode(dst)
}
