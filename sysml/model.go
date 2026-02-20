package sysml

import (
	"fmt"
	"iter"
	"reflect"
	"strings"
)

// ElementKind represents the kind of a SysML element.
type ElementKind int

const (
	KindUnknown ElementKind = iota
	KindPackage
	KindPart
	KindItem
	KindPort
	KindAttribute
	KindConnection
	KindInterface
	KindAllocation
	KindAction
	KindState
	KindTransition
	KindConstraint
	KindRequirement
	KindConcern
	KindUseCase
	KindVerification
	KindAnalysis
	KindView
	KindViewpoint
	KindComment
	KindDoc
	KindMetadata
	KindImport
	KindAlias
	KindRendering
	KindMessage
	KindDependency
	KindFlow
	KindFlowEnd
	KindEnumeration
	KindEnumerationValue
	KindCalculation
	KindControlNode
	KindOccurrence
	KindCase
	KindIncludeUseCase
	KindConjugatedPort
	KindSuccessionFlow
	KindKerMLType
	KindKerMLFeature
)

// Exported string constants for element kind comparisons.
// Use these instead of magic strings for compile-time safety.
const (
	KindPartStr             = "part"
	KindPackageStr          = "package"
	KindItemStr             = "item"
	KindPortStr             = "port"
	KindAttributeStr        = "attribute"
	KindConnectionStr       = "connection"
	KindInterfaceStr        = "interface"
	KindAllocationStr       = "allocation"
	KindActionStr           = "action"
	KindStateStr            = "state"
	KindTransitionStr       = "transition"
	KindConstraintStr       = "constraint"
	KindRequirementStr      = "requirement"
	KindConcernStr          = "concern"
	KindUseCaseStr          = "use case"
	KindVerificationStr     = "verification"
	KindAnalysisStr         = "analysis"
	KindViewStr             = "view"
	KindViewpointStr        = "viewpoint"
	KindCommentStr          = "comment"
	KindDocStr              = "doc"
	KindMetadataStr         = "metadata"
	KindImportStr           = "import"
	KindAliasStr            = "alias"
	KindRenderingStr        = "rendering"
	KindMessageStr          = "message"
	KindDependencyStr       = "dependency"
	KindFlowStr             = "flow"
	KindFlowEndStr          = "flow end"
	KindEnumerationStr      = "enumeration"
	KindEnumerationValueStr = "enumeration value"
	KindCalculationStr      = "calculation"
	KindControlNodeStr      = "control node"
	KindOccurrenceStr       = "occurrence"
	KindCaseStr             = "case"
	KindIncludeUseCaseStr   = "include use case"
	KindConjugatedPortStr   = "conjugated port"
	KindSuccessionFlowStr   = "succession flow"
	KindKerMLTypeStr        = "kerml type"
	KindKerMLFeatureStr     = "kerml feature"
)

// String returns the string representation of the element kind.
func (k ElementKind) String() string {
	switch k {
	case KindPackage:
		return "package"
	case KindPart:
		return "part"
	case KindItem:
		return "item"
	case KindPort:
		return "port"
	case KindAttribute:
		return "attribute"
	case KindConnection:
		return "connection"
	case KindInterface:
		return "interface"
	case KindAllocation:
		return "allocation"
	case KindAction:
		return "action"
	case KindState:
		return "state"
	case KindConstraint:
		return "constraint"
	case KindRequirement:
		return "requirement"
	case KindConcern:
		return "concern"
	case KindUseCase:
		return "use case"
	case KindVerification:
		return "verification"
	case KindAnalysis:
		return "analysis"
	case KindView:
		return "view"
	case KindViewpoint:
		return "viewpoint"
	case KindComment:
		return "comment"
	case KindDoc:
		return "doc"
	case KindMetadata:
		return "metadata"
	case KindImport:
		return "import"
	case KindAlias:
		return "alias"
	case KindRendering:
		return "rendering"
	case KindMessage:
		return "message"
	case KindDependency:
		return "dependency"
	case KindFlow:
		return "flow"
	case KindFlowEnd:
		return "flow end"
	case KindEnumeration:
		return "enumeration"
	case KindEnumerationValue:
		return "enumeration value"
	case KindCalculation:
		return "calculation"
	case KindTransition:
		return "transition"
	case KindControlNode:
		return "control node"
	case KindOccurrence:
		return "occurrence"
	case KindCase:
		return "case"
	case KindIncludeUseCase:
		return "include use case"
	case KindConjugatedPort:
		return "conjugated port"
	case KindSuccessionFlow:
		return "succession flow"
	case KindKerMLType:
		return "kerml type"
	case KindKerMLFeature:
		return "kerml feature"
	default:
		return "unknown"
	}
}

// Location represents a source location in the input.
type Location struct {
	Line      int
	Column    int
	EndLine   int
	EndColumn int
}

// Ref represents a reference to another element.
// It can be either resolved (pointing to an actual element) or unresolved (just a name).
type Ref[T Element] struct {
	name     string // The reference name (qualified or simple)
	resolved T      // The resolved element (nil if unresolved)
}

// NewRef creates a new unresolved reference.
func NewRef[T Element](name string) Ref[T] {
	return Ref[T]{name: name}
}

// Name returns the reference name.
func (r Ref[T]) Name() string {
	return r.name
}

// EffectiveName returns the best available reference name.
// If resolved, it prefers the resolved element's name; otherwise it falls back to unresolved name.
func (r Ref[T]) EffectiveName() string {
	if r.IsResolved() {
		if resolved := r.Resolved(); any(resolved) != nil && resolved.Name() != "" {
			return resolved.Name()
		}
	}
	return r.name
}

// Resolved returns the resolved element, or nil if unresolved.
func (r Ref[T]) Resolved() T {
	return r.resolved
}

// IsResolved returns true if the reference has been resolved.
func (r Ref[T]) IsResolved() bool {
	// Check if resolved is non-nil (for interface types)
	var zero T
	return any(r.resolved) != any(zero)
}

// Resolve sets the resolved element.
func (r *Ref[T]) Resolve(elem T) {
	r.resolved = elem
}

// Element is the base interface for all SysML elements.
type Element interface {
	// Kind returns the kind of this element.
	Kind() ElementKind

	// Name returns the name of this element, or empty string if unnamed.
	Name() string

	// QualifiedName returns the fully qualified name of this element.
	QualifiedName() string

	// Location returns the source location of this element.
	Location() Location

	// Parent returns the parent element, or nil for root elements.
	Parent() Element

	// Children returns the child elements.
	Children() []Element

	// SetParent sets the parent element.
	SetParent(parent Element)

	// Documentation returns the documentation string for this element.
	Documentation() string
}

// Definition is the interface for definition elements (e.g., part def, requirement def).
type Definition interface {
	Element
	isDefinition()
}

// Usage is the interface for usage elements (e.g., part usage, requirement usage).
type Usage interface {
	Element
	// Type returns a reference to the definition this usage is typed by.
	Type() Element
	isUsage()
}

// baseElement provides common implementation for all elements.
type baseElement struct {
	kind              ElementKind
	name              string
	declaredShortName string
	location          Location
	parent            Element
	children          []Element
	documentation     string
}

func (e *baseElement) Kind() ElementKind         { return e.kind }
func (e *baseElement) Name() string              { return e.name }
func (e *baseElement) Location() Location        { return e.location }
func (e *baseElement) Parent() Element           { return e.parent }
func (e *baseElement) Children() []Element       { return e.children }
func (e *baseElement) SetParent(p Element)       { e.parent = p }
func (e *baseElement) Documentation() string     { return e.documentation }
func (e *baseElement) DeclaredShortName() string { return e.declaredShortName }

func (e *baseElement) setDeclaredShortName(name string) {
	e.declaredShortName = name
}

func (e *baseElement) SetDocumentation(doc string) {
	e.documentation = doc
}

func (e *baseElement) QualifiedName() string {
	if e.parent == nil {
		return e.name
	}
	parentQN := e.parent.QualifiedName()
	if parentQN == "" {
		return e.name
	}
	return parentQN + "::" + e.name
}

func isNilValue(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}

func (e *baseElement) addChild(child Element, parent Element) {
	if isNilValue(child) {
		return
	}
	e.children = append(e.children, child)
	child.SetParent(parent)
}

// Package represents a SysML package.
type Package struct {
	baseElement
	IsLibrary bool

	// Typed accessors for children
	packages     []*Package
	parts        []*Part
	requirements []*Requirement
	actions      []*Action
	imports      []*Import
	items        []*Item
	states       []*State
	connections  []*Connection
	interfaces   []*Interface
	allocations  []*Allocation
	views        []*View
	viewpoints   []*Viewpoint
	calculations []*Calculation
	enumerations []*Enumeration
	constraints  []*Constraint
	dependencies []*Dependency
	docs         []*Doc
}

