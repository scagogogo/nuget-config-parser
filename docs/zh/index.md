---
layout: home

hero:
  name: "NuGet 配置解析器"
  text: "NuGet 配置文件的 Go 库"
  tagline: "轻松解析、操作和管理 NuGet 配置文件"
  image:
    src: /logo.svg
    alt: NuGet Config Parser
  actions:
    - theme: brand
      text: 开始使用
      link: /zh/guide/getting-started
    - theme: alt
      text: 查看 GitHub
      link: https://github.com/scagogogo/nuget-config-parser

features:
  - icon: 📄
    title: 配置文件解析
    details: 从文件、字符串或 io.Reader 解析 NuGet.Config 文件，具有全面的错误处理。
  
  - icon: 🔍
    title: 智能文件发现
    details: 自动在系统中查找 NuGet 配置文件，支持项目级和全局配置。
  
  - icon: 📦
    title: 包源管理
    details: 添加、移除、启用/禁用包源，完全支持协议版本和凭证。
  
  - icon: 🔐
    title: 凭证管理
    details: 安全管理私有包源的用户名/密码凭证。
  
  - icon: ⚙️
    title: 配置选项
    details: 管理全局配置选项，如代理设置、包文件夹路径等。
  
  - icon: ✏️
    title: 位置感知编辑
    details: 编辑配置文件时保持原始格式并最小化差异。
  
  - icon: 🔄
    title: 序列化支持
    details: 将配置对象转换为标准 XML 格式，具有适当的缩进。
  
  - icon: 🌐
    title: 跨平台
    details: 完全支持 Windows、Linux 和 macOS，具有平台特定的配置路径。
  
  - icon: 🧪
    title: 全面测试
    details: 经过广泛测试，具有高代码覆盖率和真实场景验证。
---

## 快速示例

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
    
    // 查找并解析配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("查找配置失败: %v", err)
    }
    
    // 显示包源
    fmt.Printf("配置文件: %s\n", configPath)
    fmt.Printf("包源数量: %d\n", len(config.PackageSources.Add))
    
    for _, source := range config.PackageSources.Add {
        fmt.Printf("- %s: %s\n", source.Key, source.Value)
    }
}
```

## 安装

```bash
go get github.com/scagogogo/nuget-config-parser
```

## 主要特性

### 🚀 易于使用
简单直观的 API，遵循 Go 最佳实践和约定。

### 🔧 功能全面
支持所有主要的 NuGet 配置功能，包括包源、凭证和全局设置。

### 📝 文档完善
为每个功能和用例提供详尽的文档和示例。

### 🎯 生产就绪
经过实战测试，具有全面的测试覆盖率和错误处理。
