package editor

import (
	"strings"
	"testing"

	"github.com/scagogogo/nuget-config-parser/pkg/parser"
)

const testConfig = `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="local" value="C:\LocalPackages" />
  </packageSources>
  <config>
    <add key="globalPackagesFolder" value="C:\packages" />
  </config>
</configuration>`

func TestNewConfigEditor(t *testing.T) {
	// 创建位置感知解析器
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromContentWithPositions([]byte(testConfig))
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 创建编辑器
	editor := NewConfigEditor(parseResult)
	if editor == nil {
		t.Fatal("编辑器为空")
	}

	// 验证配置对象
	config := editor.GetConfig()
	if config == nil {
		t.Fatal("配置对象为空")
	}

	if len(config.PackageSources.Add) != 2 {
		t.Errorf("期望2个包源，实际得到%d个", len(config.PackageSources.Add))
	}

	// 验证位置信息
	positions := editor.GetPositions()
	if len(positions) == 0 {
		t.Fatal("未跟踪到任何位置")
	}

	// 检查是否找到了关键元素
	foundPackageSources := false
	foundConfig := false
	for path := range positions {
		if strings.Contains(path, "packageSources") {
			foundPackageSources = true
		}
		if strings.Contains(path, "config") {
			foundConfig = true
		}
	}

	if !foundPackageSources {
		t.Error("未找到packageSources元素")
	}
	if !foundConfig {
		t.Error("未找到config元素")
	}
}

func TestAddPackageSource(t *testing.T) {
	// 创建位置感知解析器
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromContentWithPositions([]byte(testConfig))
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 创建编辑器
	editor := NewConfigEditor(parseResult)

	// 添加新包源
	err = editor.AddPackageSource("test-source", "https://test.com/v3/index.json", "3")
	if err != nil {
		t.Fatalf("添加包源失败: %v", err)
	}

	// 验证内存中的配置已更新
	config := editor.GetConfig()
	if len(config.PackageSources.Add) != 3 {
		t.Errorf("期望3个包源，实际得到%d个", len(config.PackageSources.Add))
	}

	// 查找新添加的包源
	found := false
	for _, source := range config.PackageSources.Add {
		if source.Key == "test-source" {
			found = true
			if source.Value != "https://test.com/v3/index.json" {
				t.Errorf("包源URL不正确，期望https://test.com/v3/index.json，实际得到%s", source.Value)
			}
			if source.ProtocolVersion != "3" {
				t.Errorf("协议版本不正确，期望3，实际得到%s", source.ProtocolVersion)
			}
			break
		}
	}

	if !found {
		t.Error("未找到新添加的包源")
	}

	// 应用编辑并验证结果
	modifiedContent, err := editor.ApplyEdits()
	if err != nil {
		t.Fatalf("应用编辑失败: %v", err)
	}

	modifiedStr := string(modifiedContent)
	if !strings.Contains(modifiedStr, "test-source") {
		t.Error("修改后的内容中未找到新包源")
	}
	if !strings.Contains(modifiedStr, "https://test.com/v3/index.json") {
		t.Error("修改后的内容中未找到新包源的URL")
	}
}

func TestUpdatePackageSourceURL(t *testing.T) {
	// 创建位置感知解析器
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromContentWithPositions([]byte(testConfig))
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 创建编辑器
	editor := NewConfigEditor(parseResult)

	// 更新包源URL
	err = editor.UpdatePackageSourceURL("local", "D:\\NewLocalPackages")
	if err != nil {
		t.Fatalf("更新包源URL失败: %v", err)
	}

	// 验证内存中的配置已更新
	config := editor.GetConfig()
	found := false
	for _, source := range config.PackageSources.Add {
		if source.Key == "local" {
			found = true
			if source.Value != "D:\\NewLocalPackages" {
				t.Errorf("包源URL未正确更新，期望D:\\NewLocalPackages，实际得到%s", source.Value)
			}
			break
		}
	}

	if !found {
		t.Error("未找到要更新的包源")
	}

	// 应用编辑并验证结果
	modifiedContent, err := editor.ApplyEdits()
	if err != nil {
		t.Fatalf("应用编辑失败: %v", err)
	}

	modifiedStr := string(modifiedContent)
	if !strings.Contains(modifiedStr, "D:\\NewLocalPackages") {
		t.Error("修改后的内容中未找到更新的URL")
	}
	if strings.Contains(modifiedStr, "C:\\LocalPackages") {
		t.Error("修改后的内容中仍包含旧的URL")
	}
}

func TestRemovePackageSource(t *testing.T) {
	// 创建位置感知解析器
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromContentWithPositions([]byte(testConfig))
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 创建编辑器
	editor := NewConfigEditor(parseResult)

	// 删除包源
	err = editor.RemovePackageSource("local")
	if err != nil {
		t.Fatalf("删除包源失败: %v", err)
	}

	// 验证内存中的配置已更新
	config := editor.GetConfig()
	if len(config.PackageSources.Add) != 1 {
		t.Errorf("期望1个包源，实际得到%d个", len(config.PackageSources.Add))
	}

	// 确保删除的包源不存在
	for _, source := range config.PackageSources.Add {
		if source.Key == "local" {
			t.Error("删除的包源仍然存在")
		}
	}

	// 应用编辑并验证结果
	modifiedContent, err := editor.ApplyEdits()
	if err != nil {
		t.Fatalf("应用编辑失败: %v", err)
	}

	modifiedStr := string(modifiedContent)
	if strings.Contains(modifiedStr, `key="local"`) {
		t.Error("修改后的内容中仍包含已删除的包源")
	}
}

func TestMultipleEdits(t *testing.T) {
	// 创建位置感知解析器
	positionAwareParser := parser.NewPositionAwareParser()
	parseResult, err := positionAwareParser.ParseFromContentWithPositions([]byte(testConfig))
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	// 创建编辑器
	editor := NewConfigEditor(parseResult)

	// 执行多个编辑操作
	err = editor.AddPackageSource("new-source", "https://new.com/v3/index.json", "3")
	if err != nil {
		t.Fatalf("添加包源失败: %v", err)
	}

	err = editor.UpdatePackageSourceURL("nuget.org", "https://updated.nuget.org/v3/index.json")
	if err != nil {
		t.Fatalf("更新包源URL失败: %v", err)
	}

	err = editor.RemovePackageSource("local")
	if err != nil {
		t.Fatalf("删除包源失败: %v", err)
	}

	// 验证内存中的配置
	config := editor.GetConfig()
	if len(config.PackageSources.Add) != 2 {
		t.Errorf("期望2个包源，实际得到%d个", len(config.PackageSources.Add))
	}

	// 应用编辑并验证结果
	modifiedContent, err := editor.ApplyEdits()
	if err != nil {
		t.Fatalf("应用编辑失败: %v", err)
	}

	modifiedStr := string(modifiedContent)

	// 验证新添加的包源
	if !strings.Contains(modifiedStr, "new-source") {
		t.Error("修改后的内容中未找到新添加的包源")
	}

	// 验证更新的URL
	if !strings.Contains(modifiedStr, "https://updated.nuget.org/v3/index.json") {
		t.Error("修改后的内容中未找到更新的URL")
	}

	// 验证删除的包源
	if strings.Contains(modifiedStr, `key="local"`) {
		t.Error("修改后的内容中仍包含已删除的包源")
	}
}
