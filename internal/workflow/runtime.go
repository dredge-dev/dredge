package workflow

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

type Runtime struct {
	Env      config.Env
	Template string
}

func GetRuntime(env config.Env, name string) (*Runtime, error) {
	if name == "" {
		return createDefaultRuntime(env)
	}
	for _, r := range env.Runtimes {
		if name == r.Name {
			runtime, err := createRuntime(env, r)
			if err != nil {
				return nil, err
			}
			return runtime, nil
		}
	}
	return nil, fmt.Errorf("Runtime %s is not defined", name)
}

func createDefaultRuntime(e config.Env) (*Runtime, error) {
	return createRuntime(e, config.Runtime{Type: "native"})
}

func createRuntime(e config.Env, r config.Runtime) (*Runtime, error) {
	if r.Type == "container" {
		return createContainerRuntime(e, r)
	} else if r.Type == "native" {
		return createNativeRuntime(e, r)
	}
	return nil, fmt.Errorf("Runtime type %s is not defined", r.Type)
}

func createContainerRuntime(e config.Env, r config.Runtime) (*Runtime, error) {
	workDir := r.GetHome()

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

	var ports []string
	for _, p := range r.Ports {
		ports = append(ports, fmt.Sprintf("-p %s", p))
	}

	return &Runtime{
		e,
		fmt.Sprintf(
			"docker run --rm %s %s %s -w %s -it %s ${cmd}",
			strings.Join(envVars, " "), strings.Join(volumes, " "), strings.Join(ports, " "), workDir, r.Image),
	}, nil
}

func createNativeRuntime(e config.Env, r config.Runtime) (*Runtime, error) {
	return &Runtime{e, "${cmd}"}, nil
}

func (r *Runtime) Execute(command string) error {
	cmd := exec.Command("/bin/bash", "-c", r.getCommand(command))
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runtime) getCommand(cmd string) string {
	command := strings.Replace(r.Template, "${cmd}", cmd, -1)

	for variable, value := range r.Env.Variables {
		command = strings.Replace(command, fmt.Sprintf("${%s}", variable), value, -1)
	}

	return command
}
