// Package nuget 提供NuGet配置文件解析的主要API
package nuget

import (
	"io"

	"github.com/scagogogo/nuget-config-parser/pkg/finder"
	"github.com/scagogogo/nuget-config-parser/pkg/manager"
	"github.com/scagogogo/nuget-config-parser/pkg/parser"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

// API 提供NuGet配置文件解析的所有功能
type API struct {
	Parser  *parser.ConfigParser
	Finder  *finder.ConfigFinder
	Manager *manager.ConfigManager
}

// NewAPI 创建新的API实例
//
// NewAPI 初始化并返回一个新的 API 实例，该实例集成了解析器、查找器和管理器的功能。
// 这是使用库的入口点，通过该函数获取的 API 实例可以执行所有 NuGet 配置相关操作。
//
// 返回值:
//   - *API: 初始化好的 API 实例
//
// 示例:
//
//	// 创建 API 实例
//	api := nuget.NewAPI()
//
//	// 然后可以使用 api 调用其他方法
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
func NewAPI() *API {
	return &API{
		Parser:  parser.NewConfigParser(),
		Finder:  finder.NewConfigFinder(),
		Manager: manager.NewConfigManager(),
	}
}

// ParseFromFile 从文件解析NuGet配置
//
// ParseFromFile 读取指定路径的文件内容，并将其解析为 NuGet 配置对象。
// 如果文件不存在或内容无法解析为有效的 NuGet 配置，将返回相应的错误。
//
// 参数:
//   - filePath: 配置文件的路径，可以是绝对路径或相对路径
//
// 返回值:
//   - *types.NuGetConfig: 解析后的配置对象，如果解析失败则为 nil
//   - error: 如果解析过程中发生错误，则返回相应的错误；如果成功则为 nil
//
// 错误:
//   - errors.ErrConfigFileNotFound: 当指定的文件不存在时
//   - errors.ErrEmptyConfigFile: 当文件存在但内容为空时
//   - errors.ErrInvalidConfigFormat: 当文件内容不是有效的 XML 时
//   - errors.ErrXMLParsing: 当 XML 解析过程中出现错误时
//   - errors.ErrMissingRequiredElement: 当配置缺少必需的元素时
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 解析配置文件
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    if errors.IsNotFoundError(err) {
//	        fmt.Println("配置文件不存在")
//	    } else if errors.IsFormatError(err) {
//	        fmt.Println("配置文件格式无效")
//	    } else {
//	        fmt.Printf("解析失败: %v\n", err)
//	    }
//	    return
//	}
//
//	// 使用解析后的配置
//	for _, source := range config.PackageSources.Add {
//	    fmt.Printf("包源: %s - %s\n", source.Key, source.Value)
//	}
func (a *API) ParseFromFile(filePath string) (*types.NuGetConfig, error) {
	return a.Parser.ParseFromFile(filePath)
}

// ParseFromString 从字符串解析NuGet配置
//
// ParseFromString 将提供的字符串内容解析为 NuGet 配置对象。
// 如果字符串内容无法解析为有效的 NuGet 配置，将返回相应的错误。
//
// 参数:
//   - content: 包含 NuGet 配置 XML 的字符串
//
// 返回值:
//   - *types.NuGetConfig: 解析后的配置对象，如果解析失败则为 nil
//   - error: 如果解析过程中发生错误，则返回相应的错误；如果成功则为 nil
//
// 错误:
//   - errors.ErrEmptyConfigFile: 当提供的字符串为空时
//   - errors.ErrInvalidConfigFormat: 当字符串内容不是有效的 XML 时
//   - errors.ErrXMLParsing: 当 XML 解析过程中出现错误时
//   - errors.ErrMissingRequiredElement: 当配置缺少必需的元素时
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// XML 配置字符串
//	configXML := `<?xml version="1.0" encoding="utf-8"?>
//	<configuration>
//	  <packageSources>
//	    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
//	  </packageSources>
//	</configuration>`
//
//	// 解析配置字符串
//	config, err := api.ParseFromString(configXML)
//	if err != nil {
//	    fmt.Printf("解析失败: %v\n", err)
//	    return
//	}
//
//	// 输出包源数量
//	fmt.Printf("包含 %d 个包源\n", len(config.PackageSources.Add))
func (a *API) ParseFromString(content string) (*types.NuGetConfig, error) {
	return a.Parser.ParseFromString(content)
}

// ParseFromReader 从io.Reader解析NuGet配置
//
// ParseFromReader 从提供的 io.Reader 读取内容并将其解析为 NuGet 配置对象。
// 这个方法适用于从非文件源（如网络响应、内存缓冲区等）读取配置。
//
// 参数:
//   - reader: 实现了 io.Reader 接口的对象，提供配置内容
//
// 返回值:
//   - *types.NuGetConfig: 解析后的配置对象，如果解析失败则为 nil
//   - error: 如果解析过程中发生错误，则返回相应的错误；如果成功则为 nil
//
// 错误:
//   - 可能返回 io.ReadAll 相关的错误
//   - errors.ErrEmptyConfigFile: 当读取的内容为空时
//   - errors.ErrInvalidConfigFormat: 当内容不是有效的 XML 时
//   - errors.ErrXMLParsing: 当 XML 解析过程中出现错误时
//   - errors.ErrMissingRequiredElement: 当配置缺少必需的元素时
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 从字符串创建 Reader
//	configXML := `<?xml version="1.0" encoding="utf-8"?>
//	<configuration>
//	  <packageSources>
//	    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
//	  </packageSources>
//	</configuration>`
//	reader := strings.NewReader(configXML)
//
//	// 解析配置
//	config, err := api.ParseFromReader(reader)
//	if err != nil {
//	    fmt.Printf("解析失败: %v\n", err)
//	    return
//	}
//
//	// 检查是否有活跃包源
//	if config.ActivePackageSource != nil {
//	    fmt.Printf("活跃包源: %s\n", config.ActivePackageSource.Add.Key)
//	}
func (a *API) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error) {
	return a.Parser.ParseFromReader(reader)
}

// FindConfigFile 查找配置文件
//
// FindConfigFile 在系统中查找第一个可用的 NuGet 配置文件，按照预定义的搜索顺序。
// 搜索顺序通常为：当前目录、上级目录、用户目录、系统目录。
// 如果设置了 NUGET_CONFIG_FILE 环境变量，会优先检查该变量指向的文件。
//
// 返回值:
//   - string: 找到的配置文件的绝对路径
//   - error: 如果未找到任何配置文件，则返回 os.ErrNotExist；否则返回 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 查找配置文件
//	configPath, err := api.FindConfigFile()
//	if err != nil {
//	    fmt.Println("未找到任何 NuGet 配置文件")
//	    return
//	}
//
//	fmt.Printf("找到配置文件: %s\n", configPath)
//
//	// 解析找到的配置文件
//	config, err := api.ParseFromFile(configPath)
//	if err != nil {
//	    fmt.Printf("解析配置文件失败: %v\n", err)
//	    return
//	}
func (a *API) FindConfigFile() (string, error) {
	return a.Finder.FindConfigFile()
}

// FindAllConfigFiles 查找所有配置文件
//
// FindAllConfigFiles 在系统中查找所有可用的 NuGet 配置文件，按照预定义的搜索顺序。
// 搜索顺序通常为：当前目录、上级目录、用户目录、系统目录。
// 如果设置了 NUGET_CONFIG_FILE 环境变量，会包含该变量指向的文件。
//
// 返回值:
//   - []string: 找到的所有配置文件的绝对路径列表。如果没有找到任何文件，则返回空切片。
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 查找所有配置文件
//	configPaths := api.FindAllConfigFiles()
//
//	if len(configPaths) == 0 {
//	    fmt.Println("未找到任何 NuGet 配置文件")
//	    return
//	}
//
//	fmt.Printf("找到 %d 个配置文件:\n", len(configPaths))
//	for i, path := range configPaths {
//	    fmt.Printf("%d. %s\n", i+1, path)
//	}
//
//	// 可以进一步解析每个配置文件
//	for _, path := range configPaths {
//	    config, err := api.ParseFromFile(path)
//	    if err != nil {
//	        fmt.Printf("解析 %s 失败: %v\n", path, err)
//	        continue
//	    }
//	    fmt.Printf("%s 包含 %d 个包源\n", path, len(config.PackageSources.Add))
//	}
func (a *API) FindAllConfigFiles() []string {
	return a.Finder.FindAllConfigFiles()
}

// FindProjectConfig 在项目目录中查找配置文件
//
// FindProjectConfig 从指定目录开始向上查找项目级 NuGet 配置文件。
// 此方法适用于需要获取项目特定配置文件的场景，会从指定目录开始，逐级向上搜索
// 直到找到配置文件或者到达文件系统根目录。
//
// 参数:
//   - startDir: 搜索的起始目录路径
//
// 返回值:
//   - string: 找到的配置文件的绝对路径
//   - error: 如果未找到任何配置文件，则返回 os.ErrNotExist；否则返回 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 从当前目录开始查找项目配置
//	projectConfigPath, err := api.FindProjectConfig(".")
//	if err != nil {
//	    fmt.Println("未找到项目级 NuGet 配置文件")
//	    return
//	}
//
//	fmt.Printf("找到项目配置文件: %s\n", projectConfigPath)
//
//	// 解析项目配置
//	config, err := api.ParseFromFile(projectConfigPath)
//	if err != nil {
//	    fmt.Printf("解析项目配置失败: %v\n", err)
//	    return
//	}
//
//	// 检查项目特定的包源
//	fmt.Println("项目包源:")
//	for _, source := range config.PackageSources.Add {
//	    fmt.Printf("  - %s: %s\n", source.Key, source.Value)
//	}
func (a *API) FindProjectConfig(startDir string) (string, error) {
	return a.Finder.FindProjectConfig(startDir)
}

// FindAndParseConfig 查找并解析配置文件
//
// FindAndParseConfig 自动查找系统中第一个可用的 NuGet 配置文件并解析它。
// 这个方法结合了查找和解析步骤，适用于只需快速获取配置内容的场景。
//
// 返回值:
//   - *types.NuGetConfig: 解析后的配置对象，如果解析失败则为 nil
//   - string: 找到并成功解析的配置文件路径
//   - error: 如果查找或解析过程中发生错误，则返回相应的错误；如果成功则为 nil
//
// 错误:
//   - errors.ErrConfigFileNotFound: 当未找到任何配置文件时
//   - 以及 ParseFromFile 可能返回的任何错误
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 查找并解析配置
//	config, configPath, err := api.FindAndParseConfig()
//	if err != nil {
//	    if errors.IsNotFoundError(err) {
//	        fmt.Println("未找到任何 NuGet 配置文件")
//	    } else {
//	        fmt.Printf("解析配置失败: %v\n", err)
//	    }
//	    return
//	}
//
//	fmt.Printf("使用配置文件: %s\n", configPath)
//
//	// 使用配置对象
//	fmt.Printf("包含 %d 个包源\n", len(config.PackageSources.Add))
//
//	if config.ActivePackageSource != nil {
//	    fmt.Printf("活跃包源: %s\n", config.ActivePackageSource.Add.Key)
//	}
func (a *API) FindAndParseConfig() (*types.NuGetConfig, string, error) {
	return a.Manager.FindAndLoadConfig()
}

// SaveConfig 保存配置到文件
//
// SaveConfig 将 NuGet 配置对象序列化为 XML 并保存到指定路径的文件中。
// 如果文件路径的父目录不存在，会自动创建。如果文件已存在，将被覆盖。
//
// 参数:
//   - config: 要保存的 NuGet 配置对象
//   - filePath: 保存的目标文件路径
//
// 返回值:
//   - error: 如果保存过程中发生错误则返回相应的错误；如果成功则为 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 创建或加载配置
//	config := api.CreateDefaultConfig()
//	// 或从现有文件加载
//	// config, err := api.ParseFromFile("/path/to/existing/NuGet.Config")
//
//	// 修改配置
//	api.AddPackageSource(config, "custom-source", "https://custom-nuget-source/v3/index.json", "3")
//
//	// 保存到文件
//	err := api.SaveConfig(config, "/path/to/save/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("配置已成功保存")
func (a *API) SaveConfig(config *types.NuGetConfig, filePath string) error {
	return a.Manager.SaveConfig(config, filePath)
}

// CreateDefaultConfig 创建默认配置
//
// CreateDefaultConfig 创建并返回一个包含默认设置的 NuGet 配置对象。
// 默认配置通常包含官方 NuGet 包源 (nuget.org)，并将其设置为活跃包源。
// 返回的配置对象可以进一步修改后保存。
//
// 返回值:
//   - *types.NuGetConfig: 包含默认设置的新配置对象
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 创建默认配置
//	config := api.CreateDefaultConfig()
//
//	// 检查默认配置的包源
//	fmt.Println("默认配置包源:")
//	for _, source := range config.PackageSources.Add {
//	    fmt.Printf("  - %s: %s\n", source.Key, source.Value)
//	}
//
//	// 添加自定义设置
//	api.AddPackageSource(config, "company-source", "https://company-nuget/v3/index.json", "3")
//	api.AddConfigOption(config, "globalPackagesFolder", "/custom/packages/path")
//
//	// 保存修改后的配置
//	err := api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	}
func (a *API) CreateDefaultConfig() *types.NuGetConfig {
	return a.Manager.CreateDefaultConfig()
}

// InitializeDefaultConfig 在指定路径创建默认配置
//
// InitializeDefaultConfig 创建一个包含默认设置的 NuGet 配置文件，并保存到指定路径。
// 此方法是创建默认配置并保存的快捷方式，等同于调用 CreateDefaultConfig 后再调用 SaveConfig。
// 如果目标文件已存在，将被覆盖。
//
// 参数:
//   - filePath: 保存默认配置的目标文件路径
//
// 返回值:
//   - error: 如果创建或保存过程中发生错误则返回相应的错误；如果成功则为 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 初始化默认配置文件
//	err := api.InitializeDefaultConfig("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("初始化默认配置失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("已创建默认配置文件")
//
//	// 现在可以加载并修改这个默认配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 添加公司内部包源
//	api.AddPackageSource(config, "internal-source", "https://internal-nuget/v3/index.json", "3")
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	}
func (a *API) InitializeDefaultConfig(filePath string) error {
	return a.Manager.InitializeDefaultConfig(filePath)
}

