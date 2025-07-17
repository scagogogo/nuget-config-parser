# 序列化

本示例演示如何将 NuGet 配置对象序列化为 XML 并处理自定义序列化场景。

## 概述

序列化涉及：
- 将配置对象转换为 XML
- 格式化 XML 输出
- 处理自定义序列化需求
- 使用 XML 模板
- 验证序列化输出

## 示例 1: 基本序列化

将配置对象转换为 XML：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 创建示例配置
    config := api.CreateDefaultConfig()
    
    // 添加一些额外的源和设置
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    api.AddPackageSource(config, "local-dev", "./packages", "")
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    api.AddCredential(config, "company-feed", "user", "pass")
    
    fmt.Println("=== 基本序列化 ===")
    
    // 序列化为 XML
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Fatalf("序列化失败: %v", err)
    }
    
    fmt.Println("序列化的 XML:")
    fmt.Println(xmlContent)
    
    // 通过解析回来验证
    fmt.Println("\n=== 验证 ===")
    parsedConfig, err := api.ParseFromString(xmlContent)
    if err != nil {
        log.Fatalf("解析序列化的 XML 失败: %v", err)
    }
    
    fmt.Printf("原始源数量: %d\n", len(config.PackageSources.Add))
    fmt.Printf("解析后源数量: %d\n", len(parsedConfig.PackageSources.Add))
    
    if len(config.PackageSources.Add) == len(parsedConfig.PackageSources.Add) {
        fmt.Println("✅ 序列化验证成功")
    } else {
        fmt.Println("❌ 序列化验证失败")
    }
}
```

## 示例 2: 格式化 XML 输出

创建格式良好的 XML，具有适当的缩进：

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
    
    // 构建综合配置
    buildComprehensiveConfig(api, config)
    
    fmt.Println("=== 格式化 XML 序列化 ===")
    
    // 带格式化的序列化
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Fatalf("序列化失败: %v", err)
    }
    
    // 显示格式化的 XML
    fmt.Println("格式化的 XML 输出:")
    fmt.Println(xmlContent)
    
    // 保存到文件并进行适当格式化
    outputFile := "FormattedConfig.xml"
    err = os.WriteFile(outputFile, []byte(xmlContent), 0644)
    if err != nil {
        log.Fatalf("写入文件失败: %v", err)
    }
    
    fmt.Printf("\n格式化的 XML 已保存到: %s\n", outputFile)
    
    // 显示文件统计信息
    info, _ := os.Stat(outputFile)
    fmt.Printf("文件大小: %d 字节\n", info.Size())
    
    // 验证输出
    validateXMLOutput(api, xmlContent)
}

func buildComprehensiveConfig(api *nuget.API, config *types.NuGetConfig) {
    // 添加多个包源
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
    
    // 添加凭证
    api.AddCredential(config, "company-stable", "employee", "secret123")
    api.AddCredential(config, "company-preview", "employee", "preview_token")
    
    // 添加配置选项
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
    
    // 禁用预览源
    api.DisablePackageSource(config, "company-preview")
    
    // 设置活跃源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
}

func validateXMLOutput(api *nuget.API, xmlContent string) {
    fmt.Println("\n=== XML 验证 ===")
    
    // 解析序列化的 XML
    parsedConfig, err := api.ParseFromString(xmlContent)
    if err != nil {
        fmt.Printf("❌ XML 验证失败: %v\n", err)
        return
    }
    
    fmt.Println("✅ XML 有效且可解析")
    
    // 检查结构完整性
    checks := []struct {
        name      string
        condition bool
    }{
        {"包源存在", len(parsedConfig.PackageSources.Add) > 0},
        {"凭证存在", parsedConfig.PackageSourceCredentials != nil},
        {"配置选项存在", parsedConfig.Config != nil && len(parsedConfig.Config.Add) > 0},
        {"活跃源已设置", parsedConfig.ActivePackageSource != nil},
        {"禁用源存在", parsedConfig.DisabledPackageSources != nil},
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

## 示例 3: 自定义 XML 模板

使用 XML 模板和自定义序列化：

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
    
    fmt.Println("=== 自定义 XML 模板 ===")
    
    // 创建配置数据
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
    
    // 使用模板生成 XML
    xmlContent, err := generateXMLFromTemplate(configData)
    if err != nil {
        log.Fatalf("模板生成失败: %v", err)
    }
    
    fmt.Println("从模板生成的 XML:")
    fmt.Println(xmlContent)
    
    // 解析生成的 XML
    config, err := api.ParseFromString(xmlContent)
    if err != nil {
        log.Fatalf("解析模板 XML 失败: %v", err)
    }
    
    fmt.Printf("\n从模板解析了 %d 个包源\n", len(config.PackageSources.Add))
    
    // 与标准序列化比较
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
        return "", fmt.Errorf("模板解析失败: %w", err)
    }
    
    var buf strings.Builder
    err = tmpl.Execute(&buf, data)
    if err != nil {
        return "", fmt.Errorf("模板执行失败: %w", err)
    }
    
    return buf.String(), nil
}

func compareWithStandardSerialization(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== 与标准序列化比较 ===")
    
    // 使用标准方法序列化
    standardXML, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("标准序列化失败: %v", err)
        return
    }
    
    fmt.Println("标准序列化:")
    fmt.Println(standardXML)
    
    // 比较长度
    fmt.Printf("\n模板 XML 长度: %d 字符\n", len(xmlTemplate))
    fmt.Printf("标准 XML 长度: %d 字符\n", len(standardXML))
}
```

