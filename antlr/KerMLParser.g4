/**
 * KerML (Kernel Modeling Language) Parser Grammar for ANTLR4
 *
 * This parser implements the KerML textual concrete syntax from OMG SysML v2.
 * Compatible with Go target for parser generation.
 *
 * Based on OMG SysML v2 specification Part 1 - Kernel Modeling Language (KerML)
 */
parser grammar KerMLParser;

options {
    tokenVocab = SysMLv2Lexer;
}

// =============================================================================
// Clause 8.2.3 Root Concrete Syntax
// =============================================================================

// Clause 8.2.3.1 Elements and Relationships Concrete Syntax

entryRuleRootNamespace
    : rootNamespace EOF
    ;

rootNamespace
    : namespaceBodyElement*
    ;

identification
    : (LT name GT)? name?
    ;

relationshipBody
    : SEMI
    | LBRACE relationshipOwnedElement* RBRACE
    ;

relationshipOwnedElement
    : ownedRelatedElement
    | ownedAnnotation
    ;

ownedRelatedElement
    : nonFeatureElement
    | featureElement
    ;

// Clause 8.2.3.2 Dependencies Concrete Syntax

dependency
    : prefixMetadataAnnotation*
      DEPENDENCY (identification? FROM)?
      qualifiedName (COMMA qualifiedName)* TO
      qualifiedName (COMMA qualifiedName)*
      relationshipBody
    ;

// Clause 8.2.3.3 Annotations Concrete Syntax

annotation
    : qualifiedName
    ;

ownedAnnotation
    : annotatingElement
    ;

annotatingElement
    : comment_
    | documentation
    | textualRepresentation
    | metadataFeature
    ;

// Comments and Documentation
comment_
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

// Textual Representation
textualRepresentation
    : (REP identification)?
      LANGUAGE STRING_VALUE REGULAR_COMMENT
    ;

// Clause 8.2.3.4 Namespaces Concrete Syntax

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

memberPrefix
    : visibilityIndicator?
    ;

visibilityIndicator
    : PUBLIC
    | PRIVATE
    | PROTECTED
    ;

namespaceMember
    : nonFeatureMember
    | namespaceFeatureMember
    ;

nonFeatureMember
    : memberPrefix memberElement
    ;

namespaceFeatureMember
    : memberPrefix featureElement
    ;

aliasMember
    : memberPrefix
      ALIAS (LT name GT)?
      name?
      FOR qualifiedName
      relationshipBody
    | memberPrefix
      ALIAS qualifiedName
      AS name
      relationshipBody
    ;

qualifiedName
    : (DOLLAR COLONCOLON)? (name COLONCOLON)* name
    ;

name
    : BASIC_NAME
    | UNRESTRICTED_NAME
    ;

// Clause 8.2.3.4.2 Imports

import_
    : visibilityIndicator?
      IMPORT ALL?
      importDeclaration relationshipBody
    ;

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

filterPackage
    : importDeclaration filterPackageMember+
    ;

filterPackageMember
    : LBRACKET ownedExpression RBRACKET
    ;

// Clause 8.2.3.4.3 Namespace Elements

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
    | class_
    | structure
    | metaclass
    | association
    | associationStructure
    | interaction
    | behavior
    | function_
    | predicate
    | multiplicity
    | package_
    | libraryPackage
    | specialization
    | conjugation
    | subclassification
    | disjoining
    | featureInverting
    | featureTyping
    | subsetting
    | redefinition
    | typeFeaturing
    ;

featureElement
    : feature
    | step
    | expression_
    | booleanExpression
    | invariant
    | connector
    | bindingConnector
    | succession
    | flow
    | successionFlow
    ;

// =============================================================================
// Clause 8.2.4 Core Concrete Syntax
// =============================================================================

// Clause 8.2.4.1 Types Concrete Syntax

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

conjugationPart
    : conjugates ownedConjugation
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

unioningPart
    : UNIONS unioning (COMMA unioning)*
    ;

