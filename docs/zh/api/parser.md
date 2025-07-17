# 解析器 API

`pkg/parser` 包提供配置文件解析功能，支持位置跟踪和详细的错误报告。

## 概述

解析器 API 负责：
- 从各种来源解析 NuGet 配置文件
- 跟踪元素位置以进行位置感知编辑
- 提供详细的错误信息用于调试
- 验证配置文件结构

## 类型

### ConfigParser

```go
type ConfigParser struct {
    DefaultConfigSearchPaths []string
    TrackPositions          bool
}
```

处理配置文件解析的主要解析器类型。

**字段:**
- `DefaultConfigSearchPaths`: 搜索配置文件的默认路径列表
- `TrackPositions`: 解析期间是否跟踪元素位置

### ParseResult

```go
type ParseResult struct {
    Config    *types.NuGetConfig          // 解析的配置
    Positions map[string]*ElementPosition // 元素位置信息
    Content   []byte                      // 原始文件内容
}
```

包含位置感知解析的结果。

**字段:**
- `Config`: 解析的 NuGet 配置对象
- `Positions`: 元素路径到其位置信息的映射
- `Content`: 原始文件内容字节

### ElementPosition

```go
type ElementPosition struct {
    TagName    string            // XML 标签名
    Attributes map[string]string // 元素属性
    Range      Range             // 文件中的元素范围
    AttrRanges map[string]Range  // 属性值范围
    Content    string            // 元素内容
    SelfClose  bool              // 是否为自闭合标签
}
```

表示 XML 元素的位置和元数据。

### Range

```go
type Range struct {
    Start Position // 开始位置
    End   Position // 结束位置
}
```

表示文件中的范围。

### Position

```go
type Position struct {
    Line   int // 行号（从1开始）
    Column int // 列号（从1开始）
    Offset int // 从文件开始的字节偏移量
}
```

表示文件中的特定位置。

## 构造函数

### NewConfigParser

```go
func NewConfigParser() *ConfigParser
```

使用默认设置创建新的配置解析器。

**返回值:**
- `*ConfigParser`: 新的解析器实例

**示例:**
```go
parser := parser.NewConfigParser()
config, err := parser.ParseFromFile("/path/to/NuGet.Config")
```

### NewPositionAwareParser

```go
func NewPositionAwareParser() *ConfigParser
```

创建启用位置跟踪的新解析器。

**返回值:**
- `*ConfigParser`: 启用位置跟踪的新解析器实例

**示例:**
```go
parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromFileWithPositions("/path/to/NuGet.Config")
```

## 解析方法

### ParseFromFile

```go
func (p *ConfigParser) ParseFromFile(filePath string) (*types.NuGetConfig, error)
```

从指定路径解析 NuGet 配置文件。

**参数:**
- `filePath` (string): 配置文件的路径

**返回值:**
- `*types.NuGetConfig`: 解析的配置对象
- `error`: 解析失败时的错误

**错误:**
- `errors.ErrConfigFileNotFound`: 文件不存在
- `errors.ErrEmptyConfigFile`: 文件为空
- `errors.ErrInvalidConfigFormat`: 无效的 XML 格式
- `*errors.ParseError`: 带有位置信息的详细解析错误

**示例:**
```go
parser := parser.NewConfigParser()
config, err := parser.ParseFromFile("/path/to/NuGet.Config")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("未找到配置文件")
    } else if errors.IsParseError(err) {
        fmt.Printf("解析错误: %v\n", err)
    }
    return
}

fmt.Printf("加载了 %d 个包源\n", len(config.PackageSources.Add))
```

### ParseFromString

```go
func (p *ConfigParser) ParseFromString(content string) (*types.NuGetConfig, error)
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

parser := parser.NewConfigParser()
config, err := parser.ParseFromString(xmlContent)
if err != nil {
    log.Fatalf("解析 XML 失败: %v", err)
}
```

### ParseFromReader

```go
func (p *ConfigParser) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error)
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

parser := parser.NewConfigParser()
config, err := parser.ParseFromReader(file)
if err != nil {
    log.Fatalf("从 reader 解析失败: %v", err)
}
```

## 位置感知解析

### ParseFromFileWithPositions

```go
func (p *ConfigParser) ParseFromFileWithPositions(filePath string) (*ParseResult, error)
```

