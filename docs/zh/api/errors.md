# 错误 API

`pkg/errors` 包定义了用于处理 NuGet 配置解析和管理错误的错误类型和工具。

## 概述

错误 API 提供：
- 针对不同失败场景的结构化错误类型
- 错误分类工具
- 带有上下文的详细错误信息
- 支持错误包装和解包

## 错误常量

### 预定义错误

```go
var (
    // ErrInvalidConfigFormat 表示无效的配置文件格式
    ErrInvalidConfigFormat = errors.New("invalid nuget config format")

    // ErrConfigFileNotFound 表示未找到配置文件
    ErrConfigFileNotFound = errors.New("nuget config file not found")

    // ErrEmptyConfigFile 表示空配置文件
    ErrEmptyConfigFile = errors.New("empty nuget config file")

    // ErrXMLParsing 表示 XML 解析错误
    ErrXMLParsing = errors.New("xml parsing error")

    // ErrMissingRequiredElement 表示配置中缺少必需元素
    ErrMissingRequiredElement = errors.New("missing required element in config")
)
```

这些预定义错误表示常见的失败场景，可以与 `errors.Is()` 一起使用进行错误检查。

## 错误类型

### ParseError

```go
type ParseError struct {
    BaseErr  error  // 基础错误
    Line     int    // 发生错误的行号
    Position int    // 发生错误的行内位置
    Context  string // 附加上下文信息
}
```

表示带有详细位置信息的解析错误。

**字段:**
- `BaseErr`: 导致解析失败的底层错误
- `Line`: 发生错误的行号（从1开始）
- `Position`: 发生错误的行内字符位置（从1开始）
- `Context`: 关于错误的附加上下文信息

**方法:**

#### Error

```go
func (e *ParseError) Error() string
```

返回带有位置信息的格式化错误消息。

**示例:**
```go
parseErr := &errors.ParseError{
    BaseErr:  errors.ErrInvalidConfigFormat,
    Line:     15,
    Position: 23,
    Context:  "无效的属性值",
}

fmt.Println(parseErr.Error())
// 输出: 第15行位置23处解析错误: 无效的属性值 - invalid nuget config format
```

#### Unwrap

```go
func (e *ParseError) Unwrap() error
```

返回基础错误，支持 `errors.Is()` 和 `errors.As()` 函数。

**示例:**
```go
parseErr := &errors.ParseError{
    BaseErr: errors.ErrInvalidConfigFormat,
    Line:    10,
    Position: 5,
    Context: "格式错误的 XML",
}

// 检查是否为格式错误
if errors.Is(parseErr, errors.ErrInvalidConfigFormat) {
    fmt.Println("这是格式错误")
}
```

## 构造函数

### NewParseError

```go
func NewParseError(baseErr error, line, position int, context string) *ParseError
```

创建带有位置信息的新解析错误。

**参数:**
- `baseErr` (error): 底层错误
- `line` (int): 发生错误的行号
- `position` (int): 发生错误的行内位置
- `context` (string): 附加上下文信息

**返回值:**
- `*ParseError`: 新的解析错误实例

**示例:**
```go
// 为无效 XML 创建解析错误
parseErr := errors.NewParseError(
    errors.ErrXMLParsing,
    25,
    10,
    "意外的结束标签",
)

fmt.Printf("解析错误: %v\n", parseErr)
```

## 错误分类函数

### IsNotFoundError

```go
func IsNotFoundError(err error) bool
```

检查错误是否表示未找到配置文件。

**参数:**
- `err` (error): 要检查的错误

**返回值:**
- `bool`: 如果错误表示文件未找到则为 true

**示例:**
```go
config, err := api.ParseFromFile("missing.config")
if err != nil {
    if errors.IsNotFoundError(err) {
        fmt.Println("未找到配置文件")
        // 创建默认配置
        config = api.CreateDefaultConfig()
    } else {
        log.Fatalf("其他错误: %v", err)
    }
}
```

