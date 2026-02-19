package low

import "testing"

func TestNormalizeGrammarName(t *testing.T) {
	if got := NormalizeGrammarName("  Requirement_Body  "); got != "requirementbody" {
		t.Fatalf("expected requirementbody, got %q", got)
	}
	if got := NormalizeGrammarName("::>"); got != "" {
		t.Fatalf("expected empty for symbols-only input, got %q", got)
	}
}

func TestResolveParserRuleName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Package", "package"},
		{"Comment", "comment"},
		{"Expression", "expression"},
		{"TypedBy", "typedBy"},
		{"RequirementBody", "requirementBody"},
	}
	for _, tc := range tests {
		got, ok := ResolveParserRuleName(tc.in)
		if !ok {
			t.Fatalf("expected to resolve %q", tc.in)
		}
		if got != tc.want {
			t.Fatalf("expected %q -> %q, got %q", tc.in, tc.want, got)
		}
	}

	if _, ok := ResolveParserRuleName("DefinitelyNotARule"); ok {
		t.Fatal("expected unknown parser rule resolution to fail")
	}
}

func TestResolveTokenName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"requirement", "REQUIREMENT"},
		{"REDEFINES", "REDEFINES"},
		{"single_line_note", "SINGLE_LINE_NOTE"},
	}
	for _, tc := range tests {
		got, ok := ResolveTokenName(tc.in)
		if !ok {
			t.Fatalf("expected to resolve %q", tc.in)
		}
		if got != tc.want {
			t.Fatalf("expected %q -> %q, got %q", tc.in, tc.want, got)
		}
	}

	if _, ok := ResolveTokenName("DefinitelyNotAToken"); ok {
		t.Fatal("expected unknown token resolution to fail")
	}
}
