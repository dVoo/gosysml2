/**
 * SysML v2 Parser Grammar for ANTLR4
 *
 * This parser implements the SysML v2 textual concrete syntax from OMG SysML v2.
 * It includes KerML rules as SysML extends KerML.
 * Compatible with Go target for parser generation.
 *
 * Based on OMG SysML v2 specification Part 2 - Systems Modeling Language (SysML)
 */
parser grammar SysMLv2Parser;

options {
    tokenVocab = SysMLv2Lexer;
}

// =============================================================================
// Entry Point and Root Namespace
// =============================================================================

entryRuleRootNamespace
    : rootNamespace EOF
    ;

rootNamespace
    : packageBodyElement*
    ;

// =============================================================================
// Clause 8.2.2.2 Elements and Relationships Textual Notation
// =============================================================================

identification
    : (LT name GT)? name?
    ;

relationshipBody
    : SEMI
    | LBRACE ownedAnnotation* RBRACE
    ;

// =============================================================================
// Clause 8.2.2.3 Dependencies Textual Notation
// =============================================================================

dependency
    : prefixMetadataAnnotation*
      DEPENDENCY dependencyDeclaration
      relationshipBody
    ;

dependencyDeclaration
    : (identification FROM)?
      qualifiedName (COMMA qualifiedName)* TO
      qualifiedName (COMMA qualifiedName)*
    ;

// =============================================================================
// Clause 8.2.2.4 Annotations Textual Notation
// =============================================================================

annotation
    : qualifiedName
    ;

ownedAnnotation
    : annotatingElement
    ;

annotatingMember
    : annotatingElement
    ;

annotatingElement
    : comment
    | documentation
    | textualRepresentation
    | metadataFeature
    ;

comment
    : (COMMENT identification
       (ABOUT annotation (COMMA annotation)*)?
      )?
      (LOCALE STRING_VALUE)?
      REGULAR_COMMENT
    ;

documentation
    : DOC identification
      (LOCALE STRING_VALUE)?
      REGULAR_COMMENT
    ;

textualRepresentation
    : (REP identification)?
      LANGUAGE STRING_VALUE REGULAR_COMMENT
    ;

// =============================================================================
// Clause 8.2.2.5 Namespaces and Packages Textual Notation
// =============================================================================

package
    : prefixMetadataMember*
      packageDeclaration packageBody
    ;

libraryPackage
    : STANDARD? LIBRARY
      prefixMetadataMember*
      packageDeclaration packageBody
    ;

packageDeclaration
    : PACKAGE identification
    ;

packageBody
    : SEMI
    | LBRACE packageBodyElement* RBRACE
    ;

packageBodyElement
    : packageMember
    | elementFilterMember
    | aliasMember
    | import_
    ;

memberPrefix
    : visibilityIndicator?
    ;

packageMember
    : memberPrefix
      (definitionElement | usageElement)
    ;

elementFilterMember
    : memberPrefix
      FILTER ownedExpression SEMI
    ;

aliasMember
    : memberPrefix
      ALIAS (LT name GT)?
      name?
      FOR qualifiedName
      relationshipBody
    ;

import_
    : visibilityIndicator
      IMPORT ALL?
      importDeclaration
      relationshipBody
    ;

// Fixed: Broken the indirect recursion by inlining filterPackage
importDeclaration
    : membershipImport
    | namespaceImport
    ;

membershipImport
    : qualifiedName (COLONCOLON STARSTAR)?
    ;

namespaceImport
    : qualifiedName COLONCOLON STAR (COLONCOLON STARSTAR)?
    | filterPackage
    ;

// Fixed: filterPackage now uses membershipImport/namespaceImport directly without going through importDeclaration
filterPackage
    : filterPackageImportPart filterPackageMember+
    ;

filterPackageImportPart
    : qualifiedName (COLONCOLON STARSTAR)?                    // membershipImport
    | qualifiedName COLONCOLON STAR (COLONCOLON STARSTAR)?    // namespaceImport (non-recursive part)
    ;

filterPackageMember
    : LBRACKET ownedExpression RBRACKET
    ;

visibilityIndicator
    : PUBLIC
    | PRIVATE
    | PROTECTED
    ;

// =============================================================================
// Clause 8.2.2.5.2 Package Elements
// =============================================================================

definitionElement
    : package
    | libraryPackage
    | annotatingElement
    | dependency
    | attributeDefinition
    | enumerationDefinition
    | occurrenceDefinition
    | individualDefinition
    | itemDefinition
    | partDefinition
    | connectionDefinition
    | flowDefinition
    | allocationDefinition
    | interfaceDefinition
    | portDefinition
    | actionDefinition
    | calculationDefinition
    | stateDefinition
    | constraintDefinition
    | requirementDefinition
    | concernDefinition
    | caseDefinition
    | analysisCaseDefinition
    | verificationCaseDefinition
    | useCaseDefinition
    | viewDefinition
    | viewpointDefinition
    | renderingDefinition
    | metadataDefinition
    | extendedDefinition
    // KerML declarations used by .kerml standard library files.
    | namespace_
    | type_
    | classifier
    | dataType
    | class
    | structure
    | metaclass
    | association
    | associationStructure
    | interaction
    | behavior
    | function_
    | predicate
    | multiplicity
    | featureSubsetting
    | step
    | expression
    | booleanExpression
    | invariant
    | feature
    | connector
    | bindingConnector
    | succession
    | kermlFlow
    | kermlSuccessionFlow
    ;

usageElement
    : nonOccurrenceUsageElement
    | occurrenceUsageElement
    ;

// =============================================================================
// Clause 8.2.2.6 Definition and Usage Textual Notation
// =============================================================================

// Definitions

basicDefinitionPrefix
    : ABSTRACT
    | VARIATION
    ;

definitionExtensionKeyword
    : prefixMetadataMember
    ;

definitionPrefix
    : basicDefinitionPrefix? definitionExtensionKeyword*
    ;

definition
    : definitionDeclaration definitionBody
    ;

definitionDeclaration
    : identification subclassificationPart?
    ;

definitionBody
    : SEMI
    | LBRACE definitionBodyItem* RBRACE
    ;

definitionBodyItem
    : definitionMember
    | variantUsageMember
    | nonOccurrenceUsageMember
    | sourceSuccessionMember? occurrenceUsageMember
    | endFeatureMember
    | aliasMember
    | import_
    ;

// Added: Support for standalone 'end' declarations like 'end end1;'
endFeatureMember
    : memberPrefix END endFeatureDeclaration usageBody
    ;

endFeatureDeclaration
    : ownedCrossFeatureMember? usageDeclaration
      (referencesKw ownedReferenceSubsetting)?
    ;

definitionMember
    : memberPrefix definitionElement
    ;

variantUsageMember
    : memberPrefix VARIANT variantUsageElement
    ;

nonOccurrenceUsageMember
    : memberPrefix nonOccurrenceUsageElement
    ;

occurrenceUsageMember
    : memberPrefix occurrenceUsageElement
    ;

structureUsageMember
    : memberPrefix structureUsageElement
    ;

behaviorUsageMember
    : memberPrefix behaviorUsageElement
    ;

// Usages

featureDirection
    : IN
    | OUT
    | INOUT
    ;

refPrefix
    : featureDirection?
      DERIVED?
      (ABSTRACT | VARIATION)?
      CONSTANT?
    ;

basicUsagePrefix
    : refPrefix
      REF?
    ;

endUsagePrefix
    : END ownedCrossFeatureMember?
    ;

ownedCrossFeatureMember
    : ownedCrossFeature
    ;

