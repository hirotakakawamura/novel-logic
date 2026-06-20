package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"novel-logic/internal/project"
	"novel-logic/internal/validate"
)

var plotCmd = &cobra.Command{Use: "plot", Short: "Plot operations"}
var plotShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show plot summary and scenes",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		fmt.Printf("summary: %s\n", d.Plot.Summary)
		for _, s := range d.Scenes {
			fmt.Printf("- %s [%s..%s] %s\n", s.ID, s.TimeStart, s.TimeEnd, s.Summary)
		}
		return nil
	},
}

var thingCmd = &cobra.Command{Use: "thing", Short: "Thing operations"}
var thingAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a new thing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringArray("tag")
		scopes, _ := cmd.Flags().GetStringArray("scope")
		return saveValidated(d, func() error {
			return d.AddThing(args[0], name, tags, scopes)
		})
	},
}

var thingScopeCmd = &cobra.Command{Use: "scope", Short: "Thing scope operations"}
var thingScopeAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add scopes to an existing thing",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		scopes, _ := cmd.Flags().GetStringArray("scope")
		if len(scopes) == 0 {
			return exitErrf(4, "at least one --scope is required")
		}
		return saveValidated(d, func() error {
			return d.AddThingScopes(args[0], scopes)
		})
	},
}

var sceneCmd = &cobra.Command{Use: "scene", Short: "Scene operations"}
var sceneAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a scene",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		summary, _ := cmd.Flags().GetString("summary")
		start, _ := cmd.Flags().GetString("time-start")
		end, _ := cmd.Flags().GetString("time-end")
		return saveValidated(d, func() error {
			return d.AddScene(args[0], summary, start, end)
		})
	},
}

var timeCmd = &cobra.Command{Use: "time", Short: "Time operations"}
var timeAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Add a time point",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		after, _ := cmd.Flags().GetString("after")
		return saveValidated(d, func() error {
			return d.AddTime(args[0], after)
		})
	},
}

var factCmd = &cobra.Command{Use: "fact", Short: "Fact operations"}
var factAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add fixed_fact or state",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		kindStr, _ := cmd.Flags().GetString("kind")
		thing, _ := cmd.Flags().GetString("thing")
		pred, _ := cmd.Flags().GetString("pred")
		scope, _ := cmd.Flags().GetString("scope")
		branch, _ := cmd.Flags().GetString("branch")
		kind := project.FactKind(kindStr)
		if err := validate.CheckPredNotThingID(d, pred); err != nil {
			return exitErr(1, err)
		}
		if kind == project.FactState {
			if err := validate.CheckForbidState(d, thing, pred); err != nil {
				return exitErr(1, err)
			}
		}
		return saveValidated(d, func() error {
			_, err := d.AddFact(kind, thing, pred, scope, branch)
			return err
		})
	},
}
var factPromoteCmd = &cobra.Command{
	Use:   "promote <id>",
	Short: "Promote fixed_fact to state",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		return saveValidated(d, func() error {
			return d.PromoteFact(args[0])
		})
	},
}

var actionCmd = &cobra.Command{Use: "action", Short: "Action operations"}
var actionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an action",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		thing, _ := cmd.Flags().GetString("thing")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		at, _ := cmd.Flags().GetString("at")
		scope, _ := cmd.Flags().GetString("scope")
		label, _ := cmd.Flags().GetString("label")
		branch, _ := cmd.Flags().GetString("branch")
		if from != "" {
			if err := validate.CheckPredNotThingID(d, from); err != nil {
				return exitErr(1, err)
			}
		}
		if err := validate.CheckPredNotThingID(d, to); err != nil {
			return exitErr(1, err)
		}
		a := project.Action{Thing: thing, From: from, To: to, At: at, Scope: scope, Label: label}
		if err := validate.CheckActionRules(d, a); err != nil {
			return exitErr(1, err)
		}
		return saveValidated(d, func() error {
			_, err := d.AddAction(thing, from, to, at, scope, label, branch)
			return err
		})
	},
}

var ruleCmd = &cobra.Command{Use: "rule", Short: "Rule operations"}
var ruleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		kindStr, _ := cmd.Flags().GetString("kind")
		thing, _ := cmd.Flags().GetString("thing")
		pred, _ := cmd.Flags().GetString("pred")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		branch, _ := cmd.Flags().GetString("branch")
		return saveValidated(d, func() error {
			_, err := d.AddRule(project.RuleKind(kindStr), thing, pred, from, to, branch)
			return err
		})
	},
}

