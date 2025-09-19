package handlers

import (
	"go-server/internal/interfaces"
	"go-server/internal/models"
	"runtime"
	"time"
)

// MetricsHandler handles metrics requests
type MetricsHandler struct {
	logger interfaces.Logger
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(logger interfaces.Logger) *MetricsHandler {
	return &MetricsHandler{logger: logger}
}

// GetAction returns the action this handler processes
func (h *MetricsHandler) GetAction() string {
	return "metrics"
}

// Handle processes the metrics request
func (h *MetricsHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling metrics request")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]any{
		"memory": map[string]any{
			"alloc_mb":        bToMb(m.Alloc),
			"total_alloc_mb":  bToMb(m.TotalAlloc),
			"sys_mb":          bToMb(m.Sys),
			"num_gc":          m.NumGC,
			"gc_cpu_fraction": m.GCCPUFraction,
		},
		"runtime": map[string]any{
			"goroutines": runtime.NumGoroutine(),
			"cpus":       runtime.NumCPU(),
		},
		"timestamp": time.Now().Unix(),
	}

	return models.NewSuccessResponse("System metrics", metrics), nil
}
