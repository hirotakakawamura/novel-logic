package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"novel-logic/internal/project"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Story branch operations",
}

var branchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		if len(d.Branches) == 0 {
			fmt.Println("main  (implicit)")
			return nil
		}
		for _, b := range d.Branches {
			line := fmt.Sprintf("%s", b.ID)
			if b.Label != "" {
				line += fmt.Sprintf("  %s", b.Label)
			}
			if b.Parent != "" {
				line += fmt.Sprintf("  parent=%s", b.Parent)
			}
			if b.ViaFork != "" {
				line += fmt.Sprintf("  via_fork=%s via_action=%s", b.ViaFork, b.ViaAction)
			}
			if d.FindMergeForBranch(b.ID) != nil {
				line += "  [merged]"
			}
			fmt.Println(line)
		}
		return nil
	},
}

var branchShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show one branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		b, _ := d.FindBranchDef(args[0])
		if b == nil {
			return exitErrf(4, "branch %q not found", args[0])
		}
		fmt.Printf("id: %s\n", b.ID)
		if b.Label != "" {
			fmt.Printf("label: %s\n", b.Label)
		}
		if b.Parent != "" {
			fmt.Printf("parent: %s\n", b.Parent)
		}
		if b.ViaFork != "" {
			fmt.Printf("via_fork: %s\n", b.ViaFork)
			fmt.Printf("via_action: %s\n", b.ViaAction)
		}
		lineage := d.BranchLineage(b.ID)
		fmt.Printf("lineage: %s\n", strings.Join(lineage, " -> "))
		fmt.Printf("active_actions: %d\n", len(d.ActiveActions(b.ID)))
		if merge := d.FindMergeForBranch(b.ID); merge != nil {
			fmt.Printf("merge: %s at %s into %s\n", merge.ID, merge.At, merge.IntoBranch)
		}
		return nil
	},
}

var branchAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a branch definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		label, _ := cmd.Flags().GetString("label")
		parent, _ := cmd.Flags().GetString("parent")
		viaFork, _ := cmd.Flags().GetString("via-fork")
		viaAction, _ := cmd.Flags().GetString("via-action")
		return saveValidated(d, func() error {
			return d.AddBranch(args[0], label, parent, viaFork, viaAction)
		})
	},
}

var branchRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a branch definition",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		return saveValidated(d, func() error {
			return d.RemoveBranch(args[0])
		})
	},
}

var forkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Fork point operations",
}

var forkAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a fork point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		parent, _ := cmd.Flags().GetString("parent")
		at, _ := cmd.Flags().GetString("at")
		scope, _ := cmd.Flags().GetString("scope")
		return saveValidated(d, func() error {
			return d.AddFork(args[0], parent, at, scope)
		})
	},
}

var forkChoiceAddCmd = &cobra.Command{
	Use:   "choice add",
	Short: "Add a fork choice (creates child branch)",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		forkID, _ := cmd.Flags().GetString("fork")
		actionID, _ := cmd.Flags().GetString("action")
		branchID, _ := cmd.Flags().GetString("branch")
		if forkID == "" || actionID == "" || branchID == "" {
			return exitErrf(4, "--fork, --action, and --branch are required")
		}
		return saveValidated(d, func() error {
			return d.AddForkChoice(forkID, actionID, branchID)
		})
	},
}

var forkShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a fork",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		f, _ := d.FindFork(args[0])
		if f == nil {
			return exitErrf(4, "fork %q not found", args[0])
		}
		fmt.Printf("id: %s\n", f.ID)
		fmt.Printf("parent_branch: %s\n", f.ParentBranch)
		fmt.Printf("at: %s\n", f.At)
		fmt.Printf("scope: %s\n", scopeOrPlot(f.Scope))
		for _, c := range f.Choices {
			fmt.Printf("  choice: action=%s -> branch=%s\n", c.Action, c.Branch)
		}
		return nil
	},
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge point operations",
}

var mergeAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a merge point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		at, _ := cmd.Flags().GetString("at")
		into, _ := cmd.Flags().GetString("into")
		scope, _ := cmd.Flags().GetString("scope")
		choices, err := parseMergeChoices(cmd)
		if err != nil {
			return exitErr(4, err)
		}
		return saveValidated(d, func() error {
			return d.AddMerge(args[0], at, scope, into, choices)
		})
	},
}

var mergeShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a merge",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		m, _ := d.FindMerge(args[0])
		if m == nil {
			return exitErrf(4, "merge %q not found", args[0])
		}
		fmt.Printf("id: %s\n", m.ID)
		fmt.Printf("at: %s\n", m.At)
		fmt.Printf("into_branch: %s\n", m.IntoBranch)
		fmt.Printf("scope: %s\n", scopeOrPlot(m.Scope))
		for _, c := range m.Choices {
			act, _ := d.FindAction(c.Action)
			to := "-"
			if act != nil {
				to = act.To
			}
			fmt.Printf("  from %s via %s (to=%s)\n", c.Branch, c.Action, to)
		}
		return nil
	},
}

func parseMergeChoices(cmd *cobra.Command) ([]project.MergeChoice, error) {
	pairs, _ := cmd.Flags().GetStringArray("choice")
	if len(pairs) == 0 {
		return nil, fmt.Errorf("at least one --choice branch:action is required")
	}
	var out []project.MergeChoice
	for _, p := range pairs {
		parts := strings.SplitN(p, ":", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid --choice %q (use branch:action)", p)
		}
		out = append(out, project.MergeChoice{Branch: parts[0], Action: parts[1]})
	}
	return out, nil
}

func initBranchCommands() {
	branchAddCmd.Flags().String("label", "", "display label")
	branchAddCmd.Flags().String("parent", "", "parent branch id")
	branchAddCmd.Flags().String("via-fork", "", "fork id that created this branch")
	branchAddCmd.Flags().String("via-action", "", "fork action id")
	branchCmd.AddCommand(branchListCmd, branchShowCmd, branchAddCmd, branchRemoveCmd)

	forkAddCmd.Flags().String("parent", project.MainBranch, "parent branch id")
	forkAddCmd.Flags().String("at", "", "time id at fork")
	forkAddCmd.Flags().String("scope", "plot", "plot or novel:<scene_id>")
	forkChoiceAddCmd.Flags().String("fork", "", "fork id")
	forkChoiceAddCmd.Flags().String("action", "", "route action id on parent branch")
	forkChoiceAddCmd.Flags().String("branch", "", "new child branch id")
	forkCmd.AddCommand(forkAddCmd, forkChoiceAddCmd, forkShowCmd)

	mergeAddCmd.Flags().String("at", "", "merge time id")
	mergeAddCmd.Flags().String("into", project.MainBranch, "branch to merge into")
	mergeAddCmd.Flags().String("scope", "plot", "plot or novel:<scene_id>")
	mergeAddCmd.Flags().StringArray("choice", nil, "from-branch merge action as branch:action (repeatable)")
	mergeCmd.AddCommand(mergeAddCmd, mergeShowCmd)

	rootCmd.AddCommand(branchCmd, forkCmd, mergeCmd)
}