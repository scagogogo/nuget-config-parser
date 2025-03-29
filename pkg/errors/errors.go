// Package errors 定义与 NuGet 配置文件解析相关的错误类型
package errors

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfigFormat 表示配置文件格式不正确的错误
	ErrInvalidConfigFormat = errors.New("invalid nuget config format")

	// ErrConfigFileNotFound 表示找不到配置文件的错误
	ErrConfigFileNotFound = errors.New("nuget config file not found")

	// ErrEmptyConfigFile 表示配置文件为空的错误
	ErrEmptyConfigFile = errors.New("empty nuget config file")

	// ErrXMLParsing 表示XML解析错误
	ErrXMLParsing = errors.New("xml parsing error")

	// ErrMissingRequiredElement 表示缺少必需元素的错误
	ErrMissingRequiredElement = errors.New("missing required element in config")
)

// ParseError 解析错误结构，提供额外上下文信息
type ParseError struct {
	// BaseErr 基础错误
	BaseErr error

	// Line 出错的行号
	Line int

	// Position 出错的位置
	Position int

	// Context 错误上下文信息
	Context string
}

// Error 格式化解析错误信息
func (e *ParseError) Error() string {
	if e.Line > 0 && e.Position > 0 {
		return fmt.Sprintf("parse error at line %d position %d: %s - %v", e.Line, e.Position, e.Context, e.BaseErr)
	}
	if e.Context != "" {
		return fmt.Sprintf("parse error: %s - %v", e.Context, e.BaseErr)
	}
	return fmt.Sprintf("parse error: %v", e.BaseErr)
}

// Unwrap 返回基础错误，支持 errors.Is 和 errors.As
func (e *ParseError) Unwrap() error {
	return e.BaseErr
}

// NewParseError 创建新的解析错误
func NewParseError(baseErr error, line, position int, context string) *ParseError {
	return &ParseError{
		BaseErr:  baseErr,
		Line:     line,
		Position: position,
		Context:  context,
	}
}

// IsNotFoundError 判断是否为找不到配置文件的错误
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrConfigFileNotFound)
}

// IsParseError 判断是否为解析错误
func IsParseError(err error) bool {
	var parseErr *ParseError
	return errors.As(err, &parseErr)
}

// IsFormatError 判断是否为格式错误
func IsFormatError(err error) bool {
	return errors.Is(err, ErrInvalidConfigFormat)
}
