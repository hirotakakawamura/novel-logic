package project

import (
	"fmt"
	"time"
)

func latestRevisionEntry(n *NovelMeta) *NovelRevision {
	if n == nil || len(n.Revisions) == 0 {
		return nil
	}
	return &n.Revisions[len(n.Revisions)-1]
}

func appendRevision(n *NovelMeta, entry NovelRevision) {
	if entry.RecordedAt.IsZero() {
		entry.RecordedAt = time.Now().UTC()
	}
	if latest := latestRevisionEntry(n); latest != nil && revisionsEqual(latest.Revision, entry.Revision) && latest.Note == entry.Note {
		n.Revision = entry.Revision
		return
	}
	n.Revisions = append(n.Revisions, entry)
	n.Revision = entry.Revision
}

// PinNovelRevision records the git commit for a novel body file.
func (d *Data) PinNovelRevision(sceneID, branch, revision, note string, allowDirty bool) (NovelRevision, error) {
	if sceneID == "" {
		return NovelRevision{}, fmt.Errorf("scene id is required")
	}
	branch = NormalizeBranch(branch)
	n, idx := d.FindNovel(sceneID, branch)
	if n == nil {
		return NovelRevision{}, fmt.Errorf("novel for scene %q on branch %q not found", sceneID, branch)
	}
	var state GitFileState
	var err error
	if revision != "" {
		state, err = ResolveGitRevisionAt(d.Root, n.BodyPath, revision)
	} else {
		state, err = ResolveGitFileState(d.Root, n.BodyPath)
	}
	if err != nil {
		return NovelRevision{}, err
	}
	if state.Dirty && !allowDirty {
		return NovelRevision{}, fmt.Errorf("working tree has uncommitted changes for %q; commit first or use --allow-dirty", n.BodyPath)
	}
	entry := NovelRevision{
		Revision:   state.Revision,
		Short:      state.Short,
		Branch:     state.Branch,
		RecordedAt: time.Now().UTC(),
		Note:       note,
		Dirty:      state.Dirty,
	}
	appendRevision(n, entry)
	d.Novels[idx] = *n
	return entry, nil
}

// NovelRevisionHints returns non-fatal suggestions for unpinned novels.
func NovelRevisionHints(d *Data) []string {
	if !isGitRepo(d.Root) {
		return nil
	}
	var hints []string
	for _, n := range d.Novels {
		if n.Revision == "" {
			hints = append(hints, fmt.Sprintf("novel %q branch %q: no git revision pinned (run: novel-logic novel revision pin %s --branch %s)", n.SceneID, NormalizeBranch(n.Branch), n.SceneID, NormalizeBranch(n.Branch)))
		}
	}
	return hints
}

// NovelRevisionIssues reports stale pins and uncommitted body edits (fatal for CI).
func NovelRevisionIssues(d *Data) []string {
	if !isGitRepo(d.Root) {
		return nil
	}
	var issues []string
	head, headErr := currentHEAD(d.Root)
	for _, n := range d.Novels {
		if n.Revision == "" {
			continue
		}
		if n.BodyPath == "" {
			continue
		}
		if !fileTrackedInGit(d.Root, n.BodyPath) {
			issues = append(issues, fmt.Sprintf("novel %q: body %s is not tracked by git", n.SceneID, n.BodyPath))
			continue
		}
		if headErr == nil && head != "" && !revisionsEqual(n.Revision, head) {
			latest, err := ResolveGitFileState(d.Root, n.BodyPath)
			if err == nil && !revisionsEqual(latest.Revision, n.Revision) {
				short := n.Revision
				if entry := latestRevisionStatic(n); entry != nil && entry.Short != "" {
					short = entry.Short
				}
				issues = append(issues, fmt.Sprintf(
					"novel %q branch %q: pinned revision %s differs from latest git commit %s for %s (run: novel-logic novel revision pin %s --branch %s)",
					n.SceneID, NormalizeBranch(n.Branch), short, latest.Short, n.BodyPath, n.SceneID, NormalizeBranch(n.Branch),
				))
			}
		}
		if workingTreeDirty(d.Root, n.BodyPath) {
			issues = append(issues, fmt.Sprintf("novel %q: working tree has uncommitted changes in %s", n.SceneID, n.BodyPath))
		}
	}
	return issues
}

func latestRevisionStatic(n NovelMeta) *NovelRevision {
	if len(n.Revisions) == 0 {
		return nil
	}
	return &n.Revisions[len(n.Revisions)-1]
}