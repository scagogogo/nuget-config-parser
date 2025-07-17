# API 参考

NuGet Config Parser 库提供了用于处理 NuGet 配置文件的全面 API。本节记录了所有公开接口、类型和方法。

## 包概述

该库分为几个包，每个包都有特定的用途：

| 包 | 描述 |
|---------|-------------|
| [`pkg/nuget`](./core.md) | 提供主要接口的核心 API 包 |
| [`pkg/parser`](./parser.md) | 配置文件解析功能 |
| [`pkg/editor`](./editor.md) | 位置感知配置编辑 |
| [`pkg/finder`](./finder.md) | 配置文件发现 |
| [`pkg/manager`](./manager.md) | 配置管理操作 |
| [`pkg/types`](./types.md) | 数据类型定义 |
| [`pkg/utils`](./utils.md) | 工具函数 |
| [`pkg/errors`](./errors.md) | 错误类型和处理 |
| [`pkg/constants`](./constants.md) | 常量和默认值 |

## 快速参考

### 核心 API

```go
import "github.com/scagogogo/nuget-config-parser/pkg/nuget"

// 创建 API 实例
api := nuget.NewAPI()

// 解析配置
config, err := api.ParseFromFile("/path/to/NuGet.Config")

// 查找配置文件
configPath, err := api.FindConfigFile()

// 修改配置
api.AddPackageSource(config, "source", "https://example.com", "3")

// 保存配置
err = api.SaveConfig(config, "/path/to/NuGet.Config")
```

### 位置感知编辑

```go
import "github.com/scagogogo/nuget-config-parser/pkg/editor"

// 使用位置跟踪解析
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")

// 创建编辑器
editor := api.CreateConfigEditor(parseResult)

// 进行更改
err = editor.AddPackageSource("new-source", "https://new.com", "3")

// 应用更改
modifiedContent, err := editor.ApplyEdits()
```

## 常用类型

### NuGetConfig

主要配置结构：

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

### PackageSource

表示单个包源：

```go
type PackageSource struct {
    Key             string `xml:"key,attr"`
    Value           string `xml:"value,attr"`
    ProtocolVersion string `xml:"protocolVersion,attr,omitempty"`
}
```

## 错误处理

该库提供结构化的错误处理：

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := api.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // 处理文件未找到
    } else if errors.IsParseError(err) {
        // 处理解析错误
    } else if errors.IsFormatError(err) {
        // 处理格式错误
    }
}
```

## 方法分类

### 解析方法

- `ParseFromFile(filePath string) (*types.NuGetConfig, error)`
- `ParseFromString(content string) (*types.NuGetConfig, error)`
- `ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)`
- `ParseFromFileWithPositions(filePath string) (*parser.ParseResult, error)`

### 查找方法

- `FindConfigFile() (string, error)`
- `FindAllConfigFiles() []string`
- `FindProjectConfig(startDir string) (string, error)`
- `FindAndParseConfig() (*types.NuGetConfig, string, error)`

### 包源方法

- `AddPackageSource(config *types.NuGetConfig, key, value, protocolVersion string)`
- `RemovePackageSource(config *types.NuGetConfig, key string) bool`
- `GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource`
- `GetAllPackageSources(config *types.NuGetConfig) []types.PackageSource`
- `EnablePackageSource(config *types.NuGetConfig, key string)`
- `DisablePackageSource(config *types.NuGetConfig, key string)`
- `IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool`

### 凭证方法

- `AddCredential(config *types.NuGetConfig, sourceKey, username, password string)`
- `RemoveCredential(config *types.NuGetConfig, sourceKey string) bool`
- `GetCredential(config *types.NuGetConfig, sourceKey string) *types.SourceCredential`

### 配置方法

- `AddConfigOption(config *types.NuGetConfig, key, value string)`
- `RemoveConfigOption(config *types.NuGetConfig, key string) bool`
- `GetConfigOption(config *types.NuGetConfig, key string) string`
- `SetActivePackageSource(config *types.NuGetConfig, key, value string)`
- `GetActivePackageSource(config *types.NuGetConfig) *types.PackageSource`

### 序列化方法

- `SaveConfig(config *types.NuGetConfig, filePath string) error`
- `SerializeToXML(config *types.NuGetConfig) (string, error)`
- `CreateDefaultConfig() *types.NuGetConfig`
- `InitializeDefaultConfig(filePath string) error`

### 编辑器方法

- `CreateConfigEditor(parseResult *parser.ParseResult) *editor.ConfigEditor`
- `AddPackageSource(key, value, protocolVersion string) error`
- `RemovePackageSource(sourceKey string) error`
- `UpdatePackageSourceURL(sourceKey, newURL string) error`
- `UpdatePackageSourceVersion(sourceKey, newVersion string) error`
- `ApplyEdits() ([]byte, error)`

## 最佳实践

### 错误处理

始终检查错误并适当处理：

```go
config, err := api.ParseFromFile(configPath)
if err != nil {
    if errors.IsNotFoundError(err) {
        // 创建默认配置
        config = api.CreateDefaultConfig()
    } else {
        return fmt.Errorf("解析配置失败: %w", err)
    }
}
```

### 资源管理

该库不需要显式的资源清理，但要注意文件操作：

```go
// 好的做法：使用 API 方法
err := api.SaveConfig(config, configPath)

// 避免：当有 API 方法可用时进行手动文件操作
```

### 性能

- 尽可能重用 API 实例
- 缓存解析的配置以供重复访问
- 使用位置感知编辑以最小化文件更改

## 线程安全

该库在设计上不是线程安全的。如果需要在并发场景中使用：

- 为每个 goroutine 创建单独的 API 实例
- 使用适当的同步机制
- 避免在没有适当锁定的情况下在 goroutine 之间共享配置对象

## 下一步

探索每个包的详细文档：

- [核心 API](./core.md) - 主要 API 接口
- [解析器](./parser.md) - 配置解析
- [编辑器](./editor.md) - 位置感知编辑
- [类型](./types.md) - 数据结构
- [示例](/zh/examples/) - 使用示例
