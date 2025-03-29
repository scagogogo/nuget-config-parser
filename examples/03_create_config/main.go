// 该示例演示了如何使用NuGet配置解析器API创建新的配置文件
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
	// 1. 创建临时目录用于保存配置
	// -----------------------------------------------
	tempDir, err := ioutil.TempDir("", "nuget-create-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	// 2. 从零创建配置
	// -----------------------------------------------
	fmt.Println("从零创建新的配置文件:")

	// 创建API实例
	api := nuget.NewAPI()

	// 创建一个新的空配置
	config := &types.NuGetConfig{
		// 初始化包源部分
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{
				// 添加官方 nuget.org 包源
				{
					Key:             "nuget.org",
					Value:           "https://api.nuget.org/v3/index.json",
					ProtocolVersion: "3",
				},
				// 添加公司内部包源
				{
					Key:   "company-internal",
					Value: "https://nuget.company.com/v3/index.json",
				},
				// 添加本地包源
				{
					Key:   "local",
					Value: "C:\\LocalPackages",
				},
			},
		},
	}

	// 3. 设置活跃包源
	// -----------------------------------------------
	// 将 nuget.org 设置为活跃包源
	if err := api.SetActivePackageSource(config, "nuget.org"); err != nil {
		log.Fatalf("设置活跃包源失败: %v", err)
	}

	// 4. 添加凭证（用于需要身份验证的包源）
	// -----------------------------------------------
	api.AddCredential(config, "company-internal", "username", "password")

	// 5. 禁用包源
	// -----------------------------------------------
	// 禁用本地包源
	api.DisablePackageSource(config, "local")

	// 6. 添加配置选项
	// -----------------------------------------------
	// 设置全局包文件夹
	api.AddConfigOption(config, "globalPackagesFolder", "%USERPROFILE%\\.nuget\\packages")
	// 设置允许的包格式版本
	api.AddConfigOption(config, "allowedPackageFormats", "nupkg")

	// 7. 保存配置到文件
	// -----------------------------------------------
	configPath := filepath.Join(tempDir, "NuGet.Config")
	err = api.SaveConfig(config, configPath)
	if err != nil {
		log.Fatalf("保存配置文件失败: %v", err)
	}

	fmt.Printf("配置已保存到: %s\n", configPath)
	// 输出示例: 配置已保存到: /tmp/nuget-create-example-123456/NuGet.Config

	// 8. 读取创建的文件并输出序列化后的 XML
	// -----------------------------------------------
	fmt.Println("\n创建的配置文件内容:")

	// 读取文件内容
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 输出文件内容
	fmt.Println(string(content))
	// 输出示例:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
	//     <add key="company-internal" value="https://nuget.company.com/v3/index.json" />
	//     <add key="local" value="C:\LocalPackages" />
	//   </packageSources>
	//   <packageSourceCredentials>
	//     <company-internal>
	//       <add key="Username" value="username" />
	//       <add key="ClearTextPassword" value="password" />
	//     </company-internal>
	//   </packageSourceCredentials>
	//   <disabledPackageSources>
	//     <add key="local" value="true" />
	//   </disabledPackageSources>
	//   <activePackageSource>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
	//   </activePackageSource>
	//   <config>
	//     <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
	//     <add key="allowedPackageFormats" value="nupkg" />
	//   </config>
	// </configuration>

	// 9. 创建并保存默认配置
	// -----------------------------------------------
	fmt.Println("\n创建默认配置文件:")

	// 创建默认配置
	defaultConfig := api.CreateDefaultConfig()
	defaultConfigPath := filepath.Join(tempDir, "Default.NuGet.Config")

	// 保存默认配置
	err = api.SaveConfig(defaultConfig, defaultConfigPath)
	if err != nil {
		log.Fatalf("保存默认配置文件失败: %v", err)
	}

	fmt.Printf("默认配置已保存到: %s\n", defaultConfigPath)
	// 输出示例: 默认配置已保存到: /tmp/nuget-create-example-123456/Default.NuGet.Config

	// 读取并输出默认配置内容
	defaultContent, err := ioutil.ReadFile(defaultConfigPath)
	if err != nil {
		log.Fatalf("读取默认配置文件失败: %v", err)
	}

	fmt.Println("默认配置文件内容:")
	fmt.Println(string(defaultContent))
	// 输出示例:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
	//   </packageSources>
	//   <activePackageSource>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
	//   </activePackageSource>
	// </configuration>
}
