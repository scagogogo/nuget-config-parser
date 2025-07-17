# Creating Configurations

This example demonstrates how to create new NuGet configuration files from scratch using the NuGet Config Parser library.

## Overview

Creating configurations involves:
- Building configuration objects programmatically
- Setting up default package sources
- Configuring authentication and credentials
- Saving configurations to files
- Initializing project-specific settings

## Example 1: Create Basic Configuration

The simplest way to create a new configuration:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Create a default configuration
    config := api.CreateDefaultConfig()
    
    fmt.Println("Created default configuration:")
    fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))
    
    // Display the default source
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf("  Protocol: v%s\n", source.ProtocolVersion)
        }
    }
    
    // Save to file
    configPath := "NuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nConfiguration saved to: %s\n", configPath)
}
```

## Example 2: Create Custom Configuration

Build a configuration with custom package sources:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
    api := nuget.NewAPI()
    
    // Create empty configuration
    config := &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{},
        },
    }
    
    // Add multiple package sources
    sources := []struct {
        key     string
        value   string
        version string
    }{
        {"nuget.org", "https://api.nuget.org/v3/index.json", "3"},
        {"company-feed", "https://nuget.company.com/v3/index.json", "3"},
        {"local-dev", "/path/to/local/packages", ""},
        {"backup-feed", "https://backup.nuget.com/api/v2", "2"},
    }
    
    fmt.Println("Creating custom configuration with multiple sources:")
    
    for _, source := range sources {
        api.AddPackageSource(config, source.key, source.value, source.version)
        fmt.Printf("Added: %s -> %s\n", source.key, source.value)
    }
    
    // Set active package source
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // Add some global configuration options
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "defaultPushSource", "https://nuget.company.com/v3/index.json")
    
    // Disable the local development source by default
    api.DisablePackageSource(config, "local-dev")
    
    // Save configuration
    configPath := "CustomNuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nCustom configuration saved to: %s\n", configPath)
    
    // Display final configuration
    displayConfiguration(config)
}

func displayConfiguration(config *types.NuGetConfig) {
    fmt.Println("\n=== Final Configuration ===")
    
    fmt.Printf("Package Sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("\nActive Source: %s\n", config.ActivePackageSource.Add.Key)
    }
    
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("\nConfiguration Options (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
        fmt.Printf("\nDisabled Sources (%d):\n", len(config.DisabledPackageSources.Add))
        for _, disabled := range config.DisabledPackageSources.Add {
            fmt.Printf("  - %s\n", disabled.Key)
        }
    }
}
```

## Example 3: Create Configuration with Credentials

