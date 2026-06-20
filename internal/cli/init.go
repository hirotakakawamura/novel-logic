package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"novel-logic/internal/template"
)

var initCmd = &cobra.Command{
	Use:   "init <path>",
	Short: "Create a new project directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		tpl, _ := cmd.Flags().GetString("template")
		force, _ := cmd.Flags().GetBool("force")

		abs, err := filepath.Abs(path)
		if err != nil {
			return exitErr(4, err)
		}
		if fi, err := os.Stat(abs); err == nil {
			if !fi.IsDir() {
				return exitErrf(4, "%s exists and is not a directory", abs)
			}
			entries, _ := os.ReadDir(abs)
			if len(entries) > 0 && !force {
				return exitErrf(4, "%s is not empty (use --force)", abs)
			}
		} else if !os.IsNotExist(err) {
			return exitErr(4, err)
		}
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return exitErr(4, err)
		}
		if err := template.Materialize(abs, tpl); err != nil {
			return exitErr(4, err)
		}
		if !quiet {
			fmt.Printf("created project at %s (template: %s)\n", abs, tpl)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().String("template", "default", "template name (default, momotaro)")
	initCmd.Flags().Bool("force", false, "overwrite non-empty directory")
	rootCmd.AddCommand(initCmd)
}