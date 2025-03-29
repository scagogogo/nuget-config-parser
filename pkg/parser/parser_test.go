package parser

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scagogogo/nuget-config-parser/pkg/errors"
	nugetTesting "github.com/scagogogo/nuget-config-parser/pkg/testing"
)

func TestNewConfigParser(t *testing.T) {
	parser := NewConfigParser()
	if parser == nil {
		t.Fatal("NewConfigParser() returned nil")
	}

	if len(parser.DefaultConfigSearchPaths) == 0 {
		t.Error("DefaultConfigSearchPaths should not be empty")
	}
}

func TestParseFromFile(t *testing.T) {
	// 创建临时文件
	validContent := nugetTesting.ValidNuGetConfig()
	validFile := nugetTesting.CreateTempFile(t, validContent)
	defer os.Remove(validFile)

	invalidContent := nugetTesting.InvalidNuGetConfig()
	invalidFile := nugetTesting.CreateTempFile(t, invalidContent)
	defer os.Remove(invalidFile)

	emptyFile := nugetTesting.CreateTempFile(t, "")
	defer os.Remove(emptyFile)

	nonExistentFile := filepath.Join(os.TempDir(), "non-existent-file.xml")
	os.Remove(nonExistentFile) // 确保不存在

	// 创建解析器
	parser := NewConfigParser()

	// 测试有效的配置文件
	t.Run("Valid config file", func(t *testing.T) {
		config, err := parser.ParseFromFile(validFile)
		if err != nil {
			t.Fatalf("ParseFromFile() error = %v", err)
		}
		if config == nil {
			t.Fatal("ParseFromFile() returned nil config")
		}
		// 验证解析结果
		if len(config.PackageSources.Add) == 0 {
			t.Error("PackageSources.Add should not be empty")
		}
	})

	// 测试无效的配置文件
	t.Run("Invalid config file", func(t *testing.T) {
		_, err := parser.ParseFromFile(invalidFile)
		if err == nil {
			t.Fatal("ParseFromFile() expected error for invalid file")
		}
		if !errors.IsFormatError(err) && !errors.IsParseError(err) {
			t.Errorf("Expected format or parse error, got %v", err)
		}
	})

	// 测试空文件
	t.Run("Empty file", func(t *testing.T) {
		_, err := parser.ParseFromFile(emptyFile)
		if err == nil {
			t.Fatal("ParseFromFile() expected error for empty file")
		}
		if err != errors.ErrEmptyConfigFile {
			t.Errorf("Expected empty file error, got %v", err)
		}
	})

	// 测试不存在的文件
	t.Run("Non-existent file", func(t *testing.T) {
		_, err := parser.ParseFromFile(nonExistentFile)
		if err == nil {
			t.Fatal("ParseFromFile() expected error for non-existent file")
		}
		if err != errors.ErrConfigFileNotFound {
			t.Errorf("Expected file not found error, got %v", err)
		}
	})
}

