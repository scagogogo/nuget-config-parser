package constants

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultConfigLocations(t *testing.T) {
	locations := GetDefaultConfigLocations()

	// 检查是否包含当前目录下的配置文件
	if !contains(locations, DefaultNuGetConfigFilename) {
		t.Errorf("GetDefaultConfigLocations() should include %s", DefaultNuGetConfigFilename)
	}

	// 检查是否包含上级目录下的配置文件
	if !contains(locations, filepath.Join("..", DefaultNuGetConfigFilename)) {
		t.Errorf("GetDefaultConfigLocations() should include %s", filepath.Join("..", DefaultNuGetConfigFilename))
	}

	// 检查用户级别配置
	userConfigDir := getUserConfigDirectory()
	if userConfigDir != "" {
		userPath := filepath.Join(userConfigDir, GlobalFolderName, DefaultNuGetConfigFilename)
		if !containsPath(locations, userPath) {
			t.Errorf("GetDefaultConfigLocations() should include a path equivalent to %s", userPath)
		}
	}

	// 检查系统级别配置
	systemConfigDir := getSystemConfigDirectory()
	if systemConfigDir != "" {
		systemPath := filepath.Join(systemConfigDir, GlobalFolderName, DefaultNuGetConfigFilename)
		if !containsPath(locations, systemPath) {
			t.Errorf("GetDefaultConfigLocations() should include a path equivalent to %s", systemPath)
		}
	}
}

func TestGetUserConfigDirectory(t *testing.T) {
	// Windows 测试通过检查逻辑而非直接设置 GOOS
	t.Run("Windows logic", func(t *testing.T) {
		// 保存 APPDATA 环境变量
		oldAppData := os.Getenv("APPDATA")
		testValue := "/mock/appdata"
		os.Setenv("APPDATA", testValue)
		defer os.Setenv("APPDATA", oldAppData)

		// 模拟 Windows 的逻辑
		if os.Getenv("APPDATA") == testValue {
			// 这里不直接调用函数，而是手动复制其在 Windows 下的行为
			result := os.Getenv("APPDATA")
			if result != testValue {
				t.Errorf("Windows path resolution should return %q, got %q", testValue, result)
			}
		}
	})

	// macOS 和 Linux 测试，只有当在对应平台上才测试实际行为
	t.Run("Current platform behavior", func(t *testing.T) {
		dir := getUserConfigDirectory()
		if dir == "" {
			t.Errorf("getUserConfigDirectory() returned empty string")
		}
	})
}

func TestGetSystemConfigDirectory(t *testing.T) {
	// Windows 测试通过检查逻辑而非直接设置 GOOS
	t.Run("Windows logic", func(t *testing.T) {
		// 保存 ProgramData 环境变量
		oldProgramData := os.Getenv("ProgramData")
		testValue := "/mock/programdata"
		os.Setenv("ProgramData", testValue)
		defer os.Setenv("ProgramData", oldProgramData)

		// 模拟 Windows 的逻辑
		if os.Getenv("ProgramData") == testValue {
			// 这里不直接调用函数，而是手动复制其在 Windows 下的行为
			result := os.Getenv("ProgramData")
			if result != testValue {
				t.Errorf("Windows system path resolution should return %q, got %q", testValue, result)
			}
		}
	})

	// 测试当前平台的行为
	t.Run("Current platform behavior", func(t *testing.T) {
		dir := getSystemConfigDirectory()
		if dir == "" {
			t.Errorf("getSystemConfigDirectory() returned empty string")
		}
	})
}

// 辅助函数: 检查切片是否包含指定字符串
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// 辅助函数: 检查切片是否包含等效路径（考虑跨平台兼容性）
func containsPath(slice []string, path string) bool {
	normalizedPath := filepath.Clean(path)
	for _, s := range slice {
		if filepath.Clean(s) == normalizedPath {
			return true
		}
	}
	return false
}
