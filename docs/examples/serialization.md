# Serialization

This example demonstrates how to serialize NuGet configuration objects to XML and handle custom serialization scenarios.

## Overview

Serialization involves:
- Converting configuration objects to XML
- Formatting XML output
- Handling custom serialization requirements
- Working with XML templates
- Validating serialized output

## Example 1: Basic Serialization

Converting configuration objects to XML:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // Create a sample configuration
    config := api.CreateDefaultConfig()
    
    // Add some additional sources and settings
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    api.AddPackageSource(config, "local-dev", "./packages", "")
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    api.AddCredential(config, "company-feed", "user", "pass")
    
    fmt.Println("=== Basic Serialization ===")
    
    // Serialize to XML
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Fatalf("Serialization failed: %v", err)
    }
    
    fmt.Println("Serialized XML:")
    fmt.Println(xmlContent)
    
    // Verify by parsing back
    fmt.Println("\n=== Verification ===")
    parsedConfig, err := api.ParseFromString(xmlContent)
    if err != nil {
        log.Fatalf("Failed to parse serialized XML: %v", err)
    }
    
    fmt.Printf("Original sources: %d\n", len(config.PackageSources.Add))
    fmt.Printf("Parsed sources: %d\n", len(parsedConfig.PackageSources.Add))
    
    if len(config.PackageSources.Add) == len(parsedConfig.PackageSources.Add) {
        fmt.Println("✅ Serialization verification successful")
    } else {
        fmt.Println("❌ Serialization verification failed")
    }
}
```

## Example 2: Formatted XML Output

Creating well-formatted XML with proper indentation:

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
    config := api.CreateDefaultConfig()
    
    // Build a comprehensive configuration
    buildComprehensiveConfig(api, config)
    
    fmt.Println("=== Formatted XML Serialization ===")
    
    // Serialize with formatting
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Fatalf("Serialization failed: %v", err)
    }
    
    // Display formatted XML
    fmt.Println("Formatted XML output:")
    fmt.Println(xmlContent)
    
    // Save to file with proper formatting
    outputFile := "FormattedConfig.xml"
    err = os.WriteFile(outputFile, []byte(xmlContent), 0644)
    if err != nil {
        log.Fatalf("Failed to write file: %v", err)
    }
    
    fmt.Printf("\nFormatted XML saved to: %s\n", outputFile)
    
    // Display file statistics
    info, _ := os.Stat(outputFile)
    fmt.Printf("File size: %d bytes\n", info.Size())
    
    // Validate the output
    validateXMLOutput(api, xmlContent)
}

func buildComprehensiveConfig(api *nuget.API, config *types.NuGetConfig) {
    // Add multiple package sources
    sources := []struct {
        key, url, version string
    }{
        {"nuget.org", "https://api.nuget.org/v3/index.json", "3"},
        {"company-stable", "https://stable.company.com/nuget", "3"},
        {"company-preview", "https://preview.company.com/nuget", "3"},
        {"local-packages", "/path/to/local/packages", ""},
    }
    
    for _, source := range sources {
        api.AddPackageSource(config, source.key, source.url, source.version)
    }
    
    // Add credentials
    api.AddCredential(config, "company-stable", "employee", "secret123")
    api.AddCredential(config, "company-preview", "employee", "preview_token")
    
    // Add configuration options
    options := map[string]string{
        "globalPackagesFolder":     "/custom/packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "Highest",
        "automaticPackageRestore": "true",
        "defaultPushSource":       "https://stable.company.com/nuget",
    }
    
    for key, value := range options {
        api.AddConfigOption(config, key, value)
    }
    
    // Disable preview source
    api.DisablePackageSource(config, "company-preview")
    
    // Set active source
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
}

func validateXMLOutput(api *nuget.API, xmlContent string) {
    fmt.Println("\n=== XML Validation ===")
    
    // Parse the serialized XML
    parsedConfig, err := api.ParseFromString(xmlContent)
    if err != nil {
        fmt.Printf("❌ XML validation failed: %v\n", err)
        return
    }
    
    fmt.Println("✅ XML is valid and parseable")
    
    // Check structure completeness
    checks := []struct {
        name      string
        condition bool
    }{
        {"Package sources present", len(parsedConfig.PackageSources.Add) > 0},
        {"Credentials present", parsedConfig.PackageSourceCredentials != nil},
        {"Config options present", parsedConfig.Config != nil && len(parsedConfig.Config.Add) > 0},
        {"Active source set", parsedConfig.ActivePackageSource != nil},
        {"Disabled sources present", parsedConfig.DisabledPackageSources != nil},
    }
    
    for _, check := range checks {
        status := "❌"
        if check.condition {
            status = "✅"
        }
        fmt.Printf("%s %s\n", status, check.name)
    }
}
```

## Example 3: Custom XML Templates

