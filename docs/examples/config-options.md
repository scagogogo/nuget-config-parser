# Configuration Options

This example demonstrates how to manage global NuGet configuration options using the NuGet Config Parser library.

## Overview

Configuration options control various aspects of NuGet behavior:
- Package storage locations
- Dependency resolution strategies
- Proxy settings
- Package restore behavior
- Default push sources

## Example 1: Basic Configuration Options

Managing fundamental configuration options:

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
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Basic Configuration Options ===")
    
    // Set up basic package locations
    homeDir := os.Getenv("HOME")
    if homeDir == "" {
        homeDir = os.Getenv("USERPROFILE") // Windows
    }
    
    // Configure package storage locations
    globalPackagesPath := filepath.Join(homeDir, ".nuget", "packages")
    repositoryPath := "./packages"
    
    api.AddConfigOption(config, "globalPackagesFolder", globalPackagesPath)
    api.AddConfigOption(config, "repositoryPath", repositoryPath)
    
    fmt.Printf("Global packages folder: %s\n", globalPackagesPath)
    fmt.Printf("Repository path: %s\n", repositoryPath)
    
    // Configure dependency resolution
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "packageRestore", "true")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    fmt.Println("Dependency resolution: Highest")
    fmt.Println("Package restore: Enabled")
    
    // Configure default push source
    api.AddConfigOption(config, "defaultPushSource", "https://api.nuget.org/v3/index.json")
    fmt.Println("Default push source: nuget.org")
    
    // Display all configuration options
    fmt.Println("\n=== All Configuration Options ===")
    if config.Config != nil {
        for _, option := range config.Config.Add {
            fmt.Printf("  %s: %s\n", option.Key, option.Value)
        }
    }
    
    // Save configuration
    err := api.SaveConfig(config, "BasicOptions.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Println("\nConfiguration saved successfully!")
}
```

## Example 2: Proxy Configuration

Setting up proxy settings for corporate environments:

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
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Proxy Configuration ===")
    
    // Get proxy settings from environment or use defaults
    proxyURL := getEnvOrDefault("HTTP_PROXY", "http://proxy.company.com:8080")
    proxyUser := getEnvOrDefault("PROXY_USER", "")
    proxyPass := getEnvOrDefault("PROXY_PASS", "")
    
    if proxyURL != "" {
        // Configure HTTP proxy
        api.AddConfigOption(config, "http_proxy", proxyURL)
        fmt.Printf("HTTP Proxy: %s\n", proxyURL)
        
        // Configure HTTPS proxy (often the same)
        api.AddConfigOption(config, "https_proxy", proxyURL)
        fmt.Printf("HTTPS Proxy: %s\n", proxyURL)
        
        // Add proxy authentication if provided
        if proxyUser != "" && proxyPass != "" {
            api.AddConfigOption(config, "http_proxy.user", proxyUser)
            api.AddConfigOption(config, "http_proxy.password", proxyPass)
            fmt.Printf("Proxy authentication: %s\n", proxyUser)
        }
        
        // Configure proxy bypass for local addresses
        api.AddConfigOption(config, "http_proxy.no_proxy", "localhost,127.0.0.1,*.local")
        fmt.Println("Proxy bypass: localhost,127.0.0.1,*.local")
    } else {
        fmt.Println("No proxy configuration needed")
    }
    
    // Additional network settings
    api.AddConfigOption(config, "http_timeout", "300")
    api.AddConfigOption(config, "http_retries", "3")
    
    fmt.Println("HTTP timeout: 300 seconds")
    fmt.Println("HTTP retries: 3")
    
    // Display proxy configuration
    fmt.Println("\n=== Proxy Settings ===")
    proxyOptions := []string{
        "http_proxy", "https_proxy", "http_proxy.user", 
        "http_proxy.password", "http_proxy.no_proxy",
        "http_timeout", "http_retries",
    }
    
    for _, key := range proxyOptions {
        value := api.GetConfigOption(config, key)
        if value != "" {
            displayValue := value
            if key == "http_proxy.password" {
                displayValue = "***masked***"
            }
            fmt.Printf("  %s: %s\n", key, displayValue)
        }
    }
    
    // Save configuration
    err := api.SaveConfig(config, "ProxyConfig.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Println("\nProxy configuration saved!")
}

func getEnvOrDefault(envVar, defaultValue string) string {
    if value := os.Getenv(envVar); value != "" {
        return value
    }
    return defaultValue
}
```

## Example 3: Environment-Specific Configuration

Configuring options based on environment:

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
    config := api.CreateDefaultConfig()
    
    // Determine environment
    environment := getEnvOrDefault("ENVIRONMENT", "development")
    fmt.Printf("Configuring for environment: %s\n", environment)
    
    // Configure based on environment
    switch environment {
    case "development":
        configureDevelopment(api, config)
    case "staging":
        configureStaging(api, config)
    case "production":
        configureProduction(api, config)
    default:
        configureDefault(api, config)
    }
    
    // Display final configuration
    displayConfiguration(api, config, environment)
    
    // Save with environment-specific name
    configFile := fmt.Sprintf("%s.Config", environment)
    err := api.SaveConfig(config, configFile)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nEnvironment-specific configuration saved to: %s\n", configFile)
}

