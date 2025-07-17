# Position-Aware Editing Examples

This example demonstrates advanced position-aware editing techniques for making precise modifications to NuGet configuration files while preserving formatting and minimizing diffs.

## Overview

Position-aware editing provides:
- Surgical modifications to specific file locations
- Preservation of original formatting and comments
- Minimal diffs for version control
- Precise control over changes

## Example 1: Basic Position-Aware Editing

Simple position-aware modifications:

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
    
    fmt.Println("=== Basic Position-Aware Editing ===")
    
    // Parse with position tracking
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse with positions: %v", err)
    }
    
    fmt.Printf("Original file size: %d bytes\n", len(parseResult.Content))
    
    // Create editor
    editor := api.CreateConfigEditor(parseResult)
    
    // Show current configuration
    config := editor.GetConfig()
    fmt.Printf("Current package sources: %d\n", len(config.PackageSources.Add))
    
    // Make precise changes
    fmt.Println("\nMaking position-aware changes...")
    
    // Add a new package source
    err = editor.AddPackageSource("company-feed", "https://nuget.company.com/v3/index.json", "3")
    if err != nil {
        log.Printf("Failed to add source: %v", err)
    } else {
        fmt.Println("✅ Added company-feed")
    }
    
    // Update an existing source URL
    err = editor.UpdatePackageSourceURL("nuget.org", "https://api.nuget.org/v3/index.json")
    if err != nil {
        log.Printf("Failed to update URL: %v", err)
    } else {
        fmt.Println("✅ Updated nuget.org URL")
    }
    
    // Apply all changes
    fmt.Println("\nApplying changes...")
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        log.Fatalf("Failed to apply edits: %v", err)
    }
    
    fmt.Printf("Modified file size: %d bytes\n", len(modifiedContent))
    fmt.Printf("Size change: %+d bytes\n", len(modifiedContent)-len(parseResult.Content))
    
    // Save the modified content
    err = os.WriteFile(configPath, modifiedContent, 0644)
    if err != nil {
        log.Fatalf("Failed to save file: %v", err)
    }
    
    fmt.Println("✅ Changes applied successfully with minimal diff!")
}
```

## Example 2: Batch Position-Aware Operations

Performing multiple operations efficiently:

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
    
    fmt.Println("=== Batch Position-Aware Operations ===")
    
    // Parse with positions
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    // Create editor
    editor := api.CreateConfigEditor(parseResult)
    
    // Define batch operations
    operations := []struct {
        name string
        op   func() error
    }{
        {"Add production feed", func() error {
            return editor.AddPackageSource("prod-feed", "https://prod.company.com/nuget", "3")
        }},
        {"Add staging feed", func() error {
            return editor.AddPackageSource("staging-feed", "https://staging.company.com/nuget", "3")
        }},
        {"Add development feed", func() error {
            return editor.AddPackageSource("dev-feed", "https://dev.company.com/nuget", "3")
        }},
        {"Update nuget.org protocol", func() error {
            return editor.UpdatePackageSourceVersion("nuget.org", "3")
        }},
        {"Remove old feed", func() error {
            return editor.RemovePackageSource("old-feed")
        }},
    }
    
    // Execute batch operations
    fmt.Println("Executing batch operations:")
    successCount := 0
    
    for _, op := range operations {
        err := op.op()
        if err != nil {
            fmt.Printf("  ❌ %s: %v\n", op.name, err)
        } else {
            fmt.Printf("  ✅ %s\n", op.name)
            successCount++
        }
    }
    
    fmt.Printf("\nCompleted %d/%d operations successfully\n", successCount, len(operations))
    
    // Apply all changes at once
    if successCount > 0 {
        fmt.Println("Applying all changes...")
        modifiedContent, err := editor.ApplyEdits()
        if err != nil {
            log.Fatalf("Failed to apply edits: %v", err)
        }
        
        // Save changes
        err = os.WriteFile(configPath, modifiedContent, 0644)
        if err != nil {
            log.Fatalf("Failed to save: %v", err)
        }
        
        fmt.Println("✅ All changes applied successfully!")
        
        // Show diff statistics
        showDiffStats(parseResult.Content, modifiedContent)
    }
}

func showDiffStats(original, modified []byte) {
    fmt.Println("\n=== Diff Statistics ===")
    fmt.Printf("Original size: %d bytes\n", len(original))
    fmt.Printf("Modified size: %d bytes\n", len(modified))
    fmt.Printf("Size change: %+d bytes\n", len(modified)-len(original))
    
    // Simple line count comparison
    originalLines := countLines(original)
    modifiedLines := countLines(modified)
    fmt.Printf("Original lines: %d\n", originalLines)
    fmt.Printf("Modified lines: %d\n", modifiedLines)
    fmt.Printf("Line change: %+d lines\n", modifiedLines-originalLines)
}

func countLines(content []byte) int {
    lines := 1
    for _, b := range content {
        if b == '\n' {
            lines++
        }
    }
    return lines
}
```

