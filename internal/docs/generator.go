package docs

import (
	"fmt"
	"time"
)

// PostmanDocGenerator generates HTML documentation from Postman collections
type PostmanDocGenerator struct {
	baseURL   string
	parser    *PostmanParser
	templates *HTMLTemplates
}

// NewPostmanDocGenerator creates a new PostmanDocGenerator
func NewPostmanDocGenerator(baseURL string) *PostmanDocGenerator {
	return &PostmanDocGenerator{
		baseURL:   baseURL,
		parser:    NewPostmanParser(),
		templates: NewHTMLTemplates(),
	}
}

// GenerateDocs generates HTML documentation from a Postman collection file
func (g *PostmanDocGenerator) GenerateDocs(collectionPath string) (string, error) {
	// Parse the Postman collection
	collection, err := g.parser.ParseCollection(collectionPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse collection: %w", err)
	}

	// Generate API documentation
	apiDoc := g.generateAPIDocumentation(collection)

	// Generate HTML
	html, err := g.templates.GenerateHTML(apiDoc)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	return html, nil
}

// GenerateDocsFromBytes generates HTML documentation from Postman collection bytes
func (g *PostmanDocGenerator) GenerateDocsFromBytes(collectionData []byte) (string, error) {
	// Parse the Postman collection
	collection, err := g.parser.ParseCollectionFromBytes(collectionData)
	if err != nil {
		return "", fmt.Errorf("failed to parse collection: %w", err)
	}

	// Generate API documentation
	apiDoc := g.generateAPIDocumentation(collection)

	// Generate HTML
	html, err := g.templates.GenerateHTML(apiDoc)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	return html, nil
}

// generateAPIDocumentation creates APIDocumentation from PostmanCollection
func (g *PostmanDocGenerator) generateAPIDocumentation(collection *PostmanCollection) *APIDocumentation {
	endpoints := g.parser.ExtractEndpoints(collection)

	return &APIDocumentation{
		Title:       collection.Info.Name,
		Description: collection.Info.Description,
		BaseURL:     g.baseURL,
		Endpoints:   endpoints,
		GeneratedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
}

// GenerateDocsFromPostman is a convenience function for backward compatibility
func GenerateDocsFromPostman(collectionPath, baseURL string) (string, error) {
	generator := NewPostmanDocGenerator(baseURL)
	return generator.GenerateDocs(collectionPath)
}
