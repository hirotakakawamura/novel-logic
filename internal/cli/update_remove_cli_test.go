package cli

import (
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

func TestFactUpdateRejectsPromoteViaUpdate(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "fact", "add",
		"--kind", "fixed", "--thing", "ally", "--pred", "companion")
	if code != 0 {
		t.Fatalf("fact add exit %d", code)
	}
	_, code = runCLI(t, "-C", dir, "fact", "update", "fact2", "--kind", "state")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1 for fixed→state via update", code)
	}
	_, code = runCLI(t, "-C", dir, "fact", "promote", "fact2")
	if code != 0 {
		t.Fatalf("promote exit %d", code)
	}
}

func TestFactUpdateRejectsDemoteToFixed(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "fact", "update", "fact1", "--kind", "fixed")
	if code != 1 {
		t.Fatalf("exit code = %d, want 1 for state→fixed demotion", code)
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

func TestThingRemoveRejectsWhenReferenced(t *testing.T) {
	dir := writeCLIProject(t)
	_, code := runCLI(t, "-C", dir, "thing", "remove", "hero")
	if code != 4 {
		t.Fatalf("exit code = %d, want 4", code)
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
