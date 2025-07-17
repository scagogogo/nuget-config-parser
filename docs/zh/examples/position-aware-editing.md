# 位置感知编辑示例

本示例演示高级位置感知编辑技术，用于对 NuGet 配置文件进行精确修改，同时保持格式并最小化差异。

## 概述

位置感知编辑提供：
- 对特定文件位置的外科手术式修改
- 保持原始格式和注释
- 版本控制的最小差异
- 对更改的精确控制

## 示例 1: 基本位置感知编辑

简单的位置感知修改：

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
    configPath := "NuGet.Config"
    
    fmt.Println("=== 基本位置感知编辑 ===")
    
    // 解析时跟踪位置
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("位置解析失败: %v", err)
    }
    
    fmt.Printf("原始文件大小: %d 字节\n", len(parseResult.Content))
    
    // 创建编辑器
    editor := api.CreateConfigEditor(parseResult)
    
    // 显示当前配置
    config := editor.GetConfig()
    fmt.Printf("当前包源: %d\n", len(config.PackageSources.Add))
    
    // 进行精确更改
    fmt.Println("\n进行位置感知更改...")
    
    // 添加新包源
    err = editor.AddPackageSource("company-feed", "https://nuget.company.com/v3/index.json", "3")
    if err != nil {
        log.Printf("添加源失败: %v", err)
    } else {
        fmt.Println("✅ 添加了 company-feed")
    }
    
    // 更新现有源 URL
    err = editor.UpdatePackageSourceURL("nuget.org", "https://api.nuget.org/v3/index.json")
    if err != nil {
        log.Printf("更新 URL 失败: %v", err)
    } else {
        fmt.Println("✅ 更新了 nuget.org URL")
    }
    
    // 应用所有更改
    fmt.Println("\n应用更改...")
    modifiedContent, err := editor.ApplyEdits()
    if err != nil {
        log.Fatalf("应用编辑失败: %v", err)
    }
    
    fmt.Printf("修改后文件大小: %d 字节\n", len(modifiedContent))
    fmt.Printf("大小变化: %+d 字节\n", len(modifiedContent)-len(parseResult.Content))
    
    // 保存修改后的内容
    err = os.WriteFile(configPath, modifiedContent, 0644)
    if err != nil {
        log.Fatalf("保存文件失败: %v", err)
    }
    
    fmt.Println("✅ 更改已成功应用，差异最小！")
}
```

## 示例 2: 批量位置感知操作

高效执行多个操作：

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
    configPath := "NuGet.Config"
    
    fmt.Println("=== 批量位置感知操作 ===")
    
    // 解析时跟踪位置
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 创建编辑器
    editor := api.CreateConfigEditor(parseResult)
    
    // 定义批量操作
    operations := []struct {
        name string
        op   func() error
    }{
        {"添加生产源", func() error {
            return editor.AddPackageSource("prod-feed", "https://prod.company.com/nuget", "3")
        }},
        {"添加预发布源", func() error {
            return editor.AddPackageSource("staging-feed", "https://staging.company.com/nuget", "3")
        }},
        {"添加开发源", func() error {
            return editor.AddPackageSource("dev-feed", "https://dev.company.com/nuget", "3")
        }},
        {"更新 nuget.org 协议", func() error {
            return editor.UpdatePackageSourceVersion("nuget.org", "3")
        }},
        {"移除旧源", func() error {
            return editor.RemovePackageSource("old-feed")
        }},
    }
    
    // 执行批量操作
    fmt.Println("执行批量操作:")
    successCount := 0
    
    for _, op := range operations {
        err := op.op()
        if err != nil {
            fmt.Printf("  ❌ %s: %v\n", op.name, err)
        } else {
            fmt.Printf("  ✅ %s\n", op.name)
            successCount++
        }
    }
    
    fmt.Printf("\n成功完成 %d/%d 操作\n", successCount, len(operations))
    
    // 一次性应用所有更改
    if successCount > 0 {
        fmt.Println("应用所有更改...")
        modifiedContent, err := editor.ApplyEdits()
        if err != nil {
            log.Fatalf("应用编辑失败: %v", err)
        }
        
        // 保存更改
        err = os.WriteFile(configPath, modifiedContent, 0644)
        if err != nil {
            log.Fatalf("保存失败: %v", err)
        }
        
        fmt.Println("✅ 所有更改已成功应用！")
        
        // 显示差异统计
        showDiffStats(parseResult.Content, modifiedContent)
    }
}

func showDiffStats(original, modified []byte) {
    fmt.Println("\n=== 差异统计 ===")
    fmt.Printf("原始大小: %d 字节\n", len(original))
    fmt.Printf("修改后大小: %d 字节\n", len(modified))
    fmt.Printf("大小变化: %+d 字节\n", len(modified)-len(original))
    
    // 简单的行数比较
    originalLines := countLines(original)
    modifiedLines := countLines(modified)
    fmt.Printf("原始行数: %d\n", originalLines)
    fmt.Printf("修改后行数: %d\n", modifiedLines)
    fmt.Printf("行数变化: %+d 行\n", modifiedLines-originalLines)
}

func countLines(content []byte) int {
    lines := 1
    for _, b := range content {
        if b == '\n' {
            lines++
        }
    }
    return lines
}
```

