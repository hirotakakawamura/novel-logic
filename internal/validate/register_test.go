package validate

import (
	"testing"

	"novel-logic/internal/project"
)

func TestCheckPredNotThingIDRejectsThingCollision(t *testing.T) {
	d := &project.Data{
		Things: []project.Thing{{ID: "hero", Tags: []string{"character"}, Scopes: []string{"plot"}}},
	}
	if err := CheckPredNotThingID(d, "hero"); err == nil {
		t.Fatal("expected collision error")
	}
}

func TestCheckActionRulesForbidTransition(t *testing.T) {
	d := &project.Data{
		Rules: []project.Rule{{
			ID: "rule1", Kind: project.RuleForbidTransition, From: "mid", To: "start",
		}},
	}
	err := CheckActionRules(d, project.Action{Thing: "hero", From: "mid", To: "start", At: "t1"})
	if err == nil {
		t.Fatal("expected forbid-transition error")
	}
}