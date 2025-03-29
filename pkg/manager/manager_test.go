package manager

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	"github.com/scagogogo/nuget-config-parser/pkg/parser"
	nugetTesting "github.com/scagogogo/nuget-config-parser/pkg/testing"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

func TestNewConfigManager(t *testing.T) {
	manager := NewConfigManager()

	if manager == nil {
		t.Fatal("NewConfigManager() returned nil")
	}

	if manager.parser == nil {
		t.Error("ConfigManager.parser is nil")
	}

	if manager.finder == nil {
		t.Error("ConfigManager.finder is nil")
	}
}

func TestLoadConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 解析为结构体以便后续比较
	p := parser.NewConfigParser()
	validConfig, err := p.ParseFromString(validConfigXML)
	if err != nil {
		t.Fatalf("Failed to parse valid config: %v", err)
	}

	// 创建 ConfigManager
	manager := NewConfigManager()

	// 保存当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// 测试读取配置
	config, err := manager.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// 检查包源数量
	if len(config.PackageSources.Add) != len(validConfig.PackageSources.Add) {
		t.Errorf("Got %d package sources, want %d", len(config.PackageSources.Add), len(validConfig.PackageSources.Add))
	}

	// 测试读取不存在的配置
	os.Remove(configPath)
	_, err = manager.LoadConfig(configPath)
	if err == nil {
		t.Error("LoadConfig() should return error when config file doesn't exist")
	}
}

func TestGetNuGetConfigFromPath(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 解析为结构体以便后续比较
	p := parser.NewConfigParser()
	validConfig, err := p.ParseFromString(validConfigXML)
	if err != nil {
		t.Fatalf("Failed to parse valid config: %v", err)
	}

	// 创建 ConfigManager
	manager := NewConfigManager()

	// 测试从指定路径读取配置
	config, err := manager.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// 检查包源数量
	if len(config.PackageSources.Add) != len(validConfig.PackageSources.Add) {
		t.Errorf("Got %d package sources, want %d", len(config.PackageSources.Add), len(validConfig.PackageSources.Add))
	}

	// 测试读取不存在的配置文件
	nonExistentPath := filepath.Join(tempDir, "non-existent.config")
	_, err = manager.LoadConfig(nonExistentPath)
	if err == nil {
		t.Error("LoadConfig() should return error for non-existent file")
	}

	// 测试读取无效的配置文件
	invalidConfigPath := filepath.Join(tempDir, "invalid.config")
	nugetTesting.CreateNuGetConfigFile(t, invalidConfigPath, "invalid xml content")
	_, err = manager.LoadConfig(invalidConfigPath)
	if err == nil {
		t.Error("LoadConfig() should return error for invalid file")
	}
}

// pathsEqual 比较两个路径是否相等，考虑平台特性
func pathsEqual(path1, path2 string) bool {
	// 清理路径
	path1 = filepath.Clean(path1)
	path2 = filepath.Clean(path2)

	// macOS 上 /private/var 和 /var 是符号链接
	if runtime.GOOS == "darwin" {
		path1 = strings.Replace(path1, "/private/var/", "/var/", 1)
		path2 = strings.Replace(path2, "/private/var/", "/var/", 1)
	}

	return path1 == path2
}

func TestFindAndLoadConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	validConfigXML := nugetTesting.ValidNuGetConfig()
	nugetTesting.CreateNuGetConfigFile(t, configPath, validConfigXML)

	// 保存当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// 创建 ConfigManager
	manager := NewConfigManager()

	// 测试查找并加载配置
	config, foundPath, err := manager.FindAndLoadConfig()
	if err != nil {
		t.Fatalf("FindAndLoadConfig() error = %v", err)
	}

	// 检查配置内容
	if config == nil {
		t.Fatal("FindAndLoadConfig() returned nil config")
	}

	// 检查找到的路径
	absConfigPath, _ := filepath.Abs(configPath)
	absFoundPath, _ := filepath.Abs(foundPath)
	if !pathsEqual(absFoundPath, absConfigPath) {
		t.Errorf("FindAndLoadConfig() found path = %q, want %q", absFoundPath, absConfigPath)
	}
}

func TestSaveConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件路径
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)

	// 创建 ConfigManager
	manager := NewConfigManager()

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
	err := manager.SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("SaveConfig() did not create file at %s", configPath)
	}

	// 读取保存的配置
	loadedConfig, err := manager.LoadConfig(configPath)
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

func TestCreateDefaultConfig(t *testing.T) {
	manager := NewConfigManager()
	config := manager.CreateDefaultConfig()

	if config == nil {
		t.Fatal("CreateDefaultConfig() returned nil")
	}

	if len(config.PackageSources.Add) == 0 {
		t.Error("Default config has no package sources")
	}

	if config.PackageSources.Add[0].Key != "nuget.org" {
		t.Errorf("Default source key = %q, want %q", config.PackageSources.Add[0].Key, "nuget.org")
	}

	if config.PackageSources.Add[0].Value != constants.DefaultPackageSource {
		t.Errorf("Default source value = %q, want %q", config.PackageSources.Add[0].Value, constants.DefaultPackageSource)
	}

	if config.ActivePackageSource == nil {
		t.Error("Default config has no active package source")
	}
}

func TestInitializeDefaultConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件路径
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)

	// 创建 ConfigManager
	manager := NewConfigManager()

	// 测试初始化默认配置
	err := manager.InitializeDefaultConfig(configPath)
	if err != nil {
		t.Fatalf("InitializeDefaultConfig() error = %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("InitializeDefaultConfig() did not create file at %s", configPath)
	}

	// 读取创建的配置
	config, err := manager.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load initialized config: %v", err)
	}

	// 检查配置内容
	if len(config.PackageSources.Add) == 0 {
		t.Error("Initialized config has no package sources")
	}

	if config.PackageSources.Add[0].Key != "nuget.org" {
		t.Errorf("Initialized source key = %q, want %q", config.PackageSources.Add[0].Key, "nuget.org")
	}
}
