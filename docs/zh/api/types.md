# 类型 API

`pkg/types` 包定义了用于表示 NuGet 配置文件的所有数据结构。这些类型直接对应于 NuGet.Config 文件的 XML 结构。

## 核心配置类型

### NuGetConfig

```go
type NuGetConfig struct {
    PackageSources             PackageSources             `xml:"packageSources"`
    PackageSourceCredentials   *PackageSourceCredentials  `xml:"packageSourceCredentials,omitempty"`
    Config                     *Config                    `xml:"config,omitempty"`
    DisabledPackageSources     *DisabledPackageSources    `xml:"disabledPackageSources,omitempty"`
    ActivePackageSource        *ActivePackageSource       `xml:"activePackageSource,omitempty"`
}
```

表示完整 NuGet 配置的根配置结构。

**字段:**
- `PackageSources`: 可用包源列表（必需）
- `PackageSourceCredentials`: 包源凭证（可选）
- `Config`: 全局配置选项（可选）
- `DisabledPackageSources`: 禁用包源列表（可选）
- `ActivePackageSource`: 当前活跃包源（可选）

**示例:**
```go
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:             "nuget.org",
                Value:           "https://api.nuget.org/v3/index.json",
                ProtocolVersion: "3",
            },
        },
    },
}
```

## 包源类型

### PackageSources

```go
type PackageSources struct {
    Clear bool            `xml:"clear,attr,omitempty"`
    Add   []PackageSource `xml:"add"`
}
```

包源定义的容器。

**字段:**
- `Clear`: 如果为 true，清除所有先前定义的包源
- `Add`: 要添加的包源列表

### PackageSource

```go
type PackageSource struct {
    Key             string `xml:"key,attr"`
    Value           string `xml:"value,attr"`
    ProtocolVersion string `xml:"protocolVersion,attr,omitempty"`
}
```

表示单个包源。

**字段:**
- `Key`: 包源的唯一标识符
- `Value`: 包源的 URL 或文件路径
- `ProtocolVersion`: NuGet 协议版本（可选，通常为 "2" 或 "3"）

**示例:**
```go
source := types.PackageSource{
    Key:             "company-feed",
    Value:           "https://nuget.company.com/v3/index.json",
    ProtocolVersion: "3",
}
```

### DisabledPackageSources

```go
type DisabledPackageSources struct {
    Add []DisabledSource `xml:"add"`
}
```

禁用包源的容器。

### DisabledSource

```go
type DisabledSource struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

表示禁用的包源。

**字段:**
- `Key`: 要禁用的包源的键
- `Value`: 通常为 "true" 表示源被禁用

### ActivePackageSource

```go
type ActivePackageSource struct {
    Add PackageSource `xml:"add"`
}
```

表示当前活跃的包源。

## 凭证类型

### PackageSourceCredentials

```go
type PackageSourceCredentials struct {
    Sources map[string]SourceCredential `xml:"-"`
}
```

包源凭证的容器。`Sources` 映射使用包源键作为键，凭证作为值。

**注意:** 此类型具有自定义 XML 编组/解组逻辑，以处理 NuGet.Config 文件中凭证的动态结构。

### SourceCredential

```go
type SourceCredential struct {
    Add []Credential `xml:"add"`
}
```

特定包源的凭证。

### Credential

```go
type Credential struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

单个凭证键值对。

**常见凭证键:**
- `Username`: 认证用户名
- `Password`: 认证密码（通常加密）
- `ClearTextPassword`: 明文密码（生产环境不推荐）

**示例:**
```go
credentials := types.SourceCredential{
    Add: []types.Credential{
        {Key: "Username", Value: "myuser"},
        {Key: "ClearTextPassword", Value: "mypassword"},
    },
}
```

## 配置选项类型

### Config

```go
type Config struct {
    Add []ConfigOption `xml:"add"`
}
```

全局配置选项的容器。

### ConfigOption

```go
type ConfigOption struct {
    Key   string `xml:"key,attr"`
    Value string `xml:"value,attr"`
}
```

单个配置选项。

**常见配置键:**
- `globalPackagesFolder`: 全局包文件夹路径
- `repositoryPath`: 包仓库路径
- `defaultPushSource`: 包发布的默认源
- `http_proxy`: HTTP 代理设置
- `http_proxy.user`: HTTP 代理用户名
- `http_proxy.password`: HTTP 代理密码

**示例:**
```go
configOptions := []types.ConfigOption{
    {Key: "globalPackagesFolder", Value: "/custom/packages/path"},
    {Key: "defaultPushSource", Value: "https://my-nuget-server.com"},
}
```

## XML 编组

所有类型都通过 Go 的 `encoding/xml` 包支持 XML 编组和解组。结构标签定义每个字段如何映射到 XML 元素和属性。

### 自定义编组

