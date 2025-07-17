# 工具 API

`pkg/utils` 包为整个 NuGet Config Parser 库中使用的文件操作、路径操作和 XML 处理提供工具函数。

## 概述

工具 API 提供：
- 文件系统操作和检查
- 跨平台路径操作
- XML 处理工具
- 字符串和数据验证辅助函数

## 文件系统操作

### FileExists

```go
func FileExists(filePath string) bool
```

检查文件是否存在且不是目录。

**参数:**
- `filePath` (string): 要检查的文件路径

**返回值:**
- `bool`: 如果文件存在且不是目录则为 true，否则为 false

**示例:**
```go
import "github.com/scagogogo/nuget-config-parser/pkg/utils"

configPath := "/path/to/NuGet.Config"
if utils.FileExists(configPath) {
    fmt.Printf("配置文件存在: %s\n", configPath)
    // 继续解析...
} else {
    fmt.Printf("配置文件未找到: %s\n", configPath)
    // 创建默认配置...
}
```

### DirExists

```go
func DirExists(dirPath string) bool
```

检查目录是否存在。

**参数:**
- `dirPath` (string): 要检查的目录路径

**返回值:**
- `bool`: 如果目录存在则为 true，否则为 false

**示例:**
```go
configDir := "/etc/nuget"
if utils.DirExists(configDir) {
    fmt.Printf("配置目录存在: %s\n", configDir)
} else {
    fmt.Printf("创建配置目录: %s\n", configDir)
    err := os.MkdirAll(configDir, 0755)
    if err != nil {
        log.Fatalf("创建目录失败: %v", err)
    }
}
```

### IsReadableFile

```go
func IsReadableFile(filePath string) bool
```

检查文件是否存在且可读。

**参数:**
- `filePath` (string): 要检查的文件路径

**返回值:**
- `bool`: 如果文件存在且可读则为 true，否则为 false

**示例:**
```go
configPath := "/path/to/NuGet.Config"
if utils.IsReadableFile(configPath) {
    // 安全读取文件
    content, err := os.ReadFile(configPath)
    if err != nil {
        log.Printf("读取文件错误: %v", err)
    }
} else {
    fmt.Printf("文件不可读: %s\n", configPath)
}
```

### IsWritableFile

```go
func IsWritableFile(filePath string) bool
```

检查文件是否可写（或如果不存在，目录是否可写）。

**参数:**
- `filePath` (string): 要检查的文件路径

**返回值:**
- `bool`: 如果文件可写则为 true，否则为 false

**示例:**
```go
configPath := "/path/to/NuGet.Config"
if utils.IsWritableFile(configPath) {
    // 安全写入文件
    err := os.WriteFile(configPath, []byte("content"), 0644)
    if err != nil {
        log.Printf("写入文件错误: %v", err)
    }
} else {
    fmt.Printf("文件不可写: %s\n", configPath)
}
```

## 路径操作

### IsAbsolutePath

```go
func IsAbsolutePath(path string) bool
```

检查路径是否为绝对路径。

**参数:**
- `path` (string): 要检查的路径

**返回值:**
- `bool`: 如果路径是绝对路径则为 true，相对路径则为 false

**示例:**
```go
// Unix/Linux 示例
fmt.Println(utils.IsAbsolutePath("/etc/nuget"))        // true
fmt.Println(utils.IsAbsolutePath("./config"))          // false
fmt.Println(utils.IsAbsolutePath("../config"))         // false

// Windows 示例
fmt.Println(utils.IsAbsolutePath("C:\\nuget\\config")) // true
fmt.Println(utils.IsAbsolutePath(".\\config"))         // false
```

### NormalizePath

```go
func NormalizePath(path string) string
```

通过清理路径并转换为操作系统特定格式来规范化文件路径。

**参数:**
- `path` (string): 要规范化的路径

**返回值:**
- `string`: 规范化的路径

**示例:**
```go
// 清理冗余路径元素
messyPath := "/etc/nuget/../nuget/./config"
cleanPath := utils.NormalizePath(messyPath)
fmt.Printf("规范化: %s\n", cleanPath)
// 输出: /etc/nuget/config

// 处理不同的分隔符
mixedPath := "/etc\\nuget/config"
normalizedPath := utils.NormalizePath(mixedPath)
fmt.Printf("规范化: %s\n", normalizedPath)
```

