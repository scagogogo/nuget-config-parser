# 包源管理

本示例演示如何使用 NuGet Config Parser 库管理 NuGet 包源，包括添加、移除、启用、禁用和配置包源。

## 概述

包源管理包括：
- 添加和移除包源
- 启用和禁用包源
- 设置活跃包源
- 管理包源优先级
- 配置包源认证
- 处理不同协议版本

## 示例 1: 基本包源操作

基本的包源添加、移除和管理：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 创建新配置或加载现有配置
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== 基本包源操作 ===")
    
    // 添加不同类型的包源
    fmt.Println("添加包源...")
    
    // 公共包源
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    fmt.Println("✅ 添加 nuget.org")
    
    // 公司内部源
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    fmt.Println("✅ 添加公司包源")
    
    // 本地文件夹源
    api.AddPackageSource(config, "local-packages", "/path/to/local/packages", "")
    fmt.Println("✅ 添加本地包源")
    
    // 网络共享源
    api.AddPackageSource(config, "network-share", "\\\\server\\packages", "")
    fmt.Println("✅ 添加网络共享源")
    
    // 显示所有包源
    fmt.Printf("\n当前包源 (%d):\n", len(config.PackageSources.Add))
    for i, source := range config.PackageSources.Add {
        protocol := source.ProtocolVersion
        if protocol == "" {
            protocol = "文件夹"
        } else {
            protocol = "v" + protocol
        }
        fmt.Printf("%d. %s (%s): %s\n", i+1, source.Key, protocol, source.Value)
    }
    
    // 设置活跃包源
    fmt.Println("\n设置活跃包源...")
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    if config.ActivePackageSource != nil {
        fmt.Printf("✅ 活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
    
    // 禁用某个包源
    fmt.Println("\n禁用包源...")
    api.DisablePackageSource(config, "network-share")
    fmt.Println("✅ 已禁用网络共享源")
    
    // 显示禁用的源
    if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
        fmt.Printf("\n禁用的源 (%d):\n", len(config.DisabledPackageSources.Add))
        for _, disabled := range config.DisabledPackageSources.Add {
            fmt.Printf("- %s\n", disabled.Key)
        }
    }
    
    // 保存配置
    err := api.SaveConfig(config, "PackageSourcesDemo.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n配置已保存到 PackageSourcesDemo.Config")
}
```

## 示例 2: 高级包源管理

管理复杂的包源场景：

```go
package main

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== 高级包源管理 ===")
    
    // 定义包源配置
    sources := []PackageSourceConfig{
        {
            Key:         "nuget.org",
            URL:         "https://api.nuget.org/v3/index.json",
            Version:     "3",
            Enabled:     true,
            Priority:    1,
            Description: "官方 NuGet 包源",
        },
        {
            Key:         "company-stable",
            URL:         "https://stable.company.com/nuget",
            Version:     "3",
            Enabled:     true,
            Priority:    2,
            Description: "公司稳定版包源",
            Username:    "company_user",
            Password:    "company_pass",
        },
        {
            Key:         "company-preview",
            URL:         "https://preview.company.com/nuget",
            Version:     "3",
            Enabled:     false, // 默认禁用预览版
            Priority:    3,
            Description: "公司预览版包源",
            Username:    "preview_user",
            Password:    "preview_pass",
        },
        {
            Key:         "local-dev",
            URL:         "./packages",
            Version:     "",
            Enabled:     true,
            Priority:    0, // 最高优先级
            Description: "本地开发包",
        },
        {
            Key:         "azure-artifacts",
            URL:         "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json",
            Version:     "3",
            Enabled:     true,
            Priority:    4,
            Description: "Azure DevOps 包源",
            Username:    "azure_user",
            Password:    "pat_token",
        },
    }
    
    // 按优先级排序并添加包源
    fmt.Println("按优先级添加包源...")
    
    // 先按优先级排序
    sortSourcesByPriority(sources)
    
    for _, sourceConfig := range sources {
        // 添加包源
        api.AddPackageSource(config, sourceConfig.Key, sourceConfig.URL, sourceConfig.Version)
        fmt.Printf("✅ 添加: %s (优先级: %d)\n", sourceConfig.Key, sourceConfig.Priority)
        
        // 添加凭证（如果需要）
        if sourceConfig.Username != "" && sourceConfig.Password != "" {
            api.AddCredential(config, sourceConfig.Key, sourceConfig.Username, sourceConfig.Password)
            fmt.Printf("   🔐 添加凭证: %s\n", sourceConfig.Key)
        }
        
        // 禁用源（如果需要）
        if !sourceConfig.Enabled {
            api.DisablePackageSource(config, sourceConfig.Key)
            fmt.Printf("   ❌ 禁用: %s\n", sourceConfig.Key)
        }
    }
    
    // 设置默认活跃源
    api.SetActivePackageSource(config, "local-dev", "./packages")
    fmt.Println("\n✅ 设置活跃源: local-dev")
    
    // 显示配置摘要
    displayAdvancedSummary(config, sources)
    
    // 演示动态源管理
    demonstrateDynamicManagement(api, config)
    
    // 保存配置
    err := api.SaveConfig(config, "AdvancedPackageSources.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n高级包源配置已保存")
}

