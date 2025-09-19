package handlers

import (
	"go-server/internal/interfaces"
	"go-server/internal/models"
)

// EchoHandler handles echo requests
type EchoHandler struct {
	logger interfaces.Logger
}

// NewEchoHandler creates a new echo handler
func NewEchoHandler(logger interfaces.Logger) *EchoHandler {
	return &EchoHandler{logger: logger}
}

// GetAction returns the action this handler processes
func (h *EchoHandler) GetAction() string {
	return "echo"
}

// Handle processes the echo request
func (h *EchoHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling echo request: %s", req.GetMessage())
	
	return models.NewSuccessResponse("Echo successful", map[string]string{
		"echoed_message": req.GetMessage(),
	}), nil
}