intersectingPart
    : INTERSECTS intersecting (COMMA intersecting)*
    ;

differencingPart
    : DIFFERENCES differencing (COMMA differencing)*
    ;

typeBody
    : SEMI
    | LBRACE typeBodyElement* RBRACE
    ;

typeBodyElement
    : nonFeatureMember
    | featureMember
    | aliasMember
    | import_
    ;

// Specialization keyword alternatives
specializes
    : COLONGT
    | SPECIALIZES
    ;

// Conjugation keyword alternatives
conjugates
    : TILDE
    | CONJUGATES
    ;

// Clause 8.2.4.1.2 Specialization

specialization
    : (SPECIALIZATION identification)?
      SUBTYPE specificType
      specializes generalType
      relationshipBody
    ;

ownedSpecialization
    : generalType
    ;

specificType
    : qualifiedName
    | ownedFeatureChain
    ;

generalType
    : qualifiedName
    | ownedFeatureChain
    ;

// Clause 8.2.4.1.3 Conjugation

conjugation
    : (CONJUGATION identification)?
      CONJUGATE
      (qualifiedName | featureChain)
      conjugates
      (qualifiedName | featureChain)
      relationshipBody
    ;

ownedConjugation
    : qualifiedName
    | featureChain
    ;

// Clause 8.2.4.1.4 Disjoining

disjoining
    : (DISJOINING identification)?
      DISJOINT
      (qualifiedName | featureChain)
      FROM
      (qualifiedName | featureChain)
      relationshipBody
    ;

ownedDisjoining
    : qualifiedName
    | featureChain
    ;

// Clause 8.2.4.1.5 Unioning, Intersecting and Differencing

unioning
    : qualifiedName
    | ownedFeatureChain
    ;

intersecting
    : qualifiedName
    | ownedFeatureChain
    ;

differencing
    : qualifiedName
    | ownedFeatureChain
    ;

// Clause 8.2.4.1.6 Feature Membership

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

// Clause 8.2.4.2 Classifiers Concrete Syntax

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

// Clause 8.2.4.2.2 Subclassification

subclassification
    : (SPECIALIZATION identification)?
      SUBCLASSIFIER qualifiedName
      specializes qualifiedName
      relationshipBody
    ;

ownedSubclassification
    : qualifiedName
    ;

// Clause 8.2.4.3 Features Concrete Syntax

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
    : CONST?
      END
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

ownedCrossFeatureMember
    : ownedCrossFeature
    ;

ownedCrossFeature
    : basicFeaturePrefix featureDeclaration
    ;

featureDirection
    : IN
    | OUT
    | INOUT
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
    : CHAINS
      (ownedFeatureChaining | featureChain)
    ;

invertingPart
    : INVERSE OF ownedFeatureInverting
    ;

typeFeaturingPart
    : FEATURED BY ownedTypeFeaturing (COMMA ownedTypeFeaturing)*
    ;

featureSpecializationPart
    : featureSpecialization+ multiplicityPart? featureSpecialization*
    | multiplicityPart featureSpecialization*
    ;

multiplicityPart
    : ownedMultiplicity
    | ownedMultiplicity?
      (ORDERED NONUNIQUE?
      | NONUNIQUE ORDERED?)
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

// DEFINED_BY keyword alternatives
definedBy
    : COLON
    | DEFINED BY
    ;

subsettings
    : subsets (COMMA ownedSubsetting)*
    ;

subsets
    : subsetsKw ownedSubsetting
    ;

// SUBSETS keyword alternatives
subsetsKw
    : COLONGT
    | SUBSETS
    ;

references
    : referencesKw ownedReferenceSubsetting
    ;

// REFERENCES keyword alternatives
referencesKw
    : COLONCOLONGT
    | REFERENCES
    ;

crosses
    : crossesKw ownedCrossSubsetting
    ;

// CROSSES keyword alternatives
crossesKw
    : EQGT
    | CROSSES
    ;

