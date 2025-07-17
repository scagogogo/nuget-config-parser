# 修改配置

本示例演示如何使用 NuGet Config Parser 库修改现有的 NuGet 配置文件。

## 概述

配置修改包括：
- 更新现有包源的URL和设置
- 添加新的包源和凭证
- 移除过时的配置项
- 批量更新多个配置
- 条件性修改基于环境

## 示例 1: 基本配置修改

修改现有配置文件的最简单方法：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 加载现有配置
    configPath := "NuGet.Config"
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    fmt.Printf("修改前包源数量: %d\n", len(config.PackageSources.Add))
    
    // 添加新的包源
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    fmt.Println("已添加公司包源")
    
    // 更新现有包源的URL
    for i, source := range config.PackageSources.Add {
        if source.Key == "nuget.org" {
            config.PackageSources.Add[i].Value = "https://api.nuget.org/v3/index.json"
            fmt.Println("已更新 nuget.org URL")
            break
        }
    }
    
    // 添加配置选项
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    api.AddConfigOption(config, "defaultPushSource", "https://nuget.company.com/v3/index.json")
    
    fmt.Printf("修改后包源数量: %d\n", len(config.PackageSources.Add))
    
    // 保存修改后的配置
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("配置修改完成并已保存")
}
```

## 示例 2: 批量配置更新

对多个配置项进行批量修改：

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
    
    // 加载配置
    config, err := api.ParseFromFile("NuGet.Config")
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    fmt.Println("=== 批量配置更新 ===")
    
    // 定义批量更新操作
    updates := []struct {
        operation string
        params    map[string]string
    }{
        {"add_source", map[string]string{"key": "staging", "url": "https://staging.company.com/nuget", "version": "3"}},
        {"add_source", map[string]string{"key": "production", "url": "https://prod.company.com/nuget", "version": "3"}},
        {"add_source", map[string]string{"key": "local-dev", "url": "./packages", "version": ""}},
        {"add_credential", map[string]string{"source": "staging", "username": "dev_user", "password": "dev_pass"}},
        {"add_credential", map[string]string{"source": "production", "username": "prod_user", "password": "prod_pass"}},
        {"add_config", map[string]string{"key": "dependencyVersion", "value": "Highest"}},
        {"add_config", map[string]string{"key": "automaticPackageRestore", "value": "true"}},
    }
    
    // 执行批量更新
    successCount := 0
    for i, update := range updates {
        fmt.Printf("执行操作 %d: %s\n", i+1, update.operation)
        
        var err error
        switch update.operation {
        case "add_source":
            api.AddPackageSource(config, update.params["key"], update.params["url"], update.params["version"])
            fmt.Printf("  ✅ 添加包源: %s\n", update.params["key"])
            
        case "add_credential":
            api.AddCredential(config, update.params["source"], update.params["username"], update.params["password"])
            fmt.Printf("  ✅ 添加凭证: %s\n", update.params["source"])
            
        case "add_config":
            api.AddConfigOption(config, update.params["key"], update.params["value"])
            fmt.Printf("  ✅ 添加配置: %s = %s\n", update.params["key"], update.params["value"])
            
        default:
            err = fmt.Errorf("未知操作: %s", update.operation)
        }
        
        if err != nil {
            fmt.Printf("  ❌ 操作失败: %v\n", err)
        } else {
            successCount++
        }
    }
    
    fmt.Printf("\n批量更新完成: %d/%d 操作成功\n", successCount, len(updates))
    
    // 设置活跃包源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // 禁用开发源（默认情况下）
    api.DisablePackageSource(config, "local-dev")
    
    // 显示最终配置摘要
    displayConfigSummary(config)
    
    // 保存配置
    err = api.SaveConfig(config, "UpdatedNuGet.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n批量更新的配置已保存到 UpdatedNuGet.Config")
}

func displayConfigSummary(config *types.NuGetConfig) {
    fmt.Println("\n=== 配置摘要 ===")
    
    fmt.Printf("包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
    
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("\n配置选项 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    if config.PackageSourceCredentials != nil && len(config.PackageSourceCredentials.Sources) > 0 {
        fmt.Printf("\n已配置凭证的源 (%d):\n", len(config.PackageSourceCredentials.Sources))
        for sourceName := range config.PackageSourceCredentials.Sources {
            fmt.Printf("  - %s\n", sourceName)
        }
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("\n活跃包源: %s\n", config.ActivePackageSource.Add.Key)
    }
}
```