Working with XML templates and custom serialization:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "text/template"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Println("=== Custom XML Templates ===")
    
    // Create configuration data
    configData := ConfigTemplateData{
        PackageSources: []PackageSourceData{
            {"nuget.org", "https://api.nuget.org/v3/index.json", "3"},
            {"company-feed", "https://nuget.company.com/v3/index.json", "3"},
        },
        GlobalPackagesFolder: "/custom/packages",
        RepositoryPath:      "./packages",
        DependencyVersion:   "Highest",
        ActiveSource:        "nuget.org",
    }
    
    // Generate XML using template
    xmlContent, err := generateXMLFromTemplate(configData)
    if err != nil {
        log.Fatalf("Template generation failed: %v", err)
    }
    
    fmt.Println("Generated XML from template:")
    fmt.Println(xmlContent)
    
    // Parse the generated XML
    config, err := api.ParseFromString(xmlContent)
    if err != nil {
        log.Fatalf("Failed to parse template XML: %v", err)
    }
    
    fmt.Printf("\nParsed %d package sources from template\n", len(config.PackageSources.Add))
    
    // Compare with standard serialization
    compareWithStandardSerialization(api, config)
}

type ConfigTemplateData struct {
    PackageSources       []PackageSourceData
    GlobalPackagesFolder string
    RepositoryPath       string
    DependencyVersion    string
    ActiveSource         string
}

type PackageSourceData struct {
    Key             string
    Value           string
    ProtocolVersion string
}

const xmlTemplate = `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
{{- range .PackageSources }}
    <add key="{{ .Key }}" value="{{ .Value }}"{{if .ProtocolVersion}} protocolVersion="{{ .ProtocolVersion }}"{{end}} />
{{- end }}
  </packageSources>
  
  <activePackageSource>
    <add key="{{ .ActiveSource }}" value="{{range .PackageSources}}{{if eq .Key $.ActiveSource}}{{.Value}}{{end}}{{end}}" />
  </activePackageSource>
  
  <config>
    <add key="globalPackagesFolder" value="{{ .GlobalPackagesFolder }}" />
    <add key="repositoryPath" value="{{ .RepositoryPath }}" />
    <add key="dependencyVersion" value="{{ .DependencyVersion }}" />
  </config>
</configuration>`

func generateXMLFromTemplate(data ConfigTemplateData) (string, error) {
    tmpl, err := template.New("nuget-config").Parse(xmlTemplate)
    if err != nil {
        return "", fmt.Errorf("template parsing failed: %w", err)
    }
    
    var buf strings.Builder
    err = tmpl.Execute(&buf, data)
    if err != nil {
        return "", fmt.Errorf("template execution failed: %w", err)
    }
    
    return buf.String(), nil
}

