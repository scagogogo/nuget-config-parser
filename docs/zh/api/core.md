# 核心 API

`pkg/nuget` 包为 NuGet Config Parser 库提供主要的 API 接口。这是大多数操作的主要入口点。

## 概述

核心 API 旨在为所有 NuGet 配置操作提供简单、统一的接口。它抽象了底层组件的复杂性，为开发者提供了清洁、易用的 API。

## API 结构

```go
type API struct {
    Parser  *parser.ConfigParser
    Finder  *finder.ConfigFinder
    Manager *manager.ConfigManager
}
```

API 结构体集成了三个核心组件：
- **Parser**: 处理配置文件解析和序列化
- **Finder**: 跨不同平台定位配置文件
- **Manager**: 管理配置修改和操作

## 构造函数

### NewAPI

```go
func NewAPI() *API
```

使用默认设置创建新的 API 实例。

**返回值:**
- `*API`: 准备使用的新 API 实例

**示例:**
```go
api := nuget.NewAPI()
```

## 解析和查找

### ParseFromFile

```go
func (a *API) ParseFromFile(filePath string) (*types.NuGetConfig, error)
```

从指定路径解析 NuGet 配置文件。

**参数:**
- `filePath` (string): 配置文件的路径

**返回值:**
- `*types.NuGetConfig`: 解析的配置对象
- `error`: 解析失败时的错误

**示例:**
```go
config, err := api.ParseFromFile("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("解析配置失败: %v", err)
}
```

### ParseFromString

```go
func (a *API) ParseFromString(content string) (*types.NuGetConfig, error)
```

从 XML 字符串解析 NuGet 配置。

**参数:**
- `content` (string): XML 内容字符串

**返回值:**
- `*types.NuGetConfig`: 解析的配置对象
- `error`: 解析失败时的错误

**示例:**
```go
xmlContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`

config, err := api.ParseFromString(xmlContent)
```

### ParseFromReader

```go
func (a *API) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)
```

从 io.Reader 解析 NuGet 配置。

**参数:**
- `reader` (io.Reader): 包含 XML 内容的 Reader

**返回值:**
- `*types.NuGetConfig`: 解析的配置对象
- `error`: 解析失败时的错误

**示例:**
```go
file, err := os.Open("NuGet.Config")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