// NewPackage creates a new Package element.
func NewPackage(name string, loc Location) *Package {
	return &Package{
		baseElement: baseElement{
			kind:     KindPackage,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		packages:     make([]*Package, 0),
		parts:        make([]*Part, 0),
		requirements: make([]*Requirement, 0),
		actions:      make([]*Action, 0),
		imports:      make([]*Import, 0),
		items:        make([]*Item, 0),
		states:       make([]*State, 0),
		connections:  make([]*Connection, 0),
		interfaces:   make([]*Interface, 0),
		allocations:  make([]*Allocation, 0),
		views:        make([]*View, 0),
		viewpoints:   make([]*Viewpoint, 0),
		calculations: make([]*Calculation, 0),
		enumerations: make([]*Enumeration, 0),
		constraints:  make([]*Constraint, 0),
		dependencies: make([]*Dependency, 0),
		docs:         make([]*Doc, 0),
	}
}

// AddChild adds a child element to the package with proper type tracking.
func (p *Package) AddChild(child Element) {
	p.baseElement.addChild(child, p)

	// Track by type for type-safe access
	switch c := child.(type) {
	case *Package:
		if c != nil {
			p.packages = append(p.packages, c)
		}
	case *Part:
		if c != nil {
			p.parts = append(p.parts, c)
		}
	case *Requirement:
		if c != nil {
			p.requirements = append(p.requirements, c)
		}
	case *Action:
		if c != nil {
			p.actions = append(p.actions, c)
		}
	case *Import:
		if c != nil {
			p.imports = append(p.imports, c)
		}
	case *Item:
		if c != nil {
			p.items = append(p.items, c)
		}
	case *State:
		if c != nil {
			p.states = append(p.states, c)
		}
	case *Connection:
		if c != nil {
			p.connections = append(p.connections, c)
		}
	case *Interface:
		if c != nil {
			p.interfaces = append(p.interfaces, c)
		}
	case *Allocation:
		if c != nil {
			p.allocations = append(p.allocations, c)
		}
	case *View:
		if c != nil {
			p.views = append(p.views, c)
		}
	case *Viewpoint:
		if c != nil {
			p.viewpoints = append(p.viewpoints, c)
		}
	case *Calculation:
		if c != nil {
			p.calculations = append(p.calculations, c)
		}
	case *Enumeration:
		if c != nil {
			p.enumerations = append(p.enumerations, c)
		}
	case *Constraint:
		if c != nil {
			p.constraints = append(p.constraints, c)
		}
	case *Dependency:
		if c != nil {
			p.dependencies = append(p.dependencies, c)
		}
	case *Doc:
		if c != nil {
			p.docs = append(p.docs, c)
		}
	}
}

// Packages returns all direct child packages.
func (p *Package) Packages() []*Package { return p.packages }

// Parts returns all direct child parts.
func (p *Package) Parts() []*Part { return p.parts }

// Requirements returns all direct child requirements.
func (p *Package) Requirements() []*Requirement { return p.requirements }

// Actions returns all direct child actions.
func (p *Package) Actions() []*Action { return p.actions }

// Imports returns all direct child imports.
func (p *Package) Imports() []*Import { return p.imports }

// AllElements returns an iterator over all direct child elements in the package.
// This provides a generic way to iterate over children without using type-specific accessors.
// For recursive traversal of all descendants, use visitor.All() or visitor.Walk().
func (p *Package) AllElements() iter.Seq[Element] {
	return func(yield func(Element) bool) {
		for _, child := range p.children {
			if !yield(child) {
				return
			}
		}
	}
}

// Items returns all direct child items.
func (p *Package) Items() []*Item { return p.items }

// States returns all direct child states.
func (p *Package) States() []*State { return p.states }

// Connections returns all direct child connections.
func (p *Package) Connections() []*Connection { return p.connections }

// Interfaces returns all direct child interfaces.
func (p *Package) Interfaces() []*Interface { return p.interfaces }

// Allocations returns all direct child allocations.
func (p *Package) Allocations() []*Allocation { return p.allocations }

// Views returns all direct child views.
func (p *Package) Views() []*View { return p.views }

// Viewpoints returns all direct child viewpoints.
func (p *Package) Viewpoints() []*Viewpoint { return p.viewpoints }

// Calculations returns all direct child calculations.
func (p *Package) Calculations() []*Calculation { return p.calculations }

// Enumerations returns all direct child enumerations.
func (p *Package) Enumerations() []*Enumeration { return p.enumerations }

// Constraints returns all direct child constraints.
func (p *Package) Constraints() []*Constraint { return p.constraints }

// Dependencies returns all direct child dependencies.
func (p *Package) Dependencies() []*Dependency { return p.dependencies }

// KerMLType represents a KerML classifier/type-style declaration.
type KerMLType struct {
	baseElement
	DeclarationKeyword string
	Specializes        []Ref[Element]
	unresolvedSupers   []string
	features           []*KerMLFeature
}

func (t *KerMLType) isDefinition() {}

// NewKerMLType creates a new KerMLType element.
func NewKerMLType(name, keyword string, loc Location) *KerMLType {
	return &KerMLType{
		baseElement: baseElement{
			kind:     KindKerMLType,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		DeclarationKeyword: keyword,
		Specializes:        make([]Ref[Element], 0),
		unresolvedSupers:   make([]string, 0),
		features:           make([]*KerMLFeature, 0),
	}
}

// AddChild adds a child element with type tracking.
func (t *KerMLType) AddChild(child Element) {
	t.baseElement.addChild(child, t)
	if feat, ok := child.(*KerMLFeature); ok && feat != nil {
		t.features = append(t.features, feat)
	}
}

// Features returns direct KerML feature children.
func (t *KerMLType) Features() []*KerMLFeature { return t.features }

// AddUnresolvedSuper adds an unresolved specialization reference.
func (t *KerMLType) AddUnresolvedSuper(name string) {
	if name == "" {
		return
	}
	t.unresolvedSupers = append(t.unresolvedSupers, name)
}

// KerMLFeature represents a feature declaration in KerML type bodies.
type KerMLFeature struct {
	baseElement
	TypeRef                 Ref[Element]
	DefaultValue            string
	SubsettedFeatures       []Element
	RedefinedFeatures       []Element
	unresolvedSubsetted     []string
	unresolvedRedefined     []string
	unresolvedTypeReference string
}

func (f *KerMLFeature) isDefinition() {}

// NewKerMLFeature creates a new KerMLFeature element.
func NewKerMLFeature(name string, loc Location) *KerMLFeature {
	return &KerMLFeature{
		baseElement: baseElement{
			kind:     KindKerMLFeature,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		SubsettedFeatures:   make([]Element, 0),
		RedefinedFeatures:   make([]Element, 0),
		unresolvedSubsetted: make([]string, 0),
		unresolvedRedefined: make([]string, 0),
	}
}

// AddUnresolvedSubsettedFeature adds unresolved subsetting target.
func (f *KerMLFeature) AddUnresolvedSubsettedFeature(name string) {
	if name == "" {
		return
	}
	f.unresolvedSubsetted = append(f.unresolvedSubsetted, name)
}

// AddUnresolvedRedefinedFeature adds unresolved redefinition target.
func (f *KerMLFeature) AddUnresolvedRedefinedFeature(name string) {
	if name == "" {
		return
	}
	f.unresolvedRedefined = append(f.unresolvedRedefined, name)
}

// Attribute represents a SysML attribute with name, type, and optional value.
type Attribute struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[Element] // Reference to the attribute type
	DefaultValue string       // Default value expression (as string for now)
	IsReadOnly   bool
	IsDerived    bool

	// Feature relationships
	SubsettedFeatures []Element // Features this attribute subsets (::>, :>)
	RedefinedFeatures []Element // Features this attribute redefines (:>>)

	unresolvedSubsettedFeatures []string
	unresolvedRedefinedFeatures []string
}

// NewAttribute creates a new Attribute element.
func NewAttribute(name string, loc Location, isDefinition bool) *Attribute {
	return &Attribute{
		baseElement: baseElement{
			kind:     KindAttribute,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:                isDefinition,
		SubsettedFeatures:           make([]Element, 0),
		RedefinedFeatures:           make([]Element, 0),
		unresolvedSubsettedFeatures: make([]string, 0),
		unresolvedRedefinedFeatures: make([]string, 0),
	}
}

func (a *Attribute) isDefinition() {}
func (a *Attribute) isUsage()      {}

// Type returns the type reference for usages.
func (a *Attribute) Type() Element {
	return a.TypeRef.Resolved()
}

// AddUnresolvedSubsettedFeature adds an unresolved subsetting/reference-subsetting name.
func (a *Attribute) AddUnresolvedSubsettedFeature(ref string) {
	if ref == "" {
		return
	}
	a.unresolvedSubsettedFeatures = append(a.unresolvedSubsettedFeatures, ref)
}

// AddUnresolvedRedefinedFeature adds an unresolved redefined feature name.
func (a *Attribute) AddUnresolvedRedefinedFeature(ref string) {
	if ref == "" {
		return
	}
	a.unresolvedRedefinedFeatures = append(a.unresolvedRedefinedFeatures, ref)
}

// Part represents a SysML part definition or usage.
type Part struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Part] // Reference to the part definition (for usages)
	Multiplicity string     // Usage multiplicity, e.g. "4", "0..1", "*"

	// Specialization (subclassification) reference for definitions using :> or "specializes"
	Specializes           Ref[*Part]
	unresolvedSpecializes string

	// Typed children
	attributes []*Attribute
	parts      []*Part
	ports      []*Port
}

func (p *Part) isDefinition() {}
func (p *Part) isUsage()      {}

// Type returns the type reference for usages.
func (p *Part) Type() Element {
	return p.TypeRef.Resolved()
}

// SetUnresolvedSpecializes sets the unresolved specialization reference.
func (p *Part) SetUnresolvedSpecializes(ref string) {
	p.unresolvedSpecializes = ref
}

// UnresolvedSpecializes returns the unresolved specialization reference.
func (p *Part) UnresolvedSpecializes() string {
	return p.unresolvedSpecializes
}

// NewPart creates a new Part element.
func NewPart(name string, loc Location, isDefinition bool) *Part {
	return &Part{
		baseElement: baseElement{
			kind:     KindPart,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		attributes:   make([]*Attribute, 0),
		parts:        make([]*Part, 0),
		ports:        make([]*Port, 0),
	}
}

// AddChild adds a child element with type tracking.
func (p *Part) AddChild(child Element) {
	p.baseElement.addChild(child, p)

	switch c := child.(type) {
	case *Attribute:
		if c != nil {
			p.attributes = append(p.attributes, c)
		}
	case *Part:
		if c != nil {
			p.parts = append(p.parts, c)
		}
	case *Port:
		if c != nil {
			p.ports = append(p.ports, c)
		}
	}
}

// Attributes returns all attributes of this part.
func (p *Part) Attributes() []*Attribute { return p.attributes }

// Parts returns all nested parts.
func (p *Part) Parts() []*Part { return p.parts }

// Ports returns all ports.
func (p *Part) Ports() []*Port { return p.ports }

// String returns a string representation of the part for debugging.
func (p *Part) String() string {
	typ := "definition"
	if !p.IsDefinition {
		typ = "usage"
	}
	multiplicity := ""
	if p.Multiplicity != "" {
		multiplicity = fmt.Sprintf("[%s]", p.Multiplicity)
	}
	specializes := ""
	if p.Specializes.IsResolved() {
		specializes = fmt.Sprintf(" -> %s", p.Specializes.Resolved().Name())
	} else if p.unresolvedSpecializes != "" {
		specializes = fmt.Sprintf(" -> %s (unresolved)", p.unresolvedSpecializes)
	}
	return fmt.Sprintf("Part<%s>{%s%s, attrs=%d, parts=%d, ports=%d}",
		typ, p.name, multiplicity+specializes, len(p.attributes), len(p.parts), len(p.ports))
}

// Port represents a SysML port definition or usage.
type Port struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Port]
	Direction    PortDirection

	// Typed children
	ports []*Port
	parts []*Part

	// Conjugated port - automatically created for each PortDefinition
	ConjugatedPort *ConjugatedPort
}

// ConjugatedPort represents a conjugated (inverted) port definition.
// Per SysML spec, every PortDefinition has an associated ConjugatedPortDefinition
// with effective name "~" + original port name.
type ConjugatedPort struct {
	baseElement
	OriginalPort           Ref[*Port] // Reference to the original port
	unresolvedOriginalPort string     // Name before resolution
}

// isDefinition marks ConjugatedPort as a definition element.
func (c *ConjugatedPort) isDefinition() {}

// GetOriginalPort returns the original port this conjugated port references.
func (c *ConjugatedPort) GetOriginalPort() *Port {
	return c.OriginalPort.Resolved()
}

// SetUnresolvedOriginalPort sets the unresolved reference to the original port.
func (c *ConjugatedPort) SetUnresolvedOriginalPort(ref string) {
	c.unresolvedOriginalPort = ref
}

// EffectiveName returns the effective name of the conjugated port (~Name).
func (c *ConjugatedPort) EffectiveName() string {
	if c.OriginalPort.IsResolved() {
		return "~" + c.OriginalPort.Resolved().Name()
	}
	// If name already has ~ prefix, return it as-is
	if len(c.name) > 0 && c.name[0] == '~' {
		return c.name
	}
	return "~" + c.name
}

// PortDirection indicates the direction of a port.
type PortDirection int

