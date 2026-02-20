package sysml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/antlr4-go/antlr/v4"
	"github.com/dVoo/gosysml2/internal/parser"
	"github.com/dVoo/gosysml2/low"
)

// childAdder is implemented by elements that can hold children.
type childAdder interface {
	AddChild(child Element)
}

// locationFromContext extracts a Location from an ANTLR context.
// This helper eliminates repetitive location extraction code in Enter* methods.
func locationFromContext(ctx interface {
	GetStart() antlr.Token
	GetStop() antlr.Token
}) Location {
	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}
	return loc
}

// ParseResult contains the result of parsing SysML input.
type ParseResult struct {
	Model  *Model      // The parsed model (nil if parsing failed completely)
	Errors *ParseError // Any errors that occurred (nil if successful)
	Tree   antlr.Tree  // The raw parse tree (for advanced use)
	Source string      // Source file path or identifier
}

// Success returns true if parsing was successful (no errors).
func (r *ParseResult) Success() bool {
	return r.Errors == nil || !r.Errors.HasErrors()
}

// ParseOption configures parsing behavior.
type ParseOption func(*parseConfig)

type parseConfig struct {
	discardTree       bool
	libraryRegistry   *LibraryRegistry
	autoLoadLibraries bool
	libraryPath       string
}

type parseRewriteHints struct {
	requireConstraintExprByPlaceholder map[string]string
	requirementBindingsByLine          map[int]map[string]string
	requirementBindingsByName          map[string][]map[string]string
}

func newParseRewriteHints() *parseRewriteHints {
	return &parseRewriteHints{
		requireConstraintExprByPlaceholder: make(map[string]string),
		requirementBindingsByLine:          make(map[int]map[string]string),
		requirementBindingsByName:          make(map[string][]map[string]string),
	}
}

var (
	requireBlockPattern = regexp.MustCompile(`(?s)require\s*\{\s*(.*?)\s*\}\s*;`)
	reqBindingPattern   = regexp.MustCompile(`(?m)(requirement\b[^\n\r;{]*?:\s*[^\[\n\r;{]+?)\s*\[([^\]\n\r]*)\](\s*(?:\{|;))`)
	identTokenPattern   = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)
	inRefLambdaPattern  = regexp.MustCompile(`\{\s*in\s+ref\s+([A-Za-z_][A-Za-z0-9_]*)\s*\{`)

	attributeKeywordNamePattern      = regexp.MustCompile(`(?m)\battribute\s+(type)\s*:`)
	aliasKeywordNamePattern          = regexp.MustCompile(`(?m)\balias\s+(multiplicity)\s+for\b`)
	attributeShortNameKeywordPattern = regexp.MustCompile(`(?m)\battribute\s*<\s*(var)\s*>`)
	refVarNamePattern                = regexp.MustCompile(`(?m)\bref\s+var(\s*(?:\[|:|::|:>|:>>|;|=))`)
	assignVarTargetPattern           = regexp.MustCompile(`(?m)\bassign\s+var(\s*:=)`)
)

func lineNumberAtOffset(input string, offset int) int {
	if offset <= 0 {
		return 1
	}
	if offset > len(input) {
		offset = len(input)
	}
	return 1 + strings.Count(input[:offset], "\n")
}

func parseBindings(bindingsText string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(bindingsText, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		if key != "" {
			result[key] = value
		}
	}
	return result
}

func extractRequirementUsageName(requirementPrefix string) string {
	colon := strings.LastIndex(requirementPrefix, ":")
	if colon < 0 {
		return ""
	}
	left := requirementPrefix[:colon]
	tokens := identTokenPattern.FindAllString(left, -1)
	if len(tokens) == 0 {
		return ""
	}
	// "requirement def X : T" is a definition form, not usage.
	if len(tokens) >= 2 && tokens[0] == "requirement" && tokens[1] == "def" {
		return ""
	}
	// For usage forms, the last identifier before ':' is the declared name.
	return tokens[len(tokens)-1]
}

func normalizeUnsupportedRequirementSyntax(input string) (string, *parseRewriteHints) {
	hints := newParseRewriteHints()
	normalized := input

	// Gap 16 compatibility: lambda parameters written as "{in ref a { ... }}"
	// are normalized to "{in a { ... }}" for current grammar support.
	normalized = inRefLambdaPattern.ReplaceAllString(normalized, `{in $1 {`)

	// Gap 19/20 compatibility: allow reserved keywords as declared names
	// in forms seen in standard library models.
	normalized = attributeKeywordNamePattern.ReplaceAllString(normalized, `attribute '$1' :`)
	normalized = aliasKeywordNamePattern.ReplaceAllString(normalized, `alias '$1' for`)
	normalized = attributeShortNameKeywordPattern.ReplaceAllString(normalized, `attribute <'$1'>`)
	normalized = refVarNamePattern.ReplaceAllString(normalized, `ref 'var'$1`)
	normalized = assignVarTargetPattern.ReplaceAllString(normalized, `assign 'var'$1`)

	// Gap 14 compatibility: "require { expr };" -> "require __gap14_constraint_N;"
	matches := requireBlockPattern.FindAllStringSubmatchIndex(normalized, -1)
	if len(matches) > 0 {
		var out strings.Builder
		last := 0
		for i, match := range matches {
			start, end := match[0], match[1]
			exprStart, exprEnd := match[2], match[3]
			out.WriteString(normalized[last:start])
			placeholder := fmt.Sprintf("__gap14_constraint_%d", i+1)
			out.WriteString("require " + placeholder + ";")
			hints.requireConstraintExprByPlaceholder[placeholder] = strings.TrimSpace(normalized[exprStart:exprEnd])
			last = end
		}
		out.WriteString(normalized[last:])
		normalized = out.String()
	}

	// Gap 15 compatibility: remove usage binding list syntax and capture values.
	matches = reqBindingPattern.FindAllStringSubmatchIndex(normalized, -1)
	if len(matches) > 0 {
		var out strings.Builder
		last := 0
		for _, match := range matches {
			start, end := match[0], match[1]
			prefixStart, prefixEnd := match[2], match[3]
			bindingsStart, bindingsEnd := match[4], match[5]
			suffixStart, suffixEnd := match[6], match[7]

			prefix := normalized[prefixStart:prefixEnd]
			bindingsText := normalized[bindingsStart:bindingsEnd]
			suffix := normalized[suffixStart:suffixEnd]

			out.WriteString(normalized[last:start])
			out.WriteString(prefix)
			out.WriteString(suffix)

			bindings := parseBindings(bindingsText)
			if len(bindings) > 0 {
				line := lineNumberAtOffset(normalized, start)
				hints.requirementBindingsByLine[line] = bindings
				if name := extractRequirementUsageName(prefix); name != "" {
					hints.requirementBindingsByName[name] = append(hints.requirementBindingsByName[name], bindings)
				}
			}
			last = end
		}
		out.WriteString(normalized[last:])
		normalized = out.String()
	}

	return normalized, hints
}

// WithDiscardTree discards the parse tree after building the model.
// This significantly reduces memory usage for large files.
func WithDiscardTree() ParseOption {
	return func(c *parseConfig) {
		c.discardTree = true
	}
}

// WithLibraryRegistry uses an existing library registry for resolving imports.
// The registry should already be loaded with the standard libraries.
func WithLibraryRegistry(reg *LibraryRegistry) ParseOption {
	return func(c *parseConfig) {
		c.libraryRegistry = reg
	}
}

// WithStandardLibrary auto-loads the standard SysML library before parsing.
// The library will be loaded from the default path (./libraries/sysml.library).
func WithStandardLibrary() ParseOption {
	return func(c *parseConfig) {
		c.autoLoadLibraries = true
	}
}

// WithLibraryPath specifies a custom path for loading standard libraries.
// Use with WithStandardLibrary() to load from a non-default location.
func WithLibraryPath(path string) ParseOption {
	return func(c *parseConfig) {
		c.libraryPath = path
	}
}

// getOrCreateRegistry returns a library registry based on configuration.
// If a registry is provided in config, it returns that.
// If autoLoadLibraries is set, it creates and populates a new registry.
// Otherwise returns nil.
func getOrCreateRegistry(cfg *parseConfig) *LibraryRegistry {
	// If registry already provided, use it
	if cfg.libraryRegistry != nil {
		return cfg.libraryRegistry
	}

	// If auto-loading is enabled, create and load registry
	if cfg.autoLoadLibraries {
		var opts []LibraryOption
		if cfg.libraryPath != "" {
			opts = append(opts, WithLibraryPaths(cfg.libraryPath))
		}
		registry := NewLibraryRegistry(opts...)
		// Load standard libraries (ignore errors - libraries may not be present)
		_ = registry.RegisterStandardLibrary()
		return registry
	}

	return nil
}

// parseWithSource parses SysML input with a specified source identifier.
func parseWithSource(input, source string, opts ...ParseOption) (result *ParseResult) {
	cfg := &parseConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	result = &ParseResult{
		Source: source,
	}

	// API-2 ergonomics: parser panics are converted into regular parse errors.
	defer func() {
		if r := recover(); r != nil {
			result.Model = nil
			result.Tree = nil
			result.Errors = &ParseError{
				Errors: []*Error{{
					Line:    0,
					Column:  0,
					Message: fmt.Sprintf("parser panic: %v", r),
					Context: "panic",
				}},
				Source: source,
				Input:  input,
			}
		}
	}()

	// Get or create library registry
	registry := getOrCreateRegistry(cfg)
	normalizedInput := input
	var rewriteHints *parseRewriteHints
	if strings.EqualFold(filepath.Ext(source), ".kerml") {
		// KerML files are parsed directly without SysML compatibility rewrites.
		rewriteHints = newParseRewriteHints()
	} else {
		normalizedInput, rewriteHints = normalizeUnsupportedRequirementSyntax(input)
	}

	// Use low-level parser
	parseOpts := make([]low.ParseOption, 0, 1)
	if strings.EqualFold(filepath.Ext(source), ".kerml") {
		parseOpts = append(parseOpts, low.WithSyntaxMode(low.SyntaxModeKerML))
	} else {
		parseOpts = append(parseOpts, low.WithSyntaxMode(low.SyntaxModeSysML))
	}
	tree, lowErrors := low.Parse(normalizedInput, parseOpts...)

	if !cfg.discardTree {
		result.Tree = tree
	}

	// Convert errors
	if lowErrors.HasErrors() {
		result.Errors = convertFromLowLevel(lowErrors, source)
	}

	// Build model from parse tree
	if tree != nil {
		result.Model = buildModel(tree, registry, rewriteHints)
	}

	return result
}

// ParseString parses a SysML string and returns a high-level model.
func ParseString(input string, opts ...ParseOption) *ParseResult {
	return parseWithSource(input, "<string>", opts...)
}

// ParseFile parses a SysML file and returns a high-level model.
func ParseFile(filename string, opts ...ParseOption) *ParseResult {
	content, err := os.ReadFile(filename)
	if err != nil {
		return &ParseResult{
			Errors: &ParseError{
				Errors: []*Error{{
					Line:    0,
					Column:  0,
					Message: fmt.Errorf("reading %s: %w", filename, err).Error(),
				}},
				Source: filename,
			},
			Source: filename,
		}
	}
	return parseWithSource(string(content), filename, opts...)
}

// ParseBytes parses SysML from a byte slice.
// This is more efficient than ParseString when you already have []byte.
func ParseBytes(input []byte, source string, opts ...ParseOption) *ParseResult {
	return parseWithSource(string(input), source, opts...)
}

// ParseReader parses SysML from an io.Reader.
// Useful when reading from network streams or other io sources.
func ParseReader(r io.Reader, source string, opts ...ParseOption) *ParseResult {
	content, err := io.ReadAll(r)
	if err != nil {
		return &ParseResult{
			Errors: &ParseError{
				Errors: []*Error{{
					Line:    0,
					Column:  0,
					Message: fmt.Errorf("reading from %s: %w", source, err).Error(),
				}},
				Source: source,
			},
			Source: source,
		}
	}
	return parseWithSource(string(content), source, opts...)
}

