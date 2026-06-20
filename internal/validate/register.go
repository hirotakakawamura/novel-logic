package validate

import (
	"fmt"

	"novel-logic/internal/project"
)

// CheckPredNotThingID rejects pred strings that match an existing thing id (§6.3.1).
func CheckPredNotThingID(d *project.Data, pred string) error {
	if pred == "" {
		return fmt.Errorf("pred is required")
	}
	if d.ThingIDs()[pred] {
		return fmt.Errorf("pred %q matches existing thing id", pred)
	}
	return nil
}

// CheckActionRules runs rule checks for a proposed action before persist.
func CheckActionRules(d *project.Data, a project.Action) error {
	if msg := checkActionRules(d, a); msg != "" {
		return fmt.Errorf("%s", msg)
	}
	return nil
}

// CheckForbidState runs forbid-state for a proposed state fact on branch.
func CheckForbidState(d *project.Data, branch, thing, pred string) error {
	if msg := checkForbidState(d, branch, thing, pred); msg != "" {
		return fmt.Errorf("%s", msg)
	}
	return nil
}