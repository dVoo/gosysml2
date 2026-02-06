# Attribute Parsing Fix Summary

## Problem

The SysML v2 parser was not correctly extracting attribute names when using the `:>>` (redefinition) syntax.

**Example:**
```sysml
attribute :>> maxLatency_ms = 200;
```

**Before fix:** Attribute had no name, only the value `200`
**After fix:** Attribute has name `maxLatency_ms` and value `200`

## Solution

Modified `gosysml2/sysml/parse.go` to extract the redefined feature name from the `featureSpecializationPart` when the `identification` doesn't contain a name.

### Changes Made

1. Added `extractRedefinitionName()` helper function that:
   - Traverses `featureSpecializationPart`
   - Looks for `redefinitions` context
   - Extracts the `qualifiedName` from the first `ownedRedefinition`

2. Modified `EnterAttributeUsage()` to:
   - First try to get name from `identification`
   - If empty, fall back to extracting from redefinitions

## Syntax Corrections for Your Example

Your original code had several syntax errors. Here are the corrections:

### 1. Requirement ID
❌ **Wrong:** `<REQ-001>`
✅ **Correct:** `<'REQ-001'>`
*Requirement IDs must be quoted strings*

### 2. Verification Definition
❌ **Wrong:** `verification case def VerifyUILatency`
✅ **Correct:** `verification def VerifyUILatency`
*Use `verification def` for verification definitions*

### 3. Verification Usage
❌ **Wrong:** `verification case verifyUILatency`
✅ **Correct:** `verification verifyUILatency`
*Use `verification` for verification usages*

### 4. Requirement Satisfaction
❌ **Wrong:** `verify REQ-001 by verifyUILatency;`
✅ **Correct:** `satisfy uiLatency by verifyUILatency;`
*Use `satisfy` to link requirements to verifications, and use the requirement name (not ID)*

## Corrected Working Example

```sysml
package ExampleReqAndVerification {

  // Reusable requirement type with custom attributes
  requirement def LatencyRequirement {
    subject system;
    attribute maxLatency_ms : Integer;
    attribute priority : Integer;
    attribute owner : String;
    attribute rationale : String;

    require constraint {
      latencyActual_ms <= maxLatency_ms
    }

    attribute latencyActual_ms : Integer;
  }

  // Concrete requirement with multiple custom attribute values
  requirement <'REQ-001'> uiLatency : LatencyRequirement {
    attribute :>> maxLatency_ms = 200;
    attribute :>> priority = 1;
    attribute :>> owner = "UX Team";
    attribute :>> rationale = "Keep UI interactions responsive under nominal load.";
    attribute :>> latencyActual_ms = 145;
  }

  // Verification definition and usage
  verification def VerifyUILatency {
    subject system;
  }

  verification verifyUILatency : VerifyUILatency;

  // Link requirement to verification
  satisfy uiLatency by verifyUILatency;
}
```

## Test Results

After the fix, attributes are correctly parsed:

```
Requirement: LatencyRequirement (def=true)
  Attributes (5):
    - maxLatency_ms
    - priority
    - owner
    - rationale
    - latencyActual_ms
  Constraints (1):
    - <unnamed>
      Expression: constraint{latencyActual_ms<=maxLatency_ms}

Requirement: 'REQ-001' (def=false)
  Attributes (5):
    - maxLatency_ms = 200
    - priority = 1
    - owner = "UX Team"
    - rationale = "Keep UI interactions responsive under nominal load."
    - latencyActual_ms = 145
```

All custom attributes are now correctly captured with both names and values in the model.
