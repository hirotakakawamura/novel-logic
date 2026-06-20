package project

import "fmt"

// BranchIssues returns validation issues for branch/fork/merge consistency.
func BranchIssues(d *Data) []string {
	var issues []string
	branchIDs := d.BranchIDs()

	for _, f := range d.Facts {
		if err := d.validateBranchRef(f.Branch); err != nil {
			issues = append(issues, fmt.Sprintf("fact %q: %v", f.ID, err))
		}
	}
	for _, a := range d.Actions {
		if err := d.validateBranchRef(a.Branch); err != nil {
			issues = append(issues, fmt.Sprintf("action %q: %v", a.ID, err))
		}
	}
	for _, r := range d.Rules {
		if err := d.validateBranchRef(r.Branch); err != nil {
			issues = append(issues, fmt.Sprintf("rule %q: %v", r.ID, err))
		}
	}
	for _, n := range d.Novels {
		if err := d.validateBranchRef(n.Branch); err != nil {
			issues = append(issues, fmt.Sprintf("novel %q branch %q: %v", n.SceneID, NormalizeBranch(n.Branch), err))
		}
	}

	seenForkSlot := map[string]string{}
	for _, fork := range d.Forks {
		slot := NormalizeBranch(fork.ParentBranch) + "|" + fork.At
		if prev, ok := seenForkSlot[slot]; ok {
			issues = append(issues, fmt.Sprintf(
				"fork.exclusive: fork %q and %q share parent %q at %q",
				prev, fork.ID, fork.ParentBranch, fork.At,
			))
		} else {
			seenForkSlot[slot] = fork.ID
		}
		if err := d.validateBranchRef(fork.ParentBranch); err != nil {
			issues = append(issues, fmt.Sprintf("fork %q: %v", fork.ID, err))
		}
		if d.TimeIndex(fork.At) < 0 {
			issues = append(issues, fmt.Sprintf("fork %q: unknown time %q", fork.ID, fork.At))
		}
		for _, c := range fork.Choices {
			if !branchIDs[c.Branch] {
				issues = append(issues, fmt.Sprintf("fork %q: unknown choice branch %q", fork.ID, c.Branch))
			}
			act, _ := d.FindAction(c.Action)
			if act == nil {
				issues = append(issues, fmt.Sprintf("fork %q: unknown choice action %q", fork.ID, c.Action))
				continue
			}
			if NormalizeBranch(act.Branch) != fork.ParentBranch {
				issues = append(issues, fmt.Sprintf("fork %q: action %q on branch %q, expected %q", fork.ID, c.Action, NormalizeBranch(act.Branch), fork.ParentBranch))
			}
			if act.At != fork.At {
				issues = append(issues, fmt.Sprintf("fork %q: action %q at %q, expected %q", fork.ID, c.Action, act.At, fork.At))
			}
		}
	}

	for _, merge := range d.Merges {
		if err := d.validateBranchRef(merge.IntoBranch); err != nil {
			issues = append(issues, fmt.Sprintf("merge %q: %v", merge.ID, err))
		}
		if d.TimeIndex(merge.At) < 0 {
			issues = append(issues, fmt.Sprintf("merge %q: unknown time %q", merge.ID, merge.At))
		}
		var toPred string
		for _, c := range merge.Choices {
			if err := d.validateBranchRef(c.Branch); err != nil {
				issues = append(issues, fmt.Sprintf("merge %q choice %q: %v", merge.ID, c.Branch, err))
			}
			act, _ := d.FindAction(c.Action)
			if act == nil {
				issues = append(issues, fmt.Sprintf("merge %q: unknown action %q for branch %q", merge.ID, c.Action, c.Branch))
				continue
			}
			if NormalizeBranch(act.Branch) != c.Branch {
				issues = append(issues, fmt.Sprintf("merge %q: action %q on branch %q, expected %q", merge.ID, c.Action, NormalizeBranch(act.Branch), c.Branch))
			}
			if act.At != merge.At {
				issues = append(issues, fmt.Sprintf("merge %q: action %q at %q, expected %q", merge.ID, c.Action, act.At, merge.At))
			}
			if toPred == "" {
				toPred = act.To
			} else if act.To != toPred {
				issues = append(issues, fmt.Sprintf("merge %q: action %q to %q differs from common to %q", merge.ID, c.Action, act.To, toPred))
			}
		}
		for _, a := range d.Actions {
			if !d.branchClosed(NormalizeBranch(a.Branch)) {
				continue
			}
			merge := d.FindMergeForBranch(a.Branch)
			if merge == nil {
				continue
			}
			if d.TimeLE(merge.At, a.At) && a.At != merge.At {
				issues = append(issues, fmt.Sprintf("merge.after_action: action %q on closed branch %q after merge %q", a.ID, a.Branch, merge.ID))
			}
		}
	}

	seenNovels := map[string]string{}
	for _, n := range d.Novels {
		key := NovelKey(n.SceneID, n.Branch)
		if prev, ok := seenNovels[key]; ok {
			issues = append(issues, fmt.Sprintf("duplicate novel: scene %q branch %q (%s and %s)", n.SceneID, NormalizeBranch(n.Branch), prev, key))
		}
		seenNovels[key] = n.SceneID
	}

	for _, msg := range d.BranchIsolatedStateIssues() {
		issues = append(issues, msg)
	}

	return issues
}