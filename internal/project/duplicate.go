package project

import "fmt"

func normalizeScope(scope string) string {
	if scope == "" {
		return "plot"
	}
	return scope
}

// FactKey is the uniqueness key for a fact within a branch.
func FactKey(kind FactKind, thing, pred, scope, branch string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", kind, thing, pred, normalizeScope(scope), NormalizeBranch(branch))
}

func (d *Data) FindFactByKey(kind FactKind, thing, pred, scope, branch string) (*Fact, int) {
	key := FactKey(kind, thing, pred, scope, branch)
	for i := range d.Facts {
		if FactKey(d.Facts[i].Kind, d.Facts[i].Thing, d.Facts[i].Pred, d.Facts[i].Scope, d.Facts[i].Branch) == key {
			return &d.Facts[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindFact(id string) (*Fact, int) {
	for i := range d.Facts {
		if d.Facts[i].ID == id {
			return &d.Facts[i], i
		}
	}
	return nil, -1
}

func duplicateFactError(existing Fact, kind FactKind, thing, pred, scope string) error {
	return registrationErrorf(
		"duplicate fact: same (kind=%s, thing=%s, pred=%s, scope=%s) already exists as %s; use: novel-logic fact update %s",
		kind, thing, pred, normalizeScope(scope), existing.ID, existing.ID,
	)
}

// ActionKey is the uniqueness key for an action (label is metadata, not part of the key).
func ActionKey(thing, from, to, at, scope, branch string) string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s", thing, from, to, at, normalizeScope(scope), NormalizeBranch(branch))
}

func (d *Data) FindActionByKey(thing, from, to, at, scope, branch string) (*Action, int) {
	key := ActionKey(thing, from, to, at, scope, branch)
	for i := range d.Actions {
		a := d.Actions[i]
		if ActionKey(a.Thing, a.From, a.To, a.At, a.Scope, a.Branch) == key {
			return &d.Actions[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindAction(id string) (*Action, int) {
	for i := range d.Actions {
		if d.Actions[i].ID == id {
			return &d.Actions[i], i
		}
	}
	return nil, -1
}

func duplicateActionError(existing Action, thing, from, to, at, scope string) error {
	return registrationErrorf(
		"duplicate action: same (thing=%s, from=%s, to=%s, at=%s, scope=%s) already exists as %s; use: novel-logic action update %s",
		thing, emptyDash(from), to, at, normalizeScope(scope), existing.ID, existing.ID,
	)
}

func emptyDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// RuleKey is the uniqueness key for a rule within a branch.
func RuleKey(kind RuleKind, thing, pred, from, to, branch string) string {
	b := NormalizeBranch(branch)
	switch kind {
	case RuleForbidState:
		return fmt.Sprintf("%s|%s|%s|%s", kind, thing, pred, b)
	case RuleForbidTransition:
		return fmt.Sprintf("%s|%s|%s|%s", kind, from, to, b)
	default:
		return string(kind) + "|" + b
	}
}

func (d *Data) FindRuleByKey(kind RuleKind, thing, pred, from, to, branch string) (*Rule, int) {
	key := RuleKey(kind, thing, pred, from, to, branch)
	for i := range d.Rules {
		r := d.Rules[i]
		if RuleKey(r.Kind, r.Thing, r.Pred, r.From, r.To, r.Branch) == key {
			return &d.Rules[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindRule(id string) (*Rule, int) {
	for i := range d.Rules {
		if d.Rules[i].ID == id {
			return &d.Rules[i], i
		}
	}
	return nil, -1
}

func duplicateRuleError(existing Rule) error {
	switch existing.Kind {
	case RuleForbidState:
		return registrationErrorf(
			"duplicate rule: same forbid-state (thing=%s, pred=%s) already exists as %s; use: novel-logic rule update %s",
			existing.Thing, existing.Pred, existing.ID, existing.ID,
		)
	case RuleForbidTransition:
		return registrationErrorf(
			"duplicate rule: same forbid-transition (from=%s, to=%s) already exists as %s; use: novel-logic rule update %s",
			existing.From, existing.To, existing.ID, existing.ID,
		)
	default:
		return registrationErrorf("duplicate rule %q already exists; use: novel-logic rule update %s", existing.ID, existing.ID)
	}
}

func duplicateThingError(id string) error {
	return registrationErrorf(
		"thing %q already exists; add scopes with: novel-logic thing scope add %s --scope <scope>; change fields with: novel-logic thing update %s",
		id, id, id,
	)
}

func duplicateTimeError(id string) error {
	return registrationErrorf("duplicate time %q already registered", id)
}

func formatFactKey(kind FactKind, thing, pred, scope string) string {
	return fmt.Sprintf("kind=%s thing=%s pred=%s scope=%s", kind, thing, pred, normalizeScope(scope))
}

func formatActionKey(thing, from, to, at, scope string) string {
	return fmt.Sprintf("thing=%s from=%s to=%s at=%s scope=%s", thing, emptyDash(from), to, at, normalizeScope(scope))
}

func formatRuleKey(r Rule) string {
	switch r.Kind {
	case RuleForbidState:
		return fmt.Sprintf("forbid-state thing=%s pred=%s", r.Thing, r.Pred)
	case RuleForbidTransition:
		return fmt.Sprintf("forbid-transition from=%s to=%s", r.From, r.To)
	default:
		return string(r.Kind)
	}
}

// DuplicateIssues returns validation issues for duplicate entity keys in loaded data.
func DuplicateIssues(d *Data) []string {
	var issues []string
	seenFacts := make(map[string]string)
	for _, f := range d.Facts {
		key := FactKey(f.Kind, f.Thing, f.Pred, f.Scope, f.Branch)
		if prev, ok := seenFacts[key]; ok {
			issues = append(issues, fmt.Sprintf("duplicate fact: %s (%s and %s)", formatFactKey(f.Kind, f.Thing, f.Pred, f.Scope), prev, f.ID))
		} else {
			seenFacts[key] = f.ID
		}
	}
	seenActions := make(map[string]string)
	for _, a := range d.Actions {
		key := ActionKey(a.Thing, a.From, a.To, a.At, a.Scope, a.Branch)
		if prev, ok := seenActions[key]; ok {
			issues = append(issues, fmt.Sprintf("duplicate action: %s (%s and %s)", formatActionKey(a.Thing, a.From, a.To, a.At, a.Scope), prev, a.ID))
		} else {
			seenActions[key] = a.ID
		}
	}
	seenRules := make(map[string]string)
	for _, r := range d.Rules {
		key := RuleKey(r.Kind, r.Thing, r.Pred, r.From, r.To, r.Branch)
		if prev, ok := seenRules[key]; ok {
			issues = append(issues, fmt.Sprintf("duplicate rule: %s (%s and %s)", formatRuleKey(r), prev, r.ID))
		} else {
			seenRules[key] = r.ID
		}
	}
	seenThings := make(map[string]bool)
	for _, t := range d.Things {
		if t.ID == "" {
			continue
		}
		if seenThings[t.ID] {
			issues = append(issues, fmt.Sprintf("duplicate thing id %q", t.ID))
		}
		seenThings[t.ID] = true
	}
	return issues
}
