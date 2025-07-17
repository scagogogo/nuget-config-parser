# Errors API

The `pkg/errors` package defines error types and utilities for handling NuGet configuration parsing and management errors.

## Overview

The Errors API provides:
- Structured error types for different failure scenarios
- Error classification utilities
- Detailed error information with context
- Support for error wrapping and unwrapping

## Error Constants

### Predefined Errors

```go
var (
    // ErrInvalidConfigFormat indicates invalid configuration file format
    ErrInvalidConfigFormat = errors.New("invalid nuget config format")

    // ErrConfigFileNotFound indicates configuration file not found
    ErrConfigFileNotFound = errors.New("nuget config file not found")

    // ErrEmptyConfigFile indicates empty configuration file
    ErrEmptyConfigFile = errors.New("empty nuget config file")

    // ErrXMLParsing indicates XML parsing error
    ErrXMLParsing = errors.New("xml parsing error")

    // ErrMissingRequiredElement indicates missing required element
    ErrMissingRequiredElement = errors.New("missing required element in config")
)
```

These predefined errors represent common failure scenarios and can be used with `errors.Is()` for error checking.

## Error Types

### ParseError

```go
type ParseError struct {
    BaseErr  error  // Base error
    Line     int    // Line number where error occurred
    Position int    // Position in line where error occurred
    Context  string // Additional context information
}
```

Represents a parsing error with detailed position information.

**Fields:**
- `BaseErr`: The underlying error that caused the parsing failure
- `Line`: Line number where the error occurred (1-based)
- `Position`: Character position in the line where the error occurred (1-based)
- `Context`: Additional context information about the error

**Methods:**

#### Error

```go
func (e *ParseError) Error() string
```

Returns a formatted error message with position information.

**Example:**
```go
parseErr := &errors.ParseError{
    BaseErr:  errors.ErrInvalidConfigFormat,
    Line:     15,
    Position: 23,
    Context:  "invalid attribute value",
}

fmt.Println(parseErr.Error())
// Output: parse error at line 15 position 23: invalid attribute value - invalid nuget config format
```

#### Unwrap

```go
func (e *ParseError) Unwrap() error
```

Returns the base error, supporting `errors.Is()` and `errors.As()` functions.

**Example:**
```go
parseErr := &errors.ParseError{
    BaseErr: errors.ErrInvalidConfigFormat,
    Line:    10,
    Position: 5,
    Context: "malformed XML",
}

// Check if it's a format error
if errors.Is(parseErr, errors.ErrInvalidConfigFormat) {
    fmt.Println("This is a format error")
}
```

## Constructor Functions

### NewParseError

```go
func NewParseError(baseErr error, line, position int, context string) *ParseError
```

Creates a new parsing error with position information.

**Parameters:**
- `baseErr` (error): The underlying error
- `line` (int): Line number where error occurred
- `position` (int): Position in line where error occurred
- `context` (string): Additional context information

**Returns:**
- `*ParseError`: New parse error instance

**Example:**
```go
// Create a parse error for invalid XML
parseErr := errors.NewParseError(
    errors.ErrXMLParsing,
    25,
    10,
    "unexpected closing tag",
)

fmt.Printf("Parse error: %v\n", parseErr)
```

## Error Classification Functions

### IsNotFoundError

```go
func IsNotFoundError(err error) bool
```

Checks if an error indicates a configuration file not found.

**Parameters:**
- `err` (error): Error to check

**Returns:**
- `bool`: True if the error indicates file not found

**Example:**
```go
config, err := api.ParseFromFile("missing.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("Configuration file not found")
        // Create default configuration
        config = api.CreateDefaultConfig()
    } else {
        log.Fatalf("Other error: %v", err)
    }
}
```

### IsParseError

```go
func IsParseError(err error) bool
```

Checks if an error is a parsing error.

**Parameters:**
- `err` (error): Error to check

**Returns:**
- `bool`: True if the error is a ParseError

**Example:**
```go
config, err := api.ParseFromString(invalidXML)
if err != nil {
    if errors.IsParseError(err) {
        var parseErr *errors.ParseError
        if errors.As(err, &parseErr) {
            fmt.Printf("Parse error at line %d: %s\n", 
                parseErr.Line, parseErr.Context)
        }
    }
}
```

### IsFormatError

```go
func IsFormatError(err error) bool
```

Checks if an error indicates invalid configuration format.

**Parameters:**
- `err` (error): Error to check

**Returns:**
- `bool`: True if the error indicates format issues

**Example:**
```go
config, err := api.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsFormatError(err) {
        fmt.Println("Invalid configuration format")
        // Handle format error
    }
}
```

## Usage Examples

### Basic Error Handling

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
    
    config, err := api.ParseFromFile("NuGet.Config")
    if err != nil {
        handleError(err)
        return
    }
    
    fmt.Printf("Successfully loaded %d package sources\n", 
        len(config.PackageSources.Add))
}

