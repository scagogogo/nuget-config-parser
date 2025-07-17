# Constants API

The `pkg/constants` package defines constants and default values used throughout the NuGet Config Parser library.

## Overview

The Constants API provides:
- Default configuration file names and paths
- Standard NuGet protocol versions
- Platform-specific configuration locations
- Helper functions for path resolution

## File and Directory Constants

### Configuration File Names

```go
const (
    // DefaultNuGetConfigFilename is the default NuGet configuration file name
    DefaultNuGetConfigFilename = "NuGet.Config"
    
    // GlobalFolderName is the global configuration folder name
    GlobalFolderName = "NuGet"
    
    // FeedNamePrefix is the package source name prefix
    FeedNamePrefix = "PackageSource"
)
```

**Usage:**
- `DefaultNuGetConfigFilename`: Standard name for NuGet configuration files
- `GlobalFolderName`: Directory name for global NuGet configurations
- `FeedNamePrefix`: Prefix used when generating automatic package source names

**Example:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// Create a configuration file path
configPath := filepath.Join("/etc", constants.GlobalFolderName, constants.DefaultNuGetConfigFilename)
// Result: "/etc/NuGet/NuGet.Config"

// Generate a package source name
sourceName := fmt.Sprintf("%s_%d", constants.FeedNamePrefix, 1)
// Result: "PackageSource_1"
```

## Package Source Constants

### Default Package Sources

```go
const (
    // DefaultPackageSource is the default package source URL
    DefaultPackageSource = "https://api.nuget.org/v3/index.json"
)
```

**Usage:**
- `DefaultPackageSource`: The official NuGet.org package source URL

**Example:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// Create a default package source
defaultSource := types.PackageSource{
    Key:   "nuget.org",
    Value: constants.DefaultPackageSource,
}
```

## Protocol Version Constants

### NuGet API Versions

```go
const (
    // NuGetV3APIProtocolVersion is the NuGet V3 API protocol version
    NuGetV3APIProtocolVersion = "3"
    
    // NuGetV2APIProtocolVersion is the NuGet V2 API protocol version  
    NuGetV2APIProtocolVersion = "2"
)
```

**Usage:**
- `NuGetV3APIProtocolVersion`: Use for modern NuGet V3 API sources
- `NuGetV2APIProtocolVersion`: Use for legacy NuGet V2 API sources

**Example:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// Create a V3 API package source
v3Source := types.PackageSource{
    Key:             "nuget.org",
    Value:           "https://api.nuget.org/v3/index.json",
    ProtocolVersion: constants.NuGetV3APIProtocolVersion,
}

// Create a V2 API package source
v2Source := types.PackageSource{
    Key:             "legacy-feed",
    Value:           "https://legacy.nuget.org/api/v2",
    ProtocolVersion: constants.NuGetV2APIProtocolVersion,
}
```

## Path Resolution Functions

### GetDefaultConfigLocations

```go
func GetDefaultConfigLocations() []string
```

Returns the default search paths for NuGet configuration files, ordered by priority.

**Returns:**
- `[]string`: List of default configuration file paths

**Search Order:**
1. Current directory: `./NuGet.Config`
2. Parent directory: `../NuGet.Config`
3. User-specific configuration directory
4. System-wide configuration directory

**Platform-specific User Paths:**
- **Windows**: `%APPDATA%\NuGet\NuGet.Config`
- **macOS**: `~/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `~/.config/NuGet/NuGet.Config` (respects `XDG_CONFIG_HOME`)

**Platform-specific System Paths:**
- **Windows**: `%ProgramData%\NuGet\NuGet.Config`
- **macOS**: `/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `/etc/NuGet/NuGet.Config`

**Example:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// Get all default configuration locations
locations := constants.GetDefaultConfigLocations()

fmt.Println("NuGet configuration file search paths:")
for i, location := range locations {
    fmt.Printf("%d. %s\n", i+1, location)
}

// Find the first existing configuration file
for _, location := range locations {
    if utils.FileExists(location) {
        fmt.Printf("Found configuration file: %s\n", location)
        break
    }
}
```

**Example Output (Linux):**
```
NuGet configuration file search paths:
1. NuGet.Config
2. ../NuGet.Config
3. /home/user/.config/NuGet/NuGet.Config
4. /etc/NuGet/NuGet.Config
```

**Example Output (Windows):**
```
NuGet configuration file search paths:
1. NuGet.Config
2. ..\NuGet.Config
3. C:\Users\user\AppData\Roaming\NuGet\NuGet.Config
4. C:\ProgramData\NuGet\NuGet.Config
```

**Example Output (macOS):**
```
NuGet configuration file search paths:
1. NuGet.Config
2. ../NuGet.Config
3. /Users/user/Library/Application Support/NuGet/NuGet.Config
4. /Library/Application Support/NuGet/NuGet.Config
```

