# Utils API

The `pkg/utils` package provides utility functions for file operations, path manipulation, and XML processing used throughout the NuGet Config Parser library.

## Overview

The Utils API provides:
- File system operations and checks
- Cross-platform path manipulation
- XML processing utilities
- String and data validation helpers

## File System Operations

### FileExists

```go
func FileExists(filePath string) bool
```

Checks if a file exists and is not a directory.

**Parameters:**
- `filePath` (string): Path to the file to check

**Returns:**
- `bool`: True if the file exists and is not a directory, false otherwise

**Example:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/utils"

configPath := "/path/to/NuGet.Config"
if utils.FileExists(configPath) {
    fmt.Printf("Configuration file exists: %s\n", configPath)
    // Proceed with parsing...
} else {
    fmt.Printf("Configuration file not found: %s\n", configPath)
    // Create default configuration...
}
```

### DirExists

```go
func DirExists(dirPath string) bool
```

Checks if a directory exists.

**Parameters:**
- `dirPath` (string): Path to the directory to check

**Returns:**
- `bool`: True if the directory exists, false otherwise

**Example:**
```go
configDir := "/etc/nuget"
if utils.DirExists(configDir) {
    fmt.Printf("Configuration directory exists: %s\n", configDir)
} else {
    fmt.Printf("Creating configuration directory: %s\n", configDir)
    err := os.MkdirAll(configDir, 0755)
    if err != nil {
        log.Fatalf("Failed to create directory: %v", err)
    }
}
```

### IsReadableFile

```go
func IsReadableFile(filePath string) bool
```

Checks if a file exists and is readable.

**Parameters:**
- `filePath` (string): Path to the file to check

**Returns:**
- `bool`: True if the file exists and is readable, false otherwise

**Example:**
```go
configPath := "/path/to/NuGet.Config"
if utils.IsReadableFile(configPath) {
    // Safe to read the file
    content, err := os.ReadFile(configPath)
    if err != nil {
        log.Printf("Error reading file: %v", err)
    }
} else {
    fmt.Printf("File is not readable: %s\n", configPath)
}
```

### IsWritableFile

```go
func IsWritableFile(filePath string) bool
```

Checks if a file is writable (or if it doesn't exist, if the directory is writable).

**Parameters:**
- `filePath` (string): Path to the file to check

**Returns:**
- `bool`: True if the file is writable, false otherwise

**Example:**
```go
configPath := "/path/to/NuGet.Config"
if utils.IsWritableFile(configPath) {
    // Safe to write to the file
    err := os.WriteFile(configPath, []byte("content"), 0644)
    if err != nil {
        log.Printf("Error writing file: %v", err)
    }
} else {
    fmt.Printf("File is not writable: %s\n", configPath)
}
```

## Path Manipulation

### IsAbsolutePath

```go
func IsAbsolutePath(path string) bool
```

Checks if a path is absolute.

**Parameters:**
- `path` (string): Path to check

**Returns:**
- `bool`: True if the path is absolute, false if relative

**Example:**
```go
// Unix/Linux examples
fmt.Println(utils.IsAbsolutePath("/etc/nuget"))        // true
fmt.Println(utils.IsAbsolutePath("./config"))          // false
fmt.Println(utils.IsAbsolutePath("../config"))         // false

// Windows examples
fmt.Println(utils.IsAbsolutePath("C:\\nuget\\config")) // true
fmt.Println(utils.IsAbsolutePath(".\\config"))         // false
```

### NormalizePath

```go
func NormalizePath(path string) string
```

Normalizes a file path by cleaning it and converting to the OS-specific format.

**Parameters:**
- `path` (string): Path to normalize

**Returns:**
- `string`: Normalized path

**Example:**
```go
// Clean up redundant path elements
messyPath := "/etc/nuget/../nuget/./config"
cleanPath := utils.NormalizePath(messyPath)
fmt.Printf("Normalized: %s\n", cleanPath)
// Output: /etc/nuget/config

