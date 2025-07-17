# Parser API

The `pkg/parser` package provides configuration file parsing functionality with support for position tracking and detailed error reporting.

## Overview

The Parser API is responsible for:
- Parsing NuGet configuration files from various sources
- Tracking element positions for position-aware editing
- Providing detailed error information for debugging
- Validating configuration file structure

## Types

### ConfigParser

```go
type ConfigParser struct {
    DefaultConfigSearchPaths []string
    TrackPositions          bool
}
```

The main parser type that handles configuration file parsing.

**Fields:**
- `DefaultConfigSearchPaths`: List of default paths to search for configuration files
- `TrackPositions`: Whether to track element positions during parsing

### ParseResult

```go
type ParseResult struct {
    Config    *types.NuGetConfig          // Parsed configuration
    Positions map[string]*ElementPosition // Element position information
    Content   []byte                      // Original file content
}
```

Contains the result of position-aware parsing.

**Fields:**
- `Config`: The parsed NuGet configuration object
- `Positions`: Map of element paths to their position information
- `Content`: Original file content as bytes

### ElementPosition

```go
type ElementPosition struct {
    TagName    string            // XML tag name
    Attributes map[string]string // Element attributes
    Range      Range             // Element range in the file
    AttrRanges map[string]Range  // Attribute value ranges
    Content    string            // Element content
    SelfClose  bool              // Whether it's a self-closing tag
}
```

Represents the position and metadata of an XML element.

### Range

```go
type Range struct {
    Start Position // Start position
    End   Position // End position
}
```

Represents a range in the file.

### Position

```go
type Position struct {
    Line   int // Line number (1-based)
    Column int // Column number (1-based)
    Offset int // Byte offset from start of file
}
```

Represents a specific position in the file.

## Constructors

### NewConfigParser

```go
func NewConfigParser() *ConfigParser
```

Creates a new configuration parser with default settings.

**Returns:**
- `*ConfigParser`: New parser instance

**Example:**
```go
parser := parser.NewConfigParser()
config, err := parser.ParseFromFile("/path/to/NuGet.Config")
```

### NewPositionAwareParser

```go
func NewPositionAwareParser() *ConfigParser
```

Creates a new parser with position tracking enabled.

**Returns:**
- `*ConfigParser`: New parser instance with position tracking

**Example:**
```go
parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromFileWithPositions("/path/to/NuGet.Config")
```

## Parsing Methods

### ParseFromFile

```go
func (p *ConfigParser) ParseFromFile(filePath string) (*types.NuGetConfig, error)
```

Parses a NuGet configuration file from the specified path.

**Parameters:**
- `filePath` (string): Path to the configuration file

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Errors:**
- `errors.ErrConfigFileNotFound`: File doesn't exist
- `errors.ErrEmptyConfigFile`: File is empty
- `errors.ErrInvalidConfigFormat`: Invalid XML format
- `*errors.ParseError`: Detailed parsing error with position information

**Example:**
```go
parser := parser.NewConfigParser()
config, err := parser.ParseFromFile("/path/to/NuGet.Config")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("Configuration file not found")
    } else if errors.IsParseError(err) {
        fmt.Printf("Parse error: %v\n", err)
    }
    return
}

fmt.Printf("Loaded %d package sources\n", len(config.PackageSources.Add))
```

### ParseFromString

```go
func (p *ConfigParser) ParseFromString(content string) (*types.NuGetConfig, error)
```

Parses a NuGet configuration from an XML string.

**Parameters:**
- `content` (string): XML content as string

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Example:**
```go
xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`

parser := parser.NewConfigParser()
config, err := parser.ParseFromString(xmlContent)
if err != nil {
    log.Fatalf("Failed to parse XML: %v", err)
}
```

### ParseFromReader

```go
func (p *ConfigParser) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)
```

Parses a NuGet configuration from an io.Reader.

**Parameters:**
- `reader` (io.Reader): Reader containing XML content

**Returns:**
- `*types.NuGetConfig`: Parsed configuration object
- `error`: Error if parsing fails

**Example:**
```go
file, err := os.Open("NuGet.Config")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