const (
	PortDirectionNone PortDirection = iota
	PortDirectionIn
	PortDirectionOut
	PortDirectionInOut
)

func (p *Port) isDefinition() {}
func (p *Port) isUsage()      {}

// Type returns the type reference for usages.
func (p *Port) Type() Element {
	return p.TypeRef.Resolved()
}

// NewPort creates a new Port element.
func NewPort(name string, loc Location, isDefinition bool) *Port {
	return &Port{
		baseElement: baseElement{
			kind:     KindPort,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		ports:        make([]*Port, 0),
		parts:        make([]*Part, 0),
	}
}

// NewConjugatedPort creates a new ConjugatedPort element.
func NewConjugatedPort(name string, loc Location) *ConjugatedPort {
	return &ConjugatedPort{
		baseElement: baseElement{
			kind:     KindConjugatedPort,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// AddChild adds a child element with type tracking.
func (p *Port) AddChild(child Element) {
	p.baseElement.addChild(child, p)

	switch c := child.(type) {
	case *Port:
		if c != nil {
			p.ports = append(p.ports, c)
		}
	case *Part:
		if c != nil {
			p.parts = append(p.parts, c)
		}
	}
}

// Ports returns nested ports.
func (p *Port) Ports() []*Port { return p.ports }

// Parts returns nested parts.
func (p *Port) Parts() []*Part { return p.parts }

// RequirementConstraint represents a constraint within a requirement (assume or require).
type RequirementConstraint struct {
	baseElement
	IsAssume   bool   // true for assume, false for require
	Expression string // The constraint expression
}

// NewRequirementConstraint creates a new constraint.
func NewRequirementConstraint(loc Location, isAssume bool, expr string) *RequirementConstraint {
	return &RequirementConstraint{
		baseElement: baseElement{
			kind:     KindConstraint,
			location: loc,
			children: make([]Element, 0),
		},
		IsAssume:   isAssume,
		Expression: expr,
	}
}

// Constraint represents a SysML constraint definition or usage.
type Constraint struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Constraint]
	Expression   string // Constraint expression

	// Typed children
	constraints []*Constraint
}

func (c *Constraint) isDefinition() {}
func (c *Constraint) isUsage()      {}

// Type returns the type reference for usages.
func (c *Constraint) Type() Element {
	return c.TypeRef.Resolved()
}

// NewConstraint creates a new Constraint element.
func NewConstraint(name string, loc Location, isDefinition bool) *Constraint {
	return &Constraint{
		baseElement: baseElement{
			kind:     KindConstraint,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		constraints:  make([]*Constraint, 0),
	}
}

// AddChild adds a child element with type tracking.
func (c *Constraint) AddChild(child Element) {
	c.baseElement.addChild(child, c)

	if nested, ok := child.(*Constraint); ok {
		if nested != nil {
			c.constraints = append(c.constraints, nested)
		}
	}
}

// Constraints returns nested constraints.
func (c *Constraint) Constraints() []*Constraint { return c.constraints }

// Requirement represents a SysML requirement definition or usage.
type Requirement struct {
	baseElement
	IsDefinition  bool
	TypeRef       Ref[*Requirement] // Reference to the requirement definition (for usages)
	RequirementID string            // Optional requirement ID (e.g., "REQ-001")
	Bindings      map[string]string // Optional usage argument bindings, e.g. [arg = value]

	// Subject
	Subject Ref[Element] // Reference to the subject element

	// Relationships with real references
	DerivedFrom []*Requirement  // Requirements this is derived from
	DerivedReqs []*Requirement  // Requirements derived from this one (inverse)
	SatisfiedBy []Element       // Elements that satisfy this requirement
	VerifiedBy  []*Verification // Verification cases that verify this

	// Constraints
	Assumptions []*RequirementConstraint // assume constraints
	Constraints []*RequirementConstraint // require constraints

	// Nested requirements
	requirements []*Requirement

	// Unresolved references (used during parsing, before resolution)
	unresolvedDerivedFrom []string
	unresolvedSatisfiedBy []string
	unresolvedVerifiedBy  []string
	unresolvedSubject     string
}

func (r *Requirement) isDefinition() {}
func (r *Requirement) isUsage()      {}

// Type returns the type reference for usages.
func (r *Requirement) Type() Element {
	return r.TypeRef.Resolved()
}

// NewRequirement creates a new Requirement element.
func NewRequirement(name string, loc Location, isDefinition bool) *Requirement {
	return &Requirement{
		baseElement: baseElement{
			kind:     KindRequirement,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:          isDefinition,
		Bindings:              make(map[string]string),
		DerivedFrom:           make([]*Requirement, 0),
		DerivedReqs:           make([]*Requirement, 0),
		SatisfiedBy:           make([]Element, 0),
		VerifiedBy:            make([]*Verification, 0),
		Assumptions:           make([]*RequirementConstraint, 0),
		Constraints:           make([]*RequirementConstraint, 0),
		requirements:          make([]*Requirement, 0),
		unresolvedDerivedFrom: make([]string, 0),
		unresolvedSatisfiedBy: make([]string, 0),
		unresolvedVerifiedBy:  make([]string, 0),
	}
}

// AddChild adds a child element with type tracking.
func (r *Requirement) AddChild(child Element) {
	r.baseElement.addChild(child, r)

	switch c := child.(type) {
	case *Requirement:
		if c != nil {
			r.requirements = append(r.requirements, c)
		}
	case *RequirementConstraint:
		if c != nil {
			if c.IsAssume {
				r.Assumptions = append(r.Assumptions, c)
			} else {
				r.Constraints = append(r.Constraints, c)
			}
		}
	}
}

// Requirements returns nested requirements.
func (r *Requirement) Requirements() []*Requirement { return r.requirements }

// Text returns the documentation text (requirement text).
func (r *Requirement) Text() string { return r.documentation }

// AddUnresolvedDerivedFrom adds an unresolved derivation reference.
func (r *Requirement) AddUnresolvedDerivedFrom(ref string) {
	r.unresolvedDerivedFrom = append(r.unresolvedDerivedFrom, ref)
}

// AddUnresolvedSatisfiedBy adds an unresolved satisfaction reference.
func (r *Requirement) AddUnresolvedSatisfiedBy(ref string) {
	r.unresolvedSatisfiedBy = append(r.unresolvedSatisfiedBy, ref)
}

// AddUnresolvedVerifiedBy adds an unresolved verification reference.
func (r *Requirement) AddUnresolvedVerifiedBy(ref string) {
	r.unresolvedVerifiedBy = append(r.unresolvedVerifiedBy, ref)
}

// SetUnresolvedSubject sets the unresolved subject reference.
func (r *Requirement) SetUnresolvedSubject(ref string) {
	r.unresolvedSubject = ref
}

// UnresolvedReferences returns all unresolved reference names for debugging.
func (r *Requirement) UnresolvedReferences() (derivedFrom, satisfiedBy, verifiedBy []string, subject string) {
	return r.unresolvedDerivedFrom, r.unresolvedSatisfiedBy, r.unresolvedVerifiedBy, r.unresolvedSubject
}

// Action represents a SysML action definition or usage.
type Action struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Action]

	// Typed children
	actions []*Action
}

func (a *Action) isDefinition() {}
func (a *Action) isUsage()      {}

// Type returns the type reference for usages.
func (a *Action) Type() Element {
	return a.TypeRef.Resolved()
}

// NewAction creates a new Action element.
func NewAction(name string, loc Location, isDefinition bool) *Action {
	return &Action{
		baseElement: baseElement{
			kind:     KindAction,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		actions:      make([]*Action, 0),
	}
}

// AddChild adds a child element with type tracking.
func (a *Action) AddChild(child Element) {
	a.baseElement.addChild(child, a)

	if c, ok := child.(*Action); ok {
		if c != nil {
			a.actions = append(a.actions, c)
		}
	}
}

// Actions returns nested actions.
func (a *Action) Actions() []*Action { return a.actions }

// Verification represents a SysML verification case definition or usage.
type Verification struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Verification]

	// Subject being verified
	Subject Ref[Element]

	// The requirement being verified (resolved reference)
	VerifiedRequirement *Requirement

	// Verification method
	Method VerificationMethod

	// Actions in the verification
	actions []*Action

	// Unresolved references
	unresolvedSubject     string
	unresolvedRequirement string
}

// VerificationMethod represents the method used for verification.
type VerificationMethod int

const (
	VerificationMethodUnspecified VerificationMethod = iota
	VerificationMethodTest
	VerificationMethodAnalysis
	VerificationMethodInspection
	VerificationMethodDemonstration
)

func (v VerificationMethod) String() string {
	switch v {
	case VerificationMethodTest:
		return "test"
	case VerificationMethodAnalysis:
		return "analysis"
	case VerificationMethodInspection:
		return "inspection"
	case VerificationMethodDemonstration:
		return "demonstration"
	default:
		return "unspecified"
	}
}

func (v *Verification) isDefinition() {}
func (v *Verification) isUsage()      {}

// Type returns the type reference for usages.
func (v *Verification) Type() Element {
	return v.TypeRef.Resolved()
}

// NewVerification creates a new Verification element.
func NewVerification(name string, loc Location, isDefinition bool) *Verification {
	return &Verification{
		baseElement: baseElement{
			kind:     KindVerification,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		actions:      make([]*Action, 0),
	}
}

// AddChild adds a child element with type tracking.
func (v *Verification) AddChild(child Element) {
	v.baseElement.addChild(child, v)

	if c, ok := child.(*Action); ok {
		if c != nil {
			v.actions = append(v.actions, c)
		}
	}
}

// Actions returns the actions in this verification.
func (v *Verification) Actions() []*Action { return v.actions }

// SetUnresolvedRequirement sets the unresolved requirement reference.
func (v *Verification) SetUnresolvedRequirement(ref string) {
	v.unresolvedRequirement = ref
}

// SetUnresolvedSubject sets the unresolved subject reference.
func (v *Verification) SetUnresolvedSubject(ref string) {
	v.unresolvedSubject = ref
}

// Concern represents a SysML concern definition or usage.
type Concern struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Concern]

	// Stakeholders (resolved references)
	Stakeholders []Element

	// Unresolved
	unresolvedStakeholders []string
}

func (c *Concern) isDefinition() {}
func (c *Concern) isUsage()      {}

// Type returns the type reference for usages.
func (c *Concern) Type() Element {
	return c.TypeRef.Resolved()
}

// NewConcern creates a new Concern element.
func NewConcern(name string, loc Location, isDefinition bool) *Concern {
	return &Concern{
		baseElement: baseElement{
			kind:     KindConcern,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:           isDefinition,
		Stakeholders:           make([]Element, 0),
		unresolvedStakeholders: make([]string, 0),
	}
}

// Text returns the documentation text (concern text).
func (c *Concern) Text() string { return c.documentation }

// UseCase represents a SysML use case definition or usage.
type UseCase struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*UseCase]

	// Subject (resolved reference)
	Subject Ref[Element]

	// Actors (resolved references)
	Actors []Element

	// Included use cases (resolved references)
	IncludedUseCases []*UseCase

	// Unresolved
	unresolvedSubject          string
	unresolvedActors           []string
	unresolvedIncludedUseCases []string
}

