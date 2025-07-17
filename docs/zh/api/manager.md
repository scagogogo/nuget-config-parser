# 管理器 API

`pkg/manager` 包提供高级配置管理操作，结合解析、查找和修改功能。

## 概述

管理器 API 负责：
- 高级配置操作
- 结合解析器和查找器功能
- 管理配置生命周期
- 为常见操作提供便捷方法

## 类型

### ConfigManager

```go
type ConfigManager struct {
    parser *parser.ConfigParser
    finder *finder.ConfigFinder
}
```

协调配置操作的主要管理器类型。

**字段:**
- `parser`: 内部配置解析器
- `finder`: 内部配置查找器

## 构造函数

### NewConfigManager

```go
func NewConfigManager() *ConfigManager
```

使用默认设置创建新的配置管理器。

**返回值:**
- `*ConfigManager`: 新的管理器实例

**示例:**
```go
manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    log.Fatalf("加载配置失败: %v", err)
}
fmt.Printf("从以下位置加载配置: %s\n", configPath)
```

## 配置加载

### LoadConfig

```go
func (m *ConfigManager) LoadConfig(filePath string) (*types.NuGetConfig, error)
```

从指定路径加载配置文件。

**参数:**
- `filePath` (string): 配置文件的路径

**返回值:**
- `*types.NuGetConfig`: 加载的配置对象
- `error`: 加载失败时的错误

**示例:**
```go
manager := manager.NewConfigManager()
config, err := manager.LoadConfig("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("加载配置失败: %v", err)
}

fmt.Printf("加载了 %d 个包源\n", len(config.PackageSources.Add))
```

### FindAndLoadConfig

```go
func (m *ConfigManager) FindAndLoadConfig() (*types.NuGetConfig, string, error)
```

查找并加载第一个可用的配置文件。

**返回值:**
- `*types.NuGetConfig`: 加载的配置对象
- `string`: 被加载的配置文件路径
- `error`: 未找到配置或加载失败时的错误

**示例:**
```go
manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    if errors.IsNotFoundError(err) {
        // 未找到配置，创建默认配置
        config = manager.CreateDefaultConfig()
        configPath = "NuGet.Config"
        err = manager.SaveConfig(config, configPath)
        if err != nil {
            log.Fatalf("创建默认配置失败: %v", err)
        }
        fmt.Printf("创建了默认配置: %s\n", configPath)
    } else {
        log.Fatalf("加载配置失败: %v", err)
    }
} else {
    fmt.Printf("加载了现有配置: %s\n", configPath)
}
```

## 配置保存

### SaveConfig

```go
func (m *ConfigManager) SaveConfig(config *types.NuGetConfig, filePath string) error
```

将配置对象保存到指定的文件路径。

**参数:**
- `config` (*types.NuGetConfig): 要保存的配置对象
- `filePath` (string): 目标文件路径

**返回值:**
- `error`: 保存失败时的错误

**示例:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// 修改配置
manager.AddPackageSource(config, "company", "https://nuget.company.com", "3")

// 保存配置
err := manager.SaveConfig(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("保存配置失败: %v", err)
}
```

## 配置创建

### CreateDefaultConfig

```go
func (m *ConfigManager) CreateDefaultConfig() *types.NuGetConfig
```

创建具有默认设置的新配置。

**返回值:**
- `*types.NuGetConfig`: 具有默认包源的新配置

**默认配置:**
- 包源: `nuget.org` 指向 `https://api.nuget.org/v3/index.json`
- 协议版本: `3`
- 活跃包源: 设置为默认源

**示例:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

fmt.Printf("默认配置有 %d 个包源\n", len(config.PackageSources.Add))
fmt.Printf("默认源: %s -> %s\n", 
    config.PackageSources.Add[0].Key, 
    config.PackageSources.Add[0].Value)
```

### InitializeDefaultConfig

```go
func (m *ConfigManager) InitializeDefaultConfig(filePath string) error
```

创建并保存默认配置到指定路径。

**参数:**
- `filePath` (string): 创建配置文件的路径

**返回值:**
- `error`: 创建或保存失败时的错误

**功能:**
- 如果父目录不存在则创建
- 生成默认配置
- 保存到指定路径

**示例:**
```go
manager := manager.NewConfigManager()

configPath := "/path/to/new/NuGet.Config"
err := manager.InitializeDefaultConfig(configPath)
if err != nil {
    log.Fatalf("初始化配置失败: %v", err)
}

fmt.Printf("在以下位置初始化了默认配置: %s\n", configPath)
```

## 包源管理

### AddPackageSource

```go
func (m *ConfigManager) AddPackageSource(config *types.NuGetConfig, key string, value string, protocolVersion string) 
```

在配置中添加或更新包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 包源的唯一标识符
- `value` (string): 包源的 URL 或路径
- `protocolVersion` (string): 协议版本（可选，可以为空）

**行为:**
- 如果存在相同键的源，则更新现有源
- 如果不存在该键的源，则添加新源
- 协议版本是可选的，可以为空

**示例:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// 添加公司包源
manager.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")

// 添加不带协议版本的本地包源
manager.AddPackageSource(config, "local", "/path/to/local/packages", "")

// 更新现有源
manager.AddPackageSource(config, "company", "https://new-nuget.company.com/v3/index.json", "3")
```

### RemovePackageSource

```go
func (m *ConfigManager) RemovePackageSource(config *types.NuGetConfig, key string) bool
```

从配置中移除包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 要移除的包源的键

**返回值:**
- `bool`: 如果找到并移除了源则为 true，否则为 false