## Platform-Specific Behavior

### Windows

```go
// User configuration directory: %APPDATA%
// System configuration directory: %ProgramData%
```

**Environment Variables:**
- `APPDATA`: User application data directory
- `ProgramData`: System-wide application data directory

**Example Paths:**
- User: `C:\Users\username\AppData\Roaming\NuGet\NuGet.Config`
- System: `C:\ProgramData\NuGet\NuGet.Config`

### macOS

```go
// User configuration directory: ~/Library/Application Support/
// System configuration directory: /Library/Application Support/
```

**Example Paths:**
- User: `/Users/username/Library/Application Support/NuGet/NuGet.Config`
- System: `/Library/Application Support/NuGet/NuGet.Config`

### Linux/Unix

```go
// User configuration directory: ~/.config/ (or $XDG_CONFIG_HOME)
// System configuration directory: /etc/
```

**Environment Variables:**
- `XDG_CONFIG_HOME`: User configuration directory (fallback to `~/.config`)

**Example Paths:**
- User: `/home/username/.config/NuGet/NuGet.Config`
- System: `/etc/NuGet/NuGet.Config`

## Usage Examples

### Finding Configuration Files

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/utils"
)

func findConfigFile() (string, error) {
    locations := constants.GetDefaultConfigLocations()
    
    for _, location := range locations {
        if utils.FileExists(location) {
            return location, nil
        }
    }
    
    return "", fmt.Errorf("no configuration file found in default locations")
}

func main() {
    configPath, err := findConfigFile()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        // Create default configuration in current directory
        defaultPath := constants.DefaultNuGetConfigFilename
        fmt.Printf("Creating default configuration: %s\n", defaultPath)
        
        // Create configuration...
    } else {
        fmt.Printf("Using configuration file: %s\n", configPath)
    }
}
```

### Creating Default Configuration

```go
package main

import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func createDefaultConfig() *types.NuGetConfig {
    // Create default package source
    defaultSource := types.PackageSource{
        Key:             "nuget.org",
        Value:           constants.DefaultPackageSource,
        ProtocolVersion: constants.NuGetV3APIProtocolVersion,
    }
    
    return &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{defaultSource},
        },
        ActivePackageSource: &types.ActivePackageSource{
            Add: defaultSource,
        },
    }
}
```

### Platform-Specific Configuration

```go
package main

import (
    "fmt"
    "runtime"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
)

func showPlatformInfo() {
    fmt.Printf("Operating System: %s\n", runtime.GOOS)
    fmt.Printf("Default config filename: %s\n", constants.DefaultNuGetConfigFilename)
    fmt.Printf("Global folder name: %s\n", constants.GlobalFolderName)
    
    locations := constants.GetDefaultConfigLocations()
    fmt.Println("\nDefault search locations:")
    for i, location := range locations {
        fmt.Printf("%d. %s\n", i+1, location)
    }
}
```

### Custom Configuration Paths

```go
package main

import (
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
)

func getCustomConfigPaths() []string {
    var paths []string
    
    // Add standard locations
    paths = append(paths, constants.GetDefaultConfigLocations()...)
    
    // Add custom locations
    customLocations := []string{
        "/opt/nuget/NuGet.Config",
        "/usr/local/etc/nuget/NuGet.Config",
        filepath.Join(os.Getenv("HOME"), "custom-nuget.config"),
    }
    
    paths = append(paths, customLocations...)
    
    return paths
}
```

## Integration with Other Packages

### With Finder Package

```go
import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/finder"
)

func createFinderWithDefaults() *finder.ConfigFinder {
    defaultPaths := constants.GetDefaultConfigLocations()
    return finder.NewConfigFinderWithPaths(defaultPaths)
}
```

### With Manager Package

```go
import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/manager"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func initializeWithDefaults() *types.NuGetConfig {
    mgr := manager.NewConfigManager()
    
    config := &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{
                {
                    Key:             "nuget.org",
                    Value:           constants.DefaultPackageSource,
                    ProtocolVersion: constants.NuGetV3APIProtocolVersion,
                },
            },
        },
    }
    
    return config
}
```

## Best Practices

1. **Use constants for consistency**: Always use the predefined constants instead of hardcoding values
2. **Respect platform differences**: Use `GetDefaultConfigLocations()` for cross-platform compatibility
3. **Check file existence**: Always verify that configuration files exist before attempting to parse them
4. **Handle missing configurations**: Provide fallbacks when default configurations are not found
5. **Use appropriate protocol versions**: Choose V3 for modern sources, V2 for legacy compatibility
6. **Follow naming conventions**: Use standard names for configuration files and directories

## Thread Safety

All constants and functions in this package are thread-safe and can be used concurrently.