// Package Source 操作

// AddPackageSource 添加包源
//
// AddPackageSource 向配置中添加一个新的 NuGet 包源或更新现有包源。
// 如果指定键名的包源已存在，将更新其 URL 和协议版本。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 包源的唯一标识符/名称
//   - value: 包源的 URL 或本地路径
//   - protocolVersion: 包源使用的协议版本（如 "2" 或 "3"），可为空
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载或创建配置
//	config := api.CreateDefaultConfig()
//
//	// 添加公共 NuGet 源
//	api.AddPackageSource(config, "nuget.org", "https://api.nuget.org/v3/index.json", "3")
//
//	// 添加公司内部源
//	api.AddPackageSource(config, "company-internal", "https://nuget.company.com/v3/index.json", "3")
//
//	// 添加本地文件夹源，不指定协议版本
//	api.AddPackageSource(config, "local", "C:\\LocalPackages", "")
//
//	// 保存配置
//	err := api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	}
func (a *API) AddPackageSource(config *types.NuGetConfig, key string, value string, protocolVersion string) {
	a.Manager.AddPackageSource(config, key, value, protocolVersion)
}

// RemovePackageSource 移除包源
//
// RemovePackageSource 从配置中移除指定键名的包源。
// 如果该包源同时也是活跃包源或在禁用包源列表中，相关条目也会被移除。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 要移除的包源标识符/名称
//
// 返回值:
//   - bool: 如果成功移除返回 true，如果指定的包源不存在返回 false
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 移除不再需要的包源
//	removed := api.RemovePackageSource(config, "old-source")
//	if removed {
//	    fmt.Println("包源已成功移除")
//	} else {
//	    fmt.Println("指定的包源不存在")
//	}
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	}
func (a *API) RemovePackageSource(config *types.NuGetConfig, key string) bool {
	return a.Manager.RemovePackageSource(config, key)
}