ownedCrossFeature
    : basicUsagePrefix usageDeclaration
    ;

usageExtensionKeyword
    : prefixMetadataMember
    ;

unextendedUsagePrefix
    : endUsagePrefix
    | basicUsagePrefix
    ;

usagePrefix
    : unextendedUsagePrefix usageExtensionKeyword*
    ;

usage
    : usageDeclaration usageCompletion
    ;

usageDeclaration
    : identification featureSpecializationPart?
    ;

usageCompletion
    : valuePart? usageBody
    ;

usageBody
    : definitionBody
    ;

valuePart
    : featureValue
    ;

featureValue
    : (EQ | COLONEQ | DEFAULT (EQ | COLONEQ)?)
      ownedExpression
    ;

// =============================================================================
// Clause 8.2.2.6.3 Reference Usages
// =============================================================================

defaultReferenceUsage
    : refPrefix usage
    ;

referenceUsage
    : (endUsagePrefix | refPrefix)
      REF usage
    ;

variantReference
    : ownedReferenceSubsetting
      featureSpecialization* usageBody
    ;

// =============================================================================
// Clause 8.2.2.6.4 Body Elements
// =============================================================================

nonOccurrenceUsageElement
    : defaultReferenceUsage
    | referenceUsage
    | attributeUsage
    | enumerationUsage
    | bindingConnectorAsUsage
    | successionAsUsage
    | extendedUsage
    ;

occurrenceUsageElement
    : structureUsageElement
    | behaviorUsageElement
    ;

structureUsageElement
    : occurrenceUsage
    | individualUsage
    | portionUsage
    | eventOccurrenceUsage
    | itemUsage
    | partUsage
    | viewUsage
    | renderingUsage
    | portUsage
    | connectionUsage
    | interfaceUsage
    | allocationUsage
    | message
    | flowUsage
    | successionFlowUsage
    ;

behaviorUsageElement
    : actionUsage
    | calculationUsage
    | stateUsage
    | constraintUsage
    | requirementUsage
    | concernUsage
    | caseUsage
    | analysisCaseUsage
    | verificationCaseUsage
    | useCaseUsage
    | viewpointUsage
    | performActionUsage
    | exhibitStateUsage
    | includeUseCaseUsage
    | assertConstraintUsage
    | satisfyRequirementUsage
    ;

variantUsageElement
    : variantReference
    | referenceUsage
    | attributeUsage
    | bindingConnectorAsUsage
    | successionAsUsage
    | occurrenceUsage
    | individualUsage
    | portionUsage
    | eventOccurrenceUsage
    | itemUsage
    | partUsage
    | viewUsage
    | renderingUsage
    | portUsage
    | connectionUsage
    | interfaceUsage
    | allocationUsage
    | message
    | flowUsage
    | successionFlowUsage
    | behaviorUsageElement
    ;

// =============================================================================
// Clause 8.2.2.6.5 Specialization
// =============================================================================

subclassificationPart
    : specializes ownedSubclassification (COMMA ownedSubclassification)*
    ;

ownedSubclassification
    : qualifiedName
    ;

featureSpecializationPart
    : featureSpecialization+ multiplicityPart? featureSpecialization*
    | multiplicityPart featureSpecialization*
    ;

featureSpecialization
    : typings
    | subsettings
    | references
    | crosses
    | redefinitions
    ;

typings
    : typedBy (COMMA ownedFeatureTyping)*
    ;

typedBy
    : definedBy ownedFeatureTyping
    ;

ownedFeatureTyping
    : qualifiedName
    | ownedFeatureChain
    | conjugatedPortTyping
    ;

subsettings
    : subsetsKw ownedSubsetting (COMMA ownedSubsetting)*
    ;

ownedSubsetting
    : qualifiedName (DOT ownedFeatureChaining)*
    ;

references
    : referencesKw ownedReferenceSubsetting
    ;

ownedReferenceSubsetting
    : qualifiedName
    | ownedFeatureChain
    ;

crosses
    : crossesKw ownedCrossSubsetting
    ;

ownedCrossSubsetting
    : qualifiedName
    | ownedFeatureChain
    ;

redefinitions
    : redefinesKw ownedRedefinition (COMMA ownedRedefinition)*
    ;

ownedRedefinition
    : qualifiedName
    | ownedFeatureChain
    ;

ownedFeatureChain
    : ownedFeatureChaining (DOT ownedFeatureChaining)+
    ;

ownedFeatureChaining
    : qualifiedName
    ;

// Keyword alternatives
specializes     : COLONGT | SPECIALIZES;
definedBy       : COLON | DEFINED BY;
subsetsKw       : COLONGT | SUBSETS;
referencesKw    : COLONCOLONGT | REFERENCES;
crossesKw       : EQGT | CROSSES;
redefinesKw     : COLONGTGT | REDEFINES;

// =============================================================================
// Clause 8.2.2.6.6 Multiplicity
// =============================================================================

// Fixed: Combined alternatives to properly handle '[*] nonunique' patterns
multiplicityPart
    : ownedMultiplicity
      (ORDERED NONUNIQUE? | NONUNIQUE ORDERED?)?
    | (ORDERED NONUNIQUE? | NONUNIQUE ORDERED?)
    ;

ownedMultiplicity
    : multiplicityRange
    ;

multiplicityRange
    : LBRACKET (multiplicityExpressionMember DOTDOT)?
      multiplicityExpressionMember RBRACKET
    ;

multiplicityExpressionMember
    : literalExpression
    | featureReferenceExpression
    ;

// =============================================================================
// Clause 8.2.2.7 Attributes Textual Notation
// =============================================================================

attributeDefinition
    : definitionPrefix ATTRIBUTE DEF definition
    ;

attributeUsage
    : usagePrefix ATTRIBUTE usage
    ;

// =============================================================================
// Clause 8.2.2.8 Enumerations Textual Notation
// =============================================================================

enumerationDefinition
    : definitionExtensionKeyword*
      ENUM DEF definitionDeclaration enumerationBody
    ;

enumerationBody
    : SEMI
    | LBRACE (annotatingMember | enumerationUsageMember)* RBRACE
    ;

enumerationUsageMember
    : memberPrefix enumeratedValue
    ;

// Fixed: Added prefix metadata support for '#Security enum ...' patterns
enumeratedValue
    : usageExtensionKeyword* ENUM? usage
    ;

enumerationUsage
    : usagePrefix ENUM usage
    ;

// =============================================================================
// Clause 8.2.2.9 Occurrences Textual Notation
// =============================================================================

occurrenceDefinitionPrefix
    : basicDefinitionPrefix?
      (INDIVIDUAL)?
      definitionExtensionKeyword*
    ;

occurrenceDefinition
    : occurrenceDefinitionPrefix OCCURRENCE DEF definition
    ;

individualDefinition
    : basicDefinitionPrefix? INDIVIDUAL
      definitionExtensionKeyword* DEF definition
    ;

// Fixed: Added optional END prefix to support 'end port', 'end item', etc.
occurrenceUsagePrefix
    : (END ownedCrossFeatureMember?)?
      basicUsagePrefix
      INDIVIDUAL?
      portionKind?
      usageExtensionKeyword*
    ;

occurrenceUsage
    : occurrenceUsagePrefix OCCURRENCE usage
    ;

individualUsage
    : basicUsagePrefix INDIVIDUAL
      usageExtensionKeyword* usage
    ;

portionUsage
    : basicUsagePrefix INDIVIDUAL?
      portionKind
      usageExtensionKeyword* usage
    ;

portionKind
    : SNAPSHOT
    | TIMESLICE
    ;

