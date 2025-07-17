# Manager API

The `pkg/manager` package provides high-level configuration management operations, combining parsing, finding, and modification capabilities.

## Overview

The Manager API is responsible for:
- High-level configuration operations
- Combining parser and finder functionality
- Managing configuration lifecycle
- Providing convenient methods for common operations

## Types

### ConfigManager

```go
type ConfigManager struct {
    parser *parser.ConfigParser
    finder *finder.ConfigFinder
}
```

The main manager type that orchestrates configuration operations.

**Fields:**
- `parser`: Internal configuration parser
- `finder`: Internal configuration finder

## Constructors

### NewConfigManager

```go
func NewConfigManager() *ConfigManager
```

Creates a new configuration manager with default settings.

**Returns:**
- `*ConfigManager`: New manager instance

**Example:**
```go
manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
fmt.Printf("Loaded config from: %s\n", configPath)
```

## Configuration Loading

### LoadConfig

```go
func (m *ConfigManager) LoadConfig(filePath string) (*types.NuGetConfig, error)
```

Loads a configuration file from the specified path.

**Parameters:**
- `filePath` (string): Path to the configuration file

**Returns:**
- `*types.NuGetConfig`: Loaded configuration object
- `error`: Error if loading fails

**Example:**
```go
manager := manager.NewConfigManager()
config, err := manager.LoadConfig("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

fmt.Printf("Loaded %d package sources\n", len(config.PackageSources.Add))
```

### FindAndLoadConfig

```go
func (m *ConfigManager) FindAndLoadConfig() (*types.NuGetConfig, string, error)
```

Finds and loads the first available configuration file.

**Returns:**
- `*types.NuGetConfig`: Loaded configuration object
- `string`: Path to the configuration file that was loaded
- `error`: Error if no configuration found or loading fails

**Example:**
```go
manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    if errors.IsNotFoundError(err) {
        // No config found, create default
        config = manager.CreateDefaultConfig()
        configPath = "NuGet.Config"
        err = manager.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("Failed to create default config: %v", err)
        }
        fmt.Printf("Created default config: %s\n", configPath)
    } else {
        log.Fatalf("Failed to load config: %v", err)
    }
} else {
    fmt.Printf("Loaded existing config: %s\n", configPath)
}
```

## Configuration Saving

### SaveConfig

```go
func (m *ConfigManager) SaveConfig(config *types.NuGetConfig, filePath string) error
```

Saves a configuration object to the specified file path.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to save
- `filePath` (string): Target file path

**Returns:**
- `error`: Error if saving fails

**Example:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// Modify the configuration
manager.AddPackageSource(config, "company", "https://nuget.company.com", "3")

// Save the configuration
err := manager.SaveConfig(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

## Configuration Creation

### CreateDefaultConfig

```go
func (m *ConfigManager) CreateDefaultConfig() *types.NuGetConfig
```

Creates a new configuration with default settings.

**Returns:**
- `*types.NuGetConfig`: New configuration with default package source

**Default Configuration:**
- Package source: `nuget.org` pointing to `https://api.nuget.org/v3/index.json`
- Protocol version: `3`
- Active package source: Set to the default source

**Example:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

fmt.Printf("Default config has %d package sources\n", len(config.PackageSources.Add))
fmt.Printf("Default source: %s -> %s\n", 
    config.PackageSources.Add[0].Key, 
    config.PackageSources.Add[0].Value)
```

### InitializeDefaultConfig

```go
func (m *ConfigManager) InitializeDefaultConfig(filePath string) error
```

Creates and saves a default configuration to the specified path.

**Parameters:**
- `filePath` (string): Path where to create the configuration file

**Returns:**
- `error`: Error if creation or saving fails

**Features:**
- Creates parent directories if they don't exist
- Generates default configuration
- Saves to the specified path

**Example:**
```go
manager := manager.NewConfigManager()

configPath := "/path/to/new/NuGet.Config"
err := manager.InitializeDefaultConfig(configPath)
if err != nil {
    log.Fatalf("Failed to initialize config: %v", err)
}

fmt.Printf("Initialized default config at: %s\n", configPath)
```

## Package Source Management

### AddPackageSource

```go
func (m *ConfigManager) AddPackageSource(config *types.NuGetConfig, key string, value string, protocolVersion string) 
```

Adds or updates a package source in the configuration.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Unique identifier for the package source
- `value` (string): URL or path to the package source
- `protocolVersion` (string): Protocol version (optional, can be empty)

**Behavior:**
- If a source with the same key exists, it updates the existing source
- If no source exists with the key, it adds a new source
- Protocol version is optional and can be empty

**Example:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// Add a company package source
manager.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")

// Add a local package source without protocol version
manager.AddPackageSource(config, "local", "/path/to/local/packages", "")

// Update an existing source
manager.AddPackageSource(config, "company", "https://new-nuget.company.com/v3/index.json", "3")
```

### RemovePackageSource

```go
func (m *ConfigManager) RemovePackageSource(config *types.NuGetConfig, key string) bool
```

Removes a package source from the configuration.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the package source to remove

**Returns:**
- `bool`: True if the source was found and removed, false otherwise

**Example:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

// Remove a package source
removed := manager.RemovePackageSource(config, "old-source")
if removed {
    fmt.Println("Package source removed successfully")
    
    // Save the updated configuration
    err := manager.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Printf("Failed to save config: %v", err)
    }
} else {
    fmt.Println("Package source not found")
}
```

### GetPackageSource

```go
func (m *ConfigManager) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource
```

Retrieves a specific package source by key.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to search
- `key` (string): Key of the package source to retrieve

**Returns:**
- `*types.PackageSource`: Package source if found, nil otherwise

**Example:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

source := manager.GetPackageSource(config, "nuget.org")
if source != nil {
    fmt.Printf("Source: %s -> %s\n", source.Key, source.Value)
    if source.ProtocolVersion != "" {
        fmt.Printf("Protocol Version: %s\n", source.ProtocolVersion)
    }
} else {
    fmt.Println("Source not found")
}
```

