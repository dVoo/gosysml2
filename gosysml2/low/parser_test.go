package low

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := `
package TestPackage {
    part def Vehicle;
}
`

	tree, errors := Parse(input)

	if errors.HasErrors() {
		t.Fatalf("parse failed: %s", errors)
	}

	if tree == nil {
		t.Fatal("parse tree is nil")
	}
}

func TestParseWithErrors(t *testing.T) {
	input := `
package Broken {
    @@@ invalid syntax
}
`

	_, errors := Parse(input)

	if !errors.HasErrors() {
		t.Error("expected parse errors")
	}

	t.Logf("Got expected errors: %s", errors)
}

func TestValidate(t *testing.T) {
	validInput := `
package Valid {
    part def A;
}
`

	invalidInput := `
package Invalid {
    @@@ not valid
}
`

	// Valid input should pass
	errors := Validate(validInput)
	if errors.HasErrors() {
		t.Errorf("valid input failed validation: %s", errors)
	}

	// Invalid input should fail
	errors = Validate(invalidInput)
	if !errors.HasErrors() {
		t.Error("invalid input should fail validation")
	}
}

func TestLexer(t *testing.T) {
	input := `package P { part x; }`

	lexer := NewLexer(input)
	tokens := lexer.AllTokens()

	if len(tokens) == 0 {
		t.Error("no tokens produced")
	}

	if lexer.HasErrors() {
		t.Errorf("lexer errors: %v", lexer.Errors())
	}

	t.Logf("Token count: %d", len(tokens))
}

func TestTokenName(t *testing.T) {
	// Test a few known token names
	name := TokenName(0)
	if name == "" {
		t.Error("token name should not be empty")
	}
	t.Logf("Token 0 name: %s", name)
}

func TestParser(t *testing.T) {
	input := `package P { part def A; part def B; }`

	parser := NewParser(input)
	tree := parser.ParseRootNamespace()

	if parser.HasErrors() {
		t.Errorf("parse errors: %s", parser.Errors())
	}

	if tree == nil {
		t.Error("parse tree is nil")
	}

	t.Logf("Token count: %d", parser.TokenCount())
}

func TestParserFromLexer(t *testing.T) {
	input := `package P { part x : A; }`

	lexer := NewLexer(input)
	parser := NewParserFromLexer(lexer)

	tree := parser.ParseRootNamespace()

	if parser.HasErrors() {
		t.Errorf("parse errors: %s", parser.Errors())
	}

	if tree == nil {
		t.Error("parse tree is nil")
	}
}

func TestParseBytes(t *testing.T) {
	input := []byte(`package ByteTest { part def X; }`)

	tree, errors := ParseBytes(input)

	if errors.HasErrors() {
		t.Fatalf("parse failed: %s", errors)
	}

	if tree == nil {
		t.Error("parse tree is nil")
	}
}

func TestErrorCollector(t *testing.T) {
	ec := NewErrorCollector("test")

	if ec.HasErrors() {
		t.Error("new collector should have no errors")
	}

	ec.SyntaxError(nil, nil, 1, 5, "test error", nil)

	if !ec.HasErrors() {
		t.Error("collector should have errors after adding one")
	}

	errors := ec.Errors()
	if len(errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(errors))
	}

	if errors[0].Line != 1 || errors[0].Column != 5 {
		t.Errorf("error location mismatch: got %d:%d", errors[0].Line, errors[0].Column)
	}

	ec.Clear()
	if ec.HasErrors() {
		t.Error("collector should be empty after clear")
	}
}

func TestParseErrors(t *testing.T) {
	pe := &ParseErrors{
		LexerErrors: []*SyntaxError{
			{Line: 1, Column: 1, Message: "lexer error", Source: "lexer"},
		},
		ParserErrors: []*SyntaxError{
			{Line: 2, Column: 2, Message: "parser error", Source: "parser"},
		},
	}

	if !pe.HasErrors() {
		t.Error("should have errors")
	}

	all := pe.All()
	if len(all) != 2 {
		t.Errorf("expected 2 errors, got %d", len(all))
	}

	errStr := pe.Error()
	if errStr == "" {
		t.Error("error string should not be empty")
	}

	t.Logf("Error string: %s", errStr)
}

func TestWithParseTree(t *testing.T) {
	input := `package P { part def A; }`

	// With parse tree (default)
	p1 := NewParser(input, WithParseTree(true))
	tree1 := p1.ParseRootNamespace()
	if tree1 == nil {
		t.Error("tree should not be nil with parse tree enabled")
	}

	// Without parse tree (validation only)
	p2 := NewParser(input, WithParseTree(false))
	p2.ParseRootNamespace()
	if p2.HasErrors() {
		t.Errorf("validation failed: %s", p2.Errors())
	}
}