### JoinPaths

```go
func JoinPaths(basePath string, paths ...string) string
```

使用操作系统特定的路径分隔符连接多个路径元素。

**参数:**
- `basePath` (string): 基础路径
- `paths` (...string): 要连接的附加路径元素

**返回值:**
- `string`: 连接的路径

**示例:**
```go
// 连接路径元素
basePath := "/etc"
configPath := utils.JoinPaths(basePath, "nuget", "NuGet.Config")
fmt.Printf("连接的路径: %s\n", configPath)
// 输出: /etc/nuget/NuGet.Config

// Windows 示例
winBase := "C:\\Users"
winPath := utils.JoinPaths(winBase, "username", ".nuget", "packages")
fmt.Printf("Windows 路径: %s\n", winPath)
// 输出: C:\Users\username\.nuget\packages

// 处理尾随斜杠
trailingSlash := "/home/user/"
result := utils.JoinPaths(trailingSlash, "nuget", "packages")
fmt.Printf("结果: %s\n", result)
// 输出: /home/user/nuget/packages
```

### ResolvePath

```go
func ResolvePath(basePath, path string) string
```

相对于基础路径解析路径。如果路径已经是绝对路径，则返回规范化的路径。

**参数:**
- `basePath` (string): 解析相对路径的基础路径
- `path` (string): 要解析的路径（可以是相对或绝对）

**返回值:**
- `string`: 解析的绝对路径

**示例:**
```go
basePath := "/etc/nuget"

// 解析相对路径
relativePath := "../packages/cache"
resolvedRelative := utils.ResolvePath(basePath, relativePath)
fmt.Printf("相对路径 '%s' 解析为: %s\n", relativePath, resolvedRelative)
// 输出: 相对路径 '../packages/cache' 解析为: /etc/packages/cache

// 处理绝对路径
absolutePath := "/var/nuget/packages"
resolvedAbsolute := utils.ResolvePath(basePath, absolutePath)
fmt.Printf("绝对路径 '%s' 保持: %s\n", absolutePath, resolvedAbsolute)
// 输出: 绝对路径 '/var/nuget/packages' 保持: /var/nuget/packages
```

### ExpandEnvVars

```go
func ExpandEnvVars(path string) string
```

展开路径字符串中的环境变量。

**参数:**
- `path` (string): 包含环境变量的路径

**返回值:**
- `string`: 展开环境变量后的路径

**支持的格式:**
- Unix/Linux/macOS: `$VAR` 或 `${VAR}`
- Windows: `%VAR%`

**示例:**
```go
// Unix/Linux/macOS 示例
unixPath := "$HOME/.nuget/packages"
expandedUnixPath := utils.ExpandEnvVars(unixPath)
fmt.Printf("展开: %s\n", expandedUnixPath)
// 输出: 展开: /home/user/.nuget/packages

bracedPath := "${HOME}/.config/NuGet/NuGet.Config"
expandedBracedPath := utils.ExpandEnvVars(bracedPath)
fmt.Printf("展开: %s\n", expandedBracedPath)

// Windows 示例
winPath := "%USERPROFILE%\\.nuget\\packages"
expandedWinPath := utils.ExpandEnvVars(winPath)
fmt.Printf("展开: %s\n", expandedWinPath)
// 输出: 展开: C:\Users\username\.nuget\packages
```

## XML 处理

### ValidateXML

```go
func ValidateXML(content []byte) error
```

验证内容是否为格式良好的 XML。

**参数:**
- `content` ([]byte): 要验证的 XML 内容

**返回值:**
- `error`: 如果 XML 格式不正确则返回错误，有效则为 nil

**示例:**
```go
// 有效 XML
validXML := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </packageSources>
</configuration>`)

err := utils.ValidateXML(validXML)
if err != nil {
    fmt.Printf("无效 XML: %v\n", err)
} else {
    fmt.Println("XML 有效")
}

