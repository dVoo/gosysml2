package sysml

// Flow represents a data or control flow in SysML.
// Flows connect elements and transfer data or control between them.
type Flow struct {
	baseElement
	IsDefinition bool
	// Flow ends (source and target)
	Source *FlowEnd
	Target *FlowEnd
	// Payload/feature information
	PayloadFeatures []Element
}

// FlowEnd represents an endpoint of a flow.
type FlowEnd struct {
	baseElement
	// Reference to the element this flow end connects to
	Reference Element
	// Feature being flowed (if any)
	Feature Element
}

// NewFlow creates a new Flow element.
func NewFlow(name string, loc Location, isDef bool) *Flow {
	return &Flow{
		baseElement: baseElement{
			name:     name,
			kind:     KindFlow,
			location: loc,
		},
		IsDefinition:    isDef,
		PayloadFeatures: make([]Element, 0),
	}
}

// NewFlowEnd creates a new FlowEnd element.
func NewFlowEnd(loc Location) *FlowEnd {
	return &FlowEnd{
		baseElement: baseElement{
			kind:     KindFlowEnd,
			location: loc,
		},
	}
}

// GetName returns the flow name.
func (f *Flow) GetName() string {
	return f.name
}

// GetKind returns the element kind.
func (f *Flow) GetKind() ElementKind {
	return f.kind
}

// GetLocation returns the source location.
func (f *Flow) GetLocation() Location {
	return f.location
}

// GetParent returns the parent element.
func (f *Flow) GetParent() Element {
	return f.parent
}

// SetParent sets the parent element.
func (f *Flow) SetParent(parent Element) {
	f.parent = parent
}

// Accept implements the visitor pattern.
func (f *Flow) Accept(v Visitor) bool {
	return v.VisitFlow(f)
}

// GetName returns the flow end name (always empty).
func (fe *FlowEnd) GetName() string {
	return ""
}

// GetKind returns the element kind.
func (fe *FlowEnd) GetKind() ElementKind {
	return fe.kind
}

// GetLocation returns the source location.
func (fe *FlowEnd) GetLocation() Location {
	return fe.location
}

// GetParent returns the parent element.
func (fe *FlowEnd) GetParent() Element {
	return fe.parent
}

// SetParent sets the parent element.
func (fe *FlowEnd) SetParent(parent Element) {
	fe.parent = parent
}

// Accept implements the visitor pattern.
func (fe *FlowEnd) Accept(v Visitor) bool {
	return v.VisitFlowEnd(fe)
}

// Flows returns all flows in the package.
func (p *Package) Flows() []*Flow {
	var flows []*Flow
	for _, child := range p.children {
		if f, ok := child.(*Flow); ok {
			flows = append(flows, f)
		}
	}
	return flows
}

// AddFlow adds a flow to the model.
func (m *Model) AddFlow(flow *Flow) {
	m.Flows = append(m.Flows, flow)
	m.Elements = append(m.Elements, flow)
}

// GetFlows returns all flows in the model.
func (m *Model) GetFlows() []*Flow {
	return m.Flows
}

// SuccessionFlow represents a succession flow usage in SysML.
// SuccessionFlowUsage is used in action bodies for control flow with succession semantics.
// It connects source and target elements with an optional guard condition.
type SuccessionFlow struct {
	baseElement
	// Source is the element that succeeds (the "from" end)
	Source Ref[Element]
	// Target is the element that is succeeded by (the "to" end)
	Target Ref[Element]
	// Guard is an optional guard condition expression
	Guard string

	// Unresolved references (used during parsing, before resolution)
	unresolvedSource string
	unresolvedTarget string
}

// isUsage marks SuccessionFlow as a usage element.
func (s *SuccessionFlow) isUsage() {}

// NewSuccessionFlow creates a new SuccessionFlow element.
func NewSuccessionFlow(name string, loc Location) *SuccessionFlow {
	return &SuccessionFlow{
		baseElement: baseElement{
			name:     name,
			kind:     KindSuccessionFlow,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedSource sets the unresolved source reference.
func (s *SuccessionFlow) SetUnresolvedSource(ref string) {
	s.unresolvedSource = ref
}

// SetUnresolvedTarget sets the unresolved target reference.
func (s *SuccessionFlow) SetUnresolvedTarget(ref string) {
	s.unresolvedTarget = ref
}

// GetName returns the succession flow name.
func (s *SuccessionFlow) GetName() string {
	return s.name
}

// GetKind returns the element kind.
func (s *SuccessionFlow) GetKind() ElementKind {
	return s.kind
}

// GetLocation returns the source location.
func (s *SuccessionFlow) GetLocation() Location {
	return s.location
}

// GetParent returns the parent element.
func (s *SuccessionFlow) GetParent() Element {
	return s.parent
}

// SetParent sets the parent element.
func (s *SuccessionFlow) SetParent(parent Element) {
	s.parent = parent
}

// Accept implements the visitor pattern.
func (s *SuccessionFlow) Accept(v Visitor) bool {
	return v.VisitSuccessionFlow(s)
}

// Type returns the type reference for usages (SuccessionFlow doesn't have a type).
func (s *SuccessionFlow) Type() Element {
	return nil
}
