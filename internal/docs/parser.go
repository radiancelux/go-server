package docs

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"
)

// PostmanParser handles parsing of Postman collections
type PostmanParser struct{}

// NewPostmanParser creates a new PostmanParser
func NewPostmanParser() *PostmanParser {
	return &PostmanParser{}
}

// ParseCollection parses a Postman collection from file
func (p *PostmanParser) ParseCollection(collectionPath string) (*PostmanCollection, error) {
	data, err := ioutil.ReadFile(collectionPath)
	if err != nil {
		return nil, err
	}

	var collection PostmanCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// ParseCollectionFromBytes parses a Postman collection from bytes
func (p *PostmanParser) ParseCollectionFromBytes(data []byte) (*PostmanCollection, error) {
	var collection PostmanCollection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, err
	}

	return &collection, nil
}

// ExtractEndpoints extracts all endpoints from a collection
func (p *PostmanParser) ExtractEndpoints(collection *PostmanCollection) []EndpointDocumentation {
	var endpoints []EndpointDocumentation
	p.extractEndpointsFromItems(collection.Item, &endpoints)
	
	// Sort endpoints by method and name
	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Method != endpoints[j].Method {
			return endpoints[i].Method < endpoints[j].Method
		}
		return endpoints[i].Name < endpoints[j].Name
	})
	
	return endpoints
}

// extractEndpointsFromItems recursively extracts endpoints from collection items
func (p *PostmanParser) extractEndpointsFromItems(items []PostmanItem, endpoints *[]EndpointDocumentation) {
	for _, item := range items {
		if item.Request != nil {
			// This is a request item
			endpoint := p.convertRequestToEndpoint(item)
			*endpoints = append(*endpoints, endpoint)
		} else if len(item.Item) > 0 {
			// This is a folder, recurse into it
			p.extractEndpointsFromItems(item.Item, endpoints)
		}
	}
}

// convertRequestToEndpoint converts a PostmanItem to EndpointDocumentation
func (p *PostmanParser) convertRequestToEndpoint(item PostmanItem) EndpointDocumentation {
	endpoint := EndpointDocumentation{
		Name:        item.Name,
		Description: item.Description,
		Method:      item.Request.Method,
		URL:         p.buildURL(item.Request.URL),
		Headers:     p.convertHeaders(item.Request.Header),
		Responses:   p.convertResponses(item.Response),
	}

	// Add body if present
	if item.Request.Body != nil && item.Request.Body.Raw != "" {
		endpoint.Body = &BodyDocumentation{
			Type:    item.Request.Body.Mode,
			Content: item.Request.Body.Raw,
		}
	}

	return endpoint
}

// buildURL constructs the full URL from PostmanURL
func (p *PostmanParser) buildURL(url *PostmanURL) string {
	if url == nil {
		return ""
	}

	// Use raw URL if available
	if url.Raw != "" {
		return url.Raw
	}

	// Build URL from components
	var result strings.Builder
	
	if url.Protocol != "" {
		result.WriteString(url.Protocol)
		result.WriteString("://")
	}
	
	if len(url.Host) > 0 {
		result.WriteString(strings.Join(url.Host, "."))
	}
	
	if len(url.Path) > 0 {
		result.WriteString("/")
		result.WriteString(strings.Join(url.Path, "/"))
	}
	
	if len(url.Query) > 0 {
		result.WriteString("?")
		var queryParts []string
		for _, q := range url.Query {
			queryParts = append(queryParts, q.Key+"="+q.Value)
		}
		result.WriteString(strings.Join(queryParts, "&"))
	}
	
	return result.String()
}

// convertHeaders converts Postman headers to documentation headers
func (p *PostmanParser) convertHeaders(headers []PostmanHeader) []HeaderDocumentation {
	var docHeaders []HeaderDocumentation
	for _, header := range headers {
		docHeaders = append(docHeaders, HeaderDocumentation{
			Name:        header.Key,
			Value:       header.Value,
			Description: "",
			Required:    false,
		})
	}
	return docHeaders
}

// convertResponses converts Postman responses to documentation responses
func (p *PostmanParser) convertResponses(responses []PostmanResponse) []ResponseDocumentation {
	var docResponses []ResponseDocumentation
	for _, response := range responses {
		docResponses = append(docResponses, ResponseDocumentation{
			Code:        response.Code,
			Status:      response.Status,
			Description: response.Name,
			Headers:     p.convertHeaders(response.Header),
			Body:        response.Body,
		})
	}
	return docResponses
}
