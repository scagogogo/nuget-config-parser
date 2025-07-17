# 配置选项

本示例演示如何使用 NuGet Config Parser 库管理全局 NuGet 配置选项。

## 概述

配置选项控制 NuGet 行为的各个方面：
- 包存储位置
- 依赖解析策略
- 代理设置
- 包还原行为
- 默认推送源

## 示例 1: 基本配置选项

管理基础配置选项：

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
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== 基本配置选项 ===")
    
    // 设置基本包位置
    homeDir := os.Getenv("HOME")
    if homeDir == "" {
        homeDir = os.Getenv("USERPROFILE") // Windows
    }
    
    // 配置包存储位置
    globalPackagesPath := filepath.Join(homeDir, ".nuget", "packages")
    repositoryPath := "./packages"
    
    api.AddConfigOption(config, "globalPackagesFolder", globalPackagesPath)
    api.AddConfigOption(config, "repositoryPath", repositoryPath)
    
    fmt.Printf("全局包文件夹: %s\n", globalPackagesPath)
    fmt.Printf("仓库路径: %s\n", repositoryPath)
    
    // 配置依赖解析
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "packageRestore", "true")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    fmt.Println("依赖解析: Highest")
    fmt.Println("包还原: 已启用")
    
    // 配置默认推送源
    api.AddConfigOption(config, "defaultPushSource", "https://api.nuget.org/v3/index.json")
    fmt.Println("默认推送源: nuget.org")
    
    // 显示所有配置选项
    fmt.Println("\n=== 所有配置选项 ===")
    if config.Config != nil {
        for _, option := range config.Config.Add {
            fmt.Printf("  %s: %s\n", option.Key, option.Value)
        }
    }
    
    // 保存配置
    err := api.SaveConfig(config, "BasicOptions.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n配置保存成功！")
}
```

## 示例 2: 代理配置

为企业环境设置代理设置：

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
    
    fmt.Println("=== 代理配置 ===")
    
    // 从环境变量或使用默认值获取代理设置
    proxyURL := getEnvOrDefault("HTTP_PROXY", "http://proxy.company.com:8080")
    proxyUser := getEnvOrDefault("PROXY_USER", "")
    proxyPass := getEnvOrDefault("PROXY_PASS", "")
    
    if proxyURL != "" {
        // 配置HTTP代理
        api.AddConfigOption(config, "http_proxy", proxyURL)
        fmt.Printf("HTTP 代理: %s\n", proxyURL)
        
        // 配置HTTPS代理（通常相同）
        api.AddConfigOption(config, "https_proxy", proxyURL)
        fmt.Printf("HTTPS 代理: %s\n", proxyURL)
        
        // 如果提供了代理认证
        if proxyUser != "" && proxyPass != "" {
            api.AddConfigOption(config, "http_proxy.user", proxyUser)
            api.AddConfigOption(config, "http_proxy.password", proxyPass)
            fmt.Printf("代理认证: %s\n", proxyUser)
        }
        
        // 配置代理绕过本地地址
        api.AddConfigOption(config, "http_proxy.no_proxy", "localhost,127.0.0.1,*.local")
        fmt.Println("代理绕过: localhost,127.0.0.1,*.local")
    } else {
        fmt.Println("无需代理配置")
    }
    
    // 其他网络设置
    api.AddConfigOption(config, "http_timeout", "300")
    api.AddConfigOption(config, "http_retries", "3")
    
    fmt.Println("HTTP 超时: 300 秒")
    fmt.Println("HTTP 重试: 3 次")
    
    // 显示代理配置
    fmt.Println("\n=== 代理设置 ===")
    proxyOptions := []string{
        "http_proxy", "https_proxy", "http_proxy.user", 
        "http_proxy.password", "http_proxy.no_proxy",
        "http_timeout", "http_retries",
    }
    
    for _, key := range proxyOptions {
        value := api.GetConfigOption(config, key)
        if value != "" {
            displayValue := value
            if key == "http_proxy.password" {
                displayValue = "***已屏蔽***"
            }
            fmt.Printf("  %s: %s\n", key, displayValue)
        }
    }
    
    // 保存配置
    err := api.SaveConfig(config, "ProxyConfig.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n代理配置已保存！")
}

func getEnvOrDefault(envVar, defaultValue string) string {
    if value := os.Getenv(envVar); value != "" {
        return value
    }
    return defaultValue
}
```

