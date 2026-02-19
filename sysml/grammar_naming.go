package sysml

import "github.com/dVoo/gosysml2/low"

// NormalizeGrammarName normalizes grammar symbol names for cross-source comparison.
func NormalizeGrammarName(name string) string {
	return low.NormalizeGrammarName(name)
}

// ResolveANTLRParserRuleName resolves a KEBNF/ANTLR-like rule name to an ANTLR parser rule.
func ResolveANTLRParserRuleName(name string) (string, bool) {
	return low.ResolveParserRuleName(name)
}

// ResolveANTLRTokenName resolves a KEBNF/ANTLR-like token name to an ANTLR symbolic token.
func ResolveANTLRTokenName(name string) (string, bool) {
	return low.ResolveTokenName(name)
}

// ParserRuleNameCandidates returns all ANTLR parser rules matching a normalized input.
func ParserRuleNameCandidates(name string) []string {
	return low.ParserRuleNameCandidates(name)
}

// TokenNameCandidates returns all ANTLR symbolic token names matching a normalized input.
func TokenNameCandidates(name string) []string {
	return low.TokenNameCandidates(name)
}
