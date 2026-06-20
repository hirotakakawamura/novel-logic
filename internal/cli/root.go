package cli

import (
	"github.com/spf13/cobra"
)

var (
	projectPath string
	quiet       bool
	verbose     bool
)

const Version = "0.1.0-mvp"

var rootCmd = &cobra.Command{
	Use:          "novel-logic",
	Short:        "Manage novel plot logic and verify consistency",
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&projectPath, "project", "C", ".", "project root directory")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-error output")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
}