# Basic Parsing

This example demonstrates the fundamental operations of parsing NuGet configuration files using the NuGet Config Parser library.

## Overview

Basic parsing involves:
- Reading configuration files from various sources
- Handling different file formats and locations
- Displaying configuration contents
- Basic error handling

## Example 1: Parse from File

The most common scenario is parsing an existing configuration file:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    // Create API instance
    api := nuget.NewAPI()
    
    // Parse configuration from file
    configPath := "NuGet.Config"
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Printf("Configuration file not found: %s\n", configPath)
            return
        }
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Display basic information
    fmt.Printf("Configuration loaded from: %s\n", configPath)
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
    
    // List all package sources
    fmt.Println("\nPackage Sources:")
    for i, source := range config.PackageSources.Add {
        fmt.Printf("%d. %s\n", i+1, source.Key)
        fmt.Printf("   URL: %s\n", source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf("   Protocol: v%s\n", source.ProtocolVersion)
        }
        fmt.Println()
    }
}
```

## Example 2: Parse from String

Sometimes you need to parse configuration from a string (e.g., from a database or API):

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // XML configuration as string
    xmlConfig := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="/path/to/local/packages" />
  </packageSources>
  <config>
    <add key="globalPackagesFolder" value="/custom/packages" />
  </config>
</configuration>`
    
    // Parse from string
    config, err := api.ParseFromString(xmlConfig)
    if err != nil {
        log.Fatalf("Failed to parse XML: %v", err)
    }
    
    fmt.Println("Configuration parsed from string successfully!")
    
    // Display package sources
    fmt.Printf("Found %d package sources:\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
    }
    
    // Display config options
    if config.Config != nil {
        fmt.Printf("\nConfiguration options (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("- %s: %s\n", option.Key, option.Value)
        }
    }
}
```

## Example 3: Parse from Reader

For streaming scenarios or when working with io.Reader:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Example 1: Parse from file using Reader
    file, err := os.Open("NuGet.Config")
    if err != nil {
        log.Printf("Could not open file: %v", err)
        
        // Example 2: Parse from string reader instead
        xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="example" value="https://example.com/nuget" />
  </packageSources>
</configuration>`
        
        reader := strings.NewReader(xmlContent)
        config, err := api.ParseFromReader(reader)
        if err != nil {
            log.Fatalf("Failed to parse from reader: %v", err)
        }
        
        fmt.Println("Parsed from string reader:")
        displayConfig(config)
        return
    }
    defer file.Close()
    
    // Parse from file reader
    config, err := api.ParseFromReader(file)
    if err != nil {
        log.Fatalf("Failed to parse from file reader: %v", err)
    }
    
    fmt.Println("Parsed from file reader:")
    displayConfig(config)
}

func displayConfig(config *types.NuGetConfig) {
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
}
```

## Example 4: Comprehensive Parsing with Error Handling

A robust example that handles various error conditions:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
    api := nuget.NewAPI()
    
    // Try to parse configuration with comprehensive error handling
    config, err := parseConfigSafely(api, "NuGet.Config")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Display comprehensive configuration information
    displayFullConfig(config)
}

func parseConfigSafely(api *nuget.API, configPath string) (*types.NuGetConfig, error) {
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Printf("Configuration file not found: %s\n", configPath)
            fmt.Println("Creating default configuration...")
            return api.CreateDefaultConfig(), nil
        }
        
        if errors.IsParseError(err) {
            return nil, fmt.Errorf("invalid configuration format: %w", err)
        }
        
        if errors.IsFormatError(err) {
            return nil, fmt.Errorf("malformed XML in configuration: %w", err)
        }
        
        return nil, fmt.Errorf("unexpected error: %w", err)
    }
    
    return config, nil
}

func displayFullConfig(config *types.NuGetConfig) {
    fmt.Println("=== NuGet Configuration ===")
    
    // Package Sources
    fmt.Printf("\nPackage Sources (%d):\n", len(config.PackageSources.Add))
    if len(config.PackageSources.Add) == 0 {
        fmt.Println("  (none)")
    } else {
        for i, source := range config.PackageSources.Add {
            fmt.Printf("%d. Key: %s\n", i+1, source.Key)
            fmt.Printf("   Value: %s\n", source.Value)
            if source.ProtocolVersion != "" {
                fmt.Printf("   Protocol Version: %s\n", source.ProtocolVersion)
            }
            fmt.Println()
        }
    }
    
    // Disabled Package Sources
    if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
        fmt.Printf("Disabled Package Sources (%d):\n", len(config.DisabledPackageSources.Add))
        for _, disabled := range config.DisabledPackageSources.Add {
            fmt.Printf("  - %s\n", disabled.Key)
        }
        fmt.Println()
    }
    
    // Active Package Source
    if config.ActivePackageSource != nil {
        fmt.Printf("Active Package Source:\n")
        fmt.Printf("  Key: %s\n", config.ActivePackageSource.Add.Key)
        fmt.Printf("  Value: %s\n", config.ActivePackageSource.Add.Value)
        fmt.Println()
    }
    
    // Package Source Credentials
    if config.PackageSourceCredentials != nil && len(config.PackageSourceCredentials.Sources) > 0 {
        fmt.Printf("Package Source Credentials (%d sources):\n", len(config.PackageSourceCredentials.Sources))
        for sourceKey := range config.PackageSourceCredentials.Sources {
            fmt.Printf("  - %s (credentials configured)\n", sourceKey)
        }
        fmt.Println()
    }
    
    // Configuration Options
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("Configuration Options (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
        fmt.Println()
    }
}
```

