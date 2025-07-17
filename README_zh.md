# NuGet Config Parser

> **📖 [完整文档和 API 参考](https://scagogogo.github.io/nuget-config-parser/)** | **[🇺🇸 English](README.md)**

[![Go CI](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/ci.yml)
[![Scheduled Tests](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml/badge.svg)](https://github.com/scagogogo/nuget-config-parser/actions/workflows/scheduled-tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scagogogo/nuget-config-parser)](https://goreportcard.com/report/github.com/scagogogo/nuget-config-parser)
[![GoDoc](https://godoc.org/github.com/scagogogo/nuget-config-parser?status.svg)](https://godoc.org/github.com/scagogogo/nuget-config-parser)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](https://scagogogo.github.io/nuget-config-parser/)

这个库提供了解析和操作 NuGet 配置文件 (NuGet.Config) 的功能。它可以帮助你在 Go 应用程序中读取、修改和创建 NuGet 配置文件，支持所有主要的 NuGet 配置功能。

## 📚 文档

### 🌐 **[在线文档](https://scagogogo.github.io/nuget-config-parser/)**

完整的文档可在线访问：**https://scagogogo.github.io/nuget-config-parser/**

文档包括：
- **📖 [入门指南](https://scagogogo.github.io/nuget-config-parser/zh/guide/getting-started)** - 逐步介绍
- **🔧 [API 参考](https://scagogogo.github.io/nuget-config-parser/zh/api/)** - 完整的 API 文档和示例
- **💡 [使用示例](https://scagogogo.github.io/nuget-config-parser/zh/examples/)** - 真实世界的使用示例
- **⚡ [最佳实践](https://scagogogo.github.io/nuget-config-parser/zh/guide/configuration)** - 推荐的模式和做法
- **🌍 多语言支持** - 提供中文和英文版本

## 📑 目录

- [功能特点](#功能特点)
- [安装](#安装)
- [快速开始](#快速开始)
- [示例](#示例)
- [API 参考](#api-参考)
- [架构](#架构)
- [贡献](#贡献)
- [许可证](#许可证)

## ✨ 功能特点

- **配置文件解析** - 解析 NuGet.Config 文件，支持从文件、字符串或 Reader 读取
- **配置文件查找** - 查找系统中的 NuGet 配置文件，支持项目级和全局配置
- **包源管理** - 添加、移除、获取包源信息
- **凭证管理** - 设置和管理包源的用户名/密码凭证
- **包源启用与禁用** - 启用/禁用包源
- **活跃包源管理** - 设置和获取活跃包源
- **配置选项管理** - 管理全局配置选项，如代理设置、包文件夹路径等
- **配置序列化** - 将配置对象序列化为标准 XML 格式
- **位置感知编辑** - 基于位置信息的精确编辑，保持原始格式，最小化diff
- **跨平台支持** - 支持 Windows、Linux 和 macOS

## 🚀 安装

使用 Go 模块安装（推荐）：

```bash
go get github.com/scagogogo/nuget-config-parser
```

## 🏁 快速开始

> 💡 **详细教程和示例请访问 [快速开始指南](https://scagogogo.github.io/nuget-config-parser/zh/guide/quick-start)**

以下是一个简单的示例，演示如何解析和使用 NuGet 配置文件：

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
        log.Fatalf("找不到配置文件: %v", err)
    }
    
    // 解析配置文件
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 显示配置信息
    fmt.Printf("配置文件: %s\n", configPath)
    fmt.Printf("包含 %d 个包源\n", len(config.PackageSources.Add))
    
    // 显示包源列表
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
        
        // 检查包源是否禁用
        if api.IsPackageSourceDisabled(config, source.Key) {
            fmt.Printf("  状态: 已禁用\n")
        } else {
            fmt.Printf("  状态: 已启用\n")
        }
    }
}
```

## 📝 示例

> 🔗 **更多示例和详细说明请访问 [示例文档](https://scagogogo.github.io/nuget-config-parser/zh/examples/)**

本项目提供了多个完整示例，展示不同的功能和用例。所有示例都位于 [examples](examples/) 目录中：

1. **[基本解析](examples/01_basic_parsing)** - 解析配置文件并访问其内容
2. **[查找配置](examples/02_search_config)** - 在系统中查找 NuGet 配置文件
3. **[创建配置](examples/03_create_config)** - 创建新的 NuGet 配置
4. **[修改配置](examples/04_modify_config)** - 修改现有的 NuGet 配置
5. **[包源管理](examples/05_package_sources)** - 包源相关操作
6. **[凭证管理](examples/06_credentials)** - 管理包源凭证
7. **[配置选项](examples/07_config_options)** - 管理全局配置选项
8. **[序列化](examples/08_serialization)** - 配置序列化和反序列化
9. **[位置感知编辑](examples/09_position_aware_editing)** - 基于位置信息的精确编辑

运行示例：

```bash
go run examples/01_basic_parsing/main.go
```

有关示例的详细说明，请参阅 [examples/README.md](examples/README.md)。

## 📚 API 参考

> 📖 **完整的 API 文档和示例：[API 参考](https://scagogogo.github.io/nuget-config-parser/zh/api/)**

### 核心 API

```go
// 创建新的 API 实例
api := nuget.NewAPI()
```

### 解析和查找

```go
// 从文件解析配置
config, err := api.ParseFromFile(filePath)

// 从字符串解析配置
config, err := api.ParseFromString(xmlContent)

// 从 io.Reader 解析配置
config, err := api.ParseFromReader(reader)

// 查找第一个可用的配置文件
configPath, err := api.FindConfigFile()

// 查找所有可用的配置文件
configPaths := api.FindAllConfigFiles()

// 在项目目录中查找配置文件
projectConfig, err := api.FindProjectConfig(startDir)

// 查找并解析配置
config, configPath, err := api.FindAndParseConfig()
```

### 包源管理

```go
// 添加或更新包源
api.AddPackageSource(config, "sourceName", "https://source-url", "3")

// 移除包源
removed := api.RemovePackageSource(config, "sourceName")

// 获取特定包源
source := api.GetPackageSource(config, "sourceName")

// 获取所有包源
sources := api.GetAllPackageSources(config)

// 设置活跃包源
err := api.SetActivePackageSource(config, "sourceName")
```

### 凭证管理

```go
// 添加凭证
api.AddCredential(config, "sourceName", "username", "password")

// 移除凭证
removed := api.RemoveCredential(config, "sourceName")
```

### 包源启用/禁用

```go
// 禁用包源
api.DisablePackageSource(config, "sourceName")

// 启用包源
enabled := api.EnablePackageSource(config, "sourceName")

// 检查包源是否禁用
disabled := api.IsPackageSourceDisabled(config, "sourceName")
```

### 配置选项

```go
// 添加或更新配置选项
api.AddConfigOption(config, "globalPackagesFolder", "/path/to/packages")

// 移除配置选项
removed := api.RemoveConfigOption(config, "globalPackagesFolder")

// 获取配置选项值
value := api.GetConfigOption(config, "globalPackagesFolder")
```

### 创建和保存

```go
// 创建默认配置
config := api.CreateDefaultConfig()

// 在指定路径创建默认配置
err := api.InitializeDefaultConfig(filePath)

// 保存配置到文件
err := api.SaveConfig(config, filePath)

// 将配置序列化为 XML 字符串
xmlString, err := api.SerializeToXML(config)

// 位置感知编辑（保持原始格式，最小化diff）
parseResult, err := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)
err = editor.AddPackageSource("new-source", "https://example.com/v3/index.json", "3")
modifiedContent, err := editor.ApplyEdits()
```

## 🏗️ 架构

该库由以下主要组件组成：

- **pkg/nuget**: 主要 API 包，提供用户接口
- **pkg/parser**: 配置解析器，负责 XML 解析
- **pkg/finder**: 配置查找器，负责查找配置文件
- **pkg/manager**: 配置管理器，负责修改配置
- **pkg/types**: 数据类型定义
- **pkg/constants**: 常量定义
- **pkg/utils**: 工具函数
- **pkg/errors**: 错误类型定义

## 🤝 贡献

欢迎贡献！如果您想为这个项目做出贡献：

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

在提交 PR 前，请确保代码通过了测试并且符合代码风格规范。

## 📄 许可证

该项目采用 MIT 许可证。有关详细信息，请参阅 [LICENSE](LICENSE) 文件。