// Handle different separators
mixedPath := "/etc\\nuget/config"
normalizedPath := utils.NormalizePath(mixedPath)
fmt.Printf("Normalized: %s\n", normalizedPath)
```

### JoinPaths

```go
func JoinPaths(basePath string, paths ...string) string
```

Joins multiple path elements using the OS-specific path separator.

**Parameters:**
- `basePath` (string): Base path
- `paths` (...string): Additional path elements to join

**Returns:**
- `string`: Joined path

**Example:**
```go
// Join path elements
basePath := "/etc"
configPath := utils.JoinPaths(basePath, "nuget", "NuGet.Config")
fmt.Printf("Joined path: %s\n", configPath)
// Output: /etc/nuget/NuGet.Config

// Windows example
winBase := "C:\\Users"
winPath := utils.JoinPaths(winBase, "username", ".nuget", "packages")
fmt.Printf("Windows path: %s\n", winPath)
// Output: C:\Users\username\.nuget\packages

// Handle trailing slashes
trailingSlash := "/home/user/"
result := utils.JoinPaths(trailingSlash, "nuget", "packages")
fmt.Printf("Result: %s\n", result)
// Output: /home/user/nuget/packages
```

### ResolvePath

```go
func ResolvePath(basePath, path string) string
```

Resolves a path relative to a base path. If the path is already absolute, returns it normalized.

**Parameters:**
- `basePath` (string): Base path for resolving relative paths
- `path` (string): Path to resolve (can be relative or absolute)

**Returns:**
- `string`: Resolved absolute path

**Example:**
```go
basePath := "/etc/nuget"

// Resolve relative path
relativePath := "../packages/cache"
resolvedRelative := utils.ResolvePath(basePath, relativePath)
fmt.Printf("Relative '%s' resolved to: %s\n", relativePath, resolvedRelative)
// Output: Relative '../packages/cache' resolved to: /etc/packages/cache

// Handle absolute path
absolutePath := "/var/nuget/packages"
resolvedAbsolute := utils.ResolvePath(basePath, absolutePath)
fmt.Printf("Absolute '%s' remains: %s\n", absolutePath, resolvedAbsolute)
// Output: Absolute '/var/nuget/packages' remains: /var/nuget/packages
```

### ExpandEnvVars

```go
func ExpandEnvVars(path string) string
```

Expands environment variables in a path string.

**Parameters:**
- `path` (string): Path containing environment variables

**Returns:**
- `string`: Path with environment variables expanded

**Supported Formats:**
- Unix/Linux/macOS: `$VAR` or `${VAR}`
- Windows: `%VAR%`

**Example:**
```go
// Unix/Linux/macOS examples
unixPath := "$HOME/.nuget/packages"
expandedUnixPath := utils.ExpandEnvVars(unixPath)
fmt.Printf("Expanded: %s\n", expandedUnixPath)
// Output: Expanded: /home/user/.nuget/packages

bracedPath := "${HOME}/.config/NuGet/NuGet.Config"
expandedBracedPath := utils.ExpandEnvVars(bracedPath)
fmt.Printf("Expanded: %s\n", expandedBracedPath)

// Windows examples
winPath := "%USERPROFILE%\\.nuget\\packages"
expandedWinPath := utils.ExpandEnvVars(winPath)
fmt.Printf("Expanded: %s\n", expandedWinPath)
// Output: Expanded: C:\Users\username\.nuget\packages
```

## XML Processing

### ValidateXML

```go
func ValidateXML(content []byte) error
```

Validates that content is well-formed XML.

**Parameters:**
- `content` ([]byte): XML content to validate

**Returns:**
- `error`: Error if XML is not well-formed, nil if valid

**Example:**
```go
// Valid XML
validXML := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`)

err := utils.ValidateXML(validXML)
if err != nil {
    fmt.Printf("Invalid XML: %v\n", err)
} else {
    fmt.Println("XML is valid")
}

// Invalid XML
invalidXML := []byte(`<configuration><packageSources><add key="test"`)
err = utils.ValidateXML(invalidXML)
if err != nil {
    fmt.Printf("Invalid XML: %v\n", err)
}
```

