// Package manager 提供管理NuGet配置文件的功能
package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	pkgErrors "github.com/scagogogo/nuget-config-parser/pkg/errors"
	"github.com/scagogogo/nuget-config-parser/pkg/finder"
	"github.com/scagogogo/nuget-config-parser/pkg/parser"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

// ConfigManager NuGet配置管理器
type ConfigManager struct {
	parser *parser.ConfigParser
	finder *finder.ConfigFinder
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		parser: parser.NewConfigParser(),
		finder: finder.NewConfigFinder(),
	}
}

// LoadConfig 加载配置文件
func (m *ConfigManager) LoadConfig(filePath string) (*types.NuGetConfig, error) {
	return m.parser.ParseFromFile(filePath)
}

// FindAndLoadConfig 查找并加载第一个可用的配置文件
func (m *ConfigManager) FindAndLoadConfig() (*types.NuGetConfig, string, error) {
	configPath, err := m.finder.FindConfigFile()
	if err != nil {
		return nil, "", pkgErrors.ErrConfigFileNotFound
	}

	config, err := m.LoadConfig(configPath)
	if err != nil {
		return nil, configPath, err
	}

	return config, configPath, nil
}

// SaveConfig 保存配置到文件
func (m *ConfigManager) SaveConfig(config *types.NuGetConfig, filePath string) error {
	return m.parser.SaveToFile(config, filePath)
}

// CreateDefaultConfig 创建默认配置
func (m *ConfigManager) CreateDefaultConfig() *types.NuGetConfig {
	// 创建包含默认源的配置
	defaultSource := types.PackageSource{
		Key:             "nuget.org",
		Value:           constants.DefaultPackageSource,
		ProtocolVersion: constants.NuGetV3APIProtocolVersion,
	}

	return &types.NuGetConfig{
		PackageSources: types.PackageSources{
			Add: []types.PackageSource{defaultSource},
		},
		ActivePackageSource: &types.ActivePackageSource{
			Add: defaultSource,
		},
	}
}

// InitializeDefaultConfig 在指定路径创建默认配置
func (m *ConfigManager) InitializeDefaultConfig(filePath string) error {
	// 检查文件目录是否存在，不存在则创建
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建默认配置
	config := m.CreateDefaultConfig()

	// 保存配置
	return m.SaveConfig(config, filePath)
}

// AddPackageSource 添加包源
func (m *ConfigManager) AddPackageSource(config *types.NuGetConfig, key string, value string, protocolVersion string) {
	// 检查是否已存在相同键的包源
	for i, source := range config.PackageSources.Add {
		if source.Key == key {
			// 更新现有包源
			config.PackageSources.Add[i].Value = value
			if protocolVersion != "" {
				config.PackageSources.Add[i].ProtocolVersion = protocolVersion
			}
			return
		}
	}

	// 添加新包源
	newSource := types.PackageSource{
		Key:   key,
		Value: value,
	}

	if protocolVersion != "" {
		newSource.ProtocolVersion = protocolVersion
	}

	config.PackageSources.Add = append(config.PackageSources.Add, newSource)
}

// RemovePackageSource 移除包源
func (m *ConfigManager) RemovePackageSource(config *types.NuGetConfig, key string) bool {
	for i, source := range config.PackageSources.Add {
		if source.Key == key {
			// 移除指定的包源
			config.PackageSources.Add = append(config.PackageSources.Add[:i], config.PackageSources.Add[i+1:]...)
			return true
		}
	}
	return false
}

// GetPackageSource 获取指定键的包源
func (m *ConfigManager) GetPackageSource(config *types.NuGetConfig, key string) *types.PackageSource {
	for _, source := range config.PackageSources.Add {
		if source.Key == key {
			return &source
		}
	}
	return nil
}

// GetAllPackageSources 获取所有包源
func (m *ConfigManager) GetAllPackageSources(config *types.NuGetConfig) []types.PackageSource {
	return config.PackageSources.Add
}

// SetActivePackageSource 设置活跃包源
func (m *ConfigManager) SetActivePackageSource(config *types.NuGetConfig, key string) error {
	// 查找包源
	var source *types.PackageSource
	for _, s := range config.PackageSources.Add {
		if s.Key == key {
			source = &s
			break
		}
	}

	if source == nil {
		return fmt.Errorf("package source with key '%s' not found", key)
	}

	// 如果 ActivePackageSource 为 nil，则初始化
	if config.ActivePackageSource == nil {
		config.ActivePackageSource = &types.ActivePackageSource{}
	}

	// 设置活跃包源
	config.ActivePackageSource.Add = *source
	return nil
}

