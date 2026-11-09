---
PLAN: "feat: definitioncheck — judge a model.Definition's Rules/MFA and a model.ModuleDefinition (the one judge for definitionc and crud)"
EXECUTOR: antigravity
REVIEWER: none
STATUS: running
SESSION: 11011781043384931358
---

> Read AGENTS.md at the repo root before any change.
> **Finish by yourself:** when the delivery check passes, commit after every stage, push and open the pull request. Never ask whether to
> proceed or finalize — nobody answers, and a session that waits ends without your work.
> Never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on your own work.

## Goal

`model` holds only the data a module declares (`Rules`, `MFA`, `ModuleDefinition`). This library holds the **judgement**: it refuses a declaration whose
`Rules` and `Allow` disagree, in one place that `definitionc` (when it generates) and `crud` (when it mounts) both call, so the rule "an allowed action
declares its rule" cannot be forgotten or judged twice differently. Isomorphic: no build tag, no `map`, `webtyp/fmt` only.

## Ecosystem search

Searched `CheckRules`, `CheckModule`, `Rules`, `ModuleDefinition` over `webtyp/*` and `veltylabs/*`: nothing judges these types yet (they are new in the `model` plan). Decision: **new library**
(api-design gate 4): the judge is not a contract (`model`), not a generator (`definitionc`) and not an engine (`crud`); two of those need it.

## Design gate

1. **Prior art** (unverified in this session): `go vet` analyzers and protobuf `protoc` validation — a declaration is judged by a tool separate from its data.
2. **Novice-name test**: `definitioncheck.CheckRules(def)` = "check this definition's rules"; `definitioncheck.CheckModule(m)` = "check this module".
3. **Ledger**: concepts +2 (two functions); judges of a declaration 2 (definitionc's own, crud's own) → 1; error texts leave `model`'s dictionary.
4. **Where**: its own library; imports only `model` and `fmt`.
5. **Deletes**: the template stub `definitioncheck.go` (`type Definitioncheck struct{}`, `New`) that `gonew` created; the `model` plan never ships a checker.

## Tools

```bash
go install webtyp.com/devflow/cmd/gotest@latest
go install webtyp.com/lang/cmd/langc@latest
```

Dependencies: `go get webtyp.com/model@v0.2.16 && go mod tidy` (`model` v0.2.16 is published and defines `Rules`, `RuleMode`, `ModuleDefinition`). No `replace`.

## APIs used

| symbol | exact signature | file |
|---|---|---|
| `model.Definition` | `Name string; Fields Fields; Wire bool; Allow Action; Rules Rules; MFA Action` | `model/definition.go` (model plan) |
| `model.Rules`, `model.RuleMode` | `Rules{Create, Read, Update, Delete RuleMode}`; `(Rules).Of(a Action) RuleMode`; `AddMyBusinessRule`=1, `HaveNoBusinessRule`=2, zero = undeclared | `model/rules.go` (model plan) |
| `model.ModuleDefinition` | `struct{ Name string; Entities []Definition }` | `model/module_definition.go` (model plan) |
| `model.Action` | `uint8` bit set; `Create, Read, Update, Delete`; `(Action).Has(want Action) bool`; `String()` → CRUD letters | `model/rbac.go:83` |
| `fmt.Err` | `func Err(msgs ...any) *Conv` (`webtyp.com/fmt`; `*Conv` is an `error`) | `fmt/error.go:11` |

If a signature differs from this table, stop and name the symbol in the PR; never adapt it by guessing.

## Stage 1 — failing tests first (`tests/rules_test.go`, `tests/module_test.go`, package `definitioncheck_test`)

```go
func patientDef(allow model.Action, rules model.Rules) model.Definition {
	return model.Definition{Name: "patient", Allow: allow, Rules: rules,
		Fields: model.Fields{{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}}}}
}
```

| Test (`tests/rules_test.go`) | setup → call → expected |
|---|---|
| `TestCheckRules_AllowedActionMustBeDeclared` | `patientDef(Create\|Update, Rules{Create: AddMyBusinessRule})` → error whose text contains `patient` and `u` |
| `TestCheckRules_DeclaredActionMustBeAllowed` | `patientDef(Create, Rules{Create: HaveNoBusinessRule, Delete: AddMyBusinessRule})` → error containing `patient` and `d` |
| `TestCheckRules_CompleteDeclarationPasses` | `patientDef(Create\|Read, Rules{Create: AddMyBusinessRule, Read: HaveNoBusinessRule})` → `nil` |
| `TestCheckRules_NothingAllowedNeedsNoRules` | `patientDef(0, Rules{})` → `nil` |
| `TestCheckRules_WireModelOnlyReadsThroughItsRule` | `Wire: true`: `Allow: Read, Rules{Read: AddMyBusinessRule}` → `nil`; `Allow: Read, Rules{Read: HaveNoBusinessRule}` → error containing the name; `Allow: Create, Rules{Create: AddMyBusinessRule}` → error; `Allow: 0` → `nil` |
| `TestCheckRules_MFAMustBeAllowed` | `Allow: Create\|Update, MFA: Update` → `nil`; `Allow: Create, MFA: Update` → error containing the name and `u`; `MFA: 0` → `nil` |

