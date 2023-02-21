package providers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/dredge-dev/dredge/internal/api"
)

var PORTS_RE = regexp.MustCompile(`0.0.0.0:([0-9]+)->([0-9]+)/tcp`)

type LocalDockerComposeProvider struct {
	Env          string
	Path         string
	Image        string
	Proto        string
	absolutePath string
}

func (l *LocalDockerComposeProvider) Init(config map[string]string) error {
	err := checkConfig(config, []string{"env", "path", "image"})
	if err != nil {
		return err
	}
	l.Env = config["env"]
	l.Path = config["path"]
	l.Image = config["image"]
	l.Proto = config["proto"]
	l.absolutePath, err = l.getAbsolutePath()
	return err
}

func (l *LocalDockerComposeProvider) Name() string {
	return "local-docker-compose"
}

func (l *LocalDockerComposeProvider) ExecuteCommand(commandName string, callbacks api.Callbacks) (interface{}, error) {
	if commandName == "get" {
		return l.Get(callbacks)
	} else if commandName == "describe" {
		return l.Describe(callbacks)
	} else if commandName == "update" {
		return l.Update(callbacks)
	}
	return nil, fmt.Errorf("could not find command %s", commandName)
}

func (l *LocalDockerComposeProvider) Get(callbacks api.Callbacks) ([]map[string]interface{}, error) {
	deploy, err := l.get()
	if err != nil {
		return nil, err
	}
	return []map[string]interface{}{deploy}, nil
}

func (l *LocalDockerComposeProvider) get() (map[string]interface{}, error) {
	instances, err := l.getInstances()
	if err != nil {
		return nil, err
	}

	version, err := l.getVersion()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":      l.Env,
		"version":   version,
		"instances": fmt.Sprintf("%d", instances),
		"type":      "container",
	}, nil
}

func (l *LocalDockerComposeProvider) Describe(c api.Callbacks) (map[string]interface{}, error) {
	inputs, err := c.RequestInput([]api.InputRequest{
		{
			Name:        "name",
			Description: "Name",
			Type:        api.Text,
		},
	})
	if err != nil {
		return nil, err
	}
	if inputs["name"] != l.Env {
		return nil, &api.NoResult{}
	}
	ret, err := l.get()
	if err != nil {
		return nil, err
	}
	ret["provider"] = l.Name()
	if ret["instances"] != "0" {
		containers := []string{}
		ps, err := l.ps()
		if err != nil {
			return nil, err
		}
		for _, line := range ps {
			parts := strings.Split(line, " ")
			containers = append(containers, parts[0])
			ports := parts[len(parts)-1]
			if PORTS_RE.MatchString(ports) {
				p := PORTS_RE.FindStringSubmatch(ports)
				if len(p) > 1 && l.Proto == "http" {
					ret["url"] = fmt.Sprintf("http://localhost:%s", p[1])
				}
			}
		}
		ret["containers"] = containers
	}
	ret["path"] = l.absolutePath
	return ret, nil
}

func (l *LocalDockerComposeProvider) Update(c api.Callbacks) (map[string]interface{}, error) {
	inputs, err := c.RequestInput([]api.InputRequest{
		{
			Name:        "name",
			Description: "name",
			Type:        api.Text,
		},
	})
	if err != nil {
		return nil, err
	}
	if inputs["name"] != l.Env {
		return nil, &api.NoResult{}
	}

	inputs, err = c.RequestInput([]api.InputRequest{
		{
			Name:        "version",
			Description: "version",
			Type:        api.Text,
		},
		{
			Name:        "instances",
			Description: "instances",
			Type:        api.Text,
		},
	})
	if err != nil {
		return nil, err
	}

	if inputs["version"] != "" {
		err = l.updateVersion(inputs["version"])
		if err != nil {
			return nil, err
		}
	}

	instances, err := strconv.Atoi(inputs["instances"])
	if err != nil {
		return nil, err
	}
	err = l.setInstances(instances, c)
	if err != nil {
		return nil, err
	}

	return l.get()
}

func (l *LocalDockerComposeProvider) setInstances(instances int, c api.Callbacks) error {
	current, err := l.getInstances()
	if err != nil {
		return err
	}
	if instances == current {
		c.Log(api.Info, "Restarting docker-compose")
		return l.restart()
	}
	if instances > current {
		c.Log(api.Info, "Starting docker-compose")
		return l.start()
	}
	c.Log(api.Info, "Stopping docker-compose")
	return l.stop()
}

func (l *LocalDockerComposeProvider) restart() error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && docker-compose restart", l.Path))
	_, err := cmd.Output()
	return err
}

func (l *LocalDockerComposeProvider) start() error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && docker-compose up -d", l.Path))
	_, err := cmd.Output()
	return err
}

func (l *LocalDockerComposeProvider) stop() error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && docker-compose down", l.Path))
	_, err := cmd.Output()
	return err
}

func (l *LocalDockerComposeProvider) getInstances() (int, error) {
	containers, err := l.ps()
	if err != nil {
		return 0, err
	}

	if len(containers) > 0 {
		return 1, nil
	}
	return 0, nil
}

func (l *LocalDockerComposeProvider) ps() ([]string, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cd %s && docker-compose ps", l.Path))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")
	if len(out) > 2 {
		return out[2:], nil
	}
	return nil, nil
}

func (l *LocalDockerComposeProvider) getVersion() (string, error) {
	b, err := os.ReadFile(l.absolutePath)
	if err != nil {
		return "", err
	}
	r, err := regexp.Compile(l.Image + `:([^\s]+)`)
	if err != nil {
		return "", err
	}
	version := r.FindString(string(b))
	if version == "" {
		return "", fmt.Errorf("version could not be determined")
	}
	return strings.TrimPrefix(version, l.Image+":"), nil
}

func (l *LocalDockerComposeProvider) updateVersion(version string) error {
	b, err := os.ReadFile(l.absolutePath)
	if err != nil {
		return err
	}
	r, err := regexp.Compile(l.Image + `:([^\s]+)`)
	if err != nil {
		return err
	}
	output := r.ReplaceAllString(string(b), l.Image+":"+version)
	return os.WriteFile(l.absolutePath, []byte(output), 0644)
}

func (l *LocalDockerComposeProvider) getAbsolutePath() (string, error) {
	fileInfo, err := os.Stat(l.Path)
	if err != nil {
		return "", err
	}
	if fileInfo.IsDir() {
		return filepath.Join(l.Path, "docker-compose.yml"), nil
	}
	return l.Path, nil
}
