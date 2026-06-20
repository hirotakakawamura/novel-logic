package project

import (
	"fmt"
	"sort"
)

func (d *Data) existingBranchIDs() map[string]bool {
	m := make(map[string]bool, len(d.Branches))
	for _, b := range d.Branches {
		m[b.ID] = true
	}
	return m
}

func (d *Data) existingForkIDs() map[string]bool {
	m := make(map[string]bool, len(d.Forks))
	for _, f := range d.Forks {
		m[f.ID] = true
	}
	return m
}

func (d *Data) existingMergeIDs() map[string]bool {
	m := make(map[string]bool, len(d.Merges))
	for _, mrg := range d.Merges {
		m[mrg.ID] = true
	}
	return m
}

// AddBranch registers a story branch. main is created automatically on load.
func (d *Data) AddBranch(id, label, parent, viaFork, viaAction string) error {
	if id == "" {
		return fmt.Errorf("branch id is required")
	}
	if id == MainBranch && parent != "" {
		return fmt.Errorf("branch %q cannot have a parent", MainBranch)
	}
	if _, idx := d.FindBranchDef(id); idx >= 0 {
		return fmt.Errorf("branch %q already exists", id)
	}
	parent = NormalizeBranch(parent)
	if parent != "" {
		if err := d.validateBranchRef(parent); err != nil {
			return err
		}
	}
	if viaFork != "" {
		fork, _ := d.FindFork(viaFork)
		if fork == nil {
			return fmt.Errorf("unknown fork %q", viaFork)
		}
	}
	if viaAction != "" {
		act, _ := d.FindAction(viaAction)
		if act == nil {
			return fmt.Errorf("unknown action %q", viaAction)
		}
	}
	d.Branches = append(d.Branches, Branch{
		ID:        id,
		Label:     label,
		Parent:    parent,
		ViaFork:   viaFork,
		ViaAction: viaAction,
	})
	return nil
}

