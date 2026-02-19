/**
 * SysML v2 Lexer Grammar for ANTLR4
 *
 * This lexer handles both KerML and SysML tokens.
 * Compatible with Go target for parser generation.
 *
 * Based on OMG SysML v2 specification.
 */
lexer grammar SysMLv2Lexer;

// =============================================================================
// Keywords - SysML (includes KerML keywords)
// =============================================================================

// SysML-specific keywords
ABOUT           : 'about';
ABSTRACT        : 'abstract';
ACCEPT          : 'accept';
ACTION          : 'action';
ACTOR           : 'actor';
AFTER           : 'after';
ALIAS           : 'alias';
ALL             : 'all';
ALLOCATE        : 'allocate';
ALLOCATION      : 'allocation';
ANALYSIS        : 'analysis';
AND             : 'and';
AS              : 'as';
ASSERT          : 'assert';
ASSIGN          : 'assign';
ASSUME          : 'assume';
AT              : 'at';
ATTRIBUTE       : 'attribute';
BIND            : 'bind';
BINDING         : 'binding';
BY              : 'by';
CALC            : 'calc';
CASE            : 'case';
COMMENT         : 'comment';
CONCERN         : 'concern';
CONNECT         : 'connect';
CONNECTION      : 'connection';
CONNECTOR       : 'connector';
CONST           : 'const';
CONSTANT        : 'constant';
CONSTRAINT      : 'constraint';
CROSSES         : 'crosses';
DECIDE          : 'decide';
DEF             : 'def';
DEFAULT         : 'default';
DEFINED         : 'defined';
DEPENDENCY      : 'dependency';
DERIVED         : 'derived';
DIFFERENCES     : 'differences';
DISJOINT        : 'disjoint';
DO              : 'do';
DOC             : 'doc';
ELSE            : 'else';
END             : 'end';
ENTRY           : 'entry';
ENUM            : 'enum';
EVENT           : 'event';
EXHIBIT         : 'exhibit';
EXIT            : 'exit';
EXPOSE          : 'expose';
FALSE           : 'false';
FEATURE         : 'feature';
FEATURED        : 'featured';
FEATURING       : 'featuring';
FILTER          : 'filter';
FIRST           : 'first';
FLOW            : 'flow';
FOR             : 'for';
FORK            : 'fork';
FRAME           : 'frame';
FROM            : 'from';
FUNCTION        : 'function';
HASTYPE         : 'hastype';
IF              : 'if';
IMPLIES         : 'implies';
IMPORT          : 'import';
IN              : 'in';
INCLUDE         : 'include';
INDIVIDUAL      : 'individual';
INOUT           : 'inout';
INTERFACE       : 'interface';
INTERSECTS      : 'intersects';
INV             : 'inv';
INVERSE         : 'inverse';
ISTYPE          : 'istype';
ITEM            : 'item';
JOIN            : 'join';
LANGUAGE        : 'language';
LIBRARY         : 'library';
LOCALE          : 'locale';
LOOP            : 'loop';
MEMBER          : 'member';
MERGE           : 'merge';
MESSAGE         : 'message';
META            : 'meta';
METADATA        : 'metadata';
MULTIPLICITY    : 'multiplicity';
NAMESPACE       : 'namespace';
NEW             : 'new';
NONUNIQUE       : 'nonunique';
NOT             : 'not';
NULL            : 'null';
OBJECTIVE       : 'objective';
OCCURRENCE      : 'occurrence';
OF              : 'of';
OR              : 'or';
ORDERED         : 'ordered';
OUT             : 'out';
PACKAGE         : 'package';
PARALLEL        : 'parallel';
PART            : 'part';
PERFORM         : 'perform';
PORTION         : 'portion';
PORT            : 'port';
PREDICATE       : 'predicate';
PRIVATE         : 'private';
PROTECTED       : 'protected';
PUBLIC          : 'public';
REDEFINES       : 'redefines';
REDEFINITION    : 'redefinition';
REF             : 'ref';
REFERENCES      : 'references';
RENDER          : 'render';
RENDERING       : 'rendering';
REP             : 'rep';
REQUIRE         : 'require';
REQUIREMENT     : 'requirement';
RETURN          : 'return';
SATISFY         : 'satisfy';
SEND            : 'send';
SNAPSHOT        : 'snapshot';
SPECIALIZATION  : 'specialization';
SPECIALIZES     : 'specializes';
STAKEHOLDER     : 'stakeholder';
STANDARD        : 'standard';
STATE           : 'state';
STEP            : 'step';
STRUCT          : 'struct';
SUBJECT         : 'subject';
SUBSET          : 'subset';
SUBSETS         : 'subsets';
SUBTYPE         : 'subtype';
SUBCLASSIFIER   : 'subclassifier';
SUCCESSION      : 'succession';
TERMINATE       : 'terminate';
THEN            : 'then';
TIMESLICE       : 'timeslice';
TO              : 'to';
TRANSITION      : 'transition';
TRUE            : 'true';
TYPE            : 'type';
TYPED           : 'typed';
TYPING          : 'typing';
UNIONS          : 'unions';
UNTIL           : 'until';
USE             : 'use';
VAR             : 'var';
VARIANT         : 'variant';
VARIATION       : 'variation';
VERIFICATION    : 'verification';
VERIFY          : 'verify';
VIA             : 'via';
VIEW            : 'view';
VIEWPOINT       : 'viewpoint';
WHEN            : 'when';
WHILE           : 'while';
XOR             : 'xor';

