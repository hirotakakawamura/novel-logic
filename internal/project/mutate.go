package project

import (
	"fmt"
	"strings"
	"time"
)

func NewID(prefix string, existing map[string]bool) string {
	for i := 1; i < 10000; i++ {
		id := fmt.Sprintf("%s%d", prefix, i)
		if !existing[id] {
			return id
		}
	}
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

func (d *Data) existingFactIDs() map[string]bool {
	m := make(map[string]bool, len(d.Facts))
	for _, f := range d.Facts {
		m[f.ID] = true
	}
	return m
}

func (d *Data) existingActionIDs() map[string]bool {
	m := make(map[string]bool, len(d.Actions))
	for _, a := range d.Actions {
		m[a.ID] = true
	}
	return m
}

func (d *Data) existingRuleIDs() map[string]bool {
	m := make(map[string]bool, len(d.Rules))
	for _, r := range d.Rules {
		m[r.ID] = true
	}
	return m
}

func MergeScope(scopes []string, scope string) []string {
	return MergeScopes(scopes, []string{scope})
}

// MergeScopes appends scopes not already present. Empty entries default to plot.
func MergeScopes(scopes []string, add []string) []string {
	if len(add) == 0 {
		add = []string{"plot"}
	}
	out := append([]string{}, scopes...)
	for _, s := range add {
		if s == "" {
			s = "plot"
		}
		found := false
		for _, existing := range out {
			if existing == s {
				found = true
				break
			}
		}
		if !found {
			out = append(out, s)
		}
	}
	return out
}

func normalizeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return []string{"plot"}
	}
	return scopes
}

func (d *Data) validateScopes(scopes []string) error {
	for _, scope := range normalizeScopes(scopes) {
		if err := validateScopeRef(d, scope); err != nil {
			return err
		}
	}
	return nil
}

func (d *Data) AddThing(id, name string, tags, scopes []string) error {
	if id == "" {
		return fmt.Errorf("thing id is required")
	}
	scopes = normalizeScopes(scopes)
	if err := d.validateScopes(scopes); err != nil {
		return err
	}
	if t, _ := d.FindThing(id); t != nil {
		return duplicateThingError(id)
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required for new thing")
	}
	d.Things = append(d.Things, Thing{
		ID:     id,
		Name:   name,
		Tags:   tags,
		Scopes: append([]string{}, scopes...),
	})
	return nil
}

// AddThingScopes appends scopes to an existing thing (or creates none — thing must exist).
func (d *Data) AddThingScopes(id string, scopes []string) error {
	if id == "" {
		return fmt.Errorf("thing id is required")
	}
	scopes = normalizeScopes(scopes)
	if err := d.validateScopes(scopes); err != nil {
		return err
	}
	t, _ := d.FindThing(id)
	if t == nil {
		return fmt.Errorf("thing %q not found", id)
	}
	t.Scopes = MergeScopes(t.Scopes, scopes)
	return nil
}

func (d *Data) AddTime(id, after string) error {
	if id == "" {
		return fmt.Errorf("time id is required")
	}
	for _, t := range d.Times {
		if t.ID == id {
			return fmt.Errorf("time %q already exists", id)
		}
	}
	for _, t := range d.Meta.TimeOrder {
		if t == id {
			return fmt.Errorf("time %q already in time_order", id)
		}
	}
	d.Times = append(d.Times, TimeEntry{ID: id})
	if after == "" {
		d.Meta.TimeOrder = append(d.Meta.TimeOrder, id)
		return nil
	}
	idx := d.TimeIndex(after)
	if idx < 0 {
		return fmt.Errorf("unknown time %q for --after", after)
	}
	order := append([]string{}, d.Meta.TimeOrder...)
	order = append(order[:idx+1], append([]string{id}, order[idx+1:]...)...)
	d.Meta.TimeOrder = order
	return nil
}

func (d *Data) AddScene(id, summary, timeStart, timeEnd string) error {
	if id == "" || summary == "" || timeStart == "" || timeEnd == "" {
		return fmt.Errorf("id, summary, time-start, and time-end are required")
	}
	for _, s := range d.Scenes {
		if s.ID == id {
			return fmt.Errorf("scene %q already exists", id)
		}
	}
	if d.TimeIndex(timeStart) < 0 {
		return fmt.Errorf("unknown time %q", timeStart)
	}
	if d.TimeIndex(timeEnd) < 0 {
		return fmt.Errorf("unknown time %q", timeEnd)
	}
	if !d.TimeLE(timeStart, timeEnd) {
		return fmt.Errorf("time_start after time_end")
	}
	d.Scenes = append(d.Scenes, Scene{
		ID:        id,
		Summary:   summary,
		TimeStart: timeStart,
		TimeEnd:   timeEnd,
	})
	return nil
}

