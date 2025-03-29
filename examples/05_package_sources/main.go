// 该示例演示了如何使用NuGet配置解析器API管理包源
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
	// 1. 创建临时目录用于示例
	// -----------------------------------------------
	tempDir, err := os.MkdirTemp("", "nuget-sources-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 创建API实例
	api := nuget.NewAPI()

	// 创建一个空配置
	config := &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{},
		},
	}

	// 2. 添加多种类型的包源
	// -----------------------------------------------
	fmt.Println("添加各种类型的包源:")

	// 2.1 添加标准的nuget.org源（V3 API）
	fmt.Println("- 添加官方NuGet源 (nuget.org)...")
	api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")

	// 2.2 添加公司内部源
	fmt.Println("- 添加公司内部源 (company-internal)...")
	api.AddPackageSource(config, "company-internal", "https://nuget.company.com/v3/index.json", "3")

	// 2.3 添加本地文件夹源
	fmt.Println("- 添加本地文件夹源 (local-packages)...")
	localPackagePath := filepath.Join("C:", "LocalPackages")
	if runtime.GOOS != "windows" {
		localPackagePath = "/tmp/LocalPackages"
	}
	api.AddPackageSource(config, "local-packages", localPackagePath, "")

	// 2.4 添加旧版NuGet源（V2 API）
	fmt.Println("- 添加旧版V2 API源 (legacy-v2)...")
	api.AddPackageSource(config, "legacy-v2", "https://packages.legacy.com/api/v2", "2")

	// 保存现有配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存配置文件失败: %v", err)
	}

	// 3. 列出所有包源
	// -----------------------------------------------
	fmt.Println("\n列出所有包源:")
	sources := api.GetAllPackageSources(config)
	for i, source := range sources {
		fmt.Printf("包源 #%d: %s (%s)", i+1, source.Key, source.Value)
		if source.ProtocolVersion != "" {
			fmt.Printf(", 协议版本: %s", source.ProtocolVersion)
		}
		fmt.Println()
	}
	// 输出示例:
	// 列出所有包源:
	// 包源 #1: nuget.org (https://api.nuget.org/v3/index.json), 协议版本: 3
	// 包源 #2: company-internal (https://nuget.company.com/v3/index.json), 协议版本: 3
	// 包源 #3: local-packages (C:\LocalPackages)
	// 包源 #4: legacy-v2 (https://packages.legacy.com/api/v2), 协议版本: 2

	// 4. 获取特定包源
	// -----------------------------------------------
	fmt.Println("\n获取特定包源:")

	specificSource := api.GetPackageSource(config, "company-internal")
	if specificSource != nil {
		fmt.Printf("找到包源: %s (%s), 协议版本: %s\n",
			specificSource.Key, specificSource.Value, specificSource.ProtocolVersion)
	} else {
		fmt.Println("未找到指定的包源")
	}
	// 输出示例:
	// 获取特定包源:
	// 找到包源: company-internal (https://nuget.company.com/v3/index.json), 协议版本: 3

	// 5. 修改现有包源
	// -----------------------------------------------
	fmt.Println("\n修改现有包源:")

	// 修改现有包源URL (更新local-packages路径)
	fmt.Println("- 修改本地包源路径...")
	api.AddPackageSource(config, "local-packages", "D:\\NuGetPackages", "")

	// 获取修改后的包源
	updatedSource := api.GetPackageSource(config, "local-packages")
	if updatedSource != nil {
		fmt.Printf("修改后的包源: %s (%s)\n",
			updatedSource.Key, updatedSource.Value)
	}
	// 输出示例:
	// 修改现有包源:
	// - 修改本地包源路径...
	// 修改后的包源: local-packages (D:\NuGetPackages)

	// 6. 禁用包源
	// -----------------------------------------------
	fmt.Println("\n禁用包源:")

	// 禁用V2旧版包源
	fmt.Println("- 禁用旧版V2源...")
	api.DisablePackageSource(config, "legacy-v2")

	// 检查是否被禁用
	isDisabled := api.IsPackageSourceDisabled(config, "legacy-v2")
	fmt.Printf("legacy-v2源已被禁用: %v\n", isDisabled)
	// 输出示例:
	// 禁用包源:
	// - 禁用旧版V2源...
	// legacy-v2源已被禁用: true

	// 7. 启用包源
	// -----------------------------------------------
	fmt.Println("\n启用包源:")

	// 重新启用被禁用的包源
	fmt.Println("- 重新启用legacy-v2源...")
	enabled := api.EnablePackageSource(config, "legacy-v2")
	fmt.Printf("源已启用: %v\n", enabled)

	// 检查是否已启用
	isDisabled = api.IsPackageSourceDisabled(config, "legacy-v2")
	fmt.Printf("legacy-v2源是否被禁用: %v\n", isDisabled)
	// 输出示例:
	// 启用包源:
	// - 重新启用legacy-v2源...
	// 源已启用: true
	// legacy-v2源是否被禁用: false

	// 8. 移除包源
	// -----------------------------------------------
	fmt.Println("\n移除包源:")

	// 移除旧版源
	fmt.Println("- 移除legacy-v2源...")
	removed := api.RemovePackageSource(config, "legacy-v2")
	fmt.Printf("源已移除: %v\n", removed)
	// 输出示例:
	// 移除包源:
	// - 移除legacy-v2源...
	// 源已移除: true

	// 9. 设置活跃包源
	// -----------------------------------------------
	fmt.Println("\n设置活跃包源:")

	// 设置公司内部源为活跃源
	fmt.Println("- 设置company-internal为活跃包源...")
	if err := api.SetActivePackageSource(config, "company-internal"); err != nil {
		log.Fatalf("设置活跃包源失败: %v", err)
	}

	// 保存配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存配置文件失败: %v", err)
	}

	// 10. 重新加载和验证配置
	// -----------------------------------------------
	fmt.Println("\n验证最终配置:")

	finalConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 验证包源数量 (应该是3个，因为移除了1个)
	fmt.Printf("最终包源数量: %d\n", len(finalConfig.PackageSources.Add))

	// 验证活跃包源
	if finalConfig.ActivePackageSource != nil {
		fmt.Printf("活跃包源: %s (%s)\n",
			finalConfig.ActivePackageSource.Add.Key,
			finalConfig.ActivePackageSource.Add.Value)
	}

	// 列出所有最终包源
	fmt.Println("\n最终包源列表:")
	for i, source := range finalConfig.PackageSources.Add {
		fmt.Printf("包源 #%d: %s (%s)\n", i+1, source.Key, source.Value)
	}
	// 输出示例:
	// 验证最终配置:
	// 最终包源数量: 3
	// 活跃包源: company-internal (https://nuget.company.com/v3/index.json)
	//
	// 最终包源列表:
	// 包源 #1: nuget.org (https://api.nuget.org/v3/index.json)
	// 包源 #2: company-internal (https://nuget.company.com/v3/index.json)
	// 包源 #3: local-packages (D:\NuGetPackages)

	// 11. 打印完整配置
	// -----------------------------------------------
	fmt.Println("\n最终配置文件内容:")
	xmlContent, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}
	fmt.Println(string(xmlContent))
}