func (u *UseCase) isDefinition() {}
func (u *UseCase) isUsage()      {}

// Type returns the type reference for usages.
func (u *UseCase) Type() Element {
	return u.TypeRef.Resolved()
}

// NewUseCase creates a new UseCase element.
func NewUseCase(name string, loc Location, isDefinition bool) *UseCase {
	return &UseCase{
		baseElement: baseElement{
			kind:     KindUseCase,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:               isDefinition,
		Actors:                     make([]Element, 0),
		IncludedUseCases:           make([]*UseCase, 0),
		unresolvedActors:           make([]string, 0),
		unresolvedIncludedUseCases: make([]string, 0),
	}
}

// AddChild adds a child element to the use case.
func (u *UseCase) AddChild(child Element) {
	u.baseElement.addChild(child, u)
}

// IncludeUseCase represents a use case inclusion relationship.
// This is a usage element that represents the "include" relationship between use cases.
type IncludeUseCase struct {
	baseElement

	// IncludedUseCase is the reference to the included use case
	IncludedUseCase Ref[*UseCase]

	// unresolvedIncludedUseCase holds the name before resolution
	unresolvedIncludedUseCase string

	// Owner is the use case that includes this (set during parsing)
	Owner Ref[*UseCase]
}

// isUsage marks IncludeUseCase as a usage element.
func (i *IncludeUseCase) isUsage() {}

// Type returns the type reference for usages (IncludeUseCase doesn't have a type).
func (i *IncludeUseCase) Type() Element {
	return nil
}

// NewIncludeUseCase creates a new IncludeUseCase element.
func NewIncludeUseCase(name string, loc Location) *IncludeUseCase {
	return &IncludeUseCase{
		baseElement: baseElement{
			kind:     KindIncludeUseCase,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedIncludedUseCase sets the unresolved reference to the included use case.
func (i *IncludeUseCase) SetUnresolvedIncludedUseCase(ref string) {
	i.unresolvedIncludedUseCase = ref
}

// AnalysisCase represents a SysML analysis case definition or usage.
type AnalysisCase struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*AnalysisCase]

	// Subject (resolved reference)
	Subject Ref[Element]

	// Return type
	ReturnType Ref[Element]

	// Unresolved
	unresolvedSubject    string
	unresolvedReturnType string
}

func (a *AnalysisCase) isDefinition() {}
func (a *AnalysisCase) isUsage()      {}

// Type returns the type reference for usages.
func (a *AnalysisCase) Type() Element {
	return a.TypeRef.Resolved()
}

// NewAnalysisCase creates a new AnalysisCase element.
func NewAnalysisCase(name string, loc Location, isDefinition bool) *AnalysisCase {
	return &AnalysisCase{
		baseElement: baseElement{
			kind:     KindAnalysis,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
	}
}

// Case represents a SysML case definition or usage.
type Case struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Case]

	// Subject (resolved reference)
	Subject Ref[Element]

	// Actors (resolved references)
	Actors []Element

	// Objectives (resolved references to requirements)
	Objectives []*Requirement

	// Unresolved references
	unresolvedSubject    string
	unresolvedActors     []string
	unresolvedObjectives []string
}

func (c *Case) isDefinition() {}
func (c *Case) isUsage()      {}

// Type returns the type reference for usages.
func (c *Case) Type() Element {
	return c.TypeRef.Resolved()
}

// NewCase creates a new Case element.
func NewCase(name string, loc Location, isDefinition bool) *Case {
	return &Case{
		baseElement: baseElement{
			kind:     KindCase,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:         isDefinition,
		Actors:               make([]Element, 0),
		Objectives:           make([]*Requirement, 0),
		unresolvedActors:     make([]string, 0),
		unresolvedObjectives: make([]string, 0),
	}
}

// SetUnresolvedSubject sets the unresolved subject reference.
func (c *Case) SetUnresolvedSubject(ref string) {
	c.unresolvedSubject = ref
}

// AddUnresolvedActor adds an unresolved actor reference.
func (c *Case) AddUnresolvedActor(ref string) {
	c.unresolvedActors = append(c.unresolvedActors, ref)
}

// AddUnresolvedObjective adds an unresolved objective reference.
func (c *Case) AddUnresolvedObjective(ref string) {
	c.unresolvedObjectives = append(c.unresolvedObjectives, ref)
}

// AddChild adds a child element to the case.
func (c *Case) AddChild(child Element) {
	c.baseElement.addChild(child, c)
}

// EnumerationValue represents a single value within an enumeration.
type EnumerationValue struct {
	baseElement
}

// NewEnumerationValue creates a new EnumerationValue element.
func NewEnumerationValue(name string, loc Location) *EnumerationValue {
	return &EnumerationValue{
		baseElement: baseElement{
			kind:     KindEnumerationValue,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// Enumeration represents a SysML enumeration definition or usage.
type Enumeration struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Enumeration]

	// Enumerated values
	values []*EnumerationValue
}

func (e *Enumeration) isDefinition() {}
func (e *Enumeration) isUsage()      {}

// Type returns the type reference for usages.
func (e *Enumeration) Type() Element {
	return e.TypeRef.Resolved()
}

// NewEnumeration creates a new Enumeration element.
func NewEnumeration(name string, loc Location, isDefinition bool) *Enumeration {
	return &Enumeration{
		baseElement: baseElement{
			kind:     KindEnumeration,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		values:       make([]*EnumerationValue, 0),
	}
}

// AddChild adds a child element with type tracking.
func (e *Enumeration) AddChild(child Element) {
	e.baseElement.addChild(child, e)

	if v, ok := child.(*EnumerationValue); ok {
		if v != nil {
			e.values = append(e.values, v)
		}
	}
}

// Values returns the enumeration values.
func (e *Enumeration) Values() []*EnumerationValue { return e.values }

// Item represents a SysML item definition or usage.
type Item struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Item]

	// Specialization (subclassification) reference for definitions using :> or "specializes"
	Specializes           Ref[*Item]
	unresolvedSpecializes string

	// Typed children
	attributes []*Attribute
	items      []*Item

	// Feature relationships for item usages
	SubsettedFeatures []Element // Features this item subsets (::>, :>, subsets)
	RedefinedFeatures []Element // Features this item redefines (:>>, redefines)

	unresolvedSubsettedFeatures []string
	unresolvedRedefinedFeatures []string
}

func (i *Item) isDefinition() {}
func (i *Item) isUsage()      {}

// Type returns the type reference for usages.
func (i *Item) Type() Element {
	return i.TypeRef.Resolved()
}

// SetUnresolvedSpecializes sets the unresolved specialization reference.
func (i *Item) SetUnresolvedSpecializes(ref string) {
	i.unresolvedSpecializes = ref
}

// UnresolvedSpecializes returns the unresolved specialization reference.
func (i *Item) UnresolvedSpecializes() string {
	return i.unresolvedSpecializes
}

// NewItem creates a new Item element.
func NewItem(name string, loc Location, isDefinition bool) *Item {
	return &Item{
		baseElement: baseElement{
			kind:     KindItem,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:                isDefinition,
		attributes:                  make([]*Attribute, 0),
		items:                       make([]*Item, 0),
		SubsettedFeatures:           make([]Element, 0),
		RedefinedFeatures:           make([]Element, 0),
		unresolvedSubsettedFeatures: make([]string, 0),
		unresolvedRedefinedFeatures: make([]string, 0),
	}
}

// AddChild adds a child element with type tracking.
func (i *Item) AddChild(child Element) {
	i.baseElement.addChild(child, i)

	switch c := child.(type) {
	case *Attribute:
		if c != nil {
			i.attributes = append(i.attributes, c)
		}
	case *Item:
		if c != nil {
			i.items = append(i.items, c)
		}
	}
}

// Attributes returns all attributes of this item.
func (i *Item) Attributes() []*Attribute { return i.attributes }

// Items returns all nested items.
func (i *Item) Items() []*Item { return i.items }

// String returns a string representation of the item for debugging.
func (i *Item) String() string {
	typ := "definition"
	if !i.IsDefinition {
		typ = "usage"
	}
	specializes := ""
	if i.Specializes.IsResolved() {
		specializes = fmt.Sprintf(" -> %s", i.Specializes.Resolved().Name())
	} else if i.unresolvedSpecializes != "" {
		specializes = fmt.Sprintf(" -> %s (unresolved)", i.unresolvedSpecializes)
	}
	return fmt.Sprintf("Item<%s>{%s%s, attrs=%d, items=%d}",
		typ, i.name, specializes, len(i.attributes), len(i.items))
}

// AddUnresolvedSubsettedFeature adds an unresolved subsetting/reference-subsetting name.
func (i *Item) AddUnresolvedSubsettedFeature(ref string) {
	if ref == "" {
		return
	}
	i.unresolvedSubsettedFeatures = append(i.unresolvedSubsettedFeatures, ref)
}

// AddUnresolvedRedefinedFeature adds an unresolved redefined feature name.
func (i *Item) AddUnresolvedRedefinedFeature(ref string) {
	if ref == "" {
		return
	}
	i.unresolvedRedefinedFeatures = append(i.unresolvedRedefinedFeatures, ref)
}

// Calculation represents a SysML calculation definition or usage.
type Calculation struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Calculation]
	ReturnType   Ref[Element] // Return type reference
	Expression   string       // Calculation expression

	unresolvedReturnType string
}

func (c *Calculation) isDefinition() {}
func (c *Calculation) isUsage()      {}

// Type returns the type reference for usages.
func (c *Calculation) Type() Element {
	return c.TypeRef.Resolved()
}

// NewCalculation creates a new Calculation element.
func NewCalculation(name string, loc Location, isDefinition bool) *Calculation {
	return &Calculation{
		baseElement: baseElement{
			kind:     KindCalculation,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
	}
}

// SetUnresolvedReturnType sets the unresolved return type reference.
func (c *Calculation) SetUnresolvedReturnType(ref string) {
	c.unresolvedReturnType = ref
}

// State represents a SysML state definition or usage.
type State struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*State]

	// Entry/do/exit actions
	EntryAction *Action
	DoAction    *Action
	ExitAction  *Action

	// Nested states and transitions
	states      []*State
	transitions []*Transition
}

func (s *State) isDefinition() {}
func (s *State) isUsage()      {}

// Type returns the type reference for usages.
func (s *State) Type() Element {
	return s.TypeRef.Resolved()
}