func TestParseFromContent(t *testing.T) {
	// 创建解析器
	parser := NewConfigParser()

	// 测试有效的 XML 内容
	t.Run("Valid XML content", func(t *testing.T) {
		content := []byte(nugetTesting.ValidNuGetConfig())
		config, err := parser.ParseFromContent(content)
		if err != nil {
			t.Fatalf("ParseFromContent() error = %v", err)
		}
		if config == nil {
			t.Fatal("ParseFromContent() returned nil config")
		}
		// 验证解析结果
		if len(config.PackageSources.Add) == 0 {
			t.Error("PackageSources.Add should not be empty")
		}
	})

	// 测试无效的 XML 内容
	t.Run("Invalid XML content", func(t *testing.T) {
		content := []byte(nugetTesting.InvalidNuGetConfig())
		_, err := parser.ParseFromContent(content)
		if err == nil {
			t.Fatal("ParseFromContent() expected error for invalid content")
		}
		// 这里可能返回格式错误或解析错误
		if !errors.IsFormatError(err) && !errors.IsParseError(err) {
			t.Errorf("Expected format or parse error, got %v", err)
		}
	})

	// 测试空内容
	t.Run("Empty content", func(t *testing.T) {
		content := []byte("")
		_, err := parser.ParseFromContent(content)
		if err == nil {
			t.Fatal("ParseFromContent() expected error for empty content")
		}
		if !errors.IsFormatError(err) {
			t.Errorf("Expected format error, got %v", err)
		}
	})

	// 测试有效但没有包源的 XML 内容
	t.Run("Valid XML without package sources", func(t *testing.T) {
		// XML 有效但没有包源
		content := []byte(`<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
  </packageSources>
</configuration>`)
		_, err := parser.ParseFromContent(content)
		if err == nil {
			t.Fatal("ParseFromContent() expected error for content without package sources")
		}
		if !errors.IsParseError(err) {
			t.Errorf("Expected parse error, got %v", err)
		}
	})

	// 测试带有 clear 标志的 XML 内容
	t.Run("Valid XML with clear flag", func(t *testing.T) {
		content := []byte(nugetTesting.EmptyNuGetConfig())
		config, err := parser.ParseFromContent(content)
		if err != nil {
			t.Fatalf("ParseFromContent() error = %v", err)
		}
		if config == nil {
			t.Fatal("ParseFromContent() returned nil config")
		}
		if !config.PackageSources.Clear {
			t.Error("PackageSources.Clear should be true")
		}
	})
}

func TestParseFromReader(t *testing.T) {
	// 创建解析器
	parser := NewConfigParser()

	// 测试有效的读取器
	t.Run("Valid reader", func(t *testing.T) {
		reader := strings.NewReader(nugetTesting.ValidNuGetConfig())
		config, err := parser.ParseFromReader(reader)
		if err != nil {
			t.Fatalf("ParseFromReader() error = %v", err)
		}
		if config == nil {
			t.Fatal("ParseFromReader() returned nil config")
		}
	})

	// 测试无效的读取器
	t.Run("Invalid reader", func(t *testing.T) {
		reader := strings.NewReader(nugetTesting.InvalidNuGetConfig())
		_, err := parser.ParseFromReader(reader)
		if err == nil {
			t.Fatal("ParseFromReader() expected error for invalid reader")
		}
	})

	// 测试读取错误
	t.Run("Reader error", func(t *testing.T) {
		reader := &errorReader{err: io.ErrUnexpectedEOF}
		_, err := parser.ParseFromReader(reader)
		if err == nil {
			t.Fatal("ParseFromReader() expected error for reader error")
		}
		if !strings.Contains(err.Error(), "failed to read from reader") {
			t.Errorf("Expected read error, got %v", err)
		}
	})
}

func TestParseFromString(t *testing.T) {
	// 创建解析器
	parser := NewConfigParser()

	// 测试有效的字符串
	t.Run("Valid string", func(t *testing.T) {
		config, err := parser.ParseFromString(nugetTesting.ValidNuGetConfig())
		if err != nil {
			t.Fatalf("ParseFromString() error = %v", err)
		}
		if config == nil {
			t.Fatal("ParseFromString() returned nil config")
		}
	})

	// 测试无效的字符串
	t.Run("Invalid string", func(t *testing.T) {
		_, err := parser.ParseFromString(nugetTesting.InvalidNuGetConfig())
		if err == nil {
			t.Fatal("ParseFromString() expected error for invalid string")
		}
	})

	// 测试空字符串
	t.Run("Empty string", func(t *testing.T) {
		_, err := parser.ParseFromString("")
		if err == nil {
			t.Fatal("ParseFromString() expected error for empty string")
		}
	})
}

