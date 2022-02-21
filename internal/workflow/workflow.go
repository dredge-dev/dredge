package workflow

import (
	"fmt"
	"os"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/pkg/browser"
)

func ExecuteWorkflow(dredgeFile *config.DredgeFile, w config.Workflow) error {
	return executeWorkflow(dredgeFile, w)
}

func executeWorkflow(dredgeFile *config.DredgeFile, workflow config.Workflow) error {
	if workflow.Import != nil {
		return importWorkflow(dredgeFile, *workflow.Import)
	}

	env := NewEnv()
	for input, description := range workflow.Inputs {
		err := env.AddInput(input, description, os.Stdin)
		if err != nil {
			return err
		}
	}

	for _, step := range workflow.Steps {
		err := executeStep(dredgeFile, workflow, step, env)
		if err != nil {
			return err
		}
	}
	return nil
}

func importWorkflow(dredgeFile *config.DredgeFile, iw config.ImportWorkflow) error {
	source := dredgeFile
	if iw.Source != "" {
		var err error
		source, err = config.GetDredgeFile(iw.Source)
		if err != nil {
			return err
		}
	}

	w := source.GetWorkflow(iw.Bucket, iw.Workflow)
	if w == nil {
		return fmt.Errorf("Could not find import workflow")
	}

	return ExecuteWorkflow(source, *w)
}

func executeStep(dredgeFile *config.DredgeFile, workflow config.Workflow, step config.Step, env Env) error {
	if step.Shell != nil {
		return executeShellStep(dredgeFile, workflow, step.Shell, env)
	} else if step.Template != nil {
		return executeTemplate(dredgeFile, workflow, step.Template, env)
	} else if step.Browser != nil {
		return openBrowser(*step.Browser)
	}
	return fmt.Errorf("No execution found for step.")
}

func executeShellStep(dredgeFile *config.DredgeFile, workflow config.Workflow, shell *config.ShellStep, env Env) error {
	runtime, err := GetRuntime(dredgeFile.Env, shell.Runtime)
	if err != nil {
		return err
	}
	return runtime.Execute(shell.Cmd)
}

func openBrowser(url string) error {
	return browser.OpenURL(url)
}
