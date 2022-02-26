package workflow

import (
	"fmt"
	"os"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/pkg/browser"
)

func ExecuteWorkflow(workflow *exec.Workflow) error {
	for input, description := range workflow.Inputs {
		err := workflow.Exec.Env.AddInput(input, description, os.Stdin)
		if err != nil {
			return err
		}
	}

	for _, step := range workflow.Steps {
		err := executeStep(workflow, step)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeStep(workflow *exec.Workflow, step config.Step) error {
	if step.Shell != nil {
		return executeShellStep(workflow, step.Shell)
	} else if step.Template != nil {
		return executeTemplate(workflow, step.Template)
	} else if step.Browser != nil {
		return openBrowser(workflow, step.Browser)
	}
	return fmt.Errorf("No execution found for step.")
}

func executeShellStep(workflow *exec.Workflow, shell *config.ShellStep) error {
	runtime, err := GetRuntime(workflow.Exec.DredgeFile.Env.Runtimes, shell.Runtime, workflow.Exec.Env)
	if err != nil {
		return err
	}
	return runtime.Execute(shell.Cmd)
}

func openBrowser(workflow *exec.Workflow, b *config.BrowserStep) error {
	return browser.OpenURL(b.Url)
}