// KerML-specific keywords (not already in SysML)
ASSOC           : 'assoc';
BEHAVIOR        : 'behavior';
BOOL            : 'bool';
CHAINS          : 'chains';
CLASS           : 'class';
CLASSIFIER      : 'classifier';
COMPOSITE       : 'composite';
CONJUGATE       : 'conjugate';
CONJUGATES      : 'conjugates';
CONJUGATION     : 'conjugation';
DATATYPE        : 'datatype';
DISJOINING      : 'disjoining';
EXPR            : 'expr';
INTERACTION     : 'interaction';
INVERTING       : 'inverting';
METACLASS       : 'metaclass';

// =============================================================================
// Operator Symbols
// =============================================================================

// Multi-character operators (must come before single-character)
COLONCOLONGT    : '::>';
COLONGTGT       : ':>>';
COLONCOLON      : '::';
COLONGT         : ':>';
COLONEQ         : ':=';
EQEQEQ          : '===';
EQEQ            : '==';
BANGEQEQ        : '!==';
BANGEQ          : '!=';
LTEQ            : '<=';
GTEQ            : '>=';
EQGT            : '=>';
MINUSGT         : '->';
DOTDOT          : '..';
DOTQUESTION     : '.?';
STARSTAR        : '**';
QUESTIONQUESTION: '??';
ATAT            : '@@';

// @ prefix for metadata annotations (allows optional space after @)
AT_ANNOTATION   : '@' [ \t]* [a-zA-Z_] [a-zA-Z0-9_]*;

// Single-character operators
LBRACE          : '{';
RBRACE          : '}';
LBRACKET        : '[';
RBRACKET        : ']';
LPAREN          : '(';
RPAREN          : ')';
LT              : '<';
GT              : '>';
SEMI            : ';';
COLON           : ':';
DOT             : '.';
COMMA           : ',';
EQ              : '=';
PLUS            : '+';
MINUS           : '-';
STAR            : '*';
SLASH           : '/';
PERCENT         : '%';
CARET           : '^';
AMP             : '&';
PIPE            : '|';
TILDE           : '~';
QUESTION        : '?';
HASH            : '#';
DOLLAR          : '$';
BANG            : '!';

// =============================================================================
// Literals and Names
// =============================================================================

// Decimal and exponential values
DECIMAL_VALUE
    : DECIMAL_DIGIT+
    ;

EXPONENTIAL_VALUE
    : DECIMAL_DIGIT+ ('e' | 'E') ('+' | '-')? DECIMAL_DIGIT+
    ;

// String literal
STRING_VALUE
    : '"' (STRING_CHARACTER | ESCAPE_SEQUENCE)* '"'
    ;

// Basic name (identifier)
BASIC_NAME
    : BASIC_INITIAL_CHARACTER BASIC_NAME_CHARACTER*
    ;

// Unrestricted name (quoted identifier)
UNRESTRICTED_NAME
    : '\'' (NAME_CHARACTER | ESCAPE_SEQUENCE)* '\''
    ;

// =============================================================================
// Comments and Whitespace
// =============================================================================

// Regular comment (used as body for Comment elements)
REGULAR_COMMENT
    : '/*' .*? '*/'
    ;

// Single-line note (ignored)
SINGLE_LINE_NOTE
    : '//' ~[\r\n]* -> channel(HIDDEN)
    ;

// Multiline note (ignored)
MULTILINE_NOTE
    : '//*' .*? '*/' -> channel(HIDDEN)
    ;

// Whitespace
WS
    : [ \t\r\n\f]+ -> channel(HIDDEN)
    ;

// =============================================================================
// Lexer Fragments
// =============================================================================

fragment DECIMAL_DIGIT
    : [0-9]
    ;

fragment BASIC_INITIAL_CHARACTER
    : [a-zA-Z_]
    ;

fragment BASIC_NAME_CHARACTER
    : BASIC_INITIAL_CHARACTER
    | DECIMAL_DIGIT
    ;

fragment STRING_CHARACTER
    : ~["\\\r\n]
    ;

fragment NAME_CHARACTER
    : ~['\\\r\n]
    ;

fragment ESCAPE_SEQUENCE
    : '\\' [fnrtv'"\\]
    ;
