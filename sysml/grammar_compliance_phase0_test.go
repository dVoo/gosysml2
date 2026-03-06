package sysml

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
)

type grammarCoverageClassification struct {
	IntentionallyIgnored []string `json:"intentionally_ignored"`
}

func TestGrammarEnterRuleCoverageClassification(t *testing.T) {
	allRules := collectParserEnterRules()
	implementedRules, err := collectModelBuilderEnterRules()
	if err != nil {
		t.Fatalf("collect modelBuilder rules: %v", err)
	}

	classPath, err := findCoverageClassificationPath()
	if err != nil {
		t.Fatalf("locate grammar coverage classification file: %v", err)
	}

	classification, err := loadCoverageClassification(classPath)
	if err != nil {
		t.Fatalf("load classification from %s: %v", classPath, err)
	}

	ignored := make(map[string]struct{}, len(classification.IntentionallyIgnored))
	for _, name := range classification.IntentionallyIgnored {
		ignored[name] = struct{}{}
	}

	var missing []string
	implementedCount := 0
	ignoredCount := 0
	for rule := range allRules {
		if _, ok := implementedRules[rule]; ok {
			implementedCount++
			continue
		}
		if _, ok := ignored[rule]; ok {
			ignoredCount++
			continue
		}
		missing = append(missing, rule)
	}

	// Guard against stale entries in the classification file.
	var stale []string
	for name := range ignored {
		if _, ok := allRules[name]; !ok {
			stale = append(stale, name)
		}
	}
	if len(stale) > 0 {
		sort.Strings(stale)
		t.Fatalf("classification contains unknown parser rules: %s", strings.Join(stale, ", "))
	}

	totalRules := len(allRules)
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf(
			"unclassified parser Enter* rules: %d/%d\nAdd each rule to modelBuilder or sysml/testdata/grammar_coverage_classification.json\n%s",
			len(missing),
			totalRules,
			strings.Join(missing, "\n"),
		)
	}

	t.Logf("grammar coverage summary: total=%d implemented=%d intentionally_ignored=%d missing=%d",
		totalRules, implementedCount, ignoredCount, len(missing))
}

func TestRepresentativeValidationParseTreeAndModelExtraction(t *testing.T) {
	testCases := []struct {
		path        string
		expectKinds []ElementKind
	}{
		{
			path:        "../validationdata/01-Parts Tree/1a-Parts Tree.sysml",
			expectKinds: []ElementKind{KindPackage, KindPart},
		},
		{
			path:        "../validationdata/08-Requirements/8-Requirements.sysml",
			expectKinds: []ElementKind{KindPackage, KindRequirement},
		},
		{
			path:        "../validationdata/09-Verification/9-Verification-simplified.sysml",
			expectKinds: []ElementKind{KindPackage, KindVerification},
		},
		{
			path:        "../validationdata/11-View and Viewpoint/11a-View-Viewpoint.sysml",
			expectKinds: []ElementKind{KindPackage, KindView, KindViewpoint},
		},
	}

	for _, tc := range testCases {
		t.Run(filepath.Base(tc.path), func(t *testing.T) {
			path := externalTestPath(t, tc.path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Skipf("validation file not found: %s", path)
			}

			result := ParseFile(path)
			if !result.Success() {
				t.Fatalf("parse failed for %s: %v", tc.path, result.Err())
			}
			if result.Tree == nil {
				t.Fatalf("parse tree is nil for %s", tc.path)
			}
			if countParseTreeNodes(result.Tree) == 0 {
				t.Fatalf("parse tree node count is 0 for %s", tc.path)
			}
			if result.Model == nil {
				t.Fatalf("model is nil for %s", tc.path)
			}

			for _, kind := range tc.expectKinds {
				elems := FindByKind(result.Model, kind)
				if len(elems) == 0 {
					t.Fatalf("expected at least one %s element in %s", kind.String(), tc.path)
				}
			}
		})
	}
}

func collectParserEnterRules() map[string]struct{} {
	rules := make(map[string]struct{})
	listenerType := reflect.TypeOf((*parser.SysMLv2ParserListener)(nil)).Elem()
	for i := 0; i < listenerType.NumMethod(); i++ {
		name := listenerType.Method(i).Name
		if strings.HasPrefix(name, "Enter") && name != "EnterEveryRule" {
			rules[name] = struct{}{}
		}
	}
	return rules
}

func collectModelBuilderEnterRules() (map[string]struct{}, error) {
	content, err := os.ReadFile("parse.go")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`func \(b \*modelBuilder\) (Enter[A-Za-z0-9_]+)\(`)
	matches := re.FindAllStringSubmatch(string(content), -1)
	rules := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		rules[m[1]] = struct{}{}
	}
	return rules, nil
}

func findCoverageClassificationPath() (string, error) {
	candidates := []string{
		filepath.Join("testdata", "grammar_coverage_classification.json"),
		filepath.Join("sysml", "testdata", "grammar_coverage_classification.json"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", os.ErrNotExist
}

func loadCoverageClassification(path string) (*grammarCoverageClassification, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c grammarCoverageClassification
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func countParseTreeNodes(tree antlr.Tree) int {
	if tree == nil {
		return 0
	}
	count := 1
	for i := 0; i < tree.GetChildCount(); i++ {
		child := tree.GetChild(i)
		if childTree, ok := child.(antlr.Tree); ok {
			count += countParseTreeNodes(childTree)
		}
	}
	return count
}
