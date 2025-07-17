# 示例

本节提供全面的示例，演示如何在各种场景中使用 NuGet Config Parser 库。每个示例都包含完整的、可运行的代码和解释。

## 概述

示例按功能和复杂性组织：

- **[基本解析](./basic-parsing.md)** - 简单的配置文件解析
- **[查找配置](./finding-configs.md)** - 在系统中定位配置文件
- **[创建配置](./creating-configs.md)** - 创建新的配置文件
- **[修改配置](./modifying-configs.md)** - 修改现有配置
- **[包源管理](./package-sources.md)** - 管理包源及其属性
- **[凭证管理](./credentials.md)** - 处理私有源的身份验证
- **[配置选项](./config-options.md)** - 管理全局 NuGet 设置
- **[序列化](./serialization.md)** - 自定义 XML 处理和验证
- **[位置感知编辑](./position-aware-editing.md)** - 高级编辑，最小化差异

## 快速开始示例

这是一个简单的示例，帮助您开始：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    // 创建 API 实例
    api := nuget.NewAPI()
    
    // 查找并解析第一个可用的配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // 如果未找到配置，创建默认配置
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        
        // 保存默认配置
        err = api.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("保存默认配置失败: %v", err)
        }
        
        fmt.Printf("创建了默认配置: %s\n", configPath)
    } else {
        fmt.Printf("找到现有配置: %s\n", configPath)
    }
    
    // 显示包源
    fmt.Printf("包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "启用"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "禁用"
        }
        
        fmt.Printf("  - %s: %s (%s)", source.Key, source.Value, status)
        if source.ProtocolVersion != "" {
            fmt.Printf(" [v%s]", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // 添加自定义包源
    api.AddPackageSource(config, "example", "https://example.com/nuget", "3")
    
    // 保存更新的配置
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存更新配置失败: %v", err)
    }
    
    fmt.Println("添加了示例包源并保存了配置")
}
```

## 常见模式

### 错误处理模式

```go
config, err := api.ParseFromFile(configPath)
if err != nil {
    if errors.IsNotFoundError(err) {
        // 处理缺失文件
        config = api.CreateDefaultConfig()
    } else if errors.IsParseError(err) {
        // 处理解析错误
        log.Fatalf("无效的配置格式: %v", err)
    } else {
        // 处理其他错误
        log.Fatalf("意外错误: %v", err)
    }
}
```

### 配置发现模式

```go
// 尝试多种方法查找配置
var config *types.NuGetConfig
var configPath string
var err error

// 1. 尝试查找并解析现有配置
config, configPath, err = api.FindAndParseConfig()
if err == nil {
    fmt.Printf("使用现有配置: %s\n", configPath)
} else {
    // 2. 尝试项目特定配置
    configPath, err = api.FindProjectConfig(".")
    if err == nil {
        config, err = api.ParseFromFile(configPath)
        if err == nil {
            fmt.Printf("使用项目配置: %s\n", configPath)
        }
    }
}

// 3. 回退到默认配置
if config == nil {
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
    fmt.Println("使用默认配置")
}
```

### 批量修改模式

```go
// 高效地进行多个更改
api.AddPackageSource(config, "feed1", "https://feed1.com", "3")
api.AddPackageSource(config, "feed2", "https://feed2.com", "3")
api.AddCredential(config, "feed1", "user1", "pass1")
api.AddCredential(config, "feed2", "user2", "pass2")
api.DisablePackageSource(config, "old-feed")

// 一次性保存所有更改
err := api.SaveConfig(config, configPath)
if err != nil {
    log.Fatalf("保存更改失败: %v", err)
}
```

### 位置感知编辑模式

```go
// 使用位置跟踪解析
parseResult, err := api.ParseFromFileWithPositions(configPath)
if err != nil {
    log.Fatalf("带位置解析失败: %v", err)
}

// 创建编辑器
editor := api.CreateConfigEditor(parseResult)

// 进行更改
editor.AddPackageSource("new-feed", "https://new-feed.com", "3")
editor.UpdatePackageSourceURL("existing-feed", "https://updated-url.com")

// 应用最小差异的更改
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatalf("应用编辑失败: %v", err)
}

// 保存修改的内容
err = os.WriteFile(configPath, modifiedContent, 0644)
if err != nil {
    log.Fatalf("保存文件失败: %v", err)
}
```

## 示例分类

### 初学者示例

非常适合开始使用库：

1. **[基本解析](./basic-parsing.md)** - 读取和显示配置文件
2. **[查找配置](./finding-configs.md)** - 在系统中定位配置文件
3. **[创建配置](./creating-configs.md)** - 创建新的配置文件

### 中级示例

用于常见的配置管理任务：

4. **[修改配置](./modifying-configs.md)** - 更新现有配置
5. **[包源管理](./package-sources.md)** - 管理包源及其属性
6. **[凭证管理](./credentials.md)** - 处理私有源的身份验证

### 高级示例

用于复杂场景和优化：

7. **[配置选项](./config-options.md)** - 管理全局 NuGet 设置
8. **[序列化](./serialization.md)** - 自定义 XML 处理和验证
9. **[位置感知编辑](./position-aware-editing.md)** - 保留格式并最小化差异

## 运行示例

所有示例都设计为自包含且可运行。要运行示例：

1. 使用示例代码创建新的 Go 文件
2. 如果需要，初始化 Go 模块：
   ```bash
   go mod init example
   go get github.com/scagogogo/nuget-config-parser
   ```
3. 运行示例：
   ```bash
   go run main.go
   ```

## 示例数据

许多示例使用示例 NuGet.Config 文件。这是一个典型的示例配置：

```xml
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="C:\LocalPackages" />
    <add key="company" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
  </packageSources>
  <packageSourceCredentials>
    <company>
      <add key="Username" value="companyuser" />
      <add key="ClearTextPassword" value="companypass" />
    </company>
  </packageSourceCredentials>
  <disabledPackageSources>
    <add key="local" value="true" />
  </disabledPackageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="C:\packages" />
    <add key="repositoryPath" value=".\packages" />
  </config>
</configuration>
```

## 演示的最佳实践

示例演示了这些最佳实践：

- **错误处理**: 针对不同场景的正确错误检查和处理
- **资源管理**: API 实例和文件操作的高效使用
- **配置验证**: 确保配置在保存前有效
- **安全性**: 凭证和敏感信息的安全处理
- **性能**: 高效的批量操作和最小文件修改
- **可维护性**: 清洁、可读的代码，具有良好的关注点分离

## 贡献示例

如果您有一个有用的示例，这里没有涵盖，请考虑为项目贡献。好的示例应该是：

- **完整**: 包括所有必要的导入和错误处理
- **专注**: 演示一个特定概念或用例
- **文档化**: 包括解释重要部分的注释
- **测试**: 验证示例与当前库版本一起工作

## 下一步

查看示例后：

1. 查看 [API 参考](/zh/api/) 获取详细的方法文档
2. 阅读 [指南](/zh/guide/) 获取概念信息
3. 探索库的源代码以了解高级使用模式
4. 考虑贡献您自己的示例或改进
