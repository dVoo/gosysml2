// Code generated from SysMLv2Parser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // SysMLv2Parser
import "github.com/antlr4-go/antlr/v4"

// BaseSysMLv2ParserListener is a complete listener for a parse tree produced by SysMLv2Parser.
type BaseSysMLv2ParserListener struct{}

var _ SysMLv2ParserListener = &BaseSysMLv2ParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSysMLv2ParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSysMLv2ParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSysMLv2ParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSysMLv2ParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterEntryRuleRootNamespace is called when production entryRuleRootNamespace is entered.
func (s *BaseSysMLv2ParserListener) EnterEntryRuleRootNamespace(ctx *EntryRuleRootNamespaceContext) {}

// ExitEntryRuleRootNamespace is called when production entryRuleRootNamespace is exited.
func (s *BaseSysMLv2ParserListener) ExitEntryRuleRootNamespace(ctx *EntryRuleRootNamespaceContext) {}

// EnterRootNamespace is called when production rootNamespace is entered.
func (s *BaseSysMLv2ParserListener) EnterRootNamespace(ctx *RootNamespaceContext) {}

// ExitRootNamespace is called when production rootNamespace is exited.
func (s *BaseSysMLv2ParserListener) ExitRootNamespace(ctx *RootNamespaceContext) {}

// EnterIdentification is called when production identification is entered.
func (s *BaseSysMLv2ParserListener) EnterIdentification(ctx *IdentificationContext) {}

// ExitIdentification is called when production identification is exited.
func (s *BaseSysMLv2ParserListener) ExitIdentification(ctx *IdentificationContext) {}

// EnterRelationshipBody is called when production relationshipBody is entered.
func (s *BaseSysMLv2ParserListener) EnterRelationshipBody(ctx *RelationshipBodyContext) {}

// ExitRelationshipBody is called when production relationshipBody is exited.
func (s *BaseSysMLv2ParserListener) ExitRelationshipBody(ctx *RelationshipBodyContext) {}

// EnterDependency is called when production dependency is entered.
func (s *BaseSysMLv2ParserListener) EnterDependency(ctx *DependencyContext) {}

// ExitDependency is called when production dependency is exited.
func (s *BaseSysMLv2ParserListener) ExitDependency(ctx *DependencyContext) {}

// EnterDependencyDeclaration is called when production dependencyDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterDependencyDeclaration(ctx *DependencyDeclarationContext) {}

// ExitDependencyDeclaration is called when production dependencyDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitDependencyDeclaration(ctx *DependencyDeclarationContext) {}

// EnterAnnotation is called when production annotation is entered.
func (s *BaseSysMLv2ParserListener) EnterAnnotation(ctx *AnnotationContext) {}

// ExitAnnotation is called when production annotation is exited.
func (s *BaseSysMLv2ParserListener) ExitAnnotation(ctx *AnnotationContext) {}

// EnterOwnedAnnotation is called when production ownedAnnotation is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedAnnotation(ctx *OwnedAnnotationContext) {}

// ExitOwnedAnnotation is called when production ownedAnnotation is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedAnnotation(ctx *OwnedAnnotationContext) {}

// EnterAnnotatingMember is called when production annotatingMember is entered.
func (s *BaseSysMLv2ParserListener) EnterAnnotatingMember(ctx *AnnotatingMemberContext) {}

// ExitAnnotatingMember is called when production annotatingMember is exited.
func (s *BaseSysMLv2ParserListener) ExitAnnotatingMember(ctx *AnnotatingMemberContext) {}

// EnterAnnotatingElement is called when production annotatingElement is entered.
func (s *BaseSysMLv2ParserListener) EnterAnnotatingElement(ctx *AnnotatingElementContext) {}

// ExitAnnotatingElement is called when production annotatingElement is exited.
func (s *BaseSysMLv2ParserListener) ExitAnnotatingElement(ctx *AnnotatingElementContext) {}

// EnterComment_ is called when production comment_ is entered.
func (s *BaseSysMLv2ParserListener) EnterComment_(ctx *Comment_Context) {}

// ExitComment_ is called when production comment_ is exited.
func (s *BaseSysMLv2ParserListener) ExitComment_(ctx *Comment_Context) {}

// EnterDocumentation is called when production documentation is entered.
func (s *BaseSysMLv2ParserListener) EnterDocumentation(ctx *DocumentationContext) {}

// ExitDocumentation is called when production documentation is exited.
func (s *BaseSysMLv2ParserListener) ExitDocumentation(ctx *DocumentationContext) {}

// EnterTextualRepresentation is called when production textualRepresentation is entered.
func (s *BaseSysMLv2ParserListener) EnterTextualRepresentation(ctx *TextualRepresentationContext) {}

// ExitTextualRepresentation is called when production textualRepresentation is exited.
func (s *BaseSysMLv2ParserListener) ExitTextualRepresentation(ctx *TextualRepresentationContext) {}

// EnterPackage_ is called when production package_ is entered.
func (s *BaseSysMLv2ParserListener) EnterPackage_(ctx *Package_Context) {}

// ExitPackage_ is called when production package_ is exited.
func (s *BaseSysMLv2ParserListener) ExitPackage_(ctx *Package_Context) {}

// EnterLibraryPackage is called when production libraryPackage is entered.
func (s *BaseSysMLv2ParserListener) EnterLibraryPackage(ctx *LibraryPackageContext) {}

// ExitLibraryPackage is called when production libraryPackage is exited.
func (s *BaseSysMLv2ParserListener) ExitLibraryPackage(ctx *LibraryPackageContext) {}

// EnterPackageDeclaration is called when production packageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterPackageDeclaration(ctx *PackageDeclarationContext) {}

// ExitPackageDeclaration is called when production packageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitPackageDeclaration(ctx *PackageDeclarationContext) {}

// EnterPackageBody is called when production packageBody is entered.
func (s *BaseSysMLv2ParserListener) EnterPackageBody(ctx *PackageBodyContext) {}

// ExitPackageBody is called when production packageBody is exited.
func (s *BaseSysMLv2ParserListener) ExitPackageBody(ctx *PackageBodyContext) {}

// EnterPackageBodyElement is called when production packageBodyElement is entered.
func (s *BaseSysMLv2ParserListener) EnterPackageBodyElement(ctx *PackageBodyElementContext) {}

// ExitPackageBodyElement is called when production packageBodyElement is exited.
func (s *BaseSysMLv2ParserListener) ExitPackageBodyElement(ctx *PackageBodyElementContext) {}

// EnterMemberPrefix is called when production memberPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterMemberPrefix(ctx *MemberPrefixContext) {}

// ExitMemberPrefix is called when production memberPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitMemberPrefix(ctx *MemberPrefixContext) {}

// EnterPackageMember is called when production packageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterPackageMember(ctx *PackageMemberContext) {}

// ExitPackageMember is called when production packageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitPackageMember(ctx *PackageMemberContext) {}

// EnterElementFilterMember is called when production elementFilterMember is entered.
func (s *BaseSysMLv2ParserListener) EnterElementFilterMember(ctx *ElementFilterMemberContext) {}

// ExitElementFilterMember is called when production elementFilterMember is exited.
func (s *BaseSysMLv2ParserListener) ExitElementFilterMember(ctx *ElementFilterMemberContext) {}

// EnterAliasMember is called when production aliasMember is entered.
func (s *BaseSysMLv2ParserListener) EnterAliasMember(ctx *AliasMemberContext) {}

// ExitAliasMember is called when production aliasMember is exited.
func (s *BaseSysMLv2ParserListener) ExitAliasMember(ctx *AliasMemberContext) {}

// EnterImport_ is called when production import_ is entered.
func (s *BaseSysMLv2ParserListener) EnterImport_(ctx *Import_Context) {}

// ExitImport_ is called when production import_ is exited.
func (s *BaseSysMLv2ParserListener) ExitImport_(ctx *Import_Context) {}

// EnterImportDeclaration is called when production importDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterImportDeclaration(ctx *ImportDeclarationContext) {}

// ExitImportDeclaration is called when production importDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitImportDeclaration(ctx *ImportDeclarationContext) {}

// EnterMembershipImport is called when production membershipImport is entered.
func (s *BaseSysMLv2ParserListener) EnterMembershipImport(ctx *MembershipImportContext) {}

// ExitMembershipImport is called when production membershipImport is exited.
func (s *BaseSysMLv2ParserListener) ExitMembershipImport(ctx *MembershipImportContext) {}

// EnterNamespaceImport is called when production namespaceImport is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceImport(ctx *NamespaceImportContext) {}

// ExitNamespaceImport is called when production namespaceImport is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceImport(ctx *NamespaceImportContext) {}

// EnterFilterPackage is called when production filterPackage is entered.
func (s *BaseSysMLv2ParserListener) EnterFilterPackage(ctx *FilterPackageContext) {}

// ExitFilterPackage is called when production filterPackage is exited.
func (s *BaseSysMLv2ParserListener) ExitFilterPackage(ctx *FilterPackageContext) {}

// EnterFilterPackageImportPart is called when production filterPackageImportPart is entered.
func (s *BaseSysMLv2ParserListener) EnterFilterPackageImportPart(ctx *FilterPackageImportPartContext) {
}

// ExitFilterPackageImportPart is called when production filterPackageImportPart is exited.
func (s *BaseSysMLv2ParserListener) ExitFilterPackageImportPart(ctx *FilterPackageImportPartContext) {
}

// EnterFilterPackageMember is called when production filterPackageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFilterPackageMember(ctx *FilterPackageMemberContext) {}

// ExitFilterPackageMember is called when production filterPackageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFilterPackageMember(ctx *FilterPackageMemberContext) {}

// EnterVisibilityIndicator is called when production visibilityIndicator is entered.
func (s *BaseSysMLv2ParserListener) EnterVisibilityIndicator(ctx *VisibilityIndicatorContext) {}

// ExitVisibilityIndicator is called when production visibilityIndicator is exited.
func (s *BaseSysMLv2ParserListener) ExitVisibilityIndicator(ctx *VisibilityIndicatorContext) {}

// EnterDefinitionElement is called when production definitionElement is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionElement(ctx *DefinitionElementContext) {}

// ExitDefinitionElement is called when production definitionElement is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionElement(ctx *DefinitionElementContext) {}

// EnterUsageElement is called when production usageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterUsageElement(ctx *UsageElementContext) {}

// ExitUsageElement is called when production usageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitUsageElement(ctx *UsageElementContext) {}

// EnterBasicDefinitionPrefix is called when production basicDefinitionPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterBasicDefinitionPrefix(ctx *BasicDefinitionPrefixContext) {}

// ExitBasicDefinitionPrefix is called when production basicDefinitionPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitBasicDefinitionPrefix(ctx *BasicDefinitionPrefixContext) {}

// EnterDefinitionExtensionKeyword is called when production definitionExtensionKeyword is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionExtensionKeyword(ctx *DefinitionExtensionKeywordContext) {
}

// ExitDefinitionExtensionKeyword is called when production definitionExtensionKeyword is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionExtensionKeyword(ctx *DefinitionExtensionKeywordContext) {
}

// EnterDefinitionPrefix is called when production definitionPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionPrefix(ctx *DefinitionPrefixContext) {}

// ExitDefinitionPrefix is called when production definitionPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionPrefix(ctx *DefinitionPrefixContext) {}

// EnterDefinition is called when production definition is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinition(ctx *DefinitionContext) {}

// ExitDefinition is called when production definition is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinition(ctx *DefinitionContext) {}

// EnterDefinitionDeclaration is called when production definitionDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionDeclaration(ctx *DefinitionDeclarationContext) {}

// ExitDefinitionDeclaration is called when production definitionDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionDeclaration(ctx *DefinitionDeclarationContext) {}

// EnterDefinitionBody is called when production definitionBody is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionBody(ctx *DefinitionBodyContext) {}

// ExitDefinitionBody is called when production definitionBody is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionBody(ctx *DefinitionBodyContext) {}

// EnterDefinitionBodyItem is called when production definitionBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionBodyItem(ctx *DefinitionBodyItemContext) {}

// ExitDefinitionBodyItem is called when production definitionBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionBodyItem(ctx *DefinitionBodyItemContext) {}

// EnterEndFeatureMember is called when production endFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterEndFeatureMember(ctx *EndFeatureMemberContext) {}

// ExitEndFeatureMember is called when production endFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitEndFeatureMember(ctx *EndFeatureMemberContext) {}

// EnterEndFeatureDeclaration is called when production endFeatureDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterEndFeatureDeclaration(ctx *EndFeatureDeclarationContext) {}

// ExitEndFeatureDeclaration is called when production endFeatureDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitEndFeatureDeclaration(ctx *EndFeatureDeclarationContext) {}

// EnterDefinitionMember is called when production definitionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinitionMember(ctx *DefinitionMemberContext) {}

// ExitDefinitionMember is called when production definitionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinitionMember(ctx *DefinitionMemberContext) {}

// EnterVariantUsageMember is called when production variantUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterVariantUsageMember(ctx *VariantUsageMemberContext) {}

// ExitVariantUsageMember is called when production variantUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitVariantUsageMember(ctx *VariantUsageMemberContext) {}

// EnterNonOccurrenceUsageMember is called when production nonOccurrenceUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterNonOccurrenceUsageMember(ctx *NonOccurrenceUsageMemberContext) {
}

// ExitNonOccurrenceUsageMember is called when production nonOccurrenceUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitNonOccurrenceUsageMember(ctx *NonOccurrenceUsageMemberContext) {
}

// EnterOccurrenceUsageMember is called when production occurrenceUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceUsageMember(ctx *OccurrenceUsageMemberContext) {}

// ExitOccurrenceUsageMember is called when production occurrenceUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceUsageMember(ctx *OccurrenceUsageMemberContext) {}

// EnterStructureUsageMember is called when production structureUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterStructureUsageMember(ctx *StructureUsageMemberContext) {}

// ExitStructureUsageMember is called when production structureUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitStructureUsageMember(ctx *StructureUsageMemberContext) {}

// EnterBehaviorUsageMember is called when production behaviorUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterBehaviorUsageMember(ctx *BehaviorUsageMemberContext) {}

// ExitBehaviorUsageMember is called when production behaviorUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitBehaviorUsageMember(ctx *BehaviorUsageMemberContext) {}

// EnterFeatureDirection is called when production featureDirection is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureDirection(ctx *FeatureDirectionContext) {}

// ExitFeatureDirection is called when production featureDirection is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureDirection(ctx *FeatureDirectionContext) {}

// EnterRefPrefix is called when production refPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterRefPrefix(ctx *RefPrefixContext) {}

// ExitRefPrefix is called when production refPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitRefPrefix(ctx *RefPrefixContext) {}

// EnterBasicUsagePrefix is called when production basicUsagePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterBasicUsagePrefix(ctx *BasicUsagePrefixContext) {}

// ExitBasicUsagePrefix is called when production basicUsagePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitBasicUsagePrefix(ctx *BasicUsagePrefixContext) {}

// EnterEndUsagePrefix is called when production endUsagePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterEndUsagePrefix(ctx *EndUsagePrefixContext) {}

// ExitEndUsagePrefix is called when production endUsagePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitEndUsagePrefix(ctx *EndUsagePrefixContext) {}

// EnterOwnedCrossFeatureMember is called when production ownedCrossFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedCrossFeatureMember(ctx *OwnedCrossFeatureMemberContext) {
}

// ExitOwnedCrossFeatureMember is called when production ownedCrossFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedCrossFeatureMember(ctx *OwnedCrossFeatureMemberContext) {
}

// EnterOwnedCrossFeature is called when production ownedCrossFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedCrossFeature(ctx *OwnedCrossFeatureContext) {}

// ExitOwnedCrossFeature is called when production ownedCrossFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedCrossFeature(ctx *OwnedCrossFeatureContext) {}

// EnterUsageExtensionKeyword is called when production usageExtensionKeyword is entered.
func (s *BaseSysMLv2ParserListener) EnterUsageExtensionKeyword(ctx *UsageExtensionKeywordContext) {}

// ExitUsageExtensionKeyword is called when production usageExtensionKeyword is exited.
func (s *BaseSysMLv2ParserListener) ExitUsageExtensionKeyword(ctx *UsageExtensionKeywordContext) {}

// EnterUnextendedUsagePrefix is called when production unextendedUsagePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterUnextendedUsagePrefix(ctx *UnextendedUsagePrefixContext) {}

