package docs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

// PostmanDocGenerator generates HTML documentation from Postman collections
type PostmanDocGenerator struct {
	baseURL string
}

// NewPostmanDocGenerator creates a new PostmanDocGenerator
func NewPostmanDocGenerator(baseURL string) *PostmanDocGenerator {
	return &PostmanDocGenerator{
		baseURL: baseURL,
	}
}

// PostmanCollection represents a Postman collection
type PostmanCollection struct {
	Info PostmanInfo `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanInfo represents collection info
type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PostmanItem represents a collection item (folder or request)
type PostmanItem struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Item        []PostmanItem  `json:"item,omitempty"`
	Request     *PostmanRequest `json:"request,omitempty"`
	Response    []PostmanResponse `json:"response,omitempty"`
	Event       []PostmanEvent `json:"event,omitempty"`
}

// PostmanRequest represents a Postman request
type PostmanRequest struct {
	Method string                 `json:"method"`
	Header []PostmanHeader        `json:"header"`
	Body   *PostmanRequestBody    `json:"body"`
	URL    PostmanURL             `json:"url"`
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
}

// PostmanResponse represents a Postman response
type PostmanResponse struct {
	Name   string                 `json:"name"`
	Status string                 `json:"status"`
	Header []PostmanHeader        `json:"header"`
	Body   string                 `json:"body"`
}

// PostmanEvent represents a Postman event (test scripts)
type PostmanEvent struct {
	Listen string `json:"listen"`
	Script PostmanScript `json:"script"`
}

// PostmanScript represents a Postman script
type PostmanScript struct {
	Exec []string `json:"exec"`
}

// EndpointDocumentation represents documentation for a single endpoint
type EndpointDocumentation struct {
	Name        string
	Method      string
	Path        string
	Description string
	Category    string
	Request     *RequestDocumentation
	Response    *ResponseDocumentation
	Tests       []string
}

// RequestDocumentation represents request documentation
type RequestDocumentation struct {
	Headers []HeaderDocumentation
	Body    *BodyDocumentation
}

// HeaderDocumentation represents header documentation
type HeaderDocumentation struct {
	Name  string
	Value string
}

// BodyDocumentation represents body documentation
type BodyDocumentation struct {
	Example string
}

// ResponseDocumentation represents response documentation
type ResponseDocumentation struct {
	Status  string
	Headers []HeaderDocumentation
	Body    string
}

// APIDocumentation represents the generated API documentation
type APIDocumentation struct {
	Title       string
	Description string
	BaseURL     string
	Endpoints   []EndpointDocumentation
	GeneratedAt time.Time
	Collection  string // Raw Postman collection JSON
}

// GenerateDocsFromPostman generates HTML documentation from a Postman collection
func GenerateDocsFromPostman(collectionPath, baseURL string) (string, error) {
	generator := NewPostmanDocGenerator(baseURL)
	return generator.GenerateDocs(collectionPath)
}

// GenerateDocs generates HTML documentation from a Postman collection file
func (g *PostmanDocGenerator) GenerateDocs(collectionPath string) (string, error) {
	// Read the Postman collection file
	collectionData, err := ioutil.ReadFile(collectionPath)
	if err != nil {
		return "", fmt.Errorf("failed to load collection: %w", err)
	}

	// Parse the collection
	var collection PostmanCollection
	if err := json.Unmarshal(collectionData, &collection); err != nil {
		return "", fmt.Errorf("failed to parse collection: %w", err)
	}

	// Generate documentation
	doc := g.parseCollection(&collection)

	// Generate HTML
	return g.generateHTML(doc)
}

// parseCollection parses the Postman collection into documentation structure
func (g *PostmanDocGenerator) parseCollection(collection *PostmanCollection) *APIDocumentation {
	// Store raw collection JSON with proper formatting
	rawCollection, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		rawCollection = []byte("{}")
	}

	doc := &APIDocumentation{
		Title:       collection.Info.Name,
		Description: collection.Info.Description,
		BaseURL:     g.baseURL,
		GeneratedAt: time.Now(),
		Endpoints:   []EndpointDocumentation{},
		Collection:  string(rawCollection),
	}

	// Process all items recursively
	g.processItems(collection.Item, "", &doc.Endpoints)

	// Sort endpoints by category and name
	sort.Slice(doc.Endpoints, func(i, j int) bool {
		if doc.Endpoints[i].Category != doc.Endpoints[j].Category {
			return doc.Endpoints[i].Category < doc.Endpoints[j].Category
		}
		return doc.Endpoints[i].Name < doc.Endpoints[j].Name
	})

	return doc
}

// processItems processes Postman items recursively
func (g *PostmanDocGenerator) processItems(items []PostmanItem, category string, endpoints *[]EndpointDocumentation) {
	for _, item := range items {
		if item.Request != nil {
			// This is a request item
			endpoint := g.processRequest(item, category)
			*endpoints = append(*endpoints, endpoint)
		} else if len(item.Item) > 0 {
			// This is a folder, process its children
			newCategory := item.Name
			if category != "" {
				newCategory = category + " / " + item.Name
			}
			g.processItems(item.Item, newCategory, endpoints)
		}
	}
}

// processRequest processes a single request item
func (g *PostmanDocGenerator) processRequest(item PostmanItem, category string) EndpointDocumentation {
	endpoint := EndpointDocumentation{
		Name:        item.Name,
		Method:      item.Request.Method,
		Path:        g.buildPath(item.Request.URL),
		Description: item.Description,
		Category:    category,
	}

	// Process request
	if item.Request != nil {
		endpoint.Request = &RequestDocumentation{
			Headers: g.processHeaders(item.Request.Header),
		}
		if item.Request.Body != nil && item.Request.Body.Raw != "" {
			endpoint.Request.Body = &BodyDocumentation{
				Example: item.Request.Body.Raw,
			}
		}
	}

	// Process response
	if len(item.Response) > 0 {
		response := item.Response[0] // Use first response as example
		endpoint.Response = &ResponseDocumentation{
			Status:  response.Status,
			Headers: g.processHeaders(response.Header),
			Body:    response.Body,
		}
	}

	// Process tests
	endpoint.Tests = g.extractTests(item.Event)

	return endpoint
}

// processHeaders processes request/response headers
func (g *PostmanDocGenerator) processHeaders(headers []PostmanHeader) []HeaderDocumentation {
	var result []HeaderDocumentation
	for _, header := range headers {
		result = append(result, HeaderDocumentation{
			Name:  header.Key,
			Value: header.Value,
		})
	}
	return result
}

// buildPath builds the request path from Postman URL
func (g *PostmanDocGenerator) buildPath(url PostmanURL) string {
	if url.Raw != "" {
		// Remove base URL if present
		path := url.Raw
		if strings.HasPrefix(path, g.baseURL) {
			path = strings.TrimPrefix(path, g.baseURL)
		}
		return path
	}
	
	// Build from components
	if len(url.Path) > 0 {
		return "/" + strings.Join(url.Path, "/")
	}
	return "/"
}

// extractTests extracts test names from Postman events
func (g *PostmanDocGenerator) extractTests(events []PostmanEvent) []string {
	var tests []string
	for _, event := range events {
		if event.Listen == "test" {
			for _, line := range event.Script.Exec {
				if strings.Contains(line, "pm.test") {
					// Extract test name from pm.test("test name", ...)
					start := strings.Index(line, `"`)
					if start != -1 {
						end := strings.Index(line[start+1:], `"`)
						if end != -1 {
							testName := line[start+1 : start+1+end]
							tests = append(tests, testName)
						}
					}
				}
			}
		}
	}
	return tests
}

