package workflow

import (
	"fmt"
	"os"
	osExec "os/exec"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

type Runtime struct {
	Env      exec.Env
	Template string
}

func GetRuntime(runtimes []config.Runtime, name string, env exec.Env) (*Runtime, error) {
	if name == "" {
		return createDefaultRuntime(env)
	}
	for _, r := range runtimes {
		if name == r.Name {
			runtime, err := createRuntime(r, env)
			if err != nil {
				return nil, err
			}
			return runtime, nil
		}
	}
	return nil, fmt.Errorf("Runtime %s is not defined", name)
}

func createDefaultRuntime(e exec.Env) (*Runtime, error) {
	return createRuntime(config.Runtime{Type: "native"}, e)
}

func createRuntime(r config.Runtime, e exec.Env) (*Runtime, error) {
	if r.Type == "container" {
		return createContainerRuntime(r, e)
	} else if r.Type == "native" {
		return createNativeRuntime(r, e)
	}
	return nil, fmt.Errorf("Runtime type %s is not defined", r.Type)
}

func createContainerRuntime(r config.Runtime, e exec.Env) (*Runtime, error) {
	workDir := r.GetHome()

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var envVars []string
	for variable, value := range e {
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

	var ports []string
	for _, p := range r.Ports {
		ports = append(ports, fmt.Sprintf("-p %s", p))
	}

	return &Runtime{
		e,
		fmt.Sprintf(
			"docker run --rm %s %s %s -w %s -it %s {{ .cmd }}",
			strings.Join(envVars, " "), strings.Join(volumes, " "), strings.Join(ports, " "), workDir, r.Image),
	}, nil
}

func createNativeRuntime(r config.Runtime, e exec.Env) (*Runtime, error) {
	return &Runtime{e, "{{ .cmd }}"}, nil
}

func (r *Runtime) Execute(command string) error {
	cmd := osExec.Command("/bin/bash", "-c", r.getCommand(command))
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runtime) getCommand(cmd string) string {
	command := strings.Replace(r.Template, "{{ .cmd }}", cmd, -1)

	for variable, value := range r.Env {
		command = strings.Replace(command, fmt.Sprintf("{{ .%s }}", variable), value, -1)
	}

	return command
}
