// Package types 定义 NuGet 配置文件的数据模型结构
package types

import (
	"encoding/xml"
	"io"
)

// NuGetConfig 表示一个完整的 NuGet 配置文件
type NuGetConfig struct {
	// PackageSources 定义可用的包源
	PackageSources PackageSources `xml:"packageSources"`

	// PackageSourceCredentials 定义包源凭证信息
	PackageSourceCredentials *PackageSourceCredentials `xml:"packageSourceCredentials,omitempty"`

	// Config 定义全局配置选项
	Config *Config `xml:"config,omitempty"`

	// DisabledPackageSources 定义被禁用的包源
	DisabledPackageSources *DisabledPackageSources `xml:"disabledPackageSources,omitempty"`

	// ActivePackageSource 定义当前活跃的包源
	ActivePackageSource *ActivePackageSource `xml:"activePackageSource,omitempty"`
}

// PackageSources 定义包源列表
type PackageSources struct {
	// Clear 如果存在并且为 true，则清除之前的所有包源
	Clear bool `xml:"clear,attr,omitempty"`

	// Add 表示添加的包源列表
	Add []PackageSource `xml:"add"`
}

// PackageSource 定义单个包源
type PackageSource struct {
	// Key 包源的唯一标识符
	Key string `xml:"key,attr"`

	// Value 包源的 URL 或路径
	Value string `xml:"value,attr"`

	// ProtocolVersion 包源使用的协议版本
	ProtocolVersion string `xml:"protocolVersion,attr,omitempty"`
}

// PackageSourceCredentials 定义包源凭证
type PackageSourceCredentials struct {
	// 键为包源名称，值为该包源的凭证
	Sources map[string]SourceCredential `xml:"-"` // 不直接序列化
}

// MarshalXML 自定义PackageSourceCredentials的XML序列化
func (p *PackageSourceCredentials) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if p == nil || len(p.Sources) == 0 {
		return nil
	}

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	// 遍历所有凭证源
	for key, cred := range p.Sources {
		// 为每个凭证源创建一个元素
		sourceElem := xml.StartElement{Name: xml.Name{Local: key}}
		if err := e.EncodeToken(sourceElem); err != nil {
			return err
		}

		// 编码凭证项
		for _, add := range cred.Add {
			addElem := xml.StartElement{Name: xml.Name{Local: "add"}}
			addElem.Attr = append(addElem.Attr,
				xml.Attr{Name: xml.Name{Local: "key"}, Value: add.Key},
				xml.Attr{Name: xml.Name{Local: "value"}, Value: add.Value},
			)
			if err := e.EncodeToken(addElem); err != nil {
				return err
			}
			if err := e.EncodeToken(xml.EndElement{Name: addElem.Name}); err != nil {
				return err
			}
		}

		if err := e.EncodeToken(xml.EndElement{Name: sourceElem.Name}); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML 自定义PackageSourceCredentials的XML反序列化
func (p *PackageSourceCredentials) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// 初始化Sources映射
	p.Sources = make(map[string]SourceCredential)

	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			if tt.Name.Local != "add" {
				// 这是一个包源名称元素
				sourceName := tt.Name.Local
				var sourceCred SourceCredential

				// 解析这个源的所有凭证
				if err := d.DecodeElement(&sourceCred, &tt); err != nil {
					return err
				}

				p.Sources[sourceName] = sourceCred
			}
		case xml.EndElement:
			if tt.Name == start.Name {
				return nil
			}
		}
	}

	return nil
}

// SourceCredential 定义源凭证信息
type SourceCredential struct {
	// Add 凭证列表，通常包含用户名和密码
	Add []Credential `xml:"add"`
}

// Credential 定义单个凭证键值对
type Credential struct {
	// Key 凭证键名，如 Username、Password、ClearTextPassword 等
	Key string `xml:"key,attr"`

	// Value 凭证值
	Value string `xml:"value,attr"`
}

// DisabledPackageSources 定义被禁用的包源
type DisabledPackageSources struct {
	// Add 表示禁用的包源列表
	Add []DisabledSource `xml:"add"`
}

// DisabledSource 定义被禁用的单个包源
type DisabledSource struct {
	// Key 包源的标识符
	Key string `xml:"key,attr"`

	// Value 通常为 "true"，表示该源被禁用
	Value string `xml:"value,attr"`
}

// ActivePackageSource 定义当前使用的包源
type ActivePackageSource struct {
	// Add 表示活跃的包源，通常只有一个
	Add PackageSource `xml:"add"`
}

// Config 定义全局配置选项
type Config struct {
	// Add 配置选项列表
	Add []ConfigOption `xml:"add"`
}

// ConfigOption 定义配置选项
type ConfigOption struct {
	// Key 配置键名
	Key string `xml:"key,attr"`

	// Value 配置值
	Value string `xml:"value,attr"`
}
