package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		baseErr  error
		line     int
		position int
		context  string
		want     string
	}{
		{
			name:     "With line and position",
			baseErr:  ErrXMLParsing,
			line:     10,
			position: 20,
			context:  "test context",
			want:     "parse error at line 10 position 20: test context - xml parsing error",
		},
		{
			name:     "Without line and position but with context",
			baseErr:  ErrXMLParsing,
			line:     0,
			position: 0,
			context:  "test context",
			want:     "parse error: test context - xml parsing error",
		},
		{
			name:     "Without line, position and context",
			baseErr:  ErrXMLParsing,
			line:     0,
			position: 0,
			context:  "",
			want:     "parse error: xml parsing error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ParseError{
				BaseErr:  tt.baseErr,
				Line:     tt.line,
				Position: tt.position,
				Context:  tt.context,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ParseError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseError_Unwrap(t *testing.T) {
	baseErr := ErrXMLParsing
	e := &ParseError{BaseErr: baseErr}

	if unwrapped := e.Unwrap(); unwrapped != baseErr {
		t.Errorf("ParseError.Unwrap() = %v, want %v", unwrapped, baseErr)
	}
}

func TestNewParseError(t *testing.T) {
	baseErr := ErrXMLParsing
	line := 10
	position := 20
	context := "test context"

	e := NewParseError(baseErr, line, position, context)

	if e.BaseErr != baseErr {
		t.Errorf("NewParseError().BaseErr = %v, want %v", e.BaseErr, baseErr)
	}
	if e.Line != line {
		t.Errorf("NewParseError().Line = %v, want %v", e.Line, line)
	}
	if e.Position != position {
		t.Errorf("NewParseError().Position = %v, want %v", e.Position, position)
	}
	if e.Context != context {
		t.Errorf("NewParseError().Context = %v, want %v", e.Context, context)
	}
}

func TestErrorPredicates(t *testing.T) {
	// Test IsNotFoundError
	t.Run("IsNotFoundError", func(t *testing.T) {
		// 测试直接使用定义的错误
		if !IsNotFoundError(ErrConfigFileNotFound) {
			t.Errorf("IsNotFoundError(ErrConfigFileNotFound) = false, want true")
		}

		// 测试包装错误
		wrappedErr := fmt.Errorf("wrapped: %w", ErrConfigFileNotFound)
		if !IsNotFoundError(wrappedErr) {
			t.Errorf("IsNotFoundError(wrappedErr) = false, want true")
		}

		// 测试其他错误
		otherErr := errors.New("other error")
		if IsNotFoundError(otherErr) {
			t.Errorf("IsNotFoundError(otherErr) = true, want false")
		}
	})

	// Test IsParseError
	t.Run("IsParseError", func(t *testing.T) {
		// 测试 ParseError 类型
		parseErr := NewParseError(ErrXMLParsing, 0, 0, "")
		if !IsParseError(parseErr) {
			t.Errorf("IsParseError(parseErr) = false, want true")
		}

		// 测试包装 ParseError
		wrappedErr := fmt.Errorf("wrapped: %w", parseErr)
		if !IsParseError(wrappedErr) {
			t.Errorf("IsParseError(wrappedErr) = false, want true")
		}

		// 测试其他错误
		otherErr := errors.New("other error")
		if IsParseError(otherErr) {
			t.Errorf("IsParseError(otherErr) = true, want false")
		}
	})

	// Test IsFormatError
	t.Run("IsFormatError", func(t *testing.T) {
		// 测试直接使用定义的错误
		if !IsFormatError(ErrInvalidConfigFormat) {
			t.Errorf("IsFormatError(ErrInvalidConfigFormat) = false, want true")
		}

		// 测试包装错误
		wrappedErr := fmt.Errorf("wrapped: %w", ErrInvalidConfigFormat)
		if !IsFormatError(wrappedErr) {
			t.Errorf("IsFormatError(wrappedErr) = false, want true")
		}

		// 测试其他错误
		otherErr := errors.New("other error")
		if IsFormatError(otherErr) {
			t.Errorf("IsFormatError(otherErr) = true, want false")
		}
	})
}
