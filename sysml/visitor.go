package sysml

import "iter"

// Visitor defines the interface for visiting SysML elements.
// Implement this interface to traverse and process a model.
type Visitor interface {
	// VisitPackage is called for each package element.
	// Return false to skip visiting children.
	VisitPackage(pkg *Package) bool

	// VisitPart is called for each part element.
	VisitPart(part *Part) bool

	// VisitRequirement is called for each requirement element.
	VisitRequirement(req *Requirement) bool

	// VisitVerification is called for each verification case element.
	VisitVerification(ver *Verification) bool

	// VisitConcern is called for each concern element.
	VisitConcern(concern *Concern) bool

	// VisitUseCase is called for each use case element.
	VisitUseCase(useCase *UseCase) bool

	// VisitAnalysisCase is called for each analysis case element.
	VisitAnalysisCase(analysis *AnalysisCase) bool

	// VisitCase is called for each case element.
	VisitCase(case_ *Case) bool

	// VisitIncludeUseCase is called for each include use case element.
	VisitIncludeUseCase(include *IncludeUseCase) bool

	// VisitAction is called for each action element.
	VisitAction(action *Action) bool

	// VisitImport is called for each import element.
	VisitImport(imp *Import) bool

	// VisitComment is called for each comment element.
	VisitComment(comment *Comment) bool

	// VisitDoc is called for each documentation element.
	VisitDoc(doc *Doc) bool

	// VisitPort is called for each port element.
	VisitPort(port *Port) bool

	// VisitConjugatedPort is called for each conjugated port element.
	VisitConjugatedPort(conj *ConjugatedPort) bool

	// VisitConnection is called for each connection element.
	VisitConnection(conn *Connection) bool

	// VisitInterface is called for each interface element.
	VisitInterface(iface *Interface) bool

	// VisitAllocation is called for each allocation element.
	VisitAllocation(alloc *Allocation) bool

	// VisitAttribute is called for each attribute element.
	VisitAttribute(attr *Attribute) bool

	// VisitItem is called for each item element.
	VisitItem(item *Item) bool

	// VisitState is called for each state element.
	VisitState(state *State) bool

	// VisitTransition is called for each transition element.
	VisitTransition(transition *Transition) bool

	// VisitCalculation is called for each calculation element.
	VisitCalculation(calc *Calculation) bool

	// VisitConstraint is called for each constraint element.
	VisitConstraint(constraint *Constraint) bool

	// VisitEnumeration is called for each enumeration element.
	VisitEnumeration(enum *Enumeration) bool

	// VisitEnumerationValue is called for each enumeration value element.
	VisitEnumerationValue(value *EnumerationValue) bool

	// VisitView is called for each view element.
	VisitView(view *View) bool

	// VisitViewpoint is called for each viewpoint element.
	VisitViewpoint(viewpoint *Viewpoint) bool

	// VisitDependency is called for each dependency element.
	VisitDependency(dep *Dependency) bool

	// VisitFlow is called for each flow element.
	VisitFlow(flow *Flow) bool

	// VisitFlowEnd is called for each flow end element.
	VisitFlowEnd(flowEnd *FlowEnd) bool

	// VisitSuccessionFlow is called for each succession flow element.
	VisitSuccessionFlow(successionFlow *SuccessionFlow) bool

	// VisitElement is called for any other element type.
	VisitElement(elem Element) bool
}

// BaseVisitor provides a default implementation of Visitor.
// Embed this in your visitor to only override methods you care about.
type BaseVisitor struct{}

