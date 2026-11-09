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
| [TestCheckRules_RevisionOnAnIntFieldPasses](../tests/revision_test.go) | CheckRules RevisionOnAnIntFieldPasses |
| [TestCheckRules_RevisionMustBeAnInteger](../tests/revision_test.go) | CheckRules RevisionMustBeAnInteger |
| [TestCheckRules_OnlyOneRevision](../tests/revision_test.go) | CheckRules OnlyOneRevision |
| [TestCheckRules_RevisionCannotBeThePrimaryKey](../tests/revision_test.go) | CheckRules RevisionCannotBeThePrimaryKey |
| [TestCheckRules_NoRevisionIsFine](../tests/revision_test.go) | CheckRules NoRevisionIsFine |