// GetPackageSource 获取包源
//
// GetPackageSource 根据键名从配置中获取特定的包源。
//
// 参数:
//   - config: NuGet 配置对象
//   - key: 要获取的包源标识符/名称
//
// 返回值:
//   - *types.PackageSource: 找到的包源对象，如果不存在则返回 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 获取特定包源
//	source := api.GetPackageSource(config, "nuget.org")
//	if source == nil {
//	    fmt.Println("未找到指定的包源")
//	    return
//	}
//
//	// 显示包源详情
//	fmt.Printf("包源: %s\n", source.Key)
//	fmt.Printf("URL: %s\n", source.Value)
//	if source.ProtocolVersion != "" {
//	    fmt.Printf("协议版本: %s\n", source.ProtocolVersion)
//	}
//
//	// 检查包源是否被禁用
//	if api.IsPackageSourceDisabled(config, source.Key) {
//	    fmt.Println("此包源当前已禁用")
//	} else {
//	    fmt.Println("此包源当前已启用")
//	}
func (a *API) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource {
	return a.Manager.GetPackageSource(config, key)
}

// GetAllPackageSources 获取所有包源
//
// GetAllPackageSources 返回配置中定义的所有包源列表。
//
// 参数:
//   - config: NuGet 配置对象
//
// 返回值:
//   - []types.PackageSource: 所有包源的切片，如果没有包源则返回空切片
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 获取所有包源
//	sources := api.GetAllPackageSources(config)
//
//	// 显示包源列表
//	fmt.Printf("找到 %d 个包源:\n", len(sources))
//	for i, source := range sources {
//	    fmt.Printf("%d. %s (%s)\n", i+1, source.Key, source.Value)
//
//	    // 检查每个包源的状态
//	    if api.IsPackageSourceDisabled(config, source.Key) {
//	        fmt.Println("   状态: 已禁用")
//	    } else {
//	        fmt.Println("   状态: 已启用")
//	    }
//	}
func (a *API) GetAllPackageSources(config *types.NuGetConfig) []types.PackageSource {
	return a.Manager.GetAllPackageSources(config)
}

