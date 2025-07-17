# Getting Started

Welcome to NuGet Config Parser! This guide will help you get up and running with the library quickly.

## What is NuGet Config Parser?

NuGet Config Parser is a Go library that provides comprehensive functionality for parsing and manipulating NuGet configuration files (NuGet.Config). It allows you to:

- Parse existing NuGet configuration files
- Create new configurations programmatically
- Modify package sources, credentials, and settings
- Find configuration files in your system
- Preserve original formatting when editing files

## Prerequisites

- Go 1.19 or later
- Basic understanding of NuGet configuration files

## Installation

Add the library to your Go project:

```bash
go get github.com/scagogogo/nuget-config-parser
```

## Your First Program

Let's create a simple program that finds and displays NuGet configuration information:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    // Create a new API instance
    api := nuget.NewAPI()
    
    // Find the first available configuration file
    configPath, err := api.FindConfigFile()
    if err != nil {
        log.Fatalf("No configuration file found: %v", err)
    }
    
    // Parse the configuration file
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("Failed to parse configuration: %v", err)
    }
    
    // Display basic information
    fmt.Printf("Configuration file: %s\n", configPath)
    fmt.Printf("Number of package sources: %d\n", len(config.PackageSources.Add))
    
    // List all package sources
    fmt.Println("\nPackage Sources:")
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
        
        // Check if source is disabled
        if api.IsPackageSourceDisabled(config, source.Key) {
            fmt.Printf("    Status: Disabled\n")
        } else {
            fmt.Printf("    Status: Enabled\n")
        }
    }
    
    // Display active package source
    if config.ActivePackageSource != nil {
        fmt.Printf("\nActive Package Source: %s\n", config.ActivePackageSource.Add.Key)
    }
}
```

## Core Concepts

### API Instance

The main entry point is the `API` struct, which provides all the functionality you need:

```go
api := nuget.NewAPI()
```

### Configuration Object

The `NuGetConfig` struct represents a complete NuGet configuration:

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

### Package Sources

Package sources are the core of NuGet configuration. Each source has:

- **Key**: A unique identifier for the source
- **Value**: The URL or path to the package source
- **ProtocolVersion**: The NuGet protocol version (optional)

## Common Operations

### Finding Configuration Files

```go
// Find the first available configuration file
configPath, err := api.FindConfigFile()

// Find all configuration files
configPaths := api.FindAllConfigFiles()

// Find project-specific configuration
projectConfig, err := api.FindProjectConfig("./my-project")
```

### Parsing Configuration

```go
// Parse from file
config, err := api.ParseFromFile("/path/to/NuGet.Config")

// Parse from string
config, err := api.ParseFromString(xmlContent)

// Parse from io.Reader
config, err := api.ParseFromReader(reader)
```

### Modifying Configuration

```go
// Add a package source
api.AddPackageSource(config, "mySource", "https://my-nuget-feed.com/v3/index.json", "3")

// Remove a package source
removed := api.RemovePackageSource(config, "mySource")

// Disable a package source
api.DisablePackageSource(config, "mySource")

// Add credentials
api.AddCredential(config, "mySource", "username", "password")
```

### Saving Configuration

```go
// Save to file
err := api.SaveConfig(config, "/path/to/NuGet.Config")

// Serialize to XML string
xmlString, err := api.SerializeToXML(config)
```

## Next Steps

Now that you understand the basics, explore these topics:

- [Installation Guide](./installation.md) - Detailed installation instructions
- [Quick Start](./quick-start.md) - More comprehensive examples
- [Configuration](./configuration.md) - Understanding NuGet configuration structure
- [Position-Aware Editing](./position-aware-editing.md) - Advanced editing features
- [API Reference](/api/) - Complete API documentation
- [Examples](/examples/) - Real-world usage examples

## Need Help?

- Check the [API Reference](/api/) for detailed method documentation
- Browse the [Examples](/examples/) for common use cases
- Visit the [GitHub repository](https://github.com/scagogogo/nuget-config-parser) for issues and discussions
