package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"novel-logic/internal/generate"
	"novel-logic/internal/lean"
	"novel-logic/internal/project"
	"novel-logic/internal/validate"
)

var (
	checkQuick      bool
	checkNoGenerate bool
	checkJobs       int
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate, generate Lean, and lake build",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		issues := validate.Run(d)
		if len(issues) > 0 {
			_ = d.RecordCheck(false, false, false, formatIssues(issues))
			if !quiet {
				fmt.Println(formatIssues(issues))
			}
			return exitErrf(1, "stage1 failed (%d issues)", len(issues))
		}
		if !quiet {
			fmt.Println("OK: stage1")
		}
		if checkQuick {
			_ = d.RecordCheck(true, true, false, "quick")
			return nil
		}
		if !checkNoGenerate {
			if err := generate.Run(d); err != nil {
				_ = d.RecordCheck(false, true, false, err.Error())
				return exitErrf(2, "generate failed: %v", err)
			}
			if !quiet {
				fmt.Println("generated logic/")
			}
		}
		tc := lean.Detect()
		if !tc.Found {
			_ = d.RecordCheck(true, true, false, "lean not found; stage2 skipped")
			if !quiet {
				fmt.Fprintln(cmd.ErrOrStderr(), "warning: lean/lake not found; stage2 skipped")
			}
			return exitErrf(5, "lean/lake not found in PATH")
		}
		logicDir := filepath.Join(d.Root, project.DirLogic)
		out, err := lean.LakeBuild(logicDir, checkJobs)
		if err != nil {
			_ = d.RecordCheck(false, true, false, trimOutput(out))
			if !quiet {
				fmt.Println(trimOutput(out))
			}
			return exitErrf(3, "stage2 failed: %v", err)
		}
		_ = d.RecordCheck(true, true, true, "ok")
		if !quiet {
			fmt.Println("OK: stage1 + stage2")
		}
		return nil
	},
}

func trimOutput(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 4000 {
		return s[len(s)-4000:]
	}
	return s
}

func init() {
	checkCmd.Flags().BoolVar(&checkQuick, "quick", false, "Stage 1 only")
	checkCmd.Flags().BoolVar(&checkNoGenerate, "no-generate", false, "skip generate, run lake build only")
	checkCmd.Flags().IntVarP(&checkJobs, "jobs", "j", 0, "lake build -j")
	rootCmd.AddCommand(checkCmd)
}