eventOccurrenceUsage
    : occurrenceUsagePrefix EVENT
      (ownedReferenceSubsetting featureSpecializationPart? usageCompletion
      | OCCURRENCE usageDeclaration? usageCompletion
      )
    ;

// Occurrence Successions

sourceSuccessionMember
    : THEN sourceSuccession
    ;

sourceSuccession
    : sourceEndMember
    ;

sourceEndMember
    : sourceEnd
    ;

sourceEnd
    : ownedMultiplicity?
    ;

// =============================================================================
// Clause 8.2.2.10 Items Textual Notation
// =============================================================================

itemDefinition
    : occurrenceDefinitionPrefix ITEM DEF definition
    ;

itemUsage
    : occurrenceUsagePrefix ITEM usage
    ;

// =============================================================================
// Clause 8.2.2.11 Parts Textual Notation
// =============================================================================

partDefinition
    : occurrenceDefinitionPrefix PART DEF definition
    ;

partUsage
    : occurrenceUsagePrefix PART usage
    ;

// =============================================================================
// Clause 8.2.2.12 Ports Textual Notation
// =============================================================================

portDefinition
    : definitionPrefix PORT DEF definition
    ;

portUsage
    : occurrenceUsagePrefix PORT usage
    ;

conjugatedPortTyping
    : TILDE qualifiedName
    ;

// =============================================================================
// Clause 8.2.2.13 Connections Textual Notation
// =============================================================================

connectionDefinition
    : occurrenceDefinitionPrefix CONNECTION DEF definition
    ;

connectionUsage
    : occurrenceUsagePrefix
      (CONNECTION usageDeclaration valuePart?
        (CONNECT connectorPart)?
      | CONNECT connectorPart
      )
      usageBody
    ;

connectorPart
    : binaryConnectorPart
    | naryConnectorPart
    ;

binaryConnectorPart
    : connectorEndMember TO connectorEndMember
    ;

naryConnectorPart
    : LPAREN connectorEndMember COMMA
      connectorEndMember
      (COMMA connectorEndMember)* RPAREN
    ;

connectorEndMember
    : connectorEnd
    ;

connectorEnd
    : ownedCrossMultiplicityMember?
      (name referencesKw)?
      ownedReferenceSubsetting
    ;

ownedCrossMultiplicityMember
    : ownedCrossMultiplicity
    ;

ownedCrossMultiplicity
    : ownedMultiplicity
    ;

// Binding Connectors

bindingConnectorAsUsage
    : usagePrefix (BINDING usageDeclaration)?
      BIND connectorEndMember
      EQ connectorEndMember
      usageBody
    ;

// Successions

successionAsUsage
    : usagePrefix (SUCCESSION usageDeclaration)?
      FIRST connectorEndMember
      THEN connectorEndMember
      usageBody
    ;

// =============================================================================
// Clause 8.2.2.14 Interfaces Textual Notation
// =============================================================================

interfaceDefinition
    : occurrenceDefinitionPrefix INTERFACE DEF
      definitionDeclaration interfaceBody
    ;

interfaceBody
    : SEMI
    | LBRACE interfaceBodyItem* RBRACE
    ;

interfaceBodyItem
    : definitionMember
    | variantUsageMember
    | interfaceNonOccurrenceUsageMember
    | sourceSuccessionMember? interfaceOccurrenceUsageMember
    | aliasMember
    | import_
    ;

interfaceNonOccurrenceUsageMember
    : memberPrefix interfaceNonOccurrenceUsageElement
    ;

interfaceNonOccurrenceUsageElement
    : referenceUsage
    | attributeUsage
    | enumerationUsage
    | bindingConnectorAsUsage
    | successionAsUsage
    ;

interfaceOccurrenceUsageMember
    : memberPrefix interfaceOccurrenceUsageElement
    ;

interfaceOccurrenceUsageElement
    : defaultInterfaceEnd
    | structureUsageElement
    | behaviorUsageElement
    ;

defaultInterfaceEnd
    : END usage
    ;

interfaceUsage
    : occurrenceUsagePrefix INTERFACE
      interfaceUsageDeclaration interfaceBody
    ;

interfaceUsageDeclaration
    : usageDeclaration valuePart? (CONNECT interfacePart)?
    | interfacePart
    ;

interfacePart
    : binaryInterfacePart
    | naryInterfacePart
    ;

binaryInterfacePart
    : interfaceEndMember TO interfaceEndMember
    ;

naryInterfacePart
    : LPAREN interfaceEndMember COMMA
      interfaceEndMember
      (COMMA interfaceEndMember)* RPAREN
    ;

interfaceEndMember
    : interfaceEnd
    ;

interfaceEnd
    : ownedCrossMultiplicityMember?
      (name referencesKw)?
      ownedReferenceSubsetting
    ;

// =============================================================================
// Clause 8.2.2.15 Allocations Textual Notation
// =============================================================================

allocationDefinition
    : occurrenceDefinitionPrefix ALLOCATION DEF definition
    ;

allocationUsage
    : occurrenceUsagePrefix
      allocationUsageDeclaration usageBody
    ;

allocationUsageDeclaration
    : ALLOCATION usageDeclaration (ALLOCATE connectorPart)?
    | ALLOCATE connectorPart
    ;

// =============================================================================
// Clause 8.2.2.16 Flows Textual Notation
// =============================================================================

flowDefinition
    : occurrenceDefinitionPrefix FLOW DEF definition
    ;

message
    : occurrenceUsagePrefix MESSAGE
      messageDeclaration definitionBody
    ;

messageDeclaration
    : usageDeclaration valuePart?
      (OF flowPayloadFeatureMember)?
      (FROM messageEventMember TO messageEventMember)?
    | messageEventMember TO messageEventMember
    ;

messageEventMember
    : messageEvent
    ;

messageEvent
    : ownedReferenceSubsetting
    ;

flowUsage
    : occurrenceUsagePrefix FLOW
      flowDeclaration definitionBody
    ;

successionFlowUsage
    : occurrenceUsagePrefix SUCCESSION FLOW
      flowDeclaration definitionBody
    ;

flowDeclaration
    : usageDeclaration valuePart?
      (OF flowPayloadFeatureMember)?
      (FROM flowEndMember TO flowEndMember)?
    | flowEndMember TO flowEndMember
    ;

flowPayloadFeatureMember
    : flowPayloadFeature
    ;

flowPayloadFeature
    : payloadFeature
    ;

payloadFeature
    : identification payloadFeatureSpecializationPart valuePart?
    | ownedFeatureTyping ownedMultiplicity?
    | ownedMultiplicity ownedFeatureTyping
    ;

payloadFeatureSpecializationPart
    : featureSpecialization+ multiplicityPart? featureSpecialization*
    | multiplicityPart featureSpecialization+
    ;

flowEndMember
    : flowEnd
    ;

// Fixed: flowEnd now directly handles dotted feature paths like A2.a or b.f.a
flowEnd
    : qualifiedName (DOT qualifiedName)*
    ;

flowEndSubsetting
    : qualifiedName
    | featureChainPrefix
    ;

featureChainPrefix
    : (ownedFeatureChaining DOT)+
      ownedFeatureChaining DOT
    ;

flowFeatureMember
    : flowFeature
    ;

flowFeature
    : flowFeatureRedefinition
    ;

flowFeatureRedefinition
    : qualifiedName
    ;

// =============================================================================
// Clause 8.2.2.17 Actions Textual Notation
// =============================================================================

actionDefinition
    : occurrenceDefinitionPrefix ACTION DEF
      definitionDeclaration actionBody
    ;

