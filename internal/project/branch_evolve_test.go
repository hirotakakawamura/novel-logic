package project

import "testing"

func TestBranchIsolatedStateIssue(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	_, err := d.AddAction("hero", "mid", "alt_only", "t3", "plot", "", "branch_a")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.AddAction("hero", "alt_only", "end", "t4", "plot", "", MainBranch)
	if err != nil {
		t.Fatal(err)
	}

	issues := d.BranchIsolatedStateIssues()
	if len(issues) == 0 {
		t.Fatal("expected branch.isolated_state issue")
	}
}

func TestPredReachableOnBranch(t *testing.T) {
	d := newTestProject(t)
	if !d.predsReachableOnBranch(MainBranch, "hero", "start", "t2", "") {
		t.Fatal("start should be reachable at t2 from state fact")
	}
	if !d.predsReachableOnBranch(MainBranch, "hero", "mid", "t3", "") {
		t.Fatal("mid should be reachable at t3 after act1")
	}
}