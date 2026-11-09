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
