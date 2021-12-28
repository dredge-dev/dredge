package workflow

import (
	"fmt"
	"strings"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/dredge-dev/dredge/internal/config"
)

func ExecuteWorkflow(dredgeFile *config.DredgeFile, cmd *cobra.Command, args []string) error {
	for _, w := range dredgeFile.Workflows {
		if cmd.Name() == w.Name {
			return execute(dredgeFile, w)
		}
	}
	return fmt.Errorf("Workflow %s is not defined", cmd.Name())
}

func execute(dredgeFile *config.DredgeFile, workflow config.Workflow) error {
	for _, step := range workflow.Steps {
		err := executeStep(dredgeFile, workflow, step)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeStep(dredgeFile *config.DredgeFile, workflow config.Workflow, step config.Step) error {
	for _, r := range dredgeFile.Env.Runtimes {
		if step.Runtime == r.Name {
			runtime, err := createRuntime(dredgeFile.Env, r)
			if err != nil {
				return err
			}
			return runtime.execute(step)
		}
	}
	return fmt.Errorf("Runtime %s is not defined", step.Runtime)
}

type Runtime struct {
	Env config.Env
	Command string
}

func createRuntime(e config.Env, r config.Runtime) (*Runtime, error) {
	if r.Type != "container" {
		return nil, fmt.Errorf("Runtime type %s is not defined (use: container)", r.Type)
	}

	workDir := "/home"

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var envVars []string
	for variable, value := range e.Variables {
		envVars = append(envVars, fmt.Sprintf("-e %s=%s", variable, value))
	}

	var volumes []string
	for _, c := range r.Cache {
		if !strings.HasPrefix(c, "/") {
			return nil, fmt.Errorf("Invalid cache path (%s): path should start with /", c)
		}
		volumes = append(volumes, fmt.Sprintf("-v %s/.dredge/cache%s:%s", currentDir, c, c))
	}
	volumes = append(volumes, fmt.Sprintf("-v %s:%s", currentDir, workDir))

	return &Runtime{
		e,
		fmt.Sprintf(
			"docker run --rm %s %s -w %s -it %s",
			strings.Join(envVars, " "), strings.Join(volumes, " "), workDir, r.Image),
	}, nil
}

func (r *Runtime) execute(step config.Step) error {
	command := r.Command + " " + step.Exec

	for variable, value := range r.Env.Variables {
		command = strings.Replace(command, fmt.Sprintf("${%s}", variable), value, -1)
	}

	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
