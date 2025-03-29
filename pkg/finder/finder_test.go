package finder

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	nugetTesting "github.com/scagogogo/nuget-config-parser/pkg/testing"
)

func TestNewConfigFinder(t *testing.T) {
	finder := NewConfigFinder()

	if finder == nil {
		t.Fatal("NewConfigFinder() returned nil")
	}

	if finder.EnvVariableName != "NUGET_CONFIG_FILE" {
		t.Errorf("EnvVariableName = %q, want %q", finder.EnvVariableName, "NUGET_CONFIG_FILE")
	}
}

func TestGetConfigFileSearchLocations(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	envConfigPath := filepath.Join(tempDir, "env-nuget.config")
	nugetTesting.CreateNuGetConfigFile(t, envConfigPath, nugetTesting.ValidNuGetConfig())

	// 设置环境变量
	cleanup := nugetTesting.SetupEnv(t, "NUGET_CONFIG_FILE", envConfigPath)
	defer cleanup()

	// 创建 finder
	finder := NewConfigFinder()

	// 测试获取配置文件搜索位置
	locations := finder.GetConfigFileSearchLocations()

	// 检查环境变量指定的位置
	if !contains(locations, envConfigPath) {
		t.Errorf("GetConfigFileSearchLocations() should include environment variable path")
	}

	// 检查默认位置
	defaultLocations := constants.GetDefaultConfigLocations()
	for _, loc := range defaultLocations {
		if !contains(locations, loc) {
			t.Errorf("GetConfigFileSearchLocations() should include default location %s", loc)
		}
	}

	// 测试环境变量指向不存在的文件
	nonExistentPath := filepath.Join(tempDir, "non-existent.config")
	cleanup = nugetTesting.SetupEnv(t, "NUGET_CONFIG_FILE", nonExistentPath)
	defer cleanup()

	locations = finder.GetConfigFileSearchLocations()
	if contains(locations, nonExistentPath) {
		t.Errorf("GetConfigFileSearchLocations() should not include non-existent environment variable path")
	}
}

func TestFindConfigFile(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	nugetTesting.CreateNuGetConfigFile(t, configPath, nugetTesting.ValidNuGetConfig())

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

	// 创建 finder
	finder := NewConfigFinder()

	// 测试查找配置文件
	found, err := finder.FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() error = %v", err)
	}

	// 检查找到的路径是否与创建的配置文件匹配
	absConfigPath, _ := filepath.Abs(configPath)
	absFound, _ := filepath.Abs(found)
	if !pathsEqual(absFound, absConfigPath) {
		t.Errorf("FindConfigFile() = %q, want %q", absFound, absConfigPath)
	}

	// 删除配置文件后测试
	os.Remove(configPath)
	_, err = finder.FindConfigFile()
	if err == nil {
		t.Errorf("FindConfigFile() should return error when no config file found")
	}

	// 测试环境变量优先级
	envConfigPath := filepath.Join(tempDir, "env-nuget.config")
	nugetTesting.CreateNuGetConfigFile(t, envConfigPath, nugetTesting.ValidNuGetConfig())
	cleanup := nugetTesting.SetupEnv(t, "NUGET_CONFIG_FILE", envConfigPath)
	defer cleanup()

	// 重新创建默认配置文件
	nugetTesting.CreateNuGetConfigFile(t, configPath, nugetTesting.ValidNuGetConfig())

	found, err = finder.FindConfigFile()
	if err != nil {
		t.Fatalf("FindConfigFile() with env var error = %v", err)
	}

	// 环境变量指定的路径应该优先
	absEnvConfigPath, _ := filepath.Abs(envConfigPath)
	absFound, _ = filepath.Abs(found)
	if !pathsEqual(absFound, absEnvConfigPath) {
		t.Errorf("FindConfigFile() with env var = %q, want %q", absFound, absEnvConfigPath)
	}
}

func TestFindAllConfigFiles(t *testing.T) {
	// 创建临时目录结构
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建多个配置文件
	configPaths := []string{
		filepath.Join(tempDir, constants.DefaultNuGetConfigFilename),
		filepath.Join(tempDir, "sub", constants.DefaultNuGetConfigFilename),
		filepath.Join(tempDir, "env-nuget.config"),
	}

	for _, path := range configPaths {
		nugetTesting.CreateNuGetConfigFile(t, path, nugetTesting.ValidNuGetConfig())
	}

	// 设置环境变量
	cleanup := nugetTesting.SetupEnv(t, "NUGET_CONFIG_FILE", configPaths[2])
	defer cleanup()

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

	// 创建 finder
	finder := NewConfigFinder()

	// 测试找到所有配置文件
	found := finder.FindAllConfigFiles()

	// 检查是否找到了所有创建的配置文件
	if len(found) < 2 {
		t.Errorf("FindAllConfigFiles() found %d files, want at least 2", len(found))
	}

	// 环境变量指定的文件应该在列表中
	foundEnvConfig := false
	absEnvConfigPath, _ := filepath.Abs(configPaths[2])
	for _, f := range found {
		absPath, _ := filepath.Abs(f)
		if pathsEqual(absPath, absEnvConfigPath) {
			foundEnvConfig = true
			break
		}
	}
	if !foundEnvConfig {
		t.Errorf("FindAllConfigFiles() should include environment variable specified config")
	}

	// 当前目录的配置文件应该在列表中
	foundCurrentConfig := false
	absCurrentConfigPath, _ := filepath.Abs(configPaths[0])
	for _, f := range found {
		absPath, _ := filepath.Abs(f)
		if pathsEqual(absPath, absCurrentConfigPath) {
			foundCurrentConfig = true
			break
		}
	}
	if !foundCurrentConfig {
		t.Errorf("FindAllConfigFiles() should include current directory config")
	}
}

