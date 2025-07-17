# Quick Start

This guide will help you get up and running with the NuGet Config Parser library quickly with practical examples.

## Prerequisites

- Go 1.19 or later
- Basic understanding of Go programming
- Familiarity with NuGet configuration files (helpful but not required)

## Installation

First, add the library to your Go project:

```bash
go get github.com/scagogogo/nuget-config-parser
```

## Basic Usage

### 1. Create Your First Program

Create a new Go file and import the library:

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
    
    // Your code here
}
```

### 2. Find and Parse Configuration

The most common operation is finding and parsing an existing configuration:

```go
func main() {
    api := nuget.NewAPI()
    
    // Find and parse in one step
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("Failed to find/parse config: %v", err)
    }
    
    fmt.Printf("Loaded config from: %s\n", configPath)
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
}
```

### 3. Create a New Configuration

If no configuration exists, create a default one:

```go
func main() {
    api := nuget.NewAPI()
    
    // Try to find existing config
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // No config found, create default
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        
        // Save the new configuration
        err = api.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("Failed to save config: %v", err)
        }
        
        fmt.Printf("Created new config: %s\n", configPath)
    } else {
        fmt.Printf("Using existing config: %s\n", configPath)
    }
}
```

### 4. Manage Package Sources

Add, remove, and manage package sources:

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // Add a custom package source
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    
    // Add a local package source
    api.AddPackageSource(config, "local-packages", "/path/to/local/packages", "")
    
    // Disable a package source
    api.DisablePackageSource(config, "local-packages")
    
    // List all package sources
    fmt.Println("Package Sources:")
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // Save changes
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

### 5. Manage Credentials

Add credentials for private package sources:

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // Add a private package source
    api.AddPackageSource(config, "private-feed", "https://private.nuget.com/v3/index.json", "3")
    
    // Add credentials for the private source
    api.AddCredential(config, "private-feed", "myusername", "mypassword")
    
    // Verify credentials were added
    credential := api.GetCredential(config, "private-feed")
    if credential != nil {
        fmt.Println("Credentials added successfully")
    }
    
    // Save configuration
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

### 6. Configuration Options

Manage global NuGet settings:

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // Set global packages folder
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages/path")
    
    // Set default push source
    api.AddConfigOption(config, "defaultPushSource", "https://my-nuget-server.com")
    
    // Set proxy settings
    api.AddConfigOption(config, "http_proxy", "http://proxy.company.com:8080")
    api.AddConfigOption(config, "http_proxy.user", "proxyuser")
    api.AddConfigOption(config, "http_proxy.password", "proxypass")
    
    // Get a configuration option
    packagesPath := api.GetConfigOption(config, "globalPackagesFolder")
    fmt.Printf("Global packages folder: %s\n", packagesPath)
    
    // Save configuration
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

## Complete Example

Here's a complete example that demonstrates multiple features:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Try to find existing configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // Create new configuration if none exists
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        fmt.Println("Created new configuration")
    } else {
        fmt.Printf("Found existing configuration: %s\n", configPath)
    }
    
    // Display current package sources
    fmt.Printf("\nCurrent package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // Add a company package source
    fmt.Println("\nAdding company package source...")
    api.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    api.AddCredential(config, "company", "employee", "secret123")
    
    // Add local development source
    fmt.Println("Adding local development source...")
    api.AddPackageSource(config, "local-dev", "/tmp/local-packages", "")
    api.DisablePackageSource(config, "local-dev") // Disabled by default
    
    // Configure global settings
    fmt.Println("Configuring global settings...")
    api.AddConfigOption(config, "globalPackagesFolder", os.ExpandEnv("$HOME/.nuget/packages"))
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // Set active package source
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // Display updated configuration
    fmt.Printf("\nUpdated package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // Show active source
    if activeSource := api.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("\nActive source: %s\n", activeSource.Key)
    }
    
    // Save the configuration
    fmt.Printf("\nSaving configuration to: %s\n", configPath)
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save configuration: %v", err)
    }
    
    fmt.Println("Configuration saved successfully!")
    
    // Optionally display the XML content
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("Failed to serialize config: %v", err)
    } else {
        fmt.Println("\nGenerated XML:")
        fmt.Println(xmlContent)
    }
}
```

## Error Handling

Always handle errors appropriately:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("NuGet.Config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // File doesn't exist
        config = api.CreateDefaultConfig()
    } else if errors.IsParseError(err) {
        // Invalid XML or format
        log.Fatalf("Invalid configuration format: %v", err)
    } else {
        // Other error
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

## Next Steps

Now that you've learned the basics:

1. Explore [Position-Aware Editing](./position-aware-editing.md) for advanced editing features
2. Check out the [Examples](/examples/) for more specific use cases
3. Read the [API Reference](/api/) for complete documentation
4. Learn about [Configuration](./configuration.md) structure and options

## Common Patterns

### Configuration Discovery Pattern

```go
// Try multiple approaches to find configuration
var config *types.NuGetConfig
var configPath string

// 1. Project-specific config
if projectConfig, err := api.FindProjectConfig("."); err == nil {
    if config, err = api.ParseFromFile(projectConfig); err == nil {
        configPath = projectConfig
    }
}

// 2. Global config
if config == nil {
    if globalConfig, err := api.FindConfigFile(); err == nil {
        if config, err = api.ParseFromFile(globalConfig); err == nil {
            configPath = globalConfig
        }
    }
}

// 3. Default config
if config == nil {
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
}
```

### Batch Operations Pattern

```go
// Make multiple changes efficiently
api.AddPackageSource(config, "feed1", "https://feed1.com", "3")
api.AddPackageSource(config, "feed2", "https://feed2.com", "3")
api.AddCredential(config, "feed1", "user1", "pass1")
api.DisablePackageSource(config, "old-feed")

// Save all changes at once
err := api.SaveConfig(config, configPath)
```

This quick start guide should get you productive with the library quickly. For more advanced features and detailed explanations, continue with the other documentation sections.
