// Package parser 实现 NuGet 配置文件的解析功能
package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/scagogogo/nuget-config-parser/pkg/constants"
	"github.com/scagogogo/nuget-config-parser/pkg/errors"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
	"github.com/scagogogo/nuget-config-parser/pkg/utils"
)

// Position 表示XML中的位置信息
type Position struct {
	Line   int // 行号（从1开始）
	Column int // 列号（从1开始）
	Offset int // 字节偏移量（从0开始）
}

// Range 表示XML中的范围信息
type Range struct {
	Start Position // 开始位置
	End   Position // 结束位置
}

// ElementPosition 记录XML元素的位置信息
type ElementPosition struct {
	TagName    string            // 标签名
	Attributes map[string]string // 属性
	Range      Range             // 元素范围
	AttrRanges map[string]Range  // 属性值的范围
	Content    string            // 元素内容
	SelfClose  bool              // 是否自闭合标签
}

// ParseResult 解析结果，包含配置和位置信息
type ParseResult struct {
	Config    *types.NuGetConfig          // 解析后的配置
	Positions map[string]*ElementPosition // 元素位置信息，key为元素路径
	Content   []byte                      // 原始内容
}

// ConfigParser NuGet 配置文件解析器
type ConfigParser struct {
	// DefaultConfigSearchPaths 配置文件搜索路径
	DefaultConfigSearchPaths []string
	// TrackPositions 是否跟踪位置信息
	TrackPositions bool
}

// NewConfigParser 创建一个新的配置解析器
func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		DefaultConfigSearchPaths: constants.GetDefaultConfigLocations(),
		TrackPositions:           false,
	}
}

