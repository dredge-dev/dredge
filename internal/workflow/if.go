package workflow

import (
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

func (workflow *Workflow) executeIfStep(ifStep *config.IfStep) error {
	cond, err := workflow.Callbacks.Template(ifStep.Cond)
	if err != nil {
		return err
	}
	if isTrue(cond) {
		return workflow.executeSteps(ifStep.Steps)
	}
	return nil
}

func isTrue(s string) bool {
	l := strings.ToLower(s)
	return l == "1" || l == "t" || l == "true" || l == "yes"
}
