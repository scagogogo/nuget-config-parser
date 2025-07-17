// Package editor 实现基于Parser位置信息的NuGet配置文件编辑功能
package editor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/scagogogo/nuget-config-parser/pkg/parser"
	"github.com/scagogogo/nuget-config-parser/pkg/types"
)

// ConfigEditor 基于Parser位置信息的配置编辑器
type ConfigEditor struct {
	parseResult *parser.ParseResult
	edits       []Edit
}

// Edit 表示一个文本编辑操作
type Edit struct {
	Range   parser.Range // 要替换的范围
	NewText string       // 新文本
	Type    string       // 编辑类型：add, update, delete
}

// NewConfigEditor 创建基于Parser的配置编辑器
func NewConfigEditor(parseResult *parser.ParseResult) *ConfigEditor {
	return &ConfigEditor{
		parseResult: parseResult,
		edits:       make([]Edit, 0),
	}
}

// GetConfig 获取配置对象
func (e *ConfigEditor) GetConfig() *types.NuGetConfig {
	return e.parseResult.Config
}

// GetPositions 获取位置信息
func (e *ConfigEditor) GetPositions() map[string]*parser.ElementPosition {
	return e.parseResult.Positions
}

// AddPackageSource 添加新的包源
func (e *ConfigEditor) AddPackageSource(key, value, protocolVersion string) error {
	// 查找packageSources元素
	packageSourcesPath := "configuration/packageSources"
	elemPos, exists := e.parseResult.Positions[packageSourcesPath]
	if !exists {
		return fmt.Errorf("未找到packageSources元素")
	}

	// 构建新的包源XML
	indent := "    " // 4个空格缩进
	newSourceXML := fmt.Sprintf("\n%s<add key=\"%s\" value=\"%s\"", indent, key, value)
	if protocolVersion != "" {
		newSourceXML += fmt.Sprintf(" protocolVersion=\"%s\"", protocolVersion)
	}
	newSourceXML += " />"

	// 在packageSources结束标签前插入
	insertPos := e.findInsertPositionBeforeEndTag(elemPos)

	edit := Edit{
		Range: parser.Range{
			Start: insertPos,
			End:   insertPos,
		},
		NewText: newSourceXML,
		Type:    "add",
	}

	e.edits = append(e.edits, edit)

	// 同时更新内存中的配置对象
	newSource := types.PackageSource{
		Key:             key,
		Value:           value,
		ProtocolVersion: protocolVersion,
	}
	e.parseResult.Config.PackageSources.Add = append(e.parseResult.Config.PackageSources.Add, newSource)

	return nil
}

