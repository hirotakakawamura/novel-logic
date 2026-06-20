package validate

import (
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func TestRunMinimalProjectOK(t *testing.T) {
	d := loadFixture(t)
	if issues := Run(d); len(issues) > 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestRunForBranchUnknown(t *testing.T) {
	d := loadFixture(t)
	issues := RunForBranch(d, "no_such_branch")
	if len(issues) == 0 {
		t.Fatal("expected branch.unknown issue")
	}
	found := false
	for _, iss := range issues {
		if iss.Code == "branch.unknown" {
			found = true
		}
	}
	if !found {
		t.Fatalf("issues = %v", issues)
	}
}

func TestRunDetectsMergeActionMismatch(t *testing.T) {
	d := loadFixture(t)
	actA, err := d.AddAction("hero", "mid", "ready_a", "t3", "plot", "", "branch_a")
	if err != nil {
		t.Fatal(err)
	}
	actB, err := d.AddAction("hero", "mid", "ready_b", "t3", "plot", "", project.MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	choices := []project.MergeChoice{
		{Branch: "branch_a", Action: actA.ID},
		{Branch: project.MainBranch, Action: actB.ID},
	}
	if err := d.AddMerge("merge_bad", "t3", "plot", project.MainBranch, choices); err != nil {
		t.Fatal(err)
	}
	issues := Run(d)
	found := false
	for _, iss := range issues {
		if iss.Code == "merge.action_mismatch" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected merge.action_mismatch, got %v", issues)
	}
}

func TestWalkthroughProjectValidates(t *testing.T) {
	root := filepath.Join("..", "..", "examples", "momotaro-walkthrough")
	d, err := project.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if issues := Run(d); len(issues) > 0 {
		t.Fatalf("walkthrough should validate clean, got %v", issues)
	}
	if issues := RunForBranch(d, "branch_dog"); len(issues) > 0 {
		t.Fatalf("branch_dog should validate clean, got %v", issues)
	}
}

func loadFixture(t *testing.T) *project.Data {
	return testfixture.LoadValidate(t)
}