func (BaseVisitor) VisitPackage(pkg *Package) bool                          { return true }
func (BaseVisitor) VisitPart(part *Part) bool                               { return true }
func (BaseVisitor) VisitRequirement(req *Requirement) bool                  { return true }
func (BaseVisitor) VisitVerification(ver *Verification) bool                { return true }
func (BaseVisitor) VisitConcern(concern *Concern) bool                      { return true }
func (BaseVisitor) VisitUseCase(useCase *UseCase) bool                      { return true }
func (BaseVisitor) VisitAnalysisCase(analysis *AnalysisCase) bool           { return true }
func (BaseVisitor) VisitCase(case_ *Case) bool                              { return true }
func (BaseVisitor) VisitIncludeUseCase(include *IncludeUseCase) bool        { return true }
func (BaseVisitor) VisitAction(action *Action) bool                         { return true }
func (BaseVisitor) VisitImport(imp *Import) bool                            { return true }
func (BaseVisitor) VisitComment(comment *Comment) bool                      { return true }
func (BaseVisitor) VisitDoc(doc *Doc) bool                                  { return true }
func (BaseVisitor) VisitPort(port *Port) bool                               { return true }
func (BaseVisitor) VisitConjugatedPort(conj *ConjugatedPort) bool           { return true }
func (BaseVisitor) VisitConnection(conn *Connection) bool                   { return true }
func (BaseVisitor) VisitInterface(iface *Interface) bool                    { return true }
func (BaseVisitor) VisitAllocation(alloc *Allocation) bool                  { return true }
func (BaseVisitor) VisitAttribute(attr *Attribute) bool                     { return true }
func (BaseVisitor) VisitItem(item *Item) bool                               { return true }
func (BaseVisitor) VisitState(state *State) bool                            { return true }
func (BaseVisitor) VisitTransition(transition *Transition) bool             { return true }
func (BaseVisitor) VisitCalculation(calc *Calculation) bool                 { return true }
func (BaseVisitor) VisitConstraint(constraint *Constraint) bool             { return true }
func (BaseVisitor) VisitEnumeration(enum *Enumeration) bool                 { return true }
func (BaseVisitor) VisitEnumerationValue(value *EnumerationValue) bool      { return true }
func (BaseVisitor) VisitView(view *View) bool                               { return true }
func (BaseVisitor) VisitViewpoint(viewpoint *Viewpoint) bool                { return true }
func (BaseVisitor) VisitDependency(dep *Dependency) bool                    { return true }
func (BaseVisitor) VisitFlow(flow *Flow) bool                               { return true }
func (BaseVisitor) VisitFlowEnd(flowEnd *FlowEnd) bool                      { return true }
func (BaseVisitor) VisitSuccessionFlow(successionFlow *SuccessionFlow) bool { return true }
func (BaseVisitor) VisitElement(elem Element) bool                          { return true }

// Visit traverses a model using the given visitor.
func Visit(model *Model, visitor Visitor) {
	for _, elem := range model.Elements {
		visitElement(elem, visitor)
	}
}

func visitElement(elem Element, visitor Visitor) {
	var continueVisit bool

	switch e := elem.(type) {
	case *Package:
		continueVisit = visitor.VisitPackage(e)
	case *Part:
		continueVisit = visitor.VisitPart(e)
	case *Requirement:
		continueVisit = visitor.VisitRequirement(e)
	case *Verification:
		continueVisit = visitor.VisitVerification(e)
	case *Concern:
		continueVisit = visitor.VisitConcern(e)
	case *UseCase:
		continueVisit = visitor.VisitUseCase(e)
	case *AnalysisCase:
		continueVisit = visitor.VisitAnalysisCase(e)
	case *Case:
		continueVisit = visitor.VisitCase(e)
	case *IncludeUseCase:
		continueVisit = visitor.VisitIncludeUseCase(e)
	case *Action:
		continueVisit = visitor.VisitAction(e)
	case *Import:
		continueVisit = visitor.VisitImport(e)
	case *Comment:
		continueVisit = visitor.VisitComment(e)
	case *Doc:
		continueVisit = visitor.VisitDoc(e)
	case *Port:
		continueVisit = visitor.VisitPort(e)
	case *ConjugatedPort:
		continueVisit = visitor.VisitConjugatedPort(e)
	case *Connection:
		continueVisit = visitor.VisitConnection(e)
	case *Interface:
		continueVisit = visitor.VisitInterface(e)
	case *Allocation:
		continueVisit = visitor.VisitAllocation(e)
	case *Attribute:
		continueVisit = visitor.VisitAttribute(e)
	case *Item:
		continueVisit = visitor.VisitItem(e)
	case *State:
		continueVisit = visitor.VisitState(e)
	case *Transition:
		continueVisit = visitor.VisitTransition(e)
	case *Calculation:
		continueVisit = visitor.VisitCalculation(e)
	case *Constraint:
		continueVisit = visitor.VisitConstraint(e)
	case *Enumeration:
		continueVisit = visitor.VisitEnumeration(e)
	case *EnumerationValue:
		continueVisit = visitor.VisitEnumerationValue(e)
	case *View:
		continueVisit = visitor.VisitView(e)
	case *Viewpoint:
		continueVisit = visitor.VisitViewpoint(e)
	case *Dependency:
		continueVisit = visitor.VisitDependency(e)
	case *Flow:
		continueVisit = visitor.VisitFlow(e)
	case *FlowEnd:
		continueVisit = visitor.VisitFlowEnd(e)
	case *SuccessionFlow:
		continueVisit = visitor.VisitSuccessionFlow(e)
	default:
		continueVisit = visitor.VisitElement(elem)
	}

	if continueVisit {
		for _, child := range elem.Children() {
			visitElement(child, visitor)
		}
	}
}

