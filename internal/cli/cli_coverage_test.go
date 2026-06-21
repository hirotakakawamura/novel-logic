package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNovelListAndShow(t *testing.T) {
	dir := writeCLIProject(t)
	gitInitCLI(t, dir)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	gitCommitCLI(t, dir, "add novel")

	out, code := runCLI(t, "-C", dir, "novel", "list")
	if code != 0 || !strings.Contains(out, "scene1") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "scene2 [no novel]") {
		t.Fatalf("output=%q", out)
	}

	out, code = runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "body:") {
		t.Fatalf("show exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "alignment:") {
		t.Fatalf("output=%q", out)
	}
}

func TestNovelShowWithPinnedRevision(t *testing.T) {
	dir := writeCLIProject(t)
	gitInitCLI(t, dir)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	gitCommitCLI(t, dir, "add novel")
	if _, code := runCLI(t, "-C", dir, "novel", "revision", "pin", "scene1"); code != 0 {
		t.Fatal("pin failed")
	}

	out, code := runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "pinned_commit:") {
		t.Fatalf("show exit %d, output=%q", code, out)
	}
}

func TestRuleShow(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-state", "--thing", "hero", "--pred", "evil"); code != 0 {
		t.Fatal("rule add failed")
	}
	out, code := runCLI(t, "-C", dir, "rule", "show", "rule1")
	if code != 0 || !strings.Contains(out, "kind: forbid-state") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "rule", "list")
	if code != 0 || !strings.Contains(out, "rule1") {
		t.Fatalf("list exit %d, output=%q", code, out)
	}
}

func TestFactListFilters(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "fact", "list", "--kind", "state", "--thing", "hero")
	if code != 0 || !strings.Contains(out, "fact1") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "fact", "list", "--kind", "fixed")
	if code != 0 || !strings.Contains(out, "(no matches)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

func TestThingListTagFilter(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "thing", "list", "--tag", "character")
	if code != 0 || !strings.Contains(out, "hero") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "thing", "list", "--tag", "nope")
	if code != 0 || !strings.Contains(out, "(no matches)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

func TestSceneShowWithNovelLayer(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	if _, code := runCLI(t, "-C", dir, "fact", "add",
		"--kind", "state", "--thing", "hero", "--pred", "calm", "--scope", "novel:scene1"); code != 0 {
		t.Fatal("fact add failed")
	}
	out, code := runCLI(t, "-C", dir, "scene", "show", "scene1")
	if code != 0 || !strings.Contains(out, "related_things") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "layers: novel(facts=1") {
		t.Fatalf("output=%q", out)
	}
}

func gitInitCLI(t *testing.T, dir string) {
	t.Helper()
	runGitCLI(t, dir, "init")
	runGitCLI(t, dir, "config", "user.email", "test@example.com")
	runGitCLI(t, dir, "config", "user.name", "Test User")
	gitCommitCLI(t, dir, "init")
}

func gitCommitCLI(t *testing.T, dir, msg string) {
	t.Helper()
	runGitCLI(t, dir, "add", ".")
	runGitCLI(t, dir, "commit", "-m", msg)
}

func TestNovelListEmpty(t *testing.T) {
	dir := writeCLIProject(t)
	out, code := runCLI(t, "-C", dir, "novel", "list")
	if code != 0 || !strings.Contains(out, "(none)") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
}

// Ensure git helpers work when body path uses slash form.
func TestNovelShowBodyPath(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	body := filepath.Join(dir, "novels", "main", "scene1.txt")
	if err := os.WriteFile(body, []byte(strings.Repeat("x", 500)), 0o644); err != nil {
		t.Fatal(err)
	}
	out, code := runCLI(t, "-C", dir, "novel", "show", "scene1")
	if code != 0 || !strings.Contains(out, "…") {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	out, code = runCLI(t, "-C", dir, "novel", "show", "scene1", "--full")
	if code != 0 || !strings.Contains(out, strings.Repeat("x", 100)) {
		t.Fatalf("full exit %d", code)
	}
}