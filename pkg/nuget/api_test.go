package nuget

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	nugetTesting "github.com/scagogogo/nuget-config-parser/pkg/testing"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

func TestNewAPI(t *testing.T) {
	api := NewAPI()

	if api == nil {
		t.Fatal("NewAPI() returned nil")
	}

	if api.Parser == nil {
		t.Error("API.Parser is nil")
	}

	if api.Finder == nil {
		t.Error("API.Finder is nil")
	}

	if api.Manager == nil {
		t.Error("API.Manager is nil")
	}
}

func TestAPIParseFromFile(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 创建 API
	api := NewAPI()

	// 测试从文件解析
	config, err := api.ParseFromFile(configPath)
	if err != nil {
		t.Fatalf("ParseFromFile() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("ParseFromFile() returned nil config")
	}

	// 检查包源数量
	if len(config.PackageSources.Add) == 0 {
		t.Error("ParseFromFile() returned config with no package sources")
	}
}

func TestAPIParseFromString(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 测试从字符串解析
	validConfigXML := nugetTesting.ValidNuGetConfig()
	config, err := api.ParseFromString(validConfigXML)
	if err != nil {
		t.Fatalf("ParseFromString() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("ParseFromString() returned nil config")
	}

	// 检查包源数量
	if len(config.PackageSources.Add) == 0 {
		t.Error("ParseFromString() returned config with no package sources")
	}

	// 测试无效字符串
	_, err = api.ParseFromString("invalid xml")
	if err == nil {
		t.Error("ParseFromString() should return error for invalid XML")
	}
}

func TestAPIParseFromReader(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 从字符串创建读取器
	validConfigXML := nugetTesting.ValidNuGetConfig()
	reader := nugetTesting.StringReader(validConfigXML)

	// 测试从读取器解析
	config, err := api.ParseFromReader(reader)
	if err != nil {
		t.Fatalf("ParseFromReader() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("ParseFromReader() returned nil config")
	}

	// 检查包源数量
	if len(config.PackageSources.Add) == 0 {
		t.Error("ParseFromReader() returned config with no package sources")
	}

	// 测试无效读取器
	invalidReader := nugetTesting.StringReader("invalid xml")
	_, err = api.ParseFromReader(invalidReader)
	if err == nil {
		t.Error("ParseFromReader() should return error for invalid XML")
	}
}

func TestAPIFindConfigFile(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 创建 API
	api := NewAPI()

	// 保存当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// 测试查找配置文件
	foundPath, err := api.FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	// 检查找到的路径
	absConfigPath, _ := filepath.Abs(configPath)
	absFoundPath, _ := filepath.Abs(foundPath)
	if !pathsEqual(absFoundPath, absConfigPath) {
		t.Errorf("FindConfigFile() = %q, want %q", absFoundPath, absConfigPath)
	}
}

func TestAPIFindAndParseConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 创建 API
	api := NewAPI()

	// 保存当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// 测试查找并解析配置
	config, foundPath, err := api.FindAndParseConfig()
	if err != nil {
		t.Fatalf("FindAndParseConfig() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("FindAndParseConfig() returned nil config")
	}

	// 检查找到的路径
	absConfigPath, _ := filepath.Abs(configPath)
	absFoundPath, _ := filepath.Abs(foundPath)
	if !pathsEqual(absFoundPath, absConfigPath) {
		t.Errorf("FindAndParseConfig() found path = %q, want %q", absFoundPath, absConfigPath)
	}
}

func TestAPISaveConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件路径
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)

	// 创建 API
	api := NewAPI()

	// 创建配置对象
	config := &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{
				{
					Key:             "test-source",
					Value:           "http://test.example.com",
					ProtocolVersion: "3",
				},
			},
		},
	}

	// 测试保存配置
	err := api.SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("SaveConfig() did not create file at %s", configPath)
	}

	// 读取保存的配置
	loadedConfig, err := api.ParseFromFile(configPath)
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	// 检查配置内容
	if len(loadedConfig.PackageSources.Add) != len(config.PackageSources.Add) {
		t.Errorf("Loaded config has %d package sources, want %d", len(loadedConfig.PackageSources.Add), len(config.PackageSources.Add))
	}

	// 检查包源属性
	if len(loadedConfig.PackageSources.Add) > 0 {
		source := loadedConfig.PackageSources.Add[0]
		if source.Key != "test-source" {
			t.Errorf("Loaded source key = %q, want %q", source.Key, "test-source")
		}
		if source.Value != "http://test.example.com" {
			t.Errorf("Loaded source value = %q, want %q", source.Value, "http://test.example.com")
		}
	}
}