actionBody
    : SEMI
    | LBRACE actionBodyItem* RBRACE
    ;

actionBodyItem
    : nonBehaviorBodyItem
    | initialNodeMember actionTargetSuccessionMember*
    | sourceSuccessionMember? actionBehaviorMember actionTargetSuccessionMember*
    | guardedSuccessionMember
    ;

nonBehaviorBodyItem
    : import_
    | aliasMember
    | definitionMember
    | variantUsageMember
    | nonOccurrenceUsageMember
    | sourceSuccessionMember? structureUsageMember
    ;

actionBehaviorMember
    : behaviorUsageMember
    | actionNodeMember
    ;

initialNodeMember
    : memberPrefix FIRST qualifiedName relationshipBody
    ;

actionNodeMember
    : memberPrefix actionNode
    ;

actionTargetSuccessionMember
    : memberPrefix actionTargetSuccession
    ;

guardedSuccessionMember
    : memberPrefix guardedSuccession
    ;

// Action Usages

actionUsage
    : occurrenceUsagePrefix ACTION
      actionUsageDeclaration actionBody
    ;

actionUsageDeclaration
    : usageDeclaration valuePart?
    ;

performActionUsage
    : occurrenceUsagePrefix PERFORM
      performActionUsageDeclaration actionBody
    ;

performActionUsageDeclaration
    : (ownedReferenceSubsetting featureSpecializationPart?
      | ACTION usageDeclaration
      )
      valuePart?
    ;

actionNode
    : controlNode
    | sendNode
    | acceptNode
    | assignmentNode
    | terminateNode
    | ifNode
    | whileLoopNode
    | forLoopNode
    ;

actionNodeUsageDeclaration
    : ACTION usageDeclaration?
    ;

actionNodePrefix
    : occurrenceUsagePrefix actionNodeUsageDeclaration?
    ;

// Control Nodes

controlNode
    : mergeNode
    | decisionNode
    | joinNode
    | forkNode
    ;

controlNodePrefix
    : refPrefix
      INDIVIDUAL?
      portionKind?
      usageExtensionKeyword*
    ;

mergeNode
    : controlNodePrefix
      MERGE usageDeclaration
      actionBody
    ;

decisionNode
    : controlNodePrefix
      DECIDE usageDeclaration
      actionBody
    ;

joinNode
    : controlNodePrefix
      JOIN usageDeclaration
      actionBody
    ;

forkNode
    : controlNodePrefix
      FORK usageDeclaration
      actionBody
    ;

// Send and Accept Action Usages

acceptNode
    : occurrenceUsagePrefix
      acceptNodeDeclaration actionBody
    ;

acceptNodeDeclaration
    : actionNodeUsageDeclaration?
      ACCEPT acceptParameterPart
    ;

acceptParameterPart
    : payloadParameterMember
      (VIA nodeParameterMember)?
    ;

payloadParameterMember
    : payloadParameter
    ;

payloadParameter
    : payloadFeature
    | identification payloadFeatureSpecializationPart? triggerValuePart
    ;

triggerValuePart
    : triggerFeatureValue
    ;

triggerFeatureValue
    : triggerExpression
    ;

triggerExpression
    : (AT | AFTER) ownedExpression
    | WHEN ownedExpression
    ;

// Fixed: Use actionNodePrefix to support 'action snd send' and made payload optional
sendNode
    : actionNodePrefix SEND
      (nodeParameterMember? senderReceiverPart)?
      actionBody
    ;

sendNodeDeclaration
    : actionNodeUsageDeclaration? SEND
      nodeParameterMember senderReceiverPart?
    ;

senderReceiverPart
    : VIA nodeParameterMember (TO nodeParameterMember)?
    | TO nodeParameterMember
    ;

nodeParameterMember
    : nodeParameter
    ;

nodeParameter
    : ownedExpression
    ;

// Assignment Action Usages

assignmentNode
    : occurrenceUsagePrefix
      assignmentNodeDeclaration actionBody
    ;

assignmentNodeDeclaration
    : actionNodeUsageDeclaration? ASSIGN
      assignmentTargetMember
      featureChainMember COLONEQ
      nodeParameterMember
    ;

assignmentTargetMember
    : assignmentTargetParameter
    ;

assignmentTargetParameter
    : (baseExpression DOT)?
    ;

featureChainMember
    : qualifiedName
    | ownedFeatureChainMember
    ;

ownedFeatureChainMember
    : ownedFeatureChain
    ;

// Terminate Action Usages

terminateNode
    : occurrenceUsagePrefix actionNodeUsageDeclaration?
      TERMINATE nodeParameterMember?
      actionBody
    ;

// Structured Control Action Usages

ifNode
    : actionNodePrefix
      IF ownedExpression
      actionBodyParameter
      (ELSE (actionBodyParameter | ifNode))?
    ;

actionBodyParameter
    : (ACTION usageDeclaration?)?
      LBRACE actionBodyItem* RBRACE
    ;

whileLoopNode
    : actionNodePrefix
      (WHILE ownedExpression | LOOP)
      actionBodyParameter
      (UNTIL ownedExpression SEMI)?
    ;

forLoopNode
    : actionNodePrefix
      FOR forVariableDeclaration
      IN ownedExpression
      actionBodyParameter
    ;

forVariableDeclaration
    : usageDeclaration
    ;

// Action Successions

actionTargetSuccession
    : (targetSuccession | guardedTargetSuccession | defaultTargetSuccession)
      usageBody
    ;

targetSuccession
    : sourceEndMember THEN connectorEndMember
    ;

guardedTargetSuccession
    : guardExpressionMember THEN transitionSuccessionMember
    ;

defaultTargetSuccession
    : ELSE transitionSuccessionMember
    ;

guardedSuccession
    : (SUCCESSION usageDeclaration)?
      FIRST featureChainMember
      guardExpressionMember
      THEN transitionSuccessionMember
      usageBody
    ;

// =============================================================================
// Clause 8.2.2.18 States Textual Notation
// =============================================================================

stateDefinition
    : occurrenceDefinitionPrefix STATE DEF
      definitionDeclaration stateDefBody
    ;

stateDefBody
    : SEMI
    | PARALLEL? LBRACE stateBodyItem* RBRACE
    ;

stateBodyItem
    : nonBehaviorBodyItem
    | sourceSuccessionMember? behaviorUsageMember targetTransitionUsageMember*
    | transitionUsageMember
    | entryActionMember entryTransitionMember*
    | doActionMember
    | exitActionMember
    ;

entryActionMember
    : memberPrefix ENTRY stateActionUsage
    ;

doActionMember
    : memberPrefix DO stateActionUsage
    ;

exitActionMember
    : memberPrefix EXIT stateActionUsage
    ;

// Fixed: entryTransitionMember should use connectorEndMember directly after THEN
// since 'entry; then S1;' has only one THEN keyword
entryTransitionMember
    : memberPrefix
      (guardedTargetSuccession | THEN connectorEndMember)
      SEMI
    ;

stateActionUsage
    : SEMI
    | statePerformActionUsage
    | stateAcceptActionUsage
    | stateSendActionUsage
    | stateAssignmentActionUsage
    ;

statePerformActionUsage
    : performActionUsageDeclaration actionBody
    ;

stateAcceptActionUsage
    : acceptNodeDeclaration actionBody
    ;

stateSendActionUsage
    : sendNodeDeclaration actionBody
    ;

stateAssignmentActionUsage
    : assignmentNodeDeclaration actionBody
    ;

transitionUsageMember
    : memberPrefix transitionUsage
    ;

