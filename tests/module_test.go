package definitioncheck_test

import (
	"strings"
	"testing"
	"webtyp.com/definitioncheck"
	"webtyp.com/model"
)

func TestCheckModule_Valid(t *testing.T) {
	m := model.ModuleDefinition{
		Name: "patient_directory",
		Entities: []model.Definition{
			patientDef(model.Create, model.Rules{Create: model.HaveNoBusinessRule}),
		},
	}
	if err := definitioncheck.CheckModule(m); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCheckModule_NameRequired(t *testing.T) {
	m := model.ModuleDefinition{Name: ""}
	if err := definitioncheck.CheckModule(m); err == nil {
		t.Error("expected error for missing name")
	}
}

func TestCheckModule_EntityTwice(t *testing.T) {
	m := model.ModuleDefinition{
		Name: "patient_directory",
		Entities: []model.Definition{
			patientDef(model.Create, model.Rules{Create: model.HaveNoBusinessRule}),
			patientDef(model.Read, model.Rules{Read: model.HaveNoBusinessRule}),
		},
	}
	err := definitioncheck.CheckModule(m)
	if err == nil {
		t.Fatal("expected error for entity twice")
	}
	if !strings.Contains(err.Error(), "patient") {
		t.Errorf("error %q does not contain 'patient'", err.Error())
	}
}

func TestCheckModule_EntityRulesChecked(t *testing.T) {
	m := model.ModuleDefinition{
		Name: "patient_directory",
		Entities: []model.Definition{
			patientDef(model.Update, model.Rules{}),
		},
	}
	err := definitioncheck.CheckModule(m)
	if err == nil {
		t.Fatal("expected error from entity rules")
	}
	if !strings.Contains(err.Error(), "patient") {
		t.Errorf("error %q does not contain 'patient'", err.Error())
	}
}