// NewPositionAwareParser 创建一个位置感知的配置解析器
func NewPositionAwareParser() *ConfigParser {
	return &ConfigParser{
		DefaultConfigSearchPaths: constants.GetDefaultConfigLocations(),
		TrackPositions:           true,
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

// ParseFromFileWithPositions 从文件解析配置并记录位置信息
func (p *ConfigParser) ParseFromFileWithPositions(filePath string) (*ParseResult, error) {
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

	return p.ParseFromContentWithPositions(data)
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

// ParseFromContentWithPositions 从内容解析配置并记录位置信息
func (p *ConfigParser) ParseFromContentWithPositions(content []byte) (*ParseResult, error) {
	// 验证内容是否为有效的XML
	if !utils.IsValidXML(string(content)) {
		return nil, errors.ErrInvalidConfigFormat
	}

	// 先进行标准解析
	var config types.NuGetConfig
	err := xml.Unmarshal(content, &config)
	if err != nil {
		return nil, errors.NewParseError(errors.ErrXMLParsing, 0, 0, fmt.Sprintf("xml.Unmarshal error: %v", err))
	}

	// 验证必需的字段
	if len(config.PackageSources.Add) == 0 {
		if !config.PackageSources.Clear {
			return nil, errors.NewParseError(errors.ErrMissingRequiredElement, 0, 0, "no package sources defined")
		}
	}

	// 跟踪位置信息
	positions, err := p.trackPositions(content)
	if err != nil {
		return nil, fmt.Errorf("failed to track positions: %w", err)
	}

	return &ParseResult{
		Config:    &config,
		Positions: positions,
		Content:   content,
	}, nil
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

// trackPositions 跟踪XML中所有元素的位置
func (p *ConfigParser) trackPositions(content []byte) (map[string]*ElementPosition, error) {
	positions := make(map[string]*ElementPosition)
	contentStr := string(content)

	var elementStack []string
	line := 1
	column := 1

	for i := 0; i < len(contentStr); i++ {
		if contentStr[i] == '\n' {
			line++
			column = 1
			continue
		}

		if contentStr[i] == '<' {
			// 查找标签结束
			tagEnd := p.findTagEnd(contentStr, i)
			if tagEnd == -1 {
				return nil, fmt.Errorf("未找到标签结束: 位置 %d", i)
			}

			tagContent := contentStr[i+1 : tagEnd]

			// 跳过注释和声明
			if strings.HasPrefix(tagContent, "!") || strings.HasPrefix(tagContent, "?") {
				p.advancePosition(contentStr, i, tagEnd+1, &line, &column)
				i = tagEnd
				continue
			}

			if strings.HasPrefix(tagContent, "/") {
				// 结束标签
				if len(elementStack) > 0 {
					elementPath := strings.Join(elementStack, "/")
					if elemPos, exists := positions[elementPath]; exists {
						elemPos.Range.End = Position{
							Line:   line,
							Column: column,
							Offset: tagEnd + 1,
						}
					}
					elementStack = elementStack[:len(elementStack)-1]
				}
			} else {
				// 开始标签
				startPos := Position{
					Line:   line,
					Column: column,
					Offset: i,
				}

				tagName, attributes, attrRanges, selfClose := p.parseTagWithRanges(tagContent, i+1)
				elementStack = append(elementStack, tagName)
				elementPath := strings.Join(elementStack, "/")

				// 为重复元素添加索引
				finalPath := elementPath
				if _, exists := positions[elementPath]; exists {
					// 如果路径已存在，添加索引
					index := 1
					for {
						indexedPath := fmt.Sprintf("%s[%d]", elementPath, index)
						if _, exists := positions[indexedPath]; !exists {
							finalPath = indexedPath
							break
						}
						index++
					}
				}

				positions[finalPath] = &ElementPosition{
					TagName:    tagName,
					Attributes: attributes,
					AttrRanges: attrRanges,
					Range: Range{
						Start: startPos,
						End:   Position{Line: line, Column: column, Offset: tagEnd + 1},
					},
					SelfClose: selfClose,
				}

				if selfClose {
					elementStack = elementStack[:len(elementStack)-1]
				}
			}

			p.advancePosition(contentStr, i, tagEnd+1, &line, &column)
			i = tagEnd
		} else {
			column++
		}
	}

	return positions, nil
}

// findTagEnd 查找标签结束位置
func (p *ConfigParser) findTagEnd(content string, start int) int {
	for i := start + 1; i < len(content); i++ {
		if content[i] == '>' {
			return i
		}
	}
	return -1
}

// advancePosition 更新位置信息
func (p *ConfigParser) advancePosition(content string, start, end int, line, column *int) {
	for i := start; i < end && i < len(content); i++ {
		if content[i] == '\n' {
			*line++
			*column = 1
		} else {
			*column++
		}
	}
}

// parseTagWithRanges 解析标签内容并记录属性范围
func (p *ConfigParser) parseTagWithRanges(tagContent string, baseOffset int) (string, map[string]string, map[string]Range, bool) {
	tagContent = strings.TrimSpace(tagContent)
	selfClose := strings.HasSuffix(tagContent, "/")
	if selfClose {
		tagContent = strings.TrimSuffix(tagContent, "/")
		tagContent = strings.TrimSpace(tagContent)
	}

	parts := strings.Fields(tagContent)
	if len(parts) == 0 {
		return "", make(map[string]string), make(map[string]Range), selfClose
	}

	tagName := parts[0]
	attributes := make(map[string]string)
	attrRanges := make(map[string]Range)

	// 解析属性并记录位置
	attrStr := strings.Join(parts[1:], " ")
	if attrStr != "" {
		// 查找属性值的位置
		re := regexp.MustCompile(`(\w+)="([^"]*)"`)
		matches := re.FindAllStringSubmatchIndex(attrStr, -1)
		for _, match := range matches {
			if len(match) >= 6 {
				attrName := attrStr[match[2]:match[3]]
				attrValue := attrStr[match[4]:match[5]]
				attributes[attrName] = attrValue

				// 记录属性值的范围（不包括引号）
				valueStart := baseOffset + strings.Index(tagContent, attrStr) + match[4]
				valueEnd := baseOffset + strings.Index(tagContent, attrStr) + match[5]
				attrRanges[attrName] = Range{
					Start: Position{Offset: valueStart},
					End:   Position{Offset: valueEnd},
				}
			}
		}
	}

	return tagName, attributes, attrRanges, selfClose
}
