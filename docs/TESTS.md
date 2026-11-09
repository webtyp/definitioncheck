# Tests

| Test | Rule |
|---|---|
| [TestCheckRules_AllowedActionMustBeDeclared](../tests/rules_test.go) | CheckRules AllowedActionMustBeDeclared |
| [TestCheckRules_DeclaredActionMustBeAllowed](../tests/rules_test.go) | CheckRules DeclaredActionMustBeAllowed |
| [TestCheckRules_CompleteDeclarationPasses](../tests/rules_test.go) | CheckRules CompleteDeclarationPasses |
| [TestCheckRules_NothingAllowedNeedsNoRules](../tests/rules_test.go) | CheckRules NothingAllowedNeedsNoRules |
| [TestCheckRules_WireModelOnlyReadsThroughItsRule](../tests/rules_test.go) | CheckRules WireModelOnlyReadsThroughItsRule |
| [TestCheckRules_MFAMustBeAllowed](../tests/rules_test.go) | CheckRules MFAMustBeAllowed |
| [TestCheckModule_Valid](../tests/module_test.go) | CheckModule Valid |
| [TestCheckModule_NameRequired](../tests/module_test.go) | CheckModule NameRequired |
| [TestCheckModule_EntityTwice](../tests/module_test.go) | CheckModule EntityTwice |
| [TestCheckModule_EntityRulesChecked](../tests/module_test.go) | CheckModule EntityRulesChecked |
