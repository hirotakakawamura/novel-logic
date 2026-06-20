package project

import "testing"

func TestFactKeyIncludesBranch(t *testing.T) {
	k1 := FactKey(FactState, "hero", "mid", "plot", MainBranch)
	k2 := FactKey(FactState, "hero", "mid", "plot", "branch_a")
	if k1 == k2 {
		t.Fatal("fact keys must differ by branch")
	}
}

func TestDuplicateIssuesDetectsSameBranch(t *testing.T) {
	d := newTestProject(t)
	dup := Fact{
		ID:     "fact_dup",
		Kind:   FactState,
		Thing:  "hero",
		Pred:   "start",
		Scope:  "plot",
		Branch: MainBranch,
	}
	d.Facts = append(d.Facts, dup)

	issues := DuplicateIssues(d)
	if len(issues) == 0 {
		t.Fatal("expected duplicate fact issue")
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