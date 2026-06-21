package project

import (
	"errors"
	"testing"
)

func TestUpdateThingFactActionRule(t *testing.T) {
	d := newTestProject(t)
	if err := d.UpdateThing("hero", "Protagonist", []string{"lead"}, true); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateFact("fact1", FactState, "hero", "origin", "plot"); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateAction("act1", "hero", "origin", "departed", "t2", "plot", "leave"); err != nil {
		t.Fatal(err)
	}
	r, err := d.AddRule(RuleForbidState, "hero", "evil", "", "", MainBranch)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateRule(r.ID, RuleForbidState, "hero", "corrupt", "", ""); err != nil {
		t.Fatal(err)
	}

	t_, _ := d.FindThing("hero")
	f, _ := d.FindFact("fact1")
	a, _ := d.FindAction("act1")
	r2, _ := d.FindRule(r.ID)
	if t_ == nil || t_.Name != "Protagonist" {
		t.Fatalf("thing = %+v", t_)
	}
	if f == nil || f.Pred != "origin" {
		t.Fatalf("fact = %+v", f)
	}
	if a == nil || a.To != "departed" || a.Label != "leave" {
		t.Fatalf("action = %+v", a)
	}
	if r2 == nil || r2.Pred != "corrupt" {
		t.Fatalf("rule = %+v", r2)
	}
}

func TestSetPlot(t *testing.T) {
	d := newTestProject(t)
	if err := d.SetPlot("New Title", "New summary"); err != nil {
		t.Fatal(err)
	}
	if d.Meta.Title != "New Title" || d.Plot.Summary != "New summary" {
		t.Fatalf("meta=%+v plot=%+v", d.Meta, d.Plot)
	}
	if err := d.SetPlot("", ""); err == nil {
		t.Fatal("expected error for empty set")
	}
}

func TestUpdateFactRejectsDemotion(t *testing.T) {
	d := newTestProject(t)
	if err := d.UpdateFact("fact1", FactState, "hero", "origin", "plot"); err != nil {
		t.Fatal(err)
	}
	err := d.UpdateFact("fact1", FactFixed, "hero", "origin", "plot")
	if err == nil {
		t.Fatal("expected demotion error")
	}
	var reg *RegistrationError
	if !errors.As(err, &reg) {
		t.Fatalf("expected RegistrationError, got %T: %v", err, err)
	}
}

func TestUpdateErrors(t *testing.T) {
	d := newTestProject(t)
	if err := d.UpdateThing("ghost", "x", nil, false); err == nil {
		t.Fatal("missing thing")
	}
	if err := d.UpdateFact("ghost", FactState, "hero", "x", "plot"); err == nil {
		t.Fatal("missing fact")
	}
	if err := d.UpdateAction("ghost", "hero", "", "x", "t2", "plot", ""); err == nil {
		t.Fatal("missing action")
	}
	if err := d.UpdateRule("ghost", RuleForbidState, "hero", "x", "", ""); err == nil {
		t.Fatal("missing rule")
	}
}