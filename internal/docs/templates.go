package docs

import (
	"html/template"
	"strings"
)

// HTMLTemplates contains all HTML templates for documentation generation
type HTMLTemplates struct {
	MainTemplate    *template.Template
	EndpointPartial *template.Template
}

// NewHTMLTemplates creates a new HTMLTemplates instance
func NewHTMLTemplates() *HTMLTemplates {
	templates := &HTMLTemplates{}
	templates.initTemplates()
	return templates
}

// initTemplates initializes all HTML templates
func (t *HTMLTemplates) initTemplates() {
	t.MainTemplate = template.Must(template.New("main").Parse(mainTemplate))
	t.EndpointPartial = template.Must(template.New("endpoint").Parse(endpointTemplate))
}

// mainTemplate is the main HTML template
const mainTemplate = `
<!DOCTYPE html>
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
            background: #f8fafc;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 0;
            margin-bottom: 30px;
            border-radius: 12px;
            text-align: center;
        }
        
        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
            font-weight: 700;
        }
        
        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }
        
        .controls {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 30px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 15px;
        }
        
        .search-box {
            flex: 1;
            min-width: 300px;
            padding: 12px 16px;
            border: 2px solid #e2e8f0;
            border-radius: 8px;
            font-size: 16px;
            transition: border-color 0.3s;
        }
        
        .search-box:focus {
            outline: none;
            border-color: #667eea;
        }
        
        .btn {
            padding: 12px 24px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
            text-decoration: none;
            display: inline-block;
        }
        
        .btn-primary {
            background: #667eea;
            color: white;
        }
        
        .btn-primary:hover {
            background: #5a67d8;
            transform: translateY(-2px);
        }
        
        .btn-secondary {
            background: #e2e8f0;
            color: #4a5568;
        }
        
        .btn-secondary:hover {
            background: #cbd5e0;
        }
        
        .endpoints {
            display: grid;
            gap: 20px;
        }
        
        .endpoint {
            background: white;
            border-radius: 12px;
            padding: 24px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            transition: transform 0.3s, box-shadow 0.3s;
        }
        
        .endpoint:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 20px rgba(0,0,0,0.15);
        }
        
        .endpoint-header {
            display: flex;
            align-items: center;
            margin-bottom: 16px;
            flex-wrap: wrap;
            gap: 12px;
        }
        
        .method-badge {
            padding: 6px 12px;
            border-radius: 6px;
            font-weight: 700;
            font-size: 14px;
            text-transform: uppercase;
        }
        
        .method-get { background: #48bb78; color: white; }
        .method-post { background: #4299e1; color: white; }
        .method-put { background: #ed8936; color: white; }
        .method-delete { background: #f56565; color: white; }
        .method-patch { background: #9f7aea; color: white; }
        
        .endpoint-title {
            font-size: 1.5rem;
            font-weight: 700;
            color: #2d3748;
            flex: 1;
        }
        
        .endpoint-url {
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            background: #f7fafc;
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 14px;
            color: #4a5568;
            border: 1px solid #e2e8f0;
        }
        
        .endpoint-description {
            color: #718096;
            margin-bottom: 20px;
            font-size: 1.1rem;
        }
        
        .section {
            margin-bottom: 20px;
        }
        
        .section-title {
            font-size: 1.2rem;
            font-weight: 600;
            color: #2d3748;
            margin-bottom: 12px;
            display: flex;
            align-items: center;
        }
        
        .section-title::before {
            content: '';
            width: 4px;
            height: 20px;
            background: #667eea;
            margin-right: 8px;
            border-radius: 2px;
        }
        
        .headers-table, .responses-table {
            width: 100%;
            border-collapse: collapse;
            background: #f7fafc;
            border-radius: 8px;
            overflow: hidden;
        }
        
        .headers-table th, .responses-table th {
            background: #e2e8f0;
            padding: 12px;
            text-align: left;
            font-weight: 600;
            color: #4a5568;
        }
        
        .headers-table td, .responses-table td {
            padding: 12px;
            border-bottom: 1px solid #e2e8f0;
        }
        
        .code-block {
            background: #2d3748;
            color: #e2e8f0;
            padding: 16px;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            overflow-x: auto;
            white-space: pre-wrap;
        }
        
        .test-section {
            background: #f7fafc;
            padding: 20px;
            border-radius: 8px;
            margin-top: 20px;
        }
        
        .test-button {
            background: #48bb78;
            color: white;
            padding: 10px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            margin-right: 10px;
            margin-bottom: 10px;
        }
        
        .test-button:hover {
            background: #38a169;
        }
        
        .copy-button {
            background: #4299e1;
            color: white;
            padding: 8px 16px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            margin-left: 10px;
        }
        
        .copy-button:hover {
            background: #3182ce;
        }
        
        .footer {
            text-align: center;
            padding: 40px 0;
            color: #718096;
            margin-top: 50px;
        }
        
        .hidden {
            display: none;
        }
        
        @media (max-width: 768px) {
            .container {
                padding: 10px;
            }
            
            .header h1 {
                font-size: 2rem;
            }
            
            .controls {
                flex-direction: column;
                align-items: stretch;
            }
            
            .search-box {
                min-width: auto;
            }
            
            .endpoint-header {
                flex-direction: column;
                align-items: flex-start;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p>{{.Description}}</p>
            <p>Base URL: <code>{{.BaseURL}}</code></p>
        </div>
        
        <div class="controls">
            <input type="text" id="searchInput" class="search-box" placeholder="Search endpoints...">
            <div>
                <button class="btn btn-primary" onclick="toggleTestMode()">Toggle Test Mode</button>
                <button class="btn btn-secondary" onclick="downloadPostmanCollection()">Download Postman Collection</button>
            </div>
        </div>
        
        <div class="endpoints" id="endpointsList">
            {{range .Endpoints}}
            <div class="endpoint" data-method="{{.Method}}" data-name="{{.Name}}">
                <div class="endpoint-header">
                    <span class="method-badge method-{{.Method | lower}}">{{.Method}}</span>
                    <h2 class="endpoint-title">{{.Name}}</h2>
                    <code class="endpoint-url">{{.URL}}</code>
                </div>
                
                {{if .Description}}
                <div class="endpoint-description">{{.Description}}</div>
                {{end}}
                
                {{if .Headers}}
                <div class="section">
                    <h3 class="section-title">Headers</h3>
                    <table class="headers-table">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Value</th>
                                <th>Required</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Headers}}
                            <tr>
                                <td><code>{{.Name}}</code></td>
                                <td>{{.Value}}</td>
                                <td>{{if .Required}}Yes{{else}}No{{end}}</td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
                {{end}}
                
                {{if .Body}}
                <div class="section">
                    <h3 class="section-title">Request Body</h3>
                    <div class="code-block">{{.Body.Content}}</div>
                </div>
                {{end}}
                
                {{if .Responses}}
                <div class="section">
                    <h3 class="section-title">Responses</h3>
                    <table class="responses-table">
                        <thead>
                            <tr>
                                <th>Status Code</th>
                                <th>Description</th>
                                <th>Body</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .Responses}}
                            <tr>
                                <td><code>{{.Code}}</code></td>
                                <td>{{.Description}}</td>
                                <td>
                                    {{if .Body}}
                                    <div class="code-block">{{.Body}}</div>
                                    {{else}}
                                    <em>No body</em>
                                    {{end}}
                                </td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
                {{end}}
                
                <div class="test-section hidden" id="test-{{.Name | replace " " "-" | lower}}">
                    <h3 class="section-title">Test This Endpoint</h3>
                    <button class="test-button" onclick="testEndpoint('{{.Method}}', '{{.URL}}')">Send Request</button>
                    <button class="copy-button" onclick="copyToClipboard('{{.URL}}')">Copy URL</button>
                    <div id="response-{{.Name | replace " " "-" | lower}}" class="response-area"></div>
                </div>
            </div>
            {{end}}
        </div>
        
        <div class="footer">
            <p>Generated on {{.GeneratedAt}}</p>
            <p>API Documentation powered by Go Server</p>
        </div>
    </div>
    
    <script>
        // Search functionality
        document.getElementById('searchInput').addEventListener('input', function(e) {
            const searchTerm = e.target.value.toLowerCase();
            const endpoints = document.querySelectorAll('.endpoint');
            
            endpoints.forEach(endpoint => {
                const name = endpoint.querySelector('.endpoint-title').textContent.toLowerCase();
                const method = endpoint.querySelector('.method-badge').textContent.toLowerCase();
                const url = endpoint.querySelector('.endpoint-url').textContent.toLowerCase();
                
                if (name.includes(searchTerm) || method.includes(searchTerm) || url.includes(searchTerm)) {
                    endpoint.style.display = 'block';
                } else {
                    endpoint.style.display = 'none';
                }
            });
        });
        
        // Toggle test mode
        function toggleTestMode() {
            const testSections = document.querySelectorAll('.test-section');
            testSections.forEach(section => {
                section.classList.toggle('hidden');
            });
        }
        
        // Test endpoint
        async function testEndpoint(method, url) {
            try {
                const response = await fetch(url, {
                    method: method,
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                
                const data = await response.text();
                const responseDiv = document.getElementById('response-' + url.split('/').pop().toLowerCase());
                responseDiv.innerHTML = '<div class="code-block">Status: ' + response.status + '\n' + data + '</div>';
            } catch (error) {
                const responseDiv = document.getElementById('response-' + url.split('/').pop().toLowerCase());
                responseDiv.innerHTML = '<div class="code-block" style="background: #f56565;">Error: ' + error.message + '</div>';
            }
        }
        
        // Copy to clipboard
        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(function() {
                alert('URL copied to clipboard!');
            });
        }
        
        // Download Postman collection
        function downloadPostmanCollection() {
            // This would need to be implemented to fetch the actual Postman collection
            alert('Postman collection download would be implemented here');
        }
    </script>
</body>
</html>
`

// endpointTemplate is a partial template for individual endpoints
const endpointTemplate = `
<div class="endpoint">
    <div class="endpoint-header">
        <span class="method-badge method-{{.Method | lower}}">{{.Method}}</span>
        <h2 class="endpoint-title">{{.Name}}</h2>
        <code class="endpoint-url">{{.URL}}</code>
    </div>
    
    {{if .Description}}
    <div class="endpoint-description">{{.Description}}</div>
    {{end}}
    
    <!-- Additional endpoint content would go here -->
</div>
`

// GenerateHTML generates HTML documentation from API documentation
func (t *HTMLTemplates) GenerateHTML(doc *APIDocumentation) (string, error) {
	var result strings.Builder
	err := t.MainTemplate.Execute(&result, doc)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
