package handlers

import (
	"fmt"
	"github.com/radiancelux/go-server/internal/interfaces"
	"github.com/radiancelux/go-server/internal/models"
)

// GreetHandler handles greeting requests
type GreetHandler struct {
	logger interfaces.Logger
}

// NewGreetHandler creates a new greet handler
func NewGreetHandler(logger interfaces.Logger) *GreetHandler {
	return &GreetHandler{logger: logger}
}

// GetAction returns the action this handler processes
func (h *GreetHandler) GetAction() string {
	return "greet"
}

// Handle processes the greet request
func (h *GreetHandler) Handle(req interfaces.APIRequest) (interfaces.APIResponse, error) {
	h.logger.Debug("Handling greet request from user %d: %s", req.GetUserID(), req.GetMessage())

	greeting := fmt.Sprintf("Hello! You said: %s", req.GetMessage())
	if req.GetUserID() > 0 {
		greeting = fmt.Sprintf("Hello User %d! You said: %s", req.GetUserID(), req.GetMessage())
	}

	return models.NewSuccessResponse("Greeting generated", map[string]string{
		"greeting": greeting,
	}), nil
}
