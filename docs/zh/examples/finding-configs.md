# 查找配置

本示例演示如何在不同平台和目录结构中定位 NuGet 配置文件。

## 概述

配置发现包括：
- 在标准位置搜索 NuGet.Config 文件
- 理解平台特定路径
- 处理项目级与全局配置
- 实现回退策略

## 示例 1: 基本配置发现

查找配置文件的最简单方法：

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
    
    // 查找第一个可用的配置文件
    configPath, err := api.FindConfigFile()
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Println("在标准位置未找到 NuGet 配置文件")
            fmt.Println("考虑创建默认配置")
            return
        }
        log.Fatalf("搜索配置时出错: %v", err)
    }
    
    fmt.Printf("找到配置文件: %s\n", configPath)
    
    // 解析找到的配置
    config, err := api.ParseFromFile(configPath)
    if err != nil {
        log.Fatalf("解析找到的配置失败: %v", err)
    }
    
    fmt.Printf("配置包含 %d 个包源\n", len(config.PackageSources.Add))
}
```

## 示例 2: 查找所有配置文件

发现搜索层次结构中所有可用的配置文件：

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 查找所有配置文件
    configPaths := api.FindAllConfigFiles()
    
    fmt.Printf("找到 %d 个配置文件:\n", len(configPaths))
    
    if len(configPaths) == 0 {
        fmt.Println("在标准位置未找到配置文件")
        displaySearchPaths()
        return
    }
    
    // 显示所有找到的配置及详细信息
    for i, configPath := range configPaths {
        fmt.Printf("\n%d. %s\n", i+1, configPath)
        
        // 检查文件是否可读
        if info, err := os.Stat(configPath); err == nil {
            fmt.Printf("   大小: %d 字节\n", info.Size())
            fmt.Printf("   修改时间: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
        }
        
        // 尝试解析并显示基本信息
        config, err := api.ParseFromFile(configPath)
        if err != nil {
            fmt.Printf("   ❌ 解析错误: %v\n", err)
        } else {
            fmt.Printf("   ✅ 有效配置\n")
            fmt.Printf("   📦 包源: %d\n", len(config.PackageSources.Add))
            
            // 显示前几个源
            for j, source := range config.PackageSources.Add {
                if j >= 3 {
                    fmt.Printf("   ... 还有 %d 个\n", len(config.PackageSources.Add)-3)
                    break
                }
                fmt.Printf("   - %s\n", source.Key)
            }
        }
    }
}

func displaySearchPaths() {
    fmt.Println("\n标准搜索位置:")
    fmt.Println("1. 当前目录: ./NuGet.Config")
    fmt.Println("2. 父目录（向上遍历）")
    
    if home := os.Getenv("HOME"); home != "" {
        fmt.Printf("3. 用户配置: %s/.config/NuGet/NuGet.Config\n", home)
    }
    
    fmt.Println("4. 系统配置: /etc/NuGet/NuGet.Config")
    fmt.Println("\n注意: 实际路径因操作系统而异")
}
```

## 示例 3: 项目特定配置发现

从特定项目目录开始查找配置文件：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
    "github.com/scagogogo/nuget-config-parser/pkg/errors"
)

func main() {
    api := nuget.NewAPI()
    
    // 获取当前工作目录
    currentDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("获取当前目录失败: %v", err)
    }
    
    fmt.Printf("从以下位置开始搜索项目配置: %s\n", currentDir)
    
    // 查找项目特定配置
    projectConfig, err := api.FindProjectConfig(currentDir)
    if err != nil {
        if errors.IsNotFoundError(err) {
            fmt.Println("未找到项目特定配置")
            
            // 回退到全局配置
            fmt.Println("搜索全局配置...")
            globalConfig, err := api.FindConfigFile()
            if err != nil {
                fmt.Println("也未找到全局配置")
                return
            }
            
            fmt.Printf("使用全局配置: %s\n", globalConfig)
            projectConfig = globalConfig
        } else {
            log.Fatalf("搜索项目配置时出错: %v", err)
        }
    } else {
        fmt.Printf("找到项目配置: %s\n", projectConfig)
    }
    
    // 显示搜索的目录层次结构
    showSearchHierarchy(currentDir)
    
    // 解析并显示配置
    config, err := api.ParseFromFile(projectConfig)
    if err != nil {
        log.Fatalf("解析配置失败: %v", err)
    }
    
    displayConfigSummary(config, projectConfig)
}

