package cli

import (
	"github.com/spf13/cobra"

	"novel-logic/internal/project"
	"novel-logic/internal/validate"
)

var thingUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing thing (name or tags)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringArray("tag")
		if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("tag") {
			return exitErrf(4, "specify at least one of --name or --tag")
		}
		return saveValidated(d, func() error {
			return d.UpdateThing(args[0], name, tags, cmd.Flags().Changed("tag"))
		})
	},
}

var factUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing fact",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		f, _ := d.FindFact(args[0])
		if f == nil {
			return exitErrf(4, "fact %q not found", args[0])
		}
		kind := f.Kind
		thing := f.Thing
		pred := f.Pred
		scope := f.Scope
		if cmd.Flags().Changed("kind") {
			kind = project.FactKind(mustString(cmd, "kind"))
		}
		if cmd.Flags().Changed("thing") {
			thing = mustString(cmd, "thing")
		}
		if cmd.Flags().Changed("pred") {
			pred = mustString(cmd, "pred")
		}
		if cmd.Flags().Changed("scope") {
			scope = mustString(cmd, "scope")
		}
		if err := validate.CheckPredNotThingID(d, pred); err != nil {
			return exitErr(1, err)
		}
		if kind == project.FactState {
			if err := validate.CheckForbidState(d, project.NormalizeBranch(f.Branch), thing, pred); err != nil {
				return exitErr(1, err)
			}
		}
		return saveValidated(d, func() error {
			return d.UpdateFact(args[0], kind, thing, pred, scope)
		})
	},
}

var actionUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing action",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		a, _ := d.FindAction(args[0])
		if a == nil {
			return exitErrf(4, "action %q not found", args[0])
		}
		thing := a.Thing
		from := a.From
		to := a.To
		at := a.At
		scope := a.Scope
		label := a.Label
		if cmd.Flags().Changed("thing") {
			thing = mustString(cmd, "thing")
		}
		if cmd.Flags().Changed("from") {
			from = mustString(cmd, "from")
		}
		if cmd.Flags().Changed("to") {
			to = mustString(cmd, "to")
		}
		if cmd.Flags().Changed("at") {
			at = mustString(cmd, "at")
		}
		if cmd.Flags().Changed("scope") {
			scope = mustString(cmd, "scope")
		}
		if cmd.Flags().Changed("label") {
			label = mustString(cmd, "label")
		}
		if from != "" {
			if err := validate.CheckPredNotThingID(d, from); err != nil {
				return exitErr(1, err)
			}
		}
		if err := validate.CheckPredNotThingID(d, to); err != nil {
			return exitErr(1, err)
		}
		updated := project.Action{Thing: thing, From: from, To: to, At: at, Scope: scope, Label: label, Branch: a.Branch}
		if err := validate.CheckActionRules(d, updated); err != nil {
			return exitErr(1, err)
		}
		return saveValidated(d, func() error {
			return d.UpdateAction(args[0], thing, from, to, at, scope, label)
		})
	},
}

var ruleUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an existing rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		r, _ := d.FindRule(args[0])
		if r == nil {
			return exitErrf(4, "rule %q not found", args[0])
		}
		kind := r.Kind
		thing := r.Thing
		pred := r.Pred
		from := r.From
		to := r.To
		if cmd.Flags().Changed("kind") {
			kind = project.RuleKind(mustString(cmd, "kind"))
		}
		if cmd.Flags().Changed("thing") {
			thing = mustString(cmd, "thing")
		}
		if cmd.Flags().Changed("pred") {
			pred = mustString(cmd, "pred")
		}
		if cmd.Flags().Changed("from") {
			from = mustString(cmd, "from")
		}
		if cmd.Flags().Changed("to") {
			to = mustString(cmd, "to")
		}
		return saveValidated(d, func() error {
			return d.UpdateRule(args[0], kind, thing, pred, from, to)
		})
	},
}

func mustString(cmd *cobra.Command, name string) string {
	v, _ := cmd.Flags().GetString(name)
	return v
}

func initUpdateCommands() {
	thingUpdateCmd.Flags().String("name", "", "display name")
	thingUpdateCmd.Flags().StringArray("tag", nil, "replace tags (repeatable)")
	thingCmd.AddCommand(thingUpdateCmd)

	factUpdateCmd.Flags().String("kind", "", "fixed or state")
	factUpdateCmd.Flags().String("thing", "", "subject thing id")
	factUpdateCmd.Flags().String("pred", "", "predicate")
	factUpdateCmd.Flags().String("scope", "", "plot or novel:<scene_id>")
	factCmd.AddCommand(factUpdateCmd)

	actionUpdateCmd.Flags().String("thing", "", "subject thing id")
	actionUpdateCmd.Flags().String("from", "", "from predicate")
	actionUpdateCmd.Flags().String("to", "", "to predicate")
	actionUpdateCmd.Flags().String("at", "", "time id")
	actionUpdateCmd.Flags().String("scope", "", "plot or novel:<scene_id>")
	actionUpdateCmd.Flags().String("label", "", "human label")
	actionCmd.AddCommand(actionUpdateCmd)

	ruleUpdateCmd.Flags().String("kind", "", "forbid-state or forbid-transition")
	ruleUpdateCmd.Flags().String("thing", "", "thing id")
	ruleUpdateCmd.Flags().String("pred", "", "predicate")
	ruleUpdateCmd.Flags().String("from", "", "from predicate")
	ruleUpdateCmd.Flags().String("to", "", "to predicate")
	ruleCmd.AddCommand(ruleUpdateCmd)
}