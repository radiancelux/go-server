// Package docs provides documentation rendering functionality.
package docs

import (
	"os"
	"strings"
)

// Converter handles markdown to HTML conversion for documentation.
type Converter struct{}

// NewConverter creates a new documentation converter.
func NewConverter() *Converter {
	return &Converter{}
}

// ConvertFile reads a markdown file and converts it to HTML.
func (c *Converter) ConvertFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return c.ConvertMarkdownToHTML(string(content)), nil
}

// ConvertMarkdownToHTML converts basic markdown to HTML.
func (c *Converter) ConvertMarkdownToHTML(markdown string) string {
	html := markdown

	// Headers
	html = strings.ReplaceAll(html, "# ", "<h1>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "#### ", "<h4>")

	// Close headers
	html = strings.ReplaceAll(html, "\n# ", "</h1>\n# ")
	html = strings.ReplaceAll(html, "\n## ", "</h2>\n## ")
	html = strings.ReplaceAll(html, "\n### ", "</h3>\n### ")
	html = strings.ReplaceAll(html, "\n#### ", "</h4>\n#### ")

	// Code blocks
	html = strings.ReplaceAll(html, "```json", "<pre><code class=\"language-json\">")
	html = strings.ReplaceAll(html, "```bash", "<pre><code class=\"language-bash\">")
	html = strings.ReplaceAll(html, "```", "</code></pre>")

	// Inline code
	html = strings.ReplaceAll(html, "`", "<code>")
	html = strings.ReplaceAll(html, "<code>", "<code>")
	html = strings.ReplaceAll(html, "</code>", "</code>")

	// Bold text
	html = strings.ReplaceAll(html, "**", "<strong>")
	html = strings.ReplaceAll(html, "**", "</strong>")

	// Line breaks
	html = strings.ReplaceAll(html, "\n", "<br>\n")

	// Wrap in HTML structure
	return c.wrapInHTML(html)
}

// wrapInHTML wraps content in a complete HTML document.
func (c *Converter) wrapInHTML(content string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>Go Server API Documentation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        h1, h2, h3, h4 { color: #333; }
        pre { background: #f5f5f5; padding: 15px; border-radius: 5px; overflow-x: auto; }
        code { background: #e8e8e8; padding: 2px 4px; border-radius: 3px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
` + content + `
</body>
</html>`
}
