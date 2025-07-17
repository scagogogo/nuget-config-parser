# Package Sources

This example demonstrates comprehensive package source management using the NuGet Config Parser library.

## Overview

Package source management involves:
- Adding and removing package sources
- Enabling and disabling sources
- Managing source priorities and protocols
- Handling different source types (HTTP, local, UNC)
- Setting active package sources

## Example 1: Basic Package Source Operations

Fundamental package source operations:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Start with default configuration
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Basic Package Source Operations ===")
    
    // Add various types of package sources
    fmt.Println("Adding package sources...")
    
    // Official NuGet.org (already exists in default config)
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    
    // Company internal feed
    api.AddPackageSource(config, "company-internal", "https://nuget.company.com/v3/index.json", "3")
    
    // Local file system source
    api.AddPackageSource(config, "local-packages", "/path/to/local/packages", "")
    
    // UNC network path
    api.AddPackageSource(config, "network-share", "\\\\server\\share\\packages", "")
    
    // Legacy V2 API source
    api.AddPackageSource(config, "legacy-feed", "https://legacy.nuget.com/api/v2", "2")
    
    // Display all sources
    fmt.Printf("\nPackage sources (%d):\n", len(config.PackageSources.Add))
    for i, source := range config.PackageSources.Add {
        fmt.Printf("%d. %s\n", i+1, source.Key)
        fmt.Printf("   URL: %s\n", source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf("   Protocol: v%s\n", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // Get specific source
    companySource := api.GetPackageSource(config, "company-internal")
    if companySource != nil {
        fmt.Printf("Found company source: %s\n", companySource.Value)
    }
    
    // Remove a source
    removed := api.RemovePackageSource(config, "legacy-feed")
    if removed {
        fmt.Println("Removed legacy feed")
    }
    
    // Save configuration
    err := api.SaveConfig(config, "PackageSources.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Println("Configuration saved successfully!")
}
```

## Example 2: Source State Management

Managing enabled/disabled states of package sources:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Package Source State Management ===")
    
    // Add multiple sources
    sources := []struct {
        key     string
        url     string
        version string
        enabled bool
    }{
        {"nuget.org", "https://api.nuget.org/v3/index.json", "3", true},
        {"company-stable", "https://stable.company.com/nuget", "3", true},
        {"company-preview", "https://preview.company.com/nuget", "3", false},
        {"local-dev", "./local-packages", "", false},
        {"backup-feed", "https://backup.nuget.com/api/v2", "2", true},
    }
    
    // Add sources and set their states
    for _, source := range sources {
        api.AddPackageSource(config, source.key, source.url, source.version)
        
        if !source.enabled {
            api.DisablePackageSource(config, source.key)
            fmt.Printf("Added and disabled: %s\n", source.key)
        } else {
            fmt.Printf("Added and enabled: %s\n", source.key)
        }
    }
    
    // Display source states
    fmt.Println("\n=== Source States ===")
    for _, source := range config.PackageSources.Add {
        status := "✅ Enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "❌ Disabled"
        }
        fmt.Printf("%s: %s\n", source.Key, status)
    }
    
    // Enable a disabled source
    fmt.Println("\nEnabling preview source...")
    api.EnablePackageSource(config, "company-preview")
    
    // Disable an enabled source temporarily
    fmt.Println("Temporarily disabling backup feed...")
    api.DisablePackageSource(config, "backup-feed")
    
    // Show updated states
    fmt.Println("\n=== Updated Source States ===")
    enabledCount := 0
    disabledCount := 0
    
    for _, source := range config.PackageSources.Add {
        if api.IsPackageSourceDisabled(config, source.Key) {
            fmt.Printf("❌ %s (disabled)\n", source.Key)
            disabledCount++
        } else {
            fmt.Printf("✅ %s (enabled)\n", source.Key)
            enabledCount++
        }
    }
    
    fmt.Printf("\nSummary: %d enabled, %d disabled\n", enabledCount, disabledCount)
    
    // Save configuration
    err := api.SaveConfig(config, "SourceStates.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

## Example 3: Active Package Source Management

Managing the active package source:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Active Package Source Management ===")
    
    // Add multiple sources
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    api.AddPackageSource(config, "local-dev", "./packages", "")
    
    // Check current active source
    activeSource := api.GetActivePackageSource(config)
    if activeSource != nil {
        fmt.Printf("Current active source: %s\n", activeSource.Key)
    } else {
        fmt.Println("No active source set")
    }
    
    // Set different active sources
    fmt.Println("\nSetting active sources...")
    
    // Set nuget.org as active
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    activeSource = api.GetActivePackageSource(config)
    fmt.Printf("Active source: %s -> %s\n", activeSource.Key, activeSource.Value)
    
    // Switch to company feed
    api.SetActivePackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json")
    activeSource = api.GetActivePackageSource(config)
    fmt.Printf("Active source: %s -> %s\n", activeSource.Key, activeSource.Value)
    
    // Demonstrate source switching based on environment
    environment := "development" // This could come from env var
    
    switch environment {
    case "development":
        api.SetActivePackageSource(config, "local-dev", "./packages")
        fmt.Println("Switched to local development source")
    case "staging":
        api.SetActivePackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json")
        fmt.Println("Switched to company staging source")
    case "production":
        api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
        fmt.Println("Switched to production source")
    }
    
    // Display final active source
    activeSource = api.GetActivePackageSource(config)
    if activeSource != nil {
        fmt.Printf("\nFinal active source: %s\n", activeSource.Key)
        fmt.Printf("URL: %s\n", activeSource.Value)
    }
    
    // Save configuration
    err := api.SaveConfig(config, "ActiveSource.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

## Example 4: Source Validation and Health Checking

Validating package sources and checking their health:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "time"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load existing configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Printf("Validating sources in: %s\n", configPath)
    fmt.Println("=== Package Source Validation ===")
    
    // Validate each source
    for _, source := range config.PackageSources.Add {
        fmt.Printf("\nValidating: %s\n", source.Key)
        fmt.Printf("URL: %s\n", source.Value)
        
        // Check source type and validate accordingly
        if isHTTPSource(source.Value) {
            validateHTTPSource(source)
        } else if isLocalPath(source.Value) {
            validateLocalSource(source)
        } else if isUNCPath(source.Value) {
            validateUNCSource(source)
        } else {
            fmt.Printf("❓ Unknown source type\n")
        }
    }
    
    // Recommend optimizations
    fmt.Println("\n=== Optimization Recommendations ===")
    recommendOptimizations(api, config)
}

func isHTTPSource(sourceURL string) bool {
    return len(sourceURL) > 4 && (sourceURL[:4] == "http" || sourceURL[:5] == "https")
}

func isLocalPath(path string) bool {
    return !isHTTPSource(path) && !isUNCPath(path)
}

func isUNCPath(path string) bool {
    return len(path) > 2 && path[:2] == "\\\\"
}

func validateHTTPSource(source types.PackageSource) {
    // Parse URL
    parsedURL, err := url.Parse(source.Value)
    if err != nil {
        fmt.Printf("❌ Invalid URL format: %v\n", err)
        return
    }
    
    fmt.Printf("   Host: %s\n", parsedURL.Host)
    fmt.Printf("   Scheme: %s\n", parsedURL.Scheme)
    
    // Check if URL is reachable
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp, err := client.Head(source.Value)
    if err != nil {
        fmt.Printf("❌ Source unreachable: %v\n", err)
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 200 {
        fmt.Printf("✅ Source is reachable (HTTP %d)\n", resp.StatusCode)
    } else {
        fmt.Printf("⚠️  Source returned HTTP %d\n", resp.StatusCode)
    }
    
    // Check protocol version compatibility
    if source.ProtocolVersion == "3" {
        fmt.Printf("✅ Using modern NuGet v3 protocol\n")
    } else if source.ProtocolVersion == "2" {
        fmt.Printf("⚠️  Using legacy NuGet v2 protocol\n")
    } else {
        fmt.Printf("❓ No protocol version specified\n")
    }
}

func validateLocalSource(source types.PackageSource) {
    // Check if path exists
    if _, err := os.Stat(source.Value); os.IsNotExist(err) {
        fmt.Printf("❌ Local path does not exist: %s\n", source.Value)
        return
    }
    
    // Check if it's a directory
    info, err := os.Stat(source.Value)
    if err != nil {
        fmt.Printf("❌ Cannot access path: %v\n", err)
        return
    }
    
    if !info.IsDir() {
        fmt.Printf("❌ Path is not a directory\n")
        return
    }
    
    fmt.Printf("✅ Local directory exists and is accessible\n")
    
    // Check for .nupkg files
    matches, err := filepath.Glob(filepath.Join(source.Value, "*.nupkg"))
    if err == nil && len(matches) > 0 {
        fmt.Printf("✅ Found %d .nupkg files\n", len(matches))
    } else {
        fmt.Printf("⚠️  No .nupkg files found in directory\n")
    }
}

func validateUNCSource(source types.PackageSource) {
    // Basic UNC path validation
    if len(source.Value) < 5 || source.Value[:2] != "\\\\" {
        fmt.Printf("❌ Invalid UNC path format\n")
        return
    }
    
    // Try to access the UNC path
    if _, err := os.Stat(source.Value); os.IsNotExist(err) {
        fmt.Printf("❌ UNC path not accessible: %s\n", source.Value)
        return
    }
    
    fmt.Printf("✅ UNC path is accessible\n")
}

func recommendOptimizations(api *nuget.API, config *types.NuGetConfig) {
    // Check for disabled sources
    disabledCount := 0
    for _, source := range config.PackageSources.Add {
        if api.IsPackageSourceDisabled(config, source.Key) {
            disabledCount++
        }
    }
    
    if disabledCount > 0 {
        fmt.Printf("💡 Consider removing %d disabled sources to reduce config complexity\n", disabledCount)
    }
    
    // Check for duplicate URLs
    urlMap := make(map[string][]string)
    for _, source := range config.PackageSources.Add {
        urlMap[source.Value] = append(urlMap[source.Value], source.Key)
    }
    
    for url, keys := range urlMap {
        if len(keys) > 1 {
            fmt.Printf("⚠️  Duplicate URL found: %s used by sources: %v\n", url, keys)
        }
    }
    
    // Check protocol versions
    v2Count := 0
    for _, source := range config.PackageSources.Add {
        if source.ProtocolVersion == "2" {
            v2Count++
        }
    }
    
    if v2Count > 0 {
        fmt.Printf("💡 Consider upgrading %d sources from v2 to v3 protocol for better performance\n", v2Count)
    }
    
    // Check for active source
    activeSource := api.GetActivePackageSource(config)
    if activeSource == nil {
        fmt.Printf("💡 Consider setting an active package source for better performance\n")
    }
}
```

## Example 5: Advanced Source Management

Complex source management scenarios:

```go
package main

import (
    "fmt"
    "log"
    "sort"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Advanced Source Management ===")
    
    // Create a source manager
    manager := NewSourceManager(api, config)
    
    // Add sources with priorities
    manager.AddSourceWithPriority("nuget.org", "https://api.nuget.org/v3/index.json", "3", 1)
    manager.AddSourceWithPriority("company-stable", "https://stable.company.com/nuget", "3", 2)
    manager.AddSourceWithPriority("company-preview", "https://preview.company.com/nuget", "3", 3)
    manager.AddSourceWithPriority("local-dev", "./packages", "", 4)
    
    // Display sources by priority
    fmt.Println("Sources by priority:")
    manager.DisplaySourcesByPriority()
    
    // Enable/disable sources based on environment
    environment := "development"
    manager.ConfigureForEnvironment(environment)
    
    // Bulk operations
    fmt.Println("\n=== Bulk Operations ===")
    
    // Disable all preview sources
    previewSources := []string{"company-preview", "nuget-preview", "dotnet-preview"}
    manager.BulkDisable(previewSources)
    
    // Enable production sources
    productionSources := []string{"nuget.org", "company-stable"}
    manager.BulkEnable(productionSources)
    
    // Display final configuration
    manager.DisplaySummary()
    
    // Save configuration
    err := api.SaveConfig(config, "AdvancedSources.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}

type SourceManager struct {
    api       *nuget.API
    config    *types.NuGetConfig
    priorities map[string]int
}

func NewSourceManager(api *nuget.API, config *types.NuGetConfig) *SourceManager {
    return &SourceManager{
        api:       api,
        config:    config,
        priorities: make(map[string]int),
    }
}

func (sm *SourceManager) AddSourceWithPriority(key, url, version string, priority int) {
    sm.api.AddPackageSource(sm.config, key, url, version)
    sm.priorities[key] = priority
    fmt.Printf("Added source: %s (priority %d)\n", key, priority)
}

func (sm *SourceManager) DisplaySourcesByPriority() {
    type SourceWithPriority struct {
        Source   types.PackageSource
        Priority int
    }
    
    var sources []SourceWithPriority
    for _, source := range sm.config.PackageSources.Add {
        priority := sm.priorities[source.Key]
        sources = append(sources, SourceWithPriority{source, priority})
    }
    
    // Sort by priority
    sort.Slice(sources, func(i, j int) bool {
        return sources[i].Priority < sources[j].Priority
    })
    
    for _, s := range sources {
        status := "enabled"
        if sm.api.IsPackageSourceDisabled(sm.config, s.Source.Key) {
            status = "disabled"
        }
        fmt.Printf("  %d. %s (%s): %s\n", s.Priority, s.Source.Key, status, s.Source.Value)
    }
}

func (sm *SourceManager) ConfigureForEnvironment(env string) {
    fmt.Printf("\nConfiguring for environment: %s\n", env)
    
    switch env {
    case "development":
        sm.api.EnablePackageSource(sm.config, "local-dev")
        sm.api.EnablePackageSource(sm.config, "company-preview")
        sm.api.SetActivePackageSource(sm.config, "local-dev", "./packages")
        fmt.Println("  - Enabled local development sources")
        
    case "staging":
        sm.api.DisablePackageSource(sm.config, "local-dev")
        sm.api.EnablePackageSource(sm.config, "company-stable")
        sm.api.SetActivePackageSource(sm.config, "company-stable", "https://stable.company.com/nuget")
        fmt.Println("  - Configured for staging environment")
        
    case "production":
        sm.api.DisablePackageSource(sm.config, "local-dev")
        sm.api.DisablePackageSource(sm.config, "company-preview")
        sm.api.EnablePackageSource(sm.config, "nuget.org")
        sm.api.SetActivePackageSource(sm.config, "nuget.org", "https://api.nuget.org/v3/index.json")
        fmt.Println("  - Configured for production environment")
    }
}

func (sm *SourceManager) BulkDisable(sourceKeys []string) {
    fmt.Printf("Bulk disabling %d sources...\n", len(sourceKeys))
    for _, key := range sourceKeys {
        if sm.api.GetPackageSource(sm.config, key) != nil {
            sm.api.DisablePackageSource(sm.config, key)
            fmt.Printf("  - Disabled: %s\n", key)
        }
    }
}

func (sm *SourceManager) BulkEnable(sourceKeys []string) {
    fmt.Printf("Bulk enabling %d sources...\n", len(sourceKeys))
    for _, key := range sourceKeys {
        if sm.api.GetPackageSource(sm.config, key) != nil {
            sm.api.EnablePackageSource(sm.config, key)
            fmt.Printf("  - Enabled: %s\n", key)
        }
    }
}

func (sm *SourceManager) DisplaySummary() {
    fmt.Println("\n=== Configuration Summary ===")
    
    totalSources := len(sm.config.PackageSources.Add)
    enabledCount := 0
    
    for _, source := range sm.config.PackageSources.Add {
        if !sm.api.IsPackageSourceDisabled(sm.config, source.Key) {
            enabledCount++
        }
    }
    
    fmt.Printf("Total sources: %d\n", totalSources)
    fmt.Printf("Enabled: %d\n", enabledCount)
    fmt.Printf("Disabled: %d\n", totalSources-enabledCount)
    
    if activeSource := sm.api.GetActivePackageSource(sm.config); activeSource != nil {
        fmt.Printf("Active source: %s\n", activeSource.Key)
    }
}
```

## Key Concepts

### Source Types

1. **HTTP/HTTPS Sources**: Remote NuGet feeds
2. **Local File System**: Local directory paths
3. **UNC Paths**: Network share locations
4. **Protocol Versions**: v2 (legacy) vs v3 (modern)

### Source States

1. **Enabled**: Source is active and will be searched
2. **Disabled**: Source exists but won't be used
3. **Active**: The primary source for operations

### Best Practices

1. **Use descriptive names**: Choose meaningful source keys
2. **Set protocol versions**: Specify v3 for modern sources
3. **Validate sources**: Check accessibility before adding
4. **Manage states**: Disable unused sources for performance
5. **Set active source**: Define a primary source for operations
6. **Environment-specific**: Configure sources per environment

## Next Steps

After mastering package source management:

1. Learn about [Credentials](./credentials.md) for authenticated sources
2. Explore [Config Options](./config-options.md) for global settings
3. Study [Position-Aware Editing](./position-aware-editing.md) for precise modifications

This guide provides comprehensive examples for managing NuGet package sources, covering everything from basic operations to advanced management scenarios.
