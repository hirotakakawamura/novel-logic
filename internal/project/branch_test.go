package project

import "testing"

func TestNormalizeBranch(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", MainBranch},
		{"main", MainBranch},
		{"branch_dog", "branch_dog"},
	}
	for _, tc := range tests {
		if got := NormalizeBranch(tc.in); got != tc.want {
			t.Errorf("NormalizeBranch(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestLoadEnsuresMainBranch(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, dir, FileProject, `title: test
time_order: [t1]
`)
	writeTestFile(t, dir, FileBranches, "[]\n")

	d, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, idx := d.FindBranchDef(MainBranch); idx < 0 {
		t.Fatal("expected implicit main branch after load")
	}
}

func TestBranchLineage(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddFork("fork1", MainBranch, "t2", "plot"); err != nil {
		t.Fatal(err)
	}
	actFork, err := d.AddAction("hero", "mid", "route_a", "t2", "plot", "pick a", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.AddForkChoice("fork1", actFork.ID, "branch_a"); err != nil {
		t.Fatal(err)
	}

	lineage := d.BranchLineage("branch_a")
	if len(lineage) != 2 || lineage[0] != MainBranch || lineage[1] != "branch_a" {
		t.Fatalf("lineage = %v, want [main branch_a]", lineage)
	}
}

func TestActiveActionsMainKeepsPostForkActions(t *testing.T) {
	d := newTestProject(t)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(d.AddFork("fork1", MainBranch, "t2", "plot"))
	actFork, err := d.AddAction("hero", "mid", "route_a", "t2", "plot", "fork", MainBranch)
	must(err)
	must(d.AddForkChoice("fork1", actFork.ID, "branch_a"))
	mainOnly, err := d.AddAction("ally", "start", "joined", "t3", "plot", "main only", MainBranch)
	must(err)
	altOnly, err := d.AddAction("ally", "start", "solo", "t3", "plot", "alt only", "branch_a")
	must(err)

	mainActs := actionIDs(d.ActiveActions(MainBranch))
	if !containsID(mainActs, "act1") {
		t.Fatalf("main should keep act1: %v", mainActs)
	}
	if !containsID(mainActs, mainOnly.ID) {
		t.Fatalf("main should include post-fork main action %q: %v", mainOnly.ID, mainActs)
	}

	altActs := actionIDs(d.ActiveActions("branch_a"))
	if containsID(altActs, mainOnly.ID) {
		t.Fatalf("branch_a should not include main-only action: %v", altActs)
	}
	if !containsID(altActs, altOnly.ID) {
		t.Fatalf("branch_a should include alt action %q: %v", altOnly.ID, altActs)
	}
	if !containsID(altActs, actFork.ID) {
		t.Fatalf("branch_a should include fork entry action: %v", altActs)
	}
}

func TestFindMergeForBranchIntoBranchNotClosed(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	actA, err := d.AddAction("hero", "route_a", "ready", "t3", "plot", "", "branch_a")
	if err != nil {
		t.Fatal(err)
	}
	actMain, err := d.AddAction("hero", "mid", "ready", "t3", "plot", "", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	choices := []MergeChoice{
		{Branch: "branch_a", Action: actA.ID},
		{Branch: MainBranch, Action: actMain.ID},
	}
	if err := d.AddMerge("merge1", "t3", "plot", MainBranch, choices); err != nil {
		t.Fatal(err)
	}
	if d.FindMergeForBranch(MainBranch) != nil {
		t.Fatal("into_branch main must not be treated as merged/closed")
	}
	if d.FindMergeForBranch("branch_a") == nil {
		t.Fatal("branch_a should have merge metadata")
	}
}

func TestAddForkExclusiveSlot(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddFork("fork1", MainBranch, "t2", "plot"); err != nil {
		t.Fatal(err)
	}
	err := d.AddFork("fork2", MainBranch, "t2", "plot")
	if err == nil {
		t.Fatal("expected error for duplicate fork slot")
	}
}

func TestMigrateNovelBodyPath(t *testing.T) {
	got := migrateNovelBodyPath("novels/scene1.txt", "scene1", MainBranch)
	want := DefaultNovelBodyPath("scene1", MainBranch)
	if got != want {
		t.Fatalf("migrate = %q, want %q", got, want)
	}
	legacy := "novels/custom.txt"
	if migrateNovelBodyPath(legacy, "scene1", MainBranch) != legacy {
		t.Fatalf("custom path should be preserved")
	}
}