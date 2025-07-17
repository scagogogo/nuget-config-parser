# Finding Configurations

This example demonstrates how to locate NuGet configuration files across different platforms and directory structures.

## Overview

Configuration discovery involves:
- Searching standard locations for NuGet.Config files
- Understanding platform-specific paths
- Handling project-level vs global configurations
- Implementing fallback strategies

## Example 1: Basic Configuration Discovery

The simplest way to find a configuration file:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    api := nuget.NewAPI()
    
    // Find the first available configuration file
    configPath, err := api.FindConfigFile()
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Println("No NuGet configuration file found in standard locations")
            fmt.Println("Consider creating a default configuration")
            return
        }
        log.Fatalf("Error searching for config: %v", err)
    }
    
    fmt.Printf("Found configuration file: %s\n", configPath)
    
    // Parse the found configuration
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("Failed to parse found config: %v", err)
    }
    
    fmt.Printf("Configuration contains %d package sources\n", len(config.PackageSources.Add))
}
```

## Example 2: Find All Configuration Files

Discover all available configuration files in the search hierarchy:

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Find all configuration files
    configPaths := api.FindAllConfigFiles()
    
    fmt.Printf("Found %d configuration files:\n", len(configPaths))
    
    if len(configPaths) == 0 {
        fmt.Println("No configuration files found in standard locations")
        displaySearchPaths()
        return
    }
    
    // Display all found configurations with details
    for i, configPath := range configPaths {
        fmt.Printf("\n%d. %s\n", i+1, configPath)
        
        // Check if file is readable
        if info, err := os.Stat(configPath); err == nil {
            fmt.Printf("   Size: %d bytes\n", info.Size())
            fmt.Printf("   Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
        }
        
        // Try to parse and show basic info
        config, err := api.ParseFromFile(configPath)
        if err != nil {
            fmt.Printf("   ❌ Parse error: %v\n", err)
        } else {
            fmt.Printf("   ✅ Valid configuration\n")
            fmt.Printf("   📦 Package sources: %d\n", len(config.PackageSources.Add))
            
            // Show first few sources
            for j, source := range config.PackageSources.Add {
                if j >= 3 {
                    fmt.Printf("   ... and %d more\n", len(config.PackageSources.Add)-3)
                    break
                }
                fmt.Printf("   - %s\n", source.Key)
            }
        }
    }
}

func displaySearchPaths() {
    fmt.Println("\nStandard search locations:")
    fmt.Println("1. Current directory: ./NuGet.Config")
    fmt.Println("2. Parent directories (walking up)")
    
    if home := os.Getenv("HOME"); home != "" {
        fmt.Printf("3. User config: %s/.config/NuGet/NuGet.Config\n", home)
    }
    
    fmt.Println("4. System config: /etc/NuGet/NuGet.Config")
    fmt.Println("\nNote: Actual paths vary by operating system")
}
```

## Example 3: Project-Specific Configuration Discovery

Find configuration files starting from a specific project directory:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    api := nuget.NewAPI()
    
    // Get current working directory
    currentDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Failed to get current directory: %v", err)
    }
    
    fmt.Printf("Searching for project configuration starting from: %s\n", currentDir)
    
    // Find project-specific configuration
    projectConfig, err := api.FindProjectConfig(currentDir)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Println("No project-specific configuration found")
            
            // Fall back to global configuration
            fmt.Println("Searching for global configuration...")
            globalConfig, err := api.FindConfigFile()
            if err != nil {
                fmt.Println("No global configuration found either")
                return
            }
            
            fmt.Printf("Using global configuration: %s\n", globalConfig)
            projectConfig = globalConfig
        } else {
            log.Fatalf("Error searching for project config: %v", err)
        }
    } else {
        fmt.Printf("Found project configuration: %s\n", projectConfig)
    }
    
    // Show the hierarchy of directories searched
    showSearchHierarchy(currentDir)
    
    // Parse and display the configuration
    config, err := api.ParseFromFile(projectConfig)
    if err != nil {
        log.Fatalf("Failed to parse configuration: %v", err)
    }
    
    displayConfigSummary(config, projectConfig)
}

