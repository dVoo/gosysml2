// Package sysml provides a high-level, developer-friendly API for working with SysML v2 models.
// This package offers convenient functions for parsing, traversing, and manipulating SysML content.
package sysml

import (
	"fmt"
	"strings"

	"github.com/dVoo/gosysml2/low"
)

// Error represents a user-friendly SysML parsing error.
type Error struct {
	Line    int
	Column  int
	Message string
	Context string // Additional context about where the error occurred
}

func (e *Error) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("line %d, column %d: %s (in %s)", e.Line, e.Column, e.Message, e.Context)
	}
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// Unwrap returns the underlying error for errors.Is/errors.As compatibility.
// Since Error is a leaf error type, it returns nil.
func (e *Error) Unwrap() error {
	return nil
}

// ParseError represents the result of a failed parse operation.
type ParseError struct {
	Items  []*Error
	Source string // Source identifier (filename or description)
	Input  string // The input that was being parsed (truncated if large)
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	if len(e.Items) == 0 {
		return "no errors"
	}

	var sb strings.Builder
	// Preallocate with estimated capacity for performance
	sb.Grow(len(e.Items)*80 + len(e.Source) + 50)
	if e.Source != "" {
		sb.WriteString(fmt.Sprintf("failed to parse %s: ", e.Source))
	} else {
		sb.WriteString("parse failed: ")
	}
	sb.WriteString(fmt.Sprintf("%d error(s)\n", len(e.Items)))

	for i, err := range e.Items {
		if i >= 10 {
			sb.WriteString(fmt.Sprintf("  ... and %d more errors\n", len(e.Items)-10))
			break
		}
		sb.WriteString("  - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}

	return sb.String()
}

// Unwrap implements Go 1.20+ multi-error unwrapping.
// This enables errors.Is(parseErr, targetErr) to check any of the contained errors.
func (e *ParseError) Unwrap() []error {
	result := make([]error, len(e.Items))
	for i, err := range e.Items {
		result[i] = err
	}
	return result
}

// Err returns nil if there are no errors, or the ParseError itself if there are errors.
// This enables idiomatic error handling:
//
//	if err := result.ParseError.Err(); err != nil { ... }
func (e *ParseError) Err() error {
	if e == nil || !e.HasErrors() {
		return nil
	}
	return e
}

// First returns the first error, or nil if there are no errors.
func (e *ParseError) First() *Error {
	if len(e.Items) == 0 {
		return nil
	}
	return e.Items[0]
}

// HasErrors returns true if there are any errors.
func (e *ParseError) HasErrors() bool {
	return len(e.Items) > 0
}

// convertFromLowLevel converts low-level errors to high-level errors.
func convertFromLowLevel(lowErrors *low.ParseErrors, source string) *ParseError {
	if lowErrors == nil || !lowErrors.HasErrors() {
		return nil
	}

	allLowErrors := lowErrors.All()
	errors := make([]*Error, len(allLowErrors))

	for i, e := range allLowErrors {
		line := e.Line - 1
		if line < 0 {
			line = -1
		}
		errors[i] = &Error{
			Line:    line,
			Column:  e.Column,
			Message: e.Message,
			Context: e.Source,
		}
	}

	return &ParseError{
		Items:  errors,
		Source: source,
	}
}

// ErrorList is a helper for collecting multiple errors.
type ErrorList struct {
	errors []*Error
}

// Add adds an error to the list.
func (l *ErrorList) Add(err *Error) {
	l.errors = append(l.errors, err)
}

// AddAt adds an error at a specific location.
func (l *ErrorList) AddAt(line, column int, message string) {
	l.errors = append(l.errors, &Error{
		Line:    line,
		Column:  column,
		Message: message,
	})
}

// Errors returns all collected errors.
func (l *ErrorList) Errors() []*Error {
	return l.errors
}

// HasErrors returns true if any errors were collected.
func (l *ErrorList) HasErrors() bool {
	return len(l.errors) > 0
}

// ToParseError converts the error list to a ParseError.
func (l *ErrorList) ToParseError(source string) *ParseError {
	if !l.HasErrors() {
		return nil
	}
	return &ParseError{
		Items:  l.errors,
		Source: source,
	}
}
