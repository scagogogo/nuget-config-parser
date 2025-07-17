# Examples

This section provides comprehensive examples demonstrating how to use the NuGet Config Parser library in various scenarios. Each example includes complete, runnable code with explanations.

## Overview

The examples are organized by functionality and complexity:

- **[Basic Parsing](./basic-parsing.md)** - Simple configuration file parsing
- **[Finding Configs](./finding-configs.md)** - Locating configuration files
- **[Creating Configs](./creating-configs.md)** - Creating new configurations
- **[Modifying Configs](./modifying-configs.md)** - Modifying existing configurations
- **[Package Sources](./package-sources.md)** - Managing package sources
- **[Credentials](./credentials.md)** - Handling authentication
- **[Config Options](./config-options.md)** - Global configuration settings
- **[Serialization](./serialization.md)** - Converting to/from XML
- **[Position-Aware Editing](./position-aware-editing.md)** - Advanced editing with minimal diffs

## Quick Start Example

Here's a simple example to get you started:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    // Create API instance
    api := nuget.NewAPI()
    
    // Find and parse the first available configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // If no config found, create a default one
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        
        // Save the default configuration
        err = api.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("Failed to save default config: %v", err)
        }
        
        fmt.Printf("Created default configuration: %s\n", configPath)
    } else {
        fmt.Printf("Found existing configuration: %s\n", configPath)
    }
    
    // Display package sources
    fmt.Printf("Package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        
        fmt.Printf("  - %s: %s (%s)", source.Key, source.Value, status)
        if source.ProtocolVersion != "" {
            fmt.Printf(" [v%s]", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // Add a custom package source
    api.AddPackageSource(config, "example", "https://example.com/nuget", "3")
    
    // Save the updated configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save updated config: %v", err)
    }
    
    fmt.Println("Added example package source and saved configuration")
}
```

## Common Patterns

### Error Handling Pattern

```go
config, err := api.ParseFromFile(configPath)
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle missing file
        config = api.CreateDefaultConfig()
    } else if errors.IsParseError(err) {
        // Handle parsing errors
        log.Fatalf("Invalid configuration format: %v", err)
    } else {
        // Handle other errors
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

### Configuration Discovery Pattern

```go
// Try multiple approaches to find configuration
var config *types.NuGetConfig
var configPath string
var err error

// 1. Try to find and parse existing config
config, configPath, err = api.FindAndParseConfig()
if err == nil {
    fmt.Printf("Using existing config: %s\n", configPath)
} else {
    // 2. Try project-specific config
    configPath, err = api.FindProjectConfig(".")
    if err == nil {
        config, err = api.ParseFromFile(configPath)
        if err == nil {
            fmt.Printf("Using project config: %s\n", configPath)
        }
    }
}

// 3. Fall back to default config
if config == nil {
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
    fmt.Println("Using default configuration")
}
```

### Batch Modification Pattern

```go
// Make multiple changes efficiently
api.AddPackageSource(config, "feed1", "https://feed1.com", "3")
api.AddPackageSource(config, "feed2", "https://feed2.com", "3")
api.AddCredential(config, "feed1", "user1", "pass1")
api.AddCredential(config, "feed2", "user2", "pass2")
api.DisablePackageSource(config, "old-feed")

// Save all changes at once
err := api.SaveConfig(config, configPath)
if err != nil {
    log.Fatalf("Failed to save changes: %v", err)
}
```

### Position-Aware Editing Pattern

```go
// Parse with position tracking
parseResult, err := api.ParseFromFileWithPositions(configPath)
if err != nil {
    log.Fatalf("Failed to parse with positions: %v", err)
}

// Create editor
editor := api.CreateConfigEditor(parseResult)

// Make changes
editor.AddPackageSource("new-feed", "https://new-feed.com", "3")
editor.UpdatePackageSourceURL("existing-feed", "https://updated-url.com")

// Apply changes with minimal diff
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatalf("Failed to apply edits: %v", err)
}

// Save the modified content
err = os.WriteFile(configPath, modifiedContent, 0644)
if err != nil {
    log.Fatalf("Failed to save file: %v", err)
}
```

## Example Categories

### Beginner Examples

Perfect for getting started with the library:

1. **[Basic Parsing](./basic-parsing.md)** - Read and display configuration files
2. **[Finding Configs](./finding-configs.md)** - Locate configuration files in your system
3. **[Creating Configs](./creating-configs.md)** - Create new configuration files

### Intermediate Examples

For common configuration management tasks:

4. **[Modifying Configs](./modifying-configs.md)** - Update existing configurations
5. **[Package Sources](./package-sources.md)** - Manage package sources and their properties
6. **[Credentials](./credentials.md)** - Handle authentication for private feeds

### Advanced Examples

For complex scenarios and optimization:

7. **[Config Options](./config-options.md)** - Manage global NuGet settings
8. **[Serialization](./serialization.md)** - Custom XML handling and validation
9. **[Position-Aware Editing](./position-aware-editing.md)** - Preserve formatting and minimize diffs

## Running the Examples

All examples are designed to be self-contained and runnable. To run an example:

1. Create a new Go file with the example code
2. Initialize a Go module if needed:
   ```bash
   go mod init example
   go get github.com/scagogogo/nuget-config-parser
   ```
3. Run the example:
   ```bash
   go run main.go
   ```

## Example Data

Many examples use sample NuGet.Config files. Here's a typical example configuration:

```xml
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="C:\LocalPackages" />
    <add key="company" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
  </packageSources>
  <packageSourceCredentials>
    <company>
      <add key="Username" value="companyuser" />
      <add key="ClearTextPassword" value="companypass" />
    </company>
  </packageSourceCredentials>
  <disabledPackageSources>
    <add key="local" value="true" />
  </disabledPackageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="C:\packages" />
    <add key="repositoryPath" value=".\packages" />
  </config>
</configuration>
```

## Best Practices Demonstrated

The examples demonstrate these best practices:

- **Error Handling**: Proper error checking and handling for different scenarios
- **Resource Management**: Efficient use of API instances and file operations
- **Configuration Validation**: Ensuring configurations are valid before saving
- **Security**: Safe handling of credentials and sensitive information
- **Performance**: Efficient batch operations and minimal file modifications
- **Maintainability**: Clean, readable code with good separation of concerns

## Contributing Examples

If you have a useful example that's not covered here, consider contributing it to the project. Good examples should be:

- **Complete**: Include all necessary imports and error handling
- **Focused**: Demonstrate one specific concept or use case
- **Documented**: Include comments explaining the important parts
- **Tested**: Verify the example works with the current library version

## Next Steps

After reviewing the examples:

1. Check the [API Reference](/api/) for detailed method documentation
2. Read the [Guide](/guide/) for conceptual information
3. Explore the library's source code for advanced usage patterns
4. Consider contributing your own examples or improvements
