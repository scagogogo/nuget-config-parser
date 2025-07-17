# Installation

This guide covers different ways to install and set up the NuGet Config Parser library in your Go project.

## Requirements

- **Go Version**: 1.19 or later
- **Operating System**: Windows, Linux, or macOS
- **Dependencies**: No external dependencies required

## Installation Methods

### Using Go Modules (Recommended)

The easiest way to install the library is using Go modules:

```bash
go get github.com/scagogogo/nuget-config-parser
```

This will download the latest version of the library and add it to your `go.mod` file.

### Specific Version

To install a specific version:

```bash
go get github.com/scagogogo/nuget-config-parser@v1.0.0
```

### Latest Development Version

To get the latest development version from the main branch:

```bash
go get github.com/scagogogo/nuget-config-parser@main
```

## Importing the Library

Once installed, import the library in your Go code:

```go
import "github.com/scagogogo/nuget-config-parser/pkg/nuget"
```

For specific packages:

```go
import (
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
    "github.com/scagogogo/nuget-config-parser/pkg/parser"
    "github.com/scagogogo/nuget-config-parser/pkg/editor"
)
```

## Verification

Create a simple test file to verify the installation:

```go
// test_installation.go
package main

import (
    "fmt"
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    fmt.Println("NuGet Config Parser installed successfully!")
    
    // Try to find a config file
    configPath, err := api.FindConfigFile()
    if err != nil {
        fmt.Println("No config file found, but library is working!")
    } else {
        fmt.Printf("Found config file: %s\n", configPath)
    }
}
```

Run the test:

```bash
go run test_installation.go
```

## Project Setup

### Initialize a New Project

If you're starting a new project:

```bash
mkdir my-nuget-project
cd my-nuget-project
go mod init my-nuget-project
go get github.com/scagogogo/nuget-config-parser
```

### Add to Existing Project

If you have an existing Go project:

```bash
cd your-existing-project
go get github.com/scagogogo/nuget-config-parser
```

## IDE Setup

### VS Code

For VS Code users, make sure you have the Go extension installed. The library should work out of the box with IntelliSense and auto-completion.

### GoLand

GoLand should automatically recognize the library after running `go get`. If you encounter issues, try:

1. File → Invalidate Caches and Restart
2. Ensure Go modules are enabled in Settings → Go → Go Modules

### Other IDEs

Most Go-compatible IDEs should work with the library automatically. Ensure your IDE supports Go modules.

## Troubleshooting

### Common Issues

#### Module Not Found

If you get a "module not found" error:

```bash
go: module github.com/scagogogo/nuget-config-parser: not found
```

Make sure you have:
1. Go modules enabled (`go mod init` if needed)
2. Internet connection for downloading
3. Correct module path

#### Version Conflicts

If you encounter version conflicts:

```bash
go mod tidy
go get github.com/scagogogo/nuget-config-parser@latest
```

#### Proxy Issues

If you're behind a corporate proxy:

```bash
go env -w GOPROXY=direct
go env -w GOSUMDB=off
```

Or configure your proxy settings:

```bash
go env -w GOPROXY=https://your-proxy.com
```

### Verification Commands

Check your Go environment:

```bash
go version
go env GOMOD
go env GOPROXY
```

List installed modules:

```bash
go list -m all | grep nuget-config-parser
```

## Build Tags

The library doesn't require any special build tags, but you can use standard Go build tags for platform-specific code:

```go
//go:build windows
// +build windows

// Windows-specific code
```

## Docker Usage

If you're using Docker, add the library to your Dockerfile:

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# The library will be downloaded with other dependencies
COPY . .
RUN go build -o main .

CMD ["./main"]
```

## Performance Considerations

The library is designed to be lightweight with no external dependencies. However, for optimal performance:

- Use the API instance efficiently (create once, reuse)
- Cache parsed configurations when possible
- Use position-aware editing for minimal file changes

## Next Steps

Now that you have the library installed:

1. Read the [Getting Started](./getting-started.md) guide
2. Try the [Quick Start](./quick-start.md) examples
3. Explore the [API Reference](/api/) for detailed documentation
4. Check out [Examples](/examples/) for real-world usage patterns
