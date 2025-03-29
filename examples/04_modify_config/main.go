// 该示例演示了如何使用NuGet配置解析器API修改现有的配置文件
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
	// 1. 创建临时目录并生成示例配置文件
	// -----------------------------------------------
	tempDir, err := os.MkdirTemp("", "nuget-modify-example-*")
	if err != nil {
		log.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir) // 程序结束后删除临时目录

	// 创建示例配置文件
	configPath := filepath.Join(tempDir, "NuGet.Config")

	// 准备一个初始NuGet配置内容
	initialConfigContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="testSource" value="https://test-source.example.com" />
  </packageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
  </config>
</configuration>`

	// 写入初始配置
	if err := os.WriteFile(configPath, []byte(initialConfigContent), 0644); err != nil {
		log.Fatalf("写入配置文件失败: %v", err)
	}

	fmt.Printf("已创建初始配置文件: %s\n", configPath)
	// 输出示例: 已创建初始配置文件: /tmp/nuget-modify-example-123456/NuGet.Config

	// 2. 使用API解析现有配置
	// -----------------------------------------------
	fmt.Println("\n正在解析现有配置文件...")

	// 创建API实例
	api := nuget.NewAPI()

	// 解析现有配置文件
	config, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	fmt.Printf("成功解析配置文件，包含 %d 个包源\n", len(config.PackageSources.Add))

	// 3. 修改配置
	// -----------------------------------------------
	fmt.Println("\n正在修改配置...")

	// 3.1 添加新的包源
	fmt.Println("添加新的包源: company-source")
	api.AddPackageSource(config, "company-source", "https://nuget.company.com/v3/index.json", "3")

	// 3.2 移除现有包源
	fmt.Println("移除现有包源: testSource")
	if removed := api.RemovePackageSource(config, "testSource"); removed {
		fmt.Println("成功移除包源 testSource")
	} else {
		fmt.Println("未找到包源 testSource")
	}

	// 3.3 修改活跃包源
	fmt.Println("修改活跃包源为: company-source")
	if err := api.SetActivePackageSource(config, "company-source"); err != nil {
		log.Fatalf("设置活跃包源失败: %v", err)
	}

	// 3.4 添加凭证
	fmt.Println("为 company-source 添加凭证")
	api.AddCredential(config, "company-source", "companyUser", "companyPassword")

	// 3.5 修改配置选项
	fmt.Println("更新全局包文件夹路径")
	api.AddConfigOption(config, "globalPackagesFolder", "/opt/nuget/packages")

	// 3.6 添加新的配置选项
	fmt.Println("添加新配置选项: dependencyVersion")
	api.AddConfigOption(config, "dependencyVersion", "Highest")

	// 4. 保存修改后的配置
	// -----------------------------------------------
	fmt.Println("\n保存修改后的配置...")

	// 创建备份
	backupPath := configPath + ".bak"
	if err := copyFile(configPath, backupPath); err != nil {
		log.Fatalf("创建配置备份失败: %v", err)
	}
	fmt.Printf("已创建配置备份: %s\n", backupPath)

	// 保存修改后的配置
	if err := api.SaveConfig(config, configPath); err != nil {
		log.Fatalf("保存修改后的配置失败: %v", err)
	}
	fmt.Printf("修改后的配置已保存到: %s\n", configPath)

	// 5. 比较修改前后的配置
	// -----------------------------------------------
	fmt.Println("\n修改前的配置内容:")
	originalContent, err := os.ReadFile(backupPath)
	if err != nil {
		log.Fatalf("读取原始配置文件失败: %v", err)
	}
	fmt.Println(string(originalContent))

	fmt.Println("\n修改后的配置内容:")
	modifiedContent, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取修改后的配置文件失败: %v", err)
	}
	fmt.Println(string(modifiedContent))
	// 输出示例:
	// 修改后的配置内容:
	// <?xml version="1.0" encoding="utf-8"?>
	// <configuration>
	//   <packageSources>
	//     <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
	//     <add key="company-source" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
	//   </packageSources>
	//   <packageSourceCredentials>
	//     <company-source>
	//       <add key="Username" value="companyUser" />
	//       <add key="ClearTextPassword" value="companyPassword" />
	//     </company-source>
	//   </packageSourceCredentials>
	//   <activePackageSource>
	//     <add key="company-source" value="https://nuget.company.com/v3/index.json" protocolVersion="3" />
	//   </activePackageSource>
	//   <config>
	//     <add key="globalPackagesFolder" value="/opt/nuget/packages" />
	//     <add key="dependencyVersion" value="Highest" />
	//   </config>
	// </configuration>

	// 6. 重新解析修改后的配置验证更改
	// -----------------------------------------------
	fmt.Println("\n验证修改后的配置:")

	modifiedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		log.Fatalf("解析修改后的配置文件失败: %v", err)
	}

	// 验证包源数量
	fmt.Printf("包源数量: %d\n", len(modifiedConfig.PackageSources.Add))

	// 验证活跃包源
	if modifiedConfig.ActivePackageSource != nil {
		fmt.Printf("活跃包源: %s\n", modifiedConfig.ActivePackageSource.Add.Key)
	}

	// 验证凭证
	if modifiedConfig.PackageSourceCredentials != nil {
		fmt.Println("包源凭证:")
		for source := range modifiedConfig.PackageSourceCredentials.Sources {
			fmt.Printf("  - %s\n", source)
		}
	}

	// 验证配置选项
	if modifiedConfig.Config != nil {
		fmt.Println("配置选项:")
		for _, option := range modifiedConfig.Config.Add {
			fmt.Printf("  - %s: %s\n", option.Key, option.Value)
		}
	}
	// 输出示例:
	// 验证修改后的配置:
	// 包源数量: 2
	// 活跃包源: company-source
	// 包源凭证:
	//   - company-source
	// 配置选项:
	//   - globalPackagesFolder: /opt/nuget/packages
	//   - dependencyVersion: Highest
}

// 辅助函数: 复制文件
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