type PackageSourceConfig struct {
    Key         string
    URL         string
    Version     string
    Enabled     bool
    Priority    int
    Description string
    Username    string
    Password    string
}

func sortSourcesByPriority(sources []PackageSourceConfig) {
    // 简单的冒泡排序，按优先级排序
    for i := 0; i < len(sources)-1; i++ {
        for j := 0; j < len(sources)-i-1; j++ {
            if sources[j].Priority > sources[j+1].Priority {
                sources[j], sources[j+1] = sources[j+1], sources[j]
            }
        }
    }
}

func displayAdvancedSummary(config *types.NuGetConfig, sources []PackageSourceConfig) {
    fmt.Println("\n=== 高级配置摘要 ===")
    
    // 显示包源及其状态
    fmt.Printf("包源配置 (%d):\n", len(config.PackageSources.Add))
    
    for _, source := range config.PackageSources.Add {
        // 查找原始配置信息
        var sourceConfig *PackageSourceConfig
        for _, sc := range sources {
            if sc.Key == source.Key {
                sourceConfig = &sc
                break
            }
        }
        
        status := "✅ 启用"
        if config.DisabledPackageSources != nil {
            for _, disabled := range config.DisabledPackageSources.Add {
                if disabled.Key == source.Key {
                    status = "❌ 禁用"
                    break
                }
            }
        }
        
        priority := "未知"
        description := "无描述"
        hasAuth := "无"
        
        if sourceConfig != nil {
            priority = fmt.Sprintf("%d", sourceConfig.Priority)
            description = sourceConfig.Description
            if sourceConfig.Username != "" {
                hasAuth = "有凭证"
            }
        }
        
        fmt.Printf("  %s [%s] (优先级: %s, 认证: %s)\n", source.Key, status, priority, hasAuth)
        fmt.Printf("    URL: %s\n", source.Value)
        fmt.Printf("    描述: %s\n", description)
        fmt.Println()
    }
    
    // 显示活跃源
    if config.ActivePackageSource != nil {
        fmt.Printf("🎯 活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
}

func demonstrateDynamicManagement(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== 动态源管理演示 ===")
    
    // 场景1: 切换到预览模式
    fmt.Println("场景1: 启用预览源...")
    
    // 启用预览源
    enablePackageSource(api, config, "company-preview")
    
    // 设置预览源为活跃源
    api.SetActivePackageSource(config, "company-preview", "https://preview.company.com/nuget")
    fmt.Println("✅ 已切换到预览模式")
    
    // 场景2: 临时禁用外部源
    fmt.Println("\n场景2: 禁用外部源（仅使用内部源）...")
    
    externalSources := []string{"nuget.org", "azure-artifacts"}
    for _, source := range externalSources {
        api.DisablePackageSource(config, source)
        fmt.Printf("❌ 禁用外部源: %s\n", source)
    }
    
    // 场景3: 添加临时源
    fmt.Println("\n场景3: 添加临时测试源...")
    
    api.AddPackageSource(config, "temp-test", "https://test.company.com/nuget", "3")
    api.AddCredential(config, "temp-test", "test_user", "test_pass")
    fmt.Println("✅ 添加临时测试源")
    
    // 场景4: 源健康检查模拟
    fmt.Println("\n场景4: 源健康检查...")
    
    healthCheck := map[string]bool{
        "nuget.org":        true,
        "company-stable":   true,
        "company-preview":  false, // 模拟不可用
        "local-dev":        true,
        "azure-artifacts":  true,
        "temp-test":        false, // 模拟不可用
    }
    
    for sourceName, isHealthy := range healthCheck {
        if isHealthy {
            fmt.Printf("✅ %s: 健康\n", sourceName)
        } else {
            fmt.Printf("❌ %s: 不可用，禁用中...\n", sourceName)
            api.DisablePackageSource(config, sourceName)
        }
    }
    
    fmt.Println("动态管理演示完成")
}

func enablePackageSource(api *nuget.API, config *types.NuGetConfig, sourceName string) {
    // 从禁用列表中移除
    if config.DisabledPackageSources != nil {
        var newDisabled []types.DisabledPackageSource
        for _, disabled := range config.DisabledPackageSources.Add {
            if disabled.Key != sourceName {
                newDisabled = append(newDisabled, disabled)
            }
        }
        config.DisabledPackageSources.Add = newDisabled
    }
}
```

## 示例 3: 企业级包源配置

企业环境中的复杂包源管理：

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
    
    fmt.Println("=== 企业级包源配置 ===")
    
    // 创建不同环境的配置
    environments := []string{"development", "staging", "production"}
    
    for _, env := range environments {
        fmt.Printf("\n配置 %s 环境...\n", env)
        
        config := createEnterpriseConfig(api, env)
        
        // 保存环境特定配置
        configPath := fmt.Sprintf("Enterprise.%s.Config", env)
        err := api.SaveConfig(config, configPath)
        if err != nil {
            log.Printf("保存 %s 配置失败: %v", env, err)
            continue
        }
        
        fmt.Printf("✅ %s 环境配置已保存到 %s\n", env, configPath)
        
        // 显示环境配置摘要
        displayEnvironmentSummary(config, env)
    }
    
    // 创建主配置文件
    createMasterConfig(api)
    
    fmt.Println("\n企业级包源配置完成")
}

func createEnterpriseConfig(api *nuget.API, environment string) *types.NuGetConfig {
    config := api.CreateDefaultConfig()
    
    // 基础包源（所有环境共有）
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    
    // 企业内部源
    api.AddPackageSource(config, "enterprise-stable", "https://nuget.enterprise.com/stable", "3")
    api.AddCredential(config, "enterprise-stable", "enterprise_user", "enterprise_pass")
    
    // 根据环境添加特定源
    switch environment {
    case "development":
        // 开发环境源
        api.AddPackageSource(config, "dev-internal", "https://dev.nuget.enterprise.com", "3")
        api.AddPackageSource(config, "local-builds", "./local-packages", "")
        api.AddPackageSource(config, "preview-feed", "https://preview.nuget.enterprise.com", "3")
        
        // 开发环境凭证
        api.AddCredential(config, "dev-internal", "dev_user", "dev_pass")
        api.AddCredential(config, "preview-feed", "preview_user", "preview_pass")
        
        // 开发环境配置
        api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
        api.AddConfigOption(config, "dependencyVersion", "Highest")
        api.AddConfigOption(config, "allowPrereleaseVersions", "true")
        
        // 设置本地构建为活跃源
        api.SetActivePackageSource(config, "local-builds", "./local-packages")
        
    case "staging":
        // 预发布环境源
        api.AddPackageSource(config, "staging-internal", "https://staging.nuget.enterprise.com", "3")
        api.AddPackageSource(config, "integration-test", "https://test.nuget.enterprise.com", "3")
        
        // 预发布环境凭证
        api.AddCredential(config, "staging-internal", "staging_user", "staging_pass")
        api.AddCredential(config, "integration-test", "test_user", "test_pass")
        
        // 预发布环境配置
        api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
        api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
        api.AddConfigOption(config, "allowPrereleaseVersions", "false")
        
        // 禁用预览源
        api.DisablePackageSource(config, "preview-feed")
        
        // 设置企业稳定源为活跃源
        api.SetActivePackageSource(config, "enterprise-stable", "https://nuget.enterprise.com/stable")
        
    case "production":
        // 生产环境源（最严格）
        api.AddPackageSource(config, "production-approved", "https://prod.nuget.enterprise.com", "3")
        
        // 生产环境凭证
        api.AddCredential(config, "production-approved", "prod_user", "prod_pass")
        
        // 生产环境配置
        api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
        api.AddConfigOption(config, "dependencyVersion", "Exact")
        api.AddConfigOption(config, "allowPrereleaseVersions", "false")
        api.AddConfigOption(config, "signatureValidationMode", "require")
        
        // 禁用所有非生产源
        api.DisablePackageSource(config, "nuget.org") // 生产环境可能不允许外部源
        
        // 设置生产批准源为活跃源
        api.SetActivePackageSource(config, "production-approved", "https://prod.nuget.enterprise.com")
    }
    
    return config
}

func displayEnvironmentSummary(config *types.NuGetConfig, environment string) {
    fmt.Printf("\n--- %s 环境摘要 ---\n", strings.ToUpper(environment))
    
    // 显示包源
    fmt.Printf("包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "启用"
        if config.DisabledPackageSources != nil {
            for _, disabled := range config.DisabledPackageSources.Add {
                if disabled.Key == source.Key {
                    status = "禁用"
                    break
                }
            }
        }
        fmt.Printf("  - %s [%s]: %s\n", source.Key, status, source.Value)
    }
    
    // 显示活跃源
    if config.ActivePackageSource != nil {
        fmt.Printf("活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
    
    // 显示配置选项
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("配置选项 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    // 显示认证源
    if config.PackageSourceCredentials != nil {
        fmt.Printf("已配置认证 (%d):\n", len(config.PackageSourceCredentials.Sources))
        for sourceName := range config.PackageSourceCredentials.Sources {
            fmt.Printf("  - %s\n", sourceName)
        }
    }
}

func createMasterConfig(api *nuget.API) {
    fmt.Println("\n创建主配置文件...")
    
    config := api.CreateDefaultConfig()
    
    // 添加所有可能的源（大部分禁用）
    allSources := map[string]struct {
        url     string
        version string
        enabled bool
    }{
        "nuget.org":             {"https://api.nuget.org/v3/index.json", "3", true},
        "enterprise-stable":     {"https://nuget.enterprise.com/stable", "3", true},
        "dev-internal":          {"https://dev.nuget.enterprise.com", "3", false},
        "staging-internal":      {"https://staging.nuget.enterprise.com", "3", false},
        "production-approved":   {"https://prod.nuget.enterprise.com", "3", false},
        "preview-feed":          {"https://preview.nuget.enterprise.com", "3", false},
        "integration-test":      {"https://test.nuget.enterprise.com", "3", false},
        "local-builds":          {"./local-packages", "", false},
    }
    
    for sourceName, sourceInfo := range allSources {
        api.AddPackageSource(config, sourceName, sourceInfo.url, sourceInfo.version)
        
        if !sourceInfo.enabled {
            api.DisablePackageSource(config, sourceName)
        }
    }
    
    // 添加通用配置
    api.AddConfigOption(config, "globalPackagesFolder", "${NUGET_PACKAGES}")
    api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    // 设置默认活跃源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // 保存主配置
    err := api.SaveConfig(config, "Enterprise.Master.Config")
    if err != nil {
        log.Printf("保存主配置失败: %v", err)
        return
    }
    
    fmt.Println("✅ 主配置文件已创建: Enterprise.Master.Config")
    fmt.Println("   包含所有环境的源定义，可根据需要启用/禁用")
}
```

## 关键概念

### 包源类型

1. **HTTP/HTTPS源** - 远程NuGet服务器
2. **本地文件夹** - 本地目录中的包
3. **网络共享** - UNC路径的包
4. **Azure Artifacts** - Azure DevOps包源

### 源管理操作

1. **添加源** - `AddPackageSource()`
2. **移除源** - 从配置中删除
3. **启用/禁用** - `DisablePackageSource()`
4. **设置活跃源** - `SetActivePackageSource()`

### 最佳实践

1. **优先级管理** - 按重要性排序源
2. **环境隔离** - 不同环境使用不同源
3. **安全认证** - 为私有源配置凭证
4. **健康监控** - 定期检查源可用性

## 下一步

掌握包源管理后：

1. 学习 [凭证管理](./credentials.md) 来处理认证
2. 探索 [配置选项](./config-options.md) 进行高级设置
3. 研究 [位置感知编辑](./position-aware-editing.md) 进行精确修改

本指南为 NuGet 包源管理提供了全面的示例，涵盖了从基本操作到企业级配置的各种场景。