func showSearchHierarchy(startDir string) {
    fmt.Println("\nSearch hierarchy (from most specific to most general):")
    
    dir := startDir
    level := 1
    
    for {
        configPath := filepath.Join(dir, "NuGet.Config")
        exists := fileExists(configPath)
        
        status := "❌"
        if exists {
            status = "✅"
        }
        
        fmt.Printf("%d. %s %s\n", level, configPath, status)
        
        parent := filepath.Dir(dir)
        if parent == dir {
            // Reached root directory
            break
        }
        
        dir = parent
        level++
        
        // Limit search depth for display
        if level > 10 {
            fmt.Println("   ... (search continues to root)")
            break
        }
    }
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func displayConfigSummary(config *types.NuGetConfig, configPath string) {
    fmt.Printf("\n=== Configuration Summary ===\n")
    fmt.Printf("File: %s\n", configPath)
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
    
    if len(config.PackageSources.Add) > 0 {
        fmt.Println("\nPackage sources:")
        for _, source := range config.PackageSources.Add {
            fmt.Printf("  - %s: %s\n", source.Key, source.Value)
        }
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("\nActive source: %s\n", config.ActivePackageSource.Add.Key)
    }
}
```

## Example 4: Cross-Platform Configuration Discovery

Handle platform-specific configuration locations:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Printf("Operating System: %s\n", runtime.GOOS)
    fmt.Printf("Architecture: %s\n", runtime.GOARCH)
    
    // Show platform-specific paths
    showPlatformPaths()
    
    // Find all configurations
    configPaths := api.FindAllConfigFiles()
    
    fmt.Printf("\nFound %d configuration files:\n", len(configPaths))
    
    for i, configPath := range configPaths {
        fmt.Printf("%d. %s\n", i+1, configPath)
        
        // Categorize the configuration
        category := categorizeConfigPath(configPath)
        fmt.Printf("   Category: %s\n", category)
        
        // Check accessibility
        if isReadable(configPath) {
            fmt.Printf("   Status: ✅ Readable\n")
        } else {
            fmt.Printf("   Status: ❌ Not readable\n")
        }
    }
    
    // Demonstrate finding and parsing
    if len(configPaths) > 0 {
        fmt.Printf("\nUsing first available configuration: %s\n", configPaths[0])
        
        config, err := api.ParseFromFile(configPaths[0])
        if err != nil {
            fmt.Printf("Failed to parse: %v\n", err)
        } else {
            fmt.Printf("Successfully parsed %d package sources\n", len(config.PackageSources.Add))
        }
    }
}

func showPlatformPaths() {
    fmt.Println("\nPlatform-specific configuration locations:")
    
    switch runtime.GOOS {
    case "windows":
        fmt.Println("User config: %APPDATA%\\NuGet\\NuGet.Config")
        fmt.Println("System config: %ProgramData%\\NuGet\\NuGet.Config")
        
        if appdata := os.Getenv("APPDATA"); appdata != "" {
            fmt.Printf("Resolved user: %s\n", filepath.Join(appdata, "NuGet", "NuGet.Config"))
        }
        
        if programdata := os.Getenv("ProgramData"); programdata != "" {
            fmt.Printf("Resolved system: %s\n", filepath.Join(programdata, "NuGet", "NuGet.Config"))
        }
        
    case "darwin":
        fmt.Println("User config: ~/Library/Application Support/NuGet/NuGet.Config")
        fmt.Println("System config: /Library/Application Support/NuGet/NuGet.Config")
        
        if home := os.Getenv("HOME"); home != "" {
            fmt.Printf("Resolved user: %s\n", filepath.Join(home, "Library", "Application Support", "NuGet", "NuGet.Config"))
        }
        
    default: // Linux and other Unix systems
        fmt.Println("User config: ~/.config/NuGet/NuGet.Config")
        fmt.Println("System config: /etc/NuGet/NuGet.Config")
        
        if home := os.Getenv("HOME"); home != "" {
            fmt.Printf("Resolved user: %s\n", filepath.Join(home, ".config", "NuGet", "NuGet.Config"))
        }
        
        // Check XDG_CONFIG_HOME
        if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
            fmt.Printf("XDG config: %s\n", filepath.Join(xdgConfig, "NuGet", "NuGet.Config"))
        }
    }
}

func categorizeConfigPath(configPath string) string {
    absPath, _ := filepath.Abs(configPath)
    
    // Check if it's in current directory or subdirectory
    if cwd, err := os.Getwd(); err == nil {
        if rel, err := filepath.Rel(cwd, absPath); err == nil && !filepath.IsAbs(rel) {
            return "Project/Local"
        }
    }
    
    // Check if it's in user directory
    if home := os.Getenv("HOME"); home != "" {
        if rel, err := filepath.Rel(home, absPath); err == nil && !filepath.IsAbs(rel) {
            return "User"
        }
    }
    
    // Check common system paths
    systemPaths := []string{"/etc", "/usr/local/etc", "/opt"}
    for _, sysPath := range systemPaths {
        if rel, err := filepath.Rel(sysPath, absPath); err == nil && !filepath.IsAbs(rel) {
            return "System"
        }
    }
    
    return "Other"
}

func isReadable(path string) bool {
    file, err := os.Open(path)
    if err != nil {
        return false
    }
    file.Close()
    return true
}
```

## Example 5: Advanced Discovery with Custom Search Paths

Implement custom search logic with additional paths:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/finder"
)

