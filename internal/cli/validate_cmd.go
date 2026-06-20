package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"novel-logic/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Run Stage 1 validation",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		branch, _ := cmd.Flags().GetString("branch")
		issues := validate.RunForBranch(d, branch)
		if len(issues) > 0 {
			if !quiet {
				fmt.Println(formatIssues(issues))
			}
			return exitErrf(1, "validation failed (%d issues)", len(issues))
		}
		hints := validate.Hints(d)
		if !quiet && (verbose || len(hints) > 0) {
			for _, h := range hints {
				fmt.Printf("[hint] %s\n", h.Message)
			}
		}
		if !quiet {
			fmt.Println("OK: stage1")
		}
		return nil
	},
}

func init() {
	validateCmd.Flags().String("branch", "", "validate a single story branch (default: all branches)")
	rootCmd.AddCommand(validateCmd)
}