## Active Package Source Management

### SetActivePackageSource

```go
func (m *ConfigManager) SetActivePackageSource(config *types.NuGetConfig, key string) error
```

Sets the active package source by key.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to modify
- `key` (string): Key of the package source to set as active

**Returns:**
- `error`: Error if the package source is not found

**Example:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// Add multiple sources
manager.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
manager.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")

// Set active source
err := manager.SetActivePackageSource(config, "company")
if err != nil {
    log.Printf("Failed to set active source: %v", err)
} else {
    fmt.Println("Active source set to 'company'")
}
```

### GetActivePackageSource

```go
func (m *ConfigManager) GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource
```

Gets the currently active package source.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration object to query

**Returns:**
- `*types.PackageSource`: Active package source if set, nil otherwise

**Example:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

activeSource := manager.GetActivePackageSource(config)
if activeSource != nil {
    fmt.Printf("Active source: %s -> %s\n", activeSource.Key, activeSource.Value)
} else {
    fmt.Println("No active source set")
}
```

## Complete Example

Here's a comprehensive example showing various manager operations:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/manager"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    // Create manager
    mgr := manager.NewConfigManager()
    
    // Try to find and load existing configuration
    config, configPath, err := mgr.FindAndLoadConfig()
    if err != nil {
        if errors.IsNotFoundError(err) {
            // No config found, create default
            fmt.Println("No configuration found, creating default...")
            config = mgr.CreateDefaultConfig()
            configPath = "NuGet.Config"
        } else {
            log.Fatalf("Failed to load config: %v", err)
        }
    } else {
        fmt.Printf("Loaded configuration from: %s\n", configPath)
    }
    
    // Display current package sources
    fmt.Printf("\nCurrent package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // Add a company package source
    fmt.Println("\nAdding company package source...")
    mgr.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    
    // Add a local development source
    fmt.Println("Adding local development source...")
    mgr.AddPackageSource(config, "local-dev", "/tmp/local-packages", "")
    
    // Set active package source
    fmt.Println("Setting active package source...")
    err = mgr.SetActivePackageSource(config, "nuget.org")
    if err != nil {
        log.Printf("Failed to set active source: %v", err)
    }
    
    // Display updated configuration
    fmt.Printf("\nUpdated package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // Show active source
    if activeSource := mgr.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("\nActive source: %s\n", activeSource.Key)
    }
    
    // Save the configuration
    fmt.Printf("\nSaving configuration to: %s\n", configPath)
    err = mgr.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save configuration: %v", err)
    }
    
    fmt.Println("Configuration saved successfully!")
}
```

## Error Handling

The manager uses standard error types from the `pkg/errors` package:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle missing configuration
        config = manager.CreateDefaultConfig()
        // ... create and save default config
    } else if errors.IsParseError(err) {
        // Handle parsing errors
        log.Printf("Parse error: %v", err)
    } else {
        // Handle other errors
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Best Practices

1. **Use manager for high-level operations**: The manager provides convenient methods for common scenarios
2. **Handle missing configurations**: Always check for `IsNotFoundError` and provide defaults
3. **Save after modifications**: Remember to save the configuration after making changes
4. **Validate package sources**: Ensure URLs are valid before adding package sources
5. **Use meaningful keys**: Choose descriptive keys for package sources
6. **Set active sources**: Consider setting an active package source for better user experience

## Thread Safety

The ConfigManager is not thread-safe. Create separate instances for concurrent use or provide appropriate synchronization.