func TestAPIPackageSourceOperations(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 创建配置对象
	config := api.CreateDefaultConfig()

	// 测试添加包源
	api.AddPackageSource(config, "test-source", "http://test.example.com", "3")

	// 检查是否添加成功
	found := false
	for _, source := range config.PackageSources.Add {
		if source.Key == "test-source" {
			found = true
			if source.Value != "http://test.example.com" {
				t.Errorf("Added source value = %q, want %q", source.Value, "http://test.example.com")
			}
			if source.ProtocolVersion != "3" {
				t.Errorf("Added source protocol version = %q, want %q", source.ProtocolVersion, "3")
			}
			break
		}
	}
	if !found {
		t.Error("AddPackageSource() did not add the package source")
	}

	// 测试获取包源
	source := api.GetPackageSource(config, "test-source")
	if source == nil {
		t.Error("GetPackageSource() returned nil for existing source")
	} else {
		if source.Value != "http://test.example.com" {
			t.Errorf("GetPackageSource() value = %q, want %q", source.Value, "http://test.example.com")
		}
	}

	// 测试获取所有包源
	sources := api.GetAllPackageSources(config)
	if len(sources) < 2 { // 默认源 + 新添加的源
		t.Errorf("GetAllPackageSources() returned %d sources, want at least 2", len(sources))
	}

	// 测试设置活跃包源
	err := api.SetActivePackageSource(config, "test-source")
	if err != nil {
		t.Errorf("SetActivePackageSource() error = %v", err)
	}
	if config.ActivePackageSource == nil || config.ActivePackageSource.Add.Key != "test-source" {
		t.Error("SetActivePackageSource() did not set the active package source")
	}

	// 测试禁用包源
	api.DisablePackageSource(config, "test-source")
	if !api.IsPackageSourceDisabled(config, "test-source") {
		t.Error("DisablePackageSource() did not disable the package source")
	}

	// 测试启用包源
	enabled := api.EnablePackageSource(config, "test-source")
	if !enabled || api.IsPackageSourceDisabled(config, "test-source") {
		t.Error("EnablePackageSource() did not enable the package source")
	}

	// 测试移除包源
	removed := api.RemovePackageSource(config, "test-source")
	if !removed {
		t.Error("RemovePackageSource() did not remove the package source")
	}
	if api.GetPackageSource(config, "test-source") != nil {
		t.Error("RemovePackageSource() did not actually remove the package source")
	}
}

func TestAPICredentialOperations(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 创建配置对象
	config := api.CreateDefaultConfig()

	// 测试添加凭证
	api.AddCredential(config, "test-source", "username", "password")

	// 检查凭证是否添加成功
	if config.PackageSourceCredentials == nil {
		t.Error("AddCredential() did not initialize PackageSourceCredentials")
	} else {
		cred, exists := config.PackageSourceCredentials.Sources["test-source"]
		if !exists {
			t.Error("AddCredential() did not add credentials for the source")
		} else {
			usernameFound := false
			passwordFound := false
			for _, c := range cred.Add {
				if c.Key == "Username" && c.Value == "username" {
					usernameFound = true
				}
				if c.Key == "ClearTextPassword" && c.Value == "password" {
					passwordFound = true
				}
			}
			if !usernameFound {
				t.Error("AddCredential() did not add username")
			}
			if !passwordFound {
				t.Error("AddCredential() did not add password")
			}
		}
	}

	// 测试移除凭证
	removed := api.RemoveCredential(config, "test-source")
	if !removed {
		t.Error("RemoveCredential() did not remove the credentials")
	}
	if config.PackageSourceCredentials != nil && len(config.PackageSourceCredentials.Sources) > 0 {
		if _, exists := config.PackageSourceCredentials.Sources["test-source"]; exists {
			t.Error("RemoveCredential() did not actually remove the credentials")
		}
	}
}

