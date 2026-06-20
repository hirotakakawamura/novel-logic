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

func TestCheckForbidStateBranchScoped(t *testing.T) {
	d := &project.Data{
		Branches: []project.Branch{{ID: project.MainBranch}, {ID: "branch_a"}},
		Rules: []project.Rule{{
			ID: "rule1", Kind: project.RuleForbidState, Thing: "hero", Pred: "bad", Branch: "branch_a",
		}},
	}
	if err := CheckForbidState(d, project.MainBranch, "hero", "bad"); err != nil {
		t.Fatalf("main branch should not inherit branch_a rule: %v", err)
	}
	if err := CheckForbidState(d, "branch_a", "hero", "bad"); err == nil {
		t.Fatal("expected forbid-state error on branch_a")
	}
}

func TestCheckActionRulesBranchScoped(t *testing.T) {
	d := &project.Data{
		Branches: []project.Branch{{ID: project.MainBranch}, {ID: "branch_a"}},
		Rules: []project.Rule{{
			ID: "rule1", Kind: project.RuleForbidTransition, From: "mid", To: "start", Branch: "branch_a",
		}},
	}
	mainAction := project.Action{Thing: "hero", From: "mid", To: "start", At: "t1", Branch: project.MainBranch}
	if err := CheckActionRules(d, mainAction); err != nil {
		t.Fatalf("main branch should not inherit branch_a rule: %v", err)
	}
	altAction := project.Action{Thing: "hero", From: "mid", To: "start", At: "t1", Branch: "branch_a"}
	if err := CheckActionRules(d, altAction); err == nil {
		t.Fatal("expected forbid-transition error on branch_a")
	}
}