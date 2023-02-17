package workflow

import (
	"bytes"
	"fmt"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

func ExecuteWorkflow(workflow *exec.Workflow) error {
	// TODO Re-arrange to get all inputs at once
	for _, input := range workflow.Inputs {
		skip, err := Template(input.Skip, workflow.Exec.Env)
		if err != nil {
			return err
		}
		if skip != "true" {
			result, err := workflow.Exec.RequestInput([]api.InputRequest{
				toInputRequest(input),
			})
			if err != nil {
				return err
			}
			workflow.Exec.Env[input.Name] = result[input.Name]
		}
	}
	return executeSteps(workflow, workflow.Steps)
}

func toInputRequest(input config.Input) api.InputRequest {
	return api.InputRequest{
		Name:         input.Name,
		Description:  input.Description,
		Type:         toInputType(input.Type),
		Values:       input.Values,
		DefaultValue: input.DefaultValue,
	}
}

func toInputType(t string) api.InputType {
	if t == config.INPUT_SELECT {
		return api.Select
	}
	return api.Text
}

func executeSteps(workflow *exec.Workflow, steps []config.Step) error {
	for _, step := range steps {
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
	} else if step.EditDredgeFile != nil {
		return executeEditDredgeFile(workflow, step.EditDredgeFile)
	} else if step.If != nil {
		return executeIfStep(workflow, step.If)
	} else if step.Execute != nil {
		return executeExecuteStep(workflow, step.Execute)
	}
	return fmt.Errorf("no execution found for step %v", step.Name)
}

func executeShellStep(workflow *exec.Workflow, shell *config.ShellStep) error {
	runtime, err := GetRuntime(workflow.Exec.Env, workflow.Exec.DredgeFile.Runtimes, shell.Runtime)
	if err != nil {
		return err
	}
	interactive := true
	var stdout *bytes.Buffer
	if shell.StdOut != "" {
		stdout = new(bytes.Buffer)
		interactive = false
	}
	var stderr *bytes.Buffer
	if shell.StdErr != "" {
		stderr = new(bytes.Buffer)
		interactive = false
	}
	err = runtime.Execute(interactive, shell.Cmd, nil, stdout, stderr)
	if err != nil {
		return err
	}
	if shell.StdOut != "" && stdout != nil {
		workflow.Exec.Env[shell.StdOut] = stdout.String()
	}
	if shell.StdErr != "" && stderr != nil {
		workflow.Exec.Env[shell.StdErr] = stderr.String()
	}
	return nil
}

func openBrowser(workflow *exec.Workflow, b *config.BrowserStep) error {
	url, err := Template(b.Url, workflow.Exec.Env)
	if err != nil {
		return err
	}
	return workflow.Exec.OpenUrl(url)
}

func executeExecuteStep(workflow *exec.Workflow, execute *config.ExecuteStep) error {
	output, err := workflow.Exec.ExecuteResourceCommand(execute.Resource, execute.Command)
	if execute.Register != "" {
		workflow.Exec.Env[execute.Register] = output.Output
	}
	return err
}
