package handlers

import (
	"net/http"
)

// SystemHandler provides health/liveness/readiness endpoints.
type SystemHandler struct {
	ready bool
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{ready: false}
}

// SetReady marks the service as ready.
func (h *SystemHandler) SetReady(v bool) {
	h.ready = v
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns OK if service is alive
// @Tags         system
// @Produce      plain
// @Success      200 {string} string "OK"
// @Router       /hc [get]
func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// LivenessCheck godoc
// @Summary      Liveness check
// @Tags         system
// @Produce      plain
// @Success      200 {string} string "ALIVE"
// @Router       /lc [get]
func (h *SystemHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ALIVE"))
}

// ReadinessCheck godoc
// @Summary      Readiness check
// @Tags         system
// @Produce      plain
// @Success      200 {string} string "READY"
// @Failure      503 {string} string "NOT READY"
// @Router       /rc [get]
func (h *SystemHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	if !h.ready {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("NOT READY"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}