### IsParseError

```go
func IsParseError(err error) bool
```

检查错误是否为解析错误。

**参数:**
- `err` (error): 要检查的错误

**返回值:**
- `bool`: 如果错误是 ParseError 则为 true

**示例:**
```go
config, err := api.ParseFromString(invalidXML)
if err != nil {
    if errors.IsParseError(err) {
        var parseErr *errors.ParseError
        if errors.As(err, &parseErr) {
            fmt.Printf("第 %d 行解析错误: %s\n", 
                parseErr.Line, parseErr.Context)
        }
    }
}
```

### IsFormatError

```go
func IsFormatError(err error) bool
```

检查错误是否表示无效的配置格式。

**参数:**
- `err` (error): 要检查的错误

**返回值:**
- `bool`: 如果错误表示格式问题则为 true

**示例:**
```go
config, err := api.ParseFromFile("invalid.config")
if err != nil {
    if errors.IsFormatError(err) {
        fmt.Println("无效的配置格式")
        // 处理格式错误
    }
}
```

## 使用示例

### 基本错误处理

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    api := nuget.NewAPI()
    
    config, err := api.ParseFromFile("NuGet.Config")
    if err != nil {
        handleError(err)
        return
    }
    
    fmt.Printf("成功加载了 %d 个包源\n", 
        len(config.PackageSources.Add))
}

func handleError(err error) {
    switch {
    case errors.IsNotFoundError(err):
        fmt.Println("未找到配置文件")
        fmt.Println("考虑创建默认配置")
        
    case errors.IsParseError(err):
        var parseErr *errors.ParseError
        if errors.As(err, &parseErr) {
            fmt.Printf("第 %d 行位置 %d 处解析错误: %s\n",
                parseErr.Line, parseErr.Position, parseErr.Context)
            fmt.Printf("底层错误: %v\n", parseErr.BaseErr)
        } else {
            fmt.Printf("解析错误: %v\n", err)
        }
        
    case errors.IsFormatError(err):
        fmt.Println("无效的配置文件格式")
        fmt.Println("请检查 XML 结构")
        
    default:
        log.Printf("意外错误: %v\n", err)
    }
}
```

### 高级错误处理

```go
func parseConfigWithRecovery(filePath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    config, err := api.ParseFromFile(filePath)
    if err != nil {
        // 尝试提供有用的错误信息
        if errors.IsNotFoundError(err) {
            return nil, fmt.Errorf("配置文件 '%s' 未找到: %w", filePath, err)
        }
        
        if errors.IsParseError(err) {
            var parseErr *errors.ParseError
            if errors.As(err, &parseErr) {
                // 提供详细的解析错误信息
                return nil, fmt.Errorf(
                    "在 %s:%d:%d 解析配置失败 - %s: %w",
                    filePath, parseErr.Line, parseErr.Position, 
                    parseErr.Context, parseErr.BaseErr)
            }
        }
        
        // 对于其他错误，用上下文包装
        return nil, fmt.Errorf("从 '%s' 加载配置失败: %w", filePath, err)
    }
    
    return config, nil
}
```

### 错误恢复策略

```go
func loadConfigWithFallback(primaryPath, fallbackPath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    // 尝试主配置
    config, err := api.ParseFromFile(primaryPath)
    if err == nil {
        return config, nil
    }
    
    // 处理主配置错误
    if errors.IsNotFoundError(err) {
        fmt.Printf("在 %s 未找到主配置，尝试备用配置...\n", primaryPath)
    } else if errors.IsParseError(err) {
        fmt.Printf("主配置有解析错误，尝试备用配置...\n")
    } else {
        return nil, fmt.Errorf("加载主配置失败: %w", err)
    }
    
    // 尝试备用配置
    config, err = api.ParseFromFile(fallbackPath)
    if err == nil {
        fmt.Printf("使用来自 %s 的备用配置\n", fallbackPath)
        return config, nil
    }
    
    // 如果备用配置也失败，创建默认配置
    if errors.IsNotFoundError(err) {
        fmt.Println("未找到配置文件，创建默认配置...")
        return api.CreateDefaultConfig(), nil
    }
    
    return nil, fmt.Errorf("加载任何配置失败: %w", err)
}
```

### 自定义错误类型

```go
// 验证失败的自定义错误
type ValidationError struct {
    Field   string
    Value   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("字段 '%s' 值 '%s' 验证错误: %s", 
        e.Field, e.Value, e.Message)
}

