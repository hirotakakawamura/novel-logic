package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	out, code := runCLI(t, "version")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(out, "novel-logic") {
		t.Fatalf("output = %q", out)
	}
}

func TestValidateOK(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "validate")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("output = %q", out)
	}
}

func TestValidateQuiet(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "-q", "validate")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected no output, got %q", out)
	}
}

func TestValidateUnknownBranch(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "validate", "--branch", "no_such")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(out, "branch.unknown") {
		t.Fatalf("output = %q", out)
	}
}

func TestCheckQuick(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "check", "--quick")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("output = %q", out)
	}
}

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

func TestGenerate(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "generate")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "generated logic/") {
		t.Fatalf("output = %q", out)
	}
	for _, name := range []string{"Project.lean", "Facts.lean", "Theorems.lean", "lakefile.toml"} {
		path := filepath.Join(dir, "logic", name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestValidateWalkthroughBranchDog(t *testing.T) {
	dir := copyWalkthroughProject(t)
	out, code := runCLI(t, "-C", dir, "validate", "--branch", "branch_dog")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "OK: stage1") {
		t.Fatalf("output = %q", out)
	}
}

func TestValidateMissingProject(t *testing.T) {
	dir := t.TempDir()
	_, code := runCLI(t, "-C", dir, "validate")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}

func TestInfo(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "info")
	if code != 0 {
		t.Fatalf("exit code = %d, output = %q", code, out)
	}
	if !strings.Contains(out, "title: test") {
		t.Fatalf("output = %q", out)
	}
}