## Example 3: Position Information Analysis

Analyzing and using position information:

```go
package main

import (
    "fmt"
    "log"
    "sort"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    configPath := "NuGet.Config"
    
    fmt.Println("=== Position Information Analysis ===")
    
    // Parse with position tracking
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    // Create editor to access positions
    editor := api.CreateConfigEditor(parseResult)
    positions := editor.GetPositions()
    
    fmt.Printf("Found %d positioned elements\n", len(positions))
    
    // Analyze position information
    analyzePositions(positions)
    
    // Find elements by line number
    findElementsByLine(positions, 10) // Find elements around line 10
    
    // Show element hierarchy
    showElementHierarchy(positions)
    
    // Demonstrate targeted editing based on positions
    performTargetedEditing(editor, positions)
}

func analyzePositions(positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== Position Analysis ===")
    
    // Group by tag name
    tagCounts := make(map[string]int)
    lineRanges := make(map[string][]int)
    
    for path, pos := range positions {
        tagCounts[pos.TagName]++
        lineRanges[pos.TagName] = append(lineRanges[pos.TagName], pos.Range.Start.Line)
    }
    
    fmt.Println("Elements by tag type:")
    for tag, count := range tagCounts {
        fmt.Printf("  %s: %d elements\n", tag, count)
    }
    
    fmt.Println("\nLine distribution:")
    for tag, lines := range lineRanges {
        if len(lines) > 0 {
            sort.Ints(lines)
            fmt.Printf("  %s: lines %d-%d\n", tag, lines[0], lines[len(lines)-1])
        }
    }
}

func findElementsByLine(positions map[string]*parser.ElementPosition, targetLine int) {
    fmt.Printf("\n=== Elements near line %d ===\n", targetLine)
    
    type ElementInfo struct {
        Path string
        Pos  *parser.ElementPosition
        Dist int
    }
    
    var nearby []ElementInfo
    
    for path, pos := range positions {
        startDist := abs(pos.Range.Start.Line - targetLine)
        endDist := abs(pos.Range.End.Line - targetLine)
        minDist := min(startDist, endDist)
        
        if minDist <= 5 { // Within 5 lines
            nearby = append(nearby, ElementInfo{path, pos, minDist})
        }
    }
    
    // Sort by distance
    sort.Slice(nearby, func(i, j int) bool {
        return nearby[i].Dist < nearby[j].Dist
    })
    
    for _, elem := range nearby {
        fmt.Printf("  %s (%s) at lines %d-%d (distance: %d)\n",
            elem.Path, elem.Pos.TagName,
            elem.Pos.Range.Start.Line, elem.Pos.Range.End.Line,
            elem.Dist)
    }
}

func showElementHierarchy(positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== Element Hierarchy ===")
    
    // Group by hierarchy level (count of '/' in path)
    levels := make(map[int][]string)
    
    for path := range positions {
        level := countChar(path, '/')
        levels[level] = append(levels[level], path)
    }
    
    // Sort levels
    var sortedLevels []int
    for level := range levels {
        sortedLevels = append(sortedLevels, level)
    }
    sort.Ints(sortedLevels)
    
    for _, level := range sortedLevels {
        fmt.Printf("Level %d:\n", level)
        sort.Strings(levels[level])
        for _, path := range levels[level] {
            indent := strings.Repeat("  ", level+1)
            pos := positions[path]
            fmt.Printf("%s%s (%s) [%d:%d-%d:%d]\n",
                indent, path, pos.TagName,
                pos.Range.Start.Line, pos.Range.Start.Column,
                pos.Range.End.Line, pos.Range.End.Column)
        }
    }
}

func performTargetedEditing(editor *editor.ConfigEditor, positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== Targeted Editing Based on Positions ===")
    
    // Find package source elements
    var packageSourcePaths []string
    for path, pos := range positions {
        if pos.TagName == "add" && strings.Contains(path, "packageSources") {
            packageSourcePaths = append(packageSourcePaths, path)
        }
    }
    
    fmt.Printf("Found %d package source elements\n", len(packageSourcePaths))
    
    // Sort by line number for predictable processing
    sort.Slice(packageSourcePaths, func(i, j int) bool {
        return positions[packageSourcePaths[i]].Range.Start.Line <
            positions[packageSourcePaths[j]].Range.Start.Line
    })
    
    // Show what we found
    for i, path := range packageSourcePaths {
        pos := positions[path]
        fmt.Printf("  %d. %s at line %d\n", i+1, path, pos.Range.Start.Line)
        
        // Show attributes if available
        if len(pos.Attributes) > 0 {
            for attrName, attrValue := range pos.Attributes {
                fmt.Printf("     %s=\"%s\"\n", attrName, attrValue)
            }
        }
    }
    
    // Perform targeted modification
    if len(packageSourcePaths) > 0 {
        fmt.Println("\nPerforming targeted modification...")
        
        // Add a new source (this will be positioned appropriately)
        err := editor.AddPackageSource("targeted-source", "https://targeted.example.com", "3")
        if err != nil {
            fmt.Printf("Failed to add targeted source: %v\n", err)
        } else {
            fmt.Println("✅ Added targeted source")
        }
    }
}

// Helper functions
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func countChar(s string, c rune) int {
    count := 0
    for _, r := range s {
        if r == c {
            count++
        }
    }
    return count
}
```