// SetActivePackageSource 设置活跃包源
//
// SetActivePackageSource 将指定的包源设置为活跃包源。
// 活跃包源通常用作默认的包源，当没有明确指定包源时使用。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 要设置为活跃的包源标识符/名称
//
// 返回值:
//   - error: 如果指定的包源不存在，则返回错误；如果成功设置则返回 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 设置活跃包源
//	err = api.SetActivePackageSource(config, "company-internal")
//	if err != nil {
//	    fmt.Printf("设置活跃包源失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("活跃包源已设置")
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	    return
//	}
func (a *API) SetActivePackageSource(config *types.NuGetConfig, key string) error {
	return a.Manager.SetActivePackageSource(config, key)
}

// 凭证操作

// AddCredential 添加包源凭证
//
// AddCredential 为指定的包源添加身份验证凭证（用户名和密码）。
// 如果包源已有凭证，将被更新。如果包源不存在，只会创建凭证信息，但不会自动添加包源。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - sourceKey: 要添加凭证的包源名称
//   - username: 身份验证用户名
//   - password: 身份验证密码（会以明文形式存储在配置文件中，使用 ClearTextPassword 键名）
//
// 注意:
//   - 配置文件中的密码是以明文方式存储的，请确保文件安全
//   - 某些 NuGet 服务器可能需要特殊的凭证类型，如 API 密钥
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 为私有包源添加凭证
//	api.AddCredential(config, "private-source", "username", "p@ssw0rd")
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("凭证已添加/更新")
//
//	// 添加包含特殊字符的凭证示例
//	api.AddCredential(config, "another-source", "user@example.com", "complex!P@ssw0rd#123")
//
//	// 使用 API 密钥作为凭证的示例
//	api.AddCredential(config, "api-source", "apikey", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx")
func (a *API) AddCredential(config *types.NuGetConfig, sourceKey string, username string, password string) {
	a.Manager.AddCredential(config, sourceKey, username, password)
}

