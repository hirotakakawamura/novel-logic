package validate

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func TestRunDetectsNovelRevisionDrift(t *testing.T) {
	dir := testfixture.WriteMinimalDir(t)
	gitInitCommitAll(t, dir, "init project")

	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.AddNovel("scene1", project.MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	if err := project.Save(d); err != nil {
		t.Fatal(err)
	}
	gitAddCommit(t, dir, "add novel")

	d, err = project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.PinNovelRevision("scene1", project.MainBranch, "", "initial pin", false); err != nil {
		t.Fatal(err)
	}
	if err := project.Save(d); err != nil {
		t.Fatal(err)
	}
	gitAddCommit(t, dir, "pin revision")

	bodyPath := filepath.Join(dir, "novels", "main", "scene1.txt")
	if err := os.WriteFile(bodyPath, []byte("revised prose\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitAddCommit(t, dir, "revise body")

	d, err = project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	issues := Run(d)
	found := false
	for _, iss := range issues {
		if iss.Code == "novel.revision_drift" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected novel.revision_drift, got %v", issues)
	}
}

func gitInitCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", msg)
}

func gitAddCommit(t *testing.T, dir, msg string) {
	t.Helper()
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", msg)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}