targetTransitionUsageMember
    : memberPrefix targetTransitionUsage
    ;

// State Usages

stateUsage
    : occurrenceUsagePrefix STATE
      actionUsageDeclaration stateUsageBody
    ;

stateUsageBody
    : SEMI
    | PARALLEL? LBRACE stateBodyItem* RBRACE
    ;

exhibitStateUsage
    : occurrenceUsagePrefix EXHIBIT
      (ownedReferenceSubsetting featureSpecializationPart?
      | STATE usageDeclaration
      )
      valuePart? stateUsageBody
    ;

// Transition Usages

transitionUsage
    : TRANSITION (usageDeclaration FIRST)?
      featureChainMember
      (triggerActionMember)?
      guardExpressionMember?
      effectBehaviorMember?
      THEN transitionSuccessionMember
      actionBody
    ;

targetTransitionUsage
    : (TRANSITION
        triggerActionMember?
        guardExpressionMember?
        effectBehaviorMember?
      | triggerActionMember
        guardExpressionMember?
        effectBehaviorMember?
      | guardExpressionMember
        effectBehaviorMember?
      )?
      THEN transitionSuccessionMember
      actionBody
    ;

triggerActionMember
    : ACCEPT triggerAction
    ;

triggerAction
    : acceptParameterPart
    ;

guardExpressionMember
    : IF ownedExpression
    ;

effectBehaviorMember
    : DO effectBehaviorUsage
    ;

effectBehaviorUsage
    : transitionPerformActionUsage
    | transitionAcceptActionUsage
    | transitionSendActionUsage
    | transitionAssignmentActionUsage
    ;

transitionPerformActionUsage
    : performActionUsageDeclaration (LBRACE actionBodyItem* RBRACE)?
    ;

transitionAcceptActionUsage
    : acceptNodeDeclaration (LBRACE actionBodyItem* RBRACE)?
    ;

transitionSendActionUsage
    : sendNodeDeclaration (LBRACE actionBodyItem* RBRACE)?
    ;

transitionAssignmentActionUsage
    : assignmentNodeDeclaration (LBRACE actionBodyItem* RBRACE)?
    ;

transitionSuccessionMember
    : transitionSuccession
    ;

transitionSuccession
    : connectorEndMember
    ;

// =============================================================================
// Clause 8.2.2.19 Calculations Textual Notation
// =============================================================================

calculationDefinition
    : occurrenceDefinitionPrefix CALC DEF
      definitionDeclaration calculationBody
    ;

calculationUsage
    : occurrenceUsagePrefix CALC
      actionUsageDeclaration calculationBody
    ;

calculationBody
    : SEMI
    | LBRACE calculationBodyPart RBRACE
    ;

calculationBodyPart
    : calculationBodyItem*
      resultExpressionMember?
    ;

calculationBodyItem
    : actionBodyItem
    | returnParameterMember
    ;

returnParameterMember
    : memberPrefix? RETURN usageElement
    ;

resultExpressionMember
    : memberPrefix? ownedExpression
    ;

// =============================================================================
// Clause 8.2.2.20 Constraints Textual Notation
// =============================================================================

constraintDefinition
    : occurrenceDefinitionPrefix CONSTRAINT DEF
      definitionDeclaration calculationBody
    ;

constraintUsage
    : occurrenceUsagePrefix CONSTRAINT
      constraintUsageDeclaration calculationBody
    ;

assertConstraintUsage
    : occurrenceUsagePrefix ASSERT NOT?
      (ownedReferenceSubsetting featureSpecializationPart?
      | CONSTRAINT constraintUsageDeclaration
      )
      calculationBody
    ;

constraintUsageDeclaration
    : usageDeclaration valuePart?
    ;

// =============================================================================
// Clause 8.2.2.21 Requirements Textual Notation
// =============================================================================

requirementDefinition
    : occurrenceDefinitionPrefix REQUIREMENT DEF
      definitionDeclaration requirementBody
    ;

requirementBody
    : SEMI
    | LBRACE requirementBodyItem* RBRACE
    ;

requirementBodyItem
    : definitionBodyItem
    | subjectMember
    | requirementConstraintMember
    | framedConcernMember
    | requirementVerificationMember
    | actorMember
    | stakeholderMember
    ;

subjectMember
    : memberPrefix subjectUsage
    ;

subjectUsage
    : SUBJECT usageExtensionKeyword* usage
    ;

requirementConstraintMember
    : memberPrefix? requirementKind requirementConstraintUsage
    ;

requirementKind
    : ASSUME
    | REQUIRE
    ;

requirementConstraintUsage
    : ownedReferenceSubsetting featureSpecializationPart? requirementBody
    | (usageExtensionKeyword* CONSTRAINT | usageExtensionKeyword+)
      constraintUsageDeclaration calculationBody
    ;

framedConcernMember
    : memberPrefix? FRAME framedConcernUsage
    ;

framedConcernUsage
    : ownedReferenceSubsetting featureSpecializationPart? calculationBody
    | (usageExtensionKeyword* CONCERN | usageExtensionKeyword+)
      calculationUsageDeclaration calculationBody
    ;

calculationUsageDeclaration
    : usageDeclaration valuePart?
    ;

actorMember
    : memberPrefix actorUsage
    ;

actorUsage
    : ACTOR usageExtensionKeyword* usage
    ;

stakeholderMember
    : memberPrefix stakeholderUsage
    ;

stakeholderUsage
    : STAKEHOLDER usageExtensionKeyword* usage
    ;

// Requirement Usages

requirementUsage
    : occurrenceUsagePrefix REQUIREMENT
      constraintUsageDeclaration requirementBody
    ;

// Fixed: Made ASSERT optional to support 'satisfy r by p;' and 'not satisfy r by p;'
satisfyRequirementUsage
    : occurrenceUsagePrefix ASSERT? NOT? SATISFY
      (ownedReferenceSubsetting featureSpecializationPart?
      | REQUIREMENT usageDeclaration
      )
      valuePart?
      (BY satisfactionSubjectMember)?
      requirementBody
    ;

satisfactionSubjectMember
    : satisfactionParameter
    ;

satisfactionParameter
    : featureChainMember
    ;

// Concerns

concernDefinition
    : occurrenceDefinitionPrefix CONCERN DEF
      definitionDeclaration requirementBody
    ;

concernUsage
    : occurrenceUsagePrefix CONCERN
      constraintUsageDeclaration requirementBody
    ;

// =============================================================================
// Clause 8.2.2.22 Cases Textual Notation
// =============================================================================

caseDefinition
    : occurrenceDefinitionPrefix CASE DEF
      definitionDeclaration caseBody
    ;

caseUsage
    : occurrenceUsagePrefix CASE
      constraintUsageDeclaration caseBody
    ;

caseBody
    : SEMI
    | LBRACE caseBodyItem* resultExpressionMember? RBRACE
    ;

caseBodyItem
    : actionBodyItem
    | subjectMember
    | actorMember
    | objectiveMember
    | returnParameterMember
    ;

objectiveMember
    : memberPrefix OBJECTIVE objectiveRequirementUsage
    ;

objectiveRequirementUsage
    : usageExtensionKeyword* constraintUsageDeclaration requirementBody
    ;

// =============================================================================
// Clause 8.2.2.23 Analysis Cases Textual Notation
// =============================================================================

analysisCaseDefinition
    : occurrenceDefinitionPrefix ANALYSIS DEF
      definitionDeclaration caseBody
    ;

analysisCaseUsage
    : occurrenceUsagePrefix ANALYSIS
      constraintUsageDeclaration caseBody
    ;

