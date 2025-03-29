// 该示例演示了如何使用NuGet配置解析器API查找配置文件
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	"github.com/scagogogo/nuget-config-parser/pkg/nuget"
)

func main() {
	// 1. 创建示例目录结构用于演示查找功能
	// -----------------------------------------------

	// 创建根临时目录
	rootDir, err := ioutil.TempDir("", "nuget-search-example-*")
	if err != nil {
		log.Fatalf("创建临时根目录失败: %v", err)
	}
	defer os.RemoveAll(rootDir) // 程序结束后删除临时目录

	// 创建项目目录结构
	projectDir := filepath.Join(rootDir, "MyProject")
	subprojectDir := filepath.Join(projectDir, "SubProject")

	// 创建目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		log.Fatalf("创建项目目录失败: %v", err)
	}
	if err := os.MkdirAll(subprojectDir, 0755); err != nil {
		log.Fatalf("创建子项目目录失败: %v", err)
	}

	// 准备配置内容
	configContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="company" value="https://nuget.mycompany.com/v3/index.json" />
  </packageSources>
</configuration>`

	// 创建根项目配置文件
	rootConfigPath := filepath.Join(rootDir, constants.DefaultNuGetConfigFilename)
	if err := ioutil.WriteFile(rootConfigPath, []byte(configContent), 0644); err != nil {
		log.Fatalf("写入根配置文件失败: %v", err)
	}

	// 创建项目级配置文件
	projectConfigPath := filepath.Join(projectDir, constants.DefaultNuGetConfigFilename)
	projectConfigContent := `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="project-source" value="https://nuget.project.com/v3/index.json" />
  </packageSources>
</configuration>`
	if err := ioutil.WriteFile(projectConfigPath, []byte(projectConfigContent), 0644); err != nil {
		log.Fatalf("写入项目配置文件失败: %v", err)
	}

	// 2. 使用API查找配置文件
	// -----------------------------------------------

	// 创建API实例
	api := nuget.NewAPI()

	// 3. 查找第一个可用的配置文件（从当前目录向上查找）
	// -----------------------------------------------
	fmt.Println("查找第一个可用的配置文件:")

	// 为了演示，我们改变当前工作目录到子项目目录
	originalDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}
	defer os.Chdir(originalDir) // 确保最后恢复原始目录

	if err := os.Chdir(subprojectDir); err != nil {
		log.Fatalf("切换到子项目目录失败: %v", err)
	}

	// 查找配置文件 (从子项目目录开始往上找)
	foundConfigPath, err := api.FindProjectConfig(".")
	if err != nil {
		fmt.Println("未找到项目配置文件")
	} else {
		fmt.Printf("找到项目配置文件: %s\n", foundConfigPath)
		// 输出示例: 找到项目配置文件: /tmp/nuget-search-example-123456/MyProject/NuGet.Config
	}

	// 4. 查找所有配置文件
	// -----------------------------------------------
	fmt.Println("\n查找所有可用的配置文件:")
	allConfigFiles := api.FindAllConfigFiles()
	for i, path := range allConfigFiles {
		fmt.Printf("配置文件 %d: %s\n", i+1, path)
	}
	// 输出示例:
	// 查找所有可用的配置文件:
	// 配置文件 1: /tmp/nuget-search-example-123456/MyProject/NuGet.Config
	// 配置文件 2: /tmp/nuget-search-example-123456/NuGet.Config
	// 配置文件 3: /Users/username/.config/NuGet/NuGet.Config

	// 5. 查找最近的配置文件并解析
	// -----------------------------------------------
	fmt.Println("\n查找并解析配置文件:")
	config, configPath, err := api.FindAndParseConfig()
	if err != nil {
		log.Fatalf("查找或解析配置失败: %v", err)
	}

	fmt.Printf("找到配置文件: %s\n", configPath)
	fmt.Printf("包含 %d 个包源\n", len(config.PackageSources.Add))

	// 6. 输出包源信息
	// -----------------------------------------------
	fmt.Println("\n解析到的包源:")
	for _, source := range config.PackageSources.Add {
		fmt.Printf("  - %s: %s\n", source.Key, source.Value)
	}
	// 输出示例:
	// 解析到的包源:
	//   - project-source: https://nuget.project.com/v3/index.json
}