## 示例 3: 位置信息分析

分析和使用位置信息：

```go
package main

import (
    "fmt"
    "log"
    "sort"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    configPath := "NuGet.Config"
    
    fmt.Println("=== 位置信息分析 ===")
    
    // 解析时跟踪位置
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 创建编辑器以访问位置
    editor := api.CreateConfigEditor(parseResult)
    positions := editor.GetPositions()
    
    fmt.Printf("找到 %d 个定位元素\n", len(positions))
    
    // 分析位置信息
    analyzePositions(positions)
    
    // 按行号查找元素
    findElementsByLine(positions, 10) // 查找第10行附近的元素
    
    // 显示元素层次结构
    showElementHierarchy(positions)
    
    // 基于位置演示目标编辑
    performTargetedEditing(editor, positions)
}

func analyzePositions(positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== 位置分析 ===")
    
    // 按标签名分组
    tagCounts := make(map[string]int)
    lineRanges := make(map[string][]int)
    
    for path, pos := range positions {
        tagCounts[pos.TagName]++
        lineRanges[pos.TagName] = append(lineRanges[pos.TagName], pos.Range.Start.Line)
    }
    
    fmt.Println("按标签类型的元素:")
    for tag, count := range tagCounts {
        fmt.Printf("  %s: %d 个元素\n", tag, count)
    }
    
    fmt.Println("\n行分布:")
    for tag, lines := range lineRanges {
        if len(lines) > 0 {
            sort.Ints(lines)
            fmt.Printf("  %s: 行 %d-%d\n", tag, lines[0], lines[len(lines)-1])
        }
    }
}

func findElementsByLine(positions map[string]*parser.ElementPosition, targetLine int) {
    fmt.Printf("\n=== 第 %d 行附近的元素 ===\n", targetLine)
    
    type ElementInfo struct {
        Path string
        Pos  *parser.ElementPosition
        Dist int
    }
    
    var nearby []ElementInfo
    
    for path, pos := range positions {
        startDist := abs(pos.Range.Start.Line - targetLine)
        endDist := abs(pos.Range.End.Line - targetLine)
        minDist := min(startDist, endDist)
        
        if minDist <= 5 { // 5行内
            nearby = append(nearby, ElementInfo{path, pos, minDist})
        }
    }
    
    // 按距离排序
    sort.Slice(nearby, func(i, j int) bool {
        return nearby[i].Dist < nearby[j].Dist
    })
    
    for _, elem := range nearby {
        fmt.Printf("  %s (%s) 在行 %d-%d (距离: %d)\n",
            elem.Path, elem.Pos.TagName,
            elem.Pos.Range.Start.Line, elem.Pos.Range.End.Line,
            elem.Dist)
    }
}

func showElementHierarchy(positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== 元素层次结构 ===")
    
    // 按层次级别分组（路径中'/'的数量）
    levels := make(map[int][]string)
    
    for path := range positions {
        level := countChar(path, '/')
        levels[level] = append(levels[level], path)
    }
    
    // 排序级别
    var sortedLevels []int
    for level := range levels {
        sortedLevels = append(sortedLevels, level)
    }
    sort.Ints(sortedLevels)
    
    for _, level := range sortedLevels {
        fmt.Printf("级别 %d:\n", level)
        sort.Strings(levels[level])
        for _, path := range levels[level] {
            indent := strings.Repeat("  ", level+1)
            pos := positions[path]
            fmt.Printf("%s%s (%s) [%d:%d-%d:%d]\n",
                indent, path, pos.TagName,
                pos.Range.Start.Line, pos.Range.Start.Column,
                pos.Range.End.Line, pos.Range.End.Column)
        }
    }
}

func performTargetedEditing(editor *editor.ConfigEditor, positions map[string]*parser.ElementPosition) {
    fmt.Println("\n=== 基于位置的目标编辑 ===")
    
    // 查找包源元素
    var packageSourcePaths []string
    for path, pos := range positions {
        if pos.TagName == "add" && strings.Contains(path, "packageSources") {
            packageSourcePaths = append(packageSourcePaths, path)
        }
    }
    
    fmt.Printf("找到 %d 个包源元素\n", len(packageSourcePaths))
    
    // 按行号排序以便可预测处理
    sort.Slice(packageSourcePaths, func(i, j int) bool {
        return positions[packageSourcePaths[i]].Range.Start.Line <
            positions[packageSourcePaths[j]].Range.Start.Line
    })
    
    // 显示我们找到的内容
    for i, path := range packageSourcePaths {
        pos := positions[path]
        fmt.Printf("  %d. %s 在行 %d\n", i+1, path, pos.Range.Start.Line)
        
        // 如果可用，显示属性
        if len(pos.Attributes) > 0 {
            for attrName, attrValue := range pos.Attributes {
                fmt.Printf("     %s=\"%s\"\n", attrName, attrValue)
            }
        }
    }
    
    // 执行目标修改
    if len(packageSourcePaths) > 0 {
        fmt.Println("\n执行目标修改...")
        
        // 添加新源（将被适当定位）
        err := editor.AddPackageSource("targeted-source", "https://targeted.example.com", "3")
        if err != nil {
            fmt.Printf("添加目标源失败: %v\n", err)
        } else {
            fmt.Println("✅ 添加了目标源")
        }
    }
}

// 辅助函数
func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func countChar(s string, c rune) int {
    count := 0
    for _, r := range s {
        if r == c {
            count++
        }
    }
    return count
}
```

