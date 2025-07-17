# Core API

The `pkg/nuget` package provides the main API interface for the NuGet Config Parser library. This is the primary entry point for most operations.

## Overview

The Core API is designed to provide a simple, unified interface for all NuGet configuration operations. It abstracts away the complexity of the underlying components and provides a clean, easy-to-use API for developers.

## API Structure

```go
type API struct {
    Parser  *parser.ConfigParser
    Finder  *finder.ConfigFinder
    Manager *manager.ConfigManager
}
```

The API struct integrates three core components:
- **Parser**: Handles configuration file parsing and serialization
- **Finder**: Locates configuration files across different platforms
- **Manager**: Manages configuration modifications and operations

## Constructor

### NewAPI

```go
func NewAPI() *API
```

Creates a new API instance with default settings.

**Returns:**
- `*API`: A new API instance ready for use

**Example:**
```go
api := nuget.NewAPI()
```

## Parsing Methods

### ParseFromFile

```go
func (a *API) ParseFromFile(filePath string) (*types.NuGetConfig, error)
```

Parses a NuGet configuration file from the specified path.

**Parameters:**
- `filePath` (string): Path to the configuration file

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Example:**
```go
config, err := api.ParseFromFile("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to parse config: %v", err)
}
```

### ParseFromString

```go
func (a *API) ParseFromString(content string) (*types.NuGetConfig, error)
```

Parses a NuGet configuration from an XML string.

**Parameters:**
- `content` (string): XML content as string

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Example:**
```go
xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`

config, err := api.ParseFromString(xmlContent)
```

### ParseFromReader

```go
func (a *API) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)
```

Parses a NuGet configuration from an io.Reader.

**Parameters:**
- `reader` (io.Reader): Reader containing XML content

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Example:**
```go
file, err := os.Open("NuGet.Config")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

config, err := api.ParseFromReader(file)
```

### ParseFromFileWithPositions

```go
func (a *API) ParseFromFileWithPositions(filePath string) (*parser.ParseResult, error)
```

Parses a configuration file while tracking element positions for position-aware editing.

**Parameters:**
- `filePath` (string): Path to the configuration file

**Returns:**
- `*parser.ParseResult`: Parse result with position information
- `error`: Error if parsing fails

**Example:**
```go
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to parse with positions: %v", err)
}

