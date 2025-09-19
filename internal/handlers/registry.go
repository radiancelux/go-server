package handlers

import "github.com/radiancelux/go-server/internal/interfaces"

// Registry manages handler registration and retrieval
type Registry struct {
	handlers map[string]interfaces.Handler
}

// NewRegistry creates a new handler registry
func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]interfaces.Handler),
	}
}

// Register adds a handler to the registry
func (r *Registry) Register(handler interfaces.Handler) {
	r.handlers[handler.GetAction()] = handler
}

// Get retrieves a handler by action
func (r *Registry) Get(action string) (interfaces.Handler, bool) {
	handler, exists := r.handlers[action]
	return handler, exists
}

// GetSupportedActions returns all supported actions
func (r *Registry) GetSupportedActions() []string {
	actions := make([]string, 0, len(r.handlers))
	for action := range r.handlers {
		actions = append(actions, action)
	}
	return actions
}
