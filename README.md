# NuGet Config Parser

[![Go CI](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml)
[![Scheduled Tests](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/nuget-config-parser)](https://goreportcard.com/report/github.com/scagogogo/nuget-config-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/nuget-config-parser?status.svg)](https://godoc.org/github.com/scagogogo/nuget-config-parser)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/nuget-config-parser/)

A comprehensive Go library for parsing and manipulating NuGet configuration files (NuGet.Config). This library helps you read, modify, and create NuGet configuration files in Go applications, supporting all major NuGet configuration features.

**[📖 Online Documentation](https://scagogogo.github.io/nuget-config-parser/)** | **[🇨🇳 中文文档](README_zh.md)**

## 📚 Documentation

Complete documentation is available online at **https://scagogogo.github.io/nuget-config-parser/**

The documentation includes:
- **Getting Started Guide** - Step-by-step introduction
- **API Reference** - Complete API documentation with examples
- **Examples** - Real-world usage examples
- **Best Practices** - Recommended patterns and practices
- **Multi-language Support** - Available in English and Chinese

## 📑 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Examples](#examples)
- [API Reference](#api-reference)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Configuration Parsing** - Parse NuGet.Config files from files, strings, or io.Reader
- **Smart File Discovery** - Find NuGet configuration files in your system, supporting project-level and global configurations
- **Package Source Management** - Add, remove, enable/disable package sources with full protocol version support
- **Credential Management** - Securely manage username/password credentials for private package sources
- **Configuration Options** - Manage global configuration options like proxy settings, package folder paths, etc.
- **Position-Aware Editing** - Edit configuration files while preserving original formatting and minimizing diffs
- **Serialization Support** - Convert configuration objects to standard XML format with proper indentation
- **Cross-Platform** - Full support for Windows, Linux, and macOS with platform-specific configuration paths

## 🚀 Installation

Install using Go modules (recommended):

```bash
go get github.com/scagogogo/nuget-config-parser
```

## 🏁 Quick Start

Here's a simple example demonstrating how to parse and use NuGet configuration files:

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
    
    // Display configuration information
    fmt.Printf("Configuration file: %s\n", configPath)
    fmt.Printf("Contains %d package sources\n", len(config.PackageSources.Add))
    
    // Display package source list
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
        
        // Check if package source is disabled
        if api.IsPackageSourceDisabled(config, source.Key) {
            fmt.Printf("  Status: Disabled\n")
        } else {
            fmt.Printf("  Status: Enabled\n")
        }
    }
}
```

## 📝 Examples

This project provides multiple complete examples demonstrating different features and use cases. All examples are located in the [examples](examples/) directory:

1. **[Basic Parsing](examples/01_basic_parsing)** - Parse configuration files and access their content
2. **[Finding Configs](examples/02_search_config)** - Find NuGet configuration files in your system
3. **[Creating Configs](examples/03_create_config)** - Create new NuGet configurations
4. **[Modifying Configs](examples/04_modify_config)** - Modify existing NuGet configurations
5. **[Package Sources](examples/05_package_sources)** - Package source related operations
6. **[Credentials](examples/06_credentials)** - Manage package source credentials
7. **[Config Options](examples/07_config_options)** - Manage global configuration options
8. **[Serialization](examples/08_serialization)** - Configuration serialization and deserialization
9. **[Position-Aware Editing](examples/09_position_aware_editing)** - Precise editing based on position information

Run examples:

```bash
go run examples/01_basic_parsing/main.go
```

For detailed example descriptions, see [examples/README.md](examples/README.md).

## 📚 API Reference

### Core API

```go
// Create new API instance
api := nuget.NewAPI()
```

### Parsing and Finding

```go
// Parse configuration from file
config, err := api.ParseFromFile(filePath)

// Parse configuration from string
config, err := api.ParseFromString(xmlContent)

// Parse configuration from io.Reader
config, err := api.ParseFromReader(reader)

// Find first available configuration file
configPath, err := api.FindConfigFile()

// Find all available configuration files
configPaths := api.FindAllConfigFiles()

// Find configuration file in project directory
projectConfig, err := api.FindProjectConfig(startDir)

// Find and parse configuration
config, configPath, err := api.FindAndParseConfig()
```

### Package Source Management

```go
// Add or update package source
api.AddPackageSource(config, "sourceName", "https://source-url", "3")

// Remove package source
removed := api.RemovePackageSource(config, "sourceName")

// Get specific package source
source := api.GetPackageSource(config, "sourceName")

// Get all package sources
sources := api.GetAllPackageSources(config)

// Enable/disable package sources
api.EnablePackageSource(config, "sourceName")
api.DisablePackageSource(config, "sourceName")
isDisabled := api.IsPackageSourceDisabled(config, "sourceName")
```

### Credential Management

```go
// Add credentials
api.AddCredential(config, "sourceName", "username", "password")

// Remove credentials
removed := api.RemoveCredential(config, "sourceName")

// Get credentials
credential := api.GetCredential(config, "sourceName")
```

### Configuration Options

```go
// Add configuration option
api.AddConfigOption(config, "globalPackagesFolder", "/custom/path")

// Remove configuration option
removed := api.RemoveConfigOption(config, "globalPackagesFolder")

// Get configuration option
value := api.GetConfigOption(config, "globalPackagesFolder")
```

### Active Package Source

```go
// Set active package source
api.SetActivePackageSource(config, "sourceName", "https://source-url")

// Get active package source
activeSource := api.GetActivePackageSource(config)
```

### Creation and Saving

```go
// Create default configuration
config := api.CreateDefaultConfig()

// Create default configuration at specified path
err := api.InitializeDefaultConfig(filePath)

// Save configuration to file
err := api.SaveConfig(config, filePath)

// Serialize configuration to XML string
xmlString, err := api.SerializeToXML(config)

// Position-aware editing (preserves original formatting, minimizes diff)
parseResult, err := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)
err = editor.AddPackageSource("new-source", "https://example.com/v3/index.json", "3")
modifiedContent, err := editor.ApplyEdits()
```

## 🏗️ Architecture

The library consists of the following main components:

- **pkg/nuget**: Main API package providing the user interface
- **pkg/parser**: Configuration parser responsible for XML parsing
- **pkg/finder**: Configuration finder responsible for locating configuration files
- **pkg/manager**: Configuration manager responsible for modifying configurations
- **pkg/editor**: Position-aware editor for precise configuration editing
- **pkg/types**: Data type definitions
- **pkg/constants**: Constant definitions
- **pkg/utils**: Utility functions
- **pkg/errors**: Error type definitions

## 📖 Documentation

Complete documentation is available online:

**🌐 [Documentation Website](https://scagogogo.github.io/nuget-config-parser/)**

The documentation includes:
- **Getting Started Guide** - Step-by-step introduction
- **API Reference** - Complete API documentation
- **Examples** - Real-world usage examples
- **Best Practices** - Recommended patterns and practices

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on:

- How to report bugs
- How to suggest new features
- How to submit pull requests
- Development setup and guidelines

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the official NuGet configuration system
- Built with Go's excellent standard library
- Thanks to all contributors and users of this library

---

**[📖 Full Documentation](https://scagogogo.github.io/nuget-config-parser/)** | **[🇨🇳 中文版本](README_zh.md)**