### FormatXML

```go
func FormatXML(content []byte) ([]byte, error)
```

Formats XML content with proper indentation.

**Parameters:**
- `content` ([]byte): XML content to format

**Returns:**
- `[]byte`: Formatted XML content
- `error`: Error if formatting fails

**Example:**
```go
// Unformatted XML
unformattedXML := []byte(`<configuration><packageSources><add key="nuget.org" value="https://api.nuget.org/v3/index.json" /></packageSources></configuration>`)

formattedXML, err := utils.FormatXML(unformattedXML)
if err != nil {
    log.Printf("Failed to format XML: %v", err)
} else {
    fmt.Println("Formatted XML:")
    fmt.Println(string(formattedXML))
}

// Output:
// <?xml version="1.0" encoding="UTF-8"?>
// <configuration>
//   <packageSources>
//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json"/>
//   </packageSources>
// </configuration>
```

### ExtractXMLElement

```go
func ExtractXMLElement(content []byte, elementName string) ([]byte, error)
```

Extracts a specific XML element from content.

**Parameters:**
- `content` ([]byte): XML content to search
- `elementName` (string): Name of the element to extract

**Returns:**
- `[]byte`: Extracted element content
- `error`: Error if element not found or extraction fails

**Example:**
```go
xmlContent := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
    <add key="local" value="/path/to/local" />
  </packageSources>
  <config>
    <add key="globalPackagesFolder" value="/packages" />
  </config>
</configuration>`)

// Extract packageSources element
packageSources, err := utils.ExtractXMLElement(xmlContent, "packageSources")
if err != nil {
    log.Printf("Failed to extract element: %v", err)
} else {
    fmt.Println("Package sources:")
    fmt.Println(string(packageSources))
}
```

## String Utilities

### IsEmptyOrWhitespace

```go
func IsEmptyOrWhitespace(s string) bool
```

Checks if a string is empty or contains only whitespace.

**Parameters:**
- `s` (string): String to check

**Returns:**
- `bool`: True if string is empty or whitespace-only

**Example:**
```go
fmt.Println(utils.IsEmptyOrWhitespace(""))           // true
fmt.Println(utils.IsEmptyOrWhitespace("   "))        // true
fmt.Println(utils.IsEmptyOrWhitespace("\t\n"))       // true
fmt.Println(utils.IsEmptyOrWhitespace("content"))    // false
fmt.Println(utils.IsEmptyOrWhitespace(" content "))  // false
```

### TrimWhitespace

```go
func TrimWhitespace(s string) string
```

Trims leading and trailing whitespace from a string.

**Parameters:**
- `s` (string): String to trim

**Returns:**
- `string`: Trimmed string

**Example:**
```go
input := "  \t  content with spaces  \n  "
trimmed := utils.TrimWhitespace(input)
fmt.Printf("Trimmed: '%s'\n", trimmed)
// Output: Trimmed: 'content with spaces'
```

### SanitizeXMLValue

```go
func SanitizeXMLValue(value string) string
```

Sanitizes a string value for safe use in XML attributes or content.

**Parameters:**
- `value` (string): Value to sanitize

**Returns:**
- `string`: Sanitized value safe for XML

**Example:**
```go
unsafeValue := `value with "quotes" & <brackets>`
safeValue := utils.SanitizeXMLValue(unsafeValue)
fmt.Printf("Sanitized: %s\n", safeValue)
// Output: Sanitized: value with &quot;quotes&quot; &amp; &lt;brackets&gt;
```

## Validation Utilities

### IsValidURL

```go
func IsValidURL(urlStr string) bool
```

Validates if a string is a valid URL.

**Parameters:**
- `urlStr` (string): URL string to validate

**Returns:**
- `bool`: True if URL is valid

**Example:**
```go
validURLs := []string{
    "https://api.nuget.org/v3/index.json",
    "http://localhost:8080/nuget",
    "file:///path/to/packages",
}

invalidURLs := []string{
    "not-a-url",
    "://missing-scheme",
    "https://",
}