func TestFindAndParseConfig(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建配置文件
	configFile := filepath.Join(tempDir, "NuGet.Config")
	nugetTesting.CreateNuGetConfigFile(t, configFile, nugetTesting.ValidNuGetConfig())

	// 保存当前目录并切换到临时目录
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Fatalf("Failed to restore directory: %v", err)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// 创建一个自定义的解析器，只搜索临时目录中的文件
	parser := &ConfigParser{
		DefaultConfigSearchPaths: []string{configFile},
	}

	// 测试查找和解析配置
	config, path, err := parser.FindAndParseConfig()
	if err != nil {
		t.Fatalf("FindAndParseConfig() error = %v", err)
	}
	if config == nil {
		t.Fatal("FindAndParseConfig() returned nil config")
	}
	if path != configFile {
		t.Errorf("FindAndParseConfig() path = %v, want %v", path, configFile)
	}

	// 测试无法找到配置
	os.Remove(configFile)
	_, _, err = parser.FindAndParseConfig()
	if err == nil {
		t.Fatal("FindAndParseConfig() expected error when no config file found")
	}
	if err != errors.ErrConfigFileNotFound {
		t.Errorf("Expected file not found error, got %v", err)
	}
}

func TestSerializeToXML(t *testing.T) {
	// 创建解析器
	parser := NewConfigParser()

	// 先解析一个有效的配置
	config, err := parser.ParseFromString(nugetTesting.ValidNuGetConfig())
	if err != nil {
		t.Fatalf("ParseFromString() error = %v", err)
	}

	// 序列化为 XML
	xmlString, err := parser.SerializeToXML(config)
	if err != nil {
		t.Fatalf("SerializeToXML() error = %v", err)
	}

	// 检查 XML 头部
	if !strings.HasPrefix(xmlString, `<?xml version="1.0" encoding="utf-8"?>`) {
		t.Errorf("SerializeToXML() should start with XML declaration")
	}

	// 再解析序列化后的 XML，确保可以正确解析
	newConfig, err := parser.ParseFromString(xmlString)
	if err != nil {
		t.Fatalf("Failed to parse serialized XML: %v", err)
	}

	// 验证数据是否一致
	if len(newConfig.PackageSources.Add) != len(config.PackageSources.Add) {
		t.Errorf("SerializeToXML() produced invalid XML: package sources count mismatch")
	}
}

func TestSaveToFile(t *testing.T) {
	// 创建临时目录
	tempDir := nugetTesting.CreateTempDir(t)
	defer os.RemoveAll(tempDir)

	// 创建解析器
	parser := NewConfigParser()

	// 先解析一个有效的配置
	config, err := parser.ParseFromString(nugetTesting.ValidNuGetConfig())
	if err != nil {
		t.Fatalf("ParseFromString() error = %v", err)
	}

	// 保存到文件
	outputFile := filepath.Join(tempDir, "output.xml")
	err = parser.SaveToFile(config, outputFile)
	if err != nil {
		t.Fatalf("SaveToFile() error = %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("SaveToFile() did not create file")
	}

	// 读取保存的文件内容
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	// 检查 XML 头部
	if !bytes.HasPrefix(content, []byte(`<?xml version="1.0" encoding="utf-8"?>`)) {
		t.Errorf("SaveToFile() should create file with XML declaration")
	}

	// 再解析保存的文件，确保可以正确解析
	newConfig, err := parser.ParseFromFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to parse saved file: %v", err)
	}

	// 验证数据是否一致
	if len(newConfig.PackageSources.Add) != len(config.PackageSources.Add) {
		t.Errorf("SaveToFile() created invalid XML: package sources count mismatch")
	}

	// 测试保存到不存在的目录
	nestedOutputFile := filepath.Join(tempDir, "nested", "deeply", "output.xml")
	err = parser.SaveToFile(config, nestedOutputFile)
	if err != nil {
		t.Fatalf("SaveToFile() error for nested directory = %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(nestedOutputFile); os.IsNotExist(err) {
		t.Fatalf("SaveToFile() did not create file in nested directory")
	}
}

// 错误读取器，用于测试读取错误的情况
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
