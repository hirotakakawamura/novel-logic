package validate

import (
	"os"
	"path/filepath"
	"testing"

	"novel-logic/internal/project"
	"novel-logic/internal/testfixture"
)

func TestRunIssueCodes(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*project.Data)
		wantCode string
	}{
		{
			name: "duplicate_action",
			mutate: func(d *project.Data) {
				d.Actions = append(d.Actions, d.Actions[0])
				d.Actions[1].ID = "act_dup"
			},
			wantCode: "duplicate",
		},
		{
			name: "duplicate_rule",
			mutate: func(d *project.Data) {
				d.Rules = append(d.Rules, project.Rule{
					ID: "rule_dup", Kind: project.RuleForbidTransition, From: "x", To: "y", Branch: project.MainBranch,
				})
				d.Rules = append(d.Rules, project.Rule{
					ID: "rule_dup2", Kind: project.RuleForbidTransition, From: "x", To: "y", Branch: project.MainBranch,
				})
			},
			wantCode: "duplicate",
		},
		{
			name: "thing_no_tag",
			mutate: func(d *project.Data) {
				d.Things[0].Tags = nil
			},
			wantCode: "thing.no_tag",
		},
		{
			name: "pred_thing_collision_action",
			mutate: func(d *project.Data) {
				d.Actions[0].To = "ally"
			},
			wantCode: "pred.thing_collision",
		},
		{
			name: "time_scene_window",
			mutate: func(d *project.Data) {
				d.Scenes[0].TimeStart = "t3"
				d.Scenes[0].TimeEnd = "t1"
			},
			wantCode: "time.scene_window",
		},
		{
			name: "ref_thing",
			mutate: func(d *project.Data) {
				d.Facts[0].Thing = "ghost"
			},
			wantCode: "ref.thing",
		},
		{
			name: "ref_time",
			mutate: func(d *project.Data) {
				d.Actions[0].At = "t99"
			},
			wantCode: "ref.time",
		},
		{
			name: "scope_invalid",
			mutate: func(d *project.Data) {
				d.Facts[0].Scope = "novel:missing_scene"
			},
			wantCode: "scope.invalid",
		},
		{
			name: "rule_violation_state",
			mutate: func(d *project.Data) {
				d.Rules = append(d.Rules, project.Rule{
					ID: "rule1", Kind: project.RuleForbidState, Thing: "hero", Pred: "evil", Branch: project.MainBranch,
				})
				d.Facts = append(d.Facts, project.Fact{
					ID: "fact_evil", Kind: project.FactState, Thing: "hero", Pred: "evil", Scope: "plot", Branch: project.MainBranch,
				})
			},
			wantCode: "rule.violation",
		},
		{
			name: "rule_incomplete",
			mutate: func(d *project.Data) {
				d.Rules = append(d.Rules, project.Rule{
					ID: "rule_bad", Kind: project.RuleForbidState, Thing: "", Pred: "", Branch: project.MainBranch,
				})
			},
			wantCode: "rule.incomplete",
		},
		{
			name: "fork_exclusive",
			mutate: func(d *project.Data) {
				d.Forks = append(d.Forks,
					project.Fork{ID: "fork_a", ParentBranch: project.MainBranch, At: "t2", Scope: "plot"},
					project.Fork{ID: "fork_b", ParentBranch: project.MainBranch, At: "t2", Scope: "plot"},
				)
			},
			wantCode: "fork.exclusive",
		},
		{
			name: "fork_invalid_unknown_time",
			mutate: func(d *project.Data) {
				d.Forks = append(d.Forks, project.Fork{ID: "fork_bad", ParentBranch: project.MainBranch, At: "t99", Scope: "plot"})
			},
			wantCode: "fork.invalid",
		},
		{
			name: "merge_invalid_unknown_action",
			mutate: func(d *project.Data) {
				d.Merges = append(d.Merges, project.Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: project.MainBranch,
					Choices: []project.MergeChoice{{Branch: project.MainBranch, Action: "no_such_act"}},
				})
			},
			wantCode: "merge.invalid",
		},
		{
			name: "branch_unknown_on_fact",
			mutate: func(d *project.Data) {
				d.Facts[0].Branch = "ghost_branch"
			},
			wantCode: "branch.unknown",
		},
		{
			name: "time_registry_mismatch",
			mutate: func(d *project.Data) {
				d.Meta.TimeOrder = append(d.Meta.TimeOrder, "ghost_time")
			},
			wantCode: "time.registry_mismatch",
		},
		{
			name: "novel_unknown_scene",
			mutate: func(d *project.Data) {
				d.Novels = append(d.Novels, project.NovelMeta{
					SceneID: "ghost_scene", Branch: project.MainBranch,
					TimeStart: "t1", TimeEnd: "t2",
					BodyPath: project.DefaultNovelBodyPath("ghost_scene", project.MainBranch),
				})
			},
			wantCode: "novel.unknown_scene",
		},
		{
			name: "novel_missing_body",
			mutate: func(d *project.Data) {
				body := project.DefaultNovelBodyPath("scene1", project.MainBranch)
				d.Novels = append(d.Novels, project.NovelMeta{
					SceneID: "scene1", Branch: project.MainBranch,
					TimeStart: "t1", TimeEnd: "t2", BodyPath: body,
				})
			},
			wantCode: "novel.missing_body",
		},
		{
			name: "novel_time_mismatch",
			mutate: func(d *project.Data) {
				body := filepath.Join(d.Root, project.DefaultNovelBodyPath("scene1", project.MainBranch))
				_ = os.MkdirAll(filepath.Dir(body), 0o755)
				_ = os.WriteFile(body, []byte("x"), 0o644)
				d.Novels = append(d.Novels, project.NovelMeta{
					SceneID: "scene1", Branch: project.MainBranch,
					TimeStart: "t9", TimeEnd: "t9", BodyPath: project.DefaultNovelBodyPath("scene1", project.MainBranch),
				})
			},
			wantCode: "novel.time_mismatch",
		},
		{
			name: "time_action_window",
			mutate: func(d *project.Data) {
				d.Actions[0].Scope = "novel:scene1"
				d.Actions[0].At = "t4"
			},
			wantCode: "time.action_window",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := testfixture.LoadMinimal(t)
			tt.mutate(d)
			issues := Run(d)
			if !hasIssueCode(issues, tt.wantCode) {
				t.Fatalf("want code %q, got %v", tt.wantCode, issues)
			}
		})
	}
}