## Example 4: Diff-Minimizing Strategies

Techniques for minimizing version control diffs:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    configPath := "NuGet.Config"
    
    fmt.Println("=== Diff-Minimizing Strategies ===")
    
    // Create backup for comparison
    originalContent, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalf("Failed to read original: %v", err)
    }
    
    // Parse with positions
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("Failed to parse: %v", err)
    }
    
    // Create editor
    editor := api.CreateConfigEditor(parseResult)
    
    // Demonstrate different editing strategies
    strategies := []struct {
        name string
        edit func() error
    }{
        {"Add single source", func() error {
            return editor.AddPackageSource("minimal-source", "https://minimal.example.com", "3")
        }},
        {"Update existing URL", func() error {
            return editor.UpdatePackageSourceURL("nuget.org", "https://api.nuget.org/v3/index.json")
        }},
        {"Update protocol version", func() error {
            return editor.UpdatePackageSourceVersion("nuget.org", "3")
        }},
    }
    
    for i, strategy := range strategies {
        fmt.Printf("\n--- Strategy %d: %s ---\n", i+1, strategy.name)
        
        // Apply single change
        err := strategy.edit()
        if err != nil {
            fmt.Printf("❌ Failed: %v\n", err)
            continue
        }
        
        // Apply and analyze diff
        modifiedContent, err := editor.ApplyEdits()
        if err != nil {
            fmt.Printf("❌ Apply failed: %v\n", err)
            continue
        }
        
        // Analyze the diff
        analyzeDiff(originalContent, modifiedContent, strategy.name)
        
        // Reset for next strategy (in real usage, you'd work with fresh editor)
        // For demo purposes, we'll continue with accumulated changes
        originalContent = modifiedContent
    }
    
    // Save final result
    finalContent, _ := editor.ApplyEdits()
    err = os.WriteFile(configPath, finalContent, 0644)
    if err != nil {
        log.Fatalf("Failed to save: %v", err)
    }
    
    fmt.Println("\n✅ All strategies applied with minimal diffs!")
}