// WalkFunc is a function that is called for each element during a walk.
// Return false to stop the walk.
type WalkFunc func(elem Element, depth int) bool

// Walk traverses a model depth-first, calling fn for each element.
func Walk(model *Model, fn WalkFunc) {
	for _, elem := range model.Elements {
		if !walkDepth(elem, fn, 0) {
			return
		}
	}
}

func walkDepth(elem Element, fn WalkFunc, depth int) bool {
	if !fn(elem, depth) {
		return false
	}
	for _, child := range elem.Children() {
		if !walkDepth(child, fn, depth+1) {
			return false
		}
	}
	return true
}

// FindAll returns all elements of type T in the model.
func FindAll[T Element](model *Model) []T {
	var result []T
	Walk(model, func(elem Element, depth int) bool {
		if t, ok := elem.(T); ok {
			result = append(result, t)
		}
		return true
	})
	return result
}

// All returns an iterator over all elements in the model (depth-first).
func All(model *Model) iter.Seq[Element] {
	return func(yield func(Element) bool) {
		for _, elem := range model.Elements {
			if !yieldAll(elem, yield) {
				return
			}
		}
	}
}

// yieldAll yields elem and all its descendants.
func yieldAll(elem Element, yield func(Element) bool) bool {
	if !yield(elem) {
		return false
	}
	for _, child := range elem.Children() {
		if !yieldAll(child, yield) {
			return false
		}
	}
	return true
}

// AllWithDepth returns an iterator with element and depth.
func AllWithDepth(model *Model) iter.Seq2[Element, int] {
	return func(yield func(Element, int) bool) {
		for _, elem := range model.Elements {
			if !yieldAllDepth(elem, 0, yield) {
				return
			}
		}
	}
}

func yieldAllDepth(elem Element, depth int, yield func(Element, int) bool) bool {
	if !yield(elem, depth) {
		return false
	}
	for _, child := range elem.Children() {
		if !yieldAllDepth(child, depth+1, yield) {
			return false
		}
	}
	return true
}

// OfKind returns an iterator over elements of a specific kind.
func OfKind(model *Model, kind ElementKind) iter.Seq[Element] {
	return func(yield func(Element) bool) {
		for elem := range All(model) {
			if elem.Kind() == kind {
				if !yield(elem) {
					return
				}
			}
		}
	}
}

// OfType returns an iterator over elements of a specific Go type.
func OfType[T Element](model *Model) iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range All(model) {
			if t, ok := elem.(T); ok {
				if !yield(t) {
					return
				}
			}
		}
	}
}

// CountAll counts all elements by kind using iterators.
func CountAll(model *Model) map[ElementKind]int {
	counts := make(map[ElementKind]int)
	for elem := range All(model) {
		counts[elem.Kind()]++
	}
	return counts
}

// Filter returns elements matching the given predicate.
func Filter(model *Model, predicate func(Element) bool) []Element {
	var result []Element
	Walk(model, func(elem Element, depth int) bool {
		if predicate(elem) {
			result = append(result, elem)
		}
		return true
	})
	return result
}

// FindByKind returns all elements of a specific kind.
func FindByKind(model *Model, kind ElementKind) []Element {
	return Filter(model, func(elem Element) bool {
		return elem.Kind() == kind
	})
}

