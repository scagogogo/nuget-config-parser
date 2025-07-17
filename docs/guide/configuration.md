# Configuration

This guide explains the structure and components of NuGet configuration files and how to work with them using the NuGet Config Parser library.

## Overview

NuGet configuration files (`NuGet.Config`) are XML files that control various aspects of NuGet behavior, including:

- Package source locations
- Authentication credentials
- Global settings and preferences
- Package restore behavior
- Proxy settings

## Configuration File Structure

A typical NuGet.Config file has the following structure:

```xml
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="C:\LocalPackages" />
  </packageSources>
  
  <packageSourceCredentials>
    <MyPrivateSource>
      <add key="Username" value="myuser" />
      <add key="ClearTextPassword" value="mypass" />
    </MyPrivateSource>
  </packageSourceCredentials>
  
  <disabledPackageSources>
    <add key="local" value="true" />
  </disabledPackageSources>
  
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  
  <config>
    <add key="globalPackagesFolder" value="C:\packages" />
    <add key="repositoryPath" value=".\packages" />
    <add key="defaultPushSource" value="https://api.nuget.org/v3/index.json" />
  </config>
</configuration>
```

## Configuration Sections

### Package Sources

The `<packageSources>` section defines where NuGet looks for packages:

```xml
<packageSources>
  <clear />  <!-- Optional: clear all inherited sources -->
  <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
  <add key="company-feed" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
  <add key="local-packages" value="C:\LocalPackages" />
</packageSources>
```

**Attributes:**
- `key`: Unique identifier for the source
- `value`: URL or file path to the package source
- `protocolVersion`: NuGet protocol version ("2" or "3")

### Package Source Credentials

The `<packageSourceCredentials>` section stores authentication information:

```xml
<packageSourceCredentials>
  <MyPrivateSource>
    <add key="Username" value="myuser" />
    <add key="ClearTextPassword" value="mypass" />
  </MyPrivateSource>
  <AnotherSource>
    <add key="Username" value="user2" />
    <add key="Password" value="encrypted_password" />
  </AnotherSource>
</packageSourceCredentials>
```

**Credential Types:**
- `Username`: Authentication username
- `Password`: Encrypted password
- `ClearTextPassword`: Plain text password (not recommended for production)

### Disabled Package Sources

The `<disabledPackageSources>` section lists temporarily disabled sources:

```xml
<disabledPackageSources>
  <add key="local-packages" value="true" />
  <add key="old-feed" value="true" />
</disabledPackageSources>
```

### Active Package Source

The `<activePackageSource>` section specifies the currently active source:

```xml
<activePackageSource>
  <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
</activePackageSource>
```

### Global Configuration

The `<config>` section contains global NuGet settings:

```xml
<config>
  <add key="globalPackagesFolder" value="C:\packages" />
  <add key="repositoryPath" value=".\packages" />
  <add key="defaultPushSource" value="https://api.nuget.org/v3/index.json" />
  <add key="dependencyVersion" value="Highest" />
  <add key="http_proxy" value="http://proxy.company.com:8080" />
  <add key="http_proxy.user" value="proxyuser" />
  <add key="http_proxy.password" value="proxypass" />
</config>
```

**Common Configuration Keys:**
- `globalPackagesFolder`: Global packages cache location
- `repositoryPath`: Project packages folder
- `defaultPushSource`: Default source for package publishing
- `dependencyVersion`: Default dependency version resolution
- `http_proxy`: HTTP proxy server
- `automaticPackageRestore`: Enable automatic package restore

## Configuration Hierarchy

NuGet uses a hierarchical configuration system where settings are inherited and can be overridden:

1. **Computer-level**: System-wide settings
2. **User-level**: User-specific settings
3. **Solution-level**: Solution-specific settings
4. **Project-level**: Project-specific settings

### Search Order

The library searches for configuration files in this order:

1. Current directory: `./NuGet.Config`
2. Parent directories (walking up the tree)
3. User configuration directory
4. System configuration directory

### Platform-Specific Locations

**Windows:**
- User: `%APPDATA%\NuGet\NuGet.Config`
- System: `%ProgramData%\NuGet\NuGet.Config`

**macOS:**
- User: `~/Library/Application Support/NuGet/NuGet.Config`
- System: `/Library/Application Support/NuGet/NuGet.Config`

