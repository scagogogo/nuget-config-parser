# 配置

本指南解释了 NuGet 配置文件的结构和组件，以及如何使用 NuGet Config Parser 库处理它们。

## 概述

NuGet 配置文件 (`NuGet.Config`) 是控制 NuGet 行为各个方面的 XML 文件，包括：

- 包源位置
- 认证凭证
- 全局设置和首选项
- 包还原行为
- 代理设置

## 配置文件结构

典型的 NuGet.Config 文件具有以下结构：

```xml
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="C:\LocalPackages" />
  </packageSources>
  
  <packageSourceCredentials>
    <MyPrivateSource>
      <add key="Username" value="myuser" />
      <add key="ClearTextPassword" value="mypass" />
    </MyPrivateSource>
  </packageSourceCredentials>
  
  <disabledPackageSources>
    <add key="local" value="true" />
  </disabledPackageSources>
  
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  
  <config>
    <add key="globalPackagesFolder" value="C:\packages" />
    <add key="repositoryPath" value=".\packages" />
    <add key="defaultPushSource" value="https://api.nuget.org/v3/index.json" />
  </config>
</configuration>
```

## 配置部分

### 包源

`<packageSources>` 部分定义 NuGet 查找包的位置：

```xml
<packageSources>
  <clear />  <!-- 可选：清除所有继承的源 -->
  <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
  <add key="company-feed" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
  <add key="local-packages" value="C:\LocalPackages" />
</packageSources>
```

**属性：**
- `key`: 源的唯一标识符
- `value`: 包源的 URL 或文件路径
- `protocolVersion`: NuGet 协议版本（"2" 或 "3"）

### 包源凭证

`<packageSourceCredentials>` 部分存储认证信息：

```xml
<packageSourceCredentials>
  <MyPrivateSource>
    <add key="Username" value="myuser" />
    <add key="ClearTextPassword" value="mypass" />
  </MyPrivateSource>
  <AnotherSource>
    <add key="Username" value="user2" />
    <add key="Password" value="encrypted_password" />
  </AnotherSource>
</packageSourceCredentials>
```

**凭证类型：**
- `Username`: 认证用户名
- `Password`: 加密密码
- `ClearTextPassword`: 明文密码（生产环境不推荐）

### 禁用的包源

`<disabledPackageSources>` 部分列出临时禁用的源：

```xml
<disabledPackageSources>
  <add key="local-packages" value="true" />
  <add key="old-feed" value="true" />
</disabledPackageSources>
```

### 活跃包源

`<activePackageSource>` 部分指定当前活跃的源：

```xml
<activePackageSource>
  <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
</activePackageSource>
```

### 全局配置

`<config>` 部分包含全局 NuGet 设置：

```xml
<config>
  <add key="globalPackagesFolder" value="C:\packages" />
  <add key="repositoryPath" value=".\packages" />
  <add key="defaultPushSource" value="https://api.nuget.org/v3/index.json" />
  <add key="dependencyVersion" value="Highest" />
  <add key="http_proxy" value="http://proxy.company.com:8080" />
  <add key="http_proxy.user" value="proxyuser" />
  <add key="http_proxy.password" value="proxypass" />
</config>
```

**常见配置键：**
- `globalPackagesFolder`: 全局包缓存位置
- `repositoryPath`: 项目包文件夹
- `defaultPushSource`: 包发布的默认源
- `dependencyVersion`: 默认依赖版本解析
- `http_proxy`: HTTP 代理服务器
- `automaticPackageRestore`: 启用自动包还原

## 配置层次结构

NuGet 使用分层配置系统，其中设置被继承并可以被覆盖：

1. **计算机级别**: 系统范围设置
2. **用户级别**: 用户特定设置
3. **解决方案级别**: 解决方案特定设置
4. **项目级别**: 项目特定设置

### 搜索顺序

库按以下顺序搜索配置文件：

1. 当前目录: `./NuGet.Config`
2. 父目录（向上遍历树）
3. 用户配置目录
4. 系统配置目录

### 平台特定位置

**Windows:**
- 用户: `%APPDATA%\NuGet\NuGet.Config`
- 系统: `%ProgramData%\NuGet\NuGet.Config`

**macOS:**
- 用户: `~/Library/Application Support/NuGet/NuGet.Config`
- 系统: `/Library/Application Support/NuGet/NuGet.Config`

**Linux:**
- 用户: `~/.config/NuGet/NuGet.Config`
- 系统: `/etc/NuGet/NuGet.Config`

## 使用配置

