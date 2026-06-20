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

func TestThingRemoveRejectsWhenReferenced(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "thing", "remove", "hero")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}

func TestForkMergeCLIEndToEnd(t *testing.T) {
	dir := writeCLIProject(t)

	mustOK := func(args ...string) {
		t.Helper()
		if _, code := runCLI(t, append([]string{"-C", dir}, args...)...); code != 0 {
			t.Fatalf("command failed: %v (exit %d)", args, code)
		}
	}

	mustOK("fork", "add", "fork1", "--parent", "main", "--at", "t2")
	mustOK("action", "add", "--thing", "hero", "--from", "mid", "--to", "route_a", "--at", "t2", "--label", "fork")

	forkAct := lastActionID(t, dir)
	mustOK("fork", "choice", "add", "--fork", "fork1", "--action", forkAct, "--branch", "branch_a")

	mustOK("action", "add", "--thing", "hero", "--from", "route_a", "--to", "merged", "--at", "t3", "--branch", "branch_a")
	actA := lastActionID(t, dir)
	mustOK("action", "add", "--thing", "hero", "--from", "mid", "--to", "merged", "--at", "t3", "--branch", "main")
	actMain := lastActionID(t, dir)

	mustOK("merge", "add", "merge1", "--at", "t3", "--into", "main",
		"--choice", "branch_a:"+actA, "--choice", "main:"+actMain)

	out, code := runCLI(t, "-C", dir, "validate")
	if code != 0 {
		t.Fatalf("validate exit %d, output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("output = %q", out)
	}

	out, code = runCLI(t, "-C", dir, "fork", "show", "fork1")
	if code != 0 || !strings.Contains(out, "branch_a") {
		t.Fatalf("fork show: code=%d output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "merge", "show", "merge1")
	if code != 0 || !strings.Contains(out, "into_branch: main") {
		t.Fatalf("merge show: code=%d output=%q", code, out)
	}
}

func TestCheckExitsWhenLeanMissing(t *testing.T) {
	dir := writeCLIProject(t)
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	out, code := runCLI(t, "-C", dir, "check")
	if code != 5 {
		t.Fatalf("exit code = %d, want 5; output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("expected stage1 OK before lean skip, output = %q", out)
	}
}

func lastActionID(t *testing.T, dir string) string {
	t.Helper()
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Actions) == 0 {
		t.Fatal("no actions")
	}
	return d.Actions[len(d.Actions)-1].ID
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