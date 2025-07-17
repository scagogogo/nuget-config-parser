# 创建配置

本示例演示如何使用 NuGet Config Parser 库从头开始创建新的 NuGet 配置文件。

## 概述

创建配置包括：
- 以编程方式构建配置对象
- 设置默认包源
- 配置认证和凭证
- 将配置保存到文件
- 初始化项目特定设置

## 示例 1: 创建基本配置

创建新配置的最简单方法：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 创建默认配置
    config := api.CreateDefaultConfig()
    
    fmt.Println("创建了默认配置:")
    fmt.Printf("包源: %d\n", len(config.PackageSources.Add))
    
    // 显示默认源
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf("  协议: v%s\n", source.ProtocolVersion)
        }
    }
    
    // 保存到文件
    configPath := "NuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("\n配置已保存到: %s\n", configPath)
}
```

## 示例 2: 创建自定义配置

构建具有自定义包源的配置：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
    api := nuget.NewAPI()
    
    // 创建空配置
    config := &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{},
        },
    }
    
    // 添加多个包源
    sources := []struct {
        key     string
        value   string
        version string
    }{
        {"nuget.org", "https://api.nuget.org/v3/index.json", "3"},
        {"company-feed", "https://nuget.company.com/v3/index.json", "3"},
        {"local-dev", "/path/to/local/packages", ""},
        {"backup-feed", "https://backup.nuget.com/api/v2", "2"},
    }
    
    fmt.Println("创建包含多个源的自定义配置:")
    
    for _, source := range sources {
        api.AddPackageSource(config, source.key, source.value, source.version)
        fmt.Printf("已添加: %s -> %s\n", source.key, source.value)
    }
    
    // 设置活跃包源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // 添加一些全局配置选项
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "defaultPushSource", "https://nuget.company.com/v3/index.json")
    
    // 默认禁用本地开发源
    api.DisablePackageSource(config, "local-dev")
    
    // 保存配置
    configPath := "CustomNuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("\n自定义配置已保存到: %s\n", configPath)
    
    // 显示最终配置
    displayConfiguration(config)
}

func displayConfiguration(config *types.NuGetConfig) {
    fmt.Println("\n=== 最终配置 ===")
    
    fmt.Printf("包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("\n活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
    
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("\n配置选项 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
        fmt.Printf("\n禁用的源 (%d):\n", len(config.DisabledPackageSources.Add))
        for _, disabled := range config.DisabledPackageSources.Add {
            fmt.Printf("  - %s\n", disabled.Key)
        }
    }
}
```

## 示例 3: 创建带凭证的配置

