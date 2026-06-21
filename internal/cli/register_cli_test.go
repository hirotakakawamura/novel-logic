package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
)

func TestFactAddRejectsForbidState(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-state", "--thing", "hero", "--pred", "bad"); code != 0 {
		t.Fatalf("rule add failed")
	}
	_, code := runCLI(t, "-C", dir, "fact", "add",
		"--kind", "state", "--thing", "hero", "--pred", "bad")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}

func TestActionAddRejectsForbidTransition(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-transition", "--from", "mid", "--to", "start"); code != 0 {
		t.Fatalf("rule add failed")
	}
	_, code := runCLI(t, "-C", dir, "action", "add",
		"--thing", "hero", "--from", "mid", "--to", "start", "--at", "t2")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}

func TestFactAddPersistsOnSuccess(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "fact", "add",
		"--kind", "fixed", "--thing", "ally", "--pred", "companion")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range d.Facts {
		if f.Thing == "ally" && f.Pred == "companion" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("fact not persisted")
	}
	path := filepath.Join(dir, project.FileFacts)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "companion") {
		t.Fatalf("facts.yaml missing new fact: %s", b)
	}
}