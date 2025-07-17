# 凭证管理

本示例演示如何使用 NuGet Config Parser 库管理 NuGet 包源的认证凭证，包括用户名/密码、API密钥和令牌认证。

## 概述

凭证管理包括：
- 为私有包源添加用户名/密码认证
- 管理API密钥和访问令牌
- 处理不同的认证方式
- 安全存储和检索凭证
- 批量凭证配置

## 示例 1: 基本凭证管理

为包源添加和管理基本认证：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== 基本凭证管理 ===")
    
    // 添加需要认证的包源
    privateSources := []struct {
        name     string
        url      string
        username string
        password string
    }{
        {"company-internal", "https://nuget.company.com/v3/index.json", "employee", "company_pass123"},
        {"azure-artifacts", "https://pkgs.dev.azure.com/myorg/_packaging/myfeed/nuget/v3/index.json", "myuser", "pat_token_here"},
        {"private-feed", "https://private.nuget.com/api/v3/index.json", "user@company.com", "secure_password"},
    }
    
    fmt.Println("添加私有包源和凭证...")
    
    for _, source := range privateSources {
        // 添加包源
        api.AddPackageSource(config, source.name, source.url, "3")
        fmt.Printf("✅ 添加包源: %s\n", source.name)
        
        // 添加凭证
        api.AddCredential(config, source.name, source.username, source.password)
        fmt.Printf("🔐 添加凭证: %s (用户: %s)\n", source.name, source.username)
    }
    
    // 添加公共源（无需凭证）
    api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    fmt.Println("✅ 添加公共源: nuget.org")
    
    // 显示配置摘要
    fmt.Printf("\n=== 配置摘要 ===\n")
    fmt.Printf("总包源数: %d\n", len(config.PackageSources.Add))
    
    // 显示所有源及其认证状态
    for _, source := range config.PackageSources.Add {
        hasCredentials := "❌ 无认证"
        
        if config.PackageSourceCredentials != nil {
            if _, exists := config.PackageSourceCredentials.Sources[source.Key]; exists {
                hasCredentials = "✅ 已配置认证"
            }
        }
        
        fmt.Printf("- %s: %s [%s]\n", source.Key, source.Value, hasCredentials)
    }
    
    // 验证凭证
    fmt.Println("\n=== 凭证验证 ===")
    for _, source := range privateSources {
        credential := api.GetCredential(config, source.name)
        if credential != nil {
            fmt.Printf("✅ %s: 用户名=%s, 密码已设置\n", source.name, credential.Username.Value)
        } else {
            fmt.Printf("❌ %s: 未找到凭证\n", source.name)
        }
    }
    
    // 保存配置
    err := api.SaveConfig(config, "CredentialsDemo.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n配置已保存到 CredentialsDemo.Config")
    fmt.Println("⚠️  注意: 生产环境中应使用更安全的凭证存储方式")
}
```

## 示例 2: 高级凭证配置

处理复杂的认证场景：

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
    config := api.CreateDefaultConfig()
    
    fmt.Println("=== 高级凭证配置 ===")
    
    // 从环境变量获取凭证
    credentialsFromEnv := getCredentialsFromEnvironment()
    
    // 配置不同类型的认证
    configureAuthentication(api, config, credentialsFromEnv)
    
    // 演示凭证更新
    demonstrateCredentialUpdates(api, config)
    
    // 演示凭证移除
    demonstrateCredentialRemoval(api, config)
    
    // 保存最终配置
    err := api.SaveConfig(config, "AdvancedCredentials.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("\n高级凭证配置完成")
}

type CredentialConfig struct {
    SourceName string
    URL        string
    AuthType   string // "basic", "token", "apikey"
    Username   string
    Password   string
    Token      string
    APIKey     string
}

func getCredentialsFromEnvironment() []CredentialConfig {
    fmt.Println("从环境变量获取凭证...")
    
    credentials := []CredentialConfig{
        {
            SourceName: "company-prod",
            URL:        "https://prod.company.com/nuget",
            AuthType:   "basic",
            Username:   getEnvOrDefault("COMPANY_NUGET_USER", "default_user"),
            Password:   getEnvOrDefault("COMPANY_NUGET_PASS", "default_pass"),
        },
        {
            SourceName: "azure-devops",
            URL:        "https://pkgs.dev.azure.com/myorg/_packaging/main/nuget/v3/index.json",
            AuthType:   "token",
            Username:   "PAT", // Personal Access Token 通常使用 PAT 作为用户名
            Token:      getEnvOrDefault("AZURE_DEVOPS_PAT", "default_pat_token"),
        },
        {
            SourceName: "github-packages",
            URL:        "https://nuget.pkg.github.com/myorg/index.json",
            AuthType:   "token",
            Username:   getEnvOrDefault("GITHUB_USERNAME", "github_user"),
            Token:      getEnvOrDefault("GITHUB_TOKEN", "github_token"),
        },
        {
            SourceName: "myget-feed",
            URL:        "https://www.myget.org/F/myfeed/api/v3/index.json",
            AuthType:   "apikey",
            APIKey:     getEnvOrDefault("MYGET_API_KEY", "myget_api_key"),
        },
    }
    
    for _, cred := range credentials {
        fmt.Printf("✅ 环境凭证: %s (%s)\n", cred.SourceName, cred.AuthType)
    }
    
    return credentials
}

func configureAuthentication(api *nuget.API, config *types.NuGetConfig, credentials []CredentialConfig) {
    fmt.Println("\n配置不同类型的认证...")
    
    for _, cred := range credentials {
        // 添加包源
        api.AddPackageSource(config, cred.SourceName, cred.URL, "3")
        fmt.Printf("📦 添加源: %s\n", cred.SourceName)
        
        // 根据认证类型配置凭证
        switch cred.AuthType {
        case "basic":
            api.AddCredential(config, cred.SourceName, cred.Username, cred.Password)
            fmt.Printf("🔐 基本认证: %s (用户: %s)\n", cred.SourceName, cred.Username)
            
        case "token":
            // 对于令牌认证，通常将令牌作为密码
            username := cred.Username
            if username == "" {
                username = "token" // 默认用户名
            }
            api.AddCredential(config, cred.SourceName, username, cred.Token)
            fmt.Printf("🎫 令牌认证: %s (用户: %s)\n", cred.SourceName, username)
            
        case "apikey":
            // API密钥通常作为密码，用户名可以是任意值
            api.AddCredential(config, cred.SourceName, "apikey", cred.APIKey)
            fmt.Printf("🔑 API密钥认证: %s\n", cred.SourceName)
            
        default:
            fmt.Printf("⚠️  未知认证类型: %s for %s\n", cred.AuthType, cred.SourceName)
        }
    }
}

func demonstrateCredentialUpdates(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== 凭证更新演示 ===")
    
    // 模拟凭证轮换场景
    updates := []struct {
        sourceName  string
        newUsername string
        newPassword string
        reason      string
    }{
        {"company-prod", "new_user", "new_secure_pass", "定期密码轮换"},
        {"azure-devops", "PAT", "new_pat_token_2024", "PAT令牌更新"},
        {"github-packages", "updated_user", "new_github_token", "GitHub令牌刷新"},
    }
    
    for _, update := range updates {
        fmt.Printf("🔄 更新凭证: %s (%s)\n", update.sourceName, update.reason)
        
        // 更新凭证（实际上是重新添加）
        api.AddCredential(config, update.sourceName, update.newUsername, update.newPassword)
        
        // 验证更新
        credential := api.GetCredential(config, update.sourceName)
        if credential != nil && credential.Username.Value == update.newUsername {
            fmt.Printf("✅ 凭证更新成功: %s\n", update.sourceName)
        } else {
            fmt.Printf("❌ 凭证更新失败: %s\n", update.sourceName)
        }
    }
}

func demonstrateCredentialRemoval(api *nuget.API, config *types.NuGetConfig) {
    fmt.Println("\n=== 凭证移除演示 ===")
    
    // 移除不再需要的凭证
    sourcesToRemove := []string{"myget-feed"}
    
    for _, sourceName := range sourcesToRemove {
        fmt.Printf("🗑️  移除凭证: %s\n", sourceName)
        
        // 移除凭证（通过重建凭证映射）
        if config.PackageSourceCredentials != nil {
            delete(config.PackageSourceCredentials.Sources, sourceName)
            fmt.Printf("✅ 已移除 %s 的凭证\n", sourceName)
        }
        
        // 验证移除
        credential := api.GetCredential(config, sourceName)
        if credential == nil {
            fmt.Printf("✅ 确认凭证已移除: %s\n", sourceName)
        } else {
            fmt.Printf("❌ 凭证移除失败: %s\n", sourceName)
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

## 示例 3: 企业凭证管理

企业环境中的凭证管理最佳实践：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Println("=== 企业凭证管理 ===")
    
    // 创建凭证管理器
    credManager := NewEnterpriseCredentialManager(api)
    
    // 加载企业凭证配置
    err := credManager.LoadEnterpriseCredentials()
    if err != nil {
        log.Fatalf("加载企业凭证失败: %v", err)
    }
    
    // 为不同环境创建配置
    environments := []string{"development", "staging", "production"}
    
    for _, env := range environments {
        fmt.Printf("\n配置 %s 环境凭证...\n", env)
        
        config, err := credManager.CreateEnvironmentConfig(env)
        if err != nil {
            log.Printf("创建 %s 环境配置失败: %v", env, err)
            continue
        }
        
        // 保存环境配置
        configPath := fmt.Sprintf("Enterprise.%s.Config", env)
        err = api.SaveConfig(config, configPath)
        if err != nil {
            log.Printf("保存 %s 配置失败: %v", env, err)
            continue
        }
        
        fmt.Printf("✅ %s 环境配置已保存\n", env)
        
        // 显示凭证摘要（不显示敏感信息）
        credManager.DisplayCredentialSummary(config, env)
    }
    
    // 演示凭证轮换
    credManager.DemonstrateCredentialRotation()
    
    fmt.Println("\n企业凭证管理完成")
}

type EnterpriseCredentialManager struct {
    api         *nuget.API
    credentials map[string]EnterpriseCredential
}

type EnterpriseCredential struct {
    SourceName  string
    URL         string
    Environment []string // 适用的环境
    AuthType    string
    Username    string
    Password    string
    Description string
    LastUpdated string
}

func NewEnterpriseCredentialManager(api *nuget.API) *EnterpriseCredentialManager {
    return &EnterpriseCredentialManager{
        api:         api,
        credentials: make(map[string]EnterpriseCredential),
    }
}

func (ecm *EnterpriseCredentialManager) LoadEnterpriseCredentials() error {
    fmt.Println("加载企业凭证配置...")
    
    // 模拟从安全存储加载凭证
    enterpriseCredentials := []EnterpriseCredential{
        {
            SourceName:  "enterprise-stable",
            URL:         "https://nuget.enterprise.com/stable",
            Environment: []string{"development", "staging", "production"},
            AuthType:    "basic",
            Username:    "enterprise_service",
            Password:    "enterprise_secure_pass_2024",
            Description: "企业稳定包源",
            LastUpdated: "2024-01-15",
        },
        {
            SourceName:  "enterprise-dev",
            URL:         "https://dev.nuget.enterprise.com",
            Environment: []string{"development"},
            AuthType:    "basic",
            Username:    "dev_service",
            Password:    "dev_pass_2024",
            Description: "企业开发包源",
            LastUpdated: "2024-01-10",
        },
        {
            SourceName:  "enterprise-staging",
            URL:         "https://staging.nuget.enterprise.com",
            Environment: []string{"staging"},
            AuthType:    "basic",
            Username:    "staging_service",
            Password:    "staging_pass_2024",
            Description: "企业预发布包源",
            LastUpdated: "2024-01-12",
        },
        {
            SourceName:  "enterprise-prod",
            URL:         "https://prod.nuget.enterprise.com",
            Environment: []string{"production"},
            AuthType:    "basic",
            Username:    "prod_service",
            Password:    "prod_secure_pass_2024",
            Description: "企业生产包源",
            LastUpdated: "2024-01-14",
        },
        {
            SourceName:  "azure-artifacts-enterprise",
            URL:         "https://pkgs.dev.azure.com/enterprise/_packaging/main/nuget/v3/index.json",
            Environment: []string{"development", "staging", "production"},
            AuthType:    "token",
            Username:    "PAT",
            Password:    "azure_pat_token_enterprise_2024",
            Description: "企业Azure DevOps包源",
            LastUpdated: "2024-01-16",
        },
    }
    
    for _, cred := range enterpriseCredentials {
        ecm.credentials[cred.SourceName] = cred
        fmt.Printf("✅ 加载凭证: %s (%s)\n", cred.SourceName, cred.Description)
    }
    
    fmt.Printf("总共加载 %d 个企业凭证\n", len(ecm.credentials))
    return nil
}

func (ecm *EnterpriseCredentialManager) CreateEnvironmentConfig(environment string) (*types.NuGetConfig, error) {
    config := ecm.api.CreateDefaultConfig()
    
    // 添加公共源
    ecm.api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
    
    // 添加适用于此环境的企业源
    addedCount := 0
    for _, cred := range ecm.credentials {
        // 检查凭证是否适用于当前环境
        if ecm.isCredentialApplicable(cred, environment) {
            // 添加包源
            ecm.api.AddPackageSource(config, cred.SourceName, cred.URL, "3")
            
            // 添加凭证
            ecm.api.AddCredential(config, cred.SourceName, cred.Username, cred.Password)
            
            addedCount++
            fmt.Printf("  ✅ 添加: %s\n", cred.SourceName)
        }
    }
    
    // 根据环境设置特定配置
    ecm.configureEnvironmentSpecificSettings(config, environment)
    
    fmt.Printf("  为 %s 环境添加了 %d 个凭证\n", environment, addedCount)
    return config, nil
}

func (ecm *EnterpriseCredentialManager) isCredentialApplicable(cred EnterpriseCredential, environment string) bool {
    for _, env := range cred.Environment {
        if env == environment {
            return true
        }
    }
    return false
}

func (ecm *EnterpriseCredentialManager) configureEnvironmentSpecificSettings(config *types.NuGetConfig, environment string) {
    switch environment {
    case "development":
        ecm.api.AddConfigOption(config, "globalPackagesFolder", "./dev-packages")
        ecm.api.AddConfigOption(config, "dependencyVersion", "Highest")
        ecm.api.SetActivePackageSource(config, "enterprise-dev", "https://dev.nuget.enterprise.com")
        
    case "staging":
        ecm.api.AddConfigOption(config, "globalPackagesFolder", "/staging/packages")
        ecm.api.AddConfigOption(config, "dependencyVersion", "HighestMinor")
        ecm.api.SetActivePackageSource(config, "enterprise-staging", "https://staging.nuget.enterprise.com")
        
    case "production":
        ecm.api.AddConfigOption(config, "globalPackagesFolder", "/prod/packages")
        ecm.api.AddConfigOption(config, "dependencyVersion", "Exact")
        ecm.api.AddConfigOption(config, "signatureValidationMode", "require")
        ecm.api.SetActivePackageSource(config, "enterprise-prod", "https://prod.nuget.enterprise.com")
        
        // 生产环境禁用外部源
        ecm.api.DisablePackageSource(config, "nuget.org")
    }
}

func (ecm *EnterpriseCredentialManager) DisplayCredentialSummary(config *types.NuGetConfig, environment string) {
    fmt.Printf("\n--- %s 环境凭证摘要 ---\n", strings.ToUpper(environment))
    
    if config.PackageSourceCredentials == nil {
        fmt.Println("无配置凭证")
        return
    }
    
    fmt.Printf("已配置凭证的源 (%d):\n", len(config.PackageSourceCredentials.Sources))
    
    for sourceName, cred := range config.PackageSourceCredentials.Sources {
        // 查找原始凭证信息
        if enterpriseCred, exists := ecm.credentials[sourceName]; exists {
            fmt.Printf("  - %s\n", sourceName)
            fmt.Printf("    描述: %s\n", enterpriseCred.Description)
            fmt.Printf("    用户: %s\n", cred.Username.Value)
            fmt.Printf("    认证类型: %s\n", enterpriseCred.AuthType)
            fmt.Printf("    最后更新: %s\n", enterpriseCred.LastUpdated)
        } else {
            fmt.Printf("  - %s (用户: %s)\n", sourceName, cred.Username.Value)
        }
    }
}

func (ecm *EnterpriseCredentialManager) DemonstrateCredentialRotation() {
    fmt.Println("\n=== 凭证轮换演示 ===")
    
    // 模拟定期凭证轮换
    rotationCandidates := []string{"enterprise-stable", "azure-artifacts-enterprise"}
    
    for _, sourceName := range rotationCandidates {
        if cred, exists := ecm.credentials[sourceName]; exists {
            fmt.Printf("🔄 轮换凭证: %s\n", sourceName)
            
            // 生成新密码（实际应用中应使用安全的密码生成器）
            newPassword := fmt.Sprintf("%s_rotated_2024", cred.Password)
            
            // 更新凭证
            cred.Password = newPassword
            cred.LastUpdated = "2024-01-20"
            ecm.credentials[sourceName] = cred
            
            fmt.Printf("✅ %s 凭证已轮换\n", sourceName)
        }
    }
    
    fmt.Println("凭证轮换完成")
}
```