// ExitUnextendedUsagePrefix is called when production unextendedUsagePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitUnextendedUsagePrefix(ctx *UnextendedUsagePrefixContext) {}

// EnterUsagePrefix is called when production usagePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterUsagePrefix(ctx *UsagePrefixContext) {}

// ExitUsagePrefix is called when production usagePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitUsagePrefix(ctx *UsagePrefixContext) {}

// EnterUsage is called when production usage is entered.
func (s *BaseSysMLv2ParserListener) EnterUsage(ctx *UsageContext) {}

// ExitUsage is called when production usage is exited.
func (s *BaseSysMLv2ParserListener) ExitUsage(ctx *UsageContext) {}

// EnterUsageDeclaration is called when production usageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterUsageDeclaration(ctx *UsageDeclarationContext) {}

// ExitUsageDeclaration is called when production usageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitUsageDeclaration(ctx *UsageDeclarationContext) {}

// EnterUsageCompletion is called when production usageCompletion is entered.
func (s *BaseSysMLv2ParserListener) EnterUsageCompletion(ctx *UsageCompletionContext) {}

// ExitUsageCompletion is called when production usageCompletion is exited.
func (s *BaseSysMLv2ParserListener) ExitUsageCompletion(ctx *UsageCompletionContext) {}

// EnterUsageBody is called when production usageBody is entered.
func (s *BaseSysMLv2ParserListener) EnterUsageBody(ctx *UsageBodyContext) {}

// ExitUsageBody is called when production usageBody is exited.
func (s *BaseSysMLv2ParserListener) ExitUsageBody(ctx *UsageBodyContext) {}

// EnterValuePart is called when production valuePart is entered.
func (s *BaseSysMLv2ParserListener) EnterValuePart(ctx *ValuePartContext) {}

// ExitValuePart is called when production valuePart is exited.
func (s *BaseSysMLv2ParserListener) ExitValuePart(ctx *ValuePartContext) {}

// EnterFeatureValue is called when production featureValue is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureValue(ctx *FeatureValueContext) {}

// ExitFeatureValue is called when production featureValue is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureValue(ctx *FeatureValueContext) {}

// EnterDefaultReferenceUsage is called when production defaultReferenceUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterDefaultReferenceUsage(ctx *DefaultReferenceUsageContext) {}

// ExitDefaultReferenceUsage is called when production defaultReferenceUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitDefaultReferenceUsage(ctx *DefaultReferenceUsageContext) {}

// EnterReferenceUsage is called when production referenceUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterReferenceUsage(ctx *ReferenceUsageContext) {}

// ExitReferenceUsage is called when production referenceUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitReferenceUsage(ctx *ReferenceUsageContext) {}

// EnterVariantReference is called when production variantReference is entered.
func (s *BaseSysMLv2ParserListener) EnterVariantReference(ctx *VariantReferenceContext) {}

// ExitVariantReference is called when production variantReference is exited.
func (s *BaseSysMLv2ParserListener) ExitVariantReference(ctx *VariantReferenceContext) {}

// EnterNonOccurrenceUsageElement is called when production nonOccurrenceUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterNonOccurrenceUsageElement(ctx *NonOccurrenceUsageElementContext) {
}

// ExitNonOccurrenceUsageElement is called when production nonOccurrenceUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitNonOccurrenceUsageElement(ctx *NonOccurrenceUsageElementContext) {
}

// EnterOccurrenceUsageElement is called when production occurrenceUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceUsageElement(ctx *OccurrenceUsageElementContext) {}

// ExitOccurrenceUsageElement is called when production occurrenceUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceUsageElement(ctx *OccurrenceUsageElementContext) {}

// EnterStructureUsageElement is called when production structureUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterStructureUsageElement(ctx *StructureUsageElementContext) {}

// ExitStructureUsageElement is called when production structureUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitStructureUsageElement(ctx *StructureUsageElementContext) {}

// EnterBehaviorUsageElement is called when production behaviorUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterBehaviorUsageElement(ctx *BehaviorUsageElementContext) {}

// ExitBehaviorUsageElement is called when production behaviorUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitBehaviorUsageElement(ctx *BehaviorUsageElementContext) {}

// EnterVariantUsageElement is called when production variantUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterVariantUsageElement(ctx *VariantUsageElementContext) {}

// ExitVariantUsageElement is called when production variantUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitVariantUsageElement(ctx *VariantUsageElementContext) {}

// EnterSubclassificationPart is called when production subclassificationPart is entered.
func (s *BaseSysMLv2ParserListener) EnterSubclassificationPart(ctx *SubclassificationPartContext) {}

// ExitSubclassificationPart is called when production subclassificationPart is exited.
func (s *BaseSysMLv2ParserListener) ExitSubclassificationPart(ctx *SubclassificationPartContext) {}

// EnterOwnedSubclassification is called when production ownedSubclassification is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedSubclassification(ctx *OwnedSubclassificationContext) {}

// ExitOwnedSubclassification is called when production ownedSubclassification is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedSubclassification(ctx *OwnedSubclassificationContext) {}

// EnterFeatureSpecializationPart is called when production featureSpecializationPart is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureSpecializationPart(ctx *FeatureSpecializationPartContext) {
}

// ExitFeatureSpecializationPart is called when production featureSpecializationPart is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureSpecializationPart(ctx *FeatureSpecializationPartContext) {
}

// EnterFeatureSpecialization is called when production featureSpecialization is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureSpecialization(ctx *FeatureSpecializationContext) {}

// ExitFeatureSpecialization is called when production featureSpecialization is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureSpecialization(ctx *FeatureSpecializationContext) {}

// EnterTypings is called when production typings is entered.
func (s *BaseSysMLv2ParserListener) EnterTypings(ctx *TypingsContext) {}

// ExitTypings is called when production typings is exited.
func (s *BaseSysMLv2ParserListener) ExitTypings(ctx *TypingsContext) {}

// EnterTypedBy is called when production typedBy is entered.
func (s *BaseSysMLv2ParserListener) EnterTypedBy(ctx *TypedByContext) {}

// ExitTypedBy is called when production typedBy is exited.
func (s *BaseSysMLv2ParserListener) ExitTypedBy(ctx *TypedByContext) {}

// EnterOwnedFeatureTyping is called when production ownedFeatureTyping is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureTyping(ctx *OwnedFeatureTypingContext) {}

// ExitOwnedFeatureTyping is called when production ownedFeatureTyping is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureTyping(ctx *OwnedFeatureTypingContext) {}

// EnterSubsettings is called when production subsettings is entered.
func (s *BaseSysMLv2ParserListener) EnterSubsettings(ctx *SubsettingsContext) {}

// ExitSubsettings is called when production subsettings is exited.
func (s *BaseSysMLv2ParserListener) ExitSubsettings(ctx *SubsettingsContext) {}

// EnterOwnedSubsetting is called when production ownedSubsetting is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedSubsetting(ctx *OwnedSubsettingContext) {}

// ExitOwnedSubsetting is called when production ownedSubsetting is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedSubsetting(ctx *OwnedSubsettingContext) {}

// EnterReferences is called when production references is entered.
func (s *BaseSysMLv2ParserListener) EnterReferences(ctx *ReferencesContext) {}

// ExitReferences is called when production references is exited.
func (s *BaseSysMLv2ParserListener) ExitReferences(ctx *ReferencesContext) {}

// EnterOwnedReferenceSubsetting is called when production ownedReferenceSubsetting is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedReferenceSubsetting(ctx *OwnedReferenceSubsettingContext) {
}

// ExitOwnedReferenceSubsetting is called when production ownedReferenceSubsetting is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedReferenceSubsetting(ctx *OwnedReferenceSubsettingContext) {
}

// EnterCrosses is called when production crosses is entered.
func (s *BaseSysMLv2ParserListener) EnterCrosses(ctx *CrossesContext) {}

// ExitCrosses is called when production crosses is exited.
func (s *BaseSysMLv2ParserListener) ExitCrosses(ctx *CrossesContext) {}

// EnterOwnedCrossSubsetting is called when production ownedCrossSubsetting is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedCrossSubsetting(ctx *OwnedCrossSubsettingContext) {}

// ExitOwnedCrossSubsetting is called when production ownedCrossSubsetting is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedCrossSubsetting(ctx *OwnedCrossSubsettingContext) {}

// EnterRedefinitions is called when production redefinitions is entered.
func (s *BaseSysMLv2ParserListener) EnterRedefinitions(ctx *RedefinitionsContext) {}

// ExitRedefinitions is called when production redefinitions is exited.
func (s *BaseSysMLv2ParserListener) ExitRedefinitions(ctx *RedefinitionsContext) {}

// EnterOwnedRedefinition is called when production ownedRedefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedRedefinition(ctx *OwnedRedefinitionContext) {}

// ExitOwnedRedefinition is called when production ownedRedefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedRedefinition(ctx *OwnedRedefinitionContext) {}

// EnterOwnedFeatureChain is called when production ownedFeatureChain is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureChain(ctx *OwnedFeatureChainContext) {}

// ExitOwnedFeatureChain is called when production ownedFeatureChain is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureChain(ctx *OwnedFeatureChainContext) {}

// EnterOwnedFeatureChaining is called when production ownedFeatureChaining is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureChaining(ctx *OwnedFeatureChainingContext) {}

// ExitOwnedFeatureChaining is called when production ownedFeatureChaining is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureChaining(ctx *OwnedFeatureChainingContext) {}

// EnterSpecializes is called when production specializes is entered.
func (s *BaseSysMLv2ParserListener) EnterSpecializes(ctx *SpecializesContext) {}

// ExitSpecializes is called when production specializes is exited.
func (s *BaseSysMLv2ParserListener) ExitSpecializes(ctx *SpecializesContext) {}

// EnterDefinedBy is called when production definedBy is entered.
func (s *BaseSysMLv2ParserListener) EnterDefinedBy(ctx *DefinedByContext) {}

// ExitDefinedBy is called when production definedBy is exited.
func (s *BaseSysMLv2ParserListener) ExitDefinedBy(ctx *DefinedByContext) {}

// EnterSubsetsKw is called when production subsetsKw is entered.
func (s *BaseSysMLv2ParserListener) EnterSubsetsKw(ctx *SubsetsKwContext) {}

// ExitSubsetsKw is called when production subsetsKw is exited.
func (s *BaseSysMLv2ParserListener) ExitSubsetsKw(ctx *SubsetsKwContext) {}

// EnterReferencesKw is called when production referencesKw is entered.
func (s *BaseSysMLv2ParserListener) EnterReferencesKw(ctx *ReferencesKwContext) {}

// ExitReferencesKw is called when production referencesKw is exited.
func (s *BaseSysMLv2ParserListener) ExitReferencesKw(ctx *ReferencesKwContext) {}

// EnterCrossesKw is called when production crossesKw is entered.
func (s *BaseSysMLv2ParserListener) EnterCrossesKw(ctx *CrossesKwContext) {}

// ExitCrossesKw is called when production crossesKw is exited.
func (s *BaseSysMLv2ParserListener) ExitCrossesKw(ctx *CrossesKwContext) {}

// EnterRedefinesKw is called when production redefinesKw is entered.
func (s *BaseSysMLv2ParserListener) EnterRedefinesKw(ctx *RedefinesKwContext) {}

// ExitRedefinesKw is called when production redefinesKw is exited.
func (s *BaseSysMLv2ParserListener) ExitRedefinesKw(ctx *RedefinesKwContext) {}

// EnterMultiplicityPart is called when production multiplicityPart is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicityPart(ctx *MultiplicityPartContext) {}

// ExitMultiplicityPart is called when production multiplicityPart is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicityPart(ctx *MultiplicityPartContext) {}

// EnterOwnedMultiplicity is called when production ownedMultiplicity is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedMultiplicity(ctx *OwnedMultiplicityContext) {}

// ExitOwnedMultiplicity is called when production ownedMultiplicity is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedMultiplicity(ctx *OwnedMultiplicityContext) {}

// EnterMultiplicityRange is called when production multiplicityRange is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicityRange(ctx *MultiplicityRangeContext) {}

// ExitMultiplicityRange is called when production multiplicityRange is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicityRange(ctx *MultiplicityRangeContext) {}

// EnterMultiplicityExpressionMember is called when production multiplicityExpressionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicityExpressionMember(ctx *MultiplicityExpressionMemberContext) {
}

// ExitMultiplicityExpressionMember is called when production multiplicityExpressionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicityExpressionMember(ctx *MultiplicityExpressionMemberContext) {
}

// EnterAttributeDefinition is called when production attributeDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterAttributeDefinition(ctx *AttributeDefinitionContext) {}

// ExitAttributeDefinition is called when production attributeDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitAttributeDefinition(ctx *AttributeDefinitionContext) {}

// EnterAttributeUsage is called when production attributeUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterAttributeUsage(ctx *AttributeUsageContext) {}

// ExitAttributeUsage is called when production attributeUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitAttributeUsage(ctx *AttributeUsageContext) {}

// EnterEnumerationDefinition is called when production enumerationDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterEnumerationDefinition(ctx *EnumerationDefinitionContext) {}

// ExitEnumerationDefinition is called when production enumerationDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitEnumerationDefinition(ctx *EnumerationDefinitionContext) {}

// EnterEnumerationBody is called when production enumerationBody is entered.
func (s *BaseSysMLv2ParserListener) EnterEnumerationBody(ctx *EnumerationBodyContext) {}

// ExitEnumerationBody is called when production enumerationBody is exited.
func (s *BaseSysMLv2ParserListener) ExitEnumerationBody(ctx *EnumerationBodyContext) {}

// EnterEnumerationUsageMember is called when production enumerationUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterEnumerationUsageMember(ctx *EnumerationUsageMemberContext) {}

// ExitEnumerationUsageMember is called when production enumerationUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitEnumerationUsageMember(ctx *EnumerationUsageMemberContext) {}

// EnterEnumeratedValue is called when production enumeratedValue is entered.
func (s *BaseSysMLv2ParserListener) EnterEnumeratedValue(ctx *EnumeratedValueContext) {}

// ExitEnumeratedValue is called when production enumeratedValue is exited.
func (s *BaseSysMLv2ParserListener) ExitEnumeratedValue(ctx *EnumeratedValueContext) {}

// EnterEnumerationUsage is called when production enumerationUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterEnumerationUsage(ctx *EnumerationUsageContext) {}

// ExitEnumerationUsage is called when production enumerationUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitEnumerationUsage(ctx *EnumerationUsageContext) {}

// EnterOccurrenceDefinitionPrefix is called when production occurrenceDefinitionPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceDefinitionPrefix(ctx *OccurrenceDefinitionPrefixContext) {
}

// ExitOccurrenceDefinitionPrefix is called when production occurrenceDefinitionPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceDefinitionPrefix(ctx *OccurrenceDefinitionPrefixContext) {
}

// EnterOccurrenceDefinition is called when production occurrenceDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceDefinition(ctx *OccurrenceDefinitionContext) {}

// ExitOccurrenceDefinition is called when production occurrenceDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceDefinition(ctx *OccurrenceDefinitionContext) {}

// EnterIndividualDefinition is called when production individualDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterIndividualDefinition(ctx *IndividualDefinitionContext) {}

// ExitIndividualDefinition is called when production individualDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitIndividualDefinition(ctx *IndividualDefinitionContext) {}

// EnterOccurrenceUsagePrefix is called when production occurrenceUsagePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceUsagePrefix(ctx *OccurrenceUsagePrefixContext) {}

// ExitOccurrenceUsagePrefix is called when production occurrenceUsagePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceUsagePrefix(ctx *OccurrenceUsagePrefixContext) {}

// EnterOccurrenceUsage is called when production occurrenceUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterOccurrenceUsage(ctx *OccurrenceUsageContext) {}

// ExitOccurrenceUsage is called when production occurrenceUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitOccurrenceUsage(ctx *OccurrenceUsageContext) {}

// EnterIndividualUsage is called when production individualUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterIndividualUsage(ctx *IndividualUsageContext) {}

// ExitIndividualUsage is called when production individualUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitIndividualUsage(ctx *IndividualUsageContext) {}

// EnterPortionUsage is called when production portionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterPortionUsage(ctx *PortionUsageContext) {}

// ExitPortionUsage is called when production portionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitPortionUsage(ctx *PortionUsageContext) {}