Build a configuration that includes authentication:

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
    
    // Add private package sources that require authentication
    privateSources := []struct {
        key      string
        url      string
        username string
        password string
    }{
        {"company-internal", "https://internal.company.com/nuget", "employee", "secret123"},
        {"partner-feed", "https://partner.company.com/nuget", "partner_user", "partner_pass"},
        {"azure-artifacts", "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json", "myuser", "pat_token"},
    }
    
    fmt.Println("Creating configuration with authenticated sources:")
    
    for _, source := range privateSources {
        // Add the package source
        api.AddPackageSource(config, source.key, source.url, "3")
        fmt.Printf("Added source: %s\n", source.key)
        
        // Add credentials for the source
        api.AddCredential(config, source.key, source.username, source.password)
        fmt.Printf("Added credentials for: %s\n", source.key)
    }
    
    // Add a public source without credentials
    api.AddPackageSource(config, "public-feed", "https://public.nuget.com/v3/index.json", "3")
    
    // Save configuration
    configPath := "AuthenticatedNuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nConfiguration with credentials saved to: %s\n", configPath)
    
    // Verify credentials were added
    fmt.Println("\nCredential verification:")
    for _, source := range privateSources {
        credential := api.GetCredential(config, source.key)
        if credential != nil {
            fmt.Printf("✅ %s has credentials configured\n", source.key)
        } else {
            fmt.Printf("❌ %s missing credentials\n", source.key)
        }
    }
    
    // Show XML output (be careful with passwords in production!)
    fmt.Println("\nGenerated XML preview:")
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("Failed to serialize: %v", err)
    } else {
        // In production, you'd want to mask passwords
        fmt.Println(xmlContent)
    }
}
```

## Example 4: Create Project-Specific Configuration

Create a configuration tailored for a specific project:

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
    
    // Get project information
    projectDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Failed to get current directory: %v", err)
    }
    
    projectName := filepath.Base(projectDir)
    
    fmt.Printf("Creating project-specific configuration for: %s\n", projectName)
    fmt.Printf("Project directory: %s\n", projectDir)
    
    // Create configuration optimized for this project
    config := api.CreateDefaultConfig()
    
    // Add project-specific package sources
    api.AddPackageSource(config, "project-local", "./packages", "")
    api.AddPackageSource(config, "project-cache", filepath.Join(projectDir, ".nuget", "cache"), "")
    
    // Configure project-specific settings
    packagesPath := filepath.Join(projectDir, "packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "globalPackagesFolder", packagesPath)
    
    // Set up development-friendly settings
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    // Add common development feeds
    api.AddPackageSource(config, "nuget-preview", "https://api.nuget.org/v3-flatcontainer", "3")
    api.AddPackageSource(config, "dotnet-core", "https://dotnetfeed.blob.core.windows.net/dotnet-core/index.json", "3")
    
    // Disable preview feeds by default
    api.DisablePackageSource(config, "nuget-preview")
    api.DisablePackageSource(config, "dotnet-core")
    
    // Set nuget.org as active source
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // Create packages directory if it doesn't exist
    if err := os.MkdirAll(packagesPath, 0755); err != nil {
        log.Printf("Warning: Failed to create packages directory: %v", err)
    }
    
    // Save configuration in project root
    configPath := filepath.Join(projectDir, "NuGet.Config")
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save project config: %v", err)
    }
    
    fmt.Printf("\nProject configuration created: %s\n", configPath)
    
    // Create a .gitignore entry for packages (if .gitignore exists)
    gitignorePath := filepath.Join(projectDir, ".gitignore")
    if _, err := os.Stat(gitignorePath); err == nil {
        addToGitignore(gitignorePath, "packages/")
        fmt.Println("Added packages/ to .gitignore")
    }
    
    // Display project configuration summary
    displayProjectSummary(config, projectName, configPath)
}

func addToGitignore(gitignorePath, entry string) {
    // Read existing .gitignore
    content, err := os.ReadFile(gitignorePath)
    if err != nil {
        return
    }
    
    // Check if entry already exists
    if strings.Contains(string(content), entry) {
        return
    }
    
    // Append entry
    file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return
    }
    defer file.Close()
    
    file.WriteString(fmt.Sprintf("\n# NuGet packages\n%s\n", entry))
}

func displayProjectSummary(config *types.NuGetConfig, projectName, configPath string) {
    fmt.Printf("\n=== Project Configuration Summary ===\n")
    fmt.Printf("Project: %s\n", projectName)
    fmt.Printf("Config file: %s\n", configPath)
    
    fmt.Printf("\nPackage Sources (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "enabled"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "disabled"
        }
        fmt.Printf("  - %s (%s): %s\n", source.Key, status, source.Value)
    }
    
    if config.Config != nil {
        fmt.Printf("\nProject Settings (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    fmt.Println("\nNext steps:")
    fmt.Println("1. Customize package sources for your project needs")
    fmt.Println("2. Add credentials for private feeds if needed")
    fmt.Println("3. Commit NuGet.Config to version control")
    fmt.Println("4. Share configuration with team members")
}
```

## Example 5: Initialize Default Configuration

Create a utility to initialize default configurations:

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
    
    // Configuration options
    configOptions := []struct {
        name        string
        path        string
        description string
    }{
        {"global", getGlobalConfigPath(), "Global user configuration"},
        {"project", "./NuGet.Config", "Project-specific configuration"},
        {"custom", "./CustomNuGet.Config", "Custom configuration file"},
    }
    
    fmt.Println("NuGet Configuration Initializer")
    fmt.Println("================================")
    
    for i, option := range configOptions {
        fmt.Printf("\n%d. %s\n", i+1, option.description)
        fmt.Printf("   Path: %s\n", option.path)
        
        // Check if configuration already exists
        if _, err := os.Stat(option.path); err == nil {
            fmt.Printf("   Status: ⚠️  Already exists\n")
            continue
        }
        
        // Create directory if needed
        dir := filepath.Dir(option.path)
        if err := os.MkdirAll(dir, 0755); err != nil {
            fmt.Printf("   Status: ❌ Failed to create directory: %v\n", err)
            continue
        }
        
        // Initialize configuration based on type
        var config *types.NuGetConfig
        
        switch option.name {
        case "global":
            config = createGlobalConfig(api)
        case "project":
            config = createProjectConfig(api)
        case "custom":
            config = createCustomConfig(api)
        default:
            config = api.CreateDefaultConfig()
        }
        
        // Save configuration
        err := api.SaveConfig(config, option.path)
        if err != nil {
            fmt.Printf("   Status: ❌ Failed to save: %v\n", err)
            continue
        }
        
        fmt.Printf("   Status: ✅ Created successfully\n")
        fmt.Printf("   Sources: %d package sources configured\n", len(config.PackageSources.Add))
    }
    
    fmt.Println("\nInitialization complete!")
    fmt.Println("\nNext steps:")
    fmt.Println("- Review and customize the generated configurations")
    fmt.Println("- Add credentials for private package sources")
    fmt.Println("- Test package restoration with your projects")
}