// NewState creates a new State element.
func NewState(name string, loc Location, isDefinition bool) *State {
	return &State{
		baseElement: baseElement{
			kind:     KindState,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		states:       make([]*State, 0),
		transitions:  make([]*Transition, 0),
	}
}

// AddChild adds a child element with type tracking.
func (s *State) AddChild(child Element) {
	s.baseElement.addChild(child, s)

	switch c := child.(type) {
	case *State:
		if c != nil {
			s.states = append(s.states, c)
		}
	case *Transition:
		if c != nil {
			s.transitions = append(s.transitions, c)
		}
	}
}

// States returns nested states.
func (s *State) States() []*State { return s.states }

// Transitions returns transitions.
func (s *State) Transitions() []*Transition { return s.transitions }

// Transition represents a SysML transition usage.
type Transition struct {
	baseElement
	Source       Ref[*State] // Source state reference
	Target       Ref[*State] // Target state reference
	TriggerExpr  string      // Trigger expression
	GuardExpr    string      // Guard expression
	EffectAction *Action     // Effect action

	unresolvedSource string
	unresolvedTarget string
}

// NewTransition creates a new Transition element.
func NewTransition(name string, loc Location) *Transition {
	return &Transition{
		baseElement: baseElement{
			kind:     KindTransition,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedSource sets the unresolved source state reference.
func (t *Transition) SetUnresolvedSource(ref string) {
	t.unresolvedSource = ref
}

// SetUnresolvedTarget sets the unresolved target state reference.
func (t *Transition) SetUnresolvedTarget(ref string) {
	t.unresolvedTarget = ref
}

// ConnectionEnd represents an endpoint of a connection.
type ConnectionEnd struct {
	baseElement
	EndRef Ref[Element] // Reference to the connected element (Part/Port)

	unresolvedEndRef string
}

// NewConnectionEnd creates a new ConnectionEnd element.
func NewConnectionEnd(name string, loc Location) *ConnectionEnd {
	return &ConnectionEnd{
		baseElement: baseElement{
			kind:     KindConnection,
			name:     name,
			location: loc,
		},
	}
}

// Connection represents a SysML connection definition or usage.
type Connection struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Connection]
	Ends         []*ConnectionEnd // Connection endpoints

	unresolvedEnds []string
}

func (c *Connection) isDefinition() {}
func (c *Connection) isUsage()      {}

// Type returns the type reference for usages.
func (c *Connection) Type() Element {
	return c.TypeRef.Resolved()
}

// NewConnection creates a new Connection element.
func NewConnection(name string, loc Location, isDefinition bool) *Connection {
	return &Connection{
		baseElement: baseElement{
			kind:     KindConnection,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:   isDefinition,
		Ends:           make([]*ConnectionEnd, 0),
		unresolvedEnds: make([]string, 0),
	}
}

// AddUnresolvedEnd adds an unresolved end reference.
func (c *Connection) AddUnresolvedEnd(ref string) {
	c.unresolvedEnds = append(c.unresolvedEnds, ref)
}

// Interface represents a SysML interface definition or usage.
type Interface struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Interface]

	// Typed children
	ports []*Port
}

func (i *Interface) isDefinition() {}
func (i *Interface) isUsage()      {}

// Type returns the type reference for usages.
func (i *Interface) Type() Element {
	return i.TypeRef.Resolved()
}

// NewInterface creates a new Interface element.
func NewInterface(name string, loc Location, isDefinition bool) *Interface {
	return &Interface{
		baseElement: baseElement{
			kind:     KindInterface,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		ports:        make([]*Port, 0),
	}
}

// AddChild adds a child element with type tracking.
func (i *Interface) AddChild(child Element) {
	i.baseElement.addChild(child, i)

	if p, ok := child.(*Port); ok {
		if p != nil {
			i.ports = append(i.ports, p)
		}
	}
}

// Ports returns all ports.
func (i *Interface) Ports() []*Port { return i.ports }

// Allocation represents a SysML allocation definition or usage.
type Allocation struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Allocation]
	Source       Ref[Element] // Allocated element
	Target       Ref[Element] // Allocating element

	unresolvedSource string
	unresolvedTarget string
}

func (a *Allocation) isDefinition() {}
func (a *Allocation) isUsage()      {}

// Type returns the type reference for usages.
func (a *Allocation) Type() Element {
	return a.TypeRef.Resolved()
}

// NewAllocation creates a new Allocation element.
func NewAllocation(name string, loc Location, isDefinition bool) *Allocation {
	return &Allocation{
		baseElement: baseElement{
			kind:     KindAllocation,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
	}
}

// SetUnresolvedSource sets the unresolved source reference.
func (a *Allocation) SetUnresolvedSource(ref string) {
	a.unresolvedSource = ref
}

// SetUnresolvedTarget sets the unresolved target reference.
func (a *Allocation) SetUnresolvedTarget(ref string) {
	a.unresolvedTarget = ref
}

// Viewpoint represents a SysML viewpoint definition or usage.
type Viewpoint struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Viewpoint]
	Concerns     []*Concern // Related concerns
	Stakeholders []Element  // Related stakeholders

	unresolvedConcerns     []string
	unresolvedStakeholders []string
}

func (v *Viewpoint) isDefinition() {}
func (v *Viewpoint) isUsage()      {}

// Type returns the type reference for usages.
func (v *Viewpoint) Type() Element {
	return v.TypeRef.Resolved()
}

