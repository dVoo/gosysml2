package sysml

import (
	"testing"
)

func TestNewComment(t *testing.T) {
	loc := Location{Line: 10, Column: 5}
	body := "This is a test comment"
	comment := NewComment(body, loc)

	if comment == nil {
		t.Fatal("NewComment returned nil")
	}

	if comment.Body != body {
		t.Errorf("Expected body '%s', got '%s'", body, comment.Body)
	}

	if comment.Kind() != KindComment {
		t.Errorf("Expected kind %v, got %v", KindComment, comment.Kind())
	}

	if comment.Name() != "" {
		t.Errorf("Expected empty name, got '%s'", comment.Name())
	}

	if comment.GetLocation() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, comment.GetLocation())
	}

	if comment.GetParent() != nil {
		t.Error("Expected nil parent for new comment")
	}

	if comment.Locale != "" {
		t.Errorf("Expected empty locale, got '%s'", comment.Locale)
	}
}

func TestNewDoc(t *testing.T) {
	loc := Location{Line: 15, Column: 8}
	body := "This is documentation"
	doc := NewDoc(body, loc)

	if doc == nil {
		t.Fatal("NewDoc returned nil")
	}

	if doc.Body != body {
		t.Errorf("Expected body '%s', got '%s'", body, doc.Body)
	}

	if doc.Kind() != KindDoc {
		t.Errorf("Expected kind %v, got %v", KindDoc, doc.Kind())
	}

	if doc.Name() != "" {
		t.Errorf("Expected empty name, got '%s'", doc.Name())
	}

	if doc.GetLocation() != loc {
		t.Errorf("Expected location %+v, got %+v", loc, doc.GetLocation())
	}

	if doc.GetParent() != nil {
		t.Error("Expected nil parent for new doc")
	}

	if doc.Locale != "" {
		t.Errorf("Expected empty locale, got '%s'", doc.Locale)
	}
}

func TestCommentParent(t *testing.T) {
	comment := NewComment("Test", Location{})
	pkg := NewPackage("TestPkg", Location{})

	comment.SetParent(pkg)

	if comment.GetParent() != pkg {
		t.Error("SetParent did not set parent correctly for comment")
	}
}

func TestDocParent(t *testing.T) {
	doc := NewDoc("Test doc", Location{})
	pkg := NewPackage("TestPkg", Location{})

	doc.SetParent(pkg)

	if doc.GetParent() != pkg {
		t.Error("SetParent did not set parent correctly for doc")
	}
}

func TestCommentAbout(t *testing.T) {
	comment := NewComment("About this element", Location{})

	// Initially empty
	if len(comment.About) != 0 {
		t.Errorf("Expected 0 about elements, got %d", len(comment.About))
	}

	// Add elements
	part := NewPart("TestPart", Location{}, true)
	attr := NewAttribute("TestAttr", Location{}, true)

	comment.AddAbout(part)
	comment.AddAbout(attr)

	if len(comment.About) != 2 {
		t.Errorf("Expected 2 about elements, got %d", len(comment.About))
	}

	if comment.About[0] != part || comment.About[1] != attr {
		t.Error("About elements not added correctly")
	}
}

func TestCommentUnresolvedAbout(t *testing.T) {
	comment := NewComment("Test", Location{})

	// Initially empty
	if len(comment.UnresolvedAbout()) != 0 {
		t.Errorf("Expected 0 unresolved about, got %d", len(comment.UnresolvedAbout()))
	}

	// Add unresolved references
	comment.AddUnresolvedAbout("ElementA")
	comment.AddUnresolvedAbout("ElementB")

	unresolved := comment.UnresolvedAbout()

	if len(unresolved) != 2 {
		t.Errorf("Expected 2 unresolved about, got %d", len(unresolved))
	}

	if unresolved[0] != "ElementA" || unresolved[1] != "ElementB" {
		t.Errorf("Unexpected unresolved about: %v", unresolved)
	}
}

