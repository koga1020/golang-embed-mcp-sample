package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/koga1020/golang-embed-mcp-sample/internal/prompts"
	"github.com/koga1020/golang-embed-mcp-sample/internal/resources"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed prompts/*
var embeddedPrompts embed.FS

//go:embed resources/*
var embeddedResources embed.FS

func main() {
	// Parse command-line options
	var promptsFilter string
	var resourcesFilter string
	flag.StringVar(&promptsFilter, "prompts", "", "Comma-separated list of prompts to load")
	flag.StringVar(&resourcesFilter, "resources", "", "Comma-separated list of resources to load")
	flag.Parse()

	log.Printf("[DEBUG] Starting embed-mcp server")

	// Parse prompt filter
	var promptFilter []string
	if promptsFilter != "" {
		promptFilter = strings.Split(promptsFilter, ",")
		for i := range promptFilter {
			promptFilter[i] = strings.TrimSpace(promptFilter[i])
		}
		log.Printf("[DEBUG] Prompt filter: %v", promptFilter)
	}

	// Parse resource filter
	var resourceFilter []string
	if resourcesFilter != "" {
		resourceFilter = strings.Split(resourcesFilter, ",")
		for i := range resourceFilter {
			resourceFilter[i] = strings.TrimSpace(resourceFilter[i])
		}
		log.Printf("[DEBUG] Resource filter: %v", resourceFilter)
	}

	// Create MCP server
	server := mcp.NewServer("embed-mcp", "v1.0.0", nil)
	log.Printf("[DEBUG] Created MCP server")

	// Initialize prompt manager
	promptManager := prompts.NewManager(embeddedPrompts)
	log.Printf("[DEBUG] Created prompt manager")

	// Register prompt handlers
	var serverPrompts []*mcp.ServerPrompt
	if len(promptFilter) > 0 {
		serverPrompts = promptManager.GetPromptsWithFilter(promptFilter)
		log.Printf("[DEBUG] Got %d filtered prompts", len(serverPrompts))
	} else {
		serverPrompts = promptManager.GetPrompts()
		log.Printf("[DEBUG] Got %d prompts (all)", len(serverPrompts))
	}

	for i, prompt := range serverPrompts {
		log.Printf("[DEBUG] Registering prompt %d: %s", i+1, prompt.Prompt.Name)
		server.AddPrompts(prompt)
	}

	// Initialize resource manager
	resourceManager := resources.NewManager(embeddedResources)
	log.Printf("[DEBUG] Created resource manager")

	// Register resource handlers
	var serverResources []*mcp.ServerResource
	if len(resourceFilter) > 0 {
		serverResources = resourceManager.GetResourcesWithFilter(resourceFilter)
		log.Printf("[DEBUG] Got %d filtered resources", len(serverResources))
	} else {
		serverResources = resourceManager.GetResources()
		log.Printf("[DEBUG] Got %d resources (all)", len(serverResources))
	}

	for i, resource := range serverResources {
		log.Printf("[DEBUG] Registering resource %d: %s", i+1, resource.Resource.Name)
		server.AddResources(resource)
	}

	log.Printf("[DEBUG] All prompts and resources registered, starting server")

	// Run server with stdio transport
	transport := mcp.NewLoggingTransport(mcp.NewStdioTransport(), os.Stderr)
	if err := server.Run(context.Background(), transport); err != nil {
		log.Printf("[ERROR] Server failed: %v", err)
		os.Exit(1)
	}
}