## 示例 4: 差异最小化策略

最小化版本控制差异的技术：

```go
package main

import (
    "fmt"
    "log"
    "os"
    "strings"
    
    "github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
    api := nuget.NewAPI()
    configPath := "NuGet.Config"
    
    fmt.Println("=== 差异最小化策略 ===")
    
    // 创建备份以便比较
    originalContent, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalf("读取原始文件失败: %v", err)
    }
    
    // 解析时跟踪位置
    parseResult, err := api.ParseFromFileWithPositions(configPath)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }
    
    // 创建编辑器
    editor := api.CreateConfigEditor(parseResult)
    
    // 演示不同的编辑策略
    strategies := []struct {
        name string
        edit func() error
    }{
        {"添加单个源", func() error {
            return editor.AddPackageSource("minimal-source", "https://minimal.example.com", "3")
        }},
        {"更新现有 URL", func() error {
            return editor.UpdatePackageSourceURL("nuget.org", "https://api.nuget.org/v3/index.json")
        }},
        {"更新协议版本", func() error {
            return editor.UpdatePackageSourceVersion("nuget.org", "3")
        }},
    }
    
    for i, strategy := range strategies {
        fmt.Printf("\n--- 策略 %d: %s ---\n", i+1, strategy.name)
        
        // 应用单个更改
        err := strategy.edit()
        if err != nil {
            fmt.Printf("❌ 失败: %v\n", err)
            continue
        }
        
        // 应用并分析差异
        modifiedContent, err := editor.ApplyEdits()
        if err != nil {
            fmt.Printf("❌ 应用失败: %v\n", err)
            continue
        }
        
        // 分析差异
        analyzeDiff(originalContent, modifiedContent, strategy.name)
        
        // 为下一个策略重置（在实际使用中，您会使用新的编辑器）
        // 为了演示目的，我们将继续累积更改
        originalContent = modifiedContent
    }
    
    // 保存最终结果
    finalContent, _ := editor.ApplyEdits()
    err = os.WriteFile(configPath, finalContent, 0644)
    if err != nil {
        log.Fatalf("保存失败: %v", err)
    }
    
    fmt.Println("\n✅ 所有策略都以最小差异应用！")
}

func analyzeDiff(original, modified []byte, strategyName string) {
    fmt.Printf("%s 的差异分析:\n", strategyName)
    
    originalLines := strings.Split(string(original), "\n")
    modifiedLines := strings.Split(string(modified), "\n")
    
    // 简单差异分析
    added := 0
    removed := 0
    changed := 0
    
    maxLines := max(len(originalLines), len(modifiedLines))
    
    for i := 0; i < maxLines; i++ {
        var origLine, modLine string
        
        if i < len(originalLines) {
            origLine = originalLines[i]
        }
        if i < len(modifiedLines) {
            modLine = modifiedLines[i]
        }
        
        if origLine == "" && modLine != "" {
            added++
        } else if origLine != "" && modLine == "" {
            removed++
        } else if origLine != modLine {
            changed++
        }
    }
    
    fmt.Printf("  添加行数: %d\n", added)
    fmt.Printf("  删除行数: %d\n", removed)
    fmt.Printf("  更改行数: %d\n", changed)
    fmt.Printf("  总差异影响: %d 行\n", added+removed+changed)
    
    // 计算差异效率
    totalLines := len(originalLines)
    diffPercentage := float64(added+removed+changed) / float64(totalLines) * 100
    fmt.Printf("  差异百分比: %.2f%%\n", diffPercentage)
    
    if diffPercentage < 5.0 {
        fmt.Printf("  ✅ 优秀的差异效率 (< 5%%)\n")
    } else if diffPercentage < 15.0 {
        fmt.Printf("  ✅ 良好的差异效率 (< 15%%)\n")
    } else {
        fmt.Printf("  ⚠️  高差异影响 (>= 15%%)\n")
    }
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

## 关键优势

### 最小差异
- 只更改文件的必要部分
- 保持原始格式和结构
- 维护注释和空白
- 减少版本控制噪音

### 精确控制
- 对修改的确切定位
- 对特定元素的外科手术式更改
- 可预测的修改行为
- 对输出的细粒度控制

### 格式保持
- 维护原始缩进
- 保留 XML 注释
- 保持一致的间距
- 尊重现有结构

## 最佳实践

1. **一次解析，多次编辑**: 对多个编辑使用单个解析结果
2. **批量操作**: 将相关更改分组在一起
3. **验证结果**: 始终验证修改后的 XML 是否有效
4. **优雅处理错误**: 检查每个操作的错误
5. **测试往返**: 确保解析 → 编辑 → 解析循环正常工作

## 常见用例

### 自动化配置更新
```go
// 以最小影响更新配置
parseResult, _ := api.ParseFromFileWithPositions(configPath)
editor := api.CreateConfigEditor(parseResult)
editor.UpdatePackageSourceURL("source", "new-url")
modifiedContent, _ := editor.ApplyEdits()
```

### CI/CD 配置管理
```go
// 应用环境特定更改
editor.AddPackageSource("ci-feed", ciURL, "3")
editor.UpdatePackageSourceVersion("nuget.org", "3")
modifiedContent, _ := editor.ApplyEdits()
```

### 配置迁移
```go
// 以保留格式迁移配置
editor.UpdatePackageSourceURL("old-source", "new-url")
editor.UpdatePackageSourceVersion("old-source", "3")
modifiedContent, _ := editor.ApplyEdits()
```

## 下一步

掌握位置感知编辑后：

1. 探索 [编辑器 API](/api/editor) 了解详细的方法文档
2. 学习 [配置结构](../guide/configuration.md)
3. 研究 [序列化](./serialization.md) 了解 XML 输出

位置感知编辑是库的最高级功能，为配置修改提供外科手术式精度，同时保持配置文件的完整性和可读性。
