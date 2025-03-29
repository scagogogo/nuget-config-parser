// Package finder 提供查找NuGet配置文件的功能
package finder

import (
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	"github.com/scagogogo/nuget-config-parser/pkg/utils"
)

// ConfigFinder NuGet配置文件查找器
type ConfigFinder struct {
	// 环境变量名，用于自定义配置文件位置
	EnvVariableName string
}

// NewConfigFinder 创建新的配置文件查找器
func NewConfigFinder() *ConfigFinder {
	return &ConfigFinder{
		EnvVariableName: "NUGET_CONFIG_FILE",
	}
}

// GetConfigFileSearchLocations 获取可能的配置文件位置列表
func (f *ConfigFinder) GetConfigFileSearchLocations() []string {
	var locations []string

	// 1. 先检查环境变量
	envPath := os.Getenv(f.EnvVariableName)
	if envPath != "" && utils.FileExists(envPath) {
		locations = append(locations, envPath)
	}

	// 2. 添加默认搜索位置
	locations = append(locations, constants.GetDefaultConfigLocations()...)

	return locations
}

// FindConfigFile 寻找第一个存在的配置文件
func (f *ConfigFinder) FindConfigFile() (string, error) {
	locations := f.GetConfigFileSearchLocations()

	for _, location := range locations {
		expandedPath := utils.ExpandEnvVars(location)
		absPath, err := filepath.Abs(expandedPath)
		if err != nil {
			continue
		}

		if utils.FileExists(absPath) {
			return absPath, nil
		}
	}

	return "", os.ErrNotExist
}

// FindAllConfigFiles 找到所有存在的配置文件
func (f *ConfigFinder) FindAllConfigFiles() []string {
	locations := f.GetConfigFileSearchLocations()
	var existingFiles []string

	for _, location := range locations {
		expandedPath := utils.ExpandEnvVars(location)
		absPath, err := filepath.Abs(expandedPath)
		if err != nil {
			continue
		}

		if utils.FileExists(absPath) {
			existingFiles = append(existingFiles, absPath)
		}
	}

	return existingFiles
}

// FindProjectConfig 在指定目录及其父目录中查找项目级配置文件
func (f *ConfigFinder) FindProjectConfig(startDir string) (string, error) {
	currentDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		configPath := filepath.Join(currentDir, constants.DefaultNuGetConfigFilename)
		if utils.FileExists(configPath) {
			return configPath, nil
		}

		// 获取父目录
		parentDir := filepath.Dir(currentDir)
		// 如果已到达根目录，则停止搜索
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", os.ErrNotExist
}

// GetUserConfigFile 获取用户级别的配置文件路径
func (f *ConfigFinder) GetUserConfigFile() string {
	userConfigDir := getUserConfigDirectory()
	if userConfigDir == "" {
		return ""
	}

	return filepath.Join(userConfigDir, constants.GlobalFolderName, constants.DefaultNuGetConfigFilename)
}

// GetMachineConfigFile 获取机器级别的配置文件路径
func (f *ConfigFinder) GetMachineConfigFile() string {
	systemConfigDir := getSystemConfigDirectory()
	if systemConfigDir == "" {
		return ""
	}

	return filepath.Join(systemConfigDir, constants.GlobalFolderName, constants.DefaultNuGetConfigFilename)
}

// getUserConfigDirectory 获取用户配置目录
func getUserConfigDirectory() string {
	switch os.Getenv("GOOS") {
	case "windows":
		return os.Getenv("APPDATA")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(homeDir, "Library", "Application Support")
	default:
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
	switch os.Getenv("GOOS") {
	case "windows":
		return os.Getenv("ProgramData")
	case "darwin":
		return "/Library/Application Support"
	default:
		return "/etc"
	}
}