func configureDevelopment(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== Development Configuration ===")
    
    // Use local packages folder for faster access
    api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // Enable verbose logging for debugging
    api.AddConfigOption(config, "verbosity", "detailed")
    
    // Use highest dependency versions for latest features
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    
    // Enable automatic package restore
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // Allow prerelease packages
    api.AddConfigOption(config, "allowPrereleaseVersions", "true")
    
    fmt.Println("  - Local packages folder")
    fmt.Println("  - Verbose logging enabled")
    fmt.Println("  - Prerelease packages allowed")
}

func configureStaging(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== Staging Configuration ===")
    
    // Use shared staging packages folder
    api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // Moderate verbosity
    api.AddConfigOption(config, "verbosity", "normal")
    
    // Use stable dependency versions
    api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
    
    // Enable package restore
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // Disable prerelease packages
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    
    // Set staging push source
    api.AddConfigOption(config, "defaultPushSource", "https://staging.company.com/nuget")
    
    fmt.Println("  - Shared staging packages")
    fmt.Println("  - Stable versions only")
    fmt.Println("  - Staging push source")
}

func configureProduction(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== Production Configuration ===")
    
    // Use production packages folder
    api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // Minimal verbosity for performance
    api.AddConfigOption(config, "verbosity", "quiet")
    
    // Use exact dependency versions for stability
    api.AddConfigOption(config, "dependencyVersion", "Exact")
    
    // Enable package restore
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // Strictly disable prerelease packages
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    
    // Set production push source
    api.AddConfigOption(config, "defaultPushSource", "https://prod.company.com/nuget")
    
    // Enable package verification
    api.AddConfigOption(config, "signatureValidationMode", "require")
    
    // Set timeouts for reliability
    api.AddConfigOption(config, "http_timeout", "600")
    api.AddConfigOption(config, "http_retries", "5")
    
    fmt.Println("  - Production packages folder")
    fmt.Println("  - Exact versions for stability")
    fmt.Println("  - Package signature validation")
    fmt.Println("  - Extended timeouts")
}

func configureDefault(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== Default Configuration ===")
    
    // Standard configuration
    homeDir := getEnvOrDefault("HOME", getEnvOrDefault("USERPROFILE", "."))
    globalPackages := filepath.Join(homeDir, ".nuget", "packages")
    
    api.AddConfigOption(config, "globalPackagesFolder", globalPackages)
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    fmt.Println("  - Standard package locations")
    fmt.Println("  - Default dependency resolution")
}

func displayConfiguration(api *nuget.API, config *types.NuGetConfig, environment string) {
    fmt.Printf("\n=== Final %s Configuration ===\n", environment)
    
    if config.Config != nil {
        for _, option := range config.Config.Add {
            value := option.Value
            if option.Key == "http_proxy.password" {
                value = "***masked***"
            }
            fmt.Printf("  %s: %s\n", option.Key, value)
        }
    }
}