## 示例 3: 条件性配置修改

基于环境或条件修改配置：

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
    
    // 获取环境信息
    environment := getEnvironment()
    fmt.Printf("当前环境: %s\n", environment)
    
    // 加载基础配置
    config, err := api.ParseFromFile("NuGet.Config")
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    fmt.Println("=== 条件性配置修改 ===")
    
    // 根据环境应用不同的修改
    switch environment {
    case "development":
        applyDevelopmentConfig(api, config)
    case "staging":
        applyStagingConfig(api, config)
    case "production":
        applyProductionConfig(api, config)
    default:
        applyDefaultConfig(api, config)
    }
    
    // 应用通用修改
    applyCommonConfig(api, config)
    
    // 清理过时配置
    cleanupObsoleteConfig(api, config)
    
    // 保存环境特定配置
    outputPath := fmt.Sprintf("NuGet.%s.Config", environment)
    err = api.SaveConfig(config, outputPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("环境特定配置已保存到: %s\n", outputPath)
    
    // 验证配置
    validateConfiguration(api, config, environment)
}

func getEnvironment() string {
    env := os.Getenv("ENVIRONMENT")
    if env == "" {
        env = os.Getenv("ASPNETCORE_ENVIRONMENT")
    }
    if env == "" {
        env = "development"
    }
    return strings.ToLower(env)
}

func applyDevelopmentConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("应用开发环境配置...")
    
    // 添加本地开发源
    api.AddPackageSource(config, "local-packages", "./packages", "")
    api.AddPackageSource(config, "dev-feed", "https://dev.company.com/nuget", "3")
    
    // 开发环境配置
    api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "allowPrereleaseVersions", "true")
    
    // 设置本地源为活跃源
    api.SetActivePackageSource(config, "local-packages", "./packages")
    
    fmt.Println("  ✅ 开发环境配置已应用")
}

func applyStagingConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("应用预发布环境配置...")
    
    // 添加预发布源
    api.AddPackageSource(config, "staging-feed", "https://staging.company.com/nuget", "3")
    
    // 预发布环境配置
    api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
    api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    
    // 添加预发布凭证
    api.AddCredential(config, "staging-feed", "staging_user", "staging_pass")
    
    // 设置预发布源为活跃源
    api.SetActivePackageSource(config, "staging-feed", "https://staging.company.com/nuget")
    
    fmt.Println("  ✅ 预发布环境配置已应用")
}

func applyProductionConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("应用生产环境配置...")
    
    // 添加生产源
    api.AddPackageSource(config, "production-feed", "https://prod.company.com/nuget", "3")
    
    // 生产环境配置
    api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
    api.AddConfigOption(config, "dependencyVersion", "Exact")
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    api.AddConfigOption(config, "signatureValidationMode", "require")
    
    // 添加生产凭证
    api.AddCredential(config, "production-feed", "prod_user", "prod_pass")
    
    // 设置生产源为活跃源
    api.SetActivePackageSource(config, "production-feed", "https://prod.company.com/nuget")
    
    fmt.Println("  ✅ 生产环境配置已应用")
}

func applyDefaultConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("应用默认配置...")
    
    // 默认配置
    api.AddConfigOption(config, "globalPackagesFolder", "~/.nuget/packages")
    api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    fmt.Println("  ✅ 默认配置已应用")
}

func applyCommonConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("应用通用配置...")
    
    // 确保 nuget.org 存在且是最新的
    found := false
    for i, source := range config.PackageSources.Add {
        if source.Key == "nuget.org" {
            config.PackageSources.Add[i].Value = "https://api.nuget.org/v3/index.json"
            config.PackageSources.Add[i].ProtocolVersion = "3"
            found = true
            break
        }
    }
    
    if !found {
        api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    }
    
    // 通用配置选项
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    fmt.Println("  ✅ 通用配置已应用")
}

func cleanupObsoleteConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("清理过时配置...")
    
    // 移除过时的包源
    obsoleteSources := []string{"old-nuget", "legacy-feed", "deprecated-source"}
    
    for _, obsoleteSource := range obsoleteSources {
        // 检查是否存在过时源
        for i, source := range config.PackageSources.Add {
            if source.Key == obsoleteSource {
                // 移除过时源
                config.PackageSources.Add = append(
                    config.PackageSources.Add[:i],
                    config.PackageSources.Add[i+1:]...)
                fmt.Printf("  ✅ 移除过时源: %s\n", obsoleteSource)
                break
            }
        }
    }
    
    // 更新过时的URL
    urlMigrations := map[string]string{
        "https://www.nuget.org/api/v2": "https://api.nuget.org/v3/index.json",
        "http://nuget.org/api/v2":     "https://api.nuget.org/v3/index.json",
    }
    
    for i, source := range config.PackageSources.Add {
        if newURL, exists := urlMigrations[source.Value]; exists {
            config.PackageSources.Add[i].Value = newURL
            config.PackageSources.Add[i].ProtocolVersion = "3"
            fmt.Printf("  ✅ 更新URL: %s -> %s\n", source.Key, newURL)
        }
    }
    
    fmt.Println("  ✅ 过时配置清理完成")
}

func validateConfiguration(api *nuget.API, config *types.NuGetConfig, environment string) {
    fmt.Println("\n=== 配置验证 ===")
    
    // 基本验证
    if len(config.PackageSources.Add) == 0 {
        fmt.Println("  ❌ 警告: 没有配置包源")
        return
    }
    
    fmt.Printf("  ✅ 包源数量: %d\n", len(config.PackageSources.Add))
    
    // 验证必需的源
    requiredSources := map[string][]string{
        "development": {"nuget.org", "local-packages"},
        "staging":     {"nuget.org", "staging-feed"},
        "production":  {"nuget.org", "production-feed"},
    }
    
    if required, exists := requiredSources[environment]; exists {
        for _, requiredSource := range required {
            found := false
            for _, source := range config.PackageSources.Add {
                if source.Key == requiredSource {
                    found = true
                    break
                }
            }
            
            if found {
                fmt.Printf("  ✅ 必需源存在: %s\n", requiredSource)
            } else {
                fmt.Printf("  ❌ 缺少必需源: %s\n", requiredSource)
            }
        }
    }
    
    // 验证活跃源
    if config.ActivePackageSource != nil {
        fmt.Printf("  ✅ 活跃源: %s\n", config.ActivePackageSource.Add.Key)
    } else {
        fmt.Println("  ⚠️  未设置活跃源")
    }
    
    fmt.Println("配置验证完成")
}
```

## 示例 4: 智能配置合并

合并多个配置文件的设置：

```go
package main

import (
    "fmt"
    "log"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Println("=== 智能配置合并 ===")
    
    // 要合并的配置文件
    configFiles := []string{
        "base.config",
        "team.config", 
        "project.config",
    }
    
    // 创建基础配置
    mergedConfig := api.CreateDefaultConfig()
    
    // 逐个合并配置文件
    for i, configFile := range configFiles {
        fmt.Printf("合并配置文件 %d: %s\n", i+1, configFile)
        
        if err := mergeConfigFile(api, mergedConfig, configFile); err != nil {
            log.Printf("合并 %s 失败: %v", configFile, err)
            continue
        }
        
        fmt.Printf("  ✅ %s 合并成功\n", configFile)
    }
    
    // 解决冲突和重复
    resolveConflicts(api, mergedConfig)
    
    // 优化配置
    optimizeConfig(api, mergedConfig)
    
    // 保存合并后的配置
    outputPath := "merged.config"
    err := api.SaveConfig(mergedConfig, outputPath)
    if err != nil {
        log.Fatalf("保存合并配置失败: %v", err)
    }
    
    fmt.Printf("合并后的配置已保存到: %s\n", outputPath)
    
    // 显示合并结果
    displayMergeResults(mergedConfig)
}