**示例:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

// 移除包源
removed := manager.RemovePackageSource(config, "old-source")
if removed {
    fmt.Println("包源移除成功")
    
    // 保存更新的配置
    err := manager.SaveConfig(config, "NuGet.Config")
    if err != nil {
        log.Printf("保存配置失败: %v", err)
    }
} else {
    fmt.Println("未找到包源")
}
```

### GetPackageSource

```go
func (m *ConfigManager) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource
```

通过键检索特定的包源。

**参数:**
- `config` (*types.NuGetConfig): 要搜索的配置对象
- `key` (string): 要检索的包源的键

**返回值:**
- `*types.PackageSource`: 如果找到则返回包源，否则为 nil

**示例:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

source := manager.GetPackageSource(config, "nuget.org")
if source != nil {
    fmt.Printf("源: %s -> %s\n", source.Key, source.Value)
    if source.ProtocolVersion != "" {
        fmt.Printf("协议版本: %s\n", source.ProtocolVersion)
    }
} else {
    fmt.Println("未找到源")
}
```

## 活跃包源管理

### SetActivePackageSource

```go
func (m *ConfigManager) SetActivePackageSource(config *types.NuGetConfig, key string) error
```

通过键设置活跃包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 要设置为活跃的包源的键

**返回值:**
- `error`: 如果未找到包源则返回错误

**示例:**
```go
manager := manager.NewConfigManager()
config := manager.CreateDefaultConfig()

// 添加多个源
manager.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
manager.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")

// 设置活跃源
err := manager.SetActivePackageSource(config, "company")
if err != nil {
    log.Printf("设置活跃源失败: %v", err)
} else {
    fmt.Println("活跃源设置为 'company'")
}
```

### GetActivePackageSource

```go
func (m *ConfigManager) GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource
```

获取当前活跃的包源。

**参数:**
- `config` (*types.NuGetConfig): 要查询的配置对象

**返回值:**
- `*types.PackageSource`: 如果设置了则返回活跃包源，否则为 nil

**示例:**
```go
manager := manager.NewConfigManager()
config, _, _ := manager.FindAndLoadConfig()

activeSource := manager.GetActivePackageSource(config)
if activeSource != nil {
    fmt.Printf("活跃源: %s -> %s\n", activeSource.Key, activeSource.Value)
} else {
    fmt.Println("未设置活跃源")
}
```

## 完整示例

这是一个展示各种管理器操作的综合示例：

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/manager"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    // 创建管理器
    mgr := manager.NewConfigManager()
    
    // 尝试查找并加载现有配置
    config, configPath, err := mgr.FindAndLoadConfig()
    if err != nil {
        if errors.IsNotFoundError(err) {
            // 未找到配置，创建默认配置
            fmt.Println("未找到配置，创建默认配置...")
            config = mgr.CreateDefaultConfig()
            configPath = "NuGet.Config"
        } else {
            log.Fatalf("加载配置失败: %v", err)
        }
    } else {
        fmt.Printf("从以下位置加载配置: %s\n", configPath)
    }
    
    // 显示当前包源
    fmt.Printf("\n当前包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // 添加公司包源
    fmt.Println("\n添加公司包源...")
    mgr.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    
    // 添加本地开发源
    fmt.Println("添加本地开发源...")
    mgr.AddPackageSource(config, "local-dev", "/tmp/local-packages", "")
    
    // 设置活跃包源
    fmt.Println("设置活跃包源...")
    err = mgr.SetActivePackageSource(config, "nuget.org")
    if err != nil {
        log.Printf("设置活跃源失败: %v", err)
    }
    
    // 显示更新的配置
    fmt.Printf("\n更新的包源 (%d):\n", len(config.PackageSources.Add))
    for _, source := range config.PackageSources.Add {
        fmt.Printf("  - %s: %s", source.Key, source.Value)
        if source.ProtocolVersion != "" {
            fmt.Printf(" (v%s)", source.ProtocolVersion)
        }
        fmt.Println()
    }
    
    // 显示活跃源
    if activeSource := mgr.GetActivePackageSource(config); activeSource != nil {
        fmt.Printf("\n活跃源: %s\n", activeSource.Key)
    }
    
    // 保存配置
    fmt.Printf("\n保存配置到: %s\n", configPath)
    err = mgr.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
    
    fmt.Println("配置保存成功！")
}
```

## 错误处理

管理器使用 `pkg/errors` 包中的标准错误类型：

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

manager := manager.NewConfigManager()
config, configPath, err := manager.FindAndLoadConfig()
if err != nil {
    if errors.IsNotFoundError(err) {
        // 处理缺失配置
        config = manager.CreateDefaultConfig()
        // ... 创建并保存默认配置
    } else if errors.IsParseError(err) {
        // 处理解析错误
        log.Printf("解析错误: %v", err)
    } else {
        // 处理其他错误
        log.Printf("意外错误: %v", err)
    }
}
```

## 最佳实践

1. **使用管理器进行高级操作**: 管理器为常见场景提供便捷方法
2. **处理缺失配置**: 始终检查 `IsNotFoundError` 并提供默认值
3. **修改后保存**: 记住在进行更改后保存配置
4. **验证包源**: 在添加包源之前确保 URL 有效
5. **使用有意义的键**: 为包源选择描述性键
6. **设置活跃源**: 考虑设置活跃包源以获得更好的用户体验

## 线程安全

ConfigManager 不是线程安全的。为并发使用创建单独的实例或提供适当的同步。
