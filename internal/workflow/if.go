package workflow

import (
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

func executeIfStep(workflow *exec.Workflow, ifStep *config.IfStep) error {
	cond, err := Template(ifStep.Cond, workflow.Exec.Env)
	if err != nil {
		return err
	}
	if isTrue(cond) {
		return executeSteps(workflow, ifStep.Steps)
	}
	return nil
}
