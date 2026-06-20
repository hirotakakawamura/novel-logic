package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/lean"
)

func TestShowListCommands(t *testing.T) {
	dir := writeCLIProject(t)
	cases := []struct {
		name    string
		args    []string
		contain string
	}{
		{"thing list", []string{"thing", "list"}, "hero"},
		{"fact list", []string{"fact", "list"}, "fact1"},
		{"action list", []string{"action", "list"}, "act1"},
		{"rule list", []string{"rule", "list"}, ""},
		{"time list", []string{"time", "list"}, "t1"},
		{"scene list", []string{"scene", "list"}, "scene1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runCLI(t, append([]string{"-C", dir}, tc.args...)...)
			if code != 0 {
				t.Fatalf("exit %d, output=%q", code, out)
			}
			if tc.contain != "" && !strings.Contains(out, tc.contain) {
				t.Fatalf("output=%q, want substring %q", out, tc.contain)
			}
			if tc.contain == "" && strings.TrimSpace(out) != "" {
				t.Fatalf("output=%q, want empty", out)
			}
		})
	}
}

func TestShowDetailCommands(t *testing.T) {
	dir := writeCLIProject(t)
	cases := []struct {
		name    string
		args    []string
		contain string
	}{
		{"thing show", []string{"thing", "show", "hero"}, "id: hero"},
		{"fact show", []string{"fact", "show", "fact1"}, "kind: state"},
		{"action show", []string{"action", "show", "act1"}, "at: t2"},
		{"scene show", []string{"scene", "show", "scene1"}, "id: scene1"},
		{"plot show", []string{"plot", "show"}, "scene1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runCLI(t, append([]string{"-C", dir}, tc.args...)...)
			if code != 0 {
				t.Fatalf("exit %d, output=%q", code, out)
			}
			if !strings.Contains(out, tc.contain) {
				t.Fatalf("output=%q", out)
			}
		})
	}
}

func TestTimelineShowsBranchActions(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "timeline", "--branch", "main", "--verbose")
	if code != 0 {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "act1") || !strings.Contains(out, "facts:") {
		t.Fatalf("output=%q", out)
	}
}

func TestStatusAndTemplateList(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "status")
	if code != 0 || !strings.Contains(out, "entities:") {
		t.Fatalf("status exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "template", "list")
	if code != 0 {
		t.Fatalf("template list exit %d", code)
	}
	for _, tpl := range []string{"default", "momotaro"} {
		if !strings.Contains(out, tpl) {
			t.Fatalf("output=%q, missing %q", out, tpl)
		}
	}
}

func TestDoctorMissingLean(t *testing.T) {
	dir := writeCLIProject(t)
	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", "")
	t.Cleanup(func() { _ = os.Setenv("PATH", oldPath) })

	_, code := runCLI(t, "-C", dir, "doctor")
	if code != 5 {
		t.Fatalf("exit %d, want 5", code)
	}
}

func TestNovelRevisionPinCLI(t *testing.T) {
	dir := writeCLIProject(t)
	gitInit(t, dir)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	gitCommitAll(t, dir, "add novel")

	out, code := runCLI(t, "-C", dir, "novel", "revision", "pin", "scene1", "--note", "test pin")
	if code != 0 {
		t.Fatalf("pin exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "pinned novel scene1") {
		t.Fatalf("output=%q", out)
	}

	out, code = runCLI(t, "-C", dir, "novel", "revision", "list", "scene1")
	if code != 0 || !strings.Contains(out, "current:") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}
}

func TestCheckStage2FailsOnBrokenLean(t *testing.T) {
	tc := lean.Detect()
	if !tc.Found {
		t.Skip("lean/lake not installed")
	}

	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "generate"); code != 0 {
		t.Fatal("generate failed")
	}
	theorems := filepath.Join(dir, "logic", "Theorems.lean")
	if err := os.WriteFile(theorems, []byte("syntax error here\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, code := runCLI(t, "-C", dir, "check", "--no-generate")
	if code != 3 {
		t.Fatalf("exit %d, want 3 for stage2 failure", code)
	}
}

func gitInit(t *testing.T, dir string) {
	t.Helper()
	runGitCLI(t, dir, "init")
	runGitCLI(t, dir, "config", "user.email", "test@example.com")
	runGitCLI(t, dir, "config", "user.name", "Test User")
	gitCommitAll(t, dir, "init")
}

func gitCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	runGitCLI(t, dir, "add", ".")
	runGitCLI(t, dir, "commit", "-m", msg)
}

func runGitCLI(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}