// FindByName returns all elements with a specific name.
func FindByName(model *Model, name string) []Element {
	return Filter(model, func(elem Element) bool {
		return elem.Name() == name
	})
}

// FindPackages returns all packages in the model.
// Deprecated: Use FindAll[*Package](model) instead.
func FindPackages(model *Model) []*Package {
	return FindAll[*Package](model)
}

// FindParts returns all parts in the model.
// Deprecated: Use FindAll[*Part](model) instead.
func FindParts(model *Model) []*Part {
	return FindAll[*Part](model)
}

// FindRequirements returns all requirements in the model.
// Deprecated: Use FindAll[*Requirement](model) instead.
func FindRequirements(model *Model) []*Requirement {
	return FindAll[*Requirement](model)
}

// FindVerifications returns all verification cases in the model.
func FindVerifications(model *Model) []*Verification {
	var result []*Verification
	Walk(model, func(elem Element, depth int) bool {
		if ver, ok := elem.(*Verification); ok {
			result = append(result, ver)
		}
		return true
	})
	return result
}

// FindConcerns returns all concerns in the model.
func FindConcerns(model *Model) []*Concern {
	var result []*Concern
	Walk(model, func(elem Element, depth int) bool {
		if concern, ok := elem.(*Concern); ok {
			result = append(result, concern)
		}
		return true
	})
	return result
}

// FindUseCases returns all use cases in the model.
func FindUseCases(model *Model) []*UseCase {
	var result []*UseCase
	Walk(model, func(elem Element, depth int) bool {
		if useCase, ok := elem.(*UseCase); ok {
			result = append(result, useCase)
		}
		return true
	})
	return result
}

// FindAnalysisCases returns all analysis cases in the model.
func FindAnalysisCases(model *Model) []*AnalysisCase {
	var result []*AnalysisCase
	Walk(model, func(elem Element, depth int) bool {
		if analysis, ok := elem.(*AnalysisCase); ok {
			result = append(result, analysis)
		}
		return true
	})
	return result
}

// FindCases returns all cases in the model.
func FindCases(model *Model) []*Case {
	var result []*Case
	Walk(model, func(elem Element, depth int) bool {
		if c, ok := elem.(*Case); ok {
			result = append(result, c)
		}
		return true
	})
	return result
}

// FindIncludeUseCases returns all include use cases in the model.
func FindIncludeUseCases(model *Model) []*IncludeUseCase {
	var result []*IncludeUseCase
	Walk(model, func(elem Element, depth int) bool {
		if inc, ok := elem.(*IncludeUseCase); ok {
			result = append(result, inc)
		}
		return true
	})
	return result
}

// FindActions returns all actions in the model.
func FindActions(model *Model) []*Action {
	var result []*Action
	Walk(model, func(elem Element, depth int) bool {
		if action, ok := elem.(*Action); ok {
			result = append(result, action)
		}
		return true
	})
	return result
}

// FindAttributes returns all attributes in the model.
func FindAttributes(model *Model) []*Attribute {
	var result []*Attribute
	Walk(model, func(elem Element, depth int) bool {
		if attr, ok := elem.(*Attribute); ok {
			result = append(result, attr)
		}
		return true
	})
	return result
}

// FindPorts returns all ports in the model.
func FindPorts(model *Model) []*Port {
	var result []*Port
	Walk(model, func(elem Element, depth int) bool {
		if port, ok := elem.(*Port); ok {
			result = append(result, port)
		}
		return true
	})
	return result
}

// FindConjugatedPorts returns all conjugated ports in the model.
func FindConjugatedPorts(model *Model) []*ConjugatedPort {
	var result []*ConjugatedPort
	Walk(model, func(elem Element, depth int) bool {
		if conj, ok := elem.(*ConjugatedPort); ok {
			result = append(result, conj)
		}
		return true
	})
	return result
}

// FindConnections returns all connections in the model.
func FindConnections(model *Model) []*Connection {
	var result []*Connection
	Walk(model, func(elem Element, depth int) bool {
		if conn, ok := elem.(*Connection); ok {
			result = append(result, conn)
		}
		return true
	})
	return result
}

