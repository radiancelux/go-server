package handlers

import (
	"go-server/internal/interfaces"
	"go-server/internal/models"
	"runtime"
	"time"
)

// StatusHandler handles status requests
type StatusHandler struct {
	logger interfaces.Logger
	port   string
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(logger interfaces.Logger, port string) *StatusHandler {
	return &StatusHandler{logger: logger, port: port}
}

// GetAction returns the action this handler processes
func (h *StatusHandler) GetAction() string {
	return "status"
}

// Handle processes the status request
func (h *StatusHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling status request")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	status := map[string]any{
		"status": "healthy",
		"server": map[string]any{
			"name":    "go-server",
			"version": "1.0.0",
			"port":    h.port,
			"uptime":  "running",
		},
		"system": map[string]any{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"goroutines": runtime.NumGoroutine(),
			"cpus":       runtime.NumCPU(),
		},
		"memory": map[string]any{
			"alloc_mb":       bToMb(m.Alloc),
			"total_alloc_mb": bToMb(m.TotalAlloc),
			"sys_mb":         bToMb(m.Sys),
			"num_gc":         m.NumGC,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	return models.NewSuccessResponse("Detailed server status", status), nil
}
