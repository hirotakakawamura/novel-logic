package project

import "testing"

func TestNewIDSkipsExisting(t *testing.T) {
	existing := map[string]bool{"fact1": true}
	if got := NewID("fact", existing); got != "fact2" {
		t.Fatalf("got %q, want fact2", got)
	}
}

func TestMergeScopesDedupesAndDefaults(t *testing.T) {
	got := MergeScopes([]string{"plot"}, []string{"", "novel:scene1", "plot"})
	want := []string{"plot", "novel:scene1"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestAddThingAndScopes(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddThing("villain", "Boss", []string{"character"}, []string{"plot"}); err != nil {
		t.Fatal(err)
	}
	if err := d.AddThingScopes("hero", []string{NovelScope("scene1")}); err != nil {
		t.Fatal(err)
	}
	t_, _ := d.FindThing("hero")
	if t_ == nil || len(t_.Scopes) < 2 {
		t.Fatalf("scopes = %v", t_.Scopes)
	}
}

func TestAddThingValidationErrors(t *testing.T) {
	d := newTestProject(t)
	cases := []struct {
		name string
		fn   func() error
	}{
		{"duplicate", func() error { return d.AddThing("hero", "", []string{"x"}, nil) }},
		{"no tags", func() error { return d.AddThing("new1", "", nil, nil) }},
		{"bad scope", func() error { return d.AddThing("new2", "", []string{"x"}, []string{"novel:ghost"}) }},
	}
	for _, tc := range cases {
		if err := tc.fn(); err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
	}
}

func TestAddTimeInsertAfter(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddTime("t1_5", "t1"); err != nil {
		t.Fatal(err)
	}
	if d.TimeIndex("t1_5") != 1 {
		t.Fatalf("time_order = %v", d.Meta.TimeOrder)
	}
}

func TestAddScene(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddScene("scene3", "close", "t3", "t4"); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := d.SceneWindow("scene3"); !ok {
		t.Fatal("scene3 not found")
	}
}

func TestAddFactPromoteAndAction(t *testing.T) {
	d := newTestProject(t)
	fixed, err := d.AddFact(FactFixed, "hero", "origin", "plot", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.PromoteFact(fixed.ID); err != nil {
		t.Fatal(err)
	}
	f, _ := d.FindFact(fixed.ID)
	if f == nil || f.Kind != FactState {
		t.Fatalf("fact = %+v", f)
	}
	act, err := d.AddAction("hero", "mid", "end", "t3", NovelScope("scene1"), "novel act", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if act.Scope != NovelScope("scene1") {
		t.Fatalf("scope = %q", act.Scope)
	}
	t_, _ := d.FindThing("hero")
	if t_ == nil || !containsScope(t_.Scopes, NovelScope("scene1")) {
		t.Fatal("EnsureThingNovelScope should add novel scope")
	}
}

func TestAddRule(t *testing.T) {
	d := newTestProject(t)
	r, err := d.AddRule(RuleForbidTransition, "", "", "mid", "evil", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if r.Kind != RuleForbidTransition {
		t.Fatalf("rule = %+v", r)
	}
}

func TestAddFactActionRuleErrors(t *testing.T) {
	d := newTestProject(t)
	if _, err := d.AddFact(FactState, "ghost", "x", "plot", MainBranch); err == nil {
		t.Fatal("unknown thing")
	}
	if _, err := d.AddAction("hero", "start", "hero", "t2", "plot", "", MainBranch); err == nil {
		t.Fatal("to pred collision")
	}
	if _, err := d.AddRule(RuleForbidState, "", "", "", "", MainBranch); err == nil {
		t.Fatal("incomplete forbid-state")
	}
	if err := d.PromoteFact("missing"); err == nil {
		t.Fatal("missing fact")
	}
}

func TestScenesContainingTimeAndPreds(t *testing.T) {
	d := newTestProject(t)
	ids := d.ScenesContainingTime("t2")
	if len(ids) != 2 {
		t.Fatalf("ids = %v", ids)
	}
	preds := d.Preds()
	if !preds["start"] || !preds["mid"] {
		t.Fatalf("preds = %v", preds)
	}
}

func TestBranchIDsAndEffectiveRules(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	ids := d.BranchIDs()
	if !ids[MainBranch] || !ids["branch_a"] {
		t.Fatalf("ids = %v", ids)
	}
	if _, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", "branch_a"); err != nil {
		t.Fatal(err)
	}
	rules := d.EffectiveRulesOnBranch("branch_a")
	if len(rules) == 0 {
		t.Fatal("expected rules on branch_a lineage")
	}
}

func TestAddBranchRemoveBranch(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveBranch("branch_a"); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveBranch(MainBranch); err == nil {
		t.Fatal("cannot remove main")
	}
}

func TestBranchIssuesUnknownBranch(t *testing.T) {
	d := newTestProject(t)
	d.Facts = append(d.Facts, Fact{
		ID: "fact_orphan", Kind: FactState, Thing: "hero", Pred: "lost",
		Scope: "plot", Branch: "nonexistent_branch",
	})
	issues := BranchIssues(d)
	found := false
	for _, iss := range issues {
		if iss.Code == "branch.unknown" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected branch.unknown, got %v", issues)
	}
}