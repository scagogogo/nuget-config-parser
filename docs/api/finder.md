# Finder API

The `pkg/finder` package provides configuration file discovery functionality across different platforms and directory structures.

## Overview

The Finder API is responsible for:
- Locating NuGet configuration files in standard locations
- Supporting platform-specific configuration paths
- Providing flexible search strategies
- Handling project-level and global configuration discovery

## Types

### ConfigFinder

```go
type ConfigFinder struct {
    SearchPaths []string
}
```

The main finder type that handles configuration file discovery.

**Fields:**
- `SearchPaths`: List of paths to search for configuration files

## Constructors

### NewConfigFinder

```go
func NewConfigFinder() *ConfigFinder
```

Creates a new configuration finder with default search paths.

**Returns:**
- `*ConfigFinder`: New finder instance with platform-specific default paths

**Example:**
```go
finder := finder.NewConfigFinder()
configPath, err := finder.FindConfigFile()
if err != nil {
    log.Printf("No configuration file found: %v", err)
} else {
    fmt.Printf("Found configuration: %s\n", configPath)
}
```

### NewConfigFinderWithPaths

```go
func NewConfigFinderWithPaths(searchPaths []string) *ConfigFinder
```

Creates a new configuration finder with custom search paths.

**Parameters:**
- `searchPaths` ([]string): Custom list of paths to search

**Returns:**
- `*ConfigFinder`: New finder instance with specified paths

**Example:**
```go
customPaths := []string{
    "/custom/path/NuGet.Config",
    "/another/path/NuGet.Config",
}
finder := finder.NewConfigFinderWithPaths(customPaths)
```

## Discovery Methods

### FindConfigFile

```go
func (f *ConfigFinder) FindConfigFile() (string, error)
```

Finds the first available NuGet configuration file.

**Returns:**
- `string`: Path to the found configuration file
- `error`: Error if no configuration file is found

**Search Order:**
1. Current directory (`./NuGet.Config`)
2. Parent directories (walking up the tree)
3. User-specific configuration directory
4. System-wide configuration directory

**Platform-specific paths:**
- **Windows**: `%APPDATA%\NuGet\NuGet.Config`, `%ProgramData%\NuGet\NuGet.Config`
- **macOS**: `~/Library/Application Support/NuGet/NuGet.Config`, `/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `~/.config/NuGet/NuGet.Config`, `/etc/NuGet/NuGet.Config`

**Example:**
```go
finder := finder.NewConfigFinder()
configPath, err := finder.FindConfigFile()
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("No NuGet configuration file found")
        // Create default configuration
    } else {
        log.Fatalf("Error searching for config: %v", err)
    }
} else {
    fmt.Printf("Using configuration file: %s\n", configPath)
}
```

### FindAllConfigFiles

```go
func (f *ConfigFinder) FindAllConfigFiles() []string
```

Finds all available NuGet configuration files in the search paths.

**Returns:**
- `[]string`: List of paths to all found configuration files

**Example:**
```go
finder := finder.NewConfigFinder()
configFiles := finder.FindAllConfigFiles()

fmt.Printf("Found %d configuration files:\n", len(configFiles))
for i, path := range configFiles {
    fmt.Printf("%d. %s\n", i+1, path)
}

// Use the first one or merge multiple configurations
if len(configFiles) > 0 {
    primaryConfig := configFiles[0]
    fmt.Printf("Using primary config: %s\n", primaryConfig)
}
```

### FindProjectConfig

```go
func (f *ConfigFinder) FindProjectConfig(startDir string) (string, error)
```

Finds a project-level configuration file starting from the specified directory.

**Parameters:**
- `startDir` (string): Starting directory for the search

**Returns:**
- `string`: Path to the found project configuration file
- `error`: Error if no project configuration file is found

**Search Strategy:**
1. Starts from `startDir`
2. Looks for `NuGet.Config` in current directory
3. Walks up parent directories until found or reaches root
4. Stops at the first configuration file found

**Example:**
```go
finder := finder.NewConfigFinder()

// Find project config starting from current directory
projectConfig, err := finder.FindProjectConfig(".")
if err != nil {
    fmt.Println("No project-specific configuration found")
} else {
    fmt.Printf("Project config: %s\n", projectConfig)
}

// Find config for a specific project
projectPath := "/path/to/my/project"
projectConfig, err = finder.FindProjectConfig(projectPath)
if err != nil {
    fmt.Printf("No config found for project at %s\n", projectPath)
} else {
    fmt.Printf("Project config: %s\n", projectConfig)
}
```

### FindGlobalConfig

```go
func (f *ConfigFinder) FindGlobalConfig() (string, error)
```

Finds the global (user-level) configuration file.

**Returns:**
- `string`: Path to the global configuration file
- `error`: Error if no global configuration file is found

**Example:**
```go
finder := finder.NewConfigFinder()
globalConfig, err := finder.FindGlobalConfig()
if err != nil {
    fmt.Println("No global configuration found")
} else {
    fmt.Printf("Global config: %s\n", globalConfig)
}
```

### FindSystemConfig

```go
func (f *ConfigFinder) FindSystemConfig() (string, error)
```

Finds the system-wide configuration file.

**Returns:**
- `string`: Path to the system configuration file
- `error`: Error if no system configuration file is found

**Example:**
```go
finder := finder.NewConfigFinder()
systemConfig, err := finder.FindSystemConfig()
if err != nil {
    fmt.Println("No system configuration found")
} else {
    fmt.Printf("System config: %s\n", systemConfig)
}
```

## Path Management

### GetDefaultSearchPaths

```go
func (f *ConfigFinder) GetDefaultSearchPaths() []string
```

Returns the default search paths for the current platform.

**Returns:**
- `[]string`: List of default search paths

**Example:**
```go
finder := finder.NewConfigFinder()
paths := finder.GetDefaultSearchPaths()

