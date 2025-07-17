# 编辑器 API

`pkg/editor` 包为 NuGet 配置文件提供位置感知的编辑功能。这允许您进行精确的修改，同时保留原始文件格式并最小化差异。

## 概述

位置感知编辑在以下情况下特别有用：
- 您想要保持原始文件格式
- 您需要最小化版本控制差异
- 您正在处理具有特定格式要求的配置文件
- 您想要保留注释和空白

## 类型

### ConfigEditor

```go
type ConfigEditor struct {
    parseResult *parser.ParseResult
    edits       []Edit
}
```

跟踪要应用于配置文件的修改的主要编辑器类型。

### Edit

```go
type Edit struct {
    Range   parser.Range // 要替换的范围
    NewText string       // 新文本内容
    Type    string       // 编辑类型: "add", "update", "delete"
}
```

表示带有位置信息的单个编辑操作。

## 构造函数

### NewConfigEditor

```go
func NewConfigEditor(parseResult *parser.ParseResult) *ConfigEditor
```

从包含位置信息的解析结果创建新的配置编辑器。

**参数:**
- `parseResult` (*parser.ParseResult): 带有位置跟踪的解析结果

**返回值:**
- `*ConfigEditor`: 新的编辑器实例

**示例:**
```go
// 使用位置跟踪解析
parseResult, err := api.ParseFromFileWithPositions("/path/to/NuGet.Config")
if err != nil {
    log.Fatal(err)
}

// 创建编辑器
editor := editor.NewConfigEditor(parseResult)
```

## 配置访问

### GetConfig

```go
func (e *ConfigEditor) GetConfig() *types.NuGetConfig
```

返回正在编辑的配置对象。

**返回值:**
- `*types.NuGetConfig`: 配置对象

**示例:**
```go
config := editor.GetConfig()
fmt.Printf("当前源: %d\n", len(config.PackageSources.Add))
```

### GetPositions

```go
func (e *ConfigEditor) GetPositions() map[string]*parser.ElementPosition
```

返回配置中所有元素的位置信息。

**返回值:**
- `map[string]*parser.ElementPosition`: 元素路径到位置信息的映射

**示例:**
```go
positions := editor.GetPositions()
for path, pos := range positions {
    fmt.Printf("元素 %s 在第 %d 行\n", path, pos.Range.Start.Line)
}
```

## 包源操作

### AddPackageSource

```go
func (e *ConfigEditor) AddPackageSource(key, value, protocolVersion string) error
```

向配置添加新的包源。

**参数:**
- `key` (string): 包源的唯一标识符
- `value` (string): 包源的 URL 或路径
- `protocolVersion` (string): 协议版本（可以为空）

**返回值:**
- `error`: 操作失败时的错误

**示例:**
```go
err := editor.AddPackageSource(
    "company-feed", 
    "https://nuget.company.com/v3/index.json", 
    "3"
)
if err != nil {
    log.Fatalf("添加包源失败: %v", err)
}
```

### RemovePackageSource

```go
func (e *ConfigEditor) RemovePackageSource(sourceKey string) error
```

从配置中移除包源。

**参数:**
- `sourceKey` (string): 要移除的包源的键

**返回值:**
- `error`: 如果未找到源或操作失败则返回错误

**示例:**
```go
err := editor.RemovePackageSource("old-feed")
if err != nil {
    log.Printf("移除包源失败: %v", err)
}
```

### UpdatePackageSourceURL

```go
func (e *ConfigEditor) UpdatePackageSourceURL(sourceKey, newURL string) error
```

更新现有包源的 URL。

**参数:**
- `sourceKey` (string): 要更新的包源的键
- `newURL` (string): 包源的新 URL

**返回值:**
- `error`: 如果未找到源或操作失败则返回错误

**示例:**
```go
err := editor.UpdatePackageSourceURL(
    "nuget.org", 
    "https://api.nuget.org/v3/index.json"
)
if err != nil {
    log.Printf("更新 URL 失败: %v", err)
}
```

### UpdatePackageSourceVersion

```go
func (e *ConfigEditor) UpdatePackageSourceVersion(sourceKey, newVersion string) error
```

更新现有包源的协议版本。

**参数:**
- `sourceKey` (string): 要更新的包源的键
- `newVersion` (string): 新的协议版本

**返回值:**
- `error`: 如果未找到源或操作失败则返回错误

**示例:**
```go
err := editor.UpdatePackageSourceVersion("my-feed", "3")
if err != nil {
    log.Printf("更新版本失败: %v", err)
}
```

## 应用更改

### ApplyEdits

```go
func (e *ConfigEditor) ApplyEdits() ([]byte, error)
```

应用所有待处理的编辑并返回修改的文件内容。

**返回值:**
- `[]byte`: 修改的文件内容
- `error`: 应用编辑失败时的错误

