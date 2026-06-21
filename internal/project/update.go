package project

import "fmt"

func (d *Data) SetPlot(title, summary string) error {
	if title == "" && summary == "" {
		return fmt.Errorf("at least one of title or summary is required")
	}
	if title != "" {
		d.Meta.Title = title
	}
	if summary != "" {
		d.Plot.Summary = summary
	}
	return nil
}

func (d *Data) UpdateThing(id, name string, tags []string, replaceTags bool) error {
	if id == "" {
		return fmt.Errorf("thing id is required")
	}
	t, _ := d.FindThing(id)
	if t == nil {
		return fmt.Errorf("thing %q not found", id)
	}
	if name != "" {
		t.Name = name
	}
	if replaceTags {
		if len(tags) == 0 {
			return fmt.Errorf("at least one tag is required")
		}
		t.Tags = append([]string{}, tags...)
	}
	return nil
}

func (d *Data) UpdateFact(id string, kind FactKind, thing, pred, scope string) error {
	if id == "" {
		return fmt.Errorf("fact id is required")
	}
	f, idx := d.FindFact(id)
	if f == nil {
		return fmt.Errorf("fact %q not found", id)
	}
	if kind != FactFixed && kind != FactState {
		return fmt.Errorf("kind must be fixed or state")
	}
	if f.Kind == FactState && kind == FactFixed {
		return registrationErrorf("cannot demote state fact %q to fixed; demotion is not allowed", id)
	}
	if f.Kind == FactFixed && kind == FactState {
		return registrationErrorf("cannot promote fixed fact %q via update; use: novel-logic fact promote %s", id, id)
	}
	if thing == "" || pred == "" {
		return fmt.Errorf("thing and pred are required")
	}
	scope = normalizeScope(scope)
	if !d.ThingIDs()[thing] {
		return fmt.Errorf("unknown thing %q", thing)
	}
	if d.ThingIDs()[pred] {
		return fmt.Errorf("pred %q matches existing thing id", pred)
	}
	if err := validateScopeRef(d, scope); err != nil {
		return err
	}
	if other, _ := d.FindFactByKey(kind, thing, pred, scope, f.Branch); other != nil && other.ID != id {
		return duplicateFactError(*other, kind, thing, pred, scope)
	}
	d.Facts[idx].Kind = kind
	d.Facts[idx].Thing = thing
	d.Facts[idx].Pred = pred
	d.Facts[idx].Scope = scope
	d.EnsureThingNovelScope(thing, scope)
	return nil
}

func (d *Data) UpdateAction(id, thing, from, to, at, scope, label string) error {
	if id == "" {
		return fmt.Errorf("action id is required")
	}
	a, idx := d.FindAction(id)
	if a == nil {
		return fmt.Errorf("action %q not found", id)
	}
	if thing == "" || to == "" || at == "" {
		return fmt.Errorf("thing, to, and at are required")
	}
	scope = normalizeScope(scope)
	if !d.ThingIDs()[thing] {
		return fmt.Errorf("unknown thing %q", thing)
	}
	if from != "" && d.ThingIDs()[from] {
		return fmt.Errorf("from pred %q matches existing thing id", from)
	}
	if d.ThingIDs()[to] {
		return fmt.Errorf("to pred %q matches existing thing id", to)
	}
	if d.TimeIndex(at) < 0 {
		return fmt.Errorf("unknown time %q", at)
	}
	if err := validateScopeRef(d, scope); err != nil {
		return err
	}
	if other, _ := d.FindActionByKey(thing, from, to, at, scope, a.Branch); other != nil && other.ID != id {
		return duplicateActionError(*other, thing, from, to, at, scope)
	}
	d.Actions[idx].Thing = thing
	d.Actions[idx].From = from
	d.Actions[idx].To = to
	d.Actions[idx].At = at
	d.Actions[idx].Scope = scope
	d.Actions[idx].Label = label
	d.EnsureThingNovelScope(thing, scope)
	return nil
}

func (d *Data) UpdateRule(id string, kind RuleKind, thing, pred, from, to string) error {
	if id == "" {
		return fmt.Errorf("rule id is required")
	}
	r, idx := d.FindRule(id)
	if r == nil {
		return fmt.Errorf("rule %q not found", id)
	}
	switch kind {
	case RuleForbidState:
		if thing == "" || pred == "" {
			return fmt.Errorf("forbid-state requires thing and pred")
		}
		from = ""
		to = ""
	case RuleForbidTransition:
		if from == "" || to == "" {
			return fmt.Errorf("forbid-transition requires from and to")
		}
		thing = ""
		pred = ""
	default:
		return fmt.Errorf("unknown rule kind %q", kind)
	}
	if other, _ := d.FindRuleByKey(kind, thing, pred, from, to, r.Branch); other != nil && other.ID != id {
		return duplicateRuleError(*other)
	}
	d.Rules[idx].Kind = kind
	d.Rules[idx].Thing = thing
	d.Rules[idx].Pred = pred
	d.Rules[idx].From = from
	d.Rules[idx].To = to
	return nil
}
