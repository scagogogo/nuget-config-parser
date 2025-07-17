# Credentials Management

This example demonstrates how to manage authentication credentials for private NuGet package sources.

## Overview

Credential management involves:
- Adding authentication for private package sources
- Managing usernames and passwords securely
- Handling different authentication methods
- Removing and updating credentials
- Best practices for credential security

## Example 1: Basic Credential Management

Adding and managing basic credentials:

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
    
    fmt.Println("=== Basic Credential Management ===")
    
    // Add private package sources
    privateSources := []struct {
        key      string
        url      string
        username string
        password string
    }{
        {"company-private", "https://nuget.company.com/v3/index.json", "employee", "company_password"},
        {"azure-artifacts", "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json", "myuser", "pat_token"},
        {"github-packages", "https://nuget.pkg.github.com/myorg/index.json", "github_user", "ghp_token"},
    }
    
    // Add sources and credentials
    for _, source := range privateSources {
        // Add the package source
        api.AddPackageSource(config, source.key, source.url, "3")
        fmt.Printf("Added source: %s\n", source.key)
        
        // Add credentials for the source
        api.AddCredential(config, source.key, source.username, source.password)
        fmt.Printf("  - Added credentials for: %s\n", source.key)
    }
    
    // Verify credentials were added
    fmt.Println("\n=== Credential Verification ===")
    for _, source := range privateSources {
        credential := api.GetCredential(config, source.key)
        if credential != nil {
            fmt.Printf("✅ %s has credentials configured\n", source.key)
            
            // Display credential details (be careful with passwords!)
            for _, cred := range credential.Add {
                if cred.Key == "Username" {
                    fmt.Printf("   Username: %s\n", cred.Value)
                } else if cred.Key == "ClearTextPassword" {
                    fmt.Printf("   Password: %s\n", maskPassword(cred.Value))
                }
            }
        } else {
            fmt.Printf("❌ %s missing credentials\n", source.key)
        }
    }
    
    // Save configuration
    err := api.SaveConfig(config, "CredentialsExample.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Println("\nConfiguration with credentials saved!")
}

func maskPassword(password string) string {
    if len(password) <= 4 {
        return "****"
    }
    return password[:2] + "****" + password[len(password)-2:]
}
```

## Example 2: Environment-Based Credentials

Using environment variables for secure credential management:

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
    
    fmt.Println("=== Environment-Based Credentials ===")
    
    // Define sources with environment variable mappings
    sources := []struct {
        key         string
        url         string
        userEnvVar  string
        passEnvVar  string
        required    bool
    }{
        {"company-feed", "https://nuget.company.com/v3/index.json", "COMPANY_NUGET_USER", "COMPANY_NUGET_PASS", true},
        {"azure-devops", "https://pkgs.dev.azure.com/myorg/_packaging/feed/nuget/v3/index.json", "AZURE_DEVOPS_USER", "AZURE_DEVOPS_PAT", true},
        {"github-packages", "https://nuget.pkg.github.com/myorg/index.json", "GITHUB_USER", "GITHUB_TOKEN", false},
        {"private-registry", "https://private.registry.com/nuget", "REGISTRY_USER", "REGISTRY_PASS", false},
    }
    
    // Process each source
    for _, source := range sources {
        // Add the package source
        api.AddPackageSource(config, source.key, source.url, "3")
        fmt.Printf("Added source: %s\n", source.key)
        
        // Get credentials from environment
        username := os.Getenv(source.userEnvVar)
        password := os.Getenv(source.passEnvVar)
        
        if username != "" && password != "" {
            // Add credentials
            api.AddCredential(config, source.key, username, password)
            fmt.Printf("  ✅ Added credentials from environment\n")
        } else if source.required {
            fmt.Printf("  ❌ Required credentials missing! Set %s and %s\n", source.userEnvVar, source.passEnvVar)
            log.Fatalf("Missing required credentials for %s", source.key)
        } else {
            fmt.Printf("  ⚠️  Optional credentials not provided\n")
            // Disable source if no credentials
            api.DisablePackageSource(config, source.key)
        }
    }
    
    // Display environment variable instructions
    fmt.Println("\n=== Environment Variable Setup ===")
    fmt.Println("To configure credentials, set these environment variables:")
    for _, source := range sources {
        status := "optional"
        if source.required {
            status = "required"
        }
        fmt.Printf("  %s (%s):\n", source.key, status)
        fmt.Printf("    export %s=\"your_username\"\n", source.userEnvVar)
        fmt.Printf("    export %s=\"your_password\"\n", source.passEnvVar)
        fmt.Println()
    }
    
    // Save configuration
    err := api.SaveConfig(config, "EnvCredentials.Config")
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Println("Configuration saved with environment-based credentials!")
}
```

