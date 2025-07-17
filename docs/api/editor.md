# Editor API

The `pkg/editor` package provides position-aware editing capabilities for NuGet configuration files. This allows you to make precise modifications while preserving the original file formatting and minimizing diffs.

## Overview

Position-aware editing is particularly useful when:
- You want to maintain the original file formatting
- You need to minimize version control diffs
- You're working with configuration files that have specific formatting requirements
- You want to preserve comments and whitespace

## Types

### ConfigEditor

```go
type ConfigEditor struct {
    parseResult *parser.ParseResult
    edits       []Edit
}
```

The main editor type that tracks modifications to be applied to a configuration file.

### Edit

```go
type Edit struct {
    Range   parser.Range // The range to replace
    NewText string       // New text content
    Type    string       // Edit type: "add", "update", "delete"
}
```

Represents a single edit operation with position information.

## Constructor

### NewConfigEditor

```go
func NewConfigEditor(parseResult *parser.ParseResult) *ConfigEditor
```

Creates a new configuration editor from a parse result that includes position information.

**Parameters:**
- `parseResult` (*parser.ParseResult): Parse result with position tracking

**Returns:**
- `*ConfigEditor`: New editor instance

**Example:**
```go
// Parse with position tracking
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatal(err)
}

// Create editor
editor := editor.NewConfigEditor(parseResult)
```

## Configuration Access

### GetConfig

```go
func (e *ConfigEditor) GetConfig() *types.NuGetConfig
```

Returns the configuration object being edited.

**Returns:**
- `*types.NuGetConfig`: The configuration object

**Example:**
```go
config := editor.GetConfig()
fmt.Printf("Current sources: %d\n", len(config.PackageSources.Add))
```

### GetPositions

```go
func (e *ConfigEditor) GetPositions() map[string]*parser.ElementPosition
```

Returns the position information for all elements in the configuration.

**Returns:**
- `map[string]*parser.ElementPosition`: Map of element paths to position information

**Example:**
```go
positions := editor.GetPositions()
for path, pos := range positions {
    fmt.Printf("Element %s at line %d\n", path, pos.Range.Start.Line)
}
```

## Package Source Operations

### AddPackageSource

```go
func (e *ConfigEditor) AddPackageSource(key, value, protocolVersion string) error
```

Adds a new package source to the configuration.

**Parameters:**
- `key` (string): Unique identifier for the package source
- `value` (string): URL or path to the package source
- `protocolVersion` (string): Protocol version (can be empty)

**Returns:**
- `error`: Error if the operation fails

**Example:**
```go
err := editor.AddPackageSource(
    "company-feed", 
    "https://nuget.company.com/v3/index.json", 
    "3"
)
if err != nil {
    log.Fatalf("Failed to add package source: %v", err)
}
```

### RemovePackageSource

```go
func (e *ConfigEditor) RemovePackageSource(sourceKey string) error
```

Removes a package source from the configuration.

**Parameters:**
- `sourceKey` (string): Key of the package source to remove

**Returns:**
- `error`: Error if the source is not found or operation fails

**Example:**
```go
err := editor.RemovePackageSource("old-feed")
if err != nil {
    log.Printf("Failed to remove package source: %v", err)
}
```

### UpdatePackageSourceURL

```go
func (e *ConfigEditor) UpdatePackageSourceURL(sourceKey, newURL string) error
```

Updates the URL of an existing package source.

**Parameters:**
- `sourceKey` (string): Key of the package source to update
- `newURL` (string): New URL for the package source

**Returns:**
- `error`: Error if the source is not found or operation fails

**Example:**
```go
err := editor.UpdatePackageSourceURL(
    "nuget.org", 
    "https://api.nuget.org/v3/index.json"
)
if err != nil {
    log.Printf("Failed to update URL: %v", err)
}
```

### UpdatePackageSourceVersion

```go
func (e *ConfigEditor) UpdatePackageSourceVersion(sourceKey, newVersion string) error
```

Updates the protocol version of an existing package source.

**Parameters:**
- `sourceKey` (string): Key of the package source to update
- `newVersion` (string): New protocol version

**Returns:**
- `error`: Error if the source is not found or operation fails

**Example:**
```go
err := editor.UpdatePackageSourceVersion("my-feed", "3")
if err != nil {
    log.Printf("Failed to update version: %v", err)
}
```

## Applying Changes

### ApplyEdits

```go
func (e *ConfigEditor) ApplyEdits() ([]byte, error)
```

Applies all pending edits and returns the modified file content.

**Returns:**
- `[]byte`: Modified file content
- `error`: Error if applying edits fails