// RemovePackageSource 删除包源
func (e *ConfigEditor) RemovePackageSource(sourceKey string) error {
	// 查找要删除的包源元素
	for path, elemPos := range e.parseResult.Positions {
		// 检查是否是packageSources下的add元素（包括带索引的）
		if (strings.Contains(path, "packageSources/add") ||
			strings.Contains(path, "packageSources/add[")) &&
			elemPos.TagName == "add" {
			if key, exists := elemPos.Attributes["key"]; exists && key == sourceKey {
				// 删除整个元素
				edit := Edit{
					Range:   elemPos.Range,
					NewText: "",
					Type:    "delete",
				}
				e.edits = append(e.edits, edit)

				// 同时更新内存中的配置对象
				e.removePackageSourceFromConfig(sourceKey)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到包源: %s", sourceKey)
}

// UpdatePackageSourceURL 更新包源的URL
func (e *ConfigEditor) UpdatePackageSourceURL(sourceKey, newURL string) error {
	return e.updatePackageSourceAttribute(sourceKey, "value", newURL)
}

// UpdatePackageSourceVersion 更新包源的协议版本
func (e *ConfigEditor) UpdatePackageSourceVersion(sourceKey, newVersion string) error {
	return e.updatePackageSourceAttribute(sourceKey, "protocolVersion", newVersion)
}

// updatePackageSourceAttribute 更新包源的属性
func (e *ConfigEditor) updatePackageSourceAttribute(sourceKey, attrName, newValue string) error {
	for path, elemPos := range e.parseResult.Positions {
		// 检查是否是packageSources下的add元素（包括带索引的）
		if (strings.Contains(path, "packageSources/add") ||
			strings.Contains(path, "packageSources/add[")) &&
			elemPos.TagName == "add" {
			if key, exists := elemPos.Attributes["key"]; exists && key == sourceKey {
				// 查找属性的位置并更新
				if attrRange, attrExists := elemPos.AttrRanges[attrName]; attrExists {
					edit := Edit{
						Range:   attrRange,
						NewText: newValue,
						Type:    "update",
					}
					e.edits = append(e.edits, edit)
				} else {
					// 属性不存在，需要添加
					return e.addAttributeToElement(elemPos, attrName, newValue)
				}

				// 更新内存中的配置对象
				e.updatePackageSourceInConfig(sourceKey, attrName, newValue)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到包源: %s", sourceKey)
}

// ApplyEdits 应用所有编辑操作，返回修改后的内容
func (e *ConfigEditor) ApplyEdits() ([]byte, error) {
	if len(e.edits) == 0 {
		return e.parseResult.Content, nil
	}

	// 按位置倒序排序，从后往前应用编辑，避免位置偏移问题
	sort.Slice(e.edits, func(i, j int) bool {
		return e.edits[i].Range.Start.Offset > e.edits[j].Range.Start.Offset
	})

	content := string(e.parseResult.Content)

	for _, edit := range e.edits {
		start := edit.Range.Start.Offset
		end := edit.Range.End.Offset

		if start < 0 || end > len(content) || start > end {
			return nil, fmt.Errorf("无效的编辑范围: start=%d, end=%d, content_len=%d", start, end, len(content))
		}

		// 应用编辑
		content = content[:start] + edit.NewText + content[end:]
	}

	return []byte(content), nil
}

// findInsertPositionBeforeEndTag 查找在结束标签前的插入位置
func (e *ConfigEditor) findInsertPositionBeforeEndTag(elemPos *parser.ElementPosition) parser.Position {
	// 查找结束标签的位置
	content := string(e.parseResult.Content)

	// 从元素开始位置向后查找结束标签
	startOffset := elemPos.Range.Start.Offset
	endOffset := elemPos.Range.End.Offset

	// 在内容中查找 </tagName> 的位置
	tagName := elemPos.TagName
	endTag := fmt.Sprintf("</%s>", tagName)

	// 从后往前查找结束标签
	for i := endOffset - 1; i >= startOffset; i-- {
		if i+len(endTag) <= len(content) && content[i:i+len(endTag)] == endTag {
			// 找到结束标签，在其前面插入
			return parser.Position{
				Line:   elemPos.Range.End.Line,
				Column: elemPos.Range.End.Column,
				Offset: i,
			}
		}
	}

	// 如果没找到结束标签，可能是自闭合标签，在结束位置前插入
	return parser.Position{
		Line:   elemPos.Range.End.Line,
		Column: elemPos.Range.End.Column,
		Offset: elemPos.Range.End.Offset - 1,
	}
}

// addAttributeToElement 向元素添加新属性
func (e *ConfigEditor) addAttributeToElement(elemPos *parser.ElementPosition, attrName, attrValue string) error {
	// 这里需要更复杂的逻辑来在正确位置插入属性
	// 简化实现：暂时返回错误
	return fmt.Errorf("添加新属性功能尚未实现")
}

// removePackageSourceFromConfig 从配置对象中移除包源
func (e *ConfigEditor) removePackageSourceFromConfig(sourceKey string) {
	sources := e.parseResult.Config.PackageSources.Add
	for i, source := range sources {
		if source.Key == sourceKey {
			e.parseResult.Config.PackageSources.Add = append(sources[:i], sources[i+1:]...)
			break
		}
	}
}

// updatePackageSourceInConfig 更新配置对象中的包源
func (e *ConfigEditor) updatePackageSourceInConfig(sourceKey, attrName, newValue string) {
	for i, source := range e.parseResult.Config.PackageSources.Add {
		if source.Key == sourceKey {
			switch attrName {
			case "value":
				e.parseResult.Config.PackageSources.Add[i].Value = newValue
			case "protocolVersion":
				e.parseResult.Config.PackageSources.Add[i].ProtocolVersion = newValue
			}
			break
		}
	}
}