// AddCredential 添加包源凭证
func (m *ConfigManager) AddCredential(config *types.NuGetConfig, sourceKey string, username string, password string) {
	// 如果 PackageSourceCredentials 为 nil，则初始化
	if config.PackageSourceCredentials == nil {
		config.PackageSourceCredentials = &types.PackageSourceCredentials{
			Sources: make(map[string]types.SourceCredential),
		}
	}

	// 创建或更新凭证
	var credentials []types.Credential

	// 添加用户名
	credentials = append(credentials, types.Credential{
		Key:   "Username",
		Value: username,
	})

	// 添加密码
	credentials = append(credentials, types.Credential{
		Key:   "ClearTextPassword",
		Value: password,
	})

	// 设置凭证
	sourceCredential := types.SourceCredential{
		Add: credentials,
	}

	config.PackageSourceCredentials.Sources[sourceKey] = sourceCredential
}

// RemoveCredential 移除包源凭证
func (m *ConfigManager) RemoveCredential(config *types.NuGetConfig, sourceKey string) bool {
	if config.PackageSourceCredentials == nil || len(config.PackageSourceCredentials.Sources) == 0 {
		return false
	}

	if _, exists := config.PackageSourceCredentials.Sources[sourceKey]; !exists {
		return false
	}

	delete(config.PackageSourceCredentials.Sources, sourceKey)
	return true
}

// DisablePackageSource 禁用包源
func (m *ConfigManager) DisablePackageSource(config *types.NuGetConfig, key string) {
	// 如果 DisabledPackageSources 为 nil，则初始化
	if config.DisabledPackageSources == nil {
		config.DisabledPackageSources = &types.DisabledPackageSources{
			Add: []types.DisabledSource{},
		}
	}

	// 检查是否已经禁用
	for i, source := range config.DisabledPackageSources.Add {
		if source.Key == key {
			// 更新为禁用状态
			config.DisabledPackageSources.Add[i].Value = "true"
			return
		}
	}

	// 添加新的禁用源
	config.DisabledPackageSources.Add = append(config.DisabledPackageSources.Add, types.DisabledSource{
		Key:   key,
		Value: "true",
	})
}

// EnablePackageSource 启用包源
func (m *ConfigManager) EnablePackageSource(config *types.NuGetConfig, key string) bool {
	if config.DisabledPackageSources == nil {
		return false
	}

	for i, source := range config.DisabledPackageSources.Add {
		if source.Key == key {
			// 从禁用列表中移除
			config.DisabledPackageSources.Add = append(
				config.DisabledPackageSources.Add[:i],
				config.DisabledPackageSources.Add[i+1:]...)
			return true
		}
	}

	return false
}

// IsPackageSourceDisabled 检查包源是否被禁用
func (m *ConfigManager) IsPackageSourceDisabled(config *types.NuGetConfig, key string) bool {
	if config.DisabledPackageSources == nil {
		return false
	}

	for _, source := range config.DisabledPackageSources.Add {
		if source.Key == key && source.Value == "true" {
			return true
		}
	}

	return false
}

// AddConfigOption 添加配置选项
func (m *ConfigManager) AddConfigOption(config *types.NuGetConfig, key string, value string) {
	// 如果 Config 为 nil，则初始化
	if config.Config == nil {
		config.Config = &types.Config{
			Add: []types.ConfigOption{},
		}
	}

	// 检查是否已存在相同键的配置选项
	for i, option := range config.Config.Add {
		if option.Key == key {
			// 更新现有选项
			config.Config.Add[i].Value = value
			return
		}
	}

	// 添加新选项
	config.Config.Add = append(config.Config.Add, types.ConfigOption{
		Key:   key,
		Value: value,
	})
}

// RemoveConfigOption 移除配置选项
func (m *ConfigManager) RemoveConfigOption(config *types.NuGetConfig, key string) bool {
	if config.Config == nil {
		return false
	}

	for i, option := range config.Config.Add {
		if option.Key == key {
			// 移除指定的配置选项
			config.Config.Add = append(config.Config.Add[:i], config.Config.Add[i+1:]...)
			return true
		}
	}

	return false
}

// GetConfigOption 获取配置选项值
func (m *ConfigManager) GetConfigOption(config *types.NuGetConfig, key string) string {
	if config.Config == nil {
		return ""
	}

	for _, option := range config.Config.Add {
		if option.Key == key {
			return option.Value
		}
	}

	return ""
}
