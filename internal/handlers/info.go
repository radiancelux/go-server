package handlers

import (
	"github.com/radiancelux/go-server/internal/interfaces"
	"github.com/radiancelux/go-server/internal/models"
)

// InfoHandler handles info requests
type InfoHandler struct {
	logger interfaces.Logger
	port   string
}

// NewInfoHandler creates a new info handler
func NewInfoHandler(logger interfaces.Logger, port string) *InfoHandler {
	return &InfoHandler{logger: logger, port: port}
}

// GetAction returns the action this handler processes
func (h *InfoHandler) GetAction() string {
	return "info"
}

// Handle processes the info request
func (h *InfoHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling info request: %s", req.GetMessage())

	return models.NewSuccessResponse("Server information", map[string]any{
		"server":     "go-server",
		"version":    "1.0.0",
		"port":       h.port,
		"user_input": req.GetMessage(),
	}), nil
}
