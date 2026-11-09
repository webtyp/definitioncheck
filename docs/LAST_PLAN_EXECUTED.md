---
PLAN: "feat: CheckRules also judges Revision (an integer, one per entity, never the primary key)"
EXECUTOR: jules
REVIEWER: none
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit after every stage, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work.
> Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work.

## Goal

`model.FieldDB.Revision` (model plan, published first) marks the column that counts a row's updates. A declaration with two of them, one on a text field, or one that is the primary key would make `crud`'s conditional update meaningless.
`CheckRules` (called by `definitionc` when it generates and by `crud` when it mounts) refuses those three at once. Nothing else changes.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
go get webtyp.com/model@v0.2.18
go mod tidy
```

## Stage 1 — failing tests first (`tests/revision_test.go`, package `definitioncheck_test`)

Use the helper `patientDef` of `tests/rules_test.go`; each test replaces its `Fields` with the given ones (a field `{Name: "revision", Type: model.Int(), DB: &model.FieldDB{Revision: true}}` is «a revision»).

| Test | setup → call → expected |
|---|---|
| `TestCheckRules_RevisionOnAnIntFieldPasses` | `Allow: Update`, `Rules{Update: HaveNoBusinessRule}`, fields `id` (PK text) + one revision → `nil` |
| `TestCheckRules_RevisionMustBeAnInteger` | the revision field has `Type: model.Text()` → error whose text contains `revision` and `integer` |
| `TestCheckRules_OnlyOneRevision` | two fields with `Revision: true` → error whose text contains `only one` |
| `TestCheckRules_RevisionCannotBeThePrimaryKey` | `DB: &model.FieldDB{PK: true, Revision: true}` on an Int field → error whose text contains `primary key` |
| `TestCheckRules_NoRevisionIsFine` | the existing complete declaration of `TestCheckRules_CompleteDeclarationPasses` → `nil` (no change) |

Add the five rows to `docs/TESTS.md` (`| [TestName](../tests/revision_test.go) | Rule |`, Rule = the name without `Test`, underscore as one space). Run `gotest`: it fails. Commit: `test: Revision declarations (red)`.

## Stage 2 — the code (`rules.go`)

Add constants next to the others:

```go
	errRevisionKind  = "a Revision field must be an integer (model.Int()): "
	errRevisionTwice = "an entity can have only one Revision field: "
	errRevisionKey   = "the primary key cannot be the Revision field: "
```

Add the function and call it at the end of `CheckRules`, just before `return nil`:

```go
// checkRevision refuses a Revision field that is not an integer, a second one, or one that is the primary key.
func checkRevision(def model.Definition) error {
	found := 0
	for _, f := range def.Fields {
		if !f.IsRevision() {
			continue
		}
		found++
		if f.Type.Storage() != model.FieldInt {
			return fmt.Err(errRevisionKind + def.Name + "." + f.Name)
		}
		if f.IsPK() {
			return fmt.Err(errRevisionKey + def.Name + "." + f.Name)
		}
		if found > 1 {
			return fmt.Err(errRevisionTwice + def.Name)
		}
	}
	return nil
}
```

(`return checkRevision(def)` replaces the final `return nil` of `CheckRules`.) `gotest` passes. Commit: `feat: judge Revision`.

## Stage 3 — translations and docs

The three new constants are extracted by `langc`: `langc sync .`, fill every empty value for all the languages of `lang.json`, fix every reuse warning (compose shared words at the call site). Zero untranslated keys, zero reuse warnings.
`README.md` and `docs/ARCHITECTURE.md`: one sentence saying `CheckRules` also judges `Revision`. Commit: `chore: translations and docs`.

## Delivery check

1. `gotest` passes. 2. `langc sync .` → zero untranslated, zero reuse warnings. 3. `git diff --stat main` lists only `rules.go`, `tests/revision_test.go`, `lang.json`, `docs/TESTS.md`, `README.md`, `docs/ARCHITECTURE.md`, `go.mod`, `go.sum`.
