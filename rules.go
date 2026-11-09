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