// RemoveCredential 移除包源凭证
//
// RemoveCredential 从配置中移除指定包源的身份验证凭证。
// 如果指定的包源没有凭证，则不做任何更改。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - sourceKey: 要移除凭证的包源名称
//
// 返回值:
//   - bool: 如果成功移除返回 true，如果指定的包源没有凭证则返回 false
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 移除不再需要的凭证
//	removed := api.RemoveCredential(config, "private-source")
//
//	if removed {
//	    fmt.Println("包源凭证已成功移除")
//
//	    // 保存修改后的配置
//	    err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	    if err != nil {
//	        fmt.Printf("保存配置失败: %v\n", err)
//	        return
//	    }
//	} else {
//	    fmt.Println("指定的包源没有凭证，无需移除")
//	}
func (a *API) RemoveCredential(config *types.NuGetConfig, sourceKey string) bool {
	return a.Manager.RemoveCredential(config, sourceKey)
}

// 禁用包源操作

// DisablePackageSource 禁用包源
//
// DisablePackageSource 在配置中将指定的包源标记为禁用状态。
// 被禁用的包源仍会在配置文件中保留，但在 NuGet 操作中不会被使用。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 要禁用的包源标识符/名称
//
// 注意:
//   - 如果指定的包源不存在，不会产生错误，但禁用信息仍会被添加到配置中
//   - 如果包源已经被禁用，此操作不会产生变化
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 临时禁用一个包源
//	api.DisablePackageSource(config, "slow-source")
//
//	// 检查是否成功禁用
//	if api.IsPackageSourceDisabled(config, "slow-source") {
//	    fmt.Println("包源已成功禁用")
//	}
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	    return
//	}
func (a *API) DisablePackageSource(config *types.NuGetConfig, key string) {
	a.Manager.DisablePackageSource(config, key)
}