// =============================================================================
// Clause 8.2.2.24 Verification Cases Textual Notation
// =============================================================================

verificationCaseDefinition
    : occurrenceDefinitionPrefix VERIFICATION DEF
      definitionDeclaration caseBody
    ;

verificationCaseUsage
    : occurrenceUsagePrefix VERIFICATION
      constraintUsageDeclaration caseBody
    ;

requirementVerificationMember
    : memberPrefix VERIFY requirementVerificationUsage
    ;

requirementVerificationUsage
    : ownedReferenceSubsetting featureSpecialization* requirementBody
    | (usageExtensionKeyword* REQUIREMENT | usageExtensionKeyword+)
      constraintUsageDeclaration requirementBody
    ;

// =============================================================================
// Clause 8.2.2.25 Use Cases Textual Notation
// =============================================================================

useCaseDefinition
    : occurrenceDefinitionPrefix USE CASE DEF
      definitionDeclaration caseBody
    ;

useCaseUsage
    : occurrenceUsagePrefix USE CASE
      constraintUsageDeclaration caseBody
    ;

includeUseCaseUsage
    : occurrenceUsagePrefix INCLUDE
      (ownedReferenceSubsetting featureSpecializationPart?
      | USE CASE usageDeclaration
      )
      valuePart?
      caseBody
    ;

// =============================================================================
// Clause 8.2.2.26 Views and Viewpoints Textual Notation
// =============================================================================

viewDefinition
    : occurrenceDefinitionPrefix VIEW DEF
      definitionDeclaration viewDefinitionBody
    ;

viewDefinitionBody
    : SEMI
    | LBRACE viewDefinitionBodyItem* RBRACE
    ;

viewDefinitionBodyItem
    : definitionBodyItem
    | elementFilterMember
    | viewRenderingMember
    ;

viewRenderingMember
    : memberPrefix RENDER viewRenderingUsage
    ;

viewRenderingUsage
    : ownedReferenceSubsetting featureSpecializationPart? usageBody
    | (usageExtensionKeyword* RENDERING | usageExtensionKeyword+) usage
    ;

// View Usages

viewUsage
    : occurrenceUsagePrefix VIEW
      usageDeclaration? valuePart?
      viewBody
    ;

viewBody
    : SEMI
    | LBRACE viewBodyItem* RBRACE
    ;

viewBodyItem
    : definitionBodyItem
    | elementFilterMember
    | viewRenderingMember
    | expose
    ;

expose
    : EXPOSE (membershipExpose | namespaceExpose)
      relationshipBody
    ;

membershipExpose
    : membershipImport
    ;

namespaceExpose
    : namespaceImport
    ;

// Viewpoints

viewpointDefinition
    : occurrenceDefinitionPrefix VIEWPOINT DEF
      definitionDeclaration requirementBody
    ;

viewpointUsage
    : occurrenceUsagePrefix VIEWPOINT
      constraintUsageDeclaration requirementBody
    ;

// Renderings

renderingDefinition
    : occurrenceDefinitionPrefix RENDERING DEF
      definition
    ;

renderingUsage
    : occurrenceUsagePrefix RENDERING
      usage
    ;

// =============================================================================
// Clause 8.2.2.27 Metadata Textual Notation
// =============================================================================

metadataDefinition
    : ABSTRACT? definitionExtensionKeyword*
      METADATA DEF definition
    ;

prefixMetadataAnnotation
    : HASH prefixMetadataUsage
    ;

prefixMetadataMember
    : HASH prefixMetadataUsage
    ;

prefixMetadataUsage
    : ownedFeatureTyping
    ;

// Fixed: Added AT_ANNOTATION support for '@Classified' style annotations
// Also added 'about' clause support for AT_ANNOTATION
metadataUsage
    : usageExtensionKeyword* (ATAT | METADATA)
      metadataUsageDeclaration
      (ABOUT annotation (COMMA annotation)*)?
      metadataBody
    | AT_ANNOTATION
      (ABOUT annotation (COMMA annotation)*)?
      metadataBody
    ;

metadataUsageDeclaration
    : (identification (COLON | TYPED BY))?
      ownedFeatureTyping
    ;

metadataBody
    : SEMI
    | LBRACE (definitionMember | metadataBodyUsageMember | aliasMember | import_)* RBRACE
    ;

metadataBodyUsageMember
    : metadataBodyUsage
    ;

metadataBodyUsage
    : REF? (COLONGTGT | REDEFINES)? ownedRedefinition
      featureSpecializationPart? valuePart?
      metadataBody
    ;

metadataFeature
    : metadataUsage
    ;

extendedDefinition
    : basicDefinitionPrefix? definitionExtensionKeyword+
      DEF definition
    ;

extendedUsage
    : unextendedUsagePrefix usageExtensionKeyword+
      usage
    ;

// =============================================================================
// Expression Support - Using ANTLR4 precedence climbing
// =============================================================================

// Main expression rule with proper precedence using ANTLR4 left-recursion
ownedExpression
    : primaryExpression                                              // base case
    | ownedExpression conditionalBinaryOperator ownedExpression      // conditional binary ops
    | ownedExpression binaryOperator ownedExpression                 // binary ops
    | unaryOperator ownedExpression                                  // unary ops
    | ownedExpression classificationTestOperator typeReference       // classification test
    | ownedExpression castOperator typeReference                     // cast
    | ownedExpression metaCastOperator typeReference                 // meta cast (e.g., x meta SysML::Type)
    | IF ownedExpression QUESTION ownedExpression ELSE ownedExpression  // conditional
    ;

conditionalBinaryOperator
    : QUESTIONQUESTION
    | OR
    | AND
    | IMPLIES
    ;

binaryOperator
    : PIPE | AMP | XOR | DOTDOT
    | EQEQ | BANGEQ | EQEQEQ | BANGEQEQ
    | LT | GT | LTEQ | GTEQ
    | PLUS | MINUS | STAR | SLASH
    | PERCENT | CARET | STARSTAR
    ;

unaryOperator
    : PLUS
    | MINUS
    | TILDE
    | NOT
    ;

classificationTestOperator
    : ISTYPE
    | HASTYPE
    | ATAT
    ;

castOperator
    : AS
    ;

metaCastOperator
    : META
    ;

typeReference
    : qualifiedName
    ;

// Primary expressions - no left recursion issues
primaryExpression
    : baseExpression primaryExpressionSuffix*
    ;

primaryExpressionSuffix
    : LBRACKET sequenceExpressionList RBRACKET              // bracket access
    | HASH LPAREN sequenceExpressionList RPAREN             // index
    | DOT qualifiedName                                      // feature chain
    | DOT bodyExpression                                     // collect
    | DOTQUESTION bodyExpression                             // select
    | MINUSGT qualifiedName argumentList                     // function operation
    | MINUSGT qualifiedName bodyExpression                   // function operation with body
    | MINUSGT qualifiedName qualifiedName                    // function operation with shorthand argument
    | argumentList                                           // method call (e.g., x.y(args))
    ;

// Base Expressions (non-recursive)
baseExpression
    : nullExpression
    | literalExpression
    | featureReferenceExpression
    | metadataAccessExpression
    | invocationExpression
    | constructorExpression
    | bodyExpression
    | sequenceExpression
    | ALL typeReference                                      // extent expression
    ;

nullExpression
    : NULL
    | LPAREN RPAREN
    ;

featureReferenceExpression
    : qualifiedName
    ;

