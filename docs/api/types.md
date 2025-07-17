# Types API

The `pkg/types` package defines all the data structures used to represent NuGet configuration files. These types correspond directly to the XML structure of NuGet.Config files.

## Core Configuration Types

### NuGetConfig

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

The root configuration structure representing a complete NuGet.Config file.

**Fields:**
- `PackageSources`: List of available package sources (required)
- `PackageSourceCredentials`: Credentials for package sources (optional)
- `Config`: Global configuration options (optional)
- `DisabledPackageSources`: List of disabled package sources (optional)
- `ActivePackageSource`: Currently active package source (optional)

**Example:**
```go
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:             "nuget.org",
                Value:           "https://api.nuget.org/v3/index.json",
                ProtocolVersion: "3",
            },
        },
    },
}
```

## Package Source Types

### PackageSources

```go
type PackageSources struct {
    Clear bool            `xml:"clear,attr,omitempty"`
    Add   []PackageSource `xml:"add"`
}
```

Container for package source definitions.

**Fields:**
- `Clear`: If true, clears all previously defined package sources
- `Add`: List of package sources to add

### PackageSource

```go
type PackageSource struct {
    Key             string `xml:"key,attr"`
    Value           string `xml:"value,attr"`
    ProtocolVersion string `xml:"protocolVersion,attr,omitempty"`
}
```

Represents a single package source.

**Fields:**
- `Key`: Unique identifier for the package source
- `Value`: URL or file path to the package source
- `ProtocolVersion`: NuGet protocol version (optional, typically "2" or "3")

**Example:**
```go
source := types.PackageSource{
    Key:             "company-feed",
    Value:           "https://nuget.company.com/v3/index.json",
    ProtocolVersion: "3",
}
```

### DisabledPackageSources

```go
type DisabledPackageSources struct {
    Add []DisabledSource `xml:"add"`
}
```

Container for disabled package sources.

### DisabledSource

```go
type DisabledSource struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

Represents a disabled package source.

**Fields:**
- `Key`: Key of the package source to disable
- `Value`: Usually "true" to indicate the source is disabled

### ActivePackageSource

```go
type ActivePackageSource struct {
    Add PackageSource `xml:"add"`
}
```

Represents the currently active package source.

## Credential Types

### PackageSourceCredentials

```go
type PackageSourceCredentials struct {
    Sources map[string]SourceCredential `xml:"-"`
}
```

Container for package source credentials. The `Sources` map uses package source keys as keys and credentials as values.

**Note:** This type has custom XML marshaling/unmarshaling logic to handle the dynamic structure of credentials in NuGet.Config files.

### SourceCredential

```go
type SourceCredential struct {
    Add []Credential `xml:"add"`
}
```

Credentials for a specific package source.

### Credential

```go
type Credential struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

A single credential key-value pair.

**Common credential keys:**
- `Username`: The username for authentication
- `Password`: The password for authentication (usually encrypted)
- `ClearTextPassword`: The password in clear text (not recommended for production)

**Example:**
```go
credentials := types.SourceCredential{
    Add: []types.Credential{
        {Key: "Username", Value: "myuser"},
        {Key: "ClearTextPassword", Value: "mypassword"},
    },
}
```

## Configuration Option Types

### Config

```go
type Config struct {
    Add []ConfigOption `xml:"add"`
}
```

Container for global configuration options.

### ConfigOption

```go
type ConfigOption struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

A single configuration option.

**Common configuration keys:**
- `globalPackagesFolder`: Path to the global packages folder
- `repositoryPath`: Path to the packages repository
- `defaultPushSource`: Default source for package publishing
- `http_proxy`: HTTP proxy settings
- `http_proxy.user`: HTTP proxy username
- `http_proxy.password`: HTTP proxy password

**Example:**
```go
configOptions := []types.ConfigOption{
    {Key: "globalPackagesFolder", Value: "/custom/packages/path"},
    {Key: "defaultPushSource", Value: "https://my-nuget-server.com"},
}
```

## XML Marshaling

All types support XML marshaling and unmarshaling through Go's `encoding/xml` package. The struct tags define how each field maps to XML elements and attributes.

### Custom Marshaling

The `PackageSourceCredentials` type implements custom XML marshaling to handle the dynamic structure where each package source has its own XML element:

```xml
<packageSourceCredentials>
  <MyPrivateSource>
    <add key="Username" value="myuser" />
    <add key="ClearTextPassword" value="mypass" />
  </MyPrivateSource>
