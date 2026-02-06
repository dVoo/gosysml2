package sysml

// PortionKind represents the type of portion for occurrences.
type PortionKind int

const (
	PortionKindNone PortionKind = iota
	PortionKindSnapshot
	PortionKindTimeslice
)

// String returns the string representation of the portion kind.
func (k PortionKind) String() string {
	switch k {
	case PortionKindSnapshot:
		return "snapshot"
	case PortionKindTimeslice:
		return "timeslice"
	default:
		return ""
	}
}

// LifeStep represents the life step of an occurrence (start or end).
type LifeStep int

const (
	LifeStepNone LifeStep = iota
	LifeStepStart
	LifeStepEnd
)

// String returns the string representation of the life step.
func (s LifeStep) String() string {
	switch s {
	case LifeStepStart:
		return "start"
	case LifeStepEnd:
		return "end"
	default:
		return ""
	}
}

// Occurrence represents a SysML occurrence for time-based modeling.
// Occurrences model time-based instances and can be:
//   - OccurrenceDefinition: defines occurrence types
//   - OccurrenceUsage: uses occurrence types
//   - IndividualDefinition: specific individual instances
//   - EventOccurrenceUsage: event-based occurrences
//   - TimeSlice: portion of an occurrence over time
//   - Snapshot: instantaneous portion of an occurrence
//   - LifeStep: start or end of an occurrence's life
//   - Portion: general portion with kind (snapshot/timeslice)
type Occurrence struct {
	baseElement
	IsDefinition bool
	IsIndividual bool
	IsEvent      bool // EventOccurrenceUsage
	IsTimeSlice  bool // TimeSlice
	IsSnapshot   bool // Snapshot
	PortionKind_ PortionKind
	LifeStep_    LifeStep

	// Typed children
	occurrences []*Occurrence
}

// NewOccurrence creates a new Occurrence element.
func NewOccurrence(name string, loc Location, isDef, isInd bool) *Occurrence {
	return &Occurrence{
		baseElement: baseElement{
			kind:     KindOccurrence,
			name:     name,
			location: loc,
			children: make([]Element, 0),
		},
		IsDefinition: isDef,
		IsIndividual: isInd,
		occurrences:  make([]*Occurrence, 0),
	}
}

// NewEventOccurrence creates a new EventOccurrenceUsage element.
func NewEventOccurrence(name string, loc Location) *Occurrence {
	occ := NewOccurrence(name, loc, false, false)
	occ.IsEvent = true
	return occ
}

// NewTimeSlice creates a new TimeSlice element.
func NewTimeSlice(name string, loc Location) *Occurrence {
	occ := NewOccurrence(name, loc, false, false)
	occ.IsTimeSlice = true
	return occ
}

// NewSnapshot creates a new Snapshot element.
func NewSnapshot(name string, loc Location) *Occurrence {
	occ := NewOccurrence(name, loc, false, false)
	occ.IsSnapshot = true
	occ.PortionKind_ = PortionKindSnapshot
	return occ
}

// NewLifeStep creates a new LifeStep element.
func NewLifeStep(name string, loc Location, step LifeStep) *Occurrence {
	occ := NewOccurrence(name, loc, false, false)
	occ.LifeStep_ = step
	return occ
}

// SetPortionKind sets the portion kind for this occurrence.
func (o *Occurrence) SetPortionKind(kind PortionKind) {
	o.PortionKind_ = kind
}

// PortionKind returns the portion kind of this occurrence.
func (o *Occurrence) PortionKind() PortionKind {
	return o.PortionKind_
}

// SetLifeStep sets the life step for this occurrence.
func (o *Occurrence) SetLifeStep(step LifeStep) {
	o.LifeStep_ = step
}

// LifeStep returns the life step of this occurrence.
func (o *Occurrence) LifeStep() LifeStep {
	return o.LifeStep_
}

// AddChild adds a child element with type tracking.
func (o *Occurrence) AddChild(child Element) {
	o.baseElement.addChild(child)

	if occ, ok := child.(*Occurrence); ok {
		o.occurrences = append(o.occurrences, occ)
	}
}

// Occurrences returns nested occurrences.
func (o *Occurrence) Occurrences() []*Occurrence {
	return o.occurrences
}

// IsFork returns false (occurrences are not control nodes).
func (o *Occurrence) IsFork() bool {
	return false
}

// IsJoin returns false (occurrences are not control nodes).
func (o *Occurrence) IsJoin() bool {
	return false
}

// IsMerge returns false (occurrences are not control nodes).
func (o *Occurrence) IsMerge() bool {
	return false
}

// IsDecision returns false (occurrences are not control nodes).
func (o *Occurrence) IsDecision() bool {
	return false
}
