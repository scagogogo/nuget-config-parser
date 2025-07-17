# 开始使用

欢迎使用 NuGet Config Parser！本指南将帮助您快速上手这个库。

## 什么是 NuGet Config Parser？

NuGet Config Parser 是一个 Go 库，提供了解析和操作 NuGet 配置文件 (NuGet.Config) 的全面功能。它允许您：

- 解析现有的 NuGet 配置文件
- 以编程方式创建新配置
- 修改包源、凭证和设置
- 在系统中查找配置文件
- 编辑文件时保持原始格式

## 前提条件

- Go 1.19 或更高版本
- 对 NuGet 配置文件的基本了解

## 安装

将库添加到您的 Go 项目：

```bash
go get github.com/scagogogo/nuget-config-parser
```

## 第一个程序

让我们创建一个简单的程序来查找和显示 NuGet 配置信息：

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
    
    // 查找第一个可用的配置文件
    configPath, err := api.FindConfigFile()
    if err != nil {
        log.Fatalf("未找到配置文件: %v", err)
    }
    
    // 解析配置文件
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 显示基本信息
    fmt.Printf("配置文件: %s\n", configPath)
    fmt.Printf("包源数量: %d\n", len(config.PackageSources.Add))
    
    // 列出所有包源
    fmt.Println("\n包源:")
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
        
        // 检查包源是否被禁用
        if api.IsPackageSourceDisabled(config, source.Key) {
            fmt.Printf("    状态: 已禁用\n")
        } else {
            fmt.Printf("    状态: 已启用\n")
        }
    }
    
    // 显示活跃包源
    if config.ActivePackageSource != nil {
        fmt.Printf("\n活跃包源: %s\n", config.ActivePackageSource.Add.Key)
    }
}
```

## 核心概念

### API 实例

主要入口点是 `API` 结构体，它提供了您需要的所有功能：

```go
api := nuget.NewAPI()
```

### 配置对象

`NuGetConfig` 结构体表示完整的 NuGet 配置：

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

### 包源

包源是 NuGet 配置的核心。每个源都有：

- **Key**: 包源的唯一标识符
- **Value**: 包源的 URL 或路径
- **ProtocolVersion**: NuGet 协议版本（可选）

## 常见操作

### 查找配置文件

```go
// 查找第一个可用的配置文件
configPath, err := api.FindConfigFile()

// 查找所有配置文件
configPaths := api.FindAllConfigFiles()

// 查找项目特定配置
projectConfig, err := api.FindProjectConfig("./my-project")
```

### 解析配置

```go
// 从文件解析
config, err := api.ParseFromFile("/path/to/NuGet.Config")

// 从字符串解析
config, err := api.ParseFromString(xmlContent)

// 从 io.Reader 解析
config, err := api.ParseFromReader(reader)
```

### 修改配置

```go
// 添加包源
api.AddPackageSource(config, "mySource", "https://my-nuget-feed.com/v3/index.json", "3")

// 移除包源
removed := api.RemovePackageSource(config, "mySource")

// 禁用包源
api.DisablePackageSource(config, "mySource")

// 添加凭证
api.AddCredential(config, "mySource", "username", "password")
```

### 保存配置

```go
// 保存到文件
err := api.SaveConfig(config, "/path/to/NuGet.Config")

// 序列化为 XML 字符串
xmlString, err := api.SerializeToXML(config)
```

## 下一步

现在您了解了基础知识，请探索这些主题：

- [安装指南](./installation.md) - 详细的安装说明
- [快速开始](./quick-start.md) - 更全面的示例
- [配置](./configuration.md) - 了解 NuGet 配置结构
- [位置感知编辑](./position-aware-editing.md) - 高级编辑功能
- [API 参考](/zh/api/) - 完整的 API 文档
- [示例](/zh/examples/) - 实际使用示例

## 需要帮助？

- 查看 [API 参考](/zh/api/) 获取详细的方法文档
- 浏览 [示例](/zh/examples/) 了解常见用例
- 访问 [GitHub 仓库](https://github.com/scagogogo/nuget-config-parser) 提出问题和讨论