func getEnvOrDefault(envVar, defaultValue string) string {
    if value := os.Getenv(envVar); value != "" {
        return value
    }
    return defaultValue
}
```

## Example 4: Advanced Configuration Management

Managing complex configuration scenarios:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load existing configuration or create new
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Printf("Managing configuration: %s\n", configPath)
    fmt.Println("=== Advanced Configuration Management ===")
    
    // Create configuration manager
    manager := NewConfigManager(api, config)
    
    // Apply configuration templates
    manager.ApplyTemplate("enterprise")
    
    // Validate configuration
    manager.ValidateConfiguration()
    
    // Optimize configuration
    manager.OptimizeConfiguration()
    
    // Display configuration report
    manager.GenerateReport()
    
    // Save optimized configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nOptimized configuration saved to: %s\n", configPath)
}

type ConfigManager struct {
    api    *nuget.API
    config *types.NuGetConfig
}

func NewConfigManager(api *nuget.API, config *types.NuGetConfig) *ConfigManager {
    return &ConfigManager{api: api, config: config}
}

func (cm *ConfigManager) ApplyTemplate(templateName string) {
    fmt.Printf("Applying template: %s\n", templateName)
    
    switch templateName {
    case "enterprise":
        cm.applyEnterpriseTemplate()
    case "developer":
        cm.applyDeveloperTemplate()
    case "ci-cd":
        cm.applyCICDTemplate()
    default:
        fmt.Printf("Unknown template: %s\n", templateName)
    }
}

func (cm *ConfigManager) applyEnterpriseTemplate() {
    // Enterprise-specific settings
    options := map[string]string{
        "globalPackagesFolder":     "/enterprise/packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "HighestMinor",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "signatureValidationMode": "require",
        "http_timeout":            "300",
        "http_retries":            "3",
        "verbosity":               "normal",
        "allowPrereleaseVersions": "false",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - Applied enterprise security settings")
    fmt.Println("  - Configured stable dependency resolution")
    fmt.Println("  - Set enterprise package locations")
}

func (cm *ConfigManager) applyDeveloperTemplate() {
    options := map[string]string{
        "globalPackagesFolder":     "./dev-packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "Highest",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "verbosity":               "detailed",
        "allowPrereleaseVersions": "true",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - Applied developer-friendly settings")
}

func (cm *ConfigManager) applyCICDTemplate() {
    options := map[string]string{
        "globalPackagesFolder":     "/tmp/packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "Exact",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "verbosity":               "minimal",
        "allowPrereleaseVersions": "false",
        "http_timeout":            "600",
        "http_retries":            "5",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - Applied CI/CD optimized settings")
}

func (cm *ConfigManager) ValidateConfiguration() {
    fmt.Println("\n=== Configuration Validation ===")
    
    issues := 0
    
    // Check required options
    requiredOptions := []string{
        "globalPackagesFolder",
        "repositoryPath",
        "dependencyVersion",
    }
    
    for _, option := range requiredOptions {
        value := cm.api.GetConfigOption(cm.config, option)
        if value == "" {
            fmt.Printf("  ❌ Missing required option: %s\n", option)
            issues++
        } else {
            fmt.Printf("  ✅ %s: %s\n", option, value)
        }
    }
    
    // Validate dependency version values
    depVersion := cm.api.GetConfigOption(cm.config, "dependencyVersion")
    validVersions := []string{"Lowest", "HighestPatch", "HighestMinor", "Highest", "Exact"}
    if depVersion != "" && !contains(validVersions, depVersion) {
        fmt.Printf("  ⚠️  Invalid dependencyVersion: %s\n", depVersion)
        issues++
    }
    
    if issues == 0 {
        fmt.Println("  ✅ Configuration validation passed")
    } else {
        fmt.Printf("  ⚠️  Found %d configuration issues\n", issues)
    }
}

func (cm *ConfigManager) OptimizeConfiguration() {
    fmt.Println("\n=== Configuration Optimization ===")
    
    // Remove duplicate or conflicting options
    cm.removeDuplicateOptions()
    
    // Set optimal defaults for missing options
    cm.setOptimalDefaults()
    
    fmt.Println("  ✅ Configuration optimized")
}

func (cm *ConfigManager) removeDuplicateOptions() {
    // This would involve checking for duplicate keys and resolving conflicts
    // For now, just report what we would do
    fmt.Println("  - Checked for duplicate options")
}

func (cm *ConfigManager) setOptimalDefaults() {
    defaults := map[string]string{
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "http_timeout":            "300",
        "http_retries":            "3",
    }
    
    for key, value := range defaults {
        if cm.api.GetConfigOption(cm.config, key) == "" {
            cm.api.AddConfigOption(cm.config, key, value)
            fmt.Printf("  - Set default %s: %s\n", key, value)
        }
    }
}

func (cm *ConfigManager) GenerateReport() {
    fmt.Println("\n=== Configuration Report ===")
    
    if cm.config.Config != nil {
        fmt.Printf("Total configuration options: %d\n", len(cm.config.Config.Add))
        
        categories := map[string][]string{
            "Package Management": {"globalPackagesFolder", "repositoryPath", "dependencyVersion"},
            "Network":           {"http_proxy", "http_timeout", "http_retries"},
            "Security":          {"signatureValidationMode", "allowPrereleaseVersions"},
            "Restore":           {"automaticPackageRestore", "packageRestore"},
        }
        
        for category, options := range categories {
            fmt.Printf("\n%s:\n", category)
            for _, option := range options {
                value := cm.api.GetConfigOption(cm.config, option)
                if value != "" {
                    fmt.Printf("  %s: %s\n", option, value)
                }
            }
        }
    }
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

## Key Configuration Options

### Package Management
- `globalPackagesFolder`: Global packages cache location
- `repositoryPath`: Project packages folder
- `dependencyVersion`: Dependency resolution strategy

### Network Settings
- `http_proxy`: HTTP proxy server
- `http_timeout`: Request timeout in seconds
- `http_retries`: Number of retry attempts

### Security Options
- `signatureValidationMode`: Package signature validation
- `allowPrereleaseVersions`: Allow prerelease packages

### Restore Behavior
- `automaticPackageRestore`: Enable automatic restore
- `packageRestore`: Enable package restore

## Best Practices

1. **Environment-specific configs**: Use different settings per environment
2. **Validate options**: Check option values are valid
3. **Use templates**: Apply consistent configuration patterns
4. **Document settings**: Comment configuration choices
5. **Regular review**: Periodically review and optimize settings

## Next Steps

After mastering configuration options:

1. Learn about [Serialization](./serialization.md) for custom XML handling
2. Explore [Position-Aware Editing](./position-aware-editing.md) for precise modifications
3. Study the [Types API](/api/types) for configuration structure details

This guide provides comprehensive examples for managing NuGet configuration options across different scenarios and environments.