// Fixed: Added AT_ANNOTATION for filter expressions like 'filter @Safety'
// Also supports qualified metadata like @SysML::PartUsage
metadataAccessExpression
    : qualifiedName DOT METADATA
    | AT_ANNOTATION (COLONCOLON qualifiedName)?
    ;

invocationExpression
    : qualifiedName argumentList
    ;

constructorExpression
    : NEW qualifiedName argumentList
    ;

bodyExpression
    : LBRACE (functionBodyPart | ownedExpression) RBRACE
    ;

sequenceExpression
    : LPAREN sequenceExpressionList RPAREN
    ;

sequenceExpressionList
    : ownedExpression (COMMA ownedExpression)*
    |
    ;

argumentList
    : LPAREN (positionalArgumentList | namedArgumentList)? RPAREN
    ;

positionalArgumentList
    : ownedExpression (COMMA ownedExpression)*
    ;

namedArgumentList
    : namedArgument (COMMA namedArgument)*
    ;

namedArgument
    : qualifiedName EQ ownedExpression
    ;

functionBodyPart
    : (typeBodyElement | returnFeatureMember)*
      resultExpressionMemberOpt?
    ;

resultExpressionMemberOpt
    : memberPrefix? ownedExpression
    ;

typeBodyElement
    : nonFeatureMember
    | featureMember
    | aliasMember
    | import_
    ;

nonFeatureMember
    : memberPrefix memberElement
    ;

featureMember
    : typeFeatureMember
    | ownedFeatureMember
    ;

typeFeatureMember
    : memberPrefix MEMBER featureElement
    ;

ownedFeatureMember
    : memberPrefix featureElement
    ;

memberElement
    : annotatingElement
    | nonFeatureElement
    ;

nonFeatureElement
    : dependency
    | namespace_
    | type_
    | classifier
    | dataType
    | class
    | structure
    | metaclass
    | association
    | associationStructure
    | interaction
    | behavior
    | function_
    | predicate
    | multiplicity
    | package
    | libraryPackage
    ;

featureElement
    : featureSubsetting
    | feature
    | step
    | expression
    | booleanExpression
    | invariant
    | connector
    | bindingConnector
    | succession
    | kermlFlow
    | kermlSuccessionFlow
    ;

returnFeatureMember
    : memberPrefix RETURN featureElement
    ;

// Literal Expressions

literalExpression
    : literalBoolean
    | literalString
    | literalInteger
    | literalReal
    | literalInfinity
    ;

literalBoolean
    : TRUE
    | FALSE
    ;

literalString
    : STRING_VALUE
    ;

literalInteger
    : DECIMAL_VALUE
    ;

literalReal
    : DECIMAL_VALUE? DOT (DECIMAL_VALUE | EXPONENTIAL_VALUE)
    | EXPONENTIAL_VALUE
    ;

literalInfinity
    : STAR
    ;

// =============================================================================
// KerML Elements (needed for completeness)
// =============================================================================

namespace_
    : prefixMetadataMember*
      namespaceDeclaration namespaceBody
    ;

namespaceDeclaration
    : NAMESPACE identification
    ;

namespaceBody
    : SEMI
    | LBRACE namespaceBodyElement* RBRACE
    ;

namespaceBodyElement
    : namespaceMember
    | aliasMember
    | import_
    ;

namespaceMember
    : nonFeatureMember
    | namespaceFeatureMember
    ;

namespaceFeatureMember
    : memberPrefix featureElement
    ;

type_
    : typePrefix TYPE
      typeDeclaration typeBody
    ;

typePrefix
    : ABSTRACT?
      prefixMetadataMember*
    ;

typeDeclaration
    : ALL? identification
      ownedMultiplicity?
      (specializationPart | conjugationPart)+
      typeRelationshipPart*
    ;

specializationPart
    : specializes ownedSpecialization (COMMA ownedSpecialization)*
    ;

ownedSpecialization
    : generalType
    ;

generalType
    : qualifiedName
    | ownedFeatureChain
    ;

conjugationPart
    : conjugates ownedConjugation
    ;

ownedConjugation
    : qualifiedName
    | featureChain
    ;

conjugates
    : TILDE
    | CONJUGATES
    ;

typeRelationshipPart
    : disjoiningPart
    | unioningPart
    | intersectingPart
    | differencingPart
    ;

disjoiningPart
    : DISJOINT FROM ownedDisjoining (COMMA ownedDisjoining)*
    ;

ownedDisjoining
    : qualifiedName
    | featureChain
    ;

unioningPart
    : UNIONS unioning (COMMA unioning)*
    ;

unioning
    : qualifiedName
    | ownedFeatureChain
    ;

intersectingPart
    : INTERSECTS intersecting (COMMA intersecting)*
    ;

intersecting
    : qualifiedName
    | ownedFeatureChain
    ;

differencingPart
    : DIFFERENCES differencing (COMMA differencing)*
    ;

differencing
    : qualifiedName
    | ownedFeatureChain
    ;

typeBody
    : SEMI
    | LBRACE typeBodyElement* RBRACE
    ;

featureChain
    : ownedFeatureChaining (DOT ownedFeatureChaining)+
    ;

classifier
    : typePrefix CLASSIFIER
      classifierDeclaration typeBody
    ;

classifierDeclaration
    : ALL? identification
      ownedMultiplicity?
      (superclassingPart | conjugationPart)?
      typeRelationshipPart*
    ;

superclassingPart
    : specializes ownedSubclassification (COMMA ownedSubclassification)*
    ;

dataType
    : typePrefix DATATYPE
      classifierDeclaration typeBody
    ;

class
    : typePrefix CLASS
      classifierDeclaration typeBody
    ;

structure
    : typePrefix STRUCT
      classifierDeclaration typeBody
    ;

metaclass
    : typePrefix METACLASS
      classifierDeclaration typeBody
    ;

association
    : typePrefix ASSOC
      classifierDeclaration typeBody
    ;

associationStructure
    : typePrefix ASSOC STRUCT
      classifierDeclaration typeBody
    ;

interaction
    : typePrefix INTERACTION
      classifierDeclaration typeBody
    ;

behavior
    : typePrefix BEHAVIOR
      classifierDeclaration typeBody
    ;

function_
    : typePrefix FUNCTION
      classifierDeclaration functionBody
    ;

functionBody
    : SEMI
    | LBRACE (functionBodyPart | ownedExpression) RBRACE
    ;

predicate
    : typePrefix PREDICATE
      classifierDeclaration functionBody
    ;

multiplicity
    : multiplicitySubset
    | multiplicityRangeDecl
    ;

multiplicitySubset
    : MULTIPLICITY identification subsetsKw ownedSubsetting
      typeBody
    ;

multiplicityRangeDecl
    : MULTIPLICITY identification multiplicityBounds
      typeBody
    ;

featureSubsetting
    : SUBSET ownedSubsetting subsettings relationshipBody
    ;

multiplicityBounds
    : LBRACKET (multiplicityExpressionMember DOTDOT)?
      multiplicityExpressionMember RBRACKET
    ;

feature
    : (featurePrefix
       (FEATURE | prefixMetadataMember)
       featureDeclaration?
      | (endFeaturePrefix | basicFeaturePrefix)
        featureDeclaration
      )
      valuePart? typeBody
    ;

endFeaturePrefix
    : CONST? END FEATURE?
    ;

basicFeaturePrefix
    : featureDirection?
      DERIVED?
      ABSTRACT?
      (COMPOSITE | PORTION)?
      (VAR | CONST)?
    ;

featurePrefix
    : (endFeaturePrefix ownedCrossFeatureMember?
      | basicFeaturePrefix
      )
      prefixMetadataMember*
    ;