// EnterPortionKind is called when production portionKind is entered.
func (s *BaseSysMLv2ParserListener) EnterPortionKind(ctx *PortionKindContext) {}

// ExitPortionKind is called when production portionKind is exited.
func (s *BaseSysMLv2ParserListener) ExitPortionKind(ctx *PortionKindContext) {}

// EnterEventOccurrenceUsage is called when production eventOccurrenceUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterEventOccurrenceUsage(ctx *EventOccurrenceUsageContext) {}

// ExitEventOccurrenceUsage is called when production eventOccurrenceUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitEventOccurrenceUsage(ctx *EventOccurrenceUsageContext) {}

// EnterSourceSuccessionMember is called when production sourceSuccessionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterSourceSuccessionMember(ctx *SourceSuccessionMemberContext) {}

// ExitSourceSuccessionMember is called when production sourceSuccessionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitSourceSuccessionMember(ctx *SourceSuccessionMemberContext) {}

// EnterSourceSuccession is called when production sourceSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterSourceSuccession(ctx *SourceSuccessionContext) {}

// ExitSourceSuccession is called when production sourceSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitSourceSuccession(ctx *SourceSuccessionContext) {}

// EnterSourceEndMember is called when production sourceEndMember is entered.
func (s *BaseSysMLv2ParserListener) EnterSourceEndMember(ctx *SourceEndMemberContext) {}

// ExitSourceEndMember is called when production sourceEndMember is exited.
func (s *BaseSysMLv2ParserListener) ExitSourceEndMember(ctx *SourceEndMemberContext) {}

// EnterSourceEnd is called when production sourceEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterSourceEnd(ctx *SourceEndContext) {}

// ExitSourceEnd is called when production sourceEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitSourceEnd(ctx *SourceEndContext) {}

// EnterItemDefinition is called when production itemDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterItemDefinition(ctx *ItemDefinitionContext) {}

// ExitItemDefinition is called when production itemDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitItemDefinition(ctx *ItemDefinitionContext) {}

// EnterItemUsage is called when production itemUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterItemUsage(ctx *ItemUsageContext) {}

// ExitItemUsage is called when production itemUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitItemUsage(ctx *ItemUsageContext) {}

// EnterPartDefinition is called when production partDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterPartDefinition(ctx *PartDefinitionContext) {}

// ExitPartDefinition is called when production partDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitPartDefinition(ctx *PartDefinitionContext) {}

// EnterPartUsage is called when production partUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterPartUsage(ctx *PartUsageContext) {}

// ExitPartUsage is called when production partUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitPartUsage(ctx *PartUsageContext) {}

// EnterPortDefinition is called when production portDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterPortDefinition(ctx *PortDefinitionContext) {}

// ExitPortDefinition is called when production portDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitPortDefinition(ctx *PortDefinitionContext) {}

// EnterPortUsage is called when production portUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterPortUsage(ctx *PortUsageContext) {}

// ExitPortUsage is called when production portUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitPortUsage(ctx *PortUsageContext) {}

// EnterConjugatedPortTyping is called when production conjugatedPortTyping is entered.
func (s *BaseSysMLv2ParserListener) EnterConjugatedPortTyping(ctx *ConjugatedPortTypingContext) {}

// ExitConjugatedPortTyping is called when production conjugatedPortTyping is exited.
func (s *BaseSysMLv2ParserListener) ExitConjugatedPortTyping(ctx *ConjugatedPortTypingContext) {}

// EnterConnectionDefinition is called when production connectionDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectionDefinition(ctx *ConnectionDefinitionContext) {}

// ExitConnectionDefinition is called when production connectionDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectionDefinition(ctx *ConnectionDefinitionContext) {}

// EnterConnectionUsage is called when production connectionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectionUsage(ctx *ConnectionUsageContext) {}

// ExitConnectionUsage is called when production connectionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectionUsage(ctx *ConnectionUsageContext) {}

// EnterConnectorPart is called when production connectorPart is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectorPart(ctx *ConnectorPartContext) {}

// ExitConnectorPart is called when production connectorPart is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectorPart(ctx *ConnectorPartContext) {}

// EnterBinaryConnectorPart is called when production binaryConnectorPart is entered.
func (s *BaseSysMLv2ParserListener) EnterBinaryConnectorPart(ctx *BinaryConnectorPartContext) {}

// ExitBinaryConnectorPart is called when production binaryConnectorPart is exited.
func (s *BaseSysMLv2ParserListener) ExitBinaryConnectorPart(ctx *BinaryConnectorPartContext) {}

// EnterNaryConnectorPart is called when production naryConnectorPart is entered.
func (s *BaseSysMLv2ParserListener) EnterNaryConnectorPart(ctx *NaryConnectorPartContext) {}

// ExitNaryConnectorPart is called when production naryConnectorPart is exited.
func (s *BaseSysMLv2ParserListener) ExitNaryConnectorPart(ctx *NaryConnectorPartContext) {}

// EnterConnectorEndMember is called when production connectorEndMember is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectorEndMember(ctx *ConnectorEndMemberContext) {}

// ExitConnectorEndMember is called when production connectorEndMember is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectorEndMember(ctx *ConnectorEndMemberContext) {}

// EnterConnectorEnd is called when production connectorEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectorEnd(ctx *ConnectorEndContext) {}

// ExitConnectorEnd is called when production connectorEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectorEnd(ctx *ConnectorEndContext) {}

// EnterOwnedCrossMultiplicityMember is called when production ownedCrossMultiplicityMember is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedCrossMultiplicityMember(ctx *OwnedCrossMultiplicityMemberContext) {
}

// ExitOwnedCrossMultiplicityMember is called when production ownedCrossMultiplicityMember is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedCrossMultiplicityMember(ctx *OwnedCrossMultiplicityMemberContext) {
}

// EnterOwnedCrossMultiplicity is called when production ownedCrossMultiplicity is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedCrossMultiplicity(ctx *OwnedCrossMultiplicityContext) {}

// ExitOwnedCrossMultiplicity is called when production ownedCrossMultiplicity is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedCrossMultiplicity(ctx *OwnedCrossMultiplicityContext) {}

// EnterBindingConnectorAsUsage is called when production bindingConnectorAsUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterBindingConnectorAsUsage(ctx *BindingConnectorAsUsageContext) {
}

// ExitBindingConnectorAsUsage is called when production bindingConnectorAsUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitBindingConnectorAsUsage(ctx *BindingConnectorAsUsageContext) {
}

// EnterSuccessionAsUsage is called when production successionAsUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterSuccessionAsUsage(ctx *SuccessionAsUsageContext) {}

// ExitSuccessionAsUsage is called when production successionAsUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitSuccessionAsUsage(ctx *SuccessionAsUsageContext) {}

// EnterInterfaceDefinition is called when production interfaceDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceDefinition(ctx *InterfaceDefinitionContext) {}

// ExitInterfaceDefinition is called when production interfaceDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceDefinition(ctx *InterfaceDefinitionContext) {}

// EnterInterfaceBody is called when production interfaceBody is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceBody(ctx *InterfaceBodyContext) {}

// ExitInterfaceBody is called when production interfaceBody is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceBody(ctx *InterfaceBodyContext) {}

// EnterInterfaceBodyItem is called when production interfaceBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceBodyItem(ctx *InterfaceBodyItemContext) {}

// ExitInterfaceBodyItem is called when production interfaceBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceBodyItem(ctx *InterfaceBodyItemContext) {}

// EnterInterfaceNonOccurrenceUsageMember is called when production interfaceNonOccurrenceUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceNonOccurrenceUsageMember(ctx *InterfaceNonOccurrenceUsageMemberContext) {
}

// ExitInterfaceNonOccurrenceUsageMember is called when production interfaceNonOccurrenceUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceNonOccurrenceUsageMember(ctx *InterfaceNonOccurrenceUsageMemberContext) {
}

// EnterInterfaceNonOccurrenceUsageElement is called when production interfaceNonOccurrenceUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceNonOccurrenceUsageElement(ctx *InterfaceNonOccurrenceUsageElementContext) {
}

// ExitInterfaceNonOccurrenceUsageElement is called when production interfaceNonOccurrenceUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceNonOccurrenceUsageElement(ctx *InterfaceNonOccurrenceUsageElementContext) {
}

// EnterInterfaceOccurrenceUsageMember is called when production interfaceOccurrenceUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceOccurrenceUsageMember(ctx *InterfaceOccurrenceUsageMemberContext) {
}

// ExitInterfaceOccurrenceUsageMember is called when production interfaceOccurrenceUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceOccurrenceUsageMember(ctx *InterfaceOccurrenceUsageMemberContext) {
}

// EnterInterfaceOccurrenceUsageElement is called when production interfaceOccurrenceUsageElement is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceOccurrenceUsageElement(ctx *InterfaceOccurrenceUsageElementContext) {
}

// ExitInterfaceOccurrenceUsageElement is called when production interfaceOccurrenceUsageElement is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceOccurrenceUsageElement(ctx *InterfaceOccurrenceUsageElementContext) {
}

// EnterDefaultInterfaceEnd is called when production defaultInterfaceEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterDefaultInterfaceEnd(ctx *DefaultInterfaceEndContext) {}

// ExitDefaultInterfaceEnd is called when production defaultInterfaceEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitDefaultInterfaceEnd(ctx *DefaultInterfaceEndContext) {}

// EnterInterfaceUsage is called when production interfaceUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceUsage(ctx *InterfaceUsageContext) {}

// ExitInterfaceUsage is called when production interfaceUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceUsage(ctx *InterfaceUsageContext) {}

// EnterInterfaceUsageDeclaration is called when production interfaceUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceUsageDeclaration(ctx *InterfaceUsageDeclarationContext) {
}

// ExitInterfaceUsageDeclaration is called when production interfaceUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceUsageDeclaration(ctx *InterfaceUsageDeclarationContext) {
}

// EnterInterfacePart is called when production interfacePart is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfacePart(ctx *InterfacePartContext) {}

// ExitInterfacePart is called when production interfacePart is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfacePart(ctx *InterfacePartContext) {}

// EnterBinaryInterfacePart is called when production binaryInterfacePart is entered.
func (s *BaseSysMLv2ParserListener) EnterBinaryInterfacePart(ctx *BinaryInterfacePartContext) {}

// ExitBinaryInterfacePart is called when production binaryInterfacePart is exited.
func (s *BaseSysMLv2ParserListener) ExitBinaryInterfacePart(ctx *BinaryInterfacePartContext) {}

// EnterNaryInterfacePart is called when production naryInterfacePart is entered.
func (s *BaseSysMLv2ParserListener) EnterNaryInterfacePart(ctx *NaryInterfacePartContext) {}

// ExitNaryInterfacePart is called when production naryInterfacePart is exited.
func (s *BaseSysMLv2ParserListener) ExitNaryInterfacePart(ctx *NaryInterfacePartContext) {}

// EnterInterfaceEndMember is called when production interfaceEndMember is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceEndMember(ctx *InterfaceEndMemberContext) {}

// ExitInterfaceEndMember is called when production interfaceEndMember is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceEndMember(ctx *InterfaceEndMemberContext) {}

// EnterInterfaceEnd is called when production interfaceEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterInterfaceEnd(ctx *InterfaceEndContext) {}

// ExitInterfaceEnd is called when production interfaceEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitInterfaceEnd(ctx *InterfaceEndContext) {}

// EnterAllocationDefinition is called when production allocationDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterAllocationDefinition(ctx *AllocationDefinitionContext) {}

// ExitAllocationDefinition is called when production allocationDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitAllocationDefinition(ctx *AllocationDefinitionContext) {}

// EnterAllocationUsage is called when production allocationUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterAllocationUsage(ctx *AllocationUsageContext) {}

// ExitAllocationUsage is called when production allocationUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitAllocationUsage(ctx *AllocationUsageContext) {}

// EnterAllocationUsageDeclaration is called when production allocationUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterAllocationUsageDeclaration(ctx *AllocationUsageDeclarationContext) {
}

// ExitAllocationUsageDeclaration is called when production allocationUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitAllocationUsageDeclaration(ctx *AllocationUsageDeclarationContext) {
}

// EnterFlowDefinition is called when production flowDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowDefinition(ctx *FlowDefinitionContext) {}

// ExitFlowDefinition is called when production flowDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowDefinition(ctx *FlowDefinitionContext) {}

// EnterMessage is called when production message is entered.
func (s *BaseSysMLv2ParserListener) EnterMessage(ctx *MessageContext) {}

// ExitMessage is called when production message is exited.
func (s *BaseSysMLv2ParserListener) ExitMessage(ctx *MessageContext) {}

// EnterMessageDeclaration is called when production messageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterMessageDeclaration(ctx *MessageDeclarationContext) {}

// ExitMessageDeclaration is called when production messageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitMessageDeclaration(ctx *MessageDeclarationContext) {}

// EnterMessageEventMember is called when production messageEventMember is entered.
func (s *BaseSysMLv2ParserListener) EnterMessageEventMember(ctx *MessageEventMemberContext) {}

// ExitMessageEventMember is called when production messageEventMember is exited.
func (s *BaseSysMLv2ParserListener) ExitMessageEventMember(ctx *MessageEventMemberContext) {}

// EnterMessageEvent is called when production messageEvent is entered.
func (s *BaseSysMLv2ParserListener) EnterMessageEvent(ctx *MessageEventContext) {}

// ExitMessageEvent is called when production messageEvent is exited.
func (s *BaseSysMLv2ParserListener) ExitMessageEvent(ctx *MessageEventContext) {}

// EnterFlowUsage is called when production flowUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowUsage(ctx *FlowUsageContext) {}

// ExitFlowUsage is called when production flowUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowUsage(ctx *FlowUsageContext) {}

// EnterSuccessionFlowUsage is called when production successionFlowUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterSuccessionFlowUsage(ctx *SuccessionFlowUsageContext) {}

// ExitSuccessionFlowUsage is called when production successionFlowUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitSuccessionFlowUsage(ctx *SuccessionFlowUsageContext) {}

// EnterFlowDeclaration is called when production flowDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowDeclaration(ctx *FlowDeclarationContext) {}

// ExitFlowDeclaration is called when production flowDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowDeclaration(ctx *FlowDeclarationContext) {}

// EnterFlowPayloadFeatureMember is called when production flowPayloadFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowPayloadFeatureMember(ctx *FlowPayloadFeatureMemberContext) {
}

// ExitFlowPayloadFeatureMember is called when production flowPayloadFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowPayloadFeatureMember(ctx *FlowPayloadFeatureMemberContext) {
}

// EnterFlowPayloadFeature is called when production flowPayloadFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowPayloadFeature(ctx *FlowPayloadFeatureContext) {}

// ExitFlowPayloadFeature is called when production flowPayloadFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowPayloadFeature(ctx *FlowPayloadFeatureContext) {}

// EnterPayloadFeature is called when production payloadFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterPayloadFeature(ctx *PayloadFeatureContext) {}

// ExitPayloadFeature is called when production payloadFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitPayloadFeature(ctx *PayloadFeatureContext) {}

// EnterPayloadFeatureSpecializationPart is called when production payloadFeatureSpecializationPart is entered.
func (s *BaseSysMLv2ParserListener) EnterPayloadFeatureSpecializationPart(ctx *PayloadFeatureSpecializationPartContext) {
}

// ExitPayloadFeatureSpecializationPart is called when production payloadFeatureSpecializationPart is exited.
func (s *BaseSysMLv2ParserListener) ExitPayloadFeatureSpecializationPart(ctx *PayloadFeatureSpecializationPartContext) {
}

// EnterFlowEndMember is called when production flowEndMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowEndMember(ctx *FlowEndMemberContext) {}

// ExitFlowEndMember is called when production flowEndMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowEndMember(ctx *FlowEndMemberContext) {}

// EnterFlowEnd is called when production flowEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowEnd(ctx *FlowEndContext) {}

// ExitFlowEnd is called when production flowEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowEnd(ctx *FlowEndContext) {}

// EnterFlowEndSubsetting is called when production flowEndSubsetting is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowEndSubsetting(ctx *FlowEndSubsettingContext) {}

// ExitFlowEndSubsetting is called when production flowEndSubsetting is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowEndSubsetting(ctx *FlowEndSubsettingContext) {}

