package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsGitRepo(t *testing.T) {
	dir := newTestProject(t).Root
	if IsGitRepo(dir) {
		t.Fatal("expected false before git init")
	}
	runGitTest(t, dir, "init")
	if !IsGitRepo(dir) {
		t.Fatal("expected true after git init")
	}
}

func TestResolveGitFileState(t *testing.T) {
	dir := newTestProject(t).Root
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "test@example.com")
	runGitTest(t, dir, "config", "user.name", "Test User")

	bodyRel := DefaultNovelBodyPath("scene1", MainBranch)
	bodyPath := filepath.Join(dir, bodyRel)
	if err := os.MkdirAll(filepath.Dir(bodyPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bodyPath, []byte("first draft\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGitTest(t, dir, "add", ".")
	runGitTest(t, dir, "commit", "-m", "add novel body")

	state, err := ResolveGitFileState(dir, bodyRel)
	if err != nil {
		t.Fatal(err)
	}
	if state.Revision == "" || state.Short == "" {
		t.Fatalf("state = %+v, want revision metadata", state)
	}
	if state.Dirty {
		t.Fatal("expected clean working tree")
	}

	if err := os.WriteFile(bodyPath, []byte("revised draft\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	state, err = ResolveGitFileState(dir, bodyRel)
	if err != nil {
		t.Fatal(err)
	}
	if !state.Dirty {
		t.Fatal("expected dirty after edit")
	}
}

func TestResolveGitFileStateNotGitRepo(t *testing.T) {
	dir := newTestProject(t).Root
	_, err := ResolveGitFileState(dir, "novels/main/scene1.txt")
	if err == nil || !strings.Contains(err.Error(), "not a git repository") {
		t.Fatalf("err = %v", err)
	}
}

func TestResolveGitFileStateUntracked(t *testing.T) {
	dir := newTestProject(t).Root
	runGitTest(t, dir, "init")
	runGitTest(t, dir, "config", "user.email", "test@example.com")
	runGitTest(t, dir, "config", "user.name", "Test User")

	bodyRel := DefaultNovelBodyPath("scene1", MainBranch)
	bodyPath := filepath.Join(dir, bodyRel)
	if err := os.MkdirAll(filepath.Dir(bodyPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bodyPath, []byte("untracked\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveGitFileState(dir, bodyRel)
	if err == nil || !strings.Contains(err.Error(), "not tracked") {
		t.Fatalf("err = %v", err)
	}
}

func runGitTest(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}