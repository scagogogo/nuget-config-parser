// 该示例演示了如何使用NuGet配置解析器API管理全局配置选项
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
	// 1. 创建临时目录和配置文件
	// -----------------------------------------------
	tempDir, err := ioutil.TempDir("", "nuget-config-options-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 创建API实例
	api := nuget.NewAPI()

	// 2. 创建基本配置
	// -----------------------------------------------
	fmt.Println("创建基本配置:")

	// 创建包含官方源的新配置
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

	// 设置活跃包源
	if err := api.SetActivePackageSource(config, "nuget.org"); err != nil {
		log.Fatalf("设置活跃包源失败: %v", err)
	}

	// 3. 添加各种配置选项
	// -----------------------------------------------
	fmt.Println("\n添加各种配置选项:")

	// 3.1 设置全局包文件夹路径
	fmt.Println("- 设置全局包文件夹路径...")
	api.AddConfigOption(config, "globalPackagesFolder", "%USERPROFILE%\\.nuget\\packages")

	// 3.2 设置代理服务器
	fmt.Println("- 设置HTTP代理...")
	api.AddConfigOption(config, "http_proxy", "http://proxy.example.com:8080")
	api.AddConfigOption(config, "no_proxy", "localhost,127.0.0.1,.local")

	// 3.3 设置依赖版本行为
	fmt.Println("- 设置依赖版本行为...")
	api.AddConfigOption(config, "dependencyVersion", "Highest")

	// 3.4 设置包确认模式
	fmt.Println("- 设置包签名验证模式...")
	api.AddConfigOption(config, "signatureValidationMode", "require")

	// 3.5 设置允许的包格式
	fmt.Println("- 设置允许的包格式...")
	api.AddConfigOption(config, "allowedPackageFormats", "nupkg")

	// 3.6 禁用源浏览器
	fmt.Println("- 禁用源浏览器...")
	api.AddConfigOption(config, "disableSourceControlIntegration", "true")

	// 3.7 添加包还原设置
	fmt.Println("- 配置包还原选项...")
	api.AddConfigOption(config, "restorePackagesPath", "packages")
	api.AddConfigOption(config, "restoreIgnoreFailedSources", "true")

	// 4. 保存配置到文件
	// -----------------------------------------------
	fmt.Println("\n保存配置到文件...")
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存配置失败: %v", err)
	}
	fmt.Printf("配置已保存到: %s\n", configPath)

	// 5. 读取配置并展示配置选项
	// -----------------------------------------------
	fmt.Println("\n读取配置中的选项:")

	loadedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	// 显示所有配置选项
	if loadedConfig.Config != nil {
		fmt.Printf("找到 %d 个配置选项\n", len(loadedConfig.Config.Add))

		fmt.Println("配置选项列表:")
		for _, option := range loadedConfig.Config.Add {
			fmt.Printf("  - %s: %s\n", option.Key, option.Value)
		}
	} else {
		fmt.Println("未找到配置选项")
	}
	// 输出示例:
	// 读取配置中的选项:
	// 找到 8 个配置选项
	// 配置选项列表:
	//   - globalPackagesFolder: %USERPROFILE%\.nuget\packages
	//   - http_proxy: http://proxy.example.com:8080
	//   - no_proxy: localhost,127.0.0.1,.local
	//   - dependencyVersion: Highest
	//   - signatureValidationMode: require
	//   - allowedPackageFormats: nupkg
	//   - disableSourceControlIntegration: true
	//   - restorePackagesPath: packages
	//   - restoreIgnoreFailedSources: true

	// 6. 获取特定配置选项的值
	// -----------------------------------------------
	fmt.Println("\n获取特定配置选项值:")

	proxyValue := api.GetConfigOption(loadedConfig, "http_proxy")
	if proxyValue != "" {
		fmt.Printf("HTTP代理设置为: %s\n", proxyValue)
	}

	dependencyVersion := api.GetConfigOption(loadedConfig, "dependencyVersion")
	if dependencyVersion != "" {
		fmt.Printf("依赖版本策略: %s\n", dependencyVersion)
	}

	// 查询不存在的选项
	nonExistentOption := api.GetConfigOption(loadedConfig, "nonExistent")
	fmt.Printf("不存在的选项值: '%s'\n", nonExistentOption)
	// 输出示例:
	// 获取特定配置选项值:
	// HTTP代理设置为: http://proxy.example.com:8080
	// 依赖版本策略: Highest
	// 不存在的选项值: ''

	// 7. 修改现有选项
	// -----------------------------------------------
	fmt.Println("\n修改现有配置选项:")

	// 修改代理设置
	fmt.Println("- 修改HTTP代理设置...")
	api.AddConfigOption(loadedConfig, "http_proxy", "http://new-proxy.example.com:3128")

	// 修改依赖版本行为
	fmt.Println("- 修改依赖版本行为...")
	api.AddConfigOption(loadedConfig, "dependencyVersion", "HighestPatch")

	// 保存修改后的配置
	if err := api.SaveConfig(loadedConfig, configPath); err != nil {
		log.Fatalf("保存修改后的配置失败: %v", err)
	}

	// 8. 移除配置选项
	// -----------------------------------------------
	fmt.Println("\n移除配置选项:")

	// 移除不再需要的选项
	fmt.Println("- 移除 'disableSourceControlIntegration' 选项...")
	removed := api.RemoveConfigOption(loadedConfig, "disableSourceControlIntegration")
	fmt.Printf("选项已移除: %v\n", removed)

	// 移除不存在的选项 (返回false)
	nonExistentRemoved := api.RemoveConfigOption(loadedConfig, "nonExistent")
	fmt.Printf("不存在的选项移除结果: %v\n", nonExistentRemoved)

	// 保存修改后的配置
	if err := api.SaveConfig(loadedConfig, configPath); err != nil {
		log.Fatalf("保存修改后的配置失败: %v", err)
	}

	// 9. 验证更改后的配置
	// -----------------------------------------------
	fmt.Println("\n验证更改后的配置:")

	finalConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("加载更新后的配置失败: %v", err)
	}

	// 获取更新后的选项
	updatedProxyValue := api.GetConfigOption(finalConfig, "http_proxy")
	fmt.Printf("更新后的HTTP代理设置: %s\n", updatedProxyValue)

	updatedDependencyVersion := api.GetConfigOption(finalConfig, "dependencyVersion")
	fmt.Printf("更新后的依赖版本策略: %s\n", updatedDependencyVersion)

	// 检查被移除的选项
	removedOptionValue := api.GetConfigOption(finalConfig, "disableSourceControlIntegration")
	if removedOptionValue == "" {
		fmt.Println("已确认: 'disableSourceControlIntegration' 选项已被移除")
	} else {
		fmt.Printf("错误: 'disableSourceControlIntegration' 选项仍存在: %s\n", removedOptionValue)
	}

	// 输出最终配置选项数量
	if finalConfig.Config != nil {
		fmt.Printf("最终配置包含 %d 个选项\n", len(finalConfig.Config.Add))
	}
	// 输出示例:
	// 验证更改后的配置:
	// 更新后的HTTP代理设置: http://new-proxy.example.com:3128
	// 更新后的依赖版本策略: HighestPatch
	// 已确认: 'disableSourceControlIntegration' 选项已被移除
	// 最终配置包含 7 个选项

	// 10. 打印完整配置
	// -----------------------------------------------
	fmt.Println("\n最终配置文件内容:")
	finalXml, err := api.SerializeToXML(finalConfig)
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}
	fmt.Println(finalXml)
	// 输出中的配置选项部分将类似于:
	// <config>
	//   <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
	//   <add key="http_proxy" value="http://new-proxy.example.com:3128" />
	//   <add key="no_proxy" value="localhost,127.0.0.1,.local" />
	//   <add key="dependencyVersion" value="HighestPatch" />
	//   <add key="signatureValidationMode" value="require" />
	//   <add key="allowedPackageFormats" value="nupkg" />
	//   <add key="restorePackagesPath" value="packages" />
	//   <add key="restoreIgnoreFailedSources" value="true" />
	// </config>
}
