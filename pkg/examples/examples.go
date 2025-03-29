// Package examples 提供NuGet配置解析器的使用示例
package examples

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

// ParseConfigExample 解析配置文件示例
func ParseConfigExample(filePath string) {
	// 创建API实例
	api := nuget.NewAPI()

	// 从文件解析配置
	config, err := api.ParseFromFile(filePath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 输出包源信息
	fmt.Println("Package Sources:")
	for _, source := range config.PackageSources.Add {
		fmt.Printf("  - %s: %s\n", source.Key, source.Value)
		if source.ProtocolVersion != "" {
			fmt.Printf("    Protocol Version: %s\n", source.ProtocolVersion)
		}
	}

	// 输出禁用的包源
	if config.DisabledPackageSources != nil && len(config.DisabledPackageSources.Add) > 0 {
		fmt.Println("\nDisabled Package Sources:")
		for _, source := range config.DisabledPackageSources.Add {
			fmt.Printf("  - %s\n", source.Key)
		}
	}

	// 输出活跃包源
	if config.ActivePackageSource != nil {
		fmt.Printf("\nActive Package Source: %s (%s)\n",
			config.ActivePackageSource.Add.Key,
			config.ActivePackageSource.Add.Value)
	}

	// 输出配置选项
	if config.Config != nil && len(config.Config.Add) > 0 {
		fmt.Println("\nConfig Options:")
		for _, option := range config.Config.Add {
			fmt.Printf("  - %s: %s\n", option.Key, option.Value)
		}
	}
}

// CreateConfigExample 创建并保存配置文件示例
func CreateConfigExample(filePath string) {
	// 创建API实例
	api := nuget.NewAPI()

	// 创建新的配置
	config := &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{
				{
					Key:             "nuget.org",
					Value:           "https://api.nuget.org/v3/index.json",
					ProtocolVersion: "3",
				},
				{
					Key:   "local",
					Value: "C:\\Packages",
				},
			},
		},
	}

	// 添加活跃包源
	api.SetActivePackageSource(config, "nuget.org")

	// 添加凭证
	api.AddCredential(config, "nuget.org", "username", "password")

	// 添加配置选项
	api.AddConfigOption(config, "dependencyVersion", "Highest")

	// 保存配置
	err := api.SaveConfig(config, filePath)
	if err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("Config saved to %s\n", filePath)
}

// FindConfigExample 查找配置文件示例
func FindConfigExample() {
	// 创建API实例
	api := nuget.NewAPI()

	// 查找配置文件
	configPath, err := api.FindConfigFile()
	if err != nil {
		fmt.Println("No config file found")
		return
	}

	fmt.Printf("Found config file: %s\n", configPath)
}

// ModifyConfigExample 修改配置文件示例
func ModifyConfigExample(filePath string) {
	// 创建API实例
	api := nuget.NewAPI()

	// 解析配置
	config, err := api.ParseFromFile(filePath)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	// 添加新包源
	api.AddPackageSource(config, "MyCustomSource", "https://my-custom-source/index.json", "3")

	// 禁用包源
	api.DisablePackageSource(config, "local")

	// 添加配置选项
	api.AddConfigOption(config, "globalPackagesFolder", filepath.Join(os.Getenv("HOME"), ".nuget", "packages"))

	// 保存修改后的配置
	tempFile := filePath + ".modified"
	err = api.SaveConfig(config, tempFile)
	if err != nil {
		log.Fatalf("Failed to save modified config: %v", err)
	}

	fmt.Printf("Modified config saved to %s\n", tempFile)
}

// ProcessAllConfigsExample 处理所有找到的配置文件
func ProcessAllConfigsExample() {
	// 创建API实例
	api := nuget.NewAPI()

	// 找到所有配置文件
	configFiles := api.FindAllConfigFiles()

	if len(configFiles) == 0 {
		fmt.Println("No config files found")
		return
	}

	fmt.Printf("Found %d config files:\n", len(configFiles))

	// 处理每个配置文件
	for i, filePath := range configFiles {
		fmt.Printf("\n=== Config %d: %s ===\n", i+1, filePath)

		config, err := api.ParseFromFile(filePath)
		if err != nil {
			fmt.Printf("Error parsing config: %v\n", err)
			continue
		}

		// 输出包源数量
		fmt.Printf("Package sources: %d\n", len(config.PackageSources.Add))

		// 输出包含凭证的包源
		if config.PackageSourceCredentials != nil {
			fmt.Printf("Sources with credentials: %d\n", len(config.PackageSourceCredentials.Sources))
		}
	}
}