## 示例 3: 环境特定配置

基于环境配置选项：

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
    config := api.CreateDefaultConfig()
    
    // 确定环境
    environment := getEnvOrDefault("ENVIRONMENT", "development")
    fmt.Printf("配置环境: %s\n", environment)
    
    // 根据环境配置
    switch environment {
    case "development":
        configureDevelopment(api, config)
    case "staging":
        configureStaging(api, config)
    case "production":
        configureProduction(api, config)
    default:
        configureDefault(api, config)
    }
    
    // 显示最终配置
    displayConfiguration(api, config, environment)
    
    // 保存环境特定名称
    configFile := fmt.Sprintf("%s.Config", environment)
    err := api.SaveConfig(config, configFile)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("\n环境特定配置已保存到: %s\n", configFile)
}

func configureDevelopment(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== 开发环境配置 ===")
    
    // 使用本地包文件夹以便快速访问
    api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // 启用详细日志记录以便调试
    api.AddConfigOption(config, "verbosity", "detailed")
    
    // 使用最高依赖版本获取最新功能
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    
    // 启用自动包还原
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // 允许预发布包
    api.AddConfigOption(config, "allowPrereleaseVersions", "true")
    
    fmt.Println("  - 本地包文件夹")
    fmt.Println("  - 启用详细日志")
    fmt.Println("  - 允许预发布包")
}

func configureStaging(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== 预发布环境配置 ===")
    
    // 使用共享预发布包文件夹
    api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // 适度详细程度
    api.AddConfigOption(config, "verbosity", "normal")
    
    // 使用稳定依赖版本
    api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
    
    // 启用包还原
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // 禁用预发布包
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    
    // 设置预发布推送源
    api.AddConfigOption(config, "defaultPushSource", "https://staging.company.com/nuget")
    
    fmt.Println("  - 共享预发布包")
    fmt.Println("  - 仅稳定版本")
    fmt.Println("  - 预发布推送源")
}

func configureProduction(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== 生产环境配置 ===")
    
    // 使用生产包文件夹
    api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // 最小详细程度以提高性能
    api.AddConfigOption(config, "verbosity", "quiet")
    
    // 使用精确依赖版本以保证稳定性
    api.AddConfigOption(config, "dependencyVersion", "Exact")
    
    // 启用包还原
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    api.AddConfigOption(config, "packageRestore", "true")
    
    // 严格禁用预发布包
    api.AddConfigOption(config, "allowPrereleaseVersions", "false")
    
    // 设置生产推送源
    api.AddConfigOption(config, "defaultPushSource", "https://prod.company.com/nuget")
    
    // 启用包验证
    api.AddConfigOption(config, "signatureValidationMode", "require")
    
    // 设置可靠性超时
    api.AddConfigOption(config, "http_timeout", "600")
    api.AddConfigOption(config, "http_retries", "5")
    
    fmt.Println("  - 生产包文件夹")
    fmt.Println("  - 精确版本以保证稳定性")
    fmt.Println("  - 包签名验证")
    fmt.Println("  - 扩展超时")
}

func configureDefault(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("=== 默认配置 ===")
    
    // 标准配置
    homeDir := getEnvOrDefault("HOME", getEnvOrDefault("USERPROFILE", "."))
    globalPackages := filepath.Join(homeDir, ".nuget", "packages")
    
    api.AddConfigOption(config, "globalPackagesFolder", globalPackages)
    api.AddConfigOption(config, "repositoryPath", "./packages")
    api.AddConfigOption(config, "dependencyVersion", "Highest")
    api.AddConfigOption(config, "automaticPackageRestore", "true")
    
    fmt.Println("  - 标准包位置")
    fmt.Println("  - 默认依赖解析")
}

