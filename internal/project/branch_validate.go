package project

import "fmt"

// BranchIssue is a structured branch/fork/merge validation finding.
type BranchIssue struct {
	Code    string
	Message string
}

// BranchIssues returns validation issues for branch/fork/merge consistency.
func BranchIssues(d *Data) []BranchIssue {
	var issues []BranchIssue
	branchIDs := d.BranchIDs()

	appendRef := func(subject string, err error) {
		if err == nil {
			return
		}
		issues = append(issues, BranchIssue{
			Code:    branchRefIssueCode(err),
			Message: fmt.Sprintf("%s: %v", subject, err),
		})
	}

	for _, f := range d.Facts {
		appendRef(fmt.Sprintf("fact %q", f.ID), d.validateBranchRef(f.Branch))
	}
	for _, a := range d.Actions {
		appendRef(fmt.Sprintf("action %q", a.ID), d.validateBranchRef(a.Branch))
	}
	for _, r := range d.Rules {
		appendRef(fmt.Sprintf("rule %q", r.ID), d.validateBranchRef(r.Branch))
	}
	for _, n := range d.Novels {
		appendRef(fmt.Sprintf("novel %q branch %q", n.SceneID, NormalizeBranch(n.Branch)), d.validateBranchRef(n.Branch))
	}

	seenForkSlot := map[string]string{}
	for _, fork := range d.Forks {
		slot := NormalizeBranch(fork.ParentBranch) + "|" + fork.At
		if prev, ok := seenForkSlot[slot]; ok {
			issues = append(issues, BranchIssue{
				Code: "fork.exclusive",
				Message: fmt.Sprintf(
					"fork %q and %q share parent %q at %q",
					prev, fork.ID, fork.ParentBranch, fork.At,
				),
			})
		} else {
			seenForkSlot[slot] = fork.ID
		}
		appendRef(fmt.Sprintf("fork %q", fork.ID), d.validateBranchRef(fork.ParentBranch))
		if d.TimeIndex(fork.At) < 0 {
			issues = append(issues, BranchIssue{
				Code:    "fork.invalid",
				Message: fmt.Sprintf("fork %q: unknown time %q", fork.ID, fork.At),
			})
		}
		for _, c := range fork.Choices {
			if !branchIDs[c.Branch] {
				issues = append(issues, BranchIssue{
					Code:    "fork.invalid",
					Message: fmt.Sprintf("fork %q: unknown choice branch %q", fork.ID, c.Branch),
				})
			}
			act, _ := d.FindAction(c.Action)
			if act == nil {
				issues = append(issues, BranchIssue{
					Code:    "fork.invalid",
					Message: fmt.Sprintf("fork %q: unknown choice action %q", fork.ID, c.Action),
				})
				continue
			}
			if NormalizeBranch(act.Branch) != fork.ParentBranch {
				issues = append(issues, BranchIssue{
					Code:    "fork.invalid",
					Message: fmt.Sprintf("fork %q: action %q on branch %q, expected %q", fork.ID, c.Action, NormalizeBranch(act.Branch), fork.ParentBranch),
				})
			}
			if act.At != fork.At {
				issues = append(issues, BranchIssue{
					Code:    "fork.invalid",
					Message: fmt.Sprintf("fork %q: action %q at %q, expected %q", fork.ID, c.Action, act.At, fork.At),
				})
			}
		}
	}

	for _, merge := range d.Merges {
		appendRef(fmt.Sprintf("merge %q", merge.ID), d.validateBranchRef(merge.IntoBranch))
		if d.TimeIndex(merge.At) < 0 {
			issues = append(issues, BranchIssue{
				Code:    "merge.invalid",
				Message: fmt.Sprintf("merge %q: unknown time %q", merge.ID, merge.At),
			})
		}
		var toPred string
		for _, c := range merge.Choices {
			appendRef(fmt.Sprintf("merge %q choice %q", merge.ID, c.Branch), d.validateBranchRef(c.Branch))
			act, _ := d.FindAction(c.Action)
			if act == nil {
				issues = append(issues, BranchIssue{
					Code:    "merge.invalid",
					Message: fmt.Sprintf("merge %q: unknown action %q for branch %q", merge.ID, c.Action, c.Branch),
				})
				continue
			}
			if NormalizeBranch(act.Branch) != c.Branch {
				issues = append(issues, BranchIssue{
					Code:    "merge.invalid",
					Message: fmt.Sprintf("merge %q: action %q on branch %q, expected %q", merge.ID, c.Action, NormalizeBranch(act.Branch), c.Branch),
				})
			}
			if act.At != merge.At {
				issues = append(issues, BranchIssue{
					Code:    "merge.invalid",
					Message: fmt.Sprintf("merge %q: action %q at %q, expected %q", merge.ID, c.Action, act.At, merge.At),
				})
			}
			if toPred == "" {
				toPred = act.To
			} else if act.To != toPred {
				issues = append(issues, BranchIssue{
					Code:    "merge.action_mismatch",
					Message: fmt.Sprintf("merge %q: action %q to %q differs from common to %q", merge.ID, c.Action, act.To, toPred),
				})
			}
		}
		for _, a := range d.Actions {
			if !d.branchClosed(NormalizeBranch(a.Branch)) {
				continue
			}
			m := d.FindMergeForBranch(a.Branch)
			if m == nil {
				continue
			}
			if d.TimeLE(m.At, a.At) && a.At != m.At {
				issues = append(issues, BranchIssue{
					Code:    "merge.after_action",
					Message: fmt.Sprintf("action %q on closed branch %q after merge %q", a.ID, a.Branch, m.ID),
				})
			}
		}
	}

	seenNovels := map[string]bool{}
	for _, n := range d.Novels {
		key := NovelKey(n.SceneID, n.Branch)
		if seenNovels[key] {
			issues = append(issues, BranchIssue{
				Code:    "novel.duplicate",
				Message: fmt.Sprintf("duplicate novel: scene %q on branch %q registered twice", n.SceneID, NormalizeBranch(n.Branch)),
			})
			continue
		}
		seenNovels[key] = true
	}

	issues = append(issues, d.BranchIsolatedStateIssues()...)
	return issues
}

func branchRefIssueCode(err error) string {
	if err == nil {
		return "branch.invalid"
	}
	msg := err.Error()
	if len(msg) >= 14 && msg[:14] == "unknown branch" {
		return "branch.unknown"
	}
	return "branch.invalid"
}