package cli

import (
	"errors"
	"fmt"
	"strings"

	"novel-logic/internal/project"
	"novel-logic/internal/validate"
)

func loadProject() (*project.Data, error) {
	return project.Load(projectPath)
}

func formatIssues(issues []validate.Issue) string {
	var b strings.Builder
	for _, iss := range issues {
		fmt.Fprintf(&b, "[%s] %s\n", iss.Code, iss.Message)
	}
	return strings.TrimRight(b.String(), "\n")
}

func requireValidate(d *project.Data) error {
	issues := validate.Run(d)
	if len(issues) == 0 {
		return nil
	}
	return exitErrf(1, "validation failed:\n%s", formatIssues(issues))
}

func saveValidated(d *project.Data, mutate func() error) error {
	if err := mutate(); err != nil {
		var reg *project.RegistrationError
		if errors.As(err, &reg) {
			return exitErr(1, err)
		}
		return exitErr(4, err)
	}
	if err := requireValidate(d); err != nil {
		return err
	}
	if err := project.Save(d); err != nil {
		return exitErr(4, err)
	}
	return nil
}
