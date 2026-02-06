package sysml

// ControlNodeKind represents the type of control node in activity diagrams.
type ControlNodeKind int

const (
	ControlNodeKindMerge ControlNodeKind = iota
	ControlNodeKindDecision
	ControlNodeKindJoin
	ControlNodeKindFork
)

// String returns the string representation of the control node kind.
func (k ControlNodeKind) String() string {
	switch k {
	case ControlNodeKindMerge:
		return "merge"
	case ControlNodeKindDecision:
		return "decision"
	case ControlNodeKindJoin:
		return "join"
	case ControlNodeKindFork:
		return "fork"
	default:
		return "unknown"
	}
}

// ControlNode represents a SysML control node for activity diagrams.
// Control nodes manage flow in activity diagrams:
//   - ForkNode: splits flow into parallel branches
//   - JoinNode: synchronizes parallel branches
//   - MergeNode: combines alternative flows
//   - DecisionNode: conditional branching (like if/else)
type ControlNode struct {
	baseElement
	NodeKind ControlNodeKind
	// Condition is the guard condition for decision nodes
	Condition string
	// IsControlNode marks this as a control node type
	IsControlNode bool
}

// NewControlNode creates a new ControlNode element.
func NewControlNode(kind ControlNodeKind, loc Location) *ControlNode {
	return &ControlNode{
		baseElement: baseElement{
			kind:     KindControlNode,
			location: loc,
			children: make([]Element, 0),
		},
		NodeKind:      kind,
		IsControlNode: true,
	}
}

// SetCondition sets the guard condition for decision nodes.
func (c *ControlNode) SetCondition(condition string) {
	c.Condition = condition
}

// GetNodeKind returns the control node kind.
func (c *ControlNode) GetNodeKind() ControlNodeKind {
	return c.NodeKind
}

// IsFork returns true if this is a fork node.
func (c *ControlNode) IsFork() bool {
	return c.NodeKind == ControlNodeKindFork
}

// IsJoin returns true if this is a join node.
func (c *ControlNode) IsJoin() bool {
	return c.NodeKind == ControlNodeKindJoin
}

// IsMerge returns true if this is a merge node.
func (c *ControlNode) IsMerge() bool {
	return c.NodeKind == ControlNodeKindMerge
}

// IsDecision returns true if this is a decision node.
func (c *ControlNode) IsDecision() bool {
	return c.NodeKind == ControlNodeKindDecision
}

// AddChild adds a child element to the control node.
// Control nodes typically don't have children in the traditional sense,
// but this method satisfies the Element interface.
func (c *ControlNode) AddChild(child Element) {
	c.baseElement.addChild(child)
}