## 示例 4: 批量序列化

序列化多个配置：

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
    
    fmt.Println("=== 批量序列化 ===")
    
    // 为多个环境创建配置
    environments := []string{"development", "staging", "production"}
    
    for _, env := range environments {
        fmt.Printf("\n创建 %s 配置...\n", env)
        
        config := createEnvironmentConfig(api, env)
        
        // 序列化配置
        xmlContent, err := api.SerializeToXML(config)
        if err != nil {
            log.Printf("序列化 %s 配置失败: %v", env, err)
            continue
        }
        
        // 保存到环境特定文件
        filename := fmt.Sprintf("NuGet.%s.Config", env)
        err = os.WriteFile(filename, []byte(xmlContent), 0644)
        if err != nil {
            log.Printf("保存 %s 配置失败: %v", env, err)
            continue
        }
        
        fmt.Printf("✅ %s 配置已保存到 %s\n", env, filename)
        
        // 显示配置摘要
        displayConfigSummary(config, env)
    }
    
    // 创建包含所有环境的主配置
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
    fmt.Printf("  %s 配置摘要:\n", environment)
    fmt.Printf("    包源: %d\n", len(config.PackageSources.Add))
    
    if config.Config != nil {
        fmt.Printf("    配置选项: %d\n", len(config.Config.Add))
    }
    
    if config.PackageSourceCredentials != nil {
        fmt.Printf("    已认证源: %d\n", len(config.PackageSourceCredentials.Sources))
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("    活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
}

func createMasterConfiguration(api *nuget.API, environments []string) {
    fmt.Println("\n=== 创建主配置 ===")
    
    masterConfig := api.CreateDefaultConfig()
    
    // 添加所有可能的源
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
    
    // 默认禁用环境特定源
    api.DisablePackageSource(masterConfig, "local-dev")
    api.DisablePackageSource(masterConfig, "preview")
    api.DisablePackageSource(masterConfig, "company-staging")
    api.DisablePackageSource(masterConfig, "company-prod")
    
    // 添加通用配置
    api.AddConfigOption(masterConfig, "globalPackagesFolder", "${NUGET_PACKAGES}")
    api.AddConfigOption(masterConfig, "dependencyVersion", "HighestMinor")
    
    // 序列化主配置
    xmlContent, err := api.SerializeToXML(masterConfig)
    if err != nil {
        log.Printf("序列化主配置失败: %v", err)
        return
    }
    
    // 保存主配置
    masterFile := "NuGet.Master.Config"
    err = os.WriteFile(masterFile, []byte(xmlContent), 0644)
    if err != nil {
        log.Printf("保存主配置失败: %v", err)
        return
    }
    
    fmt.Printf("✅ 创建主配置: %s\n", masterFile)
    fmt.Printf("   包含 %v 环境的所有源\n", environments)
}
```

## 关键概念

### 序列化过程

1. **对象到 XML**: 将配置对象转换为 XML 字符串
2. **格式化**: 应用适当的缩进和结构
3. **验证**: 确保输出是有效的 XML
4. **往返**: 通过解析序列化输出进行验证

### XML 结构

序列化的 XML 遵循标准 NuGet.Config 格式：
- `<packageSources>`: 包源定义
- `<packageSourceCredentials>`: 认证信息
- `<config>`: 全局配置选项
- `<activePackageSource>`: 活跃源选择
- `<disabledPackageSources>`: 禁用的源

### 最佳实践

1. **验证输出**: 始终验证序列化的 XML 是否可解析
2. **一致格式**: 使用适当的缩进和结构
3. **处理编码**: 确保适当的 UTF-8 编码
4. **转义值**: 正确转义 XML 特殊字符
5. **测试往返**: 验证解析 → 序列化 → 解析循环

## 常见用例

### 配置导出
```go
// 导出当前配置
config, _, _ := api.FindAndParseConfig()
xmlContent, _ := api.SerializeToXML(config)
os.WriteFile("exported-config.xml", []byte(xmlContent), 0644)
```

### 模板生成
```go
// 从模板生成配置
templateData := ConfigTemplateData{...}
xmlContent, _ := generateXMLFromTemplate(templateData)
config, _ := api.ParseFromString(xmlContent)
```

### 批量处理
```go
// 处理多个配置
for _, env := range environments {
    config := createEnvironmentConfig(api, env)
    xmlContent, _ := api.SerializeToXML(config)
    saveToFile(fmt.Sprintf("%s.config", env), xmlContent)
}
```

## 下一步

掌握序列化后：

1. 学习 [位置感知编辑](./position-aware-editing.md) 进行精确修改
2. 探索 [解析器 API](/api/parser) 了解高级解析选项
3. 研究 [配置](../guide/configuration.md) 了解结构详情

本指南为将 NuGet 配置对象序列化为 XML 提供了全面的示例，涵盖了从基本输出到复杂模板生成的各种场景。
