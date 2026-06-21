package project

import (
	"errors"
	"strings"
	"testing"
)

func TestFactKeyIncludesBranch(t *testing.T) {
	k1 := FactKey(FactState, "hero", "mid", "plot", MainBranch)
	k2 := FactKey(FactState, "hero", "mid", "plot", "branch_a")
	if k1 == k2 {
		t.Fatal("fact keys must differ by branch")
	}
}

func TestNormalizeScopeDefaultsPlot(t *testing.T) {
	if got := FactKey(FactState, "hero", "mid", "", MainBranch); got != FactKey(FactState, "hero", "mid", "plot", MainBranch) {
		t.Fatalf("empty scope should normalize to plot in keys")
	}
}

func TestDuplicateRegistrationErrorsOnAdd(t *testing.T) {
	d := newTestProject(t)

	t.Run("fact", func(t *testing.T) {
		_, err := d.AddFact(FactState, "hero", "start", "plot", MainBranch)
		assertRegistrationError(t, err,
			"duplicate fact:",
			"kind=state, thing=hero, pred=start, scope=plot",
			"fact1",
			"novel-logic fact update fact1",
		)
	})

	t.Run("action_with_empty_from", func(t *testing.T) {
		if _, err := d.AddAction("hero", "", "spawn", "t3", "plot", "", MainBranch); err != nil {
			t.Fatal(err)
		}
		_, err := d.AddAction("hero", "", "spawn", "t3", "plot", "dup", MainBranch)
		assertRegistrationError(t, err,
			"duplicate action:",
			"from=-",
			"novel-logic action update",
		)
	})

	t.Run("action", func(t *testing.T) {
		_, err := d.AddAction("hero", "start", "mid", "t2", "plot", "dup", MainBranch)
		assertRegistrationError(t, err,
			"duplicate action:",
			"thing=hero, from=start, to=mid, at=t2, scope=plot",
			"act1",
			"novel-logic action update act1",
		)
	})

	t.Run("rule_forbid_state", func(t *testing.T) {
		if _, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", MainBranch); err != nil {
			t.Fatal(err)
		}
		_, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", MainBranch)
		assertRegistrationError(t, err,
			"duplicate rule: same forbid-state (thing=hero, pred=evil)",
			"novel-logic rule update",
		)
	})

	t.Run("rule_forbid_transition", func(t *testing.T) {
		if _, err := d.AddRule(RuleForbidTransition, "", "", "mid", "evil", MainBranch); err != nil {
			t.Fatal(err)
		}
		_, err := d.AddRule(RuleForbidTransition, "", "", "mid", "evil", MainBranch)
		assertRegistrationError(t, err,
			"duplicate rule: same forbid-transition (from=mid, to=evil)",
			"novel-logic rule update",
		)
	})

	t.Run("thing", func(t *testing.T) {
		err := d.AddThing("hero", "dup", []string{"x"}, nil)
		assertRegistrationError(t, err,
			`thing "hero" already exists`,
			"novel-logic thing scope add hero",
			"novel-logic thing update hero",
		)
	})

	t.Run("time", func(t *testing.T) {
		err := d.AddTime("t1", "")
		assertRegistrationError(t, err, `duplicate time "t1" already registered`)
	})
}

