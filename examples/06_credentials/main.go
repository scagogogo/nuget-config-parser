// 该示例演示了如何使用NuGet配置解析器API管理包源凭证
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

func main() {
	// 1. 创建临时目录用于保存配置
	// -----------------------------------------------
	tempDir, err := os.MkdirTemp("", "nuget-credentials-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 创建API实例
	api := nuget.NewAPI()

	// 2. 创建一个包含多个包源的配置
	// -----------------------------------------------
	fmt.Println("创建包含多个包源的配置:")

	// 创建新的配置对象
	config := &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{
				// 公共源 - 无需凭证
				{
					Key:             "nuget.org",
					Value:           "https://api.nuget.org/v3/index.json",
					ProtocolVersion: "3",
				},
				// 公司内部源 - 需要凭证
				{
					Key:             "company-internal",
					Value:           "https://nuget.company.com/v3/index.json",
					ProtocolVersion: "3",
				},
				// 团队源 - 需要凭证
				{
					Key:             "team-repo",
					Value:           "https://nuget.team.com/v3/index.json",
					ProtocolVersion: "3",
				},
				// 合作伙伴源 - 需要凭证
				{
					Key:             "partner-repo",
					Value:           "https://partner.example.com/nuget/v3/index.json",
					ProtocolVersion: "3",
				},
			},
		},
	}

	// 3. 设置活跃包源
	// -----------------------------------------------
	fmt.Println("设置活跃包源为 company-internal...")
	if err := api.SetActivePackageSource(config, "company-internal"); err != nil {
		log.Fatalf("设置活跃包源失败: %v", err)
	}

	// 4. 添加各种凭证
	// -----------------------------------------------
	fmt.Println("\n添加各种类型的凭证:")

	// 4.1 添加基本用户名/密码凭证
	fmt.Println("- 添加公司内部源的基本凭证...")
	api.AddCredential(config, "company-internal", "company-user", "company-password")

	// 4.2 添加带有特殊字符的凭证
	fmt.Println("- 添加团队源的带特殊字符的凭证...")
	api.AddCredential(config, "team-repo", "team-user", "P@$$w0rd!*")

	// 4.3 添加API密钥作为凭证
	fmt.Println("- 添加合作伙伴源的API密钥凭证...")
	// 注意：这里仍然使用用户名/密码格式，但用户名为空，密码为API密钥
	api.AddCredential(config, "partner-repo", "", "api-key-12345abcdef")

	// 5. 保存配置
	// -----------------------------------------------
	fmt.Println("\n保存配置到文件...")
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存配置失败: %v", err)
	}

	fmt.Printf("配置已保存到: %s\n", configPath)

	// 6. 读取配置并展示凭证部分
	// -----------------------------------------------
	fmt.Println("\n读取配置检查凭证部分:")

	// 读取刚保存的配置
	loadedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	// 检查凭证是否存在
	if loadedConfig.PackageSourceCredentials == nil {
		fmt.Println("未找到凭证信息!")
	} else {
		fmt.Printf("找到 %d 个源的凭证\n",
			len(loadedConfig.PackageSourceCredentials.Sources))

		// 遍历所有凭证
		for sourceName, sourceCredentials := range loadedConfig.PackageSourceCredentials.Sources {
			fmt.Printf("\n源 '%s' 的凭证:\n", sourceName)

			// 打印每个凭证项（但隐藏实际值用于演示）
			for _, credential := range sourceCredentials.Add {
				var displayValue string
				if credential.Key == "Username" {
					// 用户名通常不需要隐藏
					displayValue = credential.Value
				} else {
					// 密码隐藏为星号
					displayValue = "********"
				}
				fmt.Printf("  - %s: %s\n", credential.Key, displayValue)
			}
		}
	}
	// 输出示例:
	// 读取配置检查凭证部分:
	// 找到 3 个源的凭证
	//
	// 源 'company-internal' 的凭证:
	//   - Username: company-user
	//   - ClearTextPassword: ********
	//
	// 源 'team-repo' 的凭证:
	//   - Username: team-user
	//   - ClearTextPassword: ********
	//
	// 源 'partner-repo' 的凭证:
	//   - Username:
	//   - ClearTextPassword: ********

	// 7. 修改现有凭证
	// -----------------------------------------------
	fmt.Println("\n修改现有凭证:")

	// 更新company-internal源的凭证
	fmt.Println("- 更新company-internal源的凭证...")
	api.AddCredential(config, "company-internal", "new-company-user", "new-company-password")

	// 重新保存配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存更新的配置失败: %v", err)
	}

	// 8. 移除凭证
	// -----------------------------------------------
	fmt.Println("\n移除凭证:")

	// 移除team-repo源的凭证
	fmt.Println("- 移除team-repo源的凭证...")
	removed := api.RemoveCredential(config, "team-repo")
	fmt.Printf("凭证已移除: %v\n", removed)

	// 保存配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存更新的配置失败: %v", err)
	}

	// 9. 验证修改后的凭证
	// -----------------------------------------------
	fmt.Println("\n验证修改后的凭证:")

	// 重新加载配置
	updatedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("读取更新后的配置失败: %v", err)
	}

	// 确认当前凭证状态
	if updatedConfig.PackageSourceCredentials != nil {
		fmt.Printf("当前有 %d 个源的凭证\n",
			len(updatedConfig.PackageSourceCredentials.Sources))

		// 列出所有源
		fmt.Println("源列表:")
		for sourceName := range updatedConfig.PackageSourceCredentials.Sources {
			fmt.Printf("  - %s\n", sourceName)
		}

		// 检查是否仍有team-repo的凭证
		if _, exists := updatedConfig.PackageSourceCredentials.Sources["team-repo"]; exists {
			fmt.Println("错误: team-repo源的凭证未被移除!")
		} else {
			fmt.Println("已确认: team-repo源的凭证已成功移除")
		}

		// 检查company-internal源的凭证是否被更新
		if creds, exists := updatedConfig.PackageSourceCredentials.Sources["company-internal"]; exists {
			for _, cred := range creds.Add {
				if cred.Key == "Username" {
					fmt.Printf("company-internal源的用户名: %s\n", cred.Value)
					if cred.Value == "new-company-user" {
						fmt.Println("已确认: company-internal源的凭证已成功更新")
					}
				}
			}
		}
	}
	// 输出示例:
	// 验证修改后的凭证:
	// 当前有 2 个源的凭证
	// 源列表:
	//   - company-internal
	//   - partner-repo
	// 已确认: team-repo源的凭证已成功移除
	// company-internal源的用户名: new-company-user
	// 已确认: company-internal源的凭证已成功更新

	// 10. 打印完整配置
	// -----------------------------------------------
	fmt.Println("\n最终配置文件内容:")
	finalXml, err := api.SerializeToXML(updatedConfig)
	if err != nil {
		log.Fatalf("序列化配置失败: %v", err)
	}
	fmt.Println(finalXml)
	// 输出示例中的凭证部分将类似于:
	// <packageSourceCredentials>
	//   <company-internal>
	//     <add key="Username" value="new-company-user" />
	//     <add key="ClearTextPassword" value="new-company-password" />
	//   </company-internal>
	//   <partner-repo>
	//     <add key="Username" value="" />
	//     <add key="ClearTextPassword" value="api-key-12345abcdef" />
	//   </partner-repo>
	// </packageSourceCredentials>
}
