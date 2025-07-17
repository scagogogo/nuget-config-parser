# Position-Aware Editing

Position-aware editing is an advanced feature that allows you to make precise modifications to NuGet configuration files while preserving the original formatting, comments, and structure. This is particularly useful when you want to minimize diffs in version control or maintain specific formatting requirements.

## Overview

Traditional configuration editing involves:
1. Parsing the entire configuration
2. Modifying the data structure
3. Serializing back to XML

This approach often results in:
- Loss of original formatting
- Removal of comments
- Large diffs in version control
- Inconsistent indentation

Position-aware editing solves these problems by:
- Tracking the exact location of each element
- Making surgical changes to specific parts
- Preserving original formatting and comments
- Generating minimal diffs

## How It Works

Position-aware editing works in three phases:

### 1. Position-Aware Parsing

First, parse the configuration file with position tracking enabled:

```go
api := nuget.NewAPI()

// Parse with position tracking
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatal(err)
}

// parseResult contains:
// - Config: The parsed configuration object
// - Positions: Map of element paths to their positions
// - Content: Original file content
```

### 2. Create Editor and Make Changes

Create an editor and specify the changes you want to make:

```go
// Create editor from parse result
editor := api.CreateConfigEditor(parseResult)

// Make changes (these are queued, not applied immediately)
err = editor.AddPackageSource("new-source", "https://example.com", "3")
err = editor.UpdatePackageSourceURL("existing-source", "https://new-url.com")
err = editor.RemovePackageSource("old-source")
```

### 3. Apply Changes

Apply all queued changes to generate the modified content:

```go
// Apply all changes and get modified content
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatal(err)
}

// Save the modified content
err = os.WriteFile("/path/to/NuGet.Config", modifiedContent, 0644)
```

## Basic Example

Here's a complete example of position-aware editing:

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
    
    configPath := "NuGet.Config"
    
    // Parse with position tracking
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse with positions: %v", err)
    }
    
    fmt.Printf("Original file size: %d bytes\n", len(parseResult.Content))
    
    // Create editor
    editor := api.CreateConfigEditor(parseResult)
    
    // Make changes
    fmt.Println("Making changes...")
    
    err = editor.AddPackageSource("company-feed", "https://nuget.company.com/v3/index.json", "3")
    if err != nil {
        log.Printf("Failed to add source: %v", err)
    }
    
    err = editor.UpdatePackageSourceURL("nuget.org", "https://api.nuget.org/v3/index.json")
    if err != nil {
        log.Printf("Failed to update URL: %v", err)
    }
    
    // Apply changes
    fmt.Println("Applying changes...")
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        log.Fatalf("Failed to apply edits: %v", err)
    }
    
    fmt.Printf("Modified file size: %d bytes\n", len(modifiedContent))
    
    // Save modified content
    err = os.WriteFile(configPath, modifiedContent, 0644)
    if err != nil {
        log.Fatalf("Failed to save file: %v", err)
    }
    
    fmt.Println("Changes applied successfully!")
}
```

## Advanced Usage

### Batch Operations

You can queue multiple operations before applying them:

```go
editor := api.CreateConfigEditor(parseResult)

// Queue multiple changes
editor.AddPackageSource("feed1", "https://feed1.com", "3")
editor.AddPackageSource("feed2", "https://feed2.com", "3")
editor.UpdatePackageSourceURL("existing", "https://new-url.com")
editor.RemovePackageSource("old-feed")

// Apply all changes at once
modifiedContent, err := editor.ApplyEdits()
```

### Inspecting Positions

You can examine the position information for debugging or analysis:

```go
parseResult, err := api.ParseFromFileWithPositions(configPath)
if err != nil {
    log.Fatal(err)
}

// Examine positions
positions := parseResult.Positions
for path, pos := range positions {
    fmt.Printf("Element: %s\n", path)
    fmt.Printf("  Tag: %s\n", pos.TagName)
    fmt.Printf("  Line: %d-%d\n", pos.Range.Start.Line, pos.Range.End.Line)
    fmt.Printf("  Attributes: %v\n", pos.Attributes)
    fmt.Println()
}
```

### Error Handling

Handle errors appropriately during editing:

```go
editor := api.CreateConfigEditor(parseResult)

err := editor.AddPackageSource("duplicate", "https://example.com", "3")
if err != nil {
    if strings.Contains(err.Error(), "already exists") {
        // Handle duplicate source
        fmt.Println("Source already exists, updating instead...")
        err = editor.UpdatePackageSourceURL("duplicate", "https://example.com")
    } else {
        log.Fatalf("Unexpected error: %v", err)
    }
}
```

## Comparison with Traditional Editing

### Traditional Approach

```go
// Traditional editing
config, err := api.ParseFromFile(configPath)
api.AddPackageSource(config, "new-source", "https://example.com", "3")
err = api.SaveConfig(config, configPath)