// generateHTML generates the HTML documentation
func (g *PostmanDocGenerator) generateHTML(doc *APIDocumentation) (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - API Documentation</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }

        .header {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 20px;
            padding: 3rem;
            margin-bottom: 2rem;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            text-align: center;
        }

        .header h1 {
            font-size: 3rem;
            font-weight: 700;
            background: linear-gradient(135deg, #667eea, #764ba2);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 1rem;
        }

        .header p {
            font-size: 1.2rem;
            color: #666;
            margin-bottom: 2rem;
        }

        .action-buttons {
            display: flex;
            gap: 1rem;
            justify-content: center;
            flex-wrap: wrap;
        }

        .btn {
            padding: 1rem 2rem;
            border: none;
            border-radius: 50px;
            font-weight: 600;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-3px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: linear-gradient(135deg, #f093fb, #f5576c);
            color: white;
        }

        .btn-secondary:hover {
            transform: translateY(-3px);
            box-shadow: 0 10px 25px rgba(240, 147, 251, 0.4);
        }

        .meta-info {
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 1.5rem;
            margin-bottom: 2rem;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }

        .meta-info p {
            margin: 0.5rem 0;
            color: #666;
        }

        .json-section {
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }

        .json-section h2 {
            font-size: 1.8rem;
            font-weight: 600;
            color: #333;
            margin-bottom: 1rem;
            text-align: center;
        }

        .json-section p {
            color: #666;
            margin-bottom: 1.5rem;
            text-align: center;
        }

        .json-container {
            background: #f8f9fa;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            padding: 1.5rem;
            position: relative;
        }

        .json-content {
            font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
            font-size: 0.9rem;
            line-height: 1.5;
            color: #333;
            white-space: pre-wrap;
            word-break: break-all;
            max-height: 500px;
            overflow-y: auto;
        }

        .copy-button {
            position: absolute;
            top: 1rem;
            right: 1rem;
            background: linear-gradient(135deg, #28a745, #20c997);
            color: white;
            border: none;
            padding: 0.5rem 1rem;
            border-radius: 25px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            font-size: 0.9rem;
        }

        .copy-button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(40, 167, 69, 0.3);
        }

        .notification {
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 1rem 1.5rem;
            border-radius: 8px;
            color: white;
            font-weight: 600;
            z-index: 1000;
            animation: slideIn 0.3s ease;
            max-width: 300px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }

        .notification.success {
            background: linear-gradient(135deg, #28a745, #20c997);
        }

        .notification.error {
            background: linear-gradient(135deg, #dc3545, #c82333);
        }

        .notification.hide {
            animation: slideOut 0.3s ease;
        }

        @keyframes slideIn {
            from { transform: translateX(100%); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }

        @keyframes slideOut {
            from { transform: translateX(0); opacity: 1; }
            to { transform: translateX(100%); opacity: 0; }
        }

        .instructions {
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }

        .instructions h3 {
            color: #333;
            margin-bottom: 1rem;
        }

        .instructions ol {
            margin-left: 1.5rem;
        }

        .instructions li {
            margin-bottom: 0.5rem;
            color: #666;
        }

        .instructions code {
            background: #f8f9fa;
            padding: 0.2rem 0.4rem;
            border-radius: 4px;
            font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
            font-size: 0.9rem;
        }

        @media (max-width: 768px) {
            .container {
                padding: 1rem;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            .action-buttons {
                flex-direction: column;
                align-items: center;
            }
            
            .json-content {
                font-size: 0.8rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>{{.Description}}</p>
            <div class="action-buttons">
                <button class="btn btn-primary" onclick="copyPostmanCollection()">
                    📋 Copy Collection JSON
                </button>
                <button class="btn btn-secondary" onclick="downloadPostmanCollection()">
                    ⬇️ Download Collection
                </button>
            </div>
        </div>

        <div class="meta-info">
            <p><strong>Base URL:</strong> {{.BaseURL}}</p>
            <p><strong>Generated:</strong> {{.GeneratedAt.Format "2006-01-02 15:04:05 MST"}}</p>
            <p><strong>Total Endpoints:</strong> {{len .Endpoints}}</p>
        </div>

        <div class="instructions">
            <h3>How to Import into Postman:</h3>
            <ol>
                <li>Click the <strong>"Copy Collection JSON"</strong> button above</li>
                <li>Open Postman and click <strong>"Import"</strong> in the top left</li>
                <li>Select <strong>"Raw text"</strong> and paste the JSON</li>
                <li>Click <strong>"Continue"</strong> and then <strong>"Import"</strong></li>
                <li>Your collection will be imported with all endpoints ready to test!</li>
            </ol>
        </div>

        <div class="json-section">
            <h2>Postman Collection JSON</h2>
            <p>Copy the JSON below and import it into Postman to test all API endpoints</p>
            <div class="json-container">
                <button class="copy-button" onclick="copyPostmanCollection()">📋 Copy</button>
                <pre class="json-content" id="json-content">{{.Collection | html}}</pre>
            </div>
        </div>
    </div>

    <script>
        // Store the Postman collection data
        const postmanCollection = {{.Collection}};
        
        function copyPostmanCollection() {
            try {
                navigator.clipboard.writeText(JSON.stringify(postmanCollection, null, 2)).then(function() {
                    showNotification('Postman collection copied to clipboard!', 'success');
                }).catch(function(err) {
                    showNotification('Failed to copy Postman collection: ' + err, 'error');
                });
            } catch (err) {
                showNotification('Clipboard API not available or failed: ' + err, 'error');
            }
        }

        function downloadPostmanCollection() {
            try {
                const filename = "{{.Title | anchorize}}-postman-collection.json";
                const element = document.createElement('a');
                element.setAttribute('href', 'data:application/json;charset=utf-8,' + encodeURIComponent(JSON.stringify(postmanCollection, null, 2)));
                element.setAttribute('download', filename);
                element.style.display = 'none';
                document.body.appendChild(element);
                element.click();
                document.body.removeChild(element);
                showNotification('Postman collection download initiated!', 'success');
            } catch (err) {
                showNotification('Failed to download Postman collection: ' + err, 'error');
            }
        }

        function showNotification(message, type) {
            const notificationContainer = document.getElementById('notification-container') || createNotificationContainer();
            const notification = document.createElement('div');
            notification.className = 'notification ' + type;
            notification.textContent = message;
            notificationContainer.appendChild(notification);

            setTimeout(() => {
                notification.classList.add('hide');
                notification.addEventListener('animationend', () => {
                    notification.remove();
                });
            }, 5000);
        }

        function createNotificationContainer() {
            const container = document.createElement('div');
            container.id = 'notification-container';
            container.style.cssText = 'position: fixed; top: 20px; right: 20px; z-index: 1000;';
            document.body.appendChild(container);
            return container;
        }
    </script>
    <div id="notification-container"></div>
</body>
</html>`

	// Create template with custom functions
	funcMap := template.FuncMap{
		"lower":      strings.ToLower,
		"anchorize":  func(s string) string { return strings.ReplaceAll(strings.ToLower(s), " ", "-") },
		"json":       func(s string) string { return s }, // Pass through JSON as-is
	}

	t, err := template.New("docs").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}