redefinitions
    : redefines (COMMA ownedRedefinition)*
    ;

redefines
    : redefinesKw ownedRedefinition
    ;

// REDEFINES keyword alternatives
redefinesKw
    : COLONGTGT
    | REDEFINES
    ;

// Clause 8.2.4.3.2 Feature Typing

featureTyping
    : (SPECIALIZATION identification)?
      TYPING qualifiedName
      definedBy generalType
      relationshipBody
    ;

ownedFeatureTyping
    : generalType
    ;

// Clause 8.2.4.3.3 Subsetting

subsetting
    : (SPECIALIZATION identification)?
      SUBSET specificType
      subsetsKw generalType
      relationshipBody
    ;

ownedSubsetting
    : generalType
    ;

ownedReferenceSubsetting
    : generalType
    ;

ownedCrossSubsetting
    : generalType
    ;

// Clause 8.2.4.3.4 Redefinition

redefinition
    : (SPECIALIZATION identification)?
      REDEFINITION specificType
      redefinesKw generalType
      relationshipBody
    ;

ownedRedefinition
    : generalType
    ;

// Clause 8.2.4.3.5 Feature Chaining

ownedFeatureChain
    : featureChain
    ;

featureChain
    : ownedFeatureChaining (DOT ownedFeatureChaining)+
    ;

ownedFeatureChaining
    : qualifiedName
    ;

// Clause 8.2.4.3.6 Feature Inverting

featureInverting
    : (INVERTING identification?)?
      INVERSE
      (qualifiedName | ownedFeatureChain)
      OF
      (qualifiedName | ownedFeatureChain)
      relationshipBody
    ;

ownedFeatureInverting
    : qualifiedName
    | ownedFeatureChain
    ;

// Clause 8.2.4.3.7 Type Featuring

typeFeaturing
    : FEATURING (identification OF)?
      qualifiedName
      BY qualifiedName
      relationshipBody
    ;

ownedTypeFeaturing
    : qualifiedName
    ;

// =============================================================================
// Clause 8.2.5 Kernel Concrete Syntax
// =============================================================================

// Clause 8.2.5.1 Data Types

dataType
    : typePrefix DATATYPE
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.2 Classes

class_
    : typePrefix CLASS
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.3 Structures

structure
    : typePrefix STRUCT
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.4 Associations

association
    : typePrefix ASSOC
      classifierDeclaration typeBody
    ;

associationStructure
    : typePrefix ASSOC STRUCT
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.5 Connectors

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
      connectorEndMember TO
      connectorEndMember
    ;

naryConnectorDeclaration
    : featureDeclaration?
      LPAREN connectorEndMember COMMA
      connectorEndMember
      (COMMA connectorEndMember)*
      RPAREN
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

// Clause 8.2.5.5.2 Binding Connectors

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

// Clause 8.2.5.5.3 Successions

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

// Clause 8.2.5.6 Behaviors

behavior
    : typePrefix BEHAVIOR
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.6.2 Steps

step
    : featurePrefix
      STEP featureDeclaration valuePart?
      typeBody
    ;

// Clause 8.2.5.7 Functions

function_
    : typePrefix FUNCTION
      classifierDeclaration functionBody
    ;

functionBody
    : SEMI
    | LBRACE functionBodyPart RBRACE
    ;

functionBodyPart
    : (typeBodyElement | returnFeatureMember)*
      resultExpressionMember?
    ;

returnFeatureMember
    : memberPrefix RETURN featureElement
    ;

resultExpressionMember
    : memberPrefix ownedExpression
    ;

// Clause 8.2.5.7.2 Expressions

expression_
    : featurePrefix
      EXPR featureDeclaration valuePart?
      functionBody
    ;

// Clause 8.2.5.7.3 Predicates

predicate
    : typePrefix PREDICATE
      classifierDeclaration functionBody
    ;

// Clause 8.2.5.7.4 Boolean Expressions and Invariants

booleanExpression
    : featurePrefix
      BOOL featureDeclaration valuePart?
      functionBody
    ;