// EnablePackageSource 启用包源
//
// EnablePackageSource 从配置中的禁用列表中移除指定的包源，使其恢复启用状态。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 要启用的包源标识符/名称
//
// 返回值:
//   - bool: 如果成功启用（即包源之前是禁用的）返回 true，如果包源不在禁用列表中返回 false
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 检查包源当前状态
//	if api.IsPackageSourceDisabled(config, "some-source") {
//	    fmt.Println("包源当前被禁用，正在启用...")
//
//	    // 启用包源
//	    enabled := api.EnablePackageSource(config, "some-source")
//
//	    if enabled {
//	        fmt.Println("包源已成功启用")
//
//	        // 保存修改后的配置
//	        err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	        if err != nil {
//	            fmt.Printf("保存配置失败: %v\n", err)
//	            return
//	        }
//	    } else {
//	        fmt.Println("启用操作失败")
//	    }
//	} else {
//	    fmt.Println("包源已经处于启用状态")
//	}
func (a *API) EnablePackageSource(config *types.NuGetConfig, key string) bool {
	return a.Manager.EnablePackageSource(config, key)
}

// IsPackageSourceDisabled 检查包源是否被禁用
//
// IsPackageSourceDisabled 检查指定的包源是否在配置中被标记为禁用状态。
//
// 参数:
//   - config: NuGet 配置对象
//   - key: 要检查的包源标识符/名称
//
// 返回值:
//   - bool: 如果包源被禁用返回 true，否则返回 false。如果包源不存在，也返回 false。
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 获取所有包源
//	sources := api.GetAllPackageSources(config)
//
//	// 检查每个包源的启用状态
//	for _, source := range sources {
//	    if api.IsPackageSourceDisabled(config, source.Key) {
//	        fmt.Printf("包源 %s 已禁用\n", source.Key)
//	    } else {
//	        fmt.Printf("包源 %s 已启用\n", source.Key)
//	    }
//	}
//
//	// 根据状态执行不同操作
//	targetSource := "some-source"
//	if api.IsPackageSourceDisabled(config, targetSource) {
//	    fmt.Printf("包源 %s 已禁用，跳过\n", targetSource)
//	} else {
//	    // 使用启用的包源
//	    fmt.Printf("使用包源 %s\n", targetSource)
//	}
func (a *API) IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool {
	return a.Manager.IsPackageSourceDisabled(config, key)
}

// 配置选项操作

// AddConfigOption 添加配置选项
//
// AddConfigOption 向配置中添加或更新全局配置选项。
// 配置选项是影响 NuGet 行为的键值对，例如代理设置、全局包文件夹路径等。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 配置选项的键名
//   - value: 配置选项的值
//
// 注意:
//   - 如果指定键名的选项已存在，将更新其值
//   - 常见的配置选项包括：globalPackagesFolder、http_proxy、no_proxy 等
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 设置全局包文件夹路径
//	api.AddConfigOption(config, "globalPackagesFolder", "/custom/nuget/packages")
//
//	// 设置代理
//	api.AddConfigOption(config, "http_proxy", "http://proxy.example.com:8080")
//
//	// 设置包版本行为
//	api.AddConfigOption(config, "dependencyVersion", "Highest")
//
//	// 保存修改后的配置
//	err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("保存配置失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("配置选项已设置")
func (a *API) AddConfigOption(config *types.NuGetConfig, key string, value string) {
	a.Manager.AddConfigOption(config, key, value)
}

