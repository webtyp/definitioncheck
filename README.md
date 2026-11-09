# definitioncheck
<img src="docs/img/badges.svg">

Judges a [`model.Definition`](https://github.com/webtyp/model)'s `Rules` and `MFA`, and a [`model.ModuleDefinition`](https://github.com/webtyp/model): the single judge called by [`definitionc`](https://github.com/webtyp/definitionc) and [`crud`](https://github.com/webtyp/crud).

## Why This Exists

In webtyp, the [`model`](https://github.com/webtyp/model) library defines data structures, but does not perform validation. Without a central judge, code generation tools like [`definitionc`](https://github.com/webtyp/definitionc) and runtime engines like [`crud`](https://github.com/webtyp/crud) would duplicate validation logic or enforce mismatched constraints.

`definitioncheck` provides the single authority that validates entity rules and module declarations according to [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

| Without definitioncheck | With definitioncheck |
|---|---|
| Generators and runtime engines write divergent rule validation. | One shared judge enforces identical rules at compile time and runtime. |
| An allowed action can easily miss a business rule declaration silently. | `CheckRules` refuses any allowed action lacking an explicit rule mode. |
| Duplicate entity names or malformed wire models slip into production. | `CheckModule` catches naming collisions and wire model misconfigurations. |

## Before You Start

- **Go version**: Go 1.24 or later.
- **Dependencies**: Only [`webtyp.com/model`](https://github.com/webtyp/model) and [`webtyp.com/fmt`](https://github.com/webtyp/fmt). No heavy database drivers, transports, or reflection.
- **Runtime target**: Fully isomorphic. Compatible with standard Go servers, TinyGo WebAssembly (browser), and Cloudflare Workers (edge).

## Step-by-Step Guide

### Step 1: Validate entity rules with CheckRules

Goal: Verify that an entity's declared rules and MFA actions strictly match its allowed CRUD operations.

```go
package main

import (
	"webtyp.com/definitioncheck"
	"webtyp.com/fmt"
	"webtyp.com/model"
)

func main() {
	patientDef := model.Definition{
		Name:  "patient",
		Allow: model.Create | model.Read,
		Rules: model.Rules{
			Create: model.AddMyBusinessRule,
			Read:   model.HaveNoBusinessRule,
		},
	}

	if err := definitioncheck.CheckRules(patientDef); err != nil {
		fmt.Println("Rule check failed:", err)
		return
	}

	fmt.Println("Patient definition rules passed!")
}
```

Expected result:
```
Patient definition rules passed!
```

### Step 2: Validate a module definition with CheckModule

Goal: Ensure a module has a valid name, contains no duplicate entity names, and all child definitions pass rule checks.

```go
package main

import (
	"webtyp.com/definitioncheck"
	"webtyp.com/fmt"
	"webtyp.com/model"
)

func main() {
	mod := model.ModuleDefinition{
		Name: "patient_directory",
		Entities: []model.Definition{
			{
				Name:  "patient",
				Allow: model.Read,
				Rules: model.Rules{Read: model.HaveNoBusinessRule},
			},
		},
	}

	if err := definitioncheck.CheckModule(mod); err != nil {
		fmt.Println("Module check failed:", err)
		return
	}

	fmt.Println("Module definition check passed!")
}
```

Expected result:
```
Module definition check passed!
```

## Documentation Map

| Order | Document | What question it answers & when to read it |
|---|---|---|
| 1 | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | How `definitioncheck` acts as the single judge and separates logic from data contracts. |
| 2 | [docs/TESTS.md](docs/TESTS.md) | Catalog of test cases and behaviors verified across rules and module checks. |