`PackageSourceCredentials` 类型实现自定义 XML 编组以处理动态结构，其中每个包源都有自己的 XML 元素：

```xml
<packageSourceCredentials>
  <MyPrivateSource>
    <add key="Username" value="myuser" />
    <add key="ClearTextPassword" value="mypass" />
  </MyPrivateSource>
</packageSourceCredentials>
```

## 使用示例

### 创建完整配置

```go
config := &types.NuGetConfig{
    PackageSources: types.PackageSources{
        Add: []types.PackageSource{
            {
                Key:             "nuget.org",
                Value:           "https://api.nuget.org/v3/index.json",
                ProtocolVersion: "3",
            },
            {
                Key:   "local",
                Value: "/path/to/local/packages",
            },
        },
    },
    PackageSourceCredentials: &types.PackageSourceCredentials{
        Sources: map[string]types.SourceCredential{
            "private-feed": {
                Add: []types.Credential{
                    {Key: "Username", Value: "user"},
                    {Key: "ClearTextPassword", Value: "pass"},
                },
            },
        },
    },
    Config: &types.Config{
        Add: []types.ConfigOption{
            {Key: "globalPackagesFolder", Value: "/custom/packages"},
        },
    },
    DisabledPackageSources: &types.DisabledPackageSources{
        Add: []types.DisabledSource{
            {Key: "local", Value: "true"},
        },
    },
    ActivePackageSource: &types.ActivePackageSource{
        Add: types.PackageSource{
            Key:   "nuget.org",
            Value: "https://api.nuget.org/v3/index.json",
        },
    },
}
```

### 使用包源

```go
// 添加新包源
newSource := types.PackageSource{
    Key:             "company-feed",
    Value:           "https://nuget.company.com/v3/index.json",
    ProtocolVersion: "3",
}
config.PackageSources.Add = append(config.PackageSources.Add, newSource)

// 查找包源
var foundSource *types.PackageSource
for i, source := range config.PackageSources.Add {
    if source.Key == "company-feed" {
        foundSource = &config.PackageSources.Add[i]
        break
    }
}

// 移除包源
for i, source := range config.PackageSources.Add {
    if source.Key == "company-feed" {
        config.PackageSources.Add = append(
            config.PackageSources.Add[:i],
            config.PackageSources.Add[i+1:]...,
        )
        break
    }
}
```

### 使用凭证

```go
// 如果为 nil 则初始化凭证
if config.PackageSourceCredentials == nil {
    config.PackageSourceCredentials = &types.PackageSourceCredentials{
        Sources: make(map[string]types.SourceCredential),
    }
}

// 为源添加凭证
config.PackageSourceCredentials.Sources["private-feed"] = types.SourceCredential{
    Add: []types.Credential{
        {Key: "Username", Value: "myuser"},
        {Key: "ClearTextPassword", Value: "mypass"},
    },
}

// 获取源的凭证
if cred, exists := config.PackageSourceCredentials.Sources["private-feed"]; exists {
    for _, c := range cred.Add {
        if c.Key == "Username" {
            fmt.Printf("用户名: %s\n", c.Value)
        }
    }
}
```

### 使用配置选项

```go
// 如果为 nil 则初始化配置
if config.Config == nil {
    config.Config = &types.Config{
        Add: []types.ConfigOption{},
    }
}

// 添加配置选项
config.Config.Add = append(config.Config.Add, types.ConfigOption{
    Key:   "globalPackagesFolder",
    Value: "/custom/packages/path",
})

// 查找配置选项
var globalPackagesFolder string
for _, option := range config.Config.Add {
    if option.Key == "globalPackagesFolder" {
        globalPackagesFolder = option.Value
        break
    }
}
```

## 验证

虽然类型本身不包含验证逻辑，但在创建或修改配置时应验证数据：

```go
func validatePackageSource(source types.PackageSource) error {
    if source.Key == "" {
        return errors.New("包源键不能为空")
    }
    if source.Value == "" {
        return errors.New("包源值不能为空")
    }
    if source.ProtocolVersion != "" && 
       source.ProtocolVersion != "2" && 
       source.ProtocolVersion != "3" {
        return errors.New("无效的协议版本")
    }
    return nil
}
```

## 线程安全

此包中的类型不是线程安全的。如果需要从多个 goroutine 访问或修改配置对象，必须提供自己的同步。

## 最佳实践

1. **初始化可选字段**: 在访问可选字段之前始终检查它们是否为 nil
2. **对可选结构使用指针**: 可选配置部分使用指针来区分空和缺失
3. **验证数据**: 在创建配置之前验证包源键、URL 和其他数据
4. **小心处理凭证**: 谨慎处理凭证，特别是明文密码
5. **使用 API**: 在可能的情况下，优先使用高级 API 方法而不是直接结构操作