</packageSourceCredentials>
```

## Usage Examples

### Creating a Complete Configuration

```go
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:             "nuget.org",
                Value:           "https://api.nuget.org/v3/index.json",
                ProtocolVersion: "3",
            },
            {
                Key:   "local",
                Value: "/path/to/local/packages",
            },
        },
    },
    PackageSourceCredentials: &types.PackageSourceCredentials{
        Sources: map[string]types.SourceCredential{
            "private-feed": {
                Add: []types.Credential{
                    {Key: "Username", Value: "user"},
                    {Key: "ClearTextPassword", Value: "pass"},
                },
            },
        },
    },
    Config: &types.Config{
        Add: []types.ConfigOption{
            {Key: "globalPackagesFolder", Value: "/custom/packages"},
        },
    },
    DisabledPackageSources: &types.DisabledPackageSources{
        Add: []types.DisabledSource{
            {Key: "local", Value: "true"},
        },
    },
    ActivePackageSource: &types.ActivePackageSource{
        Add: types.PackageSource{
            Key:   "nuget.org",
            Value: "https://api.nuget.org/v3/index.json",
        },
    },
}
```

### Working with Package Sources

```go
// Add a new package source
newSource := types.PackageSource{
    Key:             "company-feed",
    Value:           "https://nuget.company.com/v3/index.json",
    ProtocolVersion: "3",
}
config.PackageSources.Add = append(config.PackageSources.Add, newSource)

// Find a package source
var foundSource *types.PackageSource
for i, source := range config.PackageSources.Add {
    if source.Key == "company-feed" {
        foundSource = &config.PackageSources.Add[i]
        break
    }
}

// Remove a package source
for i, source := range config.PackageSources.Add {
    if source.Key == "company-feed" {
        config.PackageSources.Add = append(
            config.PackageSources.Add[:i],
            config.PackageSources.Add[i+1:]...,
        )
        break
    }
}
```

### Working with Credentials

```go
// Initialize credentials if nil
if config.PackageSourceCredentials == nil {
    config.PackageSourceCredentials = &types.PackageSourceCredentials{
        Sources: make(map[string]types.SourceCredential),
    }
}

// Add credentials for a source
config.PackageSourceCredentials.Sources["private-feed"] = types.SourceCredential{
    Add: []types.Credential{
        {Key: "Username", Value: "myuser"},
        {Key: "ClearTextPassword", Value: "mypass"},
    },
}

// Get credentials for a source
if cred, exists := config.PackageSourceCredentials.Sources["private-feed"]; exists {
    for _, c := range cred.Add {
        if c.Key == "Username" {
            fmt.Printf("Username: %s\n", c.Value)
        }
    }
}
```

### Working with Configuration Options

```go
// Initialize config if nil
if config.Config == nil {
    config.Config = &types.Config{
        Add: []types.ConfigOption{},
    }
}

// Add a configuration option
config.Config.Add = append(config.Config.Add, types.ConfigOption{
    Key:   "globalPackagesFolder",
    Value: "/custom/packages/path",
})

// Find a configuration option
var globalPackagesFolder string
for _, option := range config.Config.Add {
    if option.Key == "globalPackagesFolder" {
        globalPackagesFolder = option.Value
        break
    }
}
```

## Validation

While the types themselves don't include validation logic, you should validate data when creating or modifying configurations:

```go
func validatePackageSource(source types.PackageSource) error {
    if source.Key == "" {
        return errors.New("package source key cannot be empty")
    }
    if source.Value == "" {
        return errors.New("package source value cannot be empty")
    }
    if source.ProtocolVersion != "" && 
       source.ProtocolVersion != "2" && 
       source.ProtocolVersion != "3" {
        return errors.New("invalid protocol version")
    }
    return nil
}
```

## Thread Safety

The types in this package are not thread-safe. If you need to access or modify configuration objects from multiple goroutines, you must provide your own synchronization.

## Best Practices

1. **Initialize Optional Fields**: Always check if optional fields are nil before accessing them
2. **Use Pointers for Optional Structs**: Optional configuration sections use pointers to distinguish between empty and missing
3. **Validate Data**: Validate package source keys, URLs, and other data before creating configurations
4. **Handle Credentials Carefully**: Be cautious with credential handling, especially clear text passwords
5. **Use the API**: Prefer using the high-level API methods over direct struct manipulation when possible