func modelOrParseError(result *ParseResult, source string) (*Model, error) {
	if result == nil {
		return nil, fmt.Errorf("parse failed for %s: nil result", source)
	}
	if !result.Success() {
		if result.Errors != nil {
			return nil, result.Errors
		}
		return nil, fmt.Errorf("parse failed for %s", source)
	}
	return result.Model, nil
}

// ParseStringModel parses SysML text and returns model/error in idiomatic Go style.
func ParseStringModel(input string, opts ...ParseOption) (*Model, error) {
	result := ParseString(input, opts...)
	return modelOrParseError(result, "<string>")
}

// ParseFileModel parses a SysML file and returns model/error in idiomatic Go style.
func ParseFileModel(filename string, opts ...ParseOption) (*Model, error) {
	result := ParseFile(filename, opts...)
	return modelOrParseError(result, filename)
}

// ParseDirectory parses all .sysml and .kerml files in a directory.
// Use WithDiscardTree() option for large repositories to reduce memory.
func ParseDirectory(dir string, opts ...ParseOption) ([]*ParseResult, error) {
	var results []*ParseResult

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}
		ext := strings.ToLower(filepath.Ext(path))
		if d.IsDir() || (ext != ".sysml" && ext != ".kerml") {
			return nil
		}

		result := ParseFile(path, opts...)
		results = append(results, result)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", dir, err)
	}

	return results, nil
}

// ParseDirectoryParallel parses all .sysml and .kerml files in a directory using multiple workers.
// This can significantly speed up parsing of large repositories on multi-core machines.
// Set workers to 0 to use runtime.NumCPU().
func ParseDirectoryParallel(dir string, workers int, opts ...ParseOption) ([]*ParseResult, error) {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	// Collect files first
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}
		ext := strings.ToLower(filepath.Ext(path))
		if d.IsDir() || (ext != ".sysml" && ext != ".kerml") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory %s: %w", dir, err)
	}

	results := make([]*ParseResult, len(files))
	var wg sync.WaitGroup
	sem := make(chan struct{}, workers)

	for i, file := range files {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore

		go func(idx int, path string) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore

			results[idx] = ParseFile(path, opts...)
		}(i, file)
	}

	wg.Wait()
	return results, nil
}

// ParseDirectoryStream parses files one at a time and calls handler for each.
// This is the most memory-efficient way to process large repositories.
// The handler can process/store results and allow the GC to reclaim memory.
func ParseDirectoryStream(dir string, handler func(*ParseResult) error, opts ...ParseOption) error {
	// Always discard tree for streaming to minimize memory
	opts = append(opts, WithDiscardTree())

	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}
		ext := strings.ToLower(filepath.Ext(path))
		if d.IsDir() || (ext != ".sysml" && ext != ".kerml") {
			return nil
		}

		result := ParseFile(path, opts...)
		if err := handler(result); err != nil {
			return err
		}

		// Help GC by clearing references
		result.Model = nil
		result.Tree = nil

		return nil
	})
}

// buildModel converts the parse tree to a high-level Model.
// If registry is provided, it will be used for resolving library imports and qualified names.
func buildModel(tree parser.IEntryRuleRootNamespaceContext, registry *LibraryRegistry, rewriteHints *parseRewriteHints) *Model {
	if tree == nil {
		return nil
	}

	model := NewModel()
	builder := &modelBuilder{
		model:           model,
		libraryRegistry: registry,
		elementStack:    make([]Element, 0, 16),
		packageStack:    make([]*Package, 0, 8),
		rewriteHints:    rewriteHints,
	}
	antlr.ParseTreeWalkerDefault.Walk(builder, tree)

	// Set library registry on model for reference resolution
	if registry != nil {
		model.SetLibraryRegistry(registry)
	}

	// Build index and resolve references
	model.BuildIndex()
	model.ResolveReferences()

	return model
}

// modelBuilder walks the parse tree and builds the model.
type modelBuilder struct {
	*parser.BaseSysMLv2ParserListener
	model           *Model
	libraryRegistry *LibraryRegistry // For resolving library imports
	currentPkg      *Package
	packageStack    []*Package // Stack of packages for nested package handling
	elementStack    []Element  // Stack of elements (parts, requirements, etc.) for parent tracking
	rewriteHints    *parseRewriteHints
}

// getCurrentParent returns the current parent element for adding children.
// It checks the element stack first (for nested elements like requirements/parts),
// then falls back to currentPkg (for package-level elements).
// Returns nil if there's no valid parent (top-level elements).
func (b *modelBuilder) getCurrentParent() Element {
	if len(b.elementStack) > 0 {
		return b.elementStack[len(b.elementStack)-1]
	}
	// Return nil if currentPkg is nil to avoid interface nil pointer issues
	if b.currentPkg == nil {
		return nil
	}
	return b.currentPkg
}

// addToParent adds an element to the current parent with proper type handling.
// This helper eliminates ~200 lines of repetitive switch statements in Enter* methods.
func (b *modelBuilder) addToParent(elem Element) {
	parent := b.getCurrentParent()
	if parent == nil {
		return
	}
	if container, ok := parent.(childAdder); ok {
		container.AddChild(elem)
	}
}

func (b *modelBuilder) addOccurrenceToModel(occ *Occurrence) {
	occ.parent = b.getCurrentParent()
	b.addToParent(occ)
	b.model.AddOccurrence(occ)
}

func (b *modelBuilder) EnterPackage(ctx *parser.PackageContext) {
	name := ""
	if ctx.PackageDeclaration() != nil {
		if ident := ctx.PackageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	pkg := NewPackage(name, locationFromContext(ctx))

	if b.currentPkg != nil {
		pkg.parent = b.currentPkg
		b.currentPkg.AddChild(pkg)
	} else {
		b.model.AddPackage(pkg)
	}

	// Push current package to package stack and element stack
	b.packageStack = append(b.packageStack, b.currentPkg)
	b.elementStack = append(b.elementStack, pkg)
	b.currentPkg = pkg
}

func (b *modelBuilder) ExitPackage(ctx *parser.PackageContext) {
	// Pop from element stack
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}

	// Restore previous package from package stack
	if len(b.packageStack) > 0 {
		b.currentPkg = b.packageStack[len(b.packageStack)-1]
		b.packageStack = b.packageStack[:len(b.packageStack)-1]
	} else {
		b.currentPkg = nil
	}
}

// EnterLibraryPackage handles standard library package declarations.
func (b *modelBuilder) EnterLibraryPackage(ctx *parser.LibraryPackageContext) {
	name := ""
	if ctx.PackageDeclaration() != nil {
		if ident := ctx.PackageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	pkg := NewPackage(name, locationFromContext(ctx))
	pkg.IsLibrary = true // Mark as library package

	if b.currentPkg != nil {
		pkg.parent = b.currentPkg
		b.currentPkg.AddChild(pkg)
	} else {
		b.model.AddPackage(pkg)
	}

	// Push current package to package stack and element stack
	b.packageStack = append(b.packageStack, b.currentPkg)
	b.elementStack = append(b.elementStack, pkg)
	b.currentPkg = pkg
}

// ExitLibraryPackage handles exiting library package scope.
func (b *modelBuilder) ExitLibraryPackage(ctx *parser.LibraryPackageContext) {
	// Pop from element stack
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}

	// Restore previous package from package stack
	if len(b.packageStack) > 0 {
		b.currentPkg = b.packageStack[len(b.packageStack)-1]
		b.packageStack = b.packageStack[:len(b.packageStack)-1]
	} else {
		b.currentPkg = nil
	}
}

// EnterNamespace_ handles KerML namespace declarations.
func (b *modelBuilder) EnterNamespace_(ctx *parser.Namespace_Context) {
	name := ""
	if decl := ctx.NamespaceDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	ns := NewPackage(name, locationFromContext(ctx))
	parentPkg := b.currentPkg
	if parentPkg != nil {
		ns.parent = parentPkg
		parentPkg.AddChild(ns)
	} else {
		b.model.AddPackage(ns)
	}

	b.packageStack = append(b.packageStack, b.currentPkg)
	b.elementStack = append(b.elementStack, ns)
	b.currentPkg = ns
}

// ExitNamespace_ restores package scope when leaving a KerML namespace.
func (b *modelBuilder) ExitNamespace_(ctx *parser.Namespace_Context) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
	if len(b.packageStack) > 0 {
		b.currentPkg = b.packageStack[len(b.packageStack)-1]
		b.packageStack = b.packageStack[:len(b.packageStack)-1]
	} else {
		b.currentPkg = nil
	}
}

func (b *modelBuilder) enterKerMLType(name, shortName, keyword string, loc Location, superRefs []string) {
	elem := NewKerMLType(name, keyword, loc)
	elem.setDeclaredShortName(shortName)
	elem.parent = b.getCurrentParent()
	for _, ref := range superRefs {
		elem.AddUnresolvedSuper(ref)
	}

	if parent := b.getCurrentParent(); parent != nil {
		b.addToParent(elem)
	} else {
		b.model.Elements = append(b.model.Elements, elem)
	}
	b.elementStack = append(b.elementStack, elem)
}

func (b *modelBuilder) exitKerMLType() {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterType_(ctx *parser.Type_Context) {
	name, shortName := "", ""
	superRefs := make([]string, 0)
	if decl := ctx.TypeDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name, shortName = extractIdentificationNames(ident)
		}
		superRefs = extractTypeDeclarationSpecializations(decl)
	}
	b.enterKerMLType(name, shortName, "type", locationFromContext(ctx), superRefs)
}

func (b *modelBuilder) ExitType_(ctx *parser.Type_Context) { b.exitKerMLType() }

