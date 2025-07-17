# Modifying Configurations

This example demonstrates how to modify existing NuGet configuration files using the NuGet Config Parser library.

## Overview

Configuration modification involves:
- Loading existing configurations
- Adding, updating, or removing package sources
- Managing credentials and authentication
- Configuring global settings
- Saving changes back to files

## Example 1: Basic Configuration Modification

The simplest way to modify an existing configuration:

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
    
    // Load existing configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        if errors.IsNotFoundError(err) {
            // Create default if not found
            config = api.CreateDefaultConfig()
            configPath = "NuGet.Config"
            fmt.Println("Created new configuration")
        } else {
            log.Fatalf("Failed to load config: %v", err)
        }
    } else {
        fmt.Printf("Loaded existing configuration from: %s\n", configPath)
    }
    
    // Show current sources
    fmt.Printf("Current package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
    
    // Add a new package source
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    fmt.Println("Added company package source")
    
    // Update an existing source (if it exists)
    existingSource := api.GetPackageSource(config, "nuget.org")
    if existingSource != nil {
        api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
        fmt.Println("Updated nuget.org source")
    }
    
    // Save the modified configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("Configuration saved to: %s\n", configPath)
    
    // Show updated sources
    fmt.Printf("Updated package sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
}
```

## Example 2: Advanced Source Management

Comprehensive package source management:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Println("=== Package Source Management ===")
    
    // Add multiple sources
    sources := []struct {
        key     string
        url     string
        version string
        enabled bool
    }{
        {"company-internal", "https://internal.company.com/nuget", "3", true},
        {"company-preview", "https://preview.company.com/nuget", "3", false},
        {"local-dev", "./local-packages", "", false},
        {"backup-feed", "https://backup.nuget.com/api/v2", "2", true},
    }
    
    for _, source := range sources {
        api.AddPackageSource(config, source.key, source.url, source.version)
        fmt.Printf("Added source: %s\n", source.key)
        
        if !source.enabled {
            api.DisablePackageSource(config, source.key)
            fmt.Printf("  - Disabled: %s\n", source.key)
        }
    }
    
    // Remove old or unwanted sources
    sourcesToRemove := []string{"old-feed", "deprecated-source"}
    for _, sourceKey := range sourcesToRemove {
        if api.GetPackageSource(config, sourceKey) != nil {
            removed := api.RemovePackageSource(config, sourceKey)
            if removed {
                fmt.Printf("Removed source: %s\n", sourceKey)
            }
        }
    }
    
    // Set active package source
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    fmt.Println("Set nuget.org as active source")
    
    // Display final configuration
    fmt.Println("\n=== Final Configuration ===")
    displayPackageSources(api, config)
    
    // Save configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nConfiguration saved to: %s\n", configPath)
}

func displayPackageSources(api *nuget.API, config *types.NuGetConfig) {
    fmt.Printf("Package Sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        
        fmt.Printf("  - %s (%s): %s", source.Key, status, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" [v%s]", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    if activeSource := api.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("\nActive Source: %s\n", activeSource.Key)
    }
}
```

## Example 3: Credential Management

Managing authentication credentials for private sources:

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
    
    // Load configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Println("=== Credential Management ===")
    
    // Add private sources with credentials
    privateSources := []struct {
        key      string
        url      string
        username string
        password string
    }{
        {"company-private", "https://private.company.com/nuget", "employee", getEnvOrDefault("COMPANY_NUGET_PASSWORD", "defaultpass")},
        {"azure-artifacts", "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json", "myuser", getEnvOrDefault("AZURE_PAT", "pat_token")},
        {"github-packages", "https://nuget.pkg.github.com/myorg/index.json", "github_user", getEnvOrDefault("GITHUB_TOKEN", "ghp_token")},
    }
    
    for _, source := range privateSources {
        // Add the package source
        api.AddPackageSource(config, source.key, source.url, "3")
        fmt.Printf("Added private source: %s\n", source.key)
        
        // Add credentials
        api.AddCredential(config, source.key, source.username, source.password)
        fmt.Printf("  - Added credentials for: %s\n", source.key)
    }
    
    // Verify credentials were added
    fmt.Println("\n=== Credential Verification ===")
    for _, source := range privateSources {
        credential := api.GetCredential(config, source.key)
        if credential != nil {
            fmt.Printf("✅ %s has credentials configured\n", source.key)
            
            // Display credential info (be careful with passwords!)
            for _, cred := range credential.Add {
                if cred.Key == "Username" {
                    fmt.Printf("   Username: %s\n", cred.Value)
                }
            }
        } else {
            fmt.Printf("❌ %s missing credentials\n", source.key)
        }
    }
    
    // Remove credentials for a source
    fmt.Println("\n=== Credential Removal ===")
    sourceToRemoveCreds := "old-private-source"
    if api.GetCredential(config, sourceToRemoveCreds) != nil {
        removed := api.RemoveCredential(config, sourceToRemoveCreds)
        if removed {
            fmt.Printf("Removed credentials for: %s\n", sourceToRemoveCreds)
        }
    }
    
    // Save configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nConfiguration with credentials saved to: %s\n", configPath)
}