// 验证包源
func validatePackageSource(source *types.PackageSource) error {
    if source.Key == "" {
        return &ValidationError{
            Field:   "Key",
            Value:   source.Key,
            Message: "包源键不能为空",
        }
    }
    
    if source.Value == "" {
        return &ValidationError{
            Field:   "Value",
            Value:   source.Value,
            Message: "包源值不能为空",
        }
    }
    
    // 如果指定了协议版本则验证
    if source.ProtocolVersion != "" && 
       source.ProtocolVersion != "2" && 
       source.ProtocolVersion != "3" {
        return &ValidationError{
            Field:   "ProtocolVersion",
            Value:   source.ProtocolVersion,
            Message: "协议版本必须是 '2' 或 '3'",
        }
    }
    
    return nil
}
```

## 错误包装最佳实践

### 添加上下文

```go
func parseConfigFromPath(configPath string) (*types.NuGetConfig, error) {
    api := nuget.NewAPI()
    
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        // 用附加上下文包装错误
        return nil, fmt.Errorf("从 '%s' 解析 NuGet 配置失败: %w", 
            configPath, err)
    }
    
    return config, nil
}
```

### 保留错误类型

```go
func loadAndValidateConfig(configPath string) (*types.NuGetConfig, error) {
    config, err := parseConfigFromPath(configPath)
    if err != nil {
        // 即使在包装后也检查原始错误类型
        if errors.IsNotFoundError(err) {
            // 处理未找到的情况
            return createDefaultConfig(configPath)
        }
        
        if errors.IsParseError(err) {
            // 处理解析错误情况
            return nil, fmt.Errorf("配置文件有语法错误: %w", err)
        }
        
        return nil, err
    }
    
    // 附加验证...
    return config, nil
}
```

## 测试错误条件

```go
func TestErrorHandling(t *testing.T) {
    api := nuget.NewAPI()
    
    // 测试文件未找到
    _, err := api.ParseFromFile("nonexistent.config")
    if !errors.IsNotFoundError(err) {
        t.Errorf("期望未找到错误，得到: %v", err)
    }
    
    // 测试无效 XML
    invalidXML := "<configuration><packageSources><add key="
    _, err = api.ParseFromString(invalidXML)
    if !errors.IsParseError(err) {
        t.Errorf("期望解析错误，得到: %v", err)
    }
    
    // 测试解析错误详细信息
    var parseErr *errors.ParseError
    if errors.As(err, &parseErr) {
        if parseErr.Line <= 0 {
            t.Errorf("期望正行号，得到: %d", parseErr.Line)
        }
    }
}
```

## 最佳实践

1. **使用错误分类函数**: 始终使用 `IsNotFoundError()`、`IsParseError()` 等进行错误检查
2. **提供上下文**: 使用 `fmt.Errorf()` 和 `%w` 动词用附加上下文包装错误
3. **处理特定错误**: 为不同错误类型提供不同的处理
4. **保留错误链**: 使用错误包装来维护原始错误信息
5. **测试错误条件**: 为不同错误场景编写测试
6. **适当记录**: 为不同错误类型使用不同的日志级别
7. **提供恢复**: 为可恢复的错误实现回退策略

## 线程安全

此包中的错误类型和函数是线程安全的，可以并发使用。
