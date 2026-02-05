// Package low provides a low-level, performance-oriented API for parsing SysML v2.
// This package offers direct access to the ANTLR lexer and parser with minimal overhead.
package low

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// SyntaxError represents a single syntax error from the lexer or parser.
type SyntaxError struct {
	Line    int
	Column  int
	Message string
	Source  string // "lexer" or "parser"
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.Source, e.Line, e.Column, e.Message)
}

// Unwrap returns nil since SyntaxError doesn't wrap another error.
func (e *SyntaxError) Unwrap() error {
	return nil
}

// ErrorCollector collects syntax errors during lexing and parsing.
// It implements antlr.ErrorListener for integration with ANTLR.
type ErrorCollector struct {
	*antlr.DefaultErrorListener
	errors []*SyntaxError
	source string
}

// NewErrorCollector creates a new error collector.
func NewErrorCollector(source string) *ErrorCollector {
	return &ErrorCollector{
		errors: make([]*SyntaxError, 0, 8),
		source: source,
	}
}

// SyntaxError implements antlr.ErrorListener.
func (c *ErrorCollector) SyntaxError(
	recognizer antlr.Recognizer,
	offendingSymbol interface{},
	line, column int,
	msg string,
	e antlr.RecognitionException,
) {
	c.errors = append(c.errors, &SyntaxError{
		Line:    line,
		Column:  column,
		Message: msg,
		Source:  c.source,
	})
}

// Errors returns all collected errors.
func (c *ErrorCollector) Errors() []*SyntaxError {
	return c.errors
}

// HasErrors returns true if any errors were collected.
func (c *ErrorCollector) HasErrors() bool {
	return len(c.errors) > 0
}

// Clear removes all collected errors.
func (c *ErrorCollector) Clear() {
	c.errors = c.errors[:0]
}

// ParseErrors aggregates errors from both lexer and parser.
type ParseErrors struct {
	LexerErrors  []*SyntaxError
	ParserErrors []*SyntaxError
}

// HasErrors returns true if there are any lexer or parser errors.
func (e *ParseErrors) HasErrors() bool {
	return len(e.LexerErrors) > 0 || len(e.ParserErrors) > 0
}

// All returns all errors (lexer + parser) in order.
func (e *ParseErrors) All() []*SyntaxError {
	all := make([]*SyntaxError, 0, len(e.LexerErrors)+len(e.ParserErrors))
	all = append(all, e.LexerErrors...)
	all = append(all, e.ParserErrors...)
	return all
}

// Error implements the error interface.
func (e *ParseErrors) Error() string {
	if !e.HasErrors() {
		return ""
	}

	var sb strings.Builder
	totalErrors := len(e.LexerErrors) + len(e.ParserErrors)
	sb.Grow(totalErrors * 50)

	sb.WriteString(fmt.Sprintf("%d error(s):\n", totalErrors))

	for _, err := range e.LexerErrors {
		sb.WriteString("  ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	for _, err := range e.ParserErrors {
		sb.WriteString("  ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}

// Err returns the ParseErrors as an error if there are any errors,
// otherwise returns nil. This enables idiomatic error handling:
//
//	if err := result.Err(); err != nil {
//		return err
//	}
func (e *ParseErrors) Err() error {
	if !e.HasErrors() {
		return nil
	}
	return e
}

// Unwrap returns a joined error containing all syntax errors.
// This allows errors.Is and errors.As to work with ParseErrors.
func (e *ParseErrors) Unwrap() []error {
	all := e.All()
	errs := make([]error, len(all))
	for i, err := range all {
		errs[i] = err
	}
	return errs
}

// UnwrapError returns the first error if there's only one error,
// otherwise returns nil. This allows errors.Is to work with single errors.
func (e *ParseErrors) UnwrapError() error {
	if len(e.LexerErrors) == 1 && len(e.ParserErrors) == 0 {
		return e.LexerErrors[0]
	}
	if len(e.LexerErrors) == 0 && len(e.ParserErrors) == 1 {
		return e.ParserErrors[0]
	}
	return nil
}
