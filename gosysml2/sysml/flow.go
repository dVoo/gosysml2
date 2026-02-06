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
