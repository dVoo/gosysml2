// Code generated from antlr/SysMLv2Parser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // SysMLv2Parser
import "github.com/antlr4-go/antlr/v4"

// SysMLv2ParserListener is a complete listener for a parse tree produced by SysMLv2Parser.
type SysMLv2ParserListener interface {
	antlr.ParseTreeListener

	// EnterEntryRuleRootNamespace is called when entering the entryRuleRootNamespace production.
	EnterEntryRuleRootNamespace(c *EntryRuleRootNamespaceContext)

	// EnterRootNamespace is called when entering the rootNamespace production.
	EnterRootNamespace(c *RootNamespaceContext)

	// EnterIdentification is called when entering the identification production.
	EnterIdentification(c *IdentificationContext)

	// EnterRelationshipBody is called when entering the relationshipBody production.
	EnterRelationshipBody(c *RelationshipBodyContext)

	// EnterDependency is called when entering the dependency production.
	EnterDependency(c *DependencyContext)

	// EnterDependencyDeclaration is called when entering the dependencyDeclaration production.
	EnterDependencyDeclaration(c *DependencyDeclarationContext)

	// EnterAnnotation is called when entering the annotation production.
	EnterAnnotation(c *AnnotationContext)

	// EnterOwnedAnnotation is called when entering the ownedAnnotation production.
	EnterOwnedAnnotation(c *OwnedAnnotationContext)

	// EnterAnnotatingMember is called when entering the annotatingMember production.
	EnterAnnotatingMember(c *AnnotatingMemberContext)

	// EnterAnnotatingElement is called when entering the annotatingElement production.
	EnterAnnotatingElement(c *AnnotatingElementContext)

	// EnterComment is called when entering the comment production.
	EnterComment(c *CommentContext)

	// EnterDocumentation is called when entering the documentation production.
	EnterDocumentation(c *DocumentationContext)

	// EnterTextualRepresentation is called when entering the textualRepresentation production.
	EnterTextualRepresentation(c *TextualRepresentationContext)

	// EnterPackage is called when entering the package production.
	EnterPackage(c *PackageContext)

	// EnterLibraryPackage is called when entering the libraryPackage production.
	EnterLibraryPackage(c *LibraryPackageContext)

	// EnterPackageDeclaration is called when entering the packageDeclaration production.
	EnterPackageDeclaration(c *PackageDeclarationContext)

	// EnterPackageBody is called when entering the packageBody production.
	EnterPackageBody(c *PackageBodyContext)

	// EnterPackageBodyElement is called when entering the packageBodyElement production.
	EnterPackageBodyElement(c *PackageBodyElementContext)

	// EnterMemberPrefix is called when entering the memberPrefix production.
	EnterMemberPrefix(c *MemberPrefixContext)

	// EnterPackageMember is called when entering the packageMember production.
	EnterPackageMember(c *PackageMemberContext)

	// EnterElementFilterMember is called when entering the elementFilterMember production.
	EnterElementFilterMember(c *ElementFilterMemberContext)

	// EnterAliasMember is called when entering the aliasMember production.
	EnterAliasMember(c *AliasMemberContext)

	// EnterImport_ is called when entering the import_ production.
	EnterImport_(c *Import_Context)

	// EnterImportDeclaration is called when entering the importDeclaration production.
	EnterImportDeclaration(c *ImportDeclarationContext)

	// EnterMembershipImport is called when entering the membershipImport production.
	EnterMembershipImport(c *MembershipImportContext)

	// EnterNamespaceImport is called when entering the namespaceImport production.
	EnterNamespaceImport(c *NamespaceImportContext)

	// EnterFilterPackage is called when entering the filterPackage production.
	EnterFilterPackage(c *FilterPackageContext)

	// EnterFilterPackageImportPart is called when entering the filterPackageImportPart production.
	EnterFilterPackageImportPart(c *FilterPackageImportPartContext)

	// EnterFilterPackageMember is called when entering the filterPackageMember production.
	EnterFilterPackageMember(c *FilterPackageMemberContext)

	// EnterVisibilityIndicator is called when entering the visibilityIndicator production.
	EnterVisibilityIndicator(c *VisibilityIndicatorContext)

	// EnterDefinitionElement is called when entering the definitionElement production.
	EnterDefinitionElement(c *DefinitionElementContext)

	// EnterUsageElement is called when entering the usageElement production.
	EnterUsageElement(c *UsageElementContext)

	// EnterBasicDefinitionPrefix is called when entering the basicDefinitionPrefix production.
	EnterBasicDefinitionPrefix(c *BasicDefinitionPrefixContext)

	// EnterDefinitionExtensionKeyword is called when entering the definitionExtensionKeyword production.
	EnterDefinitionExtensionKeyword(c *DefinitionExtensionKeywordContext)

	// EnterDefinitionPrefix is called when entering the definitionPrefix production.
	EnterDefinitionPrefix(c *DefinitionPrefixContext)

	// EnterDefinition is called when entering the definition production.
	EnterDefinition(c *DefinitionContext)

	// EnterDefinitionDeclaration is called when entering the definitionDeclaration production.
	EnterDefinitionDeclaration(c *DefinitionDeclarationContext)

	// EnterDefinitionBody is called when entering the definitionBody production.
	EnterDefinitionBody(c *DefinitionBodyContext)

	// EnterDefinitionBodyItem is called when entering the definitionBodyItem production.
	EnterDefinitionBodyItem(c *DefinitionBodyItemContext)

	// EnterEndFeatureMember is called when entering the endFeatureMember production.
	EnterEndFeatureMember(c *EndFeatureMemberContext)

	// EnterEndFeatureDeclaration is called when entering the endFeatureDeclaration production.
	EnterEndFeatureDeclaration(c *EndFeatureDeclarationContext)

	// EnterDefinitionMember is called when entering the definitionMember production.
	EnterDefinitionMember(c *DefinitionMemberContext)

	// EnterVariantUsageMember is called when entering the variantUsageMember production.
	EnterVariantUsageMember(c *VariantUsageMemberContext)

	// EnterNonOccurrenceUsageMember is called when entering the nonOccurrenceUsageMember production.
	EnterNonOccurrenceUsageMember(c *NonOccurrenceUsageMemberContext)

	// EnterOccurrenceUsageMember is called when entering the occurrenceUsageMember production.
	EnterOccurrenceUsageMember(c *OccurrenceUsageMemberContext)

	// EnterStructureUsageMember is called when entering the structureUsageMember production.
	EnterStructureUsageMember(c *StructureUsageMemberContext)

	// EnterBehaviorUsageMember is called when entering the behaviorUsageMember production.
	EnterBehaviorUsageMember(c *BehaviorUsageMemberContext)

	// EnterFeatureDirection is called when entering the featureDirection production.
	EnterFeatureDirection(c *FeatureDirectionContext)

	// EnterRefPrefix is called when entering the refPrefix production.
	EnterRefPrefix(c *RefPrefixContext)

	// EnterBasicUsagePrefix is called when entering the basicUsagePrefix production.
	EnterBasicUsagePrefix(c *BasicUsagePrefixContext)

	// EnterEndUsagePrefix is called when entering the endUsagePrefix production.
	EnterEndUsagePrefix(c *EndUsagePrefixContext)

	// EnterOwnedCrossFeatureMember is called when entering the ownedCrossFeatureMember production.
	EnterOwnedCrossFeatureMember(c *OwnedCrossFeatureMemberContext)

	// EnterOwnedCrossFeature is called when entering the ownedCrossFeature production.
	EnterOwnedCrossFeature(c *OwnedCrossFeatureContext)

	// EnterUsageExtensionKeyword is called when entering the usageExtensionKeyword production.
	EnterUsageExtensionKeyword(c *UsageExtensionKeywordContext)

	// EnterUnextendedUsagePrefix is called when entering the unextendedUsagePrefix production.
	EnterUnextendedUsagePrefix(c *UnextendedUsagePrefixContext)

	// EnterUsagePrefix is called when entering the usagePrefix production.
	EnterUsagePrefix(c *UsagePrefixContext)

	// EnterUsage is called when entering the usage production.
	EnterUsage(c *UsageContext)

	// EnterUsageDeclaration is called when entering the usageDeclaration production.
	EnterUsageDeclaration(c *UsageDeclarationContext)

	// EnterUsageCompletion is called when entering the usageCompletion production.
	EnterUsageCompletion(c *UsageCompletionContext)

	// EnterUsageBody is called when entering the usageBody production.
	EnterUsageBody(c *UsageBodyContext)

	// EnterValuePart is called when entering the valuePart production.
	EnterValuePart(c *ValuePartContext)

	// EnterFeatureValue is called when entering the featureValue production.
	EnterFeatureValue(c *FeatureValueContext)

	// EnterDefaultReferenceUsage is called when entering the defaultReferenceUsage production.
	EnterDefaultReferenceUsage(c *DefaultReferenceUsageContext)

	// EnterReferenceUsage is called when entering the referenceUsage production.
	EnterReferenceUsage(c *ReferenceUsageContext)

	// EnterVariantReference is called when entering the variantReference production.
	EnterVariantReference(c *VariantReferenceContext)

	// EnterNonOccurrenceUsageElement is called when entering the nonOccurrenceUsageElement production.
	EnterNonOccurrenceUsageElement(c *NonOccurrenceUsageElementContext)

	// EnterOccurrenceUsageElement is called when entering the occurrenceUsageElement production.
	EnterOccurrenceUsageElement(c *OccurrenceUsageElementContext)

	// EnterStructureUsageElement is called when entering the structureUsageElement production.
	EnterStructureUsageElement(c *StructureUsageElementContext)

	// EnterBehaviorUsageElement is called when entering the behaviorUsageElement production.
	EnterBehaviorUsageElement(c *BehaviorUsageElementContext)

	// EnterVariantUsageElement is called when entering the variantUsageElement production.
	EnterVariantUsageElement(c *VariantUsageElementContext)

	// EnterSubclassificationPart is called when entering the subclassificationPart production.
	EnterSubclassificationPart(c *SubclassificationPartContext)

	// EnterOwnedSubclassification is called when entering the ownedSubclassification production.
	EnterOwnedSubclassification(c *OwnedSubclassificationContext)

	// EnterFeatureSpecializationPart is called when entering the featureSpecializationPart production.
	EnterFeatureSpecializationPart(c *FeatureSpecializationPartContext)

	// EnterFeatureSpecialization is called when entering the featureSpecialization production.
	EnterFeatureSpecialization(c *FeatureSpecializationContext)

	// EnterTypings is called when entering the typings production.
	EnterTypings(c *TypingsContext)

	// EnterTypedBy is called when entering the typedBy production.
	EnterTypedBy(c *TypedByContext)

	// EnterOwnedFeatureTyping is called when entering the ownedFeatureTyping production.
	EnterOwnedFeatureTyping(c *OwnedFeatureTypingContext)

	// EnterSubsettings is called when entering the subsettings production.
	EnterSubsettings(c *SubsettingsContext)

	// EnterOwnedSubsetting is called when entering the ownedSubsetting production.
	EnterOwnedSubsetting(c *OwnedSubsettingContext)

	// EnterReferences is called when entering the references production.
	EnterReferences(c *ReferencesContext)

	// EnterOwnedReferenceSubsetting is called when entering the ownedReferenceSubsetting production.
	EnterOwnedReferenceSubsetting(c *OwnedReferenceSubsettingContext)

	// EnterCrosses is called when entering the crosses production.
	EnterCrosses(c *CrossesContext)

	// EnterOwnedCrossSubsetting is called when entering the ownedCrossSubsetting production.
	EnterOwnedCrossSubsetting(c *OwnedCrossSubsettingContext)

	// EnterRedefinitions is called when entering the redefinitions production.
	EnterRedefinitions(c *RedefinitionsContext)

	// EnterOwnedRedefinition is called when entering the ownedRedefinition production.
	EnterOwnedRedefinition(c *OwnedRedefinitionContext)

	// EnterOwnedFeatureChain is called when entering the ownedFeatureChain production.
	EnterOwnedFeatureChain(c *OwnedFeatureChainContext)

	// EnterOwnedFeatureChaining is called when entering the ownedFeatureChaining production.
	EnterOwnedFeatureChaining(c *OwnedFeatureChainingContext)

	// EnterSpecializes is called when entering the specializes production.
	EnterSpecializes(c *SpecializesContext)

	// EnterDefinedBy is called when entering the definedBy production.
	EnterDefinedBy(c *DefinedByContext)

	// EnterSubsetsKw is called when entering the subsetsKw production.
	EnterSubsetsKw(c *SubsetsKwContext)

	// EnterReferencesKw is called when entering the referencesKw production.
	EnterReferencesKw(c *ReferencesKwContext)

	// EnterCrossesKw is called when entering the crossesKw production.
	EnterCrossesKw(c *CrossesKwContext)

	// EnterRedefinesKw is called when entering the redefinesKw production.
	EnterRedefinesKw(c *RedefinesKwContext)

	// EnterMultiplicityPart is called when entering the multiplicityPart production.
	EnterMultiplicityPart(c *MultiplicityPartContext)

	// EnterOwnedMultiplicity is called when entering the ownedMultiplicity production.
	EnterOwnedMultiplicity(c *OwnedMultiplicityContext)

	// EnterMultiplicityRange is called when entering the multiplicityRange production.
	EnterMultiplicityRange(c *MultiplicityRangeContext)

	// EnterMultiplicityExpressionMember is called when entering the multiplicityExpressionMember production.
	EnterMultiplicityExpressionMember(c *MultiplicityExpressionMemberContext)

	// EnterAttributeDefinition is called when entering the attributeDefinition production.
	EnterAttributeDefinition(c *AttributeDefinitionContext)

	// EnterAttributeUsage is called when entering the attributeUsage production.
	EnterAttributeUsage(c *AttributeUsageContext)

	// EnterEnumerationDefinition is called when entering the enumerationDefinition production.
	EnterEnumerationDefinition(c *EnumerationDefinitionContext)

	// EnterEnumerationBody is called when entering the enumerationBody production.
	EnterEnumerationBody(c *EnumerationBodyContext)

	// EnterEnumerationUsageMember is called when entering the enumerationUsageMember production.
	EnterEnumerationUsageMember(c *EnumerationUsageMemberContext)

	// EnterEnumeratedValue is called when entering the enumeratedValue production.
	EnterEnumeratedValue(c *EnumeratedValueContext)

	// EnterEnumerationUsage is called when entering the enumerationUsage production.
	EnterEnumerationUsage(c *EnumerationUsageContext)

	// EnterOccurrenceDefinitionPrefix is called when entering the occurrenceDefinitionPrefix production.
	EnterOccurrenceDefinitionPrefix(c *OccurrenceDefinitionPrefixContext)

	// EnterOccurrenceDefinition is called when entering the occurrenceDefinition production.
	EnterOccurrenceDefinition(c *OccurrenceDefinitionContext)

	// EnterIndividualDefinition is called when entering the individualDefinition production.
	EnterIndividualDefinition(c *IndividualDefinitionContext)

	// EnterOccurrenceUsagePrefix is called when entering the occurrenceUsagePrefix production.
	EnterOccurrenceUsagePrefix(c *OccurrenceUsagePrefixContext)

	// EnterOccurrenceUsage is called when entering the occurrenceUsage production.
	EnterOccurrenceUsage(c *OccurrenceUsageContext)

	// EnterIndividualUsage is called when entering the individualUsage production.
	EnterIndividualUsage(c *IndividualUsageContext)

	// EnterPortionUsage is called when entering the portionUsage production.
	EnterPortionUsage(c *PortionUsageContext)

	// EnterPortionKind is called when entering the portionKind production.
	EnterPortionKind(c *PortionKindContext)

	// EnterEventOccurrenceUsage is called when entering the eventOccurrenceUsage production.
	EnterEventOccurrenceUsage(c *EventOccurrenceUsageContext)

	// EnterSourceSuccessionMember is called when entering the sourceSuccessionMember production.
	EnterSourceSuccessionMember(c *SourceSuccessionMemberContext)

	// EnterSourceSuccession is called when entering the sourceSuccession production.
	EnterSourceSuccession(c *SourceSuccessionContext)

	// EnterSourceEndMember is called when entering the sourceEndMember production.
	EnterSourceEndMember(c *SourceEndMemberContext)

	// EnterSourceEnd is called when entering the sourceEnd production.
	EnterSourceEnd(c *SourceEndContext)

	// EnterItemDefinition is called when entering the itemDefinition production.
	EnterItemDefinition(c *ItemDefinitionContext)

	// EnterItemUsage is called when entering the itemUsage production.
	EnterItemUsage(c *ItemUsageContext)

	// EnterPartDefinition is called when entering the partDefinition production.
	EnterPartDefinition(c *PartDefinitionContext)

	// EnterPartUsage is called when entering the partUsage production.
	EnterPartUsage(c *PartUsageContext)

	// EnterPortDefinition is called when entering the portDefinition production.
	EnterPortDefinition(c *PortDefinitionContext)

	// EnterPortUsage is called when entering the portUsage production.
	EnterPortUsage(c *PortUsageContext)

	// EnterConjugatedPortTyping is called when entering the conjugatedPortTyping production.
	EnterConjugatedPortTyping(c *ConjugatedPortTypingContext)

	// EnterConnectionDefinition is called when entering the connectionDefinition production.
	EnterConnectionDefinition(c *ConnectionDefinitionContext)

	// EnterConnectionUsage is called when entering the connectionUsage production.
	EnterConnectionUsage(c *ConnectionUsageContext)

	// EnterConnectorPart is called when entering the connectorPart production.
	EnterConnectorPart(c *ConnectorPartContext)

	// EnterBinaryConnectorPart is called when entering the binaryConnectorPart production.
	EnterBinaryConnectorPart(c *BinaryConnectorPartContext)

	// EnterNaryConnectorPart is called when entering the naryConnectorPart production.
	EnterNaryConnectorPart(c *NaryConnectorPartContext)

	// EnterConnectorEndMember is called when entering the connectorEndMember production.
	EnterConnectorEndMember(c *ConnectorEndMemberContext)

	// EnterConnectorEnd is called when entering the connectorEnd production.
	EnterConnectorEnd(c *ConnectorEndContext)

	// EnterOwnedCrossMultiplicityMember is called when entering the ownedCrossMultiplicityMember production.
	EnterOwnedCrossMultiplicityMember(c *OwnedCrossMultiplicityMemberContext)

	// EnterOwnedCrossMultiplicity is called when entering the ownedCrossMultiplicity production.
	EnterOwnedCrossMultiplicity(c *OwnedCrossMultiplicityContext)

	// EnterBindingConnectorAsUsage is called when entering the bindingConnectorAsUsage production.
	EnterBindingConnectorAsUsage(c *BindingConnectorAsUsageContext)

	// EnterSuccessionAsUsage is called when entering the successionAsUsage production.
	EnterSuccessionAsUsage(c *SuccessionAsUsageContext)

	// EnterInterfaceDefinition is called when entering the interfaceDefinition production.
	EnterInterfaceDefinition(c *InterfaceDefinitionContext)

	// EnterInterfaceBody is called when entering the interfaceBody production.
	EnterInterfaceBody(c *InterfaceBodyContext)

	// EnterInterfaceBodyItem is called when entering the interfaceBodyItem production.
	EnterInterfaceBodyItem(c *InterfaceBodyItemContext)

	// EnterInterfaceNonOccurrenceUsageMember is called when entering the interfaceNonOccurrenceUsageMember production.
	EnterInterfaceNonOccurrenceUsageMember(c *InterfaceNonOccurrenceUsageMemberContext)

	// EnterInterfaceNonOccurrenceUsageElement is called when entering the interfaceNonOccurrenceUsageElement production.
	EnterInterfaceNonOccurrenceUsageElement(c *InterfaceNonOccurrenceUsageElementContext)

	// EnterInterfaceOccurrenceUsageMember is called when entering the interfaceOccurrenceUsageMember production.
	EnterInterfaceOccurrenceUsageMember(c *InterfaceOccurrenceUsageMemberContext)

	// EnterInterfaceOccurrenceUsageElement is called when entering the interfaceOccurrenceUsageElement production.
	EnterInterfaceOccurrenceUsageElement(c *InterfaceOccurrenceUsageElementContext)

	// EnterDefaultInterfaceEnd is called when entering the defaultInterfaceEnd production.
	EnterDefaultInterfaceEnd(c *DefaultInterfaceEndContext)

	// EnterInterfaceUsage is called when entering the interfaceUsage production.
	EnterInterfaceUsage(c *InterfaceUsageContext)

	// EnterInterfaceUsageDeclaration is called when entering the interfaceUsageDeclaration production.
	EnterInterfaceUsageDeclaration(c *InterfaceUsageDeclarationContext)

	// EnterInterfacePart is called when entering the interfacePart production.
	EnterInterfacePart(c *InterfacePartContext)

	// EnterBinaryInterfacePart is called when entering the binaryInterfacePart production.
	EnterBinaryInterfacePart(c *BinaryInterfacePartContext)

	// EnterNaryInterfacePart is called when entering the naryInterfacePart production.
	EnterNaryInterfacePart(c *NaryInterfacePartContext)

	// EnterInterfaceEndMember is called when entering the interfaceEndMember production.
	EnterInterfaceEndMember(c *InterfaceEndMemberContext)

	// EnterInterfaceEnd is called when entering the interfaceEnd production.
	EnterInterfaceEnd(c *InterfaceEndContext)

	// EnterAllocationDefinition is called when entering the allocationDefinition production.
	EnterAllocationDefinition(c *AllocationDefinitionContext)

	// EnterAllocationUsage is called when entering the allocationUsage production.
	EnterAllocationUsage(c *AllocationUsageContext)

	// EnterAllocationUsageDeclaration is called when entering the allocationUsageDeclaration production.
	EnterAllocationUsageDeclaration(c *AllocationUsageDeclarationContext)

	// EnterFlowDefinition is called when entering the flowDefinition production.
	EnterFlowDefinition(c *FlowDefinitionContext)

	// EnterMessage is called when entering the message production.
	EnterMessage(c *MessageContext)

	// EnterMessageDeclaration is called when entering the messageDeclaration production.
	EnterMessageDeclaration(c *MessageDeclarationContext)

	// EnterMessageEventMember is called when entering the messageEventMember production.
	EnterMessageEventMember(c *MessageEventMemberContext)

	// EnterMessageEvent is called when entering the messageEvent production.
	EnterMessageEvent(c *MessageEventContext)

	// EnterFlowUsage is called when entering the flowUsage production.
	EnterFlowUsage(c *FlowUsageContext)

	// EnterSuccessionFlowUsage is called when entering the successionFlowUsage production.
	EnterSuccessionFlowUsage(c *SuccessionFlowUsageContext)

	// EnterFlowDeclaration is called when entering the flowDeclaration production.
	EnterFlowDeclaration(c *FlowDeclarationContext)

	// EnterFlowPayloadFeatureMember is called when entering the flowPayloadFeatureMember production.
	EnterFlowPayloadFeatureMember(c *FlowPayloadFeatureMemberContext)

	// EnterFlowPayloadFeature is called when entering the flowPayloadFeature production.
	EnterFlowPayloadFeature(c *FlowPayloadFeatureContext)

	// EnterPayloadFeature is called when entering the payloadFeature production.
	EnterPayloadFeature(c *PayloadFeatureContext)

	// EnterPayloadFeatureSpecializationPart is called when entering the payloadFeatureSpecializationPart production.
	EnterPayloadFeatureSpecializationPart(c *PayloadFeatureSpecializationPartContext)

	// EnterFlowEndMember is called when entering the flowEndMember production.
	EnterFlowEndMember(c *FlowEndMemberContext)

	// EnterFlowEnd is called when entering the flowEnd production.
	EnterFlowEnd(c *FlowEndContext)

	// EnterFlowEndSubsetting is called when entering the flowEndSubsetting production.
	EnterFlowEndSubsetting(c *FlowEndSubsettingContext)

	// EnterFeatureChainPrefix is called when entering the featureChainPrefix production.
	EnterFeatureChainPrefix(c *FeatureChainPrefixContext)

	// EnterFlowFeatureMember is called when entering the flowFeatureMember production.
	EnterFlowFeatureMember(c *FlowFeatureMemberContext)

	// EnterFlowFeature is called when entering the flowFeature production.
	EnterFlowFeature(c *FlowFeatureContext)

	// EnterFlowFeatureRedefinition is called when entering the flowFeatureRedefinition production.
	EnterFlowFeatureRedefinition(c *FlowFeatureRedefinitionContext)

	// EnterActionDefinition is called when entering the actionDefinition production.
	EnterActionDefinition(c *ActionDefinitionContext)

	// EnterActionBody is called when entering the actionBody production.
	EnterActionBody(c *ActionBodyContext)

	// EnterActionBodyItem is called when entering the actionBodyItem production.
	EnterActionBodyItem(c *ActionBodyItemContext)

	// EnterNonBehaviorBodyItem is called when entering the nonBehaviorBodyItem production.
	EnterNonBehaviorBodyItem(c *NonBehaviorBodyItemContext)

	// EnterActionBehaviorMember is called when entering the actionBehaviorMember production.
	EnterActionBehaviorMember(c *ActionBehaviorMemberContext)

	// EnterInitialNodeMember is called when entering the initialNodeMember production.
	EnterInitialNodeMember(c *InitialNodeMemberContext)

	// EnterActionNodeMember is called when entering the actionNodeMember production.
	EnterActionNodeMember(c *ActionNodeMemberContext)

	// EnterActionTargetSuccessionMember is called when entering the actionTargetSuccessionMember production.
	EnterActionTargetSuccessionMember(c *ActionTargetSuccessionMemberContext)

	// EnterGuardedSuccessionMember is called when entering the guardedSuccessionMember production.
	EnterGuardedSuccessionMember(c *GuardedSuccessionMemberContext)

	// EnterActionUsage is called when entering the actionUsage production.
	EnterActionUsage(c *ActionUsageContext)

	// EnterActionUsageDeclaration is called when entering the actionUsageDeclaration production.
	EnterActionUsageDeclaration(c *ActionUsageDeclarationContext)

	// EnterPerformActionUsage is called when entering the performActionUsage production.
	EnterPerformActionUsage(c *PerformActionUsageContext)

	// EnterPerformActionUsageDeclaration is called when entering the performActionUsageDeclaration production.
	EnterPerformActionUsageDeclaration(c *PerformActionUsageDeclarationContext)

	// EnterActionNode is called when entering the actionNode production.
	EnterActionNode(c *ActionNodeContext)

	// EnterActionNodeUsageDeclaration is called when entering the actionNodeUsageDeclaration production.
	EnterActionNodeUsageDeclaration(c *ActionNodeUsageDeclarationContext)

	// EnterActionNodePrefix is called when entering the actionNodePrefix production.
	EnterActionNodePrefix(c *ActionNodePrefixContext)

	// EnterControlNode is called when entering the controlNode production.
	EnterControlNode(c *ControlNodeContext)

	// EnterControlNodePrefix is called when entering the controlNodePrefix production.
	EnterControlNodePrefix(c *ControlNodePrefixContext)

	// EnterMergeNode is called when entering the mergeNode production.
	EnterMergeNode(c *MergeNodeContext)

	// EnterDecisionNode is called when entering the decisionNode production.
	EnterDecisionNode(c *DecisionNodeContext)

	// EnterJoinNode is called when entering the joinNode production.
	EnterJoinNode(c *JoinNodeContext)

	// EnterForkNode is called when entering the forkNode production.
	EnterForkNode(c *ForkNodeContext)

	// EnterAcceptNode is called when entering the acceptNode production.
	EnterAcceptNode(c *AcceptNodeContext)

	// EnterAcceptNodeDeclaration is called when entering the acceptNodeDeclaration production.
	EnterAcceptNodeDeclaration(c *AcceptNodeDeclarationContext)

	// EnterAcceptParameterPart is called when entering the acceptParameterPart production.
	EnterAcceptParameterPart(c *AcceptParameterPartContext)

	// EnterPayloadParameterMember is called when entering the payloadParameterMember production.
	EnterPayloadParameterMember(c *PayloadParameterMemberContext)

	// EnterPayloadParameter is called when entering the payloadParameter production.
	EnterPayloadParameter(c *PayloadParameterContext)

	// EnterTriggerValuePart is called when entering the triggerValuePart production.
	EnterTriggerValuePart(c *TriggerValuePartContext)

	// EnterTriggerFeatureValue is called when entering the triggerFeatureValue production.
	EnterTriggerFeatureValue(c *TriggerFeatureValueContext)

	// EnterTriggerExpression is called when entering the triggerExpression production.
	EnterTriggerExpression(c *TriggerExpressionContext)

	// EnterSendNode is called when entering the sendNode production.
	EnterSendNode(c *SendNodeContext)

	// EnterSendNodeDeclaration is called when entering the sendNodeDeclaration production.
	EnterSendNodeDeclaration(c *SendNodeDeclarationContext)

	// EnterSenderReceiverPart is called when entering the senderReceiverPart production.
	EnterSenderReceiverPart(c *SenderReceiverPartContext)

	// EnterNodeParameterMember is called when entering the nodeParameterMember production.
	EnterNodeParameterMember(c *NodeParameterMemberContext)

	// EnterNodeParameter is called when entering the nodeParameter production.
	EnterNodeParameter(c *NodeParameterContext)

	// EnterAssignmentNode is called when entering the assignmentNode production.
	EnterAssignmentNode(c *AssignmentNodeContext)

	// EnterAssignmentNodeDeclaration is called when entering the assignmentNodeDeclaration production.
	EnterAssignmentNodeDeclaration(c *AssignmentNodeDeclarationContext)

	// EnterAssignmentTargetMember is called when entering the assignmentTargetMember production.
	EnterAssignmentTargetMember(c *AssignmentTargetMemberContext)

	// EnterAssignmentTargetParameter is called when entering the assignmentTargetParameter production.
	EnterAssignmentTargetParameter(c *AssignmentTargetParameterContext)

	// EnterFeatureChainMember is called when entering the featureChainMember production.
	EnterFeatureChainMember(c *FeatureChainMemberContext)

	// EnterOwnedFeatureChainMember is called when entering the ownedFeatureChainMember production.
	EnterOwnedFeatureChainMember(c *OwnedFeatureChainMemberContext)

	// EnterTerminateNode is called when entering the terminateNode production.
	EnterTerminateNode(c *TerminateNodeContext)

	// EnterIfNode is called when entering the ifNode production.
	EnterIfNode(c *IfNodeContext)

	// EnterActionBodyParameter is called when entering the actionBodyParameter production.
	EnterActionBodyParameter(c *ActionBodyParameterContext)

	// EnterWhileLoopNode is called when entering the whileLoopNode production.
	EnterWhileLoopNode(c *WhileLoopNodeContext)

	// EnterForLoopNode is called when entering the forLoopNode production.
	EnterForLoopNode(c *ForLoopNodeContext)

	// EnterForVariableDeclaration is called when entering the forVariableDeclaration production.
	EnterForVariableDeclaration(c *ForVariableDeclarationContext)

	// EnterActionTargetSuccession is called when entering the actionTargetSuccession production.
	EnterActionTargetSuccession(c *ActionTargetSuccessionContext)

	// EnterTargetSuccession is called when entering the targetSuccession production.
	EnterTargetSuccession(c *TargetSuccessionContext)

	// EnterGuardedTargetSuccession is called when entering the guardedTargetSuccession production.
	EnterGuardedTargetSuccession(c *GuardedTargetSuccessionContext)

	// EnterDefaultTargetSuccession is called when entering the defaultTargetSuccession production.
	EnterDefaultTargetSuccession(c *DefaultTargetSuccessionContext)

	// EnterGuardedSuccession is called when entering the guardedSuccession production.
	EnterGuardedSuccession(c *GuardedSuccessionContext)

	// EnterStateDefinition is called when entering the stateDefinition production.
	EnterStateDefinition(c *StateDefinitionContext)

	// EnterStateDefBody is called when entering the stateDefBody production.
	EnterStateDefBody(c *StateDefBodyContext)

	// EnterStateBodyItem is called when entering the stateBodyItem production.
	EnterStateBodyItem(c *StateBodyItemContext)

	// EnterEntryActionMember is called when entering the entryActionMember production.
	EnterEntryActionMember(c *EntryActionMemberContext)

	// EnterDoActionMember is called when entering the doActionMember production.
	EnterDoActionMember(c *DoActionMemberContext)

	// EnterExitActionMember is called when entering the exitActionMember production.
	EnterExitActionMember(c *ExitActionMemberContext)

	// EnterEntryTransitionMember is called when entering the entryTransitionMember production.
	EnterEntryTransitionMember(c *EntryTransitionMemberContext)

	// EnterStateActionUsage is called when entering the stateActionUsage production.
	EnterStateActionUsage(c *StateActionUsageContext)

	// EnterStatePerformActionUsage is called when entering the statePerformActionUsage production.
	EnterStatePerformActionUsage(c *StatePerformActionUsageContext)

	// EnterStateAcceptActionUsage is called when entering the stateAcceptActionUsage production.
	EnterStateAcceptActionUsage(c *StateAcceptActionUsageContext)

	// EnterStateSendActionUsage is called when entering the stateSendActionUsage production.
	EnterStateSendActionUsage(c *StateSendActionUsageContext)

	// EnterStateAssignmentActionUsage is called when entering the stateAssignmentActionUsage production.
	EnterStateAssignmentActionUsage(c *StateAssignmentActionUsageContext)

	// EnterTransitionUsageMember is called when entering the transitionUsageMember production.
	EnterTransitionUsageMember(c *TransitionUsageMemberContext)

	// EnterTargetTransitionUsageMember is called when entering the targetTransitionUsageMember production.
	EnterTargetTransitionUsageMember(c *TargetTransitionUsageMemberContext)

	// EnterStateUsage is called when entering the stateUsage production.
	EnterStateUsage(c *StateUsageContext)

	// EnterStateUsageBody is called when entering the stateUsageBody production.
	EnterStateUsageBody(c *StateUsageBodyContext)

	// EnterExhibitStateUsage is called when entering the exhibitStateUsage production.
	EnterExhibitStateUsage(c *ExhibitStateUsageContext)

	// EnterTransitionUsage is called when entering the transitionUsage production.
	EnterTransitionUsage(c *TransitionUsageContext)

	// EnterTargetTransitionUsage is called when entering the targetTransitionUsage production.
	EnterTargetTransitionUsage(c *TargetTransitionUsageContext)

	// EnterTriggerActionMember is called when entering the triggerActionMember production.
	EnterTriggerActionMember(c *TriggerActionMemberContext)

	// EnterTriggerAction is called when entering the triggerAction production.
	EnterTriggerAction(c *TriggerActionContext)

	// EnterGuardExpressionMember is called when entering the guardExpressionMember production.
	EnterGuardExpressionMember(c *GuardExpressionMemberContext)

	// EnterEffectBehaviorMember is called when entering the effectBehaviorMember production.
	EnterEffectBehaviorMember(c *EffectBehaviorMemberContext)

	// EnterEffectBehaviorUsage is called when entering the effectBehaviorUsage production.
	EnterEffectBehaviorUsage(c *EffectBehaviorUsageContext)

	// EnterTransitionPerformActionUsage is called when entering the transitionPerformActionUsage production.
	EnterTransitionPerformActionUsage(c *TransitionPerformActionUsageContext)

	// EnterTransitionAcceptActionUsage is called when entering the transitionAcceptActionUsage production.
	EnterTransitionAcceptActionUsage(c *TransitionAcceptActionUsageContext)

	// EnterTransitionSendActionUsage is called when entering the transitionSendActionUsage production.
	EnterTransitionSendActionUsage(c *TransitionSendActionUsageContext)

	// EnterTransitionAssignmentActionUsage is called when entering the transitionAssignmentActionUsage production.
	EnterTransitionAssignmentActionUsage(c *TransitionAssignmentActionUsageContext)

	// EnterTransitionSuccessionMember is called when entering the transitionSuccessionMember production.
	EnterTransitionSuccessionMember(c *TransitionSuccessionMemberContext)

	// EnterTransitionSuccession is called when entering the transitionSuccession production.
	EnterTransitionSuccession(c *TransitionSuccessionContext)

	// EnterCalculationDefinition is called when entering the calculationDefinition production.
	EnterCalculationDefinition(c *CalculationDefinitionContext)

	// EnterCalculationUsage is called when entering the calculationUsage production.
	EnterCalculationUsage(c *CalculationUsageContext)

	// EnterCalculationBody is called when entering the calculationBody production.
	EnterCalculationBody(c *CalculationBodyContext)

	// EnterCalculationBodyPart is called when entering the calculationBodyPart production.
	EnterCalculationBodyPart(c *CalculationBodyPartContext)

	// EnterCalculationBodyItem is called when entering the calculationBodyItem production.
	EnterCalculationBodyItem(c *CalculationBodyItemContext)

	// EnterReturnParameterMember is called when entering the returnParameterMember production.
	EnterReturnParameterMember(c *ReturnParameterMemberContext)

	// EnterResultExpressionMember is called when entering the resultExpressionMember production.
	EnterResultExpressionMember(c *ResultExpressionMemberContext)

	// EnterConstraintDefinition is called when entering the constraintDefinition production.
	EnterConstraintDefinition(c *ConstraintDefinitionContext)

	// EnterConstraintUsage is called when entering the constraintUsage production.
	EnterConstraintUsage(c *ConstraintUsageContext)

	// EnterAssertConstraintUsage is called when entering the assertConstraintUsage production.
	EnterAssertConstraintUsage(c *AssertConstraintUsageContext)

	// EnterConstraintUsageDeclaration is called when entering the constraintUsageDeclaration production.
	EnterConstraintUsageDeclaration(c *ConstraintUsageDeclarationContext)

	// EnterRequirementDefinition is called when entering the requirementDefinition production.
	EnterRequirementDefinition(c *RequirementDefinitionContext)

	// EnterRequirementBody is called when entering the requirementBody production.
	EnterRequirementBody(c *RequirementBodyContext)

	// EnterRequirementBodyItem is called when entering the requirementBodyItem production.
	EnterRequirementBodyItem(c *RequirementBodyItemContext)

	// EnterSubjectMember is called when entering the subjectMember production.
	EnterSubjectMember(c *SubjectMemberContext)

	// EnterSubjectUsage is called when entering the subjectUsage production.
	EnterSubjectUsage(c *SubjectUsageContext)

	// EnterRequirementConstraintMember is called when entering the requirementConstraintMember production.
	EnterRequirementConstraintMember(c *RequirementConstraintMemberContext)

	// EnterRequirementKind is called when entering the requirementKind production.
	EnterRequirementKind(c *RequirementKindContext)

	// EnterRequirementConstraintUsage is called when entering the requirementConstraintUsage production.
	EnterRequirementConstraintUsage(c *RequirementConstraintUsageContext)

	// EnterFramedConcernMember is called when entering the framedConcernMember production.
	EnterFramedConcernMember(c *FramedConcernMemberContext)

	// EnterFramedConcernUsage is called when entering the framedConcernUsage production.
	EnterFramedConcernUsage(c *FramedConcernUsageContext)

	// EnterCalculationUsageDeclaration is called when entering the calculationUsageDeclaration production.
	EnterCalculationUsageDeclaration(c *CalculationUsageDeclarationContext)

	// EnterActorMember is called when entering the actorMember production.
	EnterActorMember(c *ActorMemberContext)

	// EnterActorUsage is called when entering the actorUsage production.
	EnterActorUsage(c *ActorUsageContext)

	// EnterStakeholderMember is called when entering the stakeholderMember production.
	EnterStakeholderMember(c *StakeholderMemberContext)

	// EnterStakeholderUsage is called when entering the stakeholderUsage production.
	EnterStakeholderUsage(c *StakeholderUsageContext)

	// EnterRequirementUsage is called when entering the requirementUsage production.
	EnterRequirementUsage(c *RequirementUsageContext)

	// EnterSatisfyRequirementUsage is called when entering the satisfyRequirementUsage production.
	EnterSatisfyRequirementUsage(c *SatisfyRequirementUsageContext)

	// EnterSatisfactionSubjectMember is called when entering the satisfactionSubjectMember production.
	EnterSatisfactionSubjectMember(c *SatisfactionSubjectMemberContext)

	// EnterSatisfactionParameter is called when entering the satisfactionParameter production.
	EnterSatisfactionParameter(c *SatisfactionParameterContext)

	// EnterConcernDefinition is called when entering the concernDefinition production.
	EnterConcernDefinition(c *ConcernDefinitionContext)

	// EnterConcernUsage is called when entering the concernUsage production.
	EnterConcernUsage(c *ConcernUsageContext)

	// EnterCaseDefinition is called when entering the caseDefinition production.
	EnterCaseDefinition(c *CaseDefinitionContext)

	// EnterCaseUsage is called when entering the caseUsage production.
	EnterCaseUsage(c *CaseUsageContext)

	// EnterCaseBody is called when entering the caseBody production.
	EnterCaseBody(c *CaseBodyContext)

	// EnterCaseBodyItem is called when entering the caseBodyItem production.
	EnterCaseBodyItem(c *CaseBodyItemContext)

	// EnterObjectiveMember is called when entering the objectiveMember production.
	EnterObjectiveMember(c *ObjectiveMemberContext)

	// EnterObjectiveRequirementUsage is called when entering the objectiveRequirementUsage production.
	EnterObjectiveRequirementUsage(c *ObjectiveRequirementUsageContext)

	// EnterAnalysisCaseDefinition is called when entering the analysisCaseDefinition production.
	EnterAnalysisCaseDefinition(c *AnalysisCaseDefinitionContext)

	// EnterAnalysisCaseUsage is called when entering the analysisCaseUsage production.
	EnterAnalysisCaseUsage(c *AnalysisCaseUsageContext)

	// EnterVerificationCaseDefinition is called when entering the verificationCaseDefinition production.
	EnterVerificationCaseDefinition(c *VerificationCaseDefinitionContext)

	// EnterVerificationCaseUsage is called when entering the verificationCaseUsage production.
	EnterVerificationCaseUsage(c *VerificationCaseUsageContext)

	// EnterRequirementVerificationMember is called when entering the requirementVerificationMember production.
	EnterRequirementVerificationMember(c *RequirementVerificationMemberContext)

	// EnterRequirementVerificationUsage is called when entering the requirementVerificationUsage production.
	EnterRequirementVerificationUsage(c *RequirementVerificationUsageContext)

	// EnterUseCaseDefinition is called when entering the useCaseDefinition production.
	EnterUseCaseDefinition(c *UseCaseDefinitionContext)

	// EnterUseCaseUsage is called when entering the useCaseUsage production.
	EnterUseCaseUsage(c *UseCaseUsageContext)

	// EnterIncludeUseCaseUsage is called when entering the includeUseCaseUsage production.
	EnterIncludeUseCaseUsage(c *IncludeUseCaseUsageContext)

	// EnterViewDefinition is called when entering the viewDefinition production.
	EnterViewDefinition(c *ViewDefinitionContext)

	// EnterViewDefinitionBody is called when entering the viewDefinitionBody production.
	EnterViewDefinitionBody(c *ViewDefinitionBodyContext)

	// EnterViewDefinitionBodyItem is called when entering the viewDefinitionBodyItem production.
	EnterViewDefinitionBodyItem(c *ViewDefinitionBodyItemContext)

	// EnterViewRenderingMember is called when entering the viewRenderingMember production.
	EnterViewRenderingMember(c *ViewRenderingMemberContext)

	// EnterViewRenderingUsage is called when entering the viewRenderingUsage production.
	EnterViewRenderingUsage(c *ViewRenderingUsageContext)

	// EnterViewUsage is called when entering the viewUsage production.
	EnterViewUsage(c *ViewUsageContext)

	// EnterViewBody is called when entering the viewBody production.
	EnterViewBody(c *ViewBodyContext)

	// EnterViewBodyItem is called when entering the viewBodyItem production.
	EnterViewBodyItem(c *ViewBodyItemContext)

	// EnterExpose is called when entering the expose production.
	EnterExpose(c *ExposeContext)

	// EnterMembershipExpose is called when entering the membershipExpose production.
	EnterMembershipExpose(c *MembershipExposeContext)

	// EnterNamespaceExpose is called when entering the namespaceExpose production.
	EnterNamespaceExpose(c *NamespaceExposeContext)

	// EnterViewpointDefinition is called when entering the viewpointDefinition production.
	EnterViewpointDefinition(c *ViewpointDefinitionContext)

	// EnterViewpointUsage is called when entering the viewpointUsage production.
	EnterViewpointUsage(c *ViewpointUsageContext)

	// EnterRenderingDefinition is called when entering the renderingDefinition production.
	EnterRenderingDefinition(c *RenderingDefinitionContext)

	// EnterRenderingUsage is called when entering the renderingUsage production.
	EnterRenderingUsage(c *RenderingUsageContext)

	// EnterMetadataDefinition is called when entering the metadataDefinition production.
	EnterMetadataDefinition(c *MetadataDefinitionContext)

	// EnterPrefixMetadataAnnotation is called when entering the prefixMetadataAnnotation production.
	EnterPrefixMetadataAnnotation(c *PrefixMetadataAnnotationContext)

	// EnterPrefixMetadataMember is called when entering the prefixMetadataMember production.
	EnterPrefixMetadataMember(c *PrefixMetadataMemberContext)

	// EnterPrefixMetadataUsage is called when entering the prefixMetadataUsage production.
	EnterPrefixMetadataUsage(c *PrefixMetadataUsageContext)

	// EnterMetadataUsage is called when entering the metadataUsage production.
	EnterMetadataUsage(c *MetadataUsageContext)

	// EnterMetadataUsageDeclaration is called when entering the metadataUsageDeclaration production.
	EnterMetadataUsageDeclaration(c *MetadataUsageDeclarationContext)

	// EnterMetadataBody is called when entering the metadataBody production.
	EnterMetadataBody(c *MetadataBodyContext)

	// EnterMetadataBodyUsageMember is called when entering the metadataBodyUsageMember production.
	EnterMetadataBodyUsageMember(c *MetadataBodyUsageMemberContext)

	// EnterMetadataBodyUsage is called when entering the metadataBodyUsage production.
	EnterMetadataBodyUsage(c *MetadataBodyUsageContext)

	// EnterMetadataFeature is called when entering the metadataFeature production.
	EnterMetadataFeature(c *MetadataFeatureContext)

	// EnterExtendedDefinition is called when entering the extendedDefinition production.
	EnterExtendedDefinition(c *ExtendedDefinitionContext)

	// EnterExtendedUsage is called when entering the extendedUsage production.
	EnterExtendedUsage(c *ExtendedUsageContext)

	// EnterOwnedExpression is called when entering the ownedExpression production.
	EnterOwnedExpression(c *OwnedExpressionContext)

	// EnterConditionalBinaryOperator is called when entering the conditionalBinaryOperator production.
	EnterConditionalBinaryOperator(c *ConditionalBinaryOperatorContext)

	// EnterBinaryOperator is called when entering the binaryOperator production.
	EnterBinaryOperator(c *BinaryOperatorContext)

	// EnterUnaryOperator is called when entering the unaryOperator production.
	EnterUnaryOperator(c *UnaryOperatorContext)

	// EnterClassificationTestOperator is called when entering the classificationTestOperator production.
	EnterClassificationTestOperator(c *ClassificationTestOperatorContext)

	// EnterCastOperator is called when entering the castOperator production.
	EnterCastOperator(c *CastOperatorContext)

	// EnterMetaCastOperator is called when entering the metaCastOperator production.
	EnterMetaCastOperator(c *MetaCastOperatorContext)

	// EnterTypeReference is called when entering the typeReference production.
	EnterTypeReference(c *TypeReferenceContext)

	// EnterPrimaryExpression is called when entering the primaryExpression production.
	EnterPrimaryExpression(c *PrimaryExpressionContext)

	// EnterPrimaryExpressionSuffix is called when entering the primaryExpressionSuffix production.
	EnterPrimaryExpressionSuffix(c *PrimaryExpressionSuffixContext)

	// EnterBaseExpression is called when entering the baseExpression production.
	EnterBaseExpression(c *BaseExpressionContext)

	// EnterNullExpression is called when entering the nullExpression production.
	EnterNullExpression(c *NullExpressionContext)

	// EnterFeatureReferenceExpression is called when entering the featureReferenceExpression production.
	EnterFeatureReferenceExpression(c *FeatureReferenceExpressionContext)

	// EnterMetadataAccessExpression is called when entering the metadataAccessExpression production.
	EnterMetadataAccessExpression(c *MetadataAccessExpressionContext)

	// EnterInvocationExpression is called when entering the invocationExpression production.
	EnterInvocationExpression(c *InvocationExpressionContext)

	// EnterConstructorExpression is called when entering the constructorExpression production.
	EnterConstructorExpression(c *ConstructorExpressionContext)

	// EnterBodyExpression is called when entering the bodyExpression production.
	EnterBodyExpression(c *BodyExpressionContext)

	// EnterSequenceExpression is called when entering the sequenceExpression production.
	EnterSequenceExpression(c *SequenceExpressionContext)

	// EnterSequenceExpressionList is called when entering the sequenceExpressionList production.
	EnterSequenceExpressionList(c *SequenceExpressionListContext)

	// EnterArgumentList is called when entering the argumentList production.
	EnterArgumentList(c *ArgumentListContext)

	// EnterPositionalArgumentList is called when entering the positionalArgumentList production.
	EnterPositionalArgumentList(c *PositionalArgumentListContext)

	// EnterNamedArgumentList is called when entering the namedArgumentList production.
	EnterNamedArgumentList(c *NamedArgumentListContext)

	// EnterNamedArgument is called when entering the namedArgument production.
	EnterNamedArgument(c *NamedArgumentContext)

	// EnterFunctionBodyPart is called when entering the functionBodyPart production.
	EnterFunctionBodyPart(c *FunctionBodyPartContext)

	// EnterResultExpressionMemberOpt is called when entering the resultExpressionMemberOpt production.
	EnterResultExpressionMemberOpt(c *ResultExpressionMemberOptContext)

	// EnterTypeBodyElement is called when entering the typeBodyElement production.
	EnterTypeBodyElement(c *TypeBodyElementContext)

	// EnterNonFeatureMember is called when entering the nonFeatureMember production.
	EnterNonFeatureMember(c *NonFeatureMemberContext)

	// EnterFeatureMember is called when entering the featureMember production.
	EnterFeatureMember(c *FeatureMemberContext)

	// EnterTypeFeatureMember is called when entering the typeFeatureMember production.
	EnterTypeFeatureMember(c *TypeFeatureMemberContext)

	// EnterOwnedFeatureMember is called when entering the ownedFeatureMember production.
	EnterOwnedFeatureMember(c *OwnedFeatureMemberContext)

	// EnterMemberElement is called when entering the memberElement production.
	EnterMemberElement(c *MemberElementContext)

	// EnterNonFeatureElement is called when entering the nonFeatureElement production.
	EnterNonFeatureElement(c *NonFeatureElementContext)

	// EnterFeatureElement is called when entering the featureElement production.
	EnterFeatureElement(c *FeatureElementContext)

	// EnterReturnFeatureMember is called when entering the returnFeatureMember production.
	EnterReturnFeatureMember(c *ReturnFeatureMemberContext)

	// EnterLiteralExpression is called when entering the literalExpression production.
	EnterLiteralExpression(c *LiteralExpressionContext)

	// EnterLiteralBoolean is called when entering the literalBoolean production.
	EnterLiteralBoolean(c *LiteralBooleanContext)

	// EnterLiteralString is called when entering the literalString production.
	EnterLiteralString(c *LiteralStringContext)

	// EnterLiteralInteger is called when entering the literalInteger production.
	EnterLiteralInteger(c *LiteralIntegerContext)

	// EnterLiteralReal is called when entering the literalReal production.
	EnterLiteralReal(c *LiteralRealContext)

	// EnterLiteralInfinity is called when entering the literalInfinity production.
	EnterLiteralInfinity(c *LiteralInfinityContext)

	// EnterNamespace_ is called when entering the namespace_ production.
	EnterNamespace_(c *Namespace_Context)

	// EnterNamespaceDeclaration is called when entering the namespaceDeclaration production.
	EnterNamespaceDeclaration(c *NamespaceDeclarationContext)

	// EnterNamespaceBody is called when entering the namespaceBody production.
	EnterNamespaceBody(c *NamespaceBodyContext)

	// EnterNamespaceBodyElement is called when entering the namespaceBodyElement production.
	EnterNamespaceBodyElement(c *NamespaceBodyElementContext)

	// EnterNamespaceMember is called when entering the namespaceMember production.
	EnterNamespaceMember(c *NamespaceMemberContext)

	// EnterNamespaceFeatureMember is called when entering the namespaceFeatureMember production.
	EnterNamespaceFeatureMember(c *NamespaceFeatureMemberContext)

	// EnterType_ is called when entering the type_ production.
	EnterType_(c *Type_Context)

	// EnterTypePrefix is called when entering the typePrefix production.
	EnterTypePrefix(c *TypePrefixContext)

	// EnterTypeDeclaration is called when entering the typeDeclaration production.
	EnterTypeDeclaration(c *TypeDeclarationContext)

	// EnterSpecializationPart is called when entering the specializationPart production.
	EnterSpecializationPart(c *SpecializationPartContext)

	// EnterOwnedSpecialization is called when entering the ownedSpecialization production.
	EnterOwnedSpecialization(c *OwnedSpecializationContext)

	// EnterGeneralType is called when entering the generalType production.
	EnterGeneralType(c *GeneralTypeContext)

	// EnterConjugationPart is called when entering the conjugationPart production.
	EnterConjugationPart(c *ConjugationPartContext)

	// EnterOwnedConjugation is called when entering the ownedConjugation production.
	EnterOwnedConjugation(c *OwnedConjugationContext)

	// EnterConjugates is called when entering the conjugates production.
	EnterConjugates(c *ConjugatesContext)

	// EnterTypeRelationshipPart is called when entering the typeRelationshipPart production.
	EnterTypeRelationshipPart(c *TypeRelationshipPartContext)

	// EnterDisjoiningPart is called when entering the disjoiningPart production.
	EnterDisjoiningPart(c *DisjoiningPartContext)

	// EnterDisjoining is called when entering the disjoining production.
	EnterDisjoining(c *DisjoiningContext)

	// EnterOwnedDisjoining is called when entering the ownedDisjoining production.
	EnterOwnedDisjoining(c *OwnedDisjoiningContext)

	// EnterUnioningPart is called when entering the unioningPart production.
	EnterUnioningPart(c *UnioningPartContext)

	// EnterUnioning is called when entering the unioning production.
	EnterUnioning(c *UnioningContext)

	// EnterIntersectingPart is called when entering the intersectingPart production.
	EnterIntersectingPart(c *IntersectingPartContext)

	// EnterIntersecting is called when entering the intersecting production.
	EnterIntersecting(c *IntersectingContext)

	// EnterDifferencingPart is called when entering the differencingPart production.
	EnterDifferencingPart(c *DifferencingPartContext)

	// EnterDifferencing is called when entering the differencing production.
	EnterDifferencing(c *DifferencingContext)

	// EnterTypeBody is called when entering the typeBody production.
	EnterTypeBody(c *TypeBodyContext)

	// EnterFeatureChain is called when entering the featureChain production.
	EnterFeatureChain(c *FeatureChainContext)

	// EnterClassifier is called when entering the classifier production.
	EnterClassifier(c *ClassifierContext)

	// EnterSubclassifier is called when entering the subclassifier production.
	EnterSubclassifier(c *SubclassifierContext)

	// EnterClassifierDeclaration is called when entering the classifierDeclaration production.
	EnterClassifierDeclaration(c *ClassifierDeclarationContext)

	// EnterSuperclassingPart is called when entering the superclassingPart production.
	EnterSuperclassingPart(c *SuperclassingPartContext)

	// EnterDataType is called when entering the dataType production.
	EnterDataType(c *DataTypeContext)

	// EnterClass is called when entering the class production.
	EnterClass(c *ClassContext)

	// EnterStructure is called when entering the structure production.
	EnterStructure(c *StructureContext)

	// EnterMetaclass is called when entering the metaclass production.
	EnterMetaclass(c *MetaclassContext)

	// EnterAssociation is called when entering the association production.
	EnterAssociation(c *AssociationContext)

	// EnterAssociationStructure is called when entering the associationStructure production.
	EnterAssociationStructure(c *AssociationStructureContext)

	// EnterInteraction is called when entering the interaction production.
	EnterInteraction(c *InteractionContext)

	// EnterBehavior is called when entering the behavior production.
	EnterBehavior(c *BehaviorContext)

	// EnterFunction_ is called when entering the function_ production.
	EnterFunction_(c *Function_Context)

	// EnterFunctionBody is called when entering the functionBody production.
	EnterFunctionBody(c *FunctionBodyContext)

	// EnterPredicate is called when entering the predicate production.
	EnterPredicate(c *PredicateContext)

	// EnterMultiplicity is called when entering the multiplicity production.
	EnterMultiplicity(c *MultiplicityContext)

	// EnterMultiplicitySubset is called when entering the multiplicitySubset production.
	EnterMultiplicitySubset(c *MultiplicitySubsetContext)

	// EnterMultiplicityRangeDecl is called when entering the multiplicityRangeDecl production.
	EnterMultiplicityRangeDecl(c *MultiplicityRangeDeclContext)

	// EnterFeatureSubsetting is called when entering the featureSubsetting production.
	EnterFeatureSubsetting(c *FeatureSubsettingContext)

	// EnterMultiplicityBounds is called when entering the multiplicityBounds production.
	EnterMultiplicityBounds(c *MultiplicityBoundsContext)

	// EnterFeature is called when entering the feature production.
	EnterFeature(c *FeatureContext)

	// EnterEndFeaturePrefix is called when entering the endFeaturePrefix production.
	EnterEndFeaturePrefix(c *EndFeaturePrefixContext)

	// EnterBasicFeaturePrefix is called when entering the basicFeaturePrefix production.
	EnterBasicFeaturePrefix(c *BasicFeaturePrefixContext)

	// EnterFeaturePrefix is called when entering the featurePrefix production.
	EnterFeaturePrefix(c *FeaturePrefixContext)

	// EnterFeatureDeclaration is called when entering the featureDeclaration production.
	EnterFeatureDeclaration(c *FeatureDeclarationContext)

	// EnterFeatureIdentification is called when entering the featureIdentification production.
	EnterFeatureIdentification(c *FeatureIdentificationContext)

	// EnterFeatureRelationshipPart is called when entering the featureRelationshipPart production.
	EnterFeatureRelationshipPart(c *FeatureRelationshipPartContext)

	// EnterChainingPart is called when entering the chainingPart production.
	EnterChainingPart(c *ChainingPartContext)

	// EnterInvertingPart is called when entering the invertingPart production.
	EnterInvertingPart(c *InvertingPartContext)

	// EnterOwnedFeatureInverting is called when entering the ownedFeatureInverting production.
	EnterOwnedFeatureInverting(c *OwnedFeatureInvertingContext)

	// EnterTypeFeaturingPart is called when entering the typeFeaturingPart production.
	EnterTypeFeaturingPart(c *TypeFeaturingPartContext)

	// EnterOwnedTypeFeaturing is called when entering the ownedTypeFeaturing production.
	EnterOwnedTypeFeaturing(c *OwnedTypeFeaturingContext)

	// EnterStep is called when entering the step production.
	EnterStep(c *StepContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterBooleanExpression is called when entering the booleanExpression production.
	EnterBooleanExpression(c *BooleanExpressionContext)

	// EnterInvariant is called when entering the invariant production.
	EnterInvariant(c *InvariantContext)

	// EnterConnector is called when entering the connector production.
	EnterConnector(c *ConnectorContext)

	// EnterConnectorDeclaration is called when entering the connectorDeclaration production.
	EnterConnectorDeclaration(c *ConnectorDeclarationContext)

	// EnterBinaryConnectorDeclaration is called when entering the binaryConnectorDeclaration production.
	EnterBinaryConnectorDeclaration(c *BinaryConnectorDeclarationContext)

	// EnterNaryConnectorDeclaration is called when entering the naryConnectorDeclaration production.
	EnterNaryConnectorDeclaration(c *NaryConnectorDeclarationContext)

	// EnterBindingConnector is called when entering the bindingConnector production.
	EnterBindingConnector(c *BindingConnectorContext)

	// EnterBindingConnectorDeclaration is called when entering the bindingConnectorDeclaration production.
	EnterBindingConnectorDeclaration(c *BindingConnectorDeclarationContext)

	// EnterSuccession is called when entering the succession production.
	EnterSuccession(c *SuccessionContext)

	// EnterSuccessionDeclaration is called when entering the successionDeclaration production.
	EnterSuccessionDeclaration(c *SuccessionDeclarationContext)

	// EnterKermlFlow is called when entering the kermlFlow production.
	EnterKermlFlow(c *KermlFlowContext)

	// EnterKermlSuccessionFlow is called when entering the kermlSuccessionFlow production.
	EnterKermlSuccessionFlow(c *KermlSuccessionFlowContext)

	// EnterKermlFlowDeclaration is called when entering the kermlFlowDeclaration production.
	EnterKermlFlowDeclaration(c *KermlFlowDeclarationContext)

	// EnterKermlPayloadFeatureMember is called when entering the kermlPayloadFeatureMember production.
	EnterKermlPayloadFeatureMember(c *KermlPayloadFeatureMemberContext)

	// EnterKermlPayloadFeature is called when entering the kermlPayloadFeature production.
	EnterKermlPayloadFeature(c *KermlPayloadFeatureContext)

	// EnterKermlFlowEndMember is called when entering the kermlFlowEndMember production.
	EnterKermlFlowEndMember(c *KermlFlowEndMemberContext)

	// EnterKermlFlowEnd is called when entering the kermlFlowEnd production.
	EnterKermlFlowEnd(c *KermlFlowEndContext)

	// EnterKermlFlowFeatureMember is called when entering the kermlFlowFeatureMember production.
	EnterKermlFlowFeatureMember(c *KermlFlowFeatureMemberContext)

	// EnterKermlFlowFeature is called when entering the kermlFlowFeature production.
	EnterKermlFlowFeature(c *KermlFlowFeatureContext)

	// EnterKermlFlowFeatureRedefinition is called when entering the kermlFlowFeatureRedefinition production.
	EnterKermlFlowFeatureRedefinition(c *KermlFlowFeatureRedefinitionContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterName is called when entering the name production.
	EnterName(c *NameContext)

	// EnterKeywordName is called when entering the keywordName production.
	EnterKeywordName(c *KeywordNameContext)

	// ExitEntryRuleRootNamespace is called when exiting the entryRuleRootNamespace production.
	ExitEntryRuleRootNamespace(c *EntryRuleRootNamespaceContext)

	// ExitRootNamespace is called when exiting the rootNamespace production.
	ExitRootNamespace(c *RootNamespaceContext)

	// ExitIdentification is called when exiting the identification production.
	ExitIdentification(c *IdentificationContext)

	// ExitRelationshipBody is called when exiting the relationshipBody production.
	ExitRelationshipBody(c *RelationshipBodyContext)

	// ExitDependency is called when exiting the dependency production.
	ExitDependency(c *DependencyContext)

	// ExitDependencyDeclaration is called when exiting the dependencyDeclaration production.
	ExitDependencyDeclaration(c *DependencyDeclarationContext)

	// ExitAnnotation is called when exiting the annotation production.
	ExitAnnotation(c *AnnotationContext)

	// ExitOwnedAnnotation is called when exiting the ownedAnnotation production.
	ExitOwnedAnnotation(c *OwnedAnnotationContext)

	// ExitAnnotatingMember is called when exiting the annotatingMember production.
	ExitAnnotatingMember(c *AnnotatingMemberContext)

	// ExitAnnotatingElement is called when exiting the annotatingElement production.
	ExitAnnotatingElement(c *AnnotatingElementContext)

	// ExitComment is called when exiting the comment production.
	ExitComment(c *CommentContext)

	// ExitDocumentation is called when exiting the documentation production.
	ExitDocumentation(c *DocumentationContext)

	// ExitTextualRepresentation is called when exiting the textualRepresentation production.
	ExitTextualRepresentation(c *TextualRepresentationContext)

	// ExitPackage is called when exiting the package production.
	ExitPackage(c *PackageContext)

	// ExitLibraryPackage is called when exiting the libraryPackage production.
	ExitLibraryPackage(c *LibraryPackageContext)

	// ExitPackageDeclaration is called when exiting the packageDeclaration production.
	ExitPackageDeclaration(c *PackageDeclarationContext)

	// ExitPackageBody is called when exiting the packageBody production.
	ExitPackageBody(c *PackageBodyContext)

	// ExitPackageBodyElement is called when exiting the packageBodyElement production.
	ExitPackageBodyElement(c *PackageBodyElementContext)

	// ExitMemberPrefix is called when exiting the memberPrefix production.
	ExitMemberPrefix(c *MemberPrefixContext)

	// ExitPackageMember is called when exiting the packageMember production.
	ExitPackageMember(c *PackageMemberContext)

	// ExitElementFilterMember is called when exiting the elementFilterMember production.
	ExitElementFilterMember(c *ElementFilterMemberContext)

	// ExitAliasMember is called when exiting the aliasMember production.
	ExitAliasMember(c *AliasMemberContext)

	// ExitImport_ is called when exiting the import_ production.
	ExitImport_(c *Import_Context)

	// ExitImportDeclaration is called when exiting the importDeclaration production.
	ExitImportDeclaration(c *ImportDeclarationContext)

	// ExitMembershipImport is called when exiting the membershipImport production.
	ExitMembershipImport(c *MembershipImportContext)

	// ExitNamespaceImport is called when exiting the namespaceImport production.
	ExitNamespaceImport(c *NamespaceImportContext)

	// ExitFilterPackage is called when exiting the filterPackage production.
	ExitFilterPackage(c *FilterPackageContext)

	// ExitFilterPackageImportPart is called when exiting the filterPackageImportPart production.
	ExitFilterPackageImportPart(c *FilterPackageImportPartContext)

	// ExitFilterPackageMember is called when exiting the filterPackageMember production.
	ExitFilterPackageMember(c *FilterPackageMemberContext)

	// ExitVisibilityIndicator is called when exiting the visibilityIndicator production.
	ExitVisibilityIndicator(c *VisibilityIndicatorContext)

	// ExitDefinitionElement is called when exiting the definitionElement production.
	ExitDefinitionElement(c *DefinitionElementContext)

	// ExitUsageElement is called when exiting the usageElement production.
	ExitUsageElement(c *UsageElementContext)

	// ExitBasicDefinitionPrefix is called when exiting the basicDefinitionPrefix production.
	ExitBasicDefinitionPrefix(c *BasicDefinitionPrefixContext)

	// ExitDefinitionExtensionKeyword is called when exiting the definitionExtensionKeyword production.
	ExitDefinitionExtensionKeyword(c *DefinitionExtensionKeywordContext)

	// ExitDefinitionPrefix is called when exiting the definitionPrefix production.
	ExitDefinitionPrefix(c *DefinitionPrefixContext)

	// ExitDefinition is called when exiting the definition production.
	ExitDefinition(c *DefinitionContext)

	// ExitDefinitionDeclaration is called when exiting the definitionDeclaration production.
	ExitDefinitionDeclaration(c *DefinitionDeclarationContext)

	// ExitDefinitionBody is called when exiting the definitionBody production.
	ExitDefinitionBody(c *DefinitionBodyContext)

	// ExitDefinitionBodyItem is called when exiting the definitionBodyItem production.
	ExitDefinitionBodyItem(c *DefinitionBodyItemContext)

	// ExitEndFeatureMember is called when exiting the endFeatureMember production.
	ExitEndFeatureMember(c *EndFeatureMemberContext)

	// ExitEndFeatureDeclaration is called when exiting the endFeatureDeclaration production.
	ExitEndFeatureDeclaration(c *EndFeatureDeclarationContext)

	// ExitDefinitionMember is called when exiting the definitionMember production.
	ExitDefinitionMember(c *DefinitionMemberContext)

	// ExitVariantUsageMember is called when exiting the variantUsageMember production.
	ExitVariantUsageMember(c *VariantUsageMemberContext)

	// ExitNonOccurrenceUsageMember is called when exiting the nonOccurrenceUsageMember production.
	ExitNonOccurrenceUsageMember(c *NonOccurrenceUsageMemberContext)

	// ExitOccurrenceUsageMember is called when exiting the occurrenceUsageMember production.
	ExitOccurrenceUsageMember(c *OccurrenceUsageMemberContext)

	// ExitStructureUsageMember is called when exiting the structureUsageMember production.
	ExitStructureUsageMember(c *StructureUsageMemberContext)

	// ExitBehaviorUsageMember is called when exiting the behaviorUsageMember production.
	ExitBehaviorUsageMember(c *BehaviorUsageMemberContext)

	// ExitFeatureDirection is called when exiting the featureDirection production.
	ExitFeatureDirection(c *FeatureDirectionContext)

	// ExitRefPrefix is called when exiting the refPrefix production.
	ExitRefPrefix(c *RefPrefixContext)

	// ExitBasicUsagePrefix is called when exiting the basicUsagePrefix production.
	ExitBasicUsagePrefix(c *BasicUsagePrefixContext)

	// ExitEndUsagePrefix is called when exiting the endUsagePrefix production.
	ExitEndUsagePrefix(c *EndUsagePrefixContext)

	// ExitOwnedCrossFeatureMember is called when exiting the ownedCrossFeatureMember production.
	ExitOwnedCrossFeatureMember(c *OwnedCrossFeatureMemberContext)

	// ExitOwnedCrossFeature is called when exiting the ownedCrossFeature production.
	ExitOwnedCrossFeature(c *OwnedCrossFeatureContext)

	// ExitUsageExtensionKeyword is called when exiting the usageExtensionKeyword production.
	ExitUsageExtensionKeyword(c *UsageExtensionKeywordContext)

	// ExitUnextendedUsagePrefix is called when exiting the unextendedUsagePrefix production.
	ExitUnextendedUsagePrefix(c *UnextendedUsagePrefixContext)

	// ExitUsagePrefix is called when exiting the usagePrefix production.
	ExitUsagePrefix(c *UsagePrefixContext)

	// ExitUsage is called when exiting the usage production.
	ExitUsage(c *UsageContext)

	// ExitUsageDeclaration is called when exiting the usageDeclaration production.
	ExitUsageDeclaration(c *UsageDeclarationContext)

	// ExitUsageCompletion is called when exiting the usageCompletion production.
	ExitUsageCompletion(c *UsageCompletionContext)

	// ExitUsageBody is called when exiting the usageBody production.
	ExitUsageBody(c *UsageBodyContext)

	// ExitValuePart is called when exiting the valuePart production.
	ExitValuePart(c *ValuePartContext)

	// ExitFeatureValue is called when exiting the featureValue production.
	ExitFeatureValue(c *FeatureValueContext)

	// ExitDefaultReferenceUsage is called when exiting the defaultReferenceUsage production.
	ExitDefaultReferenceUsage(c *DefaultReferenceUsageContext)

	// ExitReferenceUsage is called when exiting the referenceUsage production.
	ExitReferenceUsage(c *ReferenceUsageContext)

	// ExitVariantReference is called when exiting the variantReference production.
	ExitVariantReference(c *VariantReferenceContext)

	// ExitNonOccurrenceUsageElement is called when exiting the nonOccurrenceUsageElement production.
	ExitNonOccurrenceUsageElement(c *NonOccurrenceUsageElementContext)

	// ExitOccurrenceUsageElement is called when exiting the occurrenceUsageElement production.
	ExitOccurrenceUsageElement(c *OccurrenceUsageElementContext)

	// ExitStructureUsageElement is called when exiting the structureUsageElement production.
	ExitStructureUsageElement(c *StructureUsageElementContext)

	// ExitBehaviorUsageElement is called when exiting the behaviorUsageElement production.
	ExitBehaviorUsageElement(c *BehaviorUsageElementContext)

	// ExitVariantUsageElement is called when exiting the variantUsageElement production.
	ExitVariantUsageElement(c *VariantUsageElementContext)

	// ExitSubclassificationPart is called when exiting the subclassificationPart production.
	ExitSubclassificationPart(c *SubclassificationPartContext)

	// ExitOwnedSubclassification is called when exiting the ownedSubclassification production.
	ExitOwnedSubclassification(c *OwnedSubclassificationContext)

	// ExitFeatureSpecializationPart is called when exiting the featureSpecializationPart production.
	ExitFeatureSpecializationPart(c *FeatureSpecializationPartContext)

	// ExitFeatureSpecialization is called when exiting the featureSpecialization production.
	ExitFeatureSpecialization(c *FeatureSpecializationContext)

	// ExitTypings is called when exiting the typings production.
	ExitTypings(c *TypingsContext)

	// ExitTypedBy is called when exiting the typedBy production.
	ExitTypedBy(c *TypedByContext)

	// ExitOwnedFeatureTyping is called when exiting the ownedFeatureTyping production.
	ExitOwnedFeatureTyping(c *OwnedFeatureTypingContext)

	// ExitSubsettings is called when exiting the subsettings production.
	ExitSubsettings(c *SubsettingsContext)

	// ExitOwnedSubsetting is called when exiting the ownedSubsetting production.
	ExitOwnedSubsetting(c *OwnedSubsettingContext)

	// ExitReferences is called when exiting the references production.
	ExitReferences(c *ReferencesContext)

	// ExitOwnedReferenceSubsetting is called when exiting the ownedReferenceSubsetting production.
	ExitOwnedReferenceSubsetting(c *OwnedReferenceSubsettingContext)

	// ExitCrosses is called when exiting the crosses production.
	ExitCrosses(c *CrossesContext)

	// ExitOwnedCrossSubsetting is called when exiting the ownedCrossSubsetting production.
	ExitOwnedCrossSubsetting(c *OwnedCrossSubsettingContext)

	// ExitRedefinitions is called when exiting the redefinitions production.
	ExitRedefinitions(c *RedefinitionsContext)

	// ExitOwnedRedefinition is called when exiting the ownedRedefinition production.
	ExitOwnedRedefinition(c *OwnedRedefinitionContext)

	// ExitOwnedFeatureChain is called when exiting the ownedFeatureChain production.
	ExitOwnedFeatureChain(c *OwnedFeatureChainContext)

	// ExitOwnedFeatureChaining is called when exiting the ownedFeatureChaining production.
	ExitOwnedFeatureChaining(c *OwnedFeatureChainingContext)

	// ExitSpecializes is called when exiting the specializes production.
	ExitSpecializes(c *SpecializesContext)

	// ExitDefinedBy is called when exiting the definedBy production.
	ExitDefinedBy(c *DefinedByContext)

	// ExitSubsetsKw is called when exiting the subsetsKw production.
	ExitSubsetsKw(c *SubsetsKwContext)

	// ExitReferencesKw is called when exiting the referencesKw production.
	ExitReferencesKw(c *ReferencesKwContext)

	// ExitCrossesKw is called when exiting the crossesKw production.
	ExitCrossesKw(c *CrossesKwContext)

	// ExitRedefinesKw is called when exiting the redefinesKw production.
	ExitRedefinesKw(c *RedefinesKwContext)

	// ExitMultiplicityPart is called when exiting the multiplicityPart production.
	ExitMultiplicityPart(c *MultiplicityPartContext)

	// ExitOwnedMultiplicity is called when exiting the ownedMultiplicity production.
	ExitOwnedMultiplicity(c *OwnedMultiplicityContext)

	// ExitMultiplicityRange is called when exiting the multiplicityRange production.
	ExitMultiplicityRange(c *MultiplicityRangeContext)

	// ExitMultiplicityExpressionMember is called when exiting the multiplicityExpressionMember production.
	ExitMultiplicityExpressionMember(c *MultiplicityExpressionMemberContext)

	// ExitAttributeDefinition is called when exiting the attributeDefinition production.
	ExitAttributeDefinition(c *AttributeDefinitionContext)

	// ExitAttributeUsage is called when exiting the attributeUsage production.
	ExitAttributeUsage(c *AttributeUsageContext)

	// ExitEnumerationDefinition is called when exiting the enumerationDefinition production.
	ExitEnumerationDefinition(c *EnumerationDefinitionContext)

	// ExitEnumerationBody is called when exiting the enumerationBody production.
	ExitEnumerationBody(c *EnumerationBodyContext)

	// ExitEnumerationUsageMember is called when exiting the enumerationUsageMember production.
	ExitEnumerationUsageMember(c *EnumerationUsageMemberContext)

	// ExitEnumeratedValue is called when exiting the enumeratedValue production.
	ExitEnumeratedValue(c *EnumeratedValueContext)

	// ExitEnumerationUsage is called when exiting the enumerationUsage production.
	ExitEnumerationUsage(c *EnumerationUsageContext)

	// ExitOccurrenceDefinitionPrefix is called when exiting the occurrenceDefinitionPrefix production.
	ExitOccurrenceDefinitionPrefix(c *OccurrenceDefinitionPrefixContext)

	// ExitOccurrenceDefinition is called when exiting the occurrenceDefinition production.
	ExitOccurrenceDefinition(c *OccurrenceDefinitionContext)

	// ExitIndividualDefinition is called when exiting the individualDefinition production.
	ExitIndividualDefinition(c *IndividualDefinitionContext)

	// ExitOccurrenceUsagePrefix is called when exiting the occurrenceUsagePrefix production.
	ExitOccurrenceUsagePrefix(c *OccurrenceUsagePrefixContext)

	// ExitOccurrenceUsage is called when exiting the occurrenceUsage production.
	ExitOccurrenceUsage(c *OccurrenceUsageContext)

	// ExitIndividualUsage is called when exiting the individualUsage production.
	ExitIndividualUsage(c *IndividualUsageContext)

	// ExitPortionUsage is called when exiting the portionUsage production.
	ExitPortionUsage(c *PortionUsageContext)

	// ExitPortionKind is called when exiting the portionKind production.
	ExitPortionKind(c *PortionKindContext)

	// ExitEventOccurrenceUsage is called when exiting the eventOccurrenceUsage production.
	ExitEventOccurrenceUsage(c *EventOccurrenceUsageContext)

	// ExitSourceSuccessionMember is called when exiting the sourceSuccessionMember production.
	ExitSourceSuccessionMember(c *SourceSuccessionMemberContext)

	// ExitSourceSuccession is called when exiting the sourceSuccession production.
	ExitSourceSuccession(c *SourceSuccessionContext)

	// ExitSourceEndMember is called when exiting the sourceEndMember production.
	ExitSourceEndMember(c *SourceEndMemberContext)

	// ExitSourceEnd is called when exiting the sourceEnd production.
	ExitSourceEnd(c *SourceEndContext)

	// ExitItemDefinition is called when exiting the itemDefinition production.
	ExitItemDefinition(c *ItemDefinitionContext)

	// ExitItemUsage is called when exiting the itemUsage production.
	ExitItemUsage(c *ItemUsageContext)

	// ExitPartDefinition is called when exiting the partDefinition production.
	ExitPartDefinition(c *PartDefinitionContext)

	// ExitPartUsage is called when exiting the partUsage production.
	ExitPartUsage(c *PartUsageContext)

	// ExitPortDefinition is called when exiting the portDefinition production.
	ExitPortDefinition(c *PortDefinitionContext)

	// ExitPortUsage is called when exiting the portUsage production.
	ExitPortUsage(c *PortUsageContext)

	// ExitConjugatedPortTyping is called when exiting the conjugatedPortTyping production.
	ExitConjugatedPortTyping(c *ConjugatedPortTypingContext)

	// ExitConnectionDefinition is called when exiting the connectionDefinition production.
	ExitConnectionDefinition(c *ConnectionDefinitionContext)

	// ExitConnectionUsage is called when exiting the connectionUsage production.
	ExitConnectionUsage(c *ConnectionUsageContext)

	// ExitConnectorPart is called when exiting the connectorPart production.
	ExitConnectorPart(c *ConnectorPartContext)

	// ExitBinaryConnectorPart is called when exiting the binaryConnectorPart production.
	ExitBinaryConnectorPart(c *BinaryConnectorPartContext)

	// ExitNaryConnectorPart is called when exiting the naryConnectorPart production.
	ExitNaryConnectorPart(c *NaryConnectorPartContext)

	// ExitConnectorEndMember is called when exiting the connectorEndMember production.
	ExitConnectorEndMember(c *ConnectorEndMemberContext)

	// ExitConnectorEnd is called when exiting the connectorEnd production.
	ExitConnectorEnd(c *ConnectorEndContext)

	// ExitOwnedCrossMultiplicityMember is called when exiting the ownedCrossMultiplicityMember production.
	ExitOwnedCrossMultiplicityMember(c *OwnedCrossMultiplicityMemberContext)

	// ExitOwnedCrossMultiplicity is called when exiting the ownedCrossMultiplicity production.
	ExitOwnedCrossMultiplicity(c *OwnedCrossMultiplicityContext)

	// ExitBindingConnectorAsUsage is called when exiting the bindingConnectorAsUsage production.
	ExitBindingConnectorAsUsage(c *BindingConnectorAsUsageContext)

	// ExitSuccessionAsUsage is called when exiting the successionAsUsage production.
	ExitSuccessionAsUsage(c *SuccessionAsUsageContext)

	// ExitInterfaceDefinition is called when exiting the interfaceDefinition production.
	ExitInterfaceDefinition(c *InterfaceDefinitionContext)

	// ExitInterfaceBody is called when exiting the interfaceBody production.
	ExitInterfaceBody(c *InterfaceBodyContext)

	// ExitInterfaceBodyItem is called when exiting the interfaceBodyItem production.
	ExitInterfaceBodyItem(c *InterfaceBodyItemContext)

	// ExitInterfaceNonOccurrenceUsageMember is called when exiting the interfaceNonOccurrenceUsageMember production.
	ExitInterfaceNonOccurrenceUsageMember(c *InterfaceNonOccurrenceUsageMemberContext)

	// ExitInterfaceNonOccurrenceUsageElement is called when exiting the interfaceNonOccurrenceUsageElement production.
	ExitInterfaceNonOccurrenceUsageElement(c *InterfaceNonOccurrenceUsageElementContext)

	// ExitInterfaceOccurrenceUsageMember is called when exiting the interfaceOccurrenceUsageMember production.
	ExitInterfaceOccurrenceUsageMember(c *InterfaceOccurrenceUsageMemberContext)

	// ExitInterfaceOccurrenceUsageElement is called when exiting the interfaceOccurrenceUsageElement production.
	ExitInterfaceOccurrenceUsageElement(c *InterfaceOccurrenceUsageElementContext)

	// ExitDefaultInterfaceEnd is called when exiting the defaultInterfaceEnd production.
	ExitDefaultInterfaceEnd(c *DefaultInterfaceEndContext)

	// ExitInterfaceUsage is called when exiting the interfaceUsage production.
	ExitInterfaceUsage(c *InterfaceUsageContext)

	// ExitInterfaceUsageDeclaration is called when exiting the interfaceUsageDeclaration production.
	ExitInterfaceUsageDeclaration(c *InterfaceUsageDeclarationContext)

	// ExitInterfacePart is called when exiting the interfacePart production.
	ExitInterfacePart(c *InterfacePartContext)

	// ExitBinaryInterfacePart is called when exiting the binaryInterfacePart production.
	ExitBinaryInterfacePart(c *BinaryInterfacePartContext)

	// ExitNaryInterfacePart is called when exiting the naryInterfacePart production.
	ExitNaryInterfacePart(c *NaryInterfacePartContext)

	// ExitInterfaceEndMember is called when exiting the interfaceEndMember production.
	ExitInterfaceEndMember(c *InterfaceEndMemberContext)

	// ExitInterfaceEnd is called when exiting the interfaceEnd production.
	ExitInterfaceEnd(c *InterfaceEndContext)

	// ExitAllocationDefinition is called when exiting the allocationDefinition production.
	ExitAllocationDefinition(c *AllocationDefinitionContext)

	// ExitAllocationUsage is called when exiting the allocationUsage production.
	ExitAllocationUsage(c *AllocationUsageContext)

	// ExitAllocationUsageDeclaration is called when exiting the allocationUsageDeclaration production.
	ExitAllocationUsageDeclaration(c *AllocationUsageDeclarationContext)

	// ExitFlowDefinition is called when exiting the flowDefinition production.
	ExitFlowDefinition(c *FlowDefinitionContext)

	// ExitMessage is called when exiting the message production.
	ExitMessage(c *MessageContext)

	// ExitMessageDeclaration is called when exiting the messageDeclaration production.
	ExitMessageDeclaration(c *MessageDeclarationContext)

	// ExitMessageEventMember is called when exiting the messageEventMember production.
	ExitMessageEventMember(c *MessageEventMemberContext)

	// ExitMessageEvent is called when exiting the messageEvent production.
	ExitMessageEvent(c *MessageEventContext)

	// ExitFlowUsage is called when exiting the flowUsage production.
	ExitFlowUsage(c *FlowUsageContext)

	// ExitSuccessionFlowUsage is called when exiting the successionFlowUsage production.
	ExitSuccessionFlowUsage(c *SuccessionFlowUsageContext)

	// ExitFlowDeclaration is called when exiting the flowDeclaration production.
	ExitFlowDeclaration(c *FlowDeclarationContext)

	// ExitFlowPayloadFeatureMember is called when exiting the flowPayloadFeatureMember production.
	ExitFlowPayloadFeatureMember(c *FlowPayloadFeatureMemberContext)

	// ExitFlowPayloadFeature is called when exiting the flowPayloadFeature production.
	ExitFlowPayloadFeature(c *FlowPayloadFeatureContext)

	// ExitPayloadFeature is called when exiting the payloadFeature production.
	ExitPayloadFeature(c *PayloadFeatureContext)

	// ExitPayloadFeatureSpecializationPart is called when exiting the payloadFeatureSpecializationPart production.
	ExitPayloadFeatureSpecializationPart(c *PayloadFeatureSpecializationPartContext)

	// ExitFlowEndMember is called when exiting the flowEndMember production.
	ExitFlowEndMember(c *FlowEndMemberContext)

	// ExitFlowEnd is called when exiting the flowEnd production.
	ExitFlowEnd(c *FlowEndContext)

	// ExitFlowEndSubsetting is called when exiting the flowEndSubsetting production.
	ExitFlowEndSubsetting(c *FlowEndSubsettingContext)

	// ExitFeatureChainPrefix is called when exiting the featureChainPrefix production.
	ExitFeatureChainPrefix(c *FeatureChainPrefixContext)

	// ExitFlowFeatureMember is called when exiting the flowFeatureMember production.
	ExitFlowFeatureMember(c *FlowFeatureMemberContext)

	// ExitFlowFeature is called when exiting the flowFeature production.
	ExitFlowFeature(c *FlowFeatureContext)

	// ExitFlowFeatureRedefinition is called when exiting the flowFeatureRedefinition production.
	ExitFlowFeatureRedefinition(c *FlowFeatureRedefinitionContext)

	// ExitActionDefinition is called when exiting the actionDefinition production.
	ExitActionDefinition(c *ActionDefinitionContext)

	// ExitActionBody is called when exiting the actionBody production.
	ExitActionBody(c *ActionBodyContext)

	// ExitActionBodyItem is called when exiting the actionBodyItem production.
	ExitActionBodyItem(c *ActionBodyItemContext)

	// ExitNonBehaviorBodyItem is called when exiting the nonBehaviorBodyItem production.
	ExitNonBehaviorBodyItem(c *NonBehaviorBodyItemContext)

	// ExitActionBehaviorMember is called when exiting the actionBehaviorMember production.
	ExitActionBehaviorMember(c *ActionBehaviorMemberContext)

	// ExitInitialNodeMember is called when exiting the initialNodeMember production.
	ExitInitialNodeMember(c *InitialNodeMemberContext)

	// ExitActionNodeMember is called when exiting the actionNodeMember production.
	ExitActionNodeMember(c *ActionNodeMemberContext)

	// ExitActionTargetSuccessionMember is called when exiting the actionTargetSuccessionMember production.
	ExitActionTargetSuccessionMember(c *ActionTargetSuccessionMemberContext)

	// ExitGuardedSuccessionMember is called when exiting the guardedSuccessionMember production.
	ExitGuardedSuccessionMember(c *GuardedSuccessionMemberContext)

	// ExitActionUsage is called when exiting the actionUsage production.
	ExitActionUsage(c *ActionUsageContext)

	// ExitActionUsageDeclaration is called when exiting the actionUsageDeclaration production.
	ExitActionUsageDeclaration(c *ActionUsageDeclarationContext)

	// ExitPerformActionUsage is called when exiting the performActionUsage production.
	ExitPerformActionUsage(c *PerformActionUsageContext)

	// ExitPerformActionUsageDeclaration is called when exiting the performActionUsageDeclaration production.
	ExitPerformActionUsageDeclaration(c *PerformActionUsageDeclarationContext)

	// ExitActionNode is called when exiting the actionNode production.
	ExitActionNode(c *ActionNodeContext)

	// ExitActionNodeUsageDeclaration is called when exiting the actionNodeUsageDeclaration production.
	ExitActionNodeUsageDeclaration(c *ActionNodeUsageDeclarationContext)

	// ExitActionNodePrefix is called when exiting the actionNodePrefix production.
	ExitActionNodePrefix(c *ActionNodePrefixContext)

	// ExitControlNode is called when exiting the controlNode production.
	ExitControlNode(c *ControlNodeContext)

	// ExitControlNodePrefix is called when exiting the controlNodePrefix production.
	ExitControlNodePrefix(c *ControlNodePrefixContext)

	// ExitMergeNode is called when exiting the mergeNode production.
	ExitMergeNode(c *MergeNodeContext)

	// ExitDecisionNode is called when exiting the decisionNode production.
	ExitDecisionNode(c *DecisionNodeContext)

	// ExitJoinNode is called when exiting the joinNode production.
	ExitJoinNode(c *JoinNodeContext)

	// ExitForkNode is called when exiting the forkNode production.
	ExitForkNode(c *ForkNodeContext)

	// ExitAcceptNode is called when exiting the acceptNode production.
	ExitAcceptNode(c *AcceptNodeContext)

	// ExitAcceptNodeDeclaration is called when exiting the acceptNodeDeclaration production.
	ExitAcceptNodeDeclaration(c *AcceptNodeDeclarationContext)

	// ExitAcceptParameterPart is called when exiting the acceptParameterPart production.
	ExitAcceptParameterPart(c *AcceptParameterPartContext)

	// ExitPayloadParameterMember is called when exiting the payloadParameterMember production.
	ExitPayloadParameterMember(c *PayloadParameterMemberContext)

	// ExitPayloadParameter is called when exiting the payloadParameter production.
	ExitPayloadParameter(c *PayloadParameterContext)

	// ExitTriggerValuePart is called when exiting the triggerValuePart production.
	ExitTriggerValuePart(c *TriggerValuePartContext)

	// ExitTriggerFeatureValue is called when exiting the triggerFeatureValue production.
	ExitTriggerFeatureValue(c *TriggerFeatureValueContext)

	// ExitTriggerExpression is called when exiting the triggerExpression production.
	ExitTriggerExpression(c *TriggerExpressionContext)

	// ExitSendNode is called when exiting the sendNode production.
	ExitSendNode(c *SendNodeContext)

	// ExitSendNodeDeclaration is called when exiting the sendNodeDeclaration production.
	ExitSendNodeDeclaration(c *SendNodeDeclarationContext)

	// ExitSenderReceiverPart is called when exiting the senderReceiverPart production.
	ExitSenderReceiverPart(c *SenderReceiverPartContext)

	// ExitNodeParameterMember is called when exiting the nodeParameterMember production.
	ExitNodeParameterMember(c *NodeParameterMemberContext)

	// ExitNodeParameter is called when exiting the nodeParameter production.
	ExitNodeParameter(c *NodeParameterContext)

	// ExitAssignmentNode is called when exiting the assignmentNode production.
	ExitAssignmentNode(c *AssignmentNodeContext)

	// ExitAssignmentNodeDeclaration is called when exiting the assignmentNodeDeclaration production.
	ExitAssignmentNodeDeclaration(c *AssignmentNodeDeclarationContext)

	// ExitAssignmentTargetMember is called when exiting the assignmentTargetMember production.
	ExitAssignmentTargetMember(c *AssignmentTargetMemberContext)

	// ExitAssignmentTargetParameter is called when exiting the assignmentTargetParameter production.
	ExitAssignmentTargetParameter(c *AssignmentTargetParameterContext)

	// ExitFeatureChainMember is called when exiting the featureChainMember production.
	ExitFeatureChainMember(c *FeatureChainMemberContext)

	// ExitOwnedFeatureChainMember is called when exiting the ownedFeatureChainMember production.
	ExitOwnedFeatureChainMember(c *OwnedFeatureChainMemberContext)

	// ExitTerminateNode is called when exiting the terminateNode production.
	ExitTerminateNode(c *TerminateNodeContext)

	// ExitIfNode is called when exiting the ifNode production.
	ExitIfNode(c *IfNodeContext)

	// ExitActionBodyParameter is called when exiting the actionBodyParameter production.
	ExitActionBodyParameter(c *ActionBodyParameterContext)

	// ExitWhileLoopNode is called when exiting the whileLoopNode production.
	ExitWhileLoopNode(c *WhileLoopNodeContext)

	// ExitForLoopNode is called when exiting the forLoopNode production.
	ExitForLoopNode(c *ForLoopNodeContext)

	// ExitForVariableDeclaration is called when exiting the forVariableDeclaration production.
	ExitForVariableDeclaration(c *ForVariableDeclarationContext)

	// ExitActionTargetSuccession is called when exiting the actionTargetSuccession production.
	ExitActionTargetSuccession(c *ActionTargetSuccessionContext)

	// ExitTargetSuccession is called when exiting the targetSuccession production.
	ExitTargetSuccession(c *TargetSuccessionContext)

	// ExitGuardedTargetSuccession is called when exiting the guardedTargetSuccession production.
	ExitGuardedTargetSuccession(c *GuardedTargetSuccessionContext)

	// ExitDefaultTargetSuccession is called when exiting the defaultTargetSuccession production.
	ExitDefaultTargetSuccession(c *DefaultTargetSuccessionContext)

	// ExitGuardedSuccession is called when exiting the guardedSuccession production.
	ExitGuardedSuccession(c *GuardedSuccessionContext)

	// ExitStateDefinition is called when exiting the stateDefinition production.
	ExitStateDefinition(c *StateDefinitionContext)

	// ExitStateDefBody is called when exiting the stateDefBody production.
	ExitStateDefBody(c *StateDefBodyContext)

	// ExitStateBodyItem is called when exiting the stateBodyItem production.
	ExitStateBodyItem(c *StateBodyItemContext)

	// ExitEntryActionMember is called when exiting the entryActionMember production.
	ExitEntryActionMember(c *EntryActionMemberContext)

	// ExitDoActionMember is called when exiting the doActionMember production.
	ExitDoActionMember(c *DoActionMemberContext)

	// ExitExitActionMember is called when exiting the exitActionMember production.
	ExitExitActionMember(c *ExitActionMemberContext)

	// ExitEntryTransitionMember is called when exiting the entryTransitionMember production.
	ExitEntryTransitionMember(c *EntryTransitionMemberContext)

	// ExitStateActionUsage is called when exiting the stateActionUsage production.
	ExitStateActionUsage(c *StateActionUsageContext)

	// ExitStatePerformActionUsage is called when exiting the statePerformActionUsage production.
	ExitStatePerformActionUsage(c *StatePerformActionUsageContext)

	// ExitStateAcceptActionUsage is called when exiting the stateAcceptActionUsage production.
	ExitStateAcceptActionUsage(c *StateAcceptActionUsageContext)

	// ExitStateSendActionUsage is called when exiting the stateSendActionUsage production.
	ExitStateSendActionUsage(c *StateSendActionUsageContext)

	// ExitStateAssignmentActionUsage is called when exiting the stateAssignmentActionUsage production.
	ExitStateAssignmentActionUsage(c *StateAssignmentActionUsageContext)

	// ExitTransitionUsageMember is called when exiting the transitionUsageMember production.
	ExitTransitionUsageMember(c *TransitionUsageMemberContext)

	// ExitTargetTransitionUsageMember is called when exiting the targetTransitionUsageMember production.
	ExitTargetTransitionUsageMember(c *TargetTransitionUsageMemberContext)

	// ExitStateUsage is called when exiting the stateUsage production.
	ExitStateUsage(c *StateUsageContext)

	// ExitStateUsageBody is called when exiting the stateUsageBody production.
	ExitStateUsageBody(c *StateUsageBodyContext)

	// ExitExhibitStateUsage is called when exiting the exhibitStateUsage production.
	ExitExhibitStateUsage(c *ExhibitStateUsageContext)

	// ExitTransitionUsage is called when exiting the transitionUsage production.
	ExitTransitionUsage(c *TransitionUsageContext)

	// ExitTargetTransitionUsage is called when exiting the targetTransitionUsage production.
	ExitTargetTransitionUsage(c *TargetTransitionUsageContext)

	// ExitTriggerActionMember is called when exiting the triggerActionMember production.
	ExitTriggerActionMember(c *TriggerActionMemberContext)

	// ExitTriggerAction is called when exiting the triggerAction production.
	ExitTriggerAction(c *TriggerActionContext)

	// ExitGuardExpressionMember is called when exiting the guardExpressionMember production.
	ExitGuardExpressionMember(c *GuardExpressionMemberContext)

	// ExitEffectBehaviorMember is called when exiting the effectBehaviorMember production.
	ExitEffectBehaviorMember(c *EffectBehaviorMemberContext)

	// ExitEffectBehaviorUsage is called when exiting the effectBehaviorUsage production.
	ExitEffectBehaviorUsage(c *EffectBehaviorUsageContext)

	// ExitTransitionPerformActionUsage is called when exiting the transitionPerformActionUsage production.
	ExitTransitionPerformActionUsage(c *TransitionPerformActionUsageContext)

	// ExitTransitionAcceptActionUsage is called when exiting the transitionAcceptActionUsage production.
	ExitTransitionAcceptActionUsage(c *TransitionAcceptActionUsageContext)

	// ExitTransitionSendActionUsage is called when exiting the transitionSendActionUsage production.
	ExitTransitionSendActionUsage(c *TransitionSendActionUsageContext)

	// ExitTransitionAssignmentActionUsage is called when exiting the transitionAssignmentActionUsage production.
	ExitTransitionAssignmentActionUsage(c *TransitionAssignmentActionUsageContext)

	// ExitTransitionSuccessionMember is called when exiting the transitionSuccessionMember production.
	ExitTransitionSuccessionMember(c *TransitionSuccessionMemberContext)

	// ExitTransitionSuccession is called when exiting the transitionSuccession production.
	ExitTransitionSuccession(c *TransitionSuccessionContext)

	// ExitCalculationDefinition is called when exiting the calculationDefinition production.
	ExitCalculationDefinition(c *CalculationDefinitionContext)

	// ExitCalculationUsage is called when exiting the calculationUsage production.
	ExitCalculationUsage(c *CalculationUsageContext)

	// ExitCalculationBody is called when exiting the calculationBody production.
	ExitCalculationBody(c *CalculationBodyContext)

	// ExitCalculationBodyPart is called when exiting the calculationBodyPart production.
	ExitCalculationBodyPart(c *CalculationBodyPartContext)

	// ExitCalculationBodyItem is called when exiting the calculationBodyItem production.
	ExitCalculationBodyItem(c *CalculationBodyItemContext)

	// ExitReturnParameterMember is called when exiting the returnParameterMember production.
	ExitReturnParameterMember(c *ReturnParameterMemberContext)

	// ExitResultExpressionMember is called when exiting the resultExpressionMember production.
	ExitResultExpressionMember(c *ResultExpressionMemberContext)

	// ExitConstraintDefinition is called when exiting the constraintDefinition production.
	ExitConstraintDefinition(c *ConstraintDefinitionContext)

	// ExitConstraintUsage is called when exiting the constraintUsage production.
	ExitConstraintUsage(c *ConstraintUsageContext)

	// ExitAssertConstraintUsage is called when exiting the assertConstraintUsage production.
	ExitAssertConstraintUsage(c *AssertConstraintUsageContext)

	// ExitConstraintUsageDeclaration is called when exiting the constraintUsageDeclaration production.
	ExitConstraintUsageDeclaration(c *ConstraintUsageDeclarationContext)

	// ExitRequirementDefinition is called when exiting the requirementDefinition production.
	ExitRequirementDefinition(c *RequirementDefinitionContext)

	// ExitRequirementBody is called when exiting the requirementBody production.
	ExitRequirementBody(c *RequirementBodyContext)

	// ExitRequirementBodyItem is called when exiting the requirementBodyItem production.
	ExitRequirementBodyItem(c *RequirementBodyItemContext)

	// ExitSubjectMember is called when exiting the subjectMember production.
	ExitSubjectMember(c *SubjectMemberContext)

	// ExitSubjectUsage is called when exiting the subjectUsage production.
	ExitSubjectUsage(c *SubjectUsageContext)

	// ExitRequirementConstraintMember is called when exiting the requirementConstraintMember production.
	ExitRequirementConstraintMember(c *RequirementConstraintMemberContext)

	// ExitRequirementKind is called when exiting the requirementKind production.
	ExitRequirementKind(c *RequirementKindContext)

	// ExitRequirementConstraintUsage is called when exiting the requirementConstraintUsage production.
	ExitRequirementConstraintUsage(c *RequirementConstraintUsageContext)

	// ExitFramedConcernMember is called when exiting the framedConcernMember production.
	ExitFramedConcernMember(c *FramedConcernMemberContext)

	// ExitFramedConcernUsage is called when exiting the framedConcernUsage production.
	ExitFramedConcernUsage(c *FramedConcernUsageContext)

	// ExitCalculationUsageDeclaration is called when exiting the calculationUsageDeclaration production.
	ExitCalculationUsageDeclaration(c *CalculationUsageDeclarationContext)

	// ExitActorMember is called when exiting the actorMember production.
	ExitActorMember(c *ActorMemberContext)

	// ExitActorUsage is called when exiting the actorUsage production.
	ExitActorUsage(c *ActorUsageContext)

	// ExitStakeholderMember is called when exiting the stakeholderMember production.
	ExitStakeholderMember(c *StakeholderMemberContext)

	// ExitStakeholderUsage is called when exiting the stakeholderUsage production.
	ExitStakeholderUsage(c *StakeholderUsageContext)

	// ExitRequirementUsage is called when exiting the requirementUsage production.
	ExitRequirementUsage(c *RequirementUsageContext)

	// ExitSatisfyRequirementUsage is called when exiting the satisfyRequirementUsage production.
	ExitSatisfyRequirementUsage(c *SatisfyRequirementUsageContext)

	// ExitSatisfactionSubjectMember is called when exiting the satisfactionSubjectMember production.
	ExitSatisfactionSubjectMember(c *SatisfactionSubjectMemberContext)

	// ExitSatisfactionParameter is called when exiting the satisfactionParameter production.
	ExitSatisfactionParameter(c *SatisfactionParameterContext)

	// ExitConcernDefinition is called when exiting the concernDefinition production.
	ExitConcernDefinition(c *ConcernDefinitionContext)

	// ExitConcernUsage is called when exiting the concernUsage production.
	ExitConcernUsage(c *ConcernUsageContext)

	// ExitCaseDefinition is called when exiting the caseDefinition production.
	ExitCaseDefinition(c *CaseDefinitionContext)

	// ExitCaseUsage is called when exiting the caseUsage production.
	ExitCaseUsage(c *CaseUsageContext)

	// ExitCaseBody is called when exiting the caseBody production.
	ExitCaseBody(c *CaseBodyContext)

	// ExitCaseBodyItem is called when exiting the caseBodyItem production.
	ExitCaseBodyItem(c *CaseBodyItemContext)

	// ExitObjectiveMember is called when exiting the objectiveMember production.
	ExitObjectiveMember(c *ObjectiveMemberContext)

	// ExitObjectiveRequirementUsage is called when exiting the objectiveRequirementUsage production.
	ExitObjectiveRequirementUsage(c *ObjectiveRequirementUsageContext)

	// ExitAnalysisCaseDefinition is called when exiting the analysisCaseDefinition production.
	ExitAnalysisCaseDefinition(c *AnalysisCaseDefinitionContext)

	// ExitAnalysisCaseUsage is called when exiting the analysisCaseUsage production.
	ExitAnalysisCaseUsage(c *AnalysisCaseUsageContext)

	// ExitVerificationCaseDefinition is called when exiting the verificationCaseDefinition production.
	ExitVerificationCaseDefinition(c *VerificationCaseDefinitionContext)

	// ExitVerificationCaseUsage is called when exiting the verificationCaseUsage production.
	ExitVerificationCaseUsage(c *VerificationCaseUsageContext)

	// ExitRequirementVerificationMember is called when exiting the requirementVerificationMember production.
	ExitRequirementVerificationMember(c *RequirementVerificationMemberContext)

	// ExitRequirementVerificationUsage is called when exiting the requirementVerificationUsage production.
	ExitRequirementVerificationUsage(c *RequirementVerificationUsageContext)

	// ExitUseCaseDefinition is called when exiting the useCaseDefinition production.
	ExitUseCaseDefinition(c *UseCaseDefinitionContext)

	// ExitUseCaseUsage is called when exiting the useCaseUsage production.
	ExitUseCaseUsage(c *UseCaseUsageContext)

	// ExitIncludeUseCaseUsage is called when exiting the includeUseCaseUsage production.
	ExitIncludeUseCaseUsage(c *IncludeUseCaseUsageContext)

	// ExitViewDefinition is called when exiting the viewDefinition production.
	ExitViewDefinition(c *ViewDefinitionContext)

	// ExitViewDefinitionBody is called when exiting the viewDefinitionBody production.
	ExitViewDefinitionBody(c *ViewDefinitionBodyContext)

	// ExitViewDefinitionBodyItem is called when exiting the viewDefinitionBodyItem production.
	ExitViewDefinitionBodyItem(c *ViewDefinitionBodyItemContext)

	// ExitViewRenderingMember is called when exiting the viewRenderingMember production.
	ExitViewRenderingMember(c *ViewRenderingMemberContext)

	// ExitViewRenderingUsage is called when exiting the viewRenderingUsage production.
	ExitViewRenderingUsage(c *ViewRenderingUsageContext)

	// ExitViewUsage is called when exiting the viewUsage production.
	ExitViewUsage(c *ViewUsageContext)

	// ExitViewBody is called when exiting the viewBody production.
	ExitViewBody(c *ViewBodyContext)

	// ExitViewBodyItem is called when exiting the viewBodyItem production.
	ExitViewBodyItem(c *ViewBodyItemContext)

	// ExitExpose is called when exiting the expose production.
	ExitExpose(c *ExposeContext)

	// ExitMembershipExpose is called when exiting the membershipExpose production.
	ExitMembershipExpose(c *MembershipExposeContext)

	// ExitNamespaceExpose is called when exiting the namespaceExpose production.
	ExitNamespaceExpose(c *NamespaceExposeContext)

	// ExitViewpointDefinition is called when exiting the viewpointDefinition production.
	ExitViewpointDefinition(c *ViewpointDefinitionContext)

	// ExitViewpointUsage is called when exiting the viewpointUsage production.
	ExitViewpointUsage(c *ViewpointUsageContext)

	// ExitRenderingDefinition is called when exiting the renderingDefinition production.
	ExitRenderingDefinition(c *RenderingDefinitionContext)

	// ExitRenderingUsage is called when exiting the renderingUsage production.
	ExitRenderingUsage(c *RenderingUsageContext)

	// ExitMetadataDefinition is called when exiting the metadataDefinition production.
	ExitMetadataDefinition(c *MetadataDefinitionContext)

	// ExitPrefixMetadataAnnotation is called when exiting the prefixMetadataAnnotation production.
	ExitPrefixMetadataAnnotation(c *PrefixMetadataAnnotationContext)

	// ExitPrefixMetadataMember is called when exiting the prefixMetadataMember production.
	ExitPrefixMetadataMember(c *PrefixMetadataMemberContext)

	// ExitPrefixMetadataUsage is called when exiting the prefixMetadataUsage production.
	ExitPrefixMetadataUsage(c *PrefixMetadataUsageContext)

	// ExitMetadataUsage is called when exiting the metadataUsage production.
	ExitMetadataUsage(c *MetadataUsageContext)

	// ExitMetadataUsageDeclaration is called when exiting the metadataUsageDeclaration production.
	ExitMetadataUsageDeclaration(c *MetadataUsageDeclarationContext)

	// ExitMetadataBody is called when exiting the metadataBody production.
	ExitMetadataBody(c *MetadataBodyContext)

	// ExitMetadataBodyUsageMember is called when exiting the metadataBodyUsageMember production.
	ExitMetadataBodyUsageMember(c *MetadataBodyUsageMemberContext)

	// ExitMetadataBodyUsage is called when exiting the metadataBodyUsage production.
	ExitMetadataBodyUsage(c *MetadataBodyUsageContext)

	// ExitMetadataFeature is called when exiting the metadataFeature production.
	ExitMetadataFeature(c *MetadataFeatureContext)

	// ExitExtendedDefinition is called when exiting the extendedDefinition production.
	ExitExtendedDefinition(c *ExtendedDefinitionContext)

	// ExitExtendedUsage is called when exiting the extendedUsage production.
	ExitExtendedUsage(c *ExtendedUsageContext)

	// ExitOwnedExpression is called when exiting the ownedExpression production.
	ExitOwnedExpression(c *OwnedExpressionContext)

	// ExitConditionalBinaryOperator is called when exiting the conditionalBinaryOperator production.
	ExitConditionalBinaryOperator(c *ConditionalBinaryOperatorContext)

	// ExitBinaryOperator is called when exiting the binaryOperator production.
	ExitBinaryOperator(c *BinaryOperatorContext)

	// ExitUnaryOperator is called when exiting the unaryOperator production.
	ExitUnaryOperator(c *UnaryOperatorContext)

	// ExitClassificationTestOperator is called when exiting the classificationTestOperator production.
	ExitClassificationTestOperator(c *ClassificationTestOperatorContext)

	// ExitCastOperator is called when exiting the castOperator production.
	ExitCastOperator(c *CastOperatorContext)

	// ExitMetaCastOperator is called when exiting the metaCastOperator production.
	ExitMetaCastOperator(c *MetaCastOperatorContext)

	// ExitTypeReference is called when exiting the typeReference production.
	ExitTypeReference(c *TypeReferenceContext)

	// ExitPrimaryExpression is called when exiting the primaryExpression production.
	ExitPrimaryExpression(c *PrimaryExpressionContext)

	// ExitPrimaryExpressionSuffix is called when exiting the primaryExpressionSuffix production.
	ExitPrimaryExpressionSuffix(c *PrimaryExpressionSuffixContext)

	// ExitBaseExpression is called when exiting the baseExpression production.
	ExitBaseExpression(c *BaseExpressionContext)

	// ExitNullExpression is called when exiting the nullExpression production.
	ExitNullExpression(c *NullExpressionContext)

	// ExitFeatureReferenceExpression is called when exiting the featureReferenceExpression production.
	ExitFeatureReferenceExpression(c *FeatureReferenceExpressionContext)

	// ExitMetadataAccessExpression is called when exiting the metadataAccessExpression production.
	ExitMetadataAccessExpression(c *MetadataAccessExpressionContext)

	// ExitInvocationExpression is called when exiting the invocationExpression production.
	ExitInvocationExpression(c *InvocationExpressionContext)

	// ExitConstructorExpression is called when exiting the constructorExpression production.
	ExitConstructorExpression(c *ConstructorExpressionContext)

	// ExitBodyExpression is called when exiting the bodyExpression production.
	ExitBodyExpression(c *BodyExpressionContext)

	// ExitSequenceExpression is called when exiting the sequenceExpression production.
	ExitSequenceExpression(c *SequenceExpressionContext)

	// ExitSequenceExpressionList is called when exiting the sequenceExpressionList production.
	ExitSequenceExpressionList(c *SequenceExpressionListContext)

	// ExitArgumentList is called when exiting the argumentList production.
	ExitArgumentList(c *ArgumentListContext)

	// ExitPositionalArgumentList is called when exiting the positionalArgumentList production.
	ExitPositionalArgumentList(c *PositionalArgumentListContext)

	// ExitNamedArgumentList is called when exiting the namedArgumentList production.
	ExitNamedArgumentList(c *NamedArgumentListContext)

	// ExitNamedArgument is called when exiting the namedArgument production.
	ExitNamedArgument(c *NamedArgumentContext)

	// ExitFunctionBodyPart is called when exiting the functionBodyPart production.
	ExitFunctionBodyPart(c *FunctionBodyPartContext)

	// ExitResultExpressionMemberOpt is called when exiting the resultExpressionMemberOpt production.
	ExitResultExpressionMemberOpt(c *ResultExpressionMemberOptContext)

	// ExitTypeBodyElement is called when exiting the typeBodyElement production.
	ExitTypeBodyElement(c *TypeBodyElementContext)

	// ExitNonFeatureMember is called when exiting the nonFeatureMember production.
	ExitNonFeatureMember(c *NonFeatureMemberContext)

	// ExitFeatureMember is called when exiting the featureMember production.
	ExitFeatureMember(c *FeatureMemberContext)

	// ExitTypeFeatureMember is called when exiting the typeFeatureMember production.
	ExitTypeFeatureMember(c *TypeFeatureMemberContext)

	// ExitOwnedFeatureMember is called when exiting the ownedFeatureMember production.
	ExitOwnedFeatureMember(c *OwnedFeatureMemberContext)

	// ExitMemberElement is called when exiting the memberElement production.
	ExitMemberElement(c *MemberElementContext)

	// ExitNonFeatureElement is called when exiting the nonFeatureElement production.
	ExitNonFeatureElement(c *NonFeatureElementContext)

	// ExitFeatureElement is called when exiting the featureElement production.
	ExitFeatureElement(c *FeatureElementContext)

	// ExitReturnFeatureMember is called when exiting the returnFeatureMember production.
	ExitReturnFeatureMember(c *ReturnFeatureMemberContext)

	// ExitLiteralExpression is called when exiting the literalExpression production.
	ExitLiteralExpression(c *LiteralExpressionContext)

	// ExitLiteralBoolean is called when exiting the literalBoolean production.
	ExitLiteralBoolean(c *LiteralBooleanContext)

	// ExitLiteralString is called when exiting the literalString production.
	ExitLiteralString(c *LiteralStringContext)

	// ExitLiteralInteger is called when exiting the literalInteger production.
	ExitLiteralInteger(c *LiteralIntegerContext)

	// ExitLiteralReal is called when exiting the literalReal production.
	ExitLiteralReal(c *LiteralRealContext)

	// ExitLiteralInfinity is called when exiting the literalInfinity production.
	ExitLiteralInfinity(c *LiteralInfinityContext)

	// ExitNamespace_ is called when exiting the namespace_ production.
	ExitNamespace_(c *Namespace_Context)

	// ExitNamespaceDeclaration is called when exiting the namespaceDeclaration production.
	ExitNamespaceDeclaration(c *NamespaceDeclarationContext)

	// ExitNamespaceBody is called when exiting the namespaceBody production.
	ExitNamespaceBody(c *NamespaceBodyContext)

	// ExitNamespaceBodyElement is called when exiting the namespaceBodyElement production.
	ExitNamespaceBodyElement(c *NamespaceBodyElementContext)

	// ExitNamespaceMember is called when exiting the namespaceMember production.
	ExitNamespaceMember(c *NamespaceMemberContext)

	// ExitNamespaceFeatureMember is called when exiting the namespaceFeatureMember production.
	ExitNamespaceFeatureMember(c *NamespaceFeatureMemberContext)

	// ExitType_ is called when exiting the type_ production.
	ExitType_(c *Type_Context)

	// ExitTypePrefix is called when exiting the typePrefix production.
	ExitTypePrefix(c *TypePrefixContext)

	// ExitTypeDeclaration is called when exiting the typeDeclaration production.
	ExitTypeDeclaration(c *TypeDeclarationContext)

	// ExitSpecializationPart is called when exiting the specializationPart production.
	ExitSpecializationPart(c *SpecializationPartContext)

	// ExitOwnedSpecialization is called when exiting the ownedSpecialization production.
	ExitOwnedSpecialization(c *OwnedSpecializationContext)

	// ExitGeneralType is called when exiting the generalType production.
	ExitGeneralType(c *GeneralTypeContext)

	// ExitConjugationPart is called when exiting the conjugationPart production.
	ExitConjugationPart(c *ConjugationPartContext)

	// ExitOwnedConjugation is called when exiting the ownedConjugation production.
	ExitOwnedConjugation(c *OwnedConjugationContext)

	// ExitConjugates is called when exiting the conjugates production.
	ExitConjugates(c *ConjugatesContext)

	// ExitTypeRelationshipPart is called when exiting the typeRelationshipPart production.
	ExitTypeRelationshipPart(c *TypeRelationshipPartContext)

	// ExitDisjoiningPart is called when exiting the disjoiningPart production.
	ExitDisjoiningPart(c *DisjoiningPartContext)

	// ExitDisjoining is called when exiting the disjoining production.
	ExitDisjoining(c *DisjoiningContext)

	// ExitOwnedDisjoining is called when exiting the ownedDisjoining production.
	ExitOwnedDisjoining(c *OwnedDisjoiningContext)

	// ExitUnioningPart is called when exiting the unioningPart production.
	ExitUnioningPart(c *UnioningPartContext)

	// ExitUnioning is called when exiting the unioning production.
	ExitUnioning(c *UnioningContext)

	// ExitIntersectingPart is called when exiting the intersectingPart production.
	ExitIntersectingPart(c *IntersectingPartContext)

	// ExitIntersecting is called when exiting the intersecting production.
	ExitIntersecting(c *IntersectingContext)

	// ExitDifferencingPart is called when exiting the differencingPart production.
	ExitDifferencingPart(c *DifferencingPartContext)

	// ExitDifferencing is called when exiting the differencing production.
	ExitDifferencing(c *DifferencingContext)

	// ExitTypeBody is called when exiting the typeBody production.
	ExitTypeBody(c *TypeBodyContext)

	// ExitFeatureChain is called when exiting the featureChain production.
	ExitFeatureChain(c *FeatureChainContext)

	// ExitClassifier is called when exiting the classifier production.
	ExitClassifier(c *ClassifierContext)

	// ExitSubclassifier is called when exiting the subclassifier production.
	ExitSubclassifier(c *SubclassifierContext)

	// ExitClassifierDeclaration is called when exiting the classifierDeclaration production.
	ExitClassifierDeclaration(c *ClassifierDeclarationContext)

	// ExitSuperclassingPart is called when exiting the superclassingPart production.
	ExitSuperclassingPart(c *SuperclassingPartContext)

	// ExitDataType is called when exiting the dataType production.
	ExitDataType(c *DataTypeContext)

	// ExitClass is called when exiting the class production.
	ExitClass(c *ClassContext)

	// ExitStructure is called when exiting the structure production.
	ExitStructure(c *StructureContext)

	// ExitMetaclass is called when exiting the metaclass production.
	ExitMetaclass(c *MetaclassContext)

	// ExitAssociation is called when exiting the association production.
	ExitAssociation(c *AssociationContext)

	// ExitAssociationStructure is called when exiting the associationStructure production.
	ExitAssociationStructure(c *AssociationStructureContext)

	// ExitInteraction is called when exiting the interaction production.
	ExitInteraction(c *InteractionContext)

	// ExitBehavior is called when exiting the behavior production.
	ExitBehavior(c *BehaviorContext)

	// ExitFunction_ is called when exiting the function_ production.
	ExitFunction_(c *Function_Context)

	// ExitFunctionBody is called when exiting the functionBody production.
	ExitFunctionBody(c *FunctionBodyContext)

	// ExitPredicate is called when exiting the predicate production.
	ExitPredicate(c *PredicateContext)

	// ExitMultiplicity is called when exiting the multiplicity production.
	ExitMultiplicity(c *MultiplicityContext)

	// ExitMultiplicitySubset is called when exiting the multiplicitySubset production.
	ExitMultiplicitySubset(c *MultiplicitySubsetContext)

	// ExitMultiplicityRangeDecl is called when exiting the multiplicityRangeDecl production.
	ExitMultiplicityRangeDecl(c *MultiplicityRangeDeclContext)

	// ExitFeatureSubsetting is called when exiting the featureSubsetting production.
	ExitFeatureSubsetting(c *FeatureSubsettingContext)

	// ExitMultiplicityBounds is called when exiting the multiplicityBounds production.
	ExitMultiplicityBounds(c *MultiplicityBoundsContext)

	// ExitFeature is called when exiting the feature production.
	ExitFeature(c *FeatureContext)

	// ExitEndFeaturePrefix is called when exiting the endFeaturePrefix production.
	ExitEndFeaturePrefix(c *EndFeaturePrefixContext)

	// ExitBasicFeaturePrefix is called when exiting the basicFeaturePrefix production.
	ExitBasicFeaturePrefix(c *BasicFeaturePrefixContext)

	// ExitFeaturePrefix is called when exiting the featurePrefix production.
	ExitFeaturePrefix(c *FeaturePrefixContext)

	// ExitFeatureDeclaration is called when exiting the featureDeclaration production.
	ExitFeatureDeclaration(c *FeatureDeclarationContext)

	// ExitFeatureIdentification is called when exiting the featureIdentification production.
	ExitFeatureIdentification(c *FeatureIdentificationContext)

	// ExitFeatureRelationshipPart is called when exiting the featureRelationshipPart production.
	ExitFeatureRelationshipPart(c *FeatureRelationshipPartContext)

	// ExitChainingPart is called when exiting the chainingPart production.
	ExitChainingPart(c *ChainingPartContext)

	// ExitInvertingPart is called when exiting the invertingPart production.
	ExitInvertingPart(c *InvertingPartContext)

	// ExitOwnedFeatureInverting is called when exiting the ownedFeatureInverting production.
	ExitOwnedFeatureInverting(c *OwnedFeatureInvertingContext)

	// ExitTypeFeaturingPart is called when exiting the typeFeaturingPart production.
	ExitTypeFeaturingPart(c *TypeFeaturingPartContext)

	// ExitOwnedTypeFeaturing is called when exiting the ownedTypeFeaturing production.
	ExitOwnedTypeFeaturing(c *OwnedTypeFeaturingContext)

	// ExitStep is called when exiting the step production.
	ExitStep(c *StepContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitBooleanExpression is called when exiting the booleanExpression production.
	ExitBooleanExpression(c *BooleanExpressionContext)

	// ExitInvariant is called when exiting the invariant production.
	ExitInvariant(c *InvariantContext)

	// ExitConnector is called when exiting the connector production.
	ExitConnector(c *ConnectorContext)

	// ExitConnectorDeclaration is called when exiting the connectorDeclaration production.
	ExitConnectorDeclaration(c *ConnectorDeclarationContext)

	// ExitBinaryConnectorDeclaration is called when exiting the binaryConnectorDeclaration production.
	ExitBinaryConnectorDeclaration(c *BinaryConnectorDeclarationContext)

	// ExitNaryConnectorDeclaration is called when exiting the naryConnectorDeclaration production.
	ExitNaryConnectorDeclaration(c *NaryConnectorDeclarationContext)

	// ExitBindingConnector is called when exiting the bindingConnector production.
	ExitBindingConnector(c *BindingConnectorContext)

	// ExitBindingConnectorDeclaration is called when exiting the bindingConnectorDeclaration production.
	ExitBindingConnectorDeclaration(c *BindingConnectorDeclarationContext)

	// ExitSuccession is called when exiting the succession production.
	ExitSuccession(c *SuccessionContext)

	// ExitSuccessionDeclaration is called when exiting the successionDeclaration production.
	ExitSuccessionDeclaration(c *SuccessionDeclarationContext)

	// ExitKermlFlow is called when exiting the kermlFlow production.
	ExitKermlFlow(c *KermlFlowContext)

	// ExitKermlSuccessionFlow is called when exiting the kermlSuccessionFlow production.
	ExitKermlSuccessionFlow(c *KermlSuccessionFlowContext)

	// ExitKermlFlowDeclaration is called when exiting the kermlFlowDeclaration production.
	ExitKermlFlowDeclaration(c *KermlFlowDeclarationContext)

	// ExitKermlPayloadFeatureMember is called when exiting the kermlPayloadFeatureMember production.
	ExitKermlPayloadFeatureMember(c *KermlPayloadFeatureMemberContext)

	// ExitKermlPayloadFeature is called when exiting the kermlPayloadFeature production.
	ExitKermlPayloadFeature(c *KermlPayloadFeatureContext)

	// ExitKermlFlowEndMember is called when exiting the kermlFlowEndMember production.
	ExitKermlFlowEndMember(c *KermlFlowEndMemberContext)

	// ExitKermlFlowEnd is called when exiting the kermlFlowEnd production.
	ExitKermlFlowEnd(c *KermlFlowEndContext)

	// ExitKermlFlowFeatureMember is called when exiting the kermlFlowFeatureMember production.
	ExitKermlFlowFeatureMember(c *KermlFlowFeatureMemberContext)

	// ExitKermlFlowFeature is called when exiting the kermlFlowFeature production.
	ExitKermlFlowFeature(c *KermlFlowFeatureContext)

	// ExitKermlFlowFeatureRedefinition is called when exiting the kermlFlowFeatureRedefinition production.
	ExitKermlFlowFeatureRedefinition(c *KermlFlowFeatureRedefinitionContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitName is called when exiting the name production.
	ExitName(c *NameContext)

	// ExitKeywordName is called when exiting the keywordName production.
	ExitKeywordName(c *KeywordNameContext)
}
