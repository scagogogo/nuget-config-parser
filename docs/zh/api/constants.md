# 常量 API

`pkg/constants` 包定义了在整个 NuGet Config Parser 库中使用的常量和默认值。

## 概述

常量 API 提供：
- 默认配置文件名和路径
- 标准 NuGet 协议版本
- 平台特定的配置位置
- 路径解析的辅助函数

## 文件和目录常量

### 配置文件名

```go
const (
    // DefaultNuGetConfigFilename 是默认的 NuGet 配置文件名
    DefaultNuGetConfigFilename = "NuGet.Config"
    
    // GlobalFolderName 是全局配置文件夹名
    GlobalFolderName = "NuGet"
    
    // FeedNamePrefix 是包源名称前缀
    FeedNamePrefix = "PackageSource"
)
```

**用法:**
- `DefaultNuGetConfigFilename`: NuGet 配置文件的标准名称
- `GlobalFolderName`: 全局 NuGet 配置的目录名
- `FeedNamePrefix`: 生成自动包源名称时使用的前缀

**示例:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// 创建配置文件路径
configPath := filepath.Join("/etc", constants.GlobalFolderName, constants.DefaultNuGetConfigFilename)
// 结果: "/etc/NuGet/NuGet.Config"

// 生成包源名称
sourceName := fmt.Sprintf("%s_%d", constants.FeedNamePrefix, 1)
// 结果: "PackageSource_1"
```

## 包源常量

### 默认包源

```go
const (
    // DefaultPackageSource 是默认的包源 URL
    DefaultPackageSource = "https://api.nuget.org/v3/index.json"
)
```

**用法:**
- `DefaultPackageSource`: 官方 NuGet.org 包源 URL

**示例:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// 创建默认包源
defaultSource := types.PackageSource{
    Key:   "nuget.org",
    Value: constants.DefaultPackageSource,
}
```

## 协议版本常量

### NuGet API 版本

```go
const (
    // NuGetV3APIProtocolVersion 是 NuGet V3 API 协议版本
    NuGetV3APIProtocolVersion = "3"
    
    // NuGetV2APIProtocolVersion 是 NuGet V2 API 协议版本  
    NuGetV2APIProtocolVersion = "2"
)
```

**用法:**
- `NuGetV3APIProtocolVersion`: 用于现代 NuGet V3 API 源
- `NuGetV2APIProtocolVersion`: 用于传统 NuGet V2 API 源

**示例:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// 创建 V3 API 包源
v3Source := types.PackageSource{
    Key:             "nuget.org",
    Value:           "https://api.nuget.org/v3/index.json",
    ProtocolVersion: constants.NuGetV3APIProtocolVersion,
}

// 创建 V2 API 包源
v2Source := types.PackageSource{
    Key:             "legacy-feed",
    Value:           "https://legacy.nuget.org/api/v2",
    ProtocolVersion: constants.NuGetV2APIProtocolVersion,
}
```

## 路径解析函数

### GetDefaultConfigLocations

```go
func GetDefaultConfigLocations() []string
```

返回 NuGet 配置文件的默认搜索路径，按优先级排序。

**返回值:**
- `[]string`: 默认配置文件路径列表

**搜索顺序:**
1. 当前目录: `./NuGet.Config`
2. 父目录: `../NuGet.Config`
3. 用户特定配置目录
4. 系统范围配置目录

**平台特定用户路径:**
- **Windows**: `%APPDATA%\NuGet\NuGet.Config`
- **macOS**: `~/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `~/.config/NuGet/NuGet.Config` (遵循 `XDG_CONFIG_HOME`)

**平台特定系统路径:**
- **Windows**: `%ProgramData%\NuGet\NuGet.Config`
- **macOS**: `/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `/etc/NuGet/NuGet.Config`

**示例:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/constants"

// 获取所有默认配置位置
locations := constants.GetDefaultConfigLocations()

fmt.Println("NuGet 配置文件搜索路径:")
for i, location := range locations {
    fmt.Printf("%d. %s\n", i+1, location)
}

// 查找第一个存在的配置文件
for _, location := range locations {
    if utils.FileExists(location) {
        fmt.Printf("找到配置文件: %s\n", location)
        break
    }
}
```

**示例输出 (Linux):**
```
NuGet 配置文件搜索路径:
1. NuGet.Config
2. ../NuGet.Config
3. /home/user/.config/NuGet/NuGet.Config
4. /etc/NuGet/NuGet.Config
```

**示例输出 (Windows):**
```
NuGet 配置文件搜索路径:
1. NuGet.Config
2. ..\NuGet.Config
3. C:\Users\user\AppData\Roaming\NuGet\NuGet.Config
4. C:\ProgramData\NuGet\NuGet.Config
```

**示例输出 (macOS):**
```
NuGet 配置文件搜索路径:
1. NuGet.Config
2. ../NuGet.Config
3. /Users/user/Library/Application Support/NuGet/NuGet.Config
4. /Library/Application Support/NuGet/NuGet.Config
```

## 平台特定行为

### Windows