| Test (`tests/module_test.go`) | setup → call → expected |
|---|---|
| `TestCheckModule_Valid` | `{Name:"patient_directory", Entities: [patientDef(Create, Rules{Create: HaveNoBusinessRule})]}` → `nil` |
| `TestCheckModule_NameRequired` | `{Name:""}` → error |
| `TestCheckModule_EntityTwice` | two entities named `patient` → error containing `patient` |
| `TestCheckModule_EntityRulesChecked` | entity `patientDef(Update, Rules{})` → error containing `patient` (the entity's `CheckRules` error comes through) |

Run `gotest`: it must fail to compile (symbols missing). Commit: `test: definition checks (red)`.

## Stage 2 — the code (`rules.go`, `module.go`; no build tag, no `map`)

```go
package definitioncheck

import (
	"webtyp.com/fmt"
	"webtyp.com/model"
)

const (
	errRuleUndeclared = "an allowed action must declare its rule (AddMyBusinessRule or HaveNoBusinessRule) in Rules; add it to Rules: "
	errRuleNotAllowed = "Rules declares an action that Allow does not offer; add it to Allow or remove it from Rules: "
	errMFANotAllowed  = "MFA lists an action that Allow does not offer; add it to Allow or remove it from MFA: "
	errWireRead       = "a Wire model has no table, so it can only allow Read, answered by its own rule (Rules: Read: model.AddMyBusinessRule): "
)

// CheckRules refuses a Definition whose Rules and Allow disagree: an allowed action with no rule, a rule for an action
// that is not allowed, an MFA action that is not allowed, or a Wire model that is not exactly Read with its own rule.
// It names the model and the actions (as CRUD letters).
//
//	if err := definitioncheck.CheckRules(PatientModel); err != nil { return err }
func CheckRules(def model.Definition) error {
	var undeclared, orphan model.Action
	for _, a := range [...]model.Action{model.Create, model.Read, model.Update, model.Delete} {
		declared := def.Rules.Of(a) != 0
		allowed := def.Allow.Has(a)
		if allowed && !declared {
			undeclared |= a
		}
		if declared && !allowed {
			orphan |= a
		}
	}
	if undeclared != 0 {
		return fmt.Err(errRuleUndeclared + def.Name + " (" + undeclared.String() + ")")
	}
	if orphan != 0 {
		return fmt.Err(errRuleNotAllowed + def.Name + " (" + orphan.String() + ")")
	}
	if extra := def.MFA &^ def.Allow; extra != 0 {
		return fmt.Err(errMFANotAllowed + def.Name + " (" + extra.String() + ")")
	}
	if def.Wire && def.Allow != 0 && (def.Allow != model.Read || def.Rules.Read != model.AddMyBusinessRule) {
		return fmt.Err(errWireRead + def.Name)
	}
	return nil
}
```

```go
package definitioncheck

import (
	"webtyp.com/fmt"
	"webtyp.com/model"
)

const (
	errModuleNoName = "a module definition needs a Name"
	errEntityNoName = "an entity needs a Name; module: "
	errEntityTwice  = "the module lists the same entity twice: "
)

// CheckModule refuses a module with no name, an entity with no name, an entity listed twice, or an entity whose Rules and Allow
// disagree (CheckRules).
//
//	if err := definitioncheck.CheckModule(Definition); err != nil { return err }
func CheckModule(m model.ModuleDefinition) error {
	if m.Name == "" {
		return fmt.Err(errModuleNoName)
	}
	for i, entity := range m.Entities {
		if entity.Name == "" {
			return fmt.Err(errEntityNoName + m.Name)
		}
		for _, previous := range m.Entities[:i] {
			if previous.Name == entity.Name {
				return fmt.Err(errEntityTwice + entity.Name)
			}
		}
		if err := CheckRules(entity); err != nil {
			return err
		}
	}
	return nil
}
```

`git rm definitioncheck.go` (the `gonew` template stub). `gotest` passes. Commit: `feat: CheckRules and CheckModule`.

## Stage 3 — translations

The seven error constants are extracted by `langc`. Run `langc sync .`, fill every empty value for all the languages in `lang.json` (`es, zh, hi, ar, pt, fr, de, ru`, plus the columns
`lang.Supported()` adds), and fix every reuse warning by composing shared words at the call site. Delivery requires **zero untranslated keys and zero reuse warnings**. Commit: `chore: translations`.

## Stage 4 — documentation

- `README.md` (≤ 150 lines, the four ordered sections of the AGENTS.md template): *Why This Exists* with the Without/With (without: each tool judges a declaration its own way; with: one judge); *Step-by-Step Guide* with a compilable `CheckRules` + `CheckModule` example.
- `docs/ARCHITECTURE.md`: four sections; say that `definitionc` and `crud` call it and `model` holds only the data.
- `docs/TESTS.md`: one row per test of stage 1.
- Do not link `docs/PLAN.md` from permanent docs.

## Delivery check

1. `gotest` (builds, vets, runs; a PR that does not compile, or whose tests could not run, is a failed delivery).
2. `langc sync .` → zero untranslated, zero reuse warnings.
3. `grep -n "webtyp.com/" go.mod` lists only `model`, `fmt` (and what `langc` adds): no `crud`, no `httpd`, no `storage`.
4. `git diff --stat` lists only the files of the stages.

| Stage | Files | Commit |
|---|---|---|
| 1 | `tests/rules_test.go`, `tests/module_test.go` | `test: …(red)` |
| 2 | `rules.go`, `module.go`, delete `definitioncheck.go` | `feat: …` |
| 3 | `lang.json` | `chore: translations` |
| 4 | `README.md`, `docs/ARCHITECTURE.md`, `docs/TESTS.md` | `docs: …` |
