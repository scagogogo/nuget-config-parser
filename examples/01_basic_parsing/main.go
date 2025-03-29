// 该示例演示了如何使用NuGet配置解析器API解析配置文件
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
	// 1. 创建一个临时文件用于演示解析功能
	// -----------------------------------------------

	// 准备一个示例NuGet配置内容
	configContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="localSource" value="C:\NuGet\LocalPackages" />
  </packageSources>
  <disabledPackageSources>
    <add key="localSource" value="true" />
  </disabledPackageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
  </config>
</configuration>`

	// 创建临时文件
	tempDir, err := ioutil.TempDir("", "nuget-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	// 保存配置到临时文件
	configPath := filepath.Join(tempDir, "NuGet.Config")
	if err := ioutil.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		log.Fatalf("写入配置文件失败: %v", err)
	}

	fmt.Printf("已创建临时配置文件: %s\n", configPath)
	// 输出示例: 已创建临时配置文件: /tmp/nuget-example-123456/NuGet.Config

	// 2. 使用API解析配置文件
	// -----------------------------------------------

	// 创建API实例
	api := nuget.NewAPI()

	// 从文件解析配置
	fmt.Println("\n从文件解析配置:")
	config, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 3. 访问解析后的配置内容
	// -----------------------------------------------

	// 输出包源信息
	fmt.Println("\n包源列表:")
	for _, source := range config.PackageSources.Add {
		fmt.Printf("  - %s: %s\n", source.Key, source.Value)
		if source.ProtocolVersion != "" {
			fmt.Printf("    协议版本: %s\n", source.ProtocolVersion)
		}
	}
	// 输出示例:
	// 包源列表:
	//   - nuget.org: https://api.nuget.org/v3/index.json
	//     协议版本: 3
	//   - localSource: C:\NuGet\LocalPackages

	// 输出禁用的包源
	if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
		fmt.Println("\n禁用的包源:")
		for _, source := range config.DisabledPackageSources.Add {
			fmt.Printf("  - %s\n", source.Key)
		}
	}
	// 输出示例:
	// 禁用的包源:
	//   - localSource

	// 输出活跃包源
	if config.ActivePackageSource != nil {
		fmt.Printf("\n活跃包源: %s (%s)\n",
			config.ActivePackageSource.Add.Key,
			config.ActivePackageSource.Add.Value)
	}
	// 输出示例:
	// 活跃包源: nuget.org (https://api.nuget.org/v3/index.json)

	// 输出配置选项
	if config.Config != nil && len(config.Config.Add) > 0 {
		fmt.Println("\n配置选项:")
		for _, option := range config.Config.Add {
			fmt.Printf("  - %s: %s\n", option.Key, option.Value)
		}
	}
	// 输出示例:
	// 配置选项:
	//   - globalPackagesFolder: %USERPROFILE%\.nuget\packages

	// 4. 从字符串解析配置
	// -----------------------------------------------
	fmt.Println("\n从字符串解析配置:")
	configFromString, err := api.ParseFromString(configContent)
	if err != nil {
		log.Fatalf("从字符串解析配置失败: %v", err)
	}
	fmt.Printf("成功从字符串解析配置，包含 %d 个包源\n",
		len(configFromString.PackageSources.Add))
	// 输出示例:
	// 从字符串解析配置:
	// 成功从字符串解析配置，包含 2 个包源
}