func (d *Data) AddFact(kind FactKind, thing, pred, scope, branch string) (Fact, error) {
	if kind != FactFixed && kind != FactState {
		return Fact{}, fmt.Errorf("kind must be fixed or state")
	}
	if thing == "" || pred == "" {
		return Fact{}, fmt.Errorf("thing and pred are required")
	}
	if scope == "" {
		scope = "plot"
	}
	if !d.ThingIDs()[thing] {
		return Fact{}, fmt.Errorf("unknown thing %q", thing)
	}
	if d.ThingIDs()[pred] {
		return Fact{}, fmt.Errorf("pred %q matches existing thing id", pred)
	}
	branch = NormalizeBranch(branch)
	if err := d.validateBranchRef(branch); err != nil {
		return Fact{}, err
	}
	if err := validateScopeRef(d, scope); err != nil {
		return Fact{}, err
	}
	if existing, _ := d.FindFactByKey(kind, thing, pred, scope, branch); existing != nil {
		return Fact{}, duplicateFactError(*existing, kind, thing, pred, scope)
	}
	f := Fact{
		ID:     NewID("fact", d.existingFactIDs()),
		Kind:   kind,
		Thing:  thing,
		Pred:   pred,
		Scope:  scope,
		Branch: branch,
	}
	d.Facts = append(d.Facts, f)
	d.EnsureThingNovelScope(thing, scope)
	return f, nil
}

func (d *Data) PromoteFact(id string) error {
	for i := range d.Facts {
		if d.Facts[i].ID != id {
			continue
		}
		if d.Facts[i].Kind == FactState {
			return fmt.Errorf("fact %q is already state", id)
		}
		d.Facts[i].Kind = FactState
		return nil
	}
	return fmt.Errorf("fact %q not found", id)
}

func (d *Data) AddAction(thing, from, to, at, scope, label, branch string) (Action, error) {
	if thing == "" || to == "" || at == "" {
		return Action{}, fmt.Errorf("thing, to, and at are required")
	}
	if scope == "" {
		scope = "plot"
	}
	if !d.ThingIDs()[thing] {
		return Action{}, fmt.Errorf("unknown thing %q", thing)
	}
	if from != "" && d.ThingIDs()[from] {
		return Action{}, fmt.Errorf("from pred %q matches existing thing id", from)
	}
	if d.ThingIDs()[to] {
		return Action{}, fmt.Errorf("to pred %q matches existing thing id", to)
	}
	if d.TimeIndex(at) < 0 {
		return Action{}, fmt.Errorf("unknown time %q", at)
	}
	branch = NormalizeBranch(branch)
	if err := d.validateBranchRef(branch); err != nil {
		return Action{}, err
	}
	if merge := d.FindMergeForBranch(branch); merge != nil && d.branchClosed(branch) {
		return Action{}, registrationErrorf("branch %q is closed after merge %q; register on %q", branch, merge.ID, merge.IntoBranch)
	}
	if err := validateScopeRef(d, scope); err != nil {
		return Action{}, err
	}
	if existing, _ := d.FindActionByKey(thing, from, to, at, scope, branch); existing != nil {
		return Action{}, duplicateActionError(*existing, thing, from, to, at, scope)
	}
	a := Action{
		ID:     NewID("act", d.existingActionIDs()),
		Thing:  thing,
		From:   from,
		To:     to,
		At:     at,
		Scope:  scope,
		Label:  label,
		Branch: branch,
	}
	d.Actions = append(d.Actions, a)
	d.EnsureThingNovelScope(thing, scope)
	return a, nil
}

func (d *Data) AddRule(kind RuleKind, thing, pred, from, to, branch string) (Rule, error) {
	switch kind {
	case RuleForbidState:
		if thing == "" || pred == "" {
			return Rule{}, fmt.Errorf("forbid-state requires thing and pred")
		}
	case RuleForbidTransition:
		if from == "" || to == "" {
			return Rule{}, fmt.Errorf("forbid-transition requires from and to")
		}
	default:
		return Rule{}, fmt.Errorf("unknown rule kind %q", kind)
	}
	branch = NormalizeBranch(branch)
	if err := d.validateBranchRef(branch); err != nil {
		return Rule{}, err
	}
	if existing, _ := d.FindRuleByKey(kind, thing, pred, from, to, branch); existing != nil {
		return Rule{}, duplicateRuleError(*existing)
	}
	r := Rule{
		ID:     NewID("rule", d.existingRuleIDs()),
		Kind:   kind,
		Thing:  thing,
		Pred:   pred,
		From:   from,
		To:     to,
		Branch: branch,
	}
	d.Rules = append(d.Rules, r)
	return r, nil
}

func (d *Data) RecordCheck(success, stage1, stage2 bool, message string) error {
	now := time.Now().UTC()
	d.Meta.LastCheck = &CheckResult{
		At:      now,
		Success: success,
		Stage1:  stage1,
		Stage2:  stage2,
		Message: message,
	}
	return Save(d)
}

func validateScopeRef(d *Data, scope string) error {
	if scope == "plot" || scope == "" {
		return nil
	}
	if strings.HasPrefix(scope, "novel:") {
		sid := strings.TrimPrefix(scope, "novel:")
		for _, s := range d.Scenes {
			if s.ID == sid {
				return nil
			}
		}
		return fmt.Errorf("unknown scene %q in scope", sid)
	}
	return fmt.Errorf("invalid scope %q", scope)
}
