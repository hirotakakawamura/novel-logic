package cli

import (
	"github.com/spf13/cobra"

	"novel-logic/internal/project"
)

func removeCmd(use, short string, run func(*cobra.Command, []string) error) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		RunE:  run,
	}
}

var thingRemoveCmd = removeCmd("remove <id>", "Remove a thing", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveThing(args[0])
	})
})

var thingScopeRemoveCmd = removeCmd("remove <id>", "Remove scopes from a thing", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	scopes, _ := cmd.Flags().GetStringArray("scope")
	if len(scopes) == 0 {
		return exitErrf(4, "at least one --scope is required")
	}
	return saveValidated(d, func() error {
		return d.RemoveThingScopes(args[0], scopes)
	})
})

var sceneRemoveCmd = removeCmd("remove <id>", "Remove a scene", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveScene(args[0])
	})
})

var timeRemoveCmd = removeCmd("remove <id>", "Remove a time point", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveTime(args[0])
	})
})

var factRemoveCmd = removeCmd("remove <id>", "Remove a fact", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveFact(args[0])
	})
})

var actionRemoveCmd = removeCmd("remove <id>", "Remove an action", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveAction(args[0])
	})
})

var ruleRemoveCmd = removeCmd("remove <id>", "Remove a rule", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	return saveValidated(d, func() error {
		return d.RemoveRule(args[0])
	})
})

var novelRemoveCmd = removeCmd("remove <scene_id>", "Remove novel for a scene", func(cmd *cobra.Command, args []string) error {
	d, err := loadProject()
	if err != nil {
		return exitErr(4, err)
	}
	keepBody, _ := cmd.Flags().GetBool("keep-body")
	branch, _ := cmd.Flags().GetString("branch")
	return saveValidated(d, func() error {
		return d.RemoveNovel(args[0], branch, keepBody)
	})
})

func initRemoveCommands() {
	thingRemoveCmd.Args = cobra.ExactArgs(1)
	sceneRemoveCmd.Args = cobra.ExactArgs(1)
	timeRemoveCmd.Args = cobra.ExactArgs(1)
	factRemoveCmd.Args = cobra.ExactArgs(1)
	actionRemoveCmd.Args = cobra.ExactArgs(1)
	ruleRemoveCmd.Args = cobra.ExactArgs(1)
	novelRemoveCmd.Args = cobra.ExactArgs(1)
	novelRemoveCmd.Flags().String("branch", project.MainBranch, "story branch id")
	novelRemoveCmd.Flags().Bool("keep-body", true, "keep git-tracked body file on disk")
	thingScopeRemoveCmd.Args = cobra.ExactArgs(1)

	thingScopeRemoveCmd.Flags().StringArray("scope", nil, "scope to remove (repeatable)")

	thingCmd.AddCommand(thingRemoveCmd)
	thingScopeCmd.AddCommand(thingScopeRemoveCmd)
	sceneCmd.AddCommand(sceneRemoveCmd)
	timeCmd.AddCommand(timeRemoveCmd)
	factCmd.AddCommand(factRemoveCmd)
	actionCmd.AddCommand(actionRemoveCmd)
	ruleCmd.AddCommand(ruleRemoveCmd)
	novelCmd.AddCommand(novelRemoveCmd)
}

func init() {
	initRemoveCommands()
}