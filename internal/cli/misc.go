package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"novel-logic/internal/generate"
	"novel-logic/internal/lean"
	"novel-logic/internal/project"
	"novel-logic/internal/template"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show project summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		fmt.Printf("title: %s\n", d.Meta.Title)
		fmt.Printf("path: %s\n", d.Root)
		fmt.Printf("things: %d\n", len(d.Things))
		fmt.Printf("scenes: %d\n", len(d.Scenes))
		fmt.Printf("facts: %d\n", len(d.Facts))
		fmt.Printf("actions: %d\n", len(d.Actions))
		fmt.Printf("rules: %d\n", len(d.Rules))
		fmt.Printf("times: %d\n", len(d.Meta.TimeOrder))
		if d.Meta.LastCheck != nil {
			fmt.Printf("last_check: %s success=%v stage1=%v stage2=%v\n",
				d.Meta.LastCheck.At.Format("2006-01-02 15:04:05"),
				d.Meta.LastCheck.Success, d.Meta.LastCheck.Stage1, d.Meta.LastCheck.Stage2)
		}
		return nil
	},
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose environment and project files",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("novel-logic: %s\n", Version)
		tc := lean.Detect()
		fmt.Printf("elan: %s\n", pathOrMissing(tc.Elan))
		fmt.Printf("lean: %s (%s)\n", pathOrMissing(tc.Lean), tc.Version())
		fmt.Printf("lake: %s\n", pathOrMissing(tc.Lake))
		root, _ := filepath.Abs(projectPath)
		required := []string{
			project.FileProject, project.FilePlot, project.FileThings, project.FileScenes,
			project.FileTimes, project.FileFacts, project.FileActions, project.FileRules, project.FileNovels,
		}
		recommended := []string{
			project.FileBranches, project.FileForks, project.FileMerges,
		}
		for _, f := range required {
			p := filepath.Join(root, f)
			if _, err := os.Stat(p); err != nil {
				fmt.Printf("missing: %s\n", f)
			}
		}
		for _, f := range recommended {
			p := filepath.Join(root, f)
			if _, err := os.Stat(p); err != nil {
				fmt.Printf("recommended_missing: %s\n", f)
			}
		}
		if !tc.Found {
			return exitErrf(5, "lean toolchain incomplete")
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show CLI and Lean core versions",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("novel-logic %s\n", Version)
		fmt.Printf("lean-core %s\n", generate.CoreVersion)
	},
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Template operations",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available init templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := template.List()
		if err != nil {
			return exitErr(4, err)
		}
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show health summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		tc := lean.Detect()
		fmt.Printf("lean_toolchain: %v\n", tc.Found)
		fmt.Printf("entities: things=%d scenes=%d facts=%d actions=%d rules=%d times=%d\n",
			len(d.Things), len(d.Scenes), len(d.Facts), len(d.Actions), len(d.Rules), len(d.Meta.TimeOrder))
		if d.Meta.LastCheck != nil {
			fmt.Printf("last_check_success: %v\n", d.Meta.LastCheck.Success)
			fmt.Printf("last_check_stage1: %v\n", d.Meta.LastCheck.Stage1)
			fmt.Printf("last_check_stage2: %v\n", d.Meta.LastCheck.Stage2)
		} else {
			fmt.Println("last_check: never")
		}
		return nil
	},
}

func pathOrMissing(p string) string {
	if p == "" {
		return "(not found)"
	}
	return p
}

func init() {
	templateCmd.AddCommand(templateListCmd)
	rootCmd.AddCommand(infoCmd, doctorCmd, versionCmd, templateCmd, statusCmd)
}