package low

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
)

// Parser provides low-level access to the SysML v2 parser.
// It wraps the ANTLR-generated parser with error collection.
type Parser struct {
	parser       *parser.SysMLv2Parser
	lexerErrors  *ErrorCollector
	parserErrors *ErrorCollector
	tokens       *antlr.CommonTokenStream
}

// ParseOption configures parser behavior.
type ParseOption func(*parserConfig)

type parserConfig struct {
	buildParseTree bool
}

// WithParseTree enables parse tree construction.
// Default is true. Set to false for validation-only parsing.
func WithParseTree(build bool) ParseOption {
	return func(c *parserConfig) {
		c.buildParseTree = build
	}
}

// NewParser creates a new parser from the given input string.
func NewParser(input string, opts ...ParseOption) *Parser {
	lexer := NewLexer(input)
	return NewParserFromLexer(lexer, opts...)
}

// NewParserFromBytes creates a new parser from byte slice.
func NewParserFromBytes(input []byte, opts ...ParseOption) *Parser {
	lexer := NewLexerFromBytes(input)
	return NewParserFromLexer(lexer, opts...)
}

// NewParserFromLexer creates a parser from an existing lexer.
// This allows reusing a lexer or inspecting tokens before parsing.
func NewParserFromLexer(lexer *Lexer, opts ...ParseOption) *Parser {
	cfg := &parserConfig{buildParseTree: true}
	for _, opt := range opts {
		opt(cfg)
	}

	tokens := lexer.TokenStream()
	// NOTE: Removed tokens.Fill() - let ANTLR consume tokens lazily
	// This significantly reduces memory usage for large files

	parserErrors := NewErrorCollector("parser")
	p := parser.NewSysMLv2Parser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(parserErrors)
	p.BuildParseTrees = cfg.buildParseTree

	return &Parser{
		parser:       p,
		lexerErrors:  lexer.errors,
		parserErrors: parserErrors,
		tokens:       tokens,
	}
}

// ParseRootNamespace parses the input as a complete SysML root namespace.
// This is the main entry point for parsing SysML files.
func (p *Parser) ParseRootNamespace() parser.IEntryRuleRootNamespaceContext {
	return p.parser.EntryRuleRootNamespace()
}

// ParseExpression parses the input as a SysML expression.
func (p *Parser) ParseExpression() parser.IOwnedExpressionContext {
	return p.parser.OwnedExpression()
}

// Errors returns all errors (lexer and parser) from the parse.
func (p *Parser) Errors() *ParseErrors {
	return &ParseErrors{
		LexerErrors:  p.lexerErrors.Errors(),
		ParserErrors: p.parserErrors.Errors(),
	}
}

// HasErrors returns true if any errors occurred during parsing.
func (p *Parser) HasErrors() bool {
	return p.lexerErrors.HasErrors() || p.parserErrors.HasErrors()
}

// TokenCount returns the total number of tokens in the input.
func (p *Parser) TokenCount() int {
	return len(p.tokens.GetAllTokens())
}

// Tokens returns all tokens from the input.
func (p *Parser) Tokens() []antlr.Token {
	return p.tokens.GetAllTokens()
}

// Inner returns the underlying ANTLR parser for advanced operations.
func (p *Parser) Inner() *parser.SysMLv2Parser {
	return p.parser
}

// Parse is a convenience function that parses input and returns the parse tree and errors.
func Parse(input string, opts ...ParseOption) (parser.IEntryRuleRootNamespaceContext, *ParseErrors) {
	p := NewParser(input, opts...)
	tree := p.ParseRootNamespace()
	return tree, p.Errors()
}

// ParseBytes is like Parse but accepts []byte input.
func ParseBytes(input []byte, opts ...ParseOption) (parser.IEntryRuleRootNamespaceContext, *ParseErrors) {
	p := NewParserFromBytes(input, opts...)
	tree := p.ParseRootNamespace()
	return tree, p.Errors()
}

// Validate checks if the input is valid SysML without building a full parse tree.
func Validate(input string) *ParseErrors {
	p := NewParser(input, WithParseTree(false))
	p.ParseRootNamespace()
	return p.Errors()
}

// ValidateBytes is like Validate but accepts []byte input.
func ValidateBytes(input []byte) *ParseErrors {
	p := NewParserFromBytes(input, WithParseTree(false))
	p.ParseRootNamespace()
	return p.Errors()
}
