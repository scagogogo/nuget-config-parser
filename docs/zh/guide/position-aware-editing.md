# 位置感知编辑

位置感知编辑是 NuGet Config Parser 的高级功能，允许您对配置文件进行精确修改，同时保持原始格式并最小化版本控制差异。

## 概述

位置感知编辑提供：
- **精确修改** - 对特定文件位置进行外科手术式修改
- **格式保持** - 保留原始缩进、注释和空白
- **最小差异** - 减少版本控制中的噪音
- **智能编辑** - 基于XML结构的智能修改

## 工作原理

### 解析阶段
1. **位置跟踪** - 记录每个XML元素的精确位置
2. **结构映射** - 构建元素层次结构和位置映射
3. **属性定位** - 跟踪属性的开始和结束位置

### 编辑阶段
1. **变更计划** - 计划所有修改操作
2. **冲突检测** - 检测重叠或冲突的编辑
3. **顺序应用** - 按正确顺序应用所有变更

## 基本用法

### 创建编辑器

```go
package main

import (
    "log"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    
    // 解析时跟踪位置信息
    parseResult, err := api.ParseFromFileWithPositions("NuGet.Config")
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 创建位置感知编辑器
    editor := api.CreateConfigEditor(parseResult)
    
    // 现在可以进行精确编辑
}
```

### 基本编辑操作

```go
// 添加包源
err := editor.AddPackageSource("new-source", "https://example.com/nuget", "3")

// 更新包源URL
err := editor.UpdatePackageSourceURL("existing-source", "https://new-url.com")

// 移除包源
err := editor.RemovePackageSource("old-source")

// 应用所有编辑
modifiedContent, err := editor.ApplyEdits()
```

## 高级功能

### 批量编辑

```go
// 执行多个编辑操作
operations := []func() error{
    func() error { return editor.AddPackageSource("source1", "url1", "3") },
    func() error { return editor.AddPackageSource("source2", "url2", "3") },
    func() error { return editor.UpdatePackageSourceURL("old", "new-url") },
}

for _, op := range operations {
    if err := op(); err != nil {
        log.Printf("操作失败: %v", err)
    }
}

// 一次性应用所有变更
result, err := editor.ApplyEdits()
```

### 位置信息查询

```go
// 获取所有位置信息
positions := editor.GetPositions()

// 查找特定元素的位置
for path, pos := range positions {
    fmt.Printf("元素 %s 位于 %d:%d-%d:%d\n", 
        path, 
        pos.Range.Start.Line, pos.Range.Start.Column,
        pos.Range.End.Line, pos.Range.End.Column)
}
```

## 最佳实践

### 1. 最小化编辑范围

```go
// 好的做法：只修改需要的部分
editor.UpdatePackageSourceURL("source", "new-url")

// 避免：重新创建整个配置
```

### 2. 批量操作

```go
// 好的做法：批量编辑
editor.AddPackageSource("source1", "url1", "3")
editor.AddPackageSource("source2", "url2", "3")
result, _ := editor.ApplyEdits()

// 避免：多次应用编辑
```

### 3. 错误处理

```go
if err := editor.AddPackageSource("source", "url", "3"); err != nil {
    // 处理特定错误
    switch {
    case errors.IsValidationError(err):
        log.Printf("验证错误: %v", err)
    case errors.IsConflictError(err):
        log.Printf("冲突错误: %v", err)
    default:
        log.Printf("未知错误: %v", err)
    }
}
```

## 实际应用场景

### CI/CD 配置更新

```go
func updateCIConfig(configPath, environment string) error {
    api := nuget.NewAPI()
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        return err
    }
    
    editor := api.CreateConfigEditor(parseResult)
    
    // 根据环境更新配置
    switch environment {
    case "staging":
        editor.UpdatePackageSourceURL("main", "https://staging.nuget.com")
    case "production":
        editor.UpdatePackageSourceURL("main", "https://prod.nuget.com")
    }
    
    // 应用变更并保存
    content, err := editor.ApplyEdits()
    if err != nil {
        return err
    }
    
    return os.WriteFile(configPath, content, 0644)
}
```

### 配置迁移