parser := parser.NewConfigParser()
config, err := parser.ParseFromReader(file)
if err != nil {
    log.Fatalf("Failed to parse from reader: %v", err)
}
```

## Position-Aware Parsing

### ParseFromFileWithPositions

```go
func (p *ConfigParser) ParseFromFileWithPositions(filePath string) (*ParseResult, error)
```

Parses a configuration file while tracking element positions.

**Parameters:**
- `filePath` (string): Path to the configuration file

**Returns:**
- `*ParseResult`: Parse result with position information
- `error`: Error if parsing fails

**Example:**
```go
parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to parse with positions: %v", err)
}

// Access the configuration
config := result.Config
fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))

// Access position information
for path, pos := range result.Positions {
    fmt.Printf("Element %s at line %d\n", path, pos.Range.Start.Line)
}
```

### ParseFromContentWithPositions

```go
func (p *ConfigParser) ParseFromContentWithPositions(content []byte) (*ParseResult, error)
```

Parses configuration content while tracking element positions.

**Parameters:**
- `content` ([]byte): XML content as bytes

**Returns:**
- `*ParseResult`: Parse result with position information
- `error`: Error if parsing fails

**Example:**
```go
content := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`)

parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromContentWithPositions(content)
if err != nil {
    log.Fatalf("Failed to parse content: %v", err)
}
```

## Serialization Methods

### SaveToFile

```go
func (p *ConfigParser) SaveToFile(config *types.NuGetConfig, filePath string) error
```

Serializes a configuration object to an XML file.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration to save
- `filePath` (string): Target file path

**Returns:**
- `error`: Error if saving fails

**Example:**
```go
parser := parser.NewConfigParser()
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:   "nuget.org",
                Value: "https://api.nuget.org/v3/index.json",
            },
        },
    },
}

err := parser.SaveToFile(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("Failed to save config: %v", err)
}
```

### SerializeToXML

```go
func (p *ConfigParser) SerializeToXML(config *types.NuGetConfig) (string, error)
```

Serializes a configuration object to an XML string.

**Parameters:**
- `config` (*types.NuGetConfig): Configuration to serialize

**Returns:**
- `string`: XML representation
- `error`: Error if serialization fails

**Example:**
```go
parser := parser.NewConfigParser()
xmlString, err := parser.SerializeToXML(config)
if err != nil {
    log.Fatalf("Failed to serialize: %v", err)
}

fmt.Println("Generated XML:")
fmt.Println(xmlString)
```

## Discovery Methods

### FindAndParseConfig

```go
func (p *ConfigParser) FindAndParseConfig() (*types.NuGetConfig, string, error)
```

Finds and parses the first available configuration file.

**Returns:**
- `*types.NuGetConfig`: Parsed configuration
- `string`: Path to the found configuration file
- `error`: Error if no file found or parsing fails

**Example:**
```go
parser := parser.NewConfigParser()
config, configPath, err := parser.FindAndParseConfig()
if err != nil {
    log.Fatalf("Failed to find and parse config: %v", err)
}

fmt.Printf("Using config from: %s\n", configPath)
```

## Error Handling

The parser provides detailed error information through structured error types:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := parser.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // Handle file not found
        fmt.Println("Configuration file not found")
    } else if errors.IsParseError(err) {
        // Handle parsing error with details
        parseErr := err.(*errors.ParseError)
        fmt.Printf("Parse error at line %d: %s\n", parseErr.Line, parseErr.Context)
    } else if errors.IsFormatError(err) {
        // Handle format error
        fmt.Println("Invalid configuration format")
    }
}
```

## Best Practices

1. **Use appropriate parser type**: Use `NewConfigParser()` for simple parsing, `NewPositionAwareParser()` for editing scenarios
2. **Handle errors properly**: Always check for specific error types using the error utilities
3. **Validate input**: Ensure file paths exist and content is valid before parsing
4. **Resource management**: Close file handles properly when using `ParseFromReader`
5. **Position tracking**: Only enable position tracking when needed for editing to avoid overhead

## Thread Safety

The ConfigParser is not thread-safe. Create separate instances for concurrent use or provide appropriate synchronization.