// 无效 XML
invalidXML := []byte(`<configuration><packageSources><add key="test"`)
err = utils.ValidateXML(invalidXML)
if err != nil {
    fmt.Printf("无效 XML: %v\n", err)
}
```

### FormatXML

```go
func FormatXML(content []byte) ([]byte, error)
```

使用适当的缩进格式化 XML 内容。

**参数:**
- `content` ([]byte): 要格式化的 XML 内容

**返回值:**
- `[]byte`: 格式化的 XML 内容
- `error`: 格式化失败时的错误

**示例:**
```go
// 未格式化的 XML
unformattedXML := []byte(`<configuration><packageSources><add key="nuget.org" value="https://api.nuget.org/v3/index.json" /></packageSources></configuration>`)

formattedXML, err := utils.FormatXML(unformattedXML)
if err != nil {
    log.Printf("格式化 XML 失败: %v", err)
} else {
    fmt.Println("格式化的 XML:")
    fmt.Println(string(formattedXML))
}

// 输出:
// <?xml version="1.0" encoding="UTF-8"?>
// <configuration>
//   <packageSources>
//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json"/>
//   </packageSources>
// </configuration>
```

### ExtractXMLElement

```go
func ExtractXMLElement(content []byte, elementName string) ([]byte, error)
```

从内容中提取特定的 XML 元素。

**参数:**
- `content` ([]byte): 要搜索的 XML 内容
- `elementName` (string): 要提取的元素名称

**返回值:**
- `[]byte`: 提取的元素内容
- `error`: 未找到元素或提取失败时的错误

**示例:**
```go
xmlContent := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
    <add key="local" value="/path/to/local" />
  </packageSources>
  <config>
    <add key="globalPackagesFolder" value="/packages" />
  </config>
</configuration>`)

// 提取 packageSources 元素
packageSources, err := utils.ExtractXMLElement(xmlContent, "packageSources")
if err != nil {
    log.Printf("提取元素失败: %v", err)
} else {
    fmt.Println("包源:")
    fmt.Println(string(packageSources))
}
```

## 字符串工具

### IsEmptyOrWhitespace

```go
func IsEmptyOrWhitespace(s string) bool
```

检查字符串是否为空或仅包含空白字符。

**参数:**
- `s` (string): 要检查的字符串

**返回值:**
- `bool`: 如果字符串为空或仅包含空白字符则为 true

**示例:**
```go
fmt.Println(utils.IsEmptyOrWhitespace(""))           // true
fmt.Println(utils.IsEmptyOrWhitespace("   "))        // true
fmt.Println(utils.IsEmptyOrWhitespace("\t\n"))       // true
fmt.Println(utils.IsEmptyOrWhitespace("content"))    // false
fmt.Println(utils.IsEmptyOrWhitespace(" content "))  // false
```

### TrimWhitespace

```go
func TrimWhitespace(s string) string
```

从字符串中修剪前导和尾随空白字符。

**参数:**
- `s` (string): 要修剪的字符串

**返回值:**
- `string`: 修剪后的字符串

**示例:**
```go
input := "  \t  content with spaces  \n  "
trimmed := utils.TrimWhitespace(input)
fmt.Printf("修剪后: '%s'\n", trimmed)
// 输出: 修剪后: 'content with spaces'
```

### SanitizeXMLValue

```go
func SanitizeXMLValue(value string) string
```

清理字符串值以在 XML 属性或内容中安全使用。

**参数:**
- `value` (string): 要清理的值

**返回值:**
- `string`: 对 XML 安全的清理值

**示例:**
```go
unsafeValue := `value with "quotes" & <brackets>`
safeValue := utils.SanitizeXMLValue(unsafeValue)
fmt.Printf("清理后: %s\n", safeValue)
// 输出: 清理后: value with &quot;quotes&quot; &amp; &lt;brackets&gt;
```

## 验证工具

### IsValidURL

```go
func IsValidURL(urlStr string) bool
```

验证字符串是否为有效的 URL。

**参数:**
- `urlStr` (string): 要验证的 URL 字符串

**返回值:**
- `bool`: 如果 URL 有效则为 true

