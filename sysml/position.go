package sysml

// ElementAt returns the most specific named element containing a zero-based position.
// Returns nil when no element spans the given position.
func ElementAt(model *Model, line, column int) Element {
	return elementAt(model, line, column, false)
}

// ElementAtIncludingUnnamed returns the most specific element containing a zero-based position,
// including unnamed elements such as anonymous usages.
func ElementAtIncludingUnnamed(model *Model, line, column int) Element {
	return elementAt(model, line, column, true)
}

func elementAt(model *Model, line, column int, includeUnnamed bool) Element {
	if model == nil {
		return nil
	}

	var best Element
	bestDepth := -1
	bestSpan := int(^uint(0) >> 1)

	for elem, depth := range AllWithDepth(model) {
		if elem == nil {
			continue
		}
		if !includeUnnamed && elem.Name() == "" {
			continue
		}
		loc := elem.Location()
		if !loc.Contains(line, column) {
			continue
		}
		span := locationSpanScore(loc)
		if best == nil || depth > bestDepth || (depth == bestDepth && span < bestSpan) {
			best = elem
			bestDepth = depth
			bestSpan = span
		}
	}

	return best
}

func locationSpanScore(loc Location) int {
	endLine := loc.EndLine
	endColumn := loc.EndColumn
	if endLine < 0 || endColumn < 0 {
		endLine = loc.Line
		endColumn = loc.Column
	}
	lineSpan := endLine - loc.Line
	colSpan := endColumn - loc.Column
	if lineSpan < 0 || colSpan < 0 {
		return 0
	}
	return lineSpan*100000 + colSpan
}