fmt.Println("Default search paths:")
for i, path := range paths {
    fmt.Printf("%d. %s\n", i+1, path)
}
```

### AddSearchPath

```go
func (f *ConfigFinder) AddSearchPath(path string)
```

Adds a custom search path to the finder.

**Parameters:**
- `path` (string): Path to add to the search list

**Example:**
```go
finder := finder.NewConfigFinder()
finder.AddSearchPath("/custom/config/location/NuGet.Config")
finder.AddSearchPath("/another/location/NuGet.Config")

// Now search will include custom paths
configPath, err := finder.FindConfigFile()
```

### SetSearchPaths

```go
func (f *ConfigFinder) SetSearchPaths(paths []string)
```

Sets the complete list of search paths, replacing existing ones.

**Parameters:**
- `paths` ([]string): New list of search paths

**Example:**
```go
finder := finder.NewConfigFinder()

customPaths := []string{
    "./project.config",
    "/etc/nuget/global.config",
    "/usr/local/share/nuget/system.config",
}

finder.SetSearchPaths(customPaths)
```

## Utility Methods

### ConfigExists

```go
func (f *ConfigFinder) ConfigExists(path string) bool
```

Checks if a configuration file exists at the specified path.

**Parameters:**
- `path` (string): Path to check

**Returns:**
- `bool`: True if the configuration file exists and is readable

**Example:**
```go
finder := finder.NewConfigFinder()

configPath := "/path/to/NuGet.Config"
if finder.ConfigExists(configPath) {
    fmt.Printf("Configuration exists: %s\n", configPath)
} else {
    fmt.Printf("Configuration not found: %s\n", configPath)
}
```

### ValidateConfigFile

```go
func (f *ConfigFinder) ValidateConfigFile(path string) error
```

Validates that a file is a valid NuGet configuration file.

**Parameters:**
- `path` (string): Path to the configuration file

**Returns:**
- `error`: Error if the file is not valid

**Example:**
```go
finder := finder.NewConfigFinder()

configPath := "/path/to/NuGet.Config"
err := finder.ValidateConfigFile(configPath)
if err != nil {
    fmt.Printf("Invalid configuration file: %v\n", err)
} else {
    fmt.Println("Configuration file is valid")
}
```

## Advanced Usage

### Custom Search Strategy

```go
// Create a finder with custom search strategy
func createCustomFinder() *finder.ConfigFinder {
    finder := finder.NewConfigFinder()
    
    // Add project-specific paths
    finder.AddSearchPath("./config/NuGet.Config")
    finder.AddSearchPath("./settings/NuGet.Config")
    
    // Add environment-specific paths
    if env := os.Getenv("NUGET_CONFIG_PATH"); env != "" {
        finder.AddSearchPath(env)
    }
    
    return finder
}
```

### Hierarchical Configuration Discovery

```go
// Find all configurations in hierarchy
func findConfigHierarchy(projectPath string) ([]string, error) {
    finder := finder.NewConfigFinder()
    var configs []string
    
    // Project-level config
    if projectConfig, err := finder.FindProjectConfig(projectPath); err == nil {
        configs = append(configs, projectConfig)
    }
    
    // Global config
    if globalConfig, err := finder.FindGlobalConfig(); err == nil {
        configs = append(configs, globalConfig)
    }
    
    // System config
    if systemConfig, err := finder.FindSystemConfig(); err == nil {
        configs = append(configs, systemConfig)
    }
    
    return configs, nil
}
```

## Error Handling

The finder uses standard error types from the `pkg/errors` package:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

finder := finder.NewConfigFinder()
configPath, err := finder.FindConfigFile()
if err != nil {
    if errors.IsNotFoundError(err) {
        // No configuration file found
        fmt.Println("Creating default configuration...")
        // Handle missing configuration
    } else {
        // Other error occurred
        log.Fatalf("Search error: %v", err)
    }
}
```

## Platform Differences

### Windows
- User config: `%APPDATA%\NuGet\NuGet.Config`
- System config: `%ProgramData%\NuGet\NuGet.Config`

### macOS
- User config: `~/Library/Application Support/NuGet/NuGet.Config`
- System config: `/Library/Application Support/NuGet/NuGet.Config`

### Linux/Unix
- User config: `~/.config/NuGet/NuGet.Config` (respects `XDG_CONFIG_HOME`)
- System config: `/etc/NuGet/NuGet.Config`

## Best Practices

1. **Use default finder**: Start with `NewConfigFinder()` for standard behavior
2. **Handle missing files**: Always check for `IsNotFoundError` and provide fallbacks
3. **Respect hierarchy**: Use project configs over global configs when both exist
4. **Validate paths**: Use `ConfigExists()` before attempting to parse files
5. **Custom paths**: Use `AddSearchPath()` for additional locations rather than replacing defaults
6. **Environment variables**: Consider environment-specific configuration paths

## Thread Safety

The ConfigFinder is thread-safe for read operations but not for modifications. If you need to modify search paths concurrently, provide appropriate synchronization.
