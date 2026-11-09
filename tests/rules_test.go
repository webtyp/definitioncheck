package definitioncheck_test

import (
	"strings"
	"testing"
	"webtyp.com/definitioncheck"
	"webtyp.com/model"
)

func patientDef(allow model.Action, rules model.Rules) model.Definition {
	return model.Definition{Name: "patient", Allow: allow, Rules: rules,
		Fields: model.Fields{{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}}}}
}

func patientDefWire(allow model.Action, rules model.Rules) model.Definition {
	def := patientDef(allow, rules)
	def.Wire = true
	return def
}

func patientDefMFA(allow model.Action, rules model.Rules, mfa model.Action) model.Definition {
	def := patientDef(allow, rules)
	def.MFA = mfa
	return def
}

func TestCheckRules_AllowedActionMustBeDeclared(t *testing.T) {
	def := patientDef(model.Create|model.Update, model.Rules{Create: model.AddMyBusinessRule})
	err := definitioncheck.CheckRules(def)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "patient") {
		t.Errorf("error %q does not contain 'patient'", msg)
	}
	if !strings.Contains(msg, "u") {
		t.Errorf("error %q does not contain 'u'", msg)
	}
}

func TestCheckRules_DeclaredActionMustBeAllowed(t *testing.T) {
	def := patientDef(model.Create, model.Rules{Create: model.HaveNoBusinessRule, Delete: model.AddMyBusinessRule})
	err := definitioncheck.CheckRules(def)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "patient") {
		t.Errorf("error %q does not contain 'patient'", msg)
	}
	if !strings.Contains(msg, "d") {
		t.Errorf("error %q does not contain 'd'", msg)
	}
}

func TestCheckRules_CompleteDeclarationPasses(t *testing.T) {
	def := patientDef(model.Create|model.Read, model.Rules{Create: model.AddMyBusinessRule, Read: model.HaveNoBusinessRule})
	err := definitioncheck.CheckRules(def)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCheckRules_NothingAllowedNeedsNoRules(t *testing.T) {
	def := patientDef(0, model.Rules{})
	err := definitioncheck.CheckRules(def)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCheckRules_WireModelOnlyReadsThroughItsRule(t *testing.T) {
	def1 := patientDefWire(model.Read, model.Rules{Read: model.AddMyBusinessRule})
	if err := definitioncheck.CheckRules(def1); err != nil {
		t.Errorf("expected nil error for Wire Read/AddMyBusinessRule, got %v", err)
	}

	def2 := patientDefWire(model.Read, model.Rules{Read: model.HaveNoBusinessRule})
	err2 := definitioncheck.CheckRules(def2)
	if err2 == nil {
		t.Error("expected error for Wire Read/HaveNoBusinessRule")
	} else if !strings.Contains(err2.Error(), "patient") {
		t.Errorf("error %q does not contain 'patient'", err2.Error())
	}

	def3 := patientDefWire(model.Create, model.Rules{Create: model.AddMyBusinessRule})
	if err := definitioncheck.CheckRules(def3); err == nil {
		t.Error("expected error for Wire Create")
	}

	def4 := patientDefWire(0, model.Rules{})
	if err := definitioncheck.CheckRules(def4); err != nil {
		t.Errorf("expected nil error for Wire 0, got %v", err)
	}
}

func TestCheckRules_MFAMustBeAllowed(t *testing.T) {
	def1 := patientDefMFA(model.Create|model.Update, model.Rules{Create: model.HaveNoBusinessRule, Update: model.HaveNoBusinessRule}, model.Update)
	if err := definitioncheck.CheckRules(def1); err != nil {
		t.Errorf("expected nil error for MFA Update with Allow Create|Update, got %v", err)
	}

	def2 := patientDefMFA(model.Create, model.Rules{Create: model.HaveNoBusinessRule}, model.Update)
	err2 := definitioncheck.CheckRules(def2)
	if err2 == nil {
		t.Error("expected error for MFA Update with Allow Create")
	} else {
		msg := err2.Error()
		if !strings.Contains(msg, "patient") {
			t.Errorf("error %q does not contain 'patient'", msg)
		}
		if !strings.Contains(msg, "u") {
			t.Errorf("error %q does not contain 'u'", msg)
		}
	}

	def3 := patientDefMFA(model.Create, model.Rules{Create: model.HaveNoBusinessRule}, 0)
	if err := definitioncheck.CheckRules(def3); err != nil {
		t.Errorf("expected nil error for MFA 0, got %v", err)
	}
}
