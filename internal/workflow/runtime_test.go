package workflow

import (
	"fmt"
	"os"
	"testing"

	"github.com/dredge-dev/dredge/internal/config"
	"github.com/dredge-dev/dredge/internal/exec"
	"github.com/stretchr/testify/assert"
)

func TestGetRuntime(t *testing.T) {
	home := "/home"
	wd, _ := os.Getwd()

	runtimes := []config.Runtime{
		{
			Name:  "build-container",
			Type:  "container",
			Image: "build-image:latest",
			Home:  &home,
			Cache: []string{"/go"},
			Ports: []string{"8080:8080"},
		},
		{
			Name:  "port-container",
			Type:  "container",
			Image: "port-image:latest",
			Home:  &home,
			Cache: []string{"/test"},
			Ports: []string{"{{ .PORTS }}"},
		},
		{
			Name: "native",
			Type: "native",
		},
	}

	env := exec.NewEnv()
	env["PORTS"] = "1234,80"

	tests := map[string]struct {
		name    string
		env     exec.Env
		command string
	}{
		"container": {
			name:    "build-container",
			env:     exec.NewEnv(),
			command: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/go:/go -v %s:/home -p 8080:8080 -w /home -it build-image:latest {{ .cmd }}", wd, wd),
		},
		"container with ports": {
			name:    "port-container",
			env:     env,
			command: fmt.Sprintf("docker run --rm -e PORTS=1234,80 -v %s/.dredge/cache/test:/test -v %s:/home -p 1234:1234 -p 80:80 -w /home -it port-image:latest {{ .cmd }}", wd, wd),
		},
		"container without ports": {
			name:    "port-container",
			env:     exec.NewEnv(),
			command: fmt.Sprintf("docker run --rm  -v %s/.dredge/cache/test:/test -v %s:/home  -w /home -it port-image:latest {{ .cmd }}", wd, wd),
		},
		"native": {
			name:    "native",
			env:     exec.NewEnv(),
			command: "{{ .cmd }}",
		},
		"default": {
			name:    "",
			env:     exec.NewEnv(),
			command: "{{ .cmd }}",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		runtime, err := GetRuntime(runtimes, test.name, test.env)
		assert.Nil(t, err)
		assert.Equal(t, test.command, runtime.Template)
	}
}

func TestGetCommand(t *testing.T) {
	env := exec.NewEnv()
	env["var1"] = "test"
	env["var2"] = "hello"

	tests := map[string]struct {
		runtime *Runtime
		command string
		output  string
	}{
		"empty env": {
			runtime: &Runtime{Env: exec.NewEnv(), Template: "prefix {{ .cmd }}"},
			command: "cmd",
			output:  "prefix cmd",
		},
		"env replace": {
			runtime: &Runtime{Env: env, Template: "prefix {{ .var1 }} {{ .cmd }}"},
			command: "cmd {{ .var2 }}",
			output:  "prefix test cmd hello",
		},
		"env replace in runtime": {
			runtime: &Runtime{Env: env, Template: "prefix {{ .var1 }} {{ .cmd }}"},
			command: "cmd",
			output:  "prefix test cmd",
		},
		"env replace in command": {
			runtime: &Runtime{Env: env, Template: "prefix {{ .cmd }}"},
			command: "cmd {{ .var2 }}",
			output:  "prefix cmd hello",
		},
		"command with ||": {
			runtime: &Runtime{Env: env, Template: "{{ .cmd }}"},
			command: "test || out",
			output:  "test || out",
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		assert.Equal(t, test.output, test.runtime.getCommand(test.command))
	}
}
