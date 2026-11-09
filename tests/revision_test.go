package definitioncheck_test

import (
	"strings"
	"testing"

	"webtyp.com/definitioncheck"
	"webtyp.com/model"
)

func revisionField() model.Field {
	return model.Field{Name: "revision", Type: model.Int(), DB: &model.FieldDB{Revision: true}}
}

func defWithFields(fields ...model.Field) model.Definition {
	return model.Definition{
		Name: "patient", Allow: model.Update, Rules: model.Rules{Update: model.HaveNoBusinessRule},
		Fields: model.Fields(fields),
	}
}

func idField() model.Field {
	return model.Field{Name: "id", Type: model.Text(), DB: &model.FieldDB{PK: true}}
}

func TestCheckRules_RevisionOnAnIntFieldPasses(t *testing.T) {
	if err := definitioncheck.CheckRules(defWithFields(idField(), revisionField())); err != nil {
		t.Fatalf("a revision on an integer field must pass, got %v", err)
	}
}

func TestCheckRules_RevisionMustBeAnInteger(t *testing.T) {
	text := model.Field{Name: "revision", Type: model.Text(), DB: &model.FieldDB{Revision: true}}
	err := definitioncheck.CheckRules(defWithFields(idField(), text))
	if err == nil || !strings.Contains(err.Error(), "revision") || !strings.Contains(err.Error(), "integer") {
		t.Fatalf("want an error naming the field and 'integer', got %v", err)
	}
}

func TestCheckRules_OnlyOneRevision(t *testing.T) {
	second := model.Field{Name: "version", Type: model.Int(), DB: &model.FieldDB{Revision: true}}
	err := definitioncheck.CheckRules(defWithFields(idField(), revisionField(), second))
	if err == nil || !strings.Contains(err.Error(), "only one") {
		t.Fatalf("want an error saying 'only one', got %v", err)
	}
}

func TestCheckRules_RevisionCannotBeThePrimaryKey(t *testing.T) {
	key := model.Field{Name: "id", Type: model.Int(), DB: &model.FieldDB{PK: true, Revision: true}}
	err := definitioncheck.CheckRules(defWithFields(key))
	if err == nil || !strings.Contains(err.Error(), "primary key") {
		t.Fatalf("want an error mentioning the primary key, got %v", err)
	}
}

func TestCheckRules_NoRevisionIsFine(t *testing.T) {
	def := patientDef(model.Create|model.Read, model.Rules{Create: model.AddMyBusinessRule, Read: model.HaveNoBusinessRule})
	if err := definitioncheck.CheckRules(def); err != nil {
		t.Fatalf("a declaration without a revision must still pass, got %v", err)
	}
}
