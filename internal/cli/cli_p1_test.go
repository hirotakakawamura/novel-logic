package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"novel-logic/internal/project"
)

func TestThingUpdateName(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "thing", "update", "hero", "--name", "Protagonist")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	t_, _ := d.FindThing("hero")
	if t_ == nil || t_.Name != "Protagonist" {
		t.Fatalf("thing name not updated: %+v", t_)
	}
}

func TestFactUpdatePred(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "fact", "update", "fact1", "--pred", "origin")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	f, _ := d.FindFact("fact1")
	if f == nil || f.Pred != "origin" {
		t.Fatalf("fact pred not updated: %+v", f)
	}
}

func TestActionUpdateLabel(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "action", "update", "act1", "--label", "departure")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := d.FindAction("act1")
	if a == nil || a.Label != "departure" {
		t.Fatalf("action label not updated: %+v", a)
	}
}

func TestRuleUpdatePred(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-state", "--thing", "hero", "--pred", "evil"); code != 0 {
		t.Fatal("rule add failed")
	}
	_, code := runCLI(t, "-C", dir, "rule", "update", "rule1", "--pred", "corrupt")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	r, _ := d.FindRule("rule1")
	if r == nil || r.Pred != "corrupt" {
		t.Fatalf("rule pred not updated: %+v", r)
	}
}

func TestNovelUpdateSyncsSceneWindow(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatalf("novel add failed")
	}
	scenesPath := filepath.Join(dir, project.FileScenes)
	b, err := os.ReadFile(scenesPath)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(b), "time_end: t2", "time_end: t3", 1)
	if err := os.WriteFile(scenesPath, []byte(updated), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, code := runCLI(t, "-C", dir, "novel", "update", "scene1"); code != 0 {
		t.Fatalf("novel update exit %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	n, _ := d.FindNovel("scene1", project.MainBranch)
	if n == nil || n.TimeEnd != "t3" {
		t.Fatalf("novel TimeEnd = %q, want t3", n.TimeEnd)
	}
}

func TestNovelRemoveKeepsBody(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatalf("novel add failed")
	}
	bodyAbs := filepath.Join(dir, "novels", "main", "scene1.txt")
	if _, code := runCLI(t, "-C", dir, "novel", "remove", "scene1", "--keep-body"); code != 0 {
		t.Fatalf("novel remove exit %d", code)
	}
	if _, err := os.Stat(bodyAbs); err != nil {
		t.Fatalf("body file should remain: %v", err)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, idx := d.FindNovel("scene1", project.MainBranch); idx >= 0 {
		t.Fatal("novel metadata should be removed")
	}
}

func TestActionAndRuleRemove(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-transition", "--from", "a", "--to", "b"); code != 0 {
		t.Fatal("rule add failed")
	}
	if _, code := runCLI(t, "-C", dir, "action", "remove", "act1"); code != 0 {
		t.Fatalf("action remove exit %d", code)
	}
	if _, code := runCLI(t, "-C", dir, "rule", "remove", "rule1"); code != 0 {
		t.Fatalf("rule remove exit %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Actions) != 0 || len(d.Rules) != 0 {
		t.Fatalf("actions=%d rules=%d", len(d.Actions), len(d.Rules))
	}
}

func TestFactUpdateRejectsForbidState(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "rule", "add",
		"--kind", "forbid-state", "--thing", "hero", "--pred", "blocked"); code != 0 {
		t.Fatal("rule add failed")
	}
	_, code := runCLI(t, "-C", dir, "fact", "update", "fact1", "--pred", "blocked")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
}

func TestFactRemoveSuccess(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "fact", "remove", "fact1"); code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	d, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, idx := d.FindFact("fact1"); idx >= 0 {
		t.Fatal("fact should be removed")
	}
}

func TestSceneRemoveRejectsWhenNovelRegistered(t *testing.T) {
	dir := writeCLIProject(t)
	if _, code := runCLI(t, "-C", dir, "novel", "add", "scene1", "--init"); code != 0 {
		t.Fatal("novel add failed")
	}
	_, code := runCLI(t, "-C", dir, "scene", "remove", "scene1")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}

func TestTimeRemoveRejectsWhenReferenced(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "time", "remove", "t2")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
	}
}

func TestBranchRemoveRejectsChildBranch(t *testing.T) {
	dir := writeCLIProject(t)
	mustOK := func(args ...string) {
		t.Helper()
		if _, code := runCLI(t, append([]string{"-C", dir}, args...)...); code != 0 {
			t.Fatalf("failed: %v", args)
		}
	}
	mustOK("branch", "add", "sub", "--parent", "main", "--label", "sub")
	mustOK("branch", "add", "leaf", "--parent", "sub", "--label", "leaf")

	_, code := runCLI(t, "-C", dir, "branch", "remove", "sub")
	if code != 4 {
		t.Fatalf("exit %d, want 4 when child branch exists", code)
	}
	_, code = runCLI(t, "-C", dir, "branch", "remove", "main")
	if code != 4 {
		t.Fatalf("exit %d, want 4 when removing main", code)
	}
}

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

func TestValidateHintsPlotScene(t *testing.T) {
	dir := writeCLIProject(t)
	// act1 is plot scope at t2 inside scene1/scene2 windows — should hint
	out, code := runCLI(t, "-C", dir, "validate", "--verbose")
	if code != 0 {
		t.Fatalf("exit %d, output=%q", code, out)
	}
	if !strings.Contains(out, "[hint]") && !strings.Contains(out, "action.plot_scene_hint") {
		// hints print as [hint] message without code prefix in output
		if !strings.Contains(out, "Phase B alignment") {
			t.Fatalf("expected plot scene hint, output=%q", out)
		}
	}
}