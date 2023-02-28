package workflow

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
	"github.com/dredge-dev/dredge/internal/config"
)

type Bucket struct {
	Name        string
	Description string
	Workflows   []config.Workflow
	Callbacks   api.Callbacks
}

type Workflow struct {
	Name        string
	Description string
	Inputs      []config.Input
	Steps       []config.Step
	Runtimes    []config.Runtime
	Callbacks   api.Callbacks
}

func (workflow *Workflow) Execute() error {
	// TODO Re-arrange to get all inputs at once
	for _, input := range workflow.Inputs {
		skip, err := workflow.Callbacks.Template(input.Skip)
		if err != nil {
			return err
		}
		if skip != "true" {
			result, err := workflow.Callbacks.RequestInput([]api.InputRequest{
				toInputRequest(input),
			})
			if err != nil {
				return err
			}
			// TODO Should this be moved to Exec? Should every RequestIput set the env?
			if err := workflow.Callbacks.SetEnv(input.Name, result[input.Name]); err != nil {
				return err
			}
		}
	}
	return workflow.executeSteps(workflow.Steps)
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

func (workflow *Workflow) executeSteps(steps []config.Step) error {
	for _, step := range steps {
		err := workflow.executeStep(step)
		if err != nil {
			return err
		}
	}
	return nil
}

func (workflow *Workflow) executeStep(step config.Step) error {
	if step.Shell != nil {
		return workflow.executeShellStep(step.Shell)
	} else if step.Template != nil {
		return workflow.executeTemplate(step.Template)
	} else if step.Browser != nil {
		return workflow.openBrowser(step.Browser)
	} else if step.EditDredgeFile != nil {
		return workflow.executeEditDredgeFile(step.EditDredgeFile)
	} else if step.If != nil {
		return workflow.executeIfStep(step.If)
	} else if step.Execute != nil {
		return workflow.executeExecuteStep(step.Execute)
	} else if step.Set != nil {
		return workflow.executeSetStep(step.Set)
	} else if step.Log != nil {
		return workflow.executeLogStep(step.Log)
	} else if step.Confirm != nil {
		return workflow.executeConfirmStep(step.Confirm)
	}
	return fmt.Errorf("no execution found for step %v", step.Name)
}

func (workflow *Workflow) executeShellStep(shell *config.ShellStep) error {
	runtime, err := workflow.GetRuntime(shell.Runtime)
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
		workflow.Callbacks.SetEnv(shell.StdOut, stdout.String())
	}
	if shell.StdErr != "" && stderr != nil {
		workflow.Callbacks.SetEnv(shell.StdErr, stderr.String())
	}
	return nil
}

func (workflow *Workflow) openBrowser(b *config.BrowserStep) error {
	url, err := workflow.Callbacks.Template(b.Url)
	if err != nil {
		return err
	}
	return workflow.Callbacks.OpenUrl(url)
}

func (workflow *Workflow) executeExecuteStep(execute *config.ExecuteStep) error {
	output, err := workflow.Callbacks.ExecuteResourceCommand(execute.Resource, execute.Command)
	if execute.Register != "" {
		workflow.Callbacks.SetEnv(execute.Register, output.Output)
	}
	return err
}

func (workflow *Workflow) executeSetStep(set *config.SetStep) error {
	for key, value := range *set {
		templated, err := workflow.Callbacks.Template(value)
		if err != nil {
			return err
		}
		err = workflow.Callbacks.SetEnv(key, templated)
		if err != nil {
			return err
		}
	}
	return nil
}

func (workflow *Workflow) executeLogStep(log *config.LogStep) error {
	level, err := toLogLevel(log.Level)
	if err != nil {
		return err
	}
	return workflow.Callbacks.Log(level, log.Message)
}

func toLogLevel(level string) (api.LogLevel, error) {
	levelLower := strings.ToLower(level)
	levels := []api.LogLevel{api.Fatal, api.Error, api.Warn, api.Info, api.Debug, api.Trace}
	for _, level := range levels {
		if levelLower == strings.ToLower(level.String()) {
			return level, nil
		}
	}
	return 0, fmt.Errorf("unknown log level: %s", level)
}

func (workflow *Workflow) executeConfirmStep(confirm *config.ConfirmStep) error {
	return workflow.Callbacks.Confirm(confirm.Message)
}
