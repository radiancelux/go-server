package handlers

import (
	"github.com/radiancelux/go-server/internal/interfaces"
	"github.com/radiancelux/go-server/internal/models"
	"os"
)

// ConfigHandler handles configuration requests
type ConfigHandler struct {
	logger interfaces.Logger
	port   string
}

// NewConfigHandler creates a new config handler
func NewConfigHandler(logger interfaces.Logger, port string) *ConfigHandler {
	return &ConfigHandler{logger: logger, port: port}
}

// GetAction returns the action this handler processes
func (h *ConfigHandler) GetAction() string {
	return "config"
}

// Handle processes the config request
func (h *ConfigHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling config request")

	config := map[string]any{
		"server": map[string]any{
			"port": h.port,
			"host": "localhost",
		},
		"environment": map[string]any{
			"go_env":    os.Getenv("GO_ENV"),
			"port":      os.Getenv("PORT"),
			"log_level": os.Getenv("LOG_LEVEL"),
		},
		"features": map[string]any{
			"graceful_shutdown":  true,
			"request_validation": true,
			"structured_logging": true,
			"metrics":            true,
		},
	}

	return models.NewSuccessResponse("Server configuration", config), nil
}
