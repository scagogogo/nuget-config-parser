# 查找器 API

`pkg/finder` 包提供跨不同平台和目录结构的配置文件发现功能。

## 概述

查找器 API 负责：
- 在标准位置定位 NuGet 配置文件
- 支持平台特定的配置路径
- 提供灵活的搜索策略
- 处理项目级和全局配置发现

## 类型

### ConfigFinder

```go
type ConfigFinder struct {
    SearchPaths []string
}
```

处理配置文件发现的主要查找器类型。

**字段:**
- `SearchPaths`: 搜索配置文件的路径列表

## 构造函数

### NewConfigFinder

```go
func NewConfigFinder() *ConfigFinder
```

使用默认搜索路径创建新的配置查找器。

**返回值:**
- `*ConfigFinder`: 具有平台特定默认路径的新查找器实例

**示例:**
```go
finder := finder.NewConfigFinder()
configPath, err := finder.FindConfigFile()
if err != nil {
    log.Printf("未找到配置文件: %v", err)
} else {
    fmt.Printf("找到配置: %s\n", configPath)
}
```

### NewConfigFinderWithPaths

```go
func NewConfigFinderWithPaths(searchPaths []string) *ConfigFinder
```

使用自定义搜索路径创建新的配置查找器。

**参数:**
- `searchPaths` ([]string): 要搜索的自定义路径列表

**返回值:**
- `*ConfigFinder`: 具有指定路径的新查找器实例

**示例:**
```go
customPaths := []string{
    "/custom/path/NuGet.Config",
    "/another/path/NuGet.Config",
}
finder := finder.NewConfigFinderWithPaths(customPaths)
```

## 发现方法

### FindConfigFile

```go
func (f *ConfigFinder) FindConfigFile() (string, error)
```

查找第一个可用的 NuGet 配置文件。

**返回值:**
- `string`: 找到的配置文件路径
- `error`: 未找到配置文件时的错误

**搜索顺序:**
1. 当前目录 (`./NuGet.Config`)
2. 父目录（向上遍历目录树）
3. 用户特定配置目录
4. 系统范围配置目录

**平台特定路径:**
- **Windows**: `%APPDATA%\NuGet\NuGet.Config`, `%ProgramData%\NuGet\NuGet.Config`
- **macOS**: `~/Library/Application Support/NuGet/NuGet.Config`, `/Library/Application Support/NuGet/NuGet.Config`
- **Linux**: `~/.config/NuGet/NuGet.Config`, `/etc/NuGet/NuGet.Config`

**示例:**
```go
finder := finder.NewConfigFinder()
configPath, err := finder.FindConfigFile()
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("未找到 NuGet 配置文件")
        // 创建默认配置
    } else {
        log.Fatalf("搜索配置时出错: %v", err)
    }
} else {
    fmt.Printf("使用配置文件: %s\n", configPath)
}
```

### FindAllConfigFiles

```go
func (f *ConfigFinder) FindAllConfigFiles() []string
```

查找搜索路径中所有可用的 NuGet 配置文件。

**返回值:**
- `[]string`: 所有找到的配置文件路径列表

**示例:**
```go
finder := finder.NewConfigFinder()
configFiles := finder.FindAllConfigFiles()

fmt.Printf("找到 %d 个配置文件:\n", len(configFiles))
for i, path := range configFiles {
    fmt.Printf("%d. %s\n", i+1, path)
}

// 使用第一个或合并多个配置
if len(configFiles) > 0 {
    primaryConfig := configFiles[0]
    fmt.Printf("使用主配置: %s\n", primaryConfig)
}
```

### FindProjectConfig

```go
func (f *ConfigFinder) FindProjectConfig(startDir string) (string, error)
```

从指定目录开始查找项目级配置文件。

**参数:**
- `startDir` (string): 搜索的起始目录

**返回值:**
- `string`: 找到的项目配置文件路径
- `error`: 未找到项目配置文件时的错误

**搜索策略:**
1. 从 `startDir` 开始
2. 在当前目录中查找 `NuGet.Config`
3. 向上遍历父目录直到找到或到达根目录
4. 在找到的第一个配置文件处停止

**示例:**
```go
finder := finder.NewConfigFinder()

// 从当前目录开始查找项目配置
projectConfig, err := finder.FindProjectConfig(".")
if err != nil {
    fmt.Println("未找到项目特定配置")
} else {
    fmt.Printf("项目配置: %s\n", projectConfig)
}

// 为特定项目查找配置
projectPath := "/path/to/my/project"
projectConfig, err = finder.FindProjectConfig(projectPath)
if err != nil {
    fmt.Printf("在 %s 未找到项目配置\n", projectPath)
} else {
    fmt.Printf("项目配置: %s\n", projectConfig)
}
```

### FindGlobalConfig