func TestDuplicateRegistrationErrorsOnUpdate(t *testing.T) {
	d := newTestProject(t)

	t.Run("fact", func(t *testing.T) {
		other, err := d.AddFact(FactState, "ally", "joined", "plot", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		err = d.UpdateFact("fact1", FactState, "ally", "joined", "plot")
		assertRegistrationError(t, err, "duplicate fact:", other.ID, "novel-logic fact update "+other.ID)
	})

	t.Run("action", func(t *testing.T) {
		act2, err := d.AddAction("hero", "mid", "end", "t3", "plot", "", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		err = d.UpdateAction(act2.ID, "hero", "start", "mid", "t2", "plot", "")
		assertRegistrationError(t, err, "duplicate action:", "act1", "novel-logic action update act1")
	})

	t.Run("rule_forbid_state", func(t *testing.T) {
		r1, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		r2, err := d.AddRule(RuleForbidState, "ally", "traitor", "", "", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		err = d.UpdateRule(r2.ID, RuleForbidState, "hero", "evil", "", "")
		assertRegistrationError(t, err, "duplicate rule: same forbid-state", r1.ID, "novel-logic rule update")
	})

	t.Run("rule_forbid_transition", func(t *testing.T) {
		r1, err := d.AddRule(RuleForbidTransition, "", "", "mid", "evil", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		r2, err := d.AddRule(RuleForbidTransition, "", "", "start", "mid", MainBranch)
		if err != nil {
			t.Fatal(err)
		}
		err = d.UpdateRule(r2.ID, RuleForbidTransition, "", "", "mid", "evil")
		assertRegistrationError(t, err, "duplicate rule: same forbid-transition", r1.ID, "novel-logic rule update")
	})
}

func TestDuplicateErrorHelpers(t *testing.T) {
	t.Run("duplicateFactError", func(t *testing.T) {
		err := duplicateFactError(Fact{ID: "fact1"}, FactState, "hero", "start", "")
		assertRegistrationError(t, err, "duplicate fact:", "scope=plot", "fact1")
	})

	t.Run("duplicateActionError_empty_from", func(t *testing.T) {
		err := duplicateActionError(Action{ID: "act1"}, "hero", "", "spawn", "t3", "")
		assertRegistrationError(t, err, "from=-", "act1")
	})

	t.Run("duplicateRuleError_forbid_state", func(t *testing.T) {
		err := duplicateRuleError(Rule{ID: "rule1", Kind: RuleForbidState, Thing: "hero", Pred: "evil"})
		assertRegistrationError(t, err, "forbid-state (thing=hero, pred=evil)", "rule1")
	})

	t.Run("duplicateRuleError_forbid_transition", func(t *testing.T) {
		err := duplicateRuleError(Rule{ID: "rule1", Kind: RuleForbidTransition, From: "mid", To: "evil"})
		assertRegistrationError(t, err, "forbid-transition (from=mid, to=evil)", "rule1")
	})

	t.Run("duplicateRuleError_default_kind", func(t *testing.T) {
		err := duplicateRuleError(Rule{ID: "rule_legacy", Kind: RuleKind("legacy")})
		assertRegistrationError(t, err, `duplicate rule "rule_legacy" already exists`, "novel-logic rule update rule_legacy")
	})

	t.Run("emptyDash", func(t *testing.T) {
		if got := emptyDash(""); got != "-" {
			t.Fatalf("emptyDash(\"\") = %q, want -", got)
		}
		if got := emptyDash("start"); got != "start" {
			t.Fatalf("emptyDash(\"start\") = %q", got)
		}
	})
}

func TestDuplicateIssues(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(t *testing.T, d *Data)
		wantSub string
	}{
		{
			name: "fact_same_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Facts = append(d.Facts, Fact{
					ID: "fact_dup", Kind: FactState, Thing: "hero", Pred: "start", Scope: "plot", Branch: MainBranch,
				})
			},
			wantSub: "duplicate fact: kind=state thing=hero pred=start scope=plot (fact1 and fact_dup)",
		},
		{
			name: "fact_empty_scope_normalizes",
			mutate: func(t *testing.T, d *Data) {
				d.Facts[0].Scope = ""
				d.Facts = append(d.Facts, Fact{
					ID: "fact_dup", Kind: FactState, Thing: "hero", Pred: "start", Scope: "plot", Branch: MainBranch,
				})
			},
			wantSub: "duplicate fact: kind=state thing=hero pred=start scope=plot",
		},
		{
			name: "action_with_empty_from",
			mutate: func(t *testing.T, d *Data) {
				d.Actions = append(d.Actions,
					Action{ID: "act_spawn", Thing: "hero", From: "", To: "solo", At: "t3", Scope: "plot", Branch: MainBranch},
					Action{ID: "act_dup", Thing: "hero", From: "", To: "solo", At: "t3", Scope: "plot", Branch: MainBranch},
				)
			},
			wantSub: "duplicate action: thing=hero from=- to=solo at=t3 scope=plot (act_spawn and act_dup)",
		},
		{
			name: "rule_forbid_state",
			mutate: func(t *testing.T, d *Data) {
				d.Rules = append(d.Rules,
					Rule{ID: "rule1", Kind: RuleForbidState, Thing: "hero", Pred: "evil", Branch: MainBranch},
					Rule{ID: "rule_dup", Kind: RuleForbidState, Thing: "hero", Pred: "evil", Branch: MainBranch},
				)
			},
			wantSub: "duplicate rule: forbid-state thing=hero pred=evil (rule1 and rule_dup)",
		},
		{
			name: "rule_forbid_transition",
			mutate: func(t *testing.T, d *Data) {
				d.Rules = append(d.Rules,
					Rule{ID: "rule1", Kind: RuleForbidTransition, From: "mid", To: "evil", Branch: MainBranch},
					Rule{ID: "rule_dup", Kind: RuleForbidTransition, From: "mid", To: "evil", Branch: MainBranch},
				)
			},
			wantSub: "duplicate rule: forbid-transition from=mid to=evil (rule1 and rule_dup)",
		},
		{
			name: "rule_unknown_kind",
			mutate: func(t *testing.T, d *Data) {
				d.Rules = append(d.Rules,
					Rule{ID: "rule1", Kind: RuleKind("legacy"), Branch: MainBranch},
					Rule{ID: "rule_dup", Kind: RuleKind("legacy"), Branch: MainBranch},
				)
			},
			wantSub: "duplicate rule: legacy (rule1 and rule_dup)",
		},
		{
			name: "thing_id",
			mutate: func(t *testing.T, d *Data) {
				d.Things = append(d.Things, Thing{ID: "hero", Tags: []string{"x"}, Scopes: []string{"plot"}})
			},
			wantSub: `duplicate thing id "hero"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestProject(t)
			tt.mutate(t, d)
			issues := DuplicateIssues(d)
			if len(issues) == 0 {
				t.Fatal("expected duplicate issue")
			}
			if !strings.Contains(issues[0], tt.wantSub) {
				t.Fatalf("issue = %q, want substring %q", issues[0], tt.wantSub)
			}
		})
	}
}

func TestDuplicateIssuesAllowsSameKeyOnDifferentBranches(t *testing.T) {
	d := newTestProject(t)
	d.Facts = append(d.Facts, Fact{
		ID: "fact_alt", Kind: FactState, Thing: "hero", Pred: "start", Scope: "plot", Branch: "branch_a",
	})
	if issues := DuplicateIssues(d); len(issues) != 0 {
		t.Fatalf("same key on different branches should be allowed, got %v", issues)
	}
}

func TestDuplicateIssuesSkipsEmptyThingID(t *testing.T) {
	d := newTestProject(t)
	d.Things = append(d.Things,
		Thing{ID: "", Tags: []string{"x"}, Scopes: []string{"plot"}},
		Thing{ID: "", Tags: []string{"y"}, Scopes: []string{"plot"}},
	)
	if issues := DuplicateIssues(d); len(issues) != 0 {
		t.Fatalf("empty thing ids should not duplicate, got %v", issues)
	}
}

func TestDuplicateIssuesCleanProject(t *testing.T) {
	d := newTestProject(t)
	if issues := DuplicateIssues(d); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func assertRegistrationError(t *testing.T, err error, wantSubs ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
	var reg *RegistrationError
	if !errors.As(err, &reg) {
		t.Fatalf("expected RegistrationError, got %T: %v", err, err)
	}
	for _, sub := range wantSubs {
		if !strings.Contains(reg.Error(), sub) {
			t.Fatalf("error %q missing substring %q", reg.Error(), sub)
		}
	}
}