func getEnvOrDefault(envVar, defaultValue string) string {
    if value := os.Getenv(envVar); value != "" {
        return value
    }
    return defaultValue
}
```

## Example 4: Global Configuration Options

Managing global NuGet settings:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Println("=== Global Configuration Management ===")
    
    // Set up custom package locations
    homeDir := os.Getenv("HOME")
    if homeDir == "" {
        homeDir = os.Getenv("USERPROFILE") // Windows
    }
    
    globalPackagesPath := filepath.Join(homeDir, ".nuget", "packages")
    repositoryPath := "./packages"
    
    // Configure global settings
    globalSettings := map[string]string{
        "globalPackagesFolder":     globalPackagesPath,
        "repositoryPath":          repositoryPath,
        "defaultPushSource":       "https://api.nuget.org/v3/index.json",
        "dependencyVersion":       "Highest",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
    }
    
    fmt.Println("Setting global configuration options:")
    for key, value := range globalSettings {
        api.AddConfigOption(config, key, value)
        fmt.Printf("  - %s: %s\n", key, value)
    }
    
    // Configure proxy settings (if needed)
    proxyUrl := os.Getenv("HTTP_PROXY")
    if proxyUrl != "" {
        api.AddConfigOption(config, "http_proxy", proxyUrl)
        fmt.Printf("  - http_proxy: %s\n", proxyUrl)
        
        // Add proxy credentials if available
        proxyUser := os.Getenv("PROXY_USER")
        proxyPass := os.Getenv("PROXY_PASS")
        if proxyUser != "" && proxyPass != "" {
            api.AddConfigOption(config, "http_proxy.user", proxyUser)
            api.AddConfigOption(config, "http_proxy.password", proxyPass)
            fmt.Println("  - Added proxy credentials")
        }
    }
    
    // Display current configuration options
    fmt.Println("\n=== Current Configuration Options ===")
    if config.Config != nil {
        for _, option := range config.Config.Add {
            // Mask sensitive information
            value := option.Value
            if option.Key == "http_proxy.password" || option.Key == "password" {
                value = "***masked***"
            }
            fmt.Printf("  - %s: %s\n", option.Key, value)
        }
    }
    
    // Remove obsolete settings
    obsoleteSettings := []string{"oldSetting", "deprecatedOption"}
    for _, setting := range obsoleteSettings {
        if api.GetConfigOption(config, setting) != "" {
            removed := api.RemoveConfigOption(config, setting)
            if removed {
                fmt.Printf("Removed obsolete setting: %s\n", setting)
            }
        }
    }
    
    // Save configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nGlobal configuration saved to: %s\n", configPath)
}
```

## Example 5: Batch Configuration Updates

Performing multiple configuration updates efficiently:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Println("=== Batch Configuration Updates ===")
    
    // Define all changes to make
    updates := ConfigUpdates{
        AddSources: []SourceConfig{
            {"prod-feed", "https://prod.company.com/nuget", "3", true},
            {"staging-feed", "https://staging.company.com/nuget", "3", false},
            {"dev-feed", "https://dev.company.com/nuget", "3", false},
        },
        RemoveSources: []string{"old-feed", "deprecated-source"},
        UpdateSources: []SourceUpdate{
            {"nuget.org", "https://api.nuget.org/v3/index.json", "3"},
        },
        AddCredentials: []CredentialConfig{
            {"prod-feed", "prod_user", "prod_pass"},
            {"staging-feed", "staging_user", "staging_pass"},
        },
        ConfigOptions: map[string]string{
            "globalPackagesFolder": "/custom/packages",
            "dependencyVersion":    "Highest",
            "automaticPackageRestore": "true",
        },
        ActiveSource: &SourceConfig{"nuget.org", "https://api.nuget.org/v3/index.json", "3", true},
    }
    
    // Apply all updates
    err = applyConfigUpdates(api, config, updates)
    if err != nil {
        log.Fatalf("Failed to apply updates: %v", err)
    }
    
    // Save configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nBatch updates completed and saved to: %s\n", configPath)
    
    // Display final configuration summary
    displayConfigSummary(api, config)
}

type ConfigUpdates struct {
    AddSources     []SourceConfig
    RemoveSources  []string
    UpdateSources  []SourceUpdate
    AddCredentials []CredentialConfig
    ConfigOptions  map[string]string
    ActiveSource   *SourceConfig
}

type SourceConfig struct {
    Key     string
    URL     string
    Version string
    Enabled bool
}

type SourceUpdate struct {
    Key     string
    NewURL  string
    Version string
}

type CredentialConfig struct {
    SourceKey string
    Username  string
    Password  string
}

