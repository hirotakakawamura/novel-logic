package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"novel-logic/internal/generate"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Lean project under logic/",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		if err := requireValidate(d); err != nil {
			return err
		}
		if err := generate.Run(d); err != nil {
			return exitErrf(2, "generate failed: %v", err)
		}
		if !quiet {
			fmt.Println("generated logic/")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}