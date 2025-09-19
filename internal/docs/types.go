package docs

// PostmanCollection represents a Postman collection
type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanInfo represents collection info
type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PostmanItem represents a collection item (folder or request)
type PostmanItem struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Item        []PostmanItem     `json:"item,omitempty"`
	Request     *PostmanRequest   `json:"request,omitempty"`
	Response    []PostmanResponse `json:"response,omitempty"`
	Event       []PostmanEvent    `json:"event,omitempty"`
}

// PostmanRequest represents a Postman request
type PostmanRequest struct {
	Method string              `json:"method"`
	Header []PostmanHeader     `json:"header"`
	Body   *PostmanRequestBody `json:"body"`
	URL    *PostmanURL         `json:"url"`
}

// PostmanHeader represents a request header
type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostmanRequestBody represents request body
type PostmanRequestBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

// PostmanURL represents request URL
type PostmanURL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
	Query    []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"query"`
}

// PostmanResponse represents a response example
type PostmanResponse struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Header  []PostmanHeader `json:"header"`
	Body    string `json:"body"`
}

// PostmanEvent represents a pre/post request script
type PostmanEvent struct {
	Listen string         `json:"listen"`
	Script PostmanScript  `json:"script"`
}

// PostmanScript represents a script
type PostmanScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// EndpointDocumentation represents documentation for an endpoint
type EndpointDocumentation struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Method      string                   `json:"method"`
	URL         string                   `json:"url"`
	Headers     []HeaderDocumentation    `json:"headers"`
	Body        *BodyDocumentation       `json:"body,omitempty"`
	Responses   []ResponseDocumentation  `json:"responses"`
}

// RequestDocumentation represents request documentation
type RequestDocumentation struct {
	Method  string                   `json:"method"`
	URL     string                   `json:"url"`
	Headers []HeaderDocumentation    `json:"headers"`
	Body    *BodyDocumentation       `json:"body,omitempty"`
}

// HeaderDocumentation represents header documentation
type HeaderDocumentation struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// BodyDocumentation represents body documentation
type BodyDocumentation struct {
	Type        string `json:"type"`
	Content     string `json:"content"`
	Description string `json:"description"`
}

// ResponseDocumentation represents response documentation
type ResponseDocumentation struct {
	Code        int                    `json:"code"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Headers     []HeaderDocumentation  `json:"headers"`
	Body        string                 `json:"body"`
}

// APIDocumentation represents the complete API documentation
type APIDocumentation struct {
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	BaseURL     string                     `json:"base_url"`
	Endpoints   []EndpointDocumentation    `json:"endpoints"`
	GeneratedAt string                     `json:"generated_at"`
}