**Example:**
```go
// Make several changes
editor.AddPackageSource("feed1", "https://feed1.com", "3")
editor.UpdatePackageSourceURL("feed2", "https://newfeed2.com")
editor.RemovePackageSource("old-feed")

// Apply all changes
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatalf("Failed to apply edits: %v", err)
}

// Save to file
err = os.WriteFile("/path/to/NuGet.Config", modifiedContent, 0644)
if err != nil {
    log.Fatalf("Failed to save file: %v", err)
}
```

## Complete Example

Here's a complete example showing how to use the editor:

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
    
    // Parse with position tracking
    configPath := "/path/to/NuGet.Config"
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Create editor
    editor := api.CreateConfigEditor(parseResult)
    
    // Show current configuration
    config := editor.GetConfig()
    fmt.Printf("Current package sources: %d\n", len(config.PackageSources.Add))
    
    // Make changes
    fmt.Println("Adding new package source...")
    err = editor.AddPackageSource(
        "company-internal", 
        "https://nuget.company.com/v3/index.json", 
        "3"
    )
    if err != nil {
        log.Fatalf("Failed to add source: %v", err)
    }
    
    fmt.Println("Updating existing source...")
    err = editor.UpdatePackageSourceURL(
        "nuget.org", 
        "https://api.nuget.org/v3/index.json"
    )
    if err != nil {
        log.Printf("Warning: %v", err)
    }
    
    // Apply changes
    fmt.Println("Applying changes...")
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        log.Fatalf("Failed to apply edits: %v", err)
    }
    
    // Save to file
    err = os.WriteFile(configPath, modifiedContent, 0644)
    if err != nil {
        log.Fatalf("Failed to save file: %v", err)
    }
    
    fmt.Println("Configuration updated successfully!")
    
    // Verify changes
    updatedConfig := editor.GetConfig()
    fmt.Printf("Updated package sources: %d\n", len(updatedConfig.PackageSources.Add))
}
```

## Advanced Usage

### Batch Operations

You can perform multiple operations before applying changes:

```go
// Multiple changes in one batch
editor.AddPackageSource("feed1", "https://feed1.com", "3")
editor.AddPackageSource("feed2", "https://feed2.com", "3")
editor.UpdatePackageSourceURL("existing", "https://new-url.com")
editor.RemovePackageSource("old-feed")

// Apply all at once
modifiedContent, err := editor.ApplyEdits()
```

### Error Handling

Handle errors appropriately for each operation:

```go
err := editor.AddPackageSource("duplicate", "https://example.com", "3")
if err != nil {
    if strings.Contains(err.Error(), "already exists") {
        // Handle duplicate source
        log.Printf("Source already exists, updating instead")
        err = editor.UpdatePackageSourceURL("duplicate", "https://example.com")
    } else {
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

### Position Information

Access detailed position information:

```go
positions := editor.GetPositions()
for path, elemPos := range positions {
    fmt.Printf("Element: %s\n", path)
    fmt.Printf("  Tag: %s\n", elemPos.TagName)
    fmt.Printf("  Line: %d-%d\n", elemPos.Range.Start.Line, elemPos.Range.End.Line)
    fmt.Printf("  Attributes: %v\n", elemPos.Attributes)
}
```

## Benefits of Position-Aware Editing

1. **Minimal Diffs**: Only the necessary parts of the file are changed
2. **Format Preservation**: Original indentation and formatting are maintained
3. **Comment Preservation**: Comments in the original file are preserved
4. **Precise Control**: Exact control over what gets modified
5. **Version Control Friendly**: Smaller, cleaner diffs in version control

## Limitations

1. **Adding New Attributes**: Currently limited support for adding new attributes to existing elements
2. **Complex Restructuring**: Not suitable for major structural changes to the XML
3. **Memory Usage**: Keeps the entire file content in memory during editing

## Best Practices

1. **Parse Once**: Use the same parse result for multiple edit operations
2. **Batch Changes**: Group related changes together before applying
3. **Error Handling**: Always check for errors after each operation
4. **Backup**: Consider backing up the original file before applying changes
5. **Validation**: Validate the configuration after applying changes

```go
// Good practice: batch operations
editor.AddPackageSource("feed1", "url1", "3")
editor.AddPackageSource("feed2", "url2", "3")
editor.RemovePackageSource("old")
modifiedContent, err := editor.ApplyEdits()

// Good practice: validate after changes
if err == nil {
    // Re-parse to validate
    _, err = api.ParseFromString(string(modifiedContent))
    if err != nil {
        log.Printf("Warning: Generated invalid XML: %v", err)
    }
}
```
