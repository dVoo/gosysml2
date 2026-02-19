package low

import (
	"sort"
	"strings"
	"sync"

	"github.com/dVoo/gosysml2/internal/parser"
)

var (
	namingOnce               sync.Once
	parserRuleExactLookup    map[string]string
	parserRuleCanonicalIndex map[string][]string
	tokenExactLookup         map[string]string
	tokenCanonicalIndex      map[string][]string
)

func initNamingIndexes() {
	namingOnce.Do(func() {
		parser.SysMLv2ParserInit()
		parser.SysMLv2LexerInit()

		parserRuleExactLookup = make(map[string]string)
		parserRuleCanonicalIndex = make(map[string][]string)
		tokenExactLookup = make(map[string]string)
		tokenCanonicalIndex = make(map[string][]string)

		for _, rule := range parser.SysMLv2ParserParserStaticData.RuleNames {
			if rule == "" {
				continue
			}
			parserRuleExactLookup[rule] = rule
			key := NormalizeGrammarName(rule)
			parserRuleCanonicalIndex[key] = append(parserRuleCanonicalIndex[key], rule)
		}
		for _, sym := range parser.SysMLv2LexerLexerStaticData.SymbolicNames {
			if sym == "" {
				continue
			}
			tokenExactLookup[sym] = sym
			key := NormalizeGrammarName(sym)
			tokenCanonicalIndex[key] = append(tokenCanonicalIndex[key], sym)
		}
		for k := range parserRuleCanonicalIndex {
			sort.Strings(parserRuleCanonicalIndex[k])
		}
		for k := range tokenCanonicalIndex {
			sort.Strings(tokenCanonicalIndex[k])
		}
	})
}

// NormalizeGrammarName normalizes KEBNF/ANTLR names to a shared comparison form.
// It removes non-alphanumeric characters and lowercases the result.
func NormalizeGrammarName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		if r >= 'A' && r <= 'Z' {
			b.WriteRune(r + ('a' - 'A'))
			continue
		}
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ResolveParserRuleName resolves a KEBNF/ANTLR-like name to a concrete ANTLR parser rule name.
// It supports case/underscore differences (e.g. "Package" -> "package_").
func ResolveParserRuleName(name string) (string, bool) {
	initNamingIndexes()
	if exact, ok := parserRuleExactLookup[name]; ok {
		return exact, true
	}
	candidates := parserRuleCanonicalIndex[NormalizeGrammarName(name)]
	if len(candidates) == 1 {
		return candidates[0], true
	}
	return "", false
}

// ParserRuleNameCandidates returns all ANTLR parser rule names matching a normalized input.
func ParserRuleNameCandidates(name string) []string {
	initNamingIndexes()
	candidates := parserRuleCanonicalIndex[NormalizeGrammarName(name)]
	out := make([]string, len(candidates))
	copy(out, candidates)
	return out
}

// ResolveTokenName resolves a KEBNF/ANTLR-like token name to a concrete ANTLR symbolic token name.
func ResolveTokenName(name string) (string, bool) {
	initNamingIndexes()
	if exact, ok := tokenExactLookup[name]; ok {
		return exact, true
	}
	candidates := tokenCanonicalIndex[NormalizeGrammarName(name)]
	if len(candidates) == 1 {
		return candidates[0], true
	}
	return "", false
}

// TokenNameCandidates returns all ANTLR symbolic token names matching a normalized input.
func TokenNameCandidates(name string) []string {
	initNamingIndexes()
	candidates := tokenCanonicalIndex[NormalizeGrammarName(name)]
	out := make([]string, len(candidates))
	copy(out, candidates)
	return out
}