func showSearchHierarchy(startDir string) {
    fmt.Println("\n搜索层次结构（从最具体到最一般）:")
    
    dir := startDir
    level := 1
    
    for {
        configPath := filepath.Join(dir, "NuGet.Config")
        exists := fileExists(configPath)
        
        status := "❌"
        if exists {
            status = "✅"
        }
        
        fmt.Printf("%d. %s %s\n", level, configPath, status)
        
        parent := filepath.Dir(dir)
        if parent == dir {
            // 到达根目录
            break
        }
        
        dir = parent
        level++
        
        // 限制显示深度
        if level > 10 {
            fmt.Println("   ... （搜索继续到根目录）")
            break
        }
    }
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func displayConfigSummary(config *types.NuGetConfig, configPath string) {
    fmt.Printf("\n=== 配置摘要 ===\n")
    fmt.Printf("文件: %s\n", configPath)
    fmt.Printf("包源: %d\n", len(config.PackageSources.Add))
    
    if len(config.PackageSources.Add) > 0 {
        fmt.Println("\n包源:")
        for _, source := range config.PackageSources.Add {
            fmt.Printf("  - %s: %s\n", source.Key, source.Value)
        }
    }
    
    if config.ActivePackageSource != nil {
        fmt.Printf("\n活跃源: %s\n", config.ActivePackageSource.Add.Key)
    }
}
```

## 示例 4: 跨平台配置发现

处理平台特定的配置位置：

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    fmt.Printf("操作系统: %s\n", runtime.GOOS)
    fmt.Printf("架构: %s\n", runtime.GOARCH)
    
    // 显示平台特定路径
    showPlatformPaths()
    
    // 查找所有配置
    configPaths := api.FindAllConfigFiles()
    
    fmt.Printf("\n找到 %d 个配置文件:\n", len(configPaths))
    
    for i, configPath := range configPaths {
        fmt.Printf("%d. %s\n", i+1, configPath)
        
        // 分类配置
        category := categorizeConfigPath(configPath)
        fmt.Printf("   类别: %s\n", category)
        
        // 检查可访问性
        if isReadable(configPath) {
            fmt.Printf("   状态: ✅ 可读\n")
        } else {
            fmt.Printf("   状态: ❌ 不可读\n")
        }
    }
    
    // 演示查找和解析
    if len(configPaths) > 0 {
        fmt.Printf("\n使用第一个可用配置: %s\n", configPaths[0])
        
        config, err := api.ParseFromFile(configPaths[0])
        if err != nil {
            fmt.Printf("解析失败: %v\n", err)
        } else {
            fmt.Printf("成功解析 %d 个包源\n", len(config.PackageSources.Add))
        }
    }
}

func showPlatformPaths() {
    fmt.Println("\n平台特定配置位置:")
    
    switch runtime.GOOS {
    case "windows":
        fmt.Println("用户配置: %APPDATA%\\NuGet\\NuGet.Config")
        fmt.Println("系统配置: %ProgramData%\\NuGet\\NuGet.Config")
        
        if appdata := os.Getenv("APPDATA"); appdata != "" {
            fmt.Printf("解析的用户路径: %s\n", filepath.Join(appdata, "NuGet", "NuGet.Config"))
        }
        
        if programdata := os.Getenv("ProgramData"); programdata != "" {
            fmt.Printf("解析的系统路径: %s\n", filepath.Join(programdata, "NuGet", "NuGet.Config"))
        }
        
    case "darwin":
        fmt.Println("用户配置: ~/Library/Application Support/NuGet/NuGet.Config")
        fmt.Println("系统配置: /Library/Application Support/NuGet/NuGet.Config")
        
        if home := os.Getenv("HOME"); home != "" {
            fmt.Printf("解析的用户路径: %s\n", filepath.Join(home, "Library", "Application Support", "NuGet", "NuGet.Config"))
        }
        
    default: // Linux 和其他 Unix 系统
        fmt.Println("用户配置: ~/.config/NuGet/NuGet.Config")
        fmt.Println("系统配置: /etc/NuGet/NuGet.Config")
        
        if home := os.Getenv("HOME"); home != "" {
            fmt.Printf("解析的用户路径: %s\n", filepath.Join(home, ".config", "NuGet", "NuGet.Config"))
        }
        
        // 检查 XDG_CONFIG_HOME
        if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
            fmt.Printf("XDG 配置: %s\n", filepath.Join(xdgConfig, "NuGet", "NuGet.Config"))
        }
    }
}

func categorizeConfigPath(configPath string) string {
    absPath, _ := filepath.Abs(configPath)
    
    // 检查是否在当前目录或子目录中
    if cwd, err := os.Getwd(); err == nil {
        if rel, err := filepath.Rel(cwd, absPath); err == nil && !filepath.IsAbs(rel) {
            return "项目/本地"
        }
    }
    
    // 检查是否在用户目录中
    if home := os.Getenv("HOME"); home != "" {
        if rel, err := filepath.Rel(home, absPath); err == nil && !filepath.IsAbs(rel) {
            return "用户"
        }
    }
    
    // 检查常见系统路径
    systemPaths := []string{"/etc", "/usr/local/etc", "/opt"}
    for _, sysPath := range systemPaths {
        if rel, err := filepath.Rel(sysPath, absPath); err == nil && !filepath.IsAbs(rel) {
            return "系统"
        }
    }
    
    return "其他"
}

func isReadable(path string) bool {
    file, err := os.Open(path)
    if err != nil {
        return false
    }
    file.Close()
    return true
}
```

## 关键概念

### 搜索顺序

库按以下顺序搜索配置文件：

1. **当前目录**: `./NuGet.Config`
2. **父目录**: 向上遍历目录树
3. **用户配置**: 平台特定用户目录
4. **系统配置**: 平台特定系统目录

### 平台差异

- **Windows**: 使用 `%APPDATA%` 和 `%ProgramData%`
- **macOS**: 使用 `~/Library/Application Support/` 和 `/Library/Application Support/`
- **Linux/Unix**: 使用 `~/.config/`（或 `$XDG_CONFIG_HOME`）和 `/etc/`

### 最佳实践

1. **优雅处理缺失文件**: 始终检查 `IsNotFoundError`
2. **提供回退**: 当未找到配置时有策略
3. **尊重层次结构**: 项目配置覆盖全局配置
4. **检查文件权限**: 确保文件在解析前可读
5. **使用适当的发现方法**: 根据需要选择单个文件或所有文件

## 常见模式

### 模式 1: 查找或创建

```go
config, configPath, err := api.FindAndParseConfig()
if err != nil {
    // 如果未找到则创建默认配置
    config = api.CreateDefaultConfig()
    configPath = "NuGet.Config"
    api.SaveConfig(config, configPath)
}
```

### 模式 2: 分层搜索

```go
// 首先尝试项目特定配置
if projectConfig, err := api.FindProjectConfig("."); err == nil {
    return api.ParseFromFile(projectConfig)
}

// 回退到全局配置
if globalConfig, err := api.FindConfigFile(); err == nil {
    return api.ParseFromFile(globalConfig)
}

// 最后创建默认配置
return api.CreateDefaultConfig(), nil
```

### 模式 3: 多配置合并

```go
configs := api.FindAllConfigFiles()
var mergedSources []types.PackageSource

for _, configPath := range configs {
    if config, err := api.ParseFromFile(configPath); err == nil {
        mergedSources = append(mergedSources, config.PackageSources.Add...)
    }
}
```

## 下一步

掌握配置发现后：

1. 学习 [创建配置](./creating-configs.md) 来生成新配置
2. 探索 [基本解析](./basic-parsing.md) 来理解解析详情
3. 研究 [修改配置](./modifying-configs.md) 来更新找到的配置

本指南为在不同场景和平台中查找 NuGet 配置文件提供了全面的示例。