// RemoveConfigOption 移除配置选项
//
// RemoveConfigOption 从配置中移除指定的全局配置选项。
//
// 参数:
//   - config: 要修改的 NuGet 配置对象
//   - key: 要移除的配置选项键名
//
// 返回值:
//   - bool: 如果成功移除选项返回 true，如果指定的选项不存在返回 false
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 移除不再需要的代理设置
//	removed := api.RemoveConfigOption(config, "http_proxy")
//
//	if removed {
//	    fmt.Println("代理设置已移除")
//
//	    // 保存修改后的配置
//	    err = api.SaveConfig(config, "/path/to/NuGet.Config")
//	    if err != nil {
//	        fmt.Printf("保存配置失败: %v\n", err)
//	        return
//	    }
//	} else {
//	    fmt.Println("指定的配置选项不存在")
//	}
func (a *API) RemoveConfigOption(config *types.NuGetConfig, key string) bool {
	return a.Manager.RemoveConfigOption(config, key)
}

// GetConfigOption 获取配置选项值
//
// GetConfigOption 获取指定全局配置选项的值。
//
// 参数:
//   - config: NuGet 配置对象
//   - key: 要获取的配置选项键名
//
// 返回值:
//   - string: 配置选项的值。如果指定的选项不存在，返回空字符串。
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 加载配置
//	config, err := api.ParseFromFile("/path/to/NuGet.Config")
//	if err != nil {
//	    fmt.Printf("加载配置失败: %v\n", err)
//	    return
//	}
//
//	// 获取全局包文件夹路径
//	packagesFolder := api.GetConfigOption(config, "globalPackagesFolder")
//	if packagesFolder != "" {
//	    fmt.Printf("全局包文件夹: %s\n", packagesFolder)
//	} else {
//	    fmt.Println("未设置全局包文件夹")
//	}
//
//	// 检查是否设置了代理
//	proxy := api.GetConfigOption(config, "http_proxy")
//	if proxy != "" {
//	    fmt.Printf("HTTP 代理: %s\n", proxy)
//	} else {
//	    fmt.Println("未设置 HTTP 代理")
//	}
func (a *API) GetConfigOption(config *types.NuGetConfig, key string) string {
	return a.Manager.GetConfigOption(config, key)
}

// SerializeToXML 将配置序列化为XML字符串
//
// SerializeToXML 将 NuGet 配置对象序列化为标准格式的 XML 字符串。
// 生成的 XML 包含适当的缩进和 XML 声明。
//
// 参数:
//   - config: 要序列化的 NuGet 配置对象
//
// 返回值:
//   - string: 序列化后的 XML 字符串
//   - error: 如果序列化过程中发生错误则返回相应的错误；如果成功则为 nil
//
// 示例:
//
//	api := nuget.NewAPI()
//
//	// 创建或加载配置
//	config := api.CreateDefaultConfig()
//	api.AddPackageSource(config, "custom-source", "https://custom-source/v3/index.json", "3")
//
//	// 序列化为 XML
//	xmlString, err := api.SerializeToXML(config)
//	if err != nil {
//	    fmt.Printf("序列化失败: %v\n", err)
//	    return
//	}
//
//	// 显示生成的 XML
//	fmt.Println("生成的 XML 配置:")
//	fmt.Println(xmlString)
//
//	// 也可以将生成的 XML 用于其他目的，例如传递给其他系统
//	// 或在不写入文件的情况下使用
//
//	// 解析序列化后的 XML 以验证其有效性
//	parsedConfig, err := api.ParseFromString(xmlString)
//	if err != nil {
//	    fmt.Printf("验证失败: %v\n", err)
//	    return
//	}
//
//	fmt.Println("XML 有效且可以正确解析")
func (a *API) SerializeToXML(config *types.NuGetConfig) (string, error) {
	return a.Parser.SerializeToXML(config)
}