// FindInterfaces returns all interfaces in the model.
func FindInterfaces(model *Model) []*Interface {
	var result []*Interface
	Walk(model, func(elem Element, depth int) bool {
		if iface, ok := elem.(*Interface); ok {
			result = append(result, iface)
		}
		return true
	})
	return result
}

// FindAllocations returns all allocations in the model.
func FindAllocations(model *Model) []*Allocation {
	var result []*Allocation
	Walk(model, func(elem Element, depth int) bool {
		if alloc, ok := elem.(*Allocation); ok {
			result = append(result, alloc)
		}
		return true
	})
	return result
}

// FindFlows returns all flows in the model.
func FindFlows(model *Model) []*Flow {
	var result []*Flow
	Walk(model, func(elem Element, depth int) bool {
		if flow, ok := elem.(*Flow); ok {
			result = append(result, flow)
		}
		return true
	})
	return result
}

// FindSuccessionFlows returns all succession flows in the model.
func FindSuccessionFlows(model *Model) []*SuccessionFlow {
	var result []*SuccessionFlow
	Walk(model, func(elem Element, depth int) bool {
		if sf, ok := elem.(*SuccessionFlow); ok {
			result = append(result, sf)
		}
		return true
	})
	return result
}

// FindDefinitions returns all definition elements in the model.
func FindDefinitions(model *Model) []Definition {
	var result []Definition
	Walk(model, func(elem Element, depth int) bool {
		if def, ok := elem.(Definition); ok {
			result = append(result, def)
		}
		return true
	})
	return result
}

// FindUsages returns all usage elements in the model.
func FindUsages(model *Model) []Usage {
	var result []Usage
	Walk(model, func(elem Element, depth int) bool {
		if usage, ok := elem.(Usage); ok {
			result = append(result, usage)
		}
		return true
	})
	return result
}

// FindByQualifiedName finds an element by its qualified name.
func FindByQualifiedName(model *Model, qn string) Element {
	return model.FindByQualifiedName(qn)
}

// Counter is a visitor that counts elements by kind.
type Counter struct {
	BaseVisitor
	Counts map[ElementKind]int
}

// NewCounter creates a new element counter.
func NewCounter() *Counter {
	return &Counter{
		Counts: make(map[ElementKind]int),
	}
}

func (c *Counter) VisitPackage(pkg *Package) bool {
	c.Counts[KindPackage]++
	return true
}

func (c *Counter) VisitPart(part *Part) bool {
	c.Counts[KindPart]++
	return true
}

func (c *Counter) VisitRequirement(req *Requirement) bool {
	c.Counts[KindRequirement]++
	return true
}

func (c *Counter) VisitVerification(ver *Verification) bool {
	c.Counts[KindVerification]++
	return true
}

func (c *Counter) VisitConcern(concern *Concern) bool {
	c.Counts[KindConcern]++
	return true
}

func (c *Counter) VisitUseCase(useCase *UseCase) bool {
	c.Counts[KindUseCase]++
	return true
}

func (c *Counter) VisitAnalysisCase(analysis *AnalysisCase) bool {
	c.Counts[KindAnalysis]++
	return true
}

func (c *Counter) VisitCase(case_ *Case) bool {
	c.Counts[KindCase]++
	return true
}

func (c *Counter) VisitIncludeUseCase(include *IncludeUseCase) bool {
	c.Counts[KindIncludeUseCase]++
	return true
}

func (c *Counter) VisitConjugatedPort(conj *ConjugatedPort) bool {
	c.Counts[KindConjugatedPort]++
	return true
}

func (c *Counter) VisitAction(action *Action) bool {
	c.Counts[KindAction]++
	return true
}

func (c *Counter) VisitImport(imp *Import) bool {
	c.Counts[KindImport]++
	return true
}

func (c *Counter) VisitComment(comment *Comment) bool {
	c.Counts[KindComment]++
	return true
}

func (c *Counter) VisitElement(elem Element) bool {
	c.Counts[elem.Kind()]++
	return true
}

func (c *Counter) VisitSuccessionFlow(successionFlow *SuccessionFlow) bool {
	c.Counts[KindSuccessionFlow]++
	return true
}

// Total returns the total count of all elements.
func (c *Counter) Total() int {
	total := 0
	for _, count := range c.Counts {
		total += count
	}
	return total
}
