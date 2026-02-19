package sysml

import "testing"

func TestPhase234ResolveDependencyAndCommentRefs(t *testing.T) {
	loc := Location{Line: 1, Column: 1}
	m := NewModel()
	pkg := NewPackage("P", loc)
	m.AddPackage(pkg)

	a := NewPart("A", loc, true)
	b := NewPart("B", loc, true)
	pkg.AddChild(a)
	pkg.AddChild(b)

	dep := NewDependency(loc)
	dep.AddUnresolvedClient("A")
	dep.AddUnresolvedSupplier("B")
	pkg.AddChild(dep)

	comment := NewComment("c", loc)
	comment.AddUnresolvedAbout("A")
	pkg.AddChild(comment)

	m.BuildIndex()
	m.ResolveReferences()

	if len(dep.Client) != 1 || dep.Client[0] != a {
		t.Fatalf("dependency client not resolved, got %#v", dep.Client)
	}
	if len(dep.Supplier) != 1 || dep.Supplier[0] != b {
		t.Fatalf("dependency supplier not resolved, got %#v", dep.Supplier)
	}
	if len(comment.About) != 1 || comment.About[0] != a {
		t.Fatalf("comment about not resolved, got %#v", comment.About)
	}
}

func TestPhase234ResolveNewRepresentationRefs(t *testing.T) {
	loc := Location{Line: 1, Column: 1}
	m := NewModel()
	pkg := NewPackage("P", loc)
	m.AddPackage(pkg)

	sender := NewPart("Sender", loc, true)
	receiver := NewPart("Receiver", loc, true)
	req := NewRequirement("REQ_1", loc, true)
	ver := NewVerification("V_1", loc, true)
	metaDef := NewMetadata("MetaDef", loc, true)
	metaUse := NewMetadata("MetaUse", loc, false)
	renderDef := NewRendering("RenderDef", loc, true)
	renderUse := NewRendering("RenderUse", loc, false)
	alias := NewAlias("AliasReq", loc)
	message := NewMessage("msg", loc)
	satisfy := NewSatisfyRelationship("sat", loc)
	verify := NewVerifyRelationship("ver", loc)

	metaUse.TypeRef = NewRef[*Metadata]("MetaDef")
	renderUse.TypeRef = NewRef[*Rendering]("RenderDef")
	alias.SetUnresolvedTarget("REQ_1")
	message.SetUnresolvedSender("Sender")
	message.SetUnresolvedReceiver("Receiver")
	satisfy.SetUnresolvedSatisfier("Sender")
	satisfy.SetUnresolvedRequired("REQ_1")
	verify.SetUnresolvedVerifier("V_1")
	verify.SetUnresolvedRequired("REQ_1")

	pkg.AddChild(sender)
	pkg.AddChild(receiver)
	pkg.AddChild(req)
	pkg.AddChild(ver)
	pkg.AddChild(metaDef)
	pkg.AddChild(metaUse)
	pkg.AddChild(renderDef)
	pkg.AddChild(renderUse)
	pkg.AddChild(alias)
	pkg.AddChild(message)
	pkg.AddChild(satisfy)
	pkg.AddChild(verify)

	m.BuildIndex()
	m.ResolveReferences()

	if !alias.Target.IsResolved() || alias.Target.Resolved() != req {
		t.Fatal("alias target was not resolved")
	}
	if !metaUse.TypeRef.IsResolved() || metaUse.TypeRef.Resolved() != metaDef {
		t.Fatal("metadata type reference was not resolved")
	}
	if !renderUse.TypeRef.IsResolved() || renderUse.TypeRef.Resolved() != renderDef {
		t.Fatal("rendering type reference was not resolved")
	}
	if !message.Sender.IsResolved() || message.Sender.Resolved() != sender {
		t.Fatal("message sender was not resolved")
	}
	if !message.Receiver.IsResolved() || message.Receiver.Resolved() != receiver {
		t.Fatal("message receiver was not resolved")
	}
	if !satisfy.Required.IsResolved() || satisfy.Required.Resolved() != req {
		t.Fatal("satisfy relationship requirement was not resolved")
	}
	if !verify.Required.IsResolved() || verify.Required.Resolved() != req {
		t.Fatal("verify relationship requirement was not resolved")
	}
	if len(req.SatisfiedBy) == 0 || req.SatisfiedBy[0] != sender {
		t.Fatal("requirement satisfied-by inverse relation missing")
	}
	if len(req.VerifiedBy) == 0 || req.VerifiedBy[0] != ver {
		t.Fatal("requirement verified-by inverse relation missing")
	}
}

func TestPhase234ModelTracksNestedSatisfyAndVerify(t *testing.T) {
	satModel := mustParsePhase5File(t, "../validationdata/08-Requirements/8-Requirements.sysml")
	if len(satModel.Satisfies) == 0 {
		t.Fatal("expected model.Satisfies to include package-nested satisfy relationships")
	}

	verModel := mustParsePhase5File(t, "../validationdata/09-Verification/9-Verification-simplified.sysml")
	if len(verModel.Verifies) == 0 {
		t.Fatal("expected model.Verifies to include package-nested verify relationships")
	}
}

func TestPhase234ResolveByShortNameAndDottedChain(t *testing.T) {
	loc := Location{Line: 1, Column: 1}
	m := NewModel()
	pkg := NewPackage("P", loc)
	m.AddPackage(pkg)

	req := NewRequirement("vehicleMassReq", loc, true)
	req.setDeclaredShortName("'R1'")
	pkg.AddChild(req)

	vehicle := NewPart("vehicle1_c1", loc, true)
	engine := NewPart("engine_v1", loc, true)
	vehicle.AddChild(engine)
	pkg.AddChild(vehicle)

	satisfy := NewSatisfyRelationship("sat", loc)
	satisfy.SetUnresolvedRequired("'R1'")
	satisfy.SetUnresolvedSatisfier("vehicle1_c1.engine_v1")
	pkg.AddChild(satisfy)

	m.BuildIndex()
	m.ResolveReferences()

	if !satisfy.Required.IsResolved() || satisfy.Required.Resolved() != req {
		t.Fatal("short-name requirement reference was not resolved")
	}
	if !satisfy.Satisfier.IsResolved() || satisfy.Satisfier.Resolved() != engine {
		t.Fatal("dotted-chain satisfier reference was not resolved")
	}
}

func TestPhase234RelationshipRefNamesArePreserved(t *testing.T) {
	loc := Location{Line: 1, Column: 1}

	sat := NewSatisfyRelationship("sat", loc)
	sat.SetUnresolvedSatisfier("myCar")
	sat.SetUnresolvedRequired("vehicleMassReq")
	if sat.Satisfier.Name() != "myCar" {
		t.Fatalf("expected satisfier ref name 'myCar', got %q", sat.Satisfier.Name())
	}
	if sat.Required.Name() != "vehicleMassReq" {
		t.Fatalf("expected required ref name 'vehicleMassReq', got %q", sat.Required.Name())
	}

	ver := NewVerifyRelationship("ver", loc)
	ver.SetUnresolvedVerifier("V1")
	ver.SetUnresolvedRequired("REQ_1")
	if ver.Verifier.Name() != "V1" {
		t.Fatalf("expected verifier ref name 'V1', got %q", ver.Verifier.Name())
	}
	if ver.Required.Name() != "REQ_1" {
		t.Fatalf("expected required ref name 'REQ_1', got %q", ver.Required.Name())
	}
}
