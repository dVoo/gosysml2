package sysml

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	discardTree bool
}

// WithDiscardTree discards the parse tree after building the model.
// This significantly reduces memory usage for large files.
func WithDiscardTree() ParseOption {
	return func(c *parseConfig) {
		c.discardTree = true
	}
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
	// Use low-level parser directly with bytes
	tree, lowErrors := low.ParseBytes(input)

	cfg := &parseConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	result := &ParseResult{
		Source: source,
	}

	if !cfg.discardTree {
		result.Tree = tree
	}

	if lowErrors.HasErrors() {
		result.Errors = convertFromLowLevel(lowErrors, source)
	}

	if tree != nil {
		result.Model = buildModel(tree)
	}

	return result
}

// ParseReader parses SysML from an io.Reader.
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
	return ParseBytes(content, source, opts...)
}

func parseWithSource(input, source string, opts ...ParseOption) *ParseResult {
	cfg := &parseConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Use low-level parser
	tree, lowErrors := low.Parse(input)

	result := &ParseResult{
		Source: source,
	}

	if !cfg.discardTree {
		result.Tree = tree
	}

	// Convert errors
	if lowErrors.HasErrors() {
		result.Errors = convertFromLowLevel(lowErrors, source)
	}

	// Build model from parse tree
	if tree != nil {
		result.Model = buildModel(tree)
	}

	return result
}

// ParseDirectory parses all .sysml files in a directory.
// Use WithDiscardTree() option for large repositories to reduce memory.
func ParseDirectory(dir string, opts ...ParseOption) ([]*ParseResult, error) {
	var results []*ParseResult

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", path, err)
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".sysml") {
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

// ParseDirectoryParallel parses all .sysml files in a directory using multiple workers.
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
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".sysml") {
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
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".sysml") {
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
func buildModel(tree parser.IEntryRuleRootNamespaceContext) *Model {
	if tree == nil {
		return nil
	}

	model := NewModel()
	builder := &modelBuilder{
		model:        model,
		elementStack: make([]Element, 0, 16),
		packageStack: make([]*Package, 0, 8),
	}
	antlr.ParseTreeWalkerDefault.Walk(builder, tree)

	// Build index and resolve references
	model.BuildIndex()
	model.ResolveReferences()

	return model
}

// modelBuilder walks the parse tree and builds the model.
type modelBuilder struct {
	*parser.BaseSysMLv2ParserListener
	model        *Model
	currentPkg   *Package
	packageStack []*Package // Stack of packages for nested package handling
	elementStack []Element  // Stack of elements (parts, requirements, etc.) for parent tracking
}

// getCurrentParent returns the current parent element for adding children.
// It checks the element stack first (for nested elements like requirements/parts),
// then falls back to currentPkg (for package-level elements).
// Returns nil if there's no valid parent (top-level elements).
func (b *modelBuilder) getCurrentParent() Element {
	if len(b.elementStack) > 0 {
		return b.elementStack[len(b.elementStack)-1]
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

func (b *modelBuilder) EnterPackage_(ctx *parser.Package_Context) {
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

func (b *modelBuilder) ExitPackage_(ctx *parser.Package_Context) {
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

func (b *modelBuilder) EnterPartDefinition(ctx *parser.PartDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	part := NewPart(name, locationFromContext(ctx), true)
	part.parent = b.getCurrentParent()
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
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		if ident := ctx.Usage().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	part := NewPart(name, locationFromContext(ctx), false)
	part.parent = b.getCurrentParent()
	b.addToParent(part)

	// Push part onto stack for nested elements (parts can have attributes, ports, etc.)
	b.elementStack = append(b.elementStack, part)
}

func (b *modelBuilder) ExitPartUsage(ctx *parser.PartUsageContext) {
	if len(b.elementStack) > 0 {
		b.elementStack = b.elementStack[:len(b.elementStack)-1]
	}
}

func (b *modelBuilder) EnterItemDefinition(ctx *parser.ItemDefinitionContext) {
	name := ""
	if ctx.Definition() != nil && ctx.Definition().DefinitionDeclaration() != nil {
		if ident := ctx.Definition().DefinitionDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	item := NewItem(name, locationFromContext(ctx), true)
	item.parent = b.getCurrentParent()
	b.addToParent(item)
}

func (b *modelBuilder) EnterItemUsage(ctx *parser.ItemUsageContext) {
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

	item := NewItem(name, loc, false)

	if b.currentPkg != nil {
		item.parent = b.currentPkg
		b.currentPkg.AddChild(item)
	}
}

func (b *modelBuilder) EnterImport_(ctx *parser.Import_Context) {
	loc := Location{
		Line:   ctx.GetStart().GetLine(),
		Column: ctx.GetStart().GetColumn(),
	}

	namespace := ctx.GetText()
	imp := NewImport(namespace, loc)

	if b.currentPkg != nil {
		b.currentPkg.AddChild(imp)
	} else {
		b.model.AddImport(imp)
	}
}

func (b *modelBuilder) EnterRequirementDefinition(ctx *parser.RequirementDefinitionContext) {
	name := ""
	if decl := ctx.DefinitionDeclaration(); decl != nil {
		if ident := decl.Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	req := NewRequirement(name, locationFromContext(ctx), true)
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
	if ctx.ConstraintUsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	req := NewRequirement(name, locationFromContext(ctx), false)
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
	if ctx.ConstraintUsageDeclaration() != nil {
		if ident := ctx.ConstraintUsageDeclaration().UsageDeclaration().Identification(); ident != nil {
			name = extractName(ident)
		}
	}

	ver := NewVerification(name, locationFromContext(ctx), false)
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
	if ctx.Usage() != nil && ctx.Usage().UsageDeclaration() != nil {
		usageDecl := ctx.Usage().UsageDeclaration()
		if ident := usageDecl.Identification(); ident != nil {
			name = extractName(ident)
		}
		// If no name from identification, check for redefinition
		if name == "" {
			if featSpecPart := usageDecl.FeatureSpecializationPart(); featSpecPart != nil {
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
	if ident == nil {
		return ""
	}

	// Try to get the name (use index 0 for the primary name)
	names := ident.AllName()
	if len(names) > 0 {
		return names[0].GetText()
	}

	return ""
}

// extractRedefinitionName extracts the name from a redefinition in featureSpecializationPart
func extractRedefinitionName(featSpecPart parser.IFeatureSpecializationPartContext) string {
	if featSpecPart == nil {
		return ""
	}

	// Iterate through all feature specializations
	for _, featSpec := range featSpecPart.AllFeatureSpecialization() {
		// Check if this is a redefinition
		if redef := featSpec.Redefinitions(); redef != nil {
			// Get the first owned redefinition
			ownedRedefs := redef.AllOwnedRedefinition()
			if len(ownedRedefs) > 0 {
				// Get the qualified name
				if qname := ownedRedefs[0].QualifiedName(); qname != nil {
					return qname.GetText()
				}
			}
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