构建包含认证的配置：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 从默认配置开始
    config := api.CreateDefaultConfig()
    
    // 添加需要认证的私有包源
    privateSources := []struct {
        key      string
        url      string
        username string
        password string
    }{
        {"company-internal", "https://internal.company.com/nuget", "employee", "secret123"},
        {"partner-feed", "https://partner.company.com/nuget", "partner_user", "partner_pass"},
        {"azure-artifacts", "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json", "myuser", "pat_token"},
    }
    
    fmt.Println("创建带认证源的配置:")
    
    for _, source := range privateSources {
        // 添加包源
        api.AddPackageSource(config, source.key, source.url, "3")
        fmt.Printf("已添加源: %s\n", source.key)
        
        // 为源添加凭证
        api.AddCredential(config, source.key, source.username, source.password)
        fmt.Printf("已为 %s 添加凭证\n", source.key)
    }
    
    // 添加不需要凭证的公共源
    api.AddPackageSource(config, "public-feed", "https://public.nuget.com/v3/index.json", "3")
    
    // 保存配置
    configPath := "AuthenticatedNuGet.Config"
    err := api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("\n带凭证的配置已保存到: %s\n", configPath)
    
    // 验证凭证已添加
    fmt.Println("\n凭证验证:")
    for _, source := range privateSources {
        credential := api.GetCredential(config, source.key)
        if credential != nil {
            fmt.Printf("✅ %s 已配置凭证\n", source.key)
        } else {
            fmt.Printf("❌ %s 缺少凭证\n", source.key)
        }
    }
    
    // 显示 XML 输出（生产环境中要小心密码！）
    fmt.Println("\n生成的 XML 预览:")
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("序列化失败: %v", err)
    } else {
        // 在生产环境中，您需要屏蔽密码
        fmt.Println(xmlContent)
    }
}
```

## 示例 4: 创建项目特定配置

创建针对特定项目定制的配置：

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
    
    // 获取项目信息
    projectDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("获取当前目录失败: %v", err)
    }
    
    projectName := filepath.Base(projectDir)
    
    fmt.Printf("为项目创建特定配置: %s\n", projectName)
    fmt.Printf("项目目录: %s\n", projectDir)
    
    // 创建针对此项目优化的配置
    config := api.CreateDefaultConfig()
    
    // 添加项目特定的包源
    api.AddPackageSource(config, "project-local", "./packages", "")
    api.AddPackageSource(config, "project-cache", filepath.Join(projectDir, ".nuget", "cache"), "")
    
    // 配置项目特定设置
    packagesPath := filepath.Join(projectDir, "packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "globalPackagesFolder", packagesPath)
    
    // 设置开发友好的设置
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    // 添加常见的开发源
    api.AddPackageSource(config, "nuget-preview", "https://api.nuget.org/v3-flatcontainer", "3")
    api.AddPackageSource(config, "dotnet-core", "https://dotnetfeed.blob.core.windows.net/dotnet-core/index.json", "3")
    
    // 默认禁用预览源
    api.DisablePackageSource(config, "nuget-preview")
    api.DisablePackageSource(config, "dotnet-core")
    
    // 设置 nuget.org 为活跃源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // 如果不存在则创建包目录
    if err := os.MkdirAll(packagesPath, 0755); err != nil {
        log.Printf("警告: 创建包目录失败: %v", err)
    }
    
    // 在项目根目录保存配置
    configPath := filepath.Join(projectDir, "NuGet.Config")
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存项目配置失败: %v", err)
    }
    
    fmt.Printf("\n项目配置已创建: %s\n", configPath)
    
    // 为包创建 .gitignore 条目（如果 .gitignore 存在）
    gitignorePath := filepath.Join(projectDir, ".gitignore")
    if _, err := os.Stat(gitignorePath); err == nil {
        addToGitignore(gitignorePath, "packages/")
        fmt.Println("已将 packages/ 添加到 .gitignore")
    }
    
    // 显示项目配置摘要
    displayProjectSummary(config, projectName, configPath)
}

func addToGitignore(gitignorePath, entry string) {
    // 读取现有 .gitignore
    content, err := os.ReadFile(gitignorePath)
    if err != nil {
        return
    }
    
    // 检查条目是否已存在
    if strings.Contains(string(content), entry) {
        return
    }
    
    // 追加条目
    file, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return
    }
    defer file.Close()
    
    file.WriteString(fmt.Sprintf("\n# NuGet 包\n%s\n", entry))
}

func displayProjectSummary(config *types.NuGetConfig, projectName, configPath string) {
    fmt.Printf("\n=== 项目配置摘要 ===\n")
    fmt.Printf("项目: %s\n", projectName)
    fmt.Printf("配置文件: %s\n", configPath)
    
    fmt.Printf("\n包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "已启用"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "已禁用"
        }
        fmt.Printf("  - %s (%s): %s\n", source.Key, status, source.Value)
    }
    
    if config.Config != nil {
        fmt.Printf("\n项目设置 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    fmt.Println("\n下一步:")
    fmt.Println("1. 根据项目需要自定义包源")
    fmt.Println("2. 如需要，为私有源添加凭证")
    fmt.Println("3. 将 NuGet.Config 提交到版本控制")
    fmt.Println("4. 与团队成员共享配置")
}
```

## 关键概念

### 配置结构

NuGet 配置包含：
- **包源**: 查找包的位置
- **凭证**: 私有源的认证
- **配置选项**: 全局设置和首选项
- **活跃源**: 当前选择的源
- **禁用源**: 临时禁用的源

### 最佳实践

1. **从默认开始**: 使用 `CreateDefaultConfig()` 作为基础
2. **逐步添加源**: 逐步构建配置
3. **安全处理凭证**: 小心密码存储
4. **设置适当权限**: 确保配置文件有正确权限
5. **保存前验证**: 检查配置有效性

## 常见模式

### 模式 1: 增量构建

```go
config := api.CreateDefaultConfig()
api.AddPackageSource(config, "source1", "url1", "3")
api.AddPackageSource(config, "source2", "url2", "3")
api.AddCredential(config, "source1", "user", "pass")
api.SaveConfig(config, "NuGet.Config")
```

### 模式 2: 基于模板创建

```go
func createEnterpriseConfig(api *nuget.API, companyDomain string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // 添加公司特定源
    api.AddPackageSource(config, "company", fmt.Sprintf("https://nuget.%s", companyDomain), "3")
    api.AddPackageSource(config, "company-preview", fmt.Sprintf("https://preview.nuget.%s", companyDomain), "3")
    
    // 配置企业设置
    api.AddConfigOption(config, "defaultPushSource", fmt.Sprintf("https://nuget.%s", companyDomain))
    
    return config
}
```

### 模式 3: 基于环境的配置

```go
func createConfigForEnvironment(api *nuget.API, env string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    switch env {
    case "development":
        api.AddPackageSource(config, "local", "./packages", "")
        api.AddPackageSource(config, "dev-feed", "https://dev.nuget.com", "3")
    case "staging":
        api.AddPackageSource(config, "staging-feed", "https://staging.nuget.com", "3")
    case "production":
        api.AddPackageSource(config, "prod-feed", "https://prod.nuget.com", "3")
    }
    
    return config
}
```

## 下一步

掌握配置创建后：

1. 学习 [修改配置](./modifying-configs.md) 来更新现有配置
2. 探索 [包源](./package-sources.md) 进行高级源管理
3. 研究 [凭证](./credentials.md) 进行安全认证处理

本指南为从头创建 NuGet 配置文件提供了全面的示例，涵盖了从简单更新到复杂批量操作的各种场景。