for _, url := range validURLs {
    if utils.IsValidURL(url) {
        fmt.Printf("Valid URL: %s\n", url)
    }
}

for _, url := range invalidURLs {
    if !utils.IsValidURL(url) {
        fmt.Printf("Invalid URL: %s\n", url)
    }
}
```

### IsValidPackageSourceKey

```go
func IsValidPackageSourceKey(key string) bool
```

Validates if a string is a valid package source key.

**Parameters:**
- `key` (string): Package source key to validate

**Returns:**
- `bool`: True if key is valid

**Validation Rules:**
- Not empty or whitespace-only
- Contains only alphanumeric characters, hyphens, underscores, and dots
- Does not start or end with special characters

**Example:**
```go
validKeys := []string{
    "nuget.org",
    "company-feed",
    "local_packages",
    "feed123",
}

invalidKeys := []string{
    "",
    "   ",
    "invalid key with spaces",
    "-starts-with-dash",
    "ends-with-dash-",
    "has@invalid#chars",
}

for _, key := range validKeys {
    if utils.IsValidPackageSourceKey(key) {
        fmt.Printf("Valid key: %s\n", key)
    }
}

for _, key := range invalidKeys {
    if !utils.IsValidPackageSourceKey(key) {
        fmt.Printf("Invalid key: '%s'\n", key)
    }
}
```

## Complete Usage Example

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/utils"
)

func main() {
    // File operations
    configPath := "/path/to/NuGet.Config"
    
    if utils.FileExists(configPath) {
        fmt.Printf("Configuration file exists: %s\n", configPath)
        
        if utils.IsReadableFile(configPath) {
            content, err := os.ReadFile(configPath)
            if err != nil {
                log.Printf("Error reading file: %v", err)
                return
            }
            
            // Validate XML
            if err := utils.ValidateXML(content); err != nil {
                log.Printf("Invalid XML: %v", err)
                return
            }
            
            // Format XML
            formatted, err := utils.FormatXML(content)
            if err != nil {
                log.Printf("Failed to format XML: %v", err)
            } else {
                fmt.Println("Formatted XML:")
                fmt.Println(string(formatted))
            }
        }
    } else {
        fmt.Printf("Configuration file not found: %s\n", configPath)
        
        // Create directory if needed
        configDir := utils.JoinPaths("/etc", "nuget")
        if !utils.DirExists(configDir) {
            fmt.Printf("Creating directory: %s\n", configDir)
            err := os.MkdirAll(configDir, 0755)
            if err != nil {
                log.Fatalf("Failed to create directory: %v", err)
            }
        }
    }
    
    // Path manipulation
    basePath := "/etc/nuget"
    relativePath := "../packages"
    resolvedPath := utils.ResolvePath(basePath, relativePath)
    fmt.Printf("Resolved path: %s\n", resolvedPath)
    
    // Environment variable expansion
    envPath := "$HOME/.nuget/packages"
    expandedPath := utils.ExpandEnvVars(envPath)
    fmt.Printf("Expanded path: %s\n", expandedPath)
    
    // Validation
    testURL := "https://api.nuget.org/v3/index.json"
    if utils.IsValidURL(testURL) {
        fmt.Printf("Valid URL: %s\n", testURL)
    }
    
    testKey := "nuget.org"
    if utils.IsValidPackageSourceKey(testKey) {
        fmt.Printf("Valid package source key: %s\n", testKey)
    }
}
```

## Best Practices

1. **Check file existence**: Always use `FileExists()` before attempting file operations
2. **Validate inputs**: Use validation functions for URLs and keys before processing
3. **Handle paths correctly**: Use path manipulation functions for cross-platform compatibility
4. **Sanitize XML content**: Use `SanitizeXMLValue()` for user-provided content
5. **Validate XML**: Use `ValidateXML()` before parsing to catch malformed content early
6. **Normalize paths**: Use `NormalizePath()` to clean up path strings
7. **Expand environment variables**: Use `ExpandEnvVars()` for flexible path configuration

## Thread Safety

All utility functions in this package are thread-safe and can be used concurrently.
