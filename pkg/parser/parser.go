// Package parser 实现 NuGet 配置文件的解析功能
package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	"github.com/scagogogo/nuget-config-parser/pkg/errors"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
	"github.com/scagogogo/nuget-config-parser/pkg/utils"
)

// ConfigParser NuGet 配置文件解析器
type ConfigParser struct {
	// DefaultConfigSearchPaths 配置文件搜索路径
	DefaultConfigSearchPaths []string
}

// NewConfigParser 创建一个新的配置解析器
func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		DefaultConfigSearchPaths: constants.GetDefaultConfigLocations(),
	}
}

// ParseFromFile 从文件解析配置
func (p *ConfigParser) ParseFromFile(filePath string) (*types.NuGetConfig, error) {
	// 检查文件是否存在
	if !utils.FileExists(filePath) {
		return nil, errors.ErrConfigFileNotFound
	}

	// 读取文件内容
	data, err := utils.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.ErrEmptyConfigFile
	}

	return p.ParseFromContent(data)
}

// ParseFromContent 从内容解析配置
func (p *ConfigParser) ParseFromContent(content []byte) (*types.NuGetConfig, error) {
	// 验证内容是否为有效的XML
	if !utils.IsValidXML(string(content)) {
		return nil, errors.ErrInvalidConfigFormat
	}

	// 解析XML
	var config types.NuGetConfig
	err := xml.Unmarshal(content, &config)
	if err != nil {
		return nil, errors.NewParseError(errors.ErrXMLParsing, 0, 0, fmt.Sprintf("xml.Unmarshal error: %v", err))
	}

	// 验证必需的字段
	if len(config.PackageSources.Add) == 0 {
		// 如果没有定义包源但有 clear 属性为 true，这可能是正常的情况
		if !config.PackageSources.Clear {
			return nil, errors.NewParseError(errors.ErrMissingRequiredElement, 0, 0, "no package sources defined")
		}
	}

	return &config, nil
}

// FindAndParseConfig 查找并解析配置文件
func (p *ConfigParser) FindAndParseConfig() (*types.NuGetConfig, string, error) {
	// 尝试所有默认路径
	for _, path := range p.DefaultConfigSearchPaths {
		expandedPath := utils.ExpandEnvVars(path)
		absPath, err := filepath.Abs(expandedPath)
		if err != nil {
			continue
		}

		if utils.FileExists(absPath) {
			config, err := p.ParseFromFile(absPath)
			if err == nil {
				return config, absPath, nil
			}
		}
	}

	return nil, "", errors.ErrConfigFileNotFound
}

// ParseFromReader 从io.Reader解析配置
func (p *ConfigParser) ParseFromReader(reader io.Reader) (*types.NuGetConfig, error) {
	// 读取内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from reader: %w", err)
	}

	return p.ParseFromContent(content)
}

// ParseFromString 从字符串解析配置
func (p *ConfigParser) ParseFromString(content string) (*types.NuGetConfig, error) {
	return p.ParseFromContent([]byte(content))
}

// SerializeToXML 将配置序列化为XML字符串
func (p *ConfigParser) SerializeToXML(config *types.NuGetConfig) (string, error) {
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config to XML: %w", err)
	}

	xmlHeader := `<?xml version="1.0" encoding="utf-8"?>` + "\n"
	return xmlHeader + string(data), nil
}

// SaveToFile 将配置保存到文件
func (p *ConfigParser) SaveToFile(config *types.NuGetConfig, filePath string) error {
	xmlString, err := p.SerializeToXML(config)
	if err != nil {
		return err
	}

	return utils.WriteToFile(filePath, []byte(xmlString))
}
