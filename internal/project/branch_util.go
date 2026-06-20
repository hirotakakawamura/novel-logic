package project

import "fmt"

const MainBranch = "main"

func NormalizeBranch(branch string) string {
	if branch == "" {
		return MainBranch
	}
	return branch
}

func (d *Data) BranchIDs() map[string]bool {
	m := make(map[string]bool, len(d.Branches))
	for _, b := range d.Branches {
		m[b.ID] = true
	}
	return m
}

func (d *Data) FindBranchDef(id string) (*Branch, int) {
	id = NormalizeBranch(id)
	for i := range d.Branches {
		if d.Branches[i].ID == id {
			return &d.Branches[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindFork(id string) (*Fork, int) {
	for i := range d.Forks {
		if d.Forks[i].ID == id {
			return &d.Forks[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindMerge(id string) (*Merge, int) {
	for i := range d.Merges {
		if d.Merges[i].ID == id {
			return &d.Merges[i], i
		}
	}
	return nil, -1
}

func (d *Data) FindMergeForBranch(branchID string) *Merge {
	branchID = NormalizeBranch(branchID)
	for i := range d.Merges {
		m := &d.Merges[i]
		for _, c := range m.Choices {
			if c.Branch == branchID {
				return m
			}
		}
	}
	return nil
}

func (d *Data) validateBranchRef(branch string) error {
	branch = NormalizeBranch(branch)
	if branch == MainBranch {
		if _, idx := d.FindBranchDef(MainBranch); idx < 0 {
			return fmt.Errorf("branch %q not found (run: novel-logic branch add main)", MainBranch)
		}
		return nil
	}
	if _, idx := d.FindBranchDef(branch); idx < 0 {
		return fmt.Errorf("unknown branch %q", branch)
	}
	return nil
}

func (d *Data) branchClosed(branchID string) bool {
	return d.FindMergeForBranch(branchID) != nil
}