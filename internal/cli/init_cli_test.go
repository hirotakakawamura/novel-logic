package cli

import (
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
)

func TestInitDefaultTemplate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "new-project")
	out, code := runCLI(t, "init", dir)
	if code != 0 {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	for _, name := range []string{
		project.FileProject, project.FileThings, project.FileBranches,
		project.FileForks, project.FileMerges,
	} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	out, code = runCLI(t, "-C", dir, "validate")
	if code != 0 {
		t.Fatalf("validate exit %d, output=%q", code, out)
	}
}

func TestInitRejectsNonEmptyDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "stray.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, code := runCLI(t, "init", dir)
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}

func TestInitUnknownTemplate(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "bad-tpl")
	_, code := runCLI(t, "init", dir, "--template", "no_such_template")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}