func analyzeDiff(original, modified []byte, strategyName string) {
    fmt.Printf("Diff analysis for %s:\n", strategyName)
    
    originalLines := strings.Split(string(original), "\n")
    modifiedLines := strings.Split(string(modified), "\n")
    
    // Simple diff analysis
    added := 0
    removed := 0
    changed := 0
    
    maxLines := max(len(originalLines), len(modifiedLines))
    
    for i := 0; i < maxLines; i++ {
        var origLine, modLine string
        
        if i < len(originalLines) {
            origLine = originalLines[i]
        }
        if i < len(modifiedLines) {
            modLine = modifiedLines[i]
        }
        
        if origLine == "" && modLine != "" {
            added++
        } else if origLine != "" && modLine == "" {
            removed++
        } else if origLine != modLine {
            changed++
        }
    }
    
    fmt.Printf("  Lines added: %d\n", added)
    fmt.Printf("  Lines removed: %d\n", removed)
    fmt.Printf("  Lines changed: %d\n", changed)
    fmt.Printf("  Total diff impact: %d lines\n", added+removed+changed)
    
    // Calculate diff efficiency
    totalLines := len(originalLines)
    diffPercentage := float64(added+removed+changed) / float64(totalLines) * 100
    fmt.Printf("  Diff percentage: %.2f%%\n", diffPercentage)
    
    if diffPercentage < 5.0 {
        fmt.Printf("  ✅ Excellent diff efficiency (< 5%%)\n")
    } else if diffPercentage < 15.0 {
        fmt.Printf("  ✅ Good diff efficiency (< 15%%)\n")
    } else {
        fmt.Printf("  ⚠️  High diff impact (>= 15%%)\n")
    }
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

## Key Benefits

### Minimal Diffs
- Only changes necessary parts of the file
- Preserves original formatting and structure
- Maintains comments and whitespace
- Reduces version control noise

### Precision Control
- Exact positioning of modifications
- Surgical changes to specific elements
- Predictable modification behavior
- Fine-grained control over output

### Format Preservation
- Maintains original indentation
- Preserves XML comments
- Keeps consistent spacing
- Respects existing structure

## Best Practices

1. **Parse once, edit multiple times**: Use single parse result for multiple edits
2. **Batch operations**: Group related changes together
3. **Validate results**: Always verify the modified XML is valid
4. **Handle errors gracefully**: Check each operation for errors
5. **Test round-trips**: Ensure parse → edit → parse cycles work correctly

## Common Use Cases

### Automated Configuration Updates
```go
// Update configurations with minimal impact
parseResult, _ := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)
editor.UpdatePackageSourceURL("source", "new-url")
modifiedContent, _ := editor.ApplyEdits()
```

### CI/CD Configuration Management
```go
// Apply environment-specific changes
editor.AddPackageSource("ci-feed", ciURL, "3")
editor.UpdatePackageSourceVersion("nuget.org", "3")
modifiedContent, _ := editor.ApplyEdits()
```

### Configuration Migration
```go
// Migrate configurations with preserved formatting
editor.UpdatePackageSourceURL("old-source", "new-url")
editor.UpdatePackageSourceVersion("old-source", "3")
modifiedContent, _ := editor.ApplyEdits()
```

## Next Steps

After mastering position-aware editing:

1. Explore the [Editor API](/api/editor) for detailed method documentation
2. Learn about [Configuration Structure](../guide/configuration.md)
3. Study [Serialization](./serialization.md) for understanding XML output

Position-aware editing is the most advanced feature of the library, providing surgical precision for configuration modifications while maintaining the integrity and readability of your configuration files.