func (d *Data) RemoveBranch(id string) error {
	id = NormalizeBranch(id)
	if id == MainBranch {
		return fmt.Errorf("cannot remove branch %q", MainBranch)
	}
	idx := -1
	for i := range d.Branches {
		if d.Branches[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("branch %q not found", id)
	}
	for _, b := range d.Branches {
		if NormalizeBranch(b.Parent) == id {
			return fmt.Errorf("branch %q has child %q; remove children first", id, b.ID)
		}
	}
	d.Branches = append(d.Branches[:idx], d.Branches[idx+1:]...)
	return nil
}

// AddFork registers a fork point on a parent branch.
func (d *Data) AddFork(id, parentBranch, at, scope string) error {
	if id == "" || at == "" {
		return fmt.Errorf("fork id and at are required")
	}
	parentBranch = NormalizeBranch(parentBranch)
	if err := d.validateBranchRef(parentBranch); err != nil {
		return err
	}
	if d.TimeIndex(at) < 0 {
		return fmt.Errorf("unknown time %q", at)
	}
	if scope == "" {
		scope = "plot"
	}
	if err := validateScopeRef(d, scope); err != nil {
		return err
	}
	if _, idx := d.FindFork(id); idx >= 0 {
		return fmt.Errorf("fork %q already exists", id)
	}
	for _, f := range d.Forks {
		if NormalizeBranch(f.ParentBranch) == parentBranch && f.At == at {
			return fmt.Errorf("fork already exists at %q on branch %q (%s)", at, parentBranch, f.ID)
		}
	}
	d.Forks = append(d.Forks, Fork{
		ID:           id,
		ParentBranch: parentBranch,
		At:           at,
		Scope:        scope,
	})
	return nil
}

func (d *Data) AddForkChoice(forkID, actionID, branchID string) error {
	fork, idx := d.FindFork(forkID)
	if fork == nil {
		return fmt.Errorf("fork %q not found", forkID)
	}
	if actionID == "" || branchID == "" {
		return fmt.Errorf("action and branch are required")
	}
	act, _ := d.FindAction(actionID)
	if act == nil {
		return fmt.Errorf("unknown action %q", actionID)
	}
	if NormalizeBranch(act.Branch) != fork.ParentBranch {
		return fmt.Errorf("action %q is on branch %q, expected %q", actionID, NormalizeBranch(act.Branch), fork.ParentBranch)
	}
	if act.At != fork.At {
		return fmt.Errorf("action %q at %q does not match fork at %q", actionID, act.At, fork.At)
	}
	if _, bidx := d.FindBranchDef(branchID); bidx >= 0 {
		return fmt.Errorf("branch %q already exists", branchID)
	}
	for _, c := range fork.Choices {
		if c.Action == actionID {
			return fmt.Errorf("fork %q already has choice for action %q", forkID, actionID)
		}
		if c.Branch == branchID {
			return fmt.Errorf("fork %q already has choice for branch %q", forkID, branchID)
		}
	}
	fork.Choices = append(fork.Choices, ForkChoice{Action: actionID, Branch: branchID})
	d.Forks[idx] = *fork
	d.Branches = append(d.Branches, Branch{
		ID:        branchID,
		Parent:    fork.ParentBranch,
		ViaFork:   forkID,
		ViaAction: actionID,
	})
	return nil
}

// AddMerge registers a merge point for child branches into intoBranch.
func (d *Data) AddMerge(id, at, scope, intoBranch string, choices []MergeChoice) error {
	if id == "" || at == "" {
		return fmt.Errorf("merge id and at are required")
	}
	intoBranch = NormalizeBranch(intoBranch)
	if err := d.validateBranchRef(intoBranch); err != nil {
		return err
	}
	if d.TimeIndex(at) < 0 {
		return fmt.Errorf("unknown time %q", at)
	}
	if scope == "" {
		scope = "plot"
	}
	if err := validateScopeRef(d, scope); err != nil {
		return err
	}
	if len(choices) == 0 {
		return fmt.Errorf("at least one merge choice (branch + action) is required")
	}
	if _, idx := d.FindMerge(id); idx >= 0 {
		return fmt.Errorf("merge %q already exists", id)
	}
	seen := make(map[string]bool)
	for _, c := range choices {
		if c.Branch == "" || c.Action == "" {
			return fmt.Errorf("each merge choice needs branch and action")
		}
		if seen[c.Branch] {
			return fmt.Errorf("duplicate merge choice for branch %q", c.Branch)
		}
		seen[c.Branch] = true
		if err := d.validateBranchRef(c.Branch); err != nil {
			return err
		}
		act, _ := d.FindAction(c.Action)
		if act == nil {
			return fmt.Errorf("unknown action %q", c.Action)
		}
		if NormalizeBranch(act.Branch) != c.Branch {
			return fmt.Errorf("action %q is on branch %q, expected %q", c.Action, NormalizeBranch(act.Branch), c.Branch)
		}
		if act.At != at {
			return fmt.Errorf("action %q at %q does not match merge at %q", c.Action, act.At, at)
		}
	}
	d.Merges = append(d.Merges, Merge{
		ID:         id,
		At:         at,
		Scope:      scope,
		IntoBranch: intoBranch,
		Choices:    append([]MergeChoice{}, choices...),
	})
	return nil
}

// BranchLineage returns [main, ..., branchID] from root to the given branch.
func (d *Data) BranchLineage(branchID string) []string {
	branchID = NormalizeBranch(branchID)
	var chain []string
	current := branchID
	for {
		chain = append([]string{current}, chain...)
		b, _ := d.FindBranchDef(current)
		if b == nil || NormalizeBranch(b.Parent) == "" || NormalizeBranch(b.Parent) == current {
			break
		}
		current = NormalizeBranch(b.Parent)
	}
	return chain
}

func (d *Data) forkAtForChild(parent, child string) string {
	b, _ := d.FindBranchDef(child)
	if b == nil || b.ViaFork == "" {
		return ""
	}
	fork, _ := d.FindFork(b.ViaFork)
	if fork == nil || NormalizeBranch(fork.ParentBranch) != NormalizeBranch(parent) {
		return ""
	}
	return fork.At
}

// ActiveActions returns actions effective on branchID following fork/merge boundaries.
func (d *Data) ActiveActions(branchID string) []Action {
	branchID = NormalizeBranch(branchID)
	lineage := d.BranchLineage(branchID)
	seen := make(map[string]bool)
	var out []Action

	appendAction := func(a Action) {
		if seen[a.ID] {
			return
		}
		seen[a.ID] = true
		out = append(out, a)
	}

	for i, bid := range lineage {
		forkAt := ""
		if i > 0 {
			forkAt = d.forkAtForChild(lineage[i-1], bid)
		}
		var mergeAt string
		if merge := d.FindMergeForBranch(bid); merge != nil && bid == branchID {
			mergeAt = merge.At
		}
		childFork := ""
		childEntryAction := ""
		if i < len(lineage)-1 {
			child := lineage[i+1]
			childFork = d.forkAtForChild(bid, child)
			if b, _ := d.FindBranchDef(child); b != nil {
				childEntryAction = b.ViaAction
			}
		}

		for _, a := range d.Actions {
			if NormalizeBranch(a.Branch) != bid {
				continue
			}
			if forkAt != "" && !d.TimeLE(forkAt, a.At) {
				continue
			}
			if childFork != "" && bid != branchID {
				if d.TimeLE(childFork, a.At) && !(a.At == childFork && a.ID == childEntryAction) {
					continue
				}
			}
			if mergeAt != "" && d.TimeLE(mergeAt, a.At) {
				continue
			}
			appendAction(a)
		}
	}

	if merge := d.FindMergeForBranch(branchID); merge != nil {
		into := NormalizeBranch(merge.IntoBranch)
		for _, a := range d.Actions {
			if NormalizeBranch(a.Branch) != into {
				continue
			}
			if !d.TimeLE(merge.At, a.At) {
				continue
			}
			appendAction(a)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		ti, tj := d.TimeIndex(out[i].At), d.TimeIndex(out[j].At)
		if ti != tj {
			return ti < tj
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// ActiveFacts returns facts for branchID (branch field must match).
func (d *Data) ActiveFacts(branchID string) []Fact {
	branchID = NormalizeBranch(branchID)
	var out []Fact
	for _, f := range d.Facts {
		if NormalizeBranch(f.Branch) == branchID {
			out = append(out, f)
		}
	}
	return out
}

// ActiveRules returns rules for branchID (branch field must match).
func (d *Data) ActiveRules(branchID string) []Rule {
	branchID = NormalizeBranch(branchID)
	var out []Rule
	for _, r := range d.Rules {
		if NormalizeBranch(r.Branch) == branchID {
			out = append(out, r)
		}
	}
	return out
}

// AllBranchIDs returns every branch id including implicit main.
func (d *Data) AllBranchIDs() []string {
	ids := make([]string, 0, len(d.Branches))
	for _, b := range d.Branches {
		ids = append(ids, b.ID)
	}
	if len(ids) == 0 {
		return []string{MainBranch}
	}
	sort.Strings(ids)
	return ids
}