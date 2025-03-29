package types

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestNuGetConfigXMLMarshaling(t *testing.T) {
	// 创建一个包含所有字段的配置
	config := &NuGetConfig{
		PackageSources: PackageSources{
			Clear: true,
			Add: []PackageSource{
				{
					Key:             "nuget.org",
					Value:           "https://api.nuget.org/v3/index.json",
					ProtocolVersion: "3",
				},
				{
					Key:   "local",
					Value: "C:\\Packages",
				},
			},
		},
		PackageSourceCredentials: &PackageSourceCredentials{
			Sources: map[string]SourceCredential{
				"nuget.org": {
					Add: []Credential{
						{
							Key:   "Username",
							Value: "testuser",
						},
						{
							Key:   "ClearTextPassword",
							Value: "testpass",
						},
					},
				},
			},
		},
		DisabledPackageSources: &DisabledPackageSources{
			Add: []DisabledSource{
				{
					Key:   "local",
					Value: "true",
				},
			},
		},
		ActivePackageSource: &ActivePackageSource{
			Add: PackageSource{
				Key:   "nuget.org",
				Value: "https://api.nuget.org/v3/index.json",
			},
		},
		Config: &Config{
			Add: []ConfigOption{
				{
					Key:   "globalPackagesFolder",
					Value: "%USERPROFILE%\\.nuget\\packages",
				},
			},
		},
	}

	// 序列化为 XML
	xmlData, err := xml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config to XML: %v", err)
	}

	// 反序列化为新对象
	var newConfig NuGetConfig
	err = xml.Unmarshal(xmlData, &newConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML to config: %v", err)
	}

	// 验证字段是否正确序列化和反序列化
	// 检查 PackageSources
	if newConfig.PackageSources.Clear != config.PackageSources.Clear {
		t.Errorf("PackageSources.Clear = %v, want %v", newConfig.PackageSources.Clear, config.PackageSources.Clear)
	}

	if len(newConfig.PackageSources.Add) != len(config.PackageSources.Add) {
		t.Errorf("len(PackageSources.Add) = %v, want %v", len(newConfig.PackageSources.Add), len(config.PackageSources.Add))
	} else {
		for i, source := range config.PackageSources.Add {
			if source.Key != newConfig.PackageSources.Add[i].Key {
				t.Errorf("PackageSources.Add[%d].Key = %v, want %v", i, newConfig.PackageSources.Add[i].Key, source.Key)
			}
			if source.Value != newConfig.PackageSources.Add[i].Value {
				t.Errorf("PackageSources.Add[%d].Value = %v, want %v", i, newConfig.PackageSources.Add[i].Value, source.Value)
			}
			if source.ProtocolVersion != newConfig.PackageSources.Add[i].ProtocolVersion {
				t.Errorf("PackageSources.Add[%d].ProtocolVersion = %v, want %v", i, newConfig.PackageSources.Add[i].ProtocolVersion, source.ProtocolVersion)
			}
		}
	}

	// 检查 ActivePackageSource
	if newConfig.ActivePackageSource == nil {
		t.Errorf("ActivePackageSource should not be nil")
	} else {
		if newConfig.ActivePackageSource.Add.Key != config.ActivePackageSource.Add.Key {
			t.Errorf("ActivePackageSource.Add.Key = %v, want %v", newConfig.ActivePackageSource.Add.Key, config.ActivePackageSource.Add.Key)
		}
		if newConfig.ActivePackageSource.Add.Value != config.ActivePackageSource.Add.Value {
			t.Errorf("ActivePackageSource.Add.Value = %v, want %v", newConfig.ActivePackageSource.Add.Value, config.ActivePackageSource.Add.Value)
		}
	}

	// 检查 DisabledPackageSources
	if newConfig.DisabledPackageSources == nil {
		t.Errorf("DisabledPackageSources should not be nil")
	} else if len(newConfig.DisabledPackageSources.Add) != len(config.DisabledPackageSources.Add) {
		t.Errorf("len(DisabledPackageSources.Add) = %v, want %v", len(newConfig.DisabledPackageSources.Add), len(config.DisabledPackageSources.Add))
	} else {
		for i, source := range config.DisabledPackageSources.Add {
			if source.Key != newConfig.DisabledPackageSources.Add[i].Key {
				t.Errorf("DisabledPackageSources.Add[%d].Key = %v, want %v", i, newConfig.DisabledPackageSources.Add[i].Key, source.Key)
			}
			if source.Value != newConfig.DisabledPackageSources.Add[i].Value {
				t.Errorf("DisabledPackageSources.Add[%d].Value = %v, want %v", i, newConfig.DisabledPackageSources.Add[i].Value, source.Value)
			}
		}
	}

	// 检查 Config
	if newConfig.Config == nil {
		t.Errorf("Config should not be nil")
	} else if len(newConfig.Config.Add) != len(config.Config.Add) {
		t.Errorf("len(Config.Add) = %v, want %v", len(newConfig.Config.Add), len(config.Config.Add))
	} else {
		for i, option := range config.Config.Add {
			if option.Key != newConfig.Config.Add[i].Key {
				t.Errorf("Config.Add[%d].Key = %v, want %v", i, newConfig.Config.Add[i].Key, option.Key)
			}
			if option.Value != newConfig.Config.Add[i].Value {
				t.Errorf("Config.Add[%d].Value = %v, want %v", i, newConfig.Config.Add[i].Value, option.Value)
			}
		}
	}

	// PackageSourceCredentials 需要特殊处理，因为 XML 序列化方式不同于结构体的表示方式
	// 这里我们只能测试反序列化后能否正确访问到凭证
	if newConfig.PackageSourceCredentials == nil {
		t.Errorf("PackageSourceCredentials should not be nil")
	} else {
		// 检查凭证数量
		if len(newConfig.PackageSourceCredentials.Sources) != len(config.PackageSourceCredentials.Sources) {
			t.Errorf("len(PackageSourceCredentials.Sources) = %v, want %v",
				len(newConfig.PackageSourceCredentials.Sources),
				len(config.PackageSourceCredentials.Sources))
		}

		// 检查特定源的凭证
		source, exists := newConfig.PackageSourceCredentials.Sources["nuget.org"]
		if !exists {
			t.Errorf("PackageSourceCredentials.Sources[\"nuget.org\"] does not exist")
		} else {
			// 检查凭证项数量
			originalSource := config.PackageSourceCredentials.Sources["nuget.org"]
			if len(source.Add) != len(originalSource.Add) {
				t.Errorf("len(PackageSourceCredentials.Sources[\"nuget.org\"].Add) = %v, want %v",
					len(source.Add),
					len(originalSource.Add))
			}

			// 创建映射以便于比较
			credMap := make(map[string]string)
			for _, cred := range source.Add {
				credMap[cred.Key] = cred.Value
			}

			// 验证每个凭证
			for _, origCred := range originalSource.Add {
				if val, exists := credMap[origCred.Key]; !exists {
					t.Errorf("Credential with key %s not found", origCred.Key)
				} else if val != origCred.Value {
					t.Errorf("Credential %s value = %s, want %s", origCred.Key, val, origCred.Value)
				}
			}
		}
	}
}

