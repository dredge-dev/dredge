package workflow

import (
	"fmt"
	"os"

	"github.com/dredge-dev/dredge/internal/config"
)

func ExecuteWorkflow(dredgeFile *config.DredgeFile, w config.Workflow) error {
	return executeWorkflow(dredgeFile, w)
}

func executeWorkflow(dredgeFile *config.DredgeFile, workflow config.Workflow) error {
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

func executeStep(dredgeFile *config.DredgeFile, workflow config.Workflow, step config.Step, env Env) error {
	if step.Shell != nil {
		return executeShellStep(dredgeFile, workflow, step.Shell, env)
	} else if step.Template != nil {
		return executeTemplate(dredgeFile, workflow, step.Template, env)
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
