# Architecture

## 1. Overview in One Minute

`definitioncheck` is the authoritative validator for webtyp module and entity declarations.

In the webtyp architecture, the [`model`](https://github.com/webtyp/model) library is purely declarative: it defines structs like `model.Definition`, `model.Rules`, and `model.ModuleDefinition` without business or validation logic. `definitioncheck` provides the logic ("the judge") that inspects these declarations.

Two consumers rely on this judge:
1. **[`definitionc`](https://github.com/webtyp/definitionc)**: Calls `definitioncheck` during code generation to refuse invalid module declarations before generating Go code.
2. **[`crud`](https://github.com/webtyp/crud)**: Calls `definitioncheck` at runtime when mounting modules to guarantee that registered entities comply with authorization rules.

By concentrating validation in one isomorphic library, rule invariants (such as "every allowed action must declare a rule") cannot diverge between compile-time tooling and runtime engines.

## 2. Repository Pieces & On-Disk Layout

| File | Single Responsibility | Caller / Consumer |
|---|---|---|
| [`rules.go`](../rules.go) | Validates single entity declarations (`CheckRules`): rules coverage, orphan rules, MFA alignment, Wire read-only constraints. | `definitionc`, `crud`, `CheckModule` |
| [`module.go`](../module.go) | Validates module-level constraints (`CheckModule`): module name requirement, entity name uniqueness, and delegates each entity to `CheckRules`. | `definitionc`, `crud` |
| [`lang.json`](../lang.json) | Translations for all validation error messages across 8 languages. | `webtyp.com/fmt` via `langc` |
| [`tests/rules_test.go`](../tests/rules_test.go) | Conformance test suite for entity rule checking invariants. | `gotest` |
| [`tests/module_test.go`](../tests/module_test.go) | Conformance test suite for module-level validation. | `gotest` |

## 3. Step-by-Step Flows

### Flow 1: Entity Rule Verification (`CheckRules`)

1. **Rule Coverage Check**: Iterates through standard CRUD actions (`Create`, `Read`, `Update`, `Delete`). If an action is present in `def.Allow` but not declared in `def.Rules`, returns `errRuleUndeclared` naming the entity and missing action.
2. **Orphan Rule Check**: If an action is declared in `def.Rules` but omitted from `def.Allow`, returns `errRuleNotAllowed`.
3. **MFA Alignment Check**: Computes `def.MFA &^ def.Allow`. If any MFA action is configured that is not allowed, returns `errMFANotAllowed`.
4. **Wire Model Invariant**: If `def.Wire` is true and `def.Allow` is non-zero, verifies that `def.Allow == model.Read` and `def.Rules.Read == model.AddMyBusinessRule`. Otherwise returns `errWireRead`.
5. **Success**: Returns `nil` if all checks pass.

### Flow 2: Module Verification (`CheckModule`)

1. **Module Name Check**: Ensures `m.Name` is not empty. Returns `errModuleNoName` if missing.
2. **Entity Validation Loop**: For each entity in `m.Entities`:
   - Checks that `entity.Name` is not empty; returns `errEntityNoName` if blank.
   - Scans preceding entities to ensure name uniqueness; returns `errEntityTwice` on collision.
   - Invokes `CheckRules(entity)`. If invalid, immediately returns the error.
3. **Success**: Returns `nil` when all entities and module metadata pass.

## 4. Security Boundaries & Contracts

- **Explicit Rules Only**: Every allowed CRUD action must explicitly specify either `model.AddMyBusinessRule` or `model.HaveNoBusinessRule`. Silence or omission is treated as a configuration defect, preventing unreviewed access.
- **No MFA Without Permissions**: MFA cannot be declared for actions the entity does not allow. This prevents misleading security claims in generated documentation or audits.
- **Wire Model Isolation**: Wire models represent external integrations or computed views without an underlying database table. They can only expose `Read` guarded by a custom business rule (`AddMyBusinessRule`).
- **Zero Allocations & TinyGo Safety**: The checker uses fixed arrays and bitwise operations instead of Go maps or reflection, ensuring zero runtime heap overhead and full compatibility with TinyGo WebAssembly.
- **One Revision, an integer, never the key**: `model.FieldDB.Revision` marks the column that counts a row's updates. At most one per entity, on an `Int` field, and it cannot be the primary key; otherwise `crud`'s conditional update has nothing unambiguous to compare.
