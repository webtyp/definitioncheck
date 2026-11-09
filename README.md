# definitioncheck
<img src="docs/img/badges.svg">

Judges a model.Definition's Rules and MFA and a model.ModuleDefinition: the one judge definitionc and crud call.

## Why This Exists

Without this library, each tool judges a declaration its own way (or misses rules entirely).
With this library, there is one unified judge for `definitionc` (when it generates) and `crud` (when it mounts), so the rule "an allowed action declares its rule" cannot be forgotten or judged twice differently.

## Step-by-Step Guide

```go
package main

import (
	"webtyp.com/definitioncheck"
	"webtyp.com/model"
	"fmt"
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

	m := model.ModuleDefinition{
		Name: "patient_directory",
		Entities: []model.Definition{patientDef},
	}

	if err := definitioncheck.CheckModule(m); err != nil {
		fmt.Println("Module check failed:", err)
		return
	}

	fmt.Println("Definitions check passed!")
}
```