// NewViewpoint creates a new Viewpoint element.
func NewViewpoint(name string, loc Location, isDefinition bool) *Viewpoint {
	return &Viewpoint{
		baseElement: baseElement{
			kind:     KindViewpoint,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:           isDefinition,
		Concerns:               make([]*Concern, 0),
		Stakeholders:           make([]Element, 0),
		unresolvedConcerns:     make([]string, 0),
		unresolvedStakeholders: make([]string, 0),
	}
}

// AddUnresolvedConcern adds an unresolved concern reference.
func (v *Viewpoint) AddUnresolvedConcern(ref string) {
	v.unresolvedConcerns = append(v.unresolvedConcerns, ref)
}

// AddUnresolvedStakeholder adds an unresolved stakeholder reference.
func (v *Viewpoint) AddUnresolvedStakeholder(ref string) {
	v.unresolvedStakeholders = append(v.unresolvedStakeholders, ref)
}

// View represents a SysML view definition or usage.
type View struct {
	baseElement
	IsDefinition    bool
	TypeRef         Ref[*View]
	ExposedElements []Element       // Elements exposed by this view
	Viewpoint       Ref[*Viewpoint] // Associated viewpoint
	Exposures       []ViewExposure

	unresolvedExposedElements []string
	unresolvedViewpoint       string
}

// ViewExposure describes one textual `expose ...` clause in a view body.
type ViewExposure struct {
	Namespace         string
	IsMembership      bool
	IsNamespace       bool
	IsAll             bool
	IsRecursive       bool
	FilterExpressions []string
}

func (v *View) isDefinition() {}
func (v *View) isUsage()      {}

// Type returns the type reference for usages.
func (v *View) Type() Element {
	return v.TypeRef.Resolved()
}

// NewView creates a new View element.
func NewView(name string, loc Location, isDefinition bool) *View {
	return &View{
		baseElement: baseElement{
			kind:     KindView,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition:              isDefinition,
		ExposedElements:           make([]Element, 0),
		Exposures:                 make([]ViewExposure, 0),
		unresolvedExposedElements: make([]string, 0),
	}
}

// AddUnresolvedExposedElement adds an unresolved exposed element reference.
func (v *View) AddUnresolvedExposedElement(ref string) {
	v.unresolvedExposedElements = append(v.unresolvedExposedElements, ref)
}

// SetUnresolvedViewpoint sets the unresolved viewpoint reference.
func (v *View) SetUnresolvedViewpoint(ref string) {
	v.unresolvedViewpoint = ref
}

// AddExposure records a parsed textual `expose` clause.
func (v *View) AddExposure(exposure ViewExposure) {
	v.Exposures = append(v.Exposures, exposure)
}

// Import represents a SysML import statement.
type Import struct {
	baseElement
	ImportedNamespace string
	Visibility        string   // public/private/protected/default
	IsMembership      bool     // import member(s)
	IsNamespace       bool     // import namespace/package
	IsRecursive       bool     // true for ::**
	IsAll             bool     // true for ::*
	FilterExpressions []string // filter package expressions
	ResolvedElement   Element  // The resolved imported element (if single import)
	ResolvedPackage   *Package // The resolved library package (if namespace is a library package)
	IsResolved        bool     // true if the import was successfully resolved
}

// NewImport creates a new Import element.
func NewImport(namespace string, loc Location) *Import {
	return &Import{
		baseElement: baseElement{
			kind:     KindImport,
			name:     "",
			location: loc,
		},
		ImportedNamespace: namespace,
		FilterExpressions: make([]string, 0),
	}
}

// Model represents a complete SysML model (root namespace).
type Model struct {
	// Typed top-level element collections
	Packages     []*Package
	Imports      []*Import
	Comments     []*Comment
	Dependencies []*Dependency
	Docs         []*Doc
	Flows        []*Flow
	ControlNodes []*ControlNode
	Occurrences  []*Occurrence
	Aliases      []*Alias
	Metadata     []*Metadata
	Renderings   []*Rendering
	Messages     []*Message
	Filters      []*ElementFilter
	Satisfies    []*SatisfyRelationship
	Verifies     []*VerifyRelationship

	// All top-level elements (for generic traversal)
	Elements []Element

	// Index for fast lookup by qualified name
	elementIndex map[string]Element
	// Index for resolving elements by declared short name (e.g. <'R1'>)
	shortNameIndex map[string][]Element

	// Library registry for resolving qualified names to library elements
	libraryRegistry *LibraryRegistry
}

// NewModel creates a new empty model.
func NewModel() *Model {
	return &Model{
		Packages:       make([]*Package, 0),
		Imports:        make([]*Import, 0),
		Comments:       make([]*Comment, 0),
		Dependencies:   make([]*Dependency, 0),
		Docs:           make([]*Doc, 0),
		Flows:          make([]*Flow, 0),
		ControlNodes:   make([]*ControlNode, 0),
		Occurrences:    make([]*Occurrence, 0),
		Aliases:        make([]*Alias, 0),
		Metadata:       make([]*Metadata, 0),
		Renderings:     make([]*Rendering, 0),
		Messages:       make([]*Message, 0),
		Filters:        make([]*ElementFilter, 0),
		Satisfies:      make([]*SatisfyRelationship, 0),
		Verifies:       make([]*VerifyRelationship, 0),
		Elements:       make([]Element, 0),
		elementIndex:   make(map[string]Element),
		shortNameIndex: make(map[string][]Element),
	}
}

// AddDoc adds a documentation element to the model.
func (m *Model) AddDoc(doc *Doc) {
	if doc == nil {
		return
	}
	m.Docs = append(m.Docs, doc)
	m.Elements = append(m.Elements, doc)
}

// AddDependency adds a dependency to the model.
func (m *Model) AddDependency(dep *Dependency) {
	if dep == nil {
		return
	}
	m.Dependencies = append(m.Dependencies, dep)
	m.Elements = append(m.Elements, dep)
}

// AddPackage adds a package to the model.
func (m *Model) AddPackage(pkg *Package) {
	if pkg == nil {
		return
	}
	m.Packages = append(m.Packages, pkg)
	m.Elements = append(m.Elements, pkg)
}

// AddImport adds an import to the model.
func (m *Model) AddImport(imp *Import) {
	if imp == nil {
		return
	}
	m.Imports = append(m.Imports, imp)
	m.Elements = append(m.Elements, imp)
}

// AddComment adds a comment to the model.
func (m *Model) AddComment(comment *Comment) {
	if comment == nil {
		return
	}
	m.Comments = append(m.Comments, comment)
	m.Elements = append(m.Elements, comment)
}

// AddControlNode adds a control node to the model.
func (m *Model) AddControlNode(node *ControlNode) {
	if node == nil {
		return
	}
	m.ControlNodes = append(m.ControlNodes, node)
	m.Elements = append(m.Elements, node)
}

// AddOccurrence adds an occurrence to the model.
func (m *Model) AddOccurrence(occ *Occurrence) {
	if occ == nil {
		return
	}
	m.Occurrences = append(m.Occurrences, occ)
	m.Elements = append(m.Elements, occ)
}

// AddAlias adds an alias to the model.
func (m *Model) AddAlias(alias *Alias) {
	if alias == nil {
		return
	}
	m.Aliases = append(m.Aliases, alias)
	m.Elements = append(m.Elements, alias)
}

// AddMetadata adds metadata to the model.
func (m *Model) AddMetadata(metadata *Metadata) {
	if metadata == nil {
		return
	}
	m.Metadata = append(m.Metadata, metadata)
	m.Elements = append(m.Elements, metadata)
}

// AddRendering adds rendering to the model.
func (m *Model) AddRendering(rendering *Rendering) {
	if rendering == nil {
		return
	}
	m.Renderings = append(m.Renderings, rendering)
	m.Elements = append(m.Elements, rendering)
}

// AddMessage adds message usage to the model.
func (m *Model) AddMessage(message *Message) {
	if message == nil {
		return
	}
	m.Messages = append(m.Messages, message)
	m.Elements = append(m.Elements, message)
}

// AddFilter adds an element filter to the model.
func (m *Model) AddFilter(filter *ElementFilter) {
	if filter == nil {
		return
	}
	m.Filters = append(m.Filters, filter)
	m.Elements = append(m.Elements, filter)
}

// AddSatisfy adds a satisfy relationship to the model.
func (m *Model) AddSatisfy(rel *SatisfyRelationship) {
	if rel == nil {
		return
	}
	m.Satisfies = append(m.Satisfies, rel)
	// Avoid duplicate traversal: nested relationships are already reachable
	// via parent.Children() and should not be added as extra model roots.
	if rel.Parent() == nil {
		m.Elements = append(m.Elements, rel)
	}
}

// AddVerify adds a verify relationship to the model.
func (m *Model) AddVerify(rel *VerifyRelationship) {
	if rel == nil {
		return
	}
	m.Verifies = append(m.Verifies, rel)
	// Avoid duplicate traversal: nested relationships are already reachable
	// via parent.Children() and should not be added as extra model roots.
	if rel.Parent() == nil {
		m.Elements = append(m.Elements, rel)
	}
}

// FindPackage finds a package by name.
func (m *Model) FindPackage(name string) *Package {
	for _, pkg := range m.Packages {
		if pkg.Name() == name {
			return pkg
		}
	}
	return nil
}

// SetLibraryRegistry sets the library registry for resolving qualified names.
// This enables the model to resolve references to standard library elements.
func (m *Model) SetLibraryRegistry(reg *LibraryRegistry) {
	m.libraryRegistry = reg
}

// FindByQualifiedName finds an element by its fully qualified name.
func (m *Model) FindByQualifiedName(qn string) Element {
	return m.elementIndex[qn]
}

// BuildIndex builds the element index for fast lookups.
// This should be called after parsing is complete.
func (m *Model) BuildIndex() {
	m.elementIndex = make(map[string]Element)
	m.shortNameIndex = make(map[string][]Element)
	m.Walk(func(elem Element) bool {
		qn := elem.QualifiedName()
		if qn != "" {
			m.elementIndex[qn] = elem
		}
		if snElem, ok := elem.(interface{ DeclaredShortName() string }); ok {
			sn := snElem.DeclaredShortName()
			if sn != "" {
				m.shortNameIndex[sn] = append(m.shortNameIndex[sn], elem)
			}
		}
		return true
	})
}

// ResolveReferences resolves all element references in the model.
// This should be called after parsing and BuildIndex.
func (m *Model) ResolveReferences() {
	m.Walk(func(elem Element) bool {
		switch e := elem.(type) {
		case *Requirement:
			m.resolveRequirementRefs(e)
		case *Verification:
			m.resolveVerificationRefs(e)
		case *Part:
			m.resolvePartRefs(e)
		case *Item:
			m.resolveItemRefs(e)
		case *UseCase:
			m.resolveUseCaseRefs(e)
		case *Concern:
			m.resolveConcernRefs(e)
		case *AnalysisCase:
			m.resolveAnalysisCaseRefs(e)
		case *Case:
			m.resolveCaseRefs(e)
		case *IncludeUseCase:
			m.resolveIncludeUseCaseRefs(e)
		case *Transition:
			m.resolveTransitionRefs(e)
		case *Connection:
			m.resolveConnectionRefs(e)
		case *Allocation:
			m.resolveAllocationRefs(e)
		case *View:
			m.resolveViewRefs(e)
		case *Viewpoint:
			m.resolveViewpointRefs(e)
		case *Calculation:
			m.resolveCalculationRefs(e)
		case *State:
			m.resolveStateRefs(e)
		case *Action:
			m.resolveActionRefs(e)
		case *Constraint:
			m.resolveConstraintRefs(e)
		case *Enumeration:
			m.resolveEnumerationRefs(e)
		case *Port:
			m.resolvePortRefs(e)
		case *ConjugatedPort:
			m.resolveConjugatedPortRefs(e)
		case *Attribute:
			m.resolveAttributeRefs(e)
		case *Interface:
			m.resolveInterfaceRefs(e)
		case *SuccessionFlow:
			m.resolveSuccessionFlowRefs(e)
		case *Dependency:
			m.resolveDependencyRefs(e)
		case *Comment:
			m.resolveCommentRefs(e)
		case *Alias:
			m.resolveAliasRefs(e)
		case *Metadata:
			m.resolveMetadataRefs(e)
		case *Rendering:
			m.resolveRenderingRefs(e)
		case *Message:
			m.resolveMessageRefs(e)
		case *KerMLType:
			m.resolveKerMLTypeRefs(e)
		case *KerMLFeature:
			m.resolveKerMLFeatureRefs(e)
		case *SatisfyRelationship:
			m.resolveSatisfyRelationshipRefs(e)
		case *VerifyRelationship:
			m.resolveVerifyRelationshipRefs(e)
		}
		return true
	})
}

func (m *Model) resolveRequirementRefs(r *Requirement) {
	// Resolve derived from
	for _, name := range r.unresolvedDerivedFrom {
		if elem := m.findElement(name, r); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				if req != nil {
					r.DerivedFrom = append(r.DerivedFrom, req)
				}
				// Also set inverse relationship
				if req != nil && r != nil {
					req.DerivedReqs = append(req.DerivedReqs, r)
				}
			}
		}
	}

	// Resolve satisfied by
	for _, name := range r.unresolvedSatisfiedBy {
		if elem := m.findElement(name, r); elem != nil {
			r.SatisfiedBy = append(r.SatisfiedBy, elem)
		}
	}

	// Resolve verified by
	for _, name := range r.unresolvedVerifiedBy {
		if elem := m.findElement(name, r); elem != nil {
			if ver, ok := elem.(*Verification); ok {
				if ver != nil {
					r.VerifiedBy = append(r.VerifiedBy, ver)
				}
			}
		}
	}

	// Resolve subject
	if r.unresolvedSubject != "" {
		if elem := m.findElement(r.unresolvedSubject, r); elem != nil {
			r.Subject.Resolve(elem)
		}
	}

	// Resolve type reference for usages
	if !r.IsDefinition && r.TypeRef.name != "" {
		if elem := m.findElement(r.TypeRef.name, r); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				r.TypeRef.Resolve(req)
			}
		}
	}
}

func (m *Model) resolveVerificationRefs(v *Verification) {
	// Resolve verified requirement
	if v.unresolvedRequirement != "" {
		if elem := m.findElement(v.unresolvedRequirement, v); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				v.VerifiedRequirement = req
				// Also set inverse relationship
				if req != nil && v != nil {
					req.VerifiedBy = append(req.VerifiedBy, v)
				}
			}
		}
	}

	// Resolve subject
	if v.unresolvedSubject != "" {
		if elem := m.findElement(v.unresolvedSubject, v); elem != nil {
			v.Subject.Resolve(elem)
		}
	}

	// Resolve type reference
	if !v.IsDefinition && v.TypeRef.name != "" {
		if elem := m.findElement(v.TypeRef.name, v); elem != nil {
			if ver, ok := elem.(*Verification); ok {
				v.TypeRef.Resolve(ver)
			}
		}
	}
}

func (m *Model) resolvePartRefs(p *Part) {
	// Resolve TypeRef for usages (e.g., "part x : Vehicle")
	if !p.IsDefinition && p.TypeRef.name != "" {
		if elem := m.findElement(p.TypeRef.name, p); elem != nil {
			if part, ok := elem.(*Part); ok {
				p.TypeRef.Resolve(part)
			}
		}
	}

	// Resolve Specializes reference for definitions (e.g., "part def Car :> Vehicle")
	if p.IsDefinition && p.unresolvedSpecializes != "" {
		if elem := m.findElement(p.unresolvedSpecializes, p); elem != nil {
			if part, ok := elem.(*Part); ok {
				p.Specializes.Resolve(part)
			}
		}
	}
}

func (m *Model) resolveUseCaseRefs(u *UseCase) {
	// Resolve subject
	if u.unresolvedSubject != "" {
		if elem := m.findElement(u.unresolvedSubject, u); elem != nil {
			u.Subject.Resolve(elem)
		}
	}

	// Resolve actors
	for _, name := range u.unresolvedActors {
		if elem := m.findElement(name, u); elem != nil {
			u.Actors = append(u.Actors, elem)
		}
	}

	// Resolve included use cases
	for _, name := range u.unresolvedIncludedUseCases {
		if elem := m.findElement(name, u); elem != nil {
			if uc, ok := elem.(*UseCase); ok {
				if uc != nil {
					u.IncludedUseCases = append(u.IncludedUseCases, uc)
				}
			}
		}
	}

	// Resolve type reference
	if !u.IsDefinition && u.TypeRef.name != "" {
		if elem := m.findElement(u.TypeRef.name, u); elem != nil {
			if uc, ok := elem.(*UseCase); ok {
				u.TypeRef.Resolve(uc)
			}
		}
	}
}

func (m *Model) resolveConcernRefs(c *Concern) {
	// Resolve stakeholders
	for _, name := range c.unresolvedStakeholders {
		if elem := m.findElement(name, c); elem != nil {
			c.Stakeholders = append(c.Stakeholders, elem)
		}
	}

	// Resolve type reference
	if !c.IsDefinition && c.TypeRef.name != "" {
		if elem := m.findElement(c.TypeRef.name, c); elem != nil {
			if concern, ok := elem.(*Concern); ok {
				c.TypeRef.Resolve(concern)
			}
		}
	}
}