## 安全最佳实践

### 1. 凭证存储

```go
// 不推荐：明文存储
api.AddCredential(config, "source", "user", "plaintext_password")

// 推荐：从安全存储获取
password := getFromSecureStore("source_password")
api.AddCredential(config, "source", "user", password)
```

### 2. 环境变量使用

```go
// 从环境变量获取敏感信息
username := os.Getenv("NUGET_USERNAME")
password := os.Getenv("NUGET_PASSWORD")

if username == "" || password == "" {
    log.Fatal("缺少必需的凭证环境变量")
}

api.AddCredential(config, "private-source", username, password)
```

### 3. 凭证轮换

```go
// 定期更新凭证
func rotateCredentials(api *nuget.API, config *types.NuGetConfig) {
    // 获取新凭证
    newPassword := generateSecurePassword()
    
    // 更新配置
    api.AddCredential(config, "source", "user", newPassword)
    
    // 记录轮换
    log.Printf("凭证已轮换: %s", time.Now().Format("2006-01-02"))
}
```

## 关键概念

### 认证类型

1. **基本认证** - 用户名/密码
2. **令牌认证** - API令牌或PAT
3. **API密钥** - 单一密钥认证

### 安全考虑

1. **加密存储** - 不要明文存储密码
2. **环境隔离** - 不同环境使用不同凭证
3. **定期轮换** - 定期更新密码和令牌
4. **最小权限** - 只授予必要的权限

## 下一步

掌握凭证管理后：

1. 学习 [配置选项](./config-options.md) 进行高级设置
2. 探索 [序列化](./serialization.md) 了解配置输出
3. 研究 [位置感知编辑](./position-aware-editing.md) 进行精确修改

本指南为 NuGet 凭证管理提供了全面的示例，涵盖了从基本认证到企业级安全管理的各种场景。
