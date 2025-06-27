package resources

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Manager manages embedded resources
type Manager struct {
	embeddedResources embed.FS
}

// NewManager creates a new manager with embedded resources
func NewManager(embeddedFS embed.FS) *Manager {
	log.Printf("[DEBUG] Creating resources manager")
	return &Manager{
		embeddedResources: embeddedFS,
	}
}

// GetResources returns all available resources
func (m *Manager) GetResources() []*mcp.ServerResource {
	return m.GetResourcesWithFilter(nil)
}

// GetResourcesWithFilter returns filtered resources
func (m *Manager) GetResourcesWithFilter(filter []string) []*mcp.ServerResource {
	var resources []*mcp.ServerResource

	log.Printf("[DEBUG] Loading embedded resources with filter: %v", filter)

	err := fs.WalkDir(m.embeddedResources, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			// Extract resource name (without extension)
			ext := filepath.Ext(d.Name())
			name := strings.TrimSuffix(filepath.Base(path), ext)

			// Apply filter if specified
			if filter != nil && len(filter) > 0 {
				found := false
				for _, resourceName := range filter {
					if strings.EqualFold(name, resourceName) {
						found = true
						break
					}
				}
				if !found {
					log.Printf("[DEBUG] Skipping resource %s (not in filter)", path)
					return nil
				}
			}

			// Use embedded:// scheme
			uri := fmt.Sprintf("embedded://%s", path)

			// Detect MIME type
			mimeType := getMimeType(ext)

			// Generate description
			description := fmt.Sprintf("%s resource", strings.Title(name))

			log.Printf("[DEBUG] Creating embedded resource: name=%s, uri=%s, mime=%s", path, uri, mimeType)

			resource := &mcp.ServerResource{
				Resource: &mcp.Resource{
					Name:        path,
					URI:         uri,
					MIMEType:    mimeType,
					Description: description,
				},
				Handler: m.handleResource,
			}
			resources = append(resources, resource)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Error walking embedded resources: %v", err)
		return nil
	}

	log.Printf("[DEBUG] Total embedded resources found: %d", len(resources))
	return resources
}

// handleResource handles resource read requests
func (m *Manager) handleResource(_ context.Context, _ *mcp.ServerSession, params *mcp.ReadResourceParams) (*mcp.ReadResourceResult, error) {
	log.Printf("[DEBUG] Handling resource request for URI: %s", params.URI)

	// Parse URI
	u, err := url.Parse(params.URI)
	if err != nil {
		log.Printf("[ERROR] Failed to parse URI %s: %v", params.URI, err)
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	if u.Scheme != "embedded" {
		log.Printf("[ERROR] Unsupported scheme: %s", u.Scheme)
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	// Read embedded resource
	content, err := m.readEmbeddedResource(u.Path)
	if err != nil {
		log.Printf("[ERROR] Failed to read embedded resource %s: %v", u.Path, err)
		return nil, fmt.Errorf("failed to read embedded resource: %w", err)
	}

	log.Printf("[DEBUG] Successfully read resource: %d bytes", len(content))

	// Detect MIME type
	ext := filepath.Ext(u.Path)
	mimeType := getMimeType(ext)

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      params.URI,
				MIMEType: mimeType,
				Text:     string(content),
			},
		},
	}, nil
}

// readEmbeddedResource reads an embedded resource file
func (m *Manager) readEmbeddedResource(resourcePath string) ([]byte, error) {
	// Normalize path
	cleanPath := path.Clean(strings.TrimPrefix(resourcePath, "/"))
	embeddedPath := path.Join("resources", cleanPath)

	log.Printf("[DEBUG] Reading embedded resource: %s", embeddedPath)

	content, err := m.embeddedResources.ReadFile(embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("embedded resource not found: %s", embeddedPath)
	}

	return content, nil
}

// getMimeType returns MIME type based on file extension
func getMimeType(ext string) string {
	switch strings.ToLower(ext) {
	case ".md":
		return "text/markdown"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".yaml", ".yml":
		return "application/x-yaml"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".xml":
		return "application/xml"
	case ".csv":
		return "text/csv"
	default:
		return "text/plain"
	}
}