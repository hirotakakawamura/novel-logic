package cli

import (
	"strings"
	"testing"
)

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

func TestValidateHintsPlotScene(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "validate", "--verbose")
	if code != 0 {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "[hint]") && !strings.Contains(out, "action.plot_scene_hint") {
		if !strings.Contains(out, "Phase B alignment") {
			t.Fatalf("expected plot scene hint, output=%q", out)
		}
	}
}
