package project

import (
	"fmt"
	"sort"
)

// EffectiveFactsOnBranch returns facts registered on branches in the lineage of branchID.
func (d *Data) EffectiveFactsOnBranch(branchID string) []Fact {
	lineage := map[string]bool{}
	for _, b := range d.BranchLineage(branchID) {
		lineage[b] = true
	}
	var out []Fact
	for _, f := range d.Facts {
		if lineage[NormalizeBranch(f.Branch)] {
			out = append(out, f)
		}
	}
	return out
}

// EffectiveRulesOnBranch returns rules registered on branches in the lineage of branchID.
func (d *Data) EffectiveRulesOnBranch(branchID string) []Rule {
	lineage := map[string]bool{}
	for _, b := range d.BranchLineage(branchID) {
		lineage[b] = true
	}
	var out []Rule
	for _, r := range d.Rules {
		if lineage[NormalizeBranch(r.Branch)] {
			out = append(out, r)
		}
	}
	return out
}

// predsReachableOnBranch reports whether pred is active for thing strictly before at on branchID.
func (d *Data) predsReachableOnBranch(branchID, thing, pred, at, excludeActionID string) bool {
	if pred == "" {
		return true
	}
	preds := map[string]bool{}
	for _, f := range d.EffectiveFactsOnBranch(branchID) {
		if f.Thing == thing {
			preds[f.Pred] = true
		}
	}
	atIdx := d.TimeIndex(at)
	if atIdx < 0 {
		return false
	}
	var acts []Action
	for _, a := range d.ActiveActions(branchID) {
		if a.Thing != thing || a.ID == excludeActionID {
			continue
		}
		idx := d.TimeIndex(a.At)
		if idx < 0 || idx >= atIdx {
			continue
		}
		acts = append(acts, a)
	}
	sort.Slice(acts, func(i, j int) bool {
		ti, tj := d.TimeIndex(acts[i].At), d.TimeIndex(acts[j].At)
		if ti != tj {
			return ti < tj
		}
		return acts[i].ID < acts[j].ID
	})
	for _, a := range acts {
		if a.From != "" {
			delete(preds, a.From)
		}
		preds[a.To] = true
	}
	return preds[pred]
}

// predIntroducedBranches lists branch ids that declare thing/pred via state fact or action to.
func (d *Data) predIntroducedBranches(thing, pred string) []string {
	seen := map[string]bool{}
	var out []string
	for _, a := range d.Actions {
		if a.Thing == thing && a.To == pred {
			b := NormalizeBranch(a.Branch)
			if !seen[b] {
				seen[b] = true
				out = append(out, b)
			}
		}
	}
	for _, f := range d.Facts {
		if f.Kind == FactState && f.Thing == thing && f.Pred == pred {
			b := NormalizeBranch(f.Branch)
			if !seen[b] {
				seen[b] = true
				out = append(out, b)
			}
		}
	}
	sort.Strings(out)
	return out
}

// BranchIsolatedStateIssues detects actions referencing preds exclusive to other branches.
func (d *Data) BranchIsolatedStateIssues() []string {
	var issues []string
	for _, branchID := range d.AllBranchIDs() {
		lineage := map[string]bool{}
		for _, b := range d.BranchLineage(branchID) {
			lineage[b] = true
		}
		for _, a := range d.ActiveActions(branchID) {
			if a.From == "" {
				continue
			}
			if d.predsReachableOnBranch(branchID, a.Thing, a.From, a.At, a.ID) {
				continue
			}
			intro := d.predIntroducedBranches(a.Thing, a.From)
			if len(intro) == 0 {
				continue
			}
			exclusive := true
			for _, b := range intro {
				if lineage[b] {
					exclusive = false
					break
				}
			}
			if exclusive {
				issues = append(issues, fmt.Sprintf(
					"branch.isolated_state: action %q on branch %q references from %q for %q introduced only on other branch(es) %v",
					a.ID, branchID, a.From, a.Thing, intro,
				))
			}
		}
	}
	return issues
}