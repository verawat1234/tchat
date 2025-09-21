package handlers

import (
	"encoding/json"
	"net/http"
)

type statusPayload struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Health responds to liveness probes.
func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, statusPayload{Status: "ok", Service: "authsvc"})
}

// Ready responds to readiness probes.
func Ready(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, statusPayload{Status: "ready", Service: "authsvc"})
}

func writeJSON(w http.ResponseWriter, payload statusPayload) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