editor := api.CreateConfigEditor(parseResult)
```

## Finding Methods

### FindConfigFile

```go
func (a *API) FindConfigFile() (string, error)
```

Finds the first available NuGet configuration file in the system.

**Returns:**
- `string`: Path to the found configuration file
- `error`: Error if no configuration file is found

**Example:**
```go
configPath, err := api.FindConfigFile()
if err != nil {
    log.Fatalf("No config file found: %v", err)
}
fmt.Printf("Found config: %s\n", configPath)
```

### FindAllConfigFiles

```go
func (a *API) FindAllConfigFiles() []string
```

Finds all available NuGet configuration files in the system.

**Returns:**
- `[]string`: Slice of paths to all found configuration files

**Example:**
```go
configPaths := api.FindAllConfigFiles()
fmt.Printf("Found %d config files:\n", len(configPaths))
for _, path := range configPaths {
    fmt.Printf("  - %s\n", path)
}
```

### FindProjectConfig

```go
func (a *API) FindProjectConfig(startDir string) (string, error)
```

Finds a project-level configuration file starting from the specified directory.

**Parameters:**
- `startDir` (string): Starting directory for the search

**Returns:**
- `string`: Path to the found project configuration file
- `error`: Error if no project configuration file is found

**Example:**
```go
projectConfig, err := api.FindProjectConfig("./my-project")
if err != nil {
    log.Printf("No project config found: %v", err)
} else {
    fmt.Printf("Project config: %s\n", projectConfig)
}
```

### FindAndParseConfig

```go
func (a *API) FindAndParseConfig() (*types.NuGetConfig, string, error)
```

Finds and parses the first available configuration file.

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `string`: Path to the configuration file that was parsed
- `error`: Error if no configuration file is found or parsing fails

**Example:**
```go
config, configPath, err := api.FindAndParseConfig()
if err != nil {
    log.Fatalf("Failed to find and parse config: %v", err)
}
fmt.Printf("Loaded config from: %s\n", configPath)
```

## Package Source Management

### AddPackageSource

```go
func (a *API) AddPackageSource(config *types.NuGetConfig, key, value, protocolVersion string)
```

Adds or updates a package source in the configuration.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Unique identifier for the package source
- `value` (string): URL or path to the package source
- `protocolVersion` (string): Protocol version (optional, can be empty)

**Example:**
```go
api.AddPackageSource(config, "myFeed", "https://my-nuget-feed.com/v3/index.json", "3")
```

### RemovePackageSource

```go
func (a *API) RemovePackageSource(config *types.NuGetConfig, key string) bool
```

Removes a package source from the configuration.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the package source to remove

**Returns:**
- `bool`: True if the source was found and removed, false otherwise

**Example:**
```go
removed := api.RemovePackageSource(config, "myFeed")
if removed {
    fmt.Println("Package source removed successfully")
}
```

### GetPackageSource

```go
func (a *API) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource
```

Retrieves a specific package source by key.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to search
- `key` (string): Key of the package source to retrieve

**Returns:**
- `*types.PackageSource`: Package source if found, nil otherwise

**Example:**
```go
source := api.GetPackageSource(config, "nuget.org")
if source != nil {
    fmt.Printf("Source URL: %s\n", source.Value)
}
```

### GetAllPackageSources

```go
func (a *API) GetAllPackageSources(config *types.NuGetConfig) []types.PackageSource
```

Retrieves all package sources from the configuration.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object

**Returns:**
- `[]types.PackageSource`: Slice of all package sources

**Example:**
```go
sources := api.GetAllPackageSources(config)
for _, source := range sources {
    fmt.Printf("- %s: %s\n", source.Key, source.Value)
}
```

## Package Source Status Management

### EnablePackageSource

```go
func (a *API) EnablePackageSource(config *types.NuGetConfig, key string)
```

Enables a package source by removing it from the disabled sources list.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the package source to enable

### DisablePackageSource

```go
func (a *API) DisablePackageSource(config *types.NuGetConfig, key string)
```

Disables a package source by adding it to the disabled sources list.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the package source to disable

### IsPackageSourceDisabled

```go
func (a *API) IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool
```

Checks if a package source is disabled.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to check
- `key` (string): Key of the package source to check

**Returns:**
- `bool`: True if the source is disabled, false otherwise

**Example:**
```go
if api.IsPackageSourceDisabled(config, "myFeed") {
    fmt.Println("myFeed is disabled")
    api.EnablePackageSource(config, "myFeed")
}
```

## Credential Management

### AddCredential

```go
func (a *API) AddCredential(config *types.NuGetConfig, sourceKey, username, password string)
```

Adds credentials for a package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `sourceKey` (string): Key of the package source
- `username` (string): Username for authentication
- `password` (string): Password for authentication

**Example:**
```go
api.AddCredential(config, "privateFeed", "myuser", "mypassword")
```

### RemoveCredential

```go
func (a *API) RemoveCredential(config *types.NuGetConfig, sourceKey string) bool
```

Removes credentials for a package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `sourceKey` (string): Key of the package source

**Returns:**
- `bool`: True if credentials were found and removed, false otherwise

### GetCredential

```go
func (a *API) GetCredential(config *types.NuGetConfig, sourceKey string) *types.SourceCredential
```

Retrieves credentials for a package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to search
- `sourceKey` (string): Key of the package source

**Returns:**
- `*types.SourceCredential`: Credentials if found, nil otherwise

## Configuration Options

### AddConfigOption

```go
func (a *API) AddConfigOption(config *types.NuGetConfig, key, value string)
```

Adds or updates a global configuration option.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Configuration option key
- `value` (string): Configuration option value

**Example:**
```go
api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages/path")
```

### RemoveConfigOption

```go
func (a *API) RemoveConfigOption(config *types.NuGetConfig, key string) bool
```

Removes a global configuration option.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Configuration option key to remove

**Returns:**
- `bool`: True if the option was found and removed, false otherwise

### GetConfigOption

```go
func (a *API) GetConfigOption(config *types.NuGetConfig, key string) string
```

Retrieves a global configuration option value.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to search
- `key` (string): Configuration option key

**Returns:**
- `string`: Configuration option value, empty string if not found

## Active Package Source

### SetActivePackageSource

```go
func (a *API) SetActivePackageSource(config *types.NuGetConfig, key, value string)
```

Sets the active package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the active package source
- `value` (string): URL of the active package source

### GetActivePackageSource

```go
func (a *API) GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource
```

Retrieves the active package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to search

**Returns:**
- `*types.PackageSource`: Active package source if set, nil otherwise

## Serialization and Persistence

### SaveConfig

```go
func (a *API) SaveConfig(config *types.NuGetConfig, filePath string) error
```

Saves a configuration object to a file.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to save
- `filePath` (string): Path where to save the configuration file

**Returns:**
- `error`: Error if saving fails

**Example:**
```go
err := api.SaveConfig(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

### SerializeToXML

```go
func (a *API) SerializeToXML(config *types.NuGetConfig) (string, error)
```

Serializes a configuration object to XML string.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to serialize

**Returns:**
- `string`: XML representation of the configuration
- `error`: Error if serialization fails

**Example:**
```go
xmlString, err := api.SerializeToXML(config)
if err != nil {
    log.Fatalf("Failed to serialize config: %v", err)
}
fmt.Println(xmlString)
```

## Configuration Creation

### CreateDefaultConfig

```go
func (a *API) CreateDefaultConfig() *types.NuGetConfig
```

Creates a new configuration with default settings.

**Returns:**
- `*types.NuGetConfig`: New configuration with default package source

**Example:**
```go
config := api.CreateDefaultConfig()
// config now contains nuget.org as the default source
```

### InitializeDefaultConfig

```go
func (a *API) InitializeDefaultConfig(filePath string) error
```

Creates and saves a default configuration to the specified path.

**Parameters:**
- `filePath` (string): Path where to create the configuration file

**Returns:**
- `error`: Error if creation or saving fails

**Example:**
```go
err := api.InitializeDefaultConfig("/path/to/new/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to initialize config: %v", err)
}
```

## Position-Aware Editing

### CreateConfigEditor

```go
func (a *API) CreateConfigEditor(parseResult *parser.ParseResult) *editor.ConfigEditor
```

Creates a position-aware configuration editor.

**Parameters:**
- `parseResult` (*parser.ParseResult): Parse result with position information

**Returns:**
- `*editor.ConfigEditor`: Configuration editor instance

**Example:**
```go
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatal(err)
}

editor := api.CreateConfigEditor(parseResult)
err = editor.AddPackageSource("newSource", "https://example.com", "3")
if err != nil {
    log.Fatal(err)
}

modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatal(err)
}

// Save the modified content
err = os.WriteFile("/path/to/NuGet.Config", modifiedContent, 0644)
```

## Error Handling

All methods that can fail return an error as the last return value. Use the error handling utilities from the `pkg/errors` package:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("config.xml")
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle file not found
    } else if errors.IsParseError(err) {
        // Handle parsing error
    } else {
        // Handle other errors
    }
}
```

## Best Practices

1. **Reuse API instances**: Create one API instance and reuse it throughout your application
2. **Check errors**: Always check and handle errors appropriately
3. **Use position-aware editing**: For minimal file changes, use the position-aware editing features
4. **Validate inputs**: Ensure package source keys are unique and URLs are valid
5. **Handle missing files**: Use `FindConfigFile()` or create default configurations when files don't exist