func compareWithStandardSerialization(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== Comparison with Standard Serialization ===")
    
    // Serialize using standard method
    standardXML, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("Standard serialization failed: %v", err)
        return
    }
    
    fmt.Println("Standard serialization:")
    fmt.Println(standardXML)
    
    // Compare lengths
    fmt.Printf("\nTemplate XML length: %d characters\n", len(xmlTemplate))
    fmt.Printf("Standard XML length: %d characters\n", len(standardXML))
}
```

## Example 4: Batch Serialization

Serializing multiple configurations:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Println("=== Batch Serialization ===")
    
    // Create multiple environment configurations
    environments := []string{"development", "staging", "production"}
    
    for _, env := range environments {
        fmt.Printf("\nCreating %s configuration...\n", env)
        
        config := createEnvironmentConfig(api, env)
        
        // Serialize configuration
        xmlContent, err := api.SerializeToXML(config)
        if err != nil {
            log.Printf("Failed to serialize %s config: %v", env, err)
            continue
        }
        
        // Save to environment-specific file
        filename := fmt.Sprintf("NuGet.%s.Config", env)
        err = os.WriteFile(filename, []byte(xmlContent), 0644)
        if err != nil {
            log.Printf("Failed to save %s config: %v", env, err)
            continue
        }
        
        fmt.Printf("✅ Saved %s configuration to %s\n", env, filename)
        
        // Display configuration summary
        displayConfigSummary(config, env)
    }
    
    // Create a master configuration that includes all environments
    createMasterConfiguration(api, environments)
}

func createEnvironmentConfig(api *nuget.API, environment string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    switch environment {
    case "development":
        api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
        api.AddPackageSource(config, "local-dev", "./packages", "")
        api.AddPackageSource(config, "preview", "https://api.nuget.org/v3-flatcontainer", "3")
        
        api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
        api.AddConfigOption(config, "dependencyVersion", "Highest")
        api.AddConfigOption(config, "allowPrereleaseVersions", "true")
        
        api.SetActivePackageSource(config, "local-dev", "./packages")
        
    case "staging":
        api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
        api.AddPackageSource(config, "company-staging", "https://staging.company.com/nuget", "3")
        
        api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
        api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
        api.AddConfigOption(config, "allowPrereleaseVersions", "false")
        
        api.AddCredential(config, "company-staging", "staging_user", "staging_pass")
        api.SetActivePackageSource(config, "company-staging", "https://staging.company.com/nuget")
        
    case "production":
        api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
        api.AddPackageSource(config, "company-prod", "https://prod.company.com/nuget", "3")
        
        api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
        api.AddConfigOption(config, "dependencyVersion", "Exact")
        api.AddConfigOption(config, "allowPrereleaseVersions", "false")
        api.AddConfigOption(config, "signatureValidationMode", "require")
        
        api.AddCredential(config, "company-prod", "prod_user", "prod_pass")
        api.SetActivePackageSource(config, "company-prod", "https://prod.company.com/nuget")
    }
    
    return config
}

func displayConfigSummary(config *types.NuGetConfig, environment string) {
    fmt.Printf("  %s configuration summary:\n", environment)
    fmt.Printf("    Package sources: %d\n", len(config.PackageSources.Add))
    
    if config.Config != nil {
        fmt.Printf("    Config options: %d\n", len(config.Config.Add))
    }
    
    if config.PackageSourceCredentials != nil {
        fmt.Printf("    Authenticated sources: %d\n", len(config.PackageSourceCredentials.Sources))
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("    Active source: %s\n", config.ActivePackageSource.Add.Key)
    }
}

func createMasterConfiguration(api *nuget.API, environments []string) {
    fmt.Println("\n=== Creating Master Configuration ===")
    
    masterConfig := api.CreateDefaultConfig()
    
    // Add all possible sources
    allSources := map[string]struct {
        url     string
        version string
    }{
        "nuget.org":        {"https://api.nuget.org/v3/index.json", "3"},
        "local-dev":        {"./packages", ""},
        "preview":          {"https://api.nuget.org/v3-flatcontainer", "3"},
        "company-staging":  {"https://staging.company.com/nuget", "3"},
        "company-prod":     {"https://prod.company.com/nuget", "3"},
    }
    
    for key, source := range allSources {
        api.AddPackageSource(masterConfig, key, source.url, source.version)
    }
    
    // Disable environment-specific sources by default
    api.DisablePackageSource(masterConfig, "local-dev")
    api.DisablePackageSource(masterConfig, "preview")
    api.DisablePackageSource(masterConfig, "company-staging")
    api.DisablePackageSource(masterConfig, "company-prod")
    
    // Add general configuration
    api.AddConfigOption(masterConfig, "globalPackagesFolder", "${NUGET_PACKAGES}")
    api.AddConfigOption(masterConfig, "dependencyVersion", "HighestMinor")
    
    // Serialize master configuration
    xmlContent, err := api.SerializeToXML(masterConfig)
    if err != nil {
        log.Printf("Failed to serialize master config: %v", err)
        return
    }
    
    // Save master configuration
    masterFile := "NuGet.Master.Config"
    err = os.WriteFile(masterFile, []byte(xmlContent), 0644)
    if err != nil {
        log.Printf("Failed to save master config: %v", err)
        return
    }
    
    fmt.Printf("✅ Created master configuration: %s\n", masterFile)
    fmt.Printf("   Contains all sources for %v environments\n", environments)
}
```

## Key Concepts

### Serialization Process

1. **Object to XML**: Convert configuration objects to XML strings
2. **Formatting**: Apply proper indentation and structure
3. **Validation**: Ensure output is valid XML
4. **Round-trip**: Verify by parsing serialized output

### XML Structure

The serialized XML follows the standard NuGet.Config format:
- `<packageSources>`: Package source definitions
- `<packageSourceCredentials>`: Authentication information
- `<config>`: Global configuration options
- `<activePackageSource>`: Active source selection
- `<disabledPackageSources>`: Disabled sources

### Best Practices

1. **Validate output**: Always verify serialized XML is parseable
2. **Format consistently**: Use proper indentation and structure
3. **Handle encoding**: Ensure proper UTF-8 encoding
4. **Escape values**: Properly escape XML special characters
5. **Test round-trips**: Verify parse → serialize → parse cycles

## Common Use Cases

### Configuration Export
```go
// Export current configuration
config, _, _ := api.FindAndParseConfig()
xmlContent, _ := api.SerializeToXML(config)
os.WriteFile("exported-config.xml", []byte(xmlContent), 0644)
```

### Template Generation
```go
// Generate configuration from template
templateData := ConfigTemplateData{...}
xmlContent, _ := generateXMLFromTemplate(templateData)
config, _ := api.ParseFromString(xmlContent)
```

### Batch Processing
```go
// Process multiple configurations
for _, env := range environments {
    config := createEnvironmentConfig(api, env)
    xmlContent, _ := api.SerializeToXML(config)
    saveToFile(fmt.Sprintf("%s.config", env), xmlContent)
}
```

## Next Steps

After mastering serialization:

1. Learn about [Position-Aware Editing](./position-aware-editing.md) for precise modifications
2. Explore the [Parser API](/api/parser) for advanced parsing options
3. Study [Configuration](../guide/configuration.md) for structure details

This guide provides comprehensive examples for serializing NuGet configuration objects to XML, covering various scenarios from basic output to complex template-based generation.
