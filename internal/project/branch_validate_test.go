package project

import (
	"strings"
	"testing"
)

func TestBranchIssuesCleanProject(t *testing.T) {
	d := newTestProject(t)
	if issues := BranchIssues(d); len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestBranchIssuesIncludesIsolatedState(t *testing.T) {
	d := newTestProject(t)
	if err := d.AddBranch("branch_a", "alt", MainBranch, "", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := d.AddAction("hero", "mid", "alt_only", "t3", "plot", "", "branch_a"); err != nil {
		t.Fatal(err)
	}
	if _, err := d.AddAction("hero", "alt_only", "end", "t4", "plot", "", MainBranch); err != nil {
		t.Fatal(err)
	}
	if _, ok := branchIssueMatching(BranchIssues(d), "branch.isolated_state", ""); !ok {
		t.Fatal("expected branch.isolated_state via BranchIssues")
	}
}

func TestBranchIssues(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(t *testing.T, d *Data)
		wantCode string
		wantMsg  string
	}{
		{
			name: "fact_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Facts = append(d.Facts, Fact{
					ID: "fact_bad", Kind: FactState, Thing: "hero", Pred: "lost",
					Scope: "plot", Branch: "ghost_branch",
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "action_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Actions = append(d.Actions, Action{
					ID: "act_bad", Thing: "hero", From: "start", To: "mid",
					At: "t2", Scope: "plot", Branch: "ghost_branch",
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "rule_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Rules = append(d.Rules, Rule{
					ID: "rule_bad", Kind: RuleForbidState, Thing: "hero", Pred: "evil",
					Branch: "ghost_branch",
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "novel_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Novels = append(d.Novels, NovelMeta{
					SceneID: "scene1", Branch: "ghost_branch",
					TimeStart: "t1", TimeEnd: "t2",
					BodyPath: DefaultNovelBodyPath("scene1", "ghost_branch"),
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "main_branch_missing",
			mutate: func(t *testing.T, d *Data) {
				d.Branches = nil
			},
			wantCode: "branch.invalid",
		},
		{
			name: "fork_exclusive",
			mutate: func(t *testing.T, d *Data) {
				d.Forks = append(d.Forks,
					Fork{ID: "fork_a", ParentBranch: MainBranch, At: "t2", Scope: "plot"},
					Fork{ID: "fork_b", ParentBranch: MainBranch, At: "t2", Scope: "plot"},
				)
			},
			wantCode: "fork.exclusive",
		},
		{
			name: "fork_unknown_parent_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Forks = append(d.Forks, Fork{
					ID: "fork_bad", ParentBranch: "ghost_branch", At: "t2", Scope: "plot",
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "fork_unknown_time",
			mutate: func(t *testing.T, d *Data) {
				d.Forks = append(d.Forks, Fork{
					ID: "fork_bad", ParentBranch: MainBranch, At: "t99", Scope: "plot",
				})
			},
			wantCode: "fork.invalid",
		},
		{
			name: "fork_choice_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Forks = append(d.Forks, Fork{
					ID: "fork1", ParentBranch: MainBranch, At: "t2", Scope: "plot",
					Choices: []ForkChoice{{Action: "act1", Branch: "ghost_branch"}},
				})
			},
			wantCode: "fork.invalid",
		},
		{
			name: "fork_choice_unknown_action",
			mutate: func(t *testing.T, d *Data) {
				d.Forks = append(d.Forks, Fork{
					ID: "fork1", ParentBranch: MainBranch, At: "t2", Scope: "plot",
					Choices: []ForkChoice{{Action: "no_such_act", Branch: "branch_a"}},
				})
			},
			wantCode: "fork.invalid",
		},
		{
			name: "fork_choice_action_wrong_branch",
			mutate: func(t *testing.T, d *Data) {
				if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
					t.Fatal(err)
				}
				actA, err := d.AddAction("hero", "start", "alt", "t2", "plot", "", "branch_a")
				if err != nil {
					t.Fatal(err)
				}
				d.Forks = append(d.Forks, Fork{
					ID: "fork1", ParentBranch: MainBranch, At: "t2", Scope: "plot",
					Choices: []ForkChoice{{Action: actA.ID, Branch: "branch_a"}},
				})
			},
			wantCode: "fork.invalid",
		},
		{
			name: "fork_choice_action_wrong_at",
			mutate: func(t *testing.T, d *Data) {
				if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
					t.Fatal(err)
				}
				actLate, err := d.AddAction("hero", "mid", "late", "t3", "plot", "", MainBranch)
				if err != nil {
					t.Fatal(err)
				}
				d.Forks = append(d.Forks, Fork{
					ID: "fork1", ParentBranch: MainBranch, At: "t2", Scope: "plot",
					Choices: []ForkChoice{{Action: actLate.ID, Branch: "branch_a"}},
				})
			},
			wantCode: "fork.invalid",
			wantMsg:  `expected "t2"`,
		},
		{
			name: "merge_unknown_into_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: "ghost_branch",
					Choices: []MergeChoice{{Branch: MainBranch, Action: "act1"}},
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "merge_unknown_time",
			mutate: func(t *testing.T, d *Data) {
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t99", Scope: "plot", IntoBranch: MainBranch,
					Choices: []MergeChoice{{Branch: MainBranch, Action: "act1"}},
				})
			},
			wantCode: "merge.invalid",
		},
		{
			name: "merge_choice_unknown_branch",
			mutate: func(t *testing.T, d *Data) {
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: MainBranch,
					Choices: []MergeChoice{{Branch: "ghost_branch", Action: "act1"}},
				})
			},
			wantCode: "branch.unknown",
		},
		{
			name: "merge_unknown_action",
			mutate: func(t *testing.T, d *Data) {
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: MainBranch,
					Choices: []MergeChoice{{Branch: MainBranch, Action: "no_such_act"}},
				})
			},
			wantCode: "merge.invalid",
		},
		{
			name: "merge_action_wrong_branch",
			mutate: func(t *testing.T, d *Data) {
				if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
					t.Fatal(err)
				}
				actA, err := d.AddAction("hero", "mid", "merged", "t3", "plot", "", "branch_a")
				if err != nil {
					t.Fatal(err)
				}
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: MainBranch,
					Choices: []MergeChoice{{Branch: MainBranch, Action: actA.ID}},
				})
			},
			wantCode: "merge.invalid",
		},
		{
			name: "merge_action_wrong_at",
			mutate: func(t *testing.T, d *Data) {
				actEarly, err := d.AddAction("hero", "mid", "early", "t2", "plot", "", MainBranch)
				if err != nil {
					t.Fatal(err)
				}
				d.Merges = append(d.Merges, Merge{
					ID: "merge_bad", At: "t3", Scope: "plot", IntoBranch: MainBranch,
					Choices: []MergeChoice{{Branch: MainBranch, Action: actEarly.ID}},
				})
			},
			wantCode: "merge.invalid",
		},
		{
			name: "merge_action_mismatch",
			mutate: func(t *testing.T, d *Data) {
				if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
					t.Fatal(err)
				}
				actMain, err := d.AddAction("hero", "mid", "merge_a", "t3", "plot", "", MainBranch)
				if err != nil {
					t.Fatal(err)
				}
				actA, err := d.AddAction("hero", "mid", "merge_b", "t3", "plot", "", "branch_a")
				if err != nil {
					t.Fatal(err)
				}
				if err := d.AddMerge("merge1", "t3", "plot", MainBranch, []MergeChoice{
					{Branch: MainBranch, Action: actMain.ID},
					{Branch: "branch_a", Action: actA.ID},
				}); err != nil {
					t.Fatal(err)
				}
			},
			wantCode: "merge.action_mismatch",
		},
		{
			name: "merge_after_action",
			mutate: func(t *testing.T, d *Data) {
				if err := d.AddBranch("branch_a", "", MainBranch, "", ""); err != nil {
					t.Fatal(err)
				}
				actMain, err := d.AddAction("hero", "mid", "merged", "t3", "plot", "", MainBranch)
				if err != nil {
					t.Fatal(err)
				}
				actA, err := d.AddAction("hero", "mid", "merged", "t3", "plot", "", "branch_a")
				if err != nil {
					t.Fatal(err)
				}
				if err := d.AddMerge("merge1", "t3", "plot", MainBranch, []MergeChoice{
					{Branch: MainBranch, Action: actMain.ID},
					{Branch: "branch_a", Action: actA.ID},
				}); err != nil {
					t.Fatal(err)
				}
				d.Actions = append(d.Actions, Action{
					ID: "act_late", Thing: "hero", From: "merged", To: "late",
					At: "t4", Scope: "plot", Branch: "branch_a",
				})
			},
			wantCode: "merge.after_action",
		},
		{
			name: "novel_duplicate",
			mutate: func(t *testing.T, d *Data) {
				meta := NovelMeta{
					SceneID: "scene1", Branch: MainBranch,
					TimeStart: "t1", TimeEnd: "t2",
					BodyPath: DefaultNovelBodyPath("scene1", MainBranch),
				}
				d.Novels = append(d.Novels, meta, meta)
			},
			wantCode: "novel.duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := newTestProject(t)
			tt.mutate(t, d)
			issues := BranchIssues(d)
			if _, ok := branchIssueMatching(issues, tt.wantCode, tt.wantMsg); !ok {
				t.Fatalf("want code %q msg %q, got %v", tt.wantCode, tt.wantMsg, issues)
			}
		})
	}
}

func branchIssueMatching(issues []BranchIssue, code, msgSub string) (*BranchIssue, bool) {
	for i := range issues {
		if issues[i].Code != code {
			continue
		}
		if msgSub != "" && !strings.Contains(issues[i].Message, msgSub) {
			continue
		}
		return &issues[i], true
	}
	return nil, false
}