func (m *Model) resolveAnalysisCaseRefs(a *AnalysisCase) {
	// Resolve subject
	if a.unresolvedSubject != "" {
		if elem := m.findElement(a.unresolvedSubject, a); elem != nil {
			a.Subject.Resolve(elem)
		}
	}

	// Resolve return type
	if a.unresolvedReturnType != "" {
		if elem := m.findElement(a.unresolvedReturnType, a); elem != nil {
			a.ReturnType.Resolve(elem)
		}
	}

	// Resolve type reference
	if !a.IsDefinition && a.TypeRef.name != "" {
		if elem := m.findElement(a.TypeRef.name, a); elem != nil {
			if ac, ok := elem.(*AnalysisCase); ok {
				a.TypeRef.Resolve(ac)
			}
		}
	}
}

func (m *Model) resolveCaseRefs(c *Case) {
	// Resolve subject
	if c.unresolvedSubject != "" {
		if elem := m.findElement(c.unresolvedSubject, c); elem != nil {
			c.Subject.Resolve(elem)
		}
	}

	// Resolve actors
	for _, name := range c.unresolvedActors {
		if elem := m.findElement(name, c); elem != nil {
			c.Actors = append(c.Actors, elem)
		}
	}

	// Resolve objectives (requirements)
	for _, name := range c.unresolvedObjectives {
		if elem := m.findElement(name, c); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				if req != nil {
					c.Objectives = append(c.Objectives, req)
				}
			}
		}
	}

	// Resolve type reference for usages
	if !c.IsDefinition && c.TypeRef.name != "" {
		if elem := m.findElement(c.TypeRef.name, c); elem != nil {
			if cas, ok := elem.(*Case); ok {
				c.TypeRef.Resolve(cas)
			}
		}
	}
}

func (m *Model) resolveIncludeUseCaseRefs(i *IncludeUseCase) {
	// Resolve the included use case reference
	if i.unresolvedIncludedUseCase != "" {
		if elem := m.findElement(i.unresolvedIncludedUseCase, i); elem != nil {
			if uc, ok := elem.(*UseCase); ok {
				i.IncludedUseCase.Resolve(uc)
			}
		}
	}

	// Resolve the owner reference if set
	if i.Owner.name != "" && !i.Owner.IsResolved() {
		if elem := m.findElement(i.Owner.name, i); elem != nil {
			if uc, ok := elem.(*UseCase); ok {
				i.Owner.Resolve(uc)
			}
		}
	}
}

func (m *Model) resolveItemRefs(i *Item) {
	// Resolve TypeRef for usages
	if !i.IsDefinition && i.TypeRef.name != "" {
		if elem := m.findElement(i.TypeRef.name, i); elem != nil {
			if item, ok := elem.(*Item); ok {
				i.TypeRef.Resolve(item)
			}
		}
	}

	// Resolve Specializes reference for definitions
	if i.IsDefinition && i.unresolvedSpecializes != "" {
		if elem := m.findElement(i.unresolvedSpecializes, i); elem != nil {
			if item, ok := elem.(*Item); ok {
				i.Specializes.Resolve(item)
			}
		}
	}

	for _, name := range i.unresolvedSubsettedFeatures {
		if elem := m.findElement(name, i); elem != nil {
			i.SubsettedFeatures = append(i.SubsettedFeatures, elem)
		}
	}

	for _, name := range i.unresolvedRedefinedFeatures {
		if elem := m.findElement(name, i); elem != nil {
			i.RedefinedFeatures = append(i.RedefinedFeatures, elem)
		}
	}
}

func (m *Model) resolveTransitionRefs(t *Transition) {
	// Resolve source state
	if t.unresolvedSource != "" {
		if elem := m.findElement(t.unresolvedSource, t); elem != nil {
			if state, ok := elem.(*State); ok {
				t.Source.Resolve(state)
			}
		}
	}

	// Resolve target state
	if t.unresolvedTarget != "" {
		if elem := m.findElement(t.unresolvedTarget, t); elem != nil {
			if state, ok := elem.(*State); ok {
				t.Target.Resolve(state)
			}
		}
	}
}

func (m *Model) resolveConnectionRefs(c *Connection) {
	// Resolve connection ends
	for _, endRef := range c.unresolvedEnds {
		if elem := m.findElement(endRef, c); elem != nil {
			// Create a ConnectionEnd for this reference
			end := NewConnectionEnd("", c.location)
			end.EndRef.Resolve(elem)
			c.Ends = append(c.Ends, end)
		}
	}

	// Resolve type reference
	if !c.IsDefinition && c.TypeRef.name != "" {
		if elem := m.findElement(c.TypeRef.name, c); elem != nil {
			if conn, ok := elem.(*Connection); ok {
				c.TypeRef.Resolve(conn)
			}
		}
	}
}

func (m *Model) resolveAllocationRefs(a *Allocation) {
	// Resolve source
	if a.unresolvedSource != "" {
		if elem := m.findElement(a.unresolvedSource, a); elem != nil {
			a.Source.Resolve(elem)
		}
	}

	// Resolve target
	if a.unresolvedTarget != "" {
		if elem := m.findElement(a.unresolvedTarget, a); elem != nil {
			a.Target.Resolve(elem)
		}
	}

	// Resolve type reference
	if !a.IsDefinition && a.TypeRef.name != "" {
		if elem := m.findElement(a.TypeRef.name, a); elem != nil {
			if alloc, ok := elem.(*Allocation); ok {
				a.TypeRef.Resolve(alloc)
			}
		}
	}
}

func (m *Model) resolveViewRefs(v *View) {
	// Resolve exposed elements
	for _, expRef := range v.unresolvedExposedElements {
		if elem := m.findElement(expRef, v); elem != nil {
			v.ExposedElements = append(v.ExposedElements, elem)
		}
	}

	// Resolve textual `expose ...` clauses, including namespace/member wildcards.
	for _, exposure := range v.Exposures {
		candidates := m.resolveExposureCandidates(v, exposure)
		for _, elem := range candidates {
			v.ExposedElements = append(v.ExposedElements, elem)
		}
	}
	v.ExposedElements = dedupeElements(v.ExposedElements)

	// Resolve viewpoint
	if v.unresolvedViewpoint != "" {
		if elem := m.findElement(v.unresolvedViewpoint, v); elem != nil {
			if vp, ok := elem.(*Viewpoint); ok {
				v.Viewpoint.Resolve(vp)
			}
		}
	}

	// Resolve type reference
	if !v.IsDefinition && v.TypeRef.name != "" {
		if elem := m.findElement(v.TypeRef.name, v); elem != nil {
			if view, ok := elem.(*View); ok {
				v.TypeRef.Resolve(view)
			}
		}
	}
}

func (m *Model) resolveExposureCandidates(view *View, exposure ViewExposure) []Element {
	if exposure.Namespace == "" {
		return nil
	}
	scope := m.findElement(exposure.Namespace, view)
	if scope == nil {
		return nil
	}

	var candidates []Element
	switch {
	case exposure.IsRecursive:
		candidates = collectDescendants(scope)
	case exposure.IsAll:
		candidates = append(candidates, scope.Children()...)
	default:
		candidates = append(candidates, scope)
	}
	if len(exposure.FilterExpressions) == 0 {
		return candidates
	}

	filtered := make([]Element, 0, len(candidates))
	for _, elem := range candidates {
		if matchesAnyExposureFilter(elem, exposure.FilterExpressions) {
			filtered = append(filtered, elem)
		}
	}
	return filtered
}

func collectDescendants(root Element) []Element {
	var out []Element
	var walk func(Element)
	walk = func(e Element) {
		for _, child := range e.Children() {
			if child == nil {
				continue
			}
			out = append(out, child)
			walk(child)
		}
	}
	walk(root)
	return out
}

func dedupeElements(in []Element) []Element {
	if len(in) < 2 {
		return in
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]Element, 0, len(in))
	for _, elem := range in {
		if elem == nil {
			continue
		}
		key := elem.QualifiedName()
		if key == "" {
			loc := elem.Location()
			key = fmt.Sprintf("%s@%d:%d", elem.Name(), loc.Line, loc.Column)
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, elem)
	}
	return out
}

func matchesAnyExposureFilter(elem Element, filters []string) bool {
	for _, expr := range filters {
		if matchesExposureFilter(elem, expr) {
			return true
		}
	}
	return false
}

func matchesExposureFilter(elem Element, expr string) bool {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "@") {
		expr = strings.TrimPrefix(expr, "@")
	}
	// Handle common standard-library annotations like @SysML::PartUsage.
	switch expr {
	case "SysML::PartUsage":
		part, ok := elem.(*Part)
		return ok && !part.IsDefinition
	case "SysML::PartDefinition":
		part, ok := elem.(*Part)
		return ok && part.IsDefinition
	case "SysML::Part":
		_, ok := elem.(*Part)
		return ok
	case "SysML::RequirementUsage":
		req, ok := elem.(*Requirement)
		return ok && !req.IsDefinition
	case "SysML::RequirementDefinition":
		req, ok := elem.(*Requirement)
		return ok && req.IsDefinition
	case "SysML::Requirement":
		_, ok := elem.(*Requirement)
		return ok
	case "SysML::VerificationUsage":
		ver, ok := elem.(*Verification)
		return ok && !ver.IsDefinition
	case "SysML::VerificationDefinition":
		ver, ok := elem.(*Verification)
		return ok && ver.IsDefinition
	case "SysML::Verification":
		_, ok := elem.(*Verification)
		return ok
	}
	return false
}

func (m *Model) resolveViewpointRefs(v *Viewpoint) {
	// Resolve concerns
	for _, cRef := range v.unresolvedConcerns {
		if elem := m.findElement(cRef, v); elem != nil {
			if concern, ok := elem.(*Concern); ok {
				if concern != nil {
					v.Concerns = append(v.Concerns, concern)
				}
			}
		}
	}

	// Resolve stakeholders
	for _, sRef := range v.unresolvedStakeholders {
		if elem := m.findElement(sRef, v); elem != nil {
			v.Stakeholders = append(v.Stakeholders, elem)
		}
	}

	// Resolve type reference
	if !v.IsDefinition && v.TypeRef.name != "" {
		if elem := m.findElement(v.TypeRef.name, v); elem != nil {
			if vp, ok := elem.(*Viewpoint); ok {
				v.TypeRef.Resolve(vp)
			}
		}
	}
}

func (m *Model) resolveCalculationRefs(c *Calculation) {
	// Resolve return type
	if c.unresolvedReturnType != "" {
		if elem := m.findElement(c.unresolvedReturnType, c); elem != nil {
			c.ReturnType.Resolve(elem)
		}
	}

	// Resolve type reference
	if !c.IsDefinition && c.TypeRef.name != "" {
		if elem := m.findElement(c.TypeRef.name, c); elem != nil {
			if calc, ok := elem.(*Calculation); ok {
				c.TypeRef.Resolve(calc)
			}
		}
	}
}

func (m *Model) resolveStateRefs(s *State) {
	// Resolve type reference
	if !s.IsDefinition && s.TypeRef.name != "" {
		if elem := m.findElement(s.TypeRef.name, s); elem != nil {
			if state, ok := elem.(*State); ok {
				s.TypeRef.Resolve(state)
			}
		}
	}
}