## Example 3: Credential Updates and Removal

Managing credential lifecycle:

```go
package main

import (
    "fmt"
    "log"
    
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
    
    fmt.Printf("Managing credentials in: %s\n", configPath)
    fmt.Println("=== Credential Lifecycle Management ===")
    
    // Add initial credentials
    api.AddPackageSource(config, "test-feed", "https://test.company.com/nuget", "3")
    api.AddCredential(config, "test-feed", "old_user", "old_password")
    fmt.Println("Added initial credentials")
    
    // Display current credentials
    displayCredentials(api, config, "test-feed")
    
    // Update credentials (replace existing)
    fmt.Println("\nUpdating credentials...")
    api.AddCredential(config, "test-feed", "new_user", "new_password")
    displayCredentials(api, config, "test-feed")
    
    // Add another source with credentials
    api.AddPackageSource(config, "temp-feed", "https://temp.company.com/nuget", "3")
    api.AddCredential(config, "temp-feed", "temp_user", "temp_pass")
    fmt.Println("\nAdded temporary source with credentials")
    
    // List all sources with credentials
    fmt.Println("\n=== All Sources with Credentials ===")
    listAllCredentials(api, config)
    
    // Remove credentials for specific source
    fmt.Println("\nRemoving credentials for temp-feed...")
    removed := api.RemoveCredential(config, "temp-feed")
    if removed {
        fmt.Println("✅ Credentials removed successfully")
    } else {
        fmt.Println("❌ No credentials found to remove")
    }
    
    // Verify removal
    fmt.Println("\n=== After Credential Removal ===")
    listAllCredentials(api, config)
    
    // Clean up - remove sources without credentials
    fmt.Println("\nCleaning up sources without credentials...")
    cleanupSourcesWithoutCredentials(api, config)
    
    // Save final configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    fmt.Printf("\nFinal configuration saved to: %s\n", configPath)
}

func displayCredentials(api *nuget.API, config *types.NuGetConfig, sourceKey string) {
    credential := api.GetCredential(config, sourceKey)
    if credential != nil {
        fmt.Printf("Credentials for %s:\n", sourceKey)
        for _, cred := range credential.Add {
            value := cred.Value
            if cred.Key == "ClearTextPassword" || cred.Key == "Password" {
                value = maskPassword(value)
            }
            fmt.Printf("  %s: %s\n", cred.Key, value)
        }
    } else {
        fmt.Printf("No credentials found for %s\n", sourceKey)
    }
}

func listAllCredentials(api *nuget.API, config *types.NuGetConfig) {
    credentialCount := 0
    
    for _, source := range config.PackageSources.Add {
        credential := api.GetCredential(config, source.Key)
        if credential != nil {
            credentialCount++
            fmt.Printf("  %s: Has credentials\n", source.Key)
        } else {
            fmt.Printf("  %s: No credentials\n", source.Key)
        }
    }
    
    fmt.Printf("\nTotal sources with credentials: %d\n", credentialCount)
}

func cleanupSourcesWithoutCredentials(api *nuget.API, config *types.NuGetConfig) {
    var sourcesToRemove []string
    
    for _, source := range config.PackageSources.Add {
        // Skip public sources that don't need credentials
        if source.Key == "nuget.org" {
            continue
        }
        
        credential := api.GetCredential(config, source.Key)
        if credential == nil && isPrivateSource(source.Value) {
            sourcesToRemove = append(sourcesToRemove, source.Key)
        }
    }
    
    for _, sourceKey := range sourcesToRemove {
        removed := api.RemovePackageSource(config, sourceKey)
        if removed {
            fmt.Printf("  Removed source without credentials: %s\n", sourceKey)
        }
    }
}

func isPrivateSource(url string) bool {
    // Simple heuristic to identify private sources
    privateIndicators := []string{
        "company.com",
        "internal",
        "private",
        "pkgs.dev.azure.com",
        "nuget.pkg.github.com",
    }
    
    for _, indicator := range privateIndicators {
        if strings.Contains(url, indicator) {
            return true
        }
    }
    
    return false
}

func maskPassword(password string) string {
    if len(password) <= 4 {
        return "****"
    }
    return password[:2] + "****" + password[len(password)-2:]
}
```

## Example 4: Secure Credential Handling

Best practices for secure credential management:

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
    "syscall"
    
    "golang.org/x/term"
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== Secure Credential Handling ===")
    
    // Interactive credential setup
    setupCredentialsInteractively(api, config)
    
    // Validate credentials
    validateCredentials(api, config)
    
    // Save with warning about security
    saveConfigSecurely(api, config)
}

