package workflow

import (
	"bytes"
	"fmt"
	"io"
	"os"
	osExec "os/exec"
	"sort"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
)

const dredgeDir = ".dredge"
const cacheDir = "cache"

type Runtime struct {
	Env    exec.Env
	Config config.Runtime
}

func GetRuntime(env exec.Env, runtimes []config.Runtime, name string) (*Runtime, error) {
	if name == "" {
		return &Runtime{env, config.Runtime{Type: "native"}}, nil
	}
	for _, r := range runtimes {
		if name == r.Name {
			return &Runtime{env, r}, nil
		}
	}
	return nil, fmt.Errorf("Runtime %s is not defined", name)
}

func (r *Runtime) Execute(interactive bool, command string, stdin io.Reader, stdout, stderr *bytes.Buffer) error {
	cmd, err := r.GetCommand(interactive, command)
	if err != nil {
		return err
	}
	osCmd := osExec.Command("/bin/bash", "-c", cmd)
	osCmd.Env = os.Environ()
	if stdin != nil {
		osCmd.Stdin = stdin
	} else {
		osCmd.Stdin = os.Stdin
	}
	if stdout != nil {
		osCmd.Stdout = stdout
	} else {
		osCmd.Stdout = os.Stdout
	}
	if stderr != nil {
		osCmd.Stderr = stderr
	} else {
		osCmd.Stderr = os.Stderr
	}
	return osCmd.Run()
}

func (r *Runtime) GetCommand(interactive bool, cmd string) (string, error) {
	var command string
	var err error
	if r.Config.Type == config.RUNTIME_NATIVE {
		command = cmd
	} else if r.Config.Type == config.RUNTIME_CONTAINER {
		command, err = r.getContainerCommand(interactive, cmd)
	} else {
		err = fmt.Errorf("unknown runtime type %s", r.Config.Type)
	}
	if err != nil {
		return "", err
	}
	return Template(command, r.Env)
}

func (r *Runtime) getContainerCommand(interactive bool, cmd string) (string, error) {
	workDir := r.Config.GetHome()

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var envVars []string
	for variable, value := range r.Env {
		envVars = append(envVars, fmt.Sprintf("-e %s=%s", variable, value))
	}
	sort.Strings(envVars)

	var volumes []string
	for _, c := range r.Config.Cache {
		if !strings.HasPrefix(c, "/") {
			return "", fmt.Errorf("invalid cache path (%s): path should start with /", c)
		}
		volumes = append(volumes, fmt.Sprintf("-v %s/%s/%s%s:%s", currentDir, dredgeDir, cacheDir, c, c))
	}
	if len(r.Config.GlobalCache) > 0 {
		globalCacheDir, err := getGlobalCacheDir(r.Config)
		if err != nil {
			return "", err
		}
		for _, c := range r.Config.GlobalCache {
			if !strings.HasPrefix(c, "/") {
				return "", fmt.Errorf("invalid cache path (%s): path should start with /", c)
			}
			volumes = append(volumes, fmt.Sprintf("-v %s%s:%s", globalCacheDir, c, c))
		}
	}
	volumes = append(volumes, fmt.Sprintf("-v %s:%s", currentDir, workDir))

	var ports []string
	for _, p := range r.Config.Ports {
		portsString, err := Template(p, r.Env)
		if err != nil {
			return "", err
		}
		portsParts := strings.Split(portsString, ",")
		for _, port := range portsParts {
			if len(port) > 0 {
				if strings.Contains(port, ":") {
					ports = append(ports, fmt.Sprintf("-p %s", port))
				} else {
					ports = append(ports, fmt.Sprintf("-p %s:%s", port, port))
				}
			}
		}
	}

	flags := ""
	if interactive {
		flags = "-it"
	}

	return fmt.Sprintf(
		"docker run --rm %s %s %s -w %s %s %s %s",
		strings.Join(envVars, " "), strings.Join(volumes, " "), strings.Join(ports, " "), workDir, flags, r.Config.Image, cmd), nil
}

func getGlobalCacheDir(r config.Runtime) (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	globalCacheDir := fmt.Sprintf("%s/%s/%s/%s", userHome, dredgeDir, cacheDir, r.Name)

	err = os.MkdirAll(globalCacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return globalCacheDir, nil
}