解析配置文件的同时跟踪元素位置。

**参数:**
- `filePath` (string): 配置文件的路径

**返回值:**
- `*ParseResult`: 带有位置信息的解析结果
- `error`: 解析失败时的错误

**示例:**
```go
parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("带位置解析失败: %v", err)
}

// 访问配置
config := result.Config
fmt.Printf("包源: %d\n", len(config.PackageSources.Add))

// 访问位置信息
for path, pos := range result.Positions {
    fmt.Printf("元素 %s 在第 %d 行\n", path, pos.Range.Start.Line)
}
```

### ParseFromContentWithPositions

```go
func (p *ConfigParser) ParseFromContentWithPositions(content []byte) (*ParseResult, error)
```

解析配置内容的同时跟踪元素位置。

**参数:**
- `content` ([]byte): XML 内容字节

**返回值:**
- `*ParseResult`: 带有位置信息的解析结果
- `error`: 解析失败时的错误

**示例:**
```go
content := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`)

parser := parser.NewPositionAwareParser()
result, err := parser.ParseFromContentWithPositions(content)
if err != nil {
    log.Fatalf("解析内容失败: %v", err)
}
```

## 序列化方法

### SaveToFile

```go
func (p *ConfigParser) SaveToFile(config *types.NuGetConfig, filePath string) error
```

将配置对象序列化到 XML 文件。

**参数:**
- `config` (*types.NuGetConfig): 要保存的配置
- `filePath` (string): 目标文件路径

**返回值:**
- `error`: 保存失败时的错误

**示例:**
```go
parser := parser.NewConfigParser()
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:   "nuget.org",
                Value: "https://api.nuget.org/v3/index.json",
            },
        },
    },
}

err := parser.SaveToFile(config, "/path/to/NuGet.Config")
if err != nil {
    log.Fatalf("保存配置失败: %v", err)
}
```

### SerializeToXML

```go
func (p *ConfigParser) SerializeToXML(config *types.NuGetConfig) (string, error)
```

将配置对象序列化为 XML 字符串。

**参数:**
- `config` (*types.NuGetConfig): 要序列化的配置

**返回值:**
- `string`: XML 表示
- `error`: 序列化失败时的错误

**示例:**
```go
parser := parser.NewConfigParser()
xmlString, err := parser.SerializeToXML(config)
if err != nil {
    log.Fatalf("序列化失败: %v", err)
}

fmt.Println("生成的 XML:")
fmt.Println(xmlString)
```

## 发现方法

### FindAndParseConfig

```go
func (p *ConfigParser) FindAndParseConfig() (*types.NuGetConfig, string, error)
```

查找并解析第一个可用的配置文件。

**返回值:**
- `*types.NuGetConfig`: 解析的配置
- `string`: 找到的配置文件路径
- `error`: 未找到文件或解析失败时的错误

**示例:**
```go
parser := parser.NewConfigParser()
config, configPath, err := parser.FindAndParseConfig()
if err != nil {
    log.Fatalf("查找和解析配置失败: %v", err)
}

fmt.Printf("使用来自的配置: %s\n", configPath)
```

## 错误处理

解析器通过结构化错误类型提供详细的错误信息：

```go
import "github.com/scagogogo/nuget-config-parser/pkg/errors"

config, err := parser.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        // 处理文件未找到
        fmt.Println("未找到配置文件")
    } else if errors.IsParseError(err) {
        // 处理带详细信息的解析错误
        parseErr := err.(*errors.ParseError)
        fmt.Printf("第 %d 行解析错误: %s\n", parseErr.Line, parseErr.Context)
    } else if errors.IsFormatError(err) {
        // 处理格式错误
        fmt.Println("无效的配置格式")
    }
}
```

## 最佳实践

1. **使用适当的解析器类型**: 简单解析使用 `NewConfigParser()`，编辑场景使用 `NewPositionAwareParser()`
2. **正确处理错误**: 始终使用错误工具检查特定错误类型
3. **验证输入**: 确保文件路径存在且内容在解析前有效
4. **资源管理**: 使用 `ParseFromReader` 时正确关闭文件句柄
5. **位置跟踪**: 仅在编辑需要时启用位置跟踪以避免开销

## 线程安全

ConfigParser 不是线程安全的。为并发使用创建单独的实例或提供适当的同步。