func setupCredentialsInteractively(api *nuget.API, config *types.NuGetConfig) {
    reader := bufio.NewReader(os.Stdin)
    
    fmt.Println("Setting up private package source credentials...")
    
    // Get source information
    fmt.Print("Enter source name: ")
    sourceName, _ := reader.ReadString('\n')
    sourceName = strings.TrimSpace(sourceName)
    
    fmt.Print("Enter source URL: ")
    sourceURL, _ := reader.ReadString('\n')
    sourceURL = strings.TrimSpace(sourceURL)
    
    // Add the package source
    api.AddPackageSource(config, sourceName, sourceURL, "3")
    
    // Get credentials securely
    fmt.Print("Enter username: ")
    username, _ := reader.ReadString('\n')
    username = strings.TrimSpace(username)
    
    fmt.Print("Enter password (hidden): ")
    passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatalf("Failed to read password: %v", err)
    }
    password := string(passwordBytes)
    fmt.Println() // New line after hidden input
    
    // Add credentials
    api.AddCredential(config, sourceName, username, password)
    
    fmt.Printf("✅ Credentials added for source: %s\n", sourceName)
    
    // Security reminder
    fmt.Println("\n⚠️  Security Reminder:")
    fmt.Println("  - Passwords are stored in plain text in the config file")
    fmt.Println("  - Consider using environment variables for production")
    fmt.Println("  - Ensure config file has appropriate permissions")
}

func validateCredentials(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== Credential Validation ===")
    
    for _, source := range config.PackageSources.Add {
        credential := api.GetCredential(config, source.Key)
        if credential != nil {
            fmt.Printf("Validating credentials for: %s\n", source.Key)
            
            var username, password string
            for _, cred := range credential.Add {
                if cred.Key == "Username" {
                    username = cred.Value
                } else if cred.Key == "ClearTextPassword" {
                    password = cred.Value
                }
            }
            
            // Basic validation
            if username == "" {
                fmt.Printf("  ❌ Username is empty\n")
            } else {
                fmt.Printf("  ✅ Username: %s\n", username)
            }
            
            if password == "" {
                fmt.Printf("  ❌ Password is empty\n")
            } else {
                fmt.Printf("  ✅ Password: %s\n", maskPassword(password))
            }
            
            // Additional security checks
            if len(password) < 8 {
                fmt.Printf("  ⚠️  Password is shorter than 8 characters\n")
            }
            
            if password == "password" || password == "123456" {
                fmt.Printf("  ❌ Password appears to be weak or default\n")
            }
        }
    }
}

func saveConfigSecurely(api *nuget.API, config *types.NuGetConfig) {
    configPath := "SecureCredentials.Config"
    
    // Save configuration
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
    
    // Set restrictive file permissions (Unix-like systems)
    err = os.Chmod(configPath, 0600) // Read/write for owner only
    if err != nil {
        fmt.Printf("⚠️  Could not set restrictive permissions: %v\n", err)
    } else {
        fmt.Printf("✅ Set restrictive file permissions (600)\n")
    }
    
    fmt.Printf("Configuration saved securely to: %s\n", configPath)
    
    // Security recommendations
    fmt.Println("\n=== Security Recommendations ===")
    fmt.Println("1. Use environment variables for credentials in production")
    fmt.Println("2. Ensure config files have restrictive permissions")
    fmt.Println("3. Consider using credential managers or key vaults")
    fmt.Println("4. Regularly rotate passwords and tokens")
    fmt.Println("5. Use personal access tokens instead of passwords when possible")
    fmt.Println("6. Never commit config files with credentials to version control")
}

func maskPassword(password string) string {
    if len(password) <= 4 {
        return "****"
    }
    return password[:2] + "****" + password[len(password)-2:]
}
```

## Key Concepts

### Credential Types

1. **Username/Password**: Traditional authentication
2. **Personal Access Tokens**: Modern token-based auth
3. **API Keys**: Service-specific authentication
4. **Encrypted Passwords**: NuGet-encrypted credentials

### Security Considerations

1. **Plain Text Storage**: Credentials are stored in plain text
2. **File Permissions**: Restrict access to config files
3. **Environment Variables**: Use for sensitive credentials
4. **Version Control**: Never commit credentials
5. **Token Rotation**: Regularly update credentials

### Best Practices

1. **Use Environment Variables**: For production deployments
2. **Restrict File Permissions**: Limit access to config files
3. **Use Tokens**: Prefer PATs over passwords
4. **Regular Rotation**: Update credentials periodically
5. **Credential Managers**: Use system credential stores
6. **Separate Configs**: Different configs for different environments

## Next Steps

After mastering credential management:

1. Learn about [Config Options](./config-options.md) for global settings
2. Explore [Serialization](./serialization.md) for custom XML handling
3. Study [Position-Aware Editing](./position-aware-editing.md) for precise modifications

This guide provides comprehensive examples for managing NuGet package source credentials securely and effectively.
