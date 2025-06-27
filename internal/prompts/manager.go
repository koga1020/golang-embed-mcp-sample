package prompts

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Manager manages embedded prompt templates
type Manager struct {
	embeddedPrompts embed.FS
}

// NewManager creates a new manager with embedded prompts
func NewManager(embeddedFS embed.FS) *Manager {
	log.Printf("[DEBUG] Creating prompts manager")
	return &Manager{
		embeddedPrompts: embeddedFS,
	}
}

// GetPrompts returns all available prompts
func (m *Manager) GetPrompts() []*mcp.ServerPrompt {
	return m.GetPromptsWithFilter(nil)
}

// GetPromptsWithFilter returns filtered prompts
func (m *Manager) GetPromptsWithFilter(filter []string) []*mcp.ServerPrompt {
	var prompts []*mcp.ServerPrompt

	log.Printf("[DEBUG] Loading embedded prompts with filter: %v", filter)

	err := fs.WalkDir(m.embeddedPrompts, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			// Extract prompt name (without .md extension)
			name := strings.TrimSuffix(filepath.Base(path), ".md")

			// Apply filter if specified
			if filter != nil && len(filter) > 0 {
				found := false
				for _, promptName := range filter {
					if strings.EqualFold(name, promptName) {
						found = true
						break
					}
				}
				if !found {
					log.Printf("[DEBUG] Skipping prompt %s (not in filter)", path)
					return nil
				}
			}

			description := fmt.Sprintf("%s prompt template", strings.Title(name))

			log.Printf("[DEBUG] Creating embedded prompt: name=%s, description=%s", name, description)

			prompt := &mcp.ServerPrompt{
				Prompt: &mcp.Prompt{
					Name:        name,
					Description: description,
				},
				Handler: m.handlePrompt,
			}
			prompts = append(prompts, prompt)
		}

		return nil
	})

	if err != nil {
		log.Printf("[ERROR] Error walking embedded prompts: %v", err)
		return nil
	}

	log.Printf("[DEBUG] Total embedded prompts found: %d", len(prompts))
	return prompts
}

// handlePrompt handles prompt requests
func (m *Manager) handlePrompt(_ context.Context, _ *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	log.Printf("[DEBUG] Handling prompt request for: %s", params.Name)

	content, err := m.readEmbeddedPrompt(params.Name)
	if err != nil {
		log.Printf("[ERROR] Failed to read prompt %s: %v", params.Name, err)
		return nil, fmt.Errorf("failed to read prompt: %w", err)
	}

	log.Printf("[DEBUG] Successfully read prompt: %d bytes", len(content))

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: string(content),
				},
			},
		},
	}, nil
}

// readEmbeddedPrompt reads an embedded prompt file
func (m *Manager) readEmbeddedPrompt(promptName string) ([]byte, error) {
	promptPath := fmt.Sprintf("prompts/%s.md", promptName)

	log.Printf("[DEBUG] Reading embedded prompt: %s", promptPath)

	content, err := m.embeddedPrompts.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("prompt not found: %s", promptName)
	}

	return content, nil
}