package validate

import (
	"fmt"
	"strings"

	"novel-logic/internal/project"
)

type Issue struct {
	Code    string
	Message string
}

func Run(d *project.Data) []Issue {
	return RunForBranch(d, "")
}

func RunForBranch(d *project.Data, branchID string) []Issue {
	var issues []Issue
	for _, msg := range project.DuplicateIssues(d) {
		issues = append(issues, Issue{"duplicate", msg})
	}
	for _, msg := range project.BranchIssues(d) {
		code := "branch.invalid"
		if contains(msg, "merge.after_action") {
			code = "merge.after_action"
		} else if contains(msg, "merge ") && contains(msg, "to ") {
			code = "merge.action_mismatch"
		} else if contains(msg, "fork ") {
			code = "fork.invalid"
		} else if contains(msg, "unknown branch") {
			code = "branch.unknown"
		} else if contains(msg, "body path") {
			code = "novel.branch_path"
		}
		issues = append(issues, Issue{code, msg})
	}
	for _, msg := range project.NovelBodyIssues(d) {
		issues = append(issues, Issue{"novel.missing_body", msg})
	}
	for _, msg := range project.NovelRevisionIssues(d) {
		issues = append(issues, Issue{"novel.revision_drift", msg})
	}
	thingIDs := d.ThingIDs()

	for _, t := range d.Things {
		if len(t.Tags) == 0 {
			issues = append(issues, Issue{"thing.no_tag", fmt.Sprintf("thing %q has no tags", t.ID)})
		}
		if t.ID == "" {
			issues = append(issues, Issue{"thing.empty_id", "thing with empty id"})
		}
	}

	// pred must not collide with thing id
	for _, f := range d.Facts {
		if thingIDs[f.Pred] {
			issues = append(issues, Issue{"pred.thing_collision", fmt.Sprintf("fact %q: pred %q matches thing id", f.ID, f.Pred)})
		}
		if !thingIDs[f.Thing] {
			issues = append(issues, Issue{"ref.thing", fmt.Sprintf("fact %q: unknown thing %q", f.ID, f.Thing)})
		}
		if err := checkScope(d, f.Scope); err != nil {
			issues = append(issues, Issue{"scope.invalid", fmt.Sprintf("fact %q: %v", f.ID, err)})
		}
	}

	sceneWindows := map[string]project.Scene{}
	for _, s := range d.Scenes {
		sceneWindows[s.ID] = s
		if !d.TimeLE(s.TimeStart, s.TimeEnd) {
			issues = append(issues, Issue{"time.scene_window", fmt.Sprintf("scene %q: time_start after time_end", s.ID)})
		}
	}

	actions := d.Actions
	if branchID != "" {
		actions = d.ActiveActions(branchID)
	}
	for _, a := range actions {
		if a.From != "" && thingIDs[a.From] {
			issues = append(issues, Issue{"pred.thing_collision", fmt.Sprintf("action %q: from pred %q matches thing id", a.ID, a.From)})
		}
		if thingIDs[a.To] {
			issues = append(issues, Issue{"pred.thing_collision", fmt.Sprintf("action %q: to pred %q matches thing id", a.ID, a.To)})
		}
		if !thingIDs[a.Thing] {
			issues = append(issues, Issue{"ref.thing", fmt.Sprintf("action %q: unknown thing %q", a.ID, a.Thing)})
		}
		if d.TimeIndex(a.At) < 0 {
			issues = append(issues, Issue{"ref.time", fmt.Sprintf("action %q: unknown time %q", a.ID, a.At)})
		}
		if err := checkScope(d, a.Scope); err != nil {
			issues = append(issues, Issue{"scope.invalid", fmt.Sprintf("action %q: %v", a.ID, err)})
		}
		if inWindow, msg := actionInScopeWindow(d, a, sceneWindows); !inWindow {
			issues = append(issues, Issue{"time.action_window", fmt.Sprintf("action %q: %s", a.ID, msg)})
		}
		if msg := checkActionRules(d, a); msg != "" {
			issues = append(issues, Issue{"rule.violation", fmt.Sprintf("action %q: %s", a.ID, msg)})
		}
	}

	for _, r := range d.Rules {
		if r.Kind == project.RuleForbidState {
			if r.Thing == "" || r.Pred == "" {
				issues = append(issues, Issue{"rule.incomplete", fmt.Sprintf("rule %q: forbid-state needs thing and pred", r.ID)})
			}
		}
		if r.Kind == project.RuleForbidTransition {
			if r.From == "" || r.To == "" {
				issues = append(issues, Issue{"rule.incomplete", fmt.Sprintf("rule %q: forbid-transition needs from and to", r.ID)})
			}
		}
	}

	for _, f := range d.Facts {
		if f.Kind == project.FactState {
			if msg := checkForbidState(d, f.Thing, f.Pred); msg != "" {
				issues = append(issues, Issue{"rule.violation", fmt.Sprintf("fact %q: %s", f.ID, msg)})
			}
		}
	}

	for _, n := range d.Novels {
		start, end, ok := d.SceneWindow(n.SceneID)
		if !ok {
			issues = append(issues, Issue{"novel.unknown_scene", fmt.Sprintf("novel %q: unknown scene", n.SceneID)})
			continue
		}
		if n.TimeStart != start || n.TimeEnd != end {
			issues = append(issues, Issue{"novel.time_mismatch",
				fmt.Sprintf("novel %q: time [%s..%s] differs from scene [%s..%s]",
					n.SceneID, n.TimeStart, n.TimeEnd, start, end)})
		}
	}

	return issues
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

// Hints returns non-fatal alignment suggestions (Phase A vs Phase B).
func Hints(d *project.Data) []Issue {
	var hints []Issue
	for _, a := range d.Actions {
		scope := a.Scope
		if scope == "" {
			scope = "plot"
		}
		if scope != "plot" {
			continue
		}
		scenes := d.ScenesContainingTime(a.At)
		if len(scenes) == 0 {
			continue
		}
		hints = append(hints, Issue{"action.plot_scene_hint",
			fmt.Sprintf("action %q: scope=plot but at %q is in scene(s) %v; use --scope novel:<scene> for Phase B alignment",
				a.ID, a.At, scenes)})
	}
	for _, msg := range project.NovelRevisionHints(d) {
		hints = append(hints, Issue{"novel.revision_hint", msg})
	}
	return hints
}

func checkScope(d *project.Data, scope string) error {
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

func actionInScopeWindow(d *project.Data, a project.Action, scenes map[string]project.Scene) (bool, string) {
	scope := a.Scope
	if scope == "" {
		scope = "plot"
	}
	if scope == "plot" {
		return true, ""
	}
	sid := strings.TrimPrefix(scope, "novel:")
	s, ok := scenes[sid]
	if !ok {
		return false, "scope scene not found"
	}
	if !d.TimeLE(s.TimeStart, a.At) || !d.TimeLE(a.At, s.TimeEnd) {
		return false, fmt.Sprintf("time %q outside scene window [%q,%q]", a.At, s.TimeStart, s.TimeEnd)
	}
	return true, ""
}

func checkForbidState(d *project.Data, thing, pred string) string {
	for _, r := range d.Rules {
		if r.Kind == project.RuleForbidState && r.Thing == thing && r.Pred == pred {
			return fmt.Sprintf("forbidden state %q for thing %q", pred, thing)
		}
	}
	return ""
}

func checkActionRules(d *project.Data, a project.Action) string {
	if msg := checkForbidState(d, a.Thing, a.To); msg != "" {
		return msg
	}
	if a.From != "" {
		for _, r := range d.Rules {
			if r.Kind == project.RuleForbidTransition && r.From == a.From && r.To == a.To {
				return fmt.Sprintf("forbidden transition %q -> %q", a.From, a.To)
			}
		}
	}
	return ""
}