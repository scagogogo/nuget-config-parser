---
layout: home

hero:
  name: "NuGet Config Parser"
  text: "Go Library for NuGet Configuration"
  tagline: "Parse, manipulate, and manage NuGet configuration files with ease"
  image:
    src: /logo.svg
    alt: NuGet Config Parser
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/scagogogo/nuget-config-parser

features:
  - icon: 📄
    title: Configuration Parsing
    details: Parse NuGet.Config files from files, strings, or io.Reader with comprehensive error handling.
  
  - icon: 🔍
    title: Smart File Discovery
    details: Automatically find NuGet configuration files in your system, supporting project-level and global configurations.
  
  - icon: 📦
    title: Package Source Management
    details: Add, remove, enable/disable package sources with full support for protocol versions and credentials.
  
  - icon: 🔐
    title: Credential Management
    details: Securely manage username/password credentials for private package sources.
  
  - icon: ⚙️
    title: Configuration Options
    details: Manage global configuration options like proxy settings, package folder paths, and more.
  
  - icon: ✏️
    title: Position-Aware Editing
    details: Edit configuration files while preserving original formatting and minimizing diffs.
  
  - icon: 🔄
    title: Serialization Support
    details: Convert configuration objects to standard XML format with proper indentation.
  
  - icon: 🌐
    title: Cross-Platform
    details: Full support for Windows, Linux, and macOS with platform-specific configuration paths.
  
  - icon: 🧪
    title: Comprehensive Testing
    details: Extensively tested with high code coverage and real-world scenarios.
---

## Quick Example

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
    
    // Find and parse configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("Failed to find config: %v", err)
    }
    
    // Display package sources
    fmt.Printf("Config file: %s\n", configPath)
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
    
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
    }
}
```

## Installation

```bash
go get github.com/scagogogo/nuget-config-parser
```

## Key Features

### 🚀 Easy to Use
Simple, intuitive API that follows Go best practices and conventions.

### 🔧 Comprehensive
Supports all major NuGet configuration features including package sources, credentials, and global settings.

### 📝 Well Documented
Extensive documentation with examples for every feature and use case.

### 🎯 Production Ready
Battle-tested with comprehensive test coverage and error handling.