invariant
    : featurePrefix
      INV (TRUE | FALSE)?
      featureDeclaration valuePart?
      functionBody
    ;

// =============================================================================
// Clause 8.2.5.8 Expressions Concrete Syntax
// =============================================================================

ownedExpressionReferenceMember
    : ownedExpressionReference
    ;

ownedExpressionReference
    : ownedExpressionMember
    ;

ownedExpressionMember
    : ownedExpression
    ;

ownedExpression
    : conditionalExpression
    | conditionalBinaryOperatorExpression
    | binaryOperatorExpression
    | unaryOperatorExpression
    | classificationExpression
    | metaclassificationExpression
    | extentExpression
    | primaryExpression
    ;

conditionalExpression
    : IF argumentMember QUESTION
      argumentExpressionMember ELSE
      argumentExpressionMember
    ;

conditionalBinaryOperatorExpression
    : argumentMember conditionalBinaryOperator argumentExpressionMember
    ;

conditionalBinaryOperator
    : QUESTIONQUESTION
    | OR
    | AND
    | OROR
    | ANDAND
    | IMPLIES
    ;

binaryOperatorExpression
    : argumentMember binaryOperator argumentMember
    ;

binaryOperator
    : PIPE | AMP | XOR | DOTDOT
    | EQEQ | BANGEQ | EQEQEQ | BANGEQEQ
    | LT | GT | LTEQ | GTEQ
    | PLUS | MINUS | STAR | SLASH
    | PERCENT | CARET | STARSTAR
    ;

unaryOperatorExpression
    : unaryOperator argumentMember
    ;

unaryOperator
    : PLUS
    | MINUS
    | TILDE
    | NOT
    ;

classificationExpression
    : argumentMember?
      (classificationTestOperator typeReferenceMember
      | castOperator typeResultMember
      )
    ;

classificationTestOperator
    : ISTYPE
    | HASTYPE
    | ATAT
    ;

castOperator
    : AS
    ;

metaclassificationExpression
    : metadataArgumentMember
      (classificationTestOperator typeReferenceMember
      | metaCastOperator typeResultMember
      )
    ;

metaCastOperator
    : META
    ;

extentExpression
    : ALL typeReferenceMember
    ;

argumentMember
    : argument
    ;

argument
    : argumentValue
    ;

argumentValue
    : ownedExpression
    ;

argumentExpressionMember
    : argumentExpression
    ;

argumentExpression
    : argumentExpressionValue
    ;

argumentExpressionValue
    : ownedExpressionReference
    ;

metadataArgumentMember
    : metadataArgument
    ;

metadataArgument
    : metadataValue
    ;

metadataValue
    : metadataReference
    ;

metadataReference
    : elementReferenceMember
    ;

typeReferenceMember
    : typeReference
    ;

typeResultMember
    : typeReference
    ;

typeReference
    : referenceTyping
    ;

referenceTyping
    : qualifiedName
    ;

// Clause 8.2.5.8.2 Primary Expressions

primaryExpression
    : featureChainExpression
    | nonFeatureChainPrimaryExpression
    ;

nonFeatureChainPrimaryExpression
    : bracketExpression
    | indexExpression
    | sequenceExpression
    | selectExpression
    | collectExpression
    | functionOperationExpression
    | quantityAnnotationExpression
    | baseExpression
    ;

quantityAnnotationExpression
    : primaryArgumentMember ATSIGN LBRACKET qualifiedName RBRACKET
    ;

bracketExpression
    : primaryArgumentMember LBRACKET sequenceExpressionListMember RBRACKET
    ;

indexExpression
    : primaryArgumentMember HASH
      LPAREN sequenceExpressionListMember RPAREN
    ;

sequenceExpression
    : LPAREN sequenceExpressionList RPAREN
    ;

sequenceExpressionList
    : ownedExpression COMMA?
    | sequenceOperatorExpression
    ;

