package sysml

// Dependency represents a SysML dependency relationship between elements.
// Dependencies declare that one or more client elements depend on one or more supplier elements.
type Dependency struct {
	baseElement
	// Client elements (the elements that depend on suppliers)
	Client []Element
	// Supplier elements (the elements being depended upon)
	Supplier []Element
	// Unresolved client references (stored during parsing, resolved later)
	unresolvedClient []string
	// Unresolved supplier references (stored during parsing, resolved later)
	unresolvedSupplier []string
}

// NewDependency creates a new Dependency element.
// Dependencies are typically unnamed, so this takes only a location.
func NewDependency(loc Location) *Dependency {
	return &Dependency{
		baseElement: baseElement{
			kind:     KindDependency,
			location: loc,
			children: make([]Element, 0),
		},
		Client:             make([]Element, 0),
		Supplier:           make([]Element, 0),
		unresolvedClient:   make([]string, 0),
		unresolvedSupplier: make([]string, 0),
	}
}

// Kind returns the kind of this element.
func (d *Dependency) Kind() ElementKind {
	return KindDependency
}

// Name returns the name of this element.
// Dependencies are typically unnamed, so this returns an empty string.
func (d *Dependency) Name() string {
	return ""
}

// GetKind returns the string representation of the element kind.
func (d *Dependency) GetKind() string {
	return KindDependency.String()
}

// GetLocation returns the source location of this element.
func (d *Dependency) GetLocation() Location {
	return d.location
}

// GetParent returns the parent element, or nil for root elements.
func (d *Dependency) GetParent() Element {
	return d.parent
}

// SetParent sets the parent element.
func (d *Dependency) SetParent(parent Element) {
	d.parent = parent
}

// Accept implements the visitor pattern.
func (d *Dependency) Accept(v Visitor) {
	v.VisitDependency(d)
}

// AddUnresolvedClient adds an unresolved client reference.
// This is used during parsing before references are resolved.
func (d *Dependency) AddUnresolvedClient(ref string) {
	d.unresolvedClient = append(d.unresolvedClient, ref)
}

// AddUnresolvedSupplier adds an unresolved supplier reference.
// This is used during parsing before references are resolved.
func (d *Dependency) AddUnresolvedSupplier(ref string) {
	d.unresolvedSupplier = append(d.unresolvedSupplier, ref)
}

// UnresolvedReferences returns all unresolved reference names for debugging.
func (d *Dependency) UnresolvedReferences() (client, supplier []string) {
	return d.unresolvedClient, d.unresolvedSupplier
}

// isDefinition marks Dependency as a definition element.
func (d *Dependency) isDefinition() {}

// isUsage marks Dependency as a usage element.
func (d *Dependency) isUsage() {}

func (d *Dependency) Role() ElementRole { return RoleDefinitionAndUsage }

// Type returns nil for dependencies (they don't have type references like usages).
func (d *Dependency) Type() Element {
	return nil
}
