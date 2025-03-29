// 该示例演示了如何使用NuGet配置解析器API序列化和反序列化配置
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
	// 1. 创建临时目录用于演示
	// -----------------------------------------------
	tempDir, err := ioutil.TempDir("", "nuget-serialization-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	// 创建API实例
	api := nuget.NewAPI()

	// 2. 创建一个示例配置
	// -----------------------------------------------
	fmt.Println("创建示例配置:")

	// 创建完整的配置对象，包含各种元素
	config := &types.NuGetConfig{
		// 包源设置
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{
				{
					Key:             "nuget.org",
					Value:           "https://api.nuget.org/v3/index.json",
					ProtocolVersion: "3",
				},
				{
					Key:   "local",
					Value: "C:\\LocalPackages",
				},
			},
		},

		// 活跃包源
		ActivePackageSource: &types.ActivePackageSource{
			Add: types.PackageSource{
				Key:   "nuget.org",
				Value: "https://api.nuget.org/v3/index.json",
			},
		},

		// 禁用的包源
		DisabledPackageSources: &types.DisabledPackageSources{
			Add: []types.DisabledSource{
				{
					Key:   "local",
					Value: "true",
				},
			},
		},

		// 凭证
		PackageSourceCredentials: &types.PackageSourceCredentials{
			Sources: map[string]types.SourceCredential{
				"nuget.org": {
					Add: []types.Credential{
						{
							Key:   "Username",
							Value: "user@example.com",
						},
						{
							Key:   "ClearTextPassword",
							Value: "P@ssw0rd",
						},
					},
				},
			},
		},

		// 配置选项
		Config: &types.Config{
			Add: []types.ConfigOption{
				{
					Key:   "globalPackagesFolder",
					Value: "%USERPROFILE%\\.nuget\\packages",
				},
				{
					Key:   "dependencyVersion",
					Value: "Highest",
				},
			},
		},
	}

	// 3. 将配置序列化为XML字符串
	// -----------------------------------------------
	fmt.Println("\n序列化配置为XML:")

	xmlString, err := api.SerializeToXML(config)
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}

	fmt.Println("序列化后的XML内容:")
	fmt.Println(xmlString)
	// 输出示例:
	// 序列化后的XML内容:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
	//     <add key="local" value="C:\LocalPackages" />
	//   </packageSources>
	//   <packageSourceCredentials>
	//     <nuget.org>
	//       <add key="Username" value="user@example.com" />
	//       <add key="ClearTextPassword" value="P@ssw0rd" />
	//     </nuget.org>
	//   </packageSourceCredentials>
	//   <disabledPackageSources>
	//     <add key="local" value="true" />
	//   </disabledPackageSources>
	//   <activePackageSource>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
	//   </activePackageSource>
	//   <config>
	//     <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
	//     <add key="dependencyVersion" value="Highest" />
	//   </config>
	// </configuration>

	// 4. 将序列化后的XML保存到文件
	// -----------------------------------------------
	fmt.Println("\n保存XML到文件:")

	// 创建配置文件路径
	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 保存配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存配置失败: %v", err)
	}

	fmt.Printf("配置已保存到: %s\n", configPath)

	// 5. 从XML字符串反序列化为配置对象
	// -----------------------------------------------
	fmt.Println("\n从XML字符串反序列化为配置对象:")

	// 使用ParseFromString方法反序列化
	parsedConfig, err := api.ParseFromString(xmlString)
	if err != nil {
		log.Fatalf("从XML字符串解析配置失败: %v", err)
	}

	// 验证反序列化的结果
	fmt.Printf("成功解析配置，包含 %d 个包源\n", len(parsedConfig.PackageSources.Add))

	if parsedConfig.ActivePackageSource != nil {
		fmt.Printf("活跃包源: %s\n", parsedConfig.ActivePackageSource.Add.Key)
	}

	if parsedConfig.DisabledPackageSources != nil {
		fmt.Printf("禁用的包源数量: %d\n", len(parsedConfig.DisabledPackageSources.Add))
	}

	if parsedConfig.Config != nil {
		fmt.Printf("配置选项数量: %d\n", len(parsedConfig.Config.Add))
	}

	// 输出示例:
	// 从XML字符串反序列化为配置对象:
	// 成功解析配置，包含 2 个包源
	// 活跃包源: nuget.org
	// 禁用的包源数量: 1
	// 配置选项数量: 2

	// 6. 从文件反序列化配置
	// -----------------------------------------------
	fmt.Println("\n从文件反序列化配置:")

	// 从保存的文件中加载配置
	loadedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("从文件加载配置失败: %v", err)
	}

	// 验证从文件加载的配置
	fmt.Printf("从文件加载的配置包含 %d 个包源\n", len(loadedConfig.PackageSources.Add))

	// 检查凭证是否正确加载
	if loadedConfig.PackageSourceCredentials != nil {
		fmt.Println("\n凭证信息:")
		for sourceName, creds := range loadedConfig.PackageSourceCredentials.Sources {
			fmt.Printf("源 '%s' 的凭证:\n", sourceName)
			for _, cred := range creds.Add {
				if cred.Key == "Username" {
					fmt.Printf("  用户名: %s\n", cred.Value)
				} else if cred.Key == "ClearTextPassword" {
					fmt.Printf("  密码: %s (实际应用中应该保护密码)\n", cred.Value)
				}
			}
		}
	}

	// 7. 修改配置并重新序列化
	// -----------------------------------------------
	fmt.Println("\n修改配置并重新序列化:")

	// 添加一个新的包源
	api.AddPackageSource(loadedConfig, "myCompany", "https://nuget.mycompany.com/v3/index.json", "3")

	// 添加新包源的凭证
	api.AddCredential(loadedConfig, "myCompany", "companyUser", "companyPassword")

	// 禁用一个包源
	api.DisablePackageSource(loadedConfig, "myCompany")

	// 重新序列化
	modifiedXml, err := api.SerializeToXML(loadedConfig)
	if err != nil {
		log.Fatalf("序列化修改后的配置失败: %v", err)
	}

	// 保存修改后的配置到新文件
	modifiedPath := filepath.Join(tempDir, "Modified.NuGet.Config")
	if err := ioutil.WriteFile(modifiedPath, []byte(modifiedXml), 0644); err != nil {
		log.Fatalf("保存修改后的配置失败: %v", err)
	}

	fmt.Printf("修改后的配置已保存到: %s\n", modifiedPath)

	// 8. 从流(Reader)解析配置
	// -----------------------------------------------
	fmt.Println("\n从流(Reader)解析配置:")

	// 创建一个包含配置的流
	configFile, err := os.Open(modifiedPath)
	if err != nil {
		log.Fatalf("打开配置文件失败: %v", err)
	}
	defer configFile.Close()

	// 从流解析配置
	readerConfig, err := api.ParseFromReader(configFile)
	if err != nil {
		log.Fatalf("从流解析配置失败: %v", err)
	}

	// 验证从流解析的配置
	fmt.Printf("从流解析的配置包含 %d 个包源\n", len(readerConfig.PackageSources.Add))
	fmt.Printf("从流解析的配置包含 %d 个禁用包源\n", len(readerConfig.DisabledPackageSources.Add))

	// 9. 演示空配置和特殊情况的序列化
	// -----------------------------------------------
	fmt.Println("\n特殊情况的序列化:")

	// 创建一个空配置
	emptyConfig := &types.NuGetConfig{}

	// 序列化空配置
	emptyXml, err := api.SerializeToXML(emptyConfig)
	if err != nil {
		log.Fatalf("序列化空配置失败: %v", err)
	}

	fmt.Println("空配置序列化结果:")
	fmt.Println(emptyXml)
	// 输出示例:
	// 空配置序列化结果:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources>
	//   </packageSources>
	// </configuration>

	// 10. 创建一个带有Clear标志的配置
	// -----------------------------------------------
	fmt.Println("\n带有Clear标志的配置:")

	clearConfig := &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Clear: true,
			Add: []types.PackageSource{
				{
					Key:   "onlySource",
					Value: "https://only.source.com/nuget",
				},
			},
		},
	}

	// 序列化带Clear标志的配置
	clearXml, err := api.SerializeToXML(clearConfig)
	if err != nil {
		log.Fatalf("序列化Clear配置失败: %v", err)
	}

	fmt.Println("带Clear标志的配置序列化结果:")
	fmt.Println(clearXml)
	// 输出示例:
	// 带Clear标志的配置序列化结果:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources clear="true">
	//     <add key="onlySource" value="https://only.source.com/nuget" />
	//   </packageSources>
	// </configuration>
}
