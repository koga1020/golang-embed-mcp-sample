# Embedded Resources Demo

This file demonstrates how various resource types can be embedded into Go binaries and served through MCP.

## Embed Directive

```go
//go:embed resources/*
var embeddedResources embed.FS
```

## Supported Formats

- **JSON**: Configuration and structured data
- **YAML**: Schemas and configuration files  
- **Markdown**: Documentation and templates
- **Text**: Plain text resources

## Benefits

1. **Build-time inclusion**: All resources bundled at compile time
2. **No external files**: Eliminates deployment complexity
3. **Type safety**: Go's embed.FS provides safe file access
4. **MCP integration**: Automatic MIME type detection and serving

## Access Pattern

Resources are served using the `embedded://` URI scheme:
- `embedded://config.json`
- `embedded://schema.yaml`
- `embedded://info.md`