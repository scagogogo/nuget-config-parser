package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/editor"
	"github.com/scagogogo/nuget-config-parser/pkg/parser"
)

func main() {
	fmt.Println("=== NuGet配置位置感知编辑示例 ===\n")

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "nuget-position-aware-example-")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 创建初始配置文件
	initialConfig := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="localSource" value="/tmp/NuGet/LocalPackages" />
  </packageSources>
  <disabledPackageSources>
    <add key="localSource" value="true" />
  </disabledPackageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="$HOME/.nuget/packages" />
  </config>
</configuration>`

	// 写入初始配置
	err = os.WriteFile(configPath, []byte(initialConfig), 0644)
	if err != nil {
		log.Fatalf("写入配置文件失败: %v", err)
	}

	fmt.Printf("已创建初始配置文件: %s\n\n", configPath)

	// 显示原始内容
	fmt.Println("原始配置内容:")
	fmt.Println(initialConfig)
	fmt.Println()

	// 使用位置感知解析器解析配置
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromFileWithPositions(configPath)
	if err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	fmt.Println("解析成功，位置信息:")
	for path, elemPos := range parseResult.Positions {
		fmt.Printf("  %s: 行%d-行%d, 偏移%d-%d\n",
			path,
			elemPos.Range.Start.Line,
			elemPos.Range.End.Line,
			elemPos.Range.Start.Offset,
			elemPos.Range.End.Offset)
	}
	fmt.Println()

	// 创建位置感知编辑器
	posEditor := editor.NewConfigEditor(parseResult)

	// 执行编辑操作
	fmt.Println("执行编辑操作:")

	// 1. 添加新的包源
	fmt.Println("1. 添加新包源 'company-internal'")
	err = posEditor.AddPackageSource("company-internal", "https://nuget.company.com/v3/index.json", "3")
	if err != nil {
		log.Printf("添加包源失败: %v", err)
	}

	// 2. 更新现有包源的URL
	fmt.Println("2. 更新 'localSource' 的URL")
	err = posEditor.UpdatePackageSourceURL("localSource", "/new/path/to/packages")
	if err != nil {
		log.Printf("更新包源URL失败: %v", err)
	}

	// 3. 更新协议版本
	fmt.Println("3. 更新 'nuget.org' 的协议版本")
	err = posEditor.UpdatePackageSourceVersion("nuget.org", "3")
	if err != nil {
		log.Printf("更新协议版本失败: %v", err)
	}

	fmt.Println()

	// 应用编辑并获取结果
	fmt.Println("应用编辑操作...")
	modifiedContent, err := posEditor.ApplyEdits()
	if err != nil {
		log.Fatalf("应用编辑失败: %v", err)
	}

	// 显示修改后的内容
	fmt.Println("修改后的配置内容:")
	fmt.Println(string(modifiedContent))
	fmt.Println()

	// 保存修改后的配置
	modifiedConfigPath := filepath.Join(tempDir, "NuGet.Modified.Config")
	err = os.WriteFile(modifiedConfigPath, modifiedContent, 0644)
	if err != nil {
		log.Fatalf("保存修改后的配置失败: %v", err)
	}

	fmt.Printf("修改后的配置已保存到: %s\n", modifiedConfigPath)

	// 验证修改后的配置可以正确解析
	fmt.Println("\n验证修改后的配置:")
	verifyParser := parser.NewConfigParser()
	verifiedConfig, err := verifyParser.ParseFromFile(modifiedConfigPath)
	if err != nil {
		log.Fatalf("验证修改后的配置失败: %v", err)
	}

	fmt.Printf("验证成功！修改后的配置包含 %d 个包源:\n", len(verifiedConfig.PackageSources.Add))
	for _, source := range verifiedConfig.PackageSources.Add {
		fmt.Printf("  - %s: %s", source.Key, source.Value)
		if source.ProtocolVersion != "" {
			fmt.Printf(" (协议版本: %s)", source.ProtocolVersion)
		}
		fmt.Println()
	}

	fmt.Println("\n=== 位置感知编辑示例完成 ===")
	fmt.Println("\n优势:")
	fmt.Println("- 保持原始文件格式和缩进")
	fmt.Println("- 最小化diff，只修改必要的部分")
	fmt.Println("- 精确的位置控制")
	fmt.Println("- 同时更新内存对象和文本内容")
}
