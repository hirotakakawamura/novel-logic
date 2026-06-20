package project

import (
	"fmt"
	"strings"
)

func (d *Data) RemoveThing(id string) error {
	if id == "" {
		return fmt.Errorf("thing id is required")
	}
	if _, idx := d.FindThing(id); idx < 0 {
		return fmt.Errorf("thing %q not found", id)
	}
	var refs []string
	for _, f := range d.Facts {
		if f.Thing == id {
			refs = append(refs, "fact:"+f.ID)
		}
	}
	for _, a := range d.Actions {
		if a.Thing == id {
			refs = append(refs, "action:"+a.ID)
		}
	}
	for _, r := range d.Rules {
		if r.Thing == id {
			refs = append(refs, "rule:"+r.ID)
		}
	}
	if len(refs) > 0 {
		return fmt.Errorf("thing %q is referenced by %s", id, strings.Join(refs, ", "))
	}
	d.Things = removeThingByID(d.Things, id)
	return nil
}

func removeThingByID(things []Thing, id string) []Thing {
	out := things[:0]
	for _, t := range things {
		if t.ID != id {
			out = append(out, t)
		}
	}
	return out
}

func (d *Data) RemoveThingScopes(id string, scopes []string) error {
	if id == "" {
		return fmt.Errorf("thing id is required")
	}
	if len(scopes) == 0 {
		return fmt.Errorf("at least one scope is required")
	}
	t, _ := d.FindThing(id)
	if t == nil {
		return fmt.Errorf("thing %q not found", id)
	}
	remove := map[string]bool{}
	for _, s := range scopes {
		if s == "" {
			s = "plot"
		}
		if err := validateScopeRef(d, s); err != nil {
			return err
		}
		remove[s] = true
	}
	var kept []string
	for _, s := range t.Scopes {
		if !remove[s] {
			kept = append(kept, s)
		}
	}
	if len(kept) == 0 {
		return fmt.Errorf("cannot remove all scopes from thing %q", id)
	}
	t.Scopes = kept
	return nil
}

func (d *Data) RemoveScene(id string) error {
	if id == "" {
		return fmt.Errorf("scene id is required")
	}
	found := false
	for _, s := range d.Scenes {
		if s.ID == id {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("scene %q not found", id)
	}
	novelScope := NovelScope(id)
	var refs []string
	for _, n := range d.Novels {
		if n.SceneID == id {
			refs = append(refs, "novel:"+id)
		}
	}
	for _, f := range d.Facts {
		if f.Scope == novelScope {
			refs = append(refs, "fact:"+f.ID)
		}
	}
	for _, a := range d.Actions {
		if a.Scope == novelScope {
			refs = append(refs, "action:"+a.ID)
		}
	}
	for _, t := range d.Things {
		if containsScope(t.Scopes, novelScope) {
			refs = append(refs, "thing:"+t.ID+" scope")
		}
	}
	if len(refs) > 0 {
		return fmt.Errorf("scene %q is referenced by %s (remove those first)", id, strings.Join(refs, ", "))
	}
	out := d.Scenes[:0]
	for _, s := range d.Scenes {
		if s.ID != id {
			out = append(out, s)
		}
	}
	d.Scenes = out
	return nil
}

func containsScope(scopes []string, target string) bool {
	for _, s := range scopes {
		if s == target {
			return true
		}
	}
	return false
}

func (d *Data) RemoveTime(id string) error {
	if id == "" {
		return fmt.Errorf("time id is required")
	}
	if d.TimeIndex(id) < 0 {
		return fmt.Errorf("time %q not found", id)
	}
	var refs []string
	for _, s := range d.Scenes {
		if s.TimeStart == id || s.TimeEnd == id {
			refs = append(refs, "scene:"+s.ID)
		}
	}
	for _, a := range d.Actions {
		if a.At == id {
			refs = append(refs, "action:"+a.ID)
		}
	}
	for _, n := range d.Novels {
		if n.TimeStart == id || n.TimeEnd == id {
			refs = append(refs, "novel:"+n.SceneID)
		}
	}
	if len(refs) > 0 {
		return fmt.Errorf("time %q is referenced by %s", id, strings.Join(refs, ", "))
	}
	order := d.Meta.TimeOrder[:0]
	for _, t := range d.Meta.TimeOrder {
		if t != id {
			order = append(order, t)
		}
	}
	d.Meta.TimeOrder = order
	entries := d.Times[:0]
	for _, e := range d.Times {
		if e.ID != id {
			entries = append(entries, e)
		}
	}
	d.Times = entries
	return nil
}

func (d *Data) RemoveFact(id string) error {
	if id == "" {
		return fmt.Errorf("fact id is required")
	}
	found := false
	out := d.Facts[:0]
	for _, f := range d.Facts {
		if f.ID == id {
			found = true
			continue
		}
		out = append(out, f)
	}
	if !found {
		return fmt.Errorf("fact %q not found", id)
	}
	d.Facts = out
	return nil
}

func (d *Data) RemoveAction(id string) error {
	if id == "" {
		return fmt.Errorf("action id is required")
	}
	found := false
	out := d.Actions[:0]
	for _, a := range d.Actions {
		if a.ID == id {
			found = true
			continue
		}
		out = append(out, a)
	}
	if !found {
		return fmt.Errorf("action %q not found", id)
	}
	d.Actions = out
	return nil
}

func (d *Data) RemoveRule(id string) error {
	if id == "" {
		return fmt.Errorf("rule id is required")
	}
	found := false
	out := d.Rules[:0]
	for _, r := range d.Rules {
		if r.ID == id {
			found = true
			continue
		}
		out = append(out, r)
	}
	if !found {
		return fmt.Errorf("rule %q not found", id)
	}
	d.Rules = out
	return nil
}

