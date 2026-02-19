package sysml

import "testing"

func TestPhase1RepresentationTypes(t *testing.T) {
	loc := Location{Line: 1, Column: 1}

	alias := NewAlias("AliasA", loc)
	if alias.Kind() != KindAlias {
		t.Fatalf("expected alias kind %v, got %v", KindAlias, alias.Kind())
	}
	alias.SetUnresolvedTarget("Pkg::PartA")

	metadata := NewMetadata("MetaA", loc, true)
	if metadata.Kind() != KindMetadata {
		t.Fatalf("expected metadata kind %v, got %v", KindMetadata, metadata.Kind())
	}
	annotation := NewPrefixMetadataAnnotation(loc)
	metadata.AddChild(annotation)
	if len(metadata.Annotations()) != 1 {
		t.Fatalf("expected 1 metadata annotation, got %d", len(metadata.Annotations()))
	}

	rendering := NewRendering("RenderA", loc, false)
	if rendering.Kind() != KindRendering {
		t.Fatalf("expected rendering kind %v, got %v", KindRendering, rendering.Kind())
	}
	if rendering.Type() != nil {
		t.Fatal("expected unresolved rendering type to be nil")
	}

	message := NewMessage("MessageA", loc)
	if message.Kind() != KindMessage {
		t.Fatalf("expected message kind %v, got %v", KindMessage, message.Kind())
	}
	message.SetUnresolvedSender("sender")
	message.SetUnresolvedReceiver("receiver")
	if message.Type() != nil {
		t.Fatal("expected message type to be nil")
	}

	satisfy := NewSatisfyRelationship("sat", loc)
	satisfy.SetUnresolvedSatisfier("engine")
	satisfy.SetUnresolvedRequired("REQ_1")
	if satisfy.Kind() != KindDependency {
		t.Fatalf("expected satisfy relationship kind %v, got %v", KindDependency, satisfy.Kind())
	}

	verify := NewVerifyRelationship("ver", loc)
	verify.SetUnresolvedVerifier("verify_case")
	verify.SetUnresolvedRequired("REQ_1")
	if verify.Kind() != KindDependency {
		t.Fatalf("expected verify relationship kind %v, got %v", KindDependency, verify.Kind())
	}

	filter := NewElementFilter("OnlyReqs", loc, "kind == requirement")
	if filter.Expression == "" {
		t.Fatal("expected filter expression to be set")
	}
}