// EnterFeatureChainPrefix is called when production featureChainPrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureChainPrefix(ctx *FeatureChainPrefixContext) {}

// ExitFeatureChainPrefix is called when production featureChainPrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureChainPrefix(ctx *FeatureChainPrefixContext) {}

// EnterFlowFeatureMember is called when production flowFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowFeatureMember(ctx *FlowFeatureMemberContext) {}

// ExitFlowFeatureMember is called when production flowFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowFeatureMember(ctx *FlowFeatureMemberContext) {}

// EnterFlowFeature is called when production flowFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowFeature(ctx *FlowFeatureContext) {}

// ExitFlowFeature is called when production flowFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowFeature(ctx *FlowFeatureContext) {}

// EnterFlowFeatureRedefinition is called when production flowFeatureRedefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterFlowFeatureRedefinition(ctx *FlowFeatureRedefinitionContext) {
}

// ExitFlowFeatureRedefinition is called when production flowFeatureRedefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitFlowFeatureRedefinition(ctx *FlowFeatureRedefinitionContext) {
}

// EnterActionDefinition is called when production actionDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterActionDefinition(ctx *ActionDefinitionContext) {}

// ExitActionDefinition is called when production actionDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitActionDefinition(ctx *ActionDefinitionContext) {}

// EnterActionBody is called when production actionBody is entered.
func (s *BaseSysMLv2ParserListener) EnterActionBody(ctx *ActionBodyContext) {}

// ExitActionBody is called when production actionBody is exited.
func (s *BaseSysMLv2ParserListener) ExitActionBody(ctx *ActionBodyContext) {}

// EnterActionBodyItem is called when production actionBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterActionBodyItem(ctx *ActionBodyItemContext) {}

// ExitActionBodyItem is called when production actionBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitActionBodyItem(ctx *ActionBodyItemContext) {}

// EnterNonBehaviorBodyItem is called when production nonBehaviorBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterNonBehaviorBodyItem(ctx *NonBehaviorBodyItemContext) {}

// ExitNonBehaviorBodyItem is called when production nonBehaviorBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitNonBehaviorBodyItem(ctx *NonBehaviorBodyItemContext) {}

// EnterActionBehaviorMember is called when production actionBehaviorMember is entered.
func (s *BaseSysMLv2ParserListener) EnterActionBehaviorMember(ctx *ActionBehaviorMemberContext) {}

// ExitActionBehaviorMember is called when production actionBehaviorMember is exited.
func (s *BaseSysMLv2ParserListener) ExitActionBehaviorMember(ctx *ActionBehaviorMemberContext) {}

// EnterInitialNodeMember is called when production initialNodeMember is entered.
func (s *BaseSysMLv2ParserListener) EnterInitialNodeMember(ctx *InitialNodeMemberContext) {}

// ExitInitialNodeMember is called when production initialNodeMember is exited.
func (s *BaseSysMLv2ParserListener) ExitInitialNodeMember(ctx *InitialNodeMemberContext) {}

// EnterActionNodeMember is called when production actionNodeMember is entered.
func (s *BaseSysMLv2ParserListener) EnterActionNodeMember(ctx *ActionNodeMemberContext) {}

// ExitActionNodeMember is called when production actionNodeMember is exited.
func (s *BaseSysMLv2ParserListener) ExitActionNodeMember(ctx *ActionNodeMemberContext) {}

// EnterActionTargetSuccessionMember is called when production actionTargetSuccessionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterActionTargetSuccessionMember(ctx *ActionTargetSuccessionMemberContext) {
}

// ExitActionTargetSuccessionMember is called when production actionTargetSuccessionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitActionTargetSuccessionMember(ctx *ActionTargetSuccessionMemberContext) {
}

// EnterGuardedSuccessionMember is called when production guardedSuccessionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterGuardedSuccessionMember(ctx *GuardedSuccessionMemberContext) {
}

// ExitGuardedSuccessionMember is called when production guardedSuccessionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitGuardedSuccessionMember(ctx *GuardedSuccessionMemberContext) {
}

// EnterActionUsage is called when production actionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterActionUsage(ctx *ActionUsageContext) {}

// ExitActionUsage is called when production actionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitActionUsage(ctx *ActionUsageContext) {}

// EnterActionUsageDeclaration is called when production actionUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterActionUsageDeclaration(ctx *ActionUsageDeclarationContext) {}

// ExitActionUsageDeclaration is called when production actionUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitActionUsageDeclaration(ctx *ActionUsageDeclarationContext) {}

// EnterPerformActionUsage is called when production performActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterPerformActionUsage(ctx *PerformActionUsageContext) {}

// ExitPerformActionUsage is called when production performActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitPerformActionUsage(ctx *PerformActionUsageContext) {}

// EnterPerformActionUsageDeclaration is called when production performActionUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterPerformActionUsageDeclaration(ctx *PerformActionUsageDeclarationContext) {
}

// ExitPerformActionUsageDeclaration is called when production performActionUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitPerformActionUsageDeclaration(ctx *PerformActionUsageDeclarationContext) {
}

// EnterActionNode is called when production actionNode is entered.
func (s *BaseSysMLv2ParserListener) EnterActionNode(ctx *ActionNodeContext) {}

// ExitActionNode is called when production actionNode is exited.
func (s *BaseSysMLv2ParserListener) ExitActionNode(ctx *ActionNodeContext) {}

// EnterActionNodeUsageDeclaration is called when production actionNodeUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterActionNodeUsageDeclaration(ctx *ActionNodeUsageDeclarationContext) {
}

// ExitActionNodeUsageDeclaration is called when production actionNodeUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitActionNodeUsageDeclaration(ctx *ActionNodeUsageDeclarationContext) {
}

// EnterActionNodePrefix is called when production actionNodePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterActionNodePrefix(ctx *ActionNodePrefixContext) {}

// ExitActionNodePrefix is called when production actionNodePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitActionNodePrefix(ctx *ActionNodePrefixContext) {}

// EnterControlNode is called when production controlNode is entered.
func (s *BaseSysMLv2ParserListener) EnterControlNode(ctx *ControlNodeContext) {}

// ExitControlNode is called when production controlNode is exited.
func (s *BaseSysMLv2ParserListener) ExitControlNode(ctx *ControlNodeContext) {}

// EnterControlNodePrefix is called when production controlNodePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterControlNodePrefix(ctx *ControlNodePrefixContext) {}

// ExitControlNodePrefix is called when production controlNodePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitControlNodePrefix(ctx *ControlNodePrefixContext) {}

// EnterMergeNode is called when production mergeNode is entered.
func (s *BaseSysMLv2ParserListener) EnterMergeNode(ctx *MergeNodeContext) {}

// ExitMergeNode is called when production mergeNode is exited.
func (s *BaseSysMLv2ParserListener) ExitMergeNode(ctx *MergeNodeContext) {}

// EnterDecisionNode is called when production decisionNode is entered.
func (s *BaseSysMLv2ParserListener) EnterDecisionNode(ctx *DecisionNodeContext) {}

// ExitDecisionNode is called when production decisionNode is exited.
func (s *BaseSysMLv2ParserListener) ExitDecisionNode(ctx *DecisionNodeContext) {}

// EnterJoinNode is called when production joinNode is entered.
func (s *BaseSysMLv2ParserListener) EnterJoinNode(ctx *JoinNodeContext) {}

// ExitJoinNode is called when production joinNode is exited.
func (s *BaseSysMLv2ParserListener) ExitJoinNode(ctx *JoinNodeContext) {}

// EnterForkNode is called when production forkNode is entered.
func (s *BaseSysMLv2ParserListener) EnterForkNode(ctx *ForkNodeContext) {}

// ExitForkNode is called when production forkNode is exited.
func (s *BaseSysMLv2ParserListener) ExitForkNode(ctx *ForkNodeContext) {}

// EnterAcceptNode is called when production acceptNode is entered.
func (s *BaseSysMLv2ParserListener) EnterAcceptNode(ctx *AcceptNodeContext) {}

// ExitAcceptNode is called when production acceptNode is exited.
func (s *BaseSysMLv2ParserListener) ExitAcceptNode(ctx *AcceptNodeContext) {}

// EnterAcceptNodeDeclaration is called when production acceptNodeDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterAcceptNodeDeclaration(ctx *AcceptNodeDeclarationContext) {}

// ExitAcceptNodeDeclaration is called when production acceptNodeDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitAcceptNodeDeclaration(ctx *AcceptNodeDeclarationContext) {}

// EnterAcceptParameterPart is called when production acceptParameterPart is entered.
func (s *BaseSysMLv2ParserListener) EnterAcceptParameterPart(ctx *AcceptParameterPartContext) {}

// ExitAcceptParameterPart is called when production acceptParameterPart is exited.
func (s *BaseSysMLv2ParserListener) ExitAcceptParameterPart(ctx *AcceptParameterPartContext) {}

// EnterPayloadParameterMember is called when production payloadParameterMember is entered.
func (s *BaseSysMLv2ParserListener) EnterPayloadParameterMember(ctx *PayloadParameterMemberContext) {}

// ExitPayloadParameterMember is called when production payloadParameterMember is exited.
func (s *BaseSysMLv2ParserListener) ExitPayloadParameterMember(ctx *PayloadParameterMemberContext) {}

// EnterPayloadParameter is called when production payloadParameter is entered.
func (s *BaseSysMLv2ParserListener) EnterPayloadParameter(ctx *PayloadParameterContext) {}

// ExitPayloadParameter is called when production payloadParameter is exited.
func (s *BaseSysMLv2ParserListener) ExitPayloadParameter(ctx *PayloadParameterContext) {}

// EnterTriggerValuePart is called when production triggerValuePart is entered.
func (s *BaseSysMLv2ParserListener) EnterTriggerValuePart(ctx *TriggerValuePartContext) {}

// ExitTriggerValuePart is called when production triggerValuePart is exited.
func (s *BaseSysMLv2ParserListener) ExitTriggerValuePart(ctx *TriggerValuePartContext) {}

// EnterTriggerFeatureValue is called when production triggerFeatureValue is entered.
func (s *BaseSysMLv2ParserListener) EnterTriggerFeatureValue(ctx *TriggerFeatureValueContext) {}

// ExitTriggerFeatureValue is called when production triggerFeatureValue is exited.
func (s *BaseSysMLv2ParserListener) ExitTriggerFeatureValue(ctx *TriggerFeatureValueContext) {}

// EnterTriggerExpression is called when production triggerExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterTriggerExpression(ctx *TriggerExpressionContext) {}

// ExitTriggerExpression is called when production triggerExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitTriggerExpression(ctx *TriggerExpressionContext) {}

// EnterSendNode is called when production sendNode is entered.
func (s *BaseSysMLv2ParserListener) EnterSendNode(ctx *SendNodeContext) {}

// ExitSendNode is called when production sendNode is exited.
func (s *BaseSysMLv2ParserListener) ExitSendNode(ctx *SendNodeContext) {}

// EnterSendNodeDeclaration is called when production sendNodeDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterSendNodeDeclaration(ctx *SendNodeDeclarationContext) {}

// ExitSendNodeDeclaration is called when production sendNodeDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitSendNodeDeclaration(ctx *SendNodeDeclarationContext) {}

// EnterSenderReceiverPart is called when production senderReceiverPart is entered.
func (s *BaseSysMLv2ParserListener) EnterSenderReceiverPart(ctx *SenderReceiverPartContext) {}

// ExitSenderReceiverPart is called when production senderReceiverPart is exited.
func (s *BaseSysMLv2ParserListener) ExitSenderReceiverPart(ctx *SenderReceiverPartContext) {}

// EnterNodeParameterMember is called when production nodeParameterMember is entered.
func (s *BaseSysMLv2ParserListener) EnterNodeParameterMember(ctx *NodeParameterMemberContext) {}

// ExitNodeParameterMember is called when production nodeParameterMember is exited.
func (s *BaseSysMLv2ParserListener) ExitNodeParameterMember(ctx *NodeParameterMemberContext) {}

// EnterNodeParameter is called when production nodeParameter is entered.
func (s *BaseSysMLv2ParserListener) EnterNodeParameter(ctx *NodeParameterContext) {}

// ExitNodeParameter is called when production nodeParameter is exited.
func (s *BaseSysMLv2ParserListener) ExitNodeParameter(ctx *NodeParameterContext) {}

// EnterAssignmentNode is called when production assignmentNode is entered.
func (s *BaseSysMLv2ParserListener) EnterAssignmentNode(ctx *AssignmentNodeContext) {}

// ExitAssignmentNode is called when production assignmentNode is exited.
func (s *BaseSysMLv2ParserListener) ExitAssignmentNode(ctx *AssignmentNodeContext) {}

// EnterAssignmentNodeDeclaration is called when production assignmentNodeDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterAssignmentNodeDeclaration(ctx *AssignmentNodeDeclarationContext) {
}

// ExitAssignmentNodeDeclaration is called when production assignmentNodeDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitAssignmentNodeDeclaration(ctx *AssignmentNodeDeclarationContext) {
}

// EnterAssignmentTargetMember is called when production assignmentTargetMember is entered.
func (s *BaseSysMLv2ParserListener) EnterAssignmentTargetMember(ctx *AssignmentTargetMemberContext) {}

// ExitAssignmentTargetMember is called when production assignmentTargetMember is exited.
func (s *BaseSysMLv2ParserListener) ExitAssignmentTargetMember(ctx *AssignmentTargetMemberContext) {}

// EnterAssignmentTargetParameter is called when production assignmentTargetParameter is entered.
func (s *BaseSysMLv2ParserListener) EnterAssignmentTargetParameter(ctx *AssignmentTargetParameterContext) {
}

// ExitAssignmentTargetParameter is called when production assignmentTargetParameter is exited.
func (s *BaseSysMLv2ParserListener) ExitAssignmentTargetParameter(ctx *AssignmentTargetParameterContext) {
}

// EnterFeatureChainMember is called when production featureChainMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureChainMember(ctx *FeatureChainMemberContext) {}

// ExitFeatureChainMember is called when production featureChainMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureChainMember(ctx *FeatureChainMemberContext) {}

// EnterOwnedFeatureChainMember is called when production ownedFeatureChainMember is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureChainMember(ctx *OwnedFeatureChainMemberContext) {
}

// ExitOwnedFeatureChainMember is called when production ownedFeatureChainMember is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureChainMember(ctx *OwnedFeatureChainMemberContext) {
}

// EnterTerminateNode is called when production terminateNode is entered.
func (s *BaseSysMLv2ParserListener) EnterTerminateNode(ctx *TerminateNodeContext) {}

// ExitTerminateNode is called when production terminateNode is exited.
func (s *BaseSysMLv2ParserListener) ExitTerminateNode(ctx *TerminateNodeContext) {}

// EnterIfNode is called when production ifNode is entered.
func (s *BaseSysMLv2ParserListener) EnterIfNode(ctx *IfNodeContext) {}

// ExitIfNode is called when production ifNode is exited.
func (s *BaseSysMLv2ParserListener) ExitIfNode(ctx *IfNodeContext) {}

// EnterActionBodyParameter is called when production actionBodyParameter is entered.
func (s *BaseSysMLv2ParserListener) EnterActionBodyParameter(ctx *ActionBodyParameterContext) {}

// ExitActionBodyParameter is called when production actionBodyParameter is exited.
func (s *BaseSysMLv2ParserListener) ExitActionBodyParameter(ctx *ActionBodyParameterContext) {}

// EnterWhileLoopNode is called when production whileLoopNode is entered.
func (s *BaseSysMLv2ParserListener) EnterWhileLoopNode(ctx *WhileLoopNodeContext) {}

// ExitWhileLoopNode is called when production whileLoopNode is exited.
func (s *BaseSysMLv2ParserListener) ExitWhileLoopNode(ctx *WhileLoopNodeContext) {}

// EnterForLoopNode is called when production forLoopNode is entered.
func (s *BaseSysMLv2ParserListener) EnterForLoopNode(ctx *ForLoopNodeContext) {}

// ExitForLoopNode is called when production forLoopNode is exited.
func (s *BaseSysMLv2ParserListener) ExitForLoopNode(ctx *ForLoopNodeContext) {}

// EnterForVariableDeclaration is called when production forVariableDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterForVariableDeclaration(ctx *ForVariableDeclarationContext) {}

// ExitForVariableDeclaration is called when production forVariableDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitForVariableDeclaration(ctx *ForVariableDeclarationContext) {}

