# 快速开始

本指南将通过实际示例帮助您快速上手 NuGet Config Parser 库。

## 前提条件

- Go 1.19 或更高版本
- 对 Go 编程的基本了解
- 对 NuGet 配置文件的熟悉（有帮助但非必需）

## 安装

首先，将库添加到您的 Go 项目：

```bash
go get github.com/scagogogo/nuget-config-parser
```

## 基本用法

### 1. 创建您的第一个程序

创建一个新的 Go 文件并导入库：

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
    
    // 您的代码在这里
}
```

### 2. 查找和解析配置

最常见的操作是查找和解析现有配置：

```go
func main() {
    api := nuget.NewAPI()
    
    // 一步完成查找和解析
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("查找/解析配置失败: %v", err)
    }
    
    fmt.Printf("从以下位置加载配置: %s\n", configPath)
    fmt.Printf("包源: %d\n", len(config.PackageSources.Add))
}
```

### 3. 创建新配置

如果不存在配置，创建默认配置：

```go
func main() {
    api := nuget.NewAPI()
    
    // 尝试查找现有配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // 未找到配置，创建默认配置
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        
        // 保存新配置
        err = api.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("保存配置失败: %v", err)
        }
        
        fmt.Printf("创建了新配置: %s\n", configPath)
    } else {
        fmt.Printf("使用现有配置: %s\n", configPath)
    }
}
```

### 4. 管理包源

添加、移除和管理包源：

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // 添加自定义包源
    api.AddPackageSource(config, "company-feed", "https://nuget.company.com/v3/index.json", "3")
    
    // 添加本地包源
    api.AddPackageSource(config, "local-packages", "/path/to/local/packages", "")
    
    // 禁用包源
    api.DisablePackageSource(config, "local-packages")
    
    // 列出所有包源
    fmt.Println("包源:")
    for _, source := range config.PackageSources.Add {
        status := "已启用"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "已禁用"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // 保存更改
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
}
```

### 5. 管理凭证

为私有包源添加凭证：

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // 添加私有包源
    api.AddPackageSource(config, "private-feed", "https://private.nuget.com/v3/index.json", "3")
    
    // 为私有源添加凭证
    api.AddCredential(config, "private-feed", "myusername", "mypassword")
    
    // 验证凭证已添加
    credential := api.GetCredential(config, "private-feed")
    if credential != nil {
        fmt.Println("凭证添加成功")
    }
    
    // 保存配置
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
}
```

### 6. 配置选项

管理全局 NuGet 设置：

```go
func main() {
    api := nuget.NewAPI()
    config := api.CreateDefaultConfig()
    
    // 设置全局包文件夹
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages/path")
    
    // 设置默认推送源
    api.AddConfigOption(config, "defaultPushSource", "https://my-nuget-server.com")
    
    // 设置代理设置
    api.AddConfigOption(config, "http_proxy", "http://proxy.company.com:8080")
    api.AddConfigOption(config, "http_proxy.user", "proxyuser")
    api.AddConfigOption(config, "http_proxy.password", "proxypass")
    
    // 获取配置选项
    packagesPath := api.GetConfigOption(config, "globalPackagesFolder")
    fmt.Printf("全局包文件夹: %s\n", packagesPath)
    
    // 保存配置
    err := api.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
}
```

## 完整示例

这是一个演示多个功能的完整示例：

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
    
    // 尝试查找现有配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // 如果不存在则创建新配置
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
        fmt.Println("创建了新配置")
    } else {
        fmt.Printf("找到现有配置: %s\n", configPath)
    }
    
    // 显示当前包源
    fmt.Printf("\n当前包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "已启用"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "已禁用"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // 添加公司包源
    fmt.Println("\n添加公司包源...")
    api.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    api.AddCredential(config, "company", "employee", "secret123")
    
    // 添加本地开发源
    fmt.Println("添加本地开发源...")
    api.AddPackageSource(config, "local-dev", "/tmp/local-packages", "")
    api.DisablePackageSource(config, "local-dev") // 默认禁用
    
    // 配置全局设置
    fmt.Println("配置全局设置...")
    api.AddConfigOption(config, "globalPackagesFolder", os.ExpandEnv("$HOME/.nuget/packages"))
    api.AddConfigOption(config, "repositoryPath", "./packages")
    
    // 设置活跃包源
    api.SetActivePackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json")
    
    // 显示更新的配置
    fmt.Printf("\n更新的包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        status := "已启用"
        if api.IsPackageSourceDisabled(config, source.Key) {
            status = "已禁用"
        }
        fmt.Printf("  - %s: %s (%s)\n", source.Key, source.Value, status)
    }
    
    // 显示活跃源
    if activeSource := api.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("\n活跃源: %s\n", activeSource.Key)
    }
    
    // 保存配置
    fmt.Printf("\n保存配置到: %s\n", configPath)
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("配置保存成功！")
    
    // 可选：显示 XML 内容
    xmlContent, err := api.SerializeToXML(config)
    if err != nil {
        log.Printf("序列化配置失败: %v", err)
    } else {
        fmt.Println("\n生成的 XML:")
        fmt.Println(xmlContent)
    }
}
```

## 错误处理

始终适当处理错误：

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("NuGet.Config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // 文件不存在
        config = api.CreateDefaultConfig()
    } else if errors.IsParseError(err) {
        // 无效的 XML 或格式
        log.Fatalf("无效的配置格式: %v", err)
    } else {
        // 其他错误
        log.Fatalf("意外错误: %v", err)
    }
}
```

## 下一步

现在您已经学会了基础知识：

1. 探索 [位置感知编辑](./position-aware-editing.md) 以获得高级编辑功能
2. 查看 [示例](/zh/examples/) 了解更多特定用例
3. 阅读 [API 参考](/zh/api/) 获取完整文档
4. 了解 [配置](./configuration.md) 结构和选项

## 常见模式

### 配置发现模式

```go
// 尝试多种方法查找配置
var config *types.NuGetConfig
var configPath string

// 1. 项目特定配置
if projectConfig, err := api.FindProjectConfig("."); err == nil {
    if config, err = api.ParseFromFile(projectConfig); err == nil {
        configPath = projectConfig
    }
}

// 2. 全局配置
if config == nil {
    if globalConfig, err := api.FindConfigFile(); err == nil {
        if config, err = api.ParseFromFile(globalConfig); err == nil {
            configPath = globalConfig
        }
    }
}

// 3. 默认配置
if config == nil {
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
}
```

### 批量操作模式

```go
// 高效地进行多个更改
api.AddPackageSource(config, "feed1", "https://feed1.com", "3")
api.AddPackageSource(config, "feed2", "https://feed2.com", "3")
api.AddCredential(config, "feed1", "user1", "pass1")
api.DisablePackageSource(config, "old-feed")

// 一次性保存所有更改
err := api.SaveConfig(config, configPath)
```

本快速开始指南应该让您快速上手该库。有关更高级的功能和详细说明，请继续阅读其他文档部分。