func (b *modelBuilder) EnterClassifier(ctx *parser.ClassifierContext) {
	b.enterKerMLClassifierFromDecl("classifier", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitClassifier(ctx *parser.ClassifierContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterDataType(ctx *parser.DataTypeContext) {
	b.enterKerMLClassifierFromDecl("datatype", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitDataType(ctx *parser.DataTypeContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterClass(ctx *parser.ClassContext) {
	b.enterKerMLClassifierFromDecl("class", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitClass(ctx *parser.ClassContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterStructure(ctx *parser.StructureContext) {
	b.enterKerMLClassifierFromDecl("struct", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitStructure(ctx *parser.StructureContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterMetaclass(ctx *parser.MetaclassContext) {
	b.enterKerMLClassifierFromDecl("metaclass", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitMetaclass(ctx *parser.MetaclassContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterAssociation(ctx *parser.AssociationContext) {
	b.enterKerMLClassifierFromDecl("assoc", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitAssociation(ctx *parser.AssociationContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterAssociationStructure(ctx *parser.AssociationStructureContext) {
	b.enterKerMLClassifierFromDecl("assoc struct", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitAssociationStructure(ctx *parser.AssociationStructureContext) {
	b.exitKerMLType()
}

func (b *modelBuilder) EnterInteraction(ctx *parser.InteractionContext) {
	b.enterKerMLClassifierFromDecl("interaction", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitInteraction(ctx *parser.InteractionContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterBehavior(ctx *parser.BehaviorContext) {
	b.enterKerMLClassifierFromDecl("behavior", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitBehavior(ctx *parser.BehaviorContext) { b.exitKerMLType() }

func (b *modelBuilder) EnterFunction_(ctx *parser.Function_Context) {
	b.enterKerMLClassifierFromDecl("function", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitFunction_(ctx *parser.Function_Context) { b.exitKerMLType() }

func (b *modelBuilder) EnterPredicate(ctx *parser.PredicateContext) {
	b.enterKerMLClassifierFromDecl("predicate", locationFromContext(ctx), ctx.ClassifierDeclaration())
}

func (b *modelBuilder) ExitPredicate(ctx *parser.PredicateContext) { b.exitKerMLType() }

func (b *modelBuilder) enterKerMLClassifierFromDecl(keyword string, loc Location, decl parser.IClassifierDeclarationContext) {
	name, shortName := "", ""
	superRefs := make([]string, 0)
	if decl != nil {
		if ident := decl.Identification(); ident != nil {
			name, shortName = extractIdentificationNames(ident)
		}
		superRefs = extractClassifierSpecializations(decl)
	}
	b.enterKerMLType(name, shortName, keyword, loc, superRefs)
}

func (b *modelBuilder) EnterFeature(ctx *parser.FeatureContext) {
	b.enterKerMLFeatureFromDeclaration(ctx.FeatureDeclaration(), ctx.ValuePart(), locationFromContext(ctx))
}

func (b *modelBuilder) EnterFeatureSubsetting(ctx *parser.FeatureSubsettingContext) {
	if len(b.elementStack) == 0 {
		return
	}
	parentType, ok := b.elementStack[len(b.elementStack)-1].(*KerMLType)
	if !ok || parentType == nil {
		return
	}

	feature := NewKerMLFeature(extractOwnedSubsettingText(ctx.OwnedSubsetting()), locationFromContext(ctx))
	feature.parent = parentType
	if subsettings := ctx.Subsettings(); subsettings != nil {
		for _, owned := range subsettings.AllOwnedSubsetting() {
			feature.AddUnresolvedSubsettedFeature(extractOwnedSubsettingText(owned))
		}
	}
	parentType.AddChild(feature)
}

func (b *modelBuilder) EnterStep(ctx *parser.StepContext) {
	b.enterKerMLFeatureFromDeclaration(ctx.FeatureDeclaration(), ctx.ValuePart(), locationFromContext(ctx))
}

func (b *modelBuilder) EnterExpression(ctx *parser.ExpressionContext) {
	b.enterKerMLFeatureFromDeclaration(ctx.FeatureDeclaration(), ctx.ValuePart(), locationFromContext(ctx))
}

func (b *modelBuilder) EnterBooleanExpression(ctx *parser.BooleanExpressionContext) {
	b.enterKerMLFeatureFromDeclaration(ctx.FeatureDeclaration(), ctx.ValuePart(), locationFromContext(ctx))
}

func (b *modelBuilder) EnterInvariant(ctx *parser.InvariantContext) {
	b.enterKerMLFeatureFromDeclaration(ctx.FeatureDeclaration(), ctx.ValuePart(), locationFromContext(ctx))
}

func (b *modelBuilder) enterKerMLFeatureFromDeclaration(
	decl parser.IFeatureDeclarationContext,
	valuePart parser.IValuePartContext,
	loc Location,
) {
	if len(b.elementStack) == 0 {
		return
	}
	parentType, ok := b.elementStack[len(b.elementStack)-1].(*KerMLType)
	if !ok || parentType == nil {
		return
	}

	name, shortName := extractFeatureDeclarationNames(decl)
	feature := NewKerMLFeature(name, loc)
	feature.setDeclaredShortName(shortName)
	feature.parent = parentType

	if decl != nil {
		if featSpec := decl.FeatureSpecializationPart(); featSpec != nil {
			typeRef := extractTypeReference(featSpec)
			if typeRef != "" {
				feature.TypeRef = NewRef[Element](typeRef)
				feature.unresolvedTypeReference = typeRef
			}
			for _, ref := range extractSubsettedFeatureNames(featSpec) {
				feature.AddUnresolvedSubsettedFeature(ref)
			}
			for _, ref := range extractRedefinitionNames(featSpec) {
				feature.AddUnresolvedRedefinedFeature(ref)
			}
		}
	}
	if valuePart != nil && valuePart.FeatureValue() != nil && valuePart.FeatureValue().OwnedExpression() != nil {
		feature.DefaultValue = strings.TrimSpace(valuePart.FeatureValue().OwnedExpression().GetText())
	}

	parentType.AddChild(feature)
}

func (b *modelBuilder) EnterPartDefinition(ctx *parser.PartDefinitionContext) {
	name := ""
	var subclassPart parser.ISubclassificationPartContext

	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		decl := ctx.Definition().DefinitionDeclaration()
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
		// Extract subclassification/specialization reference
		subclassPart = decl.SubclassificationPart()
	}

	part := NewPart(name, locationFromContext(ctx), true)
	part.parent = b.getCurrentParent()

	// Extract and set specialization reference
	if subclassPart != nil {
		specRef := extractSpecializationReference(subclassPart)
		if specRef != "" {
			part.SetUnresolvedSpecializes(specRef)
		}
	}

	b.addToParent(part)

	// Push part onto stack for nested elements (parts can have attributes, ports, etc.)
	b.elementStack = append(b.elementStack, part)
}

func (b *modelBuilder) ExitPartDefinition(ctx *parser.PartDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterPartUsage(ctx *parser.PartUsageContext) {
	name := ""
	typeRef := ""
	multiplicity := ""

	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		usageDecl := ctx.Usage().UsageDeclaration()
		if ident := usageDecl.Identification(); ident != nil {
			name = extractName(ident)
		}
		// Extract type reference from FeatureSpecializationPart
		if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
			typeRef = extractTypeReference(featSpecPart)
			multiplicity = extractMultiplicity(featSpecPart)
		}
	}

	part := NewPart(name, locationFromContext(ctx), false)
	part.parent = b.getCurrentParent()
	if typeRef != "" {
		part.TypeRef = NewRef[*Part](typeRef)
	}
	part.Multiplicity = multiplicity
	b.addToParent(part)

	// Push part onto stack for nested elements (parts can have attributes, ports, etc.)
	b.elementStack = append(b.elementStack, part)
}

// Debug helper - temporary for understanding grammar structure
func debugPrintFeatureSpecializationPart(featSpecPart parser.IFeatureSpecializationPartContext, prefix string) {
	if featSpecPart == nil {
		fmt.Printf("%sFeatureSpecializationPart is nil\n", prefix)
		return
	}

	featSpecs := featSpecPart.AllFeatureSpecialization()
	fmt.Printf("%sFeatureSpecializationPart has %d feature specializations\n", prefix, len(featSpecs))

	for i, featSpec := range featSpecs {
		fmt.Printf("%s  [%d] FeatureSpecialization:\n", prefix, i)

		if typings := featSpec.Typings(); typings != nil {
			fmt.Printf("%s    Has Typings\n", prefix)
			ownedTypings := typings.AllOwnedFeatureTyping()
			fmt.Printf("%s    OwnedFeatureTypings: %d\n", prefix, len(ownedTypings))
			for j, ot := range ownedTypings {
				if qname := ot.QualifiedName(); qname != nil {
					fmt.Printf("%s      [%d] QualifiedName: %s\n", prefix, j, qname.GetText())
				} else {
					fmt.Printf("%s      [%d] QualifiedName is nil\n", prefix, j)
				}
			}
		}

		if redef := featSpec.Redefinitions(); redef != nil {
			fmt.Printf("%s    Has Redefinitions\n", prefix)
		}

		if subsets := featSpec.Subsettings(); subsets != nil {
			fmt.Printf("%s    Has Subsettings\n", prefix)
		}

		if refs := featSpec.References(); refs != nil {
			fmt.Printf("%s    Has References\n", prefix)
		}

		if crosses := featSpec.Crosses(); crosses != nil {
			fmt.Printf("%s    Has Crosses\n", prefix)
		}
	}
}

func (b *modelBuilder) ExitPartUsage(ctx *parser.PartUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterItemDefinition(ctx *parser.ItemDefinitionContext) {
	name := ""
	var subclassPart parser.ISubclassificationPartContext

	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		decl := ctx.Definition().DefinitionDeclaration()
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
		// Extract subclassification/specialization reference
		subclassPart = decl.SubclassificationPart()
	}

	item := NewItem(name, locationFromContext(ctx), true)
	item.parent = b.getCurrentParent()

	// Extract and set specialization reference
	if subclassPart != nil {
		specRef := extractSpecializationReference(subclassPart)
		if specRef != "" {
			item.SetUnresolvedSpecializes(specRef)
		}
	}

	b.addToParent(item)
}

func (b *modelBuilder) EnterItemUsage(ctx *parser.ItemUsageContext) {
	name := ""
	typeRef := ""
	subsettedRefs := make([]string, 0)
	redefinedRefs := make([]string, 0)
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		usageDecl := ctx.Usage().UsageDeclaration()
		if ident := usageDecl.Identification(); ident != nil {
			name = extractName(ident)
		}
		if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
			typeRef = extractTypeReference(featSpecPart)
			subsettedRefs = extractSubsettedFeatureNames(featSpecPart)
			redefinedRefs = extractRedefinitionNames(featSpecPart)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	item := NewItem(name, loc, false)
	if typeRef != "" {
		item.TypeRef = NewRef[*Item](typeRef)
	}
	for _, ref := range subsettedRefs {
		item.AddUnresolvedSubsettedFeature(ref)
	}
	for _, ref := range redefinedRefs {
		item.AddUnresolvedRedefinedFeature(ref)
	}

	if b.currentPkg != nil {
		item.parent = b.currentPkg
		b.currentPkg.AddChild(item)
	}
}

func (b *modelBuilder) EnterImport_(ctx *parser.Import_Context) {
	loc := locationFromContext(ctx)
	namespace := ""
	imp := NewImport(namespace, loc)
	imp.Visibility = extractVisibilityIndicator(ctx.VisibilityIndicator())

	if decl := ctx.ImportDeclaration(); decl != nil {
		if membership := decl.MembershipImport(); membership != nil {
			imp.IsMembership = true
			namespace = extractQualifiedName(membership.QualifiedName())
			imp.ImportedNamespace = namespace
			if membership.COLONCOLON() != nil && membership.STARSTAR() != nil {
				imp.IsAll = true
				imp.IsRecursive = true
			}
		}
		if ns := decl.NamespaceImport(); ns != nil {
			imp.IsNamespace = true
			if ns.FilterPackage() != nil {
				filter := ns.FilterPackage()
				if filterImport := filter.FilterPackageImportPart(); filterImport != nil {
					namespace = extractQualifiedName(filterImport.QualifiedName())
					imp.ImportedNamespace = namespace
					imp.IsAll = filterImport.STAR() != nil
					imp.IsRecursive = filterImport.STARSTAR() != nil
				}
				for _, member := range filter.AllFilterPackageMember() {
					if expr := member.OwnedExpression(); expr != nil {
						imp.FilterExpressions = append(imp.FilterExpressions, expr.GetText())
					}
				}
			} else {
				namespace = extractQualifiedName(ns.QualifiedName())
				imp.ImportedNamespace = namespace
				if ns.STAR() != nil {
					imp.IsAll = true
				}
				if ns.STARSTAR() != nil {
					imp.IsAll = true
					imp.IsRecursive = true
				}
			}
		}
	}
	if imp.ImportedNamespace == "" {
		imp.ImportedNamespace = namespace
	}
	if imp.ImportedNamespace == "" {
		imp.ImportedNamespace = ctx.GetText()
	}

	// Parse import patterns and detect wildcards
	b.parseImportPattern(imp, imp.ImportedNamespace)

	// Try to resolve library import if registry is available
	if b.libraryRegistry != nil {
		b.resolveLibraryImport(imp, imp.ImportedNamespace)
	}

	if b.currentPkg != nil {
		imp.parent = b.currentPkg
		b.currentPkg.AddChild(imp)
	} else {
		b.model.AddImport(imp)
	}
}

// EnterAliasMember maps alias members to first-class Alias model elements.
func (b *modelBuilder) EnterAliasMember(ctx *parser.AliasMemberContext) {
	// Grammar (SysML/KerML): 'alias' ( '<' memberShortName = NAME '>' )? ( memberName = NAME )? 'for' ...
	var shortName, name string
	names := ctx.AllName()
	if ctx.LT() != nil {
		if len(names) > 0 {
			shortName = strings.TrimSpace(names[0].GetText())
		}
		if len(names) > 1 {
			name = strings.TrimSpace(names[1].GetText())
		}
	} else if len(names) > 0 {
		name = strings.TrimSpace(names[0].GetText())
	}

	alias := NewAlias(name, locationFromContext(ctx))
	alias.setDeclaredShortName(shortName)
	alias.SetUnresolvedTarget(extractQualifiedName(ctx.QualifiedName()))
	alias.parent = b.getCurrentParent()
	b.addToParent(alias)
	if b.currentPkg == nil {
		b.model.AddAlias(alias)
	}
}

// EnterElementFilterMember captures element filter members.
func (b *modelBuilder) EnterElementFilterMember(ctx *parser.ElementFilterMemberContext) {
	expr := strings.TrimSpace(ctx.GetText())
	filter := NewElementFilter("", locationFromContext(ctx), expr)
	filter.parent = b.getCurrentParent()
	b.addToParent(filter)
	if b.currentPkg == nil {
		b.model.AddFilter(filter)
	}
}

// EnterMembershipImport is tracked by EnterImport_ for semantic extraction.
func (b *modelBuilder) EnterMembershipImport(ctx *parser.MembershipImportContext) {}

// EnterNamespaceImport is tracked by EnterImport_ for semantic extraction.
func (b *modelBuilder) EnterNamespaceImport(ctx *parser.NamespaceImportContext) {}

// EnterFilterPackage is tracked by EnterImport_ for semantic extraction.
func (b *modelBuilder) EnterFilterPackage(ctx *parser.FilterPackageContext) {}

// parseImportPattern detects wildcard patterns in import statements.
// Sets IsAll for "::*" and IsRecursive for "::**" patterns.
func (b *modelBuilder) parseImportPattern(imp *Import, namespace string) {
	// Check for recursive wildcard (::**)
	if strings.HasSuffix(namespace, "::**") {
		imp.IsRecursive = true
		imp.IsAll = true
		return
	}

	// Check for wildcard (::*)
	if strings.HasSuffix(namespace, "::*") {
		imp.IsAll = true
		return
	}

	// Check for specific element import (contains :: but doesn't end with wildcard)
	if strings.Contains(namespace, "::") && !strings.HasSuffix(namespace, "*") {
		// This is a specific element import like "ISQ::mass"
		// IsAll and IsRecursive remain false
	}
}

// resolveLibraryImport attempts to resolve an import to a library package or element.
// If successful, sets ResolvedPackage and IsResolved fields.
func (b *modelBuilder) resolveLibraryImport(imp *Import, namespace string) {
	// Clean up the namespace for resolution
	cleanNamespace := strings.TrimSpace(namespace)

	// Remove trailing wildcards for package lookup
	pkgName := cleanNamespace
	if strings.HasSuffix(pkgName, "::**") {
		pkgName = strings.TrimSuffix(pkgName, "::**")
	} else if strings.HasSuffix(pkgName, "::*") {
		pkgName = strings.TrimSuffix(pkgName, "::*")
	}

	// Try to resolve as a package import
	if pkg, err := b.libraryRegistry.ResolveImport(pkgName); err == nil && pkg != nil {
		imp.ResolvedPackage = pkg
		imp.IsResolved = true

		// If this is a wildcard import, we could add package elements to scope
		// For now, we just record that the package was resolved
		return
	}

	// Try to resolve as a specific element import (e.g., "ISQ::mass")
	if strings.Contains(cleanNamespace, "::") && !imp.IsAll {
		if elem := b.libraryRegistry.FindElement(cleanNamespace); elem != nil {
			imp.ResolvedElement = elem
			imp.IsResolved = true
		}
	}
}

func (b *modelBuilder) EnterRequirementDefinition(ctx *parser.RequirementDefinitionContext) {
	name := ""
	shortName := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name, shortName = extractIdentificationNames(ident)
		}
	}

	req := NewRequirement(name, locationFromContext(ctx), true)
	req.setDeclaredShortName(shortName)
	req.RequirementID = shortName
	req.parent = b.getCurrentParent()
	b.addToParent(req)

	// Push requirement onto stack for nested elements
	b.elementStack = append(b.elementStack, req)
}

func (b *modelBuilder) ExitRequirementDefinition(ctx *parser.RequirementDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterRequirementUsage(ctx *parser.RequirementUsageContext) {
	name := ""
	shortName := ""
	typeRef := ""
	if ctx.ConstraintUsageDeclaration() != nil {
		if usageDecl := ctx.ConstraintUsageDeclaration().UsageDeclaration(); usageDecl != nil {
			if ident := usageDecl.Identification(); ident != nil {
				name, shortName = extractIdentificationNames(ident)
			}
			if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
				typeRef = extractTypeReference(featSpecPart)
			}
		}
	}

	req := NewRequirement(name, locationFromContext(ctx), false)
	req.setDeclaredShortName(shortName)
	req.RequirementID = shortName
	if typeRef != "" {
		req.TypeRef = NewRef[*Requirement](typeRef)
	}
	if b.rewriteHints != nil {
		line := ctx.GetStart().GetLine()
		if bindings, ok := b.rewriteHints.requirementBindingsByLine[line]; ok {
			for k, v := range bindings {
				req.Bindings[k] = v
			}
		} else if name != "" {
			if queue, ok := b.rewriteHints.requirementBindingsByName[name]; ok && len(queue) > 0 {
				for k, v := range queue[0] {
					req.Bindings[k] = v
				}
				b.rewriteHints.requirementBindingsByName[name] = queue[1:]
			}
		}
	}
	req.parent = b.getCurrentParent()
	b.addToParent(req)

	// Push requirement onto stack for nested elements
	b.elementStack = append(b.elementStack, req)
}

func (b *modelBuilder) ExitRequirementUsage(ctx *parser.RequirementUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterVerificationCaseDefinition(ctx *parser.VerificationCaseDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	ver := NewVerification(name, locationFromContext(ctx), true)
	ver.parent = b.getCurrentParent()
	b.addToParent(ver)

	// Push verification onto stack for nested elements
	b.elementStack = append(b.elementStack, ver)
}

func (b *modelBuilder) ExitVerificationCaseDefinition(ctx *parser.VerificationCaseDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterVerificationCaseUsage(ctx *parser.VerificationCaseUsageContext) {
	name := ""
	shortName := ""
	typeRef := ""
	if decl := ctx.ConstraintUsageDeclaration(); decl != nil && decl.UsageDeclaration() != nil {
		if ident := decl.UsageDeclaration().Identification(); ident != nil {
			name, shortName = extractIdentificationNames(ident)
		}
		if featSpecPart := decl.UsageDeclaration().FeatureSpecializationPart(); featSpecPart != nil {
			typeRef = extractTypeReference(featSpecPart)
		}
	}

	ver := NewVerification(name, locationFromContext(ctx), false)
	ver.setDeclaredShortName(shortName)
	if typeRef != "" {
		ver.TypeRef = NewRef[*Verification](typeRef)
	}
	ver.parent = b.getCurrentParent()
	b.addToParent(ver)

	// Push verification onto stack for nested elements
	b.elementStack = append(b.elementStack, ver)
}

func (b *modelBuilder) ExitVerificationCaseUsage(ctx *parser.VerificationCaseUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterConcernDefinition(ctx *parser.ConcernDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	concern := NewConcern(name, locationFromContext(ctx), true)
	concern.parent = b.getCurrentParent()
	b.addToParent(concern)
}

func (b *modelBuilder) EnterConcernUsage(ctx *parser.ConcernUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil && ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	concern := NewConcern(name, locationFromContext(ctx), false)
	concern.parent = b.getCurrentParent()
	b.addToParent(concern)
}

func (b *modelBuilder) EnterUseCaseDefinition(ctx *parser.UseCaseDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	useCase := NewUseCase(name, loc, true)

	if b.currentPkg != nil {
		useCase.parent = b.currentPkg
		b.currentPkg.AddChild(useCase)
	}
}

func (b *modelBuilder) EnterUseCaseUsage(ctx *parser.UseCaseUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil && ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	useCase := NewUseCase(name, loc, false)

	if b.currentPkg != nil {
		useCase.parent = b.currentPkg
		b.currentPkg.AddChild(useCase)
	}
}

func (b *modelBuilder) EnterViewDefinition(ctx *parser.ViewDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	view := NewView(name, loc, true)

	if b.currentPkg != nil {
		view.parent = b.currentPkg
		b.currentPkg.AddChild(view)
	}
}

func (b *modelBuilder) EnterViewUsage(ctx *parser.ViewUsageContext) {
	name := ""
	if ctx.UsageDeclaration() != nil {
		if ident := ctx.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	view := NewView(name, loc, false)

	if b.currentPkg != nil {
		view.parent = b.currentPkg
		b.currentPkg.AddChild(view)
	}
}

func (b *modelBuilder) EnterViewpointDefinition(ctx *parser.ViewpointDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	viewpoint := NewViewpoint(name, loc, true)

	if b.currentPkg != nil {
		viewpoint.parent = b.currentPkg
		b.currentPkg.AddChild(viewpoint)
	}
}

func (b *modelBuilder) EnterViewpointUsage(ctx *parser.ViewpointUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil && ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	viewpoint := NewViewpoint(name, loc, false)

	if b.currentPkg != nil {
		viewpoint.parent = b.currentPkg
		b.currentPkg.AddChild(viewpoint)
	}
}

func (b *modelBuilder) EnterAnalysisCaseDefinition(ctx *parser.AnalysisCaseDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	analysis := NewAnalysisCase(name, loc, true)

	if b.currentPkg != nil {
		analysis.parent = b.currentPkg
		b.currentPkg.AddChild(analysis)
	}
}

func (b *modelBuilder) EnterAnalysisCaseUsage(ctx *parser.AnalysisCaseUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	analysis := NewAnalysisCase(name, loc, false)

	if b.currentPkg != nil {
		analysis.parent = b.currentPkg
		b.currentPkg.AddChild(analysis)
	}
}

func (b *modelBuilder) EnterCaseDefinition(ctx *parser.CaseDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	c := NewCase(name, loc, true)

	if b.currentPkg != nil {
		c.parent = b.currentPkg
		b.currentPkg.AddChild(c)
	}

	// Push onto element stack for nested elements (subject, actor, objective)
	b.elementStack = append(b.elementStack, c)
}

func (b *modelBuilder) ExitCaseDefinition(ctx *parser.CaseDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterCaseUsage(ctx *parser.CaseUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	c := NewCase(name, loc, false)

	if b.currentPkg != nil {
		c.parent = b.currentPkg
		b.currentPkg.AddChild(c)
	}

	// Push onto element stack for nested elements (subject, actor, objective)
	b.elementStack = append(b.elementStack, c)
}

func (b *modelBuilder) ExitCaseUsage(ctx *parser.CaseUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterIncludeUseCaseUsage(ctx *parser.IncludeUseCaseUsageContext) {
	name := ""
	// Try to extract name from UsageDeclaration if present
	if ctx.UsageDeclaration() != nil {
		if ident := ctx.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := locationFromContext(ctx)

	include := NewIncludeUseCase(name, loc)

	// Extract the referenced use case name from OwnedReferenceSubsetting
	if ctx.OwnedReferenceSubsetting() != nil {
		refName := ctx.OwnedReferenceSubsetting().GetText()
		include.SetUnresolvedIncludedUseCase(refName)
	}

	// Add to current parent (should be a UseCase or Case)
	if len(b.elementStack) > 0 {
		parent := b.elementStack[len(b.elementStack)-1]
		include.parent = parent
		// Add to parent's children via the childAdder interface
		if container, ok := parent.(childAdder); ok {
			container.AddChild(include)
		}
	} else if b.currentPkg != nil {
		include.parent = b.currentPkg
		b.currentPkg.AddChild(include)
	}

	// Push onto element stack for nested elements
	b.elementStack = append(b.elementStack, include)
}

func (b *modelBuilder) ExitIncludeUseCaseUsage(ctx *parser.IncludeUseCaseUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterAttributeDefinition(ctx *parser.AttributeDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	attr := NewAttribute(name, loc, true)

	// Add to current parent (checks element stack first, then package)
	if len(b.elementStack) > 0 {
		// We're inside another element (Package, Part, Requirement, Verification, etc.)
		parent := b.elementStack[len(b.elementStack)-1]
		attr.parent = parent
		switch p := parent.(type) {
		case *Package:
			p.AddChild(attr)
		case *Part:
			p.AddChild(attr)
		case *Requirement:
			p.AddChild(attr)
		case *Verification:
			p.AddChild(attr)
		case *Item:
			p.AddChild(attr)
		case *Action:
			p.AddChild(attr)
		case *State:
			p.AddChild(attr)
		case *Interface:
			p.AddChild(attr)
		case *Constraint:
			p.AddChild(attr)
		default:
			// Parent type doesn't support attributes
		}
	} else if b.currentPkg != nil {
		// We're at package level (shouldn't happen if package pushed to stack)
		attr.parent = b.currentPkg
		b.currentPkg.AddChild(attr)
	}
}

func (b *modelBuilder) EnterAttributeUsage(ctx *parser.AttributeUsageContext) {
	name := ""
	typeRef := ""
	subsettedRefs := make([]string, 0)
	redefinedRefs := make([]string, 0)
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		usageDecl := ctx.Usage().UsageDeclaration()
		if ident := usageDecl.Identification(); ident != nil {
			name = extractName(ident)
		}
		if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
			typeRef = extractTypeReference(featSpecPart)
			subsettedRefs = extractSubsettedFeatureNames(featSpecPart)
			redefinedRefs = extractRedefinitionNames(featSpecPart)
		}
		// If no name from identification, check for redefinition
		if name == "" {
			if len(redefinedRefs) > 0 {
				name = redefinedRefs[0]
			} else if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
				name = extractRedefinitionName(featSpecPart)
			}
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	attr := NewAttribute(name, loc, false)
	if typeRef != "" {
		attr.TypeRef = NewRef[Element](typeRef)
	}
	for _, ref := range subsettedRefs {
		attr.AddUnresolvedSubsettedFeature(ref)
	}
	for _, ref := range redefinedRefs {
		attr.AddUnresolvedRedefinedFeature(ref)
	}

	// Extract default value if present
	if ctx.Usage() != nil && ctx.Usage().UsageCompletion() != nil {
		if valuePart := ctx.Usage().UsageCompletion().ValuePart(); valuePart != nil {
			if featureValue := valuePart.FeatureValue(); featureValue != nil {
				if expr := featureValue.OwnedExpression(); expr != nil {
					attr.DefaultValue = expr.GetText()
				}
			}
		}
	}

	// Add to current parent (checks element stack first, then package)
	parent := b.getCurrentParent()
	if parent != nil {
		attr.parent = parent
		// Type assert to call the appropriate AddChild method
		// Note: nil checks needed because getCurrentParent() can return typed nil
		switch p := parent.(type) {
		case *Package:
			if p != nil {
				p.AddChild(attr)
			}
		case *Part:
			if p != nil {
				p.AddChild(attr)
			}
		case *Requirement:
			if p != nil {
				p.AddChild(attr)
			}
		case *Verification:
			if p != nil {
				p.AddChild(attr)
			}
		case *Item:
			if p != nil {
				p.AddChild(attr)
			}
		case *Action:
			if p != nil {
				p.AddChild(attr)
			}
		case *State:
			if p != nil {
				p.AddChild(attr)
			}
		case *Interface:
			if p != nil {
				p.AddChild(attr)
			}
		case *Constraint:
			if p != nil {
				p.AddChild(attr)
			}
		default:
			// For types without specific AddChild, just set parent
			// The baseElement will track it in children
		}
	}
}

func (b *modelBuilder) EnterPortDefinition(ctx *parser.PortDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	port := NewPort(name, loc, true)

	// Add to current parent (checks element stack first, then package)
	if len(b.elementStack) > 0 {
		// We're inside another element (Package, Part, Port, etc.)
		parent := b.elementStack[len(b.elementStack)-1]
		port.parent = parent
		switch p := parent.(type) {
		case *Package:
			p.AddChild(port)
		case *Part:
			p.AddChild(port)
		case *Port:
			p.AddChild(port)
		case *Interface:
			p.AddChild(port)
		default:
			// Parent type doesn't support ports
		}
	} else if b.currentPkg != nil {
		// We're at package level (shouldn't happen if package pushed to stack)
		port.parent = b.currentPkg
		b.currentPkg.AddChild(port)
	}

	// Per SysML spec: PortDefinition always contains a ConjugatedPortDefinition
	// with effective name "~" + original port name
	conjName := "~" + name
	conjPort := NewConjugatedPort(conjName, loc)
	conjPort.parent = port
	conjPort.OriginalPort.Resolve(port)
	port.ConjugatedPort = conjPort
	port.AddChild(conjPort)

	// Push port onto stack for nested elements (ports can have nested ports)
	b.elementStack = append(b.elementStack, port)
}

func (b *modelBuilder) ExitPortDefinition(ctx *parser.PortDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterPortUsage(ctx *parser.PortUsageContext) {
	name := ""
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		if ident := ctx.Usage().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	port := NewPort(name, loc, false)

	// Add to current parent (checks element stack first, then package)
	parent := b.getCurrentParent()
	if parent != nil {
		port.parent = parent
		// Note: nil checks needed because getCurrentParent() can return typed nil
		switch p := parent.(type) {
		case *Package:
			if p != nil {
				p.AddChild(port)
			}
		case *Part:
			if p != nil {
				p.AddChild(port)
			}
		case *Port:
			if p != nil {
				p.AddChild(port)
			}
		case *Interface:
			if p != nil {
				p.AddChild(port)
			}
		default:
			// For types without specific AddChild support for ports
		}
	}

	// Push port onto stack for nested elements (ports can have nested ports)
	b.elementStack = append(b.elementStack, port)
}

func (b *modelBuilder) ExitPortUsage(ctx *parser.PortUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterConnectionDefinition(ctx *parser.ConnectionDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	conn := NewConnection(name, loc, true)

	if b.currentPkg != nil {
		conn.parent = b.currentPkg
		b.currentPkg.AddChild(conn)
	}
}

func (b *modelBuilder) EnterConnectionUsage(ctx *parser.ConnectionUsageContext) {
	name := ""
	if ctx.UsageDeclaration() != nil {
		if ident := ctx.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	conn := NewConnection(name, loc, false)

	if b.currentPkg != nil {
		conn.parent = b.currentPkg
		b.currentPkg.AddChild(conn)
	}
}

func (b *modelBuilder) EnterInterfaceDefinition(ctx *parser.InterfaceDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	iface := NewInterface(name, loc, true)

	if b.currentPkg != nil {
		iface.parent = b.currentPkg
		b.currentPkg.AddChild(iface)
	}
}

func (b *modelBuilder) EnterInterfaceUsage(ctx *parser.InterfaceUsageContext) {
	name := ""
	if ctx.InterfaceUsageDeclaration() != nil && ctx.InterfaceUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.InterfaceUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	iface := NewInterface(name, loc, false)

	if b.currentPkg != nil {
		iface.parent = b.currentPkg
		b.currentPkg.AddChild(iface)
	}
}

func (b *modelBuilder) EnterAllocationDefinition(ctx *parser.AllocationDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	alloc := NewAllocation(name, loc, true)

	if b.currentPkg != nil {
		alloc.parent = b.currentPkg
		b.currentPkg.AddChild(alloc)
	}
}

func (b *modelBuilder) EnterAllocationUsage(ctx *parser.AllocationUsageContext) {
	name := ""
	if ctx.AllocationUsageDeclaration() != nil && ctx.AllocationUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.AllocationUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	alloc := NewAllocation(name, loc, false)

	if b.currentPkg != nil {
		alloc.parent = b.currentPkg
		b.currentPkg.AddChild(alloc)
	}
}

func (b *modelBuilder) EnterStateDefinition(ctx *parser.StateDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	state := NewState(name, loc, true)

	if b.currentPkg != nil {
		state.parent = b.currentPkg
		b.currentPkg.AddChild(state)
	}
}

func (b *modelBuilder) EnterStateUsage(ctx *parser.StateUsageContext) {
	name := ""
	if ctx.ActionUsageDeclaration() != nil && ctx.ActionUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ActionUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	state := NewState(name, loc, false)

	if b.currentPkg != nil {
		state.parent = b.currentPkg
		b.currentPkg.AddChild(state)
	}
}

func (b *modelBuilder) EnterTransitionUsage(ctx *parser.TransitionUsageContext) {
	name := ""
	if ctx.UsageDeclaration() != nil {
		if ident := ctx.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	transition := NewTransition(name, loc)

	// Extract source from featureChainMember (FIRST clause)
	if ctx.FeatureChainMember() != nil {
		// Source state reference - store as unresolved for now
		sourceText := ctx.FeatureChainMember().GetText()
		transition.SetUnresolvedSource(sourceText)
	}

	// Extract target from transitionSuccessionMember (THEN clause)
	if ctx.TransitionSuccessionMember() != nil {
		targetText := ctx.TransitionSuccessionMember().GetText()
		transition.SetUnresolvedTarget(targetText)
	}

	// Extract guard expression if present
	if ctx.GuardExpressionMember() != nil {
		transition.GuardExpr = ctx.GuardExpressionMember().GetText()
	}

	// Extract trigger if present
	if ctx.TriggerActionMember() != nil {
		transition.TriggerExpr = ctx.TriggerActionMember().GetText()
	}

	if b.currentPkg != nil {
		transition.parent = b.currentPkg
		b.currentPkg.AddChild(transition)
	}
}

func (b *modelBuilder) EnterActionDefinition(ctx *parser.ActionDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	action := NewAction(name, loc, true)

	if b.currentPkg != nil {
		action.parent = b.currentPkg
		b.currentPkg.AddChild(action)
	}
}

func (b *modelBuilder) EnterActionUsage(ctx *parser.ActionUsageContext) {
	name := ""
	if ctx.ActionUsageDeclaration() != nil && ctx.ActionUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ActionUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	action := NewAction(name, loc, false)

	if b.currentPkg != nil {
		action.parent = b.currentPkg
		b.currentPkg.AddChild(action)
	}
}

func (b *modelBuilder) EnterCalculationDefinition(ctx *parser.CalculationDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	calc := NewCalculation(name, loc, true)

	if b.currentPkg != nil {
		calc.parent = b.currentPkg
		b.currentPkg.AddChild(calc)
	}
}

func (b *modelBuilder) EnterCalculationUsage(ctx *parser.CalculationUsageContext) {
	name := ""
	if ctx.ActionUsageDeclaration() != nil && ctx.ActionUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ActionUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	calc := NewCalculation(name, loc, false)

	if b.currentPkg != nil {
		calc.parent = b.currentPkg
		b.currentPkg.AddChild(calc)
	}
}

func (b *modelBuilder) EnterConstraintDefinition(ctx *parser.ConstraintDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	constraint := NewConstraint(name, loc, true)

	if b.currentPkg != nil {
		constraint.parent = b.currentPkg
		b.currentPkg.AddChild(constraint)
	}
}

func (b *modelBuilder) EnterConstraintUsage(ctx *parser.ConstraintUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil && ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	constraint := NewConstraint(name, loc, false)

	if b.currentPkg != nil {
		constraint.parent = b.currentPkg
		b.currentPkg.AddChild(constraint)
	}
}

func (b *modelBuilder) EnterAssertConstraintUsage(ctx *parser.AssertConstraintUsageContext) {
	name := ""
	if ctx.ConstraintUsageDeclaration() != nil && ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	constraint := NewConstraint(name, loc, false)
	// Assert constraints could have a flag or be distinguished somehow

	if b.currentPkg != nil {
		constraint.parent = b.currentPkg
		b.currentPkg.AddChild(constraint)
	}
}

func (b *modelBuilder) EnterEnumerationDefinition(ctx *parser.EnumerationDefinitionContext) {
	name := ""
	if ctx.DefinitionDeclaration() != nil {
		if ident := ctx.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	enum := NewEnumeration(name, loc, true)

	if b.currentPkg != nil {
		enum.parent = b.currentPkg
		b.currentPkg.AddChild(enum)
	}
}

func (b *modelBuilder) EnterEnumerationUsage(ctx *parser.EnumerationUsageContext) {
	name := ""
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		if ident := ctx.Usage().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	enum := NewEnumeration(name, loc, false)

	if b.currentPkg != nil {
		enum.parent = b.currentPkg
		b.currentPkg.AddChild(enum)
	}
}

func (b *modelBuilder) EnterEnumeratedValue(ctx *parser.EnumeratedValueContext) {
	name := ""
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		if ident := ctx.Usage().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	value := NewEnumerationValue(name, loc)

	// Add to parent (should be an Enumeration)
	if b.currentPkg != nil {
		value.parent = b.currentPkg
		b.currentPkg.AddChild(value)
	}
}

// EnterSubjectMember captures the subject of a requirement or verification.
func (b *modelBuilder) EnterSubjectMember(ctx *parser.SubjectMemberContext) {
	if len(b.elementStack) == 0 {
		return
	}

	// Get the subject name from the subjectUsage
	subjectName := ""
	if ctx.SubjectUsage() != nil {
		// Try to extract the reference from the usage
		if ctx.SubjectUsage().GetText() != "" {
			subjectName = ctx.SubjectUsage().GetText()
			// Remove "subject" keyword if present
			subjectName = strings.TrimPrefix(subjectName, "subject")
			subjectName = strings.TrimSpace(subjectName)
			// Remove semicolon if present
			subjectName = strings.TrimSuffix(subjectName, ";")
			subjectName = strings.TrimSpace(subjectName)
		}
	}

	if subjectName == "" {
		return
	}

	// Set the subject reference on the current element
	current := b.elementStack[len(b.elementStack)-1]
	switch elem := current.(type) {
	case *Requirement:
		elem.unresolvedSubject = subjectName
	case *Verification:
		elem.unresolvedSubject = subjectName
	case *Case:
		elem.SetUnresolvedSubject(subjectName)
	}
}

// EnterActorMember captures actors for requirements and cases.
func (b *modelBuilder) EnterActorMember(ctx *parser.ActorMemberContext) {
	if len(b.elementStack) == 0 {
		return
	}

	// Get the actor name from the actorUsage
	actorName := ""
	if ctx.ActorUsage() != nil {
		if ctx.ActorUsage().GetText() != "" {
			actorName = ctx.ActorUsage().GetText()
			// Remove "actor" keyword if present
			actorName = strings.TrimPrefix(actorName, "actor")
			actorName = strings.TrimSpace(actorName)
			// Remove semicolon if present
			actorName = strings.TrimSuffix(actorName, ";")
			actorName = strings.TrimSpace(actorName)
		}
	}

	if actorName == "" {
		return
	}

	// Set the actor reference on the current element
	current := b.elementStack[len(b.elementStack)-1]
	switch elem := current.(type) {
	case *Case:
		elem.AddUnresolvedActor(actorName)
	}
}

// EnterObjectiveMember captures objectives for cases.
func (b *modelBuilder) EnterObjectiveMember(ctx *parser.ObjectiveMemberContext) {
	if len(b.elementStack) == 0 {
		return
	}

	// Get the objective requirement name
	objectiveName := ""
	if ctx.ObjectiveRequirementUsage() != nil {
		if ctx.ObjectiveRequirementUsage().GetText() != "" {
			objectiveName = ctx.ObjectiveRequirementUsage().GetText()
			// Remove "objective" keyword if present
			objectiveName = strings.TrimPrefix(objectiveName, "objective")
			objectiveName = strings.TrimSpace(objectiveName)
			// Remove semicolon if present
			objectiveName = strings.TrimSuffix(objectiveName, ";")
			objectiveName = strings.TrimSpace(objectiveName)
		}
	}

	if objectiveName == "" {
		return
	}

	// Set the objective reference on the current element
	current := b.elementStack[len(b.elementStack)-1]
	switch elem := current.(type) {
	case *Case:
		elem.AddUnresolvedObjective(objectiveName)
	}
}

// EnterRequirementConstraintMember captures assume/require constraints in requirements.
func (b *modelBuilder) EnterRequirementConstraintMember(ctx *parser.RequirementConstraintMemberContext) {
	if len(b.elementStack) == 0 {
		return
	}

	current := b.elementStack[len(b.elementStack)-1]
	req, ok := current.(*Requirement)
	if !ok {
		return
	}

	// Determine if this is an assume or require constraint
	isAssume := false
	if ctx.RequirementKind() != nil {
		kindText := ctx.RequirementKind().GetText()
		isAssume = kindText == "assume"
	}

	// Extract the constraint expression
	expr := ""
	if ctx.RequirementConstraintUsage() != nil {
		expr = ctx.RequirementConstraintUsage().GetText()
	}

	if expr == "" {
		return
	}
	if b.rewriteHints != nil {
		placeholder := strings.TrimSuffix(strings.TrimSpace(expr), ";")
		if replacementExpr, ok := b.rewriteHints.requireConstraintExprByPlaceholder[placeholder]; ok {
			expr = replacementExpr
		}
	}

	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}
	if ctx.GetStop() != nil {
		loc.EndLine = ctx.GetStop().GetLine()
		loc.EndColumn = ctx.GetStop().GetColumn()
	}

	constraint := NewRequirementConstraint(loc, isAssume, expr)
	constraint.parent = req

	if isAssume {
		req.Assumptions = append(req.Assumptions, constraint)
	} else {
		req.Constraints = append(req.Constraints, constraint)
	}
}

// extractName extracts the name from an Identification context.
func extractName(ident parser.IIdentificationContext) string {
	declaredName, declaredShortName := extractIdentificationNames(ident)
	if declaredName != "" {
		return declaredName
	}
	return declaredShortName
}

// extractIdentificationNames extracts the names from an Identification context
// using grammar semantics:
//   - declaredShortName: optional '< NAME >'
//   - declaredName: optional NAME
//
// The Identification rule allows either or both.
func extractIdentificationNames(ident parser.IIdentificationContext) (declaredName, declaredShortName string) {
	if ident == nil {
		return "", ""
	}

	names := ident.AllName()
	switch len(names) {
	case 0:
		return "", ""
	case 1:
		// If '<' exists, single name is declaredShortName; otherwise declaredName.
		if ident.LT() != nil {
			return "", names[0].GetText()
		}
		return names[0].GetText(), ""
	default:
		// With both present, order in parse tree is: declaredShortName, declaredName.
		return names[1].GetText(), names[0].GetText()
	}
}

func extractQualifiedName(qname parser.IQualifiedNameContext) string {
	if qname == nil {
		return ""
	}
	return strings.TrimSpace(qname.GetText())
}

func extractOwnedReferenceSubsetting(ref parser.IOwnedReferenceSubsettingContext) string {
	if ref == nil {
		return ""
	}
	if qname := ref.QualifiedName(); qname != nil {
		return extractQualifiedName(qname)
	}
	if chain := ref.OwnedFeatureChain(); chain != nil {
		return strings.TrimSpace(chain.GetText())
	}
	return ""
}

func extractVisibilityIndicator(v parser.IVisibilityIndicatorContext) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(v.GetText())
}

// extractRedefinitionName extracts the name from a redefinition in featureSpecializationPart
func extractRedefinitionName(featSpecPart parser.IFeatureSpecializationPartContext) string {
	names := extractRedefinitionNames(featSpecPart)
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

// extractRedefinitionNames extracts all redefined feature names from a featureSpecializationPart.
func extractRedefinitionNames(featSpecPart parser.IFeatureSpecializationPartContext) []string {
	if featSpecPart == nil {
		return nil
	}

	names := make([]string, 0)

	// Iterate through all feature specializations
	for _, featSpec := range featSpecPart.AllFeatureSpecialization() {
		// Check if this is a redefinition
		if redef := featSpec.Redefinitions(); redef != nil {
			for _, ownedRedef := range redef.AllOwnedRedefinition() {
				if qname := ownedRedef.QualifiedName(); qname != nil {
					names = append(names, strings.TrimSpace(qname.GetText()))
				} else if chain := ownedRedef.OwnedFeatureChain(); chain != nil {
					names = append(names, strings.TrimSpace(chain.GetText()))
				}
			}
		}
	}

	return names
}

// extractSubsettedFeatureNames extracts all subsetted or reference-subsetted feature names.
func extractSubsettedFeatureNames(featSpecPart parser.IFeatureSpecializationPartContext) []string {
	if featSpecPart == nil {
		return nil
	}

	names := make([]string, 0)

	for _, featSpec := range featSpecPart.AllFeatureSpecialization() {
		// :> / subsets
		if subsettings := featSpec.Subsettings(); subsettings != nil {
			for _, ownedSubsetting := range subsettings.AllOwnedSubsetting() {
				if qname := ownedSubsetting.QualifiedName(); qname != nil {
					names = append(names, strings.TrimSpace(qname.GetText()))
				} else {
					names = append(names, strings.TrimSpace(ownedSubsetting.GetText()))
				}
			}
		}

		// ::> / references
		if references := featSpec.References(); references != nil {
			if ref := references.OwnedReferenceSubsetting(); ref != nil {
				names = append(names, extractOwnedReferenceSubsetting(ref))
			}
		}
	}

	return names
}

// extractTypeReference extracts the type name from a featureSpecializationPart
// This handles the "part name: Type" syntax where the type is specified via a typing relationship
func extractTypeReference(featSpecPart parser.IFeatureSpecializationPartContext) string {
	if featSpecPart == nil {
		return ""
	}

	// Iterate through all feature specializations
	featSpecs := featSpecPart.AllFeatureSpecialization()
	for _, featSpec := range featSpecs {
		// Check if this is a typing (type reference)
		if typings := featSpec.Typings(); typings != nil {
			// Try TypedBy first (for "name: Type" syntax)
			if typedBy := typings.TypedBy(); typedBy != nil {
				if ownedFeatTyping := typedBy.OwnedFeatureTyping(); ownedFeatTyping != nil {
					if qname := ownedFeatTyping.QualifiedName(); qname != nil {
						return qname.GetText()
					}
				}
			}

			// Fallback: Get all owned feature typings
			ownedTypings := typings.AllOwnedFeatureTyping()
			if len(ownedTypings) > 0 {
				// Get the qualified name from the owned feature typing
				if qname := ownedTypings[0].QualifiedName(); qname != nil {
					return qname.GetText()
				}
			}
		}

		// Fallback: some textual forms (e.g. ":>") are parsed as subsettings.
		// For high-level typing needs, use the first referenced qualified name.
		if subsettings := featSpec.Subsettings(); subsettings != nil {
			ownedSubsettings := subsettings.AllOwnedSubsetting()
			if len(ownedSubsettings) > 0 {
				if qname := ownedSubsettings[0].QualifiedName(); qname != nil {
					return qname.GetText()
				}
			}
		}
	}

	return ""
}

// extractMultiplicity extracts a usage multiplicity from a featureSpecializationPart.
// Returned values are normalized without brackets, e.g. "4", "0..1", "*".
func extractMultiplicity(featSpecPart parser.IFeatureSpecializationPartContext) string {
	if featSpecPart == nil || featSpecPart.MultiplicityPart() == nil {
		return ""
	}

	multiplicityPart := featSpecPart.MultiplicityPart()
	if multiplicityPart.OwnedMultiplicity() == nil {
		return ""
	}

	multiplicityRange := multiplicityPart.OwnedMultiplicity().MultiplicityRange()
	if multiplicityRange == nil {
		return ""
	}

	members := multiplicityRange.AllMultiplicityExpressionMember()
	switch len(members) {
	case 0:
		return ""
	case 1:
		return strings.TrimSpace(members[0].GetText())
	default:
		lower := strings.TrimSpace(members[0].GetText())
		upper := strings.TrimSpace(members[1].GetText())
		if lower != "" && upper != "" {
			return lower + ".." + upper
		}
	}

	// Fallback to direct multiplicity range text if members are unexpectedly empty.
	text := strings.TrimSpace(multiplicityRange.GetText())
	text = strings.TrimPrefix(text, "[")
	text = strings.TrimSuffix(text, "]")
	return strings.TrimSpace(text)
}

func extractClassifierSpecializations(decl parser.IClassifierDeclarationContext) []string {
	if decl == nil || decl.SuperclassingPart() == nil {
		return nil
	}
	refs := make([]string, 0)
	for _, owned := range decl.SuperclassingPart().AllOwnedSubclassification() {
		if owned == nil {
			continue
		}
		if q := owned.QualifiedName(); q != nil {
			refs = append(refs, strings.TrimSpace(q.GetText()))
		}
	}
	return refs
}

func extractTypeDeclarationSpecializations(decl parser.ITypeDeclarationContext) []string {
	if decl == nil {
		return nil
	}
	refs := make([]string, 0)
	for _, part := range decl.AllSpecializationPart() {
		if part == nil {
			continue
		}
		for _, owned := range part.AllOwnedSpecialization() {
			if owned == nil || owned.GeneralType() == nil {
				continue
			}
			gt := owned.GeneralType()
			if q := gt.QualifiedName(); q != nil {
				refs = append(refs, strings.TrimSpace(q.GetText()))
				continue
			}
			if chain := gt.OwnedFeatureChain(); chain != nil {
				refs = append(refs, strings.TrimSpace(chain.GetText()))
			}
		}
	}
	return refs
}

func extractFeatureDeclarationNames(decl parser.IFeatureDeclarationContext) (declaredName, declaredShortName string) {
	if decl == nil || decl.FeatureIdentification() == nil {
		return "", ""
	}
	ident := decl.FeatureIdentification()
	names := ident.AllName()
	switch len(names) {
	case 0:
		return "", ""
	case 1:
		if ident.LT() != nil {
			return "", names[0].GetText()
		}
		return names[0].GetText(), ""
	default:
		return names[1].GetText(), names[0].GetText()
	}
}

func extractOwnedSubsettingText(owned parser.IOwnedSubsettingContext) string {
	if owned == nil {
		return ""
	}
	if q := owned.QualifiedName(); q != nil {
		return strings.TrimSpace(q.GetText())
	}
	return strings.TrimSpace(owned.GetText())
}

// extractSpecializationReference extracts the specialization/supertype reference
// from a SubclassificationPart context. This handles the "part def X :> Y" syntax.
func extractSpecializationReference(subclassPart parser.ISubclassificationPartContext) string {
	if subclassPart == nil {
		return ""
	}

	// Get all owned subclassifications
	subclassifications := subclassPart.AllOwnedSubclassification()
	for _, subclass := range subclassifications {
		if qname := subclass.QualifiedName(); qname != nil {
			return qname.GetText()
		}
	}

	return ""
}

// MustParseString parses a SysML string and panics if there are errors.
func MustParseString(input string) *Model {
	result := ParseString(input)
	if !result.Success() {
		panic(result.Errors)
	}
	return result.Model
}

// MustParseFile parses a SysML file and panics if there are errors.
func MustParseFile(filename string) *Model {
	result := ParseFile(filename)
	if !result.Success() {
		panic(result.Errors)
	}
	return result.Model
}

// Validate checks if the input is valid SysML without building a full model.
func Validate(input string) *ParseError {
	lowErrors := low.Validate(input)
	if lowErrors.HasErrors() {
		return convertFromLowLevel(lowErrors, "<string>")
	}
	return nil
}

// ValidateFile checks if a file contains valid SysML.
func ValidateFile(filename string) *ParseError {
	content, err := os.ReadFile(filename)
	if err != nil {
		return &ParseError{
			Errors: []*Error{{
				Line:    0,
				Column:  0,
				Message: err.Error(),
			}},
			Source: filename,
		}
	}

	lowErrors := low.ValidateBytes(content)
	if lowErrors.HasErrors() {
		return convertFromLowLevel(lowErrors, filename)
	}
	return nil
}

// EnterDependency handles dependency declarations.
// Dependencies declare relationships between client and supplier elements.
func (b *modelBuilder) EnterDependency(ctx *parser.DependencyContext) {
	// Create dependency
	loc := locationFromContext(ctx)
	dep := NewDependency(loc)

	// Extract client/supplier from DependencyDeclaration
	if decl := ctx.DependencyDeclaration(); decl != nil {
		// Collect all qualified names. Grammar allows lists on both client and supplier sides.
		names := decl.AllQualifiedName()
		toIndex := len(names)
		if decl.TO() != nil {
			// For grammar shape: one or more clients before TO, suppliers after TO.
			// Split at the last client position conservatively if TO is present.
			// Minimum valid shape is client TO supplier.
			if len(names) > 1 {
				toIndex = len(names) - 1
			}
		}
		for i, qn := range names {
			name := extractQualifiedName(qn)
			if name == "" {
				continue
			}
			if i < toIndex {
				dep.AddUnresolvedClient(name)
			} else {
				dep.AddUnresolvedSupplier(name)
			}
		}
	}

	// Add to current package
	if b.currentPkg != nil {
		dep.parent = b.currentPkg
		b.currentPkg.AddChild(dep)
	}

	// Add to model's dependency list
	b.model.AddDependency(dep)

	// Push to element stack for body content (annotations, etc.)
	b.elementStack = append(b.elementStack, dep)
}

// ExitDependency handles exiting dependency scope.
func (b *modelBuilder) ExitDependency(ctx *parser.DependencyContext) {
	// Pop from element stack
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterComment handles comment declarations.
// Comments provide documentation and annotations for model elements.
func (b *modelBuilder) EnterComment(ctx *parser.CommentContext) {
	// Extract body from REGULAR_COMMENT or COMMENT
	body := ""
	if ctx.REGULAR_COMMENT() != nil {
		body = ctx.REGULAR_COMMENT().GetText()
	} else if ctx.COMMENT() != nil {
		body = ctx.COMMENT().GetText()
	}

	// Extract locale if present
	locale := ""
	if ctx.LOCALE() != nil {
		locale = ctx.LOCALE().GetText()
	}

	loc := locationFromContext(ctx)
	comment := NewComment(body, loc)
	comment.Locale = locale
	for _, annotation := range ctx.AllAnnotation() {
		if annotation == nil || annotation.QualifiedName() == nil {
			continue
		}
		name := extractQualifiedName(annotation.QualifiedName())
		if name != "" {
			comment.AddUnresolvedAbout(name)
		}
	}

	// Add to current package
	if b.currentPkg != nil {
		comment.parent = b.currentPkg
		b.currentPkg.AddChild(comment)
		b.model.AddComment(comment)
	}
}

// EnterDocumentation handles documentation declarations.
// Documentation provides inline documentation for model elements.
func (b *modelBuilder) EnterDocumentation(ctx *parser.DocumentationContext) {
	// Extract body from REGULAR_COMMENT
	body := ""
	if ctx.REGULAR_COMMENT() != nil {
		body = ctx.REGULAR_COMMENT().GetText()
	}

	// Extract locale if present
	locale := ""
	if ctx.LOCALE() != nil {
		locale = ctx.LOCALE().GetText()
	}

	loc := locationFromContext(ctx)
	doc := NewDoc(body, loc)
	doc.Locale = locale

	// Add to current package
	if b.currentPkg != nil {
		doc.parent = b.currentPkg
		b.currentPkg.AddChild(doc)
		b.model.AddDoc(doc)
	}
}

// EnterFlowDefinition handles flow definition declarations.
// Flow definitions create reusable flow specifications.
func (b *modelBuilder) EnterFlowDefinition(ctx *parser.FlowDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := locationFromContext(ctx)
	flow := NewFlow(name, loc, true)

	// Add to current package
	if b.currentPkg != nil {
		flow.parent = b.currentPkg
		b.currentPkg.AddChild(flow)
		b.model.AddFlow(flow)
	}

	// Push to element stack for nested content
	b.elementStack = append(b.elementStack, flow)
}

// ExitFlowDefinition pops the flow from the element stack.
func (b *modelBuilder) ExitFlowDefinition(ctx *parser.FlowDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterFlowUsage handles flow usage declarations.
// Flow usages instantiate flows within a specific context.
func (b *modelBuilder) EnterFlowUsage(ctx *parser.FlowUsageContext) {
	name := ""
	if ctx.FlowDeclaration() != nil && ctx.FlowDeclaration().UsageDeclaration() != nil {
		if ident := ctx.FlowDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := locationFromContext(ctx)
	flow := NewFlow(name, loc, false)

	// Add to current package
	if b.currentPkg != nil {
		flow.parent = b.currentPkg
		b.currentPkg.AddChild(flow)
		b.model.AddFlow(flow)
	}

	// Push to element stack for nested content
	b.elementStack = append(b.elementStack, flow)
}

// ExitFlowUsage pops the flow from the element stack.
func (b *modelBuilder) ExitFlowUsage(ctx *parser.FlowUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterSuccessionFlowUsage handles succession flow declarations in action bodies.
// SuccessionFlowUsage represents control flow succession from source to target elements.
func (b *modelBuilder) EnterSuccessionFlowUsage(ctx *parser.SuccessionFlowUsageContext) {
	name := ""
	if ctx.FlowDeclaration() != nil && ctx.FlowDeclaration().UsageDeclaration() != nil {
		if ident := ctx.FlowDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := locationFromContext(ctx)
	succession := NewSuccessionFlow(name, loc)

	// Extract source and target from flow declaration
	if ctx.FlowDeclaration() != nil {
		flowDecl := ctx.FlowDeclaration()

		// Check for FROM/TO syntax: from source to target
		if flowDecl.FROM() != nil && flowDecl.TO() != nil {
			// Get source and target from FlowEndMember
			flowEndMembers := flowDecl.AllFlowEndMember()
			if len(flowEndMembers) >= 2 {
				// Extract source from first flow end member
				sourceText := flowEndMembers[0].GetText()
				succession.SetUnresolvedSource(sourceText)

				// Extract target from second flow end member
				targetText := flowEndMembers[1].GetText()
				succession.SetUnresolvedTarget(targetText)
			}
		} else if flowEndMembers := flowDecl.AllFlowEndMember(); len(flowEndMembers) >= 2 {
			// Inline syntax: source to target
			sourceText := flowEndMembers[0].GetText()
			succession.SetUnresolvedSource(sourceText)

			targetText := flowEndMembers[1].GetText()
			succession.SetUnresolvedTarget(targetText)
		}
	}

	// Add to current parent (package or action)
	parent := b.getCurrentParent()
	if parent != nil {
		succession.parent = parent
		if container, ok := parent.(childAdder); ok {
			container.AddChild(succession)
		}
	}

	// Push to element stack for nested content
	b.elementStack = append(b.elementStack, succession)
}

// ExitSuccessionFlowUsage pops the succession flow from the element stack.
func (b *modelBuilder) ExitSuccessionFlowUsage(ctx *parser.SuccessionFlowUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterControlNode handles control node declarations in activity diagrams.
// Control nodes manage flow: fork (split), join (synchronize), merge (combine), decision (branch).
func (b *modelBuilder) EnterControlNode(ctx *parser.ControlNodeContext) {
	// Determine control node kind from context
	var kind ControlNodeKind
	if ctx.MergeNode() != nil {
		kind = ControlNodeKindMerge
	} else if ctx.DecisionNode() != nil {
		kind = ControlNodeKindDecision
	} else if ctx.JoinNode() != nil {
		kind = ControlNodeKindJoin
	} else if ctx.ForkNode() != nil {
		kind = ControlNodeKindFork
	} else {
		kind = ControlNodeKindMerge // default
	}

	loc := locationFromContext(ctx)
	node := NewControlNode(kind, loc)

	// Add to current package
	if b.currentPkg != nil {
		node.parent = b.currentPkg
		b.currentPkg.AddChild(node)
	}

	// Add to model
	b.model.AddControlNode(node)

	// Push to element stack
	b.elementStack = append(b.elementStack, node)
}

// ExitControlNode pops the control node from the element stack.
func (b *modelBuilder) ExitControlNode(ctx *parser.ControlNodeContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterOccurrenceDefinition handles occurrence definition declarations.
// Occurrence definitions define time-based occurrence types.
func (b *modelBuilder) EnterOccurrenceDefinition(ctx *parser.OccurrenceDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	loc := locationFromContext(ctx)
	occ := NewOccurrence(name, loc, true, false)

	// Add to current package
	if b.currentPkg != nil {
		occ.parent = b.currentPkg
		b.currentPkg.AddChild(occ)
	}

	// Add to model
	b.model.AddOccurrence(occ)

	// Push to element stack
	b.elementStack = append(b.elementStack, occ)
}

// ExitOccurrenceDefinition pops the occurrence from the element stack.
func (b *modelBuilder) ExitOccurrenceDefinition(ctx *parser.OccurrenceDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterIndividualDefinition handles individual definition declarations.
func (b *modelBuilder) EnterIndividualDefinition(ctx *parser.IndividualDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	occ := NewOccurrence(name, locationFromContext(ctx), true, true)
	b.addOccurrenceToModel(occ)
	b.elementStack = append(b.elementStack, occ)
}

func (b *modelBuilder) ExitIndividualDefinition(ctx *parser.IndividualDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterIndividualUsage handles individual usage declarations.
func (b *modelBuilder) EnterIndividualUsage(ctx *parser.IndividualUsageContext) {
	name := ""
	if usage := ctx.Usage(); usage != nil && usage.UsageDeclaration() != nil {
		if ident := usage.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	occ := NewOccurrence(name, locationFromContext(ctx), false, true)
	b.addOccurrenceToModel(occ)
	b.elementStack = append(b.elementStack, occ)
}

func (b *modelBuilder) ExitIndividualUsage(ctx *parser.IndividualUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterPortionUsage handles snapshot/timeslice portion usages.
func (b *modelBuilder) EnterPortionUsage(ctx *parser.PortionUsageContext) {
	name := ""
	if usage := ctx.Usage(); usage != nil && usage.UsageDeclaration() != nil {
		if ident := usage.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	occ := NewOccurrence(name, locationFromContext(ctx), false, ctx.INDIVIDUAL() != nil)
	if kind := ctx.PortionKind(); kind != nil {
		switch {
		case kind.SNAPSHOT() != nil:
			occ.SetPortionKind(PortionKindSnapshot)
			occ.IsSnapshot = true
		case kind.TIMESLICE() != nil:
			occ.SetPortionKind(PortionKindTimeslice)
			occ.IsTimeSlice = true
		}
	}
	b.addOccurrenceToModel(occ)
	b.elementStack = append(b.elementStack, occ)
}

func (b *modelBuilder) ExitPortionUsage(ctx *parser.PortionUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterEventOccurrenceUsage handles event occurrence usages.
func (b *modelBuilder) EnterEventOccurrenceUsage(ctx *parser.EventOccurrenceUsageContext) {
	name := ""
	if decl := ctx.UsageDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	occ := NewEventOccurrence(name, locationFromContext(ctx))
	b.addOccurrenceToModel(occ)
	b.elementStack = append(b.elementStack, occ)
}

func (b *modelBuilder) ExitEventOccurrenceUsage(ctx *parser.EventOccurrenceUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterPerformActionUsage handles perform-action usage declarations.
func (b *modelBuilder) EnterPerformActionUsage(ctx *parser.PerformActionUsageContext) {
	name := ""
	if decl := ctx.PerformActionUsageDeclaration(); decl != nil && decl.UsageDeclaration() != nil {
		if ident := decl.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	action := NewAction(name, locationFromContext(ctx), false)
	action.parent = b.getCurrentParent()
	b.addToParent(action)
	b.elementStack = append(b.elementStack, action)
}

func (b *modelBuilder) ExitPerformActionUsage(ctx *parser.PerformActionUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterExhibitStateUsage handles exhibit-state usage declarations.
func (b *modelBuilder) EnterExhibitStateUsage(ctx *parser.ExhibitStateUsageContext) {
	name := ""
	if decl := ctx.UsageDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	} else if ref := ctx.OwnedReferenceSubsetting(); ref != nil {
		name = extractOwnedReferenceSubsetting(ref)
	}
	state := NewState(name, locationFromContext(ctx), false)
	state.parent = b.getCurrentParent()
	b.addToParent(state)
	b.elementStack = append(b.elementStack, state)
}

func (b *modelBuilder) ExitExhibitStateUsage(ctx *parser.ExhibitStateUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterSatisfyRequirementUsage captures satisfy relationships as typed edges.
func (b *modelBuilder) EnterSatisfyRequirementUsage(ctx *parser.SatisfyRequirementUsageContext) {
	rel := NewSatisfyRelationship("", locationFromContext(ctx))
	if subj := ctx.SatisfactionSubjectMember(); subj != nil && subj.SatisfactionParameter() != nil {
		rel.SetUnresolvedSatisfier(strings.TrimSpace(subj.SatisfactionParameter().GetText()))
	}
	switch {
	case ctx.OwnedReferenceSubsetting() != nil:
		rel.SetUnresolvedRequired(extractOwnedReferenceSubsetting(ctx.OwnedReferenceSubsetting()))
	case ctx.UsageDeclaration() != nil && ctx.UsageDeclaration().Identification() != nil:
		rel.SetUnresolvedRequired(extractName(ctx.UsageDeclaration().Identification()))
	}
	rel.parent = b.getCurrentParent()
	b.addToParent(rel)
	b.model.AddSatisfy(rel)
}

// EnterRequirementVerificationUsage captures verify relationships as typed edges.
func (b *modelBuilder) EnterRequirementVerificationUsage(ctx *parser.RequirementVerificationUsageContext) {
	rel := NewVerifyRelationship("", locationFromContext(ctx))
	switch {
	case ctx.OwnedReferenceSubsetting() != nil:
		rel.SetUnresolvedRequired(extractOwnedReferenceSubsetting(ctx.OwnedReferenceSubsetting()))
	case ctx.ConstraintUsageDeclaration() != nil &&
		ctx.ConstraintUsageDeclaration().UsageDeclaration() != nil &&
		ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification() != nil:
		rel.SetUnresolvedRequired(extractName(ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification()))
	}

	// In common SysML forms, verify relationships are contained within a verification
	// definition/usage. Use the enclosing verification as verifier when available.
	if owner := b.getCurrentVerification(); owner != nil && owner.Name() != "" {
		rel.SetUnresolvedVerifier(owner.Name())
	} else if decl := ctx.ConstraintUsageDeclaration(); decl != nil && decl.UsageDeclaration() != nil {
		if ident := decl.UsageDeclaration().Identification(); ident != nil {
			rel.SetUnresolvedVerifier(extractName(ident))
		}
	}
	rel.parent = b.getCurrentParent()
	b.addToParent(rel)
	b.model.AddVerify(rel)
}

func (b *modelBuilder) getCurrentVerification() *Verification {
	for i := len(b.elementStack) - 1; i >= 0; i-- {
		if ver, ok := b.elementStack[i].(*Verification); ok {
			return ver
		}
	}
	return nil
}

// EnterMessage maps message usage to Message model nodes.
func (b *modelBuilder) EnterMessage(ctx *parser.MessageContext) {
	name := ""
	if decl := ctx.MessageDeclaration(); decl != nil && decl.UsageDeclaration() != nil {
		if ident := decl.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	msg := NewMessage(name, locationFromContext(ctx))
	if decl := ctx.MessageDeclaration(); decl != nil {
		events := decl.AllMessageEventMember()
		if len(events) >= 2 {
			if left := events[0].MessageEvent(); left != nil && left.OwnedReferenceSubsetting() != nil {
				msg.SetUnresolvedSender(extractOwnedReferenceSubsetting(left.OwnedReferenceSubsetting()))
			}
			if right := events[1].MessageEvent(); right != nil && right.OwnedReferenceSubsetting() != nil {
				msg.SetUnresolvedReceiver(extractOwnedReferenceSubsetting(right.OwnedReferenceSubsetting()))
			}
		}
		if payload := decl.FlowPayloadFeatureMember(); payload != nil {
			msg.Payload = payload.GetText()
		}
	}
	msg.parent = b.getCurrentParent()
	b.addToParent(msg)
	if b.currentPkg == nil {
		b.model.AddMessage(msg)
	}
}

// EnterRenderingDefinition handles rendering definitions.
func (b *modelBuilder) EnterRenderingDefinition(ctx *parser.RenderingDefinitionContext) {
	name := ""
	if def := ctx.Definition(); def != nil && def.DefinitionDeclaration() != nil {
		if ident := def.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	rendering := NewRendering(name, locationFromContext(ctx), true)
	rendering.parent = b.getCurrentParent()
	b.addToParent(rendering)
	if b.currentPkg == nil {
		b.model.AddRendering(rendering)
	}
}

// EnterRenderingUsage handles rendering usages.
func (b *modelBuilder) EnterRenderingUsage(ctx *parser.RenderingUsageContext) {
	name := ""
	if usage := ctx.Usage(); usage != nil && usage.UsageDeclaration() != nil {
		if ident := usage.UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	rendering := NewRendering(name, locationFromContext(ctx), false)
	if usage := ctx.Usage(); usage != nil && usage.UsageDeclaration() != nil && usage.UsageDeclaration().FeatureSpecializationPart() != nil {
		typeRef := extractTypeReference(usage.UsageDeclaration().FeatureSpecializationPart())
		if typeRef != "" {
			rendering.TypeRef = NewRef[*Rendering](typeRef)
		}
	}
	rendering.parent = b.getCurrentParent()
	b.addToParent(rendering)
	if b.currentPkg == nil {
		b.model.AddRendering(rendering)
	}
}

// EnterMetadataDefinition handles metadata definitions.
func (b *modelBuilder) EnterMetadataDefinition(ctx *parser.MetadataDefinitionContext) {
	name := ""
	if def := ctx.Definition(); def != nil && def.DefinitionDeclaration() != nil {
		if ident := def.DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}
	md := NewMetadata(name, locationFromContext(ctx), true)
	md.parent = b.getCurrentParent()
	b.addToParent(md)
	if b.currentPkg == nil {
		b.model.AddMetadata(md)
	}
	b.elementStack = append(b.elementStack, md)
}

func (b *modelBuilder) ExitMetadataDefinition(ctx *parser.MetadataDefinitionContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterMetadataUsage handles metadata usages.
func (b *modelBuilder) EnterMetadataUsage(ctx *parser.MetadataUsageContext) {
	name := ""
	typeRef := ""
	if decl := ctx.MetadataUsageDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
		if typing := decl.OwnedFeatureTyping(); typing != nil && typing.QualifiedName() != nil {
			typeRef = extractQualifiedName(typing.QualifiedName())
		}
	}
	md := NewMetadata(name, locationFromContext(ctx), false)
	if typeRef != "" {
		md.TypeRef = NewRef[*Metadata](typeRef)
	}
	for _, ann := range ctx.AllAnnotation() {
		if ann != nil && ann.QualifiedName() != nil {
			node := NewPrefixMetadataAnnotation(locationFromContext(ctx))
			node.SetUnresolvedMetadata(extractQualifiedName(ann.QualifiedName()))
			md.AddChild(node)
		}
	}
	md.parent = b.getCurrentParent()
	b.addToParent(md)
	if b.currentPkg == nil {
		b.model.AddMetadata(md)
	}
}

// EnterBindingConnectorAsUsage handles binding connector declarations.
// Binding connectors bind values between features.
func (b *modelBuilder) EnterBindingConnectorAsUsage(ctx *parser.BindingConnectorAsUsageContext) {
	// Binding connectors are handled as specialized connections
	// Extract the binding information from the context
	loc := locationFromContext(ctx)

	// Create a connection to represent the binding
	conn := NewConnection("", loc, false)
	conn.parent = b.getCurrentParent()

	// Add to current package
	if b.currentPkg != nil {
		b.currentPkg.AddChild(conn)
	}

	// Mark as binding connector (using first end as indicator)
	// The actual binding semantics will be resolved during reference resolution

	// Push to element stack for nested content
	b.elementStack = append(b.elementStack, conn)
}

// ExitBindingConnectorAsUsage pops the binding connector from the element stack.
func (b *modelBuilder) ExitBindingConnectorAsUsage(ctx *parser.BindingConnectorAsUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

// EnterSuccessionAsUsage handles succession declarations.
// Successions define predecessor-successor relationships.
func (b *modelBuilder) EnterSuccessionAsUsage(ctx *parser.SuccessionAsUsageContext) {
	// Successions are handled as transitions or specialized connections
	// Extract the succession information from the context
	loc := locationFromContext(ctx)

	// Create a transition to represent the succession
	trans := NewTransition("", loc)
	trans.parent = b.getCurrentParent()

	// Add to current package
	if b.currentPkg != nil {
		b.currentPkg.AddChild(trans)
	}

	// Push to element stack for nested content
	b.elementStack = append(b.elementStack, trans)
}

// ExitSuccessionAsUsage pops the succession from the element stack.
func (b *modelBuilder) ExitSuccessionAsUsage(ctx *parser.SuccessionAsUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}
