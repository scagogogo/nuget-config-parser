// Package testing 提供NuGet配置解析器测试所需的工具函数
package testing

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// CreateTempFile 创建一个临时文件并写入内容
func CreateTempFile(t *testing.T, content string) string {
	tempFile, err := os.CreateTemp("", "nuget-test-*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tempFile.Name()
}

// CreateTempDir 创建一个临时目录
func CreateTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "nuget-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tempDir
}

// SetupEnv 设置环境变量并返回恢复函数
func SetupEnv(t *testing.T, key, value string) func() {
	oldValue, exists := os.LookupEnv(key)
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}

	return func() {
		if exists {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	}
}

// ValidNuGetConfig 返回一个有效的NuGet配置XML字符串
func ValidNuGetConfig() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
    <add key="localSource" value="C:\Packages" />
  </packageSources>
  <packageSourceCredentials>
    <nuget.org>
      <add key="Username" value="testuser" />
      <add key="ClearTextPassword" value="testpass" />
    </nuget.org>
  </packageSourceCredentials>
  <disabledPackageSources>
    <add key="localSource" value="true" />
  </disabledPackageSources>
  <activePackageSource>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
  </activePackageSource>
  <config>
    <add key="globalPackagesFolder" value="%USERPROFILE%\.nuget\packages" />
  </config>
</configuration>`
}

// InvalidNuGetConfig 返回一个无效的NuGet配置XML字符串
func InvalidNuGetConfig() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3"
  </packageSources>
</configuration>`
}

// EmptyNuGetConfig 返回一个空的NuGet配置XML字符串
func EmptyNuGetConfig() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources clear="true">
  </packageSources>
</configuration>`
}

// CreateNuGetConfigFile 在指定路径创建NuGet配置文件
func CreateNuGetConfigFile(t *testing.T, path string, content string) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write config file %s: %v", path, err)
	}
}

// RemoveIfExists 如果文件存在则删除
func RemoveIfExists(t *testing.T, path string) {
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			t.Fatalf("Failed to remove file %s: %v", path, err)
		}
	}
}

// CompareFiles 比较两个文件的内容是否相同
func CompareFiles(t *testing.T, file1, file2 string) bool {
	content1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", file1, err)
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", file2, err)
	}

	return string(content1) == string(content2)
}

// StringReader 从字符串创建io.Reader
func StringReader(s string) io.Reader {
	return strings.NewReader(s)
}

// CreateNuGetConfigWithSource 创建包含单一包源的NuGet配置XML字符串
func CreateNuGetConfigWithSource(key, value string) string {
	return `<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <packageSources>
    <add key="` + key + `" value="` + value + `" />
  </packageSources>
</configuration>`
}
