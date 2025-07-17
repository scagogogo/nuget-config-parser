# API Reference

The NuGet Config Parser library provides a comprehensive API for working with NuGet configuration files. This section documents all public interfaces, types, and methods.

## Package Overview

The library is organized into several packages, each with a specific purpose:

| Package | Description |
|---------|-------------|
| [`pkg/nuget`](./core.md) | Main API package providing the primary interface |
| [`pkg/parser`](./parser.md) | Configuration file parsing functionality |
| [`pkg/editor`](./editor.md) | Position-aware configuration editing |
| [`pkg/finder`](./finder.md) | Configuration file discovery |
| [`pkg/manager`](./manager.md) | Configuration management operations |
| [`pkg/types`](./types.md) | Data type definitions |
| [`pkg/utils`](./utils.md) | Utility functions |
| [`pkg/errors`](./errors.md) | Error types and handling |
| [`pkg/constants`](./constants.md) | Constants and default values |

## Quick Reference

### Core API

```go
import "github.com/scagogogo/nuget-config-parser/pkg/nuget"

// Create API instance
api := nuget.NewAPI()

// Parse configuration
config, err := api.ParseFromFile("/path/to/NuGet.Config")

// Find configuration files
configPath, err := api.FindConfigFile()

// Modify configuration
api.AddPackageSource(config, "source", "https://example.com", "3")

// Save configuration
err = api.SaveConfig(config, "/path/to/NuGet.Config")
```

### Position-Aware Editing

```go
import "github.com/scagogogo/nuget-config-parser/pkg/editor"

// Parse with position tracking
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")

// Create editor
editor := api.CreateConfigEditor(parseResult)

// Make changes
err = editor.AddPackageSource("new-source", "https://new.com", "3")

// Apply changes
modifiedContent, err := editor.ApplyEdits()
```

## Common Types

### NuGetConfig

The main configuration structure:

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

### PackageSource

Represents a single package source:

```go
type PackageSource struct {
    Key             string `xml:"key,attr"`
    Value           string `xml:"value,attr"`
    ProtocolVersion string `xml:"protocolVersion,attr,omitempty"`
}
```

## Error Handling

The library provides structured error handling:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle file not found
    } else if errors.IsParseError(err) {
        // Handle parsing error
    } else if errors.IsFormatError(err) {
        // Handle format error
    }
}
```

## Method Categories

### Parsing Methods

- `ParseFromFile(filePath string) (*types.NuGetConfig, error)`
- `ParseFromString(content string) (*types.NuGetConfig, error)`
- `ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)`
- `ParseFromFileWithPositions(filePath string) (*parser.ParseResult, error)`

### Finding Methods

- `FindConfigFile() (string, error)`
- `FindAllConfigFiles() []string`
- `FindProjectConfig(startDir string) (string, error)`
- `FindAndParseConfig() (*types.NuGetConfig, string, error)`

### Package Source Methods

- `AddPackageSource(config *types.NuGetConfig, key, value, protocolVersion string)`
- `RemovePackageSource(config *types.NuGetConfig, key string) bool`
- `GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource`
- `GetAllPackageSources(config *types.NuGetConfig) []types.PackageSource`
- `EnablePackageSource(config *types.NuGetConfig, key string)`
- `DisablePackageSource(config *types.NuGetConfig, key string)`
- `IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool`

### Credential Methods

- `AddCredential(config *types.NuGetConfig, sourceKey, username, password string)`
- `RemoveCredential(config *types.NuGetConfig, sourceKey string) bool`
- `GetCredential(config *types.NuGetConfig, sourceKey string) *types.SourceCredential`

### Configuration Methods

- `AddConfigOption(config *types.NuGetConfig, key, value string)`
- `RemoveConfigOption(config *types.NuGetConfig, key string) bool`
- `GetConfigOption(config *types.NuGetConfig, key string) string`
- `SetActivePackageSource(config *types.NuGetConfig, key, value string)`
- `GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource`

### Serialization Methods

- `SaveConfig(config *types.NuGetConfig, filePath string) error`
- `SerializeToXML(config *types.NuGetConfig) (string, error)`
- `CreateDefaultConfig() *types.NuGetConfig`
- `InitializeDefaultConfig(filePath string) error`

### Editor Methods

- `CreateConfigEditor(parseResult *parser.ParseResult) *editor.ConfigEditor`
- `AddPackageSource(key, value, protocolVersion string) error`
- `RemovePackageSource(sourceKey string) error`
- `UpdatePackageSourceURL(sourceKey, newURL string) error`
- `UpdatePackageSourceVersion(sourceKey, newVersion string) error`
- `ApplyEdits() ([]byte, error)`

## Best Practices

### Error Handling

Always check for errors and handle them appropriately:

```go
config, err := api.ParseFromFile(configPath)
if err != nil {
    if errors.IsNotFoundError(err) {
        // Create default configuration
        config = api.CreateDefaultConfig()
    } else {
        return fmt.Errorf("failed to parse config: %w", err)
    }
}
```

### Resource Management

The library doesn't require explicit resource cleanup, but be mindful of file operations:

```go
// Good: Use the API methods
err := api.SaveConfig(config, configPath)

// Avoid: Manual file operations when API methods are available
```

### Performance

- Reuse API instances when possible
- Cache parsed configurations for repeated access
- Use position-aware editing for minimal file changes

## Thread Safety

The library is not thread-safe by design. If you need to use it in concurrent scenarios:

- Create separate API instances for each goroutine
- Use appropriate synchronization mechanisms
- Avoid sharing configuration objects between goroutines without proper locking

## Next Steps

Explore the detailed documentation for each package:

- [Core API](./core.md) - Main API interface
- [Parser](./parser.md) - Configuration parsing
- [Editor](./editor.md) - Position-aware editing
- [Types](./types.md) - Data structures
- [Examples](/examples/) - Usage examples
