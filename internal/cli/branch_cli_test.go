package cli

import (
	"strings"
	"testing"
)

func TestBranchList(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "branch", "list")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "main") {
		t.Fatalf("output = %q", out)
	}
}

func TestBranchShow(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "branch", "show", "main")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "id: main") {
		t.Fatalf("output = %q", out)
	}
}

func TestBranchRemoveRejectsChildBranch(t *testing.T) {
	dir := writeCLIProject(t)
	mustOK(t, dir, "branch", "add", "sub", "--parent", "main", "--label", "sub")
	mustOK(t, dir, "branch", "add", "leaf", "--parent", "sub", "--label", "leaf")

	_, code := runCLI(t, "-C", dir, "branch", "remove", "sub")
	if code != 4 {
		t.Fatalf("exit %d, want 4 when child branch exists", code)
	}
	_, code = runCLI(t, "-C", dir, "branch", "remove", "main")
	if code != 4 {
		t.Fatalf("exit %d, want 4 when removing main", code)
	}
}

func TestForkMergeCLIEndToEnd(t *testing.T) {
	dir := writeCLIProject(t)
	setupForkMergeCLI(t, dir)

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

func TestActionAddRejectsClosedBranchAfterMerge(t *testing.T) {
	dir := writeCLIProject(t)
	setupForkMergeCLI(t, dir)

	_, code := runCLI(t, "-C", dir, "action", "add",
		"--thing", "hero", "--from", "merged", "--to", "late", "--at", "t4", "--branch", "branch_a")
	if code != 1 {
		t.Fatalf("exit %d, want 1 when adding action on closed branch", code)
	}
}
