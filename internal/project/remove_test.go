package project

import "testing"

func TestRemoveFactActionRule(t *testing.T) {
	d := newTestProject(t)
	if err := d.RemoveFact("fact1"); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveAction("act1"); err != nil {
		t.Fatal(err)
	}
	r, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveRule(r.ID); err != nil {
		t.Fatal(err)
	}
	if len(d.Facts) != 0 || len(d.Actions) != 0 || len(d.Rules) != 0 {
		t.Fatalf("facts=%d actions=%d rules=%d", len(d.Facts), len(d.Actions), len(d.Rules))
	}
}

func TestRemoveThingBlockedByRefs(t *testing.T) {
	d := newTestProject(t)
	if err := d.RemoveThing("hero"); err == nil {
		t.Fatal("expected reference error")
	}
}

func TestRemoveThingAfterClearingRefs(t *testing.T) {
	d := newTestProject(t)
	if err := d.RemoveFact("fact1"); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveAction("act1"); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveThing("hero"); err != nil {
		t.Fatal(err)
	}
	if _, idx := d.FindThing("hero"); idx >= 0 {
		t.Fatal("hero should be removed")
	}
}

func TestRemoveThingScopes(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddThingScopes("hero", []string{NovelScope("scene1")}); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveThingScopes("hero", []string{NovelScope("scene1")}); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveThingScopes("hero", []string{"plot"}); err == nil {
		t.Fatal("cannot remove last scope")
	}
}

func TestRemoveSceneGuards(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddNovel("scene1", MainBranch, "", true); err != nil {
		t.Fatal(err)
	}
	if err := d.RemoveScene("scene1"); err == nil {
		t.Fatal("novel registered")
	}
}

func TestRemoveTimeGuards(t *testing.T) {
	d := newTestProject(t)
	if err := d.RemoveTime("t2"); err == nil {
		t.Fatal("t2 referenced by scene/action")
	}
}