featureDeclaration
    : ALL?
      (featureIdentification
        (featureSpecializationPart | conjugationPart)?
      | featureSpecializationPart
      | conjugationPart
      )
      featureRelationshipPart*
    ;

featureIdentification
    : LT name GT name?
    | name
    ;

featureRelationshipPart
    : typeRelationshipPart
    | chainingPart
    | invertingPart
    | typeFeaturingPart
    ;

chainingPart
    : CHAINS (ownedFeatureChaining | featureChain)
    ;

invertingPart
    : INVERSE OF ownedFeatureInverting
    ;

ownedFeatureInverting
    : qualifiedName
    | ownedFeatureChain
    ;

typeFeaturingPart
    : FEATURED BY ownedTypeFeaturing (COMMA ownedTypeFeaturing)*
    ;

ownedTypeFeaturing
    : qualifiedName
    ;

step
    : featurePrefix
      STEP featureDeclaration valuePart?
      typeBody
    ;

expression
    : featurePrefix
      EXPR featureDeclaration valuePart?
      functionBody
    ;

booleanExpression
    : featurePrefix
      BOOL featureDeclaration valuePart?
      functionBody
    ;

invariant
    : featurePrefix
      INV (TRUE | FALSE)?
      featureDeclaration? valuePart?
      functionBody
    ;

connector
    : featurePrefix CONNECTOR
      (featureDeclaration? valuePart?
      | connectorDeclaration
      )
      typeBody
    ;

connectorDeclaration
    : binaryConnectorDeclaration
    | naryConnectorDeclaration
    ;

binaryConnectorDeclaration
    : (featureDeclaration? FROM | ALL FROM?)?
      connectorEndMember TO connectorEndMember
    ;

naryConnectorDeclaration
    : featureDeclaration?
      LPAREN connectorEndMember COMMA
      connectorEndMember
      (COMMA connectorEndMember)*
      RPAREN
    ;

bindingConnector
    : featurePrefix BINDING
      bindingConnectorDeclaration typeBody
    ;

bindingConnectorDeclaration
    : featureDeclaration
      (OF connectorEndMember EQ connectorEndMember)?
    | ALL?
      (OF? connectorEndMember EQ connectorEndMember)?
    ;

succession
    : featurePrefix SUCCESSION
      successionDeclaration typeBody
    ;

successionDeclaration
    : featureDeclaration
      (FIRST connectorEndMember THEN connectorEndMember)?
    | ALL?
      (FIRST? connectorEndMember THEN connectorEndMember)?
    ;

kermlFlow
    : featurePrefix FLOW
      kermlFlowDeclaration typeBody
    ;

kermlSuccessionFlow
    : featurePrefix SUCCESSION FLOW
      kermlFlowDeclaration typeBody
    ;

kermlFlowDeclaration
    : featureDeclaration valuePart?
      (OF kermlPayloadFeatureMember)?
      (FROM kermlFlowEndMember TO kermlFlowEndMember)?
    | ALL?
      kermlFlowEndMember TO kermlFlowEndMember
    ;

kermlPayloadFeatureMember
    : kermlPayloadFeature
    ;

kermlPayloadFeature
    : identification payloadFeatureSpecializationPart valuePart?
    | identification valuePart
    | ownedFeatureTyping ownedMultiplicity?
    | ownedMultiplicity ownedFeatureTyping?
    ;

kermlFlowEndMember
    : kermlFlowEnd
    ;

kermlFlowEnd
    : (ownedReferenceSubsetting DOT)?
      kermlFlowFeatureMember
    ;

kermlFlowFeatureMember
    : kermlFlowFeature
    ;

kermlFlowFeature
    : kermlFlowFeatureRedefinition
    ;

kermlFlowFeatureRedefinition
    : qualifiedName
    ;

// =============================================================================
// Names and Qualified Names
// =============================================================================

qualifiedName
    : (DOLLAR COLONCOLON)? (name COLONCOLON)* name
    ;

name
    : BASIC_NAME
    | UNRESTRICTED_NAME
    | keywordName
    ;

// Allow reserved words to be reused as declared/reference names in
// syntactically unambiguous name positions (e.g. redefines behavior).
keywordName
    : ABOUT
    | ABSTRACT
    | ACCEPT
    | ACTION
    | ACTOR
    | AFTER
    | ALIAS
    | ALL
    | ALLOCATE
    | ALLOCATION
    | ANALYSIS
    | AND
    | AS
    | ASSERT
    | ASSIGN
    | ASSUME
    | AT
    | ATTRIBUTE
    | BIND
    | BINDING
    | BY
    | CALC
    | CASE
    | COMMENT
    | CONCERN
    | CONNECT
    | CONNECTION
    | CONNECTOR
    | CONST
    | CONSTANT
    | CONSTRAINT
    | CROSSES
    | DECIDE
    | DEF
    | DEFAULT
    | DEFINED
    | DEPENDENCY
    | DERIVED
    | DIFFERENCES
    | DISJOINT
    | DO
    | DOC
    | ELSE
    | END
    | ENTRY
    | ENUM
    | EVENT
    | EXHIBIT
    | EXIT
    | EXPOSE
    | FALSE
    | FEATURE
    | FEATURED
    | FEATURING
    | FILTER
    | FIRST
    | FLOW
    | FOR
    | FORK
    | FRAME
    | FROM
    | FUNCTION
    | HASTYPE
    | IF
    | IMPLIES
    | IMPORT
    | IN
    | INCLUDE
    | INDIVIDUAL
    | INOUT
    | INTERFACE
    | INTERSECTS
    | INV
    | INVERSE
    | ISTYPE
    | ITEM
    | JOIN
    | LANGUAGE
    | LIBRARY
    | LOCALE
    | LOOP
    | MEMBER
    | MERGE
    | MESSAGE
    | META
    | METADATA
    | MULTIPLICITY
    | NAMESPACE
    | NEW
    | NONUNIQUE
    | NOT
    | NULL
    | OBJECTIVE
    | OCCURRENCE
    | OF
    | OR
    | ORDERED
    | OUT
    | PACKAGE
    | PARALLEL
    | PART
    | PERFORM
    | PORTION
    | PORT
    | PREDICATE
    | PRIVATE
    | PROTECTED
    | PUBLIC
    | REDEFINES
    | REDEFINITION
    | REF
    | REFERENCES
    | RENDER
    | RENDERING
    | REP
    | REQUIRE
    | REQUIREMENT
    | RETURN
    | SATISFY
    | SEND
    | SNAPSHOT
    | SPECIALIZATION
    | SPECIALIZES
    | STAKEHOLDER
    | STANDARD
    | STATE
    | STEP
    | STRUCT
    | SUBJECT
    | SUBSETS
    | SUBTYPE
    | SUBCLASSIFIER
    | SUCCESSION
    | TERMINATE
    | THEN
    | TIMESLICE
    | TO
    | TRANSITION
    | TRUE
    | TYPE
    | TYPED
    | TYPING
    | UNIONS
    | UNTIL
    | USE
    | VAR
    | VARIANT
    | VARIATION
    | VERIFICATION
    | VERIFY
    | VIA
    | VIEW
    | VIEWPOINT
    | WHEN
    | WHILE
    | XOR
    | ASSOC
    | BEHAVIOR
    | BOOL
    | CHAINS
    | CLASS
    | CLASSIFIER
    | COMPOSITE
    | CONJUGATE
    | CONJUGATES
    | CONJUGATION
    | DATATYPE
    | DISJOINING
    | EXPR
    | INTERACTION
    | INVERTING
    | METACLASS
    ;
