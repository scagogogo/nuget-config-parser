// Package constants 定义 NuGet 配置相关的常量
package constants

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	// DefaultNuGetConfigFilename 默认的NuGet配置文件名
	DefaultNuGetConfigFilename = "NuGet.Config"

	// GlobalFolderName 全局配置文件夹名
	GlobalFolderName = "NuGet"

	// FeedNamePrefix 包源名称前缀
	FeedNamePrefix = "PackageSource"

	// DefaultPackageSource 默认包源URL
	DefaultPackageSource = "https://api.nuget.org/v3/index.json"

	// NuGetV3APIProtocolVersion NuGetV3 API协议版本
	NuGetV3APIProtocolVersion = "3"

	// NuGetV2APIProtocolVersion NuGetV2 API协议版本
	NuGetV2APIProtocolVersion = "2"
)

// GetDefaultConfigLocations 返回默认的NuGet配置文件可能的位置列表
//
// GetDefaultConfigLocations 按照 NuGet 的配置文件查找规则，返回一个包含所有可能的配置文件
// 路径的排序列表。返回的路径遵循以下顺序：
//  1. 当前目录中的 NuGet.Config
//  2. 父目录中的 NuGet.Config
//  3. 用户级别的 NuGet.Config（特定于操作系统的位置）
//  4. 系统级别的 NuGet.Config（特定于操作系统的位置）
//
// 用户和系统级别的配置文件路径因操作系统而异：
//   - Windows：%APPDATA%\NuGet 和 %ProgramData%\NuGet
//   - macOS：~/Library/Application Support/NuGet 和 /Library/Application Support/NuGet
//   - Linux：~/.config/NuGet 和 /etc/NuGet
//
// 返回值:
//   - []string: 包含所有可能的配置文件路径的切片，按照优先级排序
//
// 示例:
//
//	// 获取默认配置文件位置列表
//	locations := constants.GetDefaultConfigLocations()
//
//	// 显示所有位置
//	fmt.Println("NuGet 配置文件的可能位置:")
//	for i, location := range locations {
//	    fmt.Printf("%d. %s\n", i+1, location)
//	}
//
//	// 检查每个位置的文件是否存在
//	for _, location := range locations {
//	    if utils.FileExists(location) {
//	        fmt.Printf("找到配置文件: %s\n", location)
//	        // 读取并使用第一个找到的配置文件
//	        break
//	    }
//	}
func GetDefaultConfigLocations() []string {
	var locations []string

	// 当前目录下的配置文件
	locations = append(locations, DefaultNuGetConfigFilename)

	// 当前目录上层的配置文件
	locations = append(locations, filepath.Join("..", DefaultNuGetConfigFilename))

	// 用户级别的配置文件
	userConfigDir := getUserConfigDirectory()
	if userConfigDir != "" {
		userConfigPath := filepath.Join(userConfigDir, GlobalFolderName, DefaultNuGetConfigFilename)
		locations = append(locations, userConfigPath)
	}

	// 系统级别的配置文件
	systemConfigDir := getSystemConfigDirectory()
	if systemConfigDir != "" {
		systemConfigPath := filepath.Join(systemConfigDir, GlobalFolderName, DefaultNuGetConfigFilename)
		locations = append(locations, systemConfigPath)
	}

	return locations
}

// getUserConfigDirectory 获取用户配置目录
func getUserConfigDirectory() string {
	switch runtime.GOOS {
	case "windows":
		// Windows：%APPDATA%
		return os.Getenv("APPDATA")
	case "darwin":
		// macOS：~/Library/Application Support/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(homeDir, "Library", "Application Support")
	default:
		// Linux/Unix：~/.config/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(homeDir, ".config")
		}
		return configDir
	}
}

// getSystemConfigDirectory 获取系统配置目录
func getSystemConfigDirectory() string {
	switch runtime.GOOS {
	case "windows":
		// Windows：%ProgramData%
		return os.Getenv("ProgramData")
	case "darwin":
		// macOS：/Library/Application Support/
		return "/Library/Application Support"
	default:
		// Linux/Unix：/etc/
		return "/etc"
	}
}