sequenceOperatorExpression
    : ownedExpressionMember COMMA sequenceExpressionListMember
    ;

sequenceExpressionListMember
    : sequenceExpressionList
    ;

featureChainExpression
    : nonFeatureChainPrimaryArgumentMember DOT featureChainMember
    ;

collectExpression
    : primaryArgumentMember DOT bodyArgumentMember
    ;

selectExpression
    : primaryArgumentMember DOTQUESTION bodyArgumentMember
    ;

functionOperationExpression
    : primaryArgumentMember MINUSGT invocationTypeMember
      (bodyArgumentMember
      | functionReferenceArgumentMember
      | argumentList
      )
    ;

primaryArgumentMember
    : primaryArgument
    ;

primaryArgument
    : primaryArgumentValue
    ;

primaryArgumentValue
    : primaryExpression
    ;

nonFeatureChainPrimaryArgumentMember
    : nonFeatureChainPrimaryArgument
    ;

nonFeatureChainPrimaryArgument
    : nonFeatureChainPrimaryArgumentValue
    ;

nonFeatureChainPrimaryArgumentValue
    : nonFeatureChainPrimaryExpression
    ;

bodyArgumentMember
    : bodyArgument
    ;

bodyArgument
    : bodyArgumentValue
    ;

bodyArgumentValue
    : bodyExpression
    ;

functionReferenceArgumentMember
    : functionReferenceArgument
    ;

functionReferenceArgument
    : functionReferenceArgumentValue
    ;

functionReferenceArgumentValue
    : functionReferenceExpression
    ;

functionReferenceExpression
    : functionReferenceMember
    ;

functionReferenceMember
    : functionReference
    ;

functionReference
    : referenceTyping
    ;

invocationTypeMember
    : instantiatedTypeReference
    | ownedFeatureChainMember
    ;

featureChainMember
    : featureReferenceMember
    | ownedFeatureChainMember
    ;

ownedFeatureChainMember
    : featureChain
    ;

// Clause 8.2.5.8.3 Base Expressions

baseExpression
    : nullExpression
    | literalExpression
    | featureReferenceExpression
    | metadataAccessExpression
    | invocationExpression
    | constructorExpression
    | bodyExpression
    ;

nullExpression
    : NULL
    | LPAREN RPAREN
    ;

featureReferenceExpression
    : featureReferenceMember
    ;

featureReferenceMember
    : featureReference
    ;

featureReference
    : qualifiedName
    ;

metadataAccessExpression
    : elementReferenceMember DOT METADATA
    ;

elementReferenceMember
    : qualifiedName
    ;

invocationExpression
    : instantiatedTypeMember
      argumentList
    ;

constructorExpression
    : NEW instantiatedTypeMember
      constructorResultMember
    ;

constructorResultMember
    : constructorResult
    ;

constructorResult
    : argumentList
    ;

instantiatedTypeMember
    : instantiatedTypeReference
    | ownedFeatureChainMember
    ;

instantiatedTypeReference
    : qualifiedName
    ;

argumentList
    : LPAREN (positionalArgumentList | namedArgumentList)? RPAREN
    ;

positionalArgumentList
    : argumentMember (COMMA argumentMember)*
    ;

namedArgumentList
    : namedArgumentMember (COMMA namedArgumentMember)*
    ;

namedArgumentMember
    : namedArgument
    ;

namedArgument
    : parameterRedefinition EQ argumentValue
    ;

parameterRedefinition
    : qualifiedName
    ;

bodyExpression
    : expressionBodyMember
    ;

expressionBodyMember
    : expressionBody
    ;

expressionBody
    : LBRACE functionBodyPart RBRACE
    ;

// Clause 8.2.5.8.4 Literal Expressions

literalExpression
    : literalBoolean
    | literalString
    | literalInteger
    | literalReal
    | literalInfinity
    ;

literalBoolean
    : booleanValue
    ;

booleanValue
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
    : realValue
    ;