func applyConfigUpdates(api *nuget.API, config *types.NuGetConfig, updates ConfigUpdates) error {
    // Add new sources
    fmt.Printf("Adding %d new sources...\n", len(updates.AddSources))
    for _, source := range updates.AddSources {
        api.AddPackageSource(config, source.Key, source.URL, source.Version)
        if !source.Enabled {
            api.DisablePackageSource(config, source.Key)
        }
        fmt.Printf("  + %s\n", source.Key)
    }
    
    // Remove sources
    fmt.Printf("Removing %d sources...\n", len(updates.RemoveSources))
    for _, sourceKey := range updates.RemoveSources {
        if api.RemovePackageSource(config, sourceKey) {
            fmt.Printf("  - %s\n", sourceKey)
        }
    }
    
    // Update existing sources
    fmt.Printf("Updating %d sources...\n", len(updates.UpdateSources))
    for _, update := range updates.UpdateSources {
        api.AddPackageSource(config, update.Key, update.NewURL, update.Version)
        fmt.Printf("  ~ %s\n", update.Key)
    }
    
    // Add credentials
    fmt.Printf("Adding credentials for %d sources...\n", len(updates.AddCredentials))
    for _, cred := range updates.AddCredentials {
        api.AddCredential(config, cred.SourceKey, cred.Username, cred.Password)
        fmt.Printf("  🔐 %s\n", cred.SourceKey)
    }
    
    // Set configuration options
    fmt.Printf("Setting %d configuration options...\n", len(updates.ConfigOptions))
    for key, value := range updates.ConfigOptions {
        api.AddConfigOption(config, key, value)
        fmt.Printf("  ⚙️  %s: %s\n", key, value)
    }
    
    // Set active source
    if updates.ActiveSource != nil {
        api.SetActivePackageSource(config, updates.ActiveSource.Key, updates.ActiveSource.URL)
        fmt.Printf("Set active source: %s\n", updates.ActiveSource.Key)
    }
    
    return nil
}

func displayConfigSummary(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== Configuration Summary ===")
    
    fmt.Printf("Package Sources: %d\n", len(config.PackageSources.Add))
    enabledCount := 0
    for _, source := range config.PackageSources.Add {
        if !api.IsPackageSourceDisabled(config, source.Key) {
            enabledCount++
        }
    }
    fmt.Printf("  - Enabled: %d\n", enabledCount)
    fmt.Printf("  - Disabled: %d\n", len(config.PackageSources.Add)-enabledCount)
    
    if config.PackageSourceCredentials != nil {
        fmt.Printf("Authenticated Sources: %d\n", len(config.PackageSourceCredentials.Sources))
    }
    
    if config.Config != nil {
        fmt.Printf("Configuration Options: %d\n", len(config.Config.Add))
    }
    
    if activeSource := api.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("Active Source: %s\n", activeSource.Key)
    }
}
```

## Key Concepts

### Modification Strategies

1. **Incremental Updates**: Make small, targeted changes
2. **Batch Operations**: Group related changes together
3. **Validation**: Verify changes before saving
4. **Backup**: Consider backing up before major changes

### Best Practices

1. **Load before modify**: Always load existing configuration first
2. **Check existence**: Verify sources exist before updating/removing
3. **Handle errors**: Properly handle modification errors
4. **Save atomically**: Save all changes at once
5. **Validate results**: Ensure configuration is valid after changes

### Common Patterns

#### Pattern 1: Safe Modification
```go
// Load existing or create default
config, configPath, err := api.FindAndParseConfig()
if err != nil {
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
}

// Make changes
api.AddPackageSource(config, "new-source", "https://example.com", "3")

// Save changes
err = api.SaveConfig(config, configPath)
```

#### Pattern 2: Conditional Updates
```go
// Only update if source exists
if api.GetPackageSource(config, "existing-source") != nil {
    api.AddPackageSource(config, "existing-source", "https://new-url.com", "3")
}

// Only remove if source exists
if api.GetPackageSource(config, "old-source") != nil {
    api.RemovePackageSource(config, "old-source")
}
```

#### Pattern 3: Environment-Based Configuration
```go
environment := os.Getenv("ENVIRONMENT")
switch environment {
case "development":
    api.AddPackageSource(config, "dev-feed", "https://dev.company.com", "3")
case "staging":
    api.AddPackageSource(config, "staging-feed", "https://staging.company.com", "3")
case "production":
    api.AddPackageSource(config, "prod-feed", "https://prod.company.com", "3")
}
```

## Next Steps

After mastering configuration modification:

1. Learn about [Position-Aware Editing](./position-aware-editing.md) for minimal-diff changes
2. Explore [Package Sources](./package-sources.md) for advanced source management
3. Study [Credentials](./credentials.md) for secure authentication handling

This guide provides comprehensive examples for modifying NuGet configuration files, covering various scenarios from simple updates to complex batch operations.