```go
func (f *ConfigFinder) FindGlobalConfig() (string, error)
```

查找全局（用户级）配置文件。

**返回值:**
- `string`: 全局配置文件路径
- `error`: 未找到全局配置文件时的错误

**示例:**
```go
finder := finder.NewConfigFinder()
globalConfig, err := finder.FindGlobalConfig()
if err != nil {
    fmt.Println("未找到全局配置")
} else {
    fmt.Printf("全局配置: %s\n", globalConfig)
}
```

### FindSystemConfig

```go
func (f *ConfigFinder) FindSystemConfig() (string, error)
```

查找系统范围的配置文件。

**返回值:**
- `string`: 系统配置文件路径
- `error`: 未找到系统配置文件时的错误

**示例:**
```go
finder := finder.NewConfigFinder()
systemConfig, err := finder.FindSystemConfig()
if err != nil {
    fmt.Println("未找到系统配置")
} else {
    fmt.Printf("系统配置: %s\n", systemConfig)
}
```

## 路径管理

### GetDefaultSearchPaths

```go
func (f *ConfigFinder) GetDefaultSearchPaths() []string
```

返回当前平台的默认搜索路径。

**返回值:**
- `[]string`: 默认搜索路径列表

**示例:**
```go
finder := finder.NewConfigFinder()
paths := finder.GetDefaultSearchPaths()

fmt.Println("默认搜索路径:")
for i, path := range paths {
    fmt.Printf("%d. %s\n", i+1, path)
}
```

### AddSearchPath

```go
func (f *ConfigFinder) AddSearchPath(path string)
```

向查找器添加自定义搜索路径。

**参数:**
- `path` (string): 要添加到搜索列表的路径

**示例:**
```go
finder := finder.NewConfigFinder()
finder.AddSearchPath("/custom/config/location/NuGet.Config")
finder.AddSearchPath("/another/location/NuGet.Config")

// 现在搜索将包括自定义路径
configPath, err := finder.FindConfigFile()
```

### SetSearchPaths

```go
func (f *ConfigFinder) SetSearchPaths(paths []string)
```

设置完整的搜索路径列表，替换现有路径。

**参数:**
- `paths` ([]string): 新的搜索路径列表

**示例:**
```go
finder := finder.NewConfigFinder()

customPaths := []string{
    "./project.config",
    "/etc/nuget/global.config",
    "/usr/local/share/nuget/system.config",
}

finder.SetSearchPaths(customPaths)
```

## 工具方法

### ConfigExists

```go
func (f *ConfigFinder) ConfigExists(path string) bool
```

检查指定路径是否存在配置文件。

**参数:**
- `path` (string): 要检查的路径

**返回值:**
- `bool`: 如果配置文件存在且可读则为 true

**示例:**
```go
finder := finder.NewConfigFinder()

configPath := "/path/to/NuGet.Config"
if finder.ConfigExists(configPath) {
    fmt.Printf("配置存在: %s\n", configPath)
} else {
    fmt.Printf("配置未找到: %s\n", configPath)
}
```

### ValidateConfigFile

```go
func (f *ConfigFinder) ValidateConfigFile(path string) error
```

验证文件是否为有效的 NuGet 配置文件。

**参数:**
- `path` (string): 配置文件的路径

**返回值:**
- `error`: 如果文件无效则返回错误

**示例:**
```go
finder := finder.NewConfigFinder()

configPath := "/path/to/NuGet.Config"
err := finder.ValidateConfigFile(configPath)
if err != nil {
    fmt.Printf("无效的配置文件: %v\n", err)
} else {
    fmt.Println("配置文件有效")
}
```

## 平台差异

### Windows
- 用户配置: `%APPDATA%\NuGet\NuGet.Config`
- 系统配置: `%ProgramData%\NuGet\NuGet.Config`

### macOS
- 用户配置: `~/Library/Application Support/NuGet/NuGet.Config`
- 系统配置: `/Library/Application Support/NuGet/NuGet.Config`

### Linux/Unix
- 用户配置: `~/.config/NuGet/NuGet.Config` (遵循 `XDG_CONFIG_HOME`)
- 系统配置: `/etc/NuGet/NuGet.Config`

## 最佳实践

1. **使用默认查找器**: 从 `NewConfigFinder()` 开始以获得标准行为
2. **处理缺失文件**: 始终检查 `IsNotFoundError` 并提供回退
3. **尊重层次结构**: 当两者都存在时，使用项目配置而不是全局配置
4. **验证路径**: 在尝试解析文件之前使用 `ConfigExists()`
5. **自定义路径**: 使用 `AddSearchPath()` 添加额外位置而不是替换默认值
6. **环境变量**: 考虑环境特定的配置路径

## 线程安全

ConfigFinder 对于读取操作是线程安全的，但对于修改不是。如果需要并发修改搜索路径，请提供适当的同步。