func handleError(err error) {
    switch {
    case errors.IsNotFoundError(err):
        fmt.Println("Configuration file not found")
        fmt.Println("Consider creating a default configuration")
        
    case errors.IsParseError(err):
        var parseErr *errors.ParseError
        if errors.As(err, &parseErr) {
            fmt.Printf("Parse error at line %d, position %d: %s\n",
                parseErr.Line, parseErr.Position, parseErr.Context)
            fmt.Printf("Underlying error: %v\n", parseErr.BaseErr)
        } else {
            fmt.Printf("Parse error: %v\n", err)
        }
        
    case errors.IsFormatError(err):
        fmt.Println("Invalid configuration file format")
        fmt.Println("Please check the XML structure")
        
    default:
        log.Printf("Unexpected error: %v\n", err)
    }
}
```

### Advanced Error Handling

```go
func parseConfigWithRecovery(filePath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    config, err := api.ParseFromFile(filePath)
    if err != nil {
        // Try to provide helpful error information
        if errors.IsNotFoundError(err) {
            return nil, fmt.Errorf("configuration file '%s' not found: %w", filePath, err)
        }
        
        if errors.IsParseError(err) {
            var parseErr *errors.ParseError
            if errors.As(err, &parseErr) {
                // Provide detailed parse error information
                return nil, fmt.Errorf(
                    "failed to parse configuration at %s:%d:%d - %s: %w",
                    filePath, parseErr.Line, parseErr.Position, 
                    parseErr.Context, parseErr.BaseErr)
            }
        }
        
        // For other errors, wrap with context
        return nil, fmt.Errorf("failed to load configuration from '%s': %w", filePath, err)
    }
    
    return config, nil
}
```

### Error Recovery Strategies

```go
func loadConfigWithFallback(primaryPath, fallbackPath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    // Try primary configuration
    config, err := api.ParseFromFile(primaryPath)
    if err == nil {
        return config, nil
    }
    
    // Handle primary config errors
    if errors.IsNotFoundError(err) {
        fmt.Printf("Primary config not found at %s, trying fallback...\n", primaryPath)
    } else if errors.IsParseError(err) {
        fmt.Printf("Primary config has parse errors, trying fallback...\n")
    } else {
        return nil, fmt.Errorf("failed to load primary config: %w", err)
    }
    
    // Try fallback configuration
    config, err = api.ParseFromFile(fallbackPath)
    if err == nil {
        fmt.Printf("Using fallback configuration from %s\n", fallbackPath)
        return config, nil
    }
    
    // If fallback also fails, create default
    if errors.IsNotFoundError(err) {
        fmt.Println("No configuration files found, creating default...")
        return api.CreateDefaultConfig(), nil
    }
    
    return nil, fmt.Errorf("failed to load any configuration: %w", err)
}
```

### Custom Error Types

```go
// Custom error for validation failures
type ValidationError struct {
    Field   string
    Value   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s' with value '%s': %s", 
        e.Field, e.Value, e.Message)
}

// Validate package source
func validatePackageSource(source *types.PackageSource) error {
    if source.Key == "" {
        return &ValidationError{
            Field:   "Key",
            Value:   source.Key,
            Message: "package source key cannot be empty",
        }
    }
    
    if source.Value == "" {
        return &ValidationError{
            Field:   "Value",
            Value:   source.Value,
            Message: "package source value cannot be empty",
        }
    }
    
    // Validate protocol version if specified
    if source.ProtocolVersion != "" && 
       source.ProtocolVersion != "2" && 
       source.ProtocolVersion != "3" {
        return &ValidationError{
            Field:   "ProtocolVersion",
            Value:   source.ProtocolVersion,
            Message: "protocol version must be '2' or '3'",
        }
    }
    
    return nil
}
```

## Error Wrapping Best Practices

### Adding Context

```go
func parseConfigFromPath(configPath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        // Wrap error with additional context
        return nil, fmt.Errorf("failed to parse NuGet configuration from '%s': %w", 
            configPath, err)
    }
    
    return config, nil
}
```

### Preserving Error Types

```go
func loadAndValidateConfig(configPath string) (*types.NuGetConfig, error) {
    config, err := parseConfigFromPath(configPath)
    if err != nil {
        // Check original error type even after wrapping
        if errors.IsNotFoundError(err) {
            // Handle not found case
            return createDefaultConfig(configPath)
        }
        
        if errors.IsParseError(err) {
            // Handle parse error case
            return nil, fmt.Errorf("configuration file has syntax errors: %w", err)
        }
        
        return nil, err
    }
    
    // Additional validation...
    return config, nil
}
```

## Testing Error Conditions

```go
func TestErrorHandling(t *testing.T) {
    api := nuget.NewAPI()
    
    // Test file not found
    _, err := api.ParseFromFile("nonexistent.config")
    if !errors.IsNotFoundError(err) {
        t.Errorf("Expected not found error, got: %v", err)
    }
    
    // Test invalid XML
    invalidXML := "<configuration><packageSources><add key="
    _, err = api.ParseFromString(invalidXML)
    if !errors.IsParseError(err) {
        t.Errorf("Expected parse error, got: %v", err)
    }
    
    // Test parse error details
    var parseErr *errors.ParseError
    if errors.As(err, &parseErr) {
        if parseErr.Line <= 0 {
            t.Errorf("Expected positive line number, got: %d", parseErr.Line)
        }
    }
}
```

## Best Practices

1. **Use error classification functions**: Always use `IsNotFoundError()`, `IsParseError()`, etc. for error checking
2. **Provide context**: Wrap errors with additional context using `fmt.Errorf()` with `%w` verb
3. **Handle specific errors**: Provide different handling for different error types
4. **Preserve error chains**: Use error wrapping to maintain the original error information
5. **Test error conditions**: Write tests for different error scenarios
6. **Log appropriately**: Use different log levels for different error types
7. **Provide recovery**: Implement fallback strategies for recoverable errors

## Thread Safety

Error types and functions in this package are thread-safe and can be used concurrently.
