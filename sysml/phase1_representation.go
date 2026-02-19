package sysml

// Alias represents an alias member that names another element.
type Alias struct {
	baseElement
	Target Ref[Element]

	unresolvedTarget string
}

// NewAlias creates a new Alias element.
func NewAlias(name string, loc Location) *Alias {
	return &Alias{
		baseElement: baseElement{
			kind:     KindAlias,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedTarget sets the unresolved target reference.
func (a *Alias) SetUnresolvedTarget(ref string) {
	a.unresolvedTarget = ref
}

// Metadata represents a metadata definition or usage.
type Metadata struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Metadata]

	annotations []*PrefixMetadataAnnotation
}

func (m *Metadata) isDefinition() {}
func (m *Metadata) isUsage()      {}

// Type returns the metadata type reference for usages.
func (m *Metadata) Type() Element {
	if !m.TypeRef.IsResolved() {
		return nil
	}
	return m.TypeRef.Resolved()
}

// NewMetadata creates a new Metadata element.
func NewMetadata(name string, loc Location, isDefinition bool) *Metadata {
	return &Metadata{
		baseElement: baseElement{
			kind:     KindMetadata,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
		annotations:  make([]*PrefixMetadataAnnotation, 0),
	}
}

// AddChild adds a child element to metadata.
func (m *Metadata) AddChild(child Element) {
	m.baseElement.addChild(child, m)

	if annotation, ok := child.(*PrefixMetadataAnnotation); ok {
		m.annotations = append(m.annotations, annotation)
	}
}

// Annotations returns attached metadata annotations.
func (m *Metadata) Annotations() []*PrefixMetadataAnnotation {
	return m.annotations
}

// Rendering represents a rendering definition or usage.
type Rendering struct {
	baseElement
	IsDefinition bool
	TypeRef      Ref[*Rendering]
	Body         string
}

func (r *Rendering) isDefinition() {}
func (r *Rendering) isUsage()      {}

// Type returns the rendering type reference for usages.
func (r *Rendering) Type() Element {
	if !r.TypeRef.IsResolved() {
		return nil
	}
	return r.TypeRef.Resolved()
}

// NewRendering creates a new Rendering element.
func NewRendering(name string, loc Location, isDefinition bool) *Rendering {
	return &Rendering{
		baseElement: baseElement{
			kind:     KindRendering,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDefinition,
	}
}

// Message represents a message usage between sender and receiver elements.
type Message struct {
	baseElement
	Sender   Ref[Element]
	Receiver Ref[Element]
	Payload  string

	unresolvedSender   string
	unresolvedReceiver string
}

func (m *Message) isUsage() {}

// Type returns nil because messages are not typed usages.
func (m *Message) Type() Element {
	return nil
}

// NewMessage creates a new Message usage element.
func NewMessage(name string, loc Location) *Message {
	return &Message{
		baseElement: baseElement{
			kind:     KindMessage,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedSender sets the unresolved sender reference.
func (m *Message) SetUnresolvedSender(ref string) {
	m.unresolvedSender = ref
}

// SetUnresolvedReceiver sets the unresolved receiver reference.
func (m *Message) SetUnresolvedReceiver(ref string) {
	m.unresolvedReceiver = ref
}

// SatisfyRelationship is a typed edge from a satisfier element to a requirement.
type SatisfyRelationship struct {
	baseElement
	Satisfier Ref[Element]
	Required  Ref[*Requirement]

	unresolvedSatisfier string
	unresolvedRequired  string
}

// NewSatisfyRelationship creates a satisfy relationship.
func NewSatisfyRelationship(name string, loc Location) *SatisfyRelationship {
	return &SatisfyRelationship{
		baseElement: baseElement{
			kind:     KindDependency,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedSatisfier sets the unresolved satisfier reference.
func (r *SatisfyRelationship) SetUnresolvedSatisfier(ref string) {
	r.unresolvedSatisfier = ref
}

// SetUnresolvedRequired sets the unresolved requirement reference.
func (r *SatisfyRelationship) SetUnresolvedRequired(ref string) {
	r.unresolvedRequired = ref
}

// VerifyRelationship is a typed edge from a verification case to a requirement.
type VerifyRelationship struct {
	baseElement
	Verifier Ref[*Verification]
	Required Ref[*Requirement]

	unresolvedVerifier string
	unresolvedRequired string
}

// NewVerifyRelationship creates a verify relationship.
func NewVerifyRelationship(name string, loc Location) *VerifyRelationship {
	return &VerifyRelationship{
		baseElement: baseElement{
			kind:     KindDependency,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
	}
}

// SetUnresolvedVerifier sets the unresolved verifier reference.
func (r *VerifyRelationship) SetUnresolvedVerifier(ref string) {
	r.unresolvedVerifier = ref
}

// SetUnresolvedRequired sets the unresolved requirement reference.
func (r *VerifyRelationship) SetUnresolvedRequired(ref string) {
	r.unresolvedRequired = ref
}

// ElementFilter represents an element filter package member.
type ElementFilter struct {
	baseElement
	Expression string
}

// NewElementFilter creates a new filter package member.
func NewElementFilter(name string, loc Location, expression string) *ElementFilter {
	return &ElementFilter{
		baseElement: baseElement{
			kind:     KindUnknown,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		Expression: expression,
	}
}

// PrefixMetadataAnnotation represents a prefix metadata annotation usage.
type PrefixMetadataAnnotation struct {
	baseElement
	Metadata Ref[*Metadata]
	Values   []string

	unresolvedMetadata string
}

// NewPrefixMetadataAnnotation creates a new prefix metadata annotation.
func NewPrefixMetadataAnnotation(loc Location) *PrefixMetadataAnnotation {
	return &PrefixMetadataAnnotation{
		baseElement: baseElement{
			kind:     KindMetadata,
			location: loc,
			children: make([]Element, 0),
		},
		Values: make([]string, 0),
	}
}

// SetUnresolvedMetadata sets the unresolved metadata reference.
func (a *PrefixMetadataAnnotation) SetUnresolvedMetadata(ref string) {
	a.unresolvedMetadata = ref
}
