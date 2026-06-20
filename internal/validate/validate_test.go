package validate

import (
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
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
	t.Helper()
	dir := t.TempDir()
	files := map[string]string{
		project.FileProject: `title: fixture
time_order: [t1, t2, t3]
`,
		project.FileThings: `- id: hero
  tags: [character]
  scopes: [plot]
`,
		project.FileTimes: `- id: t1
- id: t2
- id: t3
`,
		project.FileScenes: `- id: scene1
  summary: one
  time_start: t1
  time_end: t3
`,
		project.FileBranches: `- id: main
  label: main
`,
		project.FileForks:  "[]\n",
		project.FileMerges: "[]\n",
		project.FileFacts: `- id: fact1
  kind: state
  thing: hero
  pred: start
  scope: plot
`,
		project.FileActions: `- id: act1
  thing: hero
  from: start
  to: mid
  at: t2
  scope: plot
`,
		project.FileRules:  "[]\n",
		project.FileNovels: "[]\n",
	}
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.AddBranch("branch_a", "alt", project.MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	return d
}