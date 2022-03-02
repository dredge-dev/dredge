package cmd

import (
	"fmt"

	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/dredge-dev/dredge/internal/workflow"
)

func runInitCommand(e *exec.DredgeExec, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough arguments: missing source")
	}

	de, err := e.Import(args[0])
	if err != nil {
		return err
	}

	w, err := de.GetWorkflow("", "init")
	if err != nil {
		return err
	}

	return workflow.ExecuteWorkflow(w)
}