```go
// 用户配置目录: %APPDATA%
// 系统配置目录: %ProgramData%
```

**环境变量:**
- `APPDATA`: 用户应用程序数据目录
- `ProgramData`: 系统范围应用程序数据目录

**示例路径:**
- 用户: `C:\Users\username\AppData\Roaming\NuGet\NuGet.Config`
- 系统: `C:\ProgramData\NuGet\NuGet.Config`

### macOS

```go
// 用户配置目录: ~/Library/Application Support/
// 系统配置目录: /Library/Application Support/
```

**示例路径:**
- 用户: `/Users/username/Library/Application Support/NuGet/NuGet.Config`
- 系统: `/Library/Application Support/NuGet/NuGet.Config`

### Linux/Unix

```go
// 用户配置目录: ~/.config/ (或 $XDG_CONFIG_HOME)
// 系统配置目录: /etc/
```

**环境变量:**
- `XDG_CONFIG_HOME`: 用户配置目录（回退到 `~/.config`）

**示例路径:**
- 用户: `/home/username/.config/NuGet/NuGet.Config`
- 系统: `/etc/NuGet/NuGet.Config`

## 使用示例

### 查找配置文件

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/utils"
)

func findConfigFile() (string, error) {
    locations := constants.GetDefaultConfigLocations()
    
    for _, location := range locations {
        if utils.FileExists(location) {
            return location, nil
        }
    }
    
    return "", fmt.Errorf("在默认位置未找到配置文件")
}

func main() {
    configPath, err := findConfigFile()
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        
        // 在当前目录创建默认配置
        defaultPath := constants.DefaultNuGetConfigFilename
        fmt.Printf("创建默认配置: %s\n", defaultPath)
        
        // 创建配置...
    } else {
        fmt.Printf("使用配置文件: %s\n", configPath)
    }
}
```

### 创建默认配置

```go
package main

import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func createDefaultConfig() *types.NuGetConfig {
    // 创建默认包源
    defaultSource := types.PackageSource{
        Key:             "nuget.org",
        Value:           constants.DefaultPackageSource,
        ProtocolVersion: constants.NuGetV3APIProtocolVersion,
    }
    
    return &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{defaultSource},
        },
        ActivePackageSource: &types.ActivePackageSource{
            Add: defaultSource,
        },
    }
}
```

### 平台特定配置

```go
package main

import (
    "fmt"
    "runtime"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
)

func showPlatformInfo() {
    fmt.Printf("操作系统: %s\n", runtime.GOOS)
    fmt.Printf("默认配置文件名: %s\n", constants.DefaultNuGetConfigFilename)
    fmt.Printf("全局文件夹名: %s\n", constants.GlobalFolderName)
    
    locations := constants.GetDefaultConfigLocations()
    fmt.Println("\n默认搜索位置:")
    for i, location := range locations {
        fmt.Printf("%d. %s\n", i+1, location)
    }
}
```

### 自定义配置路径

```go
package main

import (
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
)

func getCustomConfigPaths() []string {
    var paths []string
    
    // 添加标准位置
    paths = append(paths, constants.GetDefaultConfigLocations()...)
    
    // 添加自定义位置
    customLocations := []string{
        "/opt/nuget/NuGet.Config",
        "/usr/local/etc/nuget/NuGet.Config",
        filepath.Join(os.Getenv("HOME"), "custom-nuget.config"),
    }
    
    paths = append(paths, customLocations...)
    
    return paths
}
```

## 与其他包的集成

### 与 Finder 包

```go
import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/finder"
)

func createFinderWithDefaults() *finder.ConfigFinder {
    defaultPaths := constants.GetDefaultConfigLocations()
    return finder.NewConfigFinderWithPaths(defaultPaths)
}
```

### 与 Manager 包

```go
import (
    "github.com/scagogogo/nuget-config-parser/pkg/constants"
    "github.com/scagogogo/nuget-config-parser/pkg/manager"
    "github.com/scagogogo/nuget-config-parser/pkg/types"
)

func initializeWithDefaults() *types.NuGetConfig {
    mgr := manager.NewConfigManager()
    
    config := &types.NuGetConfig{
        PackageSources: types.PackageSources{
            Add: []types.PackageSource{
                {
                    Key:             "nuget.org",
                    Value:           constants.DefaultPackageSource,
                    ProtocolVersion: constants.NuGetV3APIProtocolVersion,
                },
            },
        },
    }
    
    return config
}
```

## 最佳实践

1. **使用常量保持一致性**: 始终使用预定义常量而不是硬编码值
2. **尊重平台差异**: 使用 `GetDefaultConfigLocations()` 实现跨平台兼容性
3. **检查文件存在性**: 在尝试解析配置文件之前始终验证它们是否存在
4. **处理缺失配置**: 当未找到默认配置时提供回退
5. **使用适当的协议版本**: 为现代源选择 V3，为传统兼容性选择 V2
6. **遵循命名约定**: 为配置文件和目录使用标准名称

## 线程安全

此包中的所有常量和函数都是线程安全的，可以并发使用。
