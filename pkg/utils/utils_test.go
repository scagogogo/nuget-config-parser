package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExists(t *testing.T) {
	// 测试临时文件
	tempFile, err := os.CreateTemp("", "nuget-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// 测试存在的文件
	if !FileExists(tempFile.Name()) {
		t.Errorf("FileExists(%s) = false, want true", tempFile.Name())
	}

	// 测试不存在的文件
	nonExistentFile := tempFile.Name() + ".nonexistent"
	if FileExists(nonExistentFile) {
		t.Errorf("FileExists(%s) = true, want false", nonExistentFile)
	}

	// 测试目录
	tempDir, err := os.MkdirTemp("", "nuget-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if FileExists(tempDir) {
		t.Errorf("FileExists(%s) = true for directory, want false", tempDir)
	}
}

func TestIsValidXML(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		isValid bool
	}{
		{
			name:    "Valid XML",
			xml:     `<?xml version="1.0" encoding="utf-8"?><root><child>value</child></root>`,
			isValid: true,
		},
		{
			name:    "Empty string",
			xml:     "",
			isValid: false,
		},
		{
			name:    "Invalid XML - unclosed tag",
			xml:     `<?xml version="1.0" encoding="utf-8"?><root><child>value</child>`,
			isValid: false,
		},
		{
			name:    "Invalid XML - malformed",
			xml:     `<?xml version="1.0" encoding="utf-8"?><root><child>value</xyz></root>`,
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidXML(tt.xml)
			if got != tt.isValid {
				t.Errorf("IsValidXML() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestIsAbsolutePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "Absolute path",
			path: "/absolute/path",
			want: true,
		},
		{
			name: "Relative path",
			path: "relative/path",
			want: false,
		},
		{
			name: "Current directory",
			path: ".",
			want: false,
		},
		{
			name: "Parent directory",
			path: "..",
			want: false,
		},
		{
			name: "Root directory",
			path: "/",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAbsolutePath(tt.path)
			if got != tt.want {
				t.Errorf("IsAbsolutePath(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Already normalized",
			path: "/absolute/path",
			want: "/absolute/path",
		},
		{
			name: "Path with redundant separators",
			path: "/absolute//path",
			want: "/absolute/path",
		},
		{
			name: "Path with dot",
			path: "/absolute/./path",
			want: "/absolute/path",
		},
		{
			name: "Path with double dot",
			path: "/absolute/parent/../path",
			want: "/absolute/path",
		},
		{
			name: "Relative path",
			path: "relative/path",
			want: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePath(tt.path)
			if got != tt.want {
				t.Errorf("NormalizePath(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestJoinPaths(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		paths    []string
		want     string
	}{
		{
			name:     "Join absolute base with relative",
			basePath: "/base",
			paths:    []string{"path", "to", "file"},
			want:     filepath.Join("/base", "path", "to", "file"),
		},
		{
			name:     "Join relative base with relative",
			basePath: "base",
			paths:    []string{"path", "to", "file"},
			want:     filepath.Join("base", "path", "to", "file"),
		},
		{
			name:     "No additional paths",
			basePath: "/base",
			paths:    []string{},
			want:     "/base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JoinPaths(tt.basePath, tt.paths...)
			if got != tt.want {
				t.Errorf("JoinPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		basePath string
		path     string
		want     string
	}{
		{
			name:     "Resolve relative to absolute base",
			basePath: "/base/dir",
			path:     "relative/path",
			want:     filepath.Clean("/base/dir/relative/path"),
		},
		{
			name:     "Absolute path stays unchanged",
			basePath: "/base/dir",
			path:     "/absolute/path",
			want:     "/absolute/path",
		},
		{
			name:     "Normalize path with redundant separators",
			basePath: "/base/dir",
			path:     "relative//path",
			want:     filepath.Clean("/base/dir/relative/path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolvePath(tt.basePath, tt.path)
			if got != tt.want {
				t.Errorf("ResolvePath(%s, %s) = %v, want %v", tt.basePath, tt.path, got, tt.want)
			}
		})
	}
}

func TestExpandEnvVars(t *testing.T) {
	// 设置测试环境变量
	oldEnv := os.Getenv("NUGET_TEST_VAR")
	defer os.Setenv("NUGET_TEST_VAR", oldEnv)
	os.Setenv("NUGET_TEST_VAR", "test-value")

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Path with environment variable",
			path: "/path/to/$NUGET_TEST_VAR/file",
			want: "/path/to/test-value/file",
		},
		{
			name: "Path with environment variable with braces",
			path: "/path/to/${NUGET_TEST_VAR}/file",
			want: "/path/to/test-value/file",
		},
		{
			name: "Path without environment variable",
			path: "/path/to/file",
			want: "/path/to/file",
		},
		{
			name: "Path with undefined environment variable",
			path: "/path/to/$UNDEFINED_VAR/file",
			want: "/path/to//file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpandEnvVars(tt.path)
			if got != tt.want {
				t.Errorf("ExpandEnvVars(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			name: "HTTP URL",
			str:  "http://example.com",
			want: true,
		},
		{
			name: "HTTPS URL",
			str:  "https://example.com/path",
			want: true,
		},
		{
			name: "HTTP URL with uppercase",
			str:  "HTTP://example.com",
			want: true,
		},
		{
			name: "File path",
			str:  "/path/to/file",
			want: false,
		},
		{
			name: "Empty string",
			str:  "",
			want: false,
		},
		{
			name: "String starting with different protocol",
			str:  "ftp://example.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsURL(tt.str)
			if got != tt.want {
				t.Errorf("IsURL(%s) = %v, want %v", tt.str, got, tt.want)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "nuget-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write content to the file
	content := "test content"
	if _, err := tempFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Test reading existing file
	t.Run("Existing file", func(t *testing.T) {
		data, err := ReadFile(tempFile.Name())
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if string(data) != content {
			t.Errorf("ReadFile() = %q, want %q", string(data), content)
		}
	})

	// Test reading non-existent file
	t.Run("Non-existent file", func(t *testing.T) {
		nonExistentFile := tempFile.Name() + ".nonexistent"
		_, err := ReadFile(nonExistentFile)
		if err == nil {
			t.Errorf("ReadFile() expected error, got nil")
		}
	})
}

func TestWriteToFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "nuget-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test writing to a new file
	t.Run("New file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test.txt")
		content := []byte("test content")

		err := WriteToFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteToFile() error = %v", err)
		}

		// Verify file content
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(readContent) != string(content) {
			t.Errorf("File content = %q, want %q", string(readContent), string(content))
		}
	})

	// Test writing to a new file in a non-existent directory
	t.Run("New file in new directory", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "new-dir", "test.txt")
		content := []byte("test content")

		err := WriteToFile(filePath, content)
		if err != nil {
			t.Fatalf("WriteToFile() error = %v", err)
		}

		// Verify file content
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(readContent) != string(content) {
			t.Errorf("File content = %q, want %q", string(readContent), string(content))
		}
	})

	// Test overwriting existing file
	t.Run("Overwrite existing file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "existing.txt")
		initialContent := []byte("initial content")
		newContent := []byte("new content")

		// Create initial file
		err := os.WriteFile(filePath, initialContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create initial file: %v", err)
		}

		// Overwrite with new content
		err = WriteToFile(filePath, newContent)
		if err != nil {
			t.Fatalf("WriteToFile() error = %v", err)
		}

		// Verify file content
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(readContent) != string(newContent) {
			t.Errorf("File content = %q, want %q", string(readContent), string(newContent))
		}
	})
}

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "No whitespace",
			input: "test",
			want:  "test",
		},
		{
			name:  "Leading whitespace",
			input: "  test",
			want:  "test",
		},
		{
			name:  "Trailing whitespace",
			input: "test  ",
			want:  "test",
		},
		{
			name:  "Leading and trailing whitespace",
			input: "  test  ",
			want:  "test",
		},
		{
			name:  "Multiple whitespace types",
			input: "\n\t test \n\t",
			want:  "test",
		},
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "Only whitespace",
			input: "  \t\n  ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("TrimWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Non-empty string",
			input: "test",
			want:  false,
		},
		{
			name:  "Empty string",
			input: "",
			want:  true,
		},
		{
			name:  "Only whitespace",
			input: "  \t\n  ",
			want:  true,
		},
		{
			name:  "String with whitespace",
			input: "  test  ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmpty(tt.input)
			if got != tt.want {
				t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
