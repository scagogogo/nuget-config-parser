# 基本解析

本示例演示使用 NuGet Config Parser 库解析 NuGet 配置文件的基本操作。

## 概述

基本解析包括：
- 从各种来源读取配置文件
- 处理不同的文件格式和位置
- 显示配置内容
- 基本错误处理

## 示例 1: 从文件解析

最常见的场景是解析现有配置文件：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    // 创建 API 实例
    api := nuget.NewAPI()
    
    // 从文件解析配置
    configPath := "NuGet.Config"
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Printf("未找到配置文件: %s\n", configPath)
            return
        }
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 显示基本信息
    fmt.Printf("从以下位置加载配置: %s\n", configPath)
    fmt.Printf("包源数量: %d\n", len(config.PackageSources.Add))
    
    // 列出所有包源
    fmt.Println("\n包源:")
    for i, source := range config.PackageSources.Add {
        fmt.Printf("%d. %s\n", i+1, source.Key)
        fmt.Printf("   URL: %s\n", source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf("   协议: v%s\n", source.ProtocolVersion)
        }
        fmt.Println()
    }
}
```

## 示例 2: 从字符串解析

有时您需要从字符串解析配置（例如，从数据库或 API）：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // XML 配置字符串
    xmlConfig := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="/path/to/local/packages" />
  </packageSources>
  <config>
    <add key="globalPackagesFolder" value="/custom/packages" />
  </config>
</configuration>`
    
    // 从字符串解析
    config, err := api.ParseFromString(xmlConfig)
    if err != nil {
        log.Fatalf("解析 XML 失败: %v", err)
    }
    
    fmt.Println("从字符串成功解析配置！")
    
    // 显示包源
    fmt.Printf("找到 %d 个包源:\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
    }
    
    // 显示配置选项
    if config.Config != nil {
        fmt.Printf("\n配置选项 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("- %s: %s\n", option.Key, option.Value)
        }
    }
}
```

## 示例 3: 从 Reader 解析

用于流式场景或使用 io.Reader 时：

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
    
    // 示例 1: 使用 Reader 从文件解析
    file, err := os.Open("NuGet.Config")
    if err != nil {
        log.Printf("无法打开文件: %v", err)
        
        // 示例 2: 改为从字符串 reader 解析
        xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="example" value="https://example.com/nuget" />
  </packageSources>
</configuration>`
        
        reader := strings.NewReader(xmlContent)
        config, err := api.ParseFromReader(reader)
        if err != nil {
            log.Fatalf("从 reader 解析失败: %v", err)
        }
        
        fmt.Println("从字符串 reader 解析:")
        displayConfig(config)
        return
    }
    defer file.Close()
    
    // 从文件 reader 解析
    config, err := api.ParseFromReader(file)
    if err != nil {
        log.Fatalf("从文件 reader 解析失败: %v", err)
    }
    
    fmt.Println("从文件 reader 解析:")
    displayConfig(config)
}