func TestCommentVisitor(t *testing.T) {
	comment := NewComment("Test", Location{})

	var visited bool
	visitor := &testCommentVisitor{
		onVisit: func(c *Comment) {
			visited = true
			if c != comment {
				t.Error("Visitor received wrong comment")
			}
		},
	}

	comment.Accept(visitor)

	if !visited {
		t.Error("Visitor was not called for comment")
	}
}

func TestDocVisitor(t *testing.T) {
	doc := NewDoc("Test doc", Location{})

	var visited bool
	visitor := &testDocVisitor{
		onVisit: func(d *Doc) {
			visited = true
			if d != doc {
				t.Error("Visitor received wrong doc")
			}
		},
	}

	doc.Accept(visitor)

	if !visited {
		t.Error("Visitor was not called for doc")
	}
}

type testCommentVisitor struct {
	BaseVisitor
	onVisit func(*Comment)
}

func (v *testCommentVisitor) VisitComment(c *Comment) bool {
	if v.onVisit != nil {
		v.onVisit(c)
	}
	return true
}

type testDocVisitor struct {
	BaseVisitor
	onVisit func(*Doc)
}

func (v *testDocVisitor) VisitDoc(d *Doc) bool {
	if v.onVisit != nil {
		v.onVisit(d)
	}
	return true
}

func TestModelAddComment(t *testing.T) {
	model := NewModel()
	comment := NewComment("Test comment", Location{Line: 5, Column: 10})

	model.AddComment(comment)

	if len(model.Comments) != 1 {
		t.Errorf("Expected 1 comment in model, got %d", len(model.Comments))
	}

	if len(model.Elements) < 1 {
		t.Errorf("Expected at least 1 element in model, got %d", len(model.Elements))
	}

	if model.Comments[0] != comment {
		t.Error("Comment not added correctly to model")
	}
}

func TestModelAddDoc(t *testing.T) {
	model := NewModel()
	doc := NewDoc("Test documentation", Location{Line: 5, Column: 10})

	model.AddDoc(doc)

	if len(model.Docs) != 1 {
		t.Errorf("Expected 1 doc in model, got %d", len(model.Docs))
	}

	if len(model.Elements) < 1 {
		t.Errorf("Expected at least 1 element in model, got %d", len(model.Elements))
	}

	if model.Docs[0] != doc {
		t.Error("Doc not added correctly to model")
	}
}

func TestPackageComments(t *testing.T) {
	pkg := NewPackage("TestPkg", Location{})
	comment1 := NewComment("Comment 1", Location{})
	comment2 := NewComment("Comment 2", Location{})

	pkg.AddChild(comment1)
	pkg.AddChild(comment2)

	comments := pkg.Comments()

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments in package, got %d", len(comments))
	}
}

func TestPackageDocs(t *testing.T) {
	pkg := NewPackage("TestPkg", Location{})
	doc1 := NewDoc("Doc 1", Location{})
	doc2 := NewDoc("Doc 2", Location{})

	pkg.AddChild(doc1)
	pkg.AddChild(doc2)

	docs := pkg.Docs()

	if len(docs) != 2 {
		t.Errorf("Expected 2 docs in package, got %d", len(docs))
	}
}

func TestCommentLocale(t *testing.T) {
	comment := NewComment("Test", Location{})

	// Initially empty
	if comment.Locale != "" {
		t.Errorf("Expected empty locale, got '%s'", comment.Locale)
	}

	// Set locale
	comment.Locale = "en-US"

	if comment.Locale != "en-US" {
		t.Errorf("Expected locale 'en-US', got '%s'", comment.Locale)
	}
}

func TestDocLocale(t *testing.T) {
	doc := NewDoc("Test", Location{})

	// Initially empty
	if doc.Locale != "" {
		t.Errorf("Expected empty locale, got '%s'", doc.Locale)
	}

	// Set locale
	doc.Locale = "fr-FR"

	if doc.Locale != "fr-FR" {
		t.Errorf("Expected locale 'fr-FR', got '%s'", doc.Locale)
	}
}