## Example 5: Parsing Multiple Files

Sometimes you need to parse multiple configuration files:

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    api := nuget.NewAPI()
    
    // List of configuration files to try
    configFiles := []string{
        "NuGet.Config",
        "nuget.config",
        filepath.Join("config", "NuGet.Config"),
        filepath.Join(os.Getenv("HOME"), ".nuget", "NuGet.Config"),
    }
    
    fmt.Println("Attempting to parse multiple configuration files...")
    
    var successfulConfigs []ConfigInfo
    
    for _, configPath := range configFiles {
        fmt.Printf("\nTrying: %s\n", configPath)
        
        config, err := api.ParseFromFile(configPath)
        if err != nil {
            if errors.IsNotFoundError(err) {
                fmt.Printf("  ❌ File not found\n")
            } else {
                fmt.Printf("  ❌ Parse error: %v\n", err)
            }
            continue
        }
        
        fmt.Printf("  ✅ Successfully parsed\n")
        fmt.Printf("  📦 Package sources: %d\n", len(config.PackageSources.Add))
        
        successfulConfigs = append(successfulConfigs, ConfigInfo{
            Path:   configPath,
            Config: config,
        })
    }
    
    // Summary
    fmt.Printf("\n=== Summary ===\n")
    fmt.Printf("Successfully parsed %d configuration files:\n", len(successfulConfigs))
    
    for i, info := range successfulConfigs {
        fmt.Printf("%d. %s (%d sources)\n", i+1, info.Path, len(info.Config.PackageSources.Add))
    }
    
    if len(successfulConfigs) == 0 {
        fmt.Println("No configuration files could be parsed.")
        fmt.Println("Consider creating a default configuration.")
    }
}

type ConfigInfo struct {
    Path   string
    Config *types.NuGetConfig
}
```

## Key Concepts

### Error Types

The library provides specific error types for different scenarios:

- `IsNotFoundError()`: Configuration file doesn't exist
- `IsParseError()`: Invalid XML or parsing issues
- `IsFormatError()`: Malformed configuration structure

### Configuration Structure

A parsed configuration contains:

- **PackageSources**: List of available package sources
- **DisabledPackageSources**: Sources that are disabled
- **ActivePackageSource**: Currently active source
- **PackageSourceCredentials**: Authentication information
- **Config**: Global configuration options

### Best Practices

1. **Always handle errors**: Check for specific error types
2. **Provide fallbacks**: Create default configuration when files are missing
3. **Validate input**: Ensure file paths and content are valid
4. **Display meaningful information**: Show users what was parsed
5. **Use appropriate parsing method**: Choose between file, string, or reader based on your source

## Next Steps

After mastering basic parsing:

1. Learn about [Finding Configs](./finding-configs.md) to locate configuration files
2. Explore [Creating Configs](./creating-configs.md) to generate new configurations
3. Study [Modifying Configs](./modifying-configs.md) to update existing configurations

## Common Issues

### Issue 1: File Not Found

```go
// Always check if file exists before parsing
if !utils.FileExists(configPath) {
    fmt.Printf("Configuration file does not exist: %s\n", configPath)
    // Create default or handle appropriately
}
```

### Issue 2: Invalid XML

```go
// Handle parse errors gracefully
if errors.IsParseError(err) {
    fmt.Printf("Invalid XML format: %v\n", err)
    // Consider creating a new configuration
}
```

### Issue 3: Empty Configuration

```go
// Check if configuration has content
if len(config.PackageSources.Add) == 0 {
    fmt.Println("Configuration has no package sources")
    // Add default sources if needed
}
```

This basic parsing guide provides the foundation for working with NuGet configuration files. The examples demonstrate various parsing scenarios and proper error handling techniques.