```go
func migrateConfig(oldPath, newPath string) error {
    api := nuget.NewAPI()
    parseResult, err := api.ParseFromFileWithPositions(oldPath)
    if err != nil {
        return err
    }
    
    editor := api.CreateConfigEditor(parseResult)
    
    // 更新过时的URL
    migrations := map[string]string{
        "old-nuget.com": "new-nuget.com",
        "legacy.feed.com": "modern.feed.com",
    }
    
    for oldURL, newURL := range migrations {
        // 查找并更新所有匹配的源
        config := editor.GetConfig()
        for _, source := range config.PackageSources.Add {
            if strings.Contains(source.Value, oldURL) {
                newValue := strings.Replace(source.Value, oldURL, newURL, -1)
                editor.UpdatePackageSourceURL(source.Key, newValue)
            }
        }
    }
    
    content, err := editor.ApplyEdits()
    if err != nil {
        return err
    }
    
    return os.WriteFile(newPath, content, 0644)
}
```

## 性能考虑

### 内存使用

位置感知编辑会消耗额外内存来存储位置信息：
- 小文件（< 100KB）：影响可忽略
- 大文件（> 1MB）：考虑分批处理

### 编辑效率

```go
// 高效：单次解析，多次编辑
parseResult, _ := api.ParseFromFileWithPositions(path)
editor := api.CreateConfigEditor(parseResult)

editor.AddPackageSource("source1", "url1", "3")
editor.AddPackageSource("source2", "url2", "3")
editor.UpdatePackageSourceURL("old", "new")

result, _ := editor.ApplyEdits()

// 低效：多次解析
for _, change := range changes {
    parseResult, _ := api.ParseFromFileWithPositions(path)
    editor := api.CreateConfigEditor(parseResult)
    // ... 单个编辑
}
```

## 故障排除

### 常见问题

1. **位置冲突**
   ```go
   // 检查编辑冲突
   if err := editor.ValidateEdits(); err != nil {
       log.Printf("编辑冲突: %v", err)
   }
   ```

2. **格式问题**
   ```go
   // 验证生成的XML
   content, err := editor.ApplyEdits()
   if err != nil {
       return err
   }
   
   // 验证XML有效性
   _, err = api.ParseFromString(string(content))
   if err != nil {
       log.Printf("生成的XML无效: %v", err)
   }
   ```

3. **编码问题**
   ```go
   // 确保正确的UTF-8编码
   content, err := editor.ApplyEdits()
   if err != nil {
       return err
   }
   
   if !utf8.Valid(content) {
       log.Printf("警告: 内容包含无效的UTF-8字符")
   }
   ```

## 与其他功能的集成

### 与配置发现结合

```go
// 查找配置文件
configPath, err := api.FindConfigFile()
if err != nil {
    return err
}

// 进行位置感知编辑
parseResult, err := api.ParseFromFileWithPositions(configPath)
if err != nil {
    return err
}

editor := api.CreateConfigEditor(parseResult)
// ... 编辑操作
```

### 与验证结合

```go
// 编辑前验证
config := editor.GetConfig()
if err := api.ValidateConfig(config); err != nil {
    return fmt.Errorf("编辑前验证失败: %w", err)
}

// 应用编辑
content, err := editor.ApplyEdits()
if err != nil {
    return err
}

// 编辑后验证
newConfig, err := api.ParseFromString(string(content))
if err != nil {
    return fmt.Errorf("编辑后解析失败: %w", err)
}

if err := api.ValidateConfig(newConfig); err != nil {
    return fmt.Errorf("编辑后验证失败: %w", err)
}
```

## 总结

位置感知编辑是处理配置文件的强大工具，特别适用于：
- 自动化配置管理
- CI/CD 流水线
- 配置迁移工具
- 开发工具集成

通过保持原始格式和最小化差异，它确保了配置文件的可维护性和版本控制友好性。

## 下一步

- 查看 [API 参考](../api/editor.md) 了解详细的编辑器API
- 阅读 [示例](../examples/position-aware-editing.md) 获取更多实际用例
- 探索 [配置结构](./configuration.md) 了解配置文件格式