func TestAPIConfigOptionOperations(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 创建配置对象
	config := api.CreateDefaultConfig()

	// 测试添加配置选项
	api.AddConfigOption(config, "test-option", "test-value")

	// 检查配置选项是否添加成功
	if config.Config == nil {
		t.Error("AddConfigOption() did not initialize Config")
	} else {
		optionFound := false
		for _, option := range config.Config.Add {
			if option.Key == "test-option" {
				optionFound = true
				if option.Value != "test-value" {
					t.Errorf("Added option value = %q, want %q", option.Value, "test-value")
				}
				break
			}
		}
		if !optionFound {
			t.Error("AddConfigOption() did not add the option")
		}
	}

	// 测试获取配置选项
	value := api.GetConfigOption(config, "test-option")
	if value != "test-value" {
		t.Errorf("GetConfigOption() = %q, want %q", value, "test-value")
	}

	// 测试移除配置选项
	removed := api.RemoveConfigOption(config, "test-option")
	if !removed {
		t.Error("RemoveConfigOption() did not remove the option")
	}
	value = api.GetConfigOption(config, "test-option")
	if value != "" {
		t.Errorf("RemoveConfigOption() did not actually remove the option, got value %q", value)
	}
}

func TestAPIXMLSerialization(t *testing.T) {
	// 创建 API
	api := NewAPI()

	// 创建配置对象
	config := api.CreateDefaultConfig()

	// 测试序列化为 XML
	xml, err := api.SerializeToXML(config)
	if err != nil {
		t.Fatalf("SerializeToXML() error = %v", err)
	}

	// 检查 XML 内容
	if len(xml) == 0 {
		t.Error("SerializeToXML() returned empty string")
	}
	if xml[:5] != "<?xml" {
		t.Errorf("SerializeToXML() result does not start with XML declaration: %s", xml[:10])
	}

	// 测试从序列化的 XML 解析回对象
	parsedConfig, err := api.ParseFromString(xml)
	if err != nil {
		t.Fatalf("Failed to parse serialized XML: %v", err)
	}
	if len(parsedConfig.PackageSources.Add) != len(config.PackageSources.Add) {
		t.Errorf("Parsed config has %d package sources, want %d", len(parsedConfig.PackageSources.Add), len(config.PackageSources.Add))
	}
}

// 辅助函数：检查两个路径是否等价（处理符号链接等情况）
func pathsEqual(path1, path2 string) bool {
	// 标准化路径
	path1 = filepath.Clean(path1)
	path2 = filepath.Clean(path2)

	// 直接比较
	if path1 == path2 {
		return true
	}

	// 在macOS上，/private/var和/var是等价的
	if runtime.GOOS == "darwin" {
		// 如果path1以/private/开头，尝试去掉/private/再比较
		if strings.HasPrefix(path1, "/private") {
			withoutPrivate := strings.Replace(path1, "/private", "", 1)
			if withoutPrivate == path2 {
				return true
			}
		}

		// 如果path2以/private/开头，尝试去掉/private/再比较
		if strings.HasPrefix(path2, "/private") {
			withoutPrivate := strings.Replace(path2, "/private", "", 1)
			if withoutPrivate == path1 {
				return true
			}
		}
	}

	// 可以在这里添加更多的平台特定路径比较逻辑
	return false
}