func displayConfiguration(api *nuget.API, config *types.NuGetConfig, environment string) {
    fmt.Printf("\n=== 最终 %s 配置 ===\n", environment)
    
    if config.Config != nil {
        for _, option := range config.Config.Add {
            value := option.Value
            if option.Key == "http_proxy.password" {
                value = "***已屏蔽***"
            }
            fmt.Printf("  %s: %s\n", option.Key, value)
        }
    }
}

func getEnvOrDefault(envVar, defaultValue string) string {
    if value := os.Getenv(envVar); value != "" {
        return value
    }
    return defaultValue
}
```

## 示例 4: 高级配置管理

管理复杂配置场景：

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
    
    // 加载现有配置或创建新配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    fmt.Printf("管理配置: %s\n", configPath)
    fmt.Println("=== 高级配置管理 ===")
    
    // 创建配置管理器
    manager := NewConfigManager(api, config)
    
    // 应用配置模板
    manager.ApplyTemplate("enterprise")
    
    // 验证配置
    manager.ValidateConfiguration()
    
    // 优化配置
    manager.OptimizeConfiguration()
    
    // 生成配置报告
    manager.GenerateReport()
    
    // 保存优化后的配置
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Printf("\n优化后的配置已保存到: %s\n", configPath)
}

type ConfigManager struct {
    api    *nuget.API
    config *types.NuGetConfig
}

func NewConfigManager(api *nuget.API, config *types.NuGetConfig) *ConfigManager {
    return &ConfigManager{api: api, config: config}
}

func (cm *ConfigManager) ApplyTemplate(templateName string) {
    fmt.Printf("应用模板: %s\n", templateName)
    
    switch templateName {
    case "enterprise":
        cm.applyEnterpriseTemplate()
    case "developer":
        cm.applyDeveloperTemplate()
    case "ci-cd":
        cm.applyCICDTemplate()
    default:
        fmt.Printf("未知模板: %s\n", templateName)
    }
}

func (cm *ConfigManager) applyEnterpriseTemplate() {
    // 企业特定设置
    options := map[string]string{
        "globalPackagesFolder":     "/enterprise/packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "HighestMinor",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "signatureValidationMode": "require",
        "http_timeout":            "300",
        "http_retries":            "3",
        "verbosity":               "normal",
        "allowPrereleaseVersions": "false",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - 应用企业安全设置")
    fmt.Println("  - 配置稳定依赖解析")
    fmt.Println("  - 设置企业包位置")
}

func (cm *ConfigManager) applyDeveloperTemplate() {
    options := map[string]string{
        "globalPackagesFolder":     "./dev-packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "Highest",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "verbosity":               "detailed",
        "allowPrereleaseVersions": "true",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - 应用开发者友好设置")
}

func (cm *ConfigManager) applyCICDTemplate() {
    options := map[string]string{
        "globalPackagesFolder":     "/tmp/packages",
        "repositoryPath":          "./packages",
        "dependencyVersion":       "Exact",
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "verbosity":               "minimal",
        "allowPrereleaseVersions": "false",
        "http_timeout":            "600",
        "http_retries":            "5",
    }
    
    for key, value := range options {
        cm.api.AddConfigOption(cm.config, key, value)
    }
    
    fmt.Println("  - 应用CI/CD优化设置")
}

func (cm *ConfigManager) ValidateConfiguration() {
    fmt.Println("\n=== 配置验证 ===")
    
    issues := 0
    
    // 检查必需选项
    requiredOptions := []string{
        "globalPackagesFolder",
        "repositoryPath",
        "dependencyVersion",
    }
    
    for _, option := range requiredOptions {
        value := cm.api.GetConfigOption(cm.config, option)
        if value == "" {
            fmt.Printf("  ❌ 缺少必需选项: %s\n", option)
            issues++
        } else {
            fmt.Printf("  ✅ %s: %s\n", option, value)
        }
    }
    
    // 验证依赖版本值
    depVersion := cm.api.GetConfigOption(cm.config, "dependencyVersion")
    validVersions := []string{"Lowest", "HighestPatch", "HighestMinor", "Highest", "Exact"}
    if depVersion != "" && !contains(validVersions, depVersion) {
        fmt.Printf("  ⚠️  无效的dependencyVersion: %s\n", depVersion)
        issues++
    }
    
    if issues == 0 {
        fmt.Println("  ✅ 配置验证通过")
    } else {
        fmt.Printf("  ⚠️  发现 %d 个配置问题\n", issues)
    }
}

func (cm *ConfigManager) OptimizeConfiguration() {
    fmt.Println("\n=== 配置优化 ===")
    
    // 移除重复或冲突选项
    cm.removeDuplicateOptions()
    
    // 为缺失选项设置最优默认值
    cm.setOptimalDefaults()
    
    fmt.Println("  ✅ 配置已优化")
}

func (cm *ConfigManager) removeDuplicateOptions() {
    // 这里会检查重复键并解决冲突
    // 现在只是报告我们会做什么
    fmt.Println("  - 检查重复选项")
}

func (cm *ConfigManager) setOptimalDefaults() {
    defaults := map[string]string{
        "automaticPackageRestore": "true",
        "packageRestore":          "true",
        "http_timeout":            "300",
        "http_retries":            "3",
    }
    
    for key, value := range defaults {
        if cm.api.GetConfigOption(cm.config, key) == "" {
            cm.api.AddConfigOption(cm.config, key, value)
            fmt.Printf("  - 设置默认 %s: %s\n", key, value)
        }
    }
}

func (cm *ConfigManager) GenerateReport() {
    fmt.Println("\n=== 配置报告 ===")
    
    if cm.config.Config != nil {
        fmt.Printf("总配置选项: %d\n", len(cm.config.Config.Add))
        
        categories := map[string][]string{
            "包管理": {"globalPackagesFolder", "repositoryPath", "dependencyVersion"},
            "网络":   {"http_proxy", "http_timeout", "http_retries"},
            "安全":   {"signatureValidationMode", "allowPrereleaseVersions"},
            "还原":   {"automaticPackageRestore", "packageRestore"},
        }
        
        for category, options := range categories {
            fmt.Printf("\n%s:\n", category)
            for _, option := range options {
                value := cm.api.GetConfigOption(cm.config, option)
                if value != "" {
                    fmt.Printf("  %s: %s\n", option, value)
                }
            }
        }
    }
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

## 关键配置选项

### 包管理
- `globalPackagesFolder`: 全局包缓存位置
- `repositoryPath`: 项目包文件夹
- `dependencyVersion`: 依赖解析策略

### 网络设置
- `http_proxy`: HTTP代理服务器
- `http_timeout`: 请求超时（秒）
- `http_retries`: 重试次数

### 安全选项
- `signatureValidationMode`: 包签名验证
- `allowPrereleaseVersions`: 允许预发布包

### 还原行为
- `automaticPackageRestore`: 启用自动还原
- `packageRestore`: 启用包还原

## 最佳实践

1. **环境特定配置**: 为不同环境使用不同设置
2. **验证选项**: 检查选项值是否有效
3. **使用模板**: 应用一致的配置模式
4. **文档设置**: 注释配置选择
5. **定期审查**: 定期审查和优化设置

## 下一步

掌握配置选项后：

1. 学习 [序列化](./serialization.md) 进行自定义XML处理
2. 探索 [位置感知编辑](./position-aware-editing.md) 进行精确修改
3. 研究 [类型API](/api/types) 了解配置结构详情

本指南为管理不同场景和环境中的 NuGet 配置选项提供了全面的示例。
