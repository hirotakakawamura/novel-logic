package validate

import (
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func TestPlotScopeActionSkipsTimeWindowCheck(t *testing.T) {
	d := testfixture.LoadMinimal(t)
	// act1 is plot-scoped at t2; move to t4 (outside scene1 window t1..t2).
	d.Actions[0].Scope = "plot"
	d.Actions[0].At = "t4"
	if hasIssueCode(Run(d), "time.action_window") {
		t.Fatalf("plot scope action should not trigger time.action_window, got %v", Run(d))
	}
}

func TestNovelScopeActionEnforcesTimeWindow(t *testing.T) {
	d := testfixture.LoadMinimal(t)
	d.Actions[0].Scope = "novel:scene1"
	d.Actions[0].At = "t4"
	if !hasIssueCode(Run(d), "time.action_window") {
		t.Fatalf("expected time.action_window, got %v", Run(d))
	}
}

func TestCheckActionRulesEmptyFromSkipsForbidTransition(t *testing.T) {
	d := &project.Data{
		Rules: []project.Rule{{
			ID: "rule1", Kind: project.RuleForbidTransition, From: "青年", To: "赤ちゃん",
		}},
	}
	// Empty from: initial transition to 赤ちゃん is not the forbidden 青年→赤ちゃん pair.
	a := project.Action{Thing: "hero", From: "", To: "赤ちゃん", At: "t1"}
	if err := CheckActionRules(d, a); err != nil {
		t.Fatalf("empty from should skip forbid-transition: %v", err)
	}
	withFrom := project.Action{Thing: "hero", From: "青年", To: "赤ちゃん", At: "t1"}
	if err := CheckActionRules(d, withFrom); err == nil {
		t.Fatal("expected forbid-transition error when from is set")
	}
}