**示例:**
```go
validURLs := []string{
    "https://api.nuget.org/v3/index.json",
    "http://localhost:8080/nuget",
    "file:///path/to/packages",
}

invalidURLs := []string{
    "not-a-url",
    "://missing-scheme",
    "https://",
}

for _, url := range validURLs {
    if utils.IsValidURL(url) {
        fmt.Printf("有效 URL: %s\n", url)
    }
}

for _, url := range invalidURLs {
    if !utils.IsValidURL(url) {
        fmt.Printf("无效 URL: %s\n", url)
    }
}
```

### IsValidPackageSourceKey

```go
func IsValidPackageSourceKey(key string) bool
```

验证字符串是否为有效的包源键。

**参数:**
- `key` (string): 要验证的包源键

**返回值:**
- `bool`: 如果键有效则为 true

**验证规则:**
- 不为空或仅包含空白字符
- 仅包含字母数字字符、连字符、下划线和点
- 不以特殊字符开头或结尾

**示例:**
```go
validKeys := []string{
    "nuget.org",
    "company-feed",
    "local_packages",
    "feed123",
}

invalidKeys := []string{
    "",
    "   ",
    "invalid key with spaces",
    "-starts-with-dash",
    "ends-with-dash-",
    "has@invalid#chars",
}

for _, key := range validKeys {
    if utils.IsValidPackageSourceKey(key) {
        fmt.Printf("有效键: %s\n", key)
    }
}

for _, key := range invalidKeys {
    if !utils.IsValidPackageSourceKey(key) {
        fmt.Printf("无效键: '%s'\n", key)
    }
}
```

## 完整使用示例

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/utils"
)

func main() {
    // 文件操作
    configPath := "/path/to/NuGet.Config"
    
    if utils.FileExists(configPath) {
        fmt.Printf("配置文件存在: %s\n", configPath)
        
        if utils.IsReadableFile(configPath) {
            content, err := os.ReadFile(configPath)
            if err != nil {
                log.Printf("读取文件错误: %v", err)
                return
            }
            
            // 验证 XML
            if err := utils.ValidateXML(content); err != nil {
                log.Printf("无效 XML: %v", err)
                return
            }
            
            // 格式化 XML
            formatted, err := utils.FormatXML(content)
            if err != nil {
                log.Printf("格式化 XML 失败: %v", err)
            } else {
                fmt.Println("格式化的 XML:")
                fmt.Println(string(formatted))
            }
        }
    } else {
        fmt.Printf("配置文件未找到: %s\n", configPath)
        
        // 如果需要创建目录
        configDir := utils.JoinPaths("/etc", "nuget")
        if !utils.DirExists(configDir) {
            fmt.Printf("创建目录: %s\n", configDir)
            err := os.MkdirAll(configDir, 0755)
            if err != nil {
                log.Fatalf("创建目录失败: %v", err)
            }
        }
    }
    
    // 路径操作
    basePath := "/etc/nuget"
    relativePath := "../packages"
    resolvedPath := utils.ResolvePath(basePath, relativePath)
    fmt.Printf("解析的路径: %s\n", resolvedPath)
    
    // 环境变量展开
    envPath := "$HOME/.nuget/packages"
    expandedPath := utils.ExpandEnvVars(envPath)
    fmt.Printf("展开的路径: %s\n", expandedPath)
    
    // 验证
    testURL := "https://api.nuget.org/v3/index.json"
    if utils.IsValidURL(testURL) {
        fmt.Printf("有效 URL: %s\n", testURL)
    }
    
    testKey := "nuget.org"
    if utils.IsValidPackageSourceKey(testKey) {
        fmt.Printf("有效包源键: %s\n", testKey)
    }
}
```

## 最佳实践

1. **检查文件存在性**: 在尝试文件操作之前始终使用 `FileExists()`
2. **验证输入**: 在处理之前使用验证函数检查 URL 和键
3. **正确处理路径**: 使用路径操作函数实现跨平台兼容性
4. **清理 XML 内容**: 对用户提供的内容使用 `SanitizeXMLValue()`
5. **验证 XML**: 在解析之前使用 `ValidateXML()` 尽早捕获格式错误的内容
6. **规范化路径**: 使用 `NormalizePath()` 清理路径字符串
7. **展开环境变量**: 使用 `ExpandEnvVars()` 实现灵活的路径配置

## 线程安全

此包中的所有工具函数都是线程安全的，可以并发使用。