func TestBranchIssuesCodes(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*project.Data)
		wantCode string
	}{
		{
			name: "novel_duplicate",
			mutate: func(d *project.Data) {
				meta := project.NovelMeta{
					SceneID: "scene1", Branch: project.MainBranch,
					TimeStart: "t1", TimeEnd: "t2",
					BodyPath: project.DefaultNovelBodyPath("scene1", project.MainBranch),
				}
				d.Novels = append(d.Novels, meta, meta)
			},
			wantCode: "novel.duplicate",
		},
		{
			name: "merge_after_action",
			mutate: func(d *project.Data) {
				_ = d.AddBranch("branch_a", "", project.MainBranch, "", "")
				actMain, _ := d.AddAction("hero", "mid", "merged", "t3", "plot", "", project.MainBranch)
				actA, _ := d.AddAction("hero", "mid", "merged", "t3", "plot", "", "branch_a")
				_ = d.AddMerge("merge1", "t3", "plot", project.MainBranch, []project.MergeChoice{
					{Branch: project.MainBranch, Action: actMain.ID},
					{Branch: "branch_a", Action: actA.ID},
				})
				d.Actions = append(d.Actions, project.Action{
					ID: "act_late", Thing: "hero", From: "merged", To: "late", At: "t4", Scope: "plot", Branch: "branch_a",
				})
			},
			wantCode: "merge.after_action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := testfixture.LoadMinimal(t)
			tt.mutate(d)
			issues := BranchIssuesToValidate(d)
			if !hasIssueCode(issues, tt.wantCode) {
				t.Fatalf("want code %q, got %v", tt.wantCode, issues)
			}
		})
	}
}

func BranchIssuesToValidate(d *project.Data) []Issue {
	var out []Issue
	for _, bi := range project.BranchIssues(d) {
		out = append(out, Issue{bi.Code, bi.Message})
	}
	return out
}

func hasIssueCode(issues []Issue, code string) bool {
	for _, iss := range issues {
		if iss.Code == code {
			return true
		}
	}
	return false
}