func TestStructTagsXML(t *testing.T) {
	// 检查 NuGetConfig 结构体字段的 XML 标签
	t.Run("NuGetConfig", func(t *testing.T) {
		typ := reflect.TypeOf(NuGetConfig{})

		checkFieldXMLTag(t, typ, "PackageSources", "packageSources")
		checkFieldXMLTag(t, typ, "PackageSourceCredentials", "packageSourceCredentials,omitempty")
		checkFieldXMLTag(t, typ, "Config", "config,omitempty")
		checkFieldXMLTag(t, typ, "DisabledPackageSources", "disabledPackageSources,omitempty")
		checkFieldXMLTag(t, typ, "ActivePackageSource", "activePackageSource,omitempty")
	})

	// 检查 PackageSources 结构体字段的 XML 标签
	t.Run("PackageSources", func(t *testing.T) {
		typ := reflect.TypeOf(PackageSources{})

		checkFieldXMLTag(t, typ, "Clear", "clear,attr,omitempty")
		checkFieldXMLTag(t, typ, "Add", "add")
	})

	// 检查 PackageSource 结构体字段的 XML 标签
	t.Run("PackageSource", func(t *testing.T) {
		typ := reflect.TypeOf(PackageSource{})

		checkFieldXMLTag(t, typ, "Key", "key,attr")
		checkFieldXMLTag(t, typ, "Value", "value,attr")
		checkFieldXMLTag(t, typ, "ProtocolVersion", "protocolVersion,attr,omitempty")
	})

	// 其他结构体的检查可以类似添加...
}

// 检查结构体字段的 XML 标签
func checkFieldXMLTag(t *testing.T, typ reflect.Type, fieldName, expectedTag string) {
	field, found := typ.FieldByName(fieldName)
	if !found {
		t.Errorf("Field %s not found in %s", fieldName, typ.Name())
		return
	}

	tag := field.Tag.Get("xml")
	if tag != expectedTag {
		t.Errorf("Field %s in %s has XML tag %q, want %q", fieldName, typ.Name(), tag, expectedTag)
	}
}