func TestFindProjectConfig(t *testing.T) {
	// 创建临时目录结构
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建嵌套目录
	subDir := filepath.Join(tempDir, "sub")
	subSubDir := filepath.Join(subDir, "subsub")
	if err := os.MkdirAll(subSubDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectories: %v", err)
	}

	// 在根目录创建配置文件
	rootConfigPath := filepath.Join(tempDir, constants.DefaultNuGetConfigFilename)
	nugetTesting.CreateNuGetConfigFile(t, rootConfigPath, nugetTesting.ValidNuGetConfig())

	// 创建 finder
	finder := NewConfigFinder()

	// 从子子目录开始查找
	found, err := finder.FindProjectConfig(subSubDir)
	if err != nil {
		t.Fatalf("FindProjectConfig() error = %v", err)
	}

	// 检查找到的路径是否为根目录的配置文件
	absRootConfigPath, _ := filepath.Abs(rootConfigPath)
	absFound, _ := filepath.Abs(found)
	if !pathsEqual(absFound, absRootConfigPath) {
		t.Errorf("FindProjectConfig() = %q, want %q", absFound, absRootConfigPath)
	}

	// 在子目录创建配置文件
	subConfigPath := filepath.Join(subDir, constants.DefaultNuGetConfigFilename)
	nugetTesting.CreateNuGetConfigFile(t, subConfigPath, nugetTesting.ValidNuGetConfig())

	// 从子子目录开始查找
	found, err = finder.FindProjectConfig(subSubDir)
	if err != nil {
		t.Fatalf("FindProjectConfig() error = %v", err)
	}

	// 检查找到的路径是否为子目录的配置文件（应该找最近的）
	absSubConfigPath, _ := filepath.Abs(subConfigPath)
	absFound, _ = filepath.Abs(found)
	if !pathsEqual(absFound, absSubConfigPath) {
		t.Errorf("FindProjectConfig() = %q, want %q", absFound, absSubConfigPath)
	}

	// 删除所有配置文件后测试
	os.Remove(rootConfigPath)
	os.Remove(subConfigPath)
	_, err = finder.FindProjectConfig(subSubDir)
	if err == nil {
		t.Errorf("FindProjectConfig() should return error when no config file found")
	}
}

func TestGetUserConfigFile(t *testing.T) {
	finder := NewConfigFinder()
	userConfigPath := finder.GetUserConfigFile()

	if userConfigPath == "" {
		t.Error("GetUserConfigFile() returned empty string")
	}

	if !filepath.IsAbs(userConfigPath) {
		t.Errorf("GetUserConfigFile() = %q, should be an absolute path", userConfigPath)
	}

	// 用户配置路径应该包含 NuGet 目录名
	if !containsSubstring(filepath.ToSlash(userConfigPath), constants.GlobalFolderName+"/") {
		t.Errorf("GetUserConfigFile() = %q, should contain NuGet folder", userConfigPath)
	}

	// 用户配置路径应该以配置文件名结尾
	if filepath.Base(userConfigPath) != constants.DefaultNuGetConfigFilename {
		t.Errorf("GetUserConfigFile() = %q, should end with %s", userConfigPath, constants.DefaultNuGetConfigFilename)
	}
}

func TestGetMachineConfigFile(t *testing.T) {
	finder := NewConfigFinder()
	machineConfigPath := finder.GetMachineConfigFile()

	if machineConfigPath == "" {
		t.Error("GetMachineConfigFile() returned empty string")
	}

	if !filepath.IsAbs(machineConfigPath) {
		t.Errorf("GetMachineConfigFile() = %q, should be an absolute path", machineConfigPath)
	}

	// 机器配置路径应该包含 NuGet 目录名
	if !containsSubstring(filepath.ToSlash(machineConfigPath), constants.GlobalFolderName+"/") {
		t.Errorf("GetMachineConfigFile() = %q, should contain NuGet folder", machineConfigPath)
	}

	// 机器配置路径应该以配置文件名结尾
	if filepath.Base(machineConfigPath) != constants.DefaultNuGetConfigFilename {
		t.Errorf("GetMachineConfigFile() = %q, should end with %s", machineConfigPath, constants.DefaultNuGetConfigFilename)
	}
}

// 辅助函数：检查切片是否包含指定字符串
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// 辅助函数：检查字符串是否包含子字符串
func containsSubstring(s, substr string) bool {
	return strings.Contains(s, substr)
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
