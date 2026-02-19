package sysml

import "testing"

func TestResolveANTLRParserRuleName(t *testing.T) {
	got, ok := ResolveANTLRParserRuleName("Package")
	if !ok {
		t.Fatal("expected Package to resolve")
	}
	if got != "package" {
		t.Fatalf("expected package, got %q", got)
	}
}

func TestResolveANTLRTokenName(t *testing.T) {
	got, ok := ResolveANTLRTokenName("requirement")
	if !ok {
		t.Fatal("expected requirement to resolve")
	}
	if got != "REQUIREMENT" {
		t.Fatalf("expected REQUIREMENT, got %q", got)
	}
}