func getGlobalConfigPath() string {
    if home := os.Getenv("HOME"); home != "" {
        return filepath.Join(home, ".config", "NuGet", "NuGet.Config")
    }
    return "./GlobalNuGet.Config"
}

func createGlobalConfig(api *nuget.API) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // Add common global sources
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    api.AddPackageSource(config, "microsoft", "https://pkgs.dev.azure.com/dnceng/public/_packaging/dotnet-public/nuget/v3/index.json", "3")
    
    // Global settings
    if home := os.Getenv("HOME"); home != "" {
        api.AddConfigOption(config, "globalPackagesFolder", filepath.Join(home, ".nuget", "packages"))
    }
    
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    return config
}

func createProjectConfig(api *nuget.API) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // Project-specific settings
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "packagesConfigDirectoryPath", "./packages")
    
    // Add local package source
    api.AddPackageSource(config, "local", "./packages", "")
    api.DisablePackageSource(config, "local") // Disabled by default
    
    return config
}

func createCustomConfig(api *nuget.API) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // Add various package sources for different scenarios
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    api.AddPackageSource(config, "preview", "https://api.nuget.org/v3-flatcontainer", "3")
    api.AddPackageSource(config, "local-dev", "./local-packages", "")
    
    // Disable preview by default
    api.DisablePackageSource(config, "preview")
    api.DisablePackageSource(config, "local-dev")
    
    // Custom settings
    api.AddConfigOption(config, "globalPackagesFolder", "./custom-packages")
    api.AddConfigOption(config, "defaultPushSource", "https://api.nuget.org/v3/index.json")
    
    return config
}
```

## Key Concepts

### Configuration Structure

A NuGet configuration consists of:
- **Package Sources**: Where to find packages
- **Credentials**: Authentication for private sources
- **Config Options**: Global settings and preferences
- **Active Source**: Currently selected source
- **Disabled Sources**: Sources that are temporarily disabled

### Best Practices

1. **Start with defaults**: Use `CreateDefaultConfig()` as a foundation
2. **Add sources incrementally**: Build up configuration step by step
3. **Handle credentials securely**: Be careful with password storage
4. **Set appropriate permissions**: Ensure config files have correct permissions
5. **Validate before saving**: Check configuration validity
6. **Document sources**: Use meaningful names and comments

## Common Patterns

### Pattern 1: Incremental Building

```go
config := api.CreateDefaultConfig()
api.AddPackageSource(config, "source1", "url1", "3")
api.AddPackageSource(config, "source2", "url2", "3")
api.AddCredential(config, "source1", "user", "pass")
api.SaveConfig(config, "NuGet.Config")
```

### Pattern 2: Template-Based Creation

```go
func createEnterpriseConfig(api *nuget.API, companyDomain string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // Add company-specific sources
    api.AddPackageSource(config, "company", fmt.Sprintf("https://nuget.%s", companyDomain), "3")
    api.AddPackageSource(config, "company-preview", fmt.Sprintf("https://preview.nuget.%s", companyDomain), "3")
    
    // Configure enterprise settings
    api.AddConfigOption(config, "defaultPushSource", fmt.Sprintf("https://nuget.%s", companyDomain))
    
    return config
}
```

### Pattern 3: Environment-Based Configuration

```go
func createConfigForEnvironment(api *nuget.API, env string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    switch env {
    case "development":
        api.AddPackageSource(config, "local", "./packages", "")
        api.AddPackageSource(config, "dev-feed", "https://dev.nuget.com", "3")
    case "staging":
        api.AddPackageSource(config, "staging-feed", "https://staging.nuget.com", "3")
    case "production":
        api.AddPackageSource(config, "prod-feed", "https://prod.nuget.com", "3")
    }
    
    return config
}
```

## Next Steps

After mastering configuration creation:

1. Learn about [Modifying Configs](./modifying-configs.md) to update existing configurations
2. Explore [Package Sources](./package-sources.md) for advanced source management
3. Study [Credentials](./credentials.md) for secure authentication handling

This guide provides comprehensive examples for creating NuGet configuration files from scratch, covering various scenarios and best practices.
