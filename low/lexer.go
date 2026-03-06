package low

import (
	"fmt"
	"unsafe"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
)

// Lexer provides low-level access to the SysML v2 lexer.
// It wraps the ANTLR-generated lexer with error collection.
type Lexer struct {
	lexer  *parser.SysMLv2Lexer
	errors *ErrorCollector
}

// NewLexer creates a new lexer from the given input string.
func NewLexer(input string) *Lexer {
	stream := antlr.NewInputStream(input)
	return NewLexerFromStream(stream)
}

// NewLexerFromBytes creates a new lexer from byte slice.
// This avoids a string copy when input is already in []byte form.
func NewLexerFromBytes(input []byte) *Lexer {
	stream := antlr.NewInputStream(bytesToString(input))
	return NewLexerFromStream(stream)
}

// NewLexerFromStream creates a new lexer from an ANTLR CharStream.
// This is the most flexible constructor for advanced use cases.
func NewLexerFromStream(stream antlr.CharStream) *Lexer {
	errors := NewErrorCollector("lexer")
	lexer := parser.NewSysMLv2Lexer(stream)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errors)

	return &Lexer{
		lexer:  lexer,
		errors: errors,
	}
}

// NextToken returns the next token from the input.
// Returns nil when EOF is reached.
func (l *Lexer) NextToken() antlr.Token {
	token := l.lexer.NextToken()
	if token.GetTokenType() == antlr.TokenEOF {
		return nil
	}
	return token
}

// AllTokens returns all tokens from the input.
// This consumes the entire input stream.
func (l *Lexer) AllTokens() []antlr.Token {
	tokens := make([]antlr.Token, 0, 256)
	for {
		token := l.lexer.NextToken()
		if token.GetTokenType() == antlr.TokenEOF {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}

// Errors returns any lexer errors that occurred.
func (l *Lexer) Errors() []*SyntaxError {
	return l.errors.Errors()
}

// HasErrors returns true if lexer errors occurred.
func (l *Lexer) HasErrors() bool {
	return l.errors.HasErrors()
}

// Inner returns the underlying ANTLR lexer for advanced operations.
func (l *Lexer) Inner() *parser.SysMLv2Lexer {
	return l.lexer
}

// TokenStream creates a CommonTokenStream from this lexer.
// The returned stream can be used with the Parser.
func (l *Lexer) TokenStream() *antlr.CommonTokenStream {
	stream := antlr.NewCommonTokenStream(l.lexer, antlr.TokenDefaultChannel)
	return stream
}

// TokenName returns the symbolic name for a token type.
func TokenName(tokenType int) string {
	if tokenType == antlr.TokenEOF {
		return "EOF"
	}
	if tokenType < 0 || tokenType >= len(parser.SysMLv2LexerLexerStaticData.SymbolicNames) {
		return "UNKNOWN"
	}
	name := parser.SysMLv2LexerLexerStaticData.SymbolicNames[tokenType]
	if name == "" {
		// Fall back to literal name if symbolic name is empty
		if tokenType < len(parser.SysMLv2LexerLexerStaticData.LiteralNames) {
			name = parser.SysMLv2LexerLexerStaticData.LiteralNames[tokenType]
		}
		if name == "" {
			return fmt.Sprintf("TOKEN_%d", tokenType)
		}
	}
	return name
}

func bytesToString(input []byte) string {
	if len(input) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(input), len(input))
}