func mergeConfigFile(api *nuget.API, target *types.NuGetConfig, configPath string) error {
    // 检查文件是否存在
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return fmt.Errorf("配置文件不存在: %s", configPath)
    }
    
    // 解析源配置
    sourceConfig, err := api.ParseFromFile(configPath)
    if err != nil {
        return fmt.Errorf("解析配置文件失败: %w", err)
    }
    
    // 合并包源
    for _, source := range sourceConfig.PackageSources.Add {
        // 检查是否已存在
        exists := false
        for _, existing := range target.PackageSources.Add {
            if existing.Key == source.Key {
                exists = true
                break
            }
        }
        
        if !exists {
            api.AddPackageSource(target, source.Key, source.Value, source.ProtocolVersion)
        }
    }
    
    // 合并凭证
    if sourceConfig.PackageSourceCredentials != nil {
        for sourceName, creds := range sourceConfig.PackageSourceCredentials.Sources {
            if creds.Username != nil && creds.Password != nil {
                api.AddCredential(target, sourceName, creds.Username.Value, creds.Password.Value)
            }
        }
    }
    
    // 合并配置选项
    if sourceConfig.Config != nil {
        for _, option := range sourceConfig.Config.Add {
            api.AddConfigOption(target, option.Key, option.Value)
        }
    }
    
    return nil
}

func resolveConflicts(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("解决配置冲突...")
    
    // 移除重复的包源
    seen := make(map[string]bool)
    var uniqueSources []types.PackageSource
    
    for _, source := range config.PackageSources.Add {
        key := fmt.Sprintf("%s|%s", source.Key, source.Value)
        if !seen[key] {
            seen[key] = true
            uniqueSources = append(uniqueSources, source)
        } else {
            fmt.Printf("  移除重复源: %s\n", source.Key)
        }
    }
    
    config.PackageSources.Add = uniqueSources
    
    // 解决配置选项冲突（保留最后一个值）
    if config.Config != nil {
        optionMap := make(map[string]string)
        for _, option := range config.Config.Add {
            optionMap[option.Key] = option.Value
        }
        
        // 重建配置选项列表
        config.Config.Add = nil
        for key, value := range optionMap {
            api.AddConfigOption(config, key, value)
        }
    }
    
    fmt.Println("  ✅ 冲突解决完成")
}

func optimizeConfig(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("优化配置...")
    
    // 确保 nuget.org 使用最新URL
    for i, source := range config.PackageSources.Add {
        if source.Key == "nuget.org" {
            config.PackageSources.Add[i].Value = "https://api.nuget.org/v3/index.json"
            config.PackageSources.Add[i].ProtocolVersion = "3"
            fmt.Println("  ✅ 优化 nuget.org URL")
            break
        }
    }
    
    // 设置默认配置选项
    defaultOptions := map[string]string{
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "dependencyVersion":       "HighestMinor",
    }
    
    for key, value := range defaultOptions {
        // 检查是否已存在
        exists := false
        if config.Config != nil {
            for _, option := range config.Config.Add {
                if option.Key == key {
                    exists = true
                    break
                }
            }
        }
        
        if !exists {
            api.AddConfigOption(config, key, value)
            fmt.Printf("  ✅ 添加默认选项: %s = %s\n", key, value)
        }
    }
    
    fmt.Println("  ✅ 配置优化完成")
}

func displayMergeResults(config *types.NuGetConfig) {
    fmt.Println("\n=== 合并结果 ===")
    
    fmt.Printf("总包源数: %d\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
    
    if config.Config != nil {
        fmt.Printf("\n配置选项数: %d\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
    }
    
    if config.PackageSourceCredentials != nil {
        fmt.Printf("\n已配置凭证的源数: %d\n", len(config.PackageSourceCredentials.Sources))
        for sourceName := range config.PackageSourceCredentials.Sources {
            fmt.Printf("  - %s\n", sourceName)
        }
    }
}
```

## 关键概念

### 修改策略

1. **增量修改**: 只更改需要的部分
2. **批量操作**: 一次性应用多个更改
3. **条件修改**: 基于环境或条件的修改
4. **冲突解决**: 处理重复和冲突的配置

### 最佳实践

1. **备份原始配置**: 修改前创建备份
2. **验证修改**: 确保修改后的配置有效
3. **渐进式修改**: 逐步应用复杂修改
4. **环境隔离**: 为不同环境维护不同配置

## 下一步

掌握配置修改后：

1. 学习 [包源管理](./package-sources.md) 进行高级源操作
2. 探索 [位置感知编辑](./position-aware-editing.md) 进行精确修改
3. 研究 [序列化](./serialization.md) 了解配置输出格式

本指南为修改 NuGet 配置文件提供了全面的示例，涵盖了从简单更新到复杂批量操作的各种场景。
