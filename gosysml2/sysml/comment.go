package sysml

// Comment represents a SysML comment annotation.
// Comments provide documentation and annotations for model elements.
type Comment struct {
	baseElement
	// Body is the text content of the comment
	Body string
	// Locale specifies the language/locale of the comment (e.g., "en", "fr")
	Locale string
	// About contains references to elements this comment annotates
	About []Element
	// unresolvedAbout holds string references during parsing before resolution
	unresolvedAbout []string
}

// NewComment creates a new Comment element.
func NewComment(body string, loc Location) *Comment {
	return &Comment{
		baseElement: baseElement{
			kind:     KindComment,
			location: loc,
			children: make([]Element, 0),
		},
		Body:            body,
		About:           make([]Element, 0),
		unresolvedAbout: make([]string, 0),
	}
}

// Kind returns the kind of this element.
func (c *Comment) Kind() ElementKind {
	return KindComment
}

// Name returns the name of this element.
// Comments are typically unnamed, so this returns an empty string.
func (c *Comment) Name() string {
	return ""
}

// GetKind returns the string representation of the element kind.
func (c *Comment) GetKind() string {
	return KindComment.String()
}

// GetLocation returns the source location of this element.
func (c *Comment) GetLocation() Location {
	return c.location
}

// GetParent returns the parent element, or nil for root elements.
func (c *Comment) GetParent() Element {
	return c.parent
}

// SetParent sets the parent element.
func (c *Comment) SetParent(parent Element) {
	c.parent = parent
}

// Accept implements the visitor pattern.
func (c *Comment) Accept(v Visitor) bool {
	return v.VisitComment(c)
}

// AddAbout adds an element reference to the about list.
func (c *Comment) AddAbout(elem Element) {
	c.About = append(c.About, elem)
}

// AddUnresolvedAbout adds an unresolved about reference.
// This is used during parsing before references are resolved.
func (c *Comment) AddUnresolvedAbout(ref string) {
	c.unresolvedAbout = append(c.unresolvedAbout, ref)
}

// UnresolvedAbout returns all unresolved about references for debugging.
func (c *Comment) UnresolvedAbout() []string {
	return c.unresolvedAbout
}

// Doc represents inline documentation in SysML.
// Documentation provides structured documentation for model elements.
type Doc struct {
	baseElement
	// Body is the text content of the documentation
	Body string
	// Locale specifies the language/locale of the documentation
	Locale string
}

// NewDoc creates a new Doc element.
func NewDoc(body string, loc Location) *Doc {
	return &Doc{
		baseElement: baseElement{
			kind:     KindDoc,
			location: loc,
			children: make([]Element, 0),
		},
		Body: body,
	}
}

// Kind returns the kind of this element.
func (d *Doc) Kind() ElementKind {
	return KindDoc
}

// Name returns the name of this element.
// Documentation is typically unnamed, so this returns an empty string.
func (d *Doc) Name() string {
	return ""
}

// GetKind returns the string representation of the element kind.
func (d *Doc) GetKind() string {
	return KindDoc.String()
}

// GetLocation returns the source location of this element.
func (d *Doc) GetLocation() Location {
	return d.location
}

// GetParent returns the parent element, or nil for root elements.
func (d *Doc) GetParent() Element {
	return d.parent
}

// SetParent sets the parent element.
func (d *Doc) SetParent(parent Element) {
	d.parent = parent
}

// Accept implements the visitor pattern.
func (d *Doc) Accept(v Visitor) bool {
	return v.VisitDoc(d)
}

// Comments returns all comments in the package.
func (p *Package) Comments() []*Comment {
	var comments []*Comment
	for _, child := range p.children {
		if c, ok := child.(*Comment); ok {
			comments = append(comments, c)
		}
	}
	return comments
}

// Docs returns all documentation elements in the package.
func (p *Package) Docs() []*Doc {
	var docs []*Doc
	for _, child := range p.children {
		if d, ok := child.(*Doc); ok {
			docs = append(docs, d)
		}
	}
	return docs
}
