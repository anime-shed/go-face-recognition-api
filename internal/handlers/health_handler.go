package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	logger    *logrus.Logger
	startTime time.Time
}

// NewHealthHandler creates a new health handler instance
func NewHealthHandler(logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		logger:    logger,
		startTime: time.Now(),
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Version   string    `json:"version"`
	Service   string    `json:"service"`
}

// HealthHandler handles GET /api/v1/health endpoint
func (h *HealthHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "1.0.0",
		Service:   "face-recognition-api",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ReadinessHandler handles GET /api/v1/ready endpoint
func (h *HealthHandler) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	// In a real application, you would check dependencies here
	// For now, we'll just return ready
	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now(),
		"checks": map[string]string{
			"pigo":      "ok",
			"memory":    "ok",
			"disk":      "ok",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// LivenessHandler handles GET /api/v1/live endpoint
func (h *HealthHandler) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}