### 读取配置

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 查找并解析配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    fmt.Printf("从以下位置加载配置: %s\n", configPath)
    
    // 访问包源
    for _, source := range config.PackageSources.Add {
        fmt.Printf("源: %s -> %s\n", source.Key, source.Value)
    }
    
    // 访问配置选项
    if config.Config != nil {
        for _, option := range config.Config.Add {
            fmt.Printf("设置: %s = %s\n", option.Key, option.Value)
        }
    }
}
```

### 修改配置

```go
package main

import (
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 加载现有配置
    config, configPath, err := api.FindAndParseConfig()
    if err != nil {
        // 如果未找到则创建默认配置
        config = api.CreateDefaultConfig()
        configPath = "NuGet.Config"
    }
    
    // 添加包源
    api.AddPackageSource(config, "company", "https://nuget.company.com/v3/index.json", "3")
    
    // 添加凭证
    api.AddCredential(config, "company", "myuser", "mypass")
    
    // 配置全局设置
    api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages")
    
    // 设置活跃源
    api.SetActivePackageSource(config, "company", "https://nuget.company.com/v3/index.json")
    
    // 保存配置
    err = api.SaveConfig(config, configPath)
    if err != nil {
        log.Fatalf("保存配置失败: %v", err)
    }
}
```

## 最佳实践

### 安全性

1. **避免明文密码**: 尽可能使用加密密码
2. **安全文件权限**: 确保配置文件具有适当的权限
3. **环境变量**: 对敏感信息使用环境变量
4. **凭证管理**: 考虑使用凭证管理器进行认证

### 组织

1. **分层配置**: 对项目特定设置使用项目级配置
2. **一致命名**: 为包源使用描述性名称
3. **文档**: 尽可能为配置文件添加注释
4. **版本控制**: 在版本控制中包含项目级配置

### 性能

1. **最小化源**: 仅包含必要的包源
2. **协议版本**: 为源使用适当的协议版本
3. **本地缓存**: 配置适当的缓存位置
4. **禁用未使用的源**: 禁用不需要的源

## 常见配置模式

### 企业设置

```xml
<configuration>
  <packageSources>
    <clear />
    <add key="company-internal" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
  </packageSources>
  
  <config>
    <add key="globalPackagesFolder" value="C:\CompanyPackages" />
    <add key="defaultPushSource" value="https://nuget.company.com/v3/index.json" />
    <add key="http_proxy" value="http://proxy.company.com:8080" />
  </config>
</configuration>
```

### 开发设置

```xml
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local-dev" value="./local-packages" />
    <add key="preview" value="https://api.nuget.org/v3-flatcontainer" protocolVersion="3" />
  </packageSources>
  
  <disabledPackageSources>
    <add key="preview" value="true" />
  </disabledPackageSources>
  
  <config>
    <add key="repositoryPath" value="./packages" />
    <add key="dependencyVersion" value="Highest" />
  </config>
</configuration>
```

### CI/CD 设置

```xml
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="build-artifacts" value="https://artifacts.company.com/nuget" protocolVersion="3" />
  </packageSources>
  
  <config>
    <add key="globalPackagesFolder" value="/tmp/packages" />
    <add key="automaticPackageRestore" value="true" />
  </config>
</configuration>
```

## 故障排除

### 常见问题

1. **文件未找到**: 检查文件路径和权限
2. **无效 XML**: 验证 XML 结构和编码
3. **认证失败**: 验证凭证和源 URL
4. **源冲突**: 检查重复的源键
5. **权限错误**: 确保适当的文件和目录权限

### 调试配置

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
    
    // 查找所有配置文件
    configPaths := api.FindAllConfigFiles()
    fmt.Printf("找到 %d 个配置文件:\n", len(configPaths))
    
    for i, path := range configPaths {
        fmt.Printf("%d. %s\n", i+1, path)
        
        config, err := api.ParseFromFile(path)
        if err != nil {
            if errors.IsParseError(err) {
                fmt.Printf("   解析错误: %v\n", err)
            } else {
                fmt.Printf("   错误: %v\n", err)
            }
            continue
        }
        
        fmt.Printf("   源: %d\n", len(config.PackageSources.Add))
        fmt.Printf("   设置: %d\n", len(config.Config.Add))
    }
}
```

## 下一步

- 了解 [位置感知编辑](./position-aware-editing.md) 进行高级配置修改
- 探索 [API 参考](/zh/api/) 获取详细的方法文档
- 查看 [示例](/zh/examples/) 了解实际使用场景

本配置指南提供了对 NuGet 配置文件的全面理解，以及如何使用库有效地处理它们。
