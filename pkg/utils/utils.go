// Package utils 提供解析NuGet配置文件的辅助工具函数
package utils

import (
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileExists 检查文件是否存在
//
// FileExists 检查指定路径的文件是否存在且不是目录。
//
// 参数:
//   - filePath: 要检查的文件路径
//
// 返回值:
//   - bool: 如果文件存在且不是目录则返回 true，否则返回 false
//
// 示例:
//
//	// 检查配置文件是否存在
//	configPath := "/path/to/NuGet.Config"
//	if utils.FileExists(configPath) {
//	    fmt.Printf("配置文件 %s 存在\n", configPath)
//	    // 读取文件...
//	} else {
//	    fmt.Printf("配置文件 %s 不存在\n", configPath)
//	    // 创建默认配置...
//	}
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsValidXML 检查字符串是否为有效的XML
//
// IsValidXML 验证给定的字符串是否包含有效的 XML 内容。
// 该函数会尝试解析 XML 内容，如果能够完整解析，则认为是有效的 XML。
//
// 参数:
//   - xmlStr: 要验证的 XML 字符串
//
// 返回值:
//   - bool: 如果字符串包含有效的 XML 则返回 true，否则返回 false
//     空字符串会返回 false
//
// 示例:
//
//	// 验证 XML 内容
//	configContent := `<?xml version="1.0" encoding="utf-8"?>
//	<configuration>
//	  <packageSources>
//	    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" />
//	  </packageSources>
//	</configuration>`
//
//	if utils.IsValidXML(configContent) {
//	    fmt.Println("XML 格式有效")
//	    // 进一步处理...
//	} else {
//	    fmt.Println("XML 格式无效")
//	    // 处理错误...
//	}
//
//	// 检查无效 XML
//	invalidXML := `<?xml version="1.0" encoding="utf-8"?>
//	<configuration>
//	  <packageSources>
//	    <add key="nuget.org" value="https://api.nuget.org/v3/index.json"
//	  </packageSources>
//	</configuration>`
//
//	if !utils.IsValidXML(invalidXML) {
//	    fmt.Println("XML 格式无效，缺少闭合标签")
//	}
func IsValidXML(xmlStr string) bool {
	if xmlStr == "" {
		return false
	}

	decoder := xml.NewDecoder(strings.NewReader(xmlStr))
	for {
		_, err := decoder.Token()
		if err == io.EOF {
			return true
		}
		if err != nil {
			return false
		}
	}
}

// IsAbsolutePath 检查路径是否为绝对路径
//
// IsAbsolutePath 判断给定的路径是否为绝对路径。
// 绝对路径包含从文件系统根目录开始的完整路径信息。
//
// 参数:
//   - path: 要检查的文件路径
//
// 返回值:
//   - bool: 如果路径是绝对路径则返回 true，否则返回 false
//
// 示例:
//
//	// 在不同操作系统下的绝对路径检查
//	paths := []string{
//	    "/etc/nuget/NuGet.Config",     // Unix/Linux 绝对路径
//	    "C:\\Users\\user\\NuGet.Config", // Windows 绝对路径
//	    "../NuGet.Config",              // 相对路径
//	    "NuGet.Config",                 // 相对路径
//	}
//
//	for _, path := range paths {
//	    if utils.IsAbsolutePath(path) {
//	        fmt.Printf("%s 是绝对路径\n", path)
//	    } else {
//	        fmt.Printf("%s 是相对路径\n", path)
//	    }
//	}
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// NormalizePath 标准化路径
//
// NormalizePath 清理并标准化文件路径，删除多余的分隔符和相对路径引用。
// 此函数使用系统的路径分隔符，并解析 "." 和 ".." 等路径元素。
//
// 参数:
//   - path: 要标准化的文件路径
//
// 返回值:
//   - string: 标准化后的文件路径
//
// 示例:
//
//	// 标准化各种路径
//	paths := []string{
//	    "/etc//nuget/./NuGet.Config",     // 包含多余斜杠和当前目录引用
//	    "/etc/nuget/../nuget/NuGet.Config", // 包含上级目录引用
//	    "C:\\Users\\user\\..\\user\\NuGet.Config", // Windows 路径
//	}
//
//	for _, path := range paths {
//	    normalized := utils.NormalizePath(path)
//	    fmt.Printf("原路径: %s\n标准化后: %s\n\n", path, normalized)
//	}
//
//	// 输出示例:
//	// 原路径: /etc//nuget/./NuGet.Config
//	// 标准化后: /etc/nuget/NuGet.Config
//	//
//	// 原路径: /etc/nuget/../nuget/NuGet.Config
//	// 标准化后: /etc/nuget/NuGet.Config
//	//
//	// 原路径: C:\Users\user\..\user\NuGet.Config
//	// 标准化后: C:\Users\user\NuGet.Config
func NormalizePath(path string) string {
	return filepath.Clean(path)
}

// JoinPaths 连接路径
//
// JoinPaths 使用操作系统特定的路径分隔符将多个路径元素连接成单一路径。
// 此函数会自动处理不同路径元素之间的分隔符，避免出现多余或缺失的分隔符问题。
//
// 参数:
//   - basePath: 基础路径，作为路径的起始部分
//   - paths: 可变数量的附加路径元素，将依次连接到基础路径后
//
// 返回值:
//   - string: 连接后的完整路径
//
// 示例:
//
//	// 连接路径元素
//	basePath := "/etc"
//	configPath := utils.JoinPaths(basePath, "nuget", "NuGet.Config")
//	fmt.Printf("连接后的路径: %s\n", configPath)
//	// 输出: 连接后的路径: /etc/nuget/NuGet.Config
//
//	// Windows 示例
//	winBase := "C:\\Users"
//	winPath := utils.JoinPaths(winBase, "username", ".nuget", "packages")
//	fmt.Printf("Windows 路径: %s\n", winPath)
//	// 输出: Windows 路径: C:\Users\username\.nuget\packages
//
//	// 处理尾部斜杠
//	trailingSlash := "/home/user/"
//	result := utils.JoinPaths(trailingSlash, "nuget/packages")
//	fmt.Printf("处理斜杠: %s\n", result)
//	// 输出: 处理斜杠: /home/user/nuget/packages
func JoinPaths(basePath string, paths ...string) string {
	elems := []string{basePath}
	elems = append(elems, paths...)
	return filepath.Join(elems...)
}

// ResolvePath 解析路径，如果是相对路径则根据basePath解析，否则返回原路径
//
// ResolvePath 将相对路径转换为绝对路径，如果输入已经是绝对路径则直接返回标准化后的路径。
// 此函数适用于处理配置文件中的路径引用，确保它们指向正确的位置。
//
// 参数:
//   - basePath: 基础路径，用于解析相对路径
//   - path: 要解析的路径，可以是相对路径或绝对路径
//
// 返回值:
//   - string: 解析后的绝对路径（标准化后）
//
// 示例:
//
//	// 解析配置文件中的相对路径
//	basePath := "/etc/nuget"
//
//	relativePath := "../packages/cache"
//	absolutePath := "/var/nuget/packages"
//
//	// 解析相对路径
//	resolvedRelative := utils.ResolvePath(basePath, relativePath)
//	fmt.Printf("相对路径 '%s' 解析为: %s\n", relativePath, resolvedRelative)
//	// 输出: 相对路径 '../packages/cache' 解析为: /etc/packages/cache
//
//	// 处理已经是绝对路径的情况
//	resolvedAbsolute := utils.ResolvePath(basePath, absolutePath)
//	fmt.Printf("绝对路径 '%s' 保持不变: %s\n", absolutePath, resolvedAbsolute)
//	// 输出: 绝对路径 '/var/nuget/packages' 保持不变: /var/nuget/packages
func ResolvePath(basePath, path string) string {
	if IsAbsolutePath(path) {
		return NormalizePath(path)
	}
	return NormalizePath(JoinPaths(basePath, path))
}

// ExpandEnvVars 展开路径中的环境变量
//
// ExpandEnvVars 将路径中的环境变量占位符替换为实际的环境变量值。
// 支持的占位符格式取决于操作系统：
//   - Unix/Linux/macOS: $VAR 或 ${VAR}
//   - Windows: %VAR%
//
// 参数:
//   - path: 包含环境变量的路径字符串
//
// 返回值:
//   - string: 环境变量被替换后的路径
//
// 示例:
//
//	// Unix/Linux/macOS 示例
//	unixPath := "$HOME/.nuget/packages"
//	expandedUnixPath := utils.ExpandEnvVars(unixPath)
//	fmt.Printf("展开前: %s\n展开后: %s\n", unixPath, expandedUnixPath)
//	// 输出示例: 展开前: $HOME/.nuget/packages
//	//         展开后: /home/user/.nuget/packages
//
//	// Windows 示例
//	winPath := "%USERPROFILE%\\.nuget\\packages"
//	expandedWinPath := utils.ExpandEnvVars(winPath)
//	fmt.Printf("展开前: %s\n展开后: %s\n", winPath, expandedWinPath)
//	// 输出示例: 展开前: %USERPROFILE%\.nuget\packages
//	//         展开后: C:\Users\username\.nuget\packages
//
//	// 处理不存在的环境变量
//	nonExistentVar := "${NONEXISTENT_VAR}/packages"
//	expanded := utils.ExpandEnvVars(nonExistentVar)
//	fmt.Printf("不存在的变量: %s\n", expanded)
//	// 输出: 不存在的变量: /packages
func ExpandEnvVars(path string) string {
	return os.ExpandEnv(path)
}

// IsURL 判断字符串是否为URL
//
// IsURL 检查给定的字符串是否是有效的 HTTP 或 HTTPS URL。
// 此函数通过简单地检查字符串是否以 "http://" 或 "https://" 开头来判断。
//
// 参数:
//   - str: 要检查的字符串
//
// 返回值:
//   - bool: 如果字符串是 HTTP 或 HTTPS URL 则返回 true，否则返回 false
//
// 注意:
//   - 此函数不验证 URL 的完整有效性，只检查前缀
//   - 检查时不区分大小写
//
// 示例:
//
//	// 检查各种字符串
//	urls := []string{
//	    "https://api.nuget.org/v3/index.json",
//	    "http://nuget.example.com/feed",
//	    "HTTP://mixed-case.example.com",
//	    "ftp://invalid-scheme.com",
//	    "/local/path/no/scheme",
//	    "C:\\Windows\\Path",
//	}
//
//	for _, url := range urls {
//	    if utils.IsURL(url) {
//	        fmt.Printf("%s 是有效的 HTTP/HTTPS URL\n", url)
//	    } else {
//	        fmt.Printf("%s 不是有效的 HTTP/HTTPS URL\n", url)
//	    }
//	}
//
//	// 输出示例:
//	// https://api.nuget.org/v3/index.json 是有效的 HTTP/HTTPS URL
//	// http://nuget.example.com/feed 是有效的 HTTP/HTTPS URL
//	// HTTP://mixed-case.example.com 是有效的 HTTP/HTTPS URL
//	// ftp://invalid-scheme.com 不是有效的 HTTP/HTTPS URL
//	// /local/path/no/scheme 不是有效的 HTTP/HTTPS URL
//	// C:\Windows\Path 不是有效的 HTTP/HTTPS URL
func IsURL(str string) bool {
	lowerStr := strings.ToLower(str)
	return strings.HasPrefix(lowerStr, "http://") || strings.HasPrefix(lowerStr, "https://")
}

// ReadFile 读取文件内容
//
// ReadFile 读取指定路径文件的全部内容并返回。
// 封装了标准库的 os.ReadFile 函数，提供与库其他函数一致的错误处理方式。
//
// 参数:
//   - filePath: 要读取的文件路径
//
// 返回值:
//   - []byte: 文件的全部内容，以字节切片形式返回
//   - error: 如果读取过程中发生错误则返回相应的错误；如果成功则为 nil
//
// 错误:
//   - 当文件不存在、权限不足或其他 I/O 错误时，返回相应的 error
//
// 示例:
//
//	// 读取配置文件内容
//	configPath := "/path/to/NuGet.Config"
//	content, err := utils.ReadFile(configPath)
//	if err != nil {
//	    if os.IsNotExist(err) {
//	        fmt.Printf("配置文件 %s 不存在\n", configPath)
//	    } else if os.IsPermission(err) {
//	        fmt.Printf("没有权限读取文件 %s\n", configPath)
//	    } else {
//	        fmt.Printf("读取配置文件失败: %v\n", err)
//	    }
//	    return
//	}
//
//	// 使用读取到的内容
//	fmt.Printf("读取了 %d 字节的配置数据\n", len(content))
//
//	// 转换为字符串并处理
//	configStr := string(content)
//	fmt.Printf("配置文件内容:\n%s\n", configStr)
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// WriteToFile 将内容写入文件
//
// WriteToFile 将提供的数据写入指定路径的文件中。
// 如果文件所在的目录不存在，会先创建必要的目录结构，然后再写入文件。
// 如果文件已存在，会被覆盖。
//
// 参数:
//   - filePath: 要写入的文件路径
//   - data: 要写入的数据，以字节切片形式提供
//
// 返回值:
//   - error: 如果写入过程中发生错误则返回相应的错误；如果成功则为 nil
//
// 注意:
//   - 创建的目录权限为 0755 (用户可读写执行，组和其他用户可读执行)
//   - 创建的文件权限为 0644 (用户可读写，组和其他用户可读)
//
// 示例:
//
//	// 保存配置文件
//	configPath := "/path/to/NuGet.Config"
//	configData := []byte(`<?xml version="1.0" encoding="utf-8"?>
//	<configuration>
//	  <packageSources>
//	    <add key="nuget.org" value="https://api.nuget.org/v3/index.json" protocolVersion="3" />
//	    <add key="company-internal" value="https://nuget.company.com/v3/index.json" />
//	  </packageSources>
//	</configuration>`)
//
//	err := utils.WriteToFile(configPath, configData)
//	if err != nil {
//	    if os.IsPermission(err) {
//	        fmt.Printf("没有权限写入文件 %s\n", configPath)
//	    } else {
//	        fmt.Printf("保存配置文件失败: %v\n", err)
//	    }
//	    return
//	}
//
//	fmt.Printf("配置已成功保存到 %s\n", configPath)
func WriteToFile(filePath string, data []byte) error {
	dir := filepath.Dir(filePath)

	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// TrimWhitespace 去除字符串首尾的空白字符
//
// TrimWhitespace 移除字符串开头和结尾的所有空白字符，包括空格、制表符、换行符等。
// 此函数是对 strings.TrimSpace 的封装，提供了一致的命名风格。
//
// 参数:
//   - s: 要处理的字符串
//
// 返回值:
//   - string: 去除首尾空白后的字符串
//
// 示例:
//
//	// 处理各种含有空白的字符串
//	inputs := []string{
//	    "  leading spaces",
//	    "trailing spaces  ",
//	    "  both ends  ",
//	    "\t\nwhitespace\r\n",
//	    "no whitespace",
//	    "   ",  // 只有空白
//	    "",     // 空字符串
//	}
//
//	for _, input := range inputs {
//	    trimmed := utils.TrimWhitespace(input)
//	    fmt.Printf("原始: %q\n去除空白后: %q\n\n", input, trimmed)
//	}
//
//	// 输出示例:
//	// 原始: "  leading spaces"
//	// 去除空白后: "leading spaces"
//	//
//	// 原始: "trailing spaces  "
//	// 去除空白后: "trailing spaces"
//	//
//	// 原始: "  both ends  "
//	// 去除空白后: "both ends"
//	//
//	// 原始: "\t\nwhitespace\r\n"
//	// 去除空白后: "whitespace"
//	//
//	// 原始: "no whitespace"
//	// 去除空白后: "no whitespace"
//	//
//	// 原始: "   "
//	// 去除空白后: ""
//	//
//	// 原始: ""
//	// 去除空白后: ""
func TrimWhitespace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty 检查字符串是否为空或只包含空白字符
//
// IsEmpty 判断给定的字符串在去除首尾空白后是否为空。
// 此函数可用于验证配置值或输入字段是否有实际内容。
//
// 参数:
//   - s: 要检查的字符串
//
// 返回值:
//   - bool: 如果字符串为空或仅包含空白字符则返回 true，否则返回 false
//
// 示例:
//
//	// 检查各种字符串
//	checks := []string{
//	    "",             // 空字符串
//	    "   ",          // 只有空白
//	    " \t\n\r ",     // 只有各种空白字符
//	    "hello",        // 有内容
//	    "  content  ",  // 首尾有空白但有内容
//	}
//
//	for _, s := range checks {
//	    if utils.IsEmpty(s) {
//	        fmt.Printf("%q 是空的\n", s)
//	    } else {
//	        fmt.Printf("%q 不是空的\n", s)
//	    }
//	}
//
//	// 输出示例:
//	// "" 是空的
//	// "   " 是空的
//	// " \t\n\r " 是空的
//	// "hello" 不是空的
//	// "  content  " 不是空的
//
//	// 在验证用户输入时的应用
//	userInput := "  "
//	if utils.IsEmpty(userInput) {
//	    fmt.Println("请提供有效的输入")
//	} else {
//	    fmt.Println("输入有效")
//	}
func IsEmpty(s string) bool {
	return TrimWhitespace(s) == ""
}