var novelCmd = &cobra.Command{
	Use:   "novel",
	Short: "Novel operations (body text is git-managed under novels/)",
}
var novelAddCmd = &cobra.Command{
	Use:   "add <scene_id>",
	Short: "Register a git-managed novel body file for a scene",
	Long: `Registers novel metadata in novels.yaml. Body text is not written by the CLI;
edit novels/<branch>/<scene_id>.txt (or --file path) in your editor and commit with git.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		filePath, _ := cmd.Flags().GetString("file")
		branch, _ := cmd.Flags().GetString("branch")
		initFile, _ := cmd.Flags().GetBool("init")
		pin, _ := cmd.Flags().GetBool("pin")
		note, _ := cmd.Flags().GetString("note")
		allowDirty, _ := cmd.Flags().GetBool("allow-dirty")
		return saveValidated(d, func() error {
			if err := d.AddNovel(args[0], branch, filePath, initFile); err != nil {
				return err
			}
			if pin {
				_, err := d.PinNovelRevision(args[0], branch, "", note, allowDirty)
				return err
			}
			return nil
		})
	},
}
var novelUpdateCmd = &cobra.Command{
	Use:   "update <scene_id>",
	Short: "Update novel metadata (body path or scene time sync)",
	Long:  "Updates novels.yaml only. Edit the body .txt file directly for prose changes.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		filePath, _ := cmd.Flags().GetString("file")
		branch, _ := cmd.Flags().GetString("branch")
		return saveValidated(d, func() error {
			return d.UpdateNovel(args[0], branch, filePath, cmd.Flags().Changed("file"))
		})
	},
}

func init() {
	thingAddCmd.Flags().String("name", "", "display name")
	thingAddCmd.Flags().StringArray("tag", nil, "tag (repeatable)")
	thingAddCmd.Flags().StringArray("scope", nil, "scope: plot and/or novel:<scene_id> (repeatable; default plot)")
	thingScopeAddCmd.Flags().StringArray("scope", nil, "scope to add (repeatable)")
	thingScopeCmd.AddCommand(thingScopeAddCmd)
	thingCmd.AddCommand(thingAddCmd, thingScopeCmd)

	sceneAddCmd.Flags().String("summary", "", "scene summary")
	sceneAddCmd.Flags().String("time-start", "", "start time id")
	sceneAddCmd.Flags().String("time-end", "", "end time id")
	sceneCmd.AddCommand(sceneAddCmd)

	timeAddCmd.Flags().String("after", "", "insert after time id")
	timeCmd.AddCommand(timeAddCmd)

	factAddCmd.Flags().String("kind", "", "fixed or state")
	factAddCmd.Flags().String("thing", "", "subject thing id")
	factAddCmd.Flags().String("pred", "", "predicate")
	factAddCmd.Flags().String("scope", "plot", "plot or novel:<scene_id>")
	factAddCmd.Flags().String("branch", project.MainBranch, "story branch id")
	factCmd.AddCommand(factAddCmd, factPromoteCmd)

	actionAddCmd.Flags().String("thing", "", "subject thing id")
	actionAddCmd.Flags().String("from", "", "from predicate")
	actionAddCmd.Flags().String("to", "", "to predicate")
	actionAddCmd.Flags().String("at", "", "time id")
	actionAddCmd.Flags().String("scope", "plot", "plot or novel:<scene_id>")
	actionAddCmd.Flags().String("label", "", "human label")
	actionAddCmd.Flags().String("branch", project.MainBranch, "story branch id")
	actionCmd.AddCommand(actionAddCmd)

	ruleAddCmd.Flags().String("kind", "", "forbid-state or forbid-transition")
	ruleAddCmd.Flags().String("thing", "", "thing id")
	ruleAddCmd.Flags().String("pred", "", "predicate")
	ruleAddCmd.Flags().String("from", "", "from predicate")
	ruleAddCmd.Flags().String("to", "", "to predicate")
	ruleAddCmd.Flags().String("branch", project.MainBranch, "story branch id")
	ruleCmd.AddCommand(ruleAddCmd)

	novelAddCmd.Flags().String("branch", project.MainBranch, "story branch id")
	novelAddCmd.Flags().String("file", "", "body path relative to project (default novels/<branch>/<scene_id>.txt)")
	novelAddCmd.Flags().Bool("init", true, "create empty body file if missing (for git tracking)")
	novelAddCmd.Flags().Bool("pin", false, "pin git revision after registration (requires git commit)")
	novelAddCmd.Flags().String("note", "", "note stored with git revision pin")
	novelAddCmd.Flags().Bool("allow-dirty", false, "allow revision pin with uncommitted body changes")
	novelUpdateCmd.Flags().String("branch", project.MainBranch, "story branch id")
	novelUpdateCmd.Flags().String("file", "", "new body path relative to project")
	novelCmd.AddCommand(novelAddCmd, novelUpdateCmd)

	plotCmd.AddCommand(plotShowCmd)
	rootCmd.AddCommand(plotCmd, thingCmd, sceneCmd, timeCmd, factCmd, actionCmd, ruleCmd, novelCmd)
	initShowCommands()
	initUpdateCommands()
	initNovelRevisionCommands()
	initBranchCommands()
}