// Design decision regression tests for GitHub issues #14, #16, #20.
// #14: plot scope skips time window; novel:<scene> strict (hints only for plot).
// #16: novel_extends_plot is Phase 1 — covered by docs, not tests here.
// #20: empty action from skips forbid-transition; forbid-state still checks to.
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

func TestCheckActionRulesEmptyFromStillEnforcesForbidState(t *testing.T) {
	d := &project.Data{
		Rules: []project.Rule{{
			ID: "rule1", Kind: project.RuleForbidState, Thing: "hero", Pred: "blocked",
		}},
	}
	a := project.Action{Thing: "hero", From: "", To: "blocked", At: "t1"}
	if err := CheckActionRules(d, a); err == nil {
		t.Fatal("expected forbid-state error on to even when from is empty")
	}
}

func TestRunEmptyFromSkipsForbidTransition(t *testing.T) {
	d := testfixture.LoadMinimal(t)
	d.Rules = append(d.Rules, project.Rule{
		ID: "rule_ft", Kind: project.RuleForbidTransition,
		From: "青年", To: "赤ちゃん", Branch: project.MainBranch,
	})
	d.Actions = append(d.Actions, project.Action{
		ID: "act_init", Thing: "hero", From: "", To: "赤ちゃん",
		At: "t1", Scope: "plot", Branch: project.MainBranch,
	})
	if hasIssueCode(Run(d), "rule.violation") {
		t.Fatalf("empty from should skip forbid-transition in Run(), got %v", Run(d))
	}
}

func TestRunEmptyFromEnforcesForbidState(t *testing.T) {
	d := testfixture.LoadMinimal(t)
	d.Rules = append(d.Rules, project.Rule{
		ID: "rule_fs", Kind: project.RuleForbidState,
		Thing: "hero", Pred: "blocked", Branch: project.MainBranch,
	})
	d.Actions = append(d.Actions, project.Action{
		ID: "act_bad", Thing: "hero", From: "", To: "blocked",
		At: "t1", Scope: "plot", Branch: project.MainBranch,
	})
	if !hasIssueCode(Run(d), "rule.violation") {
		t.Fatalf("expected rule.violation for forbid-state on to, got %v", Run(d))
	}
}