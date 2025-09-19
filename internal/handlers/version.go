package handlers

import (
	"github.com/radiancelux/go-server/internal/interfaces"
	"github.com/radiancelux/go-server/internal/models"
	"runtime"
)

// VersionHandler handles version requests
type VersionHandler struct {
	logger interfaces.Logger
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(logger interfaces.Logger) *VersionHandler {
	return &VersionHandler{logger: logger}
}

// GetAction returns the action this handler processes
func (h *VersionHandler) GetAction() string {
	return "version"
}

// Handle processes the version request
func (h *VersionHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling version request")

	versionInfo := map[string]any{
		"server":     "go-server",
		"version":    "1.0.0",
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"num_cpu":    runtime.NumCPU(),
	}

	return models.NewSuccessResponse("Version information", versionInfo), nil
}