// Results in:
// - Complete file rewrite
// - Loss of formatting
// - Large diffs
// - No comments preserved
```

### Position-Aware Approach

```go
// Position-aware editing
parseResult, err := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)
editor.AddPackageSource("new-source", "https://example.com", "3")
modifiedContent, err := editor.ApplyEdits()
os.WriteFile(configPath, modifiedContent, 0644)

// Results in:
// - Minimal changes
// - Preserved formatting
// - Small diffs
// - Comments maintained
```

## Use Cases

### Version Control Friendly

Position-aware editing is ideal when working with version control:

```go
// Before committing changes, use position-aware editing
// to minimize diffs and maintain readability
parseResult, err := api.ParseFromFileWithPositions("NuGet.Config")
editor := api.CreateConfigEditor(parseResult)

// Make necessary changes
editor.AddPackageSource("ci-feed", "https://ci.company.com/nuget", "3")

// Apply with minimal diff
modifiedContent, err := editor.ApplyEdits()
os.WriteFile("NuGet.Config", modifiedContent, 0644)

// Commit will show only the actual changes made
```

### Automated Configuration Management

For tools that automatically manage configurations:

```go
func updateConfigurationAutomatically(configPath string, sources []SourceConfig) error {
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        return err
    }
    
    editor := api.CreateConfigEditor(parseResult)
    
    // Apply automated changes
    for _, source := range sources {
        if source.Action == "add" {
            editor.AddPackageSource(source.Key, source.URL, source.Version)
        } else if source.Action == "update" {
            editor.UpdatePackageSourceURL(source.Key, source.URL)
        } else if source.Action == "remove" {
            editor.RemovePackageSource(source.Key)
        }
    }
    
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        return err
    }
    
    return os.WriteFile(configPath, modifiedContent, 0644)
}
```

### Configuration Templates

Maintain configuration templates with position-aware editing:

```go
func applyTemplate(templatePath, targetPath string, customizations map[string]string) error {
    // Parse template
    parseResult, err := api.ParseFromFileWithPositions(templatePath)
    if err != nil {
        return err
    }
    
    editor := api.CreateConfigEditor(parseResult)
    
    // Apply customizations
    for key, value := range customizations {
        if strings.HasPrefix(key, "source.") {
            sourceName := strings.TrimPrefix(key, "source.")
            editor.UpdatePackageSourceURL(sourceName, value)
        }
    }
    
    // Generate customized configuration
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        return err
    }
    
    return os.WriteFile(targetPath, modifiedContent, 0644)
}
```

## Limitations

### Current Limitations

1. **Adding new sections**: Limited support for adding entirely new XML sections
2. **Complex restructuring**: Not suitable for major structural changes
3. **Memory usage**: Keeps entire file content in memory during editing
4. **Attribute manipulation**: Limited support for adding new attributes to existing elements

### When to Use Traditional Editing

Use traditional editing when:
- Creating configurations from scratch
- Making major structural changes
- Formatting consistency is more important than preserving original format
- Working with very large files where memory usage is a concern

### When to Use Position-Aware Editing

Use position-aware editing when:
- Making incremental changes to existing configurations
- Working with version control
- Preserving comments and formatting is important
- Minimizing diffs is crucial
- Automating configuration updates

## Best Practices

### 1. Parse Once, Edit Multiple Times

```go
// Good: Parse once for multiple edits
parseResult, err := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)

editor.AddPackageSource("source1", "url1", "3")
editor.AddPackageSource("source2", "url2", "3")
editor.UpdatePackageSourceURL("existing", "new-url")

modifiedContent, err := editor.ApplyEdits()
```

### 2. Handle Errors Gracefully

```go
// Check each operation for errors
err := editor.AddPackageSource("new-source", "https://example.com", "3")
if err != nil {
    log.Printf("Warning: Could not add source: %v", err)
    // Continue with other operations
}
```

### 3. Validate Results

```go
// Validate the result after applying edits
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    return err
}

// Re-parse to validate
_, err = api.ParseFromString(string(modifiedContent))
if err != nil {
    return fmt.Errorf("generated invalid XML: %w", err)
}
```

### 4. Backup Before Editing

```go
// Create backup before making changes
backupPath := configPath + ".backup"
originalContent, _ := os.ReadFile(configPath)
os.WriteFile(backupPath, originalContent, 0644)

// Proceed with position-aware editing
// ...

// Remove backup if successful
os.Remove(backupPath)
```

## Next Steps

- Explore the [Editor API](/api/editor) for detailed method documentation
- Check out [Position-Aware Editing Examples](/examples/position-aware-editing) for practical scenarios
- Learn about [Configuration Structure](./configuration.md) to understand what can be edited

Position-aware editing is a powerful feature that enables precise, minimal changes to configuration files while preserving their original structure and formatting.