config, err := api.ParseFromReader(file)
```

### FindConfigFile

```go
func (a *API) FindConfigFile() (string, error)
```

查找系统中第一个可用的 NuGet 配置文件。

**返回值:**
- `string`: 找到的配置文件路径
- `error`: 未找到配置文件时的错误

**示例:**
```go
configPath, err := api.FindConfigFile()
if err != nil {
    log.Fatalf("未找到配置文件: %v", err)
}
fmt.Printf("找到配置: %s\n", configPath)
```

### FindAllConfigFiles

```go
func (a *API) FindAllConfigFiles() []string
```

查找系统中所有可用的 NuGet 配置文件。

**返回值:**
- `[]string`: 所有找到的配置文件路径切片

**示例:**
```go
configPaths := api.FindAllConfigFiles()
fmt.Printf("找到 %d 个配置文件:\n", len(configPaths))
for _, path := range configPaths {
    fmt.Printf("  - %s\n", path)
}
```

### FindAndParseConfig

```go
func (a *API) FindAndParseConfig() (*types.NuGetConfig, string, error)
```

查找并解析第一个可用的配置文件。

**返回值:**
- `*types.NuGetConfig`: 解析的配置对象
- `string`: 被解析的配置文件路径
- `error`: 未找到配置文件或解析失败时的错误

**示例:**
```go
config, configPath, err := api.FindAndParseConfig()
if err != nil {
    log.Fatalf("查找和解析配置失败: %v", err)
}
fmt.Printf("从以下位置加载配置: %s\n", configPath)
```

## 包源管理

### AddPackageSource

```go
func (a *API) AddPackageSource(config *types.NuGetConfig, key, value, protocolVersion string)
```

在配置中添加或更新包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 包源的唯一标识符
- `value` (string): 包源的 URL 或路径
- `protocolVersion` (string): 协议版本（可选，可以为空）

**示例:**
```go
api.AddPackageSource(config, "myFeed", "https://my-nuget-feed.com/v3/index.json", "3")
```

### RemovePackageSource

```go
func (a *API) RemovePackageSource(config *types.NuGetConfig, key string) bool
```

从配置中移除包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 要移除的包源的键

**返回值:**
- `bool`: 如果找到并移除了源则为 true，否则为 false

**示例:**
```go
removed := api.RemovePackageSource(config, "myFeed")
if removed {
    fmt.Println("包源移除成功")
}
```

### GetPackageSource

```go
func (a *API) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource
```

通过键检索特定的包源。

**参数:**
- `config` (*types.NuGetConfig): 要搜索的配置对象
- `key` (string): 要检索的包源的键

**返回值:**
- `*types.PackageSource`: 如果找到则返回包源，否则为 nil

**示例:**
```go
source := api.GetPackageSource(config, "nuget.org")
if source != nil {
    fmt.Printf("源 URL: %s\n", source.Value)
}
```

## 包源状态管理

### EnablePackageSource

```go
func (a *API) EnablePackageSource(config *types.NuGetConfig, key string)
```

通过从禁用源列表中移除来启用包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 要启用的包源的键

### DisablePackageSource

```go
func (a *API) DisablePackageSource(config *types.NuGetConfig, key string)
```

通过将包源添加到禁用源列表来禁用包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 要禁用的包源的键

### IsPackageSourceDisabled

```go
func (a *API) IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool
```

检查包源是否被禁用。

**参数:**
- `config` (*types.NuGetConfig): 要检查的配置对象
- `key` (string): 要检查的包源的键

**返回值:**
- `bool`: 如果源被禁用则为 true，否则为 false

**示例:**
```go
if api.IsPackageSourceDisabled(config, "myFeed") {
    fmt.Println("myFeed 已禁用")
    api.EnablePackageSource(config, "myFeed")
}
```

## 凭证管理

### AddCredential

```go
func (a *API) AddCredential(config *types.NuGetConfig, sourceKey, username, password string)
```

为包源添加凭证。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `sourceKey` (string): 包源的键
- `username` (string): 认证用户名
- `password` (string): 认证密码

**示例:**
```go
api.AddCredential(config, "privateFeed", "myuser", "mypassword")
```

### RemoveCredential

```go
func (a *API) RemoveCredential(config *types.NuGetConfig, sourceKey string) bool
```

移除包源的凭证。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `sourceKey` (string): 包源的键

**返回值:**
- `bool`: 如果找到并移除了凭证则为 true，否则为 false

### GetCredential

```go
func (a *API) GetCredential(config *types.NuGetConfig, sourceKey string) *types.SourceCredential
```

检索包源的凭证。

**参数:**
- `config` (*types.NuGetConfig): 要搜索的配置对象
- `sourceKey` (string): 包源的键

**返回值:**
- `*types.SourceCredential`: 如果找到则返回凭证，否则为 nil

## 配置选项

### AddConfigOption

```go
func (a *API) AddConfigOption(config *types.NuGetConfig, key, value string)
```

添加或更新全局配置选项。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 配置选项键
- `value` (string): 配置选项值

**示例:**
```go
api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages/path")
```

### GetConfigOption

```go
func (a *API) GetConfigOption(config *types.NuGetConfig, key string) string
```

检索全局配置选项值。

**参数:**
- `config` (*types.NuGetConfig): 要搜索的配置对象
- `key` (string): 配置选项键

**返回值:**
- `string`: 配置选项值，如果未找到则为空字符串

## 活跃包源

### SetActivePackageSource

```go
func (a *API) SetActivePackageSource(config *types.NuGetConfig, key, value string)
```

设置活跃包源。

**参数:**
- `config` (*types.NuGetConfig): 要修改的配置对象
- `key` (string): 活跃包源的键
- `value` (string): 活跃包源的 URL

### GetActivePackageSource

```go
func (a *API) GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource
```

检索活跃包源。

**参数:**
- `config` (*types.NuGetConfig): 要搜索的配置对象

**返回值:**
- `*types.PackageSource`: 如果设置了则返回活跃包源，否则为 nil

## 序列化和持久化

### SaveConfig

```go
func (a *API) SaveConfig(config *types.NuGetConfig, filePath string) error
```

将配置对象保存到文件。

**参数:**
- `config` (*types.NuGetConfig): 要保存的配置对象
- `filePath` (string): 保存配置文件的路径

**返回值:**
- `error`: 保存失败时的错误

**示例:**
```go
err := api.SaveConfig(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("保存配置失败: %v", err)
}
```

### SerializeToXML

```go
func (a *API) SerializeToXML(config *types.NuGetConfig) (string, error)
```

将配置对象序列化为 XML 字符串。

**参数:**
- `config` (*types.NuGetConfig): 要序列化的配置对象

**返回值:**
- `string`: 配置的 XML 表示
- `error`: 序列化失败时的错误

**示例:**
```go
xmlString, err := api.SerializeToXML(config)
if err != nil {
    log.Fatalf("序列化配置失败: %v", err)
}
fmt.Println(xmlString)
```

## 配置创建

### CreateDefaultConfig

```go
func (a *API) CreateDefaultConfig() *types.NuGetConfig
```

创建具有默认设置的新配置。

**返回值:**
- `*types.NuGetConfig`: 具有默认包源的新配置

**示例:**
```go
config := api.CreateDefaultConfig()
// config 现在包含 nuget.org 作为默认源
```

### InitializeDefaultConfig

```go
func (a *API) InitializeDefaultConfig(filePath string) error
```

创建并保存默认配置到指定路径。

**参数:**
- `filePath` (string): 创建配置文件的路径

**返回值:**
- `error`: 创建或保存失败时的错误

**示例:**
```go
err := api.InitializeDefaultConfig("/path/to/new/NuGet.Config")
if err != nil {
    log.Fatalf("初始化配置失败: %v", err)
}
```

## 位置感知编辑

### CreateConfigEditor

```go
func (a *API) CreateConfigEditor(parseResult *parser.ParseResult) *editor.ConfigEditor
```

创建位置感知配置编辑器。

**参数:**
- `parseResult` (*parser.ParseResult): 带有位置信息的解析结果

**返回值:**
- `*editor.ConfigEditor`: 配置编辑器实例

**示例:**
```go
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatal(err)
}

editor := api.CreateConfigEditor(parseResult)
err = editor.AddPackageSource("newSource", "https://example.com", "3")
if err != nil {
    log.Fatal(err)
}

modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatal(err)
}

// 保存修改的内容
err = os.WriteFile("/path/to/NuGet.Config", modifiedContent, 0644)
```

## 错误处理

所有可能失败的方法都将错误作为最后一个返回值。使用 `pkg/errors` 包中的错误处理工具：

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("config.xml")
if err != nil {
    if errors.IsNotFoundError(err) {
        // 处理文件未找到
    } else if errors.IsParseError(err) {
        // 处理解析错误
    } else {
        // 处理其他错误
    }
}
```

## 最佳实践

1. **重用 API 实例**: 创建一个 API 实例并在整个应用程序中重用
2. **检查错误**: 始终检查并适当处理错误
3. **使用位置感知编辑**: 对于最小文件更改，使用位置感知编辑功能
4. **验证输入**: 确保包源键是唯一的，URL 是有效的
5. **处理缺失文件**: 使用 `FindConfigFile()` 或在文件不存在时创建默认配置
6. **设置活跃源**: 考虑设置活跃包源以获得更好的用户体验