func (m *Model) resolveActionRefs(a *Action) {
	// Resolve type reference
	if !a.IsDefinition && a.TypeRef.name != "" {
		if elem := m.findElement(a.TypeRef.name, a); elem != nil {
			if action, ok := elem.(*Action); ok {
				a.TypeRef.Resolve(action)
			}
		}
	}
}

func (m *Model) resolveConstraintRefs(c *Constraint) {
	// Resolve type reference
	if !c.IsDefinition && c.TypeRef.name != "" {
		if elem := m.findElement(c.TypeRef.name, c); elem != nil {
			if constraint, ok := elem.(*Constraint); ok {
				c.TypeRef.Resolve(constraint)
			}
		}
	}
}

func (m *Model) resolveEnumerationRefs(e *Enumeration) {
	// Resolve type reference
	if !e.IsDefinition && e.TypeRef.name != "" {
		if elem := m.findElement(e.TypeRef.name, e); elem != nil {
			if enum, ok := elem.(*Enumeration); ok {
				e.TypeRef.Resolve(enum)
			}
		}
	}
}

func (m *Model) resolvePortRefs(p *Port) {
	// Resolve type reference
	if !p.IsDefinition && p.TypeRef.name != "" {
		if elem := m.findElement(p.TypeRef.name, p); elem != nil {
			if port, ok := elem.(*Port); ok {
				p.TypeRef.Resolve(port)
			}
		}
	}
}

func (m *Model) resolveConjugatedPortRefs(c *ConjugatedPort) {
	// Resolve original port reference if not already resolved
	if !c.OriginalPort.IsResolved() && c.unresolvedOriginalPort != "" {
		if elem := m.findElement(c.unresolvedOriginalPort, c); elem != nil {
			if port, ok := elem.(*Port); ok {
				c.OriginalPort.Resolve(port)
			}
		}
	}
}

func (m *Model) resolveAttributeRefs(a *Attribute) {
	// Resolve type reference
	if a.TypeRef.name != "" {
		if elem := m.findElement(a.TypeRef.name, a); elem != nil {
			a.TypeRef.Resolve(elem)
		}
	}

	for _, name := range a.unresolvedSubsettedFeatures {
		if elem := m.findElement(name, a); elem != nil {
			a.SubsettedFeatures = append(a.SubsettedFeatures, elem)
		}
	}

	for _, name := range a.unresolvedRedefinedFeatures {
		if elem := m.findElement(name, a); elem != nil {
			a.RedefinedFeatures = append(a.RedefinedFeatures, elem)
		}
	}
}

func (m *Model) resolveInterfaceRefs(i *Interface) {
	// Resolve type reference
	if !i.IsDefinition && i.TypeRef.name != "" {
		if elem := m.findElement(i.TypeRef.name, i); elem != nil {
			if iface, ok := elem.(*Interface); ok {
				i.TypeRef.Resolve(iface)
			}
		}
	}
}

func (m *Model) resolveSuccessionFlowRefs(s *SuccessionFlow) {
	// Resolve source reference
	if s.unresolvedSource != "" {
		if elem := m.findElement(s.unresolvedSource, s); elem != nil {
			s.Source.Resolve(elem)
		}
	}

	// Resolve target reference
	if s.unresolvedTarget != "" {
		if elem := m.findElement(s.unresolvedTarget, s); elem != nil {
			s.Target.Resolve(elem)
		}
	}
}

func (m *Model) resolveDependencyRefs(d *Dependency) {
	for _, name := range d.unresolvedClient {
		if elem := m.findElement(name, d); elem != nil {
			d.Client = append(d.Client, elem)
		}
	}
	for _, name := range d.unresolvedSupplier {
		if elem := m.findElement(name, d); elem != nil {
			d.Supplier = append(d.Supplier, elem)
		}
	}
}

func (m *Model) resolveCommentRefs(c *Comment) {
	for _, name := range c.unresolvedAbout {
		if elem := m.findElement(name, c); elem != nil {
			c.About = append(c.About, elem)
		}
	}
}

func (m *Model) resolveAliasRefs(a *Alias) {
	if a.unresolvedTarget == "" {
		return
	}
	if elem := m.findElement(a.unresolvedTarget, a); elem != nil {
		a.Target.Resolve(elem)
	}
}

func (m *Model) resolveMetadataRefs(md *Metadata) {
	if !md.IsDefinition && md.TypeRef.name != "" {
		if elem := m.findElement(md.TypeRef.name, md); elem != nil {
			if target, ok := elem.(*Metadata); ok {
				md.TypeRef.Resolve(target)
			}
		}
	}
	for _, annotation := range md.annotations {
		if annotation.unresolvedMetadata == "" {
			continue
		}
		if elem := m.findElement(annotation.unresolvedMetadata, md); elem != nil {
			if target, ok := elem.(*Metadata); ok {
				if target != nil {
					annotation.Metadata.Resolve(target)
				}
			}
		}
	}
}

func (m *Model) resolveRenderingRefs(r *Rendering) {
	if !r.IsDefinition && r.TypeRef.name != "" {
		if elem := m.findElement(r.TypeRef.name, r); elem != nil {
			if target, ok := elem.(*Rendering); ok {
				r.TypeRef.Resolve(target)
			}
		}
	}
}

func (m *Model) resolveMessageRefs(msg *Message) {
	if msg.unresolvedSender != "" {
		if elem := m.findElement(msg.unresolvedSender, msg); elem != nil {
			msg.Sender.Resolve(elem)
		}
	}
	if msg.unresolvedReceiver != "" {
		if elem := m.findElement(msg.unresolvedReceiver, msg); elem != nil {
			msg.Receiver.Resolve(elem)
		}
	}
}

func (m *Model) resolveKerMLTypeRefs(t *KerMLType) {
	for _, ref := range t.unresolvedSupers {
		if elem := m.findElement(ref, t.Parent()); elem != nil {
			resolved := NewRef[Element](ref)
			resolved.Resolve(elem)
			t.Specializes = append(t.Specializes, resolved)
		}
	}
}

func (m *Model) resolveKerMLFeatureRefs(f *KerMLFeature) {
	if f.unresolvedTypeReference != "" {
		if elem := m.findElement(f.unresolvedTypeReference, f.Parent()); elem != nil {
			f.TypeRef = NewRef[Element](f.unresolvedTypeReference)
			f.TypeRef.Resolve(elem)
		}
	}
	for _, ref := range f.unresolvedSubsetted {
		if elem := m.findElement(ref, f.Parent()); elem != nil {
			f.SubsettedFeatures = append(f.SubsettedFeatures, elem)
		}
	}
	for _, ref := range f.unresolvedRedefined {
		if elem := m.findElement(ref, f.Parent()); elem != nil {
			f.RedefinedFeatures = append(f.RedefinedFeatures, elem)
		}
	}
}

func (m *Model) resolveSatisfyRelationshipRefs(rel *SatisfyRelationship) {
	if rel.unresolvedSatisfier != "" {
		if elem := m.findElement(rel.unresolvedSatisfier, rel); elem != nil {
			rel.Satisfier.Resolve(elem)
		}
	}
	if rel.unresolvedRequired != "" {
		if elem := m.findElement(rel.unresolvedRequired, rel); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				rel.Required.Resolve(req)
				if rel.Satisfier.IsResolved() {
					if satisfier := rel.Satisfier.Resolved(); !isNilValue(satisfier) && req != nil {
						req.SatisfiedBy = append(req.SatisfiedBy, satisfier)
					}
				}
			}
		}
	}
}

func (m *Model) resolveVerifyRelationshipRefs(rel *VerifyRelationship) {
	if rel.unresolvedVerifier != "" {
		if elem := m.findElement(rel.unresolvedVerifier, rel); elem != nil {
			if ver, ok := elem.(*Verification); ok {
				rel.Verifier.Resolve(ver)
			}
		}
	}
	if rel.unresolvedRequired != "" {
		if elem := m.findElement(rel.unresolvedRequired, rel); elem != nil {
			if req, ok := elem.(*Requirement); ok {
				rel.Required.Resolve(req)
				if rel.Verifier.IsResolved() {
					if verifier := rel.Verifier.Resolved(); verifier != nil && req != nil {
						req.VerifiedBy = append(req.VerifiedBy, verifier)
					}
				}
			}
		}
	}
}

// findElement finds an element by name, considering scope.
// It first tries qualified name lookup in the model, then searches in parent scopes,
// and finally falls back to the library registry if available.
// User definitions in the model always take precedence over library definitions.
func (m *Model) findElement(name string, context Element) Element {
	tryByQualifiedLookup := func(candidate string) Element {
		// First try direct qualified name lookup in model (user definitions take precedence)
		if elem := m.elementIndex[candidate]; elem != nil {
			return elem
		}

		// Try relative to context
		if context != nil {
			// Walk up the parent chain
			current := context.Parent()
			for current != nil {
				qn := current.QualifiedName()
				if qn != "" {
					fullQN := qn + "::" + candidate
					if elem := m.elementIndex[fullQN]; elem != nil {
						return elem
					}
				}
				current = current.Parent()
			}
		}

		// Try as simple name in any package
		for _, pkg := range m.Packages {
			fullQN := pkg.Name() + "::" + candidate
			if elem := m.elementIndex[fullQN]; elem != nil {
				return elem
			}
		}
		return nil
	}

	if elem := tryByQualifiedLookup(name); elem != nil {
		return elem
	}

	// Support dotted feature-chain notation by mapping dots to namespace separators.
	// Example: vehicle1.engine1 -> vehicle1::engine1
	if strings.Contains(name, ".") {
		if elem := tryByQualifiedLookup(strings.ReplaceAll(name, ".", "::")); elem != nil {
			return elem
		}
	}

	// Resolve by short name (declaredShortName in grammar).
	if candidates := m.shortNameIndex[name]; len(candidates) > 0 {
		if len(candidates) == 1 || context == nil {
			return candidates[0]
		}
		// Prefer a candidate in the nearest containing scope.
		for cur := context.Parent(); cur != nil; cur = cur.Parent() {
			curQN := cur.QualifiedName()
			for _, elem := range candidates {
				if curQN == "" || strings.HasPrefix(elem.QualifiedName(), curQN+"::") {
					return elem
				}
			}
		}
		return candidates[0]
	}

	// Finally, try library registry if available
	// This allows qualified names like "ISQ::mass" or "ScalarValues::Real" to resolve
	if m.libraryRegistry != nil {
		if elem := m.libraryRegistry.FindElement(name); elem != nil {
			return elem
		}
	}

	return nil
}

// Walk visits all elements in the model depth-first.
func (m *Model) Walk(fn func(Element) bool) {
	for _, elem := range m.Elements {
		if !walkElement(elem, fn) {
			return
		}
	}
}

// WalkAll visits all elements in the model depth-first.
// Use Walk when you need short-circuit traversal.
func (m *Model) WalkAll(fn func(Element)) {
	m.Walk(func(elem Element) bool {
		fn(elem)
		return true
	})
}

func walkElement(elem Element, fn func(Element) bool) bool {
	if !fn(elem) {
		return false
	}
	for _, child := range elem.Children() {
		if !walkElement(child, fn) {
			return false
		}
	}
	return true
}