realValue
    : DECIMAL_VALUE? DOT (DECIMAL_VALUE | EXPONENTIAL_VALUE)
    | EXPONENTIAL_VALUE
    ;

literalInfinity
    : STAR
    ;

// Clause 8.2.5.9 Interactions

interaction
    : typePrefix INTERACTION
      classifierDeclaration typeBody
    ;

// Clause 8.2.5.9.2 Flows

flow
    : featurePrefix FLOW
      flowDeclaration typeBody
    ;

successionFlow
    : featurePrefix SUCCESSION FLOW
      flowDeclaration typeBody
    ;

flowDeclaration
    : featureDeclaration valuePart?
      (OF payloadFeatureMember)?
      (FROM flowEndMember TO flowEndMember)?
    | ALL?
      flowEndMember TO flowEndMember
    ;

payloadFeatureMember
    : payloadFeature
    ;

payloadFeature
    : identification payloadFeatureSpecializationPart valuePart?
    | identification valuePart
    | ownedFeatureTyping ownedMultiplicity?
    | ownedMultiplicity ownedFeatureTyping?
    ;

payloadFeatureSpecializationPart
    : featureSpecialization+ multiplicityPart? featureSpecialization*
    | multiplicityPart featureSpecialization+
    ;

flowEndMember
    : flowEnd
    ;

flowEnd
    : (ownedReferenceSubsetting DOT)?
      flowFeatureMember
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

// Clause 8.2.5.10 Feature Values

valuePart
    : featureValue
    ;

featureValue
    : (EQ | COLONEQ | DEFAULT (EQ | COLONEQ)?)
      ownedExpression
    ;

// Clause 8.2.5.11 Multiplicities

multiplicity
    : multiplicitySubset
    | multiplicityRange
    ;

multiplicitySubset
    : MULTIPLICITY identification subsets
      typeBody
    ;

multiplicityRange
    : MULTIPLICITY identification multiplicityBounds
      typeBody
    ;

ownedMultiplicity
    : ownedMultiplicityRange
    ;

ownedMultiplicityRange
    : multiplicityBounds
    ;

multiplicityBounds
    : LBRACKET (multiplicityExpressionMember DOTDOT)?
      multiplicityExpressionMember RBRACKET
    ;

multiplicityExpressionMember
    : literalExpression
    | featureReferenceExpression
    ;

// Clause 8.2.5.12 Metadata

metaclass
    : typePrefix METACLASS
      classifierDeclaration typeBody
    ;

prefixMetadataAnnotation
    : HASH prefixMetadataFeature
    ;

prefixMetadataMember
    : HASH prefixMetadataFeature
    ;

prefixMetadataFeature
    : ownedFeatureTyping
    ;

metadataFeature
    : prefixMetadataMember*
      (ATAT | METADATA)
      metadataFeatureDeclaration
      (ABOUT annotation (COMMA annotation)*)?
      metadataBody
    ;

metadataFeatureDeclaration
    : (identification (COLON | TYPED BY))?
      ownedFeatureTyping
    ;

metadataBody
    : SEMI
    | LBRACE metadataBodyElement* RBRACE
    ;

metadataBodyElement
    : nonFeatureMember
    | metadataBodyFeatureMember
    | aliasMember
    | import_
    ;

metadataBodyFeatureMember
    : metadataBodyFeature
    ;

metadataBodyFeature
    : FEATURE? (COLONGTGT | REDEFINES)? ownedRedefinition
      featureSpecializationPart? valuePart?
      metadataBody
    ;

// Clause 8.2.5.13 Packages

package_
    : prefixMetadataMember*
      packageDeclaration packageBody
    ;

libraryPackage
    : STANDARD LIBRARY
      prefixMetadataMember*
      packageDeclaration packageBody
    ;

packageDeclaration
    : PACKAGE identification
    ;

packageBody
    : SEMI
    | LBRACE (namespaceBodyElement | elementFilterMember)* RBRACE
    ;

elementFilterMember
    : memberPrefix
      FILTER ownedExpression SEMI
    ;