func main() {
    // Create a custom finder with additional search paths
    customPaths := []string{
        "./NuGet.Config",
        "./config/NuGet.Config",
        "./nuget/NuGet.Config",
        "/opt/nuget/NuGet.Config",
        "/usr/local/etc/nuget/NuGet.Config",
    }
    
    // Add environment-based paths
    if nugetHome := os.Getenv("NUGET_HOME"); nugetHome != "" {
        customPaths = append(customPaths, filepath.Join(nugetHome, "NuGet.Config"))
    }
    
    if projectRoot := os.Getenv("PROJECT_ROOT"); projectRoot != "" {
        customPaths = append(customPaths, filepath.Join(projectRoot, "NuGet.Config"))
    }
    
    // Create finder with custom paths
    configFinder := finder.NewConfigFinderWithPaths(customPaths)
    
    fmt.Println("Custom configuration search paths:")
    for i, path := range customPaths {
        exists := "❌"
        if fileExists(path) {
            exists = "✅"
        }
        fmt.Printf("%d. %s %s\n", i+1, path, exists)
    }
    
    // Find first available configuration
    configPath, err := configFinder.FindConfigFile()
    if err != nil {
        fmt.Println("\nNo configuration found in custom paths")
        
        // Fall back to standard discovery
        fmt.Println("Falling back to standard discovery...")
        api := nuget.NewAPI()
        
        if standardPath, err := api.FindConfigFile(); err == nil {
            fmt.Printf("Found standard configuration: %s\n", standardPath)
            configPath = standardPath
        } else {
            fmt.Println("No configuration found anywhere")
            return
        }
    } else {
        fmt.Printf("\nFound configuration: %s\n", configPath)
    }
    
    // Parse and use the configuration
    api := nuget.NewAPI()
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        fmt.Printf("Failed to parse configuration: %v\n", err)
        return
    }
    
    fmt.Printf("Successfully loaded configuration with %d package sources\n", len(config.PackageSources.Add))
    
    // Show configuration hierarchy
    showConfigurationHierarchy(configPath, customPaths)
}

func showConfigurationHierarchy(selectedPath string, searchPaths []string) {
    fmt.Println("\nConfiguration hierarchy (highest to lowest priority):")
    
    for i, path := range searchPaths {
        status := "Not found"
        marker := "  "
        
        if fileExists(path) {
            status = "Available"
            if path == selectedPath {
                status = "SELECTED"
                marker = "→ "
            }
        }
        
        fmt.Printf("%s%d. %s (%s)\n", marker, i+1, path, status)
    }
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}
```

## Key Concepts

### Search Order

The library searches for configuration files in this order:

1. **Current directory**: `./NuGet.Config`
2. **Parent directories**: Walking up the directory tree
3. **User configuration**: Platform-specific user directory
4. **System configuration**: Platform-specific system directory

### Platform Differences

- **Windows**: Uses `%APPDATA%` and `%ProgramData%`
- **macOS**: Uses `~/Library/Application Support/` and `/Library/Application Support/`
- **Linux/Unix**: Uses `~/.config/` (or `$XDG_CONFIG_HOME`) and `/etc/`

### Best Practices

1. **Handle missing files gracefully**: Always check for `IsNotFoundError`
2. **Provide fallbacks**: Have a strategy when no configuration is found
3. **Respect hierarchy**: Project configs override global configs
4. **Check file permissions**: Ensure files are readable before parsing
5. **Use appropriate discovery method**: Choose between single file or all files based on needs

## Common Patterns

### Pattern 1: Find or Create

```go
config, configPath, err := api.FindAndParseConfig()
if err != nil {
    // Create default if not found
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
    api.SaveConfig(config, configPath)
}
```

### Pattern 2: Hierarchical Search

```go
// Try project-specific first
if projectConfig, err := api.FindProjectConfig("."); err == nil {
    return api.ParseFromFile(projectConfig)
}

// Fall back to global
if globalConfig, err := api.FindConfigFile(); err == nil {
    return api.ParseFromFile(globalConfig)
}

// Create default as last resort
return api.CreateDefaultConfig(), nil
```

### Pattern 3: Multiple Configuration Merging

```go
configs := api.FindAllConfigFiles()
var mergedSources []types.PackageSource

for _, configPath := range configs {
    if config, err := api.ParseFromFile(configPath); err == nil {
        mergedSources = append(mergedSources, config.PackageSources.Add...)
    }
}
```

## Next Steps

After mastering configuration discovery:

1. Learn about [Creating Configs](./creating-configs.md) to generate new configurations
2. Explore [Basic Parsing](./basic-parsing.md) to understand parsing details
3. Study [Modifying Configs](./modifying-configs.md) to update found configurations

This guide provides comprehensive examples for finding NuGet configuration files across different scenarios and platforms.