**Linux:**
- User: `~/.config/NuGet/NuGet.Config`
- System: `/etc/NuGet/NuGet.Config`

## Working with Configuration

### Reading Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Find and parse configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    fmt.Printf("Loaded configuration from: %s\n", configPath)
    
    // Access package sources
    for _, source := range config.PackageSources.Add {
        fmt.Printf("Source: %s -> %s\n", source.Key, source.Value)
    }
    
    // Access configuration options
    if config.Config != nil {
        for _, option := range config.Config.Add {
            fmt.Printf("Setting: %s = %s\n", option.Key, option.Value)
        }
    }
}
```

### Modifying Configuration

```go
package main

import (
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Load existing configuration
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // Create default if not found
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    // Add package source
    api.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    
    // Add credentials
    api.AddCredential(config, "company", "myuser", "mypass")
    
    // Configure global settings
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    
    // Set active source
    api.SetActivePackageSource(config, "company", "https://nuget.company.com/v3/index.json")
    
    // Save configuration
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("Failed to save config: %v", err)
    }
}
```

## Best Practices

### Security

1. **Avoid plain text passwords**: Use encrypted passwords when possible
2. **Secure file permissions**: Ensure configuration files have appropriate permissions
3. **Environment variables**: Use environment variables for sensitive information
4. **Credential management**: Consider using credential managers for authentication

### Organization

1. **Hierarchical configuration**: Use project-level configs for project-specific settings
2. **Consistent naming**: Use descriptive names for package sources
3. **Documentation**: Comment configuration files when possible
4. **Version control**: Include project-level configs in version control

### Performance

1. **Minimize sources**: Only include necessary package sources
2. **Protocol versions**: Use appropriate protocol versions for sources
3. **Local caching**: Configure appropriate cache locations
4. **Disable unused sources**: Disable sources that aren't needed

## Common Configuration Patterns

### Enterprise Setup

```xml
<configuration>
  <packageSources>
    <clear />
    <add key="company-internal" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
  </packageSources>
  
  <config>
    <add key="globalPackagesFolder" value="C:\CompanyPackages" />
    <add key="defaultPushSource" value="https://nuget.company.com/v3/index.json" />
    <add key="http_proxy" value="http://proxy.company.com:8080" />
  </config>
</configuration>
```

### Development Setup

```xml
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local-dev" value="./local-packages" />
    <add key="preview" value="https://api.nuget.org/v3-flatcontainer" protocolVersion="3" />
  </packageSources>
  
  <disabledPackageSources>
    <add key="preview" value="true" />
  </disabledPackageSources>
  
  <config>
    <add key="repositoryPath" value="./packages" />
    <add key="dependencyVersion" value="Highest" />
  </config>
</configuration>
```

### CI/CD Setup

```xml
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="build-artifacts" value="https://artifacts.company.com/nuget" protocolVersion="3" />
  </packageSources>
  
  <config>
    <add key="globalPackagesFolder" value="/tmp/packages" />
    <add key="automaticPackageRestore" value="true" />
  </config>
</configuration>
```

## Troubleshooting

### Common Issues

1. **File not found**: Check file paths and permissions
2. **Invalid XML**: Validate XML structure and encoding
3. **Authentication failures**: Verify credentials and source URLs
4. **Source conflicts**: Check for duplicate source keys
5. **Permission errors**: Ensure proper file and directory permissions

### Debugging Configuration

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
    
    // Find all configuration files
    configPaths := api.FindAllConfigFiles()
    fmt.Printf("Found %d configuration files:\n", len(configPaths))
    
    for i, path := range configPaths {
        fmt.Printf("%d. %s\n", i+1, path)
        
        config, err := api.ParseFromFile(path)
        if err != nil {
            if errors.IsParseError(err) {
                fmt.Printf("   Parse error: %v\n", err)
            } else {
                fmt.Printf("   Error: %v\n", err)
            }
            continue
        }
        
        fmt.Printf("   Sources: %d\n", len(config.PackageSources.Add))
        fmt.Printf("   Settings: %d\n", len(config.Config.Add))
    }
}
```

## Next Steps

- Learn about [Position-Aware Editing](./position-aware-editing.md) for advanced configuration modification
- Explore the [API Reference](/api/) for detailed method documentation
- Check out [Examples](/examples/) for practical usage scenarios

This configuration guide provides a comprehensive understanding of NuGet configuration files and how to work with them effectively using the library.