func displayConfig(config *types.NuGetConfig) {
    fmt.Printf("包源: %d\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s\n", source.Key, source.Value)
    }
}
```

## 示例 4: 带错误处理的综合解析

处理各种错误条件的健壮示例：

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
    api := nuget.NewAPI()
    
    // 尝试使用综合错误处理解析配置
    config, err := parseConfigSafely(api, "NuGet.Config")
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 显示完整配置信息
    displayFullConfig(config)
}

func parseConfigSafely(api *nuget.API, configPath string) (*types.NuGetConfig, error) {
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Printf("未找到配置文件: %s\n", configPath)
            fmt.Println("创建默认配置...")
            return api.CreateDefaultConfig(), nil
        }
        
        if errors.IsParseError(err) {
            return nil, fmt.Errorf("无效的配置格式: %w", err)
        }
        
        if errors.IsFormatError(err) {
            return nil, fmt.Errorf("配置中的 XML 格式错误: %w", err)
        }
        
        return nil, fmt.Errorf("意外错误: %w", err)
    }
    
    return config, nil
}

func displayFullConfig(config *types.NuGetConfig) {
    fmt.Println("=== NuGet 配置 ===")
    
    // 包源
    fmt.Printf("\n包源 (%d):\n", len(config.PackageSources.Add))
    if len(config.PackageSources.Add) == 0 {
        fmt.Println("  (无)")
    } else {
        for i, source := range config.PackageSources.Add {
            fmt.Printf("%d. 键: %s\n", i+1, source.Key)
            fmt.Printf("   值: %s\n", source.Value)
            if source.ProtocolVersion != "" {
                fmt.Printf("   协议版本: %s\n", source.ProtocolVersion)
            }
            fmt.Println()
        }
    }
    
    // 禁用的包源
    if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
        fmt.Printf("禁用的包源 (%d):\n", len(config.DisabledPackageSources.Add))
        for _, disabled := range config.DisabledPackageSources.Add {
            fmt.Printf("  - %s\n", disabled.Key)
        }
        fmt.Println()
    }
    
    // 活跃包源
    if config.ActivePackageSource != nil {
        fmt.Printf("活跃包源:\n")
        fmt.Printf("  键: %s\n", config.ActivePackageSource.Add.Key)
        fmt.Printf("  值: %s\n", config.ActivePackageSource.Add.Value)
        fmt.Println()
    }
    
    // 包源凭证
    if config.PackageSourceCredentials != nil && len(config.PackageSourceCredentials.Sources) > 0 {
        fmt.Printf("包源凭证 (%d 个源):\n", len(config.PackageSourceCredentials.Sources))
        for sourceKey := range config.PackageSourceCredentials.Sources {
            fmt.Printf("  - %s (已配置凭证)\n", sourceKey)
        }
        fmt.Println()
    }
    
    // 配置选项
    if config.Config != nil && len(config.Config.Add) > 0 {
        fmt.Printf("配置选项 (%d):\n", len(config.Config.Add))
        for _, option := range config.Config.Add {
            fmt.Printf("  - %s: %s\n", option.Key, option.Value)
        }
        fmt.Println()
    }
}
```

## 关键概念

### 错误类型

库为不同场景提供特定的错误类型：

- `IsNotFoundError()`: 配置文件不存在
- `IsParseError()`: 无效的 XML 或解析问题
- `IsFormatError()`: 配置结构格式错误

### 配置结构

解析的配置包含：

- **PackageSources**: 可用包源列表
- **DisabledPackageSources**: 被禁用的源
- **ActivePackageSource**: 当前活跃源
- **PackageSourceCredentials**: 认证信息
- **Config**: 全局配置选项

### 最佳实践

1. **始终处理错误**: 检查特定错误类型
2. **提供回退**: 当文件缺失时创建默认配置
3. **验证输入**: 确保文件路径和内容有效
4. **显示有意义的信息**: 向用户显示解析的内容
5. **使用适当的解析方法**: 根据来源选择文件、字符串或 reader

## 下一步

掌握基本解析后：

1. 学习 [查找配置](./finding-configs.md) 来定位配置文件
2. 探索 [创建配置](./creating-configs.md) 来生成新配置
3. 研究 [修改配置](./modifying-configs.md) 来更新现有配置

## 常见问题

### 问题 1: 文件未找到

```go
// 在解析前始终检查文件是否存在
if !utils.FileExists(configPath) {
    fmt.Printf("配置文件不存在: %s\n", configPath)
    // 创建默认配置或适当处理
}
```

### 问题 2: 无效的 XML

```go
// 优雅地处理解析错误
if errors.IsParseError(err) {
    fmt.Printf("无效的 XML 格式: %v\n", err)
    // 考虑创建新配置
}
```

### 问题 3: 空配置

```go
// 检查配置是否有内容
if len(config.PackageSources.Add) == 0 {
    fmt.Println("配置没有包源")
    // 如果需要，添加默认源
}
```

本基本解析指南为使用 NuGet 配置文件提供了基础。示例演示了各种解析场景和正确的错误处理技术。
