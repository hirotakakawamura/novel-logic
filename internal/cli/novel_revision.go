package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"novel-logic/internal/project"
)

var novelRevisionCmd = &cobra.Command{
	Use:   "revision",
	Short: "Git revision pinning for novel body files",
}

var novelRevisionPinCmd = &cobra.Command{
	Use:   "pin <scene_id>",
	Short: "Pin the git commit for a novel body file",
	Long: `Records the git commit that the body file corresponds to in novels.yaml.
Use after committing prose changes so check/CI can detect revision drift.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		revision, _ := cmd.Flags().GetString("revision")
		branch, _ := cmd.Flags().GetString("branch")
		note, _ := cmd.Flags().GetString("note")
		allowDirty, _ := cmd.Flags().GetBool("allow-dirty")
		return saveValidated(d, func() error {
			pinned, err := d.PinNovelRevision(args[0], branch, revision, note, allowDirty)
			if err != nil {
				return err
			}
			if !quiet {
				fmt.Printf("pinned novel %s at %s", args[0], pinned.Revision)
				if pinned.Short != "" {
					fmt.Printf(" (%s)", pinned.Short)
				}
				if pinned.Branch != "" {
					fmt.Printf(" on %s", pinned.Branch)
				}
				if pinned.Dirty {
					fmt.Print(" [dirty]")
				}
				fmt.Println()
			}
			return nil
		})
	},
}

var novelRevisionListCmd = &cobra.Command{
	Use:   "list <scene_id>",
	Short: "List pinned git revisions for a novel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := loadProject()
		if err != nil {
			return exitErr(4, err)
		}
		branch, _ := cmd.Flags().GetString("branch")
		n, _ := d.FindNovel(args[0], branch)
		if n == nil {
			return exitErrf(4, "novel for scene %q not found", args[0])
		}
		if len(n.Revisions) == 0 {
			fmt.Println("(no revisions pinned)")
			return nil
		}
		for i, r := range n.Revisions {
			ts := ""
			if !r.RecordedAt.IsZero() {
				ts = r.RecordedAt.UTC().Format(time.RFC3339)
			}
			line := fmt.Sprintf("%d. %s", i+1, r.Revision)
			if r.Short != "" {
				line += fmt.Sprintf(" (%s)", r.Short)
			}
			if r.Branch != "" {
				line += fmt.Sprintf(" branch=%s", r.Branch)
			}
			if ts != "" {
				line += fmt.Sprintf(" at=%s", ts)
			}
			if r.Dirty {
				line += " dirty=true"
			}
			if r.Note != "" {
				line += fmt.Sprintf(" note=%q", r.Note)
			}
			fmt.Println(line)
		}
		if n.Revision != "" {
			fmt.Printf("current: %s\n", n.Revision)
		}
		return nil
	},
}

func initNovelRevisionCommands() {
	novelRevisionPinCmd.Flags().String("branch", project.MainBranch, "story branch id")
	novelRevisionPinCmd.Flags().String("revision", "", "explicit git commit SHA (default: latest commit touching the body file)")
	novelRevisionPinCmd.Flags().String("note", "", "note for this pin (e.g. PR number)")
	novelRevisionPinCmd.Flags().Bool("allow-dirty", false, "allow pinning with uncommitted working tree changes")
	novelRevisionListCmd.Flags().String("branch", project.MainBranch, "story branch id")
	novelRevisionCmd.AddCommand(novelRevisionPinCmd, novelRevisionListCmd)
	novelCmd.AddCommand(novelRevisionCmd)
}