**示例:**
```go
// 进行多个更改
editor.AddPackageSource("feed1", "https://feed1.com", "3")
editor.UpdatePackageSourceURL("feed2", "https://newfeed2.com")
editor.RemovePackageSource("old-feed")

// 应用所有更改
modifiedContent, err := editor.ApplyEdits()
if err != nil {
    log.Fatalf("应用编辑失败: %v", err)
}

// 保存到文件
err = os.WriteFile("/path/to/NuGet.Config", modifiedContent, 0644)
if err != nil {
    log.Fatalf("保存文件失败: %v", err)
}
```

## 完整示例

这是一个展示如何使用编辑器的完整示例：

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
    
    // 使用位置跟踪解析
    configPath := "/path/to/NuGet.Config"
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    // 创建编辑器
    editor := api.CreateConfigEditor(parseResult)
    
    // 显示当前配置
    config := editor.GetConfig()
    fmt.Printf("当前包源: %d\n", len(config.PackageSources.Add))
    
    // 进行更改
    fmt.Println("添加新包源...")
    err = editor.AddPackageSource(
        "company-internal", 
        "https://nuget.company.com/v3/index.json", 
        "3"
    )
    if err != nil {
        log.Fatalf("添加源失败: %v", err)
    }
    
    fmt.Println("更新现有源...")
    err = editor.UpdatePackageSourceURL(
        "nuget.org", 
        "https://api.nuget.org/v3/index.json"
    )
    if err != nil {
        log.Printf("警告: %v", err)
    }
    
    // 应用更改
    fmt.Println("应用更改...")
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        log.Fatalf("应用编辑失败: %v", err)
    }
    
    // 保存到文件
    err = os.WriteFile(configPath, modifiedContent, 0644)
    if err != nil {
        log.Fatalf("保存文件失败: %v", err)
    }
    
    fmt.Println("配置更新成功！")
    
    // 验证更改
    updatedConfig := editor.GetConfig()
    fmt.Printf("更新的包源: %d\n", len(updatedConfig.PackageSources.Add))
}
```

## 高级用法

### 批量操作

您可以在应用更改之前执行多个操作：

```go
// 一批中的多个更改
editor.AddPackageSource("feed1", "https://feed1.com", "3")
editor.AddPackageSource("feed2", "https://feed2.com", "3")
editor.UpdatePackageSourceURL("existing", "https://new-url.com")
editor.RemovePackageSource("old-feed")

// 一次性应用所有更改
modifiedContent, err := editor.ApplyEdits()
```

### 错误处理

为每个操作适当处理错误：

```go
err := editor.AddPackageSource("duplicate", "https://example.com", "3")
if err != nil {
    if strings.Contains(err.Error(), "already exists") {
        // 处理重复源
        log.Printf("源已存在，改为更新")
        err = editor.UpdatePackageSourceURL("duplicate", "https://example.com")
    } else {
        log.Fatalf("意外错误: %v", err)
    }
}
```

### 位置信息

访问详细的位置信息：

```go
positions := editor.GetPositions()
for path, elemPos := range positions {
    fmt.Printf("元素: %s\n", path)
    fmt.Printf("  标签: %s\n", elemPos.TagName)
    fmt.Printf("  行: %d-%d\n", elemPos.Range.Start.Line, elemPos.Range.End.Line)
    fmt.Printf("  属性: %v\n", elemPos.Attributes)
}
```

## 位置感知编辑的好处

1. **最小差异**: 只更改文件的必要部分
2. **格式保留**: 保持原始缩进和格式
3. **注释保留**: 保留原始文件中的注释
4. **精确控制**: 对修改内容的精确控制
5. **版本控制友好**: 在版本控制中产生更小、更清洁的差异

## 限制

1. **添加新属性**: 目前对向现有元素添加新属性的支持有限
2. **复杂重构**: 不适合对 XML 进行重大结构更改
3. **内存使用**: 在编辑期间将整个文件内容保存在内存中

## 最佳实践

1. **解析一次**: 对多个编辑操作使用相同的解析结果
2. **批量更改**: 在应用之前将相关更改组合在一起
3. **错误处理**: 始终检查每个操作后的错误
4. **备份**: 考虑在应用更改之前备份原始文件
5. **验证**: 在应用更改后验证配置

```go
// 良好实践：批量操作
editor.AddPackageSource("feed1", "url1", "3")
editor.AddPackageSource("feed2", "url2", "3")
editor.RemovePackageSource("old")
modifiedContent, err := editor.ApplyEdits()

// 良好实践：更改后验证
if err == nil {
    // 重新解析以验证
    _, err = api.ParseFromString(string(modifiedContent))
    if err != nil {
        log.Printf("警告: 生成了无效的 XML: %v", err)
    }
}
```
