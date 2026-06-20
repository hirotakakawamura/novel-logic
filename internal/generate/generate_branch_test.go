package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func branchRuleProject(t *testing.T) *project.Data {
	t.Helper()
	d := testfixture.LoadMinimal(t)
	if err := d.AddBranch("branch_a", "alt", project.MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := d.AddRule(project.RuleForbidTransition, "", "", "mid", "start", project.MainBranch); err != nil {
		t.Fatal(err)
	}
	if _, err := d.AddRule(project.RuleForbidState, "hero", "bad", "", "", "branch_a"); err != nil {
		t.Fatal(err)
	}
	return d
}

func TestRunBranchRuleSnapshot(t *testing.T) {
	d := branchRuleProject(t)
	if err := Run(d); err != nil {
		t.Fatal(err)
	}

	names := []string{"Rules.lean", "Theorems.lean"}
	for _, name := range names {
		gotPath := filepath.Join(d.Root, project.DirLogic, name)
		got, err := os.ReadFile(gotPath)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		goldenPath := filepath.Join("testdata", "branch_rule", name)
		if os.Getenv("UPDATE_GOLDEN") != "" {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}

		want, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("read golden %s: %v (run UPDATE_GOLDEN=1 go test ./internal/generate -run BranchRule)", goldenPath, err)
		}
		if string(got) != string(want) {
			t.Fatalf("snapshot mismatch for %s", name)
		}
	}

	rules, err := os.ReadFile(filepath.Join(d.Root, project.DirLogic, "Rules.lean"))
	if err != nil {
		t.Fatal(err)
	}
	body := string(rules)
	for _, want := range []string{"projectRules_main", "projectRules_branch_a", "forbiddenTransitions := [", "ThingId.hero", "PredId.bad"} {
		if !strings.Contains(body, want) {
			t.Fatalf("Rules.lean missing %q:\n%s", want, body)
		}
	}
}