// EnterActionTargetSuccession is called when production actionTargetSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterActionTargetSuccession(ctx *ActionTargetSuccessionContext) {}

// ExitActionTargetSuccession is called when production actionTargetSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitActionTargetSuccession(ctx *ActionTargetSuccessionContext) {}

// EnterTargetSuccession is called when production targetSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterTargetSuccession(ctx *TargetSuccessionContext) {}

// ExitTargetSuccession is called when production targetSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitTargetSuccession(ctx *TargetSuccessionContext) {}

// EnterGuardedTargetSuccession is called when production guardedTargetSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterGuardedTargetSuccession(ctx *GuardedTargetSuccessionContext) {
}

// ExitGuardedTargetSuccession is called when production guardedTargetSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitGuardedTargetSuccession(ctx *GuardedTargetSuccessionContext) {
}

// EnterDefaultTargetSuccession is called when production defaultTargetSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterDefaultTargetSuccession(ctx *DefaultTargetSuccessionContext) {
}

// ExitDefaultTargetSuccession is called when production defaultTargetSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitDefaultTargetSuccession(ctx *DefaultTargetSuccessionContext) {
}

// EnterGuardedSuccession is called when production guardedSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterGuardedSuccession(ctx *GuardedSuccessionContext) {}

// ExitGuardedSuccession is called when production guardedSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitGuardedSuccession(ctx *GuardedSuccessionContext) {}

// EnterStateDefinition is called when production stateDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterStateDefinition(ctx *StateDefinitionContext) {}

// ExitStateDefinition is called when production stateDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitStateDefinition(ctx *StateDefinitionContext) {}

// EnterStateDefBody is called when production stateDefBody is entered.
func (s *BaseSysMLv2ParserListener) EnterStateDefBody(ctx *StateDefBodyContext) {}

// ExitStateDefBody is called when production stateDefBody is exited.
func (s *BaseSysMLv2ParserListener) ExitStateDefBody(ctx *StateDefBodyContext) {}

// EnterStateBodyItem is called when production stateBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterStateBodyItem(ctx *StateBodyItemContext) {}

// ExitStateBodyItem is called when production stateBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitStateBodyItem(ctx *StateBodyItemContext) {}

// EnterEntryActionMember is called when production entryActionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterEntryActionMember(ctx *EntryActionMemberContext) {}

// ExitEntryActionMember is called when production entryActionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitEntryActionMember(ctx *EntryActionMemberContext) {}

// EnterDoActionMember is called when production doActionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterDoActionMember(ctx *DoActionMemberContext) {}

// ExitDoActionMember is called when production doActionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitDoActionMember(ctx *DoActionMemberContext) {}

// EnterExitActionMember is called when production exitActionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterExitActionMember(ctx *ExitActionMemberContext) {}

// ExitExitActionMember is called when production exitActionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitExitActionMember(ctx *ExitActionMemberContext) {}

// EnterEntryTransitionMember is called when production entryTransitionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterEntryTransitionMember(ctx *EntryTransitionMemberContext) {}

// ExitEntryTransitionMember is called when production entryTransitionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitEntryTransitionMember(ctx *EntryTransitionMemberContext) {}

// EnterStateActionUsage is called when production stateActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStateActionUsage(ctx *StateActionUsageContext) {}

// ExitStateActionUsage is called when production stateActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStateActionUsage(ctx *StateActionUsageContext) {}

// EnterStatePerformActionUsage is called when production statePerformActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStatePerformActionUsage(ctx *StatePerformActionUsageContext) {
}

// ExitStatePerformActionUsage is called when production statePerformActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStatePerformActionUsage(ctx *StatePerformActionUsageContext) {
}

// EnterStateAcceptActionUsage is called when production stateAcceptActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStateAcceptActionUsage(ctx *StateAcceptActionUsageContext) {}

// ExitStateAcceptActionUsage is called when production stateAcceptActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStateAcceptActionUsage(ctx *StateAcceptActionUsageContext) {}

// EnterStateSendActionUsage is called when production stateSendActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStateSendActionUsage(ctx *StateSendActionUsageContext) {}

// ExitStateSendActionUsage is called when production stateSendActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStateSendActionUsage(ctx *StateSendActionUsageContext) {}

// EnterStateAssignmentActionUsage is called when production stateAssignmentActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStateAssignmentActionUsage(ctx *StateAssignmentActionUsageContext) {
}

// ExitStateAssignmentActionUsage is called when production stateAssignmentActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStateAssignmentActionUsage(ctx *StateAssignmentActionUsageContext) {
}

// EnterTransitionUsageMember is called when production transitionUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionUsageMember(ctx *TransitionUsageMemberContext) {}

// ExitTransitionUsageMember is called when production transitionUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionUsageMember(ctx *TransitionUsageMemberContext) {}

// EnterTargetTransitionUsageMember is called when production targetTransitionUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterTargetTransitionUsageMember(ctx *TargetTransitionUsageMemberContext) {
}

// ExitTargetTransitionUsageMember is called when production targetTransitionUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitTargetTransitionUsageMember(ctx *TargetTransitionUsageMemberContext) {
}

// EnterStateUsage is called when production stateUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStateUsage(ctx *StateUsageContext) {}

// ExitStateUsage is called when production stateUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStateUsage(ctx *StateUsageContext) {}

// EnterStateUsageBody is called when production stateUsageBody is entered.
func (s *BaseSysMLv2ParserListener) EnterStateUsageBody(ctx *StateUsageBodyContext) {}

// ExitStateUsageBody is called when production stateUsageBody is exited.
func (s *BaseSysMLv2ParserListener) ExitStateUsageBody(ctx *StateUsageBodyContext) {}

// EnterExhibitStateUsage is called when production exhibitStateUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterExhibitStateUsage(ctx *ExhibitStateUsageContext) {}

// ExitExhibitStateUsage is called when production exhibitStateUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitExhibitStateUsage(ctx *ExhibitStateUsageContext) {}

// EnterTransitionUsage is called when production transitionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionUsage(ctx *TransitionUsageContext) {}

// ExitTransitionUsage is called when production transitionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionUsage(ctx *TransitionUsageContext) {}

// EnterTargetTransitionUsage is called when production targetTransitionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTargetTransitionUsage(ctx *TargetTransitionUsageContext) {}

// ExitTargetTransitionUsage is called when production targetTransitionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTargetTransitionUsage(ctx *TargetTransitionUsageContext) {}

// EnterTriggerActionMember is called when production triggerActionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterTriggerActionMember(ctx *TriggerActionMemberContext) {}

// ExitTriggerActionMember is called when production triggerActionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitTriggerActionMember(ctx *TriggerActionMemberContext) {}

// EnterTriggerAction is called when production triggerAction is entered.
func (s *BaseSysMLv2ParserListener) EnterTriggerAction(ctx *TriggerActionContext) {}

// ExitTriggerAction is called when production triggerAction is exited.
func (s *BaseSysMLv2ParserListener) ExitTriggerAction(ctx *TriggerActionContext) {}

// EnterGuardExpressionMember is called when production guardExpressionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterGuardExpressionMember(ctx *GuardExpressionMemberContext) {}

// ExitGuardExpressionMember is called when production guardExpressionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitGuardExpressionMember(ctx *GuardExpressionMemberContext) {}

// EnterEffectBehaviorMember is called when production effectBehaviorMember is entered.
func (s *BaseSysMLv2ParserListener) EnterEffectBehaviorMember(ctx *EffectBehaviorMemberContext) {}

// ExitEffectBehaviorMember is called when production effectBehaviorMember is exited.
func (s *BaseSysMLv2ParserListener) ExitEffectBehaviorMember(ctx *EffectBehaviorMemberContext) {}

// EnterEffectBehaviorUsage is called when production effectBehaviorUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterEffectBehaviorUsage(ctx *EffectBehaviorUsageContext) {}

// ExitEffectBehaviorUsage is called when production effectBehaviorUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitEffectBehaviorUsage(ctx *EffectBehaviorUsageContext) {}

// EnterTransitionPerformActionUsage is called when production transitionPerformActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionPerformActionUsage(ctx *TransitionPerformActionUsageContext) {
}

// ExitTransitionPerformActionUsage is called when production transitionPerformActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionPerformActionUsage(ctx *TransitionPerformActionUsageContext) {
}

// EnterTransitionAcceptActionUsage is called when production transitionAcceptActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionAcceptActionUsage(ctx *TransitionAcceptActionUsageContext) {
}

// ExitTransitionAcceptActionUsage is called when production transitionAcceptActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionAcceptActionUsage(ctx *TransitionAcceptActionUsageContext) {
}

// EnterTransitionSendActionUsage is called when production transitionSendActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionSendActionUsage(ctx *TransitionSendActionUsageContext) {
}

// ExitTransitionSendActionUsage is called when production transitionSendActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionSendActionUsage(ctx *TransitionSendActionUsageContext) {
}

// EnterTransitionAssignmentActionUsage is called when production transitionAssignmentActionUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionAssignmentActionUsage(ctx *TransitionAssignmentActionUsageContext) {
}

// ExitTransitionAssignmentActionUsage is called when production transitionAssignmentActionUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionAssignmentActionUsage(ctx *TransitionAssignmentActionUsageContext) {
}

// EnterTransitionSuccessionMember is called when production transitionSuccessionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionSuccessionMember(ctx *TransitionSuccessionMemberContext) {
}

// ExitTransitionSuccessionMember is called when production transitionSuccessionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionSuccessionMember(ctx *TransitionSuccessionMemberContext) {
}

// EnterTransitionSuccession is called when production transitionSuccession is entered.
func (s *BaseSysMLv2ParserListener) EnterTransitionSuccession(ctx *TransitionSuccessionContext) {}

// ExitTransitionSuccession is called when production transitionSuccession is exited.
func (s *BaseSysMLv2ParserListener) ExitTransitionSuccession(ctx *TransitionSuccessionContext) {}

// EnterCalculationDefinition is called when production calculationDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationDefinition(ctx *CalculationDefinitionContext) {}

// ExitCalculationDefinition is called when production calculationDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationDefinition(ctx *CalculationDefinitionContext) {}

// EnterCalculationUsage is called when production calculationUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationUsage(ctx *CalculationUsageContext) {}

// ExitCalculationUsage is called when production calculationUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationUsage(ctx *CalculationUsageContext) {}

// EnterCalculationBody is called when production calculationBody is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationBody(ctx *CalculationBodyContext) {}

// ExitCalculationBody is called when production calculationBody is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationBody(ctx *CalculationBodyContext) {}

// EnterCalculationBodyPart is called when production calculationBodyPart is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationBodyPart(ctx *CalculationBodyPartContext) {}

// ExitCalculationBodyPart is called when production calculationBodyPart is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationBodyPart(ctx *CalculationBodyPartContext) {}

// EnterCalculationBodyItem is called when production calculationBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationBodyItem(ctx *CalculationBodyItemContext) {}

// ExitCalculationBodyItem is called when production calculationBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationBodyItem(ctx *CalculationBodyItemContext) {}

// EnterReturnParameterMember is called when production returnParameterMember is entered.
func (s *BaseSysMLv2ParserListener) EnterReturnParameterMember(ctx *ReturnParameterMemberContext) {}

// ExitReturnParameterMember is called when production returnParameterMember is exited.
func (s *BaseSysMLv2ParserListener) ExitReturnParameterMember(ctx *ReturnParameterMemberContext) {}

// EnterResultExpressionMember is called when production resultExpressionMember is entered.
func (s *BaseSysMLv2ParserListener) EnterResultExpressionMember(ctx *ResultExpressionMemberContext) {}

// ExitResultExpressionMember is called when production resultExpressionMember is exited.
func (s *BaseSysMLv2ParserListener) ExitResultExpressionMember(ctx *ResultExpressionMemberContext) {}

// EnterConstraintDefinition is called when production constraintDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterConstraintDefinition(ctx *ConstraintDefinitionContext) {}

// ExitConstraintDefinition is called when production constraintDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitConstraintDefinition(ctx *ConstraintDefinitionContext) {}

// EnterConstraintUsage is called when production constraintUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterConstraintUsage(ctx *ConstraintUsageContext) {}

// ExitConstraintUsage is called when production constraintUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitConstraintUsage(ctx *ConstraintUsageContext) {}

// EnterAssertConstraintUsage is called when production assertConstraintUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterAssertConstraintUsage(ctx *AssertConstraintUsageContext) {}

// ExitAssertConstraintUsage is called when production assertConstraintUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitAssertConstraintUsage(ctx *AssertConstraintUsageContext) {}

// EnterConstraintUsageDeclaration is called when production constraintUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterConstraintUsageDeclaration(ctx *ConstraintUsageDeclarationContext) {
}

// ExitConstraintUsageDeclaration is called when production constraintUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitConstraintUsageDeclaration(ctx *ConstraintUsageDeclarationContext) {
}

// EnterRequirementDefinition is called when production requirementDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementDefinition(ctx *RequirementDefinitionContext) {}

// ExitRequirementDefinition is called when production requirementDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementDefinition(ctx *RequirementDefinitionContext) {}

// EnterRequirementBody is called when production requirementBody is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementBody(ctx *RequirementBodyContext) {}

// ExitRequirementBody is called when production requirementBody is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementBody(ctx *RequirementBodyContext) {}

// EnterRequirementBodyItem is called when production requirementBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementBodyItem(ctx *RequirementBodyItemContext) {}

// ExitRequirementBodyItem is called when production requirementBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementBodyItem(ctx *RequirementBodyItemContext) {}

// EnterSubjectMember is called when production subjectMember is entered.
func (s *BaseSysMLv2ParserListener) EnterSubjectMember(ctx *SubjectMemberContext) {}

// ExitSubjectMember is called when production subjectMember is exited.
func (s *BaseSysMLv2ParserListener) ExitSubjectMember(ctx *SubjectMemberContext) {}

// EnterSubjectUsage is called when production subjectUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterSubjectUsage(ctx *SubjectUsageContext) {}

// ExitSubjectUsage is called when production subjectUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitSubjectUsage(ctx *SubjectUsageContext) {}

// EnterRequirementConstraintMember is called when production requirementConstraintMember is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementConstraintMember(ctx *RequirementConstraintMemberContext) {
}

// ExitRequirementConstraintMember is called when production requirementConstraintMember is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementConstraintMember(ctx *RequirementConstraintMemberContext) {
}

// EnterRequirementKind is called when production requirementKind is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementKind(ctx *RequirementKindContext) {}

// ExitRequirementKind is called when production requirementKind is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementKind(ctx *RequirementKindContext) {}

// EnterRequirementConstraintUsage is called when production requirementConstraintUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementConstraintUsage(ctx *RequirementConstraintUsageContext) {
}

// ExitRequirementConstraintUsage is called when production requirementConstraintUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementConstraintUsage(ctx *RequirementConstraintUsageContext) {
}

// EnterFramedConcernMember is called when production framedConcernMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFramedConcernMember(ctx *FramedConcernMemberContext) {}

// ExitFramedConcernMember is called when production framedConcernMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFramedConcernMember(ctx *FramedConcernMemberContext) {}

// EnterFramedConcernUsage is called when production framedConcernUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterFramedConcernUsage(ctx *FramedConcernUsageContext) {}

// ExitFramedConcernUsage is called when production framedConcernUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitFramedConcernUsage(ctx *FramedConcernUsageContext) {}

// EnterCalculationUsageDeclaration is called when production calculationUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterCalculationUsageDeclaration(ctx *CalculationUsageDeclarationContext) {
}

// ExitCalculationUsageDeclaration is called when production calculationUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitCalculationUsageDeclaration(ctx *CalculationUsageDeclarationContext) {
}

// EnterActorMember is called when production actorMember is entered.
func (s *BaseSysMLv2ParserListener) EnterActorMember(ctx *ActorMemberContext) {}

// ExitActorMember is called when production actorMember is exited.
func (s *BaseSysMLv2ParserListener) ExitActorMember(ctx *ActorMemberContext) {}

// EnterActorUsage is called when production actorUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterActorUsage(ctx *ActorUsageContext) {}

// ExitActorUsage is called when production actorUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitActorUsage(ctx *ActorUsageContext) {}

// EnterStakeholderMember is called when production stakeholderMember is entered.
func (s *BaseSysMLv2ParserListener) EnterStakeholderMember(ctx *StakeholderMemberContext) {}

// ExitStakeholderMember is called when production stakeholderMember is exited.
func (s *BaseSysMLv2ParserListener) ExitStakeholderMember(ctx *StakeholderMemberContext) {}

// EnterStakeholderUsage is called when production stakeholderUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterStakeholderUsage(ctx *StakeholderUsageContext) {}

// ExitStakeholderUsage is called when production stakeholderUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitStakeholderUsage(ctx *StakeholderUsageContext) {}

// EnterRequirementUsage is called when production requirementUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementUsage(ctx *RequirementUsageContext) {}

// ExitRequirementUsage is called when production requirementUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementUsage(ctx *RequirementUsageContext) {}

// EnterSatisfyRequirementUsage is called when production satisfyRequirementUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterSatisfyRequirementUsage(ctx *SatisfyRequirementUsageContext) {
}

// ExitSatisfyRequirementUsage is called when production satisfyRequirementUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitSatisfyRequirementUsage(ctx *SatisfyRequirementUsageContext) {
}

// EnterSatisfactionSubjectMember is called when production satisfactionSubjectMember is entered.
func (s *BaseSysMLv2ParserListener) EnterSatisfactionSubjectMember(ctx *SatisfactionSubjectMemberContext) {
}

// ExitSatisfactionSubjectMember is called when production satisfactionSubjectMember is exited.
func (s *BaseSysMLv2ParserListener) ExitSatisfactionSubjectMember(ctx *SatisfactionSubjectMemberContext) {
}

// EnterSatisfactionParameter is called when production satisfactionParameter is entered.
func (s *BaseSysMLv2ParserListener) EnterSatisfactionParameter(ctx *SatisfactionParameterContext) {}

// ExitSatisfactionParameter is called when production satisfactionParameter is exited.
func (s *BaseSysMLv2ParserListener) ExitSatisfactionParameter(ctx *SatisfactionParameterContext) {}

// EnterConcernDefinition is called when production concernDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterConcernDefinition(ctx *ConcernDefinitionContext) {}

// ExitConcernDefinition is called when production concernDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitConcernDefinition(ctx *ConcernDefinitionContext) {}

// EnterConcernUsage is called when production concernUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterConcernUsage(ctx *ConcernUsageContext) {}

// ExitConcernUsage is called when production concernUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitConcernUsage(ctx *ConcernUsageContext) {}

// EnterCaseDefinition is called when production caseDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterCaseDefinition(ctx *CaseDefinitionContext) {}

// ExitCaseDefinition is called when production caseDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitCaseDefinition(ctx *CaseDefinitionContext) {}

// EnterCaseUsage is called when production caseUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterCaseUsage(ctx *CaseUsageContext) {}

// ExitCaseUsage is called when production caseUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitCaseUsage(ctx *CaseUsageContext) {}

// EnterCaseBody is called when production caseBody is entered.
func (s *BaseSysMLv2ParserListener) EnterCaseBody(ctx *CaseBodyContext) {}

// ExitCaseBody is called when production caseBody is exited.
func (s *BaseSysMLv2ParserListener) ExitCaseBody(ctx *CaseBodyContext) {}

// EnterCaseBodyItem is called when production caseBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterCaseBodyItem(ctx *CaseBodyItemContext) {}

// ExitCaseBodyItem is called when production caseBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitCaseBodyItem(ctx *CaseBodyItemContext) {}

// EnterObjectiveMember is called when production objectiveMember is entered.
func (s *BaseSysMLv2ParserListener) EnterObjectiveMember(ctx *ObjectiveMemberContext) {}

// ExitObjectiveMember is called when production objectiveMember is exited.
func (s *BaseSysMLv2ParserListener) ExitObjectiveMember(ctx *ObjectiveMemberContext) {}

// EnterObjectiveRequirementUsage is called when production objectiveRequirementUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterObjectiveRequirementUsage(ctx *ObjectiveRequirementUsageContext) {
}

// ExitObjectiveRequirementUsage is called when production objectiveRequirementUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitObjectiveRequirementUsage(ctx *ObjectiveRequirementUsageContext) {
}

// EnterAnalysisCaseDefinition is called when production analysisCaseDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterAnalysisCaseDefinition(ctx *AnalysisCaseDefinitionContext) {}

// ExitAnalysisCaseDefinition is called when production analysisCaseDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitAnalysisCaseDefinition(ctx *AnalysisCaseDefinitionContext) {}

// EnterAnalysisCaseUsage is called when production analysisCaseUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterAnalysisCaseUsage(ctx *AnalysisCaseUsageContext) {}

// ExitAnalysisCaseUsage is called when production analysisCaseUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitAnalysisCaseUsage(ctx *AnalysisCaseUsageContext) {}

// EnterVerificationCaseDefinition is called when production verificationCaseDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterVerificationCaseDefinition(ctx *VerificationCaseDefinitionContext) {
}

// ExitVerificationCaseDefinition is called when production verificationCaseDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitVerificationCaseDefinition(ctx *VerificationCaseDefinitionContext) {
}

// EnterVerificationCaseUsage is called when production verificationCaseUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterVerificationCaseUsage(ctx *VerificationCaseUsageContext) {}

// ExitVerificationCaseUsage is called when production verificationCaseUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitVerificationCaseUsage(ctx *VerificationCaseUsageContext) {}

// EnterRequirementVerificationMember is called when production requirementVerificationMember is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementVerificationMember(ctx *RequirementVerificationMemberContext) {
}

// ExitRequirementVerificationMember is called when production requirementVerificationMember is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementVerificationMember(ctx *RequirementVerificationMemberContext) {
}

// EnterRequirementVerificationUsage is called when production requirementVerificationUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterRequirementVerificationUsage(ctx *RequirementVerificationUsageContext) {
}

// ExitRequirementVerificationUsage is called when production requirementVerificationUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitRequirementVerificationUsage(ctx *RequirementVerificationUsageContext) {
}

// EnterUseCaseDefinition is called when production useCaseDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterUseCaseDefinition(ctx *UseCaseDefinitionContext) {}

// ExitUseCaseDefinition is called when production useCaseDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitUseCaseDefinition(ctx *UseCaseDefinitionContext) {}

// EnterUseCaseUsage is called when production useCaseUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterUseCaseUsage(ctx *UseCaseUsageContext) {}

// ExitUseCaseUsage is called when production useCaseUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitUseCaseUsage(ctx *UseCaseUsageContext) {}

// EnterIncludeUseCaseUsage is called when production includeUseCaseUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterIncludeUseCaseUsage(ctx *IncludeUseCaseUsageContext) {}

// ExitIncludeUseCaseUsage is called when production includeUseCaseUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitIncludeUseCaseUsage(ctx *IncludeUseCaseUsageContext) {}

// EnterViewDefinition is called when production viewDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterViewDefinition(ctx *ViewDefinitionContext) {}

// ExitViewDefinition is called when production viewDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitViewDefinition(ctx *ViewDefinitionContext) {}

// EnterViewDefinitionBody is called when production viewDefinitionBody is entered.
func (s *BaseSysMLv2ParserListener) EnterViewDefinitionBody(ctx *ViewDefinitionBodyContext) {}

// ExitViewDefinitionBody is called when production viewDefinitionBody is exited.
func (s *BaseSysMLv2ParserListener) ExitViewDefinitionBody(ctx *ViewDefinitionBodyContext) {}

// EnterViewDefinitionBodyItem is called when production viewDefinitionBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterViewDefinitionBodyItem(ctx *ViewDefinitionBodyItemContext) {}

// ExitViewDefinitionBodyItem is called when production viewDefinitionBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitViewDefinitionBodyItem(ctx *ViewDefinitionBodyItemContext) {}

// EnterViewRenderingMember is called when production viewRenderingMember is entered.
func (s *BaseSysMLv2ParserListener) EnterViewRenderingMember(ctx *ViewRenderingMemberContext) {}

// ExitViewRenderingMember is called when production viewRenderingMember is exited.
func (s *BaseSysMLv2ParserListener) ExitViewRenderingMember(ctx *ViewRenderingMemberContext) {}

// EnterViewRenderingUsage is called when production viewRenderingUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterViewRenderingUsage(ctx *ViewRenderingUsageContext) {}

// ExitViewRenderingUsage is called when production viewRenderingUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitViewRenderingUsage(ctx *ViewRenderingUsageContext) {}

// EnterViewUsage is called when production viewUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterViewUsage(ctx *ViewUsageContext) {}

// ExitViewUsage is called when production viewUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitViewUsage(ctx *ViewUsageContext) {}

// EnterViewBody is called when production viewBody is entered.
func (s *BaseSysMLv2ParserListener) EnterViewBody(ctx *ViewBodyContext) {}

// ExitViewBody is called when production viewBody is exited.
func (s *BaseSysMLv2ParserListener) ExitViewBody(ctx *ViewBodyContext) {}

// EnterViewBodyItem is called when production viewBodyItem is entered.
func (s *BaseSysMLv2ParserListener) EnterViewBodyItem(ctx *ViewBodyItemContext) {}

// ExitViewBodyItem is called when production viewBodyItem is exited.
func (s *BaseSysMLv2ParserListener) ExitViewBodyItem(ctx *ViewBodyItemContext) {}

// EnterExpose is called when production expose is entered.
func (s *BaseSysMLv2ParserListener) EnterExpose(ctx *ExposeContext) {}

// ExitExpose is called when production expose is exited.
func (s *BaseSysMLv2ParserListener) ExitExpose(ctx *ExposeContext) {}

// EnterMembershipExpose is called when production membershipExpose is entered.
func (s *BaseSysMLv2ParserListener) EnterMembershipExpose(ctx *MembershipExposeContext) {}

// ExitMembershipExpose is called when production membershipExpose is exited.
func (s *BaseSysMLv2ParserListener) ExitMembershipExpose(ctx *MembershipExposeContext) {}

// EnterNamespaceExpose is called when production namespaceExpose is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceExpose(ctx *NamespaceExposeContext) {}

// ExitNamespaceExpose is called when production namespaceExpose is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceExpose(ctx *NamespaceExposeContext) {}

// EnterViewpointDefinition is called when production viewpointDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterViewpointDefinition(ctx *ViewpointDefinitionContext) {}

// ExitViewpointDefinition is called when production viewpointDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitViewpointDefinition(ctx *ViewpointDefinitionContext) {}

// EnterViewpointUsage is called when production viewpointUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterViewpointUsage(ctx *ViewpointUsageContext) {}

// ExitViewpointUsage is called when production viewpointUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitViewpointUsage(ctx *ViewpointUsageContext) {}

// EnterRenderingDefinition is called when production renderingDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterRenderingDefinition(ctx *RenderingDefinitionContext) {}

// ExitRenderingDefinition is called when production renderingDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitRenderingDefinition(ctx *RenderingDefinitionContext) {}

// EnterRenderingUsage is called when production renderingUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterRenderingUsage(ctx *RenderingUsageContext) {}

// ExitRenderingUsage is called when production renderingUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitRenderingUsage(ctx *RenderingUsageContext) {}

// EnterMetadataDefinition is called when production metadataDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataDefinition(ctx *MetadataDefinitionContext) {}

// ExitMetadataDefinition is called when production metadataDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataDefinition(ctx *MetadataDefinitionContext) {}

// EnterPrefixMetadataAnnotation is called when production prefixMetadataAnnotation is entered.
func (s *BaseSysMLv2ParserListener) EnterPrefixMetadataAnnotation(ctx *PrefixMetadataAnnotationContext) {
}

// ExitPrefixMetadataAnnotation is called when production prefixMetadataAnnotation is exited.
func (s *BaseSysMLv2ParserListener) ExitPrefixMetadataAnnotation(ctx *PrefixMetadataAnnotationContext) {
}

// EnterPrefixMetadataMember is called when production prefixMetadataMember is entered.
func (s *BaseSysMLv2ParserListener) EnterPrefixMetadataMember(ctx *PrefixMetadataMemberContext) {}

// ExitPrefixMetadataMember is called when production prefixMetadataMember is exited.
func (s *BaseSysMLv2ParserListener) ExitPrefixMetadataMember(ctx *PrefixMetadataMemberContext) {}

// EnterPrefixMetadataUsage is called when production prefixMetadataUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterPrefixMetadataUsage(ctx *PrefixMetadataUsageContext) {}

// ExitPrefixMetadataUsage is called when production prefixMetadataUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitPrefixMetadataUsage(ctx *PrefixMetadataUsageContext) {}

// EnterMetadataUsage is called when production metadataUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataUsage(ctx *MetadataUsageContext) {}

// ExitMetadataUsage is called when production metadataUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataUsage(ctx *MetadataUsageContext) {}

// EnterMetadataUsageDeclaration is called when production metadataUsageDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataUsageDeclaration(ctx *MetadataUsageDeclarationContext) {
}

// ExitMetadataUsageDeclaration is called when production metadataUsageDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataUsageDeclaration(ctx *MetadataUsageDeclarationContext) {
}

// EnterMetadataBody is called when production metadataBody is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataBody(ctx *MetadataBodyContext) {}

// ExitMetadataBody is called when production metadataBody is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataBody(ctx *MetadataBodyContext) {}

// EnterMetadataBodyUsageMember is called when production metadataBodyUsageMember is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataBodyUsageMember(ctx *MetadataBodyUsageMemberContext) {
}

// ExitMetadataBodyUsageMember is called when production metadataBodyUsageMember is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataBodyUsageMember(ctx *MetadataBodyUsageMemberContext) {
}

// EnterMetadataBodyUsage is called when production metadataBodyUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataBodyUsage(ctx *MetadataBodyUsageContext) {}

// ExitMetadataBodyUsage is called when production metadataBodyUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataBodyUsage(ctx *MetadataBodyUsageContext) {}

// EnterMetadataFeature is called when production metadataFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataFeature(ctx *MetadataFeatureContext) {}

// ExitMetadataFeature is called when production metadataFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataFeature(ctx *MetadataFeatureContext) {}

// EnterExtendedDefinition is called when production extendedDefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterExtendedDefinition(ctx *ExtendedDefinitionContext) {}

// ExitExtendedDefinition is called when production extendedDefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitExtendedDefinition(ctx *ExtendedDefinitionContext) {}

// EnterExtendedUsage is called when production extendedUsage is entered.
func (s *BaseSysMLv2ParserListener) EnterExtendedUsage(ctx *ExtendedUsageContext) {}

// ExitExtendedUsage is called when production extendedUsage is exited.
func (s *BaseSysMLv2ParserListener) ExitExtendedUsage(ctx *ExtendedUsageContext) {}

// EnterOwnedExpression is called when production ownedExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedExpression(ctx *OwnedExpressionContext) {}

// ExitOwnedExpression is called when production ownedExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedExpression(ctx *OwnedExpressionContext) {}

// EnterConditionalBinaryOperator is called when production conditionalBinaryOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterConditionalBinaryOperator(ctx *ConditionalBinaryOperatorContext) {
}

// ExitConditionalBinaryOperator is called when production conditionalBinaryOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitConditionalBinaryOperator(ctx *ConditionalBinaryOperatorContext) {
}

// EnterBinaryOperator is called when production binaryOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterBinaryOperator(ctx *BinaryOperatorContext) {}

// ExitBinaryOperator is called when production binaryOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitBinaryOperator(ctx *BinaryOperatorContext) {}

// EnterUnaryOperator is called when production unaryOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterUnaryOperator(ctx *UnaryOperatorContext) {}

// ExitUnaryOperator is called when production unaryOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitUnaryOperator(ctx *UnaryOperatorContext) {}

// EnterClassificationTestOperator is called when production classificationTestOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterClassificationTestOperator(ctx *ClassificationTestOperatorContext) {
}

// ExitClassificationTestOperator is called when production classificationTestOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitClassificationTestOperator(ctx *ClassificationTestOperatorContext) {
}

// EnterCastOperator is called when production castOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterCastOperator(ctx *CastOperatorContext) {}

// ExitCastOperator is called when production castOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitCastOperator(ctx *CastOperatorContext) {}

// EnterMetaCastOperator is called when production metaCastOperator is entered.
func (s *BaseSysMLv2ParserListener) EnterMetaCastOperator(ctx *MetaCastOperatorContext) {}

// ExitMetaCastOperator is called when production metaCastOperator is exited.
func (s *BaseSysMLv2ParserListener) ExitMetaCastOperator(ctx *MetaCastOperatorContext) {}

// EnterTypeReference is called when production typeReference is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeReference(ctx *TypeReferenceContext) {}

// ExitTypeReference is called when production typeReference is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeReference(ctx *TypeReferenceContext) {}

// EnterPrimaryExpression is called when production primaryExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterPrimaryExpression(ctx *PrimaryExpressionContext) {}

// ExitPrimaryExpression is called when production primaryExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitPrimaryExpression(ctx *PrimaryExpressionContext) {}

// EnterPrimaryExpressionSuffix is called when production primaryExpressionSuffix is entered.
func (s *BaseSysMLv2ParserListener) EnterPrimaryExpressionSuffix(ctx *PrimaryExpressionSuffixContext) {
}

// ExitPrimaryExpressionSuffix is called when production primaryExpressionSuffix is exited.
func (s *BaseSysMLv2ParserListener) ExitPrimaryExpressionSuffix(ctx *PrimaryExpressionSuffixContext) {
}

// EnterBaseExpression is called when production baseExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterBaseExpression(ctx *BaseExpressionContext) {}

// ExitBaseExpression is called when production baseExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitBaseExpression(ctx *BaseExpressionContext) {}

// EnterNullExpression is called when production nullExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterNullExpression(ctx *NullExpressionContext) {}

// ExitNullExpression is called when production nullExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitNullExpression(ctx *NullExpressionContext) {}

// EnterFeatureReferenceExpression is called when production featureReferenceExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureReferenceExpression(ctx *FeatureReferenceExpressionContext) {
}

// ExitFeatureReferenceExpression is called when production featureReferenceExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureReferenceExpression(ctx *FeatureReferenceExpressionContext) {
}

// EnterMetadataAccessExpression is called when production metadataAccessExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterMetadataAccessExpression(ctx *MetadataAccessExpressionContext) {
}

// ExitMetadataAccessExpression is called when production metadataAccessExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitMetadataAccessExpression(ctx *MetadataAccessExpressionContext) {
}

// EnterInvocationExpression is called when production invocationExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterInvocationExpression(ctx *InvocationExpressionContext) {}

// ExitInvocationExpression is called when production invocationExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitInvocationExpression(ctx *InvocationExpressionContext) {}

// EnterConstructorExpression is called when production constructorExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterConstructorExpression(ctx *ConstructorExpressionContext) {}

// ExitConstructorExpression is called when production constructorExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitConstructorExpression(ctx *ConstructorExpressionContext) {}

// EnterBodyExpression is called when production bodyExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterBodyExpression(ctx *BodyExpressionContext) {}

// ExitBodyExpression is called when production bodyExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitBodyExpression(ctx *BodyExpressionContext) {}

// EnterSequenceExpression is called when production sequenceExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterSequenceExpression(ctx *SequenceExpressionContext) {}

// ExitSequenceExpression is called when production sequenceExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitSequenceExpression(ctx *SequenceExpressionContext) {}

// EnterSequenceExpressionList is called when production sequenceExpressionList is entered.
func (s *BaseSysMLv2ParserListener) EnterSequenceExpressionList(ctx *SequenceExpressionListContext) {}

// ExitSequenceExpressionList is called when production sequenceExpressionList is exited.
func (s *BaseSysMLv2ParserListener) ExitSequenceExpressionList(ctx *SequenceExpressionListContext) {}

// EnterArgumentList is called when production argumentList is entered.
func (s *BaseSysMLv2ParserListener) EnterArgumentList(ctx *ArgumentListContext) {}

// ExitArgumentList is called when production argumentList is exited.
func (s *BaseSysMLv2ParserListener) ExitArgumentList(ctx *ArgumentListContext) {}

// EnterPositionalArgumentList is called when production positionalArgumentList is entered.
func (s *BaseSysMLv2ParserListener) EnterPositionalArgumentList(ctx *PositionalArgumentListContext) {}

// ExitPositionalArgumentList is called when production positionalArgumentList is exited.
func (s *BaseSysMLv2ParserListener) ExitPositionalArgumentList(ctx *PositionalArgumentListContext) {}

// EnterNamedArgumentList is called when production namedArgumentList is entered.
func (s *BaseSysMLv2ParserListener) EnterNamedArgumentList(ctx *NamedArgumentListContext) {}

// ExitNamedArgumentList is called when production namedArgumentList is exited.
func (s *BaseSysMLv2ParserListener) ExitNamedArgumentList(ctx *NamedArgumentListContext) {}

// EnterNamedArgument is called when production namedArgument is entered.
func (s *BaseSysMLv2ParserListener) EnterNamedArgument(ctx *NamedArgumentContext) {}

// ExitNamedArgument is called when production namedArgument is exited.
func (s *BaseSysMLv2ParserListener) ExitNamedArgument(ctx *NamedArgumentContext) {}

// EnterFunctionBodyPart is called when production functionBodyPart is entered.
func (s *BaseSysMLv2ParserListener) EnterFunctionBodyPart(ctx *FunctionBodyPartContext) {}

// ExitFunctionBodyPart is called when production functionBodyPart is exited.
func (s *BaseSysMLv2ParserListener) ExitFunctionBodyPart(ctx *FunctionBodyPartContext) {}

// EnterResultExpressionMemberOpt is called when production resultExpressionMemberOpt is entered.
func (s *BaseSysMLv2ParserListener) EnterResultExpressionMemberOpt(ctx *ResultExpressionMemberOptContext) {
}

// ExitResultExpressionMemberOpt is called when production resultExpressionMemberOpt is exited.
func (s *BaseSysMLv2ParserListener) ExitResultExpressionMemberOpt(ctx *ResultExpressionMemberOptContext) {
}

// EnterTypeBodyElement is called when production typeBodyElement is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeBodyElement(ctx *TypeBodyElementContext) {}

// ExitTypeBodyElement is called when production typeBodyElement is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeBodyElement(ctx *TypeBodyElementContext) {}

// EnterNonFeatureMember is called when production nonFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterNonFeatureMember(ctx *NonFeatureMemberContext) {}

// ExitNonFeatureMember is called when production nonFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitNonFeatureMember(ctx *NonFeatureMemberContext) {}

// EnterFeatureMember is called when production featureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureMember(ctx *FeatureMemberContext) {}

// ExitFeatureMember is called when production featureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureMember(ctx *FeatureMemberContext) {}

// EnterTypeFeatureMember is called when production typeFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeFeatureMember(ctx *TypeFeatureMemberContext) {}

// ExitTypeFeatureMember is called when production typeFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeFeatureMember(ctx *TypeFeatureMemberContext) {}

// EnterOwnedFeatureMember is called when production ownedFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureMember(ctx *OwnedFeatureMemberContext) {}

// ExitOwnedFeatureMember is called when production ownedFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureMember(ctx *OwnedFeatureMemberContext) {}

// EnterMemberElement is called when production memberElement is entered.
func (s *BaseSysMLv2ParserListener) EnterMemberElement(ctx *MemberElementContext) {}

// ExitMemberElement is called when production memberElement is exited.
func (s *BaseSysMLv2ParserListener) ExitMemberElement(ctx *MemberElementContext) {}

// EnterNonFeatureElement is called when production nonFeatureElement is entered.
func (s *BaseSysMLv2ParserListener) EnterNonFeatureElement(ctx *NonFeatureElementContext) {}

// ExitNonFeatureElement is called when production nonFeatureElement is exited.
func (s *BaseSysMLv2ParserListener) ExitNonFeatureElement(ctx *NonFeatureElementContext) {}

// EnterFeatureElement is called when production featureElement is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureElement(ctx *FeatureElementContext) {}

// ExitFeatureElement is called when production featureElement is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureElement(ctx *FeatureElementContext) {}

// EnterReturnFeatureMember is called when production returnFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterReturnFeatureMember(ctx *ReturnFeatureMemberContext) {}

// ExitReturnFeatureMember is called when production returnFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitReturnFeatureMember(ctx *ReturnFeatureMemberContext) {}

// EnterLiteralExpression is called when production literalExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralExpression(ctx *LiteralExpressionContext) {}

// ExitLiteralExpression is called when production literalExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralExpression(ctx *LiteralExpressionContext) {}

// EnterLiteralBoolean is called when production literalBoolean is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralBoolean(ctx *LiteralBooleanContext) {}

// ExitLiteralBoolean is called when production literalBoolean is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralBoolean(ctx *LiteralBooleanContext) {}

// EnterLiteralString is called when production literalString is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralString(ctx *LiteralStringContext) {}

// ExitLiteralString is called when production literalString is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralString(ctx *LiteralStringContext) {}

// EnterLiteralInteger is called when production literalInteger is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralInteger(ctx *LiteralIntegerContext) {}

// ExitLiteralInteger is called when production literalInteger is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralInteger(ctx *LiteralIntegerContext) {}

// EnterLiteralReal is called when production literalReal is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralReal(ctx *LiteralRealContext) {}

// ExitLiteralReal is called when production literalReal is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralReal(ctx *LiteralRealContext) {}

// EnterLiteralInfinity is called when production literalInfinity is entered.
func (s *BaseSysMLv2ParserListener) EnterLiteralInfinity(ctx *LiteralInfinityContext) {}

// ExitLiteralInfinity is called when production literalInfinity is exited.
func (s *BaseSysMLv2ParserListener) ExitLiteralInfinity(ctx *LiteralInfinityContext) {}

// EnterNamespace_ is called when production namespace_ is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespace_(ctx *Namespace_Context) {}

// ExitNamespace_ is called when production namespace_ is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespace_(ctx *Namespace_Context) {}

// EnterNamespaceDeclaration is called when production namespaceDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceDeclaration(ctx *NamespaceDeclarationContext) {}

// ExitNamespaceDeclaration is called when production namespaceDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceDeclaration(ctx *NamespaceDeclarationContext) {}

// EnterNamespaceBody is called when production namespaceBody is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceBody(ctx *NamespaceBodyContext) {}

// ExitNamespaceBody is called when production namespaceBody is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceBody(ctx *NamespaceBodyContext) {}

// EnterNamespaceBodyElement is called when production namespaceBodyElement is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceBodyElement(ctx *NamespaceBodyElementContext) {}

// ExitNamespaceBodyElement is called when production namespaceBodyElement is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceBodyElement(ctx *NamespaceBodyElementContext) {}

// EnterNamespaceMember is called when production namespaceMember is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceMember(ctx *NamespaceMemberContext) {}

// ExitNamespaceMember is called when production namespaceMember is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceMember(ctx *NamespaceMemberContext) {}

// EnterNamespaceFeatureMember is called when production namespaceFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterNamespaceFeatureMember(ctx *NamespaceFeatureMemberContext) {}

// ExitNamespaceFeatureMember is called when production namespaceFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitNamespaceFeatureMember(ctx *NamespaceFeatureMemberContext) {}

// EnterType_ is called when production type_ is entered.
func (s *BaseSysMLv2ParserListener) EnterType_(ctx *Type_Context) {}

// ExitType_ is called when production type_ is exited.
func (s *BaseSysMLv2ParserListener) ExitType_(ctx *Type_Context) {}

// EnterTypePrefix is called when production typePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterTypePrefix(ctx *TypePrefixContext) {}

// ExitTypePrefix is called when production typePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitTypePrefix(ctx *TypePrefixContext) {}

// EnterTypeDeclaration is called when production typeDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeDeclaration(ctx *TypeDeclarationContext) {}

// ExitTypeDeclaration is called when production typeDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeDeclaration(ctx *TypeDeclarationContext) {}

// EnterSpecializationPart is called when production specializationPart is entered.
func (s *BaseSysMLv2ParserListener) EnterSpecializationPart(ctx *SpecializationPartContext) {}

// ExitSpecializationPart is called when production specializationPart is exited.
func (s *BaseSysMLv2ParserListener) ExitSpecializationPart(ctx *SpecializationPartContext) {}

// EnterOwnedSpecialization is called when production ownedSpecialization is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedSpecialization(ctx *OwnedSpecializationContext) {}

// ExitOwnedSpecialization is called when production ownedSpecialization is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedSpecialization(ctx *OwnedSpecializationContext) {}

// EnterGeneralType is called when production generalType is entered.
func (s *BaseSysMLv2ParserListener) EnterGeneralType(ctx *GeneralTypeContext) {}

// ExitGeneralType is called when production generalType is exited.
func (s *BaseSysMLv2ParserListener) ExitGeneralType(ctx *GeneralTypeContext) {}

// EnterConjugationPart is called when production conjugationPart is entered.
func (s *BaseSysMLv2ParserListener) EnterConjugationPart(ctx *ConjugationPartContext) {}

// ExitConjugationPart is called when production conjugationPart is exited.
func (s *BaseSysMLv2ParserListener) ExitConjugationPart(ctx *ConjugationPartContext) {}

// EnterOwnedConjugation is called when production ownedConjugation is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedConjugation(ctx *OwnedConjugationContext) {}

// ExitOwnedConjugation is called when production ownedConjugation is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedConjugation(ctx *OwnedConjugationContext) {}

// EnterConjugates is called when production conjugates is entered.
func (s *BaseSysMLv2ParserListener) EnterConjugates(ctx *ConjugatesContext) {}

// ExitConjugates is called when production conjugates is exited.
func (s *BaseSysMLv2ParserListener) ExitConjugates(ctx *ConjugatesContext) {}

// EnterTypeRelationshipPart is called when production typeRelationshipPart is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeRelationshipPart(ctx *TypeRelationshipPartContext) {}

// ExitTypeRelationshipPart is called when production typeRelationshipPart is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeRelationshipPart(ctx *TypeRelationshipPartContext) {}

// EnterDisjoiningPart is called when production disjoiningPart is entered.
func (s *BaseSysMLv2ParserListener) EnterDisjoiningPart(ctx *DisjoiningPartContext) {}

// ExitDisjoiningPart is called when production disjoiningPart is exited.
func (s *BaseSysMLv2ParserListener) ExitDisjoiningPart(ctx *DisjoiningPartContext) {}

// EnterOwnedDisjoining is called when production ownedDisjoining is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedDisjoining(ctx *OwnedDisjoiningContext) {}

// ExitOwnedDisjoining is called when production ownedDisjoining is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedDisjoining(ctx *OwnedDisjoiningContext) {}

// EnterUnioningPart is called when production unioningPart is entered.
func (s *BaseSysMLv2ParserListener) EnterUnioningPart(ctx *UnioningPartContext) {}

// ExitUnioningPart is called when production unioningPart is exited.
func (s *BaseSysMLv2ParserListener) ExitUnioningPart(ctx *UnioningPartContext) {}

// EnterUnioning is called when production unioning is entered.
func (s *BaseSysMLv2ParserListener) EnterUnioning(ctx *UnioningContext) {}

// ExitUnioning is called when production unioning is exited.
func (s *BaseSysMLv2ParserListener) ExitUnioning(ctx *UnioningContext) {}

// EnterIntersectingPart is called when production intersectingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterIntersectingPart(ctx *IntersectingPartContext) {}

// ExitIntersectingPart is called when production intersectingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitIntersectingPart(ctx *IntersectingPartContext) {}

// EnterIntersecting is called when production intersecting is entered.
func (s *BaseSysMLv2ParserListener) EnterIntersecting(ctx *IntersectingContext) {}

// ExitIntersecting is called when production intersecting is exited.
func (s *BaseSysMLv2ParserListener) ExitIntersecting(ctx *IntersectingContext) {}

// EnterDifferencingPart is called when production differencingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterDifferencingPart(ctx *DifferencingPartContext) {}

// ExitDifferencingPart is called when production differencingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitDifferencingPart(ctx *DifferencingPartContext) {}

// EnterDifferencing is called when production differencing is entered.
func (s *BaseSysMLv2ParserListener) EnterDifferencing(ctx *DifferencingContext) {}

// ExitDifferencing is called when production differencing is exited.
func (s *BaseSysMLv2ParserListener) ExitDifferencing(ctx *DifferencingContext) {}

// EnterTypeBody is called when production typeBody is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeBody(ctx *TypeBodyContext) {}

// ExitTypeBody is called when production typeBody is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeBody(ctx *TypeBodyContext) {}

// EnterFeatureChain is called when production featureChain is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureChain(ctx *FeatureChainContext) {}

// ExitFeatureChain is called when production featureChain is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureChain(ctx *FeatureChainContext) {}

// EnterClassifier is called when production classifier is entered.
func (s *BaseSysMLv2ParserListener) EnterClassifier(ctx *ClassifierContext) {}

// ExitClassifier is called when production classifier is exited.
func (s *BaseSysMLv2ParserListener) ExitClassifier(ctx *ClassifierContext) {}

// EnterClassifierDeclaration is called when production classifierDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterClassifierDeclaration(ctx *ClassifierDeclarationContext) {}

// ExitClassifierDeclaration is called when production classifierDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitClassifierDeclaration(ctx *ClassifierDeclarationContext) {}

// EnterSuperclassingPart is called when production superclassingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterSuperclassingPart(ctx *SuperclassingPartContext) {}

// ExitSuperclassingPart is called when production superclassingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitSuperclassingPart(ctx *SuperclassingPartContext) {}

// EnterDataType is called when production dataType is entered.
func (s *BaseSysMLv2ParserListener) EnterDataType(ctx *DataTypeContext) {}

// ExitDataType is called when production dataType is exited.
func (s *BaseSysMLv2ParserListener) ExitDataType(ctx *DataTypeContext) {}

// EnterClass_ is called when production class_ is entered.
func (s *BaseSysMLv2ParserListener) EnterClass_(ctx *Class_Context) {}

// ExitClass_ is called when production class_ is exited.
func (s *BaseSysMLv2ParserListener) ExitClass_(ctx *Class_Context) {}

// EnterStructure is called when production structure is entered.
func (s *BaseSysMLv2ParserListener) EnterStructure(ctx *StructureContext) {}

// ExitStructure is called when production structure is exited.
func (s *BaseSysMLv2ParserListener) ExitStructure(ctx *StructureContext) {}

// EnterMetaclass is called when production metaclass is entered.
func (s *BaseSysMLv2ParserListener) EnterMetaclass(ctx *MetaclassContext) {}

// ExitMetaclass is called when production metaclass is exited.
func (s *BaseSysMLv2ParserListener) ExitMetaclass(ctx *MetaclassContext) {}

// EnterAssociation is called when production association is entered.
func (s *BaseSysMLv2ParserListener) EnterAssociation(ctx *AssociationContext) {}

// ExitAssociation is called when production association is exited.
func (s *BaseSysMLv2ParserListener) ExitAssociation(ctx *AssociationContext) {}

// EnterAssociationStructure is called when production associationStructure is entered.
func (s *BaseSysMLv2ParserListener) EnterAssociationStructure(ctx *AssociationStructureContext) {}

// ExitAssociationStructure is called when production associationStructure is exited.
func (s *BaseSysMLv2ParserListener) ExitAssociationStructure(ctx *AssociationStructureContext) {}

// EnterInteraction is called when production interaction is entered.
func (s *BaseSysMLv2ParserListener) EnterInteraction(ctx *InteractionContext) {}

// ExitInteraction is called when production interaction is exited.
func (s *BaseSysMLv2ParserListener) ExitInteraction(ctx *InteractionContext) {}

// EnterBehavior is called when production behavior is entered.
func (s *BaseSysMLv2ParserListener) EnterBehavior(ctx *BehaviorContext) {}

// ExitBehavior is called when production behavior is exited.
func (s *BaseSysMLv2ParserListener) ExitBehavior(ctx *BehaviorContext) {}

// EnterFunction_ is called when production function_ is entered.
func (s *BaseSysMLv2ParserListener) EnterFunction_(ctx *Function_Context) {}

// ExitFunction_ is called when production function_ is exited.
func (s *BaseSysMLv2ParserListener) ExitFunction_(ctx *Function_Context) {}

// EnterFunctionBody is called when production functionBody is entered.
func (s *BaseSysMLv2ParserListener) EnterFunctionBody(ctx *FunctionBodyContext) {}

// ExitFunctionBody is called when production functionBody is exited.
func (s *BaseSysMLv2ParserListener) ExitFunctionBody(ctx *FunctionBodyContext) {}

// EnterPredicate is called when production predicate is entered.
func (s *BaseSysMLv2ParserListener) EnterPredicate(ctx *PredicateContext) {}

// ExitPredicate is called when production predicate is exited.
func (s *BaseSysMLv2ParserListener) ExitPredicate(ctx *PredicateContext) {}

// EnterMultiplicity is called when production multiplicity is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicity(ctx *MultiplicityContext) {}

// ExitMultiplicity is called when production multiplicity is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicity(ctx *MultiplicityContext) {}

// EnterMultiplicitySubset is called when production multiplicitySubset is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicitySubset(ctx *MultiplicitySubsetContext) {}

// ExitMultiplicitySubset is called when production multiplicitySubset is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicitySubset(ctx *MultiplicitySubsetContext) {}

// EnterMultiplicityRangeDecl is called when production multiplicityRangeDecl is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicityRangeDecl(ctx *MultiplicityRangeDeclContext) {}

// ExitMultiplicityRangeDecl is called when production multiplicityRangeDecl is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicityRangeDecl(ctx *MultiplicityRangeDeclContext) {}

// EnterMultiplicityBounds is called when production multiplicityBounds is entered.
func (s *BaseSysMLv2ParserListener) EnterMultiplicityBounds(ctx *MultiplicityBoundsContext) {}

// ExitMultiplicityBounds is called when production multiplicityBounds is exited.
func (s *BaseSysMLv2ParserListener) ExitMultiplicityBounds(ctx *MultiplicityBoundsContext) {}

// EnterFeature is called when production feature is entered.
func (s *BaseSysMLv2ParserListener) EnterFeature(ctx *FeatureContext) {}

// ExitFeature is called when production feature is exited.
func (s *BaseSysMLv2ParserListener) ExitFeature(ctx *FeatureContext) {}

// EnterEndFeaturePrefix is called when production endFeaturePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterEndFeaturePrefix(ctx *EndFeaturePrefixContext) {}

// ExitEndFeaturePrefix is called when production endFeaturePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitEndFeaturePrefix(ctx *EndFeaturePrefixContext) {}

// EnterBasicFeaturePrefix is called when production basicFeaturePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterBasicFeaturePrefix(ctx *BasicFeaturePrefixContext) {}

// ExitBasicFeaturePrefix is called when production basicFeaturePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitBasicFeaturePrefix(ctx *BasicFeaturePrefixContext) {}

// EnterFeaturePrefix is called when production featurePrefix is entered.
func (s *BaseSysMLv2ParserListener) EnterFeaturePrefix(ctx *FeaturePrefixContext) {}

// ExitFeaturePrefix is called when production featurePrefix is exited.
func (s *BaseSysMLv2ParserListener) ExitFeaturePrefix(ctx *FeaturePrefixContext) {}

// EnterFeatureDeclaration is called when production featureDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureDeclaration(ctx *FeatureDeclarationContext) {}

// ExitFeatureDeclaration is called when production featureDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureDeclaration(ctx *FeatureDeclarationContext) {}

// EnterFeatureIdentification is called when production featureIdentification is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureIdentification(ctx *FeatureIdentificationContext) {}

// ExitFeatureIdentification is called when production featureIdentification is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureIdentification(ctx *FeatureIdentificationContext) {}

// EnterFeatureRelationshipPart is called when production featureRelationshipPart is entered.
func (s *BaseSysMLv2ParserListener) EnterFeatureRelationshipPart(ctx *FeatureRelationshipPartContext) {
}

// ExitFeatureRelationshipPart is called when production featureRelationshipPart is exited.
func (s *BaseSysMLv2ParserListener) ExitFeatureRelationshipPart(ctx *FeatureRelationshipPartContext) {
}

// EnterChainingPart is called when production chainingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterChainingPart(ctx *ChainingPartContext) {}

// ExitChainingPart is called when production chainingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitChainingPart(ctx *ChainingPartContext) {}

// EnterInvertingPart is called when production invertingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterInvertingPart(ctx *InvertingPartContext) {}

// ExitInvertingPart is called when production invertingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitInvertingPart(ctx *InvertingPartContext) {}

// EnterOwnedFeatureInverting is called when production ownedFeatureInverting is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedFeatureInverting(ctx *OwnedFeatureInvertingContext) {}

// ExitOwnedFeatureInverting is called when production ownedFeatureInverting is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedFeatureInverting(ctx *OwnedFeatureInvertingContext) {}

// EnterTypeFeaturingPart is called when production typeFeaturingPart is entered.
func (s *BaseSysMLv2ParserListener) EnterTypeFeaturingPart(ctx *TypeFeaturingPartContext) {}

// ExitTypeFeaturingPart is called when production typeFeaturingPart is exited.
func (s *BaseSysMLv2ParserListener) ExitTypeFeaturingPart(ctx *TypeFeaturingPartContext) {}

// EnterOwnedTypeFeaturing is called when production ownedTypeFeaturing is entered.
func (s *BaseSysMLv2ParserListener) EnterOwnedTypeFeaturing(ctx *OwnedTypeFeaturingContext) {}

// ExitOwnedTypeFeaturing is called when production ownedTypeFeaturing is exited.
func (s *BaseSysMLv2ParserListener) ExitOwnedTypeFeaturing(ctx *OwnedTypeFeaturingContext) {}

// EnterStep is called when production step is entered.
func (s *BaseSysMLv2ParserListener) EnterStep(ctx *StepContext) {}

// ExitStep is called when production step is exited.
func (s *BaseSysMLv2ParserListener) ExitStep(ctx *StepContext) {}

// EnterExpression_ is called when production expression_ is entered.
func (s *BaseSysMLv2ParserListener) EnterExpression_(ctx *Expression_Context) {}

// ExitExpression_ is called when production expression_ is exited.
func (s *BaseSysMLv2ParserListener) ExitExpression_(ctx *Expression_Context) {}

// EnterBooleanExpression is called when production booleanExpression is entered.
func (s *BaseSysMLv2ParserListener) EnterBooleanExpression(ctx *BooleanExpressionContext) {}

// ExitBooleanExpression is called when production booleanExpression is exited.
func (s *BaseSysMLv2ParserListener) ExitBooleanExpression(ctx *BooleanExpressionContext) {}

// EnterInvariant is called when production invariant is entered.
func (s *BaseSysMLv2ParserListener) EnterInvariant(ctx *InvariantContext) {}

// ExitInvariant is called when production invariant is exited.
func (s *BaseSysMLv2ParserListener) ExitInvariant(ctx *InvariantContext) {}

// EnterConnector is called when production connector is entered.
func (s *BaseSysMLv2ParserListener) EnterConnector(ctx *ConnectorContext) {}

// ExitConnector is called when production connector is exited.
func (s *BaseSysMLv2ParserListener) ExitConnector(ctx *ConnectorContext) {}

// EnterConnectorDeclaration is called when production connectorDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterConnectorDeclaration(ctx *ConnectorDeclarationContext) {}

// ExitConnectorDeclaration is called when production connectorDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitConnectorDeclaration(ctx *ConnectorDeclarationContext) {}

// EnterBinaryConnectorDeclaration is called when production binaryConnectorDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterBinaryConnectorDeclaration(ctx *BinaryConnectorDeclarationContext) {
}

// ExitBinaryConnectorDeclaration is called when production binaryConnectorDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitBinaryConnectorDeclaration(ctx *BinaryConnectorDeclarationContext) {
}

// EnterNaryConnectorDeclaration is called when production naryConnectorDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterNaryConnectorDeclaration(ctx *NaryConnectorDeclarationContext) {
}

// ExitNaryConnectorDeclaration is called when production naryConnectorDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitNaryConnectorDeclaration(ctx *NaryConnectorDeclarationContext) {
}

// EnterBindingConnector is called when production bindingConnector is entered.
func (s *BaseSysMLv2ParserListener) EnterBindingConnector(ctx *BindingConnectorContext) {}

// ExitBindingConnector is called when production bindingConnector is exited.
func (s *BaseSysMLv2ParserListener) ExitBindingConnector(ctx *BindingConnectorContext) {}

// EnterBindingConnectorDeclaration is called when production bindingConnectorDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterBindingConnectorDeclaration(ctx *BindingConnectorDeclarationContext) {
}

// ExitBindingConnectorDeclaration is called when production bindingConnectorDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitBindingConnectorDeclaration(ctx *BindingConnectorDeclarationContext) {
}

// EnterSuccession is called when production succession is entered.
func (s *BaseSysMLv2ParserListener) EnterSuccession(ctx *SuccessionContext) {}

// ExitSuccession is called when production succession is exited.
func (s *BaseSysMLv2ParserListener) ExitSuccession(ctx *SuccessionContext) {}

// EnterSuccessionDeclaration is called when production successionDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterSuccessionDeclaration(ctx *SuccessionDeclarationContext) {}

// ExitSuccessionDeclaration is called when production successionDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitSuccessionDeclaration(ctx *SuccessionDeclarationContext) {}

// EnterKermlFlow is called when production kermlFlow is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlow(ctx *KermlFlowContext) {}

// ExitKermlFlow is called when production kermlFlow is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlow(ctx *KermlFlowContext) {}

// EnterKermlSuccessionFlow is called when production kermlSuccessionFlow is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlSuccessionFlow(ctx *KermlSuccessionFlowContext) {}

// ExitKermlSuccessionFlow is called when production kermlSuccessionFlow is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlSuccessionFlow(ctx *KermlSuccessionFlowContext) {}

// EnterKermlFlowDeclaration is called when production kermlFlowDeclaration is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowDeclaration(ctx *KermlFlowDeclarationContext) {}

// ExitKermlFlowDeclaration is called when production kermlFlowDeclaration is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowDeclaration(ctx *KermlFlowDeclarationContext) {}

// EnterKermlPayloadFeatureMember is called when production kermlPayloadFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlPayloadFeatureMember(ctx *KermlPayloadFeatureMemberContext) {
}

// ExitKermlPayloadFeatureMember is called when production kermlPayloadFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlPayloadFeatureMember(ctx *KermlPayloadFeatureMemberContext) {
}

// EnterKermlPayloadFeature is called when production kermlPayloadFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlPayloadFeature(ctx *KermlPayloadFeatureContext) {}

// ExitKermlPayloadFeature is called when production kermlPayloadFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlPayloadFeature(ctx *KermlPayloadFeatureContext) {}

// EnterKermlFlowEndMember is called when production kermlFlowEndMember is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowEndMember(ctx *KermlFlowEndMemberContext) {}

// ExitKermlFlowEndMember is called when production kermlFlowEndMember is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowEndMember(ctx *KermlFlowEndMemberContext) {}

// EnterKermlFlowEnd is called when production kermlFlowEnd is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowEnd(ctx *KermlFlowEndContext) {}

// ExitKermlFlowEnd is called when production kermlFlowEnd is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowEnd(ctx *KermlFlowEndContext) {}

// EnterKermlFlowFeatureMember is called when production kermlFlowFeatureMember is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowFeatureMember(ctx *KermlFlowFeatureMemberContext) {}

// ExitKermlFlowFeatureMember is called when production kermlFlowFeatureMember is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowFeatureMember(ctx *KermlFlowFeatureMemberContext) {}

// EnterKermlFlowFeature is called when production kermlFlowFeature is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowFeature(ctx *KermlFlowFeatureContext) {}

// ExitKermlFlowFeature is called when production kermlFlowFeature is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowFeature(ctx *KermlFlowFeatureContext) {}

// EnterKermlFlowFeatureRedefinition is called when production kermlFlowFeatureRedefinition is entered.
func (s *BaseSysMLv2ParserListener) EnterKermlFlowFeatureRedefinition(ctx *KermlFlowFeatureRedefinitionContext) {
}

// ExitKermlFlowFeatureRedefinition is called when production kermlFlowFeatureRedefinition is exited.
func (s *BaseSysMLv2ParserListener) ExitKermlFlowFeatureRedefinition(ctx *KermlFlowFeatureRedefinitionContext) {
}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BaseSysMLv2ParserListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BaseSysMLv2ParserListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterName is called when production name is entered.
func (s *BaseSysMLv2ParserListener) EnterName(ctx *NameContext) {}

// ExitName is called when production name is exited.
func (s *BaseSysMLv2ParserListener